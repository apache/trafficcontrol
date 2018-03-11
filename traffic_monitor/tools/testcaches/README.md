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

A list of parameters can be seen by running `./testcaches -h`. There are only three: the first port to use, the number of ports to use, and the number of remaps (delivery services) to serve in each fake server.

Each port is a unique fake server, with distinct incrementing stats.

When run with no parameters, it defaults to ports 40000-40999 and 1000 remaps.

Stats are served at the regular ATS `stats_over_http` endpoint, `_astats`. For example, if it's serving on port 40000, it can be reached via `curl http://localhost:40000/_astats`. It also respects the `?application=system` query parameter, and will serve only system stats (the Monitor "health check" [as opposed to the "stat check"]). For example, `curl http://localhost:40000/_astats?application=system`.
