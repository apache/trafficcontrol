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

********
Glossary
********

.. glossary::
	:sorted:

	302 content routing
		:ref:`http-cr`.

	astats (stats_over_http)
		An :abbr:`ATS (Apache Traffic Server)` plugin that allows you to monitor vitals of the :abbr:`ATS (Apache Traffic Server)` server. See :ref:`astats`.

	cache server
	cache servers
		The main function of a CDN is to proxy requests from clients to :term:`origin servers` and cache the results. To proxy, in the CDN context, is to obtain content using HTTP from an :term:`origin server` on behalf of a client. To cache is to store the results so they can be reused when other clients are requesting the same content. There are three types of proxies in use on the Internet today:

		- :term:`Reverse Proxy`: Used by Traffic Control for Edge-tier :dfn:`cache servers`.
		- :term:`Forward Proxy`: Used by Traffic Control for Mid-tier :dfn:`cache servers`.
		- Transparent Proxy: These are not used by Traffic Control. If you are interested you can learn more about transparent proxies on `wikipedia <http://en.wikipedia.org/wiki/Proxy_server#Transparent_proxy>`_.

	Cache Group
	Cache Groups
		A group of caching HTTP proxy servers that together create a combined larger cache using consistent hashing. Traffic Router treats all servers in a :dfn:`Cache Group` as though they are in the  same :term:`Physical Location`, though they are in fact only in the same general area. A :dfn:`Cache Group` has one single set of geographical coordinates even if the :term:`cache server`\ s that make up the :dfn:`Cache Group` are actually in :term:`Physical Location`\ s. The :term:`cache server`\ s in a :dfn:`Cache Group` are not aware of the other :term:`cache server`\ s in the group - there is no clustering software or communications between :term:`cache server`\ s in a :dfn:`Cache Group`.

  		There are two basic types of :dfn:`Cache Groups`: EDGE_LOC and MID_LOC ("LOC" being short for "location" - a holdover from when :dfn:`Cache Groups` were called "Cache Locations). Traffic Control is a two-tiered system, where the clients get directed to the Edge-tier (EDGE_LOC) :dfn:`Cache Group`. On cache miss, the :term:`cache server` in the Edge-tier :dfn:`Cache Group` obtains content from a Mid-tier (MID_LOC) :dfn:`Cache Group`, rather than the origin, which is shared with multiple Edge-tier :dfn:`Cache Groups`. Edge-tier :dfn:`Cache Groups` are configured to have a single "parent" :dfn:`Cache Group`, but in general Mid-tier :dfn:`Cache Groups` have many "children".

		..  Note:: Often the Edge-tier to Mid-tier relationship is based on network distance, and does not necessarily match the geographic distance.

		.. seealso:: A :dfn:`Cache Group` serves a particular part of the network as defined in the :term:`Coverage Zone File` (or :term:`Deep Coverage Zone File`, when applicable).

		Consider the example CDN in :numref:`fig-cg_hierarchy`. Here some country/province/region has been divided into quarters: Northeast, Southeast, Northwest, and Southwest. The arrows in the diagram indicate the flow of requests. If a client in the Northwest, for example, were to make a request to the :term:`Delivery Service`, it would first be directed to some :term:`cache server` in the "Northwest" Edge-tier :dfn:`Cache Group`. Should the requested content not be in cache, the Edge-tier server will select a parent from the "West" :dfn:`Cache Group` and pass the request up, caching the result for future use. All Mid-tier :dfn:`Cache Groups` (usually) answer to a single :term:`origin` that provides canonical content. If requested content is not in the Mid-tier cache, then the request will be passed up to the :term:`origin` and the result cached.

		.. _fig-cg_hierarchy:

		.. figure:: ./cg_hierarchy.*
			:align: center
			:width: 60%
			:alt: An illustration of Cache Group hierarchy

			An example CDN that shows the hierarchy between four Edge-tier :dfn:`Cache Groups`, two Mid-tier :dfn:`Cache Groups`, and one Origin

	content routing
		Directing clients (or client systems) to a particular location or device in a location for optimal delivery of content See also :ref:`http-cr` and :ref:`dns-cr`.

	Coverage Zone File
	Coverage Zone Map
		The :abbr:`CZM (Coverage Zone Map)` or :abbr:`CZF (Coverage Zone File)` is a file that maps network prefixes to :term:`Cache Groups`. Traffic Router uses the :abbr:`CZM (Coverage Zone Map)` to determine what :term:`Cache Group` is closest to the client. If the client IP address is not in this :abbr:`CZM (Coverage Zone Map)`, it falls back to geographic mapping, using a `MaxMind GeoIP2 database <https://www.maxmind.com/en/geoip2-databases>`_ to find the client's location, and the geographic coordinates from Traffic Ops for the :term:`Cache Group`. Traffic Router is inserted into the HTTP retrieval process by making it the authoritative DNS server for the domain of the CDN :term:`Delivery Service`. In the example of the :term:`reverse proxy`, the client was given the ``http://www-origin-cache.cdn.com/foo/bar/fun.html`` URL. In a Traffic Control CDN, URLs start with a routing name, which is configurable per-:term:`Delivery Service`, e.g. ``http://foo.mydeliveryservice.cdn.com/fun/example.html`` with the chosen routing name ``foo``.

		.. code-block:: json
			:caption: Example Coverage Zone File

			{ "coverageZones": {
				"cache-group-01": {
					"network6": [
						"1234:5678::/64",
						"1234:5679::/64"
					],
					"network": [
						"192.168.8.0/24",
						"192.168.9.0/24"
					]
				}
			}}


	Deep Coverage Zone File
	Deep Coverage Zone Map
		The :abbr:`DCZF (Deep Coverage Zone File)` or :abbr:`DCZM (Deep Coverage Zone Map)` maps network prefixes to "locations" - almost like the :term:`Coverage Zone File`. Location names must be unique, and within the file are simply used to group :term:`Edge-tier cache servers`. When a mapping is performed by Traffic Router, it will only look in the :abbr:`DCZF (Deep Coverage Zone File)` if the :term:`Delivery Service` to which a client is being directed makes use of :ref:`ds-deep-caching`. If the client's IP address cannot be matched by entries in this file, Traffic Router will first fall back to the regular :term:`Coverage Zone File`. Then, failing that, it will perform geographic mapping using a database provided by the :term:`Delivery Service`'s :ref:`ds-geo-provider`.

		.. code-block:: json
			:caption: Example Deep Coverage Zone File

			{ "coverageZones": {
				"cache-group-01": {
					"network6": [
						"1234:5678::/64",
						"1234:5679::/64"
					],
					"network": [
						"192.168.8.0/24",
						"192.168.9.0/24"
					],
					"caches": [
						"edge"
					]
				}
			}}

	Delivery Service
	Delivery Services
		:dfn:`Delivery Services` are often referred to as a :term:`reverse proxy` "remap rule" that exists on Edge-tier :term:`cache servers`. In most cases, a :dfn:`Delivery Service` is a one-to-one mapping to an :abbr:`FQDN (Fully Qualified Domain Name)` that is used as a hostname to deliver the content. Many options and settings regarding how to optimize the content delivery exist, which are configurable on a :dfn:`Delivery Service` basis. Some examples of these :dfn:`Delivery Service`\ settings are:

		* Cache in RAM, cache on disk, or do not cache at all.
		* Use DNS or HTTP Content routing.
		* Limits on transactions per second and bandwidth.
		* Protocol (HTTP or HTTPS).
		* Token-based authentication settings.
		* Header rewrite rules.

		Since Traffic Control version 2.1, :dfn:`Delivery Services` can optionally be linked to a :term:`Profile`, and have :term:`Parameters` associated with them. One example of a feature that uses :dfn:`Delivery Service` :term:`Parameters` is the :ref:`ds-multi-site-origin` configuration. :dfn:`Delivery Services` are also for use in allowing multiple :term:`Tenants` to coexist in a Traffic Control CDN without interfering with each other, and to keep information about their content separated.

		.. seealso:: See :ref:`delivery-services` for a more in-depth explanation of :dfn:`Delivery Services`.

	Division
	Divisions
		A group of :term:`Regions`.

	Edge
	Edge-tier
	Edge-tier cache
	Edge-tier caches
	Edge-tier cache server
	Edge-tier cache servers
		Closest to the client or end-user. The edge tier is the tier that serves the client, edge caches are caches in the edge tier. In a Traffic Control CDN the basic function of the edge cache is that of a :term:`reverse proxy`.

	Federation
	Federations
		:dfn:`Federations` allow for other ("federated") CDNs (e.g. at a different :abbr:`ISP (Internet Service Provider)`) to add a list of DNS resolvers and an :abbr:`FQDN (Fully Qualified Domain Name)` to be used in a DNS CNAME record for a :term:`Delivery Service`. When a request is made from one of the federated CDN's clients, Traffic Router will return the CNAME record configured from the federation mapping. This allows the federated CDN to serve the content without the content provider changing the URL, or having to manage multiple URLs. For example, if the external CDN was actually another :abbr:`ATC (Apache Traffic Control)`-managed CDN, then a federation mapping to direct clients toward it should use the :abbr:`FQDN (Fully Qualified Domain Name)` of a :term:`Delivery Service` on the external CDN.

		Federations only have meaning to DNS-routed :term:`Delivery Services` - HTTP-routed Delivery services should instead treat the external :abbr:`FQDN (Fully Qualified Domain Name)` as an :term:`origin` to achieve the same effect.

		.. seealso:: Federations are currently only manageable by directly using the :ref:`to-api`. The endpoints related to federations are :ref:`to-api-federations`, :ref:`to-api-federation_resolvers`, :ref:`to-api-federation_resolvers-id`, :ref:`to-api-federations-id-deliveryservices`, :ref:`to-api-federations-id-deliveryservices-id`, :ref:`to-api-federations-id-federation_resolvers`, :ref:`to-api-federations-id-users`, and :ref:`to-api-federations-id-users-id`.

	forward proxy
	forward proxies
		A forward proxy acts on behalf of the client such that the :term:`origin server` is (potentially) unaware of the proxy's existence. All Mid-tier :term:`cache server`\ s in a Traffic Control based CDN are :dfn:`forward proxies`. In a :dfn:`forward proxy` scenario, the client is explicitly configured to use the the proxy's IP address and port as a :dfn:`forward proxy`. The client always connects to the :dfn:`forward proxy` for content. The content provider does not have to change the URL the client obtains, and is (potentially) unaware of the proxy in the middle.

		..  seealso:: `ATS documentation on forward proxy <https://docs.trafficserver.apache.org/en/latest/admin/forward-proxy.en.html>`_.

		If a client uses a :dfn:`forward proxy` to request the URL ``http://www.origin.com/foo/bar/fun.html`` the resulting chain of events follows.

		#. To retrieve ``http://www.origin.com/foo/bar/fun.html``, the client sends an HTTP request to the :dfn:`forward proxy`.

			.. code-block:: http
				:caption: Client Requests Content from its :dfn:`Forward Proxy`

				GET http://www.origin.com/foo/bar/fun.html HTTP/1.1
				Host: www.origin.com

			..  Note:: In this case, the client requests the entire URL instead of just the path as is the case when using a :term:`reverse proxy` or when requesting content directly from the :term:`origin server`.

		#. The proxy verifies whether the response for ``http://www-origin-cache.cdn.com/foo/bar/fun.html`` is already in the cache. If it is not in the cache:

			#. The proxy sends the HTTP request to the :term:`origin`.

				.. code-block:: http
					:caption: The :dfn:`Forward Proxy` Requests Content from the :term:`Origin Server`

					GET /foo/bar/fun.html HTTP/1.1
					Host: www.origin.com

			#. The :term:`origin server` responds with the requested content.

				.. code-block:: http
					:caption: The :term:`Origin Server`'s Response

					HTTP/1.1 200 OK
					Date: Sun, 14 Dec 2014 23:22:44 GMT
					Server: Apache/2.2.15 (Red Hat)
					Last-Modified: Sun, 14 Dec 2014 23:18:51 GMT
					ETag: "1aa008f-2d-50a3559482cc0"
					Content-Length: 45
					Connection: close
					Content-Type: text/html; charset=UTF-8

					<!DOCTYPE html><html><body>This is a fun file</body></html>


			#. The proxy sends this on to the client, optionally adding a ``Via:`` header to indicate that the request was serviced by proxy.

				.. code-block:: http
					:caption: The :dfn:`Forward Proxy`'s Response to the Client

					HTTP/1.1 200 OK
					Date: Sun, 14 Dec 2014 23:22:44 GMT
					Last-Modified: Sun, 14 Dec 2014 23:18:51 GMT
					ETag: "1aa008f-2d-50a3559482cc0"
					Content-Length: 45
					Connection: close
					Content-Type: text/html; charset=UTF-8
					Age: 0
					Via: http/1.1 cache01.cdn.kabletown.net (ApacheTrafficServer/4.2.1 [uScSsSfUpSeN:t cCSi p sS])
					Server: ATS/4.2.1

					<!DOCTYPE html><html><body>This is a fun file</body></html>


			If, however, the requested content *was* in the cache the proxy responds to the client with the previously retrieved result

			.. code-block:: http
				:caption: The :dfn:`Forward Proxy` Sends the Cached Response

				HTTP/1.1 200 OK
				Date: Sun, 14 Dec 2014 23:22:44 GMT
				Last-Modified: Sun, 14 Dec 2014 23:18:51 GMT
				ETag: "1aa008f-2d-50a3559482cc0"
				Content-Length: 45
				Connection: close
				Content-Type: text/html; charset=UTF-8
				Age: 99711
				Via: http/1.1 cache01.cdn.kabletown.net (ApacheTrafficServer/4.2.1 [uScSsSfUpSeN:t cCSi p sS])
				Server: ATS/4.2.1

				<!DOCTYPE html><html><body>This is a fun file</body></html>

	geo localization or geo routing
		Localizing clients to the nearest caches using a geo database like the one from Maxmind.

 	Health Protocol
 		The protocol to monitor the health of all the caches. See :ref:`health-proto`.

 	localization
 		Finding location on the network, or on planet earth

	Mid
	Mid-tier
	Mid-tier cache
	Mid-tier caches
	Mid-tier cache server
	Mid-tier cache servers
		The tier above the edge tier. The mid tier does not directly serves the end-user and is used as an additional layer between the edge and the :term:`origin`. In a Traffic Control CDN the basic function of the mid cache is that of a :term:`forward proxy`.

	origin
	origins
	origin server
	origin servers
		The source of content for the CDN. Usually a redundant HTTP/1.1 webserver.

	ORT
		The "Operational Readiness Test" script that stitches the configuration configured in Traffic Portal and generated by Traffic Ops into the :term:`cache server`\ s. See :ref:`traffic-ops-ort` for more information.

	Parameter
	Parameters
		Typically refers to a line in a configuration file, but in practice can represent any arbitrary configuration option

	parent
	parents
		The :dfn:`parent(s)` of a :term:`cache server` is/are the :term:`cache server`\ (s) belonging to either the "parent" or "secondary parent" :term:`Cache Group`\ (s) of the :term:`Cache Group` to which the :term:`cache server` belongs. For example, in general it is true that an :term:`Edge-tier cache server` has one or more :dfn:`parents` which are :term:`Mid-tier cache servers`.

	Physical Location
	Physical Locations
		A pair of geographic coordinates (latitude and longitude) that is used by :term:`Cache Group`\ s to define their location. This information is used by Traffic Router to route client traffic to the geographically nearest :term:`Cache Group`.

	Profile
	Profiles
		A :dfn:`Profile` is, most generally, a group of :term:`Parameter`\ s that will be applied to a server. :dfn:`Profiles` are typically re-used by all Edge-Tier :term:`cache server`\ s within a CDN or :term:`Cache Group`. A :dfn:`Profile` will, in addition to configuration :term:`Parameter`\ s, define the CDN to which a server belongs and the "Type" of the profile - which determines some behaviors of Traffic Control components. The allowed "Types" of :dfn:`Profiles` are **not** the same as :term:`Type`\ s, and are maintained as a PostgreSQL "Enum" in :file:`traffic_ops/app/db/migrations/20170205101432_cdn_table_domain_name.go`. The only allowed values are:

		UNK_PROFILE
			A catch-all type that can be assigned to anything without imbuing it with any special meaning or behavior
		TR_PROFILE
			A Traffic Router Profile.

			.. warning:: For legacy reasons, the names of Profiles of this type *must* begin with ``CCR_`` or ``TR_``. This is **not** enforced by the :ref:`to-api` or Traffic Portal, but certain Traffic Control operations/components expect this and will fail to work otherwise!

		TM_PROFILE
			A Traffic Monitor Profile.

			.. warning:: For legacy reasons, the names of Profiles of this type *must* begin with ``RASCAL_`` or ``TM_``. This is **not** enforced by the :ref:`to-api` or Traffic Portal, but certain Traffic Control operations/components expect this and will fail to work otherwise!

		TS_PROFILE
			A Traffic Stats Profile

			.. warning:: For legacy reasons, the names of Profiles of this type *must* be ``TRAFFIC_STATS``. This is **not** enforced by the :ref:`to-api` or Traffic Portal, but certain Traffic Control operations/components expect this and will fail to work otherwise!

		TP_PROFILE
			A Traffic Portal Profile

			.. warning:: For legacy reasons, the names of Profiles of this type *must* begin with ``TRAFFIC_PORTAL``. This is **not** enforced by the :ref:`to-api` or Traffic Portal, but certain Traffic Control operations/components expect this and will fail to work otherwise!

		INFLUXDB_PROFILE
			A Profile used with `InfluxDB <https://www.influxdata.com/>`_, which is used by Traffic Stats.

			.. warning:: For legacy reasons, the names of Profiles of this type *must* begin with ``INFLUXDB``. This is **not** enforced by the :ref:`to-api` or Traffic Portal, but certain Traffic Control operations/components expect this and will fail to work otherwise!

		RIAK_PROFILE
			A Profile used for each `Riak <http://basho.com/products/riak-kv/>`_ server in a Traffic Stats cluster.

			.. warning:: For legacy reasons, the names of Profiles of this type *must* begin with ``RIAK``. This is **not** enforced by the :ref:`to-api` or Traffic Portal, but certain Traffic Control operations/components expect this and will fail to work otherwise!

		SPLUNK_PROFILE

			A Profile meant to be used with `Splunk <https://www.splunk.com/>`_ servers.

			.. warning:: For legacy reasons, the names of Profiles of this type *must* begin with ``SPLUNK``. This is **not** enforced by the :ref:`to-api` or Traffic Portal, but certain Traffic Control operations/components expect this and will fail to work otherwise!

		ORG_PROFILE
			Origin Profile.

			.. warning:: For legacy reasons, the names of Profiles of this type *must* begin with ``MSO``, or contain either ``ORG`` or ``ORIGIN`` anywhere in the name. This is **not** enforced by the :ref:`to-api` or Traffic Portal, but certain Traffic Control operations/components expect this and will fail to work otherwise!

		KAFKA_PROFILE
			A Profile for `Kafka <https://kafka.apache.org/>`_ servers.

			.. warning:: For legacy reasons, the names of Profiles of this type *must* begin with ``KAFKA``. This is **not** enforced by the :ref:`to-api` or Traffic Portal, but certain Traffic Control operations/components expect this and will fail to work otherwise!

		LOGSTASH_PROFILE
			A Profile for `Logstash <https://www.elastic.co/products/logstash>`_ servers.

			.. warning:: For legacy reasons, the names of Profiles of this type *must* begin with ``LOGSTASH_``. This is **not** enforced by the :ref:`to-api` or Traffic Portal, but certain Traffic Control operations/components expect this and will fail to work otherwise!

		ES_PROFILE
			A Profile for `ElasticSearch <https://www.elastic.co/products/elasticsearch>`_ servers.

			.. warning:: For legacy reasons, the names of Profiles of this type *must* begin with ``ELASTICSEARCH``. This is **not** enforced by the :ref:`to-api` or Traffic Portal, but certain Traffic Control operations/components expect this and will fail to work otherwise!

		ATS_PROFILE
			A Profile that can be used with either an Edge-tier or Mid-tier :term:`cache server` (but not both, in general).

			.. warning:: For legacy reasons, the names of Profiles of this type *must* begin with ``EDGE`` or ``MID``. This is **not** enforced by the :ref:`to-api` or Traffic Portal, but certain Traffic Control operations/components expect this and will fail to work otherwise!


		.. tip:: A :dfn:`Profile` of the wrong type assigned to a Traffic Control component *will* (in general) cause it to function incorrectly, regardless of the :term:`Parameters` assigned to it.

		.. danger:: Nearly all of these :dfn:`Profile` types have strict naming requirements, and it may be noted that some of said requirements are prefixes ending with ``_``, while others are either not prefixes or do not end with ``_``. This is exactly true; some requirements **need** that ``_`` and some may or may not have it. It is our suggestion, therefore, that for the time being all prefixes use the ``_`` notation to separate words, so as to avoid causing headaches remembering when that matters and when it does not.

	Queue
	Queue Updates
	Queue Server Updates
		:dfn:`Queuing Updates` is an action that signals to various ATC components - most notably :term:`cache servers` - that any configuration changes that are pending are to be applied now. Specifically, Traffic Monitor and Traffic Router are updated through a CDN :term:`Snapshot`, and *not* :dfn:`Queued Updates`. In particular, :term:`ORT` will notice that the server on which it's running has new configuration, and will request the new configuration from Traffic Ops.

		Updates may be queued on a server-by-server basis (in Traffic Portal's :ref:`tp-configure-servers` view), a Cache Group-wide basis (in Traffic Portal's :ref:`tp-configure-cache-groups` view), or on a CDN-wide basis (in Traffic Portal's :ref:`tp-cdns` view). Usually using the CDN-wide version is easiest, and unless there are special circumstances, and/or the user really knows what he or she is doing, it is recommended that the full CDN-wide :dfn:`Queue Updates` be used.

		This is similar to taking a CDN :term:`Snapshot`, but this configuration change affects only servers, and not routing.

		That seems like a vague difference because it is - in general the rule to follow is that changes to :term:`Profiles` and :term:`Parameters` requires only updates be queued, changes to the assignments of :term:`cache servers` to :term:`Delivery Services` requires both a :term:`Snapshot` *and* a :dfn:`Queue Updates`, and changes to only a :term:`Delivery Service` itself (usually) entails a :term:`Snapshot` only. These aren't exhaustive rules, and a grasp of what changes require which action(s) will take time to form. In general, when doing both :dfn:`Queuing Updates` as well as taking a CDN :term:`Snapshot`, it is advisable to first :dfn:`Queue Updates` and *then* take the :term:`Snapshot`, as otherwise Traffic Router may route clients to :term:`Edge-tier cache servers` that are not equipped to service their request(s). However, when modifying the assignment(s) of :term:`cache servers` to one or more :term:`Delivery Services`, a :term:`Snapshot` ought to be taken before updates are queued.

		.. warning:: Updates to :term:`Parameters` with certain "configFile" values may require running :term:`ORT` in a different mode, occasionally manually. Though the server may appear to no longer have pending updates in these cases, until this manual intervention is performed the configuration *will* **not** *be correct*.

	Region
	Regions
		A group of :term:`Physical Location`\ s.

	reverse proxy
	reverse proxies
		A :dfn:`reverse proxy` acts on behalf of the :term:`origin server` such that the client is (potentially) unaware it is not communicating directly with the :term:`origin`. All Edge-tier :term:`cache server`\ s in a Traffic Control CDN are :dfn:`reverse proxies`. To the end user a Traffic Control-based CDN appears as a :dfn:`reverse proxy` since it retrieves content from the :term:`origin server`, acting on behalf of that :term:`origin server`. The client requests a URL that has a hostname which resolves to the :dfn:`reverse proxy`'s IP address and, in compliance with the HTTP 1.1 specification (:rfc:`2616`), the client sends a ``Host:`` header to the :dfn:`reverse proxy` that matches the hostname in the URL. The proxy looks up this hostname in a list of mappings to find the :term:`origin` hostname; if the hostname of the ``Host:`` header is not found in the list, the proxy will send an error (usually either ``404 Not Found`` or ``503 Service Unavailable`` as appropriate) to the client. If the supplied hostname is found in this list of mappings, the proxy checks its cache, and when the content is not already present, connects to the :term:`origin` to which the requested ``Host:`` maps requests the path of the original URL, providing the :term:`origin` hostname in the ``Host`` header. The proxy then stores the URL in its cache and serves the contents to the client. When there are subsequent requests for the same URL, a caching proxy serves the content out of its cache - provided :ref:`cache-revalidation` are satisfied - thereby reducing latency and network traffic.

		.. seealso:: `The Apache Traffic Server documentation on reverse proxy <https://docs.trafficserver.apache.org/en/latest/admin/reverse-proxy-http-redirects.en.html#http-reverse-proxy>`_.

		To insert a :dfn:`reverse proxy` into a typical HTTP 1.1 request and response flow, the :dfn:`reverse proxy` needs to be told where the :term:`origin server` can be reached (and which :term:`origin` to use for a given request when it's configured to proxy requests for multiple :term:`origin`\ s). In :abbr:`ATS (Apache Traffic Server)` this is handled by adding rules to `the remap.config configuration file <https://docs.trafficserver.apache.org/en/latest/admin-guide/files/remap.config.en.html>`_. The content owner must inform the clients, by updating the URL, to receive the content from the cache and not from the :term:`origin server` directly. For example, clients might be instructed to request content from ``http://www-origin-cache.cdn.com`` which points to a :dfn:`reverse proxy` for the actual :term:`origin` located at ``http://www.origin.com``.

		Now, if the client requests ``/foo/bar/fun.html`` from the :dfn:`reverse proxy` the sequence of events is as follows. is given the URL ``http://www-origin-cache.cdn.com/foo/bar/fun.html`` (note the different hostname) and when attempting to obtain that URL, the following occurs:

		#. The client sends a DNS request to the :abbr:`LDNS (Local Domain Name Server)` to resolve the name ``www-origin-cache.cdn.com`` to an IP address.
		#. The :abbr:`LDNS (Local Domain Name Server)` finds an IP address for ``www-origin-cache.cdn.com`` e.g. ``55.44.33.22``.
		#. The client sends an HTTP request for ``/foo/bar/fun.html`` to the IP address.

			.. code-block:: http
				:caption: Client Requests Content from the :dfn:`Reverse Proxy`

				GET /foo/bar/fun.html HTTP/1.1
				Host: www-origin-cache.cdn.com

		#. The :dfn:`reverse proxy` finds out the URL of the true :term:`origin` - in the case of :abbr:`ATS (Apache Traffic Server)` this is done by looking up ``www-origin-cache.cdn.com`` in its remap rules - and finds that it is ``www.origin.com``.
		#. The proxy checks its cache to see if the response for ``GET /foo/bar/fun.html HTTP/1.1`` from ``www.origin.com`` is already in the cache.
		#. If the response is not in the cache:

			#. The proxy sends the request to the actual :term:`origin`

				.. code-block:: http
					:caption: :dfn:`Reverse Proxy` Requests Content from the :term:`Origin Server`

					GET /foo/bar/fun.html HTTP/1.1
					Host: www.origin.com

			#. The :term:`origin server` responds with the requested content

				.. code-block:: http
					:caption: Response from the :term:`Origin Server`

					HTTP/1.1 200 OK
					Date: Sun, 14 Dec 2014 23:22:44 GMT
					Server: Apache/2.2.15 (Red Hat)
					Last-Modified: Sun, 14 Dec 2014 23:18:51 GMT
					ETag: "1aa008f-2d-50a3559482cc0"
					Content-Length: 45
					Connection: close
					Content-Type: text/html; charset=UTF-8

					<!DOCTYPE html><html><body>This is a fun file</body></html>

			#. The proxy sends the response on to the client, optionally adding a ``Via:`` header to indicate that the request was serviced by proxy.

				.. code-block:: http
					:caption: Resulting Response from the :dfn:`Reverse Proxy` to the Client

					HTTP/1.1 200 OK
					Date: Sun, 14 Dec 2014 23:22:44 GMT
					Last-Modified: Sun, 14 Dec 2014 23:18:51 GMT
					ETag: "1aa008f-2d-50a3559482cc0"
					Content-Length: 45
					Connection: close
					Content-Type: text/html; charset=UTF-8
					Age: 0
					Via: http/1.1 cache01.cdn.kabletown.net (ApacheTrafficServer/4.2.1 [uScSsSfUpSeN:t cCSi p sS])
					Server: ATS/4.2.1

					<!DOCTYPE html><html><body>This is a fun file</body></html>

			If, however, the response *was* already in the cache - and still valid according to the :ref:`cache-revalidation` - the proxy responds to the client with the previously retrieved result.

			.. code-block:: http
				:caption: The :dfn:`Reverse Proxy` Provides a Cached Response

				HTTP/1.1 200 OK
				Date: Sun, 14 Dec 2014 23:22:44 GMT
				Last-Modified: Sun, 14 Dec 2014 23:18:51 GMT
				ETag: "1aa008f-2d-50a3559482cc0"
				Content-Length: 45
				Connection: close
				Content-Type: text/html; charset=UTF-8
				Age: 39711
				Via: http/1.1 cache01.cdn.kabletown.net (ApacheTrafficServer/4.2.1 [uScSsSfUpSeN:t cCSi p sS])
				Server: ATS/4.2.1

				<!DOCTYPE html><html><body>This is a fun file</body></html>

	Role
	Roles
		Permissions :dfn:`Roles` define the operations a user is allowed to perform, and are currently an ordered list of permission levels.

	Snapshot
	Snapshots
		Previously called a "CRConfig" or "CRConfig.json" (and still called such in many places), this is a rather large set of routing information generated from a CDN's configuration and topology.

	Status
	Statuses
		A :dfn:`Status` represents the current operating state of a server. The default :dfn:`Statuses` made available on initial startup of Traffic Ops are related to the :ref:`health-proto` and are explained in that section.

	Tenant
	Tenants
		Users are grouped into :dfn:`Tenants` (or :dfn:`Tenancies`) to segregate ownership of and permissions over :term:`Delivery Service`\ s and their resources. To be clear, the notion of :dfn:`Tenancy` **only** applies within the context of :term:`Delivery Service`\ s and does **not** apply permissions restrictions to any other aspect of Traffic Control.

	Type
	Types
		A :dfn:`Type` defines a type of some kind of object configured in Traffic Ops. Unfortunately, that is exactly as specific as this definition can be.
