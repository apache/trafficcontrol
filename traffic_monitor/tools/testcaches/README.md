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

# testcaches

The `testcaches` tool simulates multiple ATS caches' `_astats` endpoints.

Its primary goal is for testing the Monitor under load, but it may be useful for testing other components.

A list of parameters can be seen by running `./testcaches -h`. There are only three: the first port to use, the number
of ports to use, and the number of remaps (delivery services) to serve in each fake server.

Each port is a unique fake server, with distinct incrementing stats.

When run with no parameters, it defaults to ports 40000-40999 and 1000 remaps.

Stats are served at the regular ATS `stats_over_http` endpoint, `_astats`. For example, if it's serving on port 40000,
it can be reached via `curl http://localhost:40000/_astats`. It also respects the `?application=system` query parameter,
and will serve only system stats (the Monitor "health check" [as opposed to the "stat check"]). For
example, `curl http://localhost:40000/_astats?application=system`.

## Commands

The `testcaches` app accepts a number of commands, which manipulate the data it serves. These command are all available
via HTTP requests.

Each HTTP request is made to a fake cache at a specific port. Thus, you can modify data served by each fake cache
independently.

The commands are:

### `/cmd/setstat`

Sets how much a stat increments by every interval (currently, an interval is hard-coded to 1 second). Accepts a min and
max, and will increment by a random number between them. The min may equal the max, if a constant increment is desired.

Query Parameters:
`remap` - the remap rule to set
`stat` - the stat to set
`min` - the minimum number to increment by
`max` - the minimum to increment by

Example:
`curl -Lvsk http://localhost:4242/cmd/setstat?remap=num1.example.net&stat=out_bytes&min=10&max=25`

### `setsystem`

Sets system stats to constant values. Multiple stats may be set with a single request.

Query Parameters:
`loadavg1m` - the 1m loadavg in the `system` object.
`loadavg5m` - the 5m loadavg in the `system` object.
`loadavg10m` - the 10m loadavg in the `system` object.
`speed` - the network interface speed in the `system` object. This number is in kilobits. I.e. 20000 means 20Gbps.

Example:
`curl -sk 'http://localhost:4242/cmd/setsystem?loadavg1m=10.1&loadavg5m=27.92&loadavg10m=3.4&speed=20000' `

### `setdelay`

Sets the delay for serving all _astats requests to this fake cache. Accepts a minimum and maximum, which may be qual,
and delays the request by a random interval between them. When a delay is set, the server immediately accepts client
requests, reads headers and sets up the connection, and then delays writing out the body.

Query Parameters:
`min` - the minimum delay time, in milliseconds
`max` - the maximum delay time, in milliseconds

Example:
`curl -Lvsk 'http://localhost:4242/cmd/setdelay?min=200&max=600'`

## Docker

Build environment variables: none

Run environment variables:

- `NUM_PORTS`  - app `numPorts` argument
- `NUM_REMAPS` - app `numRemaps` argument
- `PORT_START` - app `portStart` argument
