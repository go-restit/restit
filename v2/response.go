package restit

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/go-restit/lzjson"
)

// Response is the generic response of any HTTP handle results
type Response interface {

	// StatusCode returns the HTTP status code
	StatusCode() int

	// Header returns the HTTP header of the response
	Header() http.Header

	// Body returns the body Reader / ReadCloser
	Body() io.Reader

	// String returns the body as string
	String() string

	// Raw returns the raw response data structure
	Raw() interface{}

	// JSON returns a JSON decoding the body
	JSON() (lzjson.Node, error)
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
func (r HTTPTestResponse) Body() io.Reader {
	return r.RawResponse.Body
}

// String implements Response
//
// Note: the body will not "survive" the read.
// you should use CacheResponse to wrap this before
// doing String().
func (r HTTPTestResponse) String() string {
	bytes, err := ioutil.ReadAll(r.Body())
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

// Raw implements Response
func (r HTTPTestResponse) Raw() interface{} {
	return r.RawResponse
}

// JSON implements Response
func (r HTTPTestResponse) JSON() (lzjson.Node, error) {
	reader := r.Body()
	node := lzjson.Decode(reader)
	return node, node.ParseError()
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
func (r HTTPResponse) Body() io.Reader {
	return r.RawResponse.Body
}

// String implements Response
//
// Note: the body will not "survive" the read.
// you should use CacheResponse to wrap this before
// doing String().
func (r HTTPResponse) String() string {
	bytes, err := ioutil.ReadAll(r.Body())
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

// Raw implements Response
func (r HTTPResponse) Raw() interface{} {
	return r.RawResponse
}

// JSON implements Response
func (r HTTPResponse) JSON() (lzjson.Node, error) {
	reader := r.Body()
	node := lzjson.Decode(reader)
	return node, node.ParseError()
}

// CacheResponse returns a new Response
// which Body() can be read repeatedly
func CacheResponse(r Response) Response {
	return &cachedResponse{r, nil}
}

// CachedResponse wraps a Response and cache the Body() result
// so it can be called over and over again
type cachedResponse struct {
	response     Response
	cachedReader *cachedReader
}

// StatusCode implements Response
func (cr cachedResponse) StatusCode() int {
	return cr.response.StatusCode()
}

// Header implements Response
func (cr cachedResponse) Header() http.Header {
	return cr.response.Header()
}

// Body implements Response
func (cr *cachedResponse) Body() io.Reader {
	if cr.cachedReader == nil {
		body, err := ioutil.ReadAll(cr.response.Body())
		if err == nil {
			err = io.EOF
		}
		cr.cachedReader = &cachedReader{
			body: body,
			err:  err,
		}
	}
	return cr.cachedReader.Copy()
}

// String implements Response
//
// Note: the body will not "survive" the read.
// you should use CacheResponse to wrap this before
// doing String().
func (cr *cachedResponse) String() string {
	bytes, err := ioutil.ReadAll(cr.Body())
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

// Raw implements Response
func (cr cachedResponse) Raw() interface{} {
	return cr.response.Raw()
}

// JSON implements Response
func (cr *cachedResponse) JSON() (lzjson.Node, error) {
	reader := cr.Body()
	node := lzjson.Decode(reader)
	return node, node.ParseError()
}

type cachedReader struct {
	body []byte
	err  error
	pos  int
}

func (cr *cachedReader) Read(b []byte) (n int, err error) {
	n = copy(b, cr.body[cr.pos:])
	cr.pos += n
	if cr.pos >= len(cr.body) {
		err = cr.err
	}
	return
}

func (cr *cachedReader) Copy() io.Reader {
	reader := &cachedReader{
		body: make([]byte, len(cr.body), (cap(cr.body)+1)*2),
		err:  cr.err,
	}
	copy(reader.body, cr.body)
	return reader
}
