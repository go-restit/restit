package restit

import (
	"fmt"

	"github.com/go-restit/lzjson"
	"golang.org/x/net/context"
)

// StatusCodeIs test the response status code
func StatusCodeIs(n int) Expectation {
	return Describe(
		fmt.Sprintf("status code is %d", n),
		func(ctx context.Context, resp Response) (err error) {
			if want, have := n, resp.StatusCode(); want != have {
				ctxErr := NewContextError("expected %d, got %d", want, have)
				ctxErr.Prepend("ref", "header status code")
				err = ctxErr
			}
			return
		})
}

// LengthIs test the length of a given field
func LengthIs(name string, n int) Expectation {
	return Describe(
		fmt.Sprintf("length of %#v is %d", name, n),
		func(ctx context.Context, resp Response) (err error) {
			proto, err := resp.JSON()
			if err != nil {
				return
			} else if want, have := lzjson.TypeObject, proto.Type(); want != have {
				ctxErr := NewContextError("expected root to be type %s, got %s", want, have)
				ctxErr.Prepend("ref", "response")
				err = ctxErr
				return
			}

			list := proto.Get(name)
			if want, have := lzjson.TypeArray, list.Type(); want != have {
				ctxErr := NewContextError("expected %#v to be type %s, got %s",
					name, want, have)
				ctxErr.Prepend("ref", "response."+name)
				ctxErr.Append("response", string(proto.Raw()))
				err = ctxErr
				return
			}

			if want, have := n, list.Len(); want != have {
				ctxErr := NewContextError("expected %#v to be length %#v, got %#v",
					name, want, have)
				ctxErr.Prepend("ref", "response."+name+".length")
				err = ctxErr
				return
			}

			return
		})
}

// NthTest tests a given
type NthTest struct {
	n     int
	name  string
	tests []JSONTest
}

// Nth will get the nth of a specific field in a JSON
// and test it against some JSONTest
func Nth(n int) *NthTest {
	return &NthTest{n: n, tests: make([]JSONTest, 0)}
}

// Of specify the field which should be an array
func (t *NthTest) Of(name string) *NthTest {
	t.name = name
	return t
}

// Is specify the JSONTest which the JSON node should pass
func (t *NthTest) Is(test JSONTest) *NthTest {
	t.tests = append(t.tests, test)
	return t
}

// Desc implements Expectation
func (t *NthTest) Desc() string {
	desc := ""
	return desc
}

// Do implements Expectation
func (t *NthTest) Do(ctx context.Context, resp Response) (err error) {
	root := lzjson.Decode(resp.Body())
	if root.ParseError() != nil {
		err = fmt.Errorf("error decoding body to JSON (%s)", err)
		return
	}

	field := root.Get(t.name)
	if want, have := lzjson.TypeArray, field.Type(); want != have {
		if have == lzjson.TypeUndefined {
			err = fmt.Errorf("field %#v undefined", t.name)
		} else {
			err = fmt.Errorf("field %#v is not an array, is %s (%s)",
				t.name, field.Type(), field.Raw())
		}
		return
	}

	if field.Len() <= int(t.n) {
		err = fmt.Errorf("%s does not have item %d", t.name, t.n)
		return
	}

	nth := field.GetN(int(t.n))
	for _, test := range t.tests {
		if err = test.Do(nth); err != nil {
			err = fmt.Errorf("failed \"%s\" (%s)", test.Desc(), err)
			return
		}
	}

	return
}

// JSONTest represent a test on a given lzjson.Node
type JSONTest interface {
	Desc() string
	Do(lzjson.Node) error
}

// DescribeJSON returns a JSONTest of description and do function
func DescribeJSON(desc string, do func(lzjson.Node) error) JSONTest {
	return &jsonTest{desc, do}
}

type jsonTest struct {
	desc string
	do   func(lzjson.Node) error
}

func (t *jsonTest) Desc() string {
	return t.desc
}

func (t *jsonTest) Do(node lzjson.Node) error {
	return t.do(node)
}
