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
	bodyBytes          []byte
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

// putStub transparently passes request body to Stubo
func (c *Client) putStub(scenario, args string, body []byte, headers map[string]string) ([]byte, error) {
	if scenario != "" && headers["session"] != "" {
		var s params

		path := "/stubo/api/v2/scenarios/objects/" + scenario + "/stubs?" + args

		s.path = path
		s.headers = headers
		s.method = "PUT"

		// assigning body in bytes
		s.bodyBytes = body
		// setting logger
		method := trace()
		log.WithFields(log.Fields{
			"scenario":      scenario,
			"session":       headers["session"],
			"urlPath":       path,
			"headers":       "",
			"requestMethod": "GET",
			"func":          method,
		}).Debug("Adding stub to scenario")

		return c.makeRequest(s)
	}
	return []byte(""), errors.New("api.putStub error: scenario or session not supplied")
}

// getStubList calls to Stubo's REST API
// /stubo/api/v2/scenarios/objects/{scenario_name}/stubs/detail
// returns raw response in bytes
func (c *Client) getScenarioStubs(name string) ([]byte, error) {
	if name != "" {
		path := "/stubo/api/v2/scenarios/objects/" + name + "/stubs"

		// setting logger
		method := trace()
		log.WithFields(log.Fields{
			"name":          name,
			"urlPath":       path,
			"headers":       "",
			"requestMethod": "GET",
			"func":          method,
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

		// setting logger
		method := trace()
		log.WithFields(log.Fields{
			"name":          p.name,
			"urlPath":       s.path,
			"headers":       s.headers,
			"requestMethod": s.method,
			"func":          method,
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
		// setting logger
		method := trace()
		log.WithFields(log.Fields{
			"name":          name,
			"urlPath":       path,
			"headers":       "",
			"requestMethod": "GET",
			"func":          method,
		}).Debug("Getting specified delay policy")

		return c.GetResponseBody(path)
	}
	return []byte(""), errors.New("api.getDelayPolicy error: delay policy name supplied")
}

func (c *Client) getAllDelayPolicies() ([]byte, error) {
	path := "/stubo/api/v2/delay-policy/detail"

	// setting logger
	method := trace()
	log.WithFields(log.Fields{
		"name":          "",
		"urlPath":       path,
		"headers":       "",
		"requestMethod": "GET",
		"func":          method,
	}).Debug("Getting all delay policies")

	return c.GetResponseBody(path)
}

func (c *Client) deleteDelayPolicy(name string) ([]byte, error) {
	path := "/stubo/api/v2/delay-policy/objects/" + name
	var s params
	s.path = path
	s.method = "DELETE"

	// setting logger
	method := trace()
	log.WithFields(log.Fields{
		"name":          name,
		"urlPath":       s.path,
		"headers":       "",
		"requestMethod": s.method,
		"func":          method,
	}).Debug("Deleting specified delay policy")

	return c.makeRequest(s)
}

// beginSession takes session, scenario, mode parameters. Can either
// set playback or record modes
func (c *Client) beginSession(session, scenario, mode string) ([]byte, error) {
	path := "/stubo/api/v2/scenarios/objects/" + scenario + "/action"
	var s params
	s.body = `{"begin": null, "session": "` + session + `",  "mode": "` + mode + `"}`
	s.path = path
	s.method = "POST"

	// setting logger
	method := trace()
	log.WithFields(log.Fields{
		"scenario":      scenario,
		"session":       session,
		"urlPath":       s.path,
		"headers":       "",
		"body":          s.body,
		"requestMethod": s.method,
		"func":          method,
	}).Debug("Begin session")

	return c.makeRequest(s)
}

func (c *Client) createScenario(scenario string) ([]byte, error) {
	path := "/stubo/api/v2/scenarios"
	var s params
	s.body = `{"scenario": "` + scenario + `"}`
	fmt.Println("formated body: ", s.body)
	s.path = path
	s.method = "PUT"

	// setting logger
	method := trace()
	log.WithFields(log.Fields{
		"name":          scenario,
		"urlPath":       s.path,
		"headers":       "",
		"body":          "",
		"requestMethod": s.method,
		"func":          method,
	}).Debug("Creating scenario")

	return c.makeRequest(s)
}

// getScenariosDetail gets and returns all scenarios with details
func (c *Client) getScenariosDetail() ([]byte, error) {
	path := "/stubo/api/v2/scenarios/detail"

	// setting logger
	method := trace()
	log.WithFields(log.Fields{
		"name":          "",
		"urlPath":       path,
		"headers":       "",
		"body":          "",
		"requestMethod": "",
		"func":          method,
	}).Debug("Getting scenario details")

	return c.GetResponseBody(path)
}

// getScenarios gets and returns all scenarios with details
func (c *Client) getScenarios() ([]byte, error) {
	path := "/stubo/api/v2/scenarios"

	// setting logger
	method := trace()
	log.WithFields(log.Fields{
		"name":          "",
		"urlPath":       path,
		"headers":       "",
		"body":          "",
		"requestMethod": "",
		"func":          method,
	}).Debug("Getting scenarios")

	return c.GetResponseBody(path)
}

// endSessions ends all specified scenario sessions
func (c *Client) endSessions(scenario string) ([]byte, error) {
	path := "/stubo/api/v2/scenarios/objects/" + scenario + "/action"
	var s params
	s.body = `{"end": "sessions"}`
	s.path = path
	s.method = "POST"

	// setting logger
	method := trace()
	log.WithFields(log.Fields{
		"name":          scenario,
		"urlPath":       s.path,
		"headers":       "",
		"body":          s.body,
		"requestMethod": s.method,
		"func":          method,
	}).Debug("Ending sessions")

	return c.makeRequest(s)
}

// makeRequest takes Params struct as paramateres and makes request to Stubo
// then gets response bytes and returns to caller
func (c *Client) makeRequest(s params) ([]byte, error) {
	url := StuboURI + s.path
	if s.bodyBytes == nil {
		s.bodyBytes = []byte(s.body)
	}

	// logging get transformation
	method := trace()
	log.WithFields(log.Fields{
		"func":          method,
		"url":           url,
		"body":          s.body,
		"headers":       s.headers,
		"requestMethod": s.method,
	}).Info("Transforming URL, preparing for request to Stubo")

	req, err := http.NewRequest(s.method, url, bytes.NewBuffer(s.bodyBytes))
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
			"error": err.Error(),
			"func":  method,
			"url":   url,
		}).Warn("Failed to get response from Stubo!")

		return []byte(""), err
	}
	defer resp.Body.Close()
	// reading body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// logging read error
		log.WithFields(log.Fields{
			"error": err.Error(),
			"func":  method,
			"url":   url,
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
		"func": method,
		"url":  url,
	}).Info("Transforming URL, getting response body")
	resp, err := c.HTTPClient.Get(url)

	if err != nil {
		// logging get error
		log.WithFields(log.Fields{
			"error": err.Error(),
			"func":  method,
			"url":   url,
		}).Warn("Failed to get response from Stubo!")

		return []byte(""), err
	}
	defer resp.Body.Close()
	// reading resposne body
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// logging read error
		log.WithFields(log.Fields{
			"error": err.Error(),
			"func":  method,
			"url":   url,
		}).Warn("Failed to read response from Stubo!")

		return []byte(""), err
	}
	return body, nil
}
