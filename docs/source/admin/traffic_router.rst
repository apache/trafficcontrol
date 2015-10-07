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
