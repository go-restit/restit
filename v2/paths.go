package restit

import (
	"net/url"
	"path"
)

// NewPaths return a Paths declaration of certain RESTful service
func NewPaths(base, singular, plural string) (Paths, error) {
	baseURL, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	return paths{baseURL, singular, plural}, nil
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
	base     *url.URL
	singular string
	plural   string
}

// Singular implements Paths
func (p paths) Singular(v ...string) string {
	u := *p.base
	u.Path = path.Join(append([]string{u.Path, p.singular}, v...)...)
	return u.String()
}

// Plural implements Paths
func (p paths) Plural(v ...string) string {
	u := *p.base
	u.Path = path.Join(append([]string{u.Path, p.plural}, v...)...)
	return u.String()
}
