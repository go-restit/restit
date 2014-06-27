package restit

import (
	"fmt"
	"github.com/jmcvetta/napping"
)

type dummy struct {
	Name string
}

type dummyResponse struct {
	Dummies []dummy
}

func (r *dummyResponse) Count() int {
	return len(r.Dummies)
}

func (r *dummyResponse) NthValid(n int) error {
	if r.Dummies[n].Name == "" {
		return fmt.Errorf("All dummies should have a name")
	}
	return nil
}

func (r *dummyResponse) GetNth(n int) (interface{}, error) {
	nth := r.Dummies[n]
	return &nth, nil
}

func (r *dummyResponse) Match(a interface{}, b interface{}) error {
	if a.(*dummy).Name != b.(*dummy).Name {
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
