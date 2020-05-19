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

.. _tr-overview:

**************
Traffic Router
**************
Traffic Router's function is to send clients to the most optimal :term:`cache server`. 'Optimal' in this case is based on a number of factors:

* Distance between the :term:`cache server` and the client (not necessarily measured in physical distance, but quite often in layer 3 network hops). Less network distance between the client and :term:`cache server` yields better performance and lower network load. Traffic Router helps clients connect to the best-performing :term:`cache server` for their location at the lowest network cost.

* Availability of :term:`cache servers` and the system processing/network load on the :term:`cache servers`. A common issue in Internet and television distribution scenarios is having many clients attempting to retrieve the same content at the same time. Traffic Router helps clients route around overloaded or purposely disabled :term:`cache servers`.

* Availability of content on a particular :term:`cache server`. Reusing of content through "cache hits" is the most important performance gain a CDN can offer. Traffic Router sends clients to the :term:`cache server` that is most likely to already have the desired content.

Traffic routing options are often configured at the :term:`Delivery Service` level.

.. _dns-cr:

DNS Content Routing
===================
For a DNS :term:`Delivery Service` the client might receive a URL such as ``http://video.demo1.mycdn.ciab.test/``. When the :abbr:`LDNS (Local Domain Name Server)` is resolving this ``video.demo1.mycdn.ciab.test`` hostname to an IP address, it ends at Traffic Router because it is the authoritative DNS server for ``mycdn.ciab.test`` and the domains below it, and subsequently responds with a list of IP addresses from the eligible :term:`cache servers` based on the location of the :abbr:`LDNS (Local Domain Name Server)`. When responding, Traffic Router does not know the actual client IP address or the path that the client is going to request. The decision on what :term:`cache server` IP address (or list of :term:`cache server` IP addresses) to return is solely based on the location of the :abbr:`LDNS (Local Domain Name Server)` and the health of the :term:`cache servers`. The client then connects to port 80 (HTTP) or port 443 (HTTPS) on the :term:`cache server`, and sends the ``Host: video.demo1.mycdn.ciab.test`` header. The configuration of the :term:`cache server` includes the "remap rule" ``http://video.demo1.mycdn.ciab.test http://origin.infra.ciab.test`` to map the routed name to an :term:`Origin` hostname.

.. _http-cr:

HTTP Content Routing
====================
For an HTTP :term:`Delivery Service` the client might receive a URL such as ``http://video.demo1.mycdn.ciab.test/``. The :abbr:`LDNS (Local Domain Name Server)` resolves this ``video.demo1.mycdn.ciab.test`` to an IP address, but in this case Traffic Router returns its own IP address. The client opens a connection to port 80 (HTTP) or port 443 (HTTPS) on the Traffic Router's IP address, and sends its request.

.. code-block:: http
	:caption: Example Client Request to Traffic Router

	GET / HTTP/1.1
	Host: video.demo1.mycdn.ciab.test
	Accept: */*

Traffic Router uses an HTTP ``302 Found`` response to redirect the client to the best :term:`cache server`.

.. code-block:: http
	:caption: Traffic Router Redirect to Edge-tier :term:`cache server`

	HTTP/1.1 302 Found
	Location: http://edge.demo1.mycdn.ciab.test/
	Content-Length: 0
	Date: Tue, 13 Jan 2015 20:01:41 GMT

In this case Traffic Router has access to more information when selecting a :term:`cache server` because it has a full HTTP request instead of just a hostname. Traffic Router can be configured to select a :term:`cache server` based on any of the following parts of the HTTP request:

* The client's IP address.
* The URL the client is requesting.
* All HTTP/1.1 headers.

The client follows the redirect and performs a DNS request for the IP address for ``edge.demo1.mycdn.ciab.test``, and normal HTTP steps follow, except the sending of the Host: header when connected to the cache is ``Host: edge.demo1.mycdn.ciab.test``, and the configuration of the :term:`cache server` includes the "remap rule" (e.g. ``http://edge.demo1.mycdn.ciab.test http://origin.infra.ciab.test``). Traffic Router sends all requests for the same path in a :term:`Delivery Service` to the same :term:`cache server` in a :term:`Cache Group` using consistent hashing, in this case all :term:`cache servers` in a :term:`Cache Group` are not carrying the same content, and there is a much larger combined cache in the :term:`Cache Group`. In many cases DNS content routing is the best possible option, especially in cases where the client is receiving small objects from the CDN like images and web pages. Traffic Router is redundant and horizontally scalable by adding more instances into the DNS hierarchy using NS records.
