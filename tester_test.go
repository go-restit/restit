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
	"testing"
)

func Test_Tester_Create(t *testing.T) {
	a := dummy{
		Name: "Hello Dummy",
	}
	r := dummyResponse{}
	test := Rest("Dummy", "http://foobar/dummy").
		Create(&a).
		WithResponseAs(&r).
		ExpectResultCount(1).
		ExpectResultsValid().
		ExpectResultNth(0, &a)
	test.Session = dummySession{}
	_, err := test.Run()
	if err != nil {
		t.Error(err)
	}
}

func Test_Tester_Retrieve(t *testing.T) {
	Rest("Dummy", "http://foobar/dummy").
		Retrieve("some_id")
}

func Test_Tester_Update(t *testing.T) {
	a := dummy{}
	Rest("Dummy", "http://foobar/dummy").
		Update("some_id", &a)
}

func Test_Tester_Delete(t *testing.T) {
	Rest("Dummy", "http://foobar/dummy").
		Delete("some_id")
}
