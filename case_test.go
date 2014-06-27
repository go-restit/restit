package restit

import (
	"fmt"
	"github.com/jmcvetta/napping"
	"testing"
)

func Test_Case_WithParams(t *testing.T) {

	r := napping.Request{}
	c := Case{
		Request: &r,
	}
	p := napping.Params{}
	c.WithParams(&p)
	if c.Request.Params != &p {
		t.Error("WithParams failed to set the parameter")
	}

}

func Test_Case_ExpectResultCount_0(t *testing.T) {

	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{},
	}
	c.ExpectResultCount(0)
	err := c.Expectations[0].Test(&r)
	if err != nil {
		t.Error("Unable to pass with a valid count 0")
	}

}

func Test_Case_ExpectResultCount_n(t *testing.T) {

	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{
			dummy{},
			dummy{},
			dummy{},
		},
	}
	c.ExpectResultCount(3)
	err := c.Expectations[0].Test(&r)
	if err != nil {
		t.Error("Unable to pass with a valid count n")
	}

}

func Test_Case_ExpectResultCount_err(t *testing.T) {

	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{
			dummy{},
			dummy{},
		},
	}
	c.ExpectResultCount(3)
	err := c.Expectations[0].Test(&r)
	if err == nil {
		t.Error("Unable to detect count mismatch")
	}

}

func Test_Case_ExpectResultsValid(t *testing.T) {

	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{
			dummy{Name: "Hello"},
		},
	}
	c.ExpectResultsValid()
	err := c.Expectations[0].Test(&r)
	if err != nil {
		t.Error("Unable to pass valid item")
	}

}

func Test_Case_ExpectResultsValid_err(t *testing.T) {

	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{
			dummy{},
		},
	}
	c.ExpectResultsValid()
	err := c.Expectations[0].Test(&r)
	if err == nil {
		t.Error("Unable to identify invalid item")
	}

}

func Test_Case_ExpectResultNth_match(t *testing.T) {

	d := dummy{Name: "Unique Dummy"}
	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{
			dummy{},
			d,
			dummy{},
		},
	}
	c.ExpectResultNth(1, &d)
	err := c.Expectations[0].Test(&r)
	if err != nil {
		t.Error("Expect Nth result to match given payload")
	}

}

func Test_Case_ExpectResultNth_err(t *testing.T) {

	d := dummy{Name: "Unique Dummy"}
	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{
			dummy{},
			d,
			dummy{},
		},
	}
	c.ExpectResultNth(0, &d)
	err := c.Expectations[0].Test(&r)
	if err == nil {
		t.Error("Expect Nth result to mis-match given payload")
	}

}

func Test_Case_ExpectResultsToPass_pass(t *testing.T) {

	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{
			dummy{},
		},
	}
	c.ExpectResultsToPass("Custom Test to pass",
		func(Response) error {
			return nil
		})
	err := c.Expectations[0].Test(&r)
	if err != nil {
		t.Error("Expect custom tests to pass")
	}

}

func Test_Case_ExpectResultsToPass_err(t *testing.T) {

	c := Case{}
	r := dummyResponse{
		Dummies: []dummy{
			dummy{},
		},
	}
	c.ExpectResultsToPass("Custom Test to fail",
		func(Response) error {
			return fmt.Errorf("Some error")
		})
	err := c.Expectations[0].Test(&r)
	if err == nil {
		t.Error("Expect custom tests to fail")
	}

}
