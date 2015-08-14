package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
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
	return server, client
}

func TestGetScenarios(t *testing.T) {
	server, c := testTools(200, `{"version":"1.2.3","data": [{"name","scenario1"}]}`)
	defer server.Close()
	responses, err := c.getScenarioStubs("scenario1")

	expect(t, len(responses), 1)
	expect(t, err, nil)

	correctResponse := `{"version":"1.2.3","data": [{"name","scenario1"}]}`
	expect(t, reflect.DeepEqual(correctResponse, responses[0]), true)
}
