package restit

import (
	"github.com/jmcvetta/napping"
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
	Name string
	Url  string
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
		Url:    t.Url + id,
	}
	c := Case{
		Request: &r,
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
		Url:     t.Url + id,
		Payload: payload,
	}
	c := Case{
		Request: &r,
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
		Url:    t.Url + id,
	}
	c := Case{
		Request: &r,
		Session: &s,
		Tester:  t,
	}
	return &c
}
