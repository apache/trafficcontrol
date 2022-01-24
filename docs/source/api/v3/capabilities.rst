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

.. _to-api-v3-capabilities:

****************
``capabilities``
****************
.. deprecated:: 3.1

``GET``
=======
Get all capabilities.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------+----------+---------------------------------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                                                         |
	+===========+==========+=====================================================================================================================+
	| name      | no       | Return only the capability that has this name                                                                       |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------------+
	| orderby   | no       | Choose the ordering of the results - must be the name of one of the fields of the objects in the ``response`` array |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------------+
	| sortOrder | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")                            |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------------+
	| limit     | no       | Choose the maximum number of results to return                                                                      |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------------+
	| offset    | no       | The number of results to skip before beginning to return results. Must use in conjunction with ``limit``            |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------------+
	| page      | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long and the |
	|           |          | first page is 1. If ``offset`` was defined, this query parameter has no effect. ``limit`` must be defined to make   |
	|           |          | use of ``page``.                                                                                                    |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------------+


.. code-block:: http
	:caption: Request Example

	GET /api/3.0/capabilities?name=test HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:name:        Name of the capability
:description: Describes the permissions covered by the capability.
:lastUpdated: Date and time of the last update made to this capability, in :ref:`non-rfc-datetime`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Vary: Accept-Encoding
	Transfer-Encoding: chunked
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: c18+GtX2ZI8PoCSwuAzBhl+6w3vDpKQTa/cDJC0WHxdpguOL378KBxGWW5PCSyZfJUb7wPyOL5qKMn6NNTufhg==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 15 Aug 2019 17:20:20 GMT
	Content-Length: 161

	{ "response": [
		{
			"description": "This is only a test. If this were a real capability, it might do something",
			"lastUpdated": "2019-08-15 17:18:03+00",
			"name": "test"
		}
	]}
