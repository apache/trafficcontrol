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

# Socksproxy Container for CDN-In-A-Box

Dantes socks proxy is an optional container that can be used to provide browsers and other clients the ability 
to resolve DNS queries and network connectivity directly to the tcnet bridged interface.  This is very helpful
when running the CDN-In-A-Box stack on OSX/Windows docker host that lack network bridge/ipforward support.

## Socks Container Lifecycle

```
# From infrastructure/cdn-in-a-box
alias mydc="docker compose -f $PWD/docker-compose.yml -f $PWD/optional/docker-compose.socksproxy.yml"
docker volume prune -f
mydc build 
mydc kill 
mydc rm -fv 
mydc up
```

## Socks Connection Information

- Docker Host/Port: localhost:9080
- Container Host/Port: socksproxy.infra.ciab.test:1080

## Socks Browser configuration

- Install the Foxy Proxy browser plugin.
- Add New Proxy, Manual Proxy Configuration, Host: localhost, Port: 9080, Check 'SOCKS Proxy', Select 'SOCKS v5'.
- Enable 'pre-defined and patterns' mode.

## Socks cmdline client environment
 
Some network clients support connections via socks using the socks\_proxy environment variable

```
export http_proxy=socks://localhost:9080/
export https_proxy=socks://localhost:9080/
```
