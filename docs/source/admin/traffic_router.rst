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

* CentOS 6
* 4 vCPUs
* 8GB RAM
* Successful install of Traffic Ops
* Successful install of Traffic Monitor
* Administrative access to Traffic Ops

.. Note:: Hardware requirements are generally doubled if DNSSEC is enabled

1. If no suitable profile exists, create a new profile for Traffic Router.

2. Enter the Traffic Router server into Traffic Ops, assign it to a Traffic Router profile, and ensure that its status is set to ``ONLINE``.
3. Ensure the FQDN of the Traffic Router is resolvable in DNS. This FQDN must be resolvable by the clients expected to use this CDN.
4. Install a traffic router: ``sudo yum install traffic_router``.
5. Edit ``/opt/traffic_router/conf/traffic_monitor.properties`` and specify the correct online Traffic Monitor(s) for your CDN. See :ref:`rl-tr-config-files`
	# traffic_monitor.properties: url that should normally point to this file
	traffic_monitor.properties=file:/opt/traffic_router/conf/traffic_monitor.properties

	# Frequency for reloading this file
	# traffic_monitor.properties.reload.period=60000


6. Start Tomcat: ``sudo service tomcat start``, and test lookups with dig and curl against that server.
	To restart, ``sudo service tomcat stop``, kill the traffic router process, and ``sudo service tomcat start``
	Also, crconfig previously recieved will be cached, and needs to be removed manually to actually be reloaded /opt/traffic_router/db/cr-config.json
7. Snapshot CRConfig; See :ref:`rl-snapshot-crconfig`

..  Note:: Once the CRConfig is snapshotted, live traffic will be sent to the new Traffic Routers provided that their status is set to ``ONLINE``.

8. Ensure that the parent domain (e.g.: kabletown.net) for the CDN's top level domain (e.g.: cdn.kabletown.net) contains a delegation (NS records) for the new Traffic Router, and that the value specified matches the FQDN used in step 3.

Configuring Traffic Router
==========================

.. Note:: Starting with Traffic Router 1.5, many of the configuration files under ``/opt/traffic_router/conf`` are only needed to override the default configuration values for Traffic Router. Most of the given default values will work well for any CDN. Critical values that must be changed are hostnames and credentials for communicating with other Traffic Control components such as Traffic Ops and Traffic Monitor.

.. Note:: Pre-existing installations having configuration files in ``/opt/traffic_router/conf`` will still be used and honored for Traffic Router 1.5 and onward.

For the most part, the configuration files and parameters that follow are used to get Traffic Router online and communicating with various Traffic Control components. Once Traffic Router is successfully communicating with Traffic Control, configuration is mostly performed in Traffic Ops, and is distributed throughout Traffic Control via the CRConfig snapshot process. See :ref:`rl-snapshot-crconfig` for more information. Please see the parameter documentation for Traffic Router in the Using Traffic Ops guide documented under :ref:`rl-ccr-profile` for parameters that influence the behavior of Traffic Router via the CRConfig.

.. _rl-tr-config-files:

Configuration files
-------------------

+----------------------------+-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|         File name          |                 Parameter                 |                                             Description                                             |                   Default Value                   |
+============================+===========================================+=====================================================================================================+===================================================+
| traffic_monitor.properties | traffic_monitor.bootstrap.hosts           | Traffic Monitor FQDNs and port if necessary, separated by a semicolon (;)                           | N/A                                               |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | traffic_monitor.bootstrap.local           | Use only the Traffic Monitors specified in config file                                              | false                                             |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | traffic_monitor.properties                | Path to the traffic_monitor.properties file; used internally to monitor the file for changes        | /opt/traffic_router/traffic_monitor.properties    |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | traffic_monitor.properties.reload.period  | The interval in milliseconds which Traffic Router will reload this configuration file               | 60000                                             |
+----------------------------+-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
| dns.properties             | dns.tcp.port                              | TCP port that Traffic Router will use for incoming DNS requests                                     | 53                                                |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | dns.tcp.backlog                           | Maximum length of the queue for incoming TCP connection requests                                    | 0                                                 |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | dns.udp.port                              | UDP port that Traffic Router will use for incoming DNS requests                                     | 53                                                |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | dns.max-threads                           | Maximum number of threads used to process incoming DNS requests                                     | 1000                                              |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | dns.zones.dir                             | Path to auto generated zone files for reference                                                     | /opt/traffic_router/var/auto-zones                |
+----------------------------+-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
| traffic_ops.properties     | traffic_ops.username                      | Username to access the APIs in Traffic Ops (must be in the admin role)                              | admin                                             |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | traffic_ops.password                      | Password for the user specified in traffic_ops.username                                             | N/A                                               |
+----------------------------+-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
| cache.properties           | cache.geolocation.database                | Full path to the local copy of the MaxMind geolocation binary database file                         | /opt/traffic_router/db/GeoIP2-City.mmdb           |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | cache.geolocation.database.refresh.period | The interval in milliseconds which Traffic Router will poll for a new geolocation database          | 604800000                                         |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | cache.czmap.database                      | Full path to the local copy of the coverage zone file                                               | /opt/traffic_router/db/czmap.json                 |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | cache.czmap.database.refresh.period       | The interval in milliseconds which Traffic Router will poll for a new coverage zone file            | 10800000                                          |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | cache.dczmap.database                     | Full path to the local copy of the deep coverage zone file                                          | /opt/traffic_router/db/dczmap.json                |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | cache.dczmap.database.refresh.period      | The interval in milliseconds which Traffic Router will poll for a new deep coverage zone file       | 10800000                                          |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | cache.health.json                         | Full path to the local copy of the health state                                                     | /opt/traffic_router/db/health.json                |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | cache.health.json.refresh.period          | The interval in milliseconds which Traffic Router will poll for a new health state file             | 1000                                              |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | cache.config.json                         | Full path to the local copy of the CRConfig                                                         | /opt/traffic_router/db/cr-config.json             |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | cache.config.json.refresh.period          | The interval in milliseconds which Traffic Router will poll for a new CRConfig                      | 60000                                             |
+----------------------------+-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
| log4j.properties           | various parameters                        | Configuration of log4j is documented on their site; adjust as necessary based on needs              | N/A                                               |
+----------------------------+-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+

.. _rl-tr-dnssec:

DNSSEC
======

Overview
--------
Domain Name System Security Extensions (DNSSEC) is a set of extensions to DNS that provides a cryptographic mechanism for resolvers to verify the authenticity of responses served by an authoritative DNS server.

Several RFCs (4033, 4044, 4045) describe the low level details and define the extensions, RFC 7129 provides clarification around authenticated denial of existence of records, and finally RFC 6781 describes operational best practices for administering an authoritative DNSSEC enabled DNS server. The authenticated denial of existence RFC describes how an authoritative DNS server responds in NXDOMAIN and NODATA scenarios when DNSSEC is enabled.

Traffic Router currently supports DNSSEC with NSEC, however, NSEC3 and more configurable options will be provided in the future.

Operation
---------
Upon startup or a configuration change, Traffic Router obtains keys from the keystore API in Traffic Ops which returns key signing keys (KSK) and zone signing keys (ZSK) for each delivery service that is a subdomain off the CDN's top level domain (TLD), in addition to the keys for the CDN TLD itself. Each key has timing information that allows Traffic Router to determine key validity (expiration, inception, and effective dates) in addition to the appropriate TTL to use for the DNSKEY record(s).  All TTLs are configurable parameters; see the :ref:`rl-ccr-profile` documentation for more information.

Once Traffic Router obtains the key data from the API, it converts each public key into the appropriate record types (DNSKEY, DS) to place in zones and uses the private key to sign zones. DNSKEY records are added to each delivery service's zone (e.g.: mydeliveryservice.cdn.kabletown.net) for every valid key that exists, in addition to the CDN TLD's zone. A DS record is generated from each zone's KSK and is placed in the CDN TLD's zone (e.g.: cdn.kabletown.net); the DS record for the CDN TLD must be placed in its parent zone, which is not managed by Traffic Control.

The DNSKEY to DS record relationship allows resolvers to validate signatures across zone delegation points; with Traffic Control, we control all delegation points below the CDN's TLD, **however, the DS record for the CDN TLD must be placed in the parent zone (e.g.: kabletown.net), which is not managed by Traffic Control**. As such, the DS record (available in the Traffic Ops DNSSEC administration UI) must be placed in the parent zone prior to enabling DNSSEC, and prior to generating a new CDN KSK. Based on your deployment's DNS configuration, this might be a manual process or it might be automated; either way, extreme care and diligence must be taken and knowledge of the management of the upstream zone is imperative for a successful DNSSEC deployment.

Rolling Zone Signing Keys
-------------------------
Traffic Router currently follows the zone signing key pre-publishing operational best practice described in `section 4.1.1.1 of RFC 6781`_. Once DNSSEC is enabled for a CDN in Traffic Ops, key rolls are triggered via Traffic Ops via the automated key generation process, and Traffic Router selects the active zone signing keys based on the expiration information returned from the keystore API in Traffic Ops.

.. _section 4.1.1.1 of RFC 6781: https://tools.ietf.org/html/rfc6781#section-4.1.1.1

Troubleshooting and log files
=============================
Traffic Router log files are in ``/opt/traffic_router/var/log``, and Tomcat log files are in ``/opt/tomcat/logs``. Application related logging is in ``/opt/traffic_router/var/log/traffic_router.log``, while access logs are written to ``/opt/traffic_router/var/log/access.log``.

Event Log File Format
=====================

Summary
-------

All access events to Traffic Router are logged to the file ``/opt/traffic_router/var/log/access.log``
This file grows up to 200Mb and gets rolled into older log files, 10 log files total are kept (total of up to 2Gb of logged events per traffic router)

Traffic Router logs access events in a format that largely following `ATS event logging format
<https://docs.trafficserver.apache.org/en/6.0.x/admin/event-logging-formats.en.html>`_

--------------

Sample Message
--------------

Items within brackets below are detailed under the HTTP and DNS sections
::

  144140678.000 qtype=DNS chi=192.168.10.11 ttms=789 [Fields Specific to the DNS request] rtype=CZ rloc="40.252611,58.439389" rdtl=- rerr="-" [Fields Specific to the DNS result]
  144140678.000 qtype=HTTP chi=192.168.10.11 ttms=789 [Fields Specific to the HTTP request] rtype=GEO rloc="40.252611,58.439389" rdtl=- rerr="-" [Fields Specific to the HTTP result]

.. Note:: The above message samples contain fields that are always present for every single access event to Traffic Router

**Message Format**
- Each event that is logged is a series of space separated key value pairs except for the first item.
- The first item is always the epoch in seconds with a decimal field precision of up to milliseconds
- Each key value pair is in the form of unquoted string, equals character, optionally quoted string
- Values that are quoted strings may contain space characters
- Values that are not quoted should not contains any space characters

.. Note:: Any value that is a single dash character or a dash character enclosed in quotes represents an empty value

--------

Fields Always Present
---------------------

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


**rtype meanings**

+-------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|Name         |Meaning                                                                                                                                                                 |
+=============+========================================================================================================================================================================+
|ERROR        |An internal error occurred within Traffic Router, more details may be found in the rerr field                                                                           |
+-------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|CZ           |The result was derived from Coverage Zone data based on the address in the chi field                                                                                    |
+-------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|DEEP_CZ      |The result was derived from Deep Coverage Zone data based on the address in the chi field                                                                               |
+-------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|GEO          |The result was derived from geolocation service based on the address in the chi field                                                                                   |
+-------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|MISS         |Traffic Router was unable to resolve a DNS request or find a cache for the requested resource                                                                           |
+-------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|STATIC_ROUTE |_*DNS Only*_ No DNS Delivery Service supports the hostname portion of the requested url                                                                                 |
+-------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|DS_MISS      |_*HTTP Only*_ No HTTP Delivery Service supports either this request's URL path or headers                                                                               |
+-------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|DS_REDIRECT  |The result is using the Bypass Destination configured for the matched Delivery Service when that Delivery Service is unavailable or does not have the requested resource|
+-------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|FED          |_*DNS Only*_ The result was obtained through federated coverage zone data outside of any delivery service                                                               |
+-------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|GEO_REDIRECT |The request was redirected (302) based on the National Geo blocking (Geo Limit Redirect URL) configured on the Delivery Service.                                        |
+-------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|RGALT        |The request was redirected (302) to the Regional Geo blocking URL. Regional Geo blocking is enabled on the Delivery Service and is configured through the               |
|             |regional_geoblock.polling.url setting for the Traffic Router profile.                                                                                                   |
+-------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|RGDENY       |_*DNS Only*_ The result was obtained through federated coverage zone data outside of any delivery service The request was regionally blocked because there was no rule  |
|             |for the request made.                                                                                                                                                   |
+-------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| "-"         |The request was not redirected. This is usually a result of a DNS request to the Traffic Router or an explicit denial for that request.                                 |
+-------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------+


**rdtl meanings**

+---------------------------------------+--------------------------------------------------------------------------------------------------------------------------------------------+
|Name                                   |Meaning                                                                                                                                     |
+=======================================+============================================================================================================================================+
|DS_NOT_FOUND                           |Always goes with rtypes STATIC_ROUTE and DS_MISS                                                                                            |
+---------------------------------------+--------------------------------------------------------------------------------------------------------------------------------------------+
|DS_BYPASS                              |Used Bypass Destination for Redirect of Delivery Service                                                                                    |
+---------------------------------------+--------------------------------------------------------------------------------------------------------------------------------------------+
|DS_NO_BYPASS                           |No valid Bypass Destination is configured for the matched Delivery Service and the delivery service does not have the requested resource    |
+---------------------------------------+--------------------------------------------------------------------------------------------------------------------------------------------+
|DS_CZ_ONLY                             |The selected Delivery Service only supports resource lookup based on Coverage Zone data                                                     |
+---------------------------------------+--------------------------------------------------------------------------------------------------------------------------------------------+
|DS_CLIENT_GEO_UNSUPPORTED              |Traffic Router did not find a resource supported by coverage zone data and was unable to determine the geolocation of the requesting client |
+---------------------------------------+--------------------------------------------------------------------------------------------------------------------------------------------+
|GEO_NO_CACHE_FOUND                     |Traffic Router could not find a resource via geolocation data based on the requesting client's geolocation                                  |
+---------------------------------------+--------------------------------------------------------------------------------------------------------------------------------------------+
|NO_DETAILS                             |This entry is for a standard request.                                                                                                       |
+---------------------------------------+--------------------------------------------------------------------------------------------------------------------------------------------+
|REGIONAL_GEO_ALTERNATE_WITHOUT_CACHE   |This goes with the rtype RGDENY. The URL is being regionally Geo blocked.                                                                   |
+---------------------------------------+--------------------------------------------------------------------------------------------------------------------------------------------+
|REGIONAL_GEO_NO_RULE                   |The request was blocked because there was no rule in the Delivery Service for the request.                                                  |
+---------------------------------------+--------------------------------------------------------------------------------------------------------------------------------------------+
| "-"                                   |The request was not redirected. This is usually a result of a DNS request to the Traffic Router or an explicit denial for that request.     |
+---------------------------------------+--------------------------------------------------------------------------------------------------------------------------------------------+
|DS_CZ_BACKUP_CG                        |Traffic Router found a backup cache via fallback (cr-config's edgeLocation)  / coordinates (CZF) configuration                              |
+---------------------------------------+--------------------------------------------------------------------------------------------------------------------------------------------+


---------------

HTTP Specifics
--------------

Sample Message
::

  1452197640.936 qtype=HTTP chi=69.241.53.218 url="http://foo.mm-test.jenkins.cdnlab.comcast.net/some/asset.m3u8" cqhm=GET cqhv=HTTP/1.1 rtype=GEO rloc="40.252611,58.439389" rdtl=- rerr="-" pssc=302 ttms=0 rurl="http://odol-atsec-sim-114.mm-test.jenkins.cdnlab.comcast.net:8090/some/asset.m3u8" rh="Accept: */*" rh="myheader: asdasdasdasfasg"

**Request Fields**

+-----+-----------------------------------------------------------------------------------------------------------------------------------------+-------------------------------------------+
|Name |Description                                                                                                                              |Data                                       |
+=====+=========================================================================================================================================+===========================================+
|url  |Requested URL with query string                                                                                                          |String                                     |
+-----+-----------------------------------------------------------------------------------------------------------------------------------------+-------------------------------------------+
|cqhm |Http Method                                                                                                                              |e.g GET, POST                              |
+-----+-----------------------------------------------------------------------------------------------------------------------------------------+-------------------------------------------+
|cqhv |Http Protocol Version                                                                                                                    |e.g. HTTP/1.1                              |
+-----+-----------------------------------------------------------------------------------------------------------------------------------------+-------------------------------------------+
|rh   |One or more of these key value pairs may exist in a logged event and are controlled by the configuration of the matched Delivery Service |Key value pair of the format "name: value" |
+-----+-----------------------------------------------------------------------------------------------------------------------------------------+-------------------------------------------+

**Response Fields**

+-----+----------------------------------------------------------+------------+
|Name |Description                                               |Data        |
+=====+==========================================================+============+
|rurl |The resulting url of the resource requested by the client |A URL String|
+-----+----------------------------------------------------------+------------+

------------

DNS Specifics
-------------

Sample Message
::

  144140678.000 qtype=DNS chi=192.168.10.11 ttms=123 xn=65535 fqdn=www.example.com. type=A class=IN ttl=12345 rcode=NOERROR rtype=CZ rloc="40.252611,58.439389" rdtl=- rerr="-" ans="192.168.1.2 192.168.3.4 0:0:0:0:0:ffff:c0a8:102 0:0:0:0:0:ffff:c0a8:304"

**Request Fields**

.. _qname: http://www.zytrax.com/books/dns/ch15/#qname

.. _qtype: http://www.zytrax.com/books/dns/ch15/#qtype

+------+------------------------------------------------------------------+--------------------------------------------------------+
|Name  |Description                                                       |Data                                                    |
+======+==================================================================+========================================================+
|xn    |The ID from the client DNS request header                         |a number from 0 to 65535                                |
+------+------------------------------------------------------------------+--------------------------------------------------------+
|fqdn  |The qname field from the client DNS request message (i.e. The     |A series of DNS labels/domains separated by '.'         |
|      |fully qualified domain name the client is requesting be resolved) |characters and ending with a '.' character (see qname_) |
+------+------------------------------------------------------------------+--------------------------------------------------------+
|type  |The qtype field from the client DNS request message (i.e.         |Examples are A (IpV4), AAAA (IpV6), NS (Name Service),  |
|      |the type of resolution that's requested such as IPv4, IPv6)       |  SOA (Start of Authority), and CNAME, (see qtype_)     |
+------+------------------------------------------------------------------+--------------------------------------------------------+
|class |The qclass field from the client DNS request message (i.e. The    |Either IN (Internet resource) or ANY (Traffic router    |
|      |class of resource being requested)                                |  rejects requests with any other value of class)       |
+------+------------------------------------------------------------------+--------------------------------------------------------+

**Response Fields**

+------+---------------------------------------------------------------------+-----------------------------------------------------+
|Name  | Description                                                         | Data                                                |
+======+=====================================================================+=====================================================+
|ttl   | The 'time to live' in seconds for the answer provided by Traffic    |A number from 0 to 4294967295                        |
|      | Router (clients can reliably use this answer for this long without  |                                                     |
|      | re-querying traffic router)                                         |                                                     |
+------+---------------------------------------------------------------------+-----------------------------------------------------+
|rcode | The result code for the DNS answer provided by Traffic Router       | One of NOERROR (success), NOTIMP (request is not    |
|      |                                                                     | NOTIMP (request is not  supported),                 |
|      |                                                                     | REFUSED (request is refused to be answered), or     |
|      |                                                                     | NXDOMAIN (the domain/name requested does not exist) |
+------+---------------------------------------------------------------------+-----------------------------------------------------+

.. _rl-tr-ngb:

GeoLimit Failure Redirect feature
=================================

Overview
--------
This feature is also called 'National GeoBlock' feature which is short for 'NGB' feature. In this section, the acronym 'NGB' will be used for this feature.

In the past, if the Geolimit check fails (for example, the client ip is not in the 'US' region but the geolimit is set to 'CZF + US'), the router will return 503 response; but with this feature, when the check fails, it will return 302 if the redirect url is set in the delivery service.

The Geolimit check failure has such scenarios:
1) When the GeoLimit is set to 'CZF + only', if the client ip is not in the the CZ file, the check fails
2) When the GeoLimit is set to any region, like 'CZF + US', if the client ip is not in such region, and the client ip is not in the CZ file, the check fails


Configuration
-------------
To enable the NGB feature, the DS must be configured with the proper redirect url. And the setting lays at 'Delivery Services'->Edit->'GeoLimit Redirect URL'. If no url is put in this field, the feature is disabled.

The URL has 3 kinds of formats, which have different meanings:

1. URL with no domain. If no domain is in the URL (like 'vod/dance.mp4'), the router will try to find a proper cache server within the delivery service and return the redirect url with the format like 'http://<cache server name>.<delivery service's FQDN>/<configured relative path>'

2. URL with domain that matches with the delivery service. For this URL, the router will also try to find a proper cache server within the delivery service and return the same format url as point 1.

3. URL with domain that doesn't match with the delivery service. For this URL, the router will return the configured url directly to the client.


.. _rl-deep-cache:

Deep Caching - Deep Coverage Zone Topology
==========================================

Overview
--------

Deep Caching is a feature that enables clients to be routed to the closest
possible "deep" edge caches on a per Delivery Service basis. The term "deep" is
used in the networking sense, meaning that the edge caches are located deep in
the network where the number of network hops to a client is as minimal as
possible. This deep caching topology is desirable because storing content closer
to the client gives better bandwidth savings, and sometimes the cost of
bandwidth usage in the network outweighs the cost of adding storage. While it
may not be feasible to cache an entire copy of the CDN's contents in every deep
location (for the best possible bandwidth savings), storing just a relatively
small amount of the CDN's most requested content can lead to very high bandwidth
savings.

Getting started
---------------

What you need:

#. Edge caches deployed in "deep" locations and registered in Traffic Ops
#. A Deep Coverage Zone File (DCZF) mapping these deep cache hostnames to specific network prefixes (see :ref:`rl-deep-czf` for details)
#. Deep caching parameters in the Traffic Router Profile (see :ref:`rl-ccr-profile` for details):

   * ``deepcoveragezone.polling.interval``
   * ``deepcoveragezone.polling.url``

#. Deep Caching enabled on one or more HTTP Delivery Services (i.e. ``deepCachingType`` = ALWAYS)

How it works
------------

Deep Coverage Zone routing is very similar to that of regular Coverage Zone
routing, except that the DCZF is preferred over the regular  CZF for Delivery
Services with DC (Deep Caching) enabled. If the client requests a DC-enabled
Delivery Service and their IP address gets a "hit" in the DCZF, Traffic Router
will attempt to route that client to one of the available deep caches in the
client's corresponding zone. If there are no deep caches available for a
client's request, Traffic Router will "fall back" to the regular CZF and
continue regular CZF routing from there.


.. _rl-tr-steering:

Steering feature
================

Overview
--------
A Steering delivery service is a delivery service that is used to "steer" traffic to other delivery services. A Steering delivery service will have target delivery services configured for it with weights assigned to them.  Traffic Router uses the weights to make a consistent hash ring which it then uses to make sure that requests are routed to a target based on the configured weights.  This consistent hash ring is separate from the consistent hash ring used in cache selection.

Special regular expressions called Filters can also be configured for target delivery services to pin traffic to a specific delivery service.  For example, if a filter called .*/news/.* for a target called target-ds-1 is created, any requests to traffic router with 'news' in them will be routed to target-ds-1.  This will happen regardless of the configured weights.

A client can bypass the steering functionality by providing a header called X-TC-Steering-Option with the xml_id of the target delivery service to route to.  When Traffic Router receives this header it will route to the requested target delivery service regardless of weight configuration.

Some other points of interest:

- Steering is currently only available for HTTP delivery services that are a part of the same CDN.
- A new role called STEERING has been added to the traffic ops database.  Only users with Admin or Steering privileges can modify steering assignments for a Delivery Service.
- A new API has been created in Traffic Ops under /internal.  This API is used by a Steering user to add filters and modify assignments.  (Filters can only be added via the API).
- Traffic Router uses the steering API in Traffic Ops to poll for steering assignments, the assignments are then used when routing traffic.

A couple simple use cases for steering are:

#. Migrating traffic from one delivery service to another over time.
#. Trying out new functionality for a subset of traffic with an experimental delivery service.
#. Load balancing between delivery services.



Configuration
-------------

The following needs to be completed for Steering to work correctly:

#. Two target delivery services are created in Traffic Ops.  They must both be HTTP delivery services part of the same CDN.
#. A delivery service with type STEERING is created in Traffic Ops.
#. Target delivery services are assigned to the steering delivery service using Traffic Ops.
#. A user with the role of Steering is created.
#. Using the API, the steering user assigns weights to the target delivery services.
#. If desired, the steering user can create filters for the target delivery services.

For more information see the `steering how-to guide <quick_howto/steering.html>`_.

HTTPS for Http Type Delivery Services
=====================================

Starting with version 1.7 Traffic Router added the ability to allow https traffic between itself and clients on a per http type delivery service basis.

.. Warning::
  The establishing of an HTTPS connection is much more computationally demanding than an HTTP connection.
  Since each client will in turn get redirected to ATS, Traffic Router is most always creating a new HTTPS connection for all HTTPS traffic.
  It is likely to mean that an existing Traffic Router will have some decrease in performance depending on the amount of https traffic you want to support
  As noted for DNSSEC, you may need to plan to scale Traffic Router vertically and/or horizontally to handle the new load

The summary for setting up https is to:

#. Select one of 'https', 'http and https', or 'http to https' for the delivery service 
#. Generate private keys for the delivery service using a wildcard domain such as ``*.my-delivery-service.my-cdn.example.com``
#. Obtain and import signed certificate chain
#. Snapshot CR Config

Clients may make HTTPS requests delivery services only after Traffic Router receives the certificate chain from Traffic Ops and the new CR Config.

Protocol Options
----------------

*https only*
  Traffic Router will only redirect (send a 302) to clients communicating with a secure connection, all other clients will receive a 503
*http and https*
  Traffic Router will redirect both secure and non-secure clients
*http to https*
  Traffic Router will redirect non-secure clients with a 302 and a location that is secure (i.e. starting with 'https' instead of 'http'), secure clients will remain on https
*http*
  Any secure client will get an SSL handshake error. Non-secure clients will experience the same behavior as prior to 1.7

Certificate Retrieval
---------------------

.. Warning::
  If you have https delivery services in your CDN, Traffic Router will not accept **any** connections until it is able to
  fetch certificates from Traffic Ops and load them into memory. Traffic Router does not persist certificates to the java keystore or anywhere else.

Traffic Router fetches certificates into memory:

* At startup time
* When it receives a new CR Config
* Once an hour from whenever the most recent of the last of the above occurred

.. Note::
  To adjust the frequency when Traffic Router fetches certificates add the parameter 'certificates.polling.interval' to CR Config and 
  setting it to the desired time in milliseconds.

.. Note::
  Taking a snapshot of CR Config may be used at times to avoid waiting the entire polling cycle for a new set of certificates.

.. Warning::
  If a snapshot of CR Config is made that involves a delivery service missing its certificates, Traffic Router will ignore **ALL** changes in that CR-Config
  until one of the following occurs:
  * It receives certificates for that delivery service 
  * Another snapshot of CR Config is created and the delivery service without certificates is changed so it's HTTP protocol is set to 'http'

Certificate Chain Ordering
--------------------------

The ordering of certificates within the certificate bundle matters. It must be:

#. Primary Certificate (e.g. the one created for ``*.my-delivery-service.my-cdn.example.com``)
#. Intermediate Certificate(s)
#. Root Certificate from CA (optional)

.. Warning::
  If something is wrong with the certificate chain (e.g. the order of the certificates is backwards or for the wrong domain) the
  client will get an SSL handshake.  Inspection of /opt/tomcat/logs/catalina.out is likely to yield information to reveal this.

To see the ordering of certificates you may have to manually split up your certificate chain and use openssl on each individual certificate

Suggested Way of Setting up an HTTPS Delivery Service
-----------------------------------------------------

Do the following in Traffic Ops:

#. Select one of 'https', 'http and https', or 'http to https' for the protocol field of a delivery service and click 'Save'.
#. Click 'Manage SSL Keys'.
#. Click 'Generate New Keys'.
#. Copy the contents of the Certificate Signing Request field and save it locally.
#. Click 'Load Keys'.
#. Select 'http' for the protocol field of the delivery service and click 'Save' (to avoid preventing other CR Config updates from being blocked by Traffic Router)
#. Follow your standard procedure for obtaining your signed certificate chain from a CA.
#. After receiving your certificate chain import it into Traffic Ops.
#. Edit the delivery service.
#. Restore your original choice for the protocol field and click save.
#. Click 'Manage SSL Keys'.
#. Click 'Paste Existing Keys'.
#. Paste the certificate chain into the CRT field.
#. Click 'Load Keys'.
#. Take a new snapshot of CR Config.

Once this is done you should be able to test you are getting correctly redirected by Traffic Router using curl commands to https destinations on your delivery service.

A new testing tool was created for load testing traffic router, it allows you to generate requests from your local box to multiple delivery services of a single cdn.
You can control which cdn, delivery services, how many transactions per delivery service, and how many concurrent requests.
During the test it will provide feedback about request latency and transactions per second.

While it is running it is suggested that you monitor your Traffic Router nodes for memory and CPU utilization.

Tuning Recommendations
======================

The following is an example of /opt/tomcat/bin/setenv.sh that has been tested on a multi core server running under HTTPS load test requests.
This is following the general recommendation to use the G1 garbage collector for JVM applications running on multi core machines.
In addition to using the G1 garbage collector the InitiatingHeapOccupancyPercent was lowered to run garbage collection more frequently which
improved overall throughput for Traffic Router and reduced 'Stop the World' garbage collection. Note that setting the min and max heap settings
in setenv.sh will override init scripts in /etc/init.d/tomcat.

  /opt/tomcat/bin/setenv.sh::


      #! /bin/sh
      export CATALINA_OPTS="$CATALINA_OPTS -server"
      export CATALINA_OPTS="$CATALINA_OPTS -Xms2g -Xmx2g"
      export CATALINA_OPTS="$CATALINA_OPTS -XX:+UseG1GC"
      export CATALINA_OPTS="$CATALINA_OPTS -XX:+UnlockExperimentalVMOptions"
      export CATALINA_OPTS="$CATALINA_OPTS -XX:InitiatingHeapOccupancyPercent=30"
