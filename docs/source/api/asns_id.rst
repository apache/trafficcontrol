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


.. _to-api-asns-id:

************************
``/api/1.x/asns/{{id}}``
************************
Deal with a specific ASN

``GET``
=======
:Auth. Required: Yes
:Roles Required: None
:Response Type: Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------+----------+---------------------------------------------+
	|   Name    | Required |                Description                  |
	+===========+==========+=============================================+
	|   ``id``  |   yes    | The ID of the desired ASN                   |
	+-----------+----------+---------------------------------------------+


Response Structure
------------------
:asn:          Autonomous System Numbers per APNIC for identifying a service provider
:cachegroup:   Related Cache Group name
:cachegroupId: Related Cache Group ID
:id:           An integer which uniquely identifies the ASN
:lastUpdated:  The Time / Date this server entry was last updated in ISO format

.. code-block:: json
	:caption: Response Example

	{ "response": [
		{
			"lastUpdated": "2012-09-17 21:41:22",
			"id": "28",
			"asn": "7016",
			"cachegroup": "us-pa-pittsburgh",
			"cachegroupId": "13"
		}
	]}

``PUT``
=======
Allows user to edit an ASN.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type: Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+-------------------+----------+------------------------------------------------+
	| Name              | Required |                 Description                    |
	+===================+==========+================================================+
	| ``id``            | yes      | The ID of the desired ASN                      |
	+-------------------+----------+------------------------------------------------+

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
:asn:          Autonomous System Numbers per APNIC for identifying a service provider
:cachegroup:   Related Cache Group name
:cachegroupId: Related Cache Group ID
:id:           An integer which uniquely identifies the ASN
:lastUpdated:  The Time / Date this server entry was last updated in ISO format

.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"text": "asn was updated.",
			"level": "success"
		}
	],
	"response": {
		"asn": 2,
		"cachegroup": "CDN_in_a_Box_Edge",
		"cachegroupId": 6,
		"id": 5,
		"lastUpdated": "2018-10-15 14:53:10+00"
	}}
