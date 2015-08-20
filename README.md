# Stubo proxy prototype

[![Build Status](https://travis-ci.org/rusenask/lgc.svg?branch=master)](https://travis-ci.org/rusenask/lgc)

Proxy to work with stubo API v2 (which is still under development). After setting
it up - it will translate all legacy API calls to new format REST API calls.

### Example
LGC proxy running on port 3000 and Stubo with API v2 running on port 8001.

Client calls:
 * http://localhost:3000/stubo/api/get/delay_policy?name=delay_1

This then gets translated into:
* http://localhost:8001/stubo/api/v2/delay-policy/objects/delay_1

LGC gets response and sends it back to the client.

### A little more complex example
Client calls:
* http://localhost:3000/stubo/api/begin/session?scenario=scenario_x&session=session_x&mode=record

Due to the fact that current v2 API requires user to create a scenario which then could hold session,
this API call results in two calls to stubo:
* __URL__:          http://localhost:8001/stubo/api/v2/scenarios
* __Method__:       PUT
* __Request body__:  {"scenario": "scenario_x"}

Then, after scenario is created, a second call to begin session is made:
* __URL__:          http://localhost:8001/stubo/api/v2/scenarios/objects/scenario_x/action
* __Method__:       POST
* __Request body__:  {"begin": null, "session": "session_x",  "mode": "record"}


### Requirements
go get github.com/go-zoo/bone - lightweight and lightning fast HTTP Multiplexer for Golang.

go get github.com/codegangsta/negroni - Negroni is an idiomatic approach to web middleware in Go. It is tiny, non-intrusive, and encourages use of net/http Handlers.

go get github.com/meatballhat/negroni-logrus - Negroni/Logrus middleware for merging
negroni logs with application logging. This provides additional data such as status codes,
time taken for response and latency
### Configuration

Edit conf.json.example with your stubo instance details:
{
  "StuboHost": "localhost", // your stubo hostname
  "StuboPort": "8001",  // your stubo port
  "StuboProtocol": "http", // protocol (should probably be http anyway so leave it)
  "Environment": "production",
  "debug": true
}
Rename conf.json.example to conf.json

Default LGC proxy port is 3000. You are expected to change it during server startup:
./lgc -port=":8001"
Would change it to this port. Remember to change your original stubo instance port before setting it to 8001.
Environment variable sets some logging defaults (such as format). Although you can
modify logging formatter yourself in server.go file.

Debug - when enabled outputs more information about request forming before dispatching them to stubo.


### Current legacy API translations

* exec/cmds - not present in API v2
* get/version - not present in API v2
* get/status - not present in API v2
* begin/session - __implemented__ (both playback and record modes)
* end/session - not compatible with current API v2, use "end/sessions" call
* end/sessions - __implemented__
* put/scenarios - not implemented
* get/scenarios - __implemented__
* put/stub - not present in API v2
* get/stublist - __implemented__
* put/delay_policy - not implemented
* get/delay_policy:
    + name provided - __implemented__
    + name not provided (should list all delay policies) - __implemented__
* delete/delay_policy:
    + name provided - __implemented__
    + name not provided (should delete all delay policies) - __implemented__
* get/response - not present in API v2
* delete/stubs:
    + host provided - __implemented__
    + force provided - __implemented__
* get/export - not implemented
* get/stubcount - not implemented
* put/module - not present in API v2
* get/modulelist - not present in API v2
* delete/module - not present in API v2
* delete/modules - not present in API v2
* Set Tracking Level - not present in API v2
* Blacklist a host URL - not present in API v2
* Delete Bookmark - not present in API v2
* List Bookmarks - not present in API v2
* get/stats - not present in API v2

### Compatibility
API compatibility issues:
* Need to find a way to end a specific version. Current API v2 needs scenario name to end session:
  /stubo/api/v2/scenarios/objects/{scenario_name}/action with body:
  ```javascript
  { “end”: null,
   “session”: “session_name” }
  ```
  However, the legacy API needs only session name (skipping scenario):
  stubo/api/end/session?session=session_name
  __Current solution__
  Using legacy API call to end _all_ sessions:
  stubo/api/end/sessions?scenario=scenario_name  
