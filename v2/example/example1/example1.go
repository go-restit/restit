package example1

import (
	"fmt"
	"time"

	"github.com/go-restit/restit/v2/example/server"
)

// Post is an implementation of Storable
type Post struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

// GetID implement server.Storable
func (p Post) GetID() interface{} {
	return p.ID
}

// SetID implement server.Storable
func (p *Post) SetID(id string) {
	(*p).ID = id
}

// GetType implement server.Storable
func (p Post) GetType() string {
	return "post"
}

// PatchWith implements server.Patchable.
func (p *Post) PatchWith(v interface{}) error {
	var patch Post

	switch v.(type) {
	case Post:
		patch = v.(Post)
	case *Post:
		ptr := v.(*Post)
		patch = *ptr
	default:
		return fmt.Errorf("invalid type")
	}

	if patch.ID != "" {
		p.ID = patch.ID
	}
	if patch.Title != "" {
		p.Title = patch.Title
	}
	if patch.Body != "" {
		p.Body = patch.Body
	}
	if !patch.Created.IsZero() {
		p.Created = patch.Created
	}
	if !patch.Updated.IsZero() {
		p.Updated = patch.Updated
	} else {
		p.Updated = time.Now()
	}

	return nil
}

// PostServer creates an http.Handler that handles Post
func PostServer() server.Server {
	store := server.NewStore()
	factory := server.NewFactory(Post{})
	return server.New(store, factory)
}
