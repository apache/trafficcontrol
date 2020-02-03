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

.. _to-api-types-id:

****************
``types/{{ID}}``
****************

``GET``
=======
:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-------------------------------------------------------------+
	| Name | Description                                                 |
	+======+=============================================================+
	|  ID  | The integral, unique identifier of the type being inspected |
	+------+-------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/types/48 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:description: A short description of this type
:id:          An integral, unique identifier for this type
:lastUpdated: The date and time at which this type was last updated, in ISO format
:name:        The name of this type
:useInTable:  The name of the Traffic Ops database table that contains objects which are grouped, identified, or described by this type

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: EH8jo8OrCu79Tz9xpgT3YRyKJ/p2NcTmbS3huwtqRByHz9H6qZLQjA59RIPaVSq3ZxsU6QhTaox5nBkQ9LPSAA==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 12 Dec 2018 23:50:13 GMT
	Content-Length: 168

	{ "response": [
		{
			"id": 48,
			"lastUpdated": "2018-12-12 16:26:41+00",
			"name": "TC_LOC",
			"description": "Location for Traffic Control Component Servers",
			"useInTable": "cachegroup"
		}
	]}
