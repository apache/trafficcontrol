
A simple authentication server written in go that authenticates user agains the `tm_user` table and returns a jwt representing the user, incl. its API access capabilities, derived from the user's role.

* To run:
	`go run auth.go auth.conf my-secret`

* To login:
	`curl --insecure -X POST -Lkvs --header "Content-Type:application/json" https://localhost:9004/login -d'{"username":"foo", "password":"bar"}'`