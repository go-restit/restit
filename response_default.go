package restit

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrorNthNotFound = errors.New("Nth item not found")
var ErrorNoValidator = errors.New("Validator not found. Please set validator using SetValidator")
var ErrorNoMatcher = errors.New("Matcher not found. Please set matcher using SetMatcher")
var ErrorInvalidValidator = errors.New("Validator not valid")
var ErrorInvalidMatcher = errors.New("Matcher not valid")
var ErrorInvalidType = errors.New("Response item type is not struct")
var ErrorUnexpectedItem = errors.New("At least one decoded json item is not map[string]interface{}")

type Validator func(interface{}) error
type Matcher func(interface{}, interface{}) error

// DefaultResponse is the default response type to use
type DefaultResponse map[string]interface{}

// NewResponse create a new DefaultResponse
func NewResponse(n string, t interface{}) *DefaultResponse {
	return &DefaultResponse{
		"_list_name": n,
		"_type":      t,
	}
}

// ValidateWith set the validator function.
// If no validator function is given, the test will by default fail
func (p *DefaultResponse) SetValidator(in func(interface{}) error) {
	var f Validator = in
	(*p)["_validator"] = f
}

// MatchWith set the matcher function.
// If no matcher function is given, the test will by default fail
func (p *DefaultResponse) SetMatcher(in func(interface{}, interface{}) error) {
	var f Matcher = in
	(*p)["_matcher"] = f
}

// typeListName returns the set field name of types list
func (p *DefaultResponse) typeListName() (name string) {
	var n interface{}
	var ok bool

	if n, ok = (*p)["_list_name"]; !ok {
		panic("Unable to find internal field \"_list_name\"")
	}

	if name, ok = n.(string); !ok {
		panic(fmt.Sprintf(
			"Unable to cast internal field \"%s\" to string",
			name))
	}

	return name
}

// typeType returns the defined item type
func (p *DefaultResponse) typeType() reflect.Type {
	return reflect.TypeOf((*p)["_type"])
}

// typeKind returns the reflect.Kind of defined item type
func (p *DefaultResponse) typeKind() reflect.Kind {
	return p.typeType().Kind()
}

// typeKind returns a slice of the defined item type
func (p *DefaultResponse) typeSlice() reflect.Type {
	return reflect.SliceOf(p.typeType())
}

// typePtr returns a pointer to a new variable of type
func (p *DefaultResponse) typePtr() interface{} {
	// create new type pointer
	tPtr := reflect.New(reflect.PtrTo(p.typeType()))

	// initialize with empty type value
	t := reflect.New(p.typeType())
	tPtr.Elem().Set(t)

	// return raw interface of pointer
	return tPtr.Elem().Interface()
}

func (p *DefaultResponse) getItems() (items []interface{}) {
	var ok bool
	var masked interface{}

	if masked, ok = (*p)[p.typeListName()]; !ok {
		panic(fmt.Sprintf("Field \"%s\" does not exists.",
			p.typeListName()))
	}
	if items, ok = masked.([]interface{}); !ok {
		panic(fmt.Sprintf("Unable to cast the field \"%s\" (%#v) into []interface{}",
			p.typeListName(), masked))
	}
	return
}

// Count counts the number of items in the result
func (p *DefaultResponse) Count() int {
	return len(p.getItems())
}

// NthValid test if the nth item is valid
func (p *DefaultResponse) NthValid(n int) (err error) {

	var ok bool
	var in1, in2 interface{}
	var f Validator

	if in1, ok = (*p)["_validator"]; !ok {
		return ErrorNoValidator
	}

	if f, ok = in1.(Validator); !ok {
		return ErrorInvalidValidator
	}

	if in2, err = p.GetNth(n); err != nil {
		return
	}

	return f(in2)
}

// GetNth get the Nth item in the DefaultResponse
func (p *DefaultResponse) GetNth(n int) (ret interface{}, err error) {

	var ok bool
	var nth map[string]interface{}

	// get the item list
	s := p.getItems()
	if len(s) <= n {
		err = ErrorNthNotFound
		return
	}

	// test casting the item
	if nth, ok = s[n].(map[string]interface{}); !ok {
		err = ErrorUnexpectedItem
		return
	}

	// set the value of item to the type
	t := p.typeType()
	tPtr := p.typePtr()

	v := reflect.ValueOf(tPtr).Elem()
	if v.Kind() != reflect.Struct {
		err = ErrorInvalidType
		return
	}

	// loop through all the type fields
	for i := 0; i < t.NumField(); i++ {

		var fvalue interface{}
		tf := t.Field(i)

		// get field name / tagged name for json
		fname := tf.Name
		if name := tf.Tag.Get("json"); name != "" {
			fname = name
		}

		// skip field if name is "-"
		if fname == "-" {
			continue
		}

		// check if the fname exists in the map
		if fvalue, ok = nth[fname]; !ok {
			continue // just skip
		}

		// set the field value as fvalue
		err = setField(v.Field(i), fvalue)
	}

	// if there is error, return nothing
	if err == nil {
		ret = reflect.ValueOf(tPtr).Elem().Interface()
	}
	return
}

// setField set a given field value of a typePtr (settable)
// to any value specified in v
func setField(f reflect.Value, v interface{}) (err error) {

	// handle different `v` value type
	vkind := reflect.ValueOf(v).Kind()
	if vkind == reflect.Int64 {
		vv, _ := v.(int64)
		switch f.Kind() {
		case reflect.Int:
			f.SetInt(vv)
			break
		default:
			err = fmt.Errorf("Type mismatch. The field cannot accept integer value.")
		}
		return
	} else if vkind == reflect.Float64 {
		vv, _ := v.(float64)

		switch f.Kind() {
		case reflect.Int, reflect.Int32, reflect.Int64:
			vvv := int64(vv)
			f.SetInt(vvv)
			break
		case reflect.Float32, reflect.Float64:
			f.SetFloat(vv)
			break
		default:
			err = fmt.Errorf("Type mismatch. The field cannot accept float value.")
		}
		return
	} else if vkind == reflect.String {
		vv, _ := v.(string)

		switch f.Kind() {
		case reflect.String:
			f.SetString(vv)
			break
		default:
			err = fmt.Errorf("Type mismatch. The field cannot not accept string value.")
			break
		}
		return
	}
	return
}

// Match test if the nth item is valid
func (p *DefaultResponse) Match(a interface{}, b interface{}) (err error) {

	var ok bool
	var in interface{}
	var f Matcher

	if in, ok = (*p)["_matcher"]; !ok {
		return ErrorNoMatcher
	}

	if f, ok = in.(Matcher); !ok {
		return ErrorInvalidMatcher
	}

	return f(a, b)
}

// Reset the value of this response
func (p *DefaultResponse) Reset() {
	v := DefaultResponse{
		"_list_name": (*p)["_list_name"],
		"_type":      (*p)["_type"],
		"_validator": (*p)["_validator"],
		"_matcher":   (*p)["_matcher"],
	}
	*p = v
}
