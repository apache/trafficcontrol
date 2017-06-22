A reverse proxy written in go that can front any number of microservices. It uses a rules file to map from requested host/path to microservice host/port/path.  

The API GW forwarding logic is as follow:
Find host to forward the request: Prefix match on the request path against a list of forwarding rules. The matched forwarding rule defines the target's host, and the target's authorization rules. 
Authorization: Regex match on the request path against a list of authorization rules. The matched rule defines the required capabilities to perform the HTTP method on the route. These capabilities are compared against the user's capabilities in the user's JWT

Example forward rules file:

```json
[
    { "host": "localhost", "path": "/login",               "forward": "localhost:9004", "scheme": "https", "auth": false },
    { "host": "localhost", "path": "/api/1.2/innovation/", "forward": "localhost:8004", "scheme": "http",  "auth": false },
    { "host": "localhost", "path": "/api/1.2/",            "forward": "localhost:3000", "scheme": "http",  "auth": true, "routes-file": "traffic-ops-routes.json" },
    { "host": "localhost", "path": "/internal/api/1.2/",   "forward": "localhost:3000", "scheme": "http",  "auth": true, "routes-file": "internal-routes.json" }
]
```

Example authorised routes file:
```json
[
    { "match": "/cdns/health",                        "auth": { "GET":  ["cdn-health-read"] }},
    { "match": "/cdns/capacity",                      "auth": { "GET":  ["cdn-health-read"] }},
    { "match": "/cdns/usage/overview",                "auth": { "GET":  ["cdn-stats-read"] }},
    { "match": "/cdns/name/dnsseckeys/generate",      "auth": { "GET":  ["cdn-security-keys-read"] }},
    { "match": "/cdns/name/[^\/]+/?",                 "auth": { "GET":  ["cdn-read"] }},
    { "match": "/cdns/name/[^\/]+/sslkeys",           "auth": { "GET":  ["cdn-security-keys-read"] }},
    { "match": "/cdns/name/[^\/]+/dnsseckeys",        "auth": { "GET":  ["cdn-security-keys-read"] }},
    { "match": "/cdns/name/[^\/]+/dnsseckeys/delete", "auth": { "GET":  ["cdn-security-keys-write"] }},
    { "match": "/cdns/[^\/]+/queue_update",           "auth": { "POST": ["queue-updates-write"] }},
    { "match": "/cdns/[^\/]+/snapshot",               "auth": { "PUT":  ["cdn-config-snapshot-write"] }},
    { "match": "/cdns/[^\/]+/health",                 "auth": { "GET":  ["cdn-health-read"] }},
    { "match": "/cdns/[^\/]+/?",                      "auth": { "GET":  ["cdn-read"], "PUT":  ["cdn-write"], "PATCH": ["cdn-write"], "DELETE": ["cdn-write"] }},
    { "match": "/cdns",                               "auth": { "GET":  ["cdn-read"], "POST": ["cdn-write"] }}
]
```

No restart is needed to re-read the forwarding rule file and apply; within 60 seconds of a change in the file, it will pick up the new mappings.
However, authorized routes files are not re-read. Touch the forwarding rule file to trigger an update.

**Legacy TO support**

Currently, the Mojo app reqires a valid Mojo token. As long as the Mojo code use a Mojo token for authorization, the Auth server and the API GW handle legacy authorization in the following way

* Upon every sucessful login, the auth server performs additional login against the Mojo app and recieves a Mojo token
* The Mojo token is added as a claim to the user's JWT
* Upon successive API calls, the API GW pulls the claim from the JWT and set a "mojolicious" cookie on the request

In addition, if a request contains a "mojolicious" cookie instead of an authentication bearer token, the API GW bypass JWT authentication. 
This is to support legacy code that access TO API without logging in via the new auth server.

**Before you begin**

You will need to generate a server certificate for ssl connections against webfront. In the project directory, run
~~~~
openssl req -x509 -sha256 -nodes -days 3650 -newkey rsa:2048 -keyout server.key -out server.crt
~~~~


**Run the server**

    `go run webfront.go webfront.config my-secret`

    `my-secret` is used for jwt signing


**Perform a login call (to get a token)**

    `curl -X POST --insecure -Lkvs --header "Content-Type:application/json" https://localhost:9004/login -d'{"username":"foo", "password":"bar"}'`
   
**Perform an API call (using the token
)**

    `curl --insecure -H'Authorization: Bearer <token>' -Lkvs  https://localhost:8080/ds/`

#### Legacy TO support

    1 Since Mojo app requires a valid Mojo token, the auth server performs a legacy Traffic Ops login upon ever 
