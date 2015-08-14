package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type params struct {
	path, body, method string
}

// getStubList calls to Stubo's REST API
// /stubo/api/v2/scenarios/objects/{scenario_name}/stubs/detail
// returns raw response in bytes
func getStubList(scenario string) ([]byte, error) {
	url := "http://localhost:8001/stubo/api/v2/scenarios/objects/" + scenario + "/stubs"
	return GetResponseBody(url)
}

// getDelayPolicy gets specified delay-policy
// /stubo/api/v2/delay-policy/detail
// returns raw response in bytes
func getDelayPolicy(name string) ([]byte, error) {
	url := "http://localhost:8001/stubo/api/v2/delay-policy/objects/" + name
	return GetResponseBody(url)
}

func getAllDelayPolicies() ([]byte, error) {
	url := "http://localhost:8001/stubo/api/v2/delay-policy/detail"
	return GetResponseBody(url)
}

// beginSession takes session, scenario, mode parameters. Can either
// set playback or record modes
func beginSession(session, scenario, mode string) ([]byte, error) {
	url := "http://localhost:8001/stubo/api/v2/scenarios/objects/" + scenario + "/action"
	var s params
	s.body = `{"begin": null, "session": "` + session + `",  "mode": "` + mode + `"}`
	fmt.Println("formated body for session begin: ", s.body)
	s.url = url
	s.method = "POST"
	return makeRequest(s)
}

func createScenario(scenario string) ([]byte, error) {
	url := "http://localhost:8001/stubo/api/v2/scenarios"
	var s params
	s.body = `{"scenario": "` + scenario + `"}`
	fmt.Println("formated body: ", s.body)
	s.url = url
	s.method = "PUT"
	return makeRequest(s)
}

// getScenariosDetail gets and returns all scenarios with details
func getScenariosDetail() ([]byte, error) {
	url := "http://localhost:8001/stubo/api/v2/scenarios/detail"
	return GetResponseBody(url)
}

func endSessions(scenario string) ([]byte, error) {
	url := "http://localhost:8001/stubo/api/v2/scenarios/objects/" + scenario + "/action"
	var s params
	s.body = `{"end": "sessions"}`
	fmt.Println("formated body for session begin: ", s.body)
	s.url = url
	s.method = "POST"
	return makeRequest(s)
}

func makeRequest(s params) ([]byte, error) {
	url := StuboURI + s.path
	fmt.Println("URL transformed to: ", url)
	fmt.Println("Body: ", s.body)
	var jsonStr = []byte(s.body)
	req, err := http.NewRequest(s.method, url, bytes.NewBuffer(jsonStr))
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte(""), err
	}
	defer resp.Body.Close()
	// reading body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}
	return body, nil
}

// GetResponseBody calls stubo
func GetResponseBody(path string) ([]byte, error) {
	url := StuboURI + path
	fmt.Println("Transformed to: ", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		return []byte(""), err
	}
	defer resp.Body.Close()
	// reading resposne body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		return []byte(""), err
	}
	return body, nil
}
