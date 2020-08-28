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

.. _to-api-v1-divisions-name-regions:

******************************
``divisions/{{name}}/regions``
******************************
.. deprecated:: ATCv4

``POST``
========
Creates a new :term:`Region` within the specified :term:`Division`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------+
	| Name | Description                                                |
	+======+============================================================+
	| name | The name of the division in which to create the new region |
	+------+------------------------------------------------------------+

:name: The name of the new region

.. code-block:: http
	:caption: Request Example

	POST /api/1.4/divisions/England/regions HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 27
	Content-Type: application/json

	{
		"name": "Greater_London"
	}

Response Structure
------------------
:divisionName: The name of the division which contains the new region
:divisionId:   The integral, unique identifier of the division which contains the new region
:id:           An integral, unique identifier for this region
:name:         The region name

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Thu, 06 Dec 2018 00:03:36 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: HlzhY41lFBRlLe5D0XN1w+LbU/N1WD+JXX0tzMwDFqI4VmpBLaAqzUaJqRpQdJnO2u7Z2E0b6QVOgeGRPpyUzg==
	Content-Length: 84

	{ "response": {
		"divisionName": "England",
		"divsionId": 3,
		"name": "Greater_London",
		"id": 3
	},
	"alerts": [
		{
			"level": "warning",
			"text": "This endpoint is deprecated, please use 'POST /regions' instead"
		}
	]}
