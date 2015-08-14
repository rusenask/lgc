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


###Requirements
go get github.com/go-zoo/bone - lightweight and lightning fast HTTP Multiplexer for Golang.

go get github.com/codegangsta/negroni - Negroni is an idiomatic approach to web middleware in Go. It is tiny, non-intrusive, and encourages use of net/http Handlers.

### Current legacy API translations

* "/stubo/api/get/status", "GetStatusHandler") - not implemented
* ("/stubo/api/get/response", "GetResponseHandler") - not implemented
* ("/stubo/api/get/response/.*", "GetResponseHandler") - not implemented
* ("/stubo/api/begin/session", "BeginSessionHandler") - __implemented__
* ("/stubo/api/end/session", "EndSessionHandler") - compatibility issues*
* ("/stubo/api/end/sessions", "EndSessionsHandler") - __implemented__
* ("/stubo/api/put/stub", "PutStubHandler") - not implemented
* ("/stubo/api/delete/stubs", "DeleteStubsHandler") - not implemented
* ("/stubo/api/get/stubcount", "GetStubCountHandler") - not implemented
* ("/stubo/api/get/stublist", "GetStubListHandler") - __implemented__
* ("/stubo/api/get/scenarios", "GetScenariosHandler") - not implemented
* ("/stubo/api/put/scenarios/(?P<scenario_name>[^\/]+)", "PutScenarioHandler") - not implemented
* ("/stubo/api/get/export", "GetStubExportHandler") - not implemented
* ("/stubo/api/get/stats", "GetStatsHandler") - not implemented
* ("/stubo/api/put/delay_policy", "PutDelayPolicyHandler") - not implemented
* ("/stubo/api/get/delay_policy", "GetDelayPolicyHandler") - __implemented__
* ("/stubo/api/delete/delay_policy", "DeleteDelayPolicyHandler") - not implemented
* ("/stubo/api/get/version", "GetVersionHandler") - not implemented
* ("/stubo/api/put/module", "PutModuleHandler") - not implemented
* ("/stubo/api/delete/module", "DeleteModuleHandler") - not implemented
* ("/stubo/api/delete/modules", "DeleteModulesHandler") - not implemented
* ("/stubo/api/get/modulelist", "GetModuleListHandler") - not implemented
* ("/stubo/api/put/bookmark", "PutBookmarkHandler") - not implemented, this functionality will probably be removed
* ("/stubo/api/get/bookmarks", "GetBookmarksHandler") - not implemented, this functionality will probably be removed
* ("/stubo/api/put/setting", "PutSettingHandler") - not implemented
* ("/stubo/api/get/setting", "GetSettingHandler") - not implemented
* ("/stubo/api/jump/bookmark", "JumpBookmarkHandler") - not implemented, this functionality will probably be removed
* ("/stubo/api/delete/bookmark", "DeleteBookmarkHandler") - not implemented, this functionality will probably be removed
* ("/stubo/api/import/bookmarks", "ImportBookmarksHandler") - not implemented, this functionality will probably be removed
* ("/stubo/default/execCmds", "StuboCommandHandler") - not implemented
* ("/stubo/api/exec/cmds", "StuboCommandHandler") - not implemented

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
