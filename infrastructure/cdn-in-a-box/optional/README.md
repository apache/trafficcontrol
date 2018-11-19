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

Create an alias to utilize these container(s) with the core CDN-In-A-Box stack

From the top-level directory of `cdn-in-a-box` create the following alias:

```
alias mydc='docker-compose -f docker-compose.yml -f optional/docker-compose.$NAME1.yml -f optional/docker-compose.$NAME2.yml'
```

For example, to use the vnc optional container, use the following alias:


```
alias mydc='docker-compose -f docker-compose.yml -f optional/docker-compose.vnc.yml'

```

To start the CDN-In-A-Box stack:

```
mydc build
mydc up
```
