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

# t3c

The `t3c` app generates and applies cache configuration for Traffic Control.

# Commands

The `t3c` app is divided into commands, each of which can be run as an argument to `t3c` or by themselves, e.g. `t3c apply ...` or `t3c-apply ...`.

|               command                  | description |
| -------------------------------------: | :---------- |
| [t3c-apply](./t3c-apply/README.md)     | Generate and apply cache configuration |
| [t3c-check](./t3c-check/README.md)    | Check that new config can be applied |
| [t3c-diff](./t3c-diff/README.md)       | Diff config files, like diff or git-diff but with config-specific logic |
| [t3c-request](./t3c-request/README.md) | Request data from Traffic Ops |
| [t3c-update](./t3c-update/README.md)   | Update a server's queue and reval status in Traffic Ops |
