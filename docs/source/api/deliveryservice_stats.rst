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


.. _to-api-deliveryservice_stats:

*************************
``deliveryservice_stats``
*************************
.. caution:: This page is a stub! Much of it may be missing or just downright wrong - it needs a lot of love from people with the domain knowledge required to update it.

.. versionadded:: 1.2

.. warning:: This endpoint does **NOT** respect tenancy permissions! The bug is tracked by `GitHub Issue #3187 <https://github.com/apache/trafficcontrol/issues/3187>`_.

``GET``
=======
Retrieves time-aggregated statistics on a specific :term:`Delivery Service`.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Query Parameters

	+----------------------+----------+--------------------------------------------------------------------------------------------------------------------+
	|    Name              | Required |              Description                                                                                           |
	+======================+==========+====================================================================================================================+
	| deliveryServiceName  | yes      | The name of the :term:`Delivery Service` for which statistics will be aggregated                                   |
	+----------------------+----------+--------------------------------------------------------------------------------------------------------------------+
	| metricType           | yes      | The metric type being reported - one of:                                                                           |
	|                      |          |                                                                                                                    |
	|                      |          | kbps                                                                                                               |
	|                      |          |   The total traffic rate in kilobytes per second served by the :term:`Delivery Service`                            |
	|                      |          | out_bytes                                                                                                          |
	|                      |          |   The total number of bytes sent out to clients through the :term:`Delivery Service`                               |
	|                      |          | status_4xx                                                                                                         |
	|                      |          |   The amount of requests that were serviced with 400-499 HTTP status codes                                         |
	|                      |          | status_5xx                                                                                                         |
	|                      |          |   The amount of requests that were serviced with 500-599 HTTP status codes                                         |
	|                      |          | tps_total                                                                                                          |
	|                      |          |   The total traffic rate in transactions per second served by the :term:`Delivery Service`                         |
	|                      |          | tps_2xx                                                                                                            |
	|                      |          |   The total traffic rate in transactions per second serviced with 200-299 HTTP status codes                        |
	|                      |          | tps_3xx                                                                                                            |
	|                      |          |   The total traffic rate in transactions per second serviced with 300-399 HTTP status codes                        |
	|                      |          | tps_4xx                                                                                                            |
	|                      |          |   The total traffic rate in transactions per second serviced with 400-499 HTTP status codes                        |
	|                      |          | tps_5xx                                                                                                            |
	|                      |          |   The total traffic rate in transactions per second serviced with 500-599 HTTP status codes                        |
	|                      |          |                                                                                                                    |
	+----------------------+----------+--------------------------------------------------------------------------------------------------------------------+
	| startDate            | yes      | The date and time from which statistics shall be aggregated in ISO8601 format, e.g. ``2018-08-11T12:30:00-07:00``  |
	+----------------------+----------+--------------------------------------------------------------------------------------------------------------------+
	| endDate              | yes      | The date and time until which statistics shall be aggregated in ISO8601 format, e.g. ``2018-08-12T12:30:00-07:00`` |
	+----------------------+----------+--------------------------------------------------------------------------------------------------------------------+

Response Structure
------------------
.. table:: Response Keys

	+----------------------------+---------------+-----------------------------------------------------------------------------------------+
	| Parameter                  | Type          | Description                                                                             |
	+============================+===============+=========================================================================================+
	|``source``                  | string        | The source of the data                                                                  |
	+----------------------------+---------------+-----------------------------------------------------------------------------------------+
	|``summary``                 | hash          | Summary data                                                                            |
	+----------------------------+---------------+-----------------------------------------------------------------------------------------+
	|``>totalBytes``             | float         |                                                                                         |
	+----------------------------+---------------+-----------------------------------------------------------------------------------------+
	|``>count``                  | int           |                                                                                         |
	+----------------------------+---------------+-----------------------------------------------------------------------------------------+
	|``>min``                    | float         |                                                                                         |
	+----------------------------+---------------+-----------------------------------------------------------------------------------------+
	|``>max``                    | float         |                                                                                         |
	+----------------------------+---------------+-----------------------------------------------------------------------------------------+
	|``>fifthPercentile``        | float         |                                                                                         |
	+----------------------------+---------------+-----------------------------------------------------------------------------------------+
	|``>ninetyEighthPercentile`` | float         |                                                                                         |
	+----------------------------+---------------+-----------------------------------------------------------------------------------------+
	|``>ninetyFifthPercentile``  | float         |                                                                                         |
	+----------------------------+---------------+-----------------------------------------------------------------------------------------+
	|``>average``                | float         |                                                                                         |
	+----------------------------+---------------+-----------------------------------------------------------------------------------------+
	|``>totalTransactions``      | int           |                                                                                         |
	+----------------------------+---------------+-----------------------------------------------------------------------------------------+
	|``series``                  | hash          | Series data                                                                             |
	+----------------------------+---------------+-----------------------------------------------------------------------------------------+
	|``>count``                  | int           |                                                                                         |
	+----------------------------+---------------+-----------------------------------------------------------------------------------------+
	|``>columns``                | array         |                                                                                         |
	+----------------------------+---------------+-----------------------------------------------------------------------------------------+
	|``>name``                   | string        |                                                                                         |
	+----------------------------+---------------+-----------------------------------------------------------------------------------------+
	|``>values``                 | array         |                                                                                         |
	+----------------------------+---------------+-----------------------------------------------------------------------------------------+
	|``>>time``                  | string        |                                                                                         |
	+----------------------------+---------------+-----------------------------------------------------------------------------------------+
	|``>>value``                 | float         |                                                                                         |
	+----------------------------+---------------+-----------------------------------------------------------------------------------------+

.. code-block:: json
	:caption: Response Example

	{ "response": {
		"source": "TrafficStats",
		"summary": {
			"average": 1081172.785,
			"count": 28,
			"fifthPercentile": 888827.26,
			"max": 1326680.31,
			"min": 888827.26,
			"ninetyEighthPercentile": 1324785.47,
			"ninetyFifthPercentile": 1324785.47,
			"totalBytes": 37841047.475,
			"totalTransactions": 1020202030101
		},
		"series": {
			"columns": [
				"time",
				""
			],
		"count": 60,
		"name": "kbps",
		"tags": {
			"cachegroup": "total"
		},
		"values": [
			[
				"2015-08-11T11:36:00Z",
				888827.26
			],
			[
				"2015-08-11T11:37:00Z",
				980336.563333333
			],
			[
				"2015-08-11T11:38:00Z",
				952111.975
			],
			[
				"2015-08-11T11:39:00Z",
				null
			],
			[
				"2015-08-11T11:43:00Z",
				null
			],
			[
				"2015-08-11T11:44:00Z",
				934682.943333333
			],
			[
				"2015-08-11T11:45:00Z",
				1251121.28
			],
			[
				"2015-08-11T11:46:00Z",
				1111012.99
			]
		]
	}}}
