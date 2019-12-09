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

.. _to-api-stats-summary-create:

*****************
``stats_summary/create``
*****************

``POST``
========
.. deprecated:: 1.1
	Use the ``POST`` method of stats_summary instead.


Post a stats summary for a given stat.

:Auth. Required: Yes
:Roles Required: None
:Response Type: Object

Request Structure
-----------------
:cdnName:             The CDN name for which the summary stat was taken for

	.. note:: If the ``cdn`` is equal to ``all`` it represents summary_stats across all delivery services within the given CDN

:deliveryServiceName: The :term:`Delivery Service` display name for which the summary stat was taken for

	.. note:: If the ``deliveryServiceName`` is equal to ``all`` it represents summary_stats across all delivery services within the given CDN

:statName:            Stat name summary stat represents
:statValue:           Summary stat value
:summaryTime:         Timestamp of summary, in an ISO-like format
:statDate:            Date stat was taken, in :rfc:`3339` format

	.. note:: All fields are required besides ``cdnName`` and ``deliveryServiceName`` which will default to ``all`` if not given

.. code-block:: http
	:caption: Request Example

	POST /api/1.4/stats_summary HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 113
	Content-Type: application/json

	{
        "cdnName": "CDN-in-a-Box",
        "deliveryServiceName": "all",
        "statName": "daily_maxgbps",
        "statValue": 10,
        "summaryTime": "2019-12-05 00:03:57+00",
        "statDate": "2019-12-05"
	}

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: ezxk+iP7o7KE7zpWmGc0j8nz5k+1wAzY0HiNiA2xswTQrt+N+6CgQqUV2r9G1HAsPNr0HF2PhYs/Xr7DrYOw0A==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 06 Dec 2018 02:14:45 GMT
	Content-Length: 97

	{
		"alerts": [
			{
				"text": "Stats Summary was successfully created",
				"level": "success"
			}
			{
				"level": "warning",
				"text": "This endpoint is deprecated, please use 'POST /api/1.4/stats_summary' instead"
			}
		]
	}