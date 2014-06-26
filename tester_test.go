package restit

import (
	"fmt"
	"github.com/jmcvetta/napping"
	"testing"
)

type dummy struct {
	Name string
}

type dummyResponse struct {
	Dummies []dummy
}

func (d *dummyResponse) Count() int {
	return len(d.Dummies)
}

func (d *dummyResponse) NthValid(n int) error {
	return nil
}

func (d *dummyResponse) GetNth(n int) (interface{}, error) {
	nth := d.Dummies[n]
	return &nth, nil
}

func (d *dummyResponse) Match(a interface{}, b interface{}) error {
	ptr_a := a.(*dummy)
	ptr_b := b.(*dummy)
	if (*ptr_a).Name != (*ptr_b).Name {
		return fmt.Errorf("Mismatch")
	}
	return nil
}

// dummy session with dummy send sequence
type dummySession struct {
}

func (s dummySession) Send(req *napping.Request) (
	res *napping.Response, err error) {
	var resv napping.Response
	res = &resv
	ptrResult := (*req).Result.(*dummyResponse)
	(*ptrResult).Dummies = []dummy{
		dummy{
			Name: "Hello Dummy",
		},
	}
	return
}

func Test_Tester_Create(t *testing.T) {
	a := dummy{
		Name: "Hello Dummy",
	}
	r := dummyResponse{}
	test := Rest("Dummy", "http://foobar/dummy").
		Create(&a).
		WithResponseAs(&r).
		ExpectResultCount(1).
		ExpectResultsValid().
		ExpectResultNth(0, &a)
	test.Session = dummySession{}
	_, err := test.Run()
	if err != nil {
		t.Error(err)
	}
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
