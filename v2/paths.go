package restit

import (
	"net/url"
	"path"
)

// NewNoun return a Noun declaration of certain RESTful service.
// If your RESTful API doesn't distinct between singular and plural
// noun on path, simply use the same on both parameters.
func NewNoun(singular, plural string) Noun {
	return noun{singular, plural}
}

// Noun provides simple description to noun of a RESTful service
type Noun interface {
	Singular() string
	Plural() string
}

// noun is the default implemetation of Noun
type noun struct {
	singular string
	plural   string
}

// Singular implements Noun
func (n noun) Singular() string {
	return n.singular
}

// Plural implements Noun
func (n noun) Plural() string {
	return n.plural
}

// NewPaths return a Paths declaration of certain RESTful service
func NewPaths(base string, noun Noun) (Paths, error) {
	baseURL, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	return paths{baseURL, noun}, nil
}

// Paths provides simple description to paths of a RESTful service
type Paths interface {
	// Singular returns the path string to the RESTful API
	// for single entity (e.g. "/api/article").
	// Can append additional parameters to the path.
	// (e.g. "/api/article/123")
	Singular(v ...string) string

	// Plural returns the path string to the RESTful API
	// for entire collection (e.g. "/api/articles").
	// Can append additional parameters to the path.
	// (e.g. "/api/articles/someUser")
	Plural(v ...string) string
}

// paths is the default implemetation of Paths
type paths struct {
	base *url.URL
	noun Noun
}

// Singular implements Paths
func (p paths) Singular(v ...string) string {
	u := *p.base
	u.Path = path.Join(append([]string{u.Path, p.noun.Singular()}, v...)...)
	return u.String()
}

// Plural implements Paths
func (p paths) Plural(v ...string) string {
	u := *p.base
	u.Path = path.Join(append([]string{u.Path, p.noun.Plural()}, v...)...)
	return u.String()
}
