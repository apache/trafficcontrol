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

.. _to-api-v1-cdns-metric_types-metric-start_date-start-end_date-end:

**********************************************************************
``cdns/metric_types/{{metric}}/start_date/{{start}}/end_date/{{end}}``
**********************************************************************

.. danger:: This API endpoint *does* **not** *work*. It isn't implemented in Traffic Ops, and is not expected to be added at any point in the near future. See :issue:`2309` for more information.

.. deprecated:: ATCv4
	This endpoint is deprecated, and will not be added in future API versions.


``GET``
=======
Retrieves :term:`Edge-tier` metrics of one or all :term:`Cache Groups`.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------------+----------+---------------------------+
	|       Name      | Required |        Description        |
	+=================+==========+===========================+
	| ``metric_type`` | yes      | ooff, origin_tps          |
	+-----------------+----------+---------------------------+
	| ``start``       | yes      | UNIX time, yesterday, now |
	+-----------------+----------+---------------------------+
	| ``end``         | yes      | UNIX time, yesterday, now |
	+-----------------+----------+---------------------------+

Response Structure
------------------
:stats: object

	:count:          string
	:98thPercentile: string
	:min:            string
	:max:            string
	:5thPercentile:  string
	:95thPercentile: string
	:mean:           string
	:sum:            string

:data: array

	:time:  int
	:value: number

:label: string

.. code-block:: json
	:caption: Response Example

	{ "response": [ {
		"stats": {
			"count": 1,
			"98thPercentile": 1668.03,
			"min": 1668.03,
			"max": 1668.03,
			"5thPercentile": 1668.03,
			"95thPercentile": 1668.03,
			"mean": 1668.03,
			"sum": 1668.03
		},
		"data": [
			[
				1425135900000,
				1668.03
			],
			[
				1425136200000,
				null
			]
		],
		"label": "Origin TPS"
	}]}
