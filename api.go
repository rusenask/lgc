package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type params struct {
	path, body, method string
	headers            map[string]string
}

// APIParams struct is used to pass parameters to more complex
// API functions such as resource name, host, etc..
type APIParams struct {
	name, targetHost, force string
}

// Client structure to be injected into functions to perform HTTP calls
type Client struct {
	HTTPClient *http.Client
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// New returns an error that formats as the given text.
func New(text string) error {
	return &errorString{text}
}

// getStubList calls to Stubo's REST API
// /stubo/api/v2/scenarios/objects/{scenario_name}/stubs/detail
// returns raw response in bytes
func (c *Client) getScenarioStubs(scenario string) ([]byte, error) {
	if scenario != "" {
		path := "/stubo/api/v2/scenarios/objects/" + scenario + "/stubs"

		// setting logger
		method := trace()
		log.WithFields(log.Fields{
			"scenario":      scenario,
			"urlPath":       path,
			"headers":       "",
			"requestMethod": "GET",
			"method":        method,
		}).Debug("Getting scenario stubs")

		return c.GetResponseBody(path)
	}
	return []byte(""), errors.New("api.getScenarioStubs error: scenario name not supplied")
}

// deleteScenarioStubs takes apiParams as an argument which contains
// scenario name and two optional parameters for headers:
// "force" which defaults to false and "targetHost" which can specify another
// host
func (c *Client) deleteScenarioStubs(p APIParams) ([]byte, error) {
	var s params
	// adding path
	if p.name != "" {
		path := "/stubo/api/v2/scenarios/objects/" + p.name + "/stubs"
		s.path = path
		// creating MAP for headers
		headers := make(map[string]string)
		if p.force != "" {
			headers["force"] = p.force
		}
		if p.targetHost != "" {
			headers["target_host"] = p.targetHost
		}
		s.headers = headers
		s.method = "DELETE"
		fmt.Println(s.headers)

		// setting logger
		method := trace()
		log.WithFields(log.Fields{
			"scenario":      p.name,
			"urlPath":       s.path,
			"headers":       s.headers,
			"requestMethod": s.method,
			"method":        method,
		}).Debug("Deleting scenario stubs")

		// calling delete
		return c.makeRequest(s)
	}
	return []byte(""), errors.New("api.deleteScenarioStubs error: scenario name not supplied")
}

// getDelayPolicy gets specified delay-policy
// /stubo/api/v2/delay-policy/detail
// returns raw response in bytes
func (c *Client) getDelayPolicy(name string) ([]byte, error) {
	if name != "" {
		path := "/stubo/api/v2/delay-policy/objects/" + name
		return c.GetResponseBody(path)
	}
	return []byte(""), errors.New("api.getDelayPolicy error: delay policy name supplied")
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
	var jsonStr = []byte(s.body)

	// logging get transformation
	method := trace()
	log.WithFields(log.Fields{
		"method":        method,
		"url":           url,
		"body":          s.body,
		"headers":       s.headers,
		"requestMethod": s.method,
	}).Info("Transforming URL, preparing for request to Stubo")

	req, err := http.NewRequest(s.method, url, bytes.NewBuffer(jsonStr))
	if s.headers != nil {
		for k, v := range s.headers {
			req.Header.Set(k, v)
		}
	}
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		// logging read error
		log.WithFields(log.Fields{
			"error":  err.Error(),
			"method": method,
			"url":    url,
		}).Warn("Failed to get response from Stubo!")

		return []byte(""), err
	}
	defer resp.Body.Close()
	// reading body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// logging read error
		log.WithFields(log.Fields{
			"error":  err.Error(),
			"method": method,
			"url":    url,
		}).Warn("Failed to read response from Stubo!")

		return []byte(""), err
	}
	return body, nil
}

// GetResponseBody calls stubo
func (c *Client) GetResponseBody(path string) ([]byte, error) {
	url := StuboURI + path
	// logging get transformation
	method := trace()
	log.WithFields(log.Fields{
		"method": method,
		"url":    url,
	}).Info("Transforming URL, getting response body")
	resp, err := c.HTTPClient.Get(url)

	if err != nil {
		// logging get error
		log.WithFields(log.Fields{
			"error":  err.Error(),
			"method": method,
			"url":    url,
		}).Warn("Failed to get response from Stubo!")

		return []byte(""), err
	}
	defer resp.Body.Close()
	// reading resposne body
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// logging read error
		log.WithFields(log.Fields{
			"error":  err.Error(),
			"method": method,
			"url":    url,
		}).Warn("Failed to read response from Stubo!")

		return []byte(""), err
	}
	return body, nil
}
