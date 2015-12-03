package restit_test

import (
	"encoding/json"
	"strings"
	"testing"

	restit "github.com/yookoala/restit/v2"
)

func dummyStr() string {
	return `{
    "number": 1234.56,
    "string": "foo bar",
    "arrayOfString": [
      "one",
      "two",
      "three",
      "four"
    ],
    "object": {
      "foo": "bar",
      "hello": "world",
      "answer": 42
    },
    "true": true,
    "false": false,
    "null": null
  }`
}

func TestJSON_Unmarshaler(t *testing.T) {
	str := dummyStr()
	j := &restit.JSON{}
	var umlr json.Unmarshaler = j
	if err := json.Unmarshal([]byte(str), umlr); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	}
	if want, have := str, string(j.Raw()); want != have {
		t.Errorf("\nexpected: %s\ngot: %s", want, have)
	}
}

func TestJSON_Unmarshal(t *testing.T) {
	str := dummyStr()
	j, err := restit.ReadJSON(strings.NewReader(str))
	if err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	}

	type type1 struct {
		Number        float64                `json:"number"`
		String        string                 `json:"string"`
		ArrayOfString []string               `json:"arrayOfString"`
		Object        map[string]interface{} `json:"object"`
	}
	v1 := type1{}

	if err := j.Unmarshal(&v1); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
		return
	}

	if want, have := 1234.56, v1.Number; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "foo bar", v1.String; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := 4, len(v1.ArrayOfString); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
		return
	}
	if want, have := "one", v1.ArrayOfString[0]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "two", v1.ArrayOfString[1]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "three", v1.ArrayOfString[2]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "four", v1.ArrayOfString[3]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "bar", v1.Object["foo"]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "world", v1.Object["hello"]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := float64(42), v1.Object["answer"]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

}

func TestJSON_Type(t *testing.T) {

	readJSON := func(str string) *restit.JSON {
		j, err := restit.ReadJSON(strings.NewReader(str))
		if err != nil {
			t.Errorf("unexpected error: %#v", err.Error())
			return nil
		}
		return j
	}

	if want, have := restit.TypeEmpty, readJSON("").Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}
	if want, have := restit.TypeString, readJSON(`"string"`).Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}
	if want, have := restit.TypeNumber, readJSON("1234").Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}
	if want, have := restit.TypeNumber, readJSON("-1234.567").Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}
	if want, have := restit.TypeNumber, readJSON("-1234.567E+5").Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}

	if want, have := restit.TypeObject, readJSON(`{ "foo": "bar" }`).Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}
	if want, have := restit.TypeArray, readJSON(`[ "foo", "bar" ]`).Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}

	if want, have := restit.TypeBool, readJSON("true").Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}
	if want, have := restit.TypeBool, readJSON("false").Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}
	if want, have := restit.TypeNull, readJSON("null").Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}
}

func TestJSON_Get(t *testing.T) {
	str := dummyStr()
	j, err := restit.ReadJSON(strings.NewReader(str))
	if err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	}

	if _, err := j.Get("notExists"); err == nil {
		t.Error("failed to trigger error with non-exists key")
	} else if want, have := "key not found", err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if n, err := j.Get("number"); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	} else if n == nil {
		t.Error("unexpected nil value")
	} else if want, have := restit.TypeNumber, n.Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	} else if want, have := 1234.56, n.Number(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := "", n.String(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := false, n.Bool(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := false, n.IsNull(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if n, err := j.Get("string"); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	} else if n == nil {
		t.Error("unexpected nil value")
	} else if want, have := restit.TypeString, n.Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	} else if want, have := float64(0), n.Number(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := "foo bar", n.String(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := false, n.Bool(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := false, n.IsNull(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := 7, n.Len(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	var nilJ *restit.JSON
	if n, err := j.Get("arrayOfString"); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	} else if n == nil {
		t.Error("unexpected nil value")
	} else if want, have := restit.TypeArray, n.Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	} else if want, have := 4, n.Len(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := "one", n.GetN(0).String(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := "two", n.GetN(1).String(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := "three", n.GetN(2).String(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := "four", n.GetN(3).String(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := nilJ, n.GetN(4); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := false, n.Bool(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := false, n.IsNull(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if n, err := j.Get("object"); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	} else if n == nil {
		t.Error("unexpected nil value")
	} else if want, have := restit.TypeObject, n.Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	} else if p, err := n.Get("answer"); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	} else if want, have := restit.TypeNumber, p.Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	} else if want, have := false, n.Bool(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := false, n.IsNull(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if n, err := j.Get("true"); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	} else if n == nil {
		t.Error("unexpected nil value")
	} else if want, have := restit.TypeBool, n.Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	} else if want, have := true, n.Bool(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if n, err := j.Get("false"); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	} else if n == nil {
		t.Error("unexpected nil value")
	} else if want, have := restit.TypeBool, n.Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	} else if want, have := false, n.Bool(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := false, n.IsNull(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if n, err := j.Get("null"); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	} else if n == nil {
		t.Error("unexpected nil value")
	} else if want, have := restit.TypeNull, n.Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	} else if want, have := true, n.IsNull(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}
