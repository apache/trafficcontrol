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

# VNC Client Container for CDN-In-A-Box

This container provides a basic lightweight window manager (fluxbox), firefox browser, xterm, and a few other utilities for the purpose of testing the traffic control components within the CDN-In-A-Box tcnet bridge network.  Additionally, the generated/supplied self-signed SSL CA certificate has been installed in the system trust store, which allows all applications and utilities to enjoy full SSL/TLS x509 certificate validation.

## VNC Container Lifecycle

```
# From infrastructure/cdn-in-a-box
alias mydc="docker compose -f $PWD/docker-compose.yml -f $PWD/optional/docker-compose.vnc.yml"
docker volume prune -f
mydc rm -fv 
mydc kill 
mydc build 
mydc up
```

## VNC Clients
- Tight VNC Client: https://www.tightvnc.com
- Real VNC Client: https://www.realvnc.com

## VNC Connection Information
- Docker Host/Port: localhost:5909
- Container Host/Port: vnc.infra.ciab.test:5909
- User: ciabuser
- Password: must be set via $VNC_PASSWD (See below)

## Environment Variables

Environment variables that can be set prior to starting the vnc container.

Defaults:
```
VNC_USER=ciabuser
VNC_PASSWD=Random String (change in variables.env in top-level dir)
VNC_DEPTH=32
VNC_RESOLUTION=1440x900
```
