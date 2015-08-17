package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type params struct {
	path, body, method string
}

// Client structure to be used by HTTP
type Client struct {
	HTTPClient *http.Client
}

type DelayPolicy struct {
	Name string `json:"delayPolicyRef"`
}

type ResponseData struct {
	Data []DelayPolicy `json:"data"`
}

// getStubList calls to Stubo's REST API
// /stubo/api/v2/scenarios/objects/{scenario_name}/stubs/detail
// returns raw response in bytes
func (c *Client) getScenarioStubs(scenario string) ([]byte, error) {
	fmt.Println(StuboConfig.StuboHost)
	path := "/stubo/api/v2/scenarios/objects/" + scenario + "/stubs"
	return c.GetResponseBody(path)
}

// getDelayPolicy gets specified delay-policy
// /stubo/api/v2/delay-policy/detail
// returns raw response in bytes
func (c *Client) getDelayPolicy(name string) ([]byte, error) {
	path := "/stubo/api/v2/delay-policy/objects/" + name
	return c.GetResponseBody(path)
}

func (c *Client) getAllDelayPolicies() ([]byte, error) {
	path := "/stubo/api/v2/delay-policy/detail"
	return c.GetResponseBody(path)
}

func (c *Client) deleteDelayPolicy(name string) ([]byte, error) {
	path := "/stubo/api/v2/delay-policy/objects/" + name
	var s params
	s.path = path
	s.method = "DELETE"
	return c.makeRequest(s)
}

func (c *Client) deleteAllDelayPolicies() ([]byte, error) {
	// getting all delay policy names
	allDelayPolicies, err := c.getAllDelayPolicies()
	if err != nil {
		return []byte(""), err
	}
	// Unmarshaling JSON
	var data ResponseData
	err = json.Unmarshal(allDelayPolicies, &data)
	fmt.Println(data)
	for _, dp := range data.Data {
		fmt.Println(dp.Name)
	}
	if err != nil {
		fmt.Println("error occured")
	}
	return []byte(""), nil
}

// beginSession takes session, scenario, mode parameters. Can either
// set playback or record modes
func (c *Client) beginSession(session, scenario, mode string) ([]byte, error) {
	path := "/stubo/api/v2/scenarios/objects/" + scenario + "/action"
	var s params
	s.body = `{"begin": null, "session": "` + session + `",  "mode": "` + mode + `"}`
	fmt.Println("formated body for session begin: ", s.body)
	s.path = path
	s.method = "POST"
	return c.makeRequest(s)
}

func (c *Client) createScenario(scenario string) ([]byte, error) {
	path := "/stubo/api/v2/scenarios"
	var s params
	s.body = `{"scenario": "` + scenario + `"}`
	fmt.Println("formated body: ", s.body)
	s.path = path
	s.method = "PUT"
	return c.makeRequest(s)
}

// getScenariosDetail gets and returns all scenarios with details
func (c *Client) getScenariosDetail() ([]byte, error) {
	path := "/stubo/api/v2/scenarios/detail"
	return c.GetResponseBody(path)
}

// getScenarios gets and returns all scenarios with details
func (c *Client) getScenarios() ([]byte, error) {
	path := "/stubo/api/v2/scenarios"
	return c.GetResponseBody(path)
}

// endSessions ends all specified scenario sessions
func (c *Client) endSessions(scenario string) ([]byte, error) {
	path := "/stubo/api/v2/scenarios/objects/" + scenario + "/action"
	var s params
	s.body = `{"end": "sessions"}`
	fmt.Println("formated body for session begin: ", s.body)
	s.path = path
	s.method = "POST"
	return c.makeRequest(s)
}

func (c *Client) makeRequest(s params) ([]byte, error) {
	url := StuboURI + s.path
	fmt.Println("URL transformed to: ", url)
	fmt.Println("Body: ", s.body)
	var jsonStr = []byte(s.body)
	req, err := http.NewRequest(s.method, url, bytes.NewBuffer(jsonStr))
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
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
func (c *Client) GetResponseBody(path string) ([]byte, error) {
	url := StuboURI + path
	fmt.Println("Transformed to: ", url)
	resp, err := c.HTTPClient.Get(url)
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
