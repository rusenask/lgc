package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"github.com/meatballhat/negroni-logrus"
)

// Configuration to hold stubo details
type Configuration struct {
	StuboProtocol string
	StuboHost     string
	StuboPort     string
	Environment   string
	Debug         bool
}

// StuboConfig stores target Stubo instance details (protocol, hostname, port, etc..)
var StuboConfig Configuration

// StuboURI stores URI (e.g. "http://localhost:8001")
var StuboURI string

func main() {
	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)

	// getting configuration
	file, err := os.Open("conf.json")
	if err != nil {
		log.Panic("Failed to open configuration file, quiting server.")
	}
	decoder := json.NewDecoder(file)
	StuboConfig = Configuration{}
	err = decoder.Decode(&StuboConfig)
	if err != nil {
		log.WithFields(log.Fields{"Error": err.Error()}).Panic("Failed to read configuration")
	}
	// configuring logger based on environment settings
	if StuboConfig.Environment == "production" {
		// Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.JSONFormatter{})
		// TODO: also, write to file, probably log file path should also be configurable
	} else {
		// The TextFormatter is default
		log.SetFormatter(&log.TextFormatter{})
	}

	if StuboConfig.Debug == true {
		log.SetLevel(log.DebugLevel)
		log.Info("Starting server with debug mode initiated...")
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

	client := &Client{&http.Client{}}
	mux := getRouter(HandlerHTTPClient{*client})

	n := negroni.Classic()
	n.Use(negronilogrus.NewMiddleware())
	n.UseHandler(mux)
	n.Run(*port)
}

func getRouter(h HandlerHTTPClient) *bone.Mux {
	mux := bone.New()
	mux.Post("/stubo/api/put/stub", http.HandlerFunc(h.putStubHandler))
	mux.Get("/stubo/api/get/stublist", http.HandlerFunc(h.stublistHandler))
	mux.Get("/stubo/api/delete/stubs", http.HandlerFunc(h.deleteStubsHandler))
	mux.Get("/stubo/api/put/delay_policy", http.HandlerFunc(h.putDelayPolicyHandler))
	mux.Get("/stubo/api/get/delay_policy", http.HandlerFunc(h.getDelayPolicyHandler))
	mux.Get("/stubo/api/delete/delay_policy", http.HandlerFunc(h.deleteDelayPolicyHandler))
	mux.Get("/stubo/api/begin/session", http.HandlerFunc(h.beginSessionHandler))
	mux.Get("/stubo/api/end/sessions", http.HandlerFunc(h.endSessionsHandler))
	mux.Get("/stubo/api/get/scenarios", http.HandlerFunc(h.getScenariosHandler))
	return mux
}
