..
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

.. _to-api-isos:

********
``isos``
********

``POST``
========
Generates an ISO from the requested ISO source.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: ISO:GENERATE, ISO:READ
:Response Type:  undefined - ISO image as a streaming download

Request Structure
-----------------
:dhcp: A string that specifies whether the generated system image will use DHCP IP address leasing; one of:

	yes
		DHCP will be used, and other network configuration keys need not be present in the request (and are ignored if they are)
	no
		DHCP will not be used, and the desired network configuration **must** be specified manually in the request body

:disk:          An optional string that names the block device (under ``/dev/``) used for the boot media, e.g. "sda"
:domainName:    The domain part of the system image's Fully Qualified Domain Name (FQDN)
:hostName:      The host name part of the system image's FQDN
:interfaceMtu:  A number that specifies the Maximum Transmission Unit (MTU) for the system image's network interface card - the only valid values of which I'm aware are 1500 or 9000, and this should almost always just be 1500
:interfaceName: An optional string naming the network interface to be used by the generated system image e.g. "bond0", "eth0", etc. If the special name "bond0" is used, an :abbr:`LACP (Link Aggregation Control Protocol)` binding configuration will be created and included in the system image

	.. seealso:: `The Link Aggregation Wikipedia page <https://en.wikipedia.org/wiki/Link_aggregation>`_\ .

:ip6Address:   An optional string containing the IPv6 address of the generated system image
:ip6Gateway:   An optional string specifying the IPv6 address of the generated system image's network gateway - this will be ignored if ``ipGateway`` is specified
:ipAddress:    An optional\ [1]_ string containing the IP address of the generated system image
:ipGateway:    An optional\ [1]_ string specifying the IP address of the generated system image's network gateway
:ipNetmask:    An optional\ [1]_ string specifying the subnet mask of the generated system image
:osversionDir: The name of the directory containing the ISO source

	.. seealso:: :ref:`to-api-osversions`

:rootPass: The password used by the generated system image's ``root`` user

.. code-block:: http
	:caption: Request Example

	POST /api/5.0/isos HTTP/1.1
	Host: some.trafficops.host
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 334
	Content-Type: application/json

	{
		"osversionDir": "centos72",
		"hostName": "test",
		"domainName": "quest",
		"rootPass": "twelve",
		"dhcp": "no",
		"interfaceMtu": 1500,
		"ipAddress": "1.3.3.7",
		"ipNetmask": "255.255.255.255",
		"ipGateway": "8.0.0.8",
		"ip6Address": "1::3:3:7",
		"ip6Gateway": "8::8",
		"interfaceName": "eth0",
		"disk": "hda"
	}

.. [1] This optional key is required if and only if ``dhcp`` is "no".

Response Structure
------------------
ISO image as a streaming download.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Connection: keep-alive
	Content-Disposition: attachment; filename="test-centos72_centos72-netinstall.iso"
	Content-Encoding: gzip
	Content-Type: application/download
	Date: Wed, 05 Feb 2020 21:59:15 GMT
	Set-Cookie: mojolicious=...; Path=/; Expires=Wed, 05 Feb 2020 22:59:11 GMT; Max-Age=3600; HttpOnly
	Transfer-Encoding: chunked
	Whole-Content-sha512: sLSVQGrLCQ4hGQhv2reragQHWNi2aKMcz2c/HMAH45tLcZ1LenPyOzWRcRfHUNbV4PEEKOoiTfwE2HlA+WtRIQ==
	X-Server-Name: traffic_ops_golang/
