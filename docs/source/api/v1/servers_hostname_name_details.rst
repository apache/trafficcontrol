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

.. _to-api-v1-servers-hostname-name-details:

*************************************
``servers/hostname/{{name}}/details``
*************************************
.. deprecated:: ATCv4
	Use the ``GET`` method of :ref:`to-api-servers-details` with the query parameter ``hostName`` instead.

``GET``
=======
Retrieves the details of a server.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+----------------------------------------------------+
	| Name |           Description                              |
	+======+====================================================+
	| name | The (short) hostname of the server being inspected |
	+------+----------------------------------------------------+

Response Structure
------------------
:cachegroup:       A string that is the :ref:`name of the Cache Group <cache-group-name>` to which the server belongs
:cdnName:          Name of the CDN to which the server belongs
:deliveryservices: An array of integral, unique identifiers for :term:`Delivery Services` to which this server belongs
:domainName:       The domain part of the server's :abbr:`FQDN (Fully Qualified Domain Name)`
:guid:             An identifier used to uniquely identify the server

	.. note:: This is a legacy key which only still exists for compatibility reasons - it should always be ``null``

:hostName:         The (short) hostname of the server
:httpsPort:        The port on which the server listens for incoming HTTPS connections/requests
:id:               An integral, unique identifier for this server
:iloIpAddress:     The IPv4 address of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [1]_
:iloIpGateway:     The IPv4 gateway address of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [1]_
:iloIpNetmask:     The IPv4 subnet mask of the server's :abbr:`ILO (Integrated Lights-Out)` service\ [1]_
:iloPassword:      The password of the of the server's :abbr:`ILO (Integrated Lights-Out)` service user\ [1]_ - displays as simply ``******`` if the currently logged-in user does not have the 'admin' or 'operations' :term:`Role(s) <Role>`
:iloUsername:      The user name for the server's :abbr:`ILO (Integrated Lights-Out)` service\ [1]_
:interfaceMtu:     The :abbr:`MTU (Maximum Transmission Unit)` to configured on ``interfaceName``
:interfaceName:    The name of the primary network interface used by the server
:ip6Address:       The IPv6 address and subnet mask of ``interfaceName``
:ip6Gateway:       The IPv6 address of the gateway used by ``interfaceName``
:ipAddress:        The IPv4 address of ``interfaceName``
:ipGateway:        The IPv4 address of the gateway used by ``interfaceName``
:ipNetmask:        The IPv4 subnet mask used by ``interfaceName``
:offlineReason:    A user-entered reason why the server is in ADMIN_DOWN or OFFLINE status
:physLocation:     The name of the physical location where the server resides
:profile:          The :ref:`profile-name` of the :term:`Profile` used by this server
:profileDesc:      A :ref:`profile-description` of the :term:`Profile` used by this server
:rack:             A string indicating "server rack" location
:routerHostName:   The human-readable name of the router responsible for reaching this server
:routerPortName:   The human-readable name of the port used by the router responsible for reaching this server
:status:           The status of the server

	.. seealso:: :ref:`health-proto`

:tcpPort: The port on which this server listens for incoming TCP connections

	.. note:: This is typically thought of as synonymous with "HTTP port", as the port specified by ``httpsPort`` may also be used for incoming TCP connections.

:type:       The name of the 'type' of this server
:xmppId:     An identifier to be used in XMPP communications with the server - in nearly all cases this will be the same as ``hostName``
:xmppPasswd: The password used in XMPP communications with the server

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Mon, 10 Dec 2018 17:11:53 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: ZDeQrG0D7Q3Wy3ZEUT9t21QQ9F9Yc3RR/Qr91n22UniYubdhdKnir3B+LYP5ZKkVg8ByrVPFyx6Nao0iiBTGTQ==
	Content-Length: 800

	{ "response": {
		"profile": "ATS_EDGE_TIER_CACHE",
		"xmppPasswd": "",
		"physLocation": "Apachecon North America 2018",
		"cachegroup": "CDN_in_a_Box_Edge",
		"interfaceName": "eth0",
		"id": 9,
		"tcpPort": 80,
		"httpsPort": 443,
		"ipGateway": "172.16.239.1",
		"ip6Address": "fc01:9400:1000:8::100",
		"xmppId": "edge",
		"mgmtIpNetmask": "",
		"rack": "",
		"mgmtIpGateway": "",
		"deliveryservices": [
			1
		],
		"type": "EDGE",
		"iloIpNetmask": "",
		"domainName": "infra.ciab.test",
		"iloUsername": "",
		"status": "REPORTED",
		"ipAddress": "172.16.239.100",
		"ip6Gateway": "fc01:9400:1000:8::1",
		"iloPassword": "",
		"guid": null,
		"offlineReason": "",
		"routerPortName": "",
		"ipNetmask": "255.255.255.0",
		"mgmtIpAddress": "",
		"interfaceMtu": 1500,
		"iloIpGateway": "",
		"cdnName": "CDN-in-a-Box",
		"hostName": "edge",
		"iloIpAddress": "",
		"profileDesc": "Edge Cache - Apache Traffic Server",
		"routerHostName": ""
	},
	"alerts": [
		{
			"text": "This endpoint is deprecated, please use GET /servers/details with query parameter hostName instead",
			"level": "warning"
		}
	]}

.. [1] For more information see the `Wikipedia page on Lights-Out management <https://en.wikipedia.org/wiki/Out-of-band_management>`_\ .
