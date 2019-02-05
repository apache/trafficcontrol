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

.. |graph| image:: images/graph.png
.. |info| image:: images/info.png
.. |checkmark| image:: images/good.png
.. |X| image:: images/bad.png
.. |clock| image:: images/clock-black.png

.. _to-using:

*******************
Traffic Ops - Using
*******************

.. deprecated:: 3.0
	The Traffic Ops UI is deprecated, and will be removed entirely in the next major release (4.0). A much better way to interact with the CDN is to :ref:`use Traffic Portal <usingtrafficportal>`, which is the the only UI that will be receiving updates for the foreseeable future.

The Traffic Ops Menu
====================
.. figure:: images/12m.png
	:align: center
	:alt: The Traffic Ops Landing Page

	The Traffic Ops Landing Page

The following tabs are available in the menu at the top of the Traffic Ops user interface.

.. index::
	Health Tab

Health
------
Information on the health of the system. Hover over this tab to get to the following options:

+---------------+------------------------------------------------------------------------------------------------------------------------------------+
|     Option    |                                                            Description                                                             |
+===============+====================================================================================================================================+
| Table View    | A real time view into the main performance indicators of the CDNs managed by Traffic Control.                                      |
|               | This view is sourced directly by the Traffic Monitor data and is updated every 10 seconds.                                         |
|               | This is the default screen of Traffic Ops.                                                                                         |
|               | See :ref:`health-table` for details.                                                                                               |
+---------------+------------------------------------------------------------------------------------------------------------------------------------+
| Graph View    | A real graphical time view into the main performance indicators of the CDNs managed by Traffic Control.                            |
|               | This view is sourced by the Traffic Monitor data and is updated every 10 seconds.                                                  |
|               | On loading, this screen will show a history of 24 hours of data from Traffic Stats                                                 |
|               | See :ref:`health-graph` for details.                                                                                               |
+---------------+------------------------------------------------------------------------------------------------------------------------------------+
| Server Checks | A table showing the results of the periodic check extension scripts that are run. See :ref:`server-checks`                         |
+---------------+------------------------------------------------------------------------------------------------------------------------------------+
| Daily Summary | A graph displaying the daily peaks of bandwidth, overall bytes served per day, and overall bytes served since initial installation |
|               | per CDN.                                                                                                                           |
+---------------+------------------------------------------------------------------------------------------------------------------------------------+


Delivery Services
-----------------
The main :term:`Delivery Service` table. This is where you Create/Read/Update/Delete :term:`Delivery Service`\ s of all types. Hover over to get the following sub option:

+-------------+--------------------------------------+
|    Option   |             Description              |
+=============+======================================+
| Federations | Add/Edit/Delete Federation Mappings. |
+-------------+--------------------------------------+


Servers
-------
The main Servers table. This is where you Create/Read/Update/Delete servers of all types.  Click the main tab to get to the main table, and hover over to get these sub options:

+-------------------+--------------------------------------+
|       Option      | Description                          |
+===================+======================================+
| Upload Server CSV | Bulk add of servers from a CSV file. |
+-------------------+--------------------------------------+

Parameters
----------
Parameters and Profiles can be edited here. Hover over the tab to get the following options:

+-----------------------------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|        Option               |                                                                             Description                                                                                             |
+=============================+=====================================================================================================================================================================================+
| Global Profile              | The table of global parameters. See :ref:`param-prof`. This is where you Create/Read/Update/Delete parameters in the Global profile                                                 |
+-----------------------------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| All :term:`Cache Group`\ s  | The table of all parameters *that are assigned to a cachegroup* - this may be slow to pull up, as there can be thousands of parameters.                                             |
+-----------------------------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| All Profiles                | The table of all parameters *that are assigned to a profile* - this may be slow to pull up, as there can be thousands of parameters.                                                |
+-----------------------------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Select Profile              | Select the parameter list by profile first, then get a table of just the parameters for that profile.                                                                               |
+-----------------------------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Export Profile              | Profiles can be exported from one Traffic Ops instance to another using 'Select Profile' and under the "Profile Details" dialog for the desired profile                             |
+-----------------------------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Import Profile              | Profiles can be imported from one Traffic Ops instance to another using the button "Import Profile" after using the "Export Profile" feature                                        |
+-----------------------------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Orphaned Parameters         | A table of parameters that are not associated to any profile of :term:`Cache Group`. These parameters either should be deleted or associated with a profile of :term:`Cache Group`. |
+-----------------------------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+

Tools
-----
Tools for working with Traffic Ops and it's servers. Hover over this tab to get the following options:

+--------------------+-----------------------------------------------------------------------------------------------------------------------------------+
|       Option       |                                                            Description                                                            |
+====================+===================================================================================================================================+
| Generate ISO       | Generate a bootable image for any of the servers in the Servers table (or any server for that matter). See :ref:`generate-iso`    |
+--------------------+-----------------------------------------------------------------------------------------------------------------------------------+
| Queue Updates      | Send Updates to the caches. See :ref:`queue-updates`                                                                              |
+--------------------+-----------------------------------------------------------------------------------------------------------------------------------+
| DB Dump            | Backup the Database to a .sql file.                                                                                               |
+--------------------+-----------------------------------------------------------------------------------------------------------------------------------+
| Snapshot CRConfig  | Send updates to the Traffic Monitor / Traffic Router servers.  See :ref:`queue-updates`                                           |
+--------------------+-----------------------------------------------------------------------------------------------------------------------------------+
| Invalidate Content | Invalidate or purge content from all caches in the CDN. See :ref:`purge`                                                          |
+--------------------+-----------------------------------------------------------------------------------------------------------------------------------+
| Manage DNSSEC keys | Manage DNSSEC Keys for a chosen CDN.                                                                                              |
+--------------------+-----------------------------------------------------------------------------------------------------------------------------------+


Misc
----
Miscellaneous editing options. Hover over this tab to get the following options:

+------------------------------+-------------------------------------------------------------------------------------------+
|       Option                 |                                        Description                                        |
+==============================+===========================================================================================+
| CDNs                         | Create/Read/Update/Delete CDNs                                                            |
+------------------------------+-------------------------------------------------------------------------------------------+
| :term:`Cache Group`\ s       | Create/Read/Update/Delete :term:`Cache Group`\ s                                          |
+------------------------------+-------------------------------------------------------------------------------------------+
| Users                        | Create/Read/Update/Delete users                                                           |
+------------------------------+-------------------------------------------------------------------------------------------+
| Profiles                     | Create/Read/Update/Delete profiles. See :ref:`working-with-profiles`                      |
+------------------------------+-------------------------------------------------------------------------------------------+
| Networks(ASNs)               | Create/Read/Update/Delete Autonomous System Numbers See :ref:`asn-czf`                    |
+------------------------------+-------------------------------------------------------------------------------------------+
| Hardware                     | Get detailed hardware information (note: this should be moved to a Traffic Ops Extension) |
+------------------------------+-------------------------------------------------------------------------------------------+
| Data Types                   | Create/Read/Update/Delete data types                                                      |
+------------------------------+-------------------------------------------------------------------------------------------+
| Divisions                    | Create/Read/Update/Delete divisions                                                       |
+------------------------------+-------------------------------------------------------------------------------------------+
| Regions                      | Create/Read/Update/Delete regions                                                         |
+------------------------------+-------------------------------------------------------------------------------------------+
| Physical Locations           | Create/Read/Update/Delete locations                                                       |
+------------------------------+-------------------------------------------------------------------------------------------+

.. index::
	Change Log

ChangeLog
---------
The Changelog table displays the changes that are being made to the Traffic Ops database through the Traffic Ops user interface. This tab will show the number of changes since you last visited this tab in (brackets) since the last time you visited this tab. There are currently no sub menus for this tab.


Help
----
Help for Traffic Ops and Traffic Control. Hover over this tab to get the following options:

+---------------+---------------------------------------------------------------------+
|     Option    |                             Description                             |
+===============+=====================================================================+
| About         | Traffic Ops information, such as version, database information, etc |
+---------------+---------------------------------------------------------------------+
| Release Notes | Release notes for the most recent releases of Traffic Ops           |
+---------------+---------------------------------------------------------------------+
| Logout        | Logout from Traffic Ops                                             |
+---------------+---------------------------------------------------------------------+


.. index::
	Edge Health
	Health

Health
======

.. _health-table:

The Health Table
----------------
The Health table is the default landing screen for Traffic Ops, it displays the status of the EDGE caches in a table form directly from Traffic Monitor (bypassing Traffic Stats), sorted by Mbps Out. The columns in this table are:


:Profile:          the Profile of this server or ALL, meaning this row shows data for multiple servers, and the row shows the sum of all values.
:Edge Cache Group: the edge :term:`Cache Group` short name or ALL, meaning this row shows data for multiple servers, and the row shows the sum of all values.
:Host Name:        the host name of the server or ALL, meaning this row shows data for multiple servers, and the row shows the sum of all values.
:Healthy:          indicates if this cache is healthy according to the Health Protocol. A row with ALL in any of the columns will always show a |checkmark|, this column is valid only for individual EDGE caches.
:Admin:            shows the administrative status of the server.
:Connections:      the number of connections this cache (or group of caches) has open (``ats.proxy.process.http.current_client_connections`` from ATS).
:Mbps Out:         the bandwidth being served out if this cache (or group of caches)

Since the top line has ALL, ALL, ALL, it shows the total connections and bandwidth for all caches managed by this instance of Traffic Ops.

.. _health-graph:

Graph View
----------
The Graph View shows a live view of the last 24 hours of bits per seconds served and open connections at the edge in a graph. This data is sourced from Traffic Stats. If there are 2 CDNs configured, this view will show the statistis for both, and the graphs are stacked. On the left-hand side, the totals and immediate values as well as the percentage of total possible capacity are displayed. This view is update every 10 seconds.


.. _server-checks:

Server Checks
-------------
The server checks page is intended to give an overview of the Servers managed by Traffic Control as well as their status. This data comes from `Traffic Ops extensions <traffic_ops_extensions.html>`_.

+------+-----------------------------------------------------------------------+
| Name |                 Description                                           |
+======+=======================================================================+
| ILO  | Ping the iLO interface for EDGE or MID servers                        |
+------+-----------------------------------------------------------------------+
| 10G  | Ping the IPv4 address of the EDGE or MID servers                      |
+------+-----------------------------------------------------------------------+
| 10G6 | Ping the IPv6 address of the EDGE or MID servers                      |
+------+-----------------------------------------------------------------------+
| MTU  | Ping the EDGE or MID using the configured MTU from Traffic Ops        |
+------+-----------------------------------------------------------------------+
| FQDN | DNS check that matches what the DNS servers responds with compared to |
|      | what Traffic Ops has.                                                 |
+------+-----------------------------------------------------------------------+
| DSCP | Checks the DSCP value of packets from the edge server to the Traffic  |
|      | Ops server.                                                           |
+------+-----------------------------------------------------------------------+
| RTR  | Content Router checks. Checks the health of the Content Routers.      |
|      | Checks the health of the caches using the Content Routers.            |
+------+-----------------------------------------------------------------------+
| CHR  | Cache Hit Ratio in percent.                                           |
+------+-----------------------------------------------------------------------+
| CDU  | Total Cache Disk Usage in percent.                                    |
+------+-----------------------------------------------------------------------+
| ORT  | Operational Readiness Test. Uses the ORT script on the edge and mid   |
|      | servers to determine if the configuration in Traffic Ops matches the  |
|      | configuration on the edge or mid. The user that this script runs as   |
|      | must have an ssh key on the edge servers.                             |
+------+-----------------------------------------------------------------------+

Daily Summary
-------------
Displays daily max gbps and bytes served for all CDNs.  In order for the graphs to appear, the 'daily_bw_url' and 'daily_served_url' parameters need to be be created, assigned to the global profile, and have a value of a grafana graph.  For more information on configuring grafana, see the `Traffic Stats <../traffic_stats.html>`_  section.

.. _server:

Server
======
This view shows a table of all the servers in Traffic Ops. The table columns show the most important details of the server. The **IPAddrr** column is clickable to launch an ``ssh://`` link to this server. The |graph| icon will link to a Traffic Stats graph of this server for caches, and the |info| will link to the server status pages for other server types.


Server Types
------------
These are the types of servers that can be managed in Traffic Ops:

+---------------+---------------------------------------------+
|      Name     |                 Description                 |
+===============+=============================================+
| EDGE          | Edge Cache                                  |
+---------------+---------------------------------------------+
| MID           | Mid Tier Cache                              |
+---------------+---------------------------------------------+
| ORG           | Origin                                      |
+---------------+---------------------------------------------+
| CCR           | Traffic Router                              |
+---------------+---------------------------------------------+
| RASCAL        | Rascal health polling & reporting           |
+---------------+---------------------------------------------+
| TOOLS_SERVER  | Ops hosts for managment                     |
+---------------+---------------------------------------------+
| RIAK          | Riak keystore                               |
+---------------+---------------------------------------------+
| SPLUNK        | SPLUNK indexer search head etc              |
+---------------+---------------------------------------------+
| TRAFFIC_STATS | traffic_stats server                        |
+---------------+---------------------------------------------+
| INFLUXDB      | influxDb server                             |
+---------------+---------------------------------------------+

.. _working-with-profiles:

Parameters and Profiles
=======================
Parameters are shared between profiles if the set of ``{ name, config_file, value }`` is the same. To change a value in one profile but not in others, the parameter has to be removed from the profile you want to change it in, and a new parameter entry has to be created (**Add Parameter** button at the bottom of the Parameters view), and assigned to that profile. It is easy to create new profiles from the **Misc > Profiles** view - just use the **Add/Copy Profile** button at the bottom of the profile view to copy an existing profile to a new one. Profiles can be exported from one system and imported to another using the profile view as well. It makes no sense for a parameter to not be assigned to a single profile - in that case it really has no function. To find parameters like that use the **Parameters > Orphaned Parameters** view. It is easy to create orphaned parameters by removing all profiles, or not assigning a profile directly after creating the parameter.

.. seealso:: :ref:`param-prof` in the *Configuring Traffic Ops* section.

.. _ccr-profile:

Traffic Router Profile
----------------------

+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
|                   Name                  |      Config_file       |                                                                  Description                                                                     |
+=========================================+========================+==================================================================================================================================================+
| location                                | dns.zone               | Location to store the DNS zone files in the local file system of Traffic Router.                                                                 |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| location                                | http-log4j.properties  | Location to find the log4j.properties file for Traffic Router.                                                                                   |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| location                                | dns-log4j.properties   | Location to find the dns-log4j.properties file for Traffic Router.                                                                               |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| location                                | geolocation.properties | Location to find the log4j.properties file for Traffic Router.                                                                                   |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| CDN_name                                | rascal-config.txt      | The human readable name of the CDN for this profile.                                                                                             |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| CoverageZoneJsonURL                     | CRConfig.xml           | The location (URL) to retrieve the coverage zone map file in JSON format from.                                                                   |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| ecsEnable                               | CRConfig.json          | Boolean value to enable or disable ENDS0 client subnet extensions.                                                                               |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| geolocation.polling.url                 | CRConfig.json          | The location (URL) to retrieve the geo database file from.                                                                                       |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| geolocation.polling.interval            | CRConfig.json          | How often to refresh the coverage geo location database  in ms                                                                                   |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| coveragezone.polling.interval           | CRConfig.json          | How often to refresh the coverage zone map in ms                                                                                                 |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| coveragezone.polling.url                | CRConfig.json          | The location (URL) to retrieve the coverage zone map file in JSON format from.                                                                   |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| deepcoveragezone.polling.interval       | CRConfig.json          | How often to refresh the deep coverage zone map in ms                                                                                            |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| deepcoveragezone.polling.url            | CRConfig.json          | The location (URL) to retrieve the deep coverage zone map file in JSON format from.                                                              |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| tld.soa.expire                          | CRConfig.json          | The value for the expire field the Traffic Router DNS Server will respond with on Start of Authority (SOA) records.                              |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| tld.soa.minimum                         | CRConfig.json          | The value for the minimum field the Traffic Router DNS Server will respond with on SOA records.                                                  |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| tld.soa.admin                           | CRConfig.json          | The DNS Start of Authority admin.  Should be a valid support email address for support if DNS is not working correctly.                          |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| tld.soa.retry                           | CRConfig.json          | The value for the retry field the Traffic Router DNS Server will respond with on SOA records.                                                    |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| tld.soa.refresh                         | CRConfig.json          | The TTL the Traffic Router DNS Server will respond with on A records.                                                                            |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| tld.ttls.NS                             | CRConfig.json          | The TTL the Traffic Router DNS Server will respond with on NS records.                                                                           |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| tld.ttls.SOA                            | CRConfig.json          | The TTL the Traffic Router DNS Server will respond with on SOA records.                                                                          |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| tld.ttls.AAAA                           | CRConfig.json          | The Time To Live (TTL) the Traffic Router DNS Server will respond with on AAAA records.                                                          |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| tld.ttls.A                              | CRConfig.json          | The TTL the Traffic Router DNS Server will respond with on A records.                                                                            |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| tld.ttls.DNSKEY                         | CRConfig.json          | The TTL the Traffic Router DNS Server will respond with on DNSKEY records.                                                                       |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| tld.ttls.DS                             | CRConfig.json          | The TTL the Traffic Router DNS Server will respond with on DS records.                                                                           |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| api.port                                | server.xml             | The TCP port Traffic Router listens on for API (REST) access.                                                                                    |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| api.cache-control.max-age               | CRConfig.json          | The value of the ``Cache-Control: max-age=`` header in the API responses of Traffic Router.                                                      |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| api.auth.url                            | CRConfig.json          | The API authentication URL (https://${tmHostname}/api/1.1/user/login); ${tmHostname} is a search and replace token used by Traffic Router to     |
|                                         |                        | construct the correct URL)                                                                                                                       |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| consistent.dns.routing                  | CRConfig.json          | Control whether DNS :term:`Delivery Service`\ s use consistent hashing on the edge FQDN to select caches for answers. May improve performance if |
|                                         |                        | set to true; defaults to false                                                                                                                   |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| dnssec.enabled                          | CRConfig.json          | Whether DNSSEC is enabled; this parameter is updated via the DNSSEC administration user interface.                                               |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| dnssec.allow.expired.keys               | CRConfig.json          | Allow Traffic Router to use expired DNSSEC keys to sign zones; default is true. This helps prevent DNSSEC related outages due to failed Traffic  |
|                                         |                        | Control components or connectivity issues.                                                                                                       |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| dynamic.cache.primer.enabled            | CRConfig.json          | Allow Traffic Router to attempt to prime the dynamic zone cache; defaults to true                                                                |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| dynamic.cache.primer.limit              | CRConfig.json          | Limit the number of permutations to prime when dynamic zone cache priming is enabled; defaults to 500                                            |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| keystore.maintenance.interval           | CRConfig.json          | The interval in seconds which Traffic Router will check the keystore API for new DNSSEC keys                                                     |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| keystore.api.url                        | CRConfig.json          | The keystore API URL (https://${tmHostname}/api/1.1/cdns/name/${cdnName}/dnsseckeys.json; ${tmHostname} and ${cdnName} are search and replace    |
|                                         |                        | tokens used by Traffic Router to construct the correct URL)                                                                                      |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| keystore.fetch.timeout                  | CRConfig.json          | The timeout in milliseconds for requests to the keystore API                                                                                     |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| keystore.fetch.retries                  | CRConfig.json          | The number of times Traffic Router will attempt to load keys before giving up; defaults to 5                                                     |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| keystore.fetch.wait                     | CRConfig.json          | The number of milliseconds Traffic Router will wait before a retry                                                                               |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| signaturemanager.expiration.multiplier  | CRConfig.json          | Multiplier used in conjunction with a zone's maximum TTL to calculate DNSSEC signature durations; defaults to 5                                  |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| zonemanager.threadpool.scale            | CRConfig.json          | Multiplier used to determine the number of cores to use for zone signing operations; defaults to 0.75                                            |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| zonemanager.cache.maintenance.interval  | CRConfig.json          | The interval in seconds which Traffic Router will check for zones that need to be resigned or if dynamic zones need to be expired from cache     |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| zonemanager.dynamic.response.expiration | CRConfig.json          | A string (e.g.: 300s) that defines how long a dynamic zone                                                                                       |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| DNSKEY.generation.multiplier            | CRConfig.json          | Used to deteremine when new keys need to be regenerated. Keys are regenerated if expiration is less than the generation multiplier * the TTL. If |
|                                         |                        | the parameter does not exist, the default is 10.                                                                                                 |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+
| DNSKEY.effective.multiplier             | CRConfig.json          | Used when creating an effective date for a new key set.  New keys are generated with an effective date of old key expiration - (effective        |
|                                         |                        | multiplier * TTL).  Default is 2.                                                                                                                |
+-----------------------------------------+------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------+

Tools
=====

.. index::
	ISO
	Generate ISO

.. _generate-iso:

Generate ISO
------------
Generate ISO is a tool for building custom ISOs for building caches on remote hosts. Currently it only supports Centos 7, but if you're brave and pure of heart you MIGHT be able to get it to work with other unix-like OS's.

The interface is *mostly* self-explanatory as it's got hints.

+-------------------------------+---------------------------------------------------------------------------------------------------------------------------------+
| Field                         |  Explaination                                                                                                                   |
+===============================+=================================================================================================================================+
|Choose a server from list:     | This option gets all the server names currently in the Traffic Ops database and will autofill known values.                     |
+-------------------------------+---------------------------------------------------------------------------------------------------------------------------------+
| OS Version:                   | There needs to be an _osversions.cfg_ file in the ISO directory that maps the name of a directory to a name that shows up here. |
+-------------------------------+---------------------------------------------------------------------------------------------------------------------------------+
| Hostname:                     | This is the FQDN of the server to be installed. It is required.                                                                 |
+-------------------------------+---------------------------------------------------------------------------------------------------------------------------------+
| Root password:                | If you don't put anything here it will default to the salted MD5 of "Fred". Whatever put is MD5 hashed and writte to disk.      |
+-------------------------------+---------------------------------------------------------------------------------------------------------------------------------+
| DHCP:                         | if yes, other IP settings will be ignored                                                                                       |
+-------------------------------+---------------------------------------------------------------------------------------------------------------------------------+
| IP Address:                   | Required if DHCP=no                                                                                                             |
+-------------------------------+---------------------------------------------------------------------------------------------------------------------------------+
| Netmask:                      | Required if DHCP=no                                                                                                             |
+-------------------------------+---------------------------------------------------------------------------------------------------------------------------------+
| Gateway:                      | Required if DHCP=no                                                                                                             |
+-------------------------------+---------------------------------------------------------------------------------------------------------------------------------+
| IPV6 Address:                 | Optional. /64 is assumed if prefix is omitted                                                                                   |
+-------------------------------+---------------------------------------------------------------------------------------------------------------------------------+
| IPV6 Gateway:                 | Ignored if an IPV4 gateway is specified                                                                                         |
+-------------------------------+---------------------------------------------------------------------------------------------------------------------------------+
| Network Device:               | Optional. Typical values are bond0, eth4, etc. Note: if you enter bond0, a LACP bonding config will be written                  |
+-------------------------------+---------------------------------------------------------------------------------------------------------------------------------+
| MTU:                          | If unsure, set to 1500                                                                                                          |
+-------------------------------+---------------------------------------------------------------------------------------------------------------------------------+
| Specify disk for OS install:  | Optional. Typical values are "sda".                                                                                             |
+-------------------------------+---------------------------------------------------------------------------------------------------------------------------------+


When you click the **Download ISO** button the folling occurs (all paths relative to the top level of the directory specified in _osversions.cfg_):

#. Reads /etc/resolv.conf to get a list of nameservers. This is a rather ugly hack that is in place until we get a way of configuring it in the interface.
#. Writes a file in the ks_scripts/state.out that contains directory from _osversions.cfg_ and the mkisofs string that we'll call later.
#. Writes a file in the ks_scripts/network.cfg that is a bunch of key=value pairs that set up networking.
#. Creates an MD5 hash of the password you specify and writes it to ks_scripts/password.cfg. Note that if you do not specify a password "Fred" is used. Also note that we have experienced some issues with webbrowsers autofilling that field.
#. Writes out a disk configuration file to ks_scripts/disk.cfg.
#. mkisofs is called against the directory configured in _osversions.cfg_ and an ISO is generated in memory and delivered to your webbrowser.

You now have a customized ISO that can be used to install Red Hat and derivative Linux installations with some modifications to your ks.cfg file.

Kickstart/Anaconda will mount the ISO at /mnt/stage2 during the install process (at least with 6).

You can directly include the password file anywhere in your ks.cfg file (usually in the top) by doing %include /mnt/stage2/ks_scripts/password.cfg

What we currently do is have 2 scripts, one to do hard drive configuration and one to do network configuration. Both are relatively specific to the environment they were created in, and both are *probably* wrong for other organizations, however they are currently living in the "misc" directory as examples of how to do things.

We trigger those in a %pre section in ks.cfg and they will write config files to /tmp. We will then include those files in the appropriate places using  %pre.

For example this is a section of our ks.cfg file: ::

	%include /mnt/stage2/ks_scripts/packages.txt

	%pre
		python /mnt/stage2/ks_scripts/create_network_line.py
		bash /mnt/stage2/ks_scripts/drive_config.sh
	%end

These two scripts will then run _before_ anaconda sets up it's internal structures, then a bit further up in the ks.cfg file (outside of the %pre %end block) we do an ::

	%include /mnt/stage2/ks_scripts/password.cfg
	...
	%include /tmp/network_line

	%include /tmp/drive_config
	...

This snarfs up the contents and inlines them.

If you only have one kind of hardware on your CDN it is probably best to just put the drive config right in the ks.cfg.

If you have simple networking needs (we use bonded interfaces in most, but not all locations and we have several types of hardware meaning different ethernet interface names at the OS level etc.) then something like this:

.. code-block:: bash

	#!/bin/bash
	source /mnt/stage2/ks_scripts/network.cfg
	echo "network --bootproto=static --activate --ipv6=$IPV6ADDR --ip=$IPADDR --netmask=$NETMASK --gateway=$GATEWAY --ipv6gateway=$GATEWAY --nameserver=$NAMESERVER --mtu=$MTU --hostname=$HOSTNAME" >> /tmp/network.cfg

,, Note:: that this is an example and may not work at all.

You could also put this in the %pre section. Lots of ways to solve it.

We have included the two scripts we use in the "misc" directory of the git repo:

* kickstart_create_network_line.py
* kickstart_drive_config.sh

These scripts were written to support a very narrow set of expectations and environment and are almost certainly not suitable to just drop in, but they might provide a good starting point.

.. _queue-updates:

Queue Updates and Snapshot CRConfig
-----------------------------------
When changing delivery services special care has to be taken so that Traffic Router will not send traffic to caches for delivery services that the cache doesn't know about yet. In general, when adding delivery services, or adding servers to a delivery service, it is best to update the caches before updating Traffic Router and Traffic Monitor. When deleting delivery services, or deleting server assignments to delivery services, it is best to update Traffic Router and Traffic Monitor first and then the caches. Updating the cache configuration is done through the *Queue Updates* menu, and updating Traffic Monitor and  Traffic Router config is done through the *Snapshot CRConfig* menu.

.. index::
	Cache Updates
	Queue Updates

Queue Updates
"""""""""""""
Every 15 minutes the caches should run a *syncds* to get all changes needed from Traffic Ops. The files that will be updated by the syncds job are:

- records.config
- remap.config
- parent.config
- cache.config
- hosting.config
- url\_sig\_(.*)\.config
- hdr\_rw\_(.*)\.config
- regex_revalidate.config
- ip_allow.config

A cache will only get updated when the update flag is set for it. To set the update flag, use the *Queue Updates* menu - here you can schedule updates for a whole CDN or a :term:`Cache Group`:

#. Click **Tools > Queue Updates**.
#. Select the CDN to queue updates for or select All.
#. Select the :term:`Cache Group` to queue updates for or select All.
#. Click the **Queue Updates** button.
#. When the Queue Updates for this Server? (all) window opens, click **OK**.

To schedule updates for just one cache, use the "Server Checks" page, and click the |checkmark| in the *UPD* column. The UPD column of Server Checks page will change show a |clock| when updates are pending for that cache.

.. index::
	Snapshot CRConfig

.. _snapshot-crconfig:

Snapshot CRConfig
"""""""""""""""""
Every 60 seconds Traffic Monitor will check with Traffic Ops to see if a new CRConfig snapshot exists; Traffic Monitor polls Traffic Ops for a new CRConfig, and Traffic Router polls Traffic Monitor for the same file. This is necessary to ensure that Traffic Monitor sees configuration changes first, which helps to ensure that the health and state of caches and delivery services propagates properly to Traffic Router. See :ref:`ccr-profile` for more information on the CRConfig file.

To create a new snapshot, use the *Tools > Snapshot CRConfig* menu:

	#. Click **Tools > Snapshot CRConfig**.
	#. Verify the selection of the correct CDN from the Choose CDN drop down and click **Diff CRConfig**.
		 On initial selection of this, the CRConfig Diff window says the following:

		 There is no existing CRConfig for [cdn] to diff against... Is this the first snapshot???
		 If you are not sure why you are getting this message, please do not proceed!
		 To proceed writing the snapshot anyway click the 'Write CRConfig' button below.

		 If there is an older version of the CRConfig, a window will pop up showing the differences
		 between the active CRConfig and the CRConfig about to be written.

	#. Click **Write CRConfig**.
	#. When the This will push out a new CRConfig.json. Are you sure? window opens, click **OK**.
	#. The "Successfully wrote CRConfig.json!" window opens, click **OK**.

.. Note:: Snapshotting the CDN also deletes all HTTPS certificates for every :term:`Delivery Service` which has been deleted since the last :term:`Snapshot`.

.. index::
	Invalidate Content
	Purge

.. _purge:

Invalidate Content
==================
Invalidating content on the CDN is sometimes necessary when the origin was mis-configured and something is cached in the CDN  that needs to be removed. Given the size of a typical Traffic Control CDN and the amount of content that can be cached in it, removing the content from all the caches may take a long time. To speed up content invalidation, Traffic Ops will not try to remove the content from the caches, but it makes the content inaccessible using the *regex_revalidate* ATS plugin. This forces a *revalidation* of the content, rather than a new get.

.. Note:: This method forces a HTTP *revalidation* of the content, and not a new *GET* - the origin needs to support revalidation according to the HTTP/1.1 specification, and send a ``200 OK`` or ``304 Not Modified`` as applicable.

To invalidate content:

#. Click **Tools > Invalidate Content**
#. Fill out the form fields:

	- Select the *:term:`Delivery Service`**
	- Enter the **Path Regex** - this should be a `PCRE <http://www.pcre.org/>`_ compatible regular expression for the path to match for forcing the revalidation. Be careful to only match on the content you need to remove - revalidation is an expensive operation for many origins, and a simple ``/.*`` can cause an overload condition of the origin.
	- Enter the **Time To Live** - this is how long the revalidation rule will be active for. It usually makes sense to make this the same as the ``Cache-Control`` header from the origin which sets the object time to live in cache (by ``max-age`` or ``Expires``). Entering a longer TTL here will make the caches do unnecessary work.
	- Enter the **Start Time** - this is the start time when the revalidation rule will be made active. It is pre-populated with the current time, leave as is to schedule ASAP.

#. Click the **Submit** button.


Manage DNSSEC Keys
==================
In order to support `DNSSEC <https://en.wikipedia.org/wiki/Domain_Name_System_Security_Extensions>`_ in Traffic Router, Traffic Ops provides some actions for managing DNSSEC keys for a CDN and associated :term:`Delivery Service`\ s.  DNSSEC Keys consist of a Key Signing Keys (KSK) which are used to sign other DNSKEY records as well as Zone Signing Keys (ZSK) which are used to sign other records.  DNSSEC Keys are stored in `Traffic Vault <../../overview/traffic_vault.html>`_ and should only be accessible to Traffic Ops.  Other applications needing access to this data, such as Traffic Router, must use the Traffic Ops `DNSSEC APIs <../../development/traffic_ops_api/v12/cdn.html#dnssec-keys>`_ to retrieve this information.

To Manage DNSSEC Keys:
1. Click **Tools -> Manage DNSSEC Keys**
2. Choose a CDN and click **Manage DNSSEC Keys**

	- If keys have not yet been generated for a CDN, this screen will be mostly blank with just the **CDN** and **DNSSEC Active?** fields being populated.
	- If keys have been generated for the CDN, the Manage DNSSEC Keys screen will show the TTL and Top Level Domain (TLD) :abbr:`KSK (Key Signing Key)` Expiration for the CDN as well as DS Record information which will need to be added to the parent zone of the TLD in order for DNSSEC to work.

The Manage DNSSEC Keys screen also allows a user to perform the following actions:

Activate/Deactivate DNSSEC for a CDN
------------------------------------
Fairly straight forward, this button set the **dnssec.enabled** param to either **true** or **false** on the Traffic Router profile for the CDN.  The Activate/Deactivate option is only available if DNSSEC keys exist for CDN.  In order to active DNSSEC for a CDN a user must first generate keys and then click the **Active DNSSEC** button.

Generate Keys
-------------
Generate Keys will generate DNSSEC keys for the CDN TLD as well as for each :term:`Delivery Service` in the CDN.  It is important to note that this button will create a new :abbr:`KSK (Key Signing Key)` for the TLD and, therefore, a new DS Record.  Any time a new DS Record is created, it will need to be added to the parent zone of the TLD in order for DNSSEC to work properly.  When a user clicks the **Generate Keys** button, they will be presented with a screen with the following fields:

- **CDN:** This is not editable and displays the CDN for which keys will be generated
- **ZSK Expiration (Days):**  Sets how long (in days) the Zone Signing Key will be valid for the CDN and associated :term:`Delivery Service`\ s. The default is 30 days.
- **KSK Expiration (Days):**  Sets how long (in days) the Key Signing Key will be valid for the CDN and associated :term:`Delivery Service`\ s. The default is 365 days.
- **Effective Date (GMT):** The time from which the new keys will be active.  Traffic Router will use this value to determine when to start signing with the new keys and stop signing with the old keys.

Once these fields have been correctly entered, a user can click Generate Keys.  The user will be presented with a confirmation screen to help them understand the impact of generating the keys.  If a user confirms, the keys will be generated and stored in Traffic Vault.

Regenerate KSK
--------------
Regenerate :abbr:`KSK (Key Signing Key)` will create a new Key Signing Key for the CDN TLD. A new DS Record will also be generated and need to be put into the parent zone in order for DNSSEC to work correctly. The **Regenerate KSK** button is only available if keys have already been generated for a CDN.  The intent of the button is to provide a mechanism for generating a new :abbr:`KSK (Key Signing Key)` when a previous one expires or if necessary for other reasons such as a security breach.  When a user goes to generate a new :abbr:`KSK (Key Signing Key)` they are presented with a screen with the following options:

:CDN: This is not editable and displays the CDN for which keys will be generated
:KSK Expiration (Days): Sets how long (in days) the Key Signing Key will be valid for the CDN and associated :term:`Delivery Service`\ s. The default is 365 days.
:Effective Date (GMT): The time from which the new :abbr:`KSK (Key Signing Key)` and DS Record will be active. Since generating a new :abbr:`KSK (Key Signing Key)` will generate a new DS Record that needs to be added to the parent zone, it is very important to make sure that an effective date is chosen that allows for time to get the DS Record into the parent zone. Failure to get the new DS Record into the parent zone in time could result in DNSSEC errors when Traffic Router tries to sign responses.

Once these fields have been correctly entered, a user can click Generate KSK. The user will be presented with a confirmation screen to help them understand the impact of generating the KSK.  If a user confirms, the :abbr:`KSK (Key Signing Key)` will be generated and stored in Traffic Vault.

Additionally, Traffic Ops also performs some systematic management of :abbr:`DNSSEC (DNS Security Extensions)` keys. This management is necessary to help keep keys in sync for :term:`Delivery Service`\ s in a CDN as well as to make sure keys do not expire without human intervention.

Generation of keys for new Delivery Services
--------------------------------------------
If a new :term:`Delivery Service` is created and added to a CDN that has :abbr:`DNSSEC (DNS Security Extensions)` enabled, Traffic Ops will create :abbr:`DNSSEC (DNS Security Extensions)` keys for the :term:`Delivery Service` and store them in Traffic Vault.

Regeneration of expiring keys for a Delivery Service
----------------------------------------------------
Traffic Ops has a process, controlled by :manpage:`cron(8)`, to check for expired or expiring keys and re-generate them. The process runs at 5 minute intervals to check and see if keys are expired or close to expiring (withing 10 minutes by default). If keys are expired for a :term:`Delivery Service`, Traffic Ops will regenerate new keys and store them in Traffic Vault. This process is the same for the CDN :abbr:`TLD (Top-Level Domain)` :abbr:`ZSK (Zone Signing Key)`, however Traffic Ops will not re-generate the CDN :abbr:`TLD (Top-Level Domain)` :abbr:`KSK (Key Signing Key)` systematically. The reason is that when a :abbr:`KSK (Key Signing Key)` is regenerated for the CDN :abbr:`TLD (Top-Level Domain)` then a new DS Record will also be created. The new DS Record needs to be added to the parent zone before Traffic Router attempts to sign with the new :abbr:`KSK (Key Signing Key)` in order for :abbr:`DNSSEC (DNS Security Extensions)` to work correctly. Therefore, management of the :abbr:`KSK (Key Signing Key)` needs to be a manual process.
