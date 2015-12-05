package example1_test

import (
	"fmt"
	"net/http"
	"testing"

	"golang.org/x/net/context"

	restit "github.com/yookoala/restit/v2"
	"github.com/yookoala/restit/v2/example/example1"
)

func TestServer(t *testing.T) {

	baseURL := "/dummy/api"
	noun := restit.NewNoun("post", "posts")
	h := example1.PostServer()(baseURL,
		noun.Singular(), noun.Plural())

	// create HTTPService to interact with
	// the server
	paths, _ := restit.NewPaths(baseURL, noun)
	service := restit.NewHTTPTestService(paths, h)

	// helper function to write expectations

	nthIs := func(name string, n int, test func(*restit.JSON) error) func(context.Context, restit.Response) error {
		return func(ctx context.Context, resp restit.Response) (err error) {
			proto, err := resp.JSON()
			if err != nil {
				return
			} else if want, have := restit.TypeObject, proto.Type(); want != have {
				ctxErr := restit.NewContextError("expected %s, got %s", want, have)
				ctxErr.Prepend("ref", "response")
				err = ctxErr
				return
			}
			list := proto.Get(name)
			if err != nil {
				return
			} else if want, have := restit.TypeArray, list.Type(); want != have {
				ctxErr := restit.NewContextError("expected %s, got %s", want, have)
				ctxErr.Prepend("ref", "response."+name)
				ctxErr.Append("response", string(proto.Raw()))
				err = ctxErr
				return
			} else if want, have := n, list.Len(); have <= want {
				ctxErr := restit.NewContextError("expected <= %#v, got %#v", want, have)
				ctxErr.Prepend("ref", "response."+name)
				ctxErr.Append("response", string(proto.Raw()))
				err = ctxErr
			}

			if nth := list.GetN(n); nth != nil {
				return test(nth)
			}

			ctxErr := restit.NewContextError("unexpected nil value")
			ctxErr.Prepend("ref", fmt.Sprintf("response.%s[%d]", name, n))
			ctxErr.Append("response", string(proto.Raw()))
			err = ctxErr
			return
		}
	}

	equals := func(p1 example1.Post) func(*restit.JSON) error {
		return func(j *restit.JSON) (err error) {
			p2 := example1.Post{}
			j.Unmarshal(&p2)
			if want, have := p1.ID, p2.ID; want != have {
				err = fmt.Errorf("expected %s, got %s", want, have)
				return
			} else if want, have := p1.Title, p2.Title; want != have {
				err = fmt.Errorf("expected %s, got %s", want, have)
				return
			} else if want, have := p1.Body, p2.Body; want != have {
				err = fmt.Errorf("expected %s, got %s", want, have)
				return
			} else if want, have := p1.Created, p2.Created; !want.Equal(have) {
				err = fmt.Errorf("expected %s, got %s", want, have)
				return
			} else if want, have := p1.Updated, p2.Updated; !want.Equal(have) {
				err = fmt.Errorf("expected %s, got %s", want, have)
				return
			}
			return
		}
	}

	// you can have different test
	isCreatedFrom := equals
	isUpdatedFrom := equals

	// generic expectations ---

	statusCodeIs := func(n int) restit.Expectation {
		return restit.Describe(
			fmt.Sprintf("status code is %d", n),
			func(ctx context.Context, resp restit.Response) (err error) {
				if want, have := n, resp.StatusCode(); want != have {
					ctxErr := restit.NewContextError("expected %d, got %d", want, have)
					ctxErr.Prepend("ref", "header status code")
					err = ctxErr
				}
				return
			})
	}

	lengthIs := func(name string, n int) restit.Expectation {
		return restit.Describe(
			fmt.Sprintf("list length is %d", n),
			func(ctx context.Context, resp restit.Response) (err error) {
				proto, err := resp.JSON()
				if err != nil {
					return
				} else if want, have := restit.TypeObject, proto.Type(); want != have {
					ctxErr := restit.NewContextError("expected %s, got %s", want, have)
					ctxErr.Prepend("ref", "response")
					err = ctxErr
					return
				}

				list := proto.Get(name)
				if err != nil {
					return
				} else if want, have := restit.TypeArray, list.Type(); want != have {
					ctxErr := restit.NewContextError("expected %s, got %s", want, have)
					ctxErr.Prepend("ref", "response."+name)
					ctxErr.Append("response", string(proto.Raw()))
					err = ctxErr
					return
				}

				if want, have := n, list.Len(); want != have {
					ctxErr := restit.NewContextError("expected %#v, got %#v", want, have)
					ctxErr.Prepend("ref", "response."+name+".length")
					err = ctxErr
					return
				}

				return
			})
	}

	nthItemIsCreatedFrom := func(name string, n int, payload example1.Post) restit.Expectation {
		return restit.Describe(
			fmt.Sprintf("item #%d is created from payload", n),
			nthIs(name, n, isCreatedFrom(payload)))
	}

	nthItemIsUpdatedFrom := func(name string, n int, payload example1.Post) restit.Expectation {
		return restit.Describe(
			fmt.Sprintf("item #%d is created from payload", n),
			nthIs(name, n, isUpdatedFrom(payload)))
	}

	nthItemIsEqualTo := func(name string, n int, payload example1.Post) restit.Expectation {
		return restit.Describe(
			fmt.Sprintf("item #%d is created from payload", n),
			nthIs(name, n, equals(payload)))
	}

	// generic expectations --- end

	// some dummy Post to test with
	p1 := example1.Post{
		ID:    "post-1",
		Title: "Some post content 1",
		Body:  "Some post body 1",
	}
	p1b := example1.Post{
		ID:    "post-1",
		Title: "Some post content 1b",
		Body:  "Some post body 1b",
	}
	p2 := example1.Post{
		ID:    "post-2",
		Title: "Some post content 2",
		Body:  "Some post body 2",
	}

	// the tests

	testList1 := service.List().
		Expect(statusCodeIs(http.StatusOK)).
		Expect(lengthIs(noun.Plural(), 0))
	if _, err := testList1.Do(); err != nil {
		if ctxErr, ok := err.(restit.ContextError); ok {
			t.Log(ctxErr.Log())
		}
		t.Errorf(err.Error())
	}

	testCreate1 := service.Create(p1).
		Expect(statusCodeIs(http.StatusOK)).
		Expect(lengthIs(noun.Plural(), 1)).
		Expect(nthItemIsCreatedFrom(noun.Plural(), 0, p1))
	if _, err := testCreate1.Do(); err != nil {
		if ctxErr, ok := err.(restit.ContextError); ok {
			t.Log(ctxErr.Log())
		}
		t.Errorf(err.Error())
	}

	testRetrieve1 := service.Retrieve(p1, p1.ID).
		Expect(statusCodeIs(http.StatusOK)).
		Expect(lengthIs(noun.Plural(), 1)).
		Expect(nthItemIsEqualTo(noun.Plural(), 0, p1))
	if _, err := testRetrieve1.Do(); err != nil {
		if ctxErr, ok := err.(restit.ContextError); ok {
			t.Log(ctxErr.Log())
		}
		t.Errorf(err.Error())
	}

	testCreate2 := service.Create(p2).
		Expect(statusCodeIs(http.StatusOK)).
		Expect(lengthIs(noun.Plural(), 1)).
		Expect(nthItemIsCreatedFrom(noun.Plural(), 0, p2))
	if _, err := testCreate2.Do(); err != nil {
		if ctxErr, ok := err.(restit.ContextError); ok {
			t.Log(ctxErr.Log())
		}
		t.Errorf(err.Error())
	}

	testUpdate1 := service.Update(p1b, p1.ID).
		Expect(statusCodeIs(http.StatusOK)).
		Expect(lengthIs(noun.Plural(), 1)).
		Expect(nthItemIsUpdatedFrom(noun.Plural(), 0, p1b))
	if _, err := testUpdate1.Do(); err != nil {
		if ctxErr, ok := err.(restit.ContextError); ok {
			t.Log(ctxErr.Log())
		}
		t.Errorf(err.Error())
	}

	testList2 := service.List().
		Expect(statusCodeIs(http.StatusOK)).
		Expect(lengthIs(noun.Plural(), 2))
	if _, err := testList2.Do(); err != nil {
		if ctxErr, ok := err.(restit.ContextError); ok {
			t.Log(ctxErr.Log())
		}
		t.Errorf(err.Error())
	}

	testDelete1 := service.Delete(p1b.ID).
		Expect(statusCodeIs(http.StatusOK)).
		Expect(lengthIs(noun.Plural(), 1)).
		Expect(nthItemIsEqualTo(noun.Plural(), 0, p1b))
	if _, err := testDelete1.Do(); err != nil {
		if ctxErr, ok := err.(restit.ContextError); ok {
			t.Log(ctxErr.Log())
		}
		t.Errorf(err.Error())
	}

	testList3 := service.List().
		Expect(statusCodeIs(http.StatusOK)).
		Expect(lengthIs(noun.Plural(), 1))
	if _, err := testList3.Do(); err != nil {
		if ctxErr, ok := err.(restit.ContextError); ok {
			t.Log(ctxErr.Log())
		}
		t.Errorf(err.Error())
	}

	testDelete2 := service.Delete(p2.ID).
		Expect(statusCodeIs(http.StatusOK)).
		Expect(lengthIs(noun.Plural(), 1)).
		Expect(nthItemIsEqualTo(noun.Plural(), 0, p2))
	if _, err := testDelete2.Do(); err != nil {
		if ctxErr, ok := err.(restit.ContextError); ok {
			t.Log(ctxErr.Log())
		}
		t.Errorf(err.Error())
	}

	testList4 := service.List().
		Expect(statusCodeIs(http.StatusOK)).
		Expect(lengthIs(noun.Plural(), 0))
	if _, err := testList4.Do(); err != nil {
		if ctxErr, ok := err.(restit.ContextError); ok {
			t.Log(ctxErr.Log())
		}
		t.Errorf(err.Error())
	}

}
