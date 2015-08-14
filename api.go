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

// begin/session (GET, POST)
//    query args:
//        scenario = scenario name
//        session = session name
//        mode = playback|record
//
// stubo/api/begin/session?scenario=first&session=first_1&mode=playback
func beginSession(session, scenario, mode string) []byte {
	// begin session
	return []byte("nothing yet")
}

func createScenario(scenario string) []byte {
	url := "http://localhost:8001/stubo/api/v2/scenarios"
	var s params
	s.body = "{'scenario':" + scenario + " }"
	s.url = url
	s.method = "PUT"
	makeRequest(s)
	return []byte("nothing yet")
}

func makeRequest(s params) []byte {
	fmt.Println("Transformed to: ", s.url)
	fmt.Println("Body: ", s.body)
	var jsonStr = []byte(s.body)
	req, err := http.NewRequest(s.method, s.url, bytes.NewBuffer(jsonStr))
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// reading body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
	}
	return body
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
