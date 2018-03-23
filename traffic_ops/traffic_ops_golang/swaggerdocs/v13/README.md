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

See the install documentation for [https://github.com/go-swagger/go-swagger](go-swagger)


### Generate your Documentation

The **gen_docs.sh** script will scan all the Go files in the swaggerdocs directory and extract out all of the swagger meta tags that are embedded as comments.  The output of the **gen_docs.sh** script will be the **swagger.json** spec file.

### Verifying your Documentation

Once the **swagger.json** spec file has been generated it needs to to be served over http so that you can validate it using the Swagger Editor.  

See the following steps:

*    Execute the **cors-http-server.py** (this will start a server on **http://localhost:8000**
  so that you can point to it using the [https://editor.swagger.io](Swagger Editor).  
  
  `$ ./cors-http-server.py`

*    Navigate to [https://editor.swagger.io](Swagger Editor)
    
*    Use File->Import URL then plugin **http://localhost:8000**
	* At this point the Swagger Editor will convert the **swagger.json** to yaml format and show the resulting documentation rendered as html.

	OR
	
*	 Install the [https://swagger.io/swagger-ui/](Swagger UI) yourself and run locally.
	
  

