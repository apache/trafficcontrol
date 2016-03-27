This application - webfront is a reverse proxy written in go that can front any number of microservices. It uses a rules file to map from requested host/path to microservice host/port/path.  Example rule file:

   
    [
		{"Host": "local.com", "Path" : "/8001", "Forward": "localhost:8001"},
		{"Host": "local.com", "Path" : "/8002", "Forward": "localhost:8002"},
		{"Host": "local.com", "Path" : "/8003", "Forward": "localhost:8003"},
		{"Host": "local.com", "Path" : "/8004", "Forward": "localhost:8004"},
		{"Host": "local.com", "Path" : "/8005", "Forward": "localhost:8005"},
		{"Host": "local.com", "Path" : "/8006", "Forward": "localhost:8006"},
		{"Host": "local.com", "Path" : "/8007", "Forward": "localhost:8007"}
	]


No restart is needed to re-read the rule file and apply; within 60 seconds of a change in the file, it will pick up the new mappings.

To run

	go run webfront.go -rules=rules.json -https=:9000 -https_cert=server.pem -https_key=server.key 

(or compile a binary, and run that)

To get a token:

	curl --insecure -Lkvs --header "Content-Type:application/json" -XPOST https://localhost:9000/login -d'{"username":"jvd", "password":"tootoo"}'
   
in my case that returned: 

	{"Token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQYXNzd29yZCI6InRvb3RvbyIsIlVzZXIiOiIiLCJleHAiOjE0NTg5NDg2MTl9.quCwZ5vghVBucxMxQ4fSfD84yw_yPEp9qLGGQNcHNUk"}``
   
 Example:
  
	[jvd@laika webfront (master *=)]$ curl --insecure -Lkvs --header "Content-Type:application/json" -XPOST https://localhost:9000/login -d'{"username":"jvd", "password":"tootoo"}'
	*   Trying ::1...
	* Connected to localhost (::1) port 9000 (#0)
	* TLS 1.2 connection using TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
	* Server certificate: CU
	> POST /login HTTP/1.1
	> Host: localhost:9000
	> User-Agent: curl/7.43.0
	> Accept: */*
	> Content-Type:application/json
	> Content-Length: 39
	>
	* upload completely sent off: 39 out of 39 bytes
	< HTTP/1.1 200 OK
	< Content-Type: application/json
	< Date: Thu, 24 Mar 2016 23:30:19 GMT
	< Content-Length: 157
	<
	* Connection #0 to host localhost left intact
	{"Token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQYXNzd29yZCI6InRvb3RvbyIsIlVzZXIiOiIiLCJleHAiOjE0NTg5NDg2MTl9.quCwZ5vghVBucxMxQ4fSfD84yw_yPEp9qLGGQNcHNUk"}[jvd@laika webfront (master *=)]$
 
 * To use a token: 

	curl --insecure -H'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQYXNzd29yZCI6InRvb3RvbyIsIlVzZXIiOiIiLCJleHAiOjE0NTg5NDg2MTl9.quCwZ5vghVBucxMxQ4fSfD84yw_yPEp9qLGGQNcHNUk' -Lkvs  https://localhost:9000/8003/r

Example:

   
    [jvd@laika webfront (master *%=)]$ curl --insecure -H'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQYXNzd29yZCI6InRvb3RvbyIsIlVzZXIiOiIiLCJleHAiOjE0NTg5NDg2MTl9.quCwZ5vghVBucxMxQ4fSfD84yw_yPEp9qLGGQNcHNUk' -Lkvs  https://localhost:9000/8003/r
	*   Trying ::1...
	* Connected to localhost (::1) port 9000 (#0)
	* TLS 1.2 connection using TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
	* Server certificate: CU
	> GET /8003/r HTTP/1.1
	> Host: localhost:9000
	> User-Agent: curl/7.43.0
	> Accept: */*
	> Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQYXNzd29yZCI6InRvb3RvbyIsIlVzZXIiOiIiLCJleHAiOjE0NTg5NDg2MTl9.quCwZ5vghVBucxMxQ4fSfD84yw_yPEp9qLGGQNcHNUk
	>
	< HTTP/1.1 200 OK
	< Content-Length: 24
	< Content-Type: text/plain; charset=utf-8
	< Date: Thu, 24 Mar 2016 23:34:08 GMT
	<
	Hitting 8003 with /boo1
	* Connection #0 to host localhost left intact
	[jvd@laika webfront (master *%=)]$

	[jvd@laika webfront (master=)]$ curl --insecure -H'Authorization: Bearer FAKETOKEN' -Lkvs  https://localhost:9000/8003/r *   Trying ::1...
	* Connected to localhost (::1) port 9000 (#0)
	* TLS 1.2 connection using TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
	* Server certificate: CU
	> GET /8003/r HTTP/1.1
	> Host: localhost:9000
	> User-Agent: curl/7.43.0
	> Accept: */*
	> Authorization: Bearer FAKETOKEN
	>
	< HTTP/1.1 403 Forbidden
	< Date: Thu, 24 Mar 2016 23:43:11 GMT
	< Content-Length: 0
	< Content-Type: text/plain; charset=utf-8
	<
	* Connection #0 to host localhost left intact
	[jvd@laika webfront (master=)]$
  