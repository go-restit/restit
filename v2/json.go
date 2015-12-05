package restit

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"regexp"
)

// reJSONNumber is the regular expression to match
// any JSON number values
var reJSONNum = regexp.MustCompile(`^-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+\-]?\d+)?$`)

// IsJSONNum test a string and see if it match the
// JSON definition of number
func IsJSONNum(b []byte) bool {
	return reJSONNum.Match(b)
}

// JSONType represents the different type of JSON values
// (string, number, object, array, true, false, null)
// true and false are combined as bool for obvious reason
type JSONType int

func (t JSONType) String() string {
	switch t {
	case TypeUndefined:
		return "TypeUndefined"
	case TypeString:
		return "TypeString"
	case TypeNumber:
		return "TypeNumber"
	case TypeObject:
		return "TypeObject"
	case TypeArray:
		return "TypeArray"
	case TypeBool:
		return "TypeBool"
	case TypeNull:
		return "TypeNull"
	}
	return "TypeUnknown"
}

// These constant represents different JSON value types
// as specified in http://www.json.org/
// with some exception:
// 1. true and false are combined as bool for obvious reason; and
// 2. TypeUnknown for empty strings
const (
	TypeUnknown   JSONType = -1
	TypeUndefined JSONType = iota
	TypeString
	TypeNumber
	TypeObject
	TypeArray
	TypeBool
	TypeNull
)

// ReadJSON reads a reader and store the bytes into JSON struct
func ReadJSON(reader io.Reader) (v *JSON, err error) {
	b, err := ioutil.ReadAll(reader)
	v = &JSON{b}
	return
}

// JSON is a DOM like json reading object
type JSON struct {
	raw []byte
}

// Unmarshal parses the JSON-encoded data and stores
// the result in the value pointed to by v
func (j JSON) Unmarshal(v interface{}) error {
	return json.Unmarshal(j.raw, v)
}

// UnmarshalJSON implements json.Unmarshaler
func (j *JSON) UnmarshalJSON(b []byte) error {
	j.raw = b
	return nil
}

// Type returns the JSONType of the containing JSON value
func (j JSON) Type() JSONType {

	switch {
	case j.raw == nil:
		// for nil raw, return TypeUndefined
		return TypeUndefined
	case len(j.raw) == 0:
		// for empty JSON string, return TypeUnknown
		return TypeUnknown
	case j.raw[0] == '"':
		// simply examine the first character
		// to determine the value type
		return TypeString
	case j.raw[0] == '{':
		// simply examine the first character
		// to determine the value type
		return TypeObject
	case j.raw[0] == '[':
		// simply examine the first character
		// to determine the value type
		return TypeArray
	case string(j.raw) == "true":
		fallthrough
	case string(j.raw) == "false":
		return TypeBool
	case string(j.raw) == "null":
		return TypeNull
	case IsJSONNum(j.raw):
		return TypeNumber
	}

	// return TypeUnknown for all other cases
	return TypeUnknown
}

// Get gets object's inner value.
// Only works with Object value type
func (j *JSON) Get(key string) (inner *JSON) {
	if j.Type() != TypeObject {
		inner = &JSON{nil}
		return
	}

	vmap := map[string]JSON{}
	if err := j.Unmarshal(&vmap); err != nil {
		inner = &JSON{nil} // dump the error
	} else if val, ok := vmap[key]; !ok {
		inner = &JSON{nil}
	} else {
		inner = &val
	}
	return
}

// Len gets the length of the value
// Only works with Array and String value type
func (j *JSON) Len() int {
	switch j.Type() {
	case TypeString:
		return len(string(j.raw)) - 2 // subtact the 2 " marks
	case TypeArray:
		vslice := []*JSON{}
		j.Unmarshal(&vslice)
		return len(vslice)
	}
	// default return -1 (for type mismatch)
	return -1
}

// GetN gets array's inner value.
// Only works with Array value type.
// 0 for the first item.
func (j *JSON) GetN(n int) *JSON {
	if j.Type() != TypeArray {
		return nil
	}

	vslice := []JSON{}
	j.Unmarshal(&vslice)
	if n < len(vslice) {
		return &vslice[n]
	}
	return nil
}

// Raw returns raw []byte of the JSON string
func (j *JSON) Raw() []byte {
	return j.raw
}

// String unmarshal the JSON into string then return
func (j *JSON) String() (v string) {
	j.Unmarshal(&v)
	return
}

// Number unmarshal the JSON into float64 then return
func (j *JSON) Number() (v float64) {
	j.Unmarshal(&v)
	return
}

// Bool unmarshal the JSON into bool then return
func (j *JSON) Bool() (v bool) {
	j.Unmarshal(&v)
	return
}

// IsNull tells if the JSON value is null or not
func (j *JSON) IsNull() bool {
	return j.Type() == TypeNull
}
