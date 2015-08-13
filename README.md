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
