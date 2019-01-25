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

# grovetccfg

Traffic Control configuration generator for the Grove HTTP caching proxy.

# Building

1. Install and set up a Golang development environment.
    * See https://golang.org/doc/install
2. Clone this repository into your GOPATH.
```bash
mkdir -p $GOPATH/src/github.com/apache/trafficcontrol
cd $GOPATH/src/github.com/apache/trafficcontrol
git clone https://github.com/apache/trafficcontrol/grove
```
3. Build the application
```bash
cd $GOPATH/src/github.com/apache/trafficcontrol/grove/grovetccfg
go build
```
5. Install and configure an RPM development environment
   * See https://wiki.centos.org/HowTos/SetupRpmBuildEnvironment
4. Build the RPM
```bash
./build/build_rpm.sh
```

# Running

You may use a trafficserver profile with your grove deployment but `grovetccfg` will only read the `allow_ip` and the `allow_ip6` parameters from a
traffic server profile when constructing the remap_rules file.  A sample `grove_profile.traffic_ops` file is provided to get you started in creating  a GROVE_PROFILE
type.  When you use a GROVE_PROFILE type, `grovetccfg` will read the settings from the profile and generate the `grove.cfg` file from the settings in that profile.

The `grovetccfg` tool has an RPM, but no service or config files. It must be run manually, even after installing the RPM. Consider running the tool in a cron job.

Example:

`./grovetccfg -api=1.2 -host my-http-cache -insecure -touser carpenter -topass 'walrus' -tourl https://cdn.example.net -pretty > remap.json`

Flags:

| Flag | Description |
| --- | --- |
| `api` | The Traffic Ops API version to use. The default is 1.2. If 1.3 is passed, it will use a newer and more efficient endpoint. |
| `host` | The Traffic Ops server to create configuration from. This must be a cache server in Traffic Ops. |
| `insecure` | Whether to ignore certificate errors when connecting to Traffic Ops |
| `touser` | The Traffic Ops user to use. |
| `topass` | The Traffic Ops user password. |
| `tourl` | The Traffic Ops URL, including the scheme and fully qualified domain name. |
| `pretty` | Whether to pretty-print JSON |

Exit Codes:

| Code | Description |
| --- | --- |
| 0 | Success |
| 1 | Error, see output for details |
| 2 | Error reloading service |
| 3 | Error clearing the server's update flag in Traffic Ops |
