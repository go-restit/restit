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

	orgl := post{}

	if v := server.NewPtr(orgl); v == nil {
		t.Errorf("unexpected nil value: %#v", v)
	} else if product, ok := v.(*post); !ok {
		t.Errorf("duplicate failed. returned %#v", v)
	} else if product == nil {
		t.Errorf("unexpected nil value: %#v", product)
	}

	if v := server.NewPtr(&orgl); v == nil {
		t.Errorf("unexpected nil value: %#v", v)
	} else if product, ok := v.(*post); !ok {
		t.Errorf("duplicate failed. returned %#v", v)
	} else if product == nil {
		t.Errorf("unexpected nil value: %#v", product)
	}

}

func TestNewSlicePtr(t *testing.T) {

	type post struct {
		ID  string
		Num int
	}

	orgl := post{}

	if v := server.NewSlicePtr(orgl); v == nil {
		t.Errorf("unexpected nil value: %#v", v)
	} else if product, ok := v.(*[]post); !ok {
		t.Errorf("duplicate failed. returned %#v", v)
	} else if product == nil {
		t.Errorf("unexpected nil value: %#v", product)
	}

	if v := server.NewSlicePtr(&orgl); v == nil {
		t.Errorf("unexpected nil value: %#v", v)
	} else if product, ok := v.(*[]post); !ok {
		t.Errorf("duplicate failed. returned %#v", v)
	} else if product == nil {
		t.Errorf("unexpected nil value: %#v", product)
	}

}

func TestNewFactory_Val(t *testing.T) {

	type post struct {
		ID  string
		Num int
	}

	orgl := post{}

	factory := server.NewFactory(orgl)

	if v := factory.Make(); v == nil {
		t.Errorf("unexpected nil value: %#v", v)
	} else if product, ok := v.(*post); !ok {
		t.Errorf("unexpected %#v", v)
	} else if product == nil {
		t.Errorf("unexpected nil value: %#v", product)
	}

	if v := factory.MakeSlice(); v == nil {
		t.Errorf("unexpected nil value: %#v", v)
	} else if product, ok := v.(*[]post); !ok {
		t.Errorf("unexpected %#v", v)
	} else if product == nil {
		t.Errorf("unexpected nil value: %#v", product)
	}

}

func TestNewFactory_Ptr(t *testing.T) {

	type post struct {
		ID  string
		Num int
	}

	orgl := post{}

	factory := server.NewFactory(&orgl)

	if v := factory.Make(); v == nil {
		t.Errorf("unexpected nil value: %#v", v)
	} else if product, ok := v.(*post); !ok {
		t.Errorf("unexpected %#v", v)
	} else if product == nil {
		t.Errorf("unexpected nil value: %#v", product)
	}

	if v := factory.MakeSlice(); v == nil {
		t.Errorf("unexpected nil value: %#v", v)
	} else if product, ok := v.(*[]post); !ok {
		t.Errorf("unexpected %#v", v)
	} else if product == nil {
		t.Errorf("unexpected nil value: %#v", product)
	}

}
