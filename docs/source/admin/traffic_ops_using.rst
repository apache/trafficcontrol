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

Using Traffic Ops
%%%%%%%%%%%%%%%%%


The Traffic Ops Menu
====================

.. image:: 12m.png

The following tabs are available in the menu at the top of the Traffic Ops user interface.

.. index:: 
  Health Tab

* **Health**

  Information on the health of the system. Hover over this tab to get to the following options:

  +---------------+------------------------------------------------------------------------------------------------------------------------------------+
  |     Option    |                                                            Description                                                             |
  +===============+====================================================================================================================================+
  | Table View    | A real time view into the main performance indicators of the CDNs managed by Traffic Control.                                      |
  |               | This view is sourced directly by the Traffic Monitor data and is updated every 10 seconds.                                         |
  |               | This is the default screen of Traffic Ops.                                                                                         |
  |               | See :ref:`rl-health-table` for details.                                                                                            |
  +---------------+------------------------------------------------------------------------------------------------------------------------------------+
  | Graph View    | A real graphical time view into the main performance indicators of the CDNs managed by Traffic Control.                            |
  |               | This view is sourced by the Traffic Monitor data and is updated every 10 seconds.                                                  |
  |               | On loading, this screen will show a history of 24 hours of data from Traffic Stats                                                 |
  |               | See :ref:`rl-health-graph` for details.                                                                                            |
  +---------------+------------------------------------------------------------------------------------------------------------------------------------+
  | Server Checks | A table showing the results of the periodic check extension scripts that are run. See :ref:`rl-server-checks`                      |
  +---------------+------------------------------------------------------------------------------------------------------------------------------------+
  | Daily Summary | A graph displaying the daily peaks of bandwidth, overall bytes served per day, and overall bytes served since initial installation |
  |               | per CDN.                                                                                                                           |
  +---------------+------------------------------------------------------------------------------------------------------------------------------------+

* **Delivery Services**

  The main Delivery Service table. This is where you Create/Read/Update/Delete Delivery Services of all types. There are currently no sub menus for this tab.

* **Servers**

  The main Servers table. This is where you Create/Read/Update/Delete servers of all types.  Click the main tab to get to the main table, and hover over to get these sub options:

  +-------------------+------------------------------------------------------------------------------------------+
  |       Option      |                                       Description                                        |
  +===================+==========================================================================================+
  | Upload Server CSV | Bulk add of servers from a csv file. See :ref:`rl-bulkserver`                            |
  +-------------------+------------------------------------------------------------------------------------------+

* **Parameters**

  Parameters and Profiles can be edited here. Hover over the tab to get the following options:

  +---------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------+
  |        Option       |                                                                             Description                                                                             |
  +=====================+=====================================================================================================================================================================+
  | Global Profile      | The table of global parameters. See :ref:`rl-param-prof`. This is where you Create/Read/Update/Delete parameters in the Global profile                              |
  +---------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------+
  | All Cache Groups    | TBD JvD                                                                                                                                                             |
  +---------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------+
  | All Profiles        | The table of all parameters - this may be slow to pull up, as there can be thousands of parameters.                                                                 |
  +---------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------+
  | Select Profile      | Select the parameter by Profile first, then get a table of just the parameters for that profile.                                                                    |
  +---------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------+
  | Orphaned Parameters | A table of parameters that are not associated to any profile of cache group. These parameters either should be deleted or associated with a profile of cache group. |
  +---------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------+

* **Tools**

  Tools for working with Traffic Ops and it's servers. Hover over this tab to get the following options:

  +----------------------+-----------------------------------------------------------------------------------------------------------------------------------+
  |        Option        |                                                            Description                                                            |
  +======================+===================================================================================================================================+
  | Generate ISO         | Generate a bootable image for any of the servers in the Servers table (or any server for that matter). See :ref:`rl-generate-iso` |
  +----------------------+-----------------------------------------------------------------------------------------------------------------------------------+
  | Queue Updates        | Send Updates to the caches. See :ref:`rl-queue-updates`                                                                           |
  +----------------------+-----------------------------------------------------------------------------------------------------------------------------------+
  | DB Dump              | Backup the Database to a .sql file.                                                                                               |
  +----------------------+-----------------------------------------------------------------------------------------------------------------------------------+
  | Snapshot CRConfig    | Send updates to the Traffic Monitor / Traffic Router servers.  See :ref:`rl-queue-updates`                                        |
  +----------------------+-----------------------------------------------------------------------------------------------------------------------------------+
  | Invalidate Content   | Invalidate or purge content from the CDN. See :ref:`rl-purge`                                                                     |
  +----------------------+-----------------------------------------------------------------------------------------------------------------------------------+
  | Generate DNSSEC keys | Neuman?                                                                                                                           |
  +----------------------+-----------------------------------------------------------------------------------------------------------------------------------+

* **Misc**

  Miscellaneous editing options. Hover over this tab to get the following options:

  +--------------------+-------------------------------------------------------------------------------------------+
  |       Option       |                                        Description                                        |
  +====================+===========================================================================================+
  | Cache Groups       | Create/Read/Update/Delete cache groups                                                    |
  +--------------------+-------------------------------------------------------------------------------------------+
  | Users              | Create/Read/Update/Delete users                                                           |
  +--------------------+-------------------------------------------------------------------------------------------+
  | Profiles           | Create/Read/Update/Delete profiles. See :ref:`rl-working-with-profiles`                   |
  +--------------------+-------------------------------------------------------------------------------------------+
  | Networks(ASNs)     | Create/Read/Update/Delete Autonomous System Numbers See :ref:`rl-asn-czf`                 |
  +--------------------+-------------------------------------------------------------------------------------------+
  | Hardware           | Get detailed hardware information (note: this should be moved to a Traffic Ops Extension) |
  +--------------------+-------------------------------------------------------------------------------------------+
  | Data Types         | Create/Read/Update/Delete data types                                                      |
  +--------------------+-------------------------------------------------------------------------------------------+
  | Divisions          | Create/Read/Update/Delete divisions                                                       |
  +--------------------+-------------------------------------------------------------------------------------------+
  | Regions            | Create/Read/Update/Delete regions                                                         |
  +--------------------+-------------------------------------------------------------------------------------------+
  | Physical Locations | Create/Read/Update/Delete locations                                                       |
  +--------------------+-------------------------------------------------------------------------------------------+

.. index::
  Change Log

* **ChangeLog**

  The Changelog table displays the changes that are being made to the Traffic Ops database through the Traffic Ops user interface. This tab will show the number of changes since you last visited this tab in (brackets) since the last time you visited this tab. There are currently no sub menus for this tab.

* **Help**

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

Health
======

.. _rl-health-table:

The Health Table
++++++++++++++++
The Health table is the default landing screen for Traffic Ops, it displays the status of the EDGE caches in a table form, sorted by Mbps Out. The columns in this table are:

.. |checkmark| image:: ../../../traffic_ops/app/public/images/good.png 

.. |X| image:: ../../../traffic_ops/app/public/images/bad.png
  
.. |clock| image:: ../../../traffic_ops/app/public/images/clock-black.png

* **Profile**: the Profile of this server or ALL, meaning this row shows data for multiple servers, and the row shows the sum of all values.
* **Host Name**: the host name of the server or ALL, meaning this row shows data for multiple servers, and the row shows the sum of all values.
* **Edge Cache Group**: the edge cache group short name or ALL, meaning this row shows data for multiple servers, and the row shows the sum of all values.
* **Healthy**: indicates if this cache is healthy according to the Health Protocol. A row with ALL in any of the columns will always show a |checkmark|, this column is valid only for individual EDGE caches. 
* **Admin**: shows the administrative status of the server. 
* **Connections**: the number of connections this cache (or group of caches) has open (``ats.proxy.process.http.current_client_connections`` from ATS).
* **Mbps Out**: the bandwidth being served out if this cache (or group of caches)

Since the top line has ALL, ALL, ALL, it shows the total connections and bandwidth for all caches managed by this instance of Traffic Ops.

.. _rl-health-graph:

Graph View
++++++++++
More Blah Blah


Server Checks
+++++++++++++

Daily Summary
+++++++++++++

.. _rl-server:

Server
======

Server Types
++++++++++++

.. _rl-bulkserver:

Bulk Upload Server
++++++++++++++++++

.. _rl-generate-iso:

Generate Server ISO
+++++++++++++++++++


Delivery Service
================
The fields in the Delivery Service view are:

.. Sorry for the width of this table, don't know how to make the bullet lists work otherwise. Just set your monitor to 2560*1600, and put on your glasses.

+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|                       Name                       |                                                                                                     Description                                                                                                     |
+==================================================+=====================================================================================================================================================================================================================+
| XML ID                                           | A unique string that identifies this delivery service.                                                                                                                                                              |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Content Routing Type                             | The type of content routing this delivery service will use. See :ref:`rl-ds-types`.                                                                                                                                 |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Protocol                                         | The protocol to serve this delivery service to the clients with:                                                                                                                                                    |
|                                                  |                                                                                                                                                                                                                     |
|                                                  | -  http                                                                                                                                                                                                             |
|                                                  | -  https                                                                                                                                                                                                            |
|                                                  | -  both http and https                                                                                                                                                                                              |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| DSCP Tag                                         | The DSCP value to mark IP packets to the client with.                                                                                                                                                               |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Signed URLs                                      | Use Signed URLs? See :ref:`rl-signed-urls`.                                                                                                                                                                         |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Query String Handling                            | How to treat query strings:                                                                                                                                                                                         |
|                                                  |                                                                                                                                                                                                                     |
|                                                  | - 0 use in cache key and hand up to origin -this means each unique query string Is treated as a unique URL.                                                                                                         |
|                                                  | - 1 Do not use in cache key, but pass up to origin - this means a 2 URLs that are the same except for the query string will match, and cache HIT, while the origin still sees original query string in the request. |
|                                                  | - 2 Drop at edge - this means a 2 URLs that are the same except for  the query string will match, and cache HIT, while the origin will not see original query string in the request.                                |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Geo Limit?                                       | Some services are intended to be limited by geography. The possible settings are are:                                                                                                                               |
|                                                  |                                                                                                                                                                                                                     |
|                                                  | - None - Do not limit by geography.                                                                                                                                                                                 |
|                                                  | - CZF only - If the requesting IP is not in the Coverage Zone File, do not serve the request.                                                                                                                       |
|                                                  | - CZF + US - If the requesting IP is not in the Coverage Zone File or not in the United States, do not serve the request.                                                                                           |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Bypass FQDN                                      | (for HTTP routed delivery services only) This is the FQDN Traffic Router will redirect to (with the same path) when the max Bps or Max Tps for this deliveryservice are exceeded.                                   |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Bypass Ipv4                                      | (For DNS routed delivery services only) This is the address to respond to A requests with when the the max Bps or Max Tps for this delivery service are exceeded.                                                   |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Bypass IPv6                                      | (For DNS routed delivery services only) This is the address to respond to AAAA requests with when the the max Bps or Max Tps for this delivery service are exceeded.                                                |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| IPv6 Routing Enabled?                            | When set to yes, the Traffic Router will respond to AAAA DNS requests for the tr. and edge. names of this delivery service. Otherwise, only A records will be served.                                               |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Background fetch Enabled?                        | Experimental. This enables the background_fetch plugin to fetch the whole file on seeing a range request.                                                                                                           |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Delivery Service DNS TTL                         | The Time To Live on the DNS record for the Traffic Router A and AAAA records (``tr.<deliveryservice>.<cdn-domain>``) for a HTTP delivery service *or* for the A and                                                 |
|                                                  | AAAAA records of the edge name (``edge.<deliveryservice>.<cdn-domain>``).                                                                                                                                           |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Origin Server Base URL                           | The Origin Server's base URL. This includes the protocol (http or https). Example: ``http://movies.origin.com``                                                                                                     |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| CCR profile                                      | The Traffic Router  profile for this delivery service. See :ref:`rl-ccr-profile`.                                                                                                                                   |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Maximum Bits per Second allowed globally         | The maximum bits per second this delivery service can serve across all EDGE caches before traffic will be diverted to the bypass destination. For a DNS delivery service, the Bypass Ipv4 or Ipv6  will be used     |
|                                                  | (depending on whether this was a A or AAAA request), and for HTTP delivery services the Bypass FQDN will be used.                                                                                                   |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Maximum Transactions per Second allowed globally | The maximum transactions per se this delivery service can serve across all EDGE caches before traffic will be diverted to the bypass destination. For a DNS delivery service, the Bypass Ipv4 or Ipv6  will be used |
|                                                  | (depending on whether this was a A or AAAA request), and for HTTP delivery services the Bypass FQDN will be used.                                                                                                   |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Geo Miss Default Latitude                        | Default Latitude for this delivery service. When client localization fails for bot Coverage Zone and Geo Lookup, this the client will be routed as if it was at this lat.                                           |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Geo Miss Default Longitude                       | Default Longitude for this delivery service. When client localization fails for bot Coverage Zone and Geo Lookup, this the client will be routed as if it was at this long.                                         |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Request Header Rewrite Rules                     | Header Rewrite rules for this delivery service. See :ref:`rl-header-rewrite`.                                                                                                                                       |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Long Description                                 | Long description for this delivery service. TO be consumed from the APIs by downstream tools (Portal).                                                                                                              |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Customer                                         | Customer description for this delivery service. TO be consumed from the APIs by downstream tools (Portal).                                                                                                          |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Service                                          | Service description for this delivery service. TO be consumed from the APIs by downstream tools (Portal).                                                                                                           |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Info URL                                         | Info URL  for this delivery service. TO be consumed from the APIs by downstream tools (Portal).                                                                                                                     |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Check Path                                       | A path (ex: /crossdomain.xml) to verify the connection to the origin server with. This can be used by Check Extension scripts to do periodic health checks against the delivery service.                            |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Origin Shield (Pipe Delimited String)            | Experimental. Origin Shield string. See :ref:`rl-org-shield`                                                                                                                                                        |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Active                                           | When this is set to no Traffic Router will not serve DNS or HTTP responses for this delivery service.                                                                                                               |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Last Updated                                     | (Read Only) The last time this delivery service was updated.                                                                                                                                                        |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Number of edges assigned                         | (Read Only - change by clicking the **Server Assignments** button at the bottom) The number of EDGE caches assigned to this delivery service. See :ref:`rl-assign-edges`.                                           |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Number of static DNS entries                     | (Read Only - change by clicking the **Static DNS** button at the bottom) The number of static DNS entries for this delivery service. See :ref:`rl-static-dns`.                                                      |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Example delivery URL                             | (Read Only) An example of how the delivery URL may start. This could be multiple rows if multiple HOST_REGEXP entries have been entered.                                                                            |
+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+


.. _rl-ds-types:

Delivery Service Types
++++++++++++++++++++++

+-----------------+----------------------------------------------------------------------------------------------------------------------------------------------------------------+
|       Name      |                                                                          Description                                                                           |
+=================+================================================================================================================================================================+
| HTTP            | HTTP Content Routing  - The Traffic Router DNS auth server returns its own IP address on DNS queries, and the client gets redirected to a specific cache       |
|                 | in the nearest cache group using HTTP 302.  Use this for long sessions like HLS/HDS/Smooth live streaming, where a longer setup time is not a.                 |
|                 | problem.                                                                                                                                                       |
+-----------------+----------------------------------------------------------------------------------------------------------------------------------------------------------------+
| DNS             | DNS Content Routing - The Traffic Router DNS auth server returns an edge cache IP address to the client right away. The client will find the cache quickly     |
|                 | but the Traffic Router can not route to a cache that already has this content in the cache group. Use this for smaller objects like web page images / objects. |
+-----------------+----------------------------------------------------------------------------------------------------------------------------------------------------------------+
| HTTP_NO_CACHE   | HTTP Content Routing, but the caches will not actually cache the content, they act as just proxies. The MID tier is bypassed.                                  |
+-----------------+----------------------------------------------------------------------------------------------------------------------------------------------------------------+
| HTTP_LIVE       | HTTP Content routing, but where for "standard" HTTP content routing the objects are stored on disk, for this delivery service type the objects are stored      |
|                 | on the RAM disks. Use this for linear TV. The MID tier is bypassed for this type.                                                                              |
+-----------------+----------------------------------------------------------------------------------------------------------------------------------------------------------------+
| HTTP_LIVE_NATNL | HTTP Content routing, same as HTTP_LIVE, but the MID tier is NOT bypassed.                                                                                     |
+-----------------+----------------------------------------------------------------------------------------------------------------------------------------------------------------+
| DNS_LIVE_NATNL  | DNS Content routing, ut where for "standard" DNS content routing the objects are stored on disk, for this delivery service type the objects are stored         |
|                 | on the RAM disks. Use this for linear TV. The MID tier is NOT bypassed for this type.                                                                          |
+-----------------+----------------------------------------------------------------------------------------------------------------------------------------------------------------+
| DNS_LIVE        | DNS Content routing, same as DNS_LIVE_NATIONAL, but the MID tier is bypassed.                                                                                  |
+-----------------+----------------------------------------------------------------------------------------------------------------------------------------------------------------+


.. _rl-queue-updates:

Queue Updates and Snapshot CRConfig
+++++++++++++++++++++++++++++++++++
When changing delivery services special care has to be taken so that Traffic Router will not send traffic to caches for delivery services that the cache doesn't know about yet. In general, when adding delivery services, or adding servers to a delivery service, it is best to update the caches before updating Traffic Router and Traffic Monitor. When deleting delivery services, or deleting server assignments to delivery services, it is best to update Traffic Router and Traffic Monitor first and then the caches. Updating the cache configuration is done through the *Queue Updates* menu, and updating Traffic Monitor and  Traffic Router config is done through the *Snapshot CRConfig* menu.

Queue Updates
^^^^^^^^^^^^^
Every 15 minutes the caches will run a *syncds* to get all changes needed from Traffic Ops. The files that will be updated by the syncds job are: 

- records.config
- remap.config
- parent.config
- cache.config
- hosting.config
- url\_sig\_(.*)\.config
- hdr\_rw\_(.*)\.config
- regex_revalidate.config
- ip_allow.config

A cache will only get updated when the update flag is set for it. To set the update flag, use the *Queue Updates* menu - here you can schedule updates for a whole CDN or a cache group:

  #. Click **Tools > Queue Updates**.
  #. Select the CDN to queueu uodates for, or All.
  #. Select the cache group to queue updates for, or All
  #. Click the **Queue Updates** button.
  #. When the Queue Updates for this Server? (all) window opens, click **OK**.

To schedule updates for just one cache, use the "Server Checks" page, and click the |checkmark| in the *UPD* column. The UPD column of Server Checks page will change show a |clock| when updates are pending for that cache. 


Snapshot CRConfig
^^^^^^^^^^^^^^^^^
Every 60 seconds Traffic Monitor will check with Traffic Ops to see if a new CRConfig snapshot was made. If there is a new CRCOnfig, it will apply it to both Traffic Monitor and Traffic Router. See :ref:`rl-crconfig` for more information on the CRConfig file. To create a new snapshot, use the *Tools > Snapshot CRConfig* menu:

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
  #. The Successfully wrote CRConfig.json! window opens, click **OK**.


.. _rl-purge:

Invalidate Content
==================
Invalidating content on the CDN is sometimes necessary when the origin was mis configured and something is cached in the CDN caches that needs to be removed. Given the size of a typical Traffic Control CDN and the amount of content that can be cached in it, removing the content from all the caches may take a long time. To speed up content invalidation, Traffic Ops will not try to remove the content from the caches, but it makes the content in accessible using the *regex_revalidate* ATS plugin. This forces a *revalidation* of the content, rather than a new get.

.. Note:: This method forces a HTTP *revalidation* of the content, and not a new *GET* - the origin needs to support revalidation according to the HTTP/1.1 specification, and send a ``200 OK`` or ``304 Not Modified`` as applicable. 

To invalidate content:

  1. Click **Tools > Invalidate Content**
  2. Fill out the form fields: 

    - Select the **Delivery Service**
    - Enter the **Path Regex** - this should be a `PCRE <http://www.pcre.org/>`_ compatible regular expression for the path to match for forcing the revalidation. Be careful to only match on the content you need to remove - revalidation is an expensive operation for many origins, and a simple ``/.*`` can cause an overload condition of the origin. 
    - Enter the **Time To Live** - this is how long the revalidation force will be active for. It usually makes sense to make this the same as the ``Cache-Control`` header from the origin sets the object time to live in cache (by ``max-age`` or ``Expires``). Entering a longer TTL here will make the caches do unnecessary work. 
    - Enter the **Start Time** - this is the start time when the force revalidation will be made active. Is pre populated with the current time, leave as is to schedule ASAP. 

  3. Click the **Submit** button.


