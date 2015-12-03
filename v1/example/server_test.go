package example

import (
	restit "github.com/yookoala/restit/v1"

	"net/http"
	"net/http/httptest"
	"testing"
)

type testNode struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Desc string `json:"string"`
}

func testResponse() restit.Response {
	r := restit.NewResponse("nodes", testNode{})
	r.SetValidator(func(in interface{}) error {
		return nil
	})
	r.SetMatcher(func(a interface{}, b interface{}) error {
		return nil
	})
	return r
}

func TestExampleHandler(t *testing.T) {
	s := httptest.NewServer(ExampleHandler())
	defer s.Close()
	api := restit.Rest("Nodes", s.URL+"/api/nodes")
	var err error

	// test list
	_, err = api.List().
		WithResponseAs(testResponse()).
		ExpectStatus(http.StatusOK).
		ExpectResultCount(3).
		Run()
	if err != nil {
		t.Errorf(err.Error())
	}

	// test create
	_, err = api.Create(map[string]interface{}{
		"name": "node 4",
		"desc": "example node 4",
	}).WithResponseAs(testResponse()).
		ExpectStatus(http.StatusOK).
		ExpectResultCount(1).
		Run()
	if err != nil {
		t.Errorf(err.Error())
	}

	// test retrieve
	_, err = api.Retrieve("4").
		WithResponseAs(testResponse()).
		ExpectStatus(http.StatusOK).
		ExpectResultCount(1).
		Run()

	if err != nil {
		t.Errorf(err.Error())
	}

	// test update
	_, err = api.Update("4", map[string]interface{}{
		"id":   4,
		"name": "node 4 updated",
		"desc": "example node 4 updated",
	}).WithResponseAs(testResponse()).
		ExpectStatus(http.StatusOK).
		ExpectResultCount(1).
		Run()
	if err != nil {
		t.Errorf(err.Error())
	}

	// test delete
	_, err = api.Delete("4").
		WithResponseAs(testResponse()).
		ExpectStatus(http.StatusOK).
		ExpectResultCount(1).
		Run()
	if err != nil {
		t.Errorf(err.Error())
	}

	_ = s
}
