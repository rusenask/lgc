package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

// Configuration to hold stubo details
type Configuration struct {
	StuboProtocol string
	StuboHost     string
	StuboPort     string
}

// StuboConfig stores target Stubo instance details (protocol, hostname, port, etc..)
var StuboConfig Configuration

// StuboURI stores URI (e.g. "http://localhost:8001")
var StuboURI string

// stublistHandler gets stubs, e.g.: stubo/api/get/stublist?scenario=first
func stublistHandler(w http.ResponseWriter, r *http.Request) {
	scenario, ok := r.URL.Query()["scenario"]
	if ok {
		fmt.Println("got:", r.URL.Query())
		// expecting one param - scenario
		response, err := getScenarioStubs(scenario[0])
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

// getDelayPolicyHandler - returns delay policy information, list all if
// name is not provided, e.g.: stubo/api/get/delay_policy?name=slow
func getDelayPolicyHandler(w http.ResponseWriter, r *http.Request) {
	name, ok := r.URL.Query()["name"]
	if ok {
		// name provided so looking for specific delay
		fmt.Println("got:", r.URL.Query())
		// expecting one param - scenario
		response, err := getDelayPolicy(name[0])
		// checking whether we got good response
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		// setting resposne header
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		// name is not provided, getting all delay policies
		response, err := getAllDelayPolicies()
		// checking whether we got good response
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
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
				_, err := createScenario(scenario[0])
				if err != nil {
					http.Error(w, err.Error(), 500)
				}
				// Begin session
				response, err := beginSession(session[0], scenario[0], mode[0])
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
		response, err := endSessions(scenario[0])
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
	response, err := getScenarios()
	// checking whether we got good response
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	// setting resposne header
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

func main() {
	// getting configuration
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	StuboConfig = Configuration{}
	err := decoder.Decode(&StuboConfig)
	if err != nil {
		fmt.Println("error:", err)
	}
	// looking for option args when starting App
	// like ./lgc -port=":3000" would start on port 3000
	var port = flag.String("port", ":3000", "Server port")
	flag.Parse() // parse the flag

	// assign StuboURI
	StuboURI = StuboConfig.StuboProtocol + "://" + StuboConfig.StuboHost + ":" + StuboConfig.StuboPort

	mux := bone.New()
	mux.Get("/stubo/api/get/stublist", http.HandlerFunc(stublistHandler))
	mux.Get("/stubo/api/get/delay_policy", http.HandlerFunc(getDelayPolicyHandler))
	mux.Get("/stubo/api/begin/session", http.HandlerFunc(beginSessionHandler))
	mux.Get("/stubo/api/end/sessions", http.HandlerFunc(endSessionsHandler))
	mux.Get("/stubo/api/get/scenarios", http.HandlerFunc(getScenariosHandler))
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(*port)
}
