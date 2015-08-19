package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
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

func main() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)

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

	log.WithFields(log.Fields{
		"StuboHost": StuboConfig.StuboHost,
		"StuboPort": StuboConfig.StuboPort,
		"StuboURI":  StuboURI,
		"ProxyPort": port,
	}).Info("LGC is starting")

	mux := bone.New()
	mux.Get("/stubo/api/get/stublist", http.HandlerFunc(stublistHandler))
	mux.Get("/stubo/api/delete/stubs", http.HandlerFunc(deleteStubsHandler))
	mux.Get("/stubo/api/get/delay_policy", http.HandlerFunc(getDelayPolicyHandler))
	mux.Get("/stubo/api/delete/delay_policy", http.HandlerFunc(deleteDelayPolicyHandler))
	mux.Get("/stubo/api/begin/session", http.HandlerFunc(beginSessionHandler))
	mux.Get("/stubo/api/end/sessions", http.HandlerFunc(endSessionsHandler))
	mux.Get("/stubo/api/get/scenarios", http.HandlerFunc(getScenariosHandler))
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(*port)
}
