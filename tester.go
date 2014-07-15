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
	"github.com/jmcvetta/napping"
	"log"
	"os"
)

// Create a tester for an API entry point
// name    string human-readable name of the entry point
// baseUrl string RESTful API base url
func Rest(name string, url string) *Tester {
	t := Tester{
		Name: name,
		Url:  url,
	}
	return &t
}

// Tester represents an ordinary RESTful entry point
type Tester struct {
	Name  string
	Url   string
	Trace *log.Logger
	Err   *log.Logger
}

// setup default log environment
func (t *Tester) LogDefault() *Tester {
	if t.Trace == nil {
		t.LogTraceTo(log.New(os.Stdout,
			"[TRACE] ",
			log.Ldate|log.Ltime|log.Lshortfile))
	}
	if t.Err == nil {
		t.LogErrTo(log.New(os.Stderr,
			"[ERROR] ",
			log.Ldate|log.Ltime|log.Lshortfile))
	}
	return t
}

// Add trace logger
func (t *Tester) LogTraceTo(l *log.Logger) *Tester {
	t.Trace = l
	return t
}

// Add error logger
func (t *Tester) LogErrTo(l *log.Logger) *Tester {
	t.Err = l
	return t
}

// Create Case to Create something with the payload
func (t *Tester) Create(payload interface{}) *Case {
	s := napping.Session{}
	r := napping.Request{
		Method:  "POST",
		Url:     t.Url,
		Payload: payload,
	}
	c := Case{
		Request: &r,
		Name:    "Create",
		Session: &s,
		Tester:  t,
	}
	return &c
}

// Create Case to Retrieve something with the id string
func (t *Tester) Retrieve(id string) *Case {
	s := napping.Session{}
	r := napping.Request{
		Method: "GET",
		Url:    t.Url + "/" + id,
	}
	c := Case{
		Request: &r,
		Name:    "Retrieve",
		Session: &s,
		Tester:  t,
	}
	return &c
}

// Create Case to Update something of the id with the payload
func (t *Tester) Update(
	id string, payload interface{}) *Case {
	s := napping.Session{}
	r := napping.Request{
		Method:  "PUT",
		Url:     t.Url + "/" + id,
		Payload: payload,
	}
	c := Case{
		Request: &r,
		Name:    "Update",
		Session: &s,
		Tester:  t,
	}
	return &c
}

// Create Case to Delete something of the id
func (t *Tester) Delete(id string) *Case {
	s := napping.Session{}
	r := napping.Request{
		Method: "DELETE",
		Url:    t.Url + "/" + id,
	}
	c := Case{
		Request: &r,
		Name:    "Delete",
		Session: &s,
		Tester:  t,
	}
	return &c
}
