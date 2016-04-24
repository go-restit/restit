package restit

import (
	"fmt"

	"github.com/go-restit/lzjson"
	"golang.org/x/net/context"
)

// NthTest tests a given
type NthTest struct {
	n     uint
	name  string
	tests []JSONTest
}

// Nth will get the nth of a specific field in a JSON
// and test it against some JSONTest
func Nth(n uint) *NthTest {
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
	root, err := lzjson.Decode(resp.Body())
	if err != nil {
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
