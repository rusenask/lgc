package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func testTools(code int, body string) (*httptest.Server, *Client) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	httpClient := &http.Client{Transport: tr}

	client := &Client{httpClient}
	StuboURI = "http://localhost:3000"
	return server, client
}

func TestGetScenarioStubs(t *testing.T) {
	testData := `{"version":"1.2.3","data": [{"name": "scenario1"}]}`
	server, c := testTools(200, testData)
	defer server.Close()
	name := "scenario_1"
	response, err := c.getScenarioStubs(name)
	resp := string(response)
	expect(t, len(response), 52)
	expect(t, strings.Contains(resp, "data"), true)
	expect(t, err, nil)
}

func TestDeleteScenarioStubs(t *testing.T) {
	testData := `{"version":"1.2.3","data": [{"name": "scenario1"}]}`
	server, c := testTools(200, testData)
	defer server.Close()
	var data APIParams
	data.name = "scenario_1"
	data.force = "true"
	data.targetHost = "somehost"
	response, err := c.deleteScenarioStubs(data)
	resp := string(response)
	expect(t, len(response), 52)
	expect(t, strings.Contains(resp, "data"), true)
	expect(t, err, nil)
}

func TestDeleteScenarioStubsFail(t *testing.T) {
	testData := `{"version":"1.2.3","data": [{"name": "scenario1"}]}`
	server, c := testTools(200, testData)
	defer server.Close()
	var data APIParams
	_, err := c.deleteScenarioStubs(data)
	refute(t, err, nil)
}

func TestGetDelayPolicy(t *testing.T) {
	testData := `{"version":"1.2.3","data": [{"policy": "policy_1"}]}`
	server, c := testTools(200, testData)
	defer server.Close()
	name := "scenario_1"
	response, err := c.getDelayPolicy(name)
	resp := string(response)
	expect(t, len(response), 53)
	expect(t, strings.Contains(resp, "data"), true)
	expect(t, err, nil)
}

func TestGetAllDelayPolicies(t *testing.T) {
	testData := `{"version":"1.2.3","data": [{"some: "data"}]`
	server, c := testTools(200, testData)
	defer server.Close()
	response, err := c.getAllDelayPolicies()
	resp := string(response)
	expect(t, len(response), 45)
	expect(t, strings.Contains(resp, "data"), true)
	expect(t, err, nil)
}

func TestDeleteDelayPolicy(t *testing.T) {
	testdata := `data response`
	server, c := testTools(200, testdata)
	defer server.Close()
	name := "delay_policy_name"
	response, err := c.deleteDelayPolicy(name)
	resp := string(response)
	expect(t, strings.Contains(resp, "data"), true)
	expect(t, err, nil)
}

// TestDeleteAllDelayPolicies passes stubbed response from API v2 containing
// 3 delay policies to deleteAllDelayPolicies function and expects result with
// message that all three policies were deleted. Httptest server returns 200
// for all three deletions
func TestDeleteAllDelayPolicies(t *testing.T) {
	delayPoliciesBytes := []byte(`{"version": "0.6.6",
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
																					"milliseconds": 50}]}`)
	testData := `{"version":"1.2.3","data": [{"some: "data"}]`
	server, c := testTools(200, testData)
	defer server.Close()
	response, err := c.deleteAllDelayPolicies(delayPoliciesBytes)
	resp := string(response)
	fmt.Println(resp)
	expect(t, strings.Contains(resp, "Deleted 3 delay policies: my_delay my_delay2 my_delay1"), true)
	expect(t, err, nil)
}

func TestBeginSession(t *testing.T) {
	testData := `{"version":"1.2.3","data": [{"some: "data"}]`
	server, c := testTools(200, testData)
	defer server.Close()
	response, err := c.beginSession("session", "scenario", "record")
	resp := string(response)
	expect(t, len(response), 45)
	expect(t, strings.Contains(resp, "data"), true)
	expect(t, err, nil)
}

func TestCreateScenario(t *testing.T) {
	testData := `{"version":"1.2.3","data": [{"some: "data"}]`
	server, c := testTools(201, testData)
	defer server.Close()
	response, err := c.createScenario("scenario_1")
	resp := string(response)
	expect(t, len(response), 45)
	expect(t, strings.Contains(resp, "data"), true)
	expect(t, err, nil)
}

func TestGetScenariosDetail(t *testing.T) {
	testData := `{"version":"1.2.3","data": [{"some: "data"}]`
	server, c := testTools(201, testData)
	defer server.Close()
	response, err := c.getScenariosDetail()
	resp := string(response)
	expect(t, len(response), 45)
	expect(t, strings.Contains(resp, "data"), true)
	expect(t, err, nil)
}

func TestGetScenarios(t *testing.T) {
	testData := `{"version":"1.2.3","data": [{"some: "data"}]`
	server, c := testTools(201, testData)
	defer server.Close()
	response, err := c.getScenarios()
	resp := string(response)
	expect(t, len(response), 45)
	expect(t, strings.Contains(resp, "data"), true)
	expect(t, err, nil)
}

func TestEndSessions(t *testing.T) {
	testData := `{"version":"1.2.3","data": [{"some: "data"}]`
	server, c := testTools(201, testData)
	defer server.Close()
	response, err := c.endSessions("scenario")
	resp := string(response)
	expect(t, len(response), 45)
	expect(t, strings.Contains(resp, "data"), true)
	expect(t, err, nil)
}

func TestMakeRequest(t *testing.T) {
	testData := `{"version":"1.2.3","data": [{"some: "data"}]`
	server, c := testTools(201, testData)
	defer server.Close()
	// prepare struct
	path := "/stubo/api/v2/scenarios/objects/some_scenario/action"
	var s params
	s.body = `{"end": "sessions"}`
	s.path = path
	s.method = "POST"
	response, err := c.makeRequest(s)
	resp := string(response)
	expect(t, len(response), 45)
	expect(t, strings.Contains(resp, "data"), true)
	expect(t, err, nil)
}

func TestMakeRequestFail(t *testing.T) {
	testData := `{"version":"1.2.3","data": [{"some: "data"}]`
	server, c := testTools(201, testData)
	defer server.Close()
	// prepare struct
	// path := "/stubo/api/v2/scenarios/objects/some_scenario/action"
	var s params
	StuboURI = "malformed url"
	_, err := c.makeRequest(s)
	refute(t, err, nil)
}

func TestPutStub(t *testing.T) {
	testData := `{  "version": "0.6.6",
								  "data": {
								        "message": {
								            "status": "updated", "msg": "Updated with stateful response",
								            "key": "55dc6cc1938fbef2e62d875c"}
								          }
								  }`
	server, c := testTools(201, testData)
	defer server.Close()

	scenario := "scenario1"
	args := "args=1&arg2=2"
	body := []byte("some body here")

	headers := make(map[string]string)
	headers["session"] = "session_name"
	headers["stateful"] = "true"
	// putting stub
	response, err := c.putStub(scenario, args, body, headers)
	resp := string(response)

	expect(t, strings.Contains(resp, "data"), true)
	expect(t, err, nil)
}

func TestPutStubFailNoSession(t *testing.T) {
	testData := `foo`
	server, c := testTools(200, testData)
	defer server.Close()

	scenario := "scenario1"
	args := "args=1&arg2=2"
	body := []byte("some body here")

	headers := make(map[string]string)
	// omitting session key...
	headers["stateful"] = "true"
	// putting stub
	_, err := c.putStub(scenario, args, body, headers)

	expect(t, strings.Contains(err.Error(), "scenario or session not supplied"), true)
	refute(t, err, nil)
}
