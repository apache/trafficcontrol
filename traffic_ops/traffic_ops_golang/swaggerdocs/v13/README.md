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

#### `./swaggerdocs` 
This directory contains the Go structs that glue together the Swagger 2.0 metadata that will generate the Traffic Ops API documentation using [go-swagger](https://github.com/go-swagger/go-swagger) meta tags.  The Traffic Ops API documentation is maintained by modifying the Go files in this directory that point to the **incubator-trafficcontrol/lib/go-tc/*.go** structs that render the Traffic Ops Go Proxy API's.


### Setup

* Install Docker for your platform:
[https://docs.docker.com/install](https://docs.docker.com/install)

* Install Docker Compose for your platform:
[https://docs.docker.com/compose/install](https://docs.docker.com/compose/install)

### Running the web services

The `docker-compose.yml` will start 2 services a custom http service for hosting the `swaggerspec/swagger.json` and the Swagger UI.  

To start the Swagger UI services just run:

```$ docker-compose up```

Once started navigate your browser to [http://localhost:8080](http://localhost:8080)

### Generating your Swagger Spec File

The **gen_swaggerspec.sh** script will scan all the Go files in the swaggerdocs directory and extract out all of the swagger meta tags that are embedded as comments.  The output of the **gen_swaggerspec.sh** script will be the **swaggerspec/swagger.json** spec file. 

While the Docker services are running, just re-run **gen_swaggerspec.sh** and hit refresh on the page to see the Swagger doc updates in real time.