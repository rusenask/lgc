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

func TestGetScenarios(t *testing.T) {
	testData := `{"version":"1.2.3","data": [{"name","scenario1"}]}`
	server, c := testTools(200, testData)
	defer server.Close()
	response, err := c.getScenarioStubs("scenario1")
	resp := string(response)
	// fmt.Println(strings.Contains(resp, "data"))
	// fmt.Println(len(response))
	expect(t, len(response), 51)
	expect(t, strings.Contains(resp, "data"), true)
	expect(t, err, nil)
}
