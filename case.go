// Copyright (c) 2014 Yeung Shu Hung (Koala Yeung)
//
//  This file is part of RESTit.
//
//  RESTit is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  RESTit is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//
//  Use of this source code is governed by the GPL v3 license. A copy
//  of the licence can be found in the LICENSE.md file along with RESTit.
//  If not, see <http://www.gnu.org/licenses/>.

package restit

import (
	"fmt"
	"github.com/jmcvetta/napping"
	"net/http"
	"path"
	"runtime"
)

type Case struct {
	Request      *napping.Request
	Name         string
	Session      Session
	Expectations []Expectation
	Tester       *Tester
	Result       *Result
}

// short hand to initialize and enable defaults
// even if tester is not present
func (c *Case) InitForRun() *Case {

	// Tester must be there to provide log
	if c.Tester == nil {
		c.Tester = new(Tester)
	}

	// load default logging behaviour
	c.Tester.LogDefault()

	// c.Request must be a valid napping.Request
	// this will be sent through napping.Send
	if c.Request == nil {
		c.Request = new(napping.Request)
	}

	// if result is not specified,
	// substitute nilResp as response type
	if c.Request.Result == nil {
		c.Request.Result = new(nilResp)
	}

	// trigger reset
	r, ok := c.Request.Result.(Response)
	if !ok {
		panic(fmt.Errorf(
			"The provided response %T does not implement restit.Response",
			r))
	}
	r.Reset()

	return c
}

// To actually run the test case
func (c *Case) Run() (r *Result, err error) {

	// setup default tester
	c.InitForRun()

	// get caller information
	_, file, line, _ := runtime.Caller(2)

	// send request
	res, err := c.Session.Send(c.Request)
	c.Tester.Trace.Printf("[%s:%d][%s][%s] Raw Response: \"%s\"\n",
		path.Base(file),
		line,
		c.Tester.Name,
		c.Name,
		res.RawText())
	result := Result{
		Response: res,
	}
	c.Result = &result
	r = &result

	// test each expectations
	resp := (*c).Request.Result.(Response)
	for i := 0; i < len(c.Expectations); i++ {
		err = c.Expectations[i].Test(resp)
		if err != nil {
			err = fmt.Errorf("[%s:%d][%s][%s][%s] "+
				"Failed: \"%s\"",
				path.Base(file),
				line,
				c.Tester.Name,
				c.Name,
				c.Expectations[i].Desc,
				err.Error())
			c.Tester.Err.Println(err.Error())
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

// Add a header parameter
func (c *Case) AddHeader(key, value string) *Case {
	if c.Request.Header == nil {
		c.Request.Header = &http.Header{}
	}
	c.Request.Header.Add(key, value)
	return c
}

// Set the result to the given interface{}
func (c *Case) WithResponseAs(r Response) *Case {
	c.Request.Result = r
	return c
}

// Set the query parameter
func (c *Case) WithParams(p *napping.Params) *Case {
	c.Request.Params = p
	return c
}

// Append Test to Expectations
// Tests if the result count equal to n
func (c *Case) ExpectResultCount(n int) *Case {
	c.Expectations = append(c.Expectations, Expectation{
		Desc: "Test Result Count",
		Test: func(r Response) (err error) {
			count := r.Count()
			if count != n {
				err = fmt.Errorf(
					"Result count is %d "+
						"(expected %d)",
					count, n)
			}
			return
		},
	})
	return c
}

// Append Test to Expectations
// Tests if the result count not equal to n
func (c *Case) ExpectResultCountNot(n int) *Case {
	c.Expectations = append(c.Expectations, Expectation{
		Desc: "Test Result Count",
		Test: func(r Response) (err error) {
			count := r.Count()
			if count == n {
				err = fmt.Errorf(
					"Result count is %d "+
						"(expected %d)",
					count, n)
			}
			return
		},
	})
	return c
}

// Append Test to Expectations
// Tests if the item is valid
func (c *Case) ExpectResultsValid() *Case {
	c.Expectations = append(c.Expectations, Expectation{
		Desc: "Test Results Valid",
		Test: func(r Response) (err error) {
			for i := 0; i < r.Count(); i++ {
				err = r.NthValid(i)
				if err != nil {
					err = fmt.Errorf(
						"Item %d invalid: %s",
						i, err.Error())
					return
				}
			}
			return
		},
	})
	return c
}

// Append Test to Expectation
// Tests if the nth item matches the provided one
func (c *Case) ExpectResultNth(n int, b interface{}) *Case {
	c.Expectations = append(c.Expectations, Expectation{
		Desc: fmt.Sprintf("Test #%d Result Valid", n),
		Test: func(r Response) (err error) {
			a, err := r.GetNth(n)
			if err != nil {
				return
			}
			err = r.Match(a, b)
			return
		},
	})
	return c
}

// Append Test to Expectations
// Tests if the response status is
func (c *Case) ExpectStatus(ec int) *Case {
	c.Expectations = append(c.Expectations, Expectation{
		Desc: "Test Status Code",
		Test: func(r Response) (err error) {
			rc := c.Result.Response.Status()
			if rc != ec {
				err = fmt.Errorf("Status code is %d (expected %d)",
					rc, ec)
			}
			return
		},
	})
	return c
}

// Append Custom Test to Expectation
// Allow user to inject user defined tests
func (c *Case) ExpectResultsToPass(
	desc string, test func(Response) error) *Case {
	c.Expectations = append(c.Expectations, Expectation{
		Desc: desc,
		Test: test,
	})
	return c
}

// Expection to the response in a Case
type Expectation struct {
	Desc string
	Test func(Response) error
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
