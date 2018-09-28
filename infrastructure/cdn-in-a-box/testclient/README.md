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

# Private VNC/proxy testclient for ApacheConf Demo

## To start/install the VNC/proxy container

```
alias vdc='docker-compose -f docker-compose.yml -f docker-compose.testclient.yml'
vdc build 
vdc kill && vdc rm -f 
docker volume prune -f
vdc up
```

## Install a VNC client from the apple store or with brew
- host/port `localhost:55900`
- Password for VNC session is 'demo' or whatever $VNC_PASSWD is set to in Dockerfile

## URL locations within the VNC/proxy container:
* Traffic Portal: https://trafficportal.infra.ciab.test
* Traffic Monitor: http://trafficmonitor.infra.ciab.test
* Demo1 Delivery Service: http://video.demo1.mycdn.ciab.test/index.html

## TODO:
* On both linux/osx platforms, allow connectivity with docker daemon within the container
* Generate a chopped up HLS movie
