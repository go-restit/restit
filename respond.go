package restit

// Response needed to fulfill this interface
// in order to be tested by RESTit
type TestRespond interface {
	Count() int
	NthExists(int) error
	NthValid(int) error
	NthMatches(int, *interface{}) error
}
