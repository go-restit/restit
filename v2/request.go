package restit

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// NewRequest generate a normal *http.Request with JSON-encoded payload
// as request body
func NewRequest(method, urlString string, payload interface{}) (req *http.Request, err error) {
	jbytes, err := json.Marshal(payload)
	if err != nil {
		return
	}
	return http.NewRequest(method, urlString, bytes.NewReader(jbytes))
}
