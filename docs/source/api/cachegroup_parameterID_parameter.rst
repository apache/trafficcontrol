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

.. _to-api-cachegroup-parameterID-parameter:

*****************************************
``cachegroup/{{parameter ID}}/parameter``
*****************************************
.. deprecated:: 1.1
	Use :ref:`to-api-cachegroupparameters` instead

.. caution:: This endpoint does not appear to work, and thus its use is strongly discouraged!

``GET``
=======
Extract identifying information about all :term:`Cache Groups` with a specific :term:`Parameter`

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------------------+----------+------------------------------------------------+
	|       Name       | Required | Description                                    |
	+==================+==========+================================================+
	| parameter_ID     | yes      | An integral, unique identifier for a parameter |
	+------------------+----------+------------------------------------------------+

Response Structure
------------------
:cachegroups: An array of all :term:`Cache Groups` with an associated :term:`Parameter` identifiable by the ``parameter_id`` request path parameter

	:id:   The integral, unique identifier of the :term:`Cache Group`
	:name: The human-readable name of the :term:`Cache Group`

.. code-block:: json
	:caption: Response Example

	{ "response": {
		"cachegroups": [
			{
				"name": "CDN_in_a_Box_Edge",
				"id": 7
			},
			{
				"name": "CDN_in_a_Box_Mid",
				"id": 6
			},
			{
				"name": "TRAFFIC_ANALYTICS",
				"id": 1
			},
			{
				"name": "TRAFFIC_OPS",
				"id": 2
			},
			{
				"name": "TRAFFIC_OPS_DB",
				"id": 3
			},
			{
				"name": "TRAFFIC_PORTAL",
				"id": 4
			},
			{
				"name": "TRAFFIC_STATS",
				"id": 5
			}
		]
	}}
