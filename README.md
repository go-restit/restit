RESTit
======

A Go micro-framework to help writing RESTful API integration test

Package RESTit provides helps to those who want to write an 
integration test program for their JSON-based RESTful APIs.

The aim is to make these integration readable highly re-usable,
and yet easy to modify.


Dependencies
------------

RESTit is written in Go (a.k.a. Golang). You'll have to install go
first.

You need to first install "github.com/jmcvetta/napping":

    $ go get github.com/jmcvetta/napping


Install
-------

Just like installing other golang packages:

    $ go get github.com/yookoala/restit


How to Use
----------

To use, you first have to implement the `TestRespond`
interface for the REST server response that you're 
expecting.

In Go, json contents are usually unmarshaled to structs.
What you need to do is implement 4 methods to the struct
type.

Example:

    type ExmplResp struct {
    	Status string  `json:"status"`
    	Result []Stuff `json:"result"`
    }

    func (r *ExmplResp) Count() int {
    	return len(r.Result)
    }

    func (r *ExmplResp) NthValid(n int) (err error) {
    	// some test to see nth record is valid
    	// such as:
    	if r.Result[n].StuffId == 0 {
    		return fmt.Errorf("The thing has a StuffId = 0")
    	}
    	return
    }

    func (r *ExmplResp) GetNth(n int) (item interface{}, err error) {
    	// return the nth item
    	return Result[n]
    }

    func (r *ExmplResp) Match(a interface{}, b interface{}) (match bool, err error) {

    	// cast a and b back to Stuff
    	real_a = a.(Stuff)
    	real_b = b.(Stuff)

    	// some test to check if real_a equals real_b
    	// ...
    }

Then you may write the tests.


Write Your Tests
----------------

You can then test your RESTful API like this:

    import "github.com/yookoala/restit"

    // create a tester for your stuffAPI
    // first parameter is a human readable name that will
    // appear on error second parameter is the base URL to
    // the API entry point
    stuffAPI := restit.Rest(
        "Stuff", "http://foobar:8080/api/stuffs")

    // some parameters we'll use
    var result restit.TestResult
    var test restit.TestCase
    var response ExmplResp

    // some random stuff for test
    // for example,
    stuffToCreate = Stuff{
        Name: "apple",
        Color: "red",
    }
    stuffToUpdate = Stuff{
        Name: "orange",
        Color: "orange",
    }

    // here we add some dummy security measures
    // or you may add any parameters you like
    securityInfo := napping.Params{
        "username": "valid_user",
        "token": "some_security_token",
    }


    // -------- Test Create --------
    // 1. create the stuff
    test = stuffAPI.
        Create(&stuffToCreate).
        WithParams(&securityInfo).
        WithResponseAs(&response).
        ExpectResultCount(1).
        ExpectResultsValid().
        ExpectResultNth(0, &stuffToCreate).
        ExpectResultsToPass(
            "Custom Test",
            func (r Response) error {
            // some custom test you may want to run
            // ...
        })

    result, err := test.Run()
    if err != nil {
        // you may add more verbose output for
        // inspection
        fmt.Printf("Failed creating stuff!!!!\n")
        fmt.Printf("Please inspect the Raw Response: " 
                + result.RawText())

        // or you can simply:
        panic(err)
    }

    // id of the just created stuff
    // for reference of later tests
    stuffId := response.Result[0].StuffId
    stuffToCreate.StuffId = stuffId
    stuffToUpdate.StuffId = stuffId

    // 2. test the created stuff
    test = stuffAPI.
        Retrieve(fmt.Sprintf("%d", stuffId)).
        WithResponseAs(&response).
        ExpectResultCount(1).
        ExpectResultsValid().
        ExpectResultNth(0, &stuffToCreate)
    // A short hand to just panic on any error
    result = test.RunOrPanic()


    // -------- Test Update --------
    // 1. update the stuff
    result = stuffAPI.
        Update(&stuffToUpdate,
            fmt.Sprintf("%d", stuffId)).
        WithResponseAs(&response).
        WithParams(&securityInfo).
        ExpectResultCount(1).
        ExpectResultsValid().
        ExpectResultNth(0, &stuffToUpdate).
        RunOrPanic() // Yes, you can be this lazy

    // 2. test the updated stuff
    result = stuffAPI.
        Retrieve(fmt.Sprintf("%d", stuffId)).
        WithResponseAs(&response).
        ExpectResultCount(1).
        ExpectResultsValid().
        ExpectResultNth(0, &stuffToUpdate).
        RunOrPanic()


    // -------- Test Delete --------
    // delete the stuff
    result = stuffAPI.
        Delete(fmt.Sprintf("%d", stuffId)).
        WithResponseAs(&response).
        WithParams(security).
        ExpectResultCount(1).
        ExpectResultsValid().
        ExpectResultNth(0, &stuffToUpdate).
        RunOrPanic()

    // 2. test the deleted stuff
    result = stuffAPI.
        Retrieve(fmt.Sprintf("%d", stuffId)).
        WithResponseAs(&response).
        ExpectResultStatus(404).
        RunOrPanic()


Bug Reports
-----------

To report issue, please visit the
[issue tracker](https://github.com/yookoala/restit/issues).

And of course, patches and pull requests are most welcome.
