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

.. _reference-label-tc-tr:

.. |arrow| image:: fwda.png

.. index::
  Traffic Router - Overview

Traffic Router
==============
Traffic Router's function is to send clients to the most optimal cache. Optimal in this case is based on a number of factors:

* Distance between the cache and the client (not necessarily measured in meters, but quite often in layer 3 network hops). Less network distance between the client and cache yields better performance, and lower network load. Traffic Router helps clients connect to the best performing cache for their location at the lowest network cost.

* Availability of caches and health / load on the caches. A common issue in Internet and television distribution scenarios is having many clients attempting to retrieve the same content at the same time. Traffic Router helps clients route around overloaded or down caches.

* Availability of content on a particular cache. Reusing of content through cache HITs is the most important performance gain a CDN can offer. Traffic Router sends clients to the cache that is most likely to already have the desired content.

Traffic routing options are often configured at the Delivery Service level. 

|


.. _rl-ds:

|arrow| Delivery Service
------------------------
  As discussed in the basic concepts section, the EDGE caches are configured as reverse proxies, and the Traffic Control CDN looks from the outside as a very large reverse proxy. Delivery Services are often referred to a reverse proxy remap rule. In most cases, a Delivery Service is a one to one mapping to a FQDN that is used as a hostname to deliver the content. Many options and settings regarding how to optimize the content delivery, which is configurable on a Delivery Service basis. Some examples of these Delivery Service settings are:

  * Cache in RAM, cache on disk, or do not cache at all.
  * Use DNS or HTTP Content routing (see below).
  * Limits on transactions per second and bandwidth.
  * Protocol (http or https).
  * Token based authentication settings. 
  * Header rewrite rules.

  Since Traffic Control version 2.1 deliveryservices can optionally be linked to a :ref:`rl-profile`, and have parameters associated with them. The first feature that uses deliveryservice parameters is the :ref:`rl-multi-site-origin` configuration.
  Delivery Services are also for use in allowing multi-tenants to coexist in the Traffic Control CDN without interfering with each other, and to keep information about their content separated. 

|

.. _rl-localization:

|arrow| Localization
--------------------
  Traffic Router uses a JSON input file called the *coverage zone map* to determine what *cachegroup* is closest to the client. If the client IP address is not in this coverage zone map, it falls back to *geo*, using the maxmind database to find the client's location, and the geo coordinates from Traffic Ops for the cachegroup.

|

Traffic Router is inserted into the HTTP retrieval process by making it DNS authoritative for the domain of the CDN delivery service. In the example of the reverse proxy, the client was given the ``http://www-origin-cache.cdn.com/foo/bar/fun.html`` url. In a Traffic Control CDN, URLs start with a routing name, which is configurable per-Delivery Service, e.g. ``http://foo.mydeliveryservice.cdn.com/fun/example.html`` with the chosen routing name ``foo``.

|

.. index::
  Content Routing

.. _rl-dns-cr:

|arrow| DNS Content Routing
---------------------------
  For a DNS delivery service the client might receive a URL such as ``http://foo.dsname.cdn.com/fun/example.html``. When the LDNS server is resolving this ``foo.dsname.cdn.com`` hostname to an IP address, it ends at Traffic Router because it is the authoritative DNS server for ``cdn.com`` and the domains below it, and subsequently responds with a list of IP addresses from the eligible caches based on the location of the LDNS server. When responding, Traffic Router does not know the actual client IP address or the path that the client is going to request. The decision on what cache IP address (or list of cache IP addresses) to return is solely based on the location of the LDNS server and the health of the caches. The client then connects to port 80 on the cache, and sends the ``Host: foo.dsname.cdn.com`` header. The configuration of the cache includes the remap rule ``http://foo.dsname.cdn.com http://origin.dsname.com`` to map the routed name to an origin hostname.

|

.. _rl-http-cr:

|arrow| HTTP Content Routing
----------------------------
  For an HTTP delivery service the client might receive a URL such as ``http://bar.dsname.cdn.com/fun/example.html``. The LDNS server resolves this ``bar.dsname.cdn.com`` to an IP address, but in this case Traffic Router returns its own IP address. The client opens a connection to port 80 on the Traffic Router's IP address, and sends: ::

    GET /fun/example.html HTTP/1.1
    Host: bar.dsname.cdn.com

  Traffic Router uses an HTTP 302 to redirect the client to the best cache. For example: ::

    HTTP/1.1 302 Moved Temporarily
    Server: Apache-Coyote/1.1
    Location: http://atsec-nyc-02.dsname.cdn.com/fun/example.html
    Content-Length: 0
    Date: Tue, 13 Jan 2015 20:01:41 GMT

  The information Traffic Router can consider when selecting a cache in this case is much better:

  * The client's IP address (the other side of the socket).
  * The URL path the client is requesting, excluding query string.
  * All HTTP 1.1 headers.

  The client follows the redirect and performs a DNS request for the IP address for ``atsec-nyc-02.dsname.cdn.com``, and normal HTTP steps follow, except the sending of the Host: header when connected to the cache is ``Host: atsec-nyc-02.dsname.cdn``, and the configuration of the cache includes the remap rule (e.g.``http://atsec-nyc-02.dsname.cdn  http://origin.dsname.com``).

  Traffic Router sends all requests for the same path in a delivery service to the same cache in a cache group using consistent hashing, in this case all caches in a cache group are not carrying the same content, and there is a much larger combined cache in the cache group. 

In many cases DNS content routing is the best possible option, especially in cases where the client is receiving small objects from the CDN like images and web pages. 

Traffic Router is redundant and horizontally scalable by adding more instances into the DNS hierarchy using NS records.

