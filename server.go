package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

// Configuration to hold stubo details
type Configuration struct {
	StuboHost    string
	StuboPort    string
	LgcProxyPort string
}

// stublistHandler gets stubs, e.g.: stubo/api/get/stublist?scenario=first
func stublistHandler(w http.ResponseWriter, r *http.Request) {
	scenario, ok := r.URL.Query()["scenario"]
	if ok {
		fmt.Println("got:", r.URL.Query())
		// expecting one param - scenario
		response, err := getStubList(scenario[0])
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
				_, err := createScenario(scenario[0])
				if err != nil {
					http.Error(w, err.Error(), 500)
				}
				fmt.Println(scenario)
				fmt.Println(session)
				fmt.Println(mode)
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

func main() {
	// getting configuration
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	mux := bone.New()
	mux.Get("/stubo/api/get/stublist", http.HandlerFunc(stublistHandler))
	mux.Get("/stubo/api/get/delay_policy", http.HandlerFunc(getDelayPolicyHandler))
	mux.Get("/stubo/api/begin/session", http.HandlerFunc(beginSessionHandler))
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3000")
}
