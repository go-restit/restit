package restit

import (
	"net/http"
	"net/http/httptest"
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
func NewHTTPService(paths Paths) *Service {
	return &Service{
		Paths:   paths,
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
func NewHTTPTestService(paths Paths, handler http.Handler) *Service {
	return &Service{
		Paths:   paths,
		Handler: CaseHandlerFunc(HTTPTestHandler(handler)),
	}
}

// Service provides method to interact with a RESTful service
// based on the given Paths value
type Service struct {
	Paths   Paths
	Handler CaseHandler
}

// List sends a GET request
// to plural path and examine the result
func (s Service) List(v ...string) *Case {
	req, err := NewRequest("GET", s.Paths.Plural(v...), nil)
	if err != nil {
		panic(err)
	}

	return &Case{
		Request: req,
		Handler: s.Handler,
	}
}

// Create sends a POST request (with JSON encoded payload)
// to plural path and examine the result
func (s Service) Create(payload interface{}, v ...string) *Case {
	req, err := NewRequest("POST", s.Paths.Plural(v...), payload)
	if err != nil {
		panic(err)
	}

	return &Case{
		Request: req,
		Handler: s.Handler,
	}
}

// Update sends a PUT request (wtih JSON encoded payload)
// to singular path and examine the result
func (s Service) Update(payload interface{}, v ...string) *Case {
	req, err := NewRequest("PUT", s.Paths.Singular(v...), payload)
	if err != nil {
		panic(err)
	}

	return &Case{
		Request: req,
		Handler: s.Handler,
	}
}

// Retrieve sends a GET request
// to singular path and examine the result
func (s Service) Retrieve(v ...string) *Case {
	req, err := NewRequest("GET", s.Paths.Singular(v...), nil)
	if err != nil {
		panic(err)
	}

	return &Case{
		Request: req,
		Handler: s.Handler,
	}
}

// Delete sends a DELETE request
// to singular path and examine the result
func (s Service) Delete(v ...string) *Case {
	req, err := NewRequest("DELETE", s.Paths.Singular(v...), nil)
	if err != nil {
		panic(err)
	}

	return &Case{
		Request: req,
		Handler: s.Handler,
	}
}
