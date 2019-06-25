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

.. _tr-admin:

*****************************
Traffic Router Administration
*****************************

Requirements
============
* CentOS 7
* 4 CPUs
* 8GB of RAM
* Successful install of Traffic Ops (usually on another machine)
* Successful install of Traffic Monitor (usually on another machine)
* Administrative access to Traffic Ops

.. Note:: Hardware requirements are generally doubled if :ref:`tr-DNSSEC` is enabled

Installing Traffic Router
=========================
#. If no suitable :term:`Profile` exists, create a new :term:`Profile` for Traffic Router via the :guilabel:`+` button on the :ref:`tp-profiles-page` page in Traffic Portal

	.. warning:: Traffic Ops will *only* recognize a profile as assignable to a Traffic Router if its name starts with the prefix ``ccr-``. The reason for this is a legacy limitation related to the old name for Traffic Router (Comcast Cloud Router), and will (hopefully) be rectified in the future as the old Perl parts of Traffic Ops are re-written in Go.

#. Enter the Traffic Router server into Traffic Portal on the :ref:`tp-configure-servers` page (or via the :ref:`to-api`), assign to it a Traffic Router :term:`Profile`, and ensure that its status is set to ``ONLINE``.
#. Ensure the :abbr:`FQDN (Fully Qualified Domain Name)` of the Traffic Router is resolvable in DNS. This :abbr:`FQDN (Fully Qualified Domain Name)` must be resolvable by the clients expected to use this CDN.
#. Install a Traffic Router server package, either from source or using a :file:`traffic_router-{version string}.rpm` package generated using the instructions in :ref:`dev-building`.

	.. versionchanged:: 3.0
		As of version 3.0, Traffic Router depends upon a package called ``tomcat``. This package should have been created when Traffic Router was built. If installing the ``traffic_router`` produces a depenedency error, make sure that the ``tomcat`` package is available in an accessible :manpage:`yum(8)` repository.

#. Edit :file:`/opt/traffic_router/conf/traffic_monitor.properties` and specify the correct online Traffic Monitor(s) for your CDN.

	.. seealso:: :ref:`tr-config-files`

	:file:`traffic_monitor.properties`
		URL that should normally point to this file, e.g. ``traffic_monitor.properties=file:/opt/traffic_router/conf/traffic_monitor.properties``
	:file:`traffic_monitor.properties.reload.period`
		Period to wait (in milliseconds) between reloading this file, e.g. ``traffic_monitor.properties.reload.period=60000``

#. Start Traffic Router. This is normally done by starting its :manpage:`systemd(1)` service. ``systemctl start traffic_router`` , and test DNS lookups against that server to be sure it's resolving properly. with e.g. ``dig`` or ``curl``. Also, because previously taken CDN :term:`Snapshot`\ s will be cached, they need to be removed manually to actually be reloaded. This file should be located at :file:`/opt/traffic_router/db/cr-config.json`. This should be done before starting or restarting Traffic Router.

	.. code-block:: console
		:caption: Starting and Testing Traffic Router

		[root@trafficrouter /]# systemctl start traffic_router
		[root@trafficrouter /]# dig @localhost mycdn.ciab.test

		; <<>> DiG 9.9.4-RedHat-9.9.4-72.el7 <<>> @localhost mycdn.ciab.test
		; (2 servers found)
		;; global options: +cmd
		;; Got answer:
		;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 27109
		;; flags: qr aa rd; QUERY: 1, ANSWER: 0, AUTHORITY: 1, ADDITIONAL: 0
		;; WARNING: recursion requested but not available

		;; QUESTION SECTION:
		;mycdn.ciab.test.		IN	A

		;; AUTHORITY SECTION:
		mycdn.ciab.test.	30	IN	SOA	trafficrouter.infra.ciab.test. twelve_monkeys.mycdn.ciab.test. 2019010918 28800 7200 604800 30

		;; Query time: 28 msec
		;; SERVER: ::1#53(::1)
		;; WHEN: Wed Jan 09 21:27:57 UTC 2019
		;; MSG SIZE  rcvd: 104

#. Perform a CDN :term:`Snapshot`.

	.. Note:: Once the :term:`Snapshot` is taken, live traffic will be sent to the new Traffic Routers provided that their status has been set to ``ONLINE``.

#. Ensure that the parent domain (e.g.: ``cdn.local``) for the CDN's top level domain (e.g.: ``ciab.cdn.local``) contains a delegation (Name Server records) for the new Traffic Router, and that the value specified matches the :abbr:`FQDN (Fully Qualified Domain Name)` of the Traffic Router.

Configuring Traffic Router
==========================
.. versionchanged:: 1.5
	Many of the configuration files under :file:`/opt/traffic_router/conf` are now only needed to override the default configuration values for Traffic Router. Most of the given default values will work well for any CDN. Critical values that must be changed are hostnames and credentials for communicating with other Traffic Control components such as Traffic Ops and Traffic Monitor. Pre-existing installations that store configuration files under ``/opt/traffic_router/conf`` will still be used and honored for Traffic Router 1.5 onward.

.. versionchanged:: 3.0
	Traffic Router 3.0 has been converted to a formal Tomcat instance, meaning that is now installed separately from the Tomcat servlet engine. The Traffic Router installation package contains all of the Traffic Router-specific software, configuration and startup scripts including some additional configuration files needed for Tomcat. These new configuration files can all be found in the :file:`/opt/traffic_router/conf` directory and generally serve to override Tomcat's default settings.

For the most part, the configuration files and :term:`Parameters` used by Traffic Router are used to bring it online and start communicating with various Traffic Control components. Once Traffic Router is successfully communicating with Traffic Control, configuration should mostly be performed in Traffic Portal, and will be distributed throughout Traffic Control via CDN :term:`Snapshot` process. Please see the :term:`Parameter` documentation for Traffic Router in the Using Traffic Ops guide documented under :ref:`ccr-profile` for :term:`Parameters` that influence the behavior of Traffic Router via the :term:`Snapshot`.

.. _tr-config-files:
.. table:: Traffic Router Parameters

	+----------------------------+-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	| ConfigFile                 | Parameter Name                            | Description                                                                      | Default Value                                      |
	+============================+===========================================+==================================================================================+====================================================+
	| traffic_monitor.properties | traffic_monitor.bootstrap.hosts           | Semicolon-delimited Traffic Monitor                                              | N/A                                                |
	|                            |                                           | :abbr:`FQDN (Fully Qualified Domain Name)`\ s with port numbers as necessary     |                                                    |
	|                            +-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	|                            | traffic_monitor.bootstrap.local           | Use only the Traffic Monitors specified in local configuration files             | ``false``                                          |
	|                            +-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	|                            | traffic_monitor.properties                | Path to file:`traffic_monitor.properties`; used internally to monitor the file   | ``/opt/traffic_router/traffic_monitor.properties`` |
	|                            |                                           | for changes                                                                      |                                                    |
	|                            +-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	|                            | traffic_monitor.properties.reload.period  | The interval in milliseconds for Traffic Router to wait between reloading this   | ``60000``                                          |
	|                            |                                           | configuration file                                                               |                                                    |
	+----------------------------+-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	| dns.properties             | dns.tcp.port                              | TCP port that Traffic Router will use for incoming DNS requests                  | ``53``                                             |
	|                            +-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	|                            | dns.tcp.backlog                           | Maximum length of the queue for incoming TCP connection requests                 | ``0``                                              |
	|                            +-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	|                            | dns.udp.port                              | UDP port that Traffic Router will use for incoming DNS requests                  | ``53``                                             |
	|                            +-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	|                            | dns.max-threads                           | Maximum number of threads used to process incoming DNS requests                  | ``1000``                                           |
	|                            +-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	|                            | dns.zones.dir                             | Path to automatically generated zone files for reference                         | ``/opt/traffic_router/var/auto-zones``             |
	+----------------------------+-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	| traffic_ops.properties     | traffic_ops.username                      | Username with which to access the :ref:`to-api`                                  | ``admin``                                          |
	|                            |                                           | (must have the ``admin`` :term:`Role`)                                           |                                                    |
	|                            +-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	|                            | traffic_ops.password                      | Password for the user specified in ``traffic_ops.username``                      | N/A                                                |
	+----------------------------+-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	| cache.properties           | cache.geolocation.database                | Full path to the local copy of a geographic IP mapping database                  | ``/opt/traffic_router/db/GeoIP2-City.mmdb``        |
	|                            |                                           | (usually MaxMind's GeoIP2)                                                       |                                                    |
	|                            +-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	|                            | cache.geolocation.database.refresh.period | The interval in milliseconds for Traffic Router to wait between polling for      | ``604800000``                                      |
	|                            |                                           | changes to the GeoIP2 database                                                   |                                                    |
	|                            +-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	|                            | cache.czmap.database                      | Full path to the local copy of the coverage zone file                            | ``/opt/traffic_router/db/czmap.json``              |
	|                            +-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	|                            | cache.czmap.database.refresh.period       | The interval in milliseconds for Traffic Router to wait between polling for a    | ``10800000``                                       |
	|                            |                                           | new coverage zone file                                                           |                                                    |
	|                            +-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	|                            | cache.dczmap.database                     | Full path to the local copy of the deep coverage zone file                       | ``/opt/traffic_router/db/dczmap.json``             |
	|                            +-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	|                            | cache.dczmap.database.refresh.period      | The interval in milliseconds for Traffic Router to wait between polling for a    | ``10800000``                                       |
	|                            |                                           | new deep coverage zone file                                                      |                                                    |
	|                            +-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	|                            | cache.health.json                         | Full path to the local copy of the health state                                  | ``/opt/traffic_router/db/health.json``             |
	|                            +-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	|                            | cache.health.json.refresh.period          | The interval in milliseconds which Traffic Router will poll for a new health     | ``1000``                                           |
	|                            |                                           | state file                                                                       |                                                    |
	|                            +-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	|                            | cache.config.json                         | Full path to the locally cached copy of the CDN :term:`Snapshot`                 | ``/opt/traffic_router/db/cr-config.json``          |
	|                            +-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	|                            | cache.config.json.refresh.period          | The interval in milliseconds which Traffic Router will poll for a new            | ``60000``                                          |
	|                            |                                           | :term:`Snapshot`                                                                 |                                                    |
	+----------------------------+-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	| startup.properties         | various parameters                        | This configuration is used by :manpage:`systemd(1)` to set environment variables | N/A                                                |
	|                            |                                           | when the ``traffic_router`` service is started. It primarily consists of command |                                                    |
	|                            |                                           | line settings for the Java process                                               |                                                    |
	+----------------------------+-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	| log4j.properties           | various parameters                        | Configuration of ``log4j`` is documented on                                      | N/A                                                |
	|                            |                                           | `their site <http://logging.apache.org/log4j/2.x/index.html>`_; adjust as needed |                                                    |
	+----------------------------+-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	| server.xml                 | various parameters                        | Traffic Router specific configuration for Apache Tomcat. See the Apache Tomcat   | N/A                                                |
	|                            |                                           | `documentation <https://tomcat.apache.org/tomcat-8.5-doc/index.html>`_           |                                                    |
	+----------------------------+-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+
	| web.xml                    | various parameters                        | Default settings for all Web Applications running in the Traffic Router instance | N/A                                                |
	|                            |                                           | of Tomcat                                                                        |                                                    |
	+----------------------------+-------------------------------------------+----------------------------------------------------------------------------------+----------------------------------------------------+


.. _consistent-hashing:

Consistent Hashing
==================
Traffic Router does special optimization for some requests to ensure that requests for specific content are consistently fetched from a small number (often exactly one, but dependent on :ref:`ds-initial-dispersion`) of :term:`cache servers` - thus ensuring it stays "fresh" in the cache. This is done by performing "consistent hashing" on request paths (when HTTP routing) or names requested for resolution (when DNS routing). To an extent, this behavior is configurable by modifying fields on :term:`Delivery Services`. Consistent hashing acts differently on a :term:`Delivery Service` based on how :term:`Delivery Services` of its :ref:`ds-types` route content.

- HTTP, HTTP_NO_CACHE, HTTP_LIVE, HTTP_LIVE_NATNL, DNS, DNS_LIVE, and DNS_NATNL
	These :ref:`Delivery Service Types <ds-types>` route directly to :term:`cache servers`, so consistent hashing is used to choose a :term:`cache server` to which the client will be redirected.

- STEERING and CLIENT_STEERING
	These :ref:`Delivery Service Types <ds-types>` route to "target" :term:`Delivery Services`, so consistent hashing is used to choose a "target" which will service the client request.

.. seealso:: See `the Wikipedia article on consistent hashing <http://en.wikipedia.org/wiki/Consistent_hashing>`_.

.. _pattern-based-consistenthash:

Consistent Hashing Patterns
---------------------------
.. versionadded:: 4.0

Regular expressions ("patterns") can be provided in the :ref:`ds-consistent-hashing-regex` field of an HTTP-:ref:`routed <ds-types>` Delivery Service to influence what parts of an HTTP request path are considered when performing consistent hashing. These patterns propagate to Traffic Router through :term:`Snapshots`.

.. important:: Consistent Hashing Patterns on STEERING-:ref:`ds-types` :term:`Delivery Services` will be used for Consistent Hashing - the Consistent Hashing Pattern(s) of said :term:`Delivery Service`'s target(s) will **not** be considered. If Consistent Hashing Patterns are important to the routing of content on a STEERING-:ref:`ds-types` or CLIENT_STEERING-:ref:`ds-types` :term:`Delivery Service`, they **must** be defined *on that* :term:`Delivery Service` *itself, and* **not** *on its target(s)*.

How it Works
""""""""""""
The supplied :ref:`ds-consistent-hashing-regex` is applied to the request path to extract matching elements to build a new string *before* consistent hashing is done. For example, using the pattern :regexp:`/.*?(/.*?/).*?(m3u8)` and given the request paths ``/test/path/asset.m3u8`` and ``/other/path/asset.m3u8`` the resulting string used for consistent hashing will be ``/path/m3u8``

.. seealso:: See Oracle's `documentation for the java.util.regex.Pattern <https://docs.oracle.com/javase/7/docs/api/java/util/regex/Pattern.html>`_ implementation in Java.

Testing Pattern-Based Consistent Hashing
""""""""""""""""""""""""""""""""""""""""
In order to test this feature without affecting the delivery of traffic through a CDN, there are several test tools in place.

- :ref:`tr-api`
	Several Traffic Router endpoints exist to test regular expression application against a request path, :term:`cache server` selection, and :term:`Delivery Service` selection.
- :ref:`to-api`
	The :ref:`to-api-consistenthash` endpoint will proxy request data through to one of the Traffic Router endpoints in order to test regular expression application against a request path, in the event that direct access to the :ref:`tr-api` is not possible and/or desired.
- Traffic Portal
	On the :term:`Delivery Service` creation/modification form in Traffic Portal (under :ref:`tp-services-delivery-service`), there is a :guilabel:`Test Regex` section that the user can use to validate a regular expression before saving it to a :term:`Delivery Service`.

Consistent Hash Query Parameters
--------------------------------
Normally, when performing consistent hashing for an HTTP-:ref:`routed <ds-types>` :term:`Delivery Service`, any query parameters present in the request are ignored. That is, if a client requests ``/some/path?key=value`` consistent hashing is only performed on the string '``/some/path``'. However, query parameters that are part of uniquely identifying content can be specified by adding them to the set of :ref:`ds-consistent-hashing-qparams` of a :term:`Delivery Service`. For example, suppose that the file ``/video.mp4`` is available on the :term:`origin server` in different resolutions, which are specified by the ``resolution`` query parameter. This means that ``/video.mp4?resolution=480p`` and ``/video.mp4?resolution=720p`` share a *request path*, but represent different *content*. In that case, adding ``resolution`` to the :term:`Delivery Service`'s :ref:`ds-consistent-hashing-qparams` will cause consistent hashing to be done on e.g. ``/video.mp4?resolution=480p`` instead of just ``/video.mp4`` - however if the client requests e.g. ``/video.mp4?resolution=480p&bitrate=120kbps`` consistent hashing will *only* consider ``/video.mp4?resolution=480p``.

.. note:: `Consistent Hashing Patterns`_ are applied *before* query parameters are considered - i.e. a pattern cannot match against query parameters, and need not worry about query parameters contaminating matches.

.. important:: Consistent Hash Query Parameters on the *targets* of STEERING-:ref:`ds-types` :term:`Delivery Services` will be used for Consistent Hashing - the Consistent Hash Query Parameters of said :term:`Delivery Services` themselves will **not** be considered. If Consistent Hash Query Parameters are important to the routing of content on a STEERING-:ref:`ds-types` or CLIENT_STEERING-:ref:`ds-types` :term:`Delivery Service`, they **must** be defined *on that* :term:`Delivery Service`'s' *target(s), and* **not** *on the* :term:`Delivery Service` *itself*.

.. caution:: Certain query parameters are reserved by Traffic Router for its own use, and thus cannot be present in any Consistent Hash Query Parameters. These reserved parameters are:

	 - trred
	 - format
	 - fakeClientIPAddress

.. _tr-dnssec:

DNSSEC
======
.. seealso:: `The Wikipedia page on Domain Name Security Extensions <https://en.wikipedia.org/wiki/Domain_Name_System_Security_Extensions>`_

Overview
--------
:abbr:`DNSSEC (Domain Name System Security Extensions)` is a set of extensions to DNS that provides a cryptographic mechanism for resolvers to verify the authenticity of responses served by an authoritative DNS server. Several RFCs (:rfc:`4033`, :rfc:`4044`, :rfc:`4045`) describe the low level details and define the extensions, :rfc:`7129` provides clarification around authenticated denial of existence of records, and finally :rfc:`6781` describes operational best practices for administering an authoritative :abbr:`DNSSEC (Domain Name System Security Extensions)`-enabled DNS server. The authenticated denial of existence :rfc:`7129` describes how an authoritative DNS server responds in NXDOMAIN and NODATA scenarios when :abbr:`DNSSEC (Domain Name System Security Extensions)` is enabled. Traffic Router currently supports :abbr:`DNSSEC (Domain Name System Security Extensions)` with :abbr:`NSEC (Next Secure Record)`, however, :abbr:`NSEC3 (Next Secure Record version 3)` and more configurable options are planned for the future.

Operation
---------
Upon startup or a configuration change, Traffic Router obtains keys from the 'keystore' API in Traffic Ops which returns :abbr:`KSK (Key Signing Key)`\ s and :abbr:`ZSK (Zone Signing Key)`\ s for each :term:`Delivery Service` that is a sub-domain of the CDN's :abbr:`TLD (Top Level Domain)` in addition to the keys for the CDN :abbr:`TLD (Top Level Domain)` itself. Each key has timing information that allows Traffic Router to determine key validity (expiration, inception, and effective dates) in addition to the appropriate :abbr:`TTL (Time To Live)` to use for the DNSKEY record(s). All :abbr:`TTL (Time To Live)`\ s are configurable :term:`Parameter`\ s; see the :ref:`ccr-profile` documentation for more information.

Once Traffic Router obtains the key data from the API, it converts each public key into the appropriate record types (DNSKEY, DS) to place in zones and uses the private key to sign zones. DNSKEY records are added to each :term:`Delivery Service`'s zone (e.g.: ``demo1.mycdn.ciab.test``) for every valid key that exists, in addition to the CDN :abbr:`TLD (Top Level Domain)`'s zone. A DS record is generated from each zone's :abbr:`KSK (Key Signing Key)` and is placed in the CDN :abbr:`TLD (Top Level Domain)`'s zone (e.g.: ``mycdn.ciab.test``); the DS record for the CDN :abbr:`TLD (Top Level Domain)` must be placed in its parent zone, which is not managed by Traffic Control.

The DNSKEY to DS record relationship allows resolvers to validate signatures across zone delegation points. With Traffic Control, we control all delegation points below the CDN's :abbr:`TLD (Top Level Domain)`, **however, the DS record for the CDN** :abbr:`TLD (Top Level Domain)` **must be placed in the parent zone** (e.g.: ``ciab.test``), **which is not managed by Traffic Control**. As such, the DS record must be placed in the parent zone prior to enabling :abbr:`DNSSEC (Domain Name System Security Extensions)`, and prior to generating a new CDN KSK. Based on your deployment's DNS configuration, this might be a manual process or it might be automated. Either way, extreme care and diligence must be taken and knowledge of the management of the upstream zone is imperative for a successful :abbr:`DNSSEC (Domain Name System Security Extensions)` deployment.

To enable :abbr:`DNSSEC (Domain Name System Security Extensions)` for a CDN in Traffic Portal, Go to :guilabel:`CDNs` from the sidebar and click on the desired CDN, then toggle the 'DNSSEC Enabled' field to 'true', and click on the green :guilabel:`Update` button to save the changes.

Rolling Zone Signing Keys
-------------------------
Traffic Router currently follows the :abbr:`ZSK (Zone Signing Key)` pre-publishing operational best practice described in :rfc:`6781#section-4.1.1.1`. Once :abbr:`DNSSEC (Domain Name System Security Extensions)` is enabled for a CDN in Traffic Portal, key rolls are triggered by Traffic Ops via the automated key generation process, and Traffic Router selects the active :abbr:`ZSK (Zone Signing Keys)`\ s based on the expiration information returned from the 'keystore' API of Traffic Ops.

.. _tr-logs:

Troubleshooting and Log Files
=============================
Traffic Router log files can be found under :file:`/opt/traffic_router/var/log` and :file:`/opt/tomcat/logs`. Initialization and shutdown logs are in :file:`/opt/tomcat/logs/catalina{date}.out`. Application related logging is in :file:`/opt/traffic_router/var/log/traffic_router.log`, while access logs are written to :file:`/opt/traffic_router/var/log/access.log`.

Event Log File Format
---------------------

Summary
"""""""
All access events to Traffic Router are logged to the file :file:`/opt/traffic_router/var/log/access.log`. This file grows up to 200MB and gets rolled into older log files, ten log files total are kept (total of up to 2GB of logged events per Traffic Router instance)

Traffic Router logs access events in a format that largely follows :abbr:`ATS (Apache Traffic Service)` `event logging format <https://docs.trafficserver.apache.org/en/6.0.x/admin/event-logging-formats.en.html>`_.

Message Format
""""""""""""""
- Except for the first item, each event that is logged is a series of space-separated key/value pairs.
- The first item is always the Unix epoch in seconds with a decimal field precision of up to milliseconds.
- Each key/value pair is in the form of ``unquoted_string="optionally quoted string"``
- Values that are quoted strings may contain whitespace characters.
- Values that are not quoted should not contains any whitespace characters.

.. Note:: Any value that is a single dash character or a dash character enclosed in quotes represents an empty value

.. table:: Fields Always Present

	+-------+----------------------------------------------------------------------------------+-------------------------------------------------------------------------------------+
	| Name  | Description                                                                      | Data                                                                                |
	+=======+==================================================================================+=====================================================================================+
	| qtype | Whether the request was for DNS or HTTP                                          | Always "DNS" or "HTTP"                                                              |
	+-------+----------------------------------------------------------------------------------+-------------------------------------------------------------------------------------+
	| chi   | The IP address of the requester                                                  | Depends on whether this was a DNS or HTTP request, see other sections               |
	+-------+----------------------------------------------------------------------------------+-------------------------------------------------------------------------------------+
	| rhi   | The IP address of the request source address                                     | Depends on whether this was a DNS or HTTP request, see other sections               |
	+-------+----------------------------------------------------------------------------------+-------------------------------------------------------------------------------------+
	| ttms  | The amount of time in milliseconds it took Traffic Router to process the request | A number greater than or equal to zero                                              |
	+-------+----------------------------------------------------------------------------------+-------------------------------------------------------------------------------------+
	| rtype | Routing result type                                                              | One of ERROR, CZ, DEEP_CZ, GEO, MISS, STATIC_ROUTE, DS_REDIRECT, DS_MISS, INIT, FED |
	+-------+----------------------------------------------------------------------------------+-------------------------------------------------------------------------------------+
	| rloc  | GeoLocation of result                                                            | Latitude and longitude in degrees as floating point numbers                         |
	+-------+----------------------------------------------------------------------------------+-------------------------------------------------------------------------------------+
	| rdtl  | Result details Associated with unusual conditions                                | One of DS_NOT_FOUND, DS_NO_BYPASS, DS_BYPASS, DS_CZ_ONLY, DS_CZ_BACKUP_CG           |
	+-------+----------------------------------------------------------------------------------+-------------------------------------------------------------------------------------+
	| rerr  | Message about an internal Traffic Router error                                   | String                                                                              |
	+-------+----------------------------------------------------------------------------------+-------------------------------------------------------------------------------------+

.. seealso:: If `Regional Geo-Blocking <regionalgeo-qht>`_ is enabled on the :term:`Delivery Service`, an additional field (``rgb``) will appear.

Sample Message
""""""""""""""
Items within brackets are detailed under the HTTP and DNS sections

.. code-block:: text
	:caption: Example Logfile Lines

	144140678.000 qtype=DNS chi=192.168.10.11 rhi=- ttms=789 [Fields Specific to the DNS request] rtype=CZ rloc="40.252611,58.439389" rdtl=- rerr="-" [Fields Specific to the DNS result]
	144140678.000 qtype=HTTP chi=192.168.10.11 rhi=- ttms=789 [Fields Specific to the HTTP request] rtype=GEO rloc="40.252611,58.439389" rdtl=- rerr="-" [Fields Specific to the HTTP result]

.. note:: These samples contain fields that are always present for every single access event to Traffic Router


``rtype`` Meanings
""""""""""""""""""
``-``
	The request was not redirected. This is usually a result of a DNS request to the Traffic Router or an explicit denial for that request
ANON_BLOCK
	The client's IP matched an `Anonymous Blocking <anonymous_blocking-qht>`_ rule and was blocked
CZ
	The result was derived from Coverage Zone data based on the address in the ``chi`` field
DEEP_CZ
	The result was derived from Deep Coverage Zone data based on the address in the ``chi`` field
DS_MISS
	_*HTTP Only*_ No HTTP :term:`Delivery Service`\ supports either this request's URL path or headers
DS_REDIRECT
	The result is using the Bypass Destination configured for the matched :term:`Delivery Service` when that :term:`Delivery Service` is unavailable or does not have the requested resource
ERROR
	An internal error occurred within Traffic Router, more details may be found in the ``rerr`` field
FED
	_*DNS Only*_ The result was obtained through federated coverage zone data outside of any :term:`Delivery Service`\ s
GEO
	The result was derived from geolocation service based on the address in the ``chi`` field
GEO_REDIRECT
	The request was redirected based on the National Geo blocking (Geo Limit Redirect URL) configured on the :term:`Delivery Service`
MISS
	Traffic Router was unable to resolve a DNS request or find a cache for the requested resource
RGALT
	The request was redirected to the `Regional Geo-Blocking <regionalgeo-qht>`_ URL. Regional Geo blocking is enabled on the :term:`Delivery Service` and is configured through the ``regional_geoblock.polling.url`` :term:`Parameter` on the Traffic Router :term:`Profile`
RGDENY
	_*DNS Only*_ The result was obtained through federated coverage zone data outside of any :term:`Delivery Service` - the request was regionally blocked because there was no rule for the request made
STATIC_ROUTE
	_*DNS Only*_ No DNS :term:`Delivery Service`\ supports the hostname portion of the requested URL


``rdtl`` Meanings
"""""""""""""""""
``-``
	The request was not redirected. This is usually a result of a DNS request to the Traffic Router or an explicit denial for that request
DS_BYPASS
	Used a bypass destination for redirection of the :term:`Delivery Service`
DS_CLIENT_GEO_UNSUPPORTED
	Traffic Router did not find a resource supported by coverage zone data and was unable to determine the geographic location of the requesting client
DS_CZ_BACKUP_CG
	Traffic Router found a backup cache via fall-back (through the ``edgeLocation`` field of a :term:`Snapshot`)  or via coordinates (:term:`Coverage Zone File`) configuration
DS_CZ_ONLY
	The selected :term:`Delivery Service` only supports resource lookup based on coverage zone data
DS_NO_BYPASS
	No valid bypass destination is configured for the matched :term:`Delivery Service` and the :term:`Delivery Service` does not have the requested resource
DS_NOT_FOUND
	Always goes with ``rtypes`` STATIC_ROUTE and DS_MISS
GEO_NO_CACHE_FOUND
	Traffic Router could not find a resource via geographic location data based on the requesting client's location
NO_DETAILS
	This entry is for a standard request
REGIONAL_GEO_ALTERNATE_WITHOUT_CACHE
	This goes with the ``rtype`` RGDENY. The URL is being regionally blocked
REGIONAL_GEO_NO_RULE
	The request was blocked because there was no rule in the :term:`Delivery Service` for the request

HTTP Specifics
--------------
.. code-block:: text
	:caption: Sample Message

	1452197640.936 qtype=HTTP chi=69.241.53.218 rhi=- url="http://foo.mm-test.jenkins.cdnlab.comcast.net/some/asset.m3u8" cqhm=GET cqhv=HTTP/1.1 rtype=GEO rloc="40.252611,58.439389" rdtl=- rerr="-" pssc=302 ttms=0 rurl="http://odol-atsec-sim-114.mm-test.jenkins.cdnlab.comcast.net:8090/some/asset.m3u8" rh="Accept: */*" rh="myheader: asdasdasdasfasg"

.. table:: Request Fields

	+------+--------------------------------------------------------------------------------------------------------------------------------------------------+----------------------------------------------+
	| Name | Description                                                                                                                                      | Data                                         |
	+======+==================================================================================================================================================+==============================================+
	| url  | Requested URL with query string                                                                                                                  | A URL String                                 |
	+------+--------------------------------------------------------------------------------------------------------------------------------------------------+----------------------------------------------+
	| cqhm | Http Method                                                                                                                                      | e.g ``GET``, ``POST``                        |
	+------+--------------------------------------------------------------------------------------------------------------------------------------------------+----------------------------------------------+
	| cqhv | Http Protocol Version                                                                                                                            | e.g. ``HTTP/1.1``                            |
	+------+--------------------------------------------------------------------------------------------------------------------------------------------------+----------------------------------------------+
	| rh   | One or more of these key value pairs may exist in a logged event and are controlled by the configuration of the matched :term:`Delivery Service` | Key/value pair of the format ``name: value`` |
	+------+---------------------------------------------------------------------------------------------------------------------------------------------------+---------------------------------------------+

.. table:: Response Fields

	+------+-----------------------------------------------------------+
	| Name | Description                                               |
	+======+===========================================================+
	| rurl | The resulting URL of the resource requested by the client |
	+------+-----------------------------------------------------------+

DNS Specifics
-------------
.. code-block:: text
	:caption: Sample Message

	144140678.000 qtype=DNS chi=192.168.10.11 rhi=- ttms=123 xn=65535 fqdn=www.example.com. type=A class=IN ttl=12345 rcode=NOERROR rtype=CZ rloc="40.252611,58.439389" rdtl=- rerr="-" ans="192.168.1.2 192.168.3.4 0:0:0:0:0:ffff:c0a8:102 0:0:0:0:0:ffff:c0a8:304"

.. _qname: http://www.zytrax.com/books/dns/ch15/#qname

.. _qtype: http://www.zytrax.com/books/dns/ch15/#qtype

.. table:: Request Fields

	+-------+---------------------------------------------------------------------------------+---------------------------------------------------------------------------------------------------+
	| Name  | Description                                                                     | Data                                                                                              |
	+=======+=================================================================================+===================================================================================================+
	| xn    | The ID from the client DNS request header                                       | a whole number between 0 and 65535 (inclusive)                                                    |
	+-------+---------------------------------------------------------------------------------+---------------------------------------------------------------------------------------------------+
	| rhi   | The IP address of the resolver when ENDS0 client subnet extensions are enabled. | An IPv4 or IPv6 string, or dash if request is for resolver only and no client subnet is present   |
	+-------+---------------------------------------------------------------------------------+---------------------------------------------------------------------------------------------------+
	| fqdn  | The qname field from the client DNS request message (i.e. the                   | A series of DNS labels/domains separated by '.' characters and ending with a '.' character        |
	|       | :abbr:`FQDN (Fully Qualified Domain Name)` the client is requesting be          |                                                                                                   |
	+-------+---------------------------------------------------------------------------------+---------------------------------------------------------------------------------------------------+
	| type  | The qtype field from the client DNS request message (i.e. the typeof resolution | Examples are A (IpV4), AAAA (IpV6), :abbr:`NS (Name Service)`,  :abbr:`SOA (Start of Authority)`, |
	|       | that's requested such as IPv4, IPv6)                                            | and :abbr:`CNAME (Canonical Name)`, (see qtype_)                                                  |
	+-------+---------------------------------------------------------------------------------+---------------------------------------------------------------------------------------------------+
	| class | The qclass field from the client DNS request message (i.e. the class of         | Either :abbr:`IN (Internet resource)` or ANY (Traffic Router rejects requests with any other      |
	|       | resource being requested)                                                       | value of class)                                                                                   |
	+-------+---------------------------------------------------------------------------------+---------------------------------------------------------------------------------------------------+

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

.. _deep-cache:

Deep Caching
============

Overview
--------
Deep Caching is a feature that enables clients to be routed to the closest possible "deep" Edge-tier :term:`cache server` s on a per-:term:`Delivery Service` basis. The term "deep" is used in the networking sense, meaning that the Edge-tier :term:`cache server` s are located deep in the network where the number of network hops to a client is as minimal. This deep caching topology is desirable because storing content closer to the client gives better bandwidth savings, and sometimes the cost of bandwidth usage in the network outweighs the cost of adding storage. While it may not be feasible to cache an entire copy of the CDN's contents in every deep location (for the best possible bandwidth savings), storing just a relatively small amount of the CDN's most requested content can lead to very high bandwidth savings.

What You Need
-------------
#. Edge cache deployed in "deep" locations and registered in Traffic Ops
#. A :abbr:`DCZF (Deep Coverage Zone File)` mapping these deep cache hostnames to specific network prefixes (see :ref:`deep-czf` for details)
#. Deep caching parameters in the Traffic Router Profile (see :ref:`ccr-profile` for details):

	- ``deepcoveragezone.polling.interval``
	- ``deepcoveragezone.polling.url``

#. Deep Caching enabled on one or more HTTP :term:`Delivery Service`\ s (i.e. 'Deep Caching' field on the :term:`Delivery Service` details page (under :guilabel:`Advanced Options`) set to ``ALWAYS``)

How it Works
------------
Deep Coverage Zone routing is very similar to that of regular Coverage Zone routing, except that the :abbr:`DCZF (Deep Coverage Zone File)` is preferred over the regular :abbr:`CZF (Coverage Zone File)` for :term:`Delivery Service`\ s with Deep Caching enabled. If the client requests a Deep Caching-enabled :term:`Delivery Service` and their IP address gets a "hit" in the :abbr:`DCZF (Deep Coverage Zone File)`, Traffic Router will attempt to route that client to one of the available "deep" :term:`cache server` s in the client's corresponding zone. If there are no "deep" :term:`cache server` s available for a client's request, Traffic Router will fall back to the regular :abbr:`CZF (Coverage Zone File)` and continue regular :abbr:`CZF (Coverage Zone File)` routing from there.

.. _tr-steering:

Steering Feature
================

Overview
--------
A Steering :term:`Delivery Service` is a :term:`Delivery Service` that is used to route a client to another :term:`Delivery Service`. The :ref:`Type <ds-types>` of a Steering :term:`Delivery Service` is either STEERING or CLIENT_STEERING. A Steering :term:`Delivery Service` will have target :term:`Delivery Service`\ s configured for it with weights assigned to them. Traffic Router uses the weights to make a consistent hash ring which it then uses to make sure that requests are routed to a target based on the configured weights. This consistent hash ring is separate from the consistent hash ring used in cache selection.

Special regular expressions - referred to as 'filters' - can also be configured for target :term:`Delivery Service`\ s to pin traffic to a specific :term:`Delivery Service`. For example, if the filter :regexp:`.*/news/.*` for a target called ``target-ds-1`` is created, any requests to Traffic Router with "news" in them will be routed to ``target-ds-1``. This will happen regardless of the configured weights.

Some other points of interest
"""""""""""""""""""""""""""""
- Steering is currently only available for HTTP :term:`Delivery Service`\ s that are a part of the same CDN.
- A new role called STEERING has been added to the Traffic Ops database. Only users with the Steering :term:`Role` or higher can modify steering assignments for a :term:`Delivery Service`.
- Traffic Router uses the steering endpoints of the :ref:`to-api` to poll for steering assignments, the assignments are then used when routing traffic.

A couple simple use-cases for Steering are:

- Migrating traffic from one :term:`Delivery Service` to another over time.
- Trying out new functionality for a subset of traffic with an experimental :term:`Delivery Service`.
- Load balancing between :term:`Delivery Service`\ s

The Difference Between STEERING and CLIENT_STEERING
---------------------------------------------------
The only difference between the STEERING and CLIENT_STEERING :term:`Delivery Service` :term:`Type`\ s is that CLIENT_STEERING explicitly allows a client to bypass Steering by choosing a destination :term:`Delivery Service`. A client can accomplish this by providing the ``X-TC-Steering-Option`` HTTP header with a value of the ``xml_id`` of the target :term:`Delivery Service` to which they desire to be routed. When Traffic Router receives this header it will route to the requested target :term:`Delivery Service` regardless of weight configuration. This header is ignored by STEERING :term:`Delivery Service`\ s.

Configuration
-------------
The following needs to be completed for Steering to work correctly:

#. Two target :term:`Delivery Service`\ s are created in Traffic Ops. They must both be HTTP :term:`Delivery Service`\ s part of the same CDN.
#. A :term:`Delivery Service` with type STEERING or CLIENT_STEERING is created in Traffic Portal.
#. Target :term:`Delivery Service`\ s are assigned to the Steering :term:`Delivery Service` using Traffic Portal.
#. A user with the role of Steering is created.
#. The Steering user assigns weights to the target :term:`Delivery Service`\ s.
#. If desired, the Steering user can create filters for the target :term:`Delivery Service`\ s.

.. seealso:: For more information see :ref:`steering-qht`.

HTTPS for HTTP Delivery Services
================================
.. versionadded:: 1.7
	Traffic Router now has the ability to allow HTTPS traffic between itself and clients on a per-HTTP :term:`Delivery Service` basis.

.. Note:: As of version 3.0 Traffic Router has been integrated with native OpenSSL. This makes establishing HTTPS connections to Traffic Router much less expensive than previous versions. However establishing an HTTPS connection is more computationally demanding than an HTTP connection. Since each client will in turn get redirected to an :abbr:`ATS (Apache Traffic Server)` instance, Traffic Router is most always creating a new HTTPS connection for all HTTPS traffic. It is likely to mean that an existing Traffic Router may have some decrease in performance if you wish to support a lot of HTTPS traffic. As noted for :abbr:`DNSSEC (DNS Security Extensions)`, you may need to plan to scale Traffic Router vertically and/or horizontally to handle the new load.

The HTTPS set up process is:

#. Select one of '1 - HTTPS', '2 - HTTP AND HTTPS', or '3 - HTTP TO HTTPS' for the :term:`Delivery Service`
#. Generate private keys for the :term:`Delivery Service` using a wildcard domain such as ``*.my-delivery-service.my-cdn.example.com``
#. Obtain and import signed certificate chain
#. Perform a CDN :term:`Snapshot`

Clients may make HTTPS requests to :term:`Delivery Service`\ s only after the CDN :term:`Snapshot` propagates to Traffic Router and it receives the certificate chain from Traffic Ops.

Protocol Options
----------------
HTTP
	Any secure client will get an SSL handshake error. Non-secure clients will experience the same behavior as prior to 1.7
HTTPS
	Traffic Router will only redirect (send a ``302 Found`` response) to clients communicating with a secure connection, all other clients will receive a ``503 Service Unavailable`` response
HTTP AND HTTPS
	Traffic Router will redirect both secure and non-secure clients
HTTP TO HTTPS
	Traffic Router will redirect non-secure clients with a ``302 Found`` response and a location that is secure (i.e. an ``https://`` URL instead of an ``http://`` URL), while secure clients will be redirected immediately to an appropriate target or :term:`cache server`.

Certificate Retrieval
---------------------
.. Warning:: If you have HTTPS :term:`Delivery Service`\ s in your CDN, Traffic Router will not accept **any** connections until it is able to fetch certificates from Traffic Ops and load them into memory. Traffic Router does not persist certificates to the Java Keystore or anywhere else.

Traffic Router fetches certificates into memory:

* At startup time
* When it receives a new CDN :term:`Snapshot`
* Once an hour starting whenever the most recent of the last of the above occurred

.. Note:: To adjust the frequency at which Traffic Router fetches certificates add the :term:`Parameter` ``certificates.polling.interval`` with the ConfigFile "CRConfig.json" and set it to the desired duration in milliseconds.

.. Note:: Taking a CDN :term:`Snapshot` may be used at times to avoid waiting the entire polling cycle for a new set of certificates.

.. Warning:: If a CDN :term:`Snapshot` is taken that involves a :term:`Delivery Service` missing its certificates, Traffic Router will ignore **ALL** changes in that CDN :term:`Snapshot` until one of the following occurs:

	* It receives certificates for that :term:`Delivery Service`
	* Another CDN :term:`Snapshot` is taken and the :term:`Delivery Service` without certificates is changed such that its HTTP protocol is set to 'http'

Certificate Chain Ordering
--------------------------
The ordering of certificates within the certificate bundle matters. It must be:

#. Primary Certificate (e.g. the one created for ``*.my-delivery-service.my-cdn.example.com``)
#. Intermediate Certificate(s)
#. Root Certificate from a :abbr:`CA (Certificate Authority)` (optional)

.. Warning:: If something is wrong with the certificate chain (e.g. the order of the certificates is backwards or for the wrong domain) the client will get an SSL handshake. Inspection of ``/opt/tomcat/logs/catalina.log`` is likely to yield information to reveal this.

To see the ordering of certificates you may have to manually split up your certificate chain and use :manpage:`openssl(1ssl)` on each individual certificate

Suggested Way of Setting up an HTTPS Delivery Service
-----------------------------------------------------
Assuming you have already created a :term:`Delivery Service` which you plan to modify to use HTTPS, do the following in Traffic Portal:

#. Select one of '1 - HTTPS', '2 - HTTP AND HTTPS', or '3 - HTTP TO HTTPS' for the protocol field of a :term:`Delivery Service` and click the :guilabel:`Update` button
#. Go to :menuselection:`More --> Manage SSL Keys`
#. Click on :menuselection:`More --> Generate SSL Keys`
#. Fill out the form and click on the green :guilabel:`Generate Keys` button, then confirm that you want to make these changes
#. Copy the contents of the Certificate Signing Request field and save it locally
#. Go back and select 'HTTP' for the protocol field of the :term:`Delivery Service` and click :guilabel:`Save` (to avoid preventing other CDN :term:`Snapshot` updates from being blocked by Traffic Router)
#. Follow your standard procedure for obtaining your signed certificate chain from a :abbr:`CA (Certificate Authority)`
#. After receiving your certificate chain import it into Traffic Ops
#. Edit the :term:`Delivery Service`
#. Restore your original choice for the protocol field and click :guilabel:`Save`
#. Click :menuselection:`More --> Manage SSL Keys`
#. Paste your key information into the appropriate fields
#. Click the green :guilabel:`Update Keys` button
#. Take a new CDN :term:`Snapshot`

Once this is done you should be able to verify that you are being correctly redirected by Traffic Router using e.g. :manpage:`curl(1)` commands to HTTPS destinations on your :term:`Delivery Service`.

Router Load Testing
===================
The Traffic Router load testing tool is located in the `Traffic Control repository under test/router <https://github.com/apache/trafficcontrol/tree/master/test/router>`_. It can be used to simulate a mix of HTTP and HTTPS traffic for a CDN by choosing the number of HTTP :term:`Delivery Service`\ s and the number HTTPS :term:`Delivery Service` the test will exercise.

There are 2 parts to the load test:

* A web server that makes the actual requests and takes commands to fetch data from the CDN, start the test, and return current results.
* A web page that's used to run the test and see the results.

Running the Load Tests
----------------------
#. First, clone the `Traffic Control repository <https://github.com/apache/trafficcontrol>`_.
#. You will need to make sure you have a :abbr:`CA (Certificate Authority)` file on your machine
#. The web server is a Go program, set your ``GOPATH`` environment variable appropriately (we suggest ``$HOME/go`` or ``$HOME/src``)
#. Open a terminal emulator and navigate to the ``test/router/server`` directory inside of the cloned repository
#. Execute the server binary by running ``go run server.go``
#. Using your web browser of choice, open the file ``test/router/index.html``
#. Authenticate against a Traffic Ops host - this should be a nearly instantaneous operation - you can watch the output from ``server.go`` for feedback
#. Enter the Traffic Ops host in the second form and click the button to get a list of CDN's
#. Wait for the web page to show a list of CDN's under the above form, this may take several seconds
#. The List of CDN's will display the number of HTTP- and HTTPS-capable :term:`Delivery Service`\ s that may be exercised
#. Choose the CDN you want to exercise from the drop-down menu
#. Fill out the rest of the form, enter appropriate numbers for each HTTP and HTTPS :term:`Delivery Service`\ s
#. Click :guilabel:`Run Test`
#. As the test runs the web page will occasionally report results including running time, latency, and throughput

Tuning Recommendations
======================
The following is an example of the command line parameters set in :file:`/opt/traffic_router/conf/startup.properties` that has been tested on a multi-core server running under HTTPS load test requests. This is following the general recommendation to use the G1 garbage collector for :abbr:`JVM (Java Virtual Machine)` applications running on multi-core machines. In addition to using the G1 garbage collector the ``InitiatingHeapOccupancyPercent`` was lowered to run garbage collection more frequently which improved overall throughput for Traffic Router and reduced 'Stop the World' garbage collection. Note that any environment variable settings in this file will override those set in :file:`/lib/systemd/system/traffic_router.service`.

.. code-block:: bash
	:caption: Example CATALINA_OPTS Configuration

	CATALINA_OPTS="\
	-server -Xms2g -Xmx8g \
	-Dlog4j.configuration=file://$CATALINA_BASE/conf/log4j.properties \
	-Djava.library.path=/usr/lib64 \
	-XX:+UseG1GC \
	-XX:+UnlockExperimentalVMOptions \
	-XX:InitiatingHeapOccupancyPercent=30"
