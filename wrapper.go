package ravenTree

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type WrapperResponse struct {
	*http.Response
}

// ParseBodyTo parses the HTTP response body into the provided destination (dest).
//
// This method reads the response body (w.Body) as a byte slice, then attempts
// to unmarshal the JSON content into the provided `dest` interface.
//
// The response body (w.Body) is then restored so it can be read again later, if necessary.
//
// Parameters:
//   - dest (interface{}): A pointer to the destination where the parsed JSON
//     from the response body will be stored. It can be a struct or map that matches
//     the JSON structure.
//
// Returns:
// - error: Returns an error if the reading of the body or the unmarshaling process fails.
func (w *WrapperResponse) ParseBodyTo(dest interface{}) error {
	bodyBytes, _ := io.ReadAll(w.Body)
	w.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return json.Unmarshal(bodyBytes, dest)
}

// ParseBodyToString reads the HTTP response body and returns it as a string.
//
// This method reads the entire response body (w.Body) into a byte slice, converts it
// to a string, and then restores the body so it can be read again later if needed.
//
// Returns:
// - string: The response body as a string.
func (w *WrapperResponse) ParseBodyToString() string {
	bodyBytes, _ := io.ReadAll(w.Body)
	w.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return string(bodyBytes)
}
