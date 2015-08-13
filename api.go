package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// getStubList calls to Stubo's REST API
// /stubo/api/v2/scenarios/objects/{scenario_name}/stubs/detail
// returns raw response in bytes
func getStubList(scenario string) []byte {
	url := "http://localhost:8001/stubo/api/v2/scenarios/objects/" + scenario + "/stubs"
	fmt.Println("Transformed to: ", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Got error: ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// fmt.Printf("%s", err)
		fmt.Printf("%s", err)
	}
	fmt.Println("Response body got, returning bytes")
	return body
}

// getDelayPolicy gets specified delay-policy
// /stubo/api/v2/delay-policy/detail
// returns raw response in bytes
func getDelayPolicy(name string) []byte {
	url := "http://localhost:8001/stubo/api/v2/delay-policy/objects/" + name
// GetJSONResponse calls stubo
	fmt.Println("Transformed to: ", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Got error: ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
	}
	return body
}
