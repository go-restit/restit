package example

import (
	"github.com/yookoala/restit"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExampleHandler(t *testing.T) {
	s := httptest.NewServer(ExampleHandler())
	defer s.Close()
	api := restit.Rest("Nodes", s.URL+"/api/nodes")
	var err error

	// test list
	_, err = api.List("").
		ExpectStatus(http.StatusOK).
		Run()
	if err != nil {
		t.Errorf(err.Error())
	}

	// test create
	_, err = api.Create(map[string]interface{}{
		"name": "node 4",
		"desc": "example node 4",
	}).ExpectStatus(http.StatusOK).
		Run()
	if err != nil {
		t.Errorf(err.Error())
	}

	// test retrieve
	_, err = api.Retrieve("4").
		ExpectStatus(http.StatusOK).
		Run()
	if err != nil {
		t.Errorf(err.Error())
	}

	// test update
	_, err = api.Update("4", map[string]interface{}{
		"id":   4,
		"name": "node 4 updated",
		"desc": "example node 4 updated",
	}).ExpectStatus(http.StatusAccepted).
		Run()
	if err != nil {
		t.Errorf(err.Error())
	}

	// test delete
	_, err = api.Delete("4").
		ExpectStatus(http.StatusNotFound).
		Run()
	if err != nil {
		t.Errorf(err.Error())
	}

	_ = s
}

type Resp map[string]interface{}

type Resp2 struct {
	Resp
}
