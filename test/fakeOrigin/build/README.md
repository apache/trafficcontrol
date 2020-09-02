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
# Building RPMs & Binaries
To build the RPMs and binaries you will need docker.
## Usage
```bash
docker-compose -f docker-compose.build_rpm.yml build --no-cache builder
docker-compose -f docker-compose.build_rpm.yml up --force-recreate --exit-code-from builder
```
This should dump all artifacts in the %repository root%/dist directory when it's done (like all other ATC build jobs).

## Version
If you need to manipulate the version information, you can export a few environment variables to supply your own overrides for the [VERSION](../version/VERSION) file.
* VER_MAJOR (integer)
* VER_MINOR (integer)
* VER_PATCH (integer)
* VER_DESC (short string)
* VER_COMMIT (git short hash string)
* BUILD_NUMBER (integer)
