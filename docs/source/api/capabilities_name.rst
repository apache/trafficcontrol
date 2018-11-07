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

.. _to-api-capabilities-name:

*************************
``capabilities/{{name}}``
*************************

``GET``
=======
Get a capability by name.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+-------------+----------+--------+------------------------------------+
	|    Name     | Required | Type   |          Description               |
	+=============+==========+========+====================================+
	|   ``name``  |   yes    | string | Capability name.                   |
	+-------------+----------+--------+------------------------------------+

Response Structure
------------------
:description: Describes the APIs covered by the capability.
:lastUpdated: Date and time of the last update made to this capability, in ISO format
:name:        Name of the capability

.. code-block:: json
	:caption: Response Example

	{ "response": [
		{
			"name": "cdn-read",
			"description": "View CDN configuration",
			"lastUpdated": "2017-04-02 08:22:43"
		}
	]}

``PUT``
=======
Edit a capability.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+-------------+----------+--------+------------------------------------+
	|    Name     | Required | Type   |          Description               |
	+=============+==========+========+====================================+
	|   ``name``  |   yes    | string | Capability name.                   |
	+-------------+----------+--------+------------------------------------+

.. table:: Request Data Parameters

	+-------------------+--------+----------+------------------------------------------------+
	|    Parameter      |  Type  | Required |                  Description                   |
	+===================+========+==========+================================================+
	| ``description``   | string | yes      | Describing the APIs covered by the capability  |
	+-------------------+--------+----------+------------------------------------------------+

Response Structure
------------------
:description: Describes the APIs covered by the capability.
:lastUpdated: Date and time of the last update made to this capability, in ISO format

.. code-block:: json
	:caption: Response Example

	{ "response":{
		"name": "cdn-read",
		"description": "View CDN configuration"
	},
	"alerts":[
		{
			"level": "success",
			"text": "Capability was updated."
		}
	]}

``DELETE``
==========
Delete a capability.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters**

	+-----------------+----------+------------------------------------------------+
	| Name            | Required | Description                                    |
	+=================+==========+================================================+
	| ``name``        | yes      | Capability name.                               |
	+-----------------+----------+------------------------------------------------+

Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"level": "success",
			"text": "Capability deleted."
		}
	]}
