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

.. _to-api-v3-cdns-name-health:

************************
``cdns/{{name}}/health``
************************

``GET``
=======
Retrieves the health of all :term:`Cache Groups` for a given CDN.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-------------------------------------------------------+
	| Name | Description                                           |
	+======+=======================================================+
	| name | The name of the CDN for which health will be reported |
	+------+-------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/3.0/cdns/CDN-in-a-Box/health HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...

Response Structure
------------------
:cachegroups:  An array of objects describing the health of each :term:`Cache Group`

	:name:    A string that is the :ref:`Cache Group's Name <cache-group-name>`
	:offline: The number of OFFLINE :term:`cache servers` in the :term:`Cache Group`
	:online:  The number of ONLINE :term:`cache servers` in the :term:`Cache Group`

:totalOffline: Total number of OFFLINE :term:`cache servers` across all :term:`Cache Groups` which are assigned to the CDN defined by the ``name`` request path parameter
:totalOnline:  Total number of ONLINE :term:`cache servers` across all :term:`Cache Groups` which are assigned to the CDN defined by the ``name`` request path parameter

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Encoding: gzip
	Content-Length: 108
	Content-Type: application/json
	Date: Tue, 03 Dec 2019 21:33:59 GMT
	X-Server-Name: traffic_ops_golang/
	Set-Cookie: mojolicious=...; expires=Wed, 04 Dec 2019 01:33:59 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: KpXViXeAgch58ueQqdyU8NuINBw1EUedE6Rv2ewcLUajJp6kowdbVynpwW7XiSvAyHdtClIOuT3OkhIimghzSA==

	{ "response": {
		"totalOffline": 0,
		"totalOnline": 1,
		"cachegroups": [
			{
				"offline": 0,
				"name": "CDN_in_a_Box_Edge",
				"online": 1
			}
		]
	}}
