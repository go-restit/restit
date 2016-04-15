package restit_test

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	restit "github.com/go-restit/restit/v2"
)

func TestRequest(t *testing.T) {

	// Post is the type for test
	type Post struct {
		ID      string
		Seq     int
		OwnerID string
		Created time.Time
		Updated time.Time
	}

	// repeater implements of http.HandlerFunc interace
	repeater := func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("nil request"))
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(b)
	}

	p1 := Post{
		Seq:     rand.Int(),
		OwnerID: RandString(10),
	}

	r, err := restit.NewRequest("GET", "/foo/bar", p1)
	if err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	}

	w := httptest.NewRecorder()
	repeater(w, r)

	if want, have := http.StatusOK, w.Code; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if b, err := ioutil.ReadAll(w.Body); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	} else if len(b) == 0 {
		t.Error("empty w.Body")
	} else if j, err := json.Marshal(p1); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	} else if len(j) == 0 {
		t.Error("empty json string j")
	} else if want, have := string(j), string(b); want != have {
		t.Errorf("\nexpected: %s\ngot:     %s", want, have)
	}
}
