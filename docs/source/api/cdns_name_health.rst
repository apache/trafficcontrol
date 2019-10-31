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

.. _to-api-cdns-name-health:

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

Response Structure
------------------
:cachegroups:  An array of objects describing the health of each :term:`Cache Group`

	:name:    The name of the :term:`Cache Group`
	:offline: The number of OFFLINE caches in the :term:`Cache Group`
	:online:  The number of ONLINE caches in the :term:`Cache Group`

:totalOffline: Total number of OFFLINE caches across all :term:`Cache Groups` which are assigned to the CDN defined by the ``name`` request path parameter
:totalOnline:  Total number of ONLINE caches across all :term:`Cache Groups` which are assigned to the CDN defined by the ``name`` request path parameter

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 14 Nov 2018 21:14:05 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Thu, 15 Nov 2018 01:14:05 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: KpXViXeAgch58ueQqdyU8NuINBw1EUedE6Rv2ewcLUajJp6kowdbVynpwW7XiSvAyHdtClIOuT3OkhIimghzSA==
	Content-Length: 115

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
