package ravenTree

import (
	"bytes"
	"encoding/json"
	"net/url"
	"time"
)

const (
	defaultRetryCount      = 1
	defaultMaxBackoffDelay = 10 * time.Second
)

type Options struct {
	Host            string
	Path            string
	Method          string
	Body            interface{}
	QueryParams     map[string]string
	Headers         map[string]string
	Timeout         time.Duration
	RetryCount      int
	BackoffStrategy *BackoffStrategy
	MaxBackoffDelay time.Duration
}

// bodyToBufferBody serializes the Body field into JSON and stores it in a bytes.Buffer.
//
// It marshals the `Body` field of the `Options` struct into a JSON byte slice. If an error occurs
// during the marshalling process, it returns an empty buffer and the error.
//
// Returns:
// - bytes.Buffer: A buffer containing the serialized JSON of the `Body` field.
// - error: An error that occurs during the JSON marshalling process, if any.
func (o *Options) bodyToBufferBody() (bytes.Buffer, error) {
	bytesJson, err := json.Marshal(o.Body)
	if err != nil {
		return bytes.Buffer{}, err
	}

	return *bytes.NewBuffer(bytesJson), nil
}

// buildURL constructs a full URL by combining the Host, Path, and QueryParams fields from the Options struct.
//
// It first parses the `Host` and `Path` fields into URLs. Then, it resolves the path relative to the base URL.
// If `QueryParams` are provided, they are appended to the final URL as query parameters.
//
// Returns:
// - string: The fully constructed URL as a string.
// - error: An error if there is a problem parsing the Host or Path.
func (o *Options) buildURL() (string, error) {
	baseURL, err := url.Parse(o.Host)
	if err != nil {
		return "", err
	}

	pathUrl, err := url.Parse(o.Path)
	if err != nil {
		return "", err
	}

	finalURL := baseURL.ResolveReference(pathUrl)

	if len(o.QueryParams) > 0 {
		query := finalURL.Query()
		for key, value := range o.QueryParams {
			query.Add(key, value)
		}

		finalURL.RawQuery = query.Encode()
	}

	return finalURL.String(), nil
}

// defaultOptions sets default values for any unset fields in the Options struct.
// If certain fields are not initialized by the user, this method assigns sensible defaults:
//   - BackoffStrategy: Defaults to the standard backoff strategy.
//   - MaxBackoffDelay: Defaults to 10 seconds.
//   - Timeout: Defaults to 30 seconds.
//   - RetryCount: Defaults to 1.
func (o *Options) defaultOptions() {
	if o.BackoffStrategy == nil {
		o.BackoffStrategy = BackoffDefault()
	}

	if o.MaxBackoffDelay == zero {
		o.MaxBackoffDelay = defaultMaxBackoffDelay
	}

	if o.RetryCount == zero {
		o.RetryCount = defaultRetryCount
	}
}
