package server_test

import (
	"testing"

	"github.com/yookoala/restit/v2/example/server"
)

func TestNewPtr(t *testing.T) {

	type post struct {
		ID  string
		Num int
	}

	p1 := post{}

	if p2v := server.NewPtr(p1); p2v == nil {
		t.Errorf("unexpected nil value: %#v", p2v)
	} else if _, ok := p2v.(*post); !ok {
		t.Errorf("duplicate failed. returned %#v", p2v)
	}

	if p2v := server.NewPtr(&p1); p2v == nil {
		t.Errorf("unexpected nil value: %#v", p2v)
	} else if _, ok := p2v.(*post); !ok {
		t.Errorf("duplicate failed. returned %#v", p2v)
	}

}

func TestNewSlicePtr(t *testing.T) {

	type post struct {
		ID  string
		Num int
	}

	p1 := post{}

	if p2v := server.NewSlicePtr(p1); p2v == nil {
		t.Errorf("unexpected nil value: %#v", p2v)
	} else if _, ok := p2v.(*[]post); !ok {
		t.Errorf("duplicate failed. returned %#v", p2v)
	}

	if p2v := server.NewSlicePtr(&p1); p2v == nil {
		t.Errorf("unexpected nil value: %#v", p2v)
	} else if _, ok := p2v.(*[]post); !ok {
		t.Errorf("duplicate failed. returned %#v", p2v)
	}

}

func TestNewFactory_Val(t *testing.T) {

	type post struct {
		ID  string
		Num int
	}

	p1 := post{}

	f := server.NewFactory(p1)
	if v := f.Make(); v == nil {
		t.Errorf("unexpected nil value: %#v", v)
	} else if _, ok := v.(*post); !ok {
		t.Errorf("unexpected %#v", v)
	}

	if v := f.MakeSlice(); v == nil {
		t.Errorf("unexpected nil value: %#v", v)
	} else if _, ok := v.(*[]post); !ok {
		t.Errorf("unexpected %#v", v)
	}

}

func TestNewFactory_Ptr(t *testing.T) {

	type post struct {
		ID  string
		Num int
	}

	p1 := post{}

	f := server.NewFactory(&p1)
	if v := f.Make(); v == nil {
		t.Errorf("unexpected nil value: %#v", v)
	} else if _, ok := v.(*post); !ok {
		t.Errorf("unexpected %#v", v)
	}

	if v := f.MakeSlice(); v == nil {
		t.Errorf("unexpected nil value: %#v", v)
	} else if _, ok := v.(*[]post); !ok {
		t.Errorf("unexpected %#v", v)
	}

}
