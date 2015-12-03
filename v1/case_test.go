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
	"net/url"
	"testing"
)

func Test_Case_WithResponseAs_nil(t *testing.T) {

	c := Case{
		Session: new(dummyNilSession),
	}
	c.Run()

}

func Test_Case_WithResponseAs_invalid(t *testing.T) {

	pass := false

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered in f: %#v", r)
			pass = true
			t.Logf("Invalid response type triggers error")
		}
	}()

	resp := "some invalid response"
	r := napping.Request{}
	r.Result = &resp
	c := Case{
		Request: &r,
		Session: new(dummyNilSession),
	}
	c.RunOrPanic()

	if pass != true {
		t.Errorf("Invalid response type does not raise error")
	}

}

func Test_Case_WithResponseAs_reset(t *testing.T) {

	// test reset with filled response
	r := dummyResponse{
		Dummies: []dummy{
			dummy{},
			dummy{},
			dummy{},
		},
	}

	c := Case{
		Request: &napping.Request{},
		Session: new(dummyNilSession),
	}
	t.Logf("resp: %#v", r)
	c.WithResponseAs(&r)
	c.Run()

	if len(r.Dummies) != 0 {
		t.Errorf("Case.Run() fails to trigger Response.Reset")
	}

}

func Test_Case_WithErrorAs(t *testing.T) {

	r := &napping.Request{}
	resp := dummyResponse{
		Dummies: []dummy{
			dummy{},
			dummy{},
			dummy{},
		},
	}

	c := Case{
		Request: r,
		Session: new(dummyNilSession),
	}
	c.WithErrorAs(&resp)

	if resp.Count() != 3 {
		t.Logf("Incorrect assumption, resp.Count() is not 3")
	}

	// set the resp through address in r.Error
	var a *dummyResponse
	a = (r.Error).(*dummyResponse)
	(*a).Dummies = []dummy{
		dummy{},
		dummy{},
	}

	// test if the original error variable changed as expected
	if resp.Count() != 2 {
		t.Errorf("Case.WithErrorAs() failed to set napping.Response.Error")
	}

}

func Test_Case_AddHeader(t *testing.T) {

	r := napping.Request{}
	c := Case{
		Request: &r,
	}
	c.AddHeader("foo", "bar")
	if r.Header.Get("foo") != "bar" {
		t.Error("AddHeader failed to add a header parameter")
	}

}

func Test_Case_WithParams(t *testing.T) {

	r := napping.Request{
		Params: &url.Values{},
	}
	c := Case{
		Request: &r,
	}
	p := url.Values{
		"hello": []string{"world"},
	}
	c.WithParams(&p)
	expRes := "world"
	if res := c.Request.Params.Get("hello"); res != expRes {
		t.Errorf("WithParams failed to set the parameter. "+
			"expected: %#v, got: %#v", expRes, res)
	}

}

func Test_Case_ExpectResultCount_0(t *testing.T) {

	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{},
	}
	c.ExpectResultCount(0)
	err := c.Expectations[0].Test(&r)
	if err != nil {
		t.Error("Unable to pass with a valid count 0")
	}

}

func Test_Case_ExpectResultCount_n(t *testing.T) {

	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{
			dummy{},
			dummy{},
			dummy{},
		},
	}
	c.ExpectResultCount(3)
	err := c.Expectations[0].Test(&r)
	if err != nil {
		t.Error("Unable to pass with a valid count n")
	}

}

func Test_Case_ExpectResultCount_err(t *testing.T) {

	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{
			dummy{},
			dummy{},
		},
	}
	c.ExpectResultCount(3)
	err := c.Expectations[0].Test(&r)
	if err == nil {
		t.Error("Unable to detect count mismatch")
	}

}

func Test_Case_ExpectResultCountNot_0(t *testing.T) {

	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{
			dummy{},
		},
	}
	c.ExpectResultCountNot(0)
	err := c.Expectations[0].Test(&r)
	if err != nil {
		t.Error("Unable to pass with a valid count not 0")
	}

}

func Test_Case_ExpectResultCountNot_n(t *testing.T) {

	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{},
	}
	c.ExpectResultCountNot(3)
	err := c.Expectations[0].Test(&r)
	if err != nil {
		t.Error("Unable to pass with a valid count not n")
	}

}

func Test_Case_ExpectResultCountNot_err(t *testing.T) {

	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{
			dummy{},
			dummy{},
		},
	}
	c.ExpectResultCountNot(2)
	err := c.Expectations[0].Test(&r)
	if err == nil {
		t.Error("Unable to detect count match")
	}

}

func Test_Case_ExpectResultsValid(t *testing.T) {

	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{
			dummy{Name: "Hello"},
		},
	}
	c.ExpectResultsValid()
	err := c.Expectations[0].Test(&r)
	if err != nil {
		t.Error("Unable to pass valid item")
	}

}

func Test_Case_ExpectResultsValid_err(t *testing.T) {

	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{
			dummy{},
		},
	}
	c.ExpectResultsValid()
	err := c.Expectations[0].Test(&r)
	if err == nil {
		t.Error("Unable to identify invalid item")
	}

}

func Test_Case_ExpectResultNth_match(t *testing.T) {

	d := dummy{Name: "Unique Dummy"}
	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{
			dummy{},
			d,
			dummy{},
		},
	}
	c.ExpectResultNth(1, &d)
	err := c.Expectations[0].Test(&r)
	if err != nil {
		t.Error("Expect Nth result to match given payload")
	}

}

func Test_Case_ExpectResultNth_err(t *testing.T) {

	d := dummy{Name: "Unique Dummy"}
	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{
			dummy{},
			d,
			dummy{},
		},
	}
	c.ExpectResultNth(0, &d)
	err := c.Expectations[0].Test(&r)
	if err == nil {
		t.Error("Expect Nth result to mis-match given payload")
	}

}

func Test_Case_ExpectResultsToPass_pass(t *testing.T) {

	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{
			dummy{},
		},
	}
	c.ExpectResultsToPass("Custom Test to pass",
		func(Response) error {
			return nil
		})
	err := c.Expectations[0].Test(&r)
	if err != nil {
		t.Error("Expect custom tests to pass")
	}

}

func Test_Case_ExpectResultsToPass_err(t *testing.T) {

	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{
			dummy{},
		},
	}
	c.ExpectResultsToPass("Custom Test to fail",
		func(Response) error {
			return fmt.Errorf("Some error")
		})
	err := c.Expectations[0].Test(&r)
	if err == nil {
		t.Error("Expect custom tests to fail")
	}

}
