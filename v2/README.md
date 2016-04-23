# RESTit (v2) [![Godoc][godoc-badge]][godoc] [![Travis test][travis-badge]][travis] [![Appveyor test][appveyor-badge]][appveyor]

A Go micro-framework to help writing RESTful API integration test

Package RESTit provides helps to those who want to write an
integration test program for their JSON-based RESTful APIs.

The aim is to make these integration readable highly re-usable,
and yet easy to modify.

[godoc]: https://godoc.org/github.com/go-restit/restit/v2
[godoc-badge]: https://godoc.org/github.com/go-restit/restit/v2?status.svg
[travis]: https://travis-ci.org/go-restit/restit?branch=master
[travis-badge]: https://api.travis-ci.org/go-restit/restit.svg?branch=master
[appveyor]: https://ci.appveyor.com/project/yookoala/restit?branch=master
[appveyor-badge]: https://ci.appveyor.com/api/projects/status/github/go-restit/restit?branch=master&svg=true


## Design Principles

Less is more. The main theme of RESTit v2 is:

- To reduce code for testing.
- Make tests reusable.

This version has several improvements over v1:

1. No longer need to define protocol before testing.
2. Generates [*http.Request][http.Request] directly.
   No longer depend on the [napping library][napping]. One less thing to learn.   
3. Support [httptest.ResponseRecorder][httptest.ResponseRecorder] test and
   ordinary TCP tests. Tests are nearly identical for both. More flexibility.
4. New [lzjson][lzjson] JSON decoding library allow you to examine a JSON value
   without predefined static structure. It also allow you to partially decode
   a JSON value (e.g. "posts" field in it).

[http.Request]: https://golang.org/pkg/net/http/#Request
[napping]: https://github.com/jmcvetta/napping
[lzjson]: https://github.com/go-restit/lzjson


## Usage

### Basic of Tests

RESTit (v2) supports [httptest.ResponseRecorder][httptest.ResponseRecorder]
based testing. In which approach, you don't need to run a server that listens
to TCP port. You only need an [http.Handler][http.Handler] to start testing.

[httptest.ResponseRecorder]: https://golang.org/pkg/net/http/httptest/#ResponseRecorder
[http.Handler]: https://golang.org/pkg/net/http/#Handler

```go
package post_test

import (
  "post"

  "testing"
  "net/http/httptest"

  restit "github.com/go-restit/restit/v2"
)

func TestPostAPI(t *testing.T)

  var err error
  var resp restit.Response

  // define the path for your handler
  service := restit.NewHTTPTestService("/dummy/api", post.Handler)
  token := "some_access_token"

  post1 := Post{ID: "1234", Name: "hello world", Body: "some hello world message"}
  resp, err = service.Create(post1, "/posts").
    AddQuery("access_token", token).
    AddHeader("User-Agent", "RESTit tester").
    Expect(expectation1).
    Expect(expectation2).
    Expect(expectation3).
    Do()
  if err != nil {
    t.Error(err.Error())
  }

  resp, err = service.Retrieve("/post/1234").
    AddHeader("User-Agent", "RESTit tester").
    Expect(expectation1).
    Expect(expectation2).
    Expect(expectation3).
    Do()
  if err != nil {
    t.Error(err.Error())
  }

  resp, err = service.List("/posts").
    AddHeader("User-Agent", "RESTit tester").
    Expect(expectation1).
    Expect(expectation2).
    Expect(expectation3).
    Do()
  if err != nil {
    t.Error(err.Error())
  }

  post2 := Post{ID: "1234", Name: "updated", Body: "some updated message"}
  resp, err = service.Update(post2, "/post/1234").
    AddHeader("User-Agent", "RESTit tester").
    AddQuery("access_token", token).
    Expect(expectation1).
    Expect(expectation2).
    Expect(expectation3).
    Do()
  if err != nil {
    t.Error(err.Error())
  }

  resp, err = service.Delete("/post/1234").
    AddHeader("User-Agent", "RESTit tester").
    AddQuery("access_token", token).
    Expect(expectation1).
    Expect(expectation2).
    Expect(expectation3).
    Do()
  if err != nil {
    t.Error(err.Error())
  }

}

```

We use raw [*http.Request][http.Request] in our test case. Feel free to
manipulate it before doing the test:

```go

caseCreate := service.Create(Post{Name: "hello world", Body: "some hello world message"}).
  Expect(expectation1).
  Expect(expectation2).
  Expect(expectation3)

caseCreate.Request.Header.Add("X-Custom-Header", "Hello World")
caseCreate.Do()

```

The request and response are fully examinable:

```go
package post_test

import (
  "post"

  "testing"
  "net/http/httptest"

  restit "github.com/go-restit/restit/v2"
)

func TestPostAPI(t *testing.T)

  // define the path for your handler
  service := restit.NewHTTPTestService("/dummy/api/posts", post.Handler)

  caseCreate := service.Create(Post{Name: "hello world", Body: "some hello world message"}).
    Expect(expectation1).
    Expect(expectation2).
    Expect(expectation3)

  resp, err := caseCreate.Do()
  if err != nil {
    t.Errorf("error running create %s", err)
    return
  }

  bodyBytes, err := ioutil.ReadAll(resp.Body())
  if err != nil {
    t.Errorf("error reading response body: %s", err)
    return
  }

  t.Logf("request: %#v", caseCreate.Request) // raw *htt.Request used
  t.Logf("body: %s", bodyBytes)
}

```


### Expectation

// To be written


### JSON Decoding with Ease

// To be written


## Bug Reports

To report issue, please visit the
[issue tracker](https://github.com/go-restit/restit/issues).

And of course, patches and pull requests are most welcome.
