package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
)

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

// trace returns name of the current function
func trace() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
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
		_, _, err := c.deleteDelayPolicy(dp.Name)

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
