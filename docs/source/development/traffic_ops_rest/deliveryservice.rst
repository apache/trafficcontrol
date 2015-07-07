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


.. _to-api-ds:

Delivery Service
================

**GET /api/1.1/deliveryservices.json**

  Retrieves all delivery services. See also `Using Traffic Ops - Delivery Service <http://traffic-control-cdn.net/docs/latest/admin/traffic_ops_using.html#delivery-service>`_.

  Authentication Required: Yes

  **Response Properties**

  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  |        Parameter         |  Type  |                                                             Description                                                              |
  +==========================+========+======================================================================================================================================+
  | ``active``               |  bool  | true if active, false if inactive (inact).                                                                                           |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``cacheurl``             | string | Cache URL rule to apply to this delivery service.                                                                                    |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``protocol``             | string | - 0: serve with http:// at EDGE                                                                                                      |
  |                          |        | - 1: serve with https:// at EDGE                                                                                                     |
  |                          |        | - 2: serve with both http:// and https:// at EDGE                                                                                    |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``ccrDnsTtl``            | string | The TTL of the DNS response for A or AAAA queries requesting the IP address of the tr. host.                                         |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``checkPath``            | string | The path portion of the URL to check this deliveryservice for health.                                                                |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``dnsBypassIp``          | string | The IPv4 IP to use for bypass on a DNS deliveryservice  - bypass starts when serving more than the                                   |
  |                          |        | globalMaxMbps traffic on this deliveryservice.                                                                                       |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``dnsBypassIp6``         | string | The IPv6 IP to use for bypass on a DNS deliveryservice - bypass starts when serving more than the                                    |
  |                          |        | globalMaxMbps traffic on this deliveryservice.                                                                                       |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``dnsBypassTtl``         | string | The TTL of the DNS bypass response.                                                                                                  |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``dscp``                 | string | The Differentiated Services Code Point (DSCP) with which to mark downstream (EDGE ->  customer) traffic.                             |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``edgeHeaderRewrite``    | string | The EDGE header rewrite actions to perform.                                                                                          |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``geoLimit``             | string | - 0: None - no limitations                                                                                                           |
  |                          |        | - 1: Only route on CZF file hit                                                                                                      |
  |                          |        | - 2: Only route on CZF hit or when from USA                                                                                          |
  |                          |        |                                                                                                                                      |
  |                          |        | Note that this does not prevent access to content or makes content secure; it just prevents                                          |
  |                          |        | routing to the content by Traffic Router.                                                                                            |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``globalMaxMbps``        | string | The maximum global bandwidth allowed on this deliveryservice. If exceeded, the traffic routes to the                                 |
  |                          |        | dnsByPassIp* for DNS deliveryservices and to the httpBypassFqdn for HTTP deliveryservices.                                           |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``globalMaxTps``         | string | The maximum global transactions per second allowed on this deliveryservice. When this is exceeded                                    |
  |                          |        | traffic will be sent to the dnsByPassIp* for DNS deliveryservices and to the httpBypassFqdn for                                      |
  |                          |        | HTTP deliveryservices                                                                                                                |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``headerRewrite``        | string | The EDGE header rewrite actions to perform.                                                                                          |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``httpBypassFqdn``       | string | The HTTP destination to use for bypass on an HTTP deliveryservice - bypass starts when serving more than the                         |
  |                          |        | globalMaxMbps traffic on this deliveryservice.                                                                                       |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``id``                   | string | The deliveryservice id (database row number).                                                                                        |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``infoUrl``              | string | Use this to add a URL that points to more information about that deliveryservice.                                                    |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``ipv6RoutingEnabled``   |  bool  | false: send IPv4 address of Traffic Router to client on HTTP type del.                                                               |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``longDesc``             | string | Description field 1.                                                                                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``longDesc1``            | string | Description field 2.                                                                                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``longDesc2``            | string | Description field 2.                                                                                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``matchList``            | array  | Array of matchList hashes.                                                                                                           |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``>>type``               | string | The type of MatchList (one of :ref:to-api-types use_in_table='regex').                                                               |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``>>setNumber``          | string | The set Number of the matchList.                                                                                                     |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``>>pattern``            | string | The regexp for the matchList.                                                                                                        |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``maxDnsAnswers``        | string | The maximum number of IPs to put in a A/AAAA response for a DNS deliveryservice (0 means all                                         |
  |                          |        | available).                                                                                                                          |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``missLat``              | string | The latitude to use when the client cannot be found in the CZF or the Geo lookup.                                                    |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``missLong``             | string | The longitude to use when the client cannot be found in the CZF or the Geo lookup.                                                   |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``midHeaderRewrite``     | string | The MID header rewrite actions to perform.                                                                                           |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``multiSiteOrigin``      | string | | Is the Multi Site Origin feature enabled for this delivery service. See :ref:`rl-multi-site-origin`                                |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``orgServerFqdn``        | string | The origin server base URL (FQDN when used in this instance, includes the                                                            |
  |                          |        | protocol (http:// or https://) for use in retrieving content from the origin server.                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``profileDescription``   | string | The description of the Traffic Router Profile with which this deliveryservice is associated.                                         |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``profileName``          | string | The name of the Traffic Router Profile with which this deliveryservice is associated.                                                |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``qstringIgnore``        | string | - 0: no special query string handling; it is for use in the cache-key and pass up to origin.                                         |
  |                          |        | - 1: ignore query string in cache-key, but pass it up to parent and or origin.                                                       |
  |                          |        | - 2: drop query string at edge, and do not use it in the cache-key.                                                                  |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``regexRemap``           | string | Regex Remap rule to apply to this delivery service at the Edge tier.                                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``remapText``            | string | Additional raw remap line text.                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``signed``               |  bool  | - false: token based auth (see :ref:token-based-auth) is not enabled for this deliveryservice.                                       |
  |                          |        | - true: token based auth is enabled for this deliveryservice.                                                                        |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``rangeRequestHandling`` | string | How to treat range requests:                                                                                                         |
  |                          |        |                                                                                                                                      |
  |                          |        | - 0 Do not cache (ranges requested from files taht are already cached due to a non range request will be a HIT)                      |
  |                          |        | - 1 Use the `background_fetch <https://docs.trafficserver.apache.org/en/latest/reference/plugins/background_fetch.en.html>`_ plugin. |
  |                          |        | - 2 Use the cache_range_requests plugin.                                                                                             |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``type``                 | string | The type of this deliveryservice (one of :ref:to-api-types use_in_table='deliveryservice').                                          |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``xmlId``                | string | Unique string that describes this deliveryservice.                                                                                   |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+

  **Response Example** ::

    {
      "response": [
        {
          "active": true,
          "cacheurl": null,
          "protocol": "0",
          "ccrDnsTtl": "3600",
          "checkPath": "/crossdomain.xml",
          "dnsBypassIp": "",
          "dnsBypassIp6": null,
          "dnsBypassTtl": null,
          "dscp": "40",
          "geoLimit": "0",
          "globalMaxMbps": "0",
          "globalMaxTps": "0",
          "headerRewrite": "add-header X-Powered-By: KABLETOWN [L]",
          "edgeHeaderRewrite": "add-header X-Powered-By: KABLETOWN [L]",
          "midHeaderRewrite": null,
          "httpBypassFqdn": "",
          "rangeRequestHandling": "0",
          "id": "12",
          "infoUrl": "",
          "ipv6RoutingEnabled": false,
          "longDesc": "long_desc",
          "longDesc1": "long_desc_1",
          "longDesc2": "long_desc_2",
          "matchList": [
            {
              "type": "HOST_REGEXP",
              "setNumber": "0",
              "pattern": ".*\\.images\\..*"
            }
          ],
          "maxDnsAnswers": "0",
          "missLat": "41.881944",
          "missLong": "-87.627778",
          "orgServerFqdn": "http://cdl.origin.kabletown.net",
          "profileDescription": "Comcast Content Router for cdn2.comcast.net",
          "profileName": "EDGE_CDN2",
          "qstringIgnore": "0",
          "remapText": null,
          "regexRemap": null,
          "signed": true,
          "type": "HTTP",
          "xmlId": "cdl-c2"
        },
        { .. },
        { .. }
      ],
      "version": "1.1"
    }


|

**GET /api/1.1/deliveryservices/:id.json**

  Retrieves a specific delivery service. See also `Using Traffic Ops - Delivery Service <http://traffic-control-cdn.net/docs/latest/admin/traffic_ops_using.html#delivery-service>`_.

  Authentication Required: Yes

  **Response Properties**

  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  |        Parameter         |  Type  |                                                             Description                                                              |
  +==========================+========+======================================================================================================================================+
  | ``active``               |  bool  | true if active, false if inactive (inact).                                                                                           |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``cacheurl``             | string | Cache URL rule to apply to this delivery service.                                                                                    |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``protocol``             | string | - 0: serve with http:// at EDGE                                                                                                      |
  |                          |        | - 1: serve with https:// at EDGE                                                                                                     |
  |                          |        | - 2: serve with both http:// and https:// at EDGE                                                                                    |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``ccrDnsTtl``            | string | The TTL of the DNS response for A or AAAA queries requesting the IP address of the tr. host.                                         |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``checkPath``            | string | The path portion of the URL to check this deliveryservice for health.                                                                |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``dnsBypassIp``          | string | The IPv4 IP to use for bypass on a DNS deliveryservice  - bypass starts when serving more than the                                   |
  |                          |        | globalMaxMbps traffic on this deliveryservice.                                                                                       |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``dnsBypassIp6``         | string | The IPv6 IP to use for bypass on a DNS deliveryservice - bypass starts when serving more than the                                    |
  |                          |        | globalMaxMbps traffic on this deliveryservice.                                                                                       |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``dnsBypassTtl``         | string | The TTL of the DNS bypass response.                                                                                                  |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``dscp``                 | string | The Differentiated Services Code Point (DSCP) with which to mark downstream (EDGE ->  customer) traffic.                             |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``edgeHeaderRewrite``    | string | The EDGE header rewrite actions to perform.                                                                                          |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``geoLimit``             | string | - 0: None - no limitations                                                                                                           |
  |                          |        | - 1: Only route on CZF file hit                                                                                                      |
  |                          |        | - 2: Only route on CZF hit or when from USA                                                                                          |
  |                          |        |                                                                                                                                      |
  |                          |        | Note that this does not prevent access to content or makes content secure; it just prevents                                          |
  |                          |        | routing to the content by Traffic Router.                                                                                            |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``globalMaxMbps``        | string | The maximum global bandwidth allowed on this deliveryservice. If exceeded, the traffic routes to the                                 |
  |                          |        | dnsByPassIp* for DNS deliveryservices and to the httpBypassFqdn for HTTP deliveryservices.                                           |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``globalMaxTps``         | string | The maximum global transactions per second allowed on this deliveryservice. When this is exceeded                                    |
  |                          |        | traffic will be sent to the dnsByPassIp* for DNS deliveryservices and to the httpBypassFqdn for                                      |
  |                          |        | HTTP deliveryservices                                                                                                                |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``headerRewrite``        | string | The EDGE header rewrite actions to perform.                                                                                          |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``httpBypassFqdn``       | string | The HTTP destination to use for bypass on an HTTP deliveryservice - bypass starts when serving more than the                         |
  |                          |        | globalMaxMbps traffic on this deliveryservice.                                                                                       |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``id``                   | string | The deliveryservice id (database row number).                                                                                        |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``infoUrl``              | string | Use this to add a URL that points to more information about that deliveryservice.                                                    |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``ipv6RoutingEnabled``   |  bool  | false: send IPv4 address of Traffic Router to client on HTTP type del.                                                               |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``longDesc``             | string | Description field 1.                                                                                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``longDesc1``            | string | Description field 2.                                                                                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``longDesc2``            | string | Description field 2.                                                                                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``matchList``            | array  | Array of matchList hashes.                                                                                                           |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``>>type``               | string | The type of MatchList (one of :ref:to-api-types use_in_table='regex').                                                               |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``>>setNumber``          | string | The set Number of the matchList.                                                                                                     |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``>>pattern``            | string | The regexp for the matchList.                                                                                                        |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``maxDnsAnswers``        | string | The maximum number of IPs to put in a A/AAAA response for a DNS deliveryservice (0 means all                                         |
  |                          |        | available).                                                                                                                          |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``missLat``              | string | The latitude to use when the client cannot be found in the CZF or the Geo lookup.                                                    |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``missLong``             | string | The longitude to use when the client cannot be found in the CZF or the Geo lookup.                                                   |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``midHeaderRewrite``     | string | The MID header rewrite actions to perform.                                                                                           |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``orgServerFqdn``        | string | The origin server base URL (FQDN when used in this instance, includes the                                                            |
  |                          |        | protocol (http:// or https://) for use in retrieving content from the origin server.                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``profileDescription``   | string | The description of the Traffic Router Profile with which this deliveryservice is associated.                                         |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``profileName``          | string | The name of the Traffic Router Profile with which this deliveryservice is associated.                                                |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``qstringIgnore``        | string | - 0: no special query string handling; it is for use in the cache-key and pass up to origin.                                         |
  |                          |        | - 1: ignore query string in cache-key, but pass it up to parent and or origin.                                                       |
  |                          |        | - 2: drop query string at edge, and do not use it in the cache-key.                                                                  |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``regexRemap``           | string | Regex Remap rule to apply to this delivery service at the Edge tier.                                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``remapText``            | string | Additional raw remap line text.                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``signed``               |  bool  | - false: token based auth (see :ref:token-based-auth) is not enabled for this deliveryservice.                                       |
  |                          |        | - true: token based auth is enabled for this deliveryservice.                                                                        |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``rangeRequestHandling`` | string | How to treat range requests:                                                                                                         |
  |                          |        |                                                                                                                                      |
  |                          |        | - 0 Do not cache (ranges requested from files taht are already cached due to a non range request will be a HIT)                      |
  |                          |        | - 1 Use the `background_fetch <https://docs.trafficserver.apache.org/en/latest/reference/plugins/background_fetch.en.html>`_ plugin. |
  |                          |        | - 2 Use the cache_range_requests plugin.                                                                                             |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``type``                 | string | The type of this deliveryservice (one of :ref:to-api-types use_in_table='deliveryservice').                                          |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``xmlId``                | string | Unique string that describes this deliveryservice.                                                                                   |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+

  **Response Example** ::


    {
      "response": [
        {
          "active": true,
          "cacheurl": null,
          "protocol": "0",
          "ccrDnsTtl": "3600",
          "checkPath": "/crossdomain.xml",
          "dnsBypassIp": "",
          "dnsBypassIp6": null,
          "dnsBypassTtl": null,
          "dscp": "40",
          "geoLimit": "0",
          "globalMaxMbps": "0",
          "globalMaxTps": "0",
          "headerRewrite": "add-header X-Powered-By: KABLETOWN [L]",
          "edgeHeaderRewrite": "add-header X-Powered-By: KABLETOWN [L]",
          "midHeaderRewrite": null,
          "httpBypassFqdn": "",
          "rangeRequestHandling": "0",
          "id": "12",
          "infoUrl": "",
          "ipv6RoutingEnabled": false,
          "longDesc": "long_desc",
          "longDesc1": "long_desc_1",
          "longDesc2": "long_desc_2",
          "matchList": [
            {
              "type": "HOST_REGEXP",
              "setNumber": "0",
              "pattern": ".*\\.images\\..*"
            }
          ],
          "maxDnsAnswers": "0",
          "missLat": "41.881944",
          "missLong": "-87.627778",
          "orgServerFqdn": "http://cdl.origin.kabletown.net",
          "profileDescription": "Comcast Content Router for cdn2.comcast.net",
          "profileName": "EDGE_CDN2",
          "qstringIgnore": "0",
          "remapText": null,
          "regexRemap": null,
          "signed": true,
          "type": "HTTP",
          "xmlId": "cdl-c2"
        }
      ],
      "version": "1.1"
    }

.. _to-api-ds-health:


Health
++++++
.. **GET /api/1.1/deliveryservices/:id/state.json**
.. **GET /api/1.1/deliveryservices/:id/health.json**

**GET /api/1.1/deliveryservices/:id/capacity.json**

  Retrieves the capacity percentages of a delivery service.

  Authentication Required: Yes

  **Request Route Parameters**

  +-----------------+----------+---------------------------------------------------+
  | Name            | Required | Description                                       |
  +=================+==========+===================================================+
  |id               | yes      | delivery service id.                              |
  +-----------------+----------+---------------------------------------------------+

  **Response Properties**

  +------------------------+--------+---------------------------------------------------+
  |       Parameter        |  Type  |                    Description                    |
  +========================+========+===================================================+
  | ``availablePercent``   | number | The percentage of server capacity assigned to     |
  |                        |        | the delivery service that is available.           |
  +------------------------+--------+---------------------------------------------------+
  | ``unavailablePercent`` | number | The percentage of server capacity assigned to the |
  |                        |        | delivery service that is unavailable.             |
  +------------------------+--------+---------------------------------------------------+
  | ``utilizedPercent``    | number | The percentage of server capacity assigned to the |
  |                        |        | delivery service being used.                      |
  +------------------------+--------+---------------------------------------------------+
  | ``maintenancePercent`` | number | The percentage of server capacity assigned to the |
  |                        |        | delivery service that is down for maintenance.    |
  +------------------------+--------+---------------------------------------------------+

  **Response Example** ::

    {
     "response": {
        "availablePercent": 89.0939840205533,
        "unavailablePercent": 0,
        "utilizedPercent": 10.9060020300395,
        "maintenancePercent": 0.0000139494071146245
     },
     "version": "1.1"
    }


|

**GET /api/1.1/deliveryservices/:id/routing.json**

  Retrieves the routing method percentages of a delivery service.

  Authentication Required: Yes

  **Request Route Parameters**

  +-----------------+----------+---------------------------------------------------+
  | Name            | Required | Description                                       |
  +=================+==========+===================================================+
  |id               | yes      | delivery service id.                              |
  +-----------------+----------+---------------------------------------------------+

  **Response Properties**

  +-----------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
  |    Parameter    |  Type  |                                                         Description                                                         |
  +=================+========+=============================================================================================================================+
  | ``staticRoute`` | number | The percentage of Traffic Router responses for this deliveryservice satisfied with pre-configured DNS entries.              |
  +-----------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
  | ``miss``        | number | The percentage of Traffic Router responses for this deliveryservice that were a miss (no location available for client IP). |
  +-----------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
  | ``geo``         | number | The percentage of Traffic Router responses for this deliveryservice satisfied using 3rd party geo-IP mapping.               |
  +-----------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
  | ``err``         | number | The percentage of Traffic Router requests for this deliveryservice resulting in an error.                                   |
  +-----------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
  | ``cz``          | number | The percentage of Traffic Router requests for this deliveryservice satisfied by a CZF hit.                                  |
  +-----------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
  | ``dsr``         | number | The percentage of Traffic Router requests for this deliveryservice satisfied by sending the                                 |
  |                 |        | client to the overflow CDN.                                                                                                 |
  +-----------------+--------+-----------------------------------------------------------------------------------------------------------------------------+

  **Response Example** ::

    {
     "response": {
        "staticRoute": 0,
        "miss": 0,
        "geo": 37.8855391018869,
        "err": 0,
        "cz": 62.1144608981131,
        "dsr": 0
     },
     "version": "1.1"
    }


.. _to-api-ds-metrics:

Metrics
+++++++
**GET /api/1.1/deliveryservices/:id/edge/metric_types/:metric/start_date/:start/end_date/:end/\\
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
     "version": "1.1"
    }


|

**GET /api/1.1/usage/deliveryservices/:ds/cachegroups/:name/metric_types/:metric/start_date/:start_date/\\
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

**GET /api/1.1/cdns/peakusage/:peak_usage_type/deliveryservice/:ds/cachegroup/:name/start_date/:start/\\
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

**GET /api/1.1/deliveryservices/:id/:server_type/metrics/:metric_type/:start/:end.json**

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
     "version": "1.1"
    }


.. _to-api-ds-server:

Server
++++++

**GET /api/1.1/deliveryserviceserver.json**

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
     "version": "1.1",
     "limit": 2
    }



.. _to-api-ds-sslkeys:

SSL Keys
+++++++++

**GET /api/1.1/deliveryservices/xmlId/:xmlid/sslkeys.json**

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
      "version": "1.1",
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

**GET /api/1.1/deliveryservices/hostname/:hostname/sslkeys.json**

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
      "version": "1.1",
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

**GET /api/1.1/deliveryservices/xmlId/:xmlid/sslkeys/delete.json**

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
      "version": "1.1",
      "response": "Successfully deleted ssl keys for <xml_id>"
    }


|
  
**POST /api/1.1/deliveryservices/sslkeys/generate**

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
      "version": "1.1",
      "response": "Successfully created ssl keys for ds-01"
    }

|
  
**POST /api/1.1/deliveryservices/sslkeys/add**

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
      "version": "1.1",
      "response": "Successfully added ssl keys for ds-01"
    }
