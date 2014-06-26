package restit

import (
	"github.com/jmcvetta/napping"
)

type Case struct {
	Request      *napping.Request
	Session      Session
	Expectations []Expectation
	Tester       *Tester
}

// To run the test case
func (c *Case) Run() (r *Result, err error) {
	res, err := c.Session.Send(c.Request)
	// TODO: test each expectations
	r.Response = res
	return
}

// To actually run the test case
func (c *Case) RunOrPanic() (r *Result) {
	r, err := c.Run()
	if err != nil {
		panic(err)
	}
	return
}

// Expection to the response in a Case
type Expectation struct {
}

// Test Result of a Case
type Result struct {
	Response *napping.Response
}

type Session interface {
	Send(*napping.Request) (*napping.Response, error)
}
