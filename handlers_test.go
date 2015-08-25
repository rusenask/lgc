package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-zoo/bone"
)

func setup(c Client) *bone.Mux {
	//mux router with added routes
	m := getRouter(HandlerHttpClient{c})

	return m
}

func TestStublistHandler(t *testing.T) {
	testData := `{"version":"1.2.3","data": [{"name": "scenario1"}]}`
	server, c := testTools(200, testData)
	m := setup(*c)

	defer server.Close()

	//Testing get without specifying scenario
	req, err := http.NewRequest("GET", "/stubo/api/get/stublist", nil)
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusBadRequest)
}
