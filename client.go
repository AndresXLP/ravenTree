package ravenTree

import (
	"context"
	"net/http"
	"time"
)

const (
	zero              = 0
	defaultTimeout    = 30 * time.Second
	HeaderContentType = "Content-Type"
	// MIMEApplicationJSON JavaScript Object Notation (JSON) https://www.rfc-editor.org/rfc/rfc8259
	MIMEApplicationJSON = "application/json"
)

// Tree defines the methods that any implementation of a RavenTree must provide.
//
// The Tree interface requires a single method, SendRaven, which sends options
// based on the provided Options and returns a WrapperResponse.
type Tree interface {
	// SendRaven sends a raven to a specified URL using the provided Options.
	//
	// This method constructs an HTTP request based on the given context and Options.
	//
	// By default, it sets the Content-Type header to application/json.
	//
	// It sends the request using an HTTP client and returns a
	// WrapperResponse that encapsulates the HTTP response.
	//
	// Parameters:
	//   - ctx: A context.Context to control the request's lifecycle and manage timeouts.
	//   - opt: A pointer to an Options struct that contains the necessary configuration
	//     for the request.
	//
	// Returns:
	// - WrapperResponse: A wrapper around the HTTP response.
	// - error: An error if the request fails at any point, or nil if the request is successful.
	SendRaven(ctx context.Context, opt *Options) (WrapperResponse, error)
}

type raven struct {
	client *http.Client
}

// NewRavensTree creates and returns a new instance of the RavenTree interface.
//
// This function acts as a constructor for the RavenTree implementation, returning
// an instance of `raven`, a private struct that implements the `Tree` interface.
// The returned instance includes an internal `http.Client` configured with a default timeout.
//
// Returns:
// - Tree: An object that implements the RavenTree interface with a pre-configured HTTP client.
func NewRavensTree() Tree {
	return &raven{
		client: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

func (r *raven) SendRaven(ctx context.Context, opt *Options) (WrapperResponse, error) {
	URL, err := opt.buildURL()
	if err != nil {
		return WrapperResponse{}, err
	}

	body, err := opt.bodyToBufferBody()
	if err != nil {
		return WrapperResponse{}, err
	}

	opt.defaultOptions()

	resp := &http.Response{}
	errs := ErrCollections{}

	if opt.Timeout != time.Duration(0) {
		r.client.Timeout = opt.Timeout
	}

	for i := 0; i < opt.RetryCount; i++ {
		req, err := http.NewRequestWithContext(ctx, opt.Method, URL, &body)
		if err != nil {
			return WrapperResponse{}, err
		}

		req.Header.Add(HeaderContentType, MIMEApplicationJSON)

		if len(opt.Headers) > zero {
			for key, value := range opt.Headers {
				req.Header.Add(key, value)
			}
		}

		resp, err = r.client.Do(req)
		if err != nil {
			errs.Add(err.Error())
			opt.BackoffStrategy.Next(opt.MaxBackoffDelay)
			continue
		}

		if resp.StatusCode >= http.StatusInternalServerError {
			opt.BackoffStrategy.Next(opt.MaxBackoffDelay)
			continue
		}

		errs.CleanCollection()
		break
	}

	return WrapperResponse{resp}, errs.HasError()
}
