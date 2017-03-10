A reverse proxy written in go that can front any number of microservices. It uses a rules file to map from requested host/path to microservice host/port/path.  Example rule file:

```json
	[
	  {
	    "host": "domain.com",
	    "path": "/login",
	    "forward": "localhost:9004",
	    "secure": false
	  },
	  {
	    "host": "domain.com",
	    "path": "/ds/",
	    "forward": "localhost:8081",
	    "secure": true,
	    "capabilities": {
	        "GET": "read-ds",
	        "POST": "write-ds",
	        "PUT": "write-ds",
	        "PATCH": "write-ds"
	    }
	  },
	  {
	    "host": "domain.com",
	    "path": "/cachegroups/",
	    "forward": "localhost:8082",
	    "secure": true,
	    "capabilities": {
	        "GET": "read-cg",
	        "POST": "write-cg",
	        "PUT": "write-cg",
	        "PATCH": "write-cg"
	    }
	  }
	]
```

Note: Access "capabilities" are set in the rule file as a workaround, until TO DB tables are ready.

No restart is needed to re-read the rule file and apply; within 60 seconds of a change in the file, it will pick up the new mappings.

* To run:
	`go run webfront.go webfront.conf my-secret`

* To login:
	`curl -X POST --insecure -Lkvs --header "Content-Type:application/json" https://localhost:9004/login -d'{"username":"foo", "password":"bar"}'`
   
* To use a token:
	`curl --insecure -H'Authorization: Bearer <token>' -Lkvs  https://localhost:8080/ds/`

