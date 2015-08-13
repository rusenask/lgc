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
	StuboHost string
	StuboPort string
}

// StublistHandler gets stubs, e.g.: stubo/api/get/stublist?scenario=first
func StublistHandler(w http.ResponseWriter, r *http.Request) {
	scenario, ok := r.URL.Query()["scenario"]
	if ok {
		fmt.Println("got:", r.URL.Query())
		// expecting one param - scenario
		response := getStubList(scenario[0])
		js, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		fmt.Println("Scenario name not provided.")
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
	mux.Get("/stubo/api/get/stublist", http.HandlerFunc(StublistHandler))
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3000")
}
