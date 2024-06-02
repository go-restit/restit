package restit

import (
	"fmt"
	"sort"
	"strings"
)

// ContextError defines an interface of context error
type ContextError interface {

	// Append adds a key-value pair at the end of log
	Append(key string, value interface{})

	// Prepend adds a key-value pair at the beginning of log
	Prepend(key string, value interface{})

	// Get retrieves a value of given key
	Get(key string) (value interface{})

	// Delete removes key-value pair of specific key
	Delete(key string)

	// Log returns a formated string with all key-value pair
	Log() string

	// Error returns the value of the key "message"
	Error() string
}

// NewContextError creates a default implementation of
// ContextError
func NewContextError(msg string, v ...interface{}) ContextError {
	err := &contextError{}
	err.Append("message", fmt.Sprintf(msg, v...))
	return err
}

// key value pairs in the context
type keyval struct {
	key    string
	val    interface{}
	weight int
}

// ExpandError expands errors to ContextError
func ExpandError(err error) ContextError {
	switch ctxErr := err.(type) {
	case *contextError:
		return ctxErr
	default:
		return NewContextError(err.Error())
	}
}

// contextError is the default implementation of ContextError
type contextError []keyval

// Append implements ContextError
func (ctx *contextError) Append(key string, val interface{}) {
	ctx.Delete(key)
	*ctx = append(*ctx, keyval{key, val, 0})
}

// Prepend implements ContextError
func (ctx *contextError) Prepend(key string, val interface{}) {
	ctx.Delete(key)
	*ctx = append(contextError{keyval{key, val, 0}}, *ctx...)
}

// Get implements ContextError
func (ctx *contextError) Get(key string) (val interface{}) {
	for _, kv := range *ctx {
		if kv.key == key {
			return kv.val
		}
	}
	return
}

// Delete implements ContextError
func (ctx *contextError) Delete(key string) {
	var ctx2 []keyval
	for _, kv := range *ctx {
		if kv.key != key {
			ctx2 = append(ctx2, kv)
		}
	}
	*ctx = ctx2
}

// Log implements ContextError
func (ctx *contextError) Log() string {
	sort.Sort(ctx)
	var msg []string
	for _, kv := range *ctx {
		msg = append(msg, fmt.Sprintf("%s=%#v", kv.key, kv.val))
	}
	return strings.Join(msg, " ")
}

// Error implements ContextError
func (ctx *contextError) Error() string {
	if msg := ctx.Get("message"); msg == nil {
	} else if str, ok := msg.(string); !ok {
	} else {
		return str
	}
	return "error" // dumb generic error message
}

// Len implements sort.interface
func (ctx *contextError) Len() int {
	if ctx == nil {
		return 0
	}
	return len(*ctx)
}

// Less implements sort.Interface
func (ctx *contextError) Less(i, j int) bool {
	return (*ctx)[i].weight < (*ctx)[j].weight
}

// Swap implements sort.Interface
func (ctx *contextError) Swap(i, j int) {
	(*ctx)[i], (*ctx)[j] = (*ctx)[j], (*ctx)[i]
}
