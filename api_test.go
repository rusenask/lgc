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
	response, err := c.getScenarioStubs("scenario1")
	resp := string(response)
	expect(t, len(response), 52)
	expect(t, strings.Contains(resp, "data"), true)
	expect(t, err, nil)
}

func TestGetDelayPolicy(t *testing.T) {
	testData := `{"version":"1.2.3","data": [{"policy": "policy_1"}]}`
	server, c := testTools(200, testData)
	defer server.Close()
	response, err := c.getDelayPolicy("policy_1")
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
