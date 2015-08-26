package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// HandlerHTTPClient is used to inject http.Client to handlers
type HandlerHTTPClient struct {
	http Client
}

// DelayPolicy structure for gettting delay policy references
type DelayPolicy struct {
	Name string `json:"name"`
	Ref  string `json:"delayPolicyRef"`
}

// DelayPolicyResponse structure for unmarshaling JSON structures from API v2
type DelayPolicyResponse struct {
	Data    []DelayPolicy `json:"data"`
	Version string        `json:"version"`
}

// ResponseToClient is a helper struct for artificially forming responses to clients
type ResponseToClient struct {
	Version string            `json:"data"`
	Data    map[string]string `json:"version"`
}

func httperror(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.WithFields(log.Fields{
			"url_query": r.URL.Query(),
			"url_path":  r.URL.Path,
			"error":     err.Error(),
		}).Error("Got error during HTTP request to Stubo")
	}
}

// trace returns name of the current function
func trace() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

// stublistHandler gets stubs, e.g.: stubo/api/get/stublist?scenario=first
func (h HandlerHTTPClient) stublistHandler(w http.ResponseWriter, r *http.Request) {
	scenario, ok := r.URL.Query()["scenario"]

	// setting context logger
	method := trace()
	handlersContextLogger := log.WithFields(log.Fields{
		"url_query": r.URL.Query(),
		"url_path":  r.URL.Path,
		"func":      method,
	})

	if ok {
		handlersContextLogger.Info("Got query")

		client := h.http

		// expecting one param - scenario
		response, err := client.getScenarioStubs(scenario[0])

		// checking whether we got good response
		httperror(w, r, err)

		// setting resposne header
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		// logging error
		handlersContextLogger.Warn("Scenario name was not provided")

		http.Error(w, "Scenario name not provided.", 400)
	}
}

// deleteStubsHandler deletes scenario stubs, e.g.: stubo/api/delete/stubs?scenario=first
// optional arguments host=your_host, force=true/false (defaults to false)
func (h HandlerHTTPClient) deleteStubsHandler(w http.ResponseWriter, r *http.Request) {
	scenario, ok := r.URL.Query()["scenario"]
	// setting context logger
	method := trace()
	handlersContextLogger := log.WithFields(log.Fields{
		"url_query": r.URL.Query(),
		"url_path":  r.URL.Path,
		"func":      method,
	})

	if ok {
		handlersContextLogger.Info("Got query")

		// expecting params - scenario, host, force
		client := h.http
		var p APIParams
		p.name = scenario[0]
		force, ok := r.URL.Query()["force"]
		if ok {
			p.force = force[0]
		}
		host, ok := r.URL.Query()["host"]
		if ok {
			p.targetHost = host[0]
		}
		response, err := client.deleteScenarioStubs(p)
		// checking whether we got good response
		httperror(w, r, err)
		// setting resposne header
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		msg := "Scenario name not provided."
		handlersContextLogger.Warn(msg)
		http.Error(w, msg, 400)
	}
}

func getURLHeadersArgs(expectedHeaders map[string]bool, urlQuery map[string][]string) (map[string]string, string) {
	var bufferArgs bytes.Buffer
	headers := make(map[string]string)
	// getting more headers and forming query argument
	for key, value := range urlQuery {
		// if key is in expected headers (Stubo expects these to be in headers
		// instead of URL query in API v2, transforming them..)
		if expectedHeaders[key] {
			headers[key] = value[0]
		} else {
			bufferArgs.WriteString(key + "=" + value[0] + "&")
		}
	}
	args := bufferArgs.String()
	return headers, args
}

// putStubHandler takes in POST request from client, transforms URL query arguments
// to header values and calls another function that calls Stubo API v2, returns
// response bytes without unmarshalling/marshalling them
func (h HandlerHTTPClient) putStubHandler(w http.ResponseWriter, r *http.Request) {
	urlQuery := r.URL.Query()
	// getting session name
	session, ok := urlQuery["session"]
	client := h.http

	// setting context logger
	method := trace()
	handlersContextLogger := log.WithFields(log.Fields{
		"url_query": urlQuery,
		"url_path":  r.URL.Path,
		"func":      method,
	})
	if ok {
		// session name is present, moving forward
		ScenarioSession := session[0]
		slices := strings.Split(ScenarioSession, ":")
		// check whether user has supplied scenario name as well
		if len(slices) < 2 {
			msg := "Bad request, missing session or scenario name. When under proxy, please use 'scenario:session' format in your" +
				"URL query, such as '/stubo/api/put/stub?session=scenario:session_name' "
			handlersContextLogger.Warn(msg)
			log.Warn(msg)
			http.Error(w, msg, 400)
			return
		}
		scenario := slices[0]

		// removing session from the MAP
		delete(urlQuery, "session")

		// these URL query arguments are expected and should be converted to headers
		expectedHeaders := map[string]bool{
			"ext_module":        true,
			"delay_policy":      true,
			"stateful":          true,
			"stub_created_date": true,
		}
		headers, args := getURLHeadersArgs(expectedHeaders, urlQuery)
		headers["session"] = slices[1]

		defer r.Body.Close()
		// reading resposne body
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			// logging read error
			log.WithFields(log.Fields{
				"error": err.Error(),
				"func":  method,
			}).Warn("Failed to read request body!")
		}
		// putting stub
		response, err := client.putStub(scenario, args, body, headers)
		// checking whether we got good response
		httperror(w, r, err)
		// setting resposne header
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)

	} else {
		msg := "Bad request, missing session name."
		handlersContextLogger.Warn(msg)
		http.Error(w, msg, 400)
	}
}

// getSession looks for session both in URL query and request headers
func getSession(r *http.Request) (string, bool) {
	urlQuery := r.URL.Query()
	// getting session name
	session, ok := urlQuery["session"]
	if ok {
		return session[0], true
	}
	sessionName := r.Header.Get("Stubo-Request-Session")
	if sessionName != "" {
		return sessionName, true
	}
	return "", false
}

func (h HandlerHTTPClient) getStubResponseHandler(w http.ResponseWriter, r *http.Request) {
	urlQuery := r.URL.Query()
	// getting session name
	ScenarioSession, ok := getSession(r)

	client := h.http

	// setting context logger
	method := trace()
	handlersContextLogger := log.WithFields(log.Fields{
		"url_query": urlQuery,
		"url_path":  r.URL.Path,
		"func":      method,
	})
	if ok {
		// session name is present, moving forward
		slices := strings.Split(ScenarioSession, ":")
		// check whether user has supplied scenario name as well
		if len(slices) < 2 {
			msg := "Bad request, missing session or scenario name. When under proxy, please use 'scenario:session' format in your" +
				"URL query, such as '/stubo/api/get/response?session=scenario:session_name' "
			handlersContextLogger.Warn(msg)
			log.Warn(msg)
			http.Error(w, msg, 400)
			return
		}
		scenario := slices[0]

		// removing session from the MAP
		delete(urlQuery, "session")

		// these URL query arguments are expected and should be converted to headers
		expectedHeaders := map[string]bool{}
		headers, args := getURLHeadersArgs(expectedHeaders, urlQuery)
		headers["session"] = slices[1]

		log.WithFields(log.Fields{
			"headers":  headers,
			"args":     args,
			"scenario": scenario,
		}).Info("Get response Args and Headers created...")

		defer r.Body.Close()
		// reading resposne body
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			// logging read error
			log.WithFields(log.Fields{
				"error": err.Error(),
				"func":  method,
			}).Warn("Failed to read request body!")
		}
		// Getting stubo response to request
		response, err := client.getStubResponse(scenario, args, body, headers)
		// checking whether we got good response
		httperror(w, r, err)
		// setting resposne header
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)

	} else {
		msg := "Bad request, missing session name."
		handlersContextLogger.Warn(msg)
		http.Error(w, msg, 400)
	}
}

// getDelayPolicyHandler - returns delay policy information, list all if
// name is not provided, e.g.: stubo/api/get/delay_policy?name=slow
func (h HandlerHTTPClient) getDelayPolicyHandler(w http.ResponseWriter, r *http.Request) {
	name, ok := r.URL.Query()["name"]
	client := h.http
	// setting context logger
	method := trace()
	handlersContextLogger := log.WithFields(log.Fields{
		"url_query": r.URL.Query(),
		"url_path":  r.URL.Path,
		"func":      method,
	})

	if ok {

		handlersContextLogger.Info("Got query")
		// expecting one param - scenario
		response, err := client.getDelayPolicy(name[0])
		// checking whether we got good response
		httperror(w, r, err)
		// setting resposne header
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		// name is not provided, getting all delay policies
		response, err := client.getAllDelayPolicies()
		// checking whether we got good response
		httperror(w, r, err)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

// putDelayPolicyHandler takes URL query arguments and turns them into JSON
// example query: stubo/api/put/delay_policy?name=slow&delay_type=fixed&milliseconds=1000
func (h HandlerHTTPClient) putDelayPolicyHandler(w http.ResponseWriter, r *http.Request) {
	urlQuery := r.URL.Query()
	client := h.http
	// query MAP stores key/value pairs, although all the values are of array type []
	newQuery := make(map[string]string)
	for key, value := range urlQuery {
		// taking only first argument
		newQuery[key] = value[0]
	}

	// setting context logger
	method := trace()

	// converting query MAP to JSON key/value pairs
	jsonString, err := json.Marshal(newQuery)

	httperror(w, r, err)

	handlersContextLogger := log.WithFields(log.Fields{
		"url_query": r.URL.Query(),
		"url_path":  r.URL.Path,
		"func":      method,
	})

	handlersContextLogger.Info("Got query to create new delay policy.")

	response, err := client.putDelayPolicy(jsonString)
	httperror(w, r, err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// deleteDelayPolicyHandler - deletes delay policy
// stubo/api/delete/delay_policy?name=slow
func (h HandlerHTTPClient) deleteDelayPolicyHandler(w http.ResponseWriter, r *http.Request) {
	name, ok := r.URL.Query()["name"]
	client := h.http

	// setting context logger
	method := trace()
	handlersContextLogger := log.WithFields(log.Fields{
		"url_query": r.URL.Query(),
		"url_path":  r.URL.Path,
		"func":      method,
	})

	if ok {
		handlersContextLogger.Info("Deleting specified delay policy")
		// expecting one param - name
		response, err := client.deleteDelayPolicy(name[0])
		// checking whether we got good response
		httperror(w, r, err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		handlersContextLogger.Info("Deleting all delay policies in two steps")
		delayPolicies, err := client.getAllDelayPolicies()
		httperror(w, r, err)
		if err == nil {
			handlersContextLogger.Info("Got all delay policies, deleting one by one")
			response, err := client.deleteAllDelayPolicies(delayPolicies)
			httperror(w, r, err)
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		}
	}
}

// deleteAllDelayPolicies - custom handler to delete multiple delay policies.
// This API call is not directly available through API v2 so we are taking
// response with all delay policies - unmarshalling it, getting all names
// and then deleting them one by one
func (c *Client) deleteAllDelayPolicies(dp []byte) ([]byte, error) {
	// getting all delay policy names
	allDelayPolicies := dp
	// Unmarshaling JSON
	var data DelayPolicyResponse
	err := json.Unmarshal(allDelayPolicies, &data)

	// logging
	method := trace()
	log.WithFields(log.Fields{
		"func":          method,
		"delayPolicies": data,
	}).Info("Deleting delay policies")

	if err != nil {
		return []byte(""), err
	}
	// Getting stubo version
	version := data.Version

	// Deleting delay policies
	var responses []string
	for _, dp := range data.Data {
		_, err := c.deleteDelayPolicy(dp.Name)

		if err == nil {
			responses = append(responses, dp.Name)
		} else {
			log.WithFields(log.Fields{
				"func":  method,
				"error": err.Error(),
			}).Warn("Failed to delete delay policy")
		}
	}
	// creating message for the client
	message := fmt.Sprintf("Deleted %d delay policies: ", len(responses)) + strings.Join(responses, " ")

	log.WithFields(log.Fields{
		"func":     method,
		"response": message,
	}).Info("Delay policies deleted")
	// creating structure for the response
	res := &ResponseToClient{
		Version: version,
		Data:    map[string]string{"message": message},
	}
	// encoding to JSON and returning
	respBytes, err := json.Marshal(res)
	return respBytes, err
}

// begin/session (GET, POST)
// stubo/api/begin/session?scenario=first&session=first_1&mode=playback
func (h HandlerHTTPClient) beginSessionHandler(w http.ResponseWriter, r *http.Request) {
	queryArgs, _ := url.ParseQuery(r.URL.RawQuery)

	// setting context logger
	method := trace()
	handlersContextLogger := log.WithFields(log.Fields{
		"url_query": queryArgs,
		"url_path":  r.URL.Path,
		"func":      method,
	})

	handlersContextLogger.Info("Begin session...")
	// retrieving details and validating request
	if scenario, ok := queryArgs["scenario"]; ok {
		if session, ok := queryArgs["session"]; ok {
			if mode, ok := queryArgs["mode"]; ok {
				// Create scenario. This can result in 422 (duplicate error) and this is
				// fine, since we must only ensure that it exists.
				client := h.http
				_, err := client.createScenario(scenario[0])
				httperror(w, r, err)
				// Begin session
				response, err := client.beginSession(session[0], scenario[0], mode[0])
				httperror(w, r, err)
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)
			} else {
				msg := "Bad request, missing session mode key."
				handlersContextLogger.Warn(msg)
				http.Error(w, msg, 400)
			}
		} else {
			msg := "Bad request, missing session name."
			handlersContextLogger.Warn(msg)
			http.Error(w, msg, 400)
		}
	} else {
		msg := "Bad request, missing scenario name."
		handlersContextLogger.Warn(msg)
		http.Error(w, msg, 400)
	}
}

func (h HandlerHTTPClient) endSessionsHandler(w http.ResponseWriter, r *http.Request) {

	// setting context logger
	method := trace()
	handlersContextLogger := log.WithFields(log.Fields{
		"url_query": r.URL.Query(),
		"url_path":  r.URL.Path,
		"func":      method,
	})
	scenario, ok := r.URL.Query()["scenario"]
	if ok {
		handlersContextLogger.Info("Ending session...")
		// expecting one param - scenario
		client := h.http
		response, err := client.endSessions(scenario[0])
		// checking whether we got good response
		httperror(w, r, err)
		// setting resposne header
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		msg := "Scenario name not provided."
		handlersContextLogger.Warn(msg)
		http.Error(w, msg, 400)
	}
}

func (h HandlerHTTPClient) getScenariosHandler(w http.ResponseWriter, r *http.Request) {
	client := h.http

	// setting logger
	method := trace()
	log.WithFields(log.Fields{
		"url_query": r.URL.Query(),
		"url_path":  r.URL.Path,
		"func":      method,
	}).Info("Getting scenarios")

	response, err := client.getScenarios()
	// checking whether we got good response
	httperror(w, r, err)
	// setting resposne header
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}
