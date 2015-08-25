package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-zoo/bone"
)

func setup(c Client) *bone.Mux {
	//mux router with added routes
	m := getRouter(HandlerHttpClient{c})

	return m
}

func TestStublistHandlerNoScenario(t *testing.T) {
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

func TestStublistHandler(t *testing.T) {
	testData := `{"version":"1.2.3","data": [{"name": "scenario1"}]}`
	server, c := testTools(200, testData)
	m := setup(*c)

	defer server.Close()

	//Testing get scenario stubs
	req, err := http.NewRequest("GET", "/stubo/api/get/stublist?scenario=some_name", nil)
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusOK)
}

func TestDeleteStubsHandler(t *testing.T) {
	testData := `deleted`
	server, c := testTools(200, testData)
	m := setup(*c)

	defer server.Close()

	//Testing get scenario stubs
	req, err := http.NewRequest("GET", "/stubo/api/delete/stubs?scenario=some_name", nil)
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)
	// reading resposne body
	body, err := ioutil.ReadAll(respRec.Body)

	expect(t, strings.Contains(string(body), "deleted"), true)
	expect(t, respRec.Code, http.StatusOK)
}

func TestDeleteStubsHandlerFail(t *testing.T) {
	testData := `deleted`
	server, c := testTools(200, testData)
	m := setup(*c)

	defer server.Close()

	//Testing get scenario stubs
	req, err := http.NewRequest("GET", "/stubo/api/delete/stubs", nil)
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusBadRequest)
}
