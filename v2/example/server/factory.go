package server

import "reflect"

// NewPtr takes a variable and make a new one
// of the same type
func NewPtr(v interface{}) (ret interface{}) {
	if v == nil {
		return
	}
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// create a pointer with a zero-initialized value
	// of the given type
	newPtr := reflect.New(t)
	ret = newPtr.Interface()
	return
}

// NewSlicePtr takes a variable and make a new list
// of the given type
func NewSlicePtr(v interface{}) (ret interface{}) {
	if v == nil {
		return
	}

	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// make an zero-initialized slice value of the given type
	slice := reflect.MakeSlice(reflect.SliceOf(t), 0, 0)

	// make a slice pointer variable
	slicePtr := reflect.New(slice.Type())
	slicePtr.Elem().Set(slice)

	ret = slicePtr.Interface()
	return
}

// Factory creates pointer of type and of list
type Factory interface {

	// Make returns a pointer of certain type
	Make() interface{}

	// MakeSlice returns a pointer to list of cetain type
	MakeSlice() interface{}
}

// factory is an implentation of Factory
type factory struct {
	ptrFunc      func() interface{}
	slicePtrFunc func() interface{}
}

// Make implements Factory
func (f factory) Make() interface{} {
	return f.ptrFunc()
}

// MakeSlice implements Factory
func (f factory) MakeSlice() interface{} {
	return f.slicePtrFunc()
}

// NewFactory returns Factory of the provided type
func NewFactory(v interface{}) Factory {
	return &factory{
		ptrFunc: func() interface{} {
			return NewPtr(v)
		},
		slicePtrFunc: func() interface{} {
			return NewSlicePtr(v)
		},
	}
}
