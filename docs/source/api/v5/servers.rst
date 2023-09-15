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

.. _to-api-servers:

***********
``servers``
***********

``GET``
=======
Retrieves properties of all servers across all CDNs.

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: SERVER:READ, DELIVERY-SERVICE:READ, CDN:READ, PHYSICAL-LOCATION:READ, CACHE-GROUP:READ, TYPE:READ, PROFILE:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+--------------------+----------+-------------------------------------------------------------------------------------------------------------------+
	| Name               | Required | Description                                                                                                       |
	+====================+==========+===================================================================================================================+
	| cacheGroup         | no       | Return only those servers within the :term:`Cache Group` that has this :ref:`cache-group-name`                    |
	+--------------------+----------+-------------------------------------------------------------------------------------------------------------------+
	| cachegroup         | no       | Return only those servers within the :term:`Cache Group` that has this :ref:`cache-group-id`                      |
	+--------------------+----------+-------------------------------------------------------------------------------------------------------------------+
	| cacheGroupName     | no       | Return only those servers within the :term:`Cache Group` that has this :ref:`cache-group-name`                    |
	+--------------------+----------+-------------------------------------------------------------------------------------------------------------------+
	| dsId               | no       | Return only those servers assigned to the :term:`Delivery Service` identified by this integral, unique identifier.|
	|                    |          | If the Delivery Service has a :term:`Topology` assigned to it, the :ref:`to-api-servers` endpoint will return     |
	|                    |          | each server whose :term:`Cache Group` is associated with a :term:`Topology Node` of that Topology and has the     |
	|                    |          | :term:`Server Capabilities` that are                                                                              |
	|                    |          | :term:`required by the Delivery Service <Delivery Service required capabilities>` but excluding                   |
	|                    |          | :term:`Origin Servers` that are not assigned to the Delivery Service. For more information, see                   |
	|                    |          | :ref:`multi-site-origin-qht`.                                                                                     |
	+--------------------+----------+-------------------------------------------------------------------------------------------------------------------+
	| hostName           | no       | Return only those servers that have this (short) hostname                                                         |
	+--------------------+----------+-------------------------------------------------------------------------------------------------------------------+
	| id                 | no       | Return only the server with this integral, unique identifier                                                      |
	+--------------------+----------+-------------------------------------------------------------------------------------------------------------------+
	| physicalLocation   | no       | Return only servers that have the :term:`Physical Location` with this Name                                        |
	+--------------------+----------+-------------------------------------------------------------------------------------------------------------------+
	| physicalLocationID | no       | Return only servers that have the :term:`Physical Location` with this integral, unique identifier                 |
	+--------------------+----------+-------------------------------------------------------------------------------------------------------------------+
	| status             | no       | Return only those servers with this status - see :ref:`health-proto`                                              |
	+--------------------+----------+-------------------------------------------------------------------------------------------------------------------+
	| type               | no       | Return only servers of this :term:`Type`                                                                          |
	+--------------------+----------+-------------------------------------------------------------------------------------------------------------------+
	| topology           | no       | Return only servers who belong to cacheGroups assigned to the :term:`Topology` identified by this name            |
	+--------------------+----------+-------------------------------------------------------------------------------------------------------------------+
	| sortOrder          | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")                          |
	+--------------------+----------+-------------------------------------------------------------------------------------------------------------------+
	| limit              | no       | Choose the maximum number of results to return                                                                    |
	+--------------------+----------+-------------------------------------------------------------------------------------------------------------------+
	| offset             | no       | The number of results to skip before beginning to return results. Must use in conjunction with limit              |
	+--------------------+----------+-------------------------------------------------------------------------------------------------------------------+
	| page               | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long and   |
	|                    |          | the first page is 1. If ``offset`` was defined, this query parameter has no effect. ``limit`` must be defined to  |
	|                    |          | make use of ``page``.                                                                                             |
	+--------------------+----------+-------------------------------------------------------------------------------------------------------------------+

.. deprecated:: ATCv8
	Rather than ``cachegroup`` or ``cachegroupName``, prefer ``cacheGroup`` as the other two are deprecated.

.. code-block:: http
	:caption: Request Example

	GET /api/5.0/servers?hostName=mid HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:cacheGroup: A string that is the :ref:`name of the Cache Group <cache-group-name>` to which the server belongs

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``cachegroup`` - improperly formatted camelCase.

:cacheGroupID: An integer that is the :ref:`ID of the Cache Group <cache-group-id>` to which the server belongs

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``cachegroupId`` - improperly formatted camelCase.

:cdnID: The integral, unique identifier of the CDN to which the server belongs

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``cdnId`` - improperly formatted camelCase.

:cdn: Name of the CDN to which the server belongs

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``cdnName``. It has been changed for consistency with others e.g. ``type``, ``status``, etc.

:configUpdateTime:   The last time an update was requested for this server. This field defaults to standard epoch
:configApplyTime:    The last time an update was applied for this server. This field defaults to standard epoch
:configUpdateFailed: If the last update applied for this server was applied successfully. Defaults to false.
:domainName:       The domain part of the server's :abbr:`FQDN (Fully Qualified Domain Name)`
:guid:             An identifier used to uniquely identify the server

	.. note:: This is a legacy key which only still exists for compatibility reasons - it should always be ``null``

:hostName:     The (short) hostname of the server
:httpsPort:    The port on which the server listens for incoming HTTPS connections/requests
:id:           An integral, unique identifier for this server
:iloIpAddress: The IPv4 address of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:iloIpGateway: The IPv4 gateway address of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:iloIpNetmask: The IPv4 subnet mask of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:iloPassword:  The password of the of the server's :abbr:`ILO (Integrated Lights-Out)` service user\ [#ilo]_ - displays as simply ``******`` if the currently logged-in user does not have the SECURE-SERVER:READ permission.
:iloUsername:  The user name for the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:interfaces:   A set of the network interfaces in use by the server. In most scenarios, only one will be present, but it is illegal for this set to be an empty collection.

	:ipAddresses:       A set of objects representing IP Addresses assigned to this network interface. In most scenarios, only one or two (usually one IPv4 address and one IPv6 address) will be present, but it is illegal for this set to be an empty collection.

		:address:        The actual IP address, including any mask as a CIDR-notation suffix
		:gateway:        Either the IP address of the network gateway for this address, or ``null`` to signify that no such gateway exists
		:serviceAddress: A boolean that describes whether or not the server's main service is available at this IP address. When this property is ``true``, the IP address is referred to as a "service address". It is illegal for a server to not have at least one service address. It is also illegal for a server to have more than one service address of the same address family (i.e. more than one IPv4 service address and/or more than one IPv6 address). Finally, all service addresses for a server must be contained within one interface - which is therefore sometimes referred to as the "service interface" for the server.

	:maxBandwidth:      The maximum healthy bandwidth allowed for this interface. If bandwidth exceeds this limit, Traffic Monitors will consider the entire server unhealthy - which includes *all* configured network interfaces. If this is ``null``, it has the meaning "no limit". It has no effect if ``monitor`` is not true for this interface.

		.. seealso:: :ref:`health-proto`

	:monitor:           A boolean which describes whether or not this interface should be monitored by Traffic Monitor for statistics and health consideration.
	:mtu:               The :abbr:`MTU (Maximum Transmission Unit)` of this interface. If it is ``null``, it may be assumed that the information is either not available or not applicable for this interface.
	:name:              The name of the interface. No two interfaces of the same server may share a name. It is the same as the network interface's device name on the server, e.g. ``eth0``.
	:routerPortName:    The human-readable name of the router responsible for reaching this server's interface.
	:routerPortName:    The human-readable name of the port used by the router responsible for reaching this server's interface.

:lastUpdated: The date and time at which this server description was last modified, in :RFC:`3339` format

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was given in :ref:`non-rfc-datetime`.

:mgmtIpAddress: The IPv4 address of some network interface on the server used for 'management'

	.. deprecated:: 3.0
		This field is deprecated and will be removed in a future API version. Operators should migrate this data into the ``interfaces`` property of the server.

:mgmtIpGateway: The IPv4 address of a gateway used by some network interface on the server used for 'management'

	.. deprecated:: 3.0
		This field is deprecated and will be removed in a future API version. Operators should migrate this data into the ``interfaces`` property of the server.

:mgmtIpNetmask: The IPv4 subnet mask used by some network interface on the server used for 'management'

	.. deprecated:: 3.0
		This field is deprecated and will be removed in a future API version. Operators should migrate this data into the ``interfaces`` property of the server.

:offlineReason:    A user-entered reason why the server is in ADMIN_DOWN or OFFLINE status
:physicalLocation: The name of the physical location where the server resides

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``physLocation`` - improperly formatted camelCase.

:physicalLocationID: An integral, unique identifier for the physical location where the server resides

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``physLocationId`` - improperly formatted camelCase.

:profiles: List of :ref:`profile-name` of the :term:`Profiles` used by this server

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``profileNames`` - it has been changed because now that this is the only identifying information for a :term:Profile that exists on a server, there is no need to distinguish it from, say, an ID.

:revalUpdateTime:   The last time a content invalidation/revalidation request was submitted for this server. This field defaults to standard epoch
:revalApplyTime:    The last time a content invalidation/revalidation request was applied by this server. This field defaults to standard epoch
:revalUpdateFailed: If the last content invalidation/revalidation applied for this server was applied successfully. Defaults to false.
:rack:              A string indicating "server rack" location
:status:            The :term:`Status` of the server

	.. seealso:: :ref:`health-proto`

:statusID: The integral, unique identifier of the status of this server

	.. seealso:: :ref:`health-proto`

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``statusID`` - improperly formatted camelCase.

:tcpPort: The port on which this server listens for incoming TCP connections

	.. note:: This is typically thought of as synonymous with "HTTP port", as the port specified by ``httpsPort`` may also be used for incoming TCP connections.

:type:   The name of the :term:`Type` of this server
:typeID: The integral, unique identifier of the 'type' of this server

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``typeID`` - improperly formatted camelCase.

:xmppId:     A system-generated UUID used to generate a server hashId for use in Traffic Router's consistent hashing algorithm. This value is set when a server is created and cannot be changed afterwards.
:xmppPasswd: The password used in XMPP communications with the server - displays as simply ``******`` if the currently logged-in user does not have the SECURE-SERVER:READ permission.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Tue, 19 May 2020 17:06:25 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 19 May 2020 16:06:25 GMT
	Content-Length: 538

	{ "response": [{
		"cacheGroup": "CDN_in_a_Box_Mid",
		"cacheGroupID": 6,
		"cdnID": 2,
		"cdn": "CDN-in-a-Box",
		"configUpdateTime": "1969-12-31T17:00:00-07:00",
		"configApplyTime": "1969-12-31T17:00:00-07:00",
		"configUpdateFailed": false,
		"domainName": "infra.ciab.test",
		"guid": null,
		"hostName": "mid",
		"httpsPort": 443,
		"id": 12,
		"iloIpAddress": "",
		"iloIpGateway": "",
		"iloIpNetmask": "",
		"iloPassword": "",
		"iloUsername": "",
		"lastUpdated": "2020-05-19T14:49:39Z",
		"mgmtIpAddress": "",
		"mgmtIpGateway": "",
		"mgmtIpNetmask": "",
		"offlineReason": "",
		"physicalLocation": "Apachecon North America 2018",
		"physicalLocationID": 1,
		"profiles": ["ATS_MID_TIER_CACHE"],
		"rack": "",
		"revalUpdateTime": "1969-12-31T17:00:00-07:00",
		"revalApplyTime": "1969-12-31T17:00:00-07:00",
		"revalUpdateFailed": false,
		"status": "REPORTED",
		"statusID": 3,
		"tcpPort": 80,
		"type": "MID",
		"typeID": 12,
		"xmppId": "",
		"xmppPasswd": "",
		"interfaces": [
			{
				"ipAddresses": [
					{
						"address": "172.26.0.4/16",
						"gateway": "172.26.0.1",
						"serviceAddress": true
					}
				],
				"maxBandwidth": null,
				"monitor": false,
				"mtu": 1500,
				"name": "eth0",
				"routerHostName": "",
				"routerPortName": ""
			}
		]
	}]}

``POST``
========
Allows a user to create a new server.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: SERVER:CREATE, SERVER:READ, DELIVERY-SERVICE:READ, CDN:READ, PHYSICAL-LOCATION:READ, CACHE-GROUP:READ, TYPE:READ, PROFILE:READ
:Response Type:  Object

Request Structure
-----------------
:cacheGroupID: An integer that is the :ref:`ID of the Cache Group <cache-group-id>` to which the server shall belong
:cdnID:        The integral, unique identifier of the CDN to which the server shall belong

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``cdnId`` - improperly formatted camelCase.

:domainName:   The domain part of the server's :abbr:`FQDN (Fully Qualified Domain Name)`
:hostName:     The (short) hostname of the server
:httpsPort:    An optional port number on which the server listens for incoming HTTPS connections/requests
:iloIpAddress: An optional IPv4 address of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:iloIpGateway: An optional IPv4 gateway address of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:iloIpNetmask: An optional IPv4 subnet mask of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:iloPassword:  An optional string containing the password of the of the server's :abbr:`ILO (Integrated Lights-Out)` service user\ [#ilo]_ - displays as simply ``******`` if the currently logged-in user does not have the SECURE-SERVER:READ permission.
:iloUsername:  An optional string containing the user name for the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:interfaces:   A set of the network interfaces in use by the server. In most scenarios, only one will be necessary, but it is illegal for this set to be an empty collection.

	:ipAddresses:       A set of objects representing IP Addresses assigned to this network interface. In most scenarios, only one or two (usually one IPv4 address and one IPv6 address) will be necessary, but it is illegal for this set to be an empty collection.

		:address:        The actual IP address, including any mask as a CIDR-notation suffix
		:gateway:        Either the IP address of the network gateway for this address, or ``null`` to signify that no such gateway exists
		:serviceAddress: A boolean that describes whether or not the server's main service is available at this IP address. When this property is ``true``, the IP address is referred to as a "service address". It is illegal for a server to not have at least one service address. It is also illegal for a server to have more than one service address of the same address family (i.e. more than one IPv4 service address and/or more than one IPv6 address). Finally, all service addresses for a server must be contained within one interface - which is therefore sometimes referred to as the "service interface" for the server.

	:maxBandwidth:      The maximum healthy bandwidth allowed for this interface. If bandwidth exceeds this limit, Traffic Monitors will consider the entire server unhealthy - which includes *all* configured network interfaces. If this is ``null``, it has the meaning "no limit". It has no effect if ``monitor`` is not true for this interface.

		.. seealso:: :ref:`health-proto`

	:monitor:           A boolean which describes whether or not this interface should be monitored by Traffic Monitor for statistics and health consideration.
	:mtu:               The :abbr:`MTU (Maximum Transmission Unit)` of this interface. If it is ``null``, it may be assumed that the information is either not available or not applicable for this interface.
	:name:              The name of the interface. No two interfaces of the same server may share a name. It is the same as the network interface's device name on the server, e.g. ``eth0``.
	:routerPortName:    The human-readable name of the router responsible for reaching this server's interface.
	:routerPortName:    The human-readable name of the port used by the router responsible for reaching this server's interface.

:mgmtIpAddress: The IPv4 address of some network interface on the server used for 'management'

	.. deprecated:: 3.0
		This field is deprecated and will be removed in a future API version. Operators should migrate this data into the ``interfaces`` property of the server.

:mgmtIpGateway: The IPv4 address of a gateway used by some network interface on the server used for 'management'

	.. deprecated:: 3.0
		This field is deprecated and will be removed in a future API version. Operators should migrate this data into the ``interfaces`` property of the server.

:mgmtIpNetmask: The IPv4 subnet mask used by some network interface on the server used for 'management'

	.. deprecated:: 3.0
		This field is deprecated and will be removed in a future API version. Operators should migrate this data into the ``interfaces`` property of the server.

:physicalLocationID: An integral, unique identifier for the physical location where the server resides

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``physLocationId`` - improperly formatted camelCase.

:profiles: List of :ref:`profile-name` of the :term:`Profiles` that shall be used by this server

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``profileNames`` - it has been changed because now that this is the only identifying information for a :term:Profile that exists on a server, there is no need to distinguish it from, say, an ID.

:rack:     An optional string indicating "server rack" location
:statusID: The integral, unique identifier of the status of this server

	.. seealso:: :ref:`health-proto`

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``statusId`` - improperly formatted camelCase.

:tcpPort: An optional port number on which this server listens for incoming TCP connections

	.. note:: This is typically thought of as synonymous with "HTTP port", as the port specified by ``httpsPort`` may also be used for incoming TCP connections.

:typeID: The integral, unique identifier of the 'type' of this server

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``typeId`` - improperly formatted camelCase.

:xmppId:     A system-generated UUID used to generate a server hashId for use in Traffic Router's consistent hashing algorithm. This value is set when a server is created and cannot be changed afterwards.
:xmppPasswd: An optional password used in XMPP communications with the server - displays as simply ``******`` if the currently logged-in user does not have the SECURE-SERVER:READ permission.

.. code-block:: http
	:caption: Request Example

	POST /api/5.0/servers HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 599
	Content-Type: application/json

	{
		"cacheGroupID": 6,
		"cdnID": 2,
		"domainName": "infra.ciab.test",
		"hostName": "test",
		"httpsPort": 443,
		"iloIpAddress": "",
		"iloIpGateway": "",
		"iloIpNetmask": "",
		"iloPassword": "",
		"iloUsername": "",
		"interfaces": [
			{
				"ipAddresses": [
					{
						"address": "::1",
						"gateway": "::2",
						"serviceAddress": true
					},
					{
						"address": "0.0.0.1/24",
						"gateway": "0.0.0.2",
						"serviceAddress": true
					}
				],
				"maxBandwidth": null,
				"monitor": true,
				"mtu": 1500,
				"name": "eth0",
				"routerHostName": "",
				"routerPortName": ""
			}
		],
		"interfaceMtu": 1500,
		"interfaceName": "eth0",
		"ip6Address": "::1",
		"ip6Gateway": "::2",
		"ipAddress": "0.0.0.1",
		"ipGateway": "0.0.0.2",
		"ipNetmask": "255.255.255.0",
		"mgmtIpAddress": "",
		"mgmtIpGateway": "",
		"mgmtIpNetmask": "",
		"offlineReason": "",
		"physicalLocationID": 1,
		"profiles": ["ATS_MID_TIER_CACHE"],
		"statusID": 3,
		"tcpPort": 80,
		"typeID": 12
	}

Response Structure
------------------
:cacheGroup: A string that is the :ref:`name of the Cache Group <cache-group-name>` to which the server belongs

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``cachegroup`` - improperly formatted camelCase.

:cacheGroupID: An integer that is the :ref:`ID of the Cache Group <cache-group-id>` to which the server belongs

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``cachegroupId`` - improperly formatted camelCase.

:cdnID: The integral, unique identifier of the CDN to which the server belongs

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``cdnId`` - improperly formatted camelCase.

:cdn: Name of the CDN to which the server belongs

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``cdnName``. It has been changed for consistency with others e.g. ``type``, ``status``, etc.

:configUpdateTime: The last time an update was requested for this server. This field defaults to standard epoch
:configApplyTime:  The last time an update was applied for this server. This field defaults to standard epoch
:configUpdateFailed: If the last update applied for this server was applied successfully. Defaults to false.
:domainName:       The domain part of the server's :abbr:`FQDN (Fully Qualified Domain Name)`
:guid:             An identifier used to uniquely identify the server

	.. note:: This is a legacy key which only still exists for compatibility reasons - it should always be ``null``

:hostName:     The (short) hostname of the server
:httpsPort:    The port on which the server listens for incoming HTTPS connections/requests
:id:           An integral, unique identifier for this server
:iloIpAddress: The IPv4 address of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:iloIpGateway: The IPv4 gateway address of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:iloIpNetmask: The IPv4 subnet mask of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:iloPassword:  The password of the of the server's :abbr:`ILO (Integrated Lights-Out)` service user\ [#ilo]_ - displays as simply ``******`` if the currently logged-in user does not have the SECURE-SERVER:READ permission.
:iloUsername:  The user name for the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:interfaces:   A set of the network interfaces in use by the server. In most scenarios, only one will be present, but it is illegal for this set to be an empty collection.

	:ipAddresses:       A set of objects representing IP Addresses assigned to this network interface. In most scenarios, only one or two (usually one IPv4 address and one IPv6 address) will be present, but it is illegal for this set to be an empty collection.

		:address:        The actual IP address, including any mask as a CIDR-notation suffix
		:gateway:        Either the IP address of the network gateway for this address, or ``null`` to signify that no such gateway exists
		:serviceAddress: A boolean that describes whether or not the server's main service is available at this IP address. When this property is ``true``, the IP address is referred to as a "service address". It is illegal for a server to not have at least one service address. It is also illegal for a server to have more than one service address of the same address family (i.e. more than one IPv4 service address and/or more than one IPv6 address). Finally, all service addresses for a server must be contained within one interface - which is therefore sometimes referred to as the "service interface" for the server.

	:maxBandwidth:      The maximum healthy bandwidth allowed for this interface. If bandwidth exceeds this limit, Traffic Monitors will consider the entire server unhealthy - which includes *all* configured network interfaces. If this is ``null``, it has the meaning "no limit". It has no effect if ``monitor`` is not true for this interface.

		.. seealso:: :ref:`health-proto`

	:monitor:           A boolean which describes whether or not this interface should be monitored by Traffic Monitor for statistics and health consideration.
	:mtu:               The :abbr:`MTU (Maximum Transmission Unit)` of this interface. If it is ``null``, it may be assumed that the information is either not available or not applicable for this interface.
	:name:              The name of the interface. No two interfaces of the same server may share a name. It is the same as the network interface's device name on the server, e.g. ``eth0``.
	:routerPortName:    The human-readable name of the router responsible for reaching this server's interface.
	:routerPortName:    The human-readable name of the port used by the router responsible for reaching this server's interface.

:lastUpdated: The date and time at which this server description was last modified, in :RFC:`3339` format

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was given in :ref:`non-rfc-datetime`.

:mgmtIpAddress: The IPv4 address of some network interface on the server used for 'management'

	.. deprecated:: 3.0
		This field is deprecated and will be removed in a future API version. Operators should migrate this data into the ``interfaces`` property of the server.

:mgmtIpGateway: The IPv4 address of a gateway used by some network interface on the server used for 'management'

	.. deprecated:: 3.0
		This field is deprecated and will be removed in a future API version. Operators should migrate this data into the ``interfaces`` property of the server.

:mgmtIpNetmask: The IPv4 subnet mask used by some network interface on the server used for 'management'

	.. deprecated:: 3.0
		This field is deprecated and will be removed in a future API version. Operators should migrate this data into the ``interfaces`` property of the server.

:offlineReason:    A user-entered reason why the server is in ADMIN_DOWN or OFFLINE status
:physicalLocation: The name of the :term:`Physical Location` where the server resides

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``physLocation`` - improperly formatted camelCase.

:physicalLocationID: An integral, unique identifier for the :term:`Physical Location` where the server resides

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``physLocationId`` - improperly formatted camelCase.

:profiles: List of :ref:`profile-name` of the :term:`Profiles` used by this server

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``profileNames`` - it has been changed because now that this is the only identifying information for a :term:Profile that exists on a server, there is no need to distinguish it from, say, an ID.

:revalUpdateTime: The last time a content invalidation/revalidation request was submitted for this server. This field defaults to standard epoch
:revalApplyTime:  The last time a content invalidation/revalidation request was applied by this server. This field defaults to standard epoch
:revalUpdateFailed: If the last content invalidation/revalidation applied for this server was applied successfully. Defaults to false.
:rack:            A string indicating "server rack" location
:status:          The status of the server

	.. seealso:: :ref:`health-proto`

:statusID: The integral, unique identifier of the status of this server

	.. seealso:: :ref:`health-proto`

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``statusId`` - improperly formatted camelCase.

:tcpPort: The port on which this server listens for incoming TCP connections

	.. note:: This is typically thought of as synonymous with "HTTP port", as the port specified by ``httpsPort`` may also be used for incoming TCP connections.

:type:   The name of the 'type' of this server
:typeID: The integral, unique identifier of the 'type' of this server

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was known by the name ``typeId`` - improperly formatted camelCase.

:xmppId:     A system-generated UUID used to generate a server hashId for use in Traffic Router's consistent hashing algorithm. This value is set when a server is created and cannot be changed afterwards.
:xmppPasswd: The password used in XMPP communications with the server - displays as simply ``******`` if the currently logged-in user does not have the SECURE-SERVER:READ permission.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 201 Created
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Tue, 19 May 2020 17:34:40 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 19 May 2020 16:34:40 GMT
	Content-Length: 562

	{ "alerts": [
		{
			"text": "Server created",
			"level": "success"
		}
	],
	"response": {
		"cacheGroup": "CDN_in_a_Box_Mid",
		"cacheGroupID": 6,
		"cdnID": 2,
		"cdn": "CDN-in-a-Box",
		"configUpdateTime": "1969-12-31T17:00:00-07:00",
		"configApplyTime": "1969-12-31T17:00:00-07:00",
		"configUpdateFailed": false,
		"domainName": "infra.ciab.test",
		"guid": null,
		"hostName": "test",
		"httpsPort": 443,
		"id": 14,
		"iloIpAddress": "",
		"iloIpGateway": "",
		"iloIpNetmask": "",
		"iloPassword": "",
		"iloUsername": "",
		"lastUpdated": "2020-05-19 16:34:40+00",
		"mgmtIpAddress": "",
		"mgmtIpGateway": "",
		"mgmtIpNetmask": "",
		"offlineReason": "",
		"physicalLocation": "Apachecon North America 2018",
		"physicalLocationID": 1,
		"profiles": ["ATS_MID_TIER_CACHE"],
		"rack": null,
		"revalUpdateTime": "1969-12-31T17:00:00-07:00",
		"revalApplyTime": "1969-12-31T17:00:00-07:00",
		"revalUpdateFailed": false,
		"status": "REPORTED",
		"statusID": 3,
		"tcpPort": 80,
		"type": "MID",
		"typeID": 12,
		"xmppId": null,
		"xmppPasswd": null,
		"interfaces": [
			{
				"ipAddresses": [
					{
						"address": "::1",
						"gateway": "::2",
						"serviceAddress": true
					},
					{
						"address": "0.0.0.1/24",
						"gateway": "0.0.0.2",
						"serviceAddress": true
					}
				],
				"maxBandwidth": null,
				"monitor": true,
				"mtu": 1500,
				"name": "eth0",
				"routerHostName": "",
				"routerPortName": ""
			}
		]
	}}

.. [#ilo] For more information see the `Wikipedia page on Lights-Out management <https://en.wikipedia.org/wiki/Out-of-band_management>`_\ .
