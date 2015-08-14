package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		response := getStubList(scenario[0])
		// setting resposne header
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		fmt.Println("Scenario name not provided.")
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
		response := getDelayPolicy(name[0])
		// setting resposne header
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		// name is not provided, getting all delay policies
		response := getAllDelayPolicies()
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func beginSessionHandler(w http.ResponseWriter, r *http.Request) {
	queryArgs, _ := url.ParseQuery(r.URL.RawQuery)
	fmt.Println(queryArgs)
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