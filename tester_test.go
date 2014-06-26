package restit

import (
	"github.com/jmcvetta/napping"
	"testing"
)

type dummy struct {
}

type dummySession struct {
}

func (s *dummySession) Send(req *napping.Request) (
	res *napping.Response, err error) {
	return
}

func Test_Tester_Create(t *testing.T) {
	a := dummy{}
	Rest("Dummy", "http://foobar/dummy").
		Create(&a)
}

func Test_Tester_Retrieve(t *testing.T) {
	Rest("Dummy", "http://foobar/dummy").
		Retrieve("some_id")
}

func Test_Tester_Update(t *testing.T) {
	a := dummy{}
	Rest("Dummy", "http://foobar/dummy").
		Update("some_id", &a)
}

func Test_Tester_Delete(t *testing.T) {
	Rest("Dummy", "http://foobar/dummy").
		Delete("some_id")
}
