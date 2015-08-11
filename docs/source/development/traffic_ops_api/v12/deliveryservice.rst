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


.. _to-api-v12-ds:

Delivery Service
================

**GET /api/1.2/deliveryservice_stats.json**

  Retrieves all delivery services. See also `Using Traffic Ops - Delivery Service <http://traffic-control-cdn.net/docs/latest/admin/traffic_ops_using.html#delivery-service>`_.

  Authentication Required: Yes


  Required Query Parameters: 
                             deliveryServiceName, metricType, startDate, endDate

  deliveryServiceName: 
                       (The delivery service with the desired stats)

  metricType: 
             The metric type (valid metric types: 'kbps', 'ats.proxy.process.http.current_client_connections', 'tps_total', 'tps_2xx','tps_3xx', 'tps_4xx', 'tps_5xx')

  startDate: 
             The begin date 
             (Formatted as ISO8601 (Formatted as ISO8601, for example: '2015-08-11T12:30:00-06:00')  

  endDate: 
           The end date 
           (Formatted as ISO8601, for example: '2015-08-11T13:30:00-06:00')

**GET /api/1.2/deliveryservice_stats.json**

  Example Query: http://localhost:3000/api/1.2/deliveryservice_stats.json?deliveryServiceName=yourdeliveryservice&metricType=kbps&startDate=2015-08-11T12:30:00-06:00&endDate=2015-08-11T13:30-06:00:00&interval=60s

  **Summary Properties**

  +--------------------------------------+-------+-------------------------------------------------------------------------+--+
  |        DeliveryService Stats Summary |  Type |                                                             Description |  |
  +======================================+=======+=========================================================================+==+
  | ``average``                          | float | You complete me!                                                        |  |
  +--------------------------------------+-------+-------------------------------------------------------------------------+--+
  | ``count``                            | float | You complete me!                                                        |  |
  +--------------------------------------+-------+-------------------------------------------------------------------------+--+
  | ``max``                              | float | You complete me!                                                        |  |
  +--------------------------------------+-------+-------------------------------------------------------------------------+--+
  | ``min``                              | float | You complete me!                                                        |  |
  +--------------------------------------+-------+-------------------------------------------------------------------------+--+
  | ``ninetyEighthPercentile``           | float | You complete me!                                                        |  |
  +--------------------------------------+-------+-------------------------------------------------------------------------+--+
  | ``ninetyFifthPercentile``            | float | You complete me!                                                        |  |
  +--------------------------------------+-------+-------------------------------------------------------------------------+--+
  | ``total``                            | float | You complete me!                                                        |  |
  +--------------------------------------+-------+-------------------------------------------------------------------------+--+
  | ``totalBytes``                       | float | You complete me!                                                        |  |
  +--------------------------------------+-------+-------------------------------------------------------------------------+--+
  | ``totalTransactions``                | float | You complete me!                                                        |  |
  +--------------------------------------+-------+-------------------------------------------------------------------------+--+

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

.. _to-api-v12-ds-metrics:

Metrics
+++++++
**GET /api/1.2/deliveryservices/:id/edge/metric_types/:metric/start_date/:start/end_date/:end/\\
interval/:interval/window_start/:window_start/window_end/:window_end.json**

  Retrieves edge summary metrics of all cache groups for a delivery service.

  Authentication Required: Yes

  **Request Route Parameters**

  +------------------+----------+-----------------------------------------------------------------------------+
  |       Name       | Required |                                 Description                                 |
  +==================+==========+=============================================================================+
  | ``id``           | yes      | The delivery service id.                                                    |
  +------------------+----------+-----------------------------------------------------------------------------+
  | ``metric``       | yes      | One of the following: "kbps", "tps_total", "tps_2xx", "tps_3xx", "tps_4xx", |
  |                  |          | "tps_5xx".                                                                  |
  +------------------+----------+-----------------------------------------------------------------------------+
  | ``start``        | yes      | UNIX time, yesterday, now.                                                  |
  +------------------+----------+-----------------------------------------------------------------------------+
  | ``end``          | yes      | UNIX time, yesterday, now.                                                  |
  +------------------+----------+-----------------------------------------------------------------------------+
  | ``interval``     | yes      | > 10                                                                        |
  +------------------+----------+-----------------------------------------------------------------------------+
  | ``window_start`` | yes      | UNIX time, yesterday, now.                                                  |
  +------------------+----------+-----------------------------------------------------------------------------+
  | ``window_end``   | yes      | UNIX time, yesterday, now.                                                  |
  +------------------+----------+-----------------------------------------------------------------------------+

  **Request Query Parameters**

  +-------------+----------+-------------------------------------------+
  |     Name    | Required |                Description                |
  +=============+==========+===========================================+
  | ``summary`` | no       | Flag used to return summary metrics only. |
  +-------------+----------+-------------------------------------------+

  Response Content Type: application/json


  **Response Properties**

  +-----------------+--------+-------------+
  |    Parameter    |  Type  | Description |
  +=================+========+=============+
  | ``ninetyFifth`` | number |             |
  +-----------------+--------+-------------+
  | ``average``     | int    |             |
  +-----------------+--------+-------------+
  | ``min``         | number |             |
  +-----------------+--------+-------------+
  | ``max``         | number |             |
  +-----------------+--------+-------------+
  | ``total``       | number |             |
  +-----------------+--------+-------------+

  **Response Example** ::

    {
     "response": {
        "ninetyFifth": 183982091.479,
        "average": 97444798,
        "min": 31193860.46233,
        "max": 205772883.28367,
        "total": 3643217414091.13
     },
    }


|

**GET /api/1.2/usage/deliveryservices/:ds/cachegroups/:name/metric_types/:metric/start_date/:start_date/\\
end_date/:end_date/interval/:interval.json**

  Retrieves edge metrics of one or all locations (cache groups) for a delivery service.

  Authentication Required: Yes


  **Request Route Parameters**

  +----------------------+----------+-----------------------------------------------------------------------------+
  |         Name         | Required |                                 Description                                 |
  +======================+==========+=============================================================================+
  | ``id``               | yes      | The delivery service id.                                                    |
  +----------------------+----------+-----------------------------------------------------------------------------+
  | ``cache_group_name`` | yes      | name, all.                                                                  |
  +----------------------+----------+-----------------------------------------------------------------------------+
  | ``usage_type``       | yes      | One of the following: "kbps", "tps_total", "tps_2xx", "tps_3xx", "tps_4xx", |
  |                      |          | "tps_5xx".                                                                  |
  +----------------------+----------+-----------------------------------------------------------------------------+
  | ``start``            | yes      | UNIX time, yesterday, now.                                                  |
  +----------------------+----------+-----------------------------------------------------------------------------+
  | ``end``              | yes      | UNIX time, yesterday, now.                                                  |
  +----------------------+----------+-----------------------------------------------------------------------------+
  | ``interval``         | yes      | > 10                                                                        |
  +----------------------+----------+-----------------------------------------------------------------------------+

  **Response Properties**

  +-------------------------+--------+-------------+
  |        Parameter        |  Type  | Description |
  +=========================+========+=============+
  | ``deliveryServiceName`` | string |             |
  +-------------------------+--------+-------------+
  | ``statName``            | string |             |
  +-------------------------+--------+-------------+
  | ``deliveryServiceId``   | string |             |
  +-------------------------+--------+-------------+
  | ``interval``            | int    |             |
  +-------------------------+--------+-------------+
  | ``series``              | array  |             |
  +-------------------------+--------+-------------+
  | ``>>timeBase``          | int    |             |
  +-------------------------+--------+-------------+
  | ``>>samples``           | array  |             |
  +-------------------------+--------+-------------+
  | ``end``                 | string |             |
  +-------------------------+--------+-------------+
  | ``elapsed``             | number |             |
  +-------------------------+--------+-------------+
  | ``cdnName``             | string |             |
  +-------------------------+--------+-------------+
  | ``hostName``            | string |             |
  +-------------------------+--------+-------------+
  | ``summary``             | hash   |             |
  +-------------------------+--------+-------------+
  | >``ninetyFifth``        | number |             |
  +-------------------------+--------+-------------+
  | >``average``            | int    |             |
  +-------------------------+--------+-------------+
  | >``min``                | number |             |
  +-------------------------+--------+-------------+
  | >``max``                | number |             |
  +-------------------------+--------+-------------+
  | >``total``              | number |             |
  +-------------------------+--------+-------------+
  | ``cacheGroupName``      | string |             |
  +-------------------------+--------+-------------+
  | ``start``               | string |             |
  +-------------------------+--------+-------------+

  **Response Example** ::

    TBD
     

|

**GET /api/1.2/cdns/peakusage/:peak_usage_type/deliveryservice/:ds/cachegroup/:name/start_date/:start/\\
end_date/:end/interval/:interval.json**


  Authentication Required: Yes

  **Response Properties**

  +---------------------------------+--------+-------------+
  |            Parameter            |  Type  | Description |
  +=================================+========+=============+
  | ``TotalGBytesServedSinceStart`` | number |             |
  +---------------------------------+--------+-------------+
  +---------------------------------+--------+-------------+
  | ``>>item``                      | number |             |
  +---------------------------------+--------+-------------+
  | ``>>item``                      | number |             |
  +---------------------------------+--------+-------------+
  | ``>>item``                      | number |             |
  +---------------------------------+--------+-------------+
  | ``>>item``                      | number |             |
  +---------------------------------+--------+-------------+
  | ``>>item``                      | number |             |
  +---------------------------------+--------+-------------+
  | ``>>item``                      | number |             |
  +---------------------------------+--------+-------------+

  **Response Example**

  ::
    
    TBD
 

|

**GET /api/1.2/deliveryservices/:id/:server_type/metrics/:metric_type/:start/:end.json**

  Retrieves detailed and summary metrics for MIDs or EDGEs for a delivery service.

  Authentication Required: No

  **Request Route Parameters**

  +-----------------+----------+-----------------------------------------------------------------------------+
  |       Name      | Required |                                 Description                                 |
  +=================+==========+=============================================================================+
  | ``id``          | yes      | The delivery service id.                                                    |
  +-----------------+----------+-----------------------------------------------------------------------------+
  | ``server_type`` | yes      | EDGE or MID.                                                                |
  +-----------------+----------+-----------------------------------------------------------------------------+
  | ``metric_type`` | yes      | One of the following: "kbps", "tps_total", "tps_2xx", "tps_3xx", "tps_4xx", |
  |                 |          | "tps_5xx".                                                                  |
  +-----------------+----------+-----------------------------------------------------------------------------+
  | ``start``       | yes      | UNIX time, yesterday, now.                                                  |
  +-----------------+----------+-----------------------------------------------------------------------------+
  | ``end``         | yes      | UNIX time, yesterday, now.                                                  |
  +-----------------+----------+-----------------------------------------------------------------------------+

  **Response Properties**

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

  **Response Example** ::

    {
     "response": [
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
     ],
    }


.. _to-api-v12-ds-server:

Server
++++++

**GET /api/1.2/deliveryserviceserver.json**

  Authentication Required: Yes

  **Request Query Parameters**

  +-----------+----------+----------------------------------------+
  |    Name   | Required |              Description               |
  +===========+==========+========================================+
  | ``page``  | no       | The page number for use in pagination. |
  +-----------+----------+----------------------------------------+
  | ``limit`` | no       | For use in limiting the result set.    |
  +-----------+----------+----------------------------------------+

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``lastUpdated``       | array  |                                                |
  +----------------------+--------+------------------------------------------------+
  |``server``            | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``deliveryService``   | string |                                                |
  +----------------------+--------+------------------------------------------------+


  **Response Example** ::

    {
     "page": 2,
     "orderby": "deliveryservice",
     "response": [
        {
           "lastUpdated": "2014-09-26 17:53:43",
           "server": "20",
           "deliveryService": "1"
        },
        {
           "lastUpdated": "2014-09-26 17:53:44",
           "server": "21",
           "deliveryService": "1"
        },
     ],
     "limit": 2
    }



.. _to-api-v12-ds-sslkeys:

SSL Keys
+++++++++

**GET /api/1.2/deliveryservices/xmlId/:xmlid/sslkeys.json**

  Authentication Required: Yes

  Role Required: Admin

  **Request Route Parameters**

  +-----------+----------+----------------------------------------+
  |    Name   | Required |              Description               |
  +===========+==========+========================================+
  | ``xmlId`` | yes      | xml_id of the desired delivery service |
  +-----------+----------+----------------------------------------+


  **Request Query Parameters**

  +-------------+----------+--------------------------------+
  |     Name    | Required |          Description           |
  +=============+==========+================================+
  | ``version`` | no       | The version number to retrieve |
  +-------------+----------+--------------------------------+

  **Response Properties**

  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
  |    Parameter     |  Type  |                                                               Description                                                               |
  +==================+========+=========================================================================================================================================+
  | ``crt``          | string | base64 encoded crt file for delivery service                                                                                            |
  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
  | ``csr``          | string | base64 encoded csr file for delivery service                                                                                            |
  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
  | ``key``          | string | base64 encoded private key file for delivery service                                                                                    |
  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
  | ``businessUnit`` | string | The business unit entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response |
  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
  | ``city``         | string | The city entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response          |
  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
  | ``organization`` | string | The organization entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response  |
  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
  | ``hostname``     | string | The hostname entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response      |
  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
  | ``country``      | string | The country entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response       |
  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
  | ``state``        | string | The state entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response         |
  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
  | ``version``      | string | The version of the certificate record in Riak                                                                                           |
  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+


  **Response Example** ::

    {  
      "response": {
        "certificate": {
          "crt": "crt",
          "key": "key",
          "csr": "csr"
        },
        "businessUnit": "CDN_Eng",
        "city": "Denver",
        "organization": "KableTown",
        "hostname": "foober.com",
        "country": "US",
        "state": "Colorado",
        "version": "1"
      }
    }

|

**GET /api/1.2/deliveryservices/hostname/:hostname/sslkeys.json**

  Authentication Required: Yes

  Role Required: Admin

  **Request Route Parameters**

  +--------------+----------+---------------------------------------------------+
  |     Name     | Required |                    Description                    |
  +==============+==========+===================================================+
  | ``hostname`` | yes      | pristine hostname of the desired delivery service |
  +--------------+----------+---------------------------------------------------+


  **Request Query Parameters**

  +-------------+----------+--------------------------------+
  |     Name    | Required |          Description           |
  +=============+==========+================================+
  | ``version`` | no       | The version number to retrieve |
  +-------------+----------+--------------------------------+

  **Response Properties**

  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
  |    Parameter     |  Type  |                                                               Description                                                               |
  +==================+========+=========================================================================================================================================+
  | ``crt``          | string | base64 encoded crt file for delivery service                                                                                            |
  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
  | ``csr``          | string | base64 encoded csr file for delivery service                                                                                            |
  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
  | ``key``          | string | base64 encoded private key file for delivery service                                                                                    |
  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
  | ``businessUnit`` | string | The business unit entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response |
  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
  | ``city``         | string | The city entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response          |
  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
  | ``organization`` | string | The organization entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response  |
  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
  | ``hostname``     | string | The hostname entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response      |
  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
  | ``country``      | string | The country entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response       |
  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
  | ``state``        | string | The state entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response         |
  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
  | ``version``      | string | The version of the certificate record in Riak                                                                                           |
  +------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+


  **Response Example** ::

    {  
      "response": {
        "certificate": {
          "crt": "crt",
          "key": "key",
          "csr": "csr"
        },
        "businessUnit": "CDN_Eng",
        "city": "Denver",
        "organization": "KableTown",
        "hostname": "foober.com",
        "country": "US",
        "state": "Colorado",
        "version": "1"
      }
    }

|

**GET /api/1.2/deliveryservices/xmlId/:xmlid/sslkeys/delete.json**

  Authentication Required: Yes

  Role Required: Admin

  **Request Route Parameters**

  +-----------+----------+----------------------------------------+
  |    Name   | Required |              Description               |
  +===========+==========+========================================+
  | ``xmlId`` | yes      | xml_id of the desired delivery service |
  +-----------+----------+----------------------------------------+

  **Request Query Parameters**

  +-------------+----------+--------------------------------+
  |     Name    | Required |          Description           |
  +=============+==========+================================+
  | ``version`` | no       | The version number to retrieve |
  +-------------+----------+--------------------------------+

  **Response Properties**

  +--------------+--------+------------------+
  |  Parameter   |  Type  |   Description    |
  +==============+========+==================+
  | ``response`` | string | success response |
  +--------------+--------+------------------+

  **Response Example** ::

    {  
      "response": "Successfully deleted ssl keys for <xml_id>"
    }


|
  
**POST /api/1.2/deliveryservices/sslkeys/generate**

  Generates SSL crt, csr, and private key for a delivery service

  Authentication Required: Yes
  Role Required:  Admin

  Response Content Type: application/json

  **Request Properties**


  +--------------+---------+-------------------------------------------------+
  |  Parameter   |   Type  |                   Description                   |
  +==============+=========+=================================================+
  | ``key``      | string  | xml_id of the delivery service                  |
  +--------------+---------+-------------------------------------------------+
  | ``version``  | string  | version of the keys being generated             |
  +--------------+---------+-------------------------------------------------+
  | ``hostname`` | string  | the *pristine hostname* of the delivery service |
  +--------------+---------+-------------------------------------------------+
  | ``country``  | string  |                                                 |
  +--------------+---------+-------------------------------------------------+
  | ``state``    | string  |                                                 |
  +--------------+---------+-------------------------------------------------+
  | ``city``     | string  |                                                 |
  +--------------+---------+-------------------------------------------------+
  | ``org``      | string  |                                                 |
  +--------------+---------+-------------------------------------------------+
  | ``unit``     | boolean |                                                 |
  +--------------+---------+-------------------------------------------------+


  **Request Example** ::


    {
      "key": "ds-01",
      "businessUnit": "CDN Engineering",
      "version": "3",
      "hostname": "tr.ds-01.ott.kabletown.com",
      "certificate": {
        "key": "some_key",
        "csr": "some_csr",
        "crt": "some_crt"
      },
      "country": "US",
      "organization": "Kabletown",
      "city": "Denver",
      "state": "Colorado"
    }

  **Response Properties**

  +--------------+--------+-----------------+
  |  Parameter   |  Type  |   Description   |
  +==============+========+=================+
  | ``response`` | string | response string |
  +--------------+--------+-----------------+
  | ``version``  | string | API version     |
  +--------------+--------+-----------------+


  **Response Example** ::

    {  
      "response": "Successfully created ssl keys for ds-01"
    }

|
  
**POST /api/1.2/deliveryservices/sslkeys/add**

  Allows user to add SSL crt, csr, and private key for a delivery service

  Authentication Required: Yes
  Role Required:  Admin

  **Request Properties**

  +-------------+--------+-------------------------------------+
  |  Parameter  |  Type  |             Description             |
  +=============+========+=====================================+
  | ``key``     | string | xml_id of the delivery service      |
  +-------------+--------+-------------------------------------+
  | ``version`` | string | version of the keys being generated |
  +-------------+--------+-------------------------------------+
  | ``csr``     | string |                                     |
  +-------------+--------+-------------------------------------+
  | ``crt``     | string |                                     |
  +-------------+--------+-------------------------------------+
  | ``key``     | string |                                     |
  +-------------+--------+-------------------------------------+


  **Request Example** ::


    {
      "key": "ds-01",
      "version": "1",
      "certificate": {
        "key": "some_key",
        "csr": "some_csr",
        "crt": "some_crt"
      }
    }

  **Response Properties**

  +--------------+--------+-----------------+
  |  Parameter   |  Type  |   Description   |
  +==============+========+=================+
  | ``response`` | string | response string |
  +--------------+--------+-----------------+
  | ``version``  | string | API version     |
  +--------------+--------+-----------------+


  **Response Example** ::

    {  
      "response": "Successfully added ssl keys for ds-01"
    }
