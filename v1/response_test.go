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
)

type dummy struct {
	Name string
}

type dummyResponse struct {
	Dummies []dummy
}

func (r *dummyResponse) Count() int {
	return len(r.Dummies)
}

func (r *dummyResponse) NthValid(n int) error {
	if r.Dummies[n].Name == "" {
		return fmt.Errorf("All dummies should have a name")
	}
	return nil
}

func (r *dummyResponse) GetNth(n int) (interface{}, error) {
	nth := r.Dummies[n]
	return &nth, nil
}

func (r *dummyResponse) Match(a interface{}, b interface{}) error {
	if a.(*dummy).Name != b.(*dummy).Name {
		return fmt.Errorf("Mismatch")
	}
	return nil
}

func (r *dummyResponse) Reset() {
	r.Dummies = make([]dummy, 0)
}

// dummy session with dummy send sequence
type dummySession struct {
}

func (s dummySession) Send(req *napping.Request) (
	res *napping.Response, err error) {
	var resv napping.Response
	res = &resv
	ptrResult := (*req).Result.(*dummyResponse)
	(*ptrResult).Dummies = []dummy{
		dummy{
			Name: "Hello Dummy",
		},
	}
	return
}

// dummy session for nilResp test
type dummyNilSession struct {
}

func (s dummyNilSession) Send(req *napping.Request) (
	res *napping.Response, err error) {
	res = new(napping.Response)
	// placeholder only
	return
}
