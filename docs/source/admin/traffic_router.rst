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

*****************************
Traffic Router Administration
*****************************
.. contents::
	:depth: 2
	:backlinks: top

Installing Traffic Router
==========================
The following are requirements to ensure an accurate set up:

* CentOS 7
* 4 CPUs
* 8GB of RAM
* Successful install of Traffic Ops (usually on another machine)
* Successful install of Traffic Monitor (usually on another machine)
* Administrative access to Traffic Ops

.. Note:: Hardware requirements are generally doubled if :ref:`tr-DNSSEC` is enabled

#. If no suitable profile exists, create a new profile for Traffic Router ('Configure' -> 'Profiles' -> '+').

	.. warning:: Traffic Ops will *only* recognize a profile as assignable to a Traffic Router if its name starts with the prefix ``ccr-``. The reason for this is a legacy limitation related to the old name for Traffic Router (Comcast Cloud Router), and will (hopefully) be rectified in the future as the old Perl parts of Traffic Ops are re-written in Go.

#. Enter the Traffic Router server into Traffic Portal (or via the Traffic Ops API), assign to it a Traffic Router profile, and ensure that its status is set to ``ONLINE``.
#. Ensure the Fully Qualified Domain Name (FQDN) of the Traffic Router is resolvable in DNS. This FQDN must be resolvable by the clients expected to use this CDN.
#. Install a Traffic Router server package, either from source or via ``yum install traffic_router`` if it happens to be in your repositories (check with ``yum whatprovides traffic_router``).

	.. Note:: As of Traffic Control version 3.0, Traffic Router depends upon a package called ``tomcat``. This package should have been created when you built Traffic Router. If you get an error while installing the ``traffic_router`` package make sure that the ``tomcat`` package is available in your package repositories.

#. Edit ``/opt/traffic_router/conf/traffic_monitor.properties`` and specify the correct online Traffic Monitor(s) for your CDN. See :ref:`tr-config-files`

	:traffic_monitor.properties: URL that should normally point to this file. e.x. ``traffic_monitor.properties=file:/opt/traffic_router/conf/traffic_monitor.properties``

	:traffic_monitor.properties.reload.period: Period to wait (in milliseconds) between reloading this file. e.x. ``traffic_monitor.properties.reload.period=60000``


#. Start Traffic Router. This can be done by running ``systemctl start traffic_router`` as the root user (or with ``sudo``), and test DNS lookups against that server with e.g. ``dig`` or ``curl``. To restart Traffic Router, run ``systemctl restart traffic_router`` as the root user (or with ``sudo``). Also, because previously received CRConfigs will be cached, they need to be removed manually to actually be reloaded. This file should be located at ``/opt/traffic_router/db/cr-config.json``.

#. Snapshot CRConfig; See :ref:`snapshot-crconfig`

	.. Note:: Once the CRConfig is 'snapshotted', live traffic will be sent to the new Traffic Routers provided that their status has been set to ``ONLINE``.

#. Ensure that the parent domain (e.g.: ``cdn.local``) for the CDN's top level domain (e.g.: ``ciab.cdn.local``) contains a delegation (Name Server records) for the new Traffic Router, and that the value specified matches the FQDN used in above.

Configuring Traffic Router
==========================

.. Note:: Starting with Traffic Router 1.5, many of the configuration files under ``/opt/traffic_router/conf`` are only needed to override the default configuration values for Traffic Router. Most of the given default values will work well for any CDN. Critical values that must be changed are hostnames and credentials for communicating with other Traffic Control components such as Traffic Ops and Traffic Monitor.

.. Note:: Pre-existing installations that store configuration files under ``/opt/traffic_router/conf`` will still be used and honored for Traffic Router 1.5 onward.

.. Note:: Traffic Router 3.0 has been converted to a formal Tomcat instance, meaning that is now installed separately from the Tomcat servlet engine. The Traffic Router installation package contains all of the Traffic Router-specific software, configuration and startup scripts including some additional configuration files needed for Tomcat. These new configuration files can all be found in the ``/opt/traffic_router/conf`` directory and generally serve to override Tomcat's default settings.

For the most part, the configuration files and parameters that follow are used to get Traffic Router online and communicating with various Traffic Control components. Once Traffic Router is successfully communicating with Traffic Control, configuration should mostly be performed in Traffic Portal, and will be distributed throughout Traffic Control via the CRConfig snapshot process. See :ref:`snapshot-crconfig` for more information. Please see the parameter documentation for Traffic Router in the Using Traffic Ops guide documented under :ref:`ccr-profile` for parameters that influence the behavior of Traffic Router via the CRConfig.

.. _tr-config-files:

Configuration files
-------------------

+----------------------------+-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
|         File name          |                 Parameter                 |                                        Description                                    |                   Default Value                    |
+============================+===========================================+=======================================================================================+====================================================+
| traffic_monitor.properties | traffic_monitor.bootstrap.hosts           | Semicolon-delimited Traffic Monitor FQDNs - with port numbers as necessary            | N/A                                                |
|                            +-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
|                            | traffic_monitor.bootstrap.local           | Use only the Traffic Monitors specified in local configuration files                  | ``false``                                          |
|                            +-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
|                            | traffic_monitor.properties                | Path to the ``traffic_monitor.properties`` file; used internally to monitor the file  | ``/opt/traffic_router/traffic_monitor.properties`` |
|                            |                                           | for changes                                                                           |                                                    |
|                            +-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
|                            | traffic_monitor.properties.reload.period  | The interval in milliseconds for Traffic Router to wait between reloading this        | ``60000``                                          |
|                            |                                           | configuration file                                                                    |                                                    |
+----------------------------+-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
| dns.properties             | dns.tcp.port                              | TCP port that Traffic Router will use for incoming DNS requests                       | ``53``                                             |
|                            +-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
|                            | dns.tcp.backlog                           | Maximum length of the queue for incoming TCP connection requests                      | ``0``                                              |
|                            +-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
|                            | dns.udp.port                              | UDP port that Traffic Router will use for incoming DNS requests                       | ``53``                                             |
|                            +-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
|                            | dns.max-threads                           | Maximum number of threads used to process incoming DNS requests                       | ``1000``                                           |
|                            +-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
|                            | dns.zones.dir                             | Path to automatically generated zone files for reference                              | ``/opt/traffic_router/var/auto-zones``             |
+----------------------------+-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
| traffic_ops.properties     | traffic_ops.username                      | Username with which to access the APIs in Traffic Ops (must be in the ``admin`` role) | ``admin``                                          |
|                            +-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
|                            | traffic_ops.password                      | Password for the user specified in ``traffic_ops.username``                           | N/A                                                |
+----------------------------+-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
| cache.properties           | cache.geolocation.database                | Full path to the local copy of a GeoIP2 (usually MaxMind) binary database file        | ``/opt/traffic_router/db/GeoIP2-City.mmdb``        |
|                            +-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
|                            | cache.geolocation.database.refresh.period | The interval in milliseconds for Traffic Router to wait between polling for changes   | ``604800000``                                      |
|                            |                                           | to the GeoIP2 database                                                                |                                                    |
|                            +-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
|                            | cache.czmap.database                      | Full path to the local copy of the coverage zone file                                 | ``/opt/traffic_router/db/czmap.json``              |
|                            +-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
|                            | cache.czmap.database.refresh.period       | The interval in milliseconds for Traffic Router to wait between polling for a new     | ``10800000``                                       |
|                            |                                           | coverage zone file                                                                    |                                                    |
|                            +-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
|                            | cache.dczmap.database                     | Full path to the local copy of the deep coverage zone file                            | ``/opt/traffic_router/db/dczmap.json``             |
|                            +-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
|                            | cache.dczmap.database.refresh.period      | The interval in milliseconds for Traffic Router to wait between polling for a new     | ``10800000``                                       |
|                            |                                           | deep coverage zone file                                                               |                                                    |
|                            +-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
|                            | cache.health.json                         | Full path to the local copy of the health state                                       | ``/opt/traffic_router/db/health.json``             |
|                            +-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
|                            | cache.health.json.refresh.period          | The interval in milliseconds which Traffic Router will poll for a new health state    | ``1000``                                           |
|                            |                                           | file                                                                                  |                                                    |
|                            +-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
|                            | cache.config.json                         | Full path to the local copy of the CRConfig                                           | ``/opt/traffic_router/db/cr-config.json``          |
|                            +-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
|                            | cache.config.json.refresh.period          | The interval in milliseconds which Traffic Router will poll for a new CRConfig        | ``60000``                                          |
+----------------------------+-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
| startup.properties         | various parameters                        | This configuration is used by ``systemctl`` to set environment variables when the     | N/A                                                |
|                            |                                           | ``traffic_router`` service is started. It primarily consists of command line settings |                                                    |
|                            |                                           | for the Java process                                                                  |                                                    |
+----------------------------+-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
| log4j.properties           | various parameters                        | Configuration of ``log4j`` is                                                         | N/A                                                |
|                            |                                           | `documented on their site <http://logging.apache.org/log4j/2.x/index.html>`_; adjust  |                                                    |
|                            |                                           | as needed                                                                             |                                                    |
+----------------------------+-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
| server.xml                 | various parameters                        | Traffic Router specific configuration for Apache Tomcat. See the                      | N/A                                                |
|                            |                                           | `Apache Tomcat documentation <https://tomcat.apache.org/tomcat-8.5-doc/index.html>`_. |                                                    |
+----------------------------+-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+
| web.xml                    | various parameters                        | Default settings for all Web Applications running in the Traffic Router instance of   | N/A                                                |
|                            |                                           | Tomcat                                                                                |                                                    |
+----------------------------+-------------------------------------------+---------------------------------------------------------------------------------------+----------------------------------------------------+

.. _tr-dnssec:

DNSSEC
======

Overview
--------
Domain Name System Security Extensions (DNSSEC) is a set of extensions to DNS that provides a cryptographic mechanism for resolvers to verify the authenticity of responses served by an authoritative DNS server.

Several RFCs (4033, 4044, 4045) describe the low level details and define the extensions, RFC 7129 provides clarification around authenticated denial of existence of records, and finally RFC 6781 describes operational best practices for administering an authoritative DNSSEC enabled DNS server. The authenticated denial of existence RFC describes how an authoritative DNS server responds in NXDOMAIN and NODATA scenarios when DNSSEC is enabled.

Traffic Router currently supports DNSSEC with NSEC, however, NSEC3 and more configurable options are planned for the future.

Operation
---------
Upon startup or a configuration change, Traffic Router obtains keys from the 'keystore' API in Traffic Ops which returns key signing keys (KSK) and zone signing keys (ZSK) for each Delivery Service that is a sub-domain of the CDN's Top Level Domain (TLD) in addition to the keys for the CDN TLD itself. Each key has timing information that allows Traffic Router to determine key validity (expiration, inception, and effective dates) in addition to the appropriate Time To Live (TTL) to use for the DNSKEY record(s). All TTLs are configurable parameters; see the :ref:`ccr-profile` documentation for more information.

Once Traffic Router obtains the key data from the API, it converts each public key into the appropriate record types (DNSKEY, DS) to place in zones and uses the private key to sign zones. DNSKEY records are added to each Delivery Service's zone (e.g.: mydeliveryservice.ciab.cdn.local) for every valid key that exists, in addition to the CDN TLD's zone. A DS record is generated from each zone's KSK and is placed in the CDN TLD's zone (e.g.: ciab.cdn.local); the DS record for the CDN TLD must be placed in its parent zone, which is not managed by Traffic Control.

The DNSKEY to DS record relationship allows resolvers to validate signatures across zone delegation points. With Traffic Control, we control all delegation points below the CDN's TLD, **however, the DS record for the CDN TLD must be placed in the parent zone (e.g.: cdn.local), which is not managed by Traffic Control**. As such, the DS record must be placed in the parent zone prior to enabling DNSSEC, and prior to generating a new CDN KSK. Based on your deployment's DNS configuration, this might be a manual process or it might be automated. Either way, extreme care and diligence must be taken and knowledge of the management of the upstream zone is imperative for a successful DNSSEC deployment.

To enable DNSSEC for a CDN in Traffic Portal, Go to 'CDNs' from the sidebar and click on the desired CDN, then toggle the 'DNSSEC Enabled' field to 'true', and click on the green 'Update' button to save the changes.

Rolling Zone Signing Keys
-------------------------
Traffic Router currently follows the zone signing key pre-publishing operational best practice described in `section 4.1.1.1 of RFC 6781`_. Once DNSSEC is enabled for a CDN in Traffic Portal, key rolls are triggered by Traffic Ops via the automated key generation process, and Traffic Router selects the active zone signing keys based on the expiration information returned from the 'keystore' API of Traffic Ops.

.. _section 4.1.1.1 of RFC 6781: https://tools.ietf.org/html/rfc6781#section-4.1.1.1

Troubleshooting and Log Files
=============================
Traffic Router log files can be found under ``/opt/traffic_router/var/log`` and ``/opt/tomcat/logs``. Initialization and shutdown logs are in ``/opt/tomcat/logs/catalina[date].out``. Application related logging is in ``/opt/traffic_router/var/log/traffic_router.log``, while access logs are written to ``/opt/traffic_router/var/log/access.log``.

Event Log File Format
---------------------

Summary
"""""""

All access events to Traffic Router are logged to the file ``/opt/traffic_router/var/log/access.log``
This file grows up to 200MB and gets rolled into older log files, 10 log files total are kept (total of up to 2GB of logged events per Traffic Router instance)

Traffic Router logs access events in a format that largely following `ATS event logging format
<https://docs.trafficserver.apache.org/en/6.0.x/admin/event-logging-formats.en.html>`_

Message Format
""""""""""""""
- Except for the first item, each event that is logged is a series of space-separated key/value pairs.
- The first item is always the Unix epoch in seconds with a decimal field precision of up to milliseconds.
- Each key/value pair is in the form of ``unquoted_string="optionally quoted string"``
- Values that are quoted strings may contain whitespace characters.
- Values that are not quoted should not contains any whitespace characters.

.. Note:: Any value that is a single dash character or a dash character enclosed in quotes represents an empty value

Sample Message
""""""""""""""

Items within brackets below are detailed under the HTTP and DNS sections::

  144140678.000 qtype=DNS chi=192.168.10.11 ttms=789 [Fields Specific to the DNS request] rtype=CZ rloc="40.252611,58.439389" rdtl=- rerr="-" [Fields Specific to the DNS result]
  144140678.000 qtype=HTTP chi=192.168.10.11 ttms=789 [Fields Specific to the HTTP request] rtype=GEO rloc="40.252611,58.439389" rdtl=- rerr="-" [Fields Specific to the HTTP result]

.. Note:: The above message samples contain fields that are always present for every single access event to Traffic Router


Fields Always Present
"""""""""""""""""""""

+------+---------------------------------------------------------------------------------+------------------------------------------------------------------------------------+
|Name  |Description                                                                      |Data                                                                                |
+======+=================================================================================+====================================================================================+
|qtype |Whether the request was for DNS or HTTP                                          |Always DNS or HTTP                                                                  |
+------+---------------------------------------------------------------------------------+------------------------------------------------------------------------------------+
|chi   |The IP address of the requester                                                  |Depends on whether this was a DNS or HTTP request, see below sections               |
+------+---------------------------------------------------------------------------------+------------------------------------------------------------------------------------+
|ttms  |The amount of time in milliseconds it took Traffic Router to process the request |A number greater than or equal to zero                                              |
+------+---------------------------------------------------------------------------------+------------------------------------------------------------------------------------+
|rtype |Routing Result Type                                                              |One of ERROR, CZ, DEEP_CZ, GEO, MISS, STATIC_ROUTE, DS_REDIRECT, DS_MISS, INIT, FED |
+------+---------------------------------------------------------------------------------+------------------------------------------------------------------------------------+
|rloc  |GeoLocation of result                                                            |Latitude and Longitude in Decimal Degrees                                           |
+------+---------------------------------------------------------------------------------+------------------------------------------------------------------------------------+
|rdtl  |Result Details Associated with unusual conditions                                |One of DS_NOT_FOUND, DS_NO_BYPASS, DS_BYPASS, DS_CZ_ONLY, DS_CZ_BACKUP_CG           |
+------+---------------------------------------------------------------------------------+------------------------------------------------------------------------------------+
|rerr  |Message about internal Traffic Router Error                                      |String                                                                              |
+------+---------------------------------------------------------------------------------+------------------------------------------------------------------------------------+


``rtype`` Meanings
^^^^^^^^^^^^^^^^^^

:"-":          The request was not redirected. This is usually a result of a DNS request to the Traffic Router or an explicit denial for that request

:CZ:           The result was derived from Coverage Zone data based on the address in the ``chi`` field

:DEEP_CZ:      The result was derived from Deep Coverage Zone data based on the address in the ``chi`` field

:DS_MISS:      _*HTTP Only*_ No HTTP Delivery Service supports either this request's URL path or headers

:DS_REDIRECT:  The result is using the Bypass Destination configured for the matched Delivery Service when that Delivery Service is unavailable or does not have the requested resource

:ERROR:        An internal error occurred within Traffic Router, more details may be found in the ``rerr`` field

:FED:          _*DNS Only*_ The result was obtained through federated coverage zone data outside of any Delivery Service

:GEO:          The result was derived from geolocation service based on the address in the ``chi`` field

:GEO_REDIRECT: The request was redirected (302) based on the National Geo blocking (Geo Limit Redirect URL) configured on the Delivery Service

:MISS:         Traffic Router was unable to resolve a DNS request or find a cache for the requested resource

:RGALT:        The request was redirected (302) to the Regional Geo blocking URL. Regional Geo blocking is enabled on the Delivery Service and is configured through the ``regional_geoblock.polling.url`` setting for the Traffic Router profile

:RGDENY:       _*DNS Only*_ The result was obtained through federated coverage zone data outside of any Delivery Service The request was regionally blocked because there was no rule for the request made

:STATIC_ROUTE: _*DNS Only*_ No DNS Delivery Service supports the hostname portion of the requested url




``rdtl`` Meanings
^^^^^^^^^^^^^^^^^

:"-":                                  The request was not redirected. This is usually a result of a DNS request to the Traffic Router or an explicit denial for that request

:DS_BYPASS:                            Used Bypass Destination for Redirect of Delivery Service

:DS_CLIENT_GEO_UNSUPPORTED:            Traffic Router did not find a resource supported by coverage zone data and was unable to determine the geographic location of the requesting client

:DS_CZ_BACKUP_CG:                      Traffic Router found a backup cache via fall-back (CRconfig's ``edgeLocation``)  or via coordinates (CZF) configuration

:DS_CZ_ONLY:                           The selected Delivery Service only supports resource lookup based on Coverage Zone data

:DS_NO_BYPASS:                         No valid Bypass Destination is configured for the matched Delivery Service and the Delivery Service does not have the requested resource

:DS_NOT_FOUND:                         Always goes with ``rtypes`` STATIC_ROUTE and DS_MISS

:GEO_NO_CACHE_FOUND:                   Traffic Router could not find a resource via geographic location data based on the requesting client's location

:NO_DETAILS:                           This entry is for a standard request

:REGIONAL_GEO_ALTERNATE_WITHOUT_CACHE: This goes with the ``rtype`` RGDENY. The URL is being regionally blocked

:REGIONAL_GEO_NO_RULE:                 The request was blocked because there was no rule in the Delivery Service for the request


HTTP Specifics
--------------

Sample Message
::

  1452197640.936 qtype=HTTP chi=69.241.53.218 url="http://foo.mm-test.jenkins.cdnlab.comcast.net/some/asset.m3u8" cqhm=GET cqhv=HTTP/1.1 rtype=GEO rloc="40.252611,58.439389" rdtl=- rerr="-" pssc=302 ttms=0 rurl="http://odol-atsec-sim-114.mm-test.jenkins.cdnlab.comcast.net:8090/some/asset.m3u8" rh="Accept: */*" rh="myheader: asdasdasdasfasg"

.. table:: Request Fields

	+-----+-----------------------------------------------------------------------------------------------------------------------------------------+---------------------------------------------+
	|Name |Description                                                                                                                              |Data                                         |
	+=====+=========================================================================================================================================+=============================================+
	|url  |Requested URL with query string                                                                                                          |A URL String                                 |
	+-----+-----------------------------------------------------------------------------------------------------------------------------------------+---------------------------------------------+
	|cqhm |Http Method                                                                                                                              |e.g ``GET``, ``POST``                        |
	+-----+-----------------------------------------------------------------------------------------------------------------------------------------+---------------------------------------------+
	|cqhv |Http Protocol Version                                                                                                                    |e.g. ``HTTP/1.1``                            |
	+-----+-----------------------------------------------------------------------------------------------------------------------------------------+---------------------------------------------+
	|rh   |One or more of these key value pairs may exist in a logged event and are controlled by the configuration of the matched Delivery Service |Key/value pair of the format ``name: value`` |
	+-----+-----------------------------------------------------------------------------------------------------------------------------------------+---------------------------------------------+

.. table:: Response Fields

	+-----+----------------------------------------------------------+------------+
	|Name |Description                                               |Data        |
	+=====+==========================================================+============+
	|rurl |The resulting URL of the resource requested by the client |A URL String|
	+-----+----------------------------------------------------------+------------+

------------

DNS Specifics
-------------

Sample Message
::

  144140678.000 qtype=DNS chi=192.168.10.11 ttms=123 xn=65535 fqdn=www.example.com. type=A class=IN ttl=12345 rcode=NOERROR rtype=CZ rloc="40.252611,58.439389" rdtl=- rerr="-" ans="192.168.1.2 192.168.3.4 0:0:0:0:0:ffff:c0a8:102 0:0:0:0:0:ffff:c0a8:304"

.. _qname: http://www.zytrax.com/books/dns/ch15/#qname

.. _qtype: http://www.zytrax.com/books/dns/ch15/#qtype

.. table:: Request Fields

	+------+------------------------------------------------------------------+--------------------------------------------------------+
	|Name  |Description                                                       |Data                                                    |
	+======+==================================================================+========================================================+
	|xn    |The ID from the client DNS request header                         |a whole number between 0 and 65535 (inclusive)          |
	+------+------------------------------------------------------------------+--------------------------------------------------------+
	|fqdn  |The qname field from the client DNS request message (i.e. The     |A series of DNS labels/domains separated by '.'         |
	|      |fully qualified domain name the client is requesting be resolved) |characters and ending with a '.' character (see qname_) |
	+------+------------------------------------------------------------------+--------------------------------------------------------+
	|type  |The qtype field from the client DNS request message (i.e.         |Examples are A (IpV4), AAAA (IpV6), NS (Name Service),  |
	|      |the type of resolution that's requested such as IPv4, IPv6)       |SOA (Start of Authority), and CNAME, (see qtype_)       |
	+------+------------------------------------------------------------------+--------------------------------------------------------+
	|class |The qclass field from the client DNS request message (i.e. The    |Either IN (Internet resource) or ANY (Traffic router    |
	|      |class of resource being requested)                                |rejects requests with any other value of class)         |
	+------+------------------------------------------------------------------+--------------------------------------------------------+

.. table:: Response Fields

	+------+---------------------------------------------------------------------+-----------------------------------------------------+
	|Name  | Description                                                         | Data                                                |
	+======+=====================================================================+=====================================================+
	|ttl   | The 'time to live' in seconds for the answer provided by Traffic    |A whole number between 0 and 4294967295 (inclusive)  |
	|      | Router (clients can reliably use this answer for this long without  |                                                     |
	|      | re-querying traffic router)                                         |                                                     |
	+------+---------------------------------------------------------------------+-----------------------------------------------------+
	|rcode | The result code for the DNS answer provided by Traffic Router       | One of NOERROR (success), NOTIMP (request is not    |
	|      |                                                                     | NOTIMP (request is not  supported),                 |
	|      |                                                                     | REFUSED (request is refused to be answered), or     |
	|      |                                                                     | NXDOMAIN (the domain/name requested does not exist) |
	+------+---------------------------------------------------------------------+-----------------------------------------------------+

.. _tr-ngb:

GeoLimit Failure Redirect Feature
=================================

Overview
--------

This feature is also called 'National GeoBlock' (NGB).

In the past, if the Geolimit check fails (for example, the client IP is not in the 'US' region but the Geolimit is set to 'CZF + US'), the router will respond with ``503 Service Unavailable``, but with this feature, when the check fails, it will respond with ``302 Found`` if the redirect URL is set in the Delivery Service.

The Geolimit check will fail in the following scenarios:
	- When the GeoLimit is set to 'CZF + only' and the client IP is not in the the CZ file
	- When the GeoLimit is set to any region e.g. 'CZF + US' and the client IP is not in such region, and the client IP is not in the CZ file


Configuration
-------------

To enable the NGB feature, the DS must be configured with the proper redirect URL. The setting for this can be found by clicking on 'Advanced Options' at the bottom of a Delivery Service details page, and is specified by the 'Geo Limit Redirect URL' field. An individual Delivery Service details page can be viewed by clicking on the desired Delivery Service under 'Services' -> 'Delivery Services'. If no URL is put in this field, the feature is disabled.

The URL has 3 kinds of formats, which have different meanings:

URL with no domain
	If no domain is in the URL (e.g. 'vod/dance.mp4'), Traffic Router will try to find a proper cache server within the Delivery Service and return the redirect URL in the format: ``http://[cache server name].[Delivery Service's FQDN]/[configured relative path]``

URL with domain that matches with the Delivery Service
	For this URL, Traffic Router will also try to find a proper cache server within the Delivery Service and return a redirect URL in the format: ``http://[cache server name].[Delivery Service's FQDN]/[configured relative path]``

URL with domain that doesn't match with the Delivery Service
	Traffic Router will return the configured URL directly to the client.


.. _deep-cache:

Deep Caching - Deep Coverage Zone Topology
==========================================

Overview
--------

Deep Caching is a feature that enables clients to be routed to the closest possible "deep" Edge-tier caches on a per-Delivery Service basis. The term "deep" is used in the networking sense, meaning that the Edge-tier caches are located deep in the network where the number of network hops to a client is as minimal. This deep caching topology is desirable because storing content closer to the client gives better bandwidth savings, and sometimes the cost of bandwidth usage in the network outweighs the cost of adding storage. While it may not be feasible to cache an entire copy of the CDN's contents in every deep location (for the best possible bandwidth savings), storing just a relatively small amount of the CDN's most requested content can lead to very high bandwidth savings.

Getting Started
---------------

What you need:

#. Edge caches deployed in "deep" locations and registered in Traffic Ops
#. A Deep Coverage Zone File (DCZF) mapping these deep cache hostnames to specific network prefixes (see :ref:`deep-czf` for details)
#. Deep caching parameters in the Traffic Router Profile (see :ref:`ccr-profile` for details):

   - ``deepcoveragezone.polling.interval``
   - ``deepcoveragezone.polling.url``

#. Deep Caching enabled on one or more HTTP Delivery Services (i.e. 'Deep Caching' field on the Delivery Service details page (under 'Advanced Options') set to ALWAYS)

How it Works
------------

Deep Coverage Zone routing is very similar to that of regular Coverage Zone routing, except that the DCZF is preferred over the regular CZF for Delivery Services with Deep Caching (DC) enabled. If the client requests a DC-enabled Delivery Service and their IP address gets a "hit" in the DCZF, Traffic Router will attempt to route that client to one of the available deep caches in the client's corresponding zone. If there are no deep caches available for a client's request, Traffic Router will fall back to the regular CZF and continue regular CZF routing from there.


.. _tr-steering:

Steering Feature
================

Overview
--------
A Steering Delivery Service is a Delivery Service that is used to route a client to another Delivery Service. The Type of a Steering Delivery Service is either STEERING or CLIENT_STEERING. A Steering Delivery Service will have target Delivery Services configured for it with weights assigned to them. Traffic Router uses the weights to make a consistent hash ring which it then uses to make sure that requests are routed to a target based on the configured weights. This consistent hash ring is separate from the consistent hash ring used in cache selection.

Special regular expressions - referred to as 'filters' - can also be configured for target Delivery Services to pin traffic to a specific Delivery Service. For example, if a filter called ``.*/news/.*`` for a target called 'target-ds-1' is created, any requests to Traffic Router with 'news' in them will be routed to 'target-ds-1'. This will happen regardless of the configured weights.

Some other points of interest:

- Steering is currently only available for HTTP Delivery Services that are a part of the same CDN.
- A new role called STEERING has been added to the Traffic Ops database. Only users with Admin or Steering privileges can modify steering assignments for a Delivery Service.
- A new API has been created in Traffic Ops under ``/internal``. A Steering user can either directly access this API to modify assignments, or use the Traffic Portal UI ('View Targets' under the 'More' drop-down menu on a Steering Delivery Service's details page), however a filter can only be created via the API.
- Traffic Router uses the steering API in Traffic Ops to poll for steering assignments, the assignments are then used when routing traffic.

A couple simple use-cases for Steering are:

- Migrating traffic from one Delivery Service to another over time.
- Trying out new functionality for a subset of traffic with an experimental Delivery Service.
- Load balancing between Delivery Services.


The Difference Between STEERING and CLIENT_STEERING
---------------------------------------------------

The only difference between the STEERING and CLIENT_STEERING Delivery Service Types is that CLIENT_STEERING explicitly allows a client to bypass Steering by choosing a destination Delivery Service. A client can accomplish this by providing the ``X-TC-Steering-Option`` HTTP header with a value of the ``xml_id`` of the target Delivery Service to which they desire to be routed. When Traffic Router receives this header it will route to the requested target Delivery Service regardless of weight configuration. This header is ignored by STEERING Delivery Services.


Configuration
-------------

The following needs to be completed for Steering to work correctly:

#. Two target Delivery Services are created in Traffic Ops. They must both be HTTP Delivery Services part of the same CDN.
#. A Delivery Service with type STEERING or CLIENT_STEERING is created in Traffic Portal.
#. Target Delivery Services are assigned to the Steering Delivery Service using Traffic Portal.
#. A user with the role of Steering is created.
#. The Steering user assigns weights to the target Delivery Services.
#. If desired, the Steering user can create filters for the target Delivery Services.

For more information see the `Steering how-to guide <quick_howto/steering.html>`_.

HTTPS for HTTP Delivery Services
================================

Starting with version 1.7 Traffic Router added the ability to allow HTTPS traffic between itself and clients on a per-HTTP Delivery Service basis.

.. Note:: As of version 3.0 Traffic Router has been integrated with native OpenSSL. This makes establishing HTTPS connections to Traffic Router much less expensive than previous versions. However establishing an HTTPS connection is more computationally demanding than an HTTP connection. Since each client will in turn get redirected to ATS, Traffic Router is most always creating a new HTTPS connection for all HTTPS traffic. It is likely to mean that an existing Traffic Router may have some decrease in performance if you wish to support a lot of HTTPS traffic. As noted for DNSSEC, you may need to plan to scale Traffic Router vertically and/or horizontally to handle the new load.

The HTTPS set up process is:

#. Select one of '1 - HTTPS', '2 - HTTP AND HTTPS', or '3 - HTTP TO HTTPS' for the Delivery Service
#. Generate private keys for the Delivery Service using a wildcard domain such as ``*.my-delivery-service.my-cdn.example.com``
#. Obtain and import signed certificate chain
#. Snapshot CRConfig

Clients may make HTTPS requests to Delivery Services only after Traffic Router receives the certificate chain from Traffic Ops and the new CRConfig.

Protocol Options
----------------

HTTP
	Any secure client will get an SSL handshake error. Non-secure clients will experience the same behavior as prior to 1.7
HTTPS
	Traffic Router will only redirect (send a ``302 Found`` response) to clients communicating with a secure connection, all other clients will receive a ``503 Service Unavailable`` response
HTTP AND HTTPS
	Traffic Router will redirect both secure and non-secure clients
HTTP TO HTTPS
	Traffic Router will redirect non-secure clients with a ``302 Found`` response and a location that is secure (i.e. an ``https://`` URL instead of an ``http://`` URL), while secure clients will be redirected immediately to an appropriate target or cache server.

Certificate Retrieval
---------------------

.. Warning:: If you have HTTPS Delivery Services in your CDN, Traffic Router will not accept **any** connections until it is able to fetch certificates from Traffic Ops and load them into memory. Traffic Router does not persist certificates to the Java Keystore or anywhere else.

Traffic Router fetches certificates into memory:

* At startup time
* When it receives a new CRConfig
* Once an hour starting whenever the most recent of the last of the above occurred

.. Note:: To adjust the frequency at which Traffic Router fetches certificates add the parameter ``certificates.polling.interval`` to CRConfig and set it to the desired duration in milliseconds.

.. Note:: Taking a snapshot of CRConfig may be used at times to avoid waiting the entire polling cycle for a new set of certificates.

.. Warning:: If a snapshot of CRConfig is made that involves a Delivery Service missing its certificates, Traffic Router will ignore **ALL** changes in that CRConfig until one of the following occurs:

	* It receives certificates for that Delivery Service
	* Another snapshot of CRConfig is created and the Delivery Service without certificates is changed so its HTTP protocol is set to 'http'

Certificate Chain Ordering
--------------------------

The ordering of certificates within the certificate bundle matters. It must be:

#. Primary Certificate (e.g. the one created for ``*.my-delivery-service.my-cdn.example.com``)
#. Intermediate Certificate(s)
#. Root Certificate from a Certificate Authority (CA) (optional)

.. Warning:: If something is wrong with the certificate chain (e.g. the order of the certificates is backwards or for the wrong domain) the client will get an SSL handshake. Inspection of ``/opt/tomcat/logs/catalina.log`` is likely to yield information to reveal this.

To see the ordering of certificates you may have to manually split up your certificate chain and use ``openssl`` on each individual certificate

Suggested Way of Setting up an HTTPS Delivery Service
-----------------------------------------------------

Assuming you have already created a Delivery Service which you plan to modify to use HTTPS, do the following in Traffic Portal:

#. Select one of '1 - HTTPS', '2 - HTTP AND HTTPS', or '3 - HTTP TO HTTPS' for the protocol field of a Delivery Service and click the 'Update' button
#. Under the 'More' drop-down menu, click 'Manage SSL Keys'
#. Again under the 'More' drop-down menu, click 'Generate SSL Keys'
#. Fill out the form and click on the green 'Generate Keys' button, then confirm that you want to make these changes
#. Copy the contents of the Certificate Signing Request field and save it locally
#. Go back and select 'HTTP' for the protocol field of the Delivery Service and click 'Save' (to avoid preventing other CRConfig updates from being blocked by Traffic Router)
#. Follow your standard procedure for obtaining your signed certificate chain from a CA
#. After receiving your certificate chain import it into Traffic Ops
#. Edit the Delivery Service
#. Restore your original choice for the protocol field and click save
#. Click 'Manage SSL Keys'
#. Paste your key information into the appropriate fields
#. Click the green 'Update Keys' button
#. Take a new snapshot of CRConfig

Once this is done you should be able to verify that you are being correctly redirected by Traffic Router using e.g. ``curl`` commands to HTTPS destinations on your Delivery Service.

Router Load Testing
===================

The Traffic Router load testing tool is located in the `Traffic Control repository under ``test/router`` <https://github.com/apache/trafficcontrol/tree/master/test/router>`_. It can be used to simulate a mix of HTTP and HTTPS traffic for a CDN by choosing the number of HTTP Delivery Services and the number HTTPS Delivery Services the test will exercise.

There are 2 parts to the load test:

* A web server that makes the actual requests and takes commands to fetch data from the CDN, start the test, and return current results.
* A web page that's used to run the test and see the results.

Running the Load Tests
----------------------

#. First, clone the `Traffic Control repository <https://github.com/apache/trafficcontrol>`_.
#. You will need to make sure you have a CA file on your machine
#. The web server is a Go program, set your ``GOPATH`` environment variable appropriately (we suggest ``$HOME/go`` or ``$HOME/src``)
#. Open a terminal emulator and navigate to the ``test/router/server`` directory inside of the cloned repository
#. Execute the server binary by running ``go run server.go``
#. Using your web browser of choice, open the file ``test/router/index.html``
#. Authenticate against a Traffic Ops host - this should be a nearly instantaneous operation - you can watch the output from ``server.go`` for feedback
#. Enter the Traffic Ops host in the second form and click the button to get a list of CDN's
#. Wait for the web page to show a list of CDN's under the above form, this may take several seconds
#. The List of CDN's will display the number of HTTP- and HTTPS-capable Delivery Services that may be exercised
#. Choose the CDN you want to exercise from the drop-down menu
#. Fill out the rest of the form, enter appropriate numbers for each HTTP and HTTPS delivery services
#. Click Run Test
#. As the test runs the web page will occasionally report results including running time, latency, and throughput

Tuning Recommendations
======================

The following is an example of the command line parameters set in ``/opt/traffic_router/conf/startup.properties`` that has been tested on a multi-core server running under HTTPS load test requests. This is following the general recommendation to use the G1 garbage collector for JVM applications running on multi-core machines. In addition to using the G1 garbage collector the ``InitiatingHeapOccupancyPercent`` was lowered to run garbage collection more frequently which improved overall throughput for Traffic Router and reduced 'Stop the World' garbage collection. Note that any environment variable settings in this file will override those
set in ``/lib/systemd/system/traffic_router.service``.

.. code-block:: bash

	CATALINA_OPTS="\
  	-server -Xms2g -Xmx8g \
  	-Dlog4j.configuration=file://$CATALINA_BASE/conf/log4j.properties \
  	-Djava.library.path=/usr/lib64 \
  	-XX:+UseG1GC \
  	-XX:+UnlockExperimentalVMOptions \
  	-XX:InitiatingHeapOccupancyPercent=30"
