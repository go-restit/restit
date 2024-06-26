package restit

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"
)

// CaseHandler runs a given request and return a response interface
type CaseHandler interface {
	Handle(req *http.Request) (resp Response, err error)
}

// CaseHandlerFunc implements CaseHandler
type CaseHandlerFunc func(req *http.Request) (resp Response, err error)

// Handle implements CaseHandler
func (fn CaseHandlerFunc) Handle(req *http.Request) (resp Response, err error) {
	return fn(req)
}

// Case contain all information of a single test case
type Case struct {
	Request      *http.Request
	Context      context.Context
	Handler      CaseHandler
	Expectations []Expectation
}

// AddHeader add given header key-value pair to request
func (c *Case) AddHeader(key, value string) *Case {
	c.Request.Header.Add(key, value)
	return c
}

// AddQuery add given query key-value pair to request
func (c *Case) AddQuery(key, value string) *Case {
	q := c.Request.URL.Query()
	q.Add(key, value)
	c.Request.URL.RawQuery = q.Encode()
	return c
}

// ModifyCase allows user do whatever to the case (even
// rewrite a new one) without interrupting the chaining.
func (c *Case) ModifyCase(fn func(c *Case) *Case) *Case {
	return fn(c)
}

// Expect appends an expectation to the Case
func (c *Case) Expect(exp Expectation) *Case {
	c.Expectations = append(c.Expectations, exp)
	return c
}

// Do actually execute the case with CaseHandler
// and return the result
func (c Case) Do() (resp Response, err error) {
	if c.Request == nil {
		return nil, fmt.Errorf("case.Request is nil")
	}
	if c.Context == nil {
		c.Context = context.Background()
	}
	if c.Handler == nil {
		return nil, fmt.Errorf("case.Handler is nil")
	}

	// do the request
	resp, err = c.Handler.Handle(c.Request)
	if err != nil {
		return
	}

	// wrap resulting Response with cachedResponse
	resp = CacheResponse(resp)

	// run all expectations
	for i, expect := range c.Expectations {
		if err = expect.Do(c.Context, resp); err != nil {
			var cErr ContextError
			var ok bool
			if cErr, ok = err.(ContextError); !ok {
				cErr = NewContextError(err.Error())
			}
			cErr.Prepend("desc", expect.Desc())
			cErr.Prepend("expectation", i)
			err = cErr
			return
		}
	}

	return
}

// Expectation stores procedure to run
// the expection on the result
type Expectation interface {

	// Desc returns the description string of the expectation
	Desc() string

	// Do runs the expected Response against the context
	// then return any error if failed
	Do(ctx context.Context, resp Response) (err error)
}

// defaultExpect is the default implementation
// of Expectation
type defaultExpect struct {
	desc string
	do   func(ctx context.Context, resp Response) (err error)
}

// Desc implements Expectation
func (de defaultExpect) Desc() string {
	return de.desc
}

// Do implements Expectation
func (de defaultExpect) Do(ctx context.Context, resp Response) error {
	return de.do(ctx, resp)
}

// Describe returns a default implemtation of Expectation
// by the given description and do function
func Describe(desc string, do func(ctx context.Context, resp Response) (err error)) Expectation {
	return defaultExpect{
		desc: desc,
		do:   do,
	}
}
