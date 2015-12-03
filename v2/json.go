package restit

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

// JSONType represents the different type of JSON values
// (string, number, object, array, true, false, null)
// true and false are combined as bool for obvious reason
type JSONType int

func (t JSONType) String() string {
	switch t {
	case TypeEmpty:
		return "TypeEmpty"
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
// 2. an extra TypeEmpty for empty JSON strings
const (
	TypeUnknown JSONType = -1
	TypeEmpty   JSONType = iota
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

	// for empty JSON string, return TypeEmpty
	if len(j.raw) == 0 {
		return TypeEmpty
	}

	// simply examine the first character
	// to determine the value type
	switch j.raw[0] {
	case '"':
		return TypeString
	case '-':
		fallthrough
	case '0':
		fallthrough
	case '1':
		fallthrough
	case '2':
		fallthrough
	case '3':
		fallthrough
	case '4':
		fallthrough
	case '5':
		fallthrough
	case '6':
		fallthrough
	case '7':
		fallthrough
	case '8':
		fallthrough
	case '9':
		return TypeNumber
	case '{':
		return TypeObject
	case '[':
		return TypeArray
	}

	// try to match the whole string
	// if it is not too long
	if len(j.raw) <= 5 {
		switch string(j.raw) {
		case "true":
			fallthrough
		case "false":
			return TypeBool
		case "null":
			return TypeNull
		}
	}

	// return TypeUnknown for all other cases
	return TypeUnknown
}

// Get gets object's inner value.
// Only works with Object value type
func (j *JSON) Get(key string) (inner *JSON, err error) {
	if j.Type() != TypeObject {
		err = fmt.Errorf("only support JSON of TypeObject, not %s", j.Type())
		return
	}
	vmap := map[string]JSON{}
	err = j.Unmarshal(&vmap)

	if val, ok := vmap[key]; !ok {
		err = fmt.Errorf("key not found")
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
