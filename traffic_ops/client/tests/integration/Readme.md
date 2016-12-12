# Traffic Ops Client Integration Tests

The Traffic Ops Client Integration tests are used to validate the clients responses against those from the Traffic Ops API.  In order to run the tests you will need a Traffic Ops instance with at least one of each of the following:  CDN, Delivery Service, Type, Cachegroup, User and Server.

## Running the Integration Tests
The integration tests are run using `go test`, however, there are some flags that need to be provided in order for the tests to work.  The flags are:

* toURL - The URL to Traffic Ops.  Default is "http://localhost:3000".
* toUser - The Traffic Ops user to use.  Default is "admin".
* toPass - The password of the user provided.  Deafault is "password".

Example command to run the tests: `go test -v -toUrl=https://to.kabletown.net -toUser=myUser -toPass=myPass`

*It can take serveral minutes for the integration tests to complete, so using the `-v` flag is recommended to see progress.*
