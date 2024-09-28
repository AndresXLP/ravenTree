package ravenTree

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/AndresXLP/ravenTree/utils_tests"
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
)

type ravenTreeTestSuite struct {
	suite.Suite
	underTest Tree
}

func TestRavenTreeSuite(t *testing.T) {
	utils.RunServerTesting()
	suite.Run(t, new(ravenTreeTestSuite))
}
func (suite *ravenTreeTestSuite) SetupTest() {
	suite.underTest = NewRavensTree()
}

func (suite *ravenTreeTestSuite) TestSendRaven_Success() {
	options := &Options{
		Host:        "http://localhost:8080",
		Path:        "/api/query-params",
		Method:      http.MethodGet,
		QueryParams: queryParams,
		Headers:     headers,
		Timeout:     1 * time.Second,
		RetryCount:  3,
	}

	ctx := context.Background()
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

func (suite *ravenTreeTestSuite) TestSendRaven_FailWhenTimedOut() {
	options := &Options{
		Host:       "http://localhost:8080",
		Path:       "/api/timeout",
		Method:     http.MethodGet,
		Body:       nil,
		Timeout:    1 * time.Second,
		RetryCount: 3,
	}
	ctx := context.Background()

	_, err := suite.underTest.SendRaven(ctx, options)

	suite.Error(err)
	suite.ErrorContains(err, "context deadline exceeded (Client.Timeout exceeded while awaiting headers)")
}

func (suite *ravenTreeTestSuite) TestSendRaven_SuccessWhenRetry() {
	options := &Options{
		Host:       "http://localhost:8080",
		Path:       "/api/retry",
		Method:     http.MethodGet,
		Timeout:    1 * time.Second,
		RetryCount: 3,
	}

	ctx := context.Background()
	_, err := suite.underTest.SendRaven(ctx, options)
	suite.NoError(err)
}

func (suite *ravenTreeTestSuite) TestSendRaven_FailWhenInvalidURL() {
	options := &Options{
		Host:       ":foo",
		Path:       "/api/retry",
		Method:     http.MethodGet,
		Timeout:    1 * time.Second,
		RetryCount: 3,
	}

	ctx := context.Background()
	_, err := suite.underTest.SendRaven(ctx, options)
	suite.Error(err)
}

func (suite *ravenTreeTestSuite) TestSendRaven_FailWhenInvalidPath() {
	options := &Options{
		Host:       "http://localhost:8080",
		Path:       ":foo",
		Method:     http.MethodGet,
		Timeout:    1 * time.Second,
		RetryCount: 3,
	}

	ctx := context.Background()
	_, err := suite.underTest.SendRaven(ctx, options)
	suite.Error(err)
}

func (suite *ravenTreeTestSuite) TestSendRaven_FailWhenInvalidBody() {
	options := &Options{
		Host:       "http://localhost:8080",
		Path:       "/api/retry",
		Method:     http.MethodGet,
		Body:       wrongBody,
		Timeout:    1 * time.Second,
		RetryCount: 3,
	}

	ctx := context.Background()
	_, err := suite.underTest.SendRaven(ctx, options)
	suite.Error(err)
}

func (suite *ravenTreeTestSuite) TestSendRaven_FailRequestInvalidMethod() {
	options := &Options{
		Host:       "http://localhost:8080",
		Path:       "/api/retry",
		Method:     "😰",
		Timeout:    1 * time.Second,
		RetryCount: 3,
	}

	ctx := context.Background()
	_, err := suite.underTest.SendRaven(ctx, options)
	suite.Error(err)
}
