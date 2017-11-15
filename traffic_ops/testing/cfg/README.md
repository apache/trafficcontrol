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

# Traffic Ops Configuration Test

Use this tool to compare API output between two instances of Traffic Ops API.
The test logs in to each instance using the same user/password combination,
then processes all given endpoints.  Any that get different results are
reported and written to files in the output directory (default /tmp/gofiles).

## Requirements

Two Traffic Ops instances with common login credentials.  The following
environment variables must be set:
- `TO_URL`   -- the *reference* Traffic Ops (e.g. production version)
- `TEST_URL` -- the *test* Traffic Ops (e.g. development version)
- `TO_USER`
- `TO_PASSWORD`

## Usage

```
   go test -v
```


