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

****
asns
****
List all ASNS

``GET``
=======
:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+----------------+----------+-----------------------------------------------------------------------------------------------------+
	| Parameter      | Required |                                 Description                                                         |
	+================+==========+=====================================================================================================+
	|   cachegroup   | no       | An integral, unique identifier for a Cache Group - only ANSs for this Cache Group will be returned. |
	+----------------+----------+-----------------------------------------------------------------------------------------------------+

Response Structure
------------------
:lastUpdated:  The Time / Date this server entry was last updated in ISO format
:id:           An integer which uniquely identifies the ASN
:asn:          Autonomous System Numbers per APNIC for identifying a service provider
:cachegroup:   Related Cache Group name
:cachegroupId: Related Cache Group ID

.. code-block::
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

	{ "response": {
	"asns": [
		{
			"asn": 1,
			"cachegroup": "TRAFFIC_ANALYTICS",
			"cachegroupId": 1,
			"id": 1,
			"lastUpdated": "2018-11-01 18:55:39+00"
		}
	]}}


.. versionchanged:: 1.2
	Used to contain the array in the ``response.asns`` object, changed so that ``response`` is an actual array

``POST``
========
Creates a new ASN

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type: Object

Request Structure
-----------------
.. table:: Request Data Parameters

	+-------------------+---------+----------+--------------------------------------------------------------+
	|    Parameter      |  Type   | Required |                   Description                                |
	+===================+=========+==========+==============================================================+
	| ``asn``           | integer | yes      | ASN                                                          |
	+-------------------+---------+----------+--------------------------------------------------------------+
	| ``cachegroupId``  | integer | yes      | The ID of a Cache Group to which this ASN will be assigned   |
	+-------------------+---------+----------+--------------------------------------------------------------+
	| ``cachegroup``    | string  | no       | The name of a Cache Group to which this ASN will be assigned |
	+-------------------+---------+----------+--------------------------------------------------------------+

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
