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

.. _to-api-api_capabilities_id:

************************************
``/api/1.x/api_capabilities/{{id}}``
************************************
Manages a specific API capability.

``GET``
=======
Get an API-capability mapping by id.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+-------------+----------+---------+-----------------------------------------+
	|    Name     | Required |  Type   |         Description                     |
	+=============+==========+=========+=========================================+
	|   ``id``    |   yes    | integer | A unique identifier for this capability |
	+-------------+----------+---------+-----------------------------------------+

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

	{ "response": [
		{
			"capability": "asn-read",
			"httpMethod": "GET",
			"httpRoute": "/api/*/asns",
			"id": "6",
			"lastUpdated": "2017-04-02 08:22:43"
		}
	]}

``PUT``
=======
Edit an API-capability mapping.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+-------------+----------+---------+-----------------------------------------+
	|    Name     | Required |  Type   |         Description                     |
	+=============+==========+=========+=========================================+
	|   ``id``    |   yes    | integer | A unique identifier for this capability |
	+-------------+----------+---------+-----------------------------------------+

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
		"httpMethod": "GET",
		"httpRoute": "/api/*/cdns",
		"capability": "cdn-read",
		"lastUpdated": "2017-04-02 08:22:43"
	},
	"alerts":[
		{
			"level": "success",
			"text": "API-capability mapping was updated."
		}
	]}

DELETE
======
Delete a capability.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+-------------+----------+---------+-----------------------------------------+
	|    Name     | Required |  Type   |         Description                     |
	+=============+==========+=========+=========================================+
	|   ``id``    |   yes    | integer | A unique identifier for this capability |
	+-------------+----------+---------+-----------------------------------------+

Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"level": "success",
			"text": "API-capability mapping deleted."
		}
	]}
