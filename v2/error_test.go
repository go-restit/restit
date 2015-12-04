package restit_test

import (
	"testing"

	restit "github.com/yookoala/restit/v2"
)

func TestContextError_EmptyString(t *testing.T) {
	err := restit.NewContextError("")
	var err2 error = err
	if want, have := "", err2.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestContextError_NoMessage(t *testing.T) {
	err := restit.NewContextError("")
	err.Delete("message")
	var err2 error = err
	if want, have := "error", err2.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestContextError_AppendPrepend(t *testing.T) {
	err := restit.NewContextError("dummy error")
	var err2 error = err
	if want, have := "dummy error", err2.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	err.Append("foo", "bar")
	if want, have := `message="dummy error" foo="bar"`, err.Log(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	err.Prepend("hello", "world")
	if want, have := `hello="world" message="dummy error" foo="bar"`, err.Log(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	err.Prepend("foo", "bar")
	if want, have := `foo="bar" hello="world" message="dummy error"`, err.Log(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}
