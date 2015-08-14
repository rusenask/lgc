package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type params struct {
	url, body, method string
}

// getStubList calls to Stubo's REST API
// /stubo/api/v2/scenarios/objects/{scenario_name}/stubs/detail
// returns raw response in bytes
func getStubList(scenario string) ([]byte, error) {
	url := "http://localhost:8001/stubo/api/v2/scenarios/objects/" + scenario + "/stubs"
	return GetJSONResponse(url)
}

// getDelayPolicy gets specified delay-policy
// /stubo/api/v2/delay-policy/detail
// returns raw response in bytes
func getDelayPolicy(name string) ([]byte, error) {
	url := "http://localhost:8001/stubo/api/v2/delay-policy/objects/" + name
	return GetJSONResponse(url)
}

func getAllDelayPolicies() ([]byte, error) {
	url := "http://localhost:8001/stubo/api/v2/delay-policy/detail"
	return GetJSONResponse(url)
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
	fmt.Println("URL transformed to: ", s.url)
	fmt.Println("Body: ", s.body)
	var jsonStr = []byte(s.body)
	req, err := http.NewRequest(s.method, s.url, bytes.NewBuffer(jsonStr))
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

// GetJSONResponse calls stubo
func GetJSONResponse(url string) ([]byte, error) {
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
