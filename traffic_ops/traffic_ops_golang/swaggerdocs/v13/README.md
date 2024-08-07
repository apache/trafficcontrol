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

#### `./swaggerdocs` overview

This directory contains the Go structs that glue together the Swagger 2.0 metadata that will generate the Traffic Ops API documentation using the [go-swagger](https://github.com/go-swagger/go-swagger) meta tags.  The Traffic Ops API documentation is maintained by modifying the Go files in this directory and the Go structs that they reference from here **trafficcontrol/lib/go-tc/*.go**.  These combination of these two areas of .go files will produce Swagger documentation for the Traffic Ops Go API's.

### Setup

* Install Docker for your platform:
[https://docs.docker.com/install](https://docs.docker.com/install)

* Install Docker Compose for your platform:
[https://docs.docker.com/compose/install](https://docs.docker.com/compose/install)

### Generating your Swagger Spec File

The **gen_swaggerspec.sh** script will scan all the Go files in the swaggerdocs directory and extract out all of the swagger meta tags that are embedded as comments.  The output of the **gen_swaggerspec.sh** script will be the **swaggerspec/swagger.json** spec file. 

While the Docker services are running, just re-run **gen_swaggerspec.sh** and hit refresh on the page to see the Swagger doc updates in real time.

### Running the web services

Once your `swaggerspec/swagger.json` file has been generated you will want to render it to verify it's contents with the HTTP web rendering services.

The `docker-compose.yml` will start two rendering services, a custom http service for hosting the `swaggerspec/swagger.json` and the Swagger UI.  

To start the Swagger UI services (and build them if not already built) just run:

```$ docker compose up```

NOTE: Iterative Workflow Tips:

Blow away only the local images (excluding remote ones) and bring down the container:
```$ docker compose down --rmi local```

Blow away all the images (including remote ones) and bring down the container:
```$ docker compose down --rmi all```

Once started navigate your browser to [http://localhost:8080](http://localhost:8080)

### Converting the swaggerspec/swagger.json to .rst

After you generate the `swaggerspec/swagger.json` from the steps above use the `swaggerspec` Docker Compose file to convert the `swagger.json` to .rst so that it can merged in with the existing Traffic Control documentation.

* `$ cd swaggerspec`
* `$ docker compose up` - will convert the `swagger.json` in this directory into `v13_api_docs.rst`
* `$ cp v13_api_docs.rst ../../../../../docs/source/development/traffic_ops_api`
* `$ cd ../../../../../docs`
* `$ make` - will generate all the Sphinx documentation along with the newly generated TO Swagger API 1.3 docs

NOTE: Iterative Workflow Tips:

Blow away only the local images (excluding remote ones) and bring down the container:
```$ docker compose down --rmi local```

Blow away all the images (including remote ones) and bring down the container:
```$ docker compose down --rmi all```
