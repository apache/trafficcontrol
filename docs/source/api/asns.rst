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


.. _to-api-asns:

********
``asns``
********
.. seealso:: `The Autonomous System Wikipedia page <https://en.wikipedia.org/wiki/Autonomous_system_%28Internet%29>` for an explanation of what an ASN actually is.

``GET``
=======
List all Autonomous System Numbers (ASNs).
:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+------------+----------+-----------------------------------------------------------------------------------------------------+
	| Parameter  | Required | Description                                                                                         |
	+============+==========+=====================================================================================================+
	| cachegroup | no       | An integral, unique identifier for a :term:`Cache Group` - only ANSs for this :term:`Cache Group`   |
	|            |          | will be returned.                                                                                   |
	+------------+----------+-----------------------------------------------------------------------------------------------------+
	| orderby    | no       | Choose the ordering of the results - must be the name of one of the fields of the objects in the    |
	|            |          | ``response`` array                                                                                  |
	+------------+----------+-----------------------------------------------------------------------------------------------------+
	| sortOrder  | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")            |
	+------------+----------+-----------------------------------------------------------------------------------------------------+
	| limit      | no       | Choose the maximum number of results to return                                                      |
	+------------+----------+-----------------------------------------------------------------------------------------------------+
	| offset     | no       | The number of results to skip before beginning to return results. Must use in conjunction with      |
	|            |          | limit                                                                                               |
	+------------+----------+-----------------------------------------------------------------------------------------------------+
	| page       | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are        |
	|            |          | ``limit`` long and the first page is 1. If ``offset`` was defined, this query parameter has no      |
	|            |          | effect. ``limit`` must be defined to make use of ``page``.                                          |
	+------------+----------+-----------------------------------------------------------------------------------------------------+

Response Structure
------------------
:lastUpdated:  The Time / Date this server entry was last updated in ISO format
:id:           An integer which uniquely identifies the ASN
:asn:          Autonomous System Numbers per APNIC for identifying a service provider
:cachegroup:   Related Cache Group name
:cachegroupId: Related Cache Group ID

.. versionchanged:: 1.2
	Used to contain the array in the ``response.asns`` object, changed so that ``response`` is an actual array

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: 2zeWYI/dGyCzi0ZUWXuuycLFPyL9M5nDJchC7nJMQPW3cwXTaTwf0qI3mP3G1ArZlJTk/ju6/jbUVCNcVIXX1Q==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 01 Nov 2018 18:56:38 GMT
	Content-Length: 129

	{ "response": [
		{
			"asn": 1,
			"cachegroup": "TRAFFIC_ANALYTICS",
			"cachegroupId": 1,
			"id": 1,
			"lastUpdated": "2018-11-01 18:55:39+00"
		}
	]}


``POST``
========
Creates a new Autonomous System Number (ASN).

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type: Object

Request Structure
-----------------
:asn:          The value of the new ASN
:cachegroupId: The integral, unique identifier of a Cache Group to which this ASN will be assigned
:cachegroup:   An optional field which, if present, specifies the name of a Cache Group to which this ASN will be assigned

.. note:: While this endpoint accepts the ``cachegroup`` field, sending this in the request payload has no effect except that the response will (erroneously) name the Cache Group to which the ASN was assigned. Any subsequent requests will reveal that, in fact, the Cache Group name is set by the ``cachegroupId`` field.

.. code-block:: http
	:caption: Request Example

	POST /api/1.1/asns HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 60
	Content-Type: application/x-www-form-urlencoded

	{"asn": 1, "cachegroupId": 1}

Response Structure
------------------
:lastUpdated:  The Time / Date this server entry was last updated in ISO format
:id:           An integer which uniquely identifies the ASN
:asn:          Autonomous System Numbers per APNIC for identifying a service provider
:cachegroup:   Related Cache Group name
:cachegroupId: Related Cache Group ID

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: DnM8HexH7LFkVNH8UYFe6uBQ445Ic8lRLDlOSDIuo4gjokMafxh5Ebr+CsSixNt//OxP0hoWZ+DKymSC5Hdi9Q==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 01 Nov 2018 18:57:08 GMT
	Content-Length: 175

	{ "alerts": [
		{
			"text": "asn was created.",
			"level": "success"
		}
	],
	"response": {
		"asn": 1,
		"cachegroup": "TRAFFIC_ANALYTICS",
		"cachegroupId": 1,
		"id": 2,
		"lastUpdated": "2018-11-01 18:57:08+00"
	}}
