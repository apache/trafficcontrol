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

.. _to-api-v3-servers-id:

******************
``servers/{{ID}}``
******************

``PUT``
=======
Allow user to edit a server.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------+
	| Name |                Description                  |
	+======+=============================================+
	|  ID  | The integral, unique identifier of a server |
	+------+---------------------------------------------+

:cachegroupId: An integer that is the :ref:`ID of the Cache Group <cache-group-id>` to which the server shall belong
:cdnId:        The integral, unique identifier of the CDN to which the server shall belong
:domainName:   The domain part of the server's :abbr:`FQDN (Fully Qualified Domain Name)`
:hostName:     The (short) hostname of the server
:httpsPort:    An optional port number on which the server listens for incoming HTTPS connections/requests
:iloIpAddress: An optional IPv4 address of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:iloIpGateway: An optional IPv4 gateway address of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:iloIpNetmask: An optional IPv4 subnet mask of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:iloPassword:  An optional string containing the password of the of the server's :abbr:`ILO (Integrated Lights-Out)` service user\ [#ilo]_ - displays as simply ``******`` if the currently logged-in user does not have the 'admin' or 'operations' :abbr:`Role(s) <Role>`
:iloUsername:  An optional string containing the user name for the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:interfaces:   A set of the network interfaces in use by the server. In most scenarios, only one will be necessary, but it is illegal for this set to be an empty collection.

	:ipAddresses: A set of objects representing IP Addresses assigned to this network interface. In most scenarios, only one or two (usually one IPv4 address and one IPv6 address) will be necessary, but it is illegal for this set to be an empty collection.

		:address:        The actual IP address, including any mask as a CIDR-notation suffix
		:gateway:        Either the IP address of the network gateway for this address, or ``null`` to signify that no such gateway exists
		:serviceAddress: A boolean that describes whether or not the server's main service is available at this IP address. When this property is ``true``, the IP address is referred to as a "service address". It is illegal for a server to not have at least one service address. It is also illegal for a server to have more than one service address of the same address family (i.e. more than one IPv4 service address and/or more than one IPv6 address). Finally, all service addresses for a server must be contained within one interface - which is therefore sometimes referred to as the "service interface" for the server.

	:maxBandwidth: The maximum healthy bandwidth allowed for this interface. If bandwidth exceeds this limit, Traffic Monitors will consider the entire server unhealthy - which includes *all* configured network interfaces. If this is ``null``, it has the meaning "no limit". It has no effect if ``monitor`` is not true for this interface.

		.. seealso:: :ref:`health-proto`

	:monitor: A boolean which describes whether or not this interface should be monitored by Traffic Monitor for statistics and health consideration.
	:mtu:     The :abbr:`MTU (Maximum Transmission Unit)` of this interface. If it is ``null``, it may be assumed that the information is either not available or not applicable for this interface. This unsigned integer must not be less than 1280.
	:name:    The name of the interface. No two interfaces of the same server may share a name. It is the same as the network interface's device name on the server, e.g. ``eth0``.

:mgmtIpAddress: The IPv4 address of some network interface on the server used for 'management'

	.. deprecated:: 3.0
		This field is deprecated and will be removed in a future API version. Operators should migrate this data into the ``interfaces`` property of the server.

:mgmtIpGateway: The IPv4 address of a gateway used by some network interface on the server used for 'management'

	.. deprecated:: 3.0
		This field is deprecated and will be removed in a future API version. Operators should migrate this data into the ``interfaces`` property of the server.

:mgmtIpNetmask: The IPv4 subnet mask used by some network interface on the server used for 'management'

	.. deprecated:: 3.0
		This field is deprecated and will be removed in a future API version. Operators should migrate this data into the ``interfaces`` property of the server.

:physLocationId: An integral, unique identifier for the physical location where the server resides
:profileId:      The :ref:`profile-id` the :term:`Profile` that shall be used by this server
:revalPending:   A boolean value which, if ``true`` indicates that this server has pending content invalidation/revalidation
:rack:           An optional string indicating "server rack" location
:routerHostName: An optional string containing the human-readable name of the router responsible for reaching this server
:routerPortName: An optional string containing the human-readable name of the port used by the router responsible for reaching this server
:statusId:       The integral, unique identifier of the status of this server

	.. seealso:: :ref:`health-proto`

:tcpPort: An optional port number on which this server listens for incoming TCP connections

	.. note:: This is typically thought of as synonymous with "HTTP port", as the port specified by ``httpsPort`` may also be used for incoming TCP connections.

:typeId:     The integral, unique identifier of the 'type' of this server
:updPending: A boolean value which, if ``true``, indicates that the server has updates of some kind pending, typically to be acted upon by Traffic Control Cache Config (:term:`t3c`, formerly ORT)
:xmppId:     A system-generated UUID used to generate a server hashId for use in Traffic Router's consistent hashing algorithm. This value is set when a server is created and cannot be changed afterwards.
:xmppPasswd: An optional password used in XMPP communications with the server

.. code-block:: http
	:caption: Request Example

	PUT /api/3.0/servers/14 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 599
	Content-Type: application/json

	{
		"cachegroupId": 6,
		"cdnId": 2,
		"domainName": "infra.ciab.test",
		"hostName": "quest",
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
						"serviceAddress": false
					}
				],
				"maxBandwidth": null,
				"monitor": true,
				"mtu": 1500,
				"name": "bond0"
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
		"physLocationId": 1,
		"profileId": 10,
		"routerHostName": "",
		"routerPortName": "",
		"statusId": 3,
		"tcpPort": 80,
		"typeId": 12,
		"updPending": false
	}

Response Structure
------------------
:cachegroup:   A string that is the :ref:`name of the Cache Group <cache-group-name>` to which the server belongs
:cachegroupId: An integer that is the :ref:`ID of the Cache Group <cache-group-id>` to which the server belongs
:cdnId:        The integral, unique identifier of the CDN to which the server belongs
:cdnName:      Name of the CDN to which the server belongs
:domainName:   The domain part of the server's :abbr:`FQDN (Fully Qualified Domain Name)`
:guid:         An identifier used to uniquely identify the server

	.. note:: This is a legacy key which only still exists for compatibility reasons - it should always be ``null``

:hostName:       The (short) hostname of the server
:httpsPort:      The port on which the server listens for incoming HTTPS connections/requests
:id:             An integral, unique identifier for this server
:iloIpAddress:   The IPv4 address of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:iloIpGateway:   The IPv4 gateway address of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:iloIpNetmask:   The IPv4 subnet mask of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:iloPassword:    The password of the of the server's :abbr:`ILO (Integrated Lights-Out)` service user\ [#ilo]_ - displays as simply ``******`` if the currently logged-in user does not have the 'admin' or 'operations' :abbr:`Role(s) <Role>`
:iloUsername:    The user name for the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:interfaces:     A set of the network interfaces in use by the server. In most scenarios, only one will be present, but it is illegal for this set to be an empty collection.

	:ipAddresses: A set of objects representing IP Addresses assigned to this network interface. In most scenarios, only one or two (usually one IPv4 address and one IPv6 address) will be present, but it is illegal for this set to be an empty collection.

		:address:        The actual IP address, including any mask as a CIDR-notation suffix
		:gateway:        Either the IP address of the network gateway for this address, or ``null`` to signify that no such gateway exists
		:serviceAddress: A boolean that describes whether or not the server's main service is available at this IP address. When this property is ``true``, the IP address is referred to as a "service address". It is illegal for a server to not have at least one service address. It is also illegal for a server to have more than one service address of the same address family (i.e. more than one IPv4 service address and/or more than one IPv6 address). Finally, all service addresses for a server must be contained within one interface - which is therefore sometimes referred to as the "service interface" for the server.

	:maxBandwidth: The maximum healthy bandwidth allowed for this interface. If bandwidth exceeds this limit, Traffic Monitors will consider the entire server unhealthy - which includes *all* configured network interfaces. If this is ``null``, it has the meaning "no limit". It has no effect if ``monitor`` is not true for this interface.

		.. seealso:: :ref:`health-proto`

	:monitor: A boolean which describes whether or not this interface should be monitored by Traffic Monitor for statistics and health consideration.
	:mtu:     The :abbr:`MTU (Maximum Transmission Unit)` of this interface. If it is ``null``, it may be assumed that the information is either not available or not applicable for this interface.
	:name:    The name of the interface. No two interfaces of the same server may share a name. It is the same as the network interface's device name on the server, e.g. ``eth0``.

:lastUpdated:   The date and time at which this server description was last modified
:mgmtIpAddress: The IPv4 address of some network interface on the server used for 'management'

	.. deprecated:: 3.0
		This field is deprecated and will be removed in a future API version. Operators should migrate this data into the ``interfaces`` property of the server.

:mgmtIpGateway: The IPv4 address of a gateway used by some network interface on the server used for 'management'

	.. deprecated:: 3.0
		This field is deprecated and will be removed in a future API version. Operators should migrate this data into the ``interfaces`` property of the server.

:mgmtIpNetmask: The IPv4 subnet mask used by some network interface on the server used for 'management'

	.. deprecated:: 3.0
		This field is deprecated and will be removed in a future API version. Operators should migrate this data into the ``interfaces`` property of the server.

:offlineReason:  A user-entered reason why the server is in ADMIN_DOWN or OFFLINE status
:physLocation:   The name of the :term:`Physical Location` where the server resides
:physLocationId: An integral, unique identifier for the :term:`Physical Location` where the server resides
:profile:        The :ref:`profile-name` of the :term:`Profile` used by this server
:profileDesc:    A :ref:`profile-description` of the :term:`Profile` used by this server
:profileId:      The :ref:`profile-id` the :term:`Profile` used by this server
:revalPending:   A boolean value which, if ``true`` indicates that this server has pending content invalidation/revalidation
:rack:           A string indicating "server rack" location
:routerHostName: The human-readable name of the router responsible for reaching this server
:routerPortName: The human-readable name of the port used by the router responsible for reaching this server
:status:         The status of the server

	.. seealso:: :ref:`health-proto`

:statusId: The integral, unique identifier of the status of this server

	.. seealso:: :ref:`health-proto`

:tcpPort: The port on which this server listens for incoming TCP connections

	.. note:: This is typically thought of as synonymous with "HTTP port", as the port specified by ``httpsPort`` may also be used for incoming TCP connections.

:type:       The name of the 'type' of this server
:typeId:     The integral, unique identifier of the 'type' of this server
:updPending: A boolean value which, if ``true``, indicates that the server has updates of some kind pending, typically to be acted upon by Traffic Control Cache Config (T3C, formerly ORT)
:xmppId:     A system-generated UUID used to generate a server hashId for use in Traffic Router's consistent hashing algorithm. This value is set when a server is created and cannot be changed afterwards.
:xmppPasswd: The password used in XMPP communications with the server

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Tue, 19 May 2020 17:46:33 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 19 May 2020 16:46:33 GMT
	Content-Length: 566

	{ "alerts": [
		{
			"text": "Server updated",
			"level": "success"
		}
	],
	"response": {
		"cachegroup": "CDN_in_a_Box_Mid",
		"cachegroupId": 6,
		"cdnId": 2,
		"cdnName": "CDN-in-a-Box",
		"domainName": "infra.ciab.test",
		"guid": null,
		"hostName": "quest",
		"httpsPort": 443,
		"id": 14,
		"iloIpAddress": "",
		"iloIpGateway": "",
		"iloIpNetmask": "",
		"iloPassword": "",
		"iloUsername": "",
		"lastUpdated": "2020-05-19 16:46:33+00",
		"mgmtIpAddress": "",
		"mgmtIpGateway": "",
		"mgmtIpNetmask": "",
		"offlineReason": "",
		"physLocation": "Apachecon North America 2018",
		"physLocationId": 1,
		"profile": "ATS_MID_TIER_CACHE",
		"profileDesc": "Mid Cache - Apache Traffic Server",
		"profileId": 10,
		"rack": null,
		"revalPending": false,
		"routerHostName": "",
		"routerPortName": "",
		"status": "REPORTED",
		"statusId": 3,
		"tcpPort": 80,
		"type": "MID",
		"typeId": 12,
		"updPending": false,
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
						"serviceAddress": false
					}
				],
				"maxBandwidth": null,
				"monitor": true,
				"mtu": 1500,
				"name": "bond0"
			}
		]
	}}

``DELETE``
==========
Allow user to delete server through api.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

	.. versionchanged:: 3.0
		In older versions of the API, this endpoint did not return a response object. It now returns a representation of the deleted server.

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------+
	| Name |                Description                  |
	+======+=============================================+
	|  ID  | The integral, unique identifier of a server |
	+------+---------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/3.0/servers/14 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:cachegroup:   A string that is the :ref:`name of the Cache Group <cache-group-name>` to which the server belonged
:cachegroupId: An integer that is the :ref:`ID of the Cache Group <cache-group-id>` to which the server belonged
:cdnId:        The integral, unique identifier of the CDN to which the server belonged
:cdnName:      Name of the CDN to which the server belonged
:domainName:   The domain part of the server's :abbr:`FQDN (Fully Qualified Domain Name)`
:guid:         An identifier used to uniquely identify the server

	.. note:: This is a legacy key which only still exists for compatibility reasons - it should always be ``null``

:hostName:     The (short) hostname of the server
:httpsPort:    The port on which the server listened for incoming HTTPS connections/requests
:id:           An integral, unique identifier for this server
:iloIpAddress: The IPv4 address of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:iloIpGateway: The IPv4 gateway address of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:iloIpNetmask: The IPv4 subnet mask of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:iloPassword:  The password of the of the server's :abbr:`ILO (Integrated Lights-Out)` service user\ [#ilo]_ - displays as simply ``******`` if the currently logged-in user does not have the 'admin' or 'operations' :term:`Role(s) <Role>`
:iloUsername:  The user name for the server's :abbr:`ILO (Integrated Lights-Out)` service\ [#ilo]_
:interfaces:   A set of the network interfaces that were in use by the server

	:ipAddresses: A set of objects representing IP Addresses that were assigned to this network interface

		:address:        The actual IP address, including any mask as a CIDR-notation suffix
		:gateway:        Either the IP address of the network gateway for this address, or ``null`` to signify that no such gateway exists
		:serviceAddress: A boolean that describes whether or not the server's main service is available at this IP address. When this property is ``true``, the IP address is referred to as a "service address".

	:maxBandwidth: The maximum healthy bandwidth allowed for this interface. If bandwidth exceeds this limit, Traffic Monitors would have considered the entire server unhealthy - which includes *all* configured network interfaces. If this was ``null``, it has the meaning "no limit". It had no effect if ``monitor`` was not true for this interface.

		.. seealso:: :ref:`health-proto`

	:monitor: A boolean which describes whether or not this interface should have been monitored by Traffic Monitor for statistics and health consideration
	:mtu:     The :abbr:`MTU (Maximum Transmission Unit)` of this interface. If it is ``null``, it may be assumed that the information was either not available or not applicable for this interface.
	:name:    The name of the interface. It is the same as the network interface's device name on the server, e.g. ``eth0``.

:lastUpdated:   The date and time at which this server description was last modified
:mgmtIpAddress: The IPv4 address of some network interface on the server that was used for 'management'

	.. deprecated:: 3.0
		This field is deprecated and will be removed in a future API version. Operators should migrate this data into the ``interfaces`` property of the server.

:mgmtIpGateway: The IPv4 address of a gateway used by some network interface on the server that was used for 'management'

	.. deprecated:: 3.0
		This field is deprecated and will be removed in a future API version. Operators should migrate this data into the ``interfaces`` property of the server.

:mgmtIpNetmask: The IPv4 subnet mask used by some network interface on the server that was used for 'management'

	.. deprecated:: 3.0
		This field is deprecated and will be removed in a future API version. Operators should migrate this data into the ``interfaces`` property of the server.

:offlineReason:  A user-entered reason why the server was in ADMIN_DOWN or OFFLINE status
:physLocation:   The name of the physical location where the server resided
:physLocationId: An integral, unique identifier for the physical location where the server resided
:profile:        The :ref:`profile-name` of the :term:`Profile` which was used by this server
:profileDesc:    A :ref:`profile-description` of the :term:`Profile` which was used by this server
:profileId:      The :ref:`profile-id` the :term:`Profile` which was used by this server
:revalPending:   A boolean value which, if ``true`` indicates that this server had pending content invalidation/revalidation
:rack:           A string indicating "server rack" location
:routerHostName: The human-readable name of the router responsible for reaching this server
:routerPortName: The human-readable name of the port used by the router responsible for reaching this server
:status:         The :term:`Status` of the server

	.. seealso:: :ref:`health-proto`

:statusId: The integral, unique identifier of the status of this server

	.. seealso:: :ref:`health-proto`

:tcpPort: The port on which this server listened for incoming TCP connections

	.. note:: This is typically thought of as synonymous with "HTTP port", as the port specified by ``httpsPort`` may also be used for incoming TCP connections.

:type:       The name of the :term:`Type` of this server
:typeId:     The integral, unique identifier of the 'type' of this server
:updPending: A boolean value which, if ``true``, indicates that the server had updates of some kind pending, typically to be acted upon by Traffic Ops :term:`ORT`
:xmppId:     A system-generated UUID used to generate a server hashId for use in Traffic Router's consistent hashing algorithm. This value is set when a server is created and cannot be changed afterwards.
:xmppPasswd: The password used in XMPP communications with the server

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Tue, 19 May 2020 17:50:13 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 19 May 2020 16:50:13 GMT
	Content-Length: 568

	{ "alerts": [
		{
			"text": "Server deleted",
			"level": "success"
		}
	],
	"response": {
		"cachegroup": "CDN_in_a_Box_Mid",
		"cachegroupId": 6,
		"cdnId": 2,
		"cdnName": "CDN-in-a-Box",
		"domainName": "infra.ciab.test",
		"guid": null,
		"hostName": "quest",
		"httpsPort": 443,
		"id": 14,
		"iloIpAddress": "",
		"iloIpGateway": "",
		"iloIpNetmask": "",
		"iloPassword": "",
		"iloUsername": "",
		"lastUpdated": "2020-05-19 16:46:33+00",
		"mgmtIpAddress": "",
		"mgmtIpGateway": "",
		"mgmtIpNetmask": "",
		"offlineReason": "",
		"physLocation": "Apachecon North America 2018",
		"physLocationId": 1,
		"profile": "ATS_MID_TIER_CACHE",
		"profileDesc": "Mid Cache - Apache Traffic Server",
		"profileId": 10,
		"rack": null,
		"revalPending": false,
		"routerHostName": "",
		"routerPortName": "",
		"status": "REPORTED",
		"statusId": 3,
		"tcpPort": 80,
		"type": "MID",
		"typeId": 12,
		"updPending": false,
		"xmppId": null,
		"xmppPasswd": null,
		"interfaces": [
			{
				"ipAddresses": [
					{
						"address": "0.0.0.1/24",
						"gateway": "0.0.0.2",
						"serviceAddress": false
					},
					{
						"address": "::1",
						"gateway": "::2",
						"serviceAddress": true
					}
				],
				"maxBandwidth": null,
				"monitor": true,
				"mtu": 1500,
				"name": "bond0"
			}
		]
	}}

.. [#ilo] For more information see the `Wikipedia page on Lights-Out management <https://en.wikipedia.org/wiki/Out-of-band_management>`_\ .
