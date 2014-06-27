package restit

import (
	"testing"
)

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
