# Stub-O-Matic proxy

[![Build Status](https://travis-ci.org/rusenask/lgc.svg?branch=master)](https://travis-ci.org/rusenask/lgc)

Proxy to work with Stub-O-Matic API v2 (which is still under development). After setting
it up - it will translate all legacy API calls to new format REST API calls.

### Example
LGC proxy running on port 3000 and Stubo with API v2 running on port 8001.

Client calls:
 * http://localhost:3000/stubo/api/get/delay_policy?name=delay_1

This then gets translated into:
* http://localhost:8001/stubo/api/v2/delay-policy/objects/delay_1

LGC gets response (with response code) and sends it back to the client.

### A little more complex example
Client calls:
* http://localhost:3000/stubo/api/begin/session?scenario=scenario_x&session=session_x&mode=record

Due to the fact that current v2 API requires user to create a scenario which then could hold session,
this API call results in two calls to stubo:
* __URL__:           http://localhost:8001/stubo/api/v2/scenarios
* __Method__:        PUT
* __Request body__:  {"scenario": "scenario_x"}

Then, after scenario is created, a second call to begin session is made:
* __URL__:           http://localhost:8001/stubo/api/v2/scenarios/objects/scenario_x/action
* __Method__:        POST
* __Request body__:  {"begin": null, "session": "session_x",  "mode": "record"}

### Put Stub example
Client calls (POST method):
* http://localhost:3000/stubo/api/put/stub?session=sc1:session_name&some=yes&stateful=true&additionalparam=true

It is expected that session and stateful parameters should be transformed into headers to comply with API v2 standard.
However, other parameters must remain in the URL arguments list and recorded by Stubo. Proxy transforms this into:
* __URL__:             http://localhost:8001/stubo/api/v2/scenarios/objects/sc1/stubs?some=yes&additionalparam=true
* __Method__:          PUT
* __Request body__:    remains the same, bytes get just passed to new request
* __Request headers__: session: session_name
                       stateful: true  

Stubo then sends back response and proxy passes those bytes back to the client:
```javascript
{  "version": "0.6.6",
   "data": {
        "message": {
            "status": "updated", "msg": "Updated with stateful response",
            "key": "55dc6cc1938fbef2e62d875c"}
          }
}
```

### Requirements

LGC uses vendor/ to manage dependencies. Enable vendor experiment:

GO15VENDOREXPERIMENT set to 1

"export GO15VENDOREXPERIMENT=1"

Requirements are managed by Glide - https://github.com/Masterminds/glide
You only need it if you want to add new/remove/update packages.

To install Glide on your PC, run:

brew install glide

#### Usage

> * glide create                            # Start a new workspaces
* open glide.yaml                         # and edit away!
* glide get github.com/Masterminds/cookoo # Get a package and add to glide.yaml
* glide install                           # Install packages and dependencies
#### work, work, work
* go build                                # Go tools work normally
* glide up                                # Update to newest versions of the package


github.com/go-zoo/bone - lightweight and lightning fast HTTP Multiplexer for Golang.

github.com/codegangsta/negroni - Negroni is an idiomatic approach to web middleware in Go. It is tiny, non-intrusive, and encourages use of net/http Handlers.

github.com/meatballhat/negroni-logrus - Negroni/Logrus middleware for merging
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

#### Using Docker during development

* Build container:

__docker build -t lgc .__


* Run your image:

__docker run --publish 3000:3000 --add-host="localhost:192.168.59.3" --name lgc --rm lgc__

Used flags:
* --publish 3000:3000 publishes application port, by default LGC uses port 3000 but you can map any port Here
* --add-host="localhost:192.168.59.3" - useful when running stubo on the same machine (localhost)
> During boot2docker initialization, boot2docker sets up two network adapters, one using NAT and the other a "host-only" adapter. The IP of the host on that host-only network is determined by virtualbox. You can inspect it by running ifconfig on your host and looking at the vboxnet0 interface. By default it's 192.168.59.3, but you can change it in virtualbox's global settings at Preferences -> Network -> Host-only Networks. I am able to ping that IP both from within the boot2docker VM and from within a running docker container.

* --name lgc - assigns some expected name to your container
* --rm lgc - removes container



### Current legacy API translations

* exec/cmds - not present in API v2
* get/version - not present in API v2
* get/status - not present in API v2
* begin/session - __implemented__ (both playback and record modes)
* end/session - not compatible with current API v2, use "end/sessions" call
* end/sessions - __implemented__
* put/scenarios (rename existing scenario) - not present in API v2
* get/scenarios - __implemented__
* put/stub:
    + basic insertion with scenario_name:session_name - __implemented__
    + ext_module = external module name without .py extenstion (optional) __implemented__
    + delay_policy =  delay policy name (optional) __implemented__
    + stateful = treat duplicate stubs as stateful otherwise ignore duplicates if stateful=false (default true, optional) __implemented__
    + tracking_level: full or normal (optional, overrides host or global setting) __implemented__
    + any user args will be made available to the matcher & response templates and any user exit code __implemented__
* get/stublist - __implemented__
* put/delay_policy - __implemented__
* get/delay_policy:
    + name provided - __implemented__
    + name not provided (should list all delay policies) - __implemented__
* delete/delay_policy:
    + name provided - __implemented__
    + name not provided (should delete all delay policies) - __implemented__
* get/response - __implemented__
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

### Logging

LGC uses logrus logging middleware. If "debug" mode in configuration is set to true -
debug level logs are being written as well. You can set different logging levels in
server.go


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
