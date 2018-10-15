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

.. _to-api-api_capability:

*****************************
``/api/1.x/api_capabilities``
*****************************
Deals with the capabilities that may be associated with API endpoints and methods. These capabilities are assigned to "roles", of which a user may have one or more. Capabilities support "wildcarding" or "globbing" using asterisks to group multiple routes into a single capability

``GET``
=======
Get all API-capability mappings.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+----------------+----------+--------+------------------------------------+
	|    Name        | Required | Type   |         Description                |
	+================+==========+========+====================================+
	| ``capability`` |   no     | string | Capability name.                   |
	+----------------+----------+--------+------------------------------------+

Response Structure
------------------
:capability:  Capability name
:httpMethod:  An HTTP request method, practically one of:

	- ``GET``
	- ``POST``
	- ``PUT``
	- ``PATCH``
	- ``DELETE``

:httpRoute:   The request route for which this capability applies - relative to the Traffic Ops server's URL
:id:          An integer which uniquely identifies this capability
:lastUpdated: The time at which this capability was last updated, in ISO format

.. code-block::json
	:caption: Response Example

	{ "response": [
		{
			"id": "6",
			"httpMethod": "GET",
			"httpRoute": "/api/*/asns",
			"capability": "asn-read",
			"lastUpdated": "2017-04-02 08:22:43"
		},
		{
			"id": "7",
			"httpMethod": "GET",
			"httpRoute": "/api/*/asns/*",
			"capability": "asn-read",
			"lastUpdated": "2017-04-02 08:22:43"
		}
	]}

``POST``
========
Create an API-capability mapping.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"

Request Structure
-----------------
.. table:: Request Data Parameters

	+----------------+----------+--------+--------------------------------------------------+
	|    Name        | Required | Type   |                Description                       |
	+================+==========+========+==================================================+
	| ``httpMethod`` | yes      | string | One of: 'GET', 'POST', 'PUT', 'PATCH', 'DELETE'. |
	+----------------+----------+--------+--------------------------------------------------+
	| ``httpRoute``  | yes      | string | API route.                                       |
	+----------------+----------+--------+--------------------------------------------------+
	| ``capability`` | yes      | string | Capability name                                  |
	+----------------+----------+--------+--------------------------------------------------+

Response Structure
------------------
:capability:  Capability name
:httpMethod:  An HTTP request method, practically one of:

	- ``GET``
	- ``POST``
	- ``PUT``
	- ``PATCH``
	- ``DELETE``

:httpRoute:   The request route for which this capability applies - relative to the Traffic Ops server's URL
:id:          An integer which uniquely identifies this capability
:lastUpdated: The time at which this capability was last updated, in ISO format

.. code-block:: json
	:caption: Response Example

	{ "response": {
		"id": "6",
		"httpMethod": "POST",
		"httpRoute": "/api/*/cdns",
		"capability": "cdn-write",
		"lastUpdated": "2017-04-02 08:22:43"
	},
	"alerts":[
		{
			"level": "success",
			"text": "API-capability mapping was created."
		}
	]}
