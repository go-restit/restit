package restit_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	restit "github.com/go-restit/restit/v2"
)

// dummyTestSuite to test dummyServiceHandler below
func dummyTestSuite(service *restit.Service, baseURL string) (err error) {

	type post struct {
		ID      string    `json:"id"`
		Name    string    `json:"name"`
		Created time.Time `json:"created"`
		Updated time.Time `json:"updated"`
	}

	err = func() (err error) {
		randPath := RandString(20)
		randPayload := RandString(20)
		testCase := service.NewCase("GET", randPayload, "/"+randPath)
		resp, err := testCase.Do()
		if err != nil {
			return
		}
		randPayloadJSON, _ := json.Marshal(randPayload)

		body, _ := ioutil.ReadAll(resp.Body())
		expected := fmt.Sprintf(`{"URL":"/dummyAPI/posts/%s","method":"GET","payload":%#v}`+"\n",
			randPath, string(randPayloadJSON))

		//`{"URL":"/dummyAPI/posts/` + randPath +
		//`","method":"GET","payload":"` + randPayloadJSON + `"}` + "\n"

		if want, have := expected, string(body); want != have {
			err = fmt.Errorf("[NewCase][request]\nexpected: %s\ngot:      %s", want, have)
		}
		return
	}()
	if err != nil {
		return
	}

	err = func() (err error) {
		randPath := RandString(20)
		testCase := service.List(randPath)
		if req := testCase.Request; req == nil {
			err = fmt.Errorf("[list][request] empty")
		} else if want, have := baseURL+"/"+randPath, req.URL.String(); want != have {
			err = fmt.Errorf("[list][request.URL] expected %#v, got %#v", want, have)
		} else if want, have := "GET", req.Method; want != have {
			err = fmt.Errorf("[list][request.Method] expected %#v, got %#v", want, have)
		}
		return
	}()
	if err != nil {
		return
	}

	err = func() (err error) {
		testCase := service.Create(post{ID: RandString(20), Name: RandString(20)})
		if req := testCase.Request; req == nil {
			err = fmt.Errorf("[create][request] empty")
		} else if want, have := baseURL, req.URL.String(); want != have {
			err = fmt.Errorf("[create][request.URL] expected %#v, got %#v", want, have)
		} else if want, have := "POST", req.Method; want != have {
			err = fmt.Errorf("[create][request.Method] expected %#v, got %#v", want, have)
		}
		return
	}()
	if err != nil {
		return
	}

	err = func() (err error) {
		testCase := service.Create(post{ID: RandString(20), Name: RandString(20)}, "foo", "bar")
		if req := testCase.Request; req == nil {
			err = fmt.Errorf("[create][request] empty")
		} else if want, have := baseURL+"/foo/bar", req.URL.String(); want != have {
			err = fmt.Errorf("[create][request.URL] expected %#v, got %#v", want, have)
		} else if want, have := "POST", req.Method; want != have {
			err = fmt.Errorf("[create][request.Method] expected %#v, got %#v", want, have)
		}
		return
	}()
	if err != nil {
		return
	}

	err = func() (err error) {
		postID := RandString(20)
		testCase := service.Update(post{ID: postID, Name: RandString(20)}, postID)
		if req := testCase.Request; req == nil {
			err = fmt.Errorf("[update][request] empty")
		} else if want, have := baseURL+"/"+postID, req.URL.String(); want != have {
			err = fmt.Errorf("[update][request.URL] expected %#v, got %#v", want, have)
		} else if want, have := "PUT", req.Method; want != have {
			err = fmt.Errorf("[update][request.Method] expected %#v, got %#v", want, have)
		}
		return
	}()
	if err != nil {
		return
	}

	err = func() (err error) {
		postID := RandString(20)
		testCase := service.Retrieve(postID)
		if req := testCase.Request; req == nil {
			err = fmt.Errorf("[retrieve][request] empty")
		} else if want, have := baseURL+"/"+postID, req.URL.String(); want != have {
			err = fmt.Errorf("[retrieve][request.URL] expected %#v, got %#v", want, have)
		} else if want, have := "GET", req.Method; want != have {
			err = fmt.Errorf("[retrieve][request.Method] expected %#v, got %#v", want, have)
		}
		return
	}()
	if err != nil {
		return
	}

	err = func() (err error) {
		postID := RandString(20)
		testCase := service.Delete(postID)
		if req := testCase.Request; req == nil {
			err = fmt.Errorf("[delete][request] empty")
		} else if want, have := baseURL+"/"+postID, req.URL.String(); want != have {
			err = fmt.Errorf("[delete][request.URL] expected %#v, got %#v", want, have)
		} else if want, have := "DELETE", req.Method; want != have {
			err = fmt.Errorf("[delete][request.Method] expected %#v, got %#v", want, have)
		}
		return
	}()
	if err != nil {
		return
	}

	return
}

// dummyServiceHandler provides a dummy service to test with
func dummyServiceHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		resp := json.NewEncoder(w)
		b, _ := ioutil.ReadAll(r.Body)
		resp.Encode(map[string]interface{}{
			"URL":     r.URL.String(),
			"method":  r.Method,
			"payload": string(b),
		})
	})
}

func TestService(t *testing.T) {
	handler := dummyServiceHandler()
	baseURL, _ := url.Parse("/dummyAPI/posts")
	service := &restit.Service{
		BaseURL: baseURL,
		Handler: restit.CaseHandlerFunc(restit.HTTPTestHandler(handler)),
	}

	// run test suite against the dummy service
	if err := dummyTestSuite(service, baseURL.String()); err != nil {
		t.Errorf("failed running test suite: %s", err.Error())
	}

}

func TestHTTPService(t *testing.T) {

	// create a real HTTP server with httptest.Server
	handler := dummyServiceHandler()
	baseURL := "/dummyAPI/posts"
	testServer := httptest.NewServer(handler)
	defer testServer.Close()

	// create HTTPService
	service := restit.NewHTTPService(testServer.URL + baseURL)

	// run test suite against the dummy service
	if err := dummyTestSuite(service, testServer.URL+baseURL); err != nil {
		t.Errorf("failed running test suite: %s", err.Error())
	}
}

func TestHTTPTestService(t *testing.T) {
	handler := dummyServiceHandler()
	baseURL := "/dummyAPI/posts"

	// create HTTPTestService
	service := restit.NewHTTPTestService(baseURL, handler)

	// run test suite against the dummy service
	if err := dummyTestSuite(service, baseURL); err != nil {
		t.Errorf("failed running test suite: %s", err.Error())
	}
}
