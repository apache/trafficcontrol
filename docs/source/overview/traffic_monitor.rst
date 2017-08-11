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

.. _reference-label-tc-tm:

.. index::
	Traffic Monitor - Overview

.. |arrow| image:: fwda.png

Traffic Monitor
===============
Traffic Monitor is an HTTP service that monitors the caches in a CDN for a variety of metrics. These metrics are for use in determining the overall health of a given cache and the related delivery services. A given CDN can operate a number of Traffic Monitors, from a number of geographically diverse locations, to prevent false positives caused by network problems at a given site.

Traffic Monitors operate independently, but use the state of other Traffic Monitors in conjunction with their own state, to provide a consistent view of CDN cache health to upstream applications such as Traffic Router. Health Protocol governs the cache and Delivery Service availability. 

Traffic Monitor provides a view into CDN health using several RESTful JSON endpoints, which are consumed by other Traffic Monitors and upstream components such as Traffic Router. Traffic Monitor is also responsible for serving the overall CDN configuration to Traffic Router, which ensures that the configuration of these two critical components remain synchronized as operational and health related changes propagate through the CDN.


.. _rl-astats:

|arrow| Cache Monitoring
-------------------------
	Traffic Monitor polls all caches configured with a status of ``REPORTED`` or ``ADMIN_DOWN`` at an interval specified as a configuration parameter in Traffic Ops. If the cache is set to ``ADMIN_DOWN`` it is marked as unavailable but still polled for availability and statistics. If the cache is explicitly configured with a status of ``ONLINE`` or ``OFFLINE``, it is not polled by Traffic Monitor and presented to Traffic Router as configured, regardless of actual availability.

	Traffic Monitor makes HTTP requests at regular intervals to a special URL on each EDGE cache and consumes the JSON output. The special URL is a plugin running on the Apache Traffic Server (ATS) caches called astats, which is restricted to Traffic Monitor only. The astats plugin provides insight into application and system performance, such as:

	- Throughput (e.g. bytes in, bytes out, etc).
	- Transactions (e.g. number of 2xx, 3xx, 4xx responses, etc).
	- Connections (e.g. from clients, to parents, origins, etc).
	- Cache performance (e.g.: hits, misses, refreshes, etc).
	- Storage performance (e.g.: writes, reads, frags, directories, etc).
	- System performance (e.g: load average, network interface throughput, etc).

	Many of the application level statistics are available at the global or aggregate level, some at the Delivery Service (remap rule) level. Traffic Monitor uses the system level performance to determine the overall health of the cache by evaluating network throughput and load against values configured in Traffic Ops. Traffic Monitor also uses throughput and transaction statistics at the remap rule level to determine Delivery Service health.

If astats is unavailable due to a network related issue or the system statistics have exceeded the configured thresholds, Traffic Monitor will mark the cache as unavailable. If the delivery service statistics exceed the configured thresholds, the delivery service is marked as unavailable, and Traffic Router will start sending clients to the overflow destinations for that delivery service, but the cache remains available to serve other content, 

.. seealso:: For more information on ATS Statistics, see the `ATS documentation <https://docs.trafficserver.apache.org/en/latest/index.html>`_

.. _rl-health-proto:

|arrow| Health Protocol 
-----------------------
	Redundant Traffic Monitor servers operate independently from each other but take the state of other Traffic Monitors into account when asked for health state information. In the above overview of cache monitoring, the behavior of Traffic Monitor pertains only to how an individual instance detects and handles failures. The Health Protocol adds another dimension to the health state of the CDN by merging the states of all Traffic Monitors into one, and then taking the *optimistic* approach when dealing with a cache or Delivery Service that might have been marked as unavailable by this particular instance or a peer instance of Traffic Monitor.

	Upon startup or configuration change in Traffic Ops, in addition to caches, Traffic Monitor begins polling its peer Traffic Monitors whose state is set to ``ONLINE``. Each ``ONLINE`` Traffic Monitor polls all of its peers at a configurable interval and saves the peer's state for later use. When polling its peers, Traffic Monitor asks for the raw health state from each respective peer, which is strictly that instance's view of the CDN's health. When any ``ONLINE`` Traffic Monitor is asked for CDN health by an upstream component, such as Traffic Router, the component gets the health protocol influenced version of CDN health (non-raw view).

	In operation of the health protocol, Traffic Monitor takes all health states from all peers, including the locally known health state, and serves an optimistic outlook to the requesting client. This means that, for example, if three of the four Traffic Monitors see a given cache or Delivery Service as exceeding its thresholds and unavailable, it is still considered available.  Only if all Traffic Monitors agree that the given object is unavailable is that state propagated to upstream components. This optimistic approach to the Health Protocol is counter to the "fail fast" philosophy, but serves well for large networks with complicated geography and or routing. The optimistic Health Protocol allows network failures or latency to occur without affecting overall traffic routing, as Traffic Monitors can and do have a different view of the network when deployed in geographically diverse locations. Short polling intervals of both the caches and Traffic Monitor peers help to reduce customer impact of outages.

It is not uncommon for a cache to be marked unavailable by Traffic Monitor - in fact, it is business as usual for many CDNs. A hot video asset may cause a single cache (say cache-03) to get close to it's interface capacity, the health protocol "kicks in", and Traffic Monitor marks cache-03 as unavailable. New clients want to see the same asset, and now, Traffic Router will send these customers to another cache (say cache-01) in the same cachegroup. The load is now shared between cache-01 and cache-03. As clients finish watching the asset on cache-03, it will drop below the threshold and gets marked available again, and new clients will now go back to cache-03 again. 

It is less common for a delivery service to be marked unavailable by Traffic Monitor, the delivery service thresholds are usually used for overflow situations at extreme peaks to protect other delivery services in the CDN from getting impacted.

