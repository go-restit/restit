package restit_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	restit "github.com/yookoala/restit/v2"
)

// dummyTestSuite to test dummyServiceHandler below
func dummyTestSuite(service *restit.Service, paths restit.Paths) (err error) {

	type post struct {
		ID      string    `json:"id"`
		Name    string    `json:"name"`
		Created time.Time `json:"created"`
		Updated time.Time `json:"updated"`
	}

	err = func() (err error) {
		randPath := []string{RandString(20), RandString(20)}
		testCase := service.List(randPath...)
		if req := testCase.Request; req == nil {
			err = fmt.Errorf("[list][request] empty")
		} else if want, have := paths.Plural(randPath...), req.URL.String(); want != have {
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
		} else if want, have := paths.Plural(), req.URL.String(); want != have {
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
		} else if want, have := paths.Plural("foo", "bar"), req.URL.String(); want != have {
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
		} else if want, have := paths.Singular(postID), req.URL.String(); want != have {
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
		testCase := service.Retrieve(post{ID: postID, Name: RandString(20)}, postID)
		if req := testCase.Request; req == nil {
			err = fmt.Errorf("[retrieve][request] empty")
		} else if want, have := paths.Singular(postID), req.URL.String(); want != have {
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
		} else if want, have := paths.Singular(postID), req.URL.String(); want != have {
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
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	}
	return handler
}

func TestService(t *testing.T) {
	noun := restit.NewNoun("post", "posts")
	paths, err := restit.NewPaths("/dummyAPI", noun)
	if err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	}

	service := &restit.Service{
		Paths: paths,
	}

	// run test suite against the dummy service
	if err := dummyTestSuite(service, paths); err != nil {
		t.Errorf("failed running test suite: %s", err.Error())
	}

}

func TestHTTPService(t *testing.T) {

	// create a real HTTP server with httptest.Server
	handler := dummyServiceHandler()
	testServer := httptest.NewServer(handler)
	defer testServer.Close()

	noun := restit.NewNoun("post", "posts")
	paths, err := restit.NewPaths(testServer.URL+"/dummyAPI", noun)
	if err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	}

	// create HTTPService
	service := restit.NewHTTPService(paths)

	// run test suite against the dummy service
	if err := dummyTestSuite(service, paths); err != nil {
		t.Errorf("failed running test suite: %s", err.Error())
	}
}

func TestHTTPTestService(t *testing.T) {
	handler := dummyServiceHandler()

	noun := restit.NewNoun("post", "posts")
	paths, err := restit.NewPaths("/dummyAPI", noun)
	if err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	}

	// create HTTPTestService
	service := restit.NewHTTPTestService(paths, handler)

	// run test suite against the dummy service
	if err := dummyTestSuite(service, paths); err != nil {
		t.Errorf("failed running test suite: %s", err.Error())
	}
}
