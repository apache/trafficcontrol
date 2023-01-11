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
# Traffic Portal Test Automation
This directory contains integration tests for Traffic Portal.

## Prerequisites
* [Node](https://nodejs.org) version 16 or later.

## Building and Running
To build and run the tests, one can use the `npm` (or `pnpm`) scripts.

* Install the dependencies with `npm install` (or `pnpm install`)
* Run the webdriver, either in the background or in a separate terminal -
	because it's a long-running process - with `npm run start-webdriver` (or
	`pnpm run start-webdriver`)
* Run the tests with `npm test` (or `pnpm test`)

## Command Line Parameters
The tests can accept a few command line parameters - which can be separated
from the `npm` flags with `--`.

| Flag                            | Description                                                                                          |
| ------------------------------- | :--------------------------------------------------------------------------------------------------: |
| params.baseUrl                  | Environment test run on. Tests are written for cdn-in-a-box only. Do not run on other environment                                   |
| capabilities.shardTestFiles     | Input `true` or `false` to turn on or off parallelization. If the value is false, maxInstances will always count as 1. The default value in the config file = false                            |
| capabilities.maxInstances       | Input number of Chromium instances that your machine can handle. Test will fail if local machine cannot handle a lot of Chromium instances. The default value = 1    |

### Example
```bash
npm test -- --params.baseUrl https://localhost --capabilities.shardTestFiles true
```
