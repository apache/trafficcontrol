.. 
.. Copyright 2015 Comcast Cable Communications Management, LLC
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


.. _to-api-v12-cache-stats:

Cache Statistics
===========================

.. _to-api-v12-cache-stats-route:

/api/1.2/cache_stats
++++++++++++++++++++

**GET /api/1.2/cache_stats.json**

  Retrieves statistics about the CDN. 

  Authentication Required: Yes

  
  **Query Parameters**

  +------------------+-------------------------------------------------------------------------------------------------------------------+
  |  Query Parameter | Description                                                                                                       |
  +==================+===================================================================================================================+
  | ``cdnName``      | The name that was configured in the Parameters database table as 'CDN_name'                                       |
  +------------------+-------------------------------------------------------------------------------------------------------------------+
  | ``metricType``   | The metric type (valid metric types: 'ats.proxy.process.http.current_client_connections', 'bandwidth', 'maxKbps') |
  +------------------+-------------------------------------------------------------------------------------------------------------------+
  | ``startDate``    | The begin date (Formatted as ISO8601 (Formatted as ISO8601, for example: '2015-08-11T12:30:00-06:00')             |
  +------------------+-------------------------------------------------------------------------------------------------------------------+
  | ``endDate``      | The end date (Formatted as ISO8601 (Formatted as ISO8601, for example: '2015-08-13T12:30:00-06:00')               |
  +------------------+-------------------------------------------------------------------------------------------------------------------+

  Required Query Parameters: 
                             cdnName, metricType, startDate, endDate


**GET /api/1.2/cache_stats.json**

  Example Query: http://localhost:3000/api/1.2/cache_stats.json?cdnName=yourcdn&metricType=bandwidth&startDate=2015-08-11T12:30:00Z&endDate=2015-08-11T13:30:00Z&interval=60s

  **Summary Properties**

  +-------------------------------+-------+------------------+
  | DeliveryService Stats Summary |  Type | Description      |
  +===============================+=======+==================+
  | ``average``                   | float | You complete me! |
  +-------------------------------+-------+------------------+
  | ``count``                     | float | You complete me! |
  +-------------------------------+-------+------------------+
  | ``max``                       | float | You complete me! |
  +-------------------------------+-------+------------------+
  | ``min``                       | float | You complete me! |
  +-------------------------------+-------+------------------+
  | ``ninetyEighthPercentile``    | float | You complete me! |
  +-------------------------------+-------+------------------+
  | ``ninetyFifthPercentile``     | float | You complete me! |
  +-------------------------------+-------+------------------+
  | ``total``                     | float | You complete me! |
  +-------------------------------+-------+------------------+
  | ``totalBytes``                | float | You complete me! |
  +-------------------------------+-------+------------------+
  | ``totalTransactions``         | float | You complete me! |
  +-------------------------------+-------+------------------+

  **Response Example** ::

                {
                    "response": {
                        "query": {
                            "language": {
                                "influxdbDatabaseName": "cache_stats",
                                "influxdbSeriesQuery": "SELECT sum(value)*1000/6 FROM \"bandwidth\" WHERE 
                                          time > '2015-08-10T16:40:00-06:00' 
                                          AND time < '2015-08-10T17:10:00-06:00' 
                                          AND cdn = 'yourcdn' GROUP BY time(60s), cdn ORDER BY asc",
                                "influxdbSummaryQuery": "SELECT mean(value), 
                                                                percentile(value, 5), 
                                                                percentile(value, 95), 
                                                                percentile(value, 98), 
                                                                min(value), 
                                                                max(value), 
                                                                sum(value), 
                                                                count(value) FROM \"bandwidth\" 
                                                                WHERE cdn = 'over-the-top' 
                                                                AND time > '2015-08-10T16:40:00-06:00' 
                                                                AND time < '2015-08-10T17:10:00-06:00' GROUP BY time(60s), cdn"
                            },
                            "parameters": {
                                "cdnName": "over-the-top",
                                "endDate": "2015-08-10T17:10:00-06:00",
                                "interval": "60s",
                                "limit": null,
                                "metricType": "bandwidth",
                                "offset": null,
                                "orderby": null,
                                "startDate": "2015-08-10T16:40:00-06:00"
                            }
                        },
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
                    }
                }

|
