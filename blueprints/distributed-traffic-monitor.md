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
# Distributed Traffic Monitor

## Problem Description
Currently, TM polls all caches in a CDN. As CDNs grow, this becomes a major
pain point as TM is limited by the amount of bandwidth and CPU it requires to
receive and process data from every cache on the CDN, and scaling vertically by
running it on better hardware is only feasible up to a certain point. Also, the
performance of a cache observed by a TM which is very far away from it does not
always reflect the performance observed by clients that are actually using the
cache (because the clients are typically much closer to it).

## Proposed Change
TM should have the ability to poll only a subset of caches in a CDN and peer
with other TMs which are monitoring other subsets in order to get a full view
of the CDN's health. This would allow us to run TM in a more distributed manner
across the CDN, giving us a view of cache health that is closer to what clients
actually observe and enabling us to scale TM horizontally. Additionally, we
would like to have the option to disable _stat polling_ in order for these
distributed TMs to focus on _health polling_.

### Traffic Portal Impact
This proposal does not require any TP changes.

### Traffic Ops Impact
This proposal might have limited impact on TO. The existing TO API endpoints
already provide the data that TM will need to run in a distributed manner, and
any changes made to TM APIs that TO uses will remain backwards-compatible.
However, TO may need to be updated if it uses any stat-polling-related TM APIs
so that it only requests from TMs that have stat-polling enabled.

### t3c Impact
This proposal does not require `t3c` changes. Note: the `tc-health-client`
periodically polls a random TM to get cache health states, and because
distributed TMs will still serve the cache health states of all caches in a
CDN, there will be no impact to the `tc-health-client`. It can continue to poll
any random TM and still get all the cache health data for the entire CDN.

### Traffic Monitor Impact
TM will gain at least two more configuration options:
- `distributed_polling_enabled` (default: false) - when set to true, TM will
  run in _distributed mode_ (more details on this below). When set to false, TM
  will run in its legacy, normal mode.
- `stat_polling_disabled` (default: false) - when set to true, TM will _not_ do
  stat polling for caches. When set to false, TM will do stat polling for
  caches (legacy, normal behavior). Initially, this must be set to true if
  `distributed_polling_enabled` is also set to true. In a later phase of
  development, we will add the ability to enable stat polling in distributed
  mode.

Note: these are configuration options as opposed to profile parameters because
we currently do not have the capability to have per-profile monitoring.json
snapshots (or per-TM-server configuration in one snapshot).

To use _distributed mode_, generally all TMs in the CDN need to be running in
distributed mode (if they're taking part in the health protocol). It should
still be possible to run TMs in the _legacy_ (non-distributed) mode in order to
provide cache stat polling (which is important for Traffic Stats), but they
should not be set to `ONLINE` in order to keep them from interfering with the
health protocol.

While in _distributed mode_, a TM instance will only monitor a subset of
cachegroups in its given CDN. The number of cachegroups each TM will monitor
depends on the number of cachegroups that contain TM servers for the CDN. These
will be referred to as "TM groups." A TM group contains 1 to many TM servers,
and a CDN can have 1 to many TM groups. If there are N TM groups, each TM group
will monitor roughly 1/N of the cachegroups in the CDN. Each TM in the group
will monitor all of caches in that 1/N portion of cachegroups that the TM group
is responsible for. For example, if there are 10 cachegroups and 3 TM groups:
- TM group 1 monitors cachegroups 1-4
- TM group 2 monitors cachegroups 5-7
- TM group 3 monitors cachegroups 8-10

Because every TM can serve the health state of every cache, distributed TMs
will need to peer not only with their own group members but also with other
groups as well. However, instead of simultaneously requesting cache health
states from all out-of-group peers, each distributed TM will simultaneously
request cache health states from 1 TM in every other TM group, alternating
between group members in a deterministic, round-robin fashion. For this
out-of-group peering, a new TM API route will be added that returns only the
cache health states for caches that the TM group is responsible for polling.

A safety feature will be added to TM (while running in distributed mode) to
ensure that all cachegroups are polled by at least 1 TM group, and an
additional profile parameter override will be available in order to manually
assign cachegroups to TM groups for polling.

### Traffic Router Impact
This proposal should have no impact on TR.

### Traffic Stats Impact
Because we will be able to disable stats polling on TM, TS will need to poll
TMs that actually have stats polling enabled. TMs with polling enabled should
be given a specific server status (other than `ONLINE`), which TS will be
configured to poll, and that might mean creating a new server status
specifically for that purpose.

### Traffic Vault Impact
This proposal has no impact on Traffic Vault.

### Documentation Impact
Any new configuration options added to TM should be documented, and the steps
necessary to run TM in a distributed manner as well as how it works should be
described in some form of documentation (probably the TM admin docs).

### Testing Impact
New TM unit and integration tests should be added where applicable. It would
also be recommended to run both types of TMs in production (distributed and
non-distributed) and compare the reported cache health states between both
types. This would help discover any issues with running TM in a distributed
manner using data from a production environment. However, TR should still get
health states from the non-distributed TMs until we are confident in the health
states reported by distributed TMs.

### Performance Impact
This proposal allows TM to be scaled horizontally, so operators can increase
the number of TM groups in order to get the desired amount of load per TM.

### Security Impact
This proposal does not have much impact on security, but allowing TM to scale
horizontally means that there may be more firewall rules that will need to be
applied to any new TM servers that are deployed. However, TM will not need any
_new_ ports opened, assuming the same `httpListener` and `httpsListener`
configuration is used.

### Upgrade Impact
TMs running in a distributed manner can be upgraded in the same way that
non-distributed TMs are upgraded today. For instance, we would likely upgrade
the `OFFLINE` TMs, then set the upgraded TMs to `ONLINE` while simultaneously
setting the old TMs to `OFFLINE`.

### Operations Impact
There should be little impact on operations other than the effort necessary to
provision and deploy new TM servers to run in a distributed manner. Existing
automation can still be used for upgrades, configuration, etc., but automation
may need a way to differentiate between non-distributed and distributed TMs
within the same environment so that both types are configured differently.

Troubleshooting distributed TMs might be more difficult than non-distributed
TMs as there will be more servers involved. However, the health of a cache
should always be determined by the same TMs (assuming no new TM groups are
added to the system), so it would be best to investigate the TM servers in the
"authoritative" TM group for the cache under investigation. To help aid this
kind of troubleshooting, we may want TM to have an API that returns information
about which TM groups it thinks are currently monitoring which cache groups.

### Developer Impact
Developers should know that once this change is implemented, there will be two
different "run modes" for TM -- distributed and non-distributed. TM will do
certain things differently in the distributed mode compared to the
non-distributed mode even though the vast majority of things will be the same.
Therefore, developers will need to take care to ensure the proper behavior is
followed depending on which "run mode" TM is in.

Also, because this proposal will allow TMs to monitor only a subset of caches,
it may make it easier to set up a development environment using production-like
data and caches. It is somewhat infeasible for most TM development environments
to poll an entire, large CDN, but with distributed TM groups, developers could
essentially choose how many caches they want their local TM to poll.

## Alternatives

- Cache Self-Monitoring: Make caches monitor themselves by using remap rules,
  essentially replacing TM's Cache Health Monitoring. The
  [Proof-of-Concept](https://github.com/apache/trafficcontrol/pull/4529) has
  more details.

## Dependencies
This proposal does not intend to add any new dependencies.

## References
The following mailing list threads were related to this blueprint:
- [Proposal: Distributed Health Monitoring](https://lists.apache.org/thread.html/rf3307f824c0f82892cbb0fea74a5c6a274c8ea4f303d125e8f1212da%40%3Cdev.trafficcontrol.apache.org%3E)
- [Distributed Traffic Monitor Feedback/Requirements](https://lists.apache.org/thread.html/rf985a2b9e8a440d396a0097a71882919bff5b3cb5f8d6c3a53143162%40%3Cdev.trafficcontrol.apache.org%3E)
