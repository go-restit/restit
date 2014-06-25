package restit

type TestRespond interface {
	Count() int
	NthExists(int) error
	NthValid(int) error
	NthMatches(int, *interface{}) error
}
