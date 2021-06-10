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

.. _to-api-v3-deliveryservices-id-servers-eligible:

********************************************
``deliveryservices/{{ID}}/servers/eligible``
********************************************

.. caution:: This endpoint may not work as advertised, and its use is therefore discouraged!

``GET``
=======
Retrieves properties of :term:`Edge-tier cache servers` eligible for assignment to a particular :term:`Delivery Service`. Eligibility is determined based on the following properties:

- The name of the server's :term:`Type` must match one of the glob patterns ``EDGE*``, ``ORG*``
- The server and :term:`Delivery Service` must belong to the same CDN
- If the :term:`Delivery Service` has :ref:`ds-required-capabilities`, an :term:`Edge-tier cache server` must have all of those defined capabilities

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------------------------------------------------------+
	| Name | Description                                                                                 |
	+======+=============================================================================================+
	| ID   | The integral, unique identifier of the Delivery service for which servers will be displayed |
	+------+---------------------------------------------------------------------------------------------+

Response Structure
------------------
:cachegroup:     A string which is the :ref:`Name of the Cache Group <cache-group-name>` to which the server belongs
:cachegroupId:   An integer that is the :ref:`ID of the Cache Group <cache-group-id>` to which the server belongs
:cdnId:          An integral, unique identifier the CDN to which the server belongs
:cdnName:        The name of the CDN to which the server belongs
:domainName:     The domain name part of the :abbr:`FQDN (Fully Qualified Domain Name)` of the server
:guid:           Optionally represents an identifier used to uniquely identify the server
:hostName:       The (short) hostname of the server
:httpsPort:      The port on which the server listens for incoming HTTPS requests - 443 in most cases
:id:             An integral, unique identifier for the server
:iloIpAddress:   The IPv4 address of the lights-out-management port\ [#ilowikipedia]_
:iloIpGateway:   The IPv4 gateway address of the lights-out-management port\ [#ilowikipedia]_
:iloIpNetmask:   The IPv4 subnet mask of the lights-out-management port\ [#ilowikipedia]_
:iloPassword:    The password of the of the lights-out-management user - displays as ``******`` unless the requesting user has the 'admin' role)\ [#ilowikipedia]_
:iloUsername:    The user name for lights-out-management\ [#ilowikipedia]_
:interfaces:     An array of interface and IP address information

	:max_bandwidth:  The maximum allowed bandwidth for this interface to be considered "healthy" by Traffic Monitor. This has no effect if `monitor` is not true. Values are in kb/s. The value `null` means "no limit".
	:monitor:        A boolean indicating if Traffic Monitor should monitor this interface
	:mtu:            The :abbr:`MTU (Maximum Transmission Unit)` to configure for ``interfaceName``

		.. seealso:: `The Wikipedia article on Maximum Transmission Unit <https://en.wikipedia.org/wiki/Maximum_transmission_unit>`_

	:name:           The network interface name used by the server.

	:ipAddresses:    An array of the IP address information for the interface

		:address:       The IPv4 or IPv6 address and subnet mask of the server - applicable for the interface ``name``
		:gateway:       The IPv4 or IPv6 gateway address of the server - applicable for the interface ``name``
		:service_address:  A boolean determining if content will be routed to the IP address

:lastUpdated:    The time and date at which this server was last updated, in :ref:`non-rfc-datetime`
:mgmtIpAddress:  The IPv4 address of the server's management port
:mgmtIpGateway:  The IPv4 gateway of the server's management port
:mgmtIpNetmask:  The IPv4 subnet mask of the server's management port
:offlineReason:  A user-entered reason why the server is in ADMIN_DOWN or OFFLINE status (will be empty if not offline)
:physLocation:   The name of the :term:`Physical Location` at which the server resides
:physLocationId: An integral, unique identifier for the :term:`Physical Location` at which the server resides
:profile:        The :ref:`profile-name` of the :term:`Profile` assigned to this server
:profileDesc:    A :ref:`profile-description` of the :term:`Profile` assigned to this server
:profileId:      The :ref:`profile-id` of the :term:`Profile` assigned to this server
:rack:           A string indicating "rack" location
:routerHostName: The human-readable name of the router
:routerPortName: The human-readable name of the router port
:status:         The Status of the server

	.. seealso:: :ref:`health-proto`

:statusId:       An integral, unique identifier for the status of the server

	.. seealso:: :ref:`health-proto`

:tcpPort:        The default port on which the main application listens for incoming TCP connections - 80 in most cases
:type:           The name of the :term:`Type` of this server
:typeId:         An integral, unique identifier for the :term:`Type` of this server
:updPending:     ``true`` if the server has updates pending, ``false`` otherwise

.. code-block:: json
	:caption: Response Example

	{ "response": [
		{
			"cachegroup": "CDN_in_a_Box_Edge",
			"cachegroupId": 7,
			"cdnId": 2,
			"cdnName": "CDN-in-a-Box",
			"domainName": "infra.ciab.test",
			"guid": null,
			"hostName": "edge",
			"httpsPort": 443,
			"id": 10,
			"iloIpAddress": "",
			"iloIpGateway": "",
			"iloIpNetmask": "",
			"iloPassword": "",
			"iloUsername": "",
			"lastUpdated": "2018-10-30 16:01:12+00",
			"mgmtIpAddress": "",
			"mgmtIpGateway": "",
			"mgmtIpNetmask": "",
			"offlineReason": "",
			"physLocation": "Apachecon North America 2018",
			"physLocationId": 1,
			"profile": "ATS_EDGE_TIER_CACHE",
			"profileDesc": "Edge Cache - Apache Traffic Server",
			"profileId": 9,
			"rack": "",
			"routerHostName": "",
			"routerPortName": "",
			"status": "REPORTED",
			"statusId": 3,
			"tcpPort": 80,
			"type": "EDGE",
			"typeId": 11,
			"updPending": false,
			"interfaces": [{
				"ipAddresses": [
					{
						"address": "172.16.239.100",
						"gateway": "172.16.239.1",
						"service_address": true
					},
					{
						"address": "fc01:9400:1000:8::100",
						"gateway": "fc01:9400:1000:8::1",
						"service_address": true
					}
				],
				"max_bandwidth": 0,
				"monitor": true,
				"mtu": 1500,
				"name": "eth0"
			}]
		}
	]}

.. [#ilowikipedia] See `the Wikipedia article on Out-of-Band Management <https://en.wikipedia.org/wiki/Out-of-band_management>`_ for more information.
