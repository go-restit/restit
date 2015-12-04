package restit_test

import (
	"testing"

	restit "github.com/yookoala/restit/v2"
)

func TestPaths_PathOnlyURL(t *testing.T) {
	paths, err := restit.NewPaths("/foo/bar", "box", "boxes")
	if err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
		return
	}

	if want, have := "/foo/bar/box", paths.Singular(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "/foo/bar/box/hello/world", paths.Singular("hello", "world"); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if want, have := "/foo/bar/boxes", paths.Plural(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "/foo/bar/boxes/hello/world", paths.Plural("hello", "world"); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestPaths_PathOnlyURL_WithTrailSlash(t *testing.T) {
	paths, err := restit.NewPaths("/foo/bar/", "box", "boxes")
	if err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
		return
	}

	if want, have := "/foo/bar/box", paths.Singular(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "/foo/bar/box/hello/world", paths.Singular("hello", "world"); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if want, have := "/foo/bar/boxes", paths.Plural(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "/foo/bar/boxes/hello/world", paths.Plural("hello", "world"); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestPaths_FullURL(t *testing.T) {
	paths, err := restit.NewPaths("http://localhost:1234/foo/bar", "box", "boxes")
	if err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
		return
	}

	if want, have := "http://localhost:1234/foo/bar/box", paths.Singular(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "http://localhost:1234/foo/bar/box/hello/world", paths.Singular("hello", "world"); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if want, have := "http://localhost:1234/foo/bar/boxes", paths.Plural(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "http://localhost:1234/foo/bar/boxes/hello/world", paths.Plural("hello", "world"); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestPaths_FullURL_WithTrailSlash(t *testing.T) {
	paths, err := restit.NewPaths("http://localhost:1234/foo/bar/", "box", "boxes")
	if err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
		return
	}

	if want, have := "http://localhost:1234/foo/bar/box", paths.Singular(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "http://localhost:1234/foo/bar/box/hello/world", paths.Singular("hello", "world"); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if want, have := "http://localhost:1234/foo/bar/boxes", paths.Plural(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "http://localhost:1234/foo/bar/boxes/hello/world", paths.Plural("hello", "world"); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}
