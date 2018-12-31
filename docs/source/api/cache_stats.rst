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


.. _to-api-cache_stats:

***************
``cache_stats``
***************
.. caution:: This page is a stub! Much of it may be missing or just downright wrong - it needs a lot of love from people with the domain knowledge required to update it.

Retrieves detailed, aggregated statistics for caches configured in Traffic Ops.

.. versionadded:: 1.2

.. seealso:: This gives an aggregate of statistics for *all caches* within a particular CDN and time range. For statistics basic statistics from all caches regardless of CDN and at the current time, use :ref:`to-api-caches-stats`.

``GET``
-------
Retrieves statistics about the caches within the CDN

:Auth. Required: Yes
:Roles Required: None
:Response Type: Object

Request Structure
-----------------
.. table:: Request Query Parameters

	+------------+----------+-------------------------------------------------------------------------------------------------------------------+
	|    Name    | Required |              Description                                                                                          |
	+============+==========+===================================================================================================================+
	| cdnName    | yes      | The name of a CDN. Results will represent caches within this CDN                                                  |
	+------------+----------+-------------------------------------------------------------------------------------------------------------------+
	| metricType | yes      | The metric type (valid metric types: 'ats.proxy.process.http.current_client_connections', 'bandwidth', 'maxKbps') |
	+------------+----------+-------------------------------------------------------------------------------------------------------------------+
	| startDate  | yes      | The begin date for data aggregation in ISO format, e.g. '2015-08-11T12:30:00-06:00'                               |
	+------------+----------+-------------------------------------------------------------------------------------------------------------------+
	| endDate    | yes      | The end date for data aggregation in ISO format, e.g. '2015-08-12T12:30:00-06:00'                                 |
	+------------+----------+-------------------------------------------------------------------------------------------------------------------+

Response Structure
------------------
:series:  A collection of tabular data and its descriptors

	:columns: An array of names, in order, of the columns of the table. The first is a label for the first entry in each "value", and so on.
	:count:   The total number of data points in the "values" array
	:name:    The name of the metric which was aggregated
	:values:  An array of the actual data points. Each of which is itself an array of properties, which are labeled by the "columns" array. This can be thought of as the data's rows

		:time:  The time in ISO format at which this datum was collected
		:value: The value of the datum. Its meaning is dependent upon "name" - and by extension the ``metricType`` request query parameter

:summary: A summary of the data contained in the "series" object

	:average:                The arithmetic mean across all data point values
	:count:                  The total number of data points in the "series.values" array
	:fifthPercentile:        The right-hand threshold value for the 5\ :sup:`th` percentile
	:max:                    The maximum of the requested metric values
	:min:                    The minimum of the requested metric values
	:ninetyEighthPercentile: The right-hand threshold value for the 98\ :sup:`th` percentile
	:ninetyFifthPercentile:  The right-hand threshold value for the 95\ :sup:`th` percentile


.. code-block:: json
	:caption: Response Example

	{ "response": {
		"series": {
			"columns": [
				"time",
				""
			],
			"count": 29,
			"name": "bandwidth",
			"tags": {
				"cdn": "over-the-top"
			},
			"values": [
				[
					"2015-08-10T22:40:00Z",
					229340299720
				],
				[
					"2015-08-10T22:41:00Z",
					224309221713.334
				],
				[
					"2015-08-10T22:42:00Z",
					229551834168.334
				],
				[
					"2015-08-10T22:43:00Z",
					225179658876.667
				],
				[
					"2015-08-10T22:44:00Z",
					230443968275
				]
			]
		},
		"summary": {
			"average": 970410.295,
			"count": 1376041798,
			"fifthPercentile": 202.03,
			"max": 3875441.02,
			"min": 0,
			"ninetyEighthPercentile": 2957940.93,
			"ninetyFifthPercentile": 2366728.63
		}
	}}
