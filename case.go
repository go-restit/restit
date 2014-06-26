package restit

import (
	"fmt"
	"github.com/jmcvetta/napping"
)

type Case struct {
	Request      *napping.Request
	Session      Session
	Expectations []Expectation
	Tester       *Tester
}

// To actually run the test case
func (c *Case) Run() (r *Result, err error) {
	res, err := c.Session.Send(c.Request)
	result := Result {
		Response: res,
	}
	r = &result

	// test each expectations
	resp := (*c).Request.Result.(Response)
	for i := 0; i < len(c.Expectations); i++ {
		err = c.Expectations[i].Test(&resp)
		if err != nil {
			err = fmt.Errorf("Failed in %s: %s",
				c.Expectations[i].Desc,
				err.Error())
			return
		}
	}

	return
}

// To run the test case and panic on error
func (c *Case) RunOrPanic() (r *Result) {
	r, err := c.Run()
	if err != nil {
		panic(err)
	}
	return
}

// Set the result to the given interface{}
func (c *Case) WithResponseAs(r interface{}) (*Case) {
	c.Request.Result = r
	return c
}

// Append Test to Expectations
// Tests if the result count equal to n
func (c *Case) ExpectResultCount(n int) (*Case) {
	c.Expectations = append(c.Expectations, Expectation{
		Desc: "Expect Result Count",
		Test: func(r *Response) (err error) {
			count := (*r).Count()
			if (count != n) {
				err = fmt.Errorf(
					"Result count is %d (expected %d)",
					count, n)
			}
			return
		},
	})
	return c
}

// Append Test to Expectations
// Tests if the item is valid
func (c *Case) ExpectResultsValid(n int) (*Case) {
	c.Expectations = append(c.Expectations, Expectation{
		Desc: fmt.Sprintf("Expect #%d result valid", n),
		Test: func(r *Response) (err error) {
			for i:=0; i<(*r).Count(); i++ {
				err = (*r).NthValid(i)
				return
			}
			return
		},
	})
	return c
}

// Append Test to Expectation
// Tests if the nth item matches the provided one
func (c *Case) ExpectResultNth(n int, b interface{}) (*Case) {
	c.Expectations = append(c.Expectations, Expectation{
		Desc: fmt.Sprintf("Expect #%d result valid", n),
		Test: func(r *Response) (err error) {
			a, err := (*r).GetNth(n)
			if err != nil {
				return
			}
			err = (*r).Match(a, b)
			return
		},
	})
	return c
}

// Append Custom Test to Expectation
// Allow user to inject user defined tests
func (c *Case) ExpectResultToPass(
	desc string,
	test func(*Response)(error)) (*Case) {
	c.Expectations = append(c.Expectations, Expectation{
		Desc: desc,
		Test: test,
	})
	return c
}

// Expection to the response in a Case
type Expectation struct{
	Desc string
	Test func(*Response) (error)
}

// Test Result of a Case
type Result struct {
	Response *napping.Response
}

// Wrap the napping.Session in Session
// to make unit testing easier
type Session interface {
	Send(*napping.Request) (*napping.Response, error)
}
