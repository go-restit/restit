package restit

import (
	"encoding/json"
	"errors"
	"testing"
)

var testJsonList = `{
	"code": 200,
	"status": "success",
	"items": [
		{
			"id": 1,
			"title": "item 1",
			"desc": "item 1 desc"
		},
		{
			"id": 2,
			"title": "item 2",
			"desc": "item 2 desc",
			"Hello": "world",
			"Foo": "invisible"
		},
		{
			"id": 3,
			"title": "item 3",
			"desc": "item 2 desc"
		}
	]
}`

type testItem struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Hello string
	Foo   string `json:"-"`
}

func TestDefaultResponse(t *testing.T) {
	var r Response = &DefaultResponse{}
	_ = r
	t.Log("*DefaultResponse implements Response")
}

func TestDefaultResponseResp(t *testing.T) {

	// define a response with item type as testItem
	// and item list field named "items"
	p := NewResponse("items", testItem{})

	// set validator
	p.SetValidator(func(in interface{}) error {
		a := in.(testItem)
		if a.Title == "" {
			return errors.New("Incorrect item")
		}
		return nil
	})

	// set matcher
	p.SetMatcher(func(a interface{}, b interface{}) (err error) {
		aItem := a.(testItem)
		bItem := b.(testItem)
		if aItem.Title != bItem.Title {
			err = errors.New("Title mismatch")
		} else if aItem.Desc != bItem.Desc {
			err = errors.New("Desc mismatch")
		} else if aItem.Hello != bItem.Hello {
			err = errors.New("Hello mismatch")
		} else if aItem.Foo != bItem.Foo {
			err = errors.New("Foo mismatch")
		}
		return
	})

	// decode testJsonList with *DefaultResponse
	json.Unmarshal([]byte(testJsonList), &p)

	// test result
	nth, err := p.GetNth(1)
	if err != nil {
		t.Errorf(err.Error())
	}
	item := nth.(testItem)
	t.Logf("Successfully cast item to testItem:\n%#v", item)

	// test if items match the json input
	if err := p.Match(item, testItem{
		ID:    2,
		Title: "item 2",
		Desc:  "item 2 desc",
		Hello: "world",
	}); err != nil {
		t.Errorf("Failed to set JSON value correctly")
	} else {
		t.Logf("JSON value set correctly")
	}
}
