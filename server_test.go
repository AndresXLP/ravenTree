package ravenTree_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/AndresXLP/ravenTree"
	"github.com/charmbracelet/log"
)

func RunServerTesting() {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/timeout", geTimedOut)
	mux.HandleFunc("GET /api/retry", getRetry)
	mux.HandleFunc("GET /api/query-params", getWithQueryParams)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start the server in a goroutine
	go func() {
		log.Print("Starting server on :8080\n")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Could not start server: %s\n", err)
		}
	}()
}

func geTimedOut(w http.ResponseWriter, r *http.Request) {
	// Simulate a 3-second delay (timeout)
	time.Sleep(3 * time.Second)

	w.Header().Set(ravenTree.HeaderContentType, ravenTree.MIMEApplicationJSON)
	w.WriteHeader(http.StatusOK)
	// Write an empty JSON response
	_, _ = w.Write([]byte(`{}`))
}

var try = 0

var mu sync.Mutex // To prevent race conditions when updating 'try'

func getRetry(w http.ResponseWriter, r *http.Request) {

	w.Header().Set(ravenTree.HeaderContentType, ravenTree.MIMEApplicationJSON)

	// Lock to prevent race conditions in case of concurrent requests
	mu.Lock()
	if try == 3 {
		log.Infof("Successful request on the %drd attempt", try)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`)) // Empty JSON response
		return
	}

	if try == 0 {
		log.Info("Initial attempt")
	} else {
		log.Infof("Trying to retry %d", try)
	}

	try += 1
	mu.Unlock()

	time.Sleep(3 * time.Second) // Simulate delay for retries

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{}`)) // Empty JSON response
}

func getWithQueryParams(w http.ResponseWriter, r *http.Request) {
	// Extract query parameters
	email := r.URL.Query().Get("email")
	username := r.URL.Query().Get("username")
	test := r.URL.Query().Get("test")
	// Create the response struct
	resp := response{
		Data: []string{email, username, test},
	}

	w.Header().Set(ravenTree.HeaderContentType, ravenTree.MIMEApplicationJSON)
	w.WriteHeader(http.StatusOK)

	// Encode the response into JSON and write it to the response writer
	_ = json.NewEncoder(w).Encode(resp)
}
