package example1

import (
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

// PostServer creates an http.Handler that handles Post
func PostServer() server.Server {
	store := server.NewStore()
	factory := server.NewFactory(Post{})
	return server.New(store, factory)
}
