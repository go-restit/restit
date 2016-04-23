package restit

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
)

// HTTPHandler implements CaseHandlerFunc
func HTTPHandler(req *http.Request) (resp Response, err error) {
	c := &http.Client{}
	rawResp, err := c.Do(req)
	if err != nil {
		return
	}

	resp = HTTPResponse{rawResp}
	return
}

// NewHTTPService create a normal HTTP service to a real
// HTTP server
func NewHTTPService(rawURL string) *Service {
	baseURL, _ := url.Parse(rawURL)
	return &Service{
		BaseURL: baseURL,
		Handler: CaseHandlerFunc(HTTPHandler),
	}
}

// HTTPTestHandler implements CaseHandlerFunc
func HTTPTestHandler(handler http.Handler) func(*http.Request) (Response, error) {
	return func(req *http.Request) (resp Response, err error) {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		resp = HTTPTestResponse{w}
		return
	}
}

// NewHTTPTestService create a dummy service based on httptest.Recorder
// as request handler
func NewHTTPTestService(rawURL string, handler http.Handler) *Service {
	baseURL, _ := url.Parse(rawURL)
	return &Service{
		BaseURL: baseURL,
		Handler: CaseHandlerFunc(HTTPTestHandler(handler)),
	}
}

// Service provides method to interact with a RESTful service
// based on the given Paths value
type Service struct {
	BaseURL *url.URL
	Handler CaseHandler
}

// NewCase creates a new Case struct with
func (s Service) NewCase(method string, payload interface{}, paths ...string) *Case {

	// formulate request URL
	requestURL, err := url.Parse(s.BaseURL.String())
	if err != nil {
		panic(err)
	}
	if len(paths) > 0 {
		requestURL.Path = path.Join(append([]string{requestURL.Path}, paths...)...)
	}

	// formulate request
	req, err := NewRequest(method, requestURL.String(), payload)
	if err != nil {
		panic(err)
	}

	return &Case{
		Request: req,
		Handler: s.Handler,
	}
}

// List sends a GET request
// to plural path and examine the result
func (s Service) List(paths ...string) *Case {
	return s.NewCase("GET", nil, paths...)
}

// Create sends a POST request (with JSON encoded payload)
// to plural path and examine the result
func (s Service) Create(payload interface{}, paths ...string) *Case {
	return s.NewCase("POST", payload, paths...)
}

// Update sends a PUT request (wtih JSON encoded payload)
// to singular path and examine the result
func (s Service) Update(payload interface{}, paths ...string) *Case {
	return s.NewCase("PUT", payload, paths...)
}

// Retrieve sends a GET request
// to singular path and examine the result
func (s Service) Retrieve(paths ...string) *Case {
	return s.NewCase("GET", nil, paths...)
}

// Delete sends a DELETE request
// to singular path and examine the result
func (s Service) Delete(paths ...string) *Case {
	return s.NewCase("DELETE", nil, paths...)
}
