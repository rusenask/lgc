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
	m := getRouter(HandlerHTTPClient{c})

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

	//Testing delete scenario stubs
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

	//Testing delete scenario stubs
	req, err := http.NewRequest("GET", "/stubo/api/delete/stubs", nil)
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusBadRequest)
}

func TestPutStubsHandlerFail(t *testing.T) {
	testData := `inserted`
	server, c := testTools(201, testData)
	m := setup(*c)

	defer server.Close()

	//Testing put stub
	req, err := http.NewRequest("POST", "/stubo/api/put/stub", nil)
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusBadRequest)
}

func TestPutStubsHandlerSuccess(t *testing.T) {
	testData := `inserted`
	server, c := testTools(201, testData)
	m := setup(*c)

	defer server.Close()

	//Testing put stub
	req, err := http.NewRequest("POST", "/stubo/api/put/stub?session=scenario:session",
		strings.NewReader("anything here, proxy doesn't unmarshall it anyway"))
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusOK)
}

func TestPutStubsHandlerNoScenario(t *testing.T) {
	testData := `inserted`
	server, c := testTools(201, testData)
	m := setup(*c)

	defer server.Close()

	//Testing put stub
	req, err := http.NewRequest("POST", "/stubo/api/put/stub?session=session",
		strings.NewReader("anything here, proxy doesn't unmarshall it anyway"))
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusBadRequest)
}

func TestPutStubsHandlerMultipleHeaders(t *testing.T) {
	testData := `inserted`
	server, c := testTools(201, testData)
	m := setup(*c)

	defer server.Close()

	//Testing put stub
	req, err := http.NewRequest("POST", "/stubo/api/put/stub?session=scenario:session&valued=2&value=4&ext_module=some_module",
		strings.NewReader("anything here, proxy doesn't unmarshall it anyway"))
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusOK)
}

func TestPutDelayPolicyHandler(t *testing.T) {
	testData := `inserted`
	server, c := testTools(201, testData)
	m := setup(*c)

	defer server.Close()

	//Testing get specific delay policy
	req, err := http.NewRequest("GET", "/stubo/api/put/delay_policy?name=slow&delay_type=fixed&milliseconds=1000", nil)
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusOK)
}

func TestPutDelayPolicyHandlerNoParams(t *testing.T) {
	testData := `inserted`
	server, c := testTools(201, testData)
	m := setup(*c)

	defer server.Close()

	//Testing get specific delay policy
	req, err := http.NewRequest("GET", "/stubo/api/put/delay_policy", nil)
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusOK)
}

func TestGetDelayPolicyHandler(t *testing.T) {
	testData := `delay`
	server, c := testTools(200, testData)
	m := setup(*c)

	defer server.Close()

	//Testing get specific delay policy
	req, err := http.NewRequest("GET", "/stubo/api/get/delay_policy?name=somename", nil)
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusOK)
}

func TestGetAllDelayPolicyHandler(t *testing.T) {
	testData := `delay`
	server, c := testTools(200, testData)
	m := setup(*c)

	defer server.Close()

	//Testing get all delay policies
	req, err := http.NewRequest("GET", "/stubo/api/get/delay_policy", nil)
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusOK)
}

func TestDeleteDelayPolicyHandler(t *testing.T) {
	testData := `delay`
	server, c := testTools(200, testData)
	m := setup(*c)

	defer server.Close()

	//Testing get all delay policies
	req, err := http.NewRequest("GET", "/stubo/api/delete/delay_policy?name=some_delay", nil)
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusOK)
}

func TestDeleteAllDelayPoliciesHandler(t *testing.T) {
	testData := `{"version": "0.6.6",
																 "data": [
																					{"delay_type":
																					 "fixed",
																			  	 "delayPolicyRef": "/stubo/api/v2/delay-policy/objects/my_delay",
																					 "name": "my_delay",
																					 "milliseconds": 50},
																					{"delay_type": "fixed",
																					"delayPolicyRef":
																					"/stubo/api/v2/delay-policy/objects/my_delay2",
																					"name": "my_delay2", "milliseconds": 50},
																					{"delay_type": "fixed",
																					"delayPolicyRef": "/stubo/api/v2/delay-policy/objects/my_delay1",
																					"name": "my_delay1",
																					"milliseconds": 50}]}`
	server, c := testTools(200, testData)
	m := setup(*c)

	defer server.Close()

	//Testing get all delay policies
	req, err := http.NewRequest("GET", "/stubo/api/delete/delay_policy", nil)
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusOK)
}

func TestBeginSessionHandler(t *testing.T) {
	testData := `begin session`
	server, c := testTools(200, testData)
	m := setup(*c)

	defer server.Close()

	//Testing get all delay policies
	req, err := http.NewRequest("GET", "/stubo/api/begin/session?scenario=scenario_x&session=session_x&mode=record", nil)
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusOK)
}

func TestBeginSessionHandlerMissingSession(t *testing.T) {
	testData := `begin session`
	server, c := testTools(200, testData)
	m := setup(*c)

	defer server.Close()

	//Testing get all delay policies
	req, err := http.NewRequest("GET", "/stubo/api/begin/session?scenario=scenario_x&mode=record", nil)
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusBadRequest)
}

func TestBeginSessionHandlerMissingScenario(t *testing.T) {
	testData := `begin session`
	server, c := testTools(200, testData)
	m := setup(*c)

	defer server.Close()

	//Testing get all delay policies
	req, err := http.NewRequest("GET", "/stubo/api/begin/session?session=session&mode=record", nil)
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusBadRequest)
}

func TestBeginSessionHandlerMissingMode(t *testing.T) {
	testData := `begin session`
	server, c := testTools(200, testData)
	m := setup(*c)

	defer server.Close()

	//Testing get all delay policies
	req, err := http.NewRequest("GET", "/stubo/api/begin/session?session=session&scenario=sce1", nil)
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusBadRequest)
}

func TestEndSessionsHandler(t *testing.T) {
	testData := `begin session`
	server, c := testTools(200, testData)
	m := setup(*c)

	defer server.Close()

	//Testing get all delay policies
	req, err := http.NewRequest("GET", "/stubo/api/end/sessions?scenario=sce1", nil)
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusOK)
}

func TestEndSessionsHandlerMissingScenario(t *testing.T) {
	testData := ``
	server, c := testTools(200, testData)
	m := setup(*c)

	defer server.Close()

	//Testing get all delay policies
	req, err := http.NewRequest("GET", "/stubo/api/end/sessions", nil)
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusBadRequest)
}

func TestGetScenariosHandler(t *testing.T) {
	testData := `scenarios`
	server, c := testTools(200, testData)
	m := setup(*c)

	defer server.Close()

	//Testing get all delay policies
	req, err := http.NewRequest("GET", "/stubo/api/get/scenarios", nil)
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)
	// reading resposne body
	body, err := ioutil.ReadAll(respRec.Body)

	expect(t, strings.Contains(string(body), "scenarios"), true)

	expect(t, respRec.Code, http.StatusOK)
}

func TestGetStubResponseHandler(t *testing.T) {
	testData := `Some response`
	server, c := testTools(200, testData)
	m := setup(*c)

	defer server.Close()

	req, err := http.NewRequest("POST", "/stubo/api/get/response?session=sce:x",
		strings.NewReader("anything here, proxy doesn't unmarshall it anyway"))
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)
	// reading resposne body
	body, err := ioutil.ReadAll(respRec.Body)

	expect(t, strings.Contains(string(body), "Some response"), true)

	expect(t, respRec.Code, http.StatusOK)
}

func TestGetStubResponseHanderMissingSession(t *testing.T) {
	testData := `Some response`
	server, c := testTools(200, testData)
	m := setup(*c)

	defer server.Close()

	req, err := http.NewRequest("POST", "/stubo/api/get/response",
		strings.NewReader("anything here, proxy doesn't unmarshall it anyway"))
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)
	// reading resposne body
	_, err = ioutil.ReadAll(respRec.Body)

	expect(t, respRec.Code, http.StatusBadRequest)
}

func TestGetStubResponseHanderMissingScenarioPrefix(t *testing.T) {
	testData := `Some response`
	server, c := testTools(200, testData)
	m := setup(*c)

	defer server.Close()

	req, err := http.NewRequest("POST", "/stubo/api/get/response?session=x",
		strings.NewReader("anything here, proxy doesn't unmarshall it anyway"))
	// no error is expected
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)
	// reading resposne body
	_, err = ioutil.ReadAll(respRec.Body)

	expect(t, respRec.Code, http.StatusBadRequest)
}
