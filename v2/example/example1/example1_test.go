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

	// creates a http.Handler of a dummy RESTful API service
	// that handle requests to "/dummy/api/posts" and "/dummy/api/post/{id}"
	h := example1.PostServer()("/dummy/api", "post", "posts")

	// create HTTPService to interact the http.Handler through
	// httptest.ResponseWriter testing routine
	service := restit.NewHTTPTestService("/dummy/api", h)

	// helper function to write expectations
	equals := func(p1 example1.Post) func(lzjson.Node) error {
		return func(j lzjson.Node) (err error) {
			p2 := example1.Post{}
			j.Unmarshal(&p2)
			if want, have := p1.ID, p2.ID; want != have {
				err = fmt.Errorf("ID expected %s, got %s", want, have)
				return
			} else if want, have := p1.Title, p2.Title; want != have {
				err = fmt.Errorf("Title expected %s, got %s", want, have)
				return
			} else if want, have := p1.Body, p2.Body; want != have {
				err = fmt.Errorf("Body expected %s, got %s", want, have)
				return
			} else if want, have := p1.Created, p2.Created; !want.Equal(have) {
				err = fmt.Errorf("Created expected %s, got %s", want, have)
				return
			} else if want, have := p1.Updated, p2.Updated; !want.Equal(have) {
				err = fmt.Errorf("Updated expected %s, got %s", want, have)
				return
			}
			return
		}
	}

	// helper function to write expectations
	isPatchedWith := func(patch example1.Post) func(lzjson.Node) error {
		return func(j lzjson.Node) (err error) {
			stored := example1.Post{}
			j.Unmarshal(&stored)
			if want, have := patch.ID, stored.ID; want != "" && want != have {
				err = fmt.Errorf("ID expected %s, got %s", want, have)
				return
			} else if want, have := patch.Title, stored.Title; want != "" && want != have {
				err = fmt.Errorf("Title expected %s, got %s", want, have)
				return
			} else if want, have := patch.Body, stored.Body; want != "" && want != have {
				err = fmt.Errorf("Body expected %s, got %s", want, have)
				return
			} else if want, have := patch.Created, stored.Created; !want.IsZero() && !want.Equal(have) {
				err = fmt.Errorf("Created expected %s, got %s", want, have)
				return
			} else if want, have := patch.Updated, stored.Updated; !want.IsZero() && !want.Equal(have) {
				err = fmt.Errorf("Updated expected %s, got %s", want, have)
				return
			}
			return
		}
	}

	// we're reusing `equals` here but you may have different test function
	// for different case
	isCreatedFrom := equals
	isUpdatedFrom := equals

	// test listing before creating anything (should be empty)
	testList1 := service.List("/posts").
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 0))
	if _, err := testList1.Do(); err != nil {
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

	// test create and retrieve p1
	p1 := example1.Post{
		ID:    "post-1",
		Title: "Some post content 1",
		Body:  "Some post body 1",
	}
	testCreate1 := service.Create(p1, "/posts").
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 1)).
		Expect(restit.Nth(0).Of("posts").Is(restit.DescribeJSON(
			"item #0 returned is created from payload", isCreatedFrom(p1))))
	if resp, err := testCreate1.Do(); err != nil {
		t.Logf("raw response:\n%s\n", resp.Raw())
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

	testRetrieve1 := service.Retrieve("/post/" + p1.ID).
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 1)).
		Expect(restit.Nth(0).Of("posts").Is(restit.DescribeJSON(
			"item #0 retrieved is equal to p1", isCreatedFrom(p1))))
	if resp, err := testRetrieve1.Do(); err != nil {
		t.Logf("raw response:\n%s\n", resp.Raw())
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

	// test create and retrieve p2
	p2 := example1.Post{
		ID:    "post-2",
		Title: "Some post content 2",
		Body:  "Some post body 2",
	}
	testCreate2 := service.Create(p2, "/posts").
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 1)).
		Expect(restit.Nth(0).Of("posts").Is(restit.DescribeJSON(
			"item #0 retrieved is equal to p2", isCreatedFrom(p2))))
	if resp, err := testCreate2.Do(); err != nil {
		t.Logf("raw response:\n%s\n", resp.Raw())
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

	testRetrieve2 := service.Retrieve("/post/" + p2.ID).
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 1)).
		Expect(restit.Nth(0).Of("posts").Is(restit.DescribeJSON(
			"item #0 retrieved is equal to p2", isCreatedFrom(p2))))
	if resp, err := testRetrieve2.Do(); err != nil {
		t.Logf("raw response:\n%s\n", resp.Raw())
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

	// test updating p1 with p1b
	p1b := example1.Post{
		ID:    "post-1",
		Title: "Some post content 1b",
		Body:  "Some post body 1b",
	}
	testUpdate1 := service.Update(p1b, "/post/"+p1.ID).
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 1)).
		Expect(restit.Nth(0).Of("posts").Is(restit.DescribeJSON(
			"item #0 retrieved is updated from p1b", isUpdatedFrom(p1b))))
	if resp, err := testUpdate1.Do(); err != nil {
		t.Logf("raw response:\n%s\n", resp.Raw())
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

	// test listing after all the creation
	testList2 := service.List("/posts").
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 2)).
		Expect(restit.Nth(0).Of("posts").Is(restit.DescribeJSON(
			"item #0 retrieved is updated from p1b", isUpdatedFrom(p1b)))).
		Expect(restit.Nth(1).Of("posts").Is(restit.DescribeJSON(
			"item #1 retrieved is created from p2", isCreatedFrom(p2))))
	if resp, err := testList2.Do(); err != nil {
		t.Logf("raw response:\n%s\n", resp.Raw())
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

	// test patching p1 with p1c
	p1c := example1.Post{
		Title: "Some post content 1c",
		Body:  "Some post body 1c",
	}
	testPatch1 := service.Patch(p1c, "/post/"+p1.ID).
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 1)).
		Expect(restit.Nth(0).Of("posts").Is(restit.DescribeJSON(
			"item #0 retrieved is patched with p1c", isPatchedWith(p1c))))
	if resp, err := testPatch1.Do(); err != nil {
		t.Logf("raw response:\n%s\n", resp.Raw())
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

	// test deleting p1
	testDelete1 := service.Delete("/post/" + p1.ID).
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 1)).
		Expect(restit.Nth(0).Of("posts").Is(restit.DescribeJSON(
			"item #0 retrieved is patched with p1c", isPatchedWith(p1c))))
	if resp, err := testDelete1.Do(); err != nil {
		t.Logf("raw response:\n%s\n", resp.Raw())
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

	// test listing after deleting p1
	testList3 := service.List("/posts").
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 1))
	if resp, err := testList3.Do(); err != nil {
		t.Logf("raw response:\n%s\n", resp.Raw())
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

	// test deleting p2
	testDelete2 := service.Delete("/post/" + p2.ID).
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 1)).
		Expect(restit.Nth(0).Of("posts").Is(restit.DescribeJSON(
			"item #0 returned is equal to p2", equals(p2))))
	if resp, err := testDelete2.Do(); err != nil {
		t.Logf("raw response:\n%s\n", resp.Raw())
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

	// test listing after deleting p2 (should be empty)
	testList4 := service.List("/posts").
		Expect(restit.StatusCodeIs(http.StatusOK)).
		Expect(restit.LengthIs("posts", 0))
	if resp, err := testList4.Do(); err != nil {
		t.Logf("raw response:\n%s\n", resp.Raw())
		t.Log(err.(restit.ContextError).Log())
		t.Errorf(err.Error())
	}

}
