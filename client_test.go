package ravenTree_test

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/AndresXLP/ravenTree"
	"github.com/stretchr/testify/suite"
)

type response struct {
	Data []string `json:"data"`
}

var (
	queryParams = map[string]string{
		"email":    "test@test.com",
		"username": "tester",
		"test":     "true",
	}

	headers = map[string]string{
		"Authorization": "Bearer token",
	}

	responsePrettyStringExpected = "{\"data\":[\"test@test.com\",\"tester\",\"true\"]}\n"

	responseQueryParamsExpected = response{Data: []string{
		"test@test.com", "tester", "true",
	}}

	wrongBody = struct {
		Data chan struct{}
	}{
		Data: make(chan struct{}),
	}

	ctx = context.Background()
)

type ravenTreeTestSuite struct {
	suite.Suite
	underTest ravenTree.Tree
}

func TestRavenTreeSuite(t *testing.T) {
	suite.Run(t, new(ravenTreeTestSuite))
}
func (suite *ravenTreeTestSuite) SetupTest() {
	suite.underTest = ravenTree.NewRavensTree()
}

func (suite *ravenTreeTestSuite) TestSendRaven_SuccessWithDefaultOptions() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Check request path
		suite.Equal("/api/default", r.URL.Path)
		// Check request method
		suite.Equal(http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		return
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	options := &ravenTree.Options{
		Host:   server.URL,
		Path:   "/api/default",
		Method: http.MethodGet,
	}

	resp, err := suite.underTest.SendRaven(ctx, options)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *ravenTreeTestSuite) TestSendRaven_SuccessWithHeadersAndQueryParams() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Check request path
		suite.Equal("/api/query-params", r.URL.Path)
		// Check request method
		suite.Equal(http.MethodGet, r.Method)
		// Check expected header
		suite.Equal(headers["Authorization"], r.Header.Get("Authorization"))

		// Extract query parameters
		email := r.URL.Query().Get("email")
		username := r.URL.Query().Get("username")
		test := r.URL.Query().Get("test")

		// Create the response struct
		resp := response{
			Data: []string{email, username, test},
		}

		w.WriteHeader(http.StatusOK)

		// Encode the response into JSON and write it to the response writer
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	options := &ravenTree.Options{
		Host:        server.URL,
		Path:        "/api/query-params",
		Method:      http.MethodGet,
		QueryParams: queryParams,
		Headers:     headers,
	}

	resp, err := suite.underTest.SendRaven(ctx, options)

	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	data := response{}
	err = resp.ParseBodyTo(&data)

	suite.NoError(err)
	suite.Equal(responseQueryParamsExpected, data)

	stringData := resp.ParseBodyToString()
	suite.Equal(responsePrettyStringExpected, stringData)
}

func (suite *ravenTreeTestSuite) TestSendRaven_SuccessWhenRetryWithoutBackoffStrategy() {
	try := 1
	mu := &sync.Mutex{} // To prevent race conditions when updating 'try'

	handler := func(w http.ResponseWriter, r *http.Request) {
		// Check request path
		suite.Equal("/api/retry", r.URL.Path)
		// Check request method
		suite.Equal(http.MethodGet, r.Method)

		mu.Lock()
		if try == 4 {
			log.Printf("Successful request on attempt # %d", try)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`)) // Empty JSON response
			return
		}

		log.Printf("Attempt # %d", try)
		try++
		mu.Unlock()

		time.Sleep(3 * time.Second) // Simulate delay for retries
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	options := &ravenTree.Options{
		Host:       server.URL,
		Path:       "/api/retry",
		Method:     http.MethodGet,
		Timeout:    1 * time.Second,
		RetryCount: 5,
	}

	since := time.Now()
	_, err := suite.underTest.SendRaven(ctx, options)
	duration := time.Since(since)
	expectedDuration := 3 * time.Second
	suite.NoError(err)
	suite.True(duration > expectedDuration && duration < expectedDuration+50*time.Millisecond)
}

func (suite *ravenTreeTestSuite) TestSendRaven_SuccessWhenRetryWithBackoffLineal() {
	try := 1
	mu := &sync.Mutex{} // To prevent race conditions when updating 'try'

	handler := func(w http.ResponseWriter, r *http.Request) {
		// Check request path
		suite.Equal("/api/retry", r.URL.Path)
		// Check request method
		suite.Equal(http.MethodGet, r.Method)

		mu.Lock()
		if try == 4 {
			log.Printf("Successful request on attempt # %d", try)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`)) // Empty JSON response
			return
		}

		log.Printf("Attempt # %d", try)
		try++
		mu.Unlock()

		time.Sleep(3 * time.Second) // Simulate delay for retries
		return
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	options := &ravenTree.Options{
		Host:            server.URL,
		Path:            "/api/retry",
		Method:          http.MethodGet,
		Timeout:         1 * time.Second,
		RetryCount:      5,
		BackoffStrategy: ravenTree.BackoffLineal(),
		MaxBackoffDelay: 3 * time.Second,
	}

	since := time.Now()
	_, err := suite.underTest.SendRaven(ctx, options)
	duration := time.Since(since)
	expectedDuration := 9 * time.Second
	suite.NoError(err)
	suite.True(duration > expectedDuration && duration < expectedDuration+50*time.Millisecond)
}

func (suite *ravenTreeTestSuite) TestSendRaven_SuccessWhenRetryWithBackoffExponential() {
	try := 1
	mu := &sync.Mutex{} // To prevent race conditions when updating 'try'

	handler := func(w http.ResponseWriter, r *http.Request) {
		// Check request path
		suite.Equal("/api/retry", r.URL.Path)
		// Check request method
		suite.Equal(http.MethodGet, r.Method)

		mu.Lock()
		if try == 5 {
			log.Printf("Successful request on attempt # %d", try)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`)) // Empty JSON response
			return
		}

		log.Printf("Attempt # %d", try)
		try++
		mu.Unlock()

		time.Sleep(3 * time.Second) // Simulate delay for retries
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	options := &ravenTree.Options{
		Host:            server.URL,
		Path:            "/api/retry",
		Method:          http.MethodGet,
		Timeout:         1 * time.Second,
		RetryCount:      5,
		BackoffStrategy: ravenTree.BackoffExponential(),
		MaxBackoffDelay: 15 * time.Second,
	}

	since := time.Now()
	_, err := suite.underTest.SendRaven(ctx, options)
	duration := time.Since(since)
	expectedDuration := 19 * time.Second
	suite.NoError(err)
	suite.True(duration > expectedDuration && duration < expectedDuration+50*time.Millisecond)
}

func (suite *ravenTreeTestSuite) TestSendRaven_FailWhenTimedOut() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Check request path
		suite.Equal("/api/timeout", r.URL.Path)
		// Check request method
		suite.Equal(http.MethodGet, r.Method)

		time.Sleep(2 * time.Second)
		return
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	options := &ravenTree.Options{
		Host:    server.URL,
		Path:    "/api/timeout",
		Method:  http.MethodGet,
		Timeout: 1 * time.Second,
	}

	_, err := suite.underTest.SendRaven(ctx, options)

	suite.Error(err)
	suite.ErrorContains(err, "context deadline exceeded (Client.Timeout exceeded while awaiting headers)")
}

func (suite *ravenTreeTestSuite) TestSendRaven_FailWhenInvalidURL() {
	options := &ravenTree.Options{
		Host: ":foo",
	}

	_, err := suite.underTest.SendRaven(ctx, options)
	suite.Error(err)
}

func (suite *ravenTreeTestSuite) TestSendRaven_FailWhenInvalidPath() {
	options := &ravenTree.Options{
		Host: "http://localhost:8080",
		Path: ":foo",
	}

	_, err := suite.underTest.SendRaven(ctx, options)
	suite.Error(err)
}

func (suite *ravenTreeTestSuite) TestSendRaven_FailRequestInvalidMethod() {
	options := &ravenTree.Options{
		Host:   "http://localhost:8080",
		Path:   "/api",
		Method: "😰",
	}

	_, err := suite.underTest.SendRaven(ctx, options)
	suite.Error(err)
}

func (suite *ravenTreeTestSuite) TestSendRaven_FailWhenInvalidBody() {
	options := &ravenTree.Options{
		Host:   "http://localhost:8080",
		Path:   "/api/retry",
		Method: http.MethodGet,
		Body:   wrongBody,
	}

	_, err := suite.underTest.SendRaven(ctx, options)
	suite.Error(err)
}

func (suite *ravenTreeTestSuite) TestSendRaven_FailWhenInternalServerError() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	options := &ravenTree.Options{
		Host:   server.URL,
		Path:   "/api/error",
		Method: http.MethodGet,
	}

	resp, err := suite.underTest.SendRaven(ctx, options)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, resp.StatusCode)
}
