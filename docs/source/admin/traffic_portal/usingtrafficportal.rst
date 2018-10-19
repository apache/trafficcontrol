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

.. _usingtrafficportal:

**********************
Traffic Portal - Using
**********************
Traffic Portal is the official Traffic Control UI. Traffic Portal typically runs on a different machine than Traffic Ops, and works by using the Traffic Ops API. The following high-level items are available in the Traffic Portal menu.

.. figure:: ./images/tp_menu.png
	:width: 55%
	:align: center
	:alt: The Traffic Portal Landing Page

	Traffic Portal Start Page

Dashboard
=========

The Dashboard is the default landing page for Traffic Portal. It provides a real-time view into the main performance indicators of the CDNs managed by Traffic Control. It also displays various statistics about the overall health of your CDN.

Current Bandwidth
	The current bandwidth of all of your CDNs.

Current Connections
	The current number of connections to all of your CDNs.

Healthy Caches
	Displays the number of healthy caches across all CDNs.  Click the link to view the healthy caches on the cache stats page.

Unhealthy Caches
	Displays the number of unhealthy caches across all CDNs.  Click the link to view the unhealthy caches on the cache stats page.

Online Caches
	Displays the number of cache servers with ONLINE status. Traffic Monitor will not monitor the state of ONLINE servers. For more information, see :ref:`health-proto`.

Reported Caches
	Displays the number of cache servers with REPORTED status. For more information, see :ref:`health-proto`.

Offline Caches
	Displays the number of cache servers with OFFLINE status. For more information, see :ref:`health-proto`.

Admin Down Caches
	Displays the number of caches with ADMIN_DOWN status. For more information, see :ref:`health-proto`.

Each component of this view is updated on the intervals defined in the ``tp.domain.com/traffic_portal_properties.json`` configuration file.

CDNs
====

A table of CDNs with the following columns:

:Name:           The name of the CDN
:Domain:         The CDN's Top-Level Domain (TLD)
:DNSSEC Enabled: 'true' if :ref:`tr-dnssec` is enabled on this CDN, 'false' otherwise.

CDN management includes the ability to (where applicable):

- create a new CDN
- update an existing CDN
- delete an existing CDN
- queue/clear updates on all servers in a CDN
- diff CDN snapshots
- create a CDN snapshot
- manage a CDN's DNSSEC keys
- manage a CDN's federations
- view Delivery Services of a CDN
- view CDN profiles
- view servers within a CDN

Monitor
=======

The 'Monitor' section of Traffic Portal is used to display statistics regarding the various cache servers within all CDNs visible to the user. It retrieves this information through the Traffic Ops API from Traffic Monitor instances.

.. figure:: ./images/tp_menu_monitor.png
	:align: center
	:alt: The Traffic Portal 'Monitor' Menu

	The 'Monitor' Menu


Cache Checks
------------
A real-time view into the status of each cache.

The cache checks page is intended to give an overview of the caches managed by Traffic Control as well as their status.

:Hostname: Cache host name
:Profile:  The name of the profile applied to the cache
:Status:   The status of the cache (one of: ONLINE, REPORTED, ADMIN_DOWN, OFFLINE)
:UPD:      Configuration updates pending for an EDGE or MID
:RVL:      Content invalidation requests are pending for this server and/or its parent(s)
:ILO:      Ping the iLO interface for EDGE or MID servers
:10G:      Ping the IPv4 address of the EDGE or MID servers
:FQDN:     DNS check that matches what the DNS servers responds with compared to what Traffic Ops has
:DSCP:     Checks the DSCP value of packets from the EDGE server to the Traffic Ops server
:10G6:     Ping the IPv6 address of the EDGE or MID servers
:MTU:      Ping the EDGE or MID using the configured MTU from Traffic Ops
:RTR:      Content Router checks. Checks the health of the Content Routers. Checks the health of the caches using the Content Routers
:CHR:      Cache Hit Ratio percent
:CDU:      Total Cache Disk Usage percent
:ORT:      Operational Readiness Test - uses the ORT script on the EDGE and MID servers to determine if the configuration in Traffic Ops matches the configuration on the EDGE or MID. The user that this script runs as must have an SSH key on the EDGE servers.


Cache Stats
-----------
A table showing the results of the periodic check extension scripts that are run. These can be grouped by Cache Group and/or Profile.

:Profile:     Name of the profile applied to the Edge-tier or Mid-tier cache server
:Host:        'ALL' for entries grouped by cache group, or the hostname of a particular cache server
:Cache Group: Name of the Cache Group to which this server belongs, or the name of the Cache Group that is grouped for entries grouped by Cache Group
:Healthy:     True/False as determined by Traffic Monitor (See :ref:`health-proto`)
:Status:      Status of the cache or Cache Group
:Connections: Number of connections to this cache server or Cache Group
:MbpsOut:     Data flow outward (toward client) in Megabits per second

Services
========
'Services' groups the functionality to modify Delivery Services - for those users with the necessary permissions - or make Requests for such changes - for uses without necessary permissions.

.. figure:: images/tp_table_ds_requests.png
	:align: center
	:alt: An example table of Delivery Service Requests

	Table of Delivery Service Requests

Delivery Services
-----------------
This page contains a table displaying all Delivery Services visible to the user. Each entry in this table has the following fields:

:Key (XML ID): A unique string that identifies this Delivery Service
:Tenant: The tenant to which the Delivery Service is assigned
:Origin: The Origin Server's base URL. This includes the protocol (HTTP or HTTPS). Example: ``http://movies.origin.com``
:Active: When this is set to 'false', Traffic Router will not serve DNS or HTTP responses for this Delivery Service
:Type: The type of content routing this Delivery Service will use

	.. seealso:: :ref:`ds-types`

:Protocol: The protocol which which this Delivery Service serves clients. Its value is one of:

	HTTP
		Only insecure requests will be serviced
	HTTPS
		Only secure requests will be serviced
	HTTP and HTTPS
		Both secure and insecure requests will be serviced
	HTTP to HTTPS
		Insecure requests will be redirected to secure locations and secure requests are serviced normally

:CDN: The CDN to which the Delivery Service belongs
:IPv6 Enabled: When set to 'true', the Traffic Router will respond to AAAA DNS requests for the routed name of this Delivery Service, Otherwise, only A records will be served
:DSCP: The Differentiated Services Code Point (DSCP) value with which to mark IP packets sent to the client
:Signing Algorithm: See :ref:`signed-urls`
:Query String Handling: Describes how the Delivery Service treats query strings. It has one of the following possible values:

	USE
		The query string will be used in the Apache Traffic Server (ATS) 'cache key' and is passed in requests to the origin (each unique query string is treated as a unique URL)
	IGNORE
		The query string will *not* be used in the ATS 'cache key', but *will* be passed in requests to the origin
	DROP
		The query string is stripped from the request URL at the Edge-tier cache, and so is not used in the ATS 'cache key', and is not passed in requests to the origin

	.. seealso:: :ref:`qstring-handling`

:Last Updated: Timestamp when the Delivery Service was last updated.                                                                 |

Delivery Service management includes the ability to (where applicable):

- create a new Delivery Service
- clone an existing Delivery Service
- update an existing Delivery Service
- delete an existing Delivery Service
- compare Delivery Services
- manage Delivery Service SSL keys
- manage Delivery Service URL signature keys
- manage Delivery Service URI signing keys
- view and assign Delivery Service servers
- create, update and delete Delivery Service regular expressions
- view and create Delivery Service invalidate content jobs
- manage steering targets

Delivery Service Requests
-------------------------
If enabled in the ``tp.domain.com/traffic_portal_properties.json``, all Delivery Service changes (create, update and delete) are captured as a Delivery Service Request and must be reviewed before fulfillment/deployment.

:Delivery Service: A unique string that identifies the Delivery Service that with which the request is associated. This unique string is also known (and ofter referred to within documentation and source code) as a 'Delivery Service key' or 'XML ID'.                                                  |
:Type: The type of Delivery Service Request: 'create', 'update', or 'delete' according to what was requested
:Status: The status of the Delivery Service Request. Has the following possible values:

	draft
		The Delivery Service Request is *not* ready for review and fulfillment
	submitted
		The Delivery Service Request is ready for review and fulfillment
	rejected
		The Delivery Service Request has been rejected and cannot be modified
	pending
		The Delivery Service Request has been fulfilled but the changes have yet to be deployed
	complete
		The Delivery Service Request has been fulfilled and the changes have been deployed

:Author: The user responsible for creating the Delivery Service Request
:Assignee: The user responsible for fulfilling the Delivery Service Request. Currently, the operations role or above is required to assign Delivery Service Requests
:Last Edited By: The last user to edit the Delivery Service Request
:Created: Relative time indicating when the Delivery Service Request was created
:Actions: Actions that can be performed on a Delivery Service Request. The following actions are provided:

	fulfill
		Implement the changes captured in the Delivery Service Request
	reject
		Reject the changes captured in the Delivery Service Request
	delete
		Delete the Delivery Service Request

Delivery service request management includes the ability to (where applicable):

- create a new delivery service request
- update an existing delivery service request
- delete an existing delivery service request
- update the status of a delivery service request
- assign a delivery service request
- reject a delivery service request
- fulfill a delivery service request
- complete a delivery service request

.. seealso:: :ref:`ds_requests`

Configure
=========
Interfaces for managing the various components of Traffic Control and how they interact are grouped under 'Configure'.

.. figure:: ./images/tp_menu_configure.png
	:align: center
	:alt: The 'Configure' Menu

	The 'Configure' Menu

Servers
-------
A table of all servers (of all kinds) across all Delivery Services visible to the user, with functionality to create, update, and delete them. It has the following columns:

:UPD: 'true' when updates to the server's configuration are pending, 'false' otherwise
:Host: The hostname of the server
:Domain: The server's domain. (The FQDN of the server is given by 'Host.Domain')
:IP: The server's IPv4 address
:IPv6: The server's IPv6 address
:Status: The server's status (see :ref:`health-proto`)
:Type: The type of server e.g. EDGE for an Edge-tier cache
:Profile: The name of the server's profile
:CDN: The name of the CDN to which this server is assigned (if any)
:Cache Group: The name of the Cache Group to which this server belongs
:ILO: If not empty, this is the IPv4 address of the server's Integrated Lights-Out (ILO) interface

	.. seealso:: `Hewlett Packard ILO Wikipedia Page <https://en.wikipedia.org/wiki/HP_Integrated_Lights-Out>`_

Server management includes the ability to (where applicable):

- create a new server
- update an existing server
- delete an existing server
- queue/clear updates on a server
- update server status
- view server delivery services
- view server configuration files
- clone delivery service assignments
- assign delivery services to server


Profiles
--------
A table of all profiles. From here you can see parameters, servers and Delivery Services assigned to each profile, as well as the ability to create, update, delete, import and export profiles.
Each entry in the table has these fields:

:Name:             The name of the profile
:Type:             The type of this profile, which indicates the kinds of objects to which the profile may be assigned
:Routing Disabled: For profiles applied to cache servers (Edge-tier or Mid-tier) this indicates that Traffic Router will refuse to provide routes to these machines
:Description:      A user-defined description of the profile, typically indicating its purpose
:CDN:              The CDN to which this profile is restricted. To use the same profile across multiple CDNs, clone the profile and change the clone's CDN field.

Profile management includes the ability to (where applicable):

- create a new profile
- update an existing profile
- delete an existing profile
- clone a profile
- export a profile
- view profile parameters
- view profile delivery services
- view profile servers

.. seealso:: :ref:`working-with-profiles`


Parameters
----------
Allows for the creation, update, and deletion of parameters, as well as modification of their assignment to servers and Delivery Services.
This page displays a table of parameters with the following columns:

:Name:        The name of the parameter
:Config File: The configuration file where this parameter is stored, possibly the special value ``location``, indicating that this parameter actually names the location of a configuration file rather than its contents, or ``package`` to indicate that this parameter specifies a package to be installed rather than anything to do with configuration files
:Value:       The value of the parameter. The meaning of this depends on the value of 'Config File'
:Secure:      When this is 'true', a user requesting to see this parameter will see the value ``********`` instead of its actual value if the user's permission level isn't 'admin'
:Profiles:    The number of profiles currently using this parameter

Parameter management includes the ability to (where applicable):

- create a new parameter
- update an existing parameter
- delete an existing parameter
- view parameter profiles


.. _tp-configure-types:

Types
-----
'Types' groups Delivery Services, servers and Cache Groups for various purposes. Each entry in the table shown on this page has the following fields:

:Name:         The name of the Type
:Use In Table: States the use of this Type, e.g. ``server`` indicates this is a Type assigned to servers
:Description:  A short, usually user-defined, description of the Type

Type management includes the ability to (where applicable):

- create a new type
- update an existing type
- delete an existing type
- view delivery services assigned to a type
- view servers assigned to a type
- view cache groups assigned to a type


Statuses
--------
A table of all possible server statuses, with the ability to create, update, and delete statuses. This page shows a table of statuses with the following columns:

:Name:        The name of this status
:Description: A short, usually user-defined, description of this status

Status management includes the ability to (where applicable):

- create a new status
- update an existing status
- delete an existing status
- view status servers


Topology
========
'Topology' groups views and functionality that deal with how CDNs and their Traffic Control components are grouped and distributed, both on a logical level as well as a physical level.

.. figure:: ./images/tp_menu_topology.png
	:align: center

	'Topology' Menu

Cache Groups
------------
'Cache Groups' are sets of cache servers, typically grouped by geographic proximity. This menu allows user to add or remove caches from Cache Groups as well as creating, updating and deleting Cache Groups themselves. Each entry in the table of Cache Groups on this page has the following fields:

:Name:       The full name of this Cache Group
:Short Name: A shorter, more human-friendly name for this Cache Group
:Type:       The Type of this Cache Group (see :ref:`tp-configure-types`)
:Latitude:   A geographic latitude assigned to this Cache Group
:Longitude:  A geographic longitude assigned to this Cache Group

Cache group management includes the ability to (where applicable):

- create a new cache group
- update an existing cache group
- delete an existing cache group
- queue/clear updates for all servers in a cache group
- view cache group ASNs
- view and assign cache group parameters
- view cache group servers


Coordinates
-----------
'Coordinates' allows a label to be given to a set of geographic coordinates for ease of use. Each entry in the table on this page has the following fields:

:Name:      The name of this coordinate pair
:Latitude:  The geographic latitude part of the coordinate pair
:Longitude: The geographic longitude part of the coordinate pair

Coordination management includes the ability to (where applicable):

- create a new coordinate pair
- update an existing coordinate pair
- delete an existing coordinate pair


Phys Locations
--------------
A table of physical locations which may be assigned to servers and Cache Groups, typically for the purpose of optimizing client routing. Here they can be created, updated deleted and assigned. Each entry has the following columns:

:Name:       The full name of the physical location
:Short Name: A shorter, more human-friendly name for this physical location
:Address:    The location's street address (street number and name)
:City:       The city within which the location resides
:State:      The state within which the location's city lies
:Region:     The Region to which this physical location has been assigned

Physical location management includes the ability to (where applicable):

- create a new physical location
- update an existing physical location
- delete an existing physical location
- view physical location servers


Divisions
---------
Here Divisions may be created and deleted, and their constituent Regions may be viewed. Each entry in the table on this page has the following fields:

:Name: The name of the Division

Division management includes the ability to (where applicable):


Regions
-------
Regions are groups of Cache Groups, and are themselves grouped into Divisions. Each entry in the table on this page has the following fields:

:Name:     The name of this Region
:Division: The Division to which this Region is assigned

Region management includes the ability to (where applicable):

- create a new Region
- update an existing Region
- delete an existing Region
- view Region physical locations


ASNs
----
Manage Autonomous System Numbers (ASNs). Each entry in the table on this page has the following fields:

:ASN:         The actual ASN
:Cache Group: The Cache Group to which this ASN is assigned

ASN management includes the ability to (where applicable):

- create a new ASN
- update an existing ASN
- delete an existing ASN

.. seealso:: `Autonomous System (Internet) Wikipedia Page <https://en.wikipedia.org/wiki/Autonomous_system_(Internet)>`_


Tools
=====

.. figure:: ./images/tp_menu_tools.png
	:align: center
	:alt: The 'Tools' Menu

	The 'Tools' Menu

'Tools' contains various tools that don't directly relate to manipulating Traffic Control components or their groupings.

Invalidate Content
------------------
Here, specific assets can be invalidated in all caches of a Delivery Service, forcing content to be updated from the origin. Specifically, this *doesn't* mean that cache servers will immediately remove items from their caches, but rather will fetch new copies whenever a request is made matching the 'Asset URL' regular expression. This behavior persists until the Invalidate Content Job's Time To Live (TTL) expires. Each entry in the table on this page has the following fields:

:Delivery Service: The Delivery Service to which to apply this Invalidate Content Job
:Asset URL:        A URL or regular expression which describes the asset(s) to be invalidated
:Parameters:       So far, the only use for this is setting a TTL over which the Invalidate Content Job shall remain active
:Start:            An effective start time until which the job is delayed
:Created By:       The user name of the person who created this Invalidate Content Job

Invalidate content includes the ability to (where applicable):

- create a new invalidate content job

Generate ISO
------------
Generates a boot-able system image for any of the servers in the Servers table (or any server for that matter). Currently it only supports CentOS 6 or 7, but if you're brave and pure of heart you MIGHT be able to get it to work with other Unix-like Operating Systems. The interface is *mostly* self-explanatory, but here is a short explanation of the fields in that form.

Copy Server Attributes From
	Optional. This option lets the user choose a server from the Traffic Ops database and will auto-fill the other fields as much as possible based on that server's properties
OS Version
	This list is populated by modifying the ``osversions.cfg`` file on the Traffic Ops server. This file maps OS names to the name of a directory under ``app/public/iso/`` directory within the Traffic Ops install directory
Hostname
	The desired hostname of the resultant system
Domain
	The desired domain name of the resultant system
DHCP
	If this is 'no' the IP settings of the system must be specified, and the following extra fields will appear:

		IP Address
			The resultant system's IPv4 Address
		Network Subnet
			The system's network subnet mask
		Network Gateway
			The system's network gateway's IPv4 Address

Network MTU
	The system's network's Maximum Transmission Unit (MTU). Despite being a text field, this can only be 1500 or 9000 - it should almost always be 1500

		.. seealso:: `The Maximum transmission unit Wikipedia Page <https://en.wikipedia.org/wiki/Maximum_transmission_unit>`_

Disk for OS Install
	The disk on which to install the base system. A reasonable default is ``sda`` (the ``/dev/`` prefix is not necessary)
Root Password
	The password to be used for the root user. Input is MD5 hashed before being written to disk
Confirm Root Password
	Repeat the 'Root Password' to be sure it's right
Interface Name
	Optional. The name of the resultant system's network interface. Typical values are bond0, eth4, etc. If bond0 is entered, a Link Aggregation Control Protocol bonding configuration will be written

		.. seealso:: `The Link aggregation Wikipedia Page <https://en.wikipedia.org/wiki/Link_aggregation>`_

Stream ISO
	If this is 'yes', then the download will start immediately as the ISO is written directly to the socket connection to Traffic Ops. If this is 'no', then the download will begin only *after* the ISO has finished being generated. For almost all use cases, this should be 'yes'.


User Admin
==========

This section offers administrative functionality for users and their permissions.

.. figure:: ./images/tp_menu_user_admin.png
	:align: center
	:alt: The 'User Admin' Menu

	The 'User Admin' Menu

User
----
This page lists all the users that are visible to the user (so, for 'admin' users, all users will appear here). Each entry in the table on this page has the following fields:

:Full Name: The user's full, real name
:Username:  The user's username
:Email:     The user's email address
:Tenant:    The user's Tenant
:Role:      The user's Role

User management includes the ability to (where applicable):

- register a new user
- create a new user
- update an existing user
- view delivery services visible to a user


Tenants
-------
A 'Tenant' essentially groups users with the Delivery Services about which they're allowed to know. Each entry in the table on this page has the following entries:

:Name:   The name of the Tenant
:Active: If 'true' users of this Tenant group are allowed to login and have active Delivery Services
:Parent: The parent of this Tenant. The default is the 'root' Tenant, which has no users.

Tenant management includes the ability to (where applicable):

- create a new tenant
- update an existing tenant
- delete an existing tenant
- view users assigned to a tenant
- view delivery services assigned to a tenant

Roles
-----
'Roles' grant a user permissions to do certain things. Each entry in the table on this page has the following fields:

:Name:            The name of the role
:Privilege Level: The privilege level of this role. This is a whole number that actually controls what a user is allowed to do. Higher numbers correspond to higher permission levels
:Description:     A short description of the role and what it is allowed to do

Role management includes the ability to (where applicable):

- view all roles
- create new roles

.. note:: Roles cannot be deleted through the Traffic Portal UI

Other
=====
Custom menu items. By default, this contains only a link to the Traffic Control documentation.

.. figure:: ./images/tp_menu_other.png
	:align: center
	:alt: The 'Other' Menu

	The 'Other' Menu

Docs
----
This is just a link to `the Traffic Control Documentation <https://trafficcontrol.apache.org>`_.

Custom Menu Items
-----------------
This section is configurable in the ``tp.domain.com/traffic_portal_properties.json`` configuration file, in the ``customMenu`` section.
