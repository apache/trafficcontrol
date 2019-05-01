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

.. _delivery-services:

*****************
Delivery Services
*****************
"Delivery Services" are a very important construct in :abbr:`ATC (Apache Traffic Control)`. At their most basic, they are a source of content and a set of :term:`cache server`\ s and configuration options used to distribute that content.

.. _ds-objects:

Delivery Service Objects
========================
Delivery Services are modeled several times over, in the Traffic Ops database, in Traffic Portal forms and tables, in the legacy Perl Traffic Ops codebase, and several times for various :ref:`to-api` versions in the new Go Traffic Ops codebase. Go-specific data structures can be found in `the project's GoDoc documentation <https://godoc.org/github.com/apache/trafficcontrol/lib/go-tc#DeliveryServiceNullableV11>`_. Rather than application-specific definitions, what follows is an attempt at consolidating all of the different properties and names of properties of Delivery Service objects throughout the :abbr:`ATC (Apache Traffic Control)` suite. The names of these fields are typically chosen as the most human-readable and/or most commonly-used names for the fields, and when reading please note that in many cases these names will appear camelCased or snake_cased to be machine-readable. Any aliases of these fields that are not merely case transformations of the indicated, canonical names will be noted in a table of aliases.

.. seealso:: The API reference for Delivery Service-related endpoints such as :ref:`to-api-deliveryservices` contains definitions of the Delivery Service object(s) returned and/or accepted by those endpoints.

Active
------
Whether or not this Delivery Service is active on the CDN and can be served. When a Delivery Service is not "active", Traffic Router will not be made aware of its existence - i.e. it will not appear in CDN :term:`Snapshot`\ s. Setting a Delivery Service to be "active" (or "inactive") will require that a new :term:`Snapshot` be taken.

Anonymous Blocking
------------------
Enables/Disables blocking of anonymized IP address - proxies, :abbr:`TOR (The Onion Ring)` exit nodes, etc - for this Delivery Service. Set to true to enable blocking of anonymous IPs for this Delivery Service.

.. table:: Aliases

	+--------------------------+-----------------------------------------------------------------------------+-----------------------------------------------------------------------------------------+
	| Name                     | Use(s)                                                                      | Type(s)                                                                                 |
	+==========================+=============================================================================+=========================================================================================+
	| anonymousBlockingEnabled | Traffic Ops client and server Go code, :ref:`to-api` requests and responses | usually unchanged (boolean), but sometimes as a string containing a boolean e.g. in the |
	|                          |                                                                             | response of a ``GET`` request to :ref:`to-api-cdns-name-snapshot`                       |
	+--------------------------+-----------------------------------------------------------------------------+-----------------------------------------------------------------------------------------+

.. note:: Anonymous Blocking requires an anonymous IP address database from the Delivery Service's Geolocation Provider. E.g. `MaxMind's Anonymous IP Database <https://www.maxmind.com/en/solutions/geoip2-enterprise-product-suite/anonymous-ip-database>`_ when MaxMind is used as the Geolocation Provider.

.. seealso:: The :ref:`anonymous_blocking-qht` "Quick-How-To" guide.

Cache URL Expression
--------------------
.. deprecated:: 3.0
	This feature is no longer supported by :abbr:`ATS (Apache Traffic Server)` and consequently it will be removed from Traffic Control in the future.

Manipulates the cache key of the incoming requests. Normally, the cache key is the :term:`origin` domain. This can be changed so that multiple services can share a cache key, can also be used to preserve cached content if service origin is changed.

.. warning:: This field provides access to a feature that was only present in :abbr:`ATS (Apache Traffic Server)` 6.X and earlier. As :term:`cache server`\ s must now use :abbr:`ATS (Apache Traffic Server)` 7.1.X, this field **must** be blank unless all :term:`cache servers` can be guaranteed to use that older :abbr:`ATS (Apache Traffic Server)` version (**NOT** recommended).

CDN
---
A CDN to which this Delivery Service belongs. Only servers and :term:`Cache Group`\ s within this CDN are available to route content for this Delivery Service. Additionally, only Traffic Routers assigned to this CDN will perform said routing. Most often ``cdn``/``CDN`` refers to the *name* of the CDN to which the Delivery Service belongs, but occasionally (most notably in the payloads and/or query parameters of certain :ref:`to-api` endpoints) it actually refers to the *integral, unique identifier* of said CDN.

Check Path
----------
A request path on the :term:`origin server` which is used to by certain :ref:`Traffic Ops Extensions <admin-to-ext-script>` to indicate the "health" of the :term:`origin`.

Deep Caching
------------
Controls the :ref:`deep-cache` feature of Traffic Router when serving content for this Delivery Service. This should always be represented by one of two values:

ALWAYS
	This Delivery Service will always use :ref:`deep-cache`
NEVER
	This Delivery Service will never use :ref:`deep-cache`

.. impl-detail:: Traffic Ops and Traffic Ops client Go code use an empty string as the name of the enumeration member that represents "NEVER".

Display Name
------------
The "name" of the Delivery Service. Since nearly any use of a string-based identification method for Delivery Services (e.g. in Traffic Portal tables) uses xml_id_, this is of limited use. For that reason and for consistency's sake it is suggested that this be the same as the xml_id_. However, unlike the xml_id_, this can contain any UTF-8 characters without restriction.

DNS Bypass CNAME
----------------
When the limits placed on this Delivery Service by the `Global Max Mbps`_ and/or `Global Max Tps`_ are exceeded, a DNS-:ref:`Routed <ds-types>` Delivery Service will direct excess traffic to the host referred to by this :abbr:`CNAME (Canonical Name)` record.

.. note:: IPv6 traffic will be redirected if and only if `IPv6 Routing Enabled`_ is "true" for this Delivery Service.

DNS Bypass IP
-------------
When the limits placed on this Delivery Service by the `Global Max Mbps`_ and/or `Global Max Tps`_ are exceeded, a DNS-:ref:`Routed <ds-types>` Delivery Service will direct excess IPv4 traffic to this IPv4 address.

DNS Bypass IPv6
---------------
When the limits placed on this Delivery Service by the `Global Max Mbps`_ and/or `Global Max Tps`_ are exceeded, a DNS-:ref:`Routed <ds-types>` Delivery Service will direct excess IPv6 traffic to this IPv6 address.

.. note:: This requires an accompanying configuration of `IPv6 Routing Enabled`_ such that IPv6 traffic is allowed at all.

DNS Bypass TTL
--------------
When the limits placed on this Delivery Service by the `Global Max Mbps`_ and/or `Global Max Tps`_ are exceeded, a DNS-:ref:`Routed <ds-types>` Delivery Service will direct excess traffic to their `DNS Bypass IP`_, `DNS Bypass IPv6`_, or `DNS Bypass CNAME`_.

DNS TTL
-------
The :abbr:`TTL (Time To Live)` on the DNS record for the Traffic Router A and AAAA records. DNS-:ref:`Routed <ds-types>` Delivery Services will send this :abbr:`TTL (Time To Live)` along with their record responses to clients requesting access to this Delivery Service. Setting too high or too low will result in poor caching performance.

.. table:: Aliases

	+-------------+--------------------------------------------------------------------------------------+---------------------------------------------+
	| Name        | Use(s)                                                                               | Type(s)                                     |
	+=============+======================================================================================+=============================================+
	| CCR DNS TTL | In Delivery Service objects returned by the :ref:`to-api`                            | unchanged (``int``, ``integer`` etc.)       |
	+-------------+--------------------------------------------------------------------------------------+---------------------------------------------+
	| CCR TTL     | Legacy Traffic Ops UI, documentation for older Traffic Control versions              | unchanged (``int``, ``integer`` etc.)       |
	+-------------+--------------------------------------------------------------------------------------+---------------------------------------------+
	| ttl         | In CDN :term:`Snapshot` structures, where it is displayed on a per-record-type-basis | map of record type names to integral values |
	+-------------+--------------------------------------------------------------------------------------+---------------------------------------------+

.. _ds-dscp:

DSCP
----
The :abbr:`DSCP (Differentiated Services Code Point)` which will be used to mark IP packets as they are sent out of the CDN to the client.

.. seealso:: `The Differentiated Services Wikipedia article <https://en.wikipedia.org/wiki/Differentiated_services>`_.

.. danger:: The :abbr:`DSCP (Differentiated Services Code Point)` setting in Traffic Portal is *only* for setting traffic towards the client, and gets applied *after* the initial TCP handshake is complete and the HTTP request has been received. Before that the cache can't determine what Delivery Service is being requested, and consequently can't know what :abbr:`DSCP (Differentiated Services Code Point)` to apply. Therefore, the :abbr:`DSCP (Differentiated Services Code Point)` feature can not be used for security settings; the IP packets that form the TCP handshake are not going to be :abbr:`DSCP (Differentiated Services Code Point)`-marked.

.. impl-detail:: DSCP settings only apply on :term:`cache servers` that run :abbr:`Apache Traffic Server`. The implementation uses the `ATS Header Rewrite Plugin <https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/header_rewrite.en.html>`_ to create a rule that will mark traffic bound outward from the CDN to the client.

Edge Header Rewrite Rules
-------------------------
This field in general contains the contents of the a configuration file used by the `ATS Header Rewrite Plugin <https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/header_rewrite.en.html>`_ when serving content for this Delivery Service - on :term:`Edge-tier cache server`\ s.

.. tip:: Because this ultimately is the contents of an :abbr:`ATS (Apache Traffic Server)` configuration file, it can make use of the :ref:`ort-special-strings`.

Fair-Queuing Pacing Rate Bps
----------------------------
The maximum bytes per second a :term:`cache server` will deliver on any single TCP connection. This uses the Linux kernel’s Fair-Queuing :manpage:`setsockopt(2)` (``SO_MAX_PACING_RATE``) to limit the rate of delivery. Traffic exceeding this speed will only be rate-limited and not diverted. This option requires extra configuration on all :term:`cache servers` assigned to this Delivery Service - specifically, the line ``net.core.default_qdisc = fq`` must exist in :file:`/etc/sysctl.conf`.

.. seealso:: :manpage:`tc-fq_codel(8)`

.. table:: Aliases

	+--------------+---------------------------------------------------------------------------------+---------------------------------------+
	| Name         | Use(s)                                                                          | Type(s)                               |
	+==============+=================================================================================+=======================================+
	| FQPacingRate | Traffic Ops source code, Delivery Service objects returned by the :ref:`to-api` | unchanged (``int``, ``integer`` etc.) |
	+--------------+---------------------------------------------------------------------------------+---------------------------------------+

Geo Limit
---------
Limits access to a Delivery Service by geographic location. The only practical difference between this and `Regional Geoblocking`_ is the configuration method; as opposed to `Regional Geoblocking`_, GeoLimit configuration is handled by country-wide codes and the :term:`Coverage Zone File`. When a client is denied access to a requested resource on an HTTP-:ref:`Routed <ds-types>` Delivery Service, they will receive a ``503 Service Unavailable`` instead of the usual ``302 Found`` response - unless `Geo Limit Redirect URL`_ is defined, in which case a ``302 Found`` response pointing to that URL will be returned by Traffic Router. If the Delivery Service is a DNS-:ref:`Routed <ds-types>` Delivery Service, the IP address of the *resolver* for the client DNS request is what is checked. If the IP address of this resolver is found to be in a restricted location, the Traffic Router will respond with an ``NXDOMAIN`` response, causing the name resolution to fail. This is nearly always an integral, unique identifier for a behavior set to be followed by Traffic Router. The defined values are:

0
	Geographic access limiting is not enabled, and content served by this Delivery Service will be accessible regardless of the clients geographic location. (Aliased as "0 - None" in Traffic Portal forms)
1
	A client will be allowed to request content if and only if their IP address is found by Traffic Router within the :term:`Coverage Zone File`. Otherwise, access will be denied. (Aliased as "1 - CZF Only" in Traffic Portal forms)
2
	A client will be allowed to request content if their IP address is found by Traffic Router within the :term:`Coverage Zone File`, or if looking up the client's IP address in the Geographic IP mapping database provided by `Geolocation Provider`_ indicates the client resides in a country that is found in the `Geo Limit Countries`_ array. (Aliased as "2 - CZF + Country Code(s)" in Traffic Portal forms - formerly was known as "CZF + US" when only the US country code was supported)

.. warning:: The definitions of each integral, unique identifier are hidden in implementations in each :abbr:`ATC (Apache Traffic Control)` component. Different components will handle invalid values differently, and there's no actual enforcement that the stored integral, unique identifier actually be within the representable range.

.. table:: Aliases

	+------------------+---------------------------------------------------------------------------+------------------------------------------------------------------------------------------------+
	| Name             | Use(s)                                                                    | Type(s)                                                                                        |
	+==================+===========================================================================+================================================================================================+
	| coverageZoneOnly | In CDN :term:`Snapshot` structures, especially in :ref:`to-api` responses | A boolean which, if ``true``, tells Traffic Router to only service requests when the client IP |
	|                  |                                                                           | address is found in the :term:`Coverage Zone File`                                             |
	+------------------+---------------------------------------------------------------------------+------------------------------------------------------------------------------------------------+

.. danger:: Geographic access limiting is **not** sufficient to guarantee access is properly restricted. The limiting is implemented by Traffic Router, which means that direct requests to :term:`Edge-tier cache server`\ s will bypass it entirely.

Geo Limit Countries
-------------------
When `Geo Limit`_ is being used with this Delivery Service (and is set to exactly ``2``), this is optionally a list of country codes to which access to content provided by the Delivery Service will be restricted. Normally, this is a comma-delimited string of said country codes. When creating a Delivery Service with this field or modifying the Geo Limit Countries field on an existing Delivery Service, any amount of whitespace between country codes is permissible, as it will be removed on submission, but responses from the :ref:`to-api` should never include such whitespace.

.. table:: Aliases

	+------------------+---------------------------------------------------------------------------+------------------------------------------------------------------------------------------------+
	| Name             | Use(s)                                                                    | Type(s)                                                                                        |
	+==================+===========================================================================+================================================================================================+
	| geoEnabled       | In CDN :term:`Snapshot` structures, especially in :ref:`to-api` responses | An array of objects each having the key "countryCode" that is a string containing an allowed   |
	|                  |                                                                           | country code - one should exist for each allowed country code                                  |
	+------------------+---------------------------------------------------------------------------+------------------------------------------------------------------------------------------------+

Geo Limit Redirect URL
----------------------
If `Geo Limit`_ is being used with this Delivery Service, this is optionally a URL to which clients will be redirected when Traffic Router determines that they are not in a geographic zone that permits their access to the Delivery Service content. This changes the response from Traffic Router from ``503 Service Unavailable`` to ``302 Found`` with a provided location that will be this URL. There is no restriction on the provided URL; it may even be the path to a resource served by this Delivery Service. In fact, this field need not even be a full URL, it can be a relative path. Both of these cases are handled specially by Traffic Router.

- If the provided URL is a resource served by the Delivery Service (e.g. if the client requests ``http://cdn.dsXMLID.somedomain.example.com/index.html`` but are denied access by `Geo Limit`_ and the Geo Limit Redirect URL is something like ``http://cdn.dsXMLID.somedomain.example.com/help.php``), Traffic Router will find an appropriate :term:`Edge-tier cache server` and redirect the client, ignoring Geo Limit restrictions *for this request only*.
- If the provided "URL" is actually a relative path, it will be considered *relative to the requested Delivery Service :abbr:`FQDN (Fully Qualified Domain Name)`*. This means that e.g. if the client requests ``http://cdn.dsXMLID.somedomain.example.com/index.html`` but are denied access by `Geo Limit`_ and the Geo Limit Redirect URL is something like ``/help.php``, Traffic Router will find an appropriate :term:`Edge-tier cache server` and redirect the client to it as though they had requested ``http://cdn.dsXMLID.somedomain.example.com/help.php``, ignoring `Geo Limit`_ restrictions *for this request only*.

.. table:: Aliases

	+---------------------------------+----------------------------------------------------------------+-------------------------------------------------------------------------------------------------+
	| Name                            | Use(s)                                                         | Type(s)                                                                                         |
	+=================================+================================================================+=================================================================================================+
	| :abbr:`NGB (National GeoBlock)` | Older documentation, in Traffic Router comments and error logs | unchanged (``string``, ``String`` etc.)                                                         |
	+---------------------------------+----------------------------------------------------------------+-------------------------------------------------------------------------------------------------+
	| geoRedirectURLType              | Internally in Traffic Router                                   | A ``String`` that describes whether or not the actual Geo Limit Redirect URL is relative to the |
	|                                 |                                                                | Delivery Service base :abbr:`FQDN (Fully Qualified Domain Name)`. Should be one of:             |
	|                                 |                                                                |                                                                                                 |
	|                                 |                                                                | INVALID_URL                                                                                     |
	|                                 |                                                                |     The Geo Limit Redirect URL has not yet been parsed, or an error occurred during parsing     |
	|                                 |                                                                | DS_URL                                                                                          |
	|                                 |                                                                |     The Geo Limit Redirect URL is served by this Delivery Service                               |
	|                                 |                                                                | NOT_DS_URL                                                                                      |
	|                                 |                                                                |     The Geo Limit Redirect URL is external to this Delivery Service                             |
	+---------------------------------+----------------------------------------------------------------+-------------------------------------------------------------------------------------------------+

.. note:: The use of a redirect URL relies on the ability of Traffic Router to redirect the client using HTTP ``302 Found`` responses. As such, this field has no effect on DNS-:ref:`Routed <ds-types>` Delivery Services.

Geolocation Provider
--------------------
This is nearly always the integral, unique identifier of a provider for a database that maps IP addresses to geographic locations. Less frequently, this may be accompanied by the actual name of the provider. Only two values are possible at the time of this writing:

0: MaxMind
	IP address to geographic location mapping will be provided by a `MaxMind GeoIP2 database <https://www.maxmind.com/en/geoip2-databases>`_.
1: Neustar
	IP address to geographic location mapping will be provided by a `Neustar GeoPoint IP address database <https://www.security.neustar/digital-performance/ip-intelligence/ip-address-data>`_.

	.. warning:: It's not clear whether Neustar databases are actually supported; this is an old option and compatibility may have been broken over time.

.. table:: Aliases

	+-------------+-------------------------------------------------------------------------------+-----------------------------------------+
	| Name        | Use(s)                                                                        | Type(s)                                 |
	+=============+===============================================================================+=========================================+
	| geoProvider | Traffic Ops and Traffic Ops client code, :ref:`to-api` requests and responses | unchanged (integral, unique identifier) |
	+-------------+-------------------------------------------------------------------------------+-----------------------------------------+

Geo Miss Default Latitude
-------------------------
Default Latitude for this Delivery Service. When the geographic location of the client cannot be determined, they will be routed as if they were at this latitude.

.. table:: Aliases

	+---------+--------------------------------------------------------+---------------------+
	| Name    | Use(s)                                                 | Type(s)             |
	+---------+--------------------------------------------------------+---------------------+
	| missLat | In :ref:`to-api` responses and Traffic Ops source code | unchanged (numeric) |
	+---------+--------------------------------------------------------+---------------------+

Geo Miss Default Longitude
--------------------------
Default Longitude for this Delivery Service. When the geographic location of the client cannot be determined, they will be routed as if they were at this longitude.

.. table:: Aliases

	+----------+--------------------------------------------------------+---------------------+
	| Name     | Use(s)                                                 | Type(s)             |
	+----------+--------------------------------------------------------+---------------------+
	| missLong | In :ref:`to-api` responses and Traffic Ops source code | unchanged (numeric) |
	+----------+--------------------------------------------------------+---------------------+

Global Max Mbps
---------------
The maximum :abbr:`Mbps (Megabits per second)` this Delivery Service can serve across all :term:`Edge-tier cache server`\ s before traffic will be diverted to the bypass destination. For a DNS-:ref:`Routed <ds-types>` Delivery Service, the `DNS Bypass IP`_ or `DNS Bypass IPv6`_ will be used (depending on whether this was a A or AAAA request), and for HTTP-:ref:`Routed <ds-types>` Delivery Services the `HTTP Bypass FQDN`_ will be used.

.. table:: Aliases

	+--------------------+--------------------------------------------------------------------------------------+------------------------------------------------------------------------------------------------------------------+
	| Name               | Use(s)                                                                               | Type(s)                                                                                                          |
	+====================+======================================================================================+==================================================================================================================+
	| totalKbpsThreshold | In :ref:`to-api` responses - most notably :ref:`to-api-cdns-name-configs-monitoring` | unchanged (numeric), but converted from :abbr:`Mbps (Megabits per second)` to :abbr:`Kbps (kilobits per second)` |
	+--------------------+--------------------------------------------------------------------------------------+------------------------------------------------------------------------------------------------------------------+

Global Max TPS
--------------
The maximum :abbr:`TPS (Transactions per Second)` this Delivery Service can serve across all :term:`Edge-tier cache server`\ s before traffic will be diverted to the bypass destination. For a DNS-:ref:`Routed <ds-types>` Delivery Service, the `DNS Bypass IP`_ or `DNS Bypass IPv6`_ will be used (depending on whether this was a A or AAAA request), and for HTTP-:ref:`Routed <ds-types>` Delivery Services the `HTTP Bypass FQDN`_ will be used.

.. table:: Aliases

	+-------------------+--------------------------------------------------------------------------------------+---------------------+
	| Name              | Use(s)                                                                               | Type(s)             |
	+===================+======================================================================================+=====================+
	| totalTpsThreshold | In :ref:`to-api` responses - most notably :ref:`to-api-cdns-name-configs-monitoring` | unchanged (numeric) |
	+-------------------+--------------------------------------------------------------------------------------+---------------------+

HTTP Bypass FQDN
----------------
When the limits placed on this Delivery Service by the `Global Max Mbps`_ and/or `Global Max Tps`_ are exceeded, an HTTP-:ref:`Routed <ds-types>` Delivery Service will direct excess traffic to this :abbr:`Fully Qualified Domain Name`.

IPv6 Routing Enabled
--------------------
A boolean value that controls whether or not clients using IPv6 can be routed to this Delivery Service by Traffic Router. When creating a Delivery Service in Traffic Portal, this will default to "true".

Info URL
--------
This should be a URL (though neither the :ref:`to-api` nor the Traffic Ops Database in any way enforce the validity of said URL) to which administrators or others may refer for further information regarding a Delivery Service - e.g. a related JIRA ticket.

Initial Dispersion
------------------
The number of :term:`cache servers` between which traffic requesting the same object will be randomly split - meaning that if 4 clients all request the same object (one after another), then if this is above 4 there is a possibility that all 4 are cache misses, necessitating a fresh pull of the content from the next-highest level in the CDN. For most use-cases, this should be ``1``.

Logs Enabled
------------
A boolean switch that can be toggled to enable/disable logging for a Delivery Service.

.. note:: This doesn't actually do anything. It was part of the functionality for a planned Traffic Control component named "Traffic Logs" - which was never created.

Long Description
----------------
Free text field that has no strictly defined purpose, but it is suggested that it contain a short description of the Delivery Service and its purpose.

.. table::

	+----------+---------------------------------------------------------+-----------------------------------------+
	| Name     | Use(s)                                                  | Type(s)                                 |
	+==========+=========================================================+=========================================+
	| longDesc | Traffic Control source code and :ref:`to-api` responses | unchanged (``string``, ``String`` etc.) |
	+----------+---------------------------------------------------------+-----------------------------------------+

Long Description 2
------------------
Free text field that has no strictly defined purpose.

.. table::

	+----------------------------+---------------------------------------------------------+-----------------------------------------+
	| Name                       | Use(s)                                                  | Type(s)                                 |
	+============================+=========================================================+=========================================+
	| longDesc1\ [#cardinality]_ | Traffic Control source code and :ref:`to-api` responses | unchanged (``string``, ``String`` etc.) |
	+----------------------------+---------------------------------------------------------+-----------------------------------------+

Long Description 3
------------------
Free text field that has no strictly defined purpose.

.. table::

	+----------------------------+---------------------------------------------------------+-----------------------------------------+
	| Name                       | Use(s)                                                  | Type(s)                                 |
	+============================+=========================================================+=========================================+
	| longDesc2\ [#cardinality]_ | Traffic Control source code and :ref:`to-api` responses | unchanged (``string``, ``String`` etc.) |
	+----------------------------+---------------------------------------------------------+-----------------------------------------+

Max DNS Answers
---------------
The maximum number of :term:`Edge-tier cache server` IP addresses that the Traffic Router will include in responses to DNS requests for DNS-:ref:`Routed <ds-types>` Delivery Services. The :ref:`to-api` restricts this value to the range [1, 15], but no matching restraints are placed on the actual data as stored in the Traffic Ops Database. When provided, the :term:`cache server` IP addresses included are rotated in each response to spread traffic evenly. Ideally this number will reflect the amount of traffic - e.g. ``1`` for a trial Delivery Service with very little traffic, ``2`` for a small production Delivery Service. Add 1 for every 20 :abbr:`Gbps (Gigabits per second)` of traffic you expect at peak.

Mid Header Rewrite Rules
------------------------
This field in general contains the contents of the a configuration file used by the `ATS Header Rewrite Plugin <https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/header_rewrite.en.html>`_ when serving content for this Delivery Service - on :term:`Mid-tier cache server`\ s.

.. tip:: Because this ultimately is the contents of an :abbr:`ATS (Apache Traffic Server)` configuration file, it can make use of the :ref:`ort-special-strings`.

Origin Server Base URL
----------------------
The Origin Server’s base URL which includes the protocol (http or https). Example: ``http://movies.origin.com``. Must not include paths, query parameters, document fragment identifiers, or username/password URL fields.

.. table:: Aliases

	+---------------+------------------------------------------------------------+----------------------------------------------+
	| Name          | Use(s)                                                     | Type(s)                                      |
	+===============+============================================================+==============================================+
	| orgServerFqdn | :ref:`to-api` responses and in Traffic Control source code | unchanged (usually ``str``, ``string`` etc.) |
	+---------------+------------------------------------------------------------+----------------------------------------------+

Origin Shield
-------------
An experimental feature that allows administrators to list additional forward proxies that sit between the :term:`Mid-tier` and the :term:`origin`. In most scenarios, this is represented (and required to be input) as a pipe (``|``)-delimited string.

Profile
-------
Either the name of a :term:`Profile` used by this Delivery Service, or an integral, unique identifier for said :term:`Profile`.

.. table:: Aliases

	+-------------+------------------------------------------------------------------------------------------------+----------------------------------------------------------------------------------------+
	| Name        | Use(s)                                                                                         | Type(s)                                                                                |
	+=============+================================================================================================+========================================================================================+
	| profileId   | In Traffic Control source code and some :ref:`to-api` responses dealing with Delivery Services | Unlike the more general "Profile", this is *always* an integral, unique identifier     |
	+-------------+------------------------------------------------------------------------------------------------+----------------------------------------------------------------------------------------+
	| profileName | In Traffic Control source code and some :ref:`to-api` responses dealing with Delivery Services | Unlike the more general "Profile", this is *always* a name (``str``, ``string``, etc.) |
	+-------------+------------------------------------------------------------------------------------------------+----------------------------------------------------------------------------------------+

Protocol
--------
The protocol with which to serve content from this Delivery Service. This defines the way the Delivery Service will handle client requests that are either HTTP or HTTPS, which is distinct from what protocols are used to direct traffic. For example, this can be used to direct clients to only request content using HTTP, or to allow clients to use either HTTP or HTTPS, etc. Normally, this will be the name of the protocol handling, but occasionally this will appear as the integral, unique identifier of the protocol handling instead. The integral, unique identifiers and their associated names and meanings are:

0: HTTP
	This Delivery Service will only accept unsecured HTTP requests. Requests made with HTTPS will fail.
1: HTTPS
	This Delivery Service will only accept secured HTTPS requests. Requests made with HTTP will fail.
2: HTTP AND HTTPS
	This Delivery Service will accept both unsecured HTTP requests and secured HTTPS requests.
3: HTTP TO HTTPS
	When this Delivery Service is using HTTP :ref:`Content Routing <ds-types>` unsecured HTTP requests will be met with a response that indicates to the client that further requests must use HTTPS.

	.. note:: If any other type of :ref:`Content Routing <ds-types>` is used, this functionality cannot be used. In those cases, a protocol setting of ``3``/"HTTP TO HTTPS" will result in the same behavior as ``1``/"HTTPS". This behavior is tracked by `GitHub Issue #3221 <https://github.com/apache/trafficcontrol/issues/3221>`_


.. warning:: The definitions of each integral, unique identifier are hidden in implementations in each :abbr:`ATC (Apache Traffic Control)` component. Different components will handle invalid values differently, and there's no actual enforcement that the stored integral, unique identifier actually be within the representable range.

.. table:: Aliases

	+----------+-------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| Name     | Use(s)                  | Type(s)                                                                                                                                                             |
	+==========+=========================+=====================================================================================================================================================================+
	| Protocol | CDN :term:`Snapshot`\ s | An object containing the key ``"acceptHttps"`` that is a string containing a boolean that expresses whether Traffic Router should accept HTTPS requests for this    |
	|          |                         | Delivery Service, and the key ``"redirectToHttps"`` that is also a string containing a boolean which expresses whether or not Traffic Router should redirect HTTP   |
	|          |                         | requests to HTTPS URLs. Optionally, the key ``"acceptHttp"`` may also appear, once again a string containing a boolean that expresses whether or not Traffic Router |
	|          |                         | should accept unsecured HTTP requests - this is implicitly treated as ``"true"`` by Traffic Router when it is not present.                                          |
	+----------+-------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------+

Query String Handling
---------------------
Describes how query strings should be handled by the :term:`Edge-tier cache server`\ s when serving content for this Delivery Service. This is nearly always expressed as an integral, unique identifier for each behavior, though in Traffic Portal a more descriptive value is typically used, or at least provided in addition to the integral, unique identifier. The allowed values and their meanings are:

0
	For the purposes of caching, :term:`Edge-tier cache server`\ s will consider URLs unique if and only if they are unique up to and including any and all query parameters. They will also pass the query parameters in their own requests to :term:`Mid-tier cache server`\ s (which in turn will exhibit the same caching behavior and pass the query parameters in requests to the :term:`origin`). (Aliased as "USE" in Traffic Portal tables, and "0 - use qstring in cache key, and pass up" in Traffic Portal forms)
1
	For the purposes of caching, neither :term:`Edge-tier` nor :term:`Mid-tier cache server`\ s will consider the query parameter string when determining if a URL is stored in cache. However, the query string will still be passed in upstream requests to :term:`Mid-tier cache server`\ s and in turn the :term:`origin`. (Aliased as "IGNORE" in Traffic Portal tables and "1 - ignore in cache key, and pass up" in Traffic Portal forms)
2
	The query parameter string will be stripped from URLs immediately when the request is received by an :term:`Edge-tier cache server`. This means it is never considered for the purposes of caching unique URLs and will not be passed in upstream requests. (Aliased as "DROP" in Traffic Portal tables and "2 - drop at edge" in Traffic Portal forms)

	.. warning:: The implementation of dropping query parameter strings at the :term:`Edge-tier` uses a `Regex Remap Expression`_ and thus Delivery Services with this type of query string handling cannot make use of `Regex Remap Expression`_\ s.

.. table:: Aliases

	+------------------+------------------------------------------------------------+-----------------------------------------------------------------------------------------+
	| Name             | Use(s)                                                     | Type(s)                                                                                 |
	+==================+============================================================+=========================================================================================+
	| Qstring Handling | Traffic Portal tables                                      | One of the Traffic Portal value aliases "USE" (``0``), "IGNORE" (``1``), "DROP" (``2``) |
	+------------------+------------------------------------------------------------+-----------------------------------------------------------------------------------------+
	| qstringIgnore    | Traffic Ops Go/Perl code, :ref:`to-api` requests/responses | unchanged (integral, unique identifier)                                                 |
	+------------------+------------------------------------------------------------+-----------------------------------------------------------------------------------------+

Range Request Handling
----------------------
Describes how HTTP "Range Requests" should be handled by the Delivery Service at the :term:`Edge-tier`. This is nearly always an integral, unique identifier for the behavior set required of the :term:`Edge-tier cache server`\ s. The valid values and their respective meanings are:

0
	Do not cache Range Requests at all. (Aliased as "0 - Don't cache" in Traffic Portal forms)

		.. note:: This is not retroactive - when modifying an existing Delivery Services to have this value for "Range Request Handling", ranges requested from files that are already cached due to a non-range request will be served out of cache for as long as the Cache-Control headers allow.

1
	Use the `background_fetch <https://docs.trafficserver.apache.org/en/7.1.x/admin-guide/plugins/background_fetch.en.html>`_ plugin to service the range request while caching the whole object. (Aliased as "1 - Use background_fetch plugin" in Traffic Portal forms)
2
	Use the `cache_range_requests <https://github.com/apache/trafficserver/tree/7.1.x/plugins/experimental/cache_range_requests>`_ plugin to cache ranges as unique objects. (Aliased as "2 - Use cache_range_requests plugin" in Traffic Portal forms)

.. note:: Range Request Handling can only be implemented on :term:`cache server`\ s using :abbr:`ATS (Apache Traffic Server)` because of its dependence on :abbr:`ATS (Apache Traffic Server)` plugins. The value may be set on any Delivery Service, but will have no effect when the :term:`cache server`\ s that ultimately end up serving the content are e.g. Grove, Nginx, etc.

.. warning:: The definitions of each integral, unique identifier are hidden in implementations in each :abbr:`ATC (Apache Traffic Control)` component. Different components will handle invalid values differently, and there's no actual enforcement that the stored integral, unique identifier actually be within the representable range.

.. _ds-raw-remap:

Raw Remap Text
--------------
For HTTP and DNS-:ref:`Routed <ds-types>` Delivery Services, this will be added to the end of a line in the `remap.config ATS configuration file <https://docs.trafficserver.apache.org/en/7.1.x/admin-guide/files/remap.config.en.html>`_ line on the cache verbatim. For ANY_MAP-:ref:`Type <ds-types>` Delivery Services this must be defined.

.. tip:: Because this ultimately is a raw line of content in a configuration file, it can make use of the :ref:`ort-special-strings`. Of particular note is the :ref:`ort-remap-override` template string.

.. note:: This field **must** be defined on ANY_MAP-`Type`_ Delivery Services, but is otherwise optional.

.. table:: Aliases

	+-----------+-----------------------------------------------------------------+---------------------------------------+
	| Name      | Use(s)                                                          | Type(s)                               |
	+===========+=================================================================+=======================================+
	| remapText | In Traffic Ops source code and :ref:`to-api` requests/responses | unchanged (``text``, ``string`` etc.) |
	+-----------+-----------------------------------------------------------------+---------------------------------------+

Regex remap expression
----------------------
Allows remapping of incoming requests URL using regular expressions to search and replace text. In a more literal sense, this is the raw contents of a configuration file used by the `ATS regex_remap plugin  <https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/regex_remap.en.html>`_.

.. caution:: This field is not validated by Traffic Ops to be correct syntactically, and can cause Traffic Server to not start if invalid. Please use with caution.

.. warning:: Regex remap expressions are incompatible with `Query String Handling`_ being set to ``2``. The behavior of a :term:`cache server` under that configuration is undefined.

Regional Geoblocking
--------------------
A boolean value that defines whether or not :ref:`Regional Geoblocking <regionalgeo-qht>` is active on this Delivery Service. The actual configuration of :ref:`Regional Geoblocking <regionalgeo-qht>` is done in the :term:`Profile` used by the Traffic Router serving the Delivery Service. Rules for this Delivery Service may exist, but they will not actually be used unless this field is ``true``.

.. tip:: :ref:`Regional Geoblocking <regionalgeo-qht>` is configured primarily with respect to Canadian postal codes, so unless specifically Canadian regions should be allowed/disallowed to access content, `Geo Limit`_ is probably a better setting for controlling access to content according to geographic location.

Routing Name
------------
The smallest DNS zone used to create an :abbr:`FQDN (Fully Qualified Domain Name)` used by clients to request content. All together, the constructed :abbr:`FQDN (Fully Qualified Domain Name)` looks like: :file:`{Delivery Service Routing Name}.{Delivery Service xml_id}.{CDN Subdomain}.{CDN Domain}.{Top-Level Domain}`\ [#xmlValid]_.

Signing Algorithm
-----------------
URLs/URIs may be signed using one of two algorithms before a request for the content to which they refer is sent to the :term:`origin` (which in practice can be any upstream network). At the time of this writing, this field is restricted within the Traffic Ops Database to one of two values (or ``NULL``/"None", to indicate no signing should be done).

.. seealso:: The url_sig `README <https://github.com/apache/trafficserver/blob/master/plugins/experimental/url_sig/README>`_.

url_sig
	URL signing will be implemented in this Delivery Service using the `url_sig Apache Traffic Server plugin <https://docs.trafficserver.apache.org/en/7.1.x/admin-guide/plugins/url_sig.en.html>`_. (Aliased as "URL Signature Keys" in Traffic Portal forms)
uri_signing
	URL signing will be implemented in this Delivery Service using an algorithm based on a work-in-progress RFC specification draft. (Aliased as "URI Signing Keys" in Traffic Portal forms)

.. table:: Aliases

	+--------+------------------------------------------------------------------------------------------+---------------------------------------------------------------------------------------------+
	| Name   | Use(s)                                                                                   | Type(s)                                                                                     |
	+========+==========================================================================================+=============================================================================================+
	| Signed | In all components prior to Traffic Control v2.2. Some endpoints in early versions of the | A boolean value where ``true`` was the same as "url_sig" in current versions, and ``false`` |
	|        | :ref:`to-api` will still return this field instead of "signingAlgorithm".                | indicated URL signing would not be done for the Delivery Service.                           |
	+--------+------------------------------------------------------------------------------------------+---------------------------------------------------------------------------------------------+

Keys for either algorithm can be generated within :ref:`Traffic Portal <tp-services-delivery-service>`.

Tenant
------
The :term:`Tenant` who owns this Delivery Service. They (and their parents, if any) are the only ones allowed to make changes to this Delivery Service. Typically, ``tenant``/``Tenant`` refers to the *name* of the owning :term:`Tenant`, but occasionally (most notably in the payloads and/or query parameters of certain :ref:`to-api` requests) it actually refers to the *integral, unique identifier* of said :term:`Tenant`.

.. table:: Aliases

	+----------+----------------------------------------------+--------------------------------------------------------+
	| Name     | Use(s)                                       | Type(s)                                                |
	+==========+==============================================+========================================================+
	| TenantID | Go code and :ref:`to-api` requests/responses | Integral, unique identifier (``bigint``, ``int`` etc.) |
	+----------+----------------------------------------------+--------------------------------------------------------+

Traffic Router Additional Response Headers
------------------------------------------
List of HTTP header ``{{name}}:{{value}}`` pairs separated by ``__RETURN__`` or simply on separate lines. Listed pairs will be included in all HTTP responses from Traffic Router for HTTP-:ref:`Routed <ds-types>` Delivery Services.

.. deprecated:: 4.0
	The use of ``__RETURN__`` as a substitute for a real newline is unnecessary and the ability to do so will be removed in the future.

.. table:: Aliases

	+-------------------+----------------------------------------------------------------------------------------+-----------------------------+
	| Name              | Use(s)                                                                                 | Type(s)                     |
	+===================+========================================================================================+=============================+
	| trResponseHeaders | Traffic Control source code and Delivery Service objects returned by the :ref:`to-api` | unchanged (``string`` etc.) |
	+-------------------+----------------------------------------------------------------------------------------+-----------------------------+

Traffic Router Log Request Headers
----------------------------------
List of HTTP header names separated by ``__RETURN__`` or simply on separate lines. Listed pairs will be logged for all HTTP requests to Traffic Router for HTTP-:ref:`Routed <ds-types>` Delivery Services.

.. deprecated:: 4.0
	The use of ``__RETURN__`` as a substitute for a real newline is unnecessary and the ability to do so will be removed in the future.

.. table:: Aliases

	+------------------+----------------------------------------------------------------------------------------+-----------------------------+
	| Name             | Use(s)                                                                                 | Type(s)                     |
	+==================+========================================================================================+=============================+
	| trRequestHeaders | Traffic Control source code and Delivery Service objects returned by the :ref:`to-api` | unchanged (``string`` etc.) |
	+------------------+----------------------------------------------------------------------------------------+-----------------------------+

.. _ds-types:

Type
----
Defines the content routing method used by the Delivery Service. In most cases this is an integral, unique identifier that corresponds to an enumeration of the Delivery Service Types. In other cases, this the actual name of said type.

The "Type" of a Delivery Service can mean several things. First, it can be used to refer to the "routing type" of Delivery Service. This is one of:

.. tip:: The only way to get the integral, unique identifier of a :term:`Type` of Delivery Service is to look at the database after it has been generated; these are non-deterministic and cannot be guaranteed to have any particular value, or even consistent values. This can be done directly or, preferably, using the :ref:`to-api-types` endpoint. Unfortunately, knowing the name of the :term:`Type` is rarely enough for many applications. The ``useInColumn`` values of these :term:`Types` will be ``deliveryservice``.

DNS
	Delivery Services of this routing type are routed by Traffic Router by providing DNS records that provide the IP addresses of :term:`cache servers` when clients look up the full Delivery Service :abbr:`FQDN (Fully Qualified Domain Name)`.
HTTP
	The Traffic Router(s) responsible for routing this Delivery Service will still answer DNS requests for the Delivery Service :abbr:`FQDN (Fully Qualified Domain Name)`, but will provide its own IP address. The client then directs its HTTP request to the Traffic Router, which will use an `HTTP redirection response <https://developer.mozilla.org/en-US/docs/Web/HTTP/Status#Redirection_messages>`_ to direct the client to a :term:`cache server`.

More generally, though, Delivery Services have a Type that defines not only how traffic is routed, but also how content is cached and semantically defines what "content" means in the context of a given Delivery Service.

ANY_MAP
	This is a special kind of Delivery Service that should only be used when control over the clients is guaranteed, and very fine control over the :abbr:`ATS (Apache Traffic Server)` `remap.config  <https://docs.trafficserver.apache.org/en/7.1.x/admin-guide/files/remap.config.en.html>`_ line for this Delivery Service is required. ANY_MAP is not known to Traffic Router. It is not routed in any way. For Delivery Services of this type, the "Raw Remap Text" field **must** be defined, as it is the only configuration generated by Traffic Control. The only way for a client to utilize delivery through an ANY_MAP service is by knowing in advance the IP address of one or more :term:`Edge-tier cache server`\ s and make the appropriate request(s).
DNS
	Uses DNS content routing. Delivers content normally. This is the recommended Type for delivering smaller objects like web page images.
DNS_LIVE
	Uses DNS Content routing, but optimizes caching for live video streaming. Specifically, the configuration generated for :term:`cache servers` responsible for serving content for this Delivery Service will not cache that content on storage disks. Instead, they will make use of RAM block devices dedicated to ATS - as specified by the special ``RAM_Drive_Prefix`` and ``RAM_Drive_Letters`` :term:`Parameters`. Also, any :term:`Mid-tier` of caching is bypassed.
DNS_LIVE_NATNL
	Works exactly the same as DNS_LIVE, but is optimized for delivery of live video content across a wide physical area. What this means is that the :term:`Mid-tier` of caching is **not** bypassed, unlike DNS_LIVE. The :term:`Mid-tier` will also use block RAM devices.
HTTP
	Uses HTTP content routing, delivers content normally. This is the recommended Type for delivering larger objects like video streams.
HTTP_LIVE
	Uses HTTP Content routing, but optimizes caching for live video streaming. Specifically, the configuration generated for :term:`cache servers` responsible for serving content for this Delivery Service will not cache that content on storage disks. Instead, they will make use of RAM block devices dedicated to ATS - as specified by the special ``RAM_Drive_Prefix`` and ``RAM_Drive_Letters`` :term:`Parameters`. Also, any :term:`Mid-tier` of caching is bypassed.
HTTP_LIVE_NATNL
	Works exactly the same as HTTP_LIVE, but is optimized for delivery of live video content across a wide physical area. What this means is that the :term:`Mid-tier` of caching is **not** bypassed, unlike HTTP_LIVE. The :term:`Mid-tier` will also use block RAM devices.
HTTP_NO_CACHE
	Uses HTTP Content Routing, but :term:`cache servers` will not actually cache the delivered content - they act as just proxies. This will bypass any existing :term:`Mid-tier` entirely (as it's totally useless when content is not being cached).

STEERING
	This is a sort of "meta" Delivery Service. It is used for directing clients to one of a set of Delivery Services, rather than delivering content directly itself. The Delivery Services to which a STEERING Delivery Service routes clients are referred to as "targets". Targets in general have an associated "value" and can be of several :term:`Types` that define the meaning of the value - these being:

	STEERING_ORDER
		The value of a STEERING_ORDER target sets a strict order of preference. In cases where a response to a client contains multiple Delivery Services, those targets with a lower "value" appear earlier than those with a higher "value". In cases where two or more targets share the same value, they each have an equal chance of being presented to the client - effectively spreading traffic evenly across them.
	STEERING_WEIGHT
		The values of STEERING_WEIGHT targets are interpreted as "weights", which define how likely it is that any given client will be routed to a specific Delivery Service - effectively this determines the spread of traffic across each target.
	STEERING_GEO_ORDER
		These targets behave exactly like STEERING_ORDER targets, but Delivery Services are grouped according to the "locations" of their :term:`origins`. Before choosing a Delivery Service to which to direct the client, Traffic Router will first create a subset of choices by eliminating all but the closest location to the client as possibilities. Once this subset is created, the values of the targets establish a strict precedence ordering, just like STEERING_ORDER targets.
	STEERING_GEO_WEIGHT
		These targets behave exactly like STEERING_WEIGHT targets, but Delivery Services are grouped according to the "locations" of their :term:`origins`. Before choosing a Delivery Service to which to direct the client, Traffic Router will first create a subset of choices by eliminating all but the closest location to the client as possibilities. Once this subset is chosen, the values of the targets establish the likelihood that any given target within the subset will be chosen for the client - effectively determining the spread of traffic across targets within that subset.

	The targets of a Delivery Service may be set using :ref:`the appropriate section of Traffic Portal <tp-services-delivery-service>` or via the :ref:`to-api-steering-id-targets` and :ref:`to-api-steering-id-targets-targetID` :ref:`to-api` endpoints.

	.. important:: To make use of the STEERING_GEO_ORDER and/or STEERING_GEO_WEIGHT target types, it is first necessary to ensure that at least the "primary" :term:`origin` of the :term:`Delivery Service` has an associated geographic coordinate pair. This can be done either from the :ref:`tp-configure-origins` page in Traffic Portal, or using the :ref:`to-api-origins` :ref:`to-api` endpoint.

	.. seealso:: For more information on setting up a STEERING (or CLIENT_STEERING) Delivery Service, see :ref:`steering-qht`.

	.. seealso:: For implementation details about how Traffic Router routes STEERING (and CLIENT_STEERING) Delivery Services, see :ref:`tr-steering`.

CLIENT_STEERING
	A CLIENT_STEERING Delivery Service is exactly like STEERING except that it provides clients with methods of bypassing the weights, orders, and localizations of targets in order to choose any arbitrary target at will. When utilizing these methods, the client will either directly choose a target immediately or request a list of all available targets from Traffic Router and then choose one to which to send a subsequent request for actual content.

.. note:: "Steering" is also commonly used to collectively refer to either of the kinds of Delivery Services that can participate in steering behavior (STEERING and CLIENT_STEERING).

.. table:: Aliases

	+----------------------+-------------------------------------------------+-----------------------------------------------------------------+
	| Name                 | Use(s)                                          | Type(s)                                                         |
	+======================+=================================================+=================================================================+
	| Content Routing Type | Traffic Portal forms                            | The name of any of the Delivery Service `Type`_\ s (``string``) |
	+----------------------+-------------------------------------------------+-----------------------------------------------------------------+
	| TypeID               | In Go code and :ref:`to-api` requests/responses | Integral, unique identifier (``bigint``, ``int`` etc.)          |
	+----------------------+-------------------------------------------------+-----------------------------------------------------------------+

Use Multi-Site Origin Feature
-----------------------------
A boolean value that indicates whether or not this Delivery Service uses :ref:`multi-site-origin`. There are very few good reasons for this to not be ``false``.

.. table:: Aliases

	+-----------------+-------------------------------------------------+----------------------------------------+
	| Name            | Use(s)                                          | Type(s)                                |
	+=================+=================================================+========================================+
	| multiSiteOrigin | In Go code and :ref:`to-api` requests/responses | unchanged (``bool``, ``boolean`` etc.) |
	+-----------------+-------------------------------------------------+----------------------------------------+

.. _ds-xmlid:

xml_id
------
A text-based unique identifier for a Delivery Service. Many :ref:`to-api` endpoints and internal :abbr:`ATC (Apache Traffic Control)` functions use this to uniquely identify a Delivery Service as opposed to the historically favored "ID". This string will become a part of the CDN service domain, which all together looks like: :file:`{Delivery Service Routing Name}.{Delivery Service xml_id}.{CDN Subdomain}.{CDN Domain}.{Top-Level Domain}`. Must be all lowercase, no spaces or special characters, but may contain dashes/hyphens\ [#xmlValid]_.

.. table:: Aliases

	+------+---------------------------------+------------------------+
	| Name | Use(s)                          | Type(s)                |
	+======+=================================+========================+
	| Key  | Traffic Portal tables and forms | unchanged (``string``) |
	+------+---------------------------------+------------------------+

.. parent-selection:

Parent Selection
----------------

Parameters in the Edge (child) profile that influence this feature:

+-----------------------------------------------+----------------+---------------+-------------------------------------------------------+
|                      Name                     |    Filename    |    Default    |                      Description                      |
+===============================================+================+===============+=======================================================+
| CONFIG proxy.config.                          | records.config | INT 1         | enable parent selection.  This is a required setting. |
| http.parent_proxy_routing_enable              |                |               |                                                       |
+-----------------------------------------------+----------------+---------------+-------------------------------------------------------+
| CONFIG proxy.config.                          | records.config | INT 1         | required for parent selection.                        |
| url_remap.remap_required                      |                |               |                                                       |
+-----------------------------------------------+----------------+---------------+-------------------------------------------------------+
| CONFIG proxy.config.                          | records.config | INT 0         | See                                                   |
| http.no_dns_just_forward_to_parent            |                |               |                                                       |
+-----------------------------------------------+----------------+---------------+-------------------------------------------------------+
| CONFIG proxy.config.                          | records.config | INT 1         |                                                       |
| http.uncacheable_requests_bypass_parent       |                |               |                                                       |
+-----------------------------------------------+----------------+---------------+-------------------------------------------------------+
| CONFIG proxy.config.                          | records.config | INT 1         |                                                       |
| http.parent_proxy_routing_enable              |                |               |                                                       |
+-----------------------------------------------+----------------+---------------+-------------------------------------------------------+
| CONFIG proxy.config.                          | records.config | INT 300       |                                                       |
| http.parent_proxy.retry_time                  |                |               |                                                       |
+-----------------------------------------------+----------------+---------------+-------------------------------------------------------+
| CONFIG proxy.config.                          | records.config | INT 10        |                                                       |
| http.parent_proxy.fail_threshold              |                |               |                                                       |
+-----------------------------------------------+----------------+---------------+-------------------------------------------------------+
| CONFIG proxy.config.                          | records.config | INT 4         |                                                       |
| http.parent_proxy.total_connect_attempts      |                |               |                                                       |
+-----------------------------------------------+----------------+---------------+-------------------------------------------------------+
| CONFIG proxy.config.                          | records.config | INT 2         |                                                       |
| http.parent_proxy.per_parent_connect_attempts |                |               |                                                       |
+-----------------------------------------------+----------------+---------------+-------------------------------------------------------+
| CONFIG proxy.config.                          | records.config | INT 30        |                                                       |
| http.parent_proxy.connect_attempts_timeout    |                |               |                                                       |
+-----------------------------------------------+----------------+---------------+-------------------------------------------------------+
| CONFIG proxy.config.                          | records.config | INT 0         |                                                       |
| http.forward.proxy_auth_to_parent             |                |               |                                                       |
+-----------------------------------------------+----------------+---------------+-------------------------------------------------------+
| CONFIG proxy.config.                          | records.config | INT 0         |                                                       |
| http.parent_proxy_routing_enable              |                |               |                                                       |
+-----------------------------------------------+----------------+---------------+-------------------------------------------------------+
| CONFIG proxy.config.                          | records.config | STRING        |                                                       |
| http.parent_proxy.file                        |                | parent.config |                                                       |
+-----------------------------------------------+----------------+---------------+-------------------------------------------------------+
| CONFIG proxy.config.                          | records.config | INT 3         |                                                       |
| http.parent_proxy.connect_attempts_timeout    |                |               |                                                       |
+-----------------------------------------------+----------------+---------------+-------------------------------------------------------+
| algorithm                                     | parent.config  | urlhash       | The algorithm to use.                                 |
+-----------------------------------------------+----------------+---------------+-------------------------------------------------------+


Parameters in the Mid (parent) profile that influence this feature:

+----------------+---------------+---------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|      Name      |    Filename   | Default |                                                                                 Description                                                                              |
+================+===============+=========+==========================================================================================================================================================================+
| domain_name    | CRConfig.json | -       | Only parents with the same value as the edge are going to be used as parents (to keep separation between CDNs)                                                           |
+----------------+---------------+---------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| weight         | parent.config | 1.0     | The weight of this parent, translates to the number of replicas in the consistent hash ring. This parameter only has effect with algorithm at the client set to          |
|                |               |         | "consistent_hash"                                                                                                                                                        |
+----------------+---------------+---------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| port           | parent.config | 80      | The port this parent is listening on as a forward proxy.                                                                                                                 |
+----------------+---------------+---------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| use_ip_address | parent.config | 0       | 1 means use IP(v4) address of this parent in the parent.config, 0 means use the host_name.domain_name concatenation.                                                     |
+----------------+---------------+---------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------+

.. _qstring-handling:

Qstring Handling
----------------

Delivery services have a Query String Handling option that, when set to ignore, will automatically add a regex remap to that delivery service's config.  There may be times this is not preferred, or there may be requirements for one delivery service or server(s) to behave differently.  When this is required, the psel.qstring_handling parameter can be set in either the delivery service profile or the server profile, but it is important to note that the server profile will override ALL delivery services assigned to servers with this profile parameter.  If the parameter is not set for the server profile but is present for the :term:`Delivery Service` profile, this will override the setting in the delivery service.  A value of "ignore" will not result in the addition of regex remap configuration.

+-----------------------+---------------+---------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|      Name             |    Filename   | Default |                                                                                    Description                                                                    |
+=======================+===============+=========+===================================================================================================================================================================+
| psel.qstring_handling | parent.config | -       | Sets qstring handling without the use of regex remap for a delivery service when assigned to a delivery service profile, and overrides qstring handling for all   |
|                       |               |         | :term:`Delivery Service`\ s for associated servers when assigned to a server profile. Value must be "consider" or "ignore".                                       |
+-----------------------+---------------+---------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------+

.. _multi-site-origin:

Multi Site Origin
-----------------

.. Note:: The configuration of this feature changed significantly between ATS version 5 and >= 6. Some configuration in Traffic Control is different as well. This documentation assumes ATS 6 or higher. See :ref:`multi-site-origin-qht` for more details.

Normally, the mid servers are not aware of any redundancy at the origin layer. With Multi Site Origin enabled this changes - Traffic Server (and Traffic Ops) are now made aware of the fact there are multiple origins, and can be configured to do more advanced failover and loadbalancing actions. A prerequisite for MSO to work is that the multiple origin sites serve identical content with identical paths, and both are configured to serve the same origin hostname as is configured in the deliveryservice `Origin Server Base URL` field. See the `Apache Traffic Server docs <https://docs.trafficserver.apache.org/en/latest/admin-guide/files/parent.config.en.html>`_ for more information on that cache's implementation.

With This feature enabled, origin servers (or origin server VIP names for a site) are going to be entered as servers in to the Traiffic Ops UI. Server type is "ORG".

Parameters in the mid profile that influence this feature:

+--------------------------------------------------------------------------+----------------+------------+----------------------------------------------------------------------------------------------------+
|                                   Name                                   |    Filename    |  Default   |                                            Description                                             |
+==========================================================================+================+============+====================================================================================================+
| CONFIG proxy.config. http.parent_proxy_routing_enable                    | records.config | INT 1      | enable parent selection.  This is a required setting.                                              |
+--------------------------------------------------------------------------+----------------+------------+----------------------------------------------------------------------------------------------------+
| CONFIG proxy.config. url_remap.remap_required                            | records.config | INT 1      | required for parent selection.                                                                     |
+--------------------------------------------------------------------------+----------------+------------+----------------------------------------------------------------------------------------------------+


Parameters in the deliveryservice profile that influence this feature:

+---------------------------------------------+----------------+-----------------+---------------------------------------------------------------------------------------------------------------------------------+
|                                   Name      |    Filename    |  Default        |                                                                         Description                                             |
+=============================================+================+=================+=================================================================================================================================+
| mso.parent_retry                            | parent.config  | \-              | Either ``simple_retry``, ``dead_server_retry`` or ``both``.                                                                     |
+---------------------------------------------+----------------+-----------------+---------------------------------------------------------------------------------------------------------------------------------+
| mso.algorithm                               | parent.config  | consistent_hash | The algorithm to use. ``consisten_hash``, ``strict``, ``true``, ``false``, or ``latched``.                                      |
|                                             |                |                 |                                                                                                                                 |
|                                             |                |                 | - ``consisten_hash`` - spreads requests across multiple parents simultaneously based on hash of content URL.                    |
|                                             |                |                 | - ``strict`` - strict Round Robin spreads requests across multiple parents simultaneously based on order of requests.           |
|                                             |                |                 | - ``true`` - same as strict, but ensures that requests from the same IP always go to the same parent if available.              |
|                                             |                |                 | - ``false`` - uses only a single parent at any given time and switches to a new parent only if the current parent fails.        |
|                                             |                |                 | - ``latched`` - same as false, but now, a failed parent will not be retried.                                                    |
+---------------------------------------------+----------------+-----------------+---------------------------------------------------------------------------------------------------------------------------------+
| mso.unavailable_server_retry_response_codes | parent.config  | "503"           | Quoted, comma separated list of HTTP status codes that count as a unavailable_server_retry_response_code.                       |
+---------------------------------------------+----------------+-----------------+---------------------------------------------------------------------------------------------------------------------------------+
| mso.max_unavailable_server_retries          | parent.config  | 1               | How many times an unavailable server will be retried.                                                                           |
+---------------------------------------------+----------------+-----------------+---------------------------------------------------------------------------------------------------------------------------------+
| mso.simple_retry_response_codes             | parent.config  | "404"           | Quoted, comma separated list of HTTP status codes that count as a simple retry response code.                                   |
+---------------------------------------------+----------------+-----------------+---------------------------------------------------------------------------------------------------------------------------------+
| mso.max_simple_retries                      | parent.config  | 1               | How many times a simple retry will be done.                                                                                     |
+---------------------------------------------+----------------+-----------------+---------------------------------------------------------------------------------------------------------------------------------+



see :ref:`multi-site-origin-qht` for a *quick how to* on this feature.

.. _regex-remap:

Regex Remap Expression
----------------------
The regex remap expression allows to to use a regex and resulting match group(s) in order to modify the request URIs that are sent to origin. For example: ::

	^/original/(.*) http://origin.example.com/remapped/$1

.. Note:: If **Query String Handling** is set to ``2 Drop at edge``, then you will not be allowed to save a regex remap expression, as dropping query strings actually relies on a regex remap of its own. However, if there is a need to both drop query strings **and** remap request URIs, this can be accomplished by setting **Query String Handling** to ``1 Do not use in cache key, but pass up to origin``, and then using a custom regex remap expression to do the necessary remapping, while simultaneously dropping query strings. The following example will capture the original request URI up to, but not including, the query string and then forward to a remapped URI: ::

	^/([^?]*).* http://origin.example.com/remapped/$1


.. _ds-regexp:

Delivery Service Regexp
-----------------------
This table defines how requests are matched to the delivery service. There are 3 type of entries possible here:

+---------------+----------------------------------------------------------------------+--------------+-----------+
|      Name     |                             Description                              |   DS Type    |   Status  |
+===============+======================================================================+==============+===========+
| HOST_REGEXP   | This is the regular expresion to match the host part of the URL.     | DNS and HTTP | Supported |
+---------------+----------------------------------------------------------------------+--------------+-----------+
| PATH_REGEXP   | This is the regular expresion to match the path part of the URL.     | HTTP         | Beta      |
+---------------+----------------------------------------------------------------------+--------------+-----------+
| HEADER_REGEXP | This is the regular expresion to match on any header in the request. | HTTP         | Beta      |
+---------------+----------------------------------------------------------------------+--------------+-----------+

The **Order** entry defines the order in which the regular expressions get evaluated. To support ``CNAMES`` from domains outside of the Traffic Control top level DNS domain, enter multiple ``HOST_REGEXP`` lines.

.. Note:: In most cases is is sufficient to have just one entry in this table that has a ``HOST_REGEXP`` Type, and Order ``0``. For the *movies* delivery service in the Kabletown CDN, the entry is simply single ``HOST_REGEXP`` set to ``.*\.movies\..*``. This will match every url that has a hostname that ends with ``movies.cdn1.kabletown.net``, since ``cdn1.kabletown.net`` is the Kabletown CDN's DNS domain.

.. index::
	Static DNS Entries

.. _static-dns:

Static DNS Entries
------------------
Static DNS entries allow you to create other names *under* the delivery service domain. You can enter any valid hostname, and create a CNAME, A or AAAA record for it by clicking the **Static DNS** button at the bottom of the delivery service details screen.

.. index::
	Server Assignments

.. _assign-edges:

Server Assignments
------------------
Click the **Server Assignments** button at the bottom of the screen to assign servers to this delivery service.  Servers can be selected by drilling down in a tree, starting at the profile, then the :term:`Cache Group`, and then the individual servers. Traffic Router will only route traffic for this delivery service to servers that are assigned to it.


.. _asn-czf:

The Coverage Zone File and ASN Table
------------------------------------
The Coverage Zone File (CZF) should contain a cachegroup name to network prefix mapping in the form:

.. code-block:: json

	{
		"coverageZones": {
			"cache-group-01": {
				"coordinates": {
					"latitude":  1.1,
					"longitude": 2.2
				},
				"network6": [
					"1234:5678::/64",
					"1234:5679::/64"
				],
				"network": [
					"192.168.8.0/24",
					"192.168.9.0/24"
				]
			},
			"cache-group-02": {
				"coordinates": {
					"latitude":  3.3,
					"longitude": 4.4
				},
				"network6": [
					"1234:567a::/64",
					"1234:567b::/64"
				],
				"network": [
					"192.168.4.0/24",
					"192.168.5.0/24"
				]
			}
		}
	}

The CZF is an input to the Traffic Control CDN, and as such does not get generated by Traffic Ops, but rather, it gets consumed by Traffic Router. Some popular IP management systems output a very similar file to the CZF but in stead of a cachegroup an ASN will be listed. Traffic Ops has the "Networks (ASNs)" view to aid with the conversion of files like that to a Traffic Control CZF file; this table is not used anywhere in Traffic Ops, but can be used to script the conversion using the API.

The script that generates the CZF file is not part of Traffic Control, since it is different for each situation.

.. note:: The ``"coordinates"`` section is optional and may be used by Traffic Router for localization in the case of a CZF "hit" where the zone name does not map to a :term:`Cache Group` name in Traffic Ops (i.e. Traffic Router will route to the closest :term:`Cache Group`\ (s) geographically).

.. _deep-czf:

The Deep Coverage Zone File
---------------------------
The Deep Coverage Zone File (DCZF) format is similar to the CZF format but adds a ``caches`` list under each ``deepCoverageZone``:

.. code-block:: json

	{
		"deepCoverageZones": {
			"location-01": {
				"coordinates": {
					"latitude":  5.5,
					"longitude": 6.6
				},
				"network6": [
					"1234:5678::/64",
					"1234:5679::/64"
				],
				"network": [
					"192.168.8.0/24",
					"192.168.9.0/24"
				],
				"caches": [
					"edge-01",
					"edge-02"
				]
			},
			"location-02": {
				"coordinates": {
					"latitude":  7.7,
					"longitude": 8.8
				},
				"network6": [
					"1234:567a::/64",
					"1234:567b::/64"
				],
				"network": [
					"192.168.4.0/24",
					"192.168.5.0/24"
				],
				"caches": [
					"edge-02",
					"edge-03"
				]
			}
		}
	}

Each entry in the ``caches`` list is the hostname of an edge cache registered in Traffic Ops which will be used for "deep" caching in that Deep Coverage Zone. Unlike a regular CZF, coverage zones in the DCZF do not map to a :term:`Cache Group` in Traffic Ops, so currently the deep coverage zone name only needs to be unique.

If the Traffic Router gets a DCZF "hit" for a requested :term:`Delivery Service` that has Deep Caching enabled, the client will be routed to an available "deep" cache from that zone's ``caches`` list.

.. note:: The ``"coordinates"`` section is optional.

.. [#xmlValid] Some things to consider when choosing an xml_id and routing name: the name should be descriptive and unique, but as brief as possible to avoid creating a monstrous :abbr:`FQDN (Fully Qualified Domain Name)`. Also, because these are combined to form an :abbr:`FQDN (Fully Qualified Domain Name)`, they should not contain any characters that are illegal for a DNS subdomain, e.g. ``.`` (period/dot). Finally, the restrictions on what characters are allowable (especially in xml_id) are, in general, **NOT** enforced by the :ref:`to-api`, so take care that the name is appropriate. See :rfc:`1035` for exact guidelines.
.. [#cardinality] In source code and :ref:`to-api` responses, the "Long Description" fields of a Delivery Service are "0-indexed" - hence the names differing slightly from the ones displayed in user-friendly UIs.
