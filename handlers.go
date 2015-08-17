package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

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

func httperror(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// stublistHandler gets stubs, e.g.: stubo/api/get/stublist?scenario=first
func stublistHandler(w http.ResponseWriter, r *http.Request) {
	scenario, ok := r.URL.Query()["scenario"]
	if ok {
		fmt.Println("got:", r.URL.Query())
		client := &Client{&http.Client{}}
		// expecting one param - scenario
		response, err := client.getScenarioStubs(scenario[0])
		// checking whether we got good response
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		// setting resposne header
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		http.Error(w, "Scenario name not provided.", 400)
	}
}

// deleteStubsHandler deletes scenario stubs, e.g.: stubo/api/delete/stubs?scenario=first
// optional arguments host=your_host, force=true/false (defaults to false)
func deleteStubsHandler(w http.ResponseWriter, r *http.Request) {
	scenario, ok := r.URL.Query()["scenario"]
	if ok {
		// expecting one param - scenario
		client := &Client{&http.Client{}}
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
		httperror(w, err)
		// setting resposne header
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		http.Error(w, "Scenario name not provided.", 400)
	}
}

// getDelayPolicyHandler - returns delay policy information, list all if
// name is not provided, e.g.: stubo/api/get/delay_policy?name=slow
func getDelayPolicyHandler(w http.ResponseWriter, r *http.Request) {
	name, ok := r.URL.Query()["name"]
	client := &Client{&http.Client{}}
	if ok {
		// name provided so looking for specific delay
		fmt.Println("got:", r.URL.Query())
		// expecting one param - scenario
		response, err := client.getDelayPolicy(name[0])
		// checking whether we got good response
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		// setting resposne header
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		// name is not provided, getting all delay policies
		response, err := client.getAllDelayPolicies()
		// checking whether we got good response
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

// deleteDelayPolicyHandler - deletes delay policy
// stubo/api/delete/delay_policy?name=slow
func deleteDelayPolicyHandler(w http.ResponseWriter, r *http.Request) {
	name, ok := r.URL.Query()["name"]
	client := &Client{&http.Client{}}
	if ok {
		// expecting one param - name
		response, err := client.deleteDelayPolicy(name[0])
		// checking whether we got good response
		httperror(w, err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		fmt.Println("Deleting all delay policies")
		delayPolicies, err := client.getAllDelayPolicies()
		httperror(w, err)
		if err == nil {
			response, err := client.deleteAllDelayPolicies(delayPolicies)
			httperror(w, err)
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
	fmt.Println(data)
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
		}
	}
	// creating message for the client
	message := fmt.Sprintf("Deleted %d delay policies: ", len(responses)) + strings.Join(responses, " ")

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
func beginSessionHandler(w http.ResponseWriter, r *http.Request) {
	queryArgs, _ := url.ParseQuery(r.URL.RawQuery)
	// retrieving details and validating request
	if scenario, ok := queryArgs["scenario"]; ok {
		if session, ok := queryArgs["session"]; ok {
			if mode, ok := queryArgs["mode"]; ok {
				// Create scenario. This can result in 422 (duplicate error) and this is
				// fine, since we must only ensure that it exists.
				client := &Client{&http.Client{}}
				_, err := client.createScenario(scenario[0])
				if err != nil {
					http.Error(w, err.Error(), 500)
				}
				// Begin session
				response, err := client.beginSession(session[0], scenario[0], mode[0])
				if err != nil {
					http.Error(w, err.Error(), 500)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)
			} else {
				http.Error(w, "Bad request, missing session mode key.", 400)
			}
		} else {
			http.Error(w, "Bad request, missing session name.", 400)
		}
	} else {
		http.Error(w, "Bad request, missing scenario name.", 400)
	}
}

func endSessionsHandler(w http.ResponseWriter, r *http.Request) {
	scenario, ok := r.URL.Query()["scenario"]
	if ok {
		fmt.Println("got:", r.URL.Query())
		// expecting one param - scenario
		client := &Client{&http.Client{}}
		response, err := client.endSessions(scenario[0])
		// checking whether we got good response
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		// setting resposne header
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		http.Error(w, "Scenario name not provided.", 400)
	}
}

func getScenariosHandler(w http.ResponseWriter, r *http.Request) {
	client := &Client{&http.Client{}}
	response, err := client.getScenarios()
	// checking whether we got good response
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	// setting resposne header
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}
