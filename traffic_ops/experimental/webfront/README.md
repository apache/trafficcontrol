<!--
    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.
-->

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
`go run webfront.go webfront.config my-secret`
`secret` is used for jwt signing


* To login:
`curl -X POST --insecure -Lkvs --header "Content-Type:application/json" https://localhost:9004/login -d'{"username":"foo", "password":"bar"}'`
   
* To use a token:
`curl --insecure -H'Authorization: Bearer <token>' -Lkvs  https://localhost:8080/ds/`

