package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// getStubList calls to Stubo's REST API
// /stubo/api/v2/scenarios/objects/{scenario_name}/stubs/detail
func getStubList(scenario string) string {
	fmt.Println("got scenario to list stubs:", scenario)
	resp, err := http.Get("http://localhost:8001/stubo/api/v2/scenarios/objects/" + scenario + "/stubs")
	if err != nil {
		fmt.Println("Got error: ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
	}
	// fmt.Printf("%s\n", string(body))
	return string(body)
}
