.. 
.. Copyright 2015 Comcast Cable Communications Management, LLC
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

.. _reference-label-tc-tm:

.. index::
	Traffic Monitor - Overview

.. |arrow| image:: fwda.png

Traffic Monitor
===============
Traffic Monitor is a Java/Tomcat application that monitors the caches in a CDN for a variety of metrics. These metrics are for use in determining the overall health of a given cache and the related delivery services. A given CDN can operate a number of Traffic Monitors, from a number of geographically diverse locations, to prevent false positives caused by network problems at a given site.

Traffic Monitors operate independently, but use the state of other Traffic Monitors in conjunction with their own state, to provide a consistent view of CDN cache health to upstream applications such as Traffic Router. Health Protocol governs the cache and Delivery Service availability. 

Traffic Monitor provides a view into CDN health using several RESTful JSON endpoints, which are consumed by other Traffic Monitors and upstream components such as Traffic Router. Traffic Monitor is also responsible for serving the overall CDN configuration to Traffic Router, which ensures that the configuration of these two critical components remain synchronized as operational and health related changes propagate through the CDN.


|arrow| Cache Monitoring
-------------------------
Traffic Monitor currently polls all caches configured with a status of ``REPORTED`` or ``ADMIN_DOWN`` at an interval specified as a configuration parameter in Traffic Ops. If the cache is set to ``ADMIN_DOWN`` it is marked as unavailable but still polled for availability and statistics. If the cache is configured with a status of ``ONLINE`` or ``OFFLINE``, it is not polled by Traffic Monitor and assumes that its current state matches its configured status.

Traffic Monitor makes HTTP requests at regular intervals to a special URL on each EDGE cache and consumes the JSON output. The special URL is a plugin running on the Apache Traffic Server (ATS) caches called astats, which is restricted to Traffic Monitor only. The astats plugin provides insight into application and system performance, such as, but not limited to:

- Throughput (e.g. bytes in, bytes out, etc).
- Transactions (e.g. number of 2xx, 3xx, 4xx responses, etc).
- Connections (e.g. from clients, to parents, origins, etc).
- Cache performance (e.g.: hits, misses, refreshes, etc).
- Storage performance (e.g.: writes, reads, frags, directories, etc).
- System performance (e.g: load average, network interface throughput, etc).

Many of the application level statistics are available at the global, aggregate level, or at the Delivery Service (remap rule) level. Traffic Monitor uses the system level performance to determine the overall health of the cache by evaluating network throughput and load against values configured in Traffic Ops. Traffic Monitor also uses throughput and transaction statistics at the remap rule level to determine Delivery Service health.

If astats is unavailable due to a network related issue, or the system or Delivery Service statistics have exceeded the configured thresholds; any of these scenarios disable the respective object on the CDN. If astats is unavailable or the system's thresholds are exceeded, the entire cache is unavailable. Exceeding thresholds of a Delivery Service causes its disablement across all caches in the CDN. If all caches are unavailable for a given Delivery Service, or if the Delivery Service is unavailable due to exceeding thresholds, Traffic Router stops routing traffic to the Delivery Service.

.. seealso:: For more information on ATS Statistics, see the `ATS documentation <https://docs.trafficserver.apache.org/en/latest/index.html>`_

.. index::
	Health Protocol

|arrow| Health Protocol 
-----------------------
Traffic Monitors operate independently but take the state of other Traffic Monitors into account when asked for health state information. So far, the behaviors of Traffic Monitor pertain only to how an individual instance detects and handles failures. The Health Protocol adds another dimension to the health state of the CDN by merging the states of all Traffic Monitors into one, then taking the *optimistic* approach when dealing with a cache or Delivery Service that might have been marked as unavailable by this particular instance or a peer instance of Traffic Monitor.

.. that last sentence doesn't make sense to me. It's verbose and yet I don't know what it's trying to convey.

Upon startup or configuration change in Traffic Ops, in addition to caches, Traffic Monitor begins polling its peer Traffic Monitors whose state is set to ``ONLINE``. Each ``ONLINE`` Traffic Monitor polls all of its peers at a configurable interval and saves the peer's state for later use. When polling its peers, Traffic Monitor asks for the raw health state from each respective peer, which is strictly that instance's view of the CDN's health. When any ``ONLINE`` Traffic Monitor is asked for CDN health by an upstream component, such as Traffic Router, the component gets the health protocol influenced version of CDN health (non-raw view).

In operation of the health protocol, Traffic Monitor takes all health states from all peers, including the locally known health state, and serves an optimistic outlook to the requesting client. This means that, for example, if three of the four Traffic Monitors see a given cache or Delivery Service as exceeding its thresholds and unavailable, it is still considered available.  Only until all Traffic Monitors agree that the given object is unavailable is that state propagated to upstream components. This optimistic approach to the Health Protocol is counter to the "fail fast" philosophy, but serves well for large networks with complicated geography and or routing. The optimistic Health Protocol allows network failures or latency to occur without affecting overall traffic routing, as Traffic Monitors can and do have a different view of the network when deployed in geographically diverse locations. Short polling intervals of both the caches and Traffic Monitor peers help to reduce customer impact of outages.