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

# cache-config

Traffic Control cache configuration is done via the `t3c` app and its ecosystem of sub-apps.

These are provided in the RPM `trafficcontrol-cache-config`.

To apply Traffic Control configuration and changes to caches, users will typically run `t3c` periodically via `cron` or some other system automation mechanism. See [t3c](./t3c/README.md).

The `t3c` app is an ecosystem of apps that work together, similar to `git` and other Linux tools. The `t3c` app itself has commands to proxy the other apps, as well as a mode to generate and apply the entire configuration.

## Documentation
Each sub-command provides a README.md file. If your system has Pandoc and GNU
`make` and `date` (UNIX `date` untested, might work - UNIX `make` definitely
won't), they can also be transformed into Linux/UNIX "manual pages" using
`make` (or explicitly `make man`) and reStructuredText documentation using
`make rst`.
