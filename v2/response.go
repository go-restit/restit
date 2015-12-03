package restit

import (
	"io"
	"net/http"
	"net/http/httptest"
)

// Response is the generic response of any HTTP handle results
type Response interface {

	// StatusCode returns the HTTP status code
	StatusCode() int

	// Header returns the HTTP header of the response
	Header() http.Header

	// Body returns the body Reader / ReadCloser
	Body() io.ReadCloser

	// Raw returns the raw response data structure
	Raw() interface{}

	// JSON returns a JSON decoding the body
	JSON() (*JSON, error)
}

// dummyCloser warps a Reader and implements ReadCloser
type dummyCloser struct {
	io.Reader
}

// Close implements ReadCloser
func (dummyCloser) Close() error {
	return nil
}

// HTTPTestResponse wraps a *httptest.ResponseRecorder
// and implements Response interface for it
type HTTPTestResponse struct {
	RawResponse *httptest.ResponseRecorder
}

// StatusCode implements Response
func (r HTTPTestResponse) StatusCode() int {
	return r.RawResponse.Code
}

// Header implements Response
func (r HTTPTestResponse) Header() http.Header {
	return r.RawResponse.Header()
}

// Body implements Response
func (r HTTPTestResponse) Body() io.ReadCloser {
	return dummyCloser{r.RawResponse.Body}
}

// Raw implements Response
func (r HTTPTestResponse) Raw() interface{} {
	return r.RawResponse
}

// JSON implements Response
func (r HTTPTestResponse) JSON() (*JSON, error) {
	reader := r.Body()
	defer reader.Close()
	return ReadJSON(reader)
}

// HTTPResponse wraps a *http.Response
// and implements Response interface for it
type HTTPResponse struct {
	RawResponse *http.Response
}

// StatusCode implements Response
func (r HTTPResponse) StatusCode() int {
	return r.RawResponse.StatusCode
}

// Header implements Response
func (r HTTPResponse) Header() http.Header {
	return r.RawResponse.Header
}

// Body implements Response
func (r HTTPResponse) Body() io.ReadCloser {
	return r.RawResponse.Body
}

// Raw implements Response
func (r HTTPResponse) Raw() interface{} {
	return r.RawResponse
}

// JSON implements Response
func (r HTTPResponse) JSON() (*JSON, error) {
	reader := r.Body()
	defer reader.Close()
	return ReadJSON(reader)
}
