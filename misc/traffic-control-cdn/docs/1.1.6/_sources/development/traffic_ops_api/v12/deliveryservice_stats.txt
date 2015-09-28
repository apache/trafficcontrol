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


.. _to-api-v12-ds-stats:

Delivery Service Statistics
===========================

.. _to-api-v12-ds-stats-route:

/api/1.2/deliveryservice_stats
++++++++++++++++++++++++++++++

**GET /api/1.2/deliveryservice_stats.json**

  Retrieves statistics on the delivery services. See also `Using Traffic Ops - Delivery Service <http://traffic-control-cdn.net/docs/latest/admin/traffic_ops_using.html#delivery-service>`_.

  Authentication Required: Yes


  Required Query Parameters: 
                             deliveryServiceName, metricType, startDate, endDate

  **Query Parameters**

  +--------------------------------+----------------------------------------------------------------------------------------------------------------------------------------------+
  |  DeliveryService Stats Summary | Description                                                                                                                                  |
  +================================+==============================================================================================================================================+
  | ``deliveryServiceName``        | The delivery service with the desired stats                                                                                                  |
  +--------------------------------+----------------------------------------------------------------------------------------------------------------------------------------------+
  | ``metricType``                 | The metric type (valid metric types: 'kbps', 'out_bytes', 'status_4xx', 'status_5xx', tps_total', 'tps_2xx','tps_3xx', 'tps_4xx', 'tps_5xx') |
  +--------------------------------+----------------------------------------------------------------------------------------------------------------------------------------------+
  | ``startDate``                  | The begin date (Formatted as ISO8601, for example: '2015-08-11T12:30:00-06:00')                                                              |
  +--------------------------------+----------------------------------------------------------------------------------------------------------------------------------------------+
  | ``endDate``                    | The end date (Formatted as ISO8601, for example: '2015-08-12T12:30:00-06:00')                                                                |
  +--------------------------------+----------------------------------------------------------------------------------------------------------------------------------------------+

**GET /api/1.2/deliveryservice_stats.json**

  Example Query: http://localhost:3000/api/1.2/deliveryservice_stats.json?deliveryServiceName=yourdeliveryservice&metricType=kbps&startDate=2015-08-11T12:30:00-06:00&endDate=2015-08-11T13:30-06:00:00&interval=60s

  **Summary Properties**

  +--------------------------------+-------+------------------+
  |  DeliveryService Stats Summary |  Type | Description      |
  +================================+=======+==================+
  | ``average``                    | float | You complete me! |
  +--------------------------------+-------+------------------+
  | ``count``                      | float | You complete me! |
  +--------------------------------+-------+------------------+
  | ``max``                        | float | You complete me! |
  +--------------------------------+-------+------------------+
  | ``min``                        | float | You complete me! |
  +--------------------------------+-------+------------------+
  | ``ninetyEighthPercentile``     | float | You complete me! |
  +--------------------------------+-------+------------------+
  | ``ninetyFifthPercentile``      | float | You complete me! |
  +--------------------------------+-------+------------------+
  | ``total``                      | float | You complete me! |
  +--------------------------------+-------+------------------+
  | ``totalBytes``                 | float | You complete me! |
  +--------------------------------+-------+------------------+
  | ``totalTransactions``          | float | You complete me! |
  +--------------------------------+-------+------------------+

  **Response Example** ::

                {
                    "response": {
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
                        "version": "1.2",
                        "query": {
                            "language": {
                                "influxdbDatabaseName": "deliveryservice_stats",
                                "influxdbSeriesQuery": "SELECT sum(value)/count(value) FROM kbps WHERE cachegroup = 'total' 
                                                        AND deliveryservice = 'cim-jitp' 
                                                        AND time >='2015-08-11T11:30:00Z' 
                                                        AND time <= '2015-08-11T12:30:00Z' GROUP BY time(60s), cachegroup",
                                "influxdbSummaryQuery": "SELECT mean(value), percentile(value, 5), percentile(value, 95), 
                                                                percentile(value, 98), min(value), max(value), 
                                                        count(value) FROM kbps WHERE time >= '2015-08-11T11:30:00Z' 
                                                        AND time <= '2015-08-11T12:30:00Z' 
                                                        AND cachegroup = 'total' 
                                                        AND deliveryservice = 'cim-jitp'"
                            },
                            "parameters": {
                                "deliveryServiceName": "yourdeliveryservicename",
                                "endDate": "2015-08-11T12:30:00Z",
                                "exclude": null,
                                "interval": "60s",
                                "limit": null,
                                "metricType": "kbps",
                                "offset": null,
                                "orderby": null,
                                "startDate": "2015-08-11T11:30:00Z"
                            }
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
                        }
                    }
                }


|
