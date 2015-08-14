# Stubo proxy prototype

Proxy to work with stubo API v2 (which is still under development). After setting
it up - it will translate all legacy API calls to new format REST API calls.

Example:
LGC proxy running on port 3000 and Stubo with API v2 running on port 8001.

Client calls:
http://localhost:3000/stubo/api/get/delay_policy?name=delay_1
This then gets translated into:
http://localhost:8001/stubo/api/v2/delay-policy/objects/delay_1
LGC gets response and sends it back to the client.


Requirements
go get github.com/go-zoo/bone - lightweight and lightning fast HTTP Multiplexer for Golang.

go get github.com/codegangsta/negroni - Negroni is an idiomatic approach to web middleware in Go. It is tiny, non-intrusive, and encourages use of net/http Handlers.


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
