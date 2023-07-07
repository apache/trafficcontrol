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
# Add Varnish Cache Support

## Problem Description

Currently Traffic Control uses Traffic Server as the underlying cache server. We can expand that by introducing Varnish cache as an option for the cache server used with its great performance, robustness and modularity.

## Proposed Change

From a high level point of view, ATS operates based on configuration files that describe in details how it should work and interact with other servers in the cache hierarchy. These configuration files are managed and generated using `t3c` components that utilize Traffic Ops APIs to get profiles and parameters data required for the configuration files. The proposed change is to use the same data fetched from Traffic Ops APIs to generate configuration files for Varnish cache with almost the same functionality.

Note: the changes should not affect existing components but rather build on them.

### Traffic Portal Impact

n/a

### Traffic Ops Impact

n/a

#### REST API Impact

n/a

#### Client Impact

n/a

#### Data Model / Database Impact

- A new profile type will be added `VARNISH_PROFILE` but it won't affect the model or other structures in `lib/go-tc`
- `DeliveryService` structs contain fields related to ATS like `remapText`. It will be parsed and translated to Varnish configuration.

### ORT Impact

- `go-varnish` package will be developed to handle generating configuration files for Varnish, Hitch, `varnishncsa` and `logrotate`. For detailed description of mapping configuration files from ATS to Varnish refer to [Varnish Support](https://github.com/apache/trafficcontrol/wiki/Varnish-Support) wiki.
- New subcommands will be added to `t3c` including `t3c-varnish-generate` that works similar to `t3c-generate` but instead will depend on `go-varnish` package to generate different configuration files and `t3c-varnish-apply` that handles applying these configuration files, restarting and reloading different services running.
- `go-atscfg` will be refactored to export some of its functionality to be reusable from `go-varnish`. So, instead of rewriting the logic of which IPs are allowed for specific HTTP requests, it could be separated and exported in a function that both packages utilize.

### Traffic Monitor Impact

New statistics parser will be added to Traffic Monitor to handle data coming from Varnish cache statistics endpoint. There is no `VMOD` that exposes Varnish statistics so a service that keeps polling `varnishstat` will be developed.

### Traffic Router Impact

n/a

### Traffic Stats Impact

n/a

### Traffic Vault Impact

n/a

### Documentation Impact

New documentation will be needed for how to setup Varnish with TC and what is the differences between Varnish and ATS.

### Testing Impact

In addition to unit tests and integration tests, [`varnishtest`](https://varnish-cache.org/docs/trunk/reference/varnishtest.html) could be used to test Varnish cache is operating as expected.

### Performance Impact

For current components there should be no performance impact. However, between Traffic Server and Varnish it isn't clear yet what the difference in performance will be.

### Security Impact

n/a

### Upgrade Impact

n/a

### Operations Impact

n/a

### Developer Impact

n/a

## Alternatives

n/a

## Dependencies

- Varnish and its utilities (`varnishtest`, `varnishstat`, `varnishncsa`, ...).
- Hitch to manage SSL part of the server.
- `logrotate` to manage logs generated from `varnishncsa`.

## References

- https://github.com/apache/trafficcontrol/wiki/Varnish-Support
- https://varnish-cache.org/docs/trunk/reference/
- https://varnish-cache.org/vmods/
- https://github.com/apache/trafficcontrol/wiki/Varnish-Support
