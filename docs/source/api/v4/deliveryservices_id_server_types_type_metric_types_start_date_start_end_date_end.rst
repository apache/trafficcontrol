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

.. _to-api-deliveryservices-id-server_types-type-metric_types-start_date-start-end_date-end:

****************************************************************************************************
``deliveryservices/{{id}}/server_types/{{type}}/metric_types/start_date/{{start}}/end_date/{{end}}``
****************************************************************************************************

.. danger:: This endpoint doesn't appear to work, and so its use is strongly discouraged! The below documentation cannot be verified, therefore it may be inaccurate and/or incomplete!

``GET``
=======
Retrieves detailed and summary metrics for Mid-tier and Edge-tier caches assigned to a :term:`Delivery Service`.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------------------+----------+-----------------------------------------------------------------------------+
	|       Name       | Required |                                 Description                                 |
	+==================+==========+=============================================================================+
	| ``id``           | yes      | The delivery service id.                                                    |
	+------------------+----------+-----------------------------------------------------------------------------+
	| ``server_types`` | yes      | EDGE or MID.                                                                |
	+------------------+----------+-----------------------------------------------------------------------------+
	| ``metric_types`` | yes      | One of the following: "kbps", "tps", "tps_2xx", "tps_3xx", "tps_4xx",       |
	|                  |          | "tps_5xx".                                                                  |
	+------------------+----------+-----------------------------------------------------------------------------+
	| ``start_date``   | yes      | UNIX time                                                                   |
	+------------------+----------+-----------------------------------------------------------------------------+
	| ``end_date``     | yes      | UNIX time                                                                   |
	+------------------+----------+-----------------------------------------------------------------------------+

.. table:: Request Query Parameters

	+------------------+----------+-----------------------------------------------------------------------------+
	|       Name       | Required |                                 Description                                 |
	+==================+==========+=============================================================================+
	| ``stats``        | no       | Flag used to return only summary metrics                                    |
	+------------------+----------+-----------------------------------------------------------------------------+

Response Structure
------------------
+----------------------+--------+-------------+
|      Parameter       |  Type  | Description |
+======================+========+=============+
| ``stats``            | hash   |             |
+----------------------+--------+-------------+
| ``>>count``          | int    |             |
+----------------------+--------+-------------+
| ``>>98thPercentile`` | number |             |
+----------------------+--------+-------------+
| ``>>min``            | number |             |
+----------------------+--------+-------------+
| ``>>max``            | number |             |
+----------------------+--------+-------------+
| ``>>5thPercentile``  | number |             |
+----------------------+--------+-------------+
| ``>>95thPercentile`` | number |             |
+----------------------+--------+-------------+
| ``>>median``         | number |             |
+----------------------+--------+-------------+
| ``>>mean``           | number |             |
+----------------------+--------+-------------+
| ``>>stddev``         | number |             |
+----------------------+--------+-------------+
| ``>>sum``            | number |             |
+----------------------+--------+-------------+
| ``data``             | array  |             |
+----------------------+--------+-------------+
| ``>>item``           | array  |             |
+----------------------+--------+-------------+
| ``>>time``           | number |             |
+----------------------+--------+-------------+
| ``>>value``          | number |             |
+----------------------+--------+-------------+
| ``label``            | string |             |
+----------------------+--------+-------------+

.. code-block:: json
	:caption: Response Example

	{ "response": [
		{
			"stats": {
				"count": 988,
				"98thPercentile": 16589105.55958,
				"min": 3185442.975,
				"max": 17124754.257,
				"5thPercentile": 3901253.95445,
				"95thPercentile": 16013210.034,
				"median": 8816895.576,
				"mean": 8995846.31741194,
				"stddev": 3941169.83683573,
				"sum": 333296106.060112
			},
			"data": [
				[
					1414303200000,
					12923518.466
				],
				[
					1414303500000,
					12625139.65
				]
			],
			"label": "MID Kbps"
		}
	]}


