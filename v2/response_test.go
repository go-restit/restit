package restit_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	restit "github.com/go-restit/restit/v2"
)

func TestResponse_httptest(t *testing.T) {
	msg := RandString(10)
	requestID := RandString(10)

	w := httptest.NewRecorder()
	w.Header().Set("X-Request-ID", requestID)
	w.Write([]byte(msg))
	w.WriteHeader(http.StatusOK)
	w.Flush()

	var resp restit.Response = &restit.HTTPTestResponse{RawResponse: w}
	if want, have := http.StatusOK, resp.StatusCode(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := requestID, resp.Header().Get("X-Request-ID"); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if b, err := ioutil.ReadAll(resp.Body()); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	} else if want, have := msg, string(b); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if want, have := w, resp.Raw(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestResponse_httptest_JSON(t *testing.T) {
	msg := dummyJSONStr()
	requestID := RandString(10)

	w := httptest.NewRecorder()
	w.Header().Set("X-Request-ID", requestID)
	w.Write([]byte(msg))
	w.WriteHeader(http.StatusOK)
	w.Flush()

	var resp restit.Response = &restit.HTTPTestResponse{RawResponse: w}
	if want, have := http.StatusOK, resp.StatusCode(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := requestID, resp.Header().Get("X-Request-ID"); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	// test JSON result
	if j, err := resp.JSON(); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	} else if err := dummyJSONTest(j); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	}
}

func TestResponse_http(t *testing.T) {
	msg := RandString(10)
	requestID := RandString(10)

	// a simple repeater to test with
	var repeater http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {

		if r == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("nil request"))
			return
		}

		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("method not supported"))
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Add("x-Request-ID", requestID)
		w.WriteHeader(http.StatusAccepted)
		w.Write(b)
	}
	srv := httptest.NewServer(repeater)

	// POST request to the test server
	rawResp, err := http.Post(srv.URL, "text/plain", strings.NewReader(msg))
	if err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	}

	var resp restit.Response = &restit.HTTPResponse{RawResponse: rawResp}
	if want, have := http.StatusAccepted, resp.StatusCode(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := requestID, resp.Header().Get("X-Request-ID"); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if b, err := ioutil.ReadAll(resp.Body()); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	} else if want, have := msg, string(b); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if want, have := rawResp, resp.Raw(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestResponse_http_JSON(t *testing.T) {
	msg := dummyJSONStr()
	requestID := RandString(10)

	// a simple repeater to test with
	var repeater http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {

		if r == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("nil request"))
			return
		}

		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("method not supported"))
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Add("x-Request-ID", requestID)
		w.WriteHeader(http.StatusAccepted)
		w.Write(b)
	}
	srv := httptest.NewServer(repeater)

	// POST request to the test server
	rawResp, err := http.Post(srv.URL, "text/plain", strings.NewReader(msg))
	if err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	}

	var resp restit.Response = &restit.HTTPResponse{RawResponse: rawResp}
	if want, have := http.StatusAccepted, resp.StatusCode(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := requestID, resp.Header().Get("X-Request-ID"); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	// test JSON result
	if j, err := resp.JSON(); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	} else if err := dummyJSONTest(j); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	}
}

func TestResponse_cached(t *testing.T) {

	msg := RandString(10)
	requestID := RandString(10)

	w := httptest.NewRecorder()
	w.Header().Set("X-Request-ID", requestID)
	w.Write([]byte(msg))
	w.WriteHeader(http.StatusOK)
	w.Flush()

	resp := restit.CacheResponse(&restit.HTTPTestResponse{RawResponse: w})

	// read once
	if result, err := ioutil.ReadAll(resp.Body()); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	} else if want, have := msg, string(result); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	// read twice
	if result, err := ioutil.ReadAll(resp.Body()); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	} else if want, have := msg, string(result); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

}

func TestResponse_cached_JSON(t *testing.T) {
	msg := dummyJSONStr()
	requestID := RandString(10)

	w := httptest.NewRecorder()
	w.Header().Set("X-Request-ID", requestID)
	w.Write([]byte(msg))
	w.WriteHeader(http.StatusOK)
	w.Flush()

	resp := restit.CacheResponse(&restit.HTTPTestResponse{RawResponse: w})

	if want, have := http.StatusOK, resp.StatusCode(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := requestID, resp.Header().Get("X-Request-ID"); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	// test JSON result once
	if j, err := resp.JSON(); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	} else if err := dummyJSONTest(j); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	}

	// test JSON result twice
	if j, err := resp.JSON(); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	} else if err := dummyJSONTest(j); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	}
}
