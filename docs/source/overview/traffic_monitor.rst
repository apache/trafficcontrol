..
..
.. Licensed under the Apache License, Version 2.0 (the "License");
.. you may not use this file except in compliance with the License.
.. You may obtain a copy of the License at
..
..     http://www.apache.org/licenses/LICENSE-2.0
..
.. Unless required by applicable law or agreed to in writing, software
.. distributed under the License is distributed on an "AS IS" BASIS,
.. WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
.. See the License for the specific language governing permissions and
.. limitations under the License.
..

.. _tm-overview:

***************
Traffic Monitor
***************
Traffic Monitor is an HTTP service that monitors the :term:`cache servers` in a CDN for a variety of metrics. These metrics are for use in determining the overall "health" of a given :term:`cache server` and the related :term:`Delivery Services`. A given CDN can operate a number of Traffic Monitors, from a number of geographically diverse locations, to prevent false positives caused by network problems at a given site. Traffic Monitors operate independently, but use the state of other Traffic Monitors in conjunction with their own state to provide a consistent view of CDN :term:`cache server` health to downstream applications such as :ref:`tr-overview`. `Health Protocol`_ governs the :term:`cache server` and :term:`Delivery Service` availability. Traffic Monitor provides a view into CDN health using several RESTful JSON endpoints, which are consumed by other Traffic Monitors and downstream components such as :ref:`tr-overview`. Traffic Monitor is also responsible for serving the overall CDN configuration to :ref:`tr-overview`, which ensures that the configuration of these two critical components remain synchronized as operational and health related changes propagate through the CDN.

.. _astats:

Cache Monitoring
================
Traffic Monitor polls all :term:`cache servers` configured with a status of ``REPORTED`` or ``ADMIN_DOWN`` at an interval specified as a configuration parameter in :ref:`to-overview`. If the :term:`cache server` is set to ``ADMIN_DOWN`` it is marked as unavailable but still polled for availability and statistics. If the :term:`cache server` is explicitly configured with a status of ``ONLINE`` or ``OFFLINE``, it is not polled by Traffic Monitor and presented to :ref:`tr-overview` as configured, regardless of actual availability. Traffic Monitor makes HTTP requests at regular intervals to a special URL on each Edge-tier :term:`cache server` and consumes the JSON output. The special URL is served by a plugin running on the :abbr:`ATS (Apache Traffic Server)` :term:`cache servers` called `"astats" <https://github.com/apache/trafficcontrol/tree/master/traffic_server/plugins/astats_over_http>`_, which is restricted to Traffic Monitor only. The astats plugin provides insight into application and system performance, such as:

- Throughput (e.g. bytes in, bytes out, etc).
- Transactions (e.g. number of 2xx, 3xx, 4xx responses, etc).
- Connections (e.g. from clients, to parents, origins, etc).
- Cache performance (e.g.: hits, misses, refreshes, etc).
- Storage performance (e.g.: writes, reads, frags, directories, etc).
- System performance (e.g: load average, network interface throughput, etc).

Many of the application-level statistics are available at the global or aggregate level, some at the :term:`Delivery Service` level. Traffic Monitor uses the system-level performance to determine the overall health of the :term:`cache server` by evaluating network throughput and load against values configured in :ref:`to-overview`. Traffic Monitor also uses throughput and transaction statistics at the :term:`Delivery Service` level to determine :term:`Delivery Service` health. If astats is unavailable due to a network-related issue or the system statistics have exceeded the configured thresholds, Traffic Monitor will mark the :term:`cache server` as unavailable. If the :term:`Delivery Service` statistics exceed the configured thresholds, the :term:`Delivery Service` is marked as unavailable, and :ref:`tr-overview` will start sending clients to the overflow destinations for that :term:`Delivery Service`, but the :term:`cache server` remains available to serve other content.

.. seealso:: For more information on :abbr:`ATS (Apache Traffic Server)` statistics, see the `ATS documentation <https://docs.trafficserver.apache.org/en/7.1.x/index.html>`_

.. _health-proto:

Health Protocol
===============

Optimistic Health Protocol
--------------------------
Redundant Traffic Monitor servers operate independently from each other but take the state of other Traffic Monitors into account when asked for health state information. In `Cache Monitoring`_, the behavior of a single Traffic Monitor instance is described. The :dfn:`Health Protocol` adds another dimension to the health state of the CDN by merging the states of all Traffic Monitors into one, and then taking the *optimistic* approach when dealing with a :term:`cache server` or :term:`Delivery Service` that might have been marked as unavailable by this particular instance or a peer instance of Traffic Monitor. Upon startup or configuration change in :ref:`to-overview`, in addition to :term:`cache servers`, Traffic Monitor begins polling its peer Traffic Monitors whose state is set to ``ONLINE``. Each ``ONLINE`` Traffic Monitor polls all of its peers at a configurable interval and saves the peer's state for later use. When polling its peers, Traffic Monitor asks for the raw health state from each respective peer, which is strictly that instance's view of the CDN's health. When any ``ONLINE`` Traffic Monitor is asked for CDN health by a downstream component, such as :ref:`tr-overview`, the component gets the Health Protocol-influenced version of CDN health (non-raw view). In operation of the Health Protocol, Traffic Monitor takes all health states from all peers, including the locally known health state, and serves an optimistic outlook to the requesting client. This means that, for example, if three of the four Traffic Monitors see a given :term:`cache server` or :term:`Delivery Service` as exceeding its thresholds and unavailable, it is still considered available. Only if all Traffic Monitors agree that the given object is unavailable is that state propagated to downstream components. This optimistic approach to the Health Protocol is counter to the "fail fast" philosophy, but serves well for large networks with complicated geography and/or routing. The optimistic Health Protocol allows network failures or latency to occur without affecting overall traffic routing, as Traffic Monitors can and do have a different view of the network when deployed in geographically diverse locations.

Optimistic Quorum
-----------------
In order to prevent split-brain monitoring scenarios, a minimum of three Traffic Monitors are required to properly monitor a given CDN and the optimistic quorum feature should be enabled. If three or more Traffic Monitors are set to ``ONLINE``, the optimistic quorum can be employed by setting the ``peer_optimistic_quorum_min`` property in ``traffic_monitor.cfg`` to a value greater than zero. This value represents the minimum number of peers that must be available in order to participate in the `Optimistic Health Protocol`_. If Traffic Monitor detects that the number of available peers is less than this number, Traffic Monitor withdraws itself from participation in the health protocol by serving 503s for cache health state calls until connectivity is restored.

The optimistic quorum prevents invalid state propagation caused by a Traffic Monitor losing connectivity to the network and consequently marking all peers and caches as unavailable. When connectivity is restored, a race between peering recovery and polling from Traffic Routers begins. If Traffic Router were to poll a Traffic Monitor that has no available peers and optimistic quorum is not enabled or cannot be used (i.e.: too few Traffic Monitors), the Traffic Monitor will serve its local state only until peer connectivity is restored. If Traffic Router polls the Traffic Monitor when in this state, that is, prior to regaining peering, negative cache states caused by the lack of connectivity would be consumed and directly impact which caches are available for consideration for routing, until the Traffic Router polls a Traffic Monitor that has good state, or peering is restored. For this reason, it is recommended to run a minimum of three Traffic Monitors, with ``peer_optimistic_quorum_min`` set to a value of 1 or greater. Note that this value cannot exceed the number of peers of any given Traffic Monitor; that is, a value of 2 is the maximum value that can be used when three Traffic Monitors are in use. If this number exceeds the number of peers, the Traffic Monitor will always serve 503s and an error will be logged.

Protocol Engagement
-------------------
Short polling intervals of both the :term:`cache servers` and Traffic Monitor peers help to reduce customer impact of outages. It is not uncommon for a :term:`cache server` to be marked unavailable by Traffic Monitor - in fact, it is business as usual for many CDNs. Should a widely requested video asset cause a single :term:`cache server` to get close to its interface capacity, the Health Protocol will "kick in," and Traffic Monitor marks the :term:`cache server` as unavailable. New clients want to see the same asset, and now :ref:`tr-overview` will send these customers to another :term:`cache server` in the same :term:`Cache Group`. The load is now shared between the two :term:`cache servers`. As clients finish watching the asset on the overloaded :term:`cache server`, it will drop below the threshold and gets marked available again, and new clients will begin to be directed to it once more. It is less common for a :term:`Delivery Service` to be marked unavailable by Traffic Monitor. The :term:`Delivery Service` thresholds are usually used for overflow situations at extreme peaks to protect other :term:`Delivery Services` in the CDN from being impacted.
