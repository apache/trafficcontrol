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

.. _glossary:

Glossary
========

.. glossary::

	302 content routing
		:ref:`rl-http-cr`.

	astats (stats_over_http)
		An ATS plugin that allows you to monitor vitals of the ATS server. See :ref:`rl-astats`.

	cache
		A caching proxy server. See :ref:`rl-caching_proxy`.

	cachegroup
		A group of caches that together create a combined larger cache using consistent hashing. See :ref:`rl-cachegroup`.

	consistent hashing
		See `the Wikipedia article <http://en.wikipedia.org/wiki/Consistent_hashing>`_; Traffic Control uses consistent hashing when using :ref:`rl-http-cr` for the edge tier and when selecting parents in the mid tier.

	content routing
		Directing clients (or client systems) to a particular location or device in a location for optimal delivery of content See also :ref:`rl-http-cr` and :ref:`rl-dns-cr`.

	coverage zone map
		The coverage zone map (czm) or coverage zone file (zcf) is a file that maps network prefixes to cachegroups. See :ref:`rl-localization`.

	delivery service
		A grouping of content in the CDN, usually a determined by the URL hostname. See :ref:`rl-ds`.

	edge (tier or cache)
		Closest to the client or end-user. The edge tier is the tier that serves the client, edge caches are caches in the edge tier. In a Traffic Control CDN the basic function of the edge cache is that of a :ref:`rl-rev-proxy`.  See also :ref:`rl-cachegroup`.

	(traffic ops) extension 
		Using *extensions*, Traffic Ops be extended to use proprietary checks or monitoring sources. See :ref:`rl-trops-ext`.

	forward proxy
		A proxy that works that acts like it is the client to the origin. See :ref:`rl-fwd-proxy`.

	geo localization or geo routing
		Localizing clients to the nearest caches using a geo database like the one from Maxmind. 

 	health protocol
 		The protocol to monitor the health of all the caches. See :ref:`rl-health-proto`. 

 	localization
 		Finding location on the network, or on planet earth. See :ref:`rl-localization`.

	mid (tier or cache)
		The tier above the edge tier. The mid tier does not directly serves the end-user and is used as an additional layer between the edge and the origin. In a Traffic Control CDN the basic function of the mid cache is that of a :ref:`rl-fwd-proxy`. See also :ref:`rl-cachegroup`.

	origin
		The source of content for the CDN. Usually a redundant HTTP/1.1 webserver.

	parent (cache or cachegroup)
		The (group of) cache(s) in the higher tier.  See :ref:`rl-cachegroup`.

	profile
		A group of settings (parameters) that will be applied to a server. See :ref:`rl-profile`.

	reverse proxy
		A proxy that acts like it is the origin to the client. See :ref:`rl-rev-proxy`.



