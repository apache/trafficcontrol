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

# Traffic Ops Compare

Use this tool to compare API output between two instances of Traffic Ops API.
It logs in to each instance and then processes all given endpoints.  Any that get
different results are reported and written to files in the output directory
(default ./results).

## Requirements

Two Traffic Ops instances with login credentials (possibly different).  The following
environment variables must be set:

- `TO_URL`   -- the *reference* Traffic Ops (e.g. production version)
- `TO_USER`  -- the username for `TO_URL`
- `TO_PASSWORD`  -- the password for `TO_URL`

- `TEST_URL` -- the *test* Traffic Ops (e.g. development version)

These are optional:

- `TEST_USER`  -- the username for `TO_URL` (default -- same as `TO_USER`)
- `TEST_PASSWORD`  -- the password for `TO_URL` (default -- same as `TO_PASSWORD`)

## Usage

```
   go run compare.go [-results <dir>] [-route <API route] [-file <file of routes>] [-snapshot]
```

Options:

- `-results <dir>` -- directory to write difference results
- `-route <route>` -- a specific route to compare
- `-file <file>`   -- file containing routes to check (-route takes precedence)
- `-snapshot`      -- compare snapshot for each CDN
