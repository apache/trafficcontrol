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

.. _to-api-v2-servers-id:

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

:cachegroupId:     An integer that is the :ref:`ID of the Cache Group <cache-group-id>` to which the server shall belong
:cdnId:            The integral, unique identifier of the CDN to which the server shall belong
:configUpdateTime: The last time an update was requested for this server.

    .. note:: To maintain backwards compatibility, the ``updPending`` boolean flag will trump this value. However, it is advised to no longer use the ``upd_pending`` flag and is preferred to use this timestamp instead. 

:domainName:   The domain part of the server's :abbr:`FQDN (Fully Qualified Domain Name)`
:hostName:     The (short) hostname of the server
:httpsPort:    An optional port number on which the server listens for incoming HTTPS connections/requests
:iloIpAddress: An optional IPv4 address of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [1]_
:iloIpGateway: An optional IPv4 gateway address of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [1]_
:iloIpNetmask: An optional IPv4 subnet mask of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [1]_
:iloPassword:  An optional string containing the password of the of the server's :abbr:`ILO (Integrated Lights-Out)` service user\ [1]_ - displays as simply ``******`` if the currently logged-in user does not have the 'admin' or 'operations' :abbr:`Role(s) <Role>`
:iloUsername:  An optional string containing the user name for the server's :abbr:`ILO (Integrated Lights-Out)` service\ [1]_
:interfaceMtu: The :abbr:`MTU (Maximum Transmission Unit)` configured on ``interfaceName``

	.. note:: In virtually all cases this ought to be 1500. Further note that the only acceptable values are 1500 and 9000.

:interfaceName:   The name of the primary network interface used by the server
:ip6Address:      An optional IPv6 address and subnet mask of ``interfaceName``
:ip6IsService:    An optional boolean value which if ``true`` indicates that the IPv6 address will be used for routing content.  Defaults to ``true``.
:ip6Gateway:      An optional IPv6 address of the gateway used by ``interfaceName``
:ipAddress:       The IPv4 address of ``interfaceName``
:ipIsService:     An optional boolean value which if ``true`` indicates that the IPv4 address will be used for routing content.  Defaults to ``true``.
:ipGateway:       The IPv4 address of the gateway used by ``interfaceName``
:ipNetmask:       The IPv4 subnet mask used by ``interfaceName``
:mgmtIpAddress:   An optional IPv4 address of some network interface on the server used for 'management'
:mgmtIpGateway:   An optional IPv4 address of a gateway used by some network interface on the server used for 'management'
:mgmtIpNetmask:   An optional IPv4 subnet mask used by some network interface on the server used for 'management'
:physLocationId:  An integral, unique identifier for the physical location where the server resides
:profileId:       The :ref:`profile-id` the :term:`Profile` that shall be used by this server
:revalPending:    A boolean value which, if ``true`` indicates that this server has pending content invalidation/revalidation
:revalUpdateTime: The last time a content invalidation/revalidation request was submitted for this server. This field defaults to standard epoch

    .. note:: To maintain backwards compatibility, the ``revalPending`` boolean flag will trump this value. However, it is advised to no longer use the ``revalPending`` flag and is preferred to use this timestamp instead.

:rack:           An optional string indicating "server rack" location
:routerHostName: An optional string containing the human-readable name of the router responsible for reaching this server
:routerPortName: An optional string containing the human-readable name of the port used by the router responsible for reaching this server
:statusId:       The integral, unique identifier of the status of this server

	.. seealso:: :ref:`health-proto`

:tcpPort: An optional port number on which this server listens for incoming TCP connections

	.. note:: This is typically thought of as synonymous with "HTTP port", as the port specified by ``httpsPort`` may also be used for incoming TCP connections.

:typeId:     The integral, unique identifier of the 'type' of this server
:updPending: A boolean value which, if ``true``, indicates that the server has updates of some kind pending, typically to be acted upon by Traffic Ops ORT
:xmppId:     A system-generated UUID used to generate a server hashId for use in Traffic Router's consistent hashing algorithm. This value is set when a server is created and cannot be changed afterwards.
:xmppPasswd: An optional password used in XMPP communications with the server

.. code-block:: http
	:caption: Request Example

	PUT /api/2.0/servers/13 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 599
	Content-Type: application/json

	{
		"cachegroupId": 6,
		"cdnId": 2,
		"configUpdateTime": "2022-02-18T13:52:47.129174-07:00",
		"domainName": "infra.ciab.test",
		"hostName": "quest",
		"httpsPort": 443,
		"iloIpAddress": "",
		"iloIpGateway": "",
		"iloIpNetmask": "",
		"iloPassword": "",
		"iloUsername": "",
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
		"ipIsService": true,
		"ip6IsService": true
	}

Response Structure
------------------
:cachegroup:       A string that is the :ref:`name of the Cache Group <cache-group-name>` to which the server belongs
:cachegroupId:     An integer that is the :ref:`ID of the Cache Group <cache-group-id>` to which the server belongs
:cdnId:            The integral, unique identifier of the CDN to which the server belongs
:cdnName:          Name of the CDN to which the server belongs
:configUpdateTime: The last time an update was requested for this server. This field defaults to standard epoch
:configApplyTime:  The last time an update was applied for this server. This field defaults to standard epoch
:domainName:       The domain part of the server's :abbr:`FQDN (Fully Qualified Domain Name)`
:guid:             An identifier used to uniquely identify the server

	.. note:: This is a legacy key which only still exists for compatibility reasons - it should always be ``null``

:hostName:       The (short) hostname of the server
:httpsPort:      The port on which the server listens for incoming HTTPS connections/requests
:id:             An integral, unique identifier for this server
:iloIpAddress:   The IPv4 address of the server's Integrated Lights-Out (ILO) service\ [1]_
:iloIpGateway:   The IPv4 gateway address of the server's ILO service\ [1]_
:iloIpNetmask:   The IPv4 subnet mask of the server's ILO service\ [1]_
:iloPassword:    The password of the of the server's ILO service user\ [1]_ - displays as simply ``******`` if the currently logged-in user does not have the 'admin' or 'operations' role(s)
:iloUsername:    The user name for the server's ILO service\ [1]_
:interfaceMtu:   The Maximum Transmission Unit (MTU) to configured on ``interfaceName``
:interfaceName:  The name of the primary network interface used by the server
:ip6Address:     The IPv6 address and subnet mask of ``interfaceName``
:ip6IsService:   A boolean value which if ``true`` indicates that the IPv6 address will be used for routing content.
:ip6Gateway:     The IPv6 address of the gateway used by ``interfaceName``
:ipAddress:      The IPv4 address of ``interfaceName``
:ipIsService:    A boolean value which if ``true`` indicates that the IPv4 address will be used for routing content.
:ipGateway:      The IPv4 address of the gateway used by ``interfaceName``
:ipNetmask:      The IPv4 subnet mask used by ``interfaceName``
:lastUpdated:    The date and time at which this server description was last modified
:mgmtIpAddress:  The IPv4 address of some network interface on the server used for 'management'
:mgmtIpGateway:  The IPv4 address of a gateway used by some network interface on the server used for 'management'
:mgmtIpNetmask:  The IPv4 subnet mask used by some network interface on the server used for 'management'
:offlineReason:  A user-entered reason why the server is in ADMIN_DOWN or OFFLINE status
:physLocation:   The name of the physical location where the server resides
:physLocationId: An integral, unique identifier for the physical location where the server resides
:profile:        The :ref:`profile-name` of the :term:`Profile` used by this server
:profileDesc:    A :ref:`profile-description` of the :term:`Profile` used by this server
:profileId:      The :ref:`profile-id` the :term:`Profile` used by this server
:revalPending:   A boolean value which, if ``true`` indicates that this server has pending content invalidation/revalidation

    .. note:: While not officially deprecated, this is based on the values corresponding to ``revalUpdateTime`` and ``revalApplyTime``. It is preferred to use the timestamp fields going forward as this will likely be deprecated in the future.

:revalUpdateTime: The last time a content invalidation/revalidation request was submitted for this server. This field defaults to standard epoch
:revalApplyTime:  The last time a content invalidation/revalidation request was applied by this server. This field defaults to standard epoch
:rack:            A string indicating "server rack" location
:routerHostName:  The human-readable name of the router responsible for reaching this server
:routerPortName:  The human-readable name of the port used by the router responsible for reaching this server
:status:          The status of the server

	.. seealso:: :ref:`health-proto`

:statusId: The integral, unique identifier of the status of this server

	.. seealso:: :ref:`health-proto`

:tcpPort: The port on which this server listens for incoming TCP connections

	.. note:: This is typically thought of as synonymous with "HTTP port", as the port specified by ``httpsPort`` may also be used for incoming TCP connections.

:type:       The name of the 'type' of this server
:typeId:     The integral, unique identifier of the 'type' of this server
:updPending: A boolean value which, if ``true``, indicates that the server has updates of some kind pending, typically to be acted upon by Traffic Control Cache Config (T3C, formerly ORT)

    .. note:: While not officially deprecated, this is based on the values corresponding to ``configUpdateTime`` and ``configApplyTime``. It is preferred to use the timestamp fields going forward as this will likely be deprecated in the future.

:xmppId:     A system-generated UUID used to generate a server hashId for use in Traffic Router's consistent hashing algorithm. This value is set when a server is created and cannot be changed afterwards.
:xmppPasswd: The password used in XMPP communications with the server

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: 9lGAMCCC9I/bOpuBSyf3ACffjHeRuXCTuxrA/oU78uYzW5FeFTq5PHSSnsnqKG5E0vWg0Rko0CwguGeNc9IT0w==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 10 Dec 2018 17:58:57 GMT
	Content-Length: 848

	{ "alerts": [
		{
			"text": "server was updated.",
			"level": "success"
		}
	],
	"response": {
		"cachegroup": null,
		"cachegroupId": 6,
		"cdnId": 2,
		"cdnName": null,
		"configUpdateTime": "2022-02-28T15:44:15.895145-07:00",
		"configApplyTime": "2022-02-18T13:52:47.129174-07:00",
		"domainName": "infra.ciab.test",
		"guid": null,
		"hostName": "quest",
		"httpsPort": 443,
		"id": 13,
		"iloIpAddress": "",
		"iloIpGateway": "",
		"iloIpNetmask": "",
		"iloPassword": "",
		"iloUsername": "",
		"interfaceMtu": 1500,
		"interfaceName": "eth0",
		"ip6Address": "::1",
		"ip6Gateway": "::2",
		"ipAddress": "0.0.0.1",
		"ipGateway": "0.0.0.2",
		"ipNetmask": "255.255.255.0",
		"lastUpdated": "2018-12-10 17:58:57+00",
		"mgmtIpAddress": "",
		"mgmtIpGateway": "",
		"mgmtIpNetmask": "",
		"offlineReason": "",
		"physLocation": null,
		"physLocationId": 1,
		"profile": null,
		"profileDesc": null,
		"profileId": 10,
		"rack": null,
		"revalPending": null,
		"revalUpdateTime": "1969-12-31T17:00:00-07:00",
		"revalApplyTime": "1969-12-31T17:00:00-07:00",
		"routerHostName": "",
		"routerPortName": "",
		"status": null,
		"statusId": 3,
		"tcpPort": 80,
		"type": "",
		"typeId": 12,
		"updPending": true,
		"xmppId": null,
		"xmppPasswd": null,
		"ipIsService": true,
		"ip6IsService": true
	}}

``DELETE``
==========
Allow user to delete server through api.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

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

	DELETE /api/2.0/servers/13 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: JZdjKJYWN9w9NF6VE/rVkGUqecycKB2ABkkI4LNDmgpJLwu53bRHAA+4uWrow0zuba/4MSEhHKshutziypSxPg==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 10 Dec 2018 18:23:21 GMT
	Content-Length: 61

	{ "alerts": [
		{
			"text": "server was deleted.",
			"level": "success"
		}
	]}

.. [1] For more information see the `Wikipedia page on Lights-Out management <https://en.wikipedia.org/wiki/Out-of-band_management>`_\ .
