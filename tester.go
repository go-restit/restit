package restit

import (
	"fmt"
	"github.com/jmcvetta/napping"
)


// Tester represents an ordinary RESTful entry point
type Tester struct {
	BaseUrl string
}


// Test to create a record with the payload provided.
// Expects the response to carry a valid record the newly created record.
func (t *Tester) TestCreate(
	payload interface{}, result TestRespond) (resp *napping.Response, err error) {

	// make the REST create request
	resp, err = napping.Post(t.BaseUrl,
		&payload, &result, nil)

	// test: has to be exactly 1 result
	count := result.Count()
	if count != 1 {
		err = fmt.Errorf("Bad response in TestCreate. "+
			"There are %d results (expecting 1)",
			count)
		return
	}

	// test: test the result
	err = result.NthValid(0)
	if err != nil {
		return
	}
	err = result.NthMatches(0, &payload)
	if err != nil {
		return
	}

	return
}


// Test to retrieve 1 record. Then try to match the record with
// the NthMatches interface method
// Expects the response to carry 1 and only 1 valid record
func (t *Tester) TestRetrieveOne(
	id string, payload interface{}, result TestRespond) (resp *napping.Response, err error) {

	// REST retrieve record with id
	p := napping.Params{} // empty payload for retrieve
	resp, err = napping.Get(t.BaseUrl+"/"+id,
		&p, &result, nil)
	if err != nil {
		return
	}

	// test: has to be exactly 1 result
	count := result.Count()
	if count != 1 {
		err = fmt.Errorf("Bad response in TestRetrieveOne. "+
			"There are %d results (expecting 1)",
			count)
		return
	}

	// test: test the result
	err = result.NthValid(0)
	if err != nil {
		return
	}
	err = result.NthMatches(0, &payload)
	if err != nil {
		return
	}

	return
}


// Test to update 1 record.
// Expects the response to carry the newly updated record.
func (t *Tester) TestUpdate(
	id string, payload interface{}, result TestRespond) (resp *napping.Response, err error) {

	// REST update record (of given id) with the payload
	resp, err = napping.Put(t.BaseUrl+"/"+id,
		&payload, &result, nil)
	if err != nil {
		return
	}

	// test: has to be exactly 1 result
	count := result.Count()
	if count != 1 {
		err = fmt.Errorf("Bad response in TestUpdate. "+
			"There are %d results (expecting 1)",
			count)
		return
	}

	// test: test the result
	err = result.NthValid(0)
	if err != nil {
		return
	}
	err = result.NthMatches(0, &payload)
	if err != nil {
		return
	}

	return
}


// Test to delete 1 record.
// Expects the response to carry the just deleted record.
func (t *Tester) TestDelete(
	id string, payload interface{}, result TestRespond) (resp *napping.Response, err error) {

	// REST delete record of given id
	resp, err = napping.Delete(t.BaseUrl+"/"+id,
		&result, nil)
	if err != nil {
		return
	}

	// test: has to be exactly 1 result
	count := result.Count()
	if count != 1 {
		err = fmt.Errorf("Bad response in TestUpdate. "+
			"There are %d results (expecting 1)",
			count)
		return
	}

	// test: test the result
	err = result.NthValid(0)
	if err != nil {
		return
	}
	err = result.NthMatches(0, &payload)
	if err != nil {
		return
	}

	return
}
