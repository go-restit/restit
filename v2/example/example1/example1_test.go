package example1_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-restit/lzjson"
	restit "github.com/go-restit/restit/v2"
	"github.com/go-restit/restit/v2/example/example1"
)

func TestServer(t *testing.T) {

	baseURL := "/dummy/api"
	h := example1.PostServer()(baseURL, "post", "posts")

	// create HTTPService to interact with
	// the server
	service := restit.NewHTTPTestService(baseURL, h)

	// helper function to write expectations

	equals := func(p1 example1.Post) func(lzjson.Node) error {
		return func(j lzjson.Node) (err error) {
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

	testList1 := service.List("/posts").
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 0))
	if _, err := testList1.Do(); err != nil {
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

	testCreate1 := service.Create(p1, "/posts").
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 1)).
		Expect(restit.Nth(0).Of("posts").Is(restit.DescribeJSON(
			"item #0 returned is created from payload", isCreatedFrom(p1))))
	if _, err := testCreate1.Do(); err != nil {
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

	testRetrieve1 := service.Retrieve("/post/" + p1.ID).
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 1)).
		Expect(restit.Nth(0).Of("posts").Is(restit.DescribeJSON(
			"item #0 retrieved is equal to p1", isCreatedFrom(p1))))
	if _, err := testRetrieve1.Do(); err != nil {
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

	testCreate2 := service.Create(p2, "/posts").
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 1)).
		Expect(restit.Nth(0).Of("posts").Is(restit.DescribeJSON(
			"item #0 retrieved is equal to p2", isCreatedFrom(p2))))
	if _, err := testCreate2.Do(); err != nil {
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

	testUpdate1 := service.Update(p1b, "/post/"+p1.ID).
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 1)).
		Expect(restit.Nth(0).Of("posts").Is(restit.DescribeJSON(
			"item #0 retrieved is equal to p1b", isUpdatedFrom(p1b))))
	if _, err := testUpdate1.Do(); err != nil {
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

	testList2 := service.List("/posts").
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 2))
	if _, err := testList2.Do(); err != nil {
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

	testDelete1 := service.Delete("/post/" + p1b.ID).
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 1)).
		Expect(restit.Nth(0).Of("posts").Is(restit.DescribeJSON(
			"item #0 returned is equal to p1b", isUpdatedFrom(p1b))))
	if _, err := testDelete1.Do(); err != nil {
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

	testList3 := service.List("/posts").
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 1))
	if _, err := testList3.Do(); err != nil {
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

	testDelete2 := service.Delete("/post/" + p2.ID).
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 1)).
		Expect(restit.Nth(0).Of("posts").Is(restit.DescribeJSON(
			"item #0 returned is equal to p2", equals(p2))))
	if _, err := testDelete2.Do(); err != nil {
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

	testList4 := service.List("/posts").
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 0))
	if _, err := testList4.Do(); err != nil {
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

}
