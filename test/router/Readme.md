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

# Traffic Router Tests

## Load Test

You can simulate a mix of HTTP and HTTPS traffic for a CDN by choosing
the number of HTTP delivery services and the number HTTPS delivery services the test will
exercise.

There are 2 parts to the load test.

* A web server that makes the actual requests and takes commands to
fetch data from the CDN, start the test, and return current results.
* A web page that's used to run the test and see the results.

### Running the Load Test

1. You will need to make sure you have a CA file on your filesystem
2. The web server is a go program, set your GOPATH env variable appropriately
3. Open a terminal program and navigate to the traffic_control/test/router/server directory
4. execute the command `go run server.go`
5. Open the file traffic_control/test/router/index.html
6. Authenticate against a Traffic Ops host, should be an instantaneous operation, you can watch the output from server.go for feedback
7. Enter the Traffic Ops host in the second form and click the button to get a list of CDN's
8. Wait for the web page to show a list of CDN's under the above form, this may take several seconds
9. The List of CDN's will display the number of Http and Https capable delivery services that may be exercised
10. Choose the CDN you want to exercise from the dropdown
11. Fill out the rest of the form, enter appropriate numbers for each http and https delivery services
12. Click Run Test
13. As the test runs the web page will occaisionally report results including running time, latency, and throughput