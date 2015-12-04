package server_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/yookoala/restit/v2/example/server"
)

type post struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

// GetID implement server.Storable
func (p post) GetID() interface{} {
	return p.ID
}

// SetID implement server.Storable
func (p *post) SetID(id string) {
	(*p).ID = id
}

// GetType implement server.Storable
func (p post) GetType() string {
	return "post"
}

// Equal test if the post is equal to the given object
func (p post) Equal(v interface{}) (err error) {

	var prev post

	switch v.(type) {
	case post:
		prev = v.(post)
	case *post:
		prevPtr := v.(*post)
		prev = *prevPtr
	default:
		err = fmt.Errorf("invalid type")
		return
	}

	if want, have := prev.ID, p.ID; want != have {
		err = fmt.Errorf("id mismatch. expected %#v, got %#v", want, have)
		return
	}

	if want, have := prev.Title, p.Title; want != have {
		err = fmt.Errorf("title mismatch. expected %#v, got %#v", want, have)
		return
	}

	if want, have := prev.Body, p.Body; want != have {
		err = fmt.Errorf("body mismatch. expected %#v, got %#v", want, have)
		return
	}

	if want, have := prev.Created, p.Created; !want.Equal(have) {
		err = fmt.Errorf("created time mismatch. expected %#v, got %#v", want, have)
		return
	}

	if want, have := prev.Updated, p.Updated; !want.Equal(have) {
		err = fmt.Errorf("updated time mismatch. expected %#v, got %#v", want, have)
		return
	}

	return
}

func TestStore(t *testing.T) {

	p1 := post{
		ID:      "p1",
		Title:   "ABC",
		Body:    "Hello body",
		Created: time.Now(),
		Updated: time.Now(),
	}

	store := server.NewStore()

	// test append then retrieve
	store.Put(&p1)
	if now := store.Get(p1.GetType(), p1.ID); now == nil {
		t.Error("unexpected nil value")
	} else if err := p1.Equal(now); err != nil {
		t.Errorf("unexpected error %#v", err.Error())
	}

	// test overwritting the stored item
	p2 := p1
	p2.Title = "ABC 2"
	p2.Body = "Hello body 2"
	store.Put(&p2)
	if now := store.Get(p1.GetType(), p1.ID); now == nil {
		t.Error("unexpected nil value")
	} else if err := p1.Equal(now); err == nil {
		t.Errorf("did not trigger error")
	} else if err := p2.Equal(now); err != nil {
		t.Errorf("unexpected error %#v", err.Error())
	}

	// test listing
	if list := store.List(p1.GetType()); list == nil {
		t.Error("unexpected nil value")
	} else if want, have := 1, len(list); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := p2, list[0].(*post); want.Equal(have) != nil {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if err := want.Equal(have); err != nil {
		t.Errorf("err: %#v", err.Error())
	}

	// test appending more item
	p3 := p1
	p3.ID = "p3"
	store.Put(&p3)
	if list := store.List(p1.GetType()); list == nil {
		t.Error("unexpected nil value")
	} else if want, have := 2, len(list); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	// test deleteing the stored item
	store.Delete(p2.GetType(), p2.ID)
	if now := store.Get(p1.GetType(), p1.ID); now != nil {
		t.Errorf("delete failed, got %#v", now)
	} else if list := store.List(p1.GetType()); list == nil {
		t.Error("unexpected nil value")
	} else if want, have := 1, len(list); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

}
