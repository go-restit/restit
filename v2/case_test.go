package restit_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/context"

	restit "github.com/yookoala/restit/v2"
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
		resp = &restit.HTTPResponse{hResp}
		go func() {
			chDone <- 1
		}()
		return
	}
	return
}

func TestCase_EmptyRequest(t *testing.T) {
	testHandler, _ := getTestHandler()
	testCase := restit.Case{
		Handler: testHandler,
	}

	resp, err := testCase.Do()
	if err == nil {
		t.Errorf("unable to trigger error")
	} else if want, have := "Request is nil", err.Error(); want != have {
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
	} else if want, have := "Handler is nil", err.Error(); want != have {
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
	} else if want, have := "expectation=1 desc=\"dummy test 2\" msg=\"dummy error\"", err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}
