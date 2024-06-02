package restit_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/context"

	restit "github.com/go-restit/restit/v2"
)

type ctxKey int

const (
	dummyKey ctxKey = iota
)

func getTestHandler() (fn restit.CaseHandlerFunc, done <-chan int) {
	chDone := make(chan int)
	done = chDone
	fn = func(req *http.Request) (resp restit.Response, err error) {
		hResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       req.Body,
		}
		resp = &restit.HTTPResponse{RawResponse: hResp}
		go func() {
			chDone <- 1
		}()
		return
	}
	return
}

func TestCase_AddHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "/foo/bar", nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	testCase := &restit.Case{
		Request: req,
	}
	if want, have := "", testCase.Request.Header.Get("hello"); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	testCase.AddHeader("hello", "world")
	if want, have := "world", testCase.Request.Header.Get("hello"); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestCase_AddQuery(t *testing.T) {
	req, err := http.NewRequest("GET", "/foo/bar?foo=bar", nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	testCase := &restit.Case{
		Request: req,
	}
	if want, have := "", testCase.Request.URL.Query().Get("hello"); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	testCase.AddQuery("hello", "world")
	if want, have := "world", testCase.Request.URL.Query().Get("hello"); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
		t.Logf("testCase.Request.URL: %#v", testCase.Request.URL)
	}
}

func TestCase_EmptyRequest(t *testing.T) {
	testHandler, _ := getTestHandler()
	testCase := restit.Case{
		Handler: testHandler,
	}

	resp, err := testCase.Do()
	if err == nil {
		t.Errorf("unable to trigger error")
	} else if want, have := "case.Request is nil", err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	var nilResp restit.Response
	if want, have := nilResp, resp; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestCase_EmptyHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/foo/bar", nil)
	if err != nil {
		t.Errorf("unexpected error %#v", err.Error())
		return
	}

	testCase := restit.Case{
		Request: req,
	}
	resp, err := testCase.Do()
	if err == nil {
		t.Errorf("unable to trigger error")
	} else if want, have := "case.Handler is nil", err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	var nilResp restit.Response
	if want, have := nilResp, resp; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestCase_EmptyContext(t *testing.T) {
	req, err := http.NewRequest("GET", "/foo/bar",
		strings.NewReader(dummyJSONStr()))
	if err != nil {
		t.Errorf("unexpected error %#v", err.Error())
		return
	}

	expRun := false
	testHandler, handlerDone := getTestHandler()
	testCase := restit.Case{
		Request: req,
		Context: context.WithValue(
			context.Background(), dummyKey, "hello dummy"),
		Handler: testHandler,
		Expectations: []restit.Expectation{
			restit.Describe("dummy test 1",
				func(ctx context.Context, resp restit.Response) (err error) {
					expRun = true
					return
				}),
		},
	}

	// run the case
	resp, err := testCase.Do()
	if err != nil {
		t.Errorf("unexpected error %#v", err.Error())
		return
	}

	// test the response *JSON
	if j, err := resp.JSON(); err != nil {
		t.Errorf("unexpected error %#v", err.Error())
	} else if err := dummyJSONTest(j); err != nil {
		t.Errorf("unexpected error %#v", err.Error())
	}

	// test if the handler is run
	d, err := time.ParseDuration("1s")
	if err != nil {
		t.Errorf("unexpected error %#v", err.Error())
		return
	}
	timeout := time.After(d)
	select {
	case <-handlerDone:
	case <-timeout:
		t.Error("handler did not run")
		return
	}

	// test if the expectation is run
	if want, have := true, expRun; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestCase_WithContext(t *testing.T) {
	req, err := http.NewRequest("GET", "/foo/bar",
		strings.NewReader(dummyJSONStr()))
	if err != nil {
		t.Errorf("unexpected error %#v", err.Error())
		return
	}

	expRun := false
	expHasContext := false

	testHandler, handlerDone := getTestHandler()
	testCase := restit.Case{
		Request: req,
		Context: context.WithValue(
			context.Background(), dummyKey, "hello dummy"),
		Handler: testHandler,
		Expectations: []restit.Expectation{
			restit.Describe("dummy test 1",
				func(ctx context.Context, resp restit.Response) (err error) {
					expRun = true
					if want, have := "hello dummy", ctx.Value(dummyKey); want != have {
						err = fmt.Errorf("ctx.Value(dummyKey) expected %#v, got %#v", want, have)
						return
					}
					expHasContext = true
					return
				}),
		},
	}

	// run the case
	resp, err := testCase.Do()
	if err != nil {
		t.Errorf("unexpected error %#v", err.Error())
		return
	}

	if want, have := true, expHasContext; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	// test the response *JSON
	if j, err := resp.JSON(); err != nil {
		t.Errorf("unexpected error %#v", err.Error())
	} else if err := dummyJSONTest(j); err != nil {
		t.Errorf("unexpected error %#v", err.Error())
	}

	// test if the handler is run
	d, err := time.ParseDuration("1s")
	if err != nil {
		t.Errorf("unexpected error %#v", err.Error())
		return
	}
	timeout := time.After(d)
	select {
	case <-handlerDone:
	case <-timeout:
		t.Error("handler did not run")
		return
	}

	// test if the expectation is run
	if want, have := true, expRun; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestCase_ExpectationErr(t *testing.T) {
	req, err := http.NewRequest("GET", "/foo/bar",
		strings.NewReader(dummyJSONStr()))
	if err != nil {
		t.Errorf("unexpected error %#v", err.Error())
		return
	}

	testHandler, _ := getTestHandler()
	testCase := restit.Case{
		Request: req,
		Context: context.WithValue(
			context.Background(), dummyKey, "hello dummy"),
		Handler: testHandler,
		Expectations: []restit.Expectation{
			restit.Describe("dummy test 1",
				func(ctx context.Context, resp restit.Response) (err error) {
					return
				}),
			restit.Describe("dummy test 2",
				func(ctx context.Context, resp restit.Response) (err error) {
					return fmt.Errorf("dummy error")
				}),
			restit.Describe("dummy test 3",
				func(ctx context.Context, resp restit.Response) (err error) {
					return fmt.Errorf("never should have run this")
				}),
		},
	}

	if _, err := testCase.Do(); err == nil {
		t.Error("failed to trigger error")
	} else if want, have := "dummy error", err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if ctxErr, ok := err.(restit.ContextError); !ok {
		t.Errorf("expected restit.ContextError, got %#v", err)
	} else if want, have := `expectation=1 desc="dummy test 2" message="dummy error"`, ctxErr.Log(); want != have {
		t.Errorf("\nexpected: %s\ngot:      %s", want, have)
	}
}

func TestCase_ExpectationContextErr(t *testing.T) {
	req, err := http.NewRequest("GET", "/foo/bar",
		strings.NewReader(dummyJSONStr()))
	if err != nil {
		t.Errorf("unexpected error %#v", err.Error())
		return
	}

	testHandler, _ := getTestHandler()
	testCase := restit.Case{
		Request: req,
		Context: context.WithValue(
			context.Background(), dummyKey, "hello dummy"),
		Handler: testHandler,
		Expectations: []restit.Expectation{
			restit.Describe("dummy test 1",
				func(ctx context.Context, resp restit.Response) (err error) {
					return
				}),
			restit.Describe("dummy test 2",
				func(ctx context.Context, resp restit.Response) (err error) {
					ctxErr := restit.NewContextError("dummy error")
					ctxErr.Append("foo", "bar")
					err = ctxErr
					return
				}),
			restit.Describe("dummy test 3",
				func(ctx context.Context, resp restit.Response) (err error) {
					return fmt.Errorf("never should have run this")
				}),
		},
	}

	if _, err := testCase.Do(); err == nil {
		t.Error("failed to trigger error")
	} else if want, have := "dummy error", err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if ctxErr, ok := err.(restit.ContextError); !ok {
		t.Errorf("expected restit.ContextError, got %#v", err)
	} else if want, have := `expectation=1 desc="dummy test 2" message="dummy error" foo="bar"`, ctxErr.Log(); want != have {
		t.Errorf("\nexpected: %s\ngot:      %s", want, have)
	}
}

func TestCase_Describe(t *testing.T) {
	fnRun := false
	str := RandString(20)
	fn := func(ctx context.Context, resp restit.Response) (err error) {
		fnRun = true
		return
	}

	desc := restit.Describe(str, fn)
	desc.Do(context.Background(), nil)
	if want, have := str, desc.Desc(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := true, fnRun; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestCase_ModifyCase(t *testing.T) {
	var urlResult string

	testVal := RandString(20)
	r, err := http.NewRequest("GET", "http://example.com/foo/bar", nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	c := &restit.Case{
		Request: r,
	}
	c.
		AddQuery("hello", testVal).
		ModifyCase(func(c *restit.Case) *restit.Case {
			urlResult = c.Request.URL.String()
			return c
		})

	if want, have := "http://example.com/foo/bar?hello="+testVal, urlResult; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestCase_Expect(t *testing.T) {
	c := &restit.Case{}
	str := RandString(20)
	c.Expect(restit.Describe(str, nil))
	if want, have := 1, len(c.Expectations); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := str, c.Expectations[0].Desc(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}
