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

*****************************
Traffic Router Administration
*****************************
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

..  Note:: The ``CDN_name`` parameter with a config file name of ``rascal-config.txt`` must exist in the profile and the value should be the name of the CDN for which this Traffic Router will be authoritative. This same parameter will be mapped to all profiles that participate in this CDN (edges, mids, Traffic Monitors, etc). See :ref:`rl-param-prof` for more information.

2. Enter the Traffic Router server into Traffic Ops, assign it to a Traffic Router profile, and ensure that its status is set to ``ONLINE``.
3. Ensure the FQDN of the Traffic Monitor is resolvable in DNS. This FQDN must be resolvable by the clients expected to use this CDN.
4. Install a traffic router: ``sudo yum install traffic_router``.
5. Edit ``/opt/traffic_router/conf/traffic_monitor.properties`` and specify the correct online Traffic Monitor(s) for your CDN. See :ref:`rl-tr-config-files`
	# traffic_monitor.properties: url that should normally point to this file
	traffic_monitor.properties=file:/opt/traffic_router/conf/traffic_monitor.properties

	# Frequency for reloading this file
	# traffic_monitor.properties.reload.period=60000
   

6. Start Tomcat: ``sudo service tomcat start``, and test lookups with dig and curl against that server.
7. Snapshot CRConfig; See :ref:`rl-snapshot-crconfig`

..  Note:: Once the CRConfig is snapshotted, live traffic will be sent to the new Traffic Routers provided that their status is set to ``ONLINE``.

8. Ensure that the parent domain (e.g.: kabletown.net) for the CDN's top level domain (e.g.: cdn.kabletown.net) contains a delegation (NS records) for the new Traffic Router, and that the value specified matches the FQDN used in step 3.

Configuring Traffic Router
==========================

By default, Traffic Router installs all configuration files under ``/opt/traffic_router/conf``. For the most part, the configuration files and parameters that follow are used to get Traffic Router online and communicating with various Traffic Control components. Once Traffic Router is successfully communicating with Traffic Control, configuration is mostly performed in Traffic Ops, and is distributed throughout Traffic Control via the CRConfig snapshot process. See :ref:`rl-snapshot-crconfig` for more information. Please see the parameter documentation for Traffic Router in the Using Traffic Ops guide documented under :ref:`rl-ccr-profile` for parameters that influence the behavior of Traffic Router via the CRConfig.

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
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | dns.routing.name                          | The label (A/AAAA) Traffic Router will use for the entry point for a DNS delivery service           | edge (e.g.: edge.mydeliveryservice.kabletown.net) |
+----------------------------+-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
| traffic_ops.properties     | traffic_ops.username                      | Username to access the APIs in Traffic Ops (must be in the admin role)                              | admin                                             |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | traffic_ops.password                      | Password for the user specified in traffic_ops.username                                             | N/A                                               |
+----------------------------+-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
| http.properties            | http.routing.name                         | The label (A/AAAA) Traffic Router will use for the entry point for an HTTP delivery service         | tr (e.g.: tr.mydeliveryservice.kabletown.net)     |
+----------------------------+-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
| cache.properties           | cache.geolocation.database                | Full path to the local copy of the MaxMind geolocation binary database file                         | /opt/traffic_router/db/GeoIP2-City.mmdb           |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | cache.geolocation.database.refresh.period | The interval in milliseconds which Traffic Router will poll for a new geolocation database          | 604800000                                         |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | cache.czmap.database                      | Full path to the local copy of the coverage zone file                                               | /opt/traffic_router/db/czmap.json                 |
|                            +-------------------------------------------+-----------------------------------------------------------------------------------------------------+---------------------------------------------------+
|                            | cache.czmap.database.refresh.period       | The interval in milliseconds which Traffic Router will poll for a new coverage zone file            | 10800000                                          |
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

+------+---------------------------------------------------------------------------------+---------------------------------------------------------------------------+
|Name  |Description                                                                      |Data                                                                       |
+======+=================================================================================+===========================================================================+
|qtype |Whether the request was for DNS or HTTP                                          |Always DNS or HTTP                                                         |
+------+---------------------------------------------------------------------------------+---------------------------------------------------------------------------+
|chi   |The IP address of the requester                                                  |Depends on whether this was a DNS or HTTP request, see below sections      |
+------+---------------------------------------------------------------------------------+---------------------------------------------------------------------------+
|ttms  |The amount of time in milliseconds it took Traffic Router to process the request |A number greater than or equal to zero                                     |
+------+---------------------------------------------------------------------------------+---------------------------------------------------------------------------+
|rtype |Routing Result Type                                                              |One of ERROR, CZ, GEO, MISS, STATIC_ROUTE, DS_REDIRECT, DS_MISS, INIT, FED |
+------+---------------------------------------------------------------------------------+---------------------------------------------------------------------------+
|rloc  |GeoLocation of result                                                            |Latitude and Longitude in Decimal Degrees                                  |
+------+---------------------------------------------------------------------------------+---------------------------------------------------------------------------+
|rdtl  |Result Details Associated with unusual conditions                                |One of DS_NOT_FOUND, DS_NO_BYPASS, DS_BYPASS, DS_CZ_ONLY                   |
+------+---------------------------------------------------------------------------------+---------------------------------------------------------------------------+
|rerr  |Message about internal Traffic Router Error                                      |String                                                                     |
+------+---------------------------------------------------------------------------------+---------------------------------------------------------------------------+

**rtype meanings**

+-------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|Name         |Meaning                                                                                                                                                                 |
+=============+========================================================================================================================================================================+
|ERROR        |An internal error occurred within Traffic Router, more details may be found in the rerr field                                                                           |
+-------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|CZ           |The result was derived from Coverage Zone data based on the address in the chi field                                                                                    |
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

**rdtl meanings**

+--------------------------+--------------------------------------------------------------------------------------------------------------------------------------------+
|Name                      |Meaning                                                                                                                                     |
+==========================+============================================================================================================================================+
|DS_NOT_FOUND              |Always goes with rtypes STATIC_ROUTE and DS_MISS                                                                                            |
+--------------------------+--------------------------------------------------------------------------------------------------------------------------------------------+
|DS_BYPASS                 |Used Bypass Destination for Redirect of Delivery Service                                                                                    |
+--------------------------+--------------------------------------------------------------------------------------------------------------------------------------------+
|DS_NO_BYPASS              |No valid Bypass Destination is configured for the matched Delivery Service and the delivery service does not have the requested resource    |
+--------------------------+--------------------------------------------------------------------------------------------------------------------------------------------+
|DS_CZ_ONLY                |The selected Delivery Service only supports resource lookup based on Coverage Zone data                                                     |
+--------------------------+--------------------------------------------------------------------------------------------------------------------------------------------+
|DS_CLIENT_GEO_UNSUPPORTED |Traffic Router did not find a resource supported by coverage zone data and was unable to determine the geolocation of the requesting client |
+--------------------------+--------------------------------------------------------------------------------------------------------------------------------------------+
|GEO_NO_CACHE_FOUND        |Traffic Router could not find a resource via geolocation data based on the requesting client's geolocation                                  |
+--------------------------+--------------------------------------------------------------------------------------------------------------------------------------------+

---------------

HTTP Specifics
--------------

Sample Message
::
  1452197640.936 qtype=HTTP chi=69.241.53.218 url="http://ccr.mm-test.jenkins.cdnlab.comcast.net/some/asset.m3u8" cqhm=GET cqhv=HTTP/1.1 rtype=GEO rloc="40.252611,58.439389" rdtl=- rerr="-" pssc=302 ttms=0 rurl="http://odol-atsec-sim-114.mm-test.jenkins.cdnlab.comcast.net:8090/some/asset.m3u8" rh="Accept: */*" rh="myheader: asdasdasdasfasg"

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

