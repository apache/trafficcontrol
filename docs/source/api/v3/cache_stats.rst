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


.. _to-api-v3-cache_stats:

***************
``cache_stats``
***************
Retrieves detailed, aggregated statistics for caches in a specific CDN.

.. seealso:: This gives an aggregate of statistics for *all caches* within a particular CDN and time range. For statistics basic statistics from all caches regardless of CDN and at the current time, use :ref:`to-api-v3-caches-stats`.

``GET``
-------
Retrieves statistics about the caches within the CDN

:Auth. Required: Yes
:Roles Required: None
:Response Type: Object

Request Structure
-----------------
.. table:: Request Query Parameters

	+---------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	|    Name             | Required          | Description                                                                                                                                                                               |
	+=====================+===================+===========================================================================================================================================================================================+
	| cdnName             | yes               | The name of a CDN. Results will represent caches within this CDN                                                                                                                          |
	+---------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| endDate             | yes               | The date and time until which statistics shall be aggregated in :rfc:`3339` format (with or without sub-second precision), the number of nanoseconds since the Unix                       |
	|                     |                   | Epoch, or in the same, proprietary format as the ``lastUpdated`` fields prevalent throughout the Traffic Ops API                                                                          |
	+---------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| exclude             | no                | Either "series" to omit the data series from the result, or "summary" to omit the summary data from the result - directly corresponds to fields in the                                    |
	|                     |                   | `Response Structure`_                                                                                                                                                                     |
	+---------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| interval            | no                | Specifies the interval within which data will be "bucketed"; e.g. when requesting data from 2019-07-25T00:00:00Z to 2019-07-25T23:59:59Z with an interval of "1m",                        |
	|                     |                   | the resulting data series (assuming it is not excluded) should contain                                                                                                                    |
	|                     |                   | :math:`24\frac{\mathrm{hours}}{\mathrm{day}}\times60\frac{\mathrm{minutes}}{\mathrm{hour}}\times1\mathrm{day}\times1\frac{\mathrm{minute}}{\mathrm{data point}}=1440\mathrm{data\;points}`|
	|                     |                   | The allowed values for this parameter are valid InfluxQL duration literal strings matching :regexp:`^\d+[mhdw]$`                                                                          |
	|                     |                   |                                                                                                                                                                                           |
	+---------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| limit               | no                | A natural number indicating the maximum amount of data points should be returned in the ``series`` object                                                                                 |
	+---------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| metricType          | yes               | The metric type being reported - one of: 'connections', 'bandwidth', 'maxkbps'                                                                                                            |
	+---------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| offset              | no                | A natural number of data points to drop from the beginning of the returned data set                                                                                                       |
	+---------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| orderby             | no                | Though one struggles to imagine why, this can be used to specify "time" to sort data points by their "time" (which is the default behavior)                                               |
	+---------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| startDate           | yes               | The date and time from which statistics shall be aggregated in :rfc:`3339` format (with or without sub-second precision), the number of nanoseconds since the Unix                        |
	|                     |                   | Epoch, or in the same, proprietary format as the ``lastUpdated`` fields prevalent throughout the Traffic Ops API                                                                          |
	+---------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+

.. _v3-cache_stats-get-request-example:
.. code-block:: http
	:caption: Request Example

	GET /api/3.0/cache_stats?cdnName=CDN&endDate=2019-10-28T20:49:00Z&metricType=bandwidth&startDate=2019-10-28T20:45:00Z HTTP/1.1
	User-Agent: python-requests/2.20.1
	Accept-Encoding: gzip, deflate
	Accept: application/json;timestamp=unix, application/json;timestamp=rfc;q=0.9, application/json;q=0.8, */*;q=0.7
	Connection: keep-alive
	Cookie: mojolicious=...

Content Format
""""""""""""""
It's important to note in :ref:`v3-cache_stats-get-request-example` the use of a complex "Accept" header. This endpoint accepts two special media types in the "Accept" header that instruct it on how to format the timestamps associated with the returned data. Specifically, Traffic Ops will recognize the special, optional, non-standard parameter of :mimetype:`application/json`: ``timestamp``. The values of this parameter are restricted to one of

rfc
	Returned timestamps will be formatted according to :rfc:`3339` (no sub-second precision).
unix
	Returned timestamps will be formatted as the number of nanoseconds since the Unix Epoch (midnight on January 1\ :sup:`st` 1970 UTC).

	.. impl-detail:: The endpoint passes back nanoseconds, specifically, because that is the form used both by InfluxDB, which is used to store the data being served, and Go's standard library. Clients may need to convert the value to match their own standard libraries - e.g. the :js:class:`Date` class in Javascript expects milliseconds.

The default behavior - when only e.g. :mimetype:`application/json` or :mimetype:`*/*` is given - is to use :rfc:`3339` formatting. It will, however, respect quality parameters. It is suggested that clients request timestamps they can handle specifically, rather than relying on this default behavior, as it **is subject to change** and is in fact **expected to invert in the next major release** as string-based time formats become deprecated.

.. seealso:: For more information on the "Accept" HTTP header, consult `its dedicated page on MDN <https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept>`_.

Response Structure
------------------
:series: An object containing the actual data series and information necessary for working with it.

	:columns: This is an array of names of the columns of the data contained in the "values" array - should always be ``["time", "sum_count"]``
	:count:   The number of data points contained in the "values" array
	:name:    The name of the data set. Should always match :samp:`{metric}.ds.1min` where ``metric`` is the requested ``metricType``
	:values:  The actual array of data points. Each represents a length of time specified by the ``interval`` query parameter

		:time:  The time at which the measurement was taken. This corresponds to the *beginning* of the interval. This time comes in the format of either an :rfc:`3339`-formatted string, or a number containing the number of nanoseconds since the Unix Epoch depending on the "Accept" header sent by the client, according to the rules outlined in `Content Format`_.
		:value: The value of the requested ``metricType`` at the time given by ``time``. This will always be a floating point number, unless no data is available for the data interval, in which case it will be ``null``

:summary: A summary of the data contained in the "series" object

	:average:                The arithmetic mean of the data's values
	:count:                  The number of measurements taken within the requested timespan. This is, in general, **not** the same as the ``count`` field of the ``series`` object, as it reflects the number of underlying, un-"bucketed" data points, and is therefore dependent on the implementation of Traffic Stats.
	:fifthPercentile:        Data points with values less than or equal to this number constitute the "bottom" 5% of the data set
	:max:                    The maximum value that can be found in the requested data set
	:min:                    The minimum value that can be found in the requested data set
	:ninetyEighthPercentile: Data points with values greater than or equal to this number constitute the "top" 2% of the data set
	:ninetyFifthPercentile:  Data points with values greater than or equal to this number constitute the "top" 5% of the data set


.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: p4asf1n7fXGtgpW/dWgolJWdXjwDcCjyvjOPFqkckbgoXGUHEj5/wlz7brlQ48t3ZnOWCqOlbsu2eSiBssBtUQ==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 28 Oct 2019 20:49:51 GMT

	{ "response": {
		"series": {
			"columns": [
				"time",
				"sum_count"
			],
			"count": 4,
			"name": "bandwidth.cdn.1min",
			"tags": {
				"cdn": "CDN-in-a-Box"
			},
			"values": [
				[
					1572295500000000000,
					null
				],
				[
					1572295560000000000,
					113.66666666666666
				],
				[
					1572295620000000000,
					108.83333333333334
				],
				[
					1572295680000000000,
					113
				]
			]
		},
		"summary": {
			"average": 111.83333333333333,
			"count": 3,
			"fifthPercentile": 0,
			"max": 113.66666666666666,
			"min": 108.83333333333334,
			"ninetyEighthPercentile": 113.66666666666666,
			"ninetyFifthPercentile": 113.66666666666666
		}
	}}
