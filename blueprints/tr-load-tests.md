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
# Traffic Router Ultimate Test Harness <!-- a concise title for the new feature this blueprint will describe -->

## Problem Description
<!--
*What* is being asked for?
*Why* is this necessary?
*How* will this be used?
-->
As the entrypoint of an Apache Traffic Control CDN, Traffic Router is the most exposed, most critical component, and it must be able to route traffic quickly and at a very high rate. Although Traffic Router has met this requirement in the past, if Traffic Router undergoes no performance tests as it evolves, there is no guarantee that its performance will not decline with time.

Important times to run performance tests:

* When traffic volume to your CDN changes significantly
* When making harware changes to the server hosting the Traffic Router instance
* When upgrading Traffic Router to a new Apache Traffic Control version
* When developing or maintaining a Traffic Router feature that could have an impact on any aspect of Traffic Router performance, including before, during, and after pull request review
* When a commit that modifies Traffic Router is pushed to a GitHub branch, for all Apache Traffic Control branches

A load test for Delivery Services exists in the project at [`/test/router`](https://github.com/apache/trafficcontrol/tree/RELEASE-6.0.1/test/router), but
- It has not been maintained over time, currently does not work
- Only tests HTTP Delivery Services
- Does not support testing Coverage Zone Maps
- Does not use the TO Client Library for requests to Traffic Ops
- Prompts only for the number of requests to make to Delivery Services, not a length of time to run the test
- Does not fail if some minimum threshold of requests per second is not met
- Is not configurable in other ways, such as length of paths generated for HTTP requests to Delivery Service, which client IP address to use

## Proposed Change
<!--
*How* will this be implemented (at a high level)?
-->
The Traffic Router Ultimate Test Harness will include an end-to-end performance test suite verify that the features of Traffic Router meet expected performance thresholds, as well as additional end-to-end tests of other Traffic Router features.

The TR Ultimate Test Harness may extend [`/test/router`](https://github.com/apache/trafficcontrol/tree/RELEASE-6.0.1/test/router) where possible, but it should not limit itself for that secondary goal.

### Traffic Portal Impact
<!--
*How* will this impact Traffic Portal?
What new UI changes will be required?
Will entirely new pages/views be necessary?
Will a new field be added to an existing form?
How will the user interact with the new UI changes?
-->
No Traffic Portal impact is anticipated.

### Traffic Ops Impact
<!--
*How* will this impact Traffic Ops (at a high level)?
-->
No Traffic Ops impact is anticipated.

#### REST API Impact
<!--
*How* will this impact the Traffic Ops REST API?

What new endpoints will be required?
How will existing endpoints be changed?
What will the requests and responses look like?
What fields are required or optional?
What are the defaults for optional fields?
What are the validation constraints?
-->
No Traffic Ops REST API impact is anticipated.

#### Client Impact
<!--
*How* will this impact Traffic Ops REST API clients (Go, Python, Java)?

If new endpoints are required, will corresponding client methods be added?
-->
Clients importing the `github.com/apache/trafficcontrol/v8/lib/go-tc` package will optionally be able to import a constant for `X-MM-Client-IP`, a request header Traffic Router to specify to Traffic Router the IP address to use to geolocate that client:
https://github.com/apache/trafficcontrol/blob/1ed2964d16618aeebef142b01a538336a44d07dd/traffic_router/core/src/main/java/org/apache/traffic_control/traffic_router/core/request/HTTPRequest.java#L29

Additionally, a struct used to unmarshall a Coverage Zone File could be placed in `lib/go-tc`.

#### Data Model / Database Impact
<!--
*How* will this impact the Traffic Ops data model?
*How* will this impact the Traffic Ops database schema?

What changes to the lib/go-tc structs will be required?
What new tables and columns will be required?
How will existing tables and columns be changed?
What are the column data types and modifiers?
What are the FK references and constraints?
-->
No Data Model impact is anticipated.

### Cache Config Impact
<!--
*How* will this impact ORT?
-->
No Cache Config impact is anticipated.

### Traffic Monitor Impact
<!--
*How* will this impact Traffic Monitor?

Will new profile parameters be required?
-->
No Traffic Monitor impact is anticipated.

### Traffic Router Impact
<!--
*How* will this impact Traffic Router?

Will new profile parameters be required?
How will the CRConfig be changed?
How will changes in Traffic Ops data be reflected in the CRConfig?
Will Traffic Router remain backwards-compatible with old CRConfigs?
Will old Traffic Routers remain forwards-compatible with new CRConfigs?
-->
The addition of the TR Ultimate Test Harness themselves will not change Traffic Router functionality in any way. For visibility, however, the TR Ultimate Test Harness should reside in a directory within the `traffic_router` directory. This will be the first time since 545929f7cc that Golang sources will exist in the `traffic_router` directory, so any assumption that all sources within the `traffic_router` directory directly impact Traffic Router's ability to compile should be abandoned.

The TR Ultimate Test Harness should not be included in the Traffic Router RPM, as it is meant to be run on a host separate from Traffic Routers.

### Traffic Stats Impact
<!--
*How* will this impact Traffic Stats?
-->
No Traffic Stats impact is anticipated.

### Traffic Vault Impact
<!--
*How* will this impact Traffic Vault?

Will there be any new data stored in or removed from Riak?
Will there be any changes to the Riak requests and responses?
-->
No Traffic Vault impact is anticipated.

### Documentation Impact
<!--
*How* will this impact the documentation?

What new documentation will be required?
What existing documentation will need to be updated?
-->
Instructions for using the Traffic Router Ultimate Test Harness should be added to the documentation. This should include:
- Small rationale for inclusion of TR Ultimate Test Harness
- Setup instructions
    * The permissions that a user running the Traffic Router Test Harness should have:
        - CDN snapshots
        - CDN information
        - information about Traffic Router-type Servers in those CDNs
        - Type information
        - Delivery Service information
    - Documentation of each option
    - Example commands

### Testing Impact
<!--
*How* will this impact testing?

What is the high-level test plan?
How should this be tested?
Can this be tested within the existing test frameworks?
How should the existing frameworks be enhanced in order to test this properly?
-->
#### Load Tests
The Router Ultimate Test Harness should include a load test for HTTP-routed Delivery Services and for DNS-routed Delivery Services

##### Load Test Options

| Option                          | Description                                                                                                                                                                                             | Delivery Service Type | Default                     |
|---------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------|-----------------------------|
| IPv4 TR addresses only          | Test IPv4 Traffic Router addresses only                                                                                                                                                                 | HTTP, DNS             | False                       |
| IPv6 TR addresses only          | Test IPv6 Traffic Router addresses only                                                                                                                                                                 | HTTP, DNS             | False                       |
| CDN name                        | The name of a CDN to search for Delivery Services                                                                                                                                                       | HTTP, DNS             | all                         |
| Delivery Service name           | The name (XMLID) of a Delivery Service to use for tests                                                                                                                                                 | HTTP, DNS             | None                        |
| Traffic Router name             | Instead of iterating through Traffic Routers, test only a specific Traffic Router, identified by hostname.                                                                                              | HTTP, DNS             | all                         |
| Client IP address               | If provided, Traffic Router will use the value of the `X-MM-Client-IP` request header as the IP address that Traffic Router's geolocation considers. This option should you specify such an IP address. | HTTP, DNS             | None                        |
| Use coverage zone map           | Whether to use an IP address from the Traffic Router's Coverage Zone File                                                                                                                               | HTTP, DNS             | False                       |
| Coverage zone location          | The coverage zone location to use (implies *Use coverage zone map*)                                                                                                                                     | HTTP, DNS             | None                        |
| *Requests per second* threshold | The minimum number of requests per second a Traffic Router must successfully respond to                                                                                                                 | HTTP, DNS             | 8000 for HTTP, 7200 for DNS |
| Benchmark time                  | The duration of each load test, in seconds                                                                                                                                                              | HTTP, DNS             | 300                         |
| Thread count                    | The number of threads to spawn for each test                                                                                                                                                            | HTTP, DNS             | 12                          |
| Path count                      | The number of paths to generate for use in requests to Delivery Services                                                                                                                                | HTTP                  | 10000                       |
| Maximum path length             | The maximum string length for each generated path                                                                                                                                                       | HTTP                  | 100                         |
| Use location header             | Whether the HTTP HTTP Delivery service should redirect the user or server the routing information as a JSON response.                                                                                   | HTTP                  | True                        |

These options will be structured in a config file:

```jsonc
{
    "all": {
        "cdn_name": "Kabletown"
    },
    "http": {
        "ipv4_only": true,
        /* more options */
        "path_count": 5000
    },
    "dns": {
        "delivery_service_name": "static"
    }
}
```

Options that should apply to a specific type of test should go under a key named that test type, while options applying to all tests can go under the `"all"` key.

#### Other Tests

Additionally, the TR Ultimate Test Harness should provide the ability to verify that a DNS-routed Delivery Service assigned to a Federation resolves to that Federation's CNAME, rather than a Cache's IP address, depending on the IP address of the client querying Traffic Router.

### Automation Impact

A GitHub Action to run the tests that the TR Ultimate Test Harness should be added, but only if it consistently meets a constant *requests per second* threshold. Traffic Router will perform better on some GitHub Actions runners than others, so this should be tested after writing the GitHub Action.

If a meaningful *requests per second* threshold cannot be found for GitHub Actions runners, we may consider trying again in the future, in case consistency of Traffic Router performance on GitHub Actions runners improves.

### Performance Impact
<!--
*How* will this impact performance?

Are the changes expected to improve performance in any way?
Is there anything particularly CPU, network, or storage-intensive to be aware of?
What are the known bottlenecks to be aware of that may need to be addressed?
-->
By increasing our attention to Traffic Router's performance, Traffic Router's performance should not decrease, and its performance may increase.

### Security Impact
<!--
*How* will this impact overall security?

Are there any security risks to be aware of?
What privilege level is required for these changes?
Do these changes increase the attack surface (e.g. new untrusted input)?
How will untrusted input be validated?
If these changes are used maliciously or improperly, what could go wrong?
Will these changes adhere to multi-tenancy?
Will data be protected in transit (e.g. via HTTPS or TLS)?
Will these changes require sensitive data that should be encrypted at rest?
Will these changes require handling of any secrets?
Will new SQL queries properly use parameter binding?
-->
If Traffic Router is vulnerable to denial-of-service attacks relating to HTTP requests or DNS queries, there is potential for the TR Ultimate Test Harness to uncover such vulnerabilities.

### Upgrade Impact
<!--
*How* will this impact the upgrade of an existing system?

Will a database migration be required?
Do the various components need to be upgraded in a specific order?
Will this affect the ability to rollback an upgrade?
Are there any special steps to be followed before an upgrade can be done?
Are there any special steps to be followed during the upgrade?
Are there any special steps to be followed after the upgrade is complete?
-->
An Apache Traffic Control administrator may choose to use the TR Ultimate Test Harness to verify that they can get as good of Traffic Router performance in a new Apache Traffic Control version as they can in the version they are upgrading from. If, according to their testing, performance decreases in the newer Traffic Router, the administrator may choose to delay upgrading until they are able to attain the same level of performance, either by changing Traffic Router configuration or waiting for an even newer Traffic Router version to be released.

### Operations Impact
<!--
*How* will this impact overall operation of the system?

Will the changes make it harder to operate the system?
Will the changes introduce new configuration that will need to be managed?
Can the changes be easily automated?
Do the changes have known limitations or risks that operators should be made aware of?
Will the changes introduce new steps to be followed for existing operations?
-->
Needless to say, the Traffic Router Ultimate Test Harness should not be run against a Traffic Router that is simultaneously being used to route production traffic. In order to avoid this. the Traffic Router Ultimate Test Harness should only be run in non-production environments.

### Developer Impact
<!--
*How* will this impact other developers?

Will it make it easier to set up a development environment?
Will it make the code easier to maintain?
What do other developers need to know about these changes?
Are the changes straightforward, or will new developer instructions be necessary?
-->
When developing Traffic Router features, especially ones found or anticipated to affect Traffic Router performance, a developer may choose to test the result of the new feature on Traffic Router's performance using the Traffic Router Ultimate Test Harness.

## Alternatives
<!--
What are some of the alternative solutions for this problem?
What are the pros and cons of each approach?
What design trade-offs were made and why?
-->
`wrk` is an alternative for HTTP load testing.

`flamethrower` is an alternative for DNS load testing.

## Dependencies
<!--
Are there any significant new dependencies that will be required?
How were the dependencies assessed and chosen?
How will the new dependencies be managed?
Are the dependencies required at build-time, run-time, or both?
-->
No additional dependencies are anticipated. If additional Go dependencies are required, those dependencies will be added to the Apache Traffic Control [`go.mod`](https://github.com/apache/trafficcontrol/blob/RELEASE-6.0.0/go.mod) file.

## References
<!--
Include any references to external links here.
-->
- Comment in the [#traffic-control-traffic_router](https://the-asf.slack.com/archives/C023PVBLWE4) channel of the ASF Slack where [@limited](https://github.com/limited) recommends `wrk` and `flamethrower`: https://the-asf.slack.com/archives/C023PVBLWE4/p1631639255008200
