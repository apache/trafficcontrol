
A simple authentication server written in go that authenticates user against the `tm_user` table and returns a jwt access token representing the user, incl. its API access capabilities, derived from the user's role.

Note that the authentication server is designed to work in conjunction with the "webfront" server, that acts as an API GW. Once you obtain an access token from the auth service you can use it with "webfront" to authenticate your API calls. See [webfront documentation](../webfront/README.md)

#### Legacy TO support

Legacy TO authorization code requires any API call to pass a mojolicious access token in its access control headers.
Untill this code is deprecated, the Auth server and the API GW handle legacy authorization in the following way:
Upon every sucessful login the auth server performs additional login against legacy TO (mojolicious app) and recieves a legacy TO authentication token.
This token is passed back on the user's JWT, and used by the API GW to set access control headers upon any consecutive API calls.

**Before you begin**

You will need to generate a server certificate for ssl connections against webfront. In the project directory, run
~~~~
openssl req -x509 -sha256 -nodes -days 3650 -newkey rsa:2048 -keyout server.key -out server.crt
~~~~

**Run the server**

	`go run auth.go auth.config my-secret`

	`my-secret` is used for jwt signing

**Perform a login call (to get a token)**

	`curl --insecure -X POST -Lkvs --header "Content-Type:application/json" https://localhost:9004/login -d'{"username":"username", "password":"password"}'`

See [webfront documentation](../webfront/README.md) for using this token in your API calls against the webfront server. 

Note that webfront forwad login calls to the auth server. In real-world scanarios login calls are done against webfront (API GW) and not directly against the auth server. Login calls via webfront do not require a token