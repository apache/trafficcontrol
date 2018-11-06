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

.. _to-api-deliveryservices-id-unassigned_servers:

******************************************
deliveryservices/{{ID}}/unassigned_servers
******************************************

.. caution:: This route does not appear to work properly, and its use is strongly discouraged! Also note that the documentation here is not being updated as a result of this, and may contain out-of-date and/or erroneous information.

``GET``
=======
Retrieves properties of Edge-tier servers not assigned to a Delivery Service.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [1]_
:Response Type:  Array

.. table:: Request Route Parameters

	+-----------------+----------+---------------------------------------------------+
	| Name            | Required | Description                                       |
	+=================+==========+===================================================+
	| ``id``          | yes      | Delivery service ID.                              |
	+-----------------+----------+---------------------------------------------------+

.. table:: Response Properties

	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	|     Parameter      |  Type  |                                                Description                                                 |
	+====================+========+============================================================================================================+
	| ``cachegroup``     | string | The cache group name                                                                                       |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``cachegroupId``   | string | The cache group id.                                                                                        |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``cdnId``          | string | Id of the CDN to which the server belongs to.                                                              |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``cdnName``        | string | Name of the CDN to which the server belongs to.                                                            |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``domainName``     | string | The domain name part of the FQDN of the cache.                                                             |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``guid``           | string | An identifier used to uniquely identify the server.                                                        |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``hostName``       | string | The host name part of the cache.                                                                           |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``httpsPort``      | string | The HTTPS port on which the main application listens (443 in most cases).                                  |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``id``             | string | The server id (database row number).                                                                       |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``iloIpAddress``   | string | The IPv4 address of the lights-out-management port.                                                        |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``iloIpGateway``   | string | The IPv4 gateway address of the lights-out-management port.                                                |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``iloIpNetmask``   | string | The IPv4 netmask of the lights-out-management port.                                                        |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``iloPassword``    | string | The password of the of the lights-out-management user (displays as ****** unless you are an 'admin' user). |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``iloUsername``    | string | The user name for lights-out-management.                                                                   |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``interfaceMtu``   | string | The Maximum Transmission Unit (MTU) to configure for ``interfaceName``.                                    |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``interfaceName``  | string | The network interface name used for serving traffic.                                                       |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``ip6Address``     | string | The IPv6 address/netmask for ``interfaceName``.                                                            |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``ip6Gateway``     | string | The IPv6 gateway for ``interfaceName``.                                                                    |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``ipAddress``      | string | The IPv4 address for ``interfaceName``.                                                                    |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``ipGateway``      | string | The IPv4 gateway for ``interfaceName``.                                                                    |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``ipNetmask``      | string | The IPv4 netmask for ``interfaceName``.                                                                    |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``lastUpdated``    | string | The Time and Date for the last update for this server.                                                     |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``mgmtIpAddress``  | string | The IPv4 address of the management port (optional).                                                        |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``mgmtIpGateway``  | string | The IPv4 gateway of the management port (optional).                                                        |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``mgmtIpNetmask``  | string | The IPv4 netmask of the management port (optional).                                                        |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``offlineReason``  | string | A user-entered reason why the server is in ADMIN_DOWN or OFFLINE status.                                   |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``physLocation``   | string | The physical location name (see :ref:`to-api-v11-phys-loc`).                                               |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``physLocationId`` | string | The physical location id (see :ref:`to-api-v11-phys-loc`).                                                 |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``profile``        | string | The assigned profile name (see :ref:`to-api-v11-profile`).                                                 |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``profileDesc``    | string | The assigned profile description (see :ref:`to-api-v11-profile`).                                          |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``profileId``      | string | The assigned profile Id (see :ref:`to-api-v11-profile`).                                                   |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``rack``           | string | A string indicating rack location.                                                                         |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``routerHostName`` | string | The human readable name of the router.                                                                     |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``routerPortName`` | string | The human readable name of the router port.                                                                |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``status``         | string | The Status string (See :ref:`to-api-v11-status`).                                                          |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``statusId``       | string | The Status id (See :ref:`to-api-v11-status`).                                                              |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``tcpPort``        | string | The default TCP port on which the main application listens (80 for a cache in most cases).                 |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``type``           | string | The name of the type of this server (see :ref:`to-api-v11-type`).                                          |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``typeId``         | string | The id of the type of this server (see :ref:`to-api-v11-type`).                                            |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+
	| ``updPending``     |  bool  |                                                                                                            |
	+--------------------+--------+------------------------------------------------------------------------------------------------------------+

.. code-block:: json
	:caption: Response Example

	 {
			"response": [
					{
							"cachegroup": "us-il-chicago",
							"cachegroupId": "3",
							"cdnId": "3",
							"cdnName": "CDN-1",
							"domainName": "chi.kabletown.net",
							"guid": null,
							"hostName": "atsec-chi-00",
							"id": "19",
							"iloIpAddress": "172.16.2.6",
							"iloIpGateway": "172.16.2.1",
							"iloIpNetmask": "255.255.255.0",
							"iloPassword": "********",
							"iloUsername": "",
							"interfaceMtu": "9000",
							"interfaceName": "bond0",
							"ip6Address": "2033:D0D0:3300::2:2/64",
							"ip6Gateway": "2033:D0D0:3300::2:1",
							"ipAddress": "10.10.2.2",
							"ipGateway": "10.10.2.1",
							"ipNetmask": "255.255.255.0",
							"lastUpdated": "2015-03-08 15:57:32",
							"mgmtIpAddress": "",
							"mgmtIpGateway": "",
							"mgmtIpNetmask": "",
							"offlineReason": "N/A",
							"physLocation": "plocation-chi-1",
							"physLocationId": "9",
							"profile": "EDGE1_CDN1_421_SSL",
							"profileDesc": "EDGE1_CDN1_421_SSL profile",
							"profileId": "12",
							"rack": "RR 119.02",
							"routerHostName": "rtr-chi.kabletown.net",
							"routerPortName": "2",
							"status": "ONLINE",
							"statusId": "6",
							"tcpPort": "80",
							"httpsPort": "443",
							"type": "EDGE",
							"typeId": "3",
							"updPending": false
					},
				]
		}

.. [1] Users with the roles "admin" and/or "operations" will be able to see servers not assigned to *any* given Delivery Service, whereas any other user will only be able to see the servers not assigned to Delivery Services their Tenant is allowed to see.

