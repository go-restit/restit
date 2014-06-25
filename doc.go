// Copyright 2014 Koala Yeung. All rights reserved.
// Use of this source code is governed by a GPL v3
// license that can be found in the LICENSE.md file.

/*
	Package RESTit provides helps to those who
	want to write an integration test program
	for their RESTful APIs.

	To use, you first have to implement the
	`TestRespond` interface.

	Example:

		type ExmplResp struct {
			Status string  `json:"status"`
			Result []Stuff `json:"result"`
		}

		func (r *ExmplResp) Count() int {
			return len(r.Result)
		}

		func (r *ExmplResp) NthExists(n int) (err error) {
			if n < 0 || n > r.Count() {
				err = fmt.Errorf("Nth item (%d) not exist. Length = %d",
					n, len(r.Result))
			}
			return
		}

		func (r *ExmplResp) NthValid(n int) (err error) {

			// check if the item exists
			err = r.NthExists(n)
			if err != nil {
				return
			}

			// test: the id should not be 0
			if r.Result[n].StuffId == 0 {
				return fmt.Errorf("The thing has a StuffId = 0")
			}

			return
		}

		func (r *ExmplResp) NthMatches(n int, comp *interface{}) (err error) {

			// check if the item exists
			err = r.NthExists(n)
			if err != nil {
				return
			}

			// check if the item match the payload
			stuff := r.Result[n]
			cptr := (*comp).(*map[string]string)
			c := *cptr
			if stuff.Name != c["name"] {
				err = fmt.Errorf("Name is \"%s\" (expected \"%s\")",
				stuff.Name, c["name"])
				return
			}

			return
		}

	Then you can test your RESTful API like this:

			// create a tester for your stuff
			tester := restit.Tester{
				BaseUrl: "http://foobar:8080/api/stuffs",
			}


			// -------- Test Create --------
			// 1. create the stuff
			resp, err = tester.TestCreate(&stuffToCreate, &result)
			if err != nil {
				fmt.Printf("Raw: %s\n", resp.RawText())
				panic(err)
			}
			stuffId := result.Result[0].StuffId // id of the created stuff

			// 2. test the created stuff
			_, err = tester.TestRetrieveOne(fmt.Sprintf("%d", stuffId), &stuffToCreate, &result)
			if err != nil {
				fmt.Printf("Raw: %s\n", resp.RawText())
				panic(err)
			}


			// -------- Test Update --------
			// 1. update the stuff
			resp, err = tester.TestUpdate(fmt.Sprintf("%d", stuffId), &stuffToUpdate, &result)
			if err != nil {
				fmt.Printf("Raw: %s\n", resp.RawText())
				panic(err)
			}

			// 2. test the updated stuff
			_, err = tester.TestRetrieveOne(fmt.Sprintf("%d", stuffId), &stuffToUpdate, &result)
			if err != nil {
				fmt.Printf("Raw: %s\n", resp.RawText())
				panic(err)
			}


			// -------- Test Delete --------
			// delete the stuff
			_, err = tester.TestDelete(fmt.Sprintf("%d", stuffId), &stuffToUpdate, &result)
			if err != nil {
				fmt.Printf("Raw: %s\n", resp.RawText())
				panic(err)
			}

*/
package restit
