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

# Traffic Portal

An AngularJS client served from a lightweight Node.js web server. Traffic Portal was designed to consume the Traffic Ops API.

Installation / configuration instructions may be found in [the `build/` directory](./build/)

## Server Options
`Usage: node /path/to/server.js [-c CONFIG]`

`-c CONFIG`
    Specify a configuration file to use at path `CONFIG`, rather than just one of `conf/config.js`, `conf/configDev.js`, or the RPM install location (`/etc/traffic_portal/conf/config.js`)
