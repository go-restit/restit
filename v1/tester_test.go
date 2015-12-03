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
	"log"
	"os"
	"testing"
)

func Test_Tester_Log(t *testing.T) {
	lt := log.New(os.Stdout, "[TRACE] ",
		log.Ldate|log.Ltime|log.Lshortfile)
	le := log.New(os.Stdout, "[ERROR] ",
		log.Ldate|log.Ltime|log.Lshortfile)
	test := Rest("Dummy", "http://foobar/dummy").
		LogTraceTo(lt).
		LogErrTo(le)
	if test.Trace != lt {
		t.Error("Failed to set LogTrace with LogTraceTo")
	}
	if test.Err != le {
		t.Error("Failed to set LogErr with LogErrTo")
	}
}

func Test_Tester_List_Emtpy(t *testing.T) {
	u := "http://foobar/dummy"
	expected := u

	test := Rest("Dummy", u).
		List()
	if test.Request.Url != expected {
		t.Errorf("Unexpected generated URL. Expected \"%s\" but get \"%s\"",
			u,
			test.Request.Url)
	}
}

func Test_Tester_List_Single(t *testing.T) {
	u := "http://foobar/dummy"
	expected := u + "/hello"

	test := Rest("Dummy", u).
		List("hello")
	if test.Request.Url != expected {
		t.Errorf("Unexpected generated URL. Expected \"%s\" but get \"%s\"",
			expected,
			test.Request.Url)
	}
}

func Test_Tester_List_Multiple(t *testing.T) {
	u := "http://foobar/dummy"
	expected := u + "/hello/world"

	test := Rest("Dummy", u).
		List("hello", "world")
	if test.Request.Url != "http://foobar/dummy/hello/world" {
		t.Errorf("Unexpected generated URL. Expected \"%s\" but get \"%s\"",
			expected,
			test.Request.Url)
	}
}

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
