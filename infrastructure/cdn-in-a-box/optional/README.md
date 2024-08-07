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

## CDN-In-A-Box Optional Container(s)

Create an alias to utilize these container(s) with the core CDN-In-A-Box stack. Note, that the exposed port(s) have been moved to an optional docker compose file to allow for concurrent CiaB instances.

From the top-level directory of `cdn-in-a-box` create the following alias:

```
alias mydc="docker compose "` \
        `"-f $PWD/docker-compose.yml "` \
        `"-f $PWD/docker-compose.expose-ports.yml "` \
        `"-f $PWD/optional/docker-compose.$NAME1.yml "` \
        `"-f $PWD/optional/docker-compose.$NAME1.expose-ports.yml "` \
        `"-f $PWD/optional/docker-compose.$NAME2.yml "` \
        `"-f $PWD/optional/docker-compose.$NAME2.expose-ports.yml "
```

For example, to add the socksproxy and vnc optional container(s), use the following alias:


```
alias mydc="docker compose "` \
        `"-f $PWD/docker-compose.yml "` \
        `"-f $PWD/docker-compose.expose-ports.yml "` \
        `"-f $PWD/optional/docker-compose.socksproxy.yml "` \
        `"-f $PWD/optional/docker-compose.socksproxy.expose-ports.yml "` \
        `"-f $PWD/optional/docker-compose.vnc.yml "` \
        `"-f $PWD/optional/docker-compose.vnc.expose-ports.yml "
```

To start the CDN-In-A-Box stack:

```
mydc build
mydc rm -fv
mydc up
```
