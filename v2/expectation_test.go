package restit_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/go-restit/lzjson"
	"github.com/go-restit/restit/v2"
	"golang.org/x/net/context"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

func TestStatusCodeIs(t *testing.T) {
	statusCode := rand.Int()
	resp := restit.HTTPResponse{
		RawResponse: &http.Response{
			StatusCode: statusCode,
		},
	}

	emptyCtx := context.Background()
	if err := restit.StatusCodeIs(statusCode).Do(emptyCtx, resp); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if err := restit.StatusCodeIs(statusCode-1).Do(emptyCtx, resp); err == nil {
		t.Errorf("expected to trigger error but didn't")
	} else if want, have := fmt.Sprintf("expected %d, got %d", statusCode-1, statusCode), err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestLengthIs(t *testing.T) {

	resp := restit.CacheResponse(restit.HTTPResponse{
		RawResponse: &http.Response{
			Body: ioutil.NopCloser(strings.NewReader(`{"hello": ["hello 1", "hello 2"]}`)),
		},
	})

	emptyCtx := context.Background()
	if err := restit.LengthIs("hello", 2).Do(emptyCtx, resp); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if err := restit.LengthIs("hello", 1).Do(emptyCtx, resp); err == nil {
		t.Errorf("should trigger error but didn't")
	} else if want, have := `expected "hello" to be length 1, got 2`, err.Error(); want != have {
		t.Errorf("\nexpected: %s\ngot:      %s", want, have)
	}
}

func TestNthTest(t *testing.T) {

	resp := restit.CacheResponse(restit.HTTPResponse{
		RawResponse: &http.Response{
			Body: io.NopCloser(strings.NewReader(`{"hello": ["hello 1", "hello 2"]}`)),
		},
	})

	emptyCtx := context.Background()

	test1Run := false
	test1 := restit.Nth(0).Of("hello").Is(restit.DescribeJSON("hello 1", func(node lzjson.Node) (err error) {
		if want, have := "hello 1", node.String(); want != have {
			err = fmt.Errorf("expected %#v, got %#v", want, have)
		}
		test1Run = true
		return
	}))
	if err := test1.Do(emptyCtx, resp); err != nil {
		t.Errorf("unexpected test1 error: %s", err)
	}
	if test1Run == false {
		t.Errorf("test1 never run")
	}

	test2Run := false
	test2 := restit.Nth(1).Of("hello").Is(restit.DescribeJSON("hello 1", func(node lzjson.Node) (err error) {
		if want, have := "hello 1", node.String(); want != have {
			err = fmt.Errorf("expected %#v, got %#v", want, have)
		}
		test2Run = true
		return
	}))
	if err := test2.Do(emptyCtx, resp); err == nil {
		t.Errorf("test2 should trigger error but didn't")
	} else if want, have := `failed "hello 1" (expected "hello 1", got "hello 2")`, err.Error(); want != have {
		t.Errorf("\nexpected: %s\ngot:      %s", want, have)
	}
	if test2Run == false {
		t.Errorf("test2 never run")
	}

}
