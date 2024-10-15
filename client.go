package ravenTree

import (
	"context"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/gojektech/heimdall/v6/httpclient"
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

type raven struct{}

// NewRavensTree creates and returns a new instance of the RavenTree interface.
//
// It returns an instance of `raven`, a private struct that implements the `Tree` interface.
// This function acts as a constructor for the `Tree` implementation.
//
// Returns:
// - Tree: An object that implements the RavenTree interface.
func NewRavensTree() Tree {
	return &raven{}
}

const (
	zero = 0

	HeaderContentType = "Content-Type"

	// MIMEApplicationJSON JavaScript Object Notation (JSON) https://www.rfc-editor.org/rfc/rfc8259
	MIMEApplicationJSON = "application/json"
)

func (u *raven) SendRaven(ctx context.Context, opt *Options) (WrapperResponse, error) {
	log.SetOutput(os.Stdout)

	var opts []httpclient.Option

	if opt.Timeout > zero {
		opts = append(opts, httpclient.WithHTTPTimeout(opt.Timeout))
	}

	if opt.RetryCount > zero {
		opts = append(opts, httpclient.WithRetryCount(opt.RetryCount))
	}

	URL, err := opt.buildURL()
	if err != nil {
		return WrapperResponse{}, err
	}

	body, err := opt.bodyToBufferBody()
	if err != nil {
		return WrapperResponse{}, err
	}

	client := httpclient.NewClient(opts...)
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

	log.Info("Raven Send to...",
		"URL", URL,
		"Method", opt.Method,
		"Body", opt.Body,
	)

	resp, err := client.Do(req)

	return WrapperResponse{resp}, err
}
