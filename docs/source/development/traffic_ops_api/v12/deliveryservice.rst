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


.. _to-api-v12-ds:

Delivery Service
================

.. _to-api-v12-ds-route:

/api/1.2/deliveryservices
+++++++++++++++++++++++++

**GET /api/1.2/deliveryservices**

  Retrieves all delivery services. See also `Using Traffic Ops - Delivery Service <http://trafficcontrol.apache.org/docs/latest/admin/traffic_ops_using.html#delivery-service>`_.

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  |        Parameter         |  Type  |                                                             Description                                                              |
  +==========================+========+======================================================================================================================================+
  | ``active``               |  bool  | true if active, false if inactive.                                                                                                   |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``cacheurl``             | string | Cache URL rule to apply to this delivery service.                                                                                    |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``ccrDnsTtl``            | string | The TTL of the DNS response for A or AAAA queries requesting the IP address of the tr. host.                                         |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``cdnId``                | string | Id of the CDN to which the delivery service belongs to.                                                                              |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``cdnName``              | string | Name of the CDN to which the delivery service belongs to.                                                                            |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``checkPath``            | string | The path portion of the URL to check this deliveryservice for health.                                                                |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``displayName``          | string | The display name of the delivery service.                                                                                            |
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
  | ``geoLimitRedirectUrl``  | string |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``geoLimit``             | string | - 0: None - no limitations                                                                                                           |
  |                          |        | - 1: Only route on CZF file hit                                                                                                      |
  |                          |        | - 2: Only route on CZF hit or when from USA                                                                                          |
  |                          |        |                                                                                                                                      |
  |                          |        | Note that this does not prevent access to content or makes content secure; it just prevents                                          |
  |                          |        | routing to the content by Traffic Router.                                                                                            |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``geoLimitCountries``    | string |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``geoProvider``          | string |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``globalMaxMbps``        | string | The maximum global bandwidth allowed on this deliveryservice. If exceeded, the traffic routes to the                                 |
  |                          |        | dnsByPassIp* for DNS deliveryservices and to the httpBypassFqdn for HTTP deliveryservices.                                           |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``globalMaxTps``         | string | The maximum global transactions per second allowed on this deliveryservice. When this is exceeded                                    |
  |                          |        | traffic will be sent to the dnsByPassIp* for DNS deliveryservices and to the httpBypassFqdn for                                      |
  |                          |        | HTTP deliveryservices                                                                                                                |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``httpBypassFqdn``       | string | The HTTP destination to use for bypass on an HTTP deliveryservice - bypass starts when serving more than the                         |
  |                          |        | globalMaxMbps traffic on this deliveryservice.                                                                                       |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``id``                   | string | The deliveryservice id (database row number).                                                                                        |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``infoUrl``              | string | Use this to add a URL that points to more information about that deliveryservice.                                                    |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``initialDispersion``    | string |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``ipv6RoutingEnabled``   |  bool  | false: send IPv4 address of Traffic Router to client on HTTP type del.                                                               |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``lastUpdated``          | string |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``logsEnabled``          |  bool  |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``longDesc``             | string | Description field 1.                                                                                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``longDesc1``            | string | Description field 2.                                                                                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``longDesc2``            | string | Description field 2.                                                                                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``>>type``               | string | The type of MatchList (one of :ref:to-api-v11-types use_in_table='regex').                                                           |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``>>setNumber``          | string | The set Number of the matchList.                                                                                                     |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``>>pattern``            | string | The regexp for the matchList.                                                                                                        |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``maxDnsAnswers``        | string | The maximum number of IPs to put in a A/AAAA response for a DNS deliveryservice (0 means all                                         |
  |                          |        | available).                                                                                                                          |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``midHeaderRewrite``     | string | The MID header rewrite actions to perform.                                                                                           |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``missLat``              | string | The latitude to use when the client cannot be found in the CZF or the Geo lookup.                                                    |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``missLong``             | string | The longitude to use when the client cannot be found in the CZF or the Geo lookup.                                                   |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``multiSiteOrigin``      |  bool  | Is the Multi Site Origin feature enabled for this delivery service (0=false, 1=true). See :ref:`rl-multi-site-origin`                |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``multiSiteOriginAlgor`` |  bool  | Is the Multi Site Origin feature enabled for this delivery service (0=false, 1=true). See :ref:`rl-multi-site-origin`                |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``orgServerFqdn``        | string | The origin server base URL (FQDN when used in this instance, includes the                                                            |
  |                          |        | protocol (http:// or https://) for use in retrieving content from the origin server.                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``originShield``         | string |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``profileDescription``   | string | The description of the Traffic Router Profile with which this deliveryservice is associated.                                         |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``profileId``            | string | The id of the Traffic Router Profile with which this deliveryservice is associated.                                                  |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``profileName``          | string | The name of the Traffic Router Profile with which this deliveryservice is associated.                                                |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``protocol``             | string | - 0: serve with http:// at EDGE                                                                                                      |
  |                          |        | - 1: serve with https:// at EDGE                                                                                                     |
  |                          |        | - 2: serve with both http:// and https:// at EDGE                                                                                    |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``qstringIgnore``        | string | - 0: no special query string handling; it is for use in the cache-key and pass up to origin.                                         |
  |                          |        | - 1: ignore query string in cache-key, but pass it up to parent and or origin.                                                       |
  |                          |        | - 2: drop query string at edge, and do not use it in the cache-key.                                                                  |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``rangeRequestHandling`` | string | How to treat range requests:                                                                                                         |
  |                          |        |                                                                                                                                      |
  |                          |        | - 0 Do not cache (ranges requested from files taht are already cached due to a non range request will be a HIT)                      |
  |                          |        | - 1 Use the `background_fetch <https://docs.trafficserver.apache.org/en/latest/reference/plugins/background_fetch.en.html>`_ plugin. |
  |                          |        | - 2 Use the cache_range_requests plugin.                                                                                             |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``regexRemap``           | string | Regex Remap rule to apply to this delivery service at the Edge tier.                                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``regionalGeoBlocking``  |  bool  | Regex Remap rule to apply to this delivery service at the Edge tier.                                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``remapText``            | string | Additional raw remap line text.                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``signed``               |  bool  | - false: token based auth (see :ref:token-based-auth) is not enabled for this deliveryservice.                                       |
  |                          |        | - true: token based auth is enabled for this deliveryservice.                                                                        |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``sslKeyVersion``        | string |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``trRequestHeaders``     | string |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``trResponseHeaders``    | string |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``type``                 | string | The type of this deliveryservice (one of :ref:to-api-v11-types use_in_table='deliveryservice').                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``typeId``               | string | The type of this deliveryservice (one of :ref:to-api-v11-types use_in_table='deliveryservice').                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``xmlId``                | string | Unique string that describes this deliveryservice.                                                                                   |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+

  **Response Example** ::

    {
      "response": [
        {
            "active": true,
            "cacheurl": null,
            "ccrDnsTtl": "3600",
            "cdnId": "2",
            "cdnName": "over-the-top",
            "checkPath": "",
            "displayName": "My Cool Delivery Service",
            "dnsBypassCname": "",
            "dnsBypassIp": "",
            "dnsBypassIp6": "",
            "dnsBypassTtl": "30",
            "dscp": "40",
            "edgeHeaderRewrite": null,
            "exampleURLs": [
                "http://edge.foo-ds.foo.bar.net"
            ],
            "geoLimit": "0",
            "geoLimitCountries": null,
            "geoLimitRedirectURL": null,
            "geoProvider": "0",
            "globalMaxMbps": null,
            "globalMaxTps": "0",
            "httpBypassFqdn": "",
            "id": "442",
            "infoUrl": "",
            "initialDispersion": "1",
            "ipv6RoutingEnabled": true,
            "lastUpdated": "2016-01-26 08:49:35",
            "logsEnabled": false,
            "longDesc": "",
            "longDesc1": "",
            "longDesc2": "",
            "matchList": [
                {
                    "pattern": ".*\\.foo-ds\\..*",
                    "setNumber": "0",
                    "type": "HOST_REGEXP"
                }
            ],
            "maxDnsAnswers": "0",
            "midHeaderRewrite": null,
            "missLat": "41.881944",
            "missLong": "-87.627778",
            "multiSiteOrigin": false,
            "multiSiteOriginAlgorithm": null,
            "orgServerFqdn": "http://baz.boo.net",
            "originShield": null,
            "profileDescription": "Content Router for over-the-top",
            "profileId": "5",
            "profileName": "ROUTER_TOP",
            "protocol": "0",
            "qstringIgnore": "1",
            "rangeRequestHandling": "0",
            "regexRemap": null,
            "regionalGeoBlocking": false,
            "remapText": null,
            "signed": false,
            "sslKeyVersion": "0",
            "trRequestHeaders": null,
            "trResponseHeaders": "Access-Control-Allow-Origin: *",
            "type": "HTTP",
            "typeId": "8",
            "xmlId": "foo-ds"
        }
        { .. },
        { .. }
      ]
    }

|


**GET /api/1.2/deliveryservices/:id**

  Retrieves a specific delivery service. See also `Using Traffic Ops - Delivery Service <http://trafficcontrol.apache.org/docs/latest/admin/traffic_ops_using.html#delivery-service>`_.

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  |        Parameter         |  Type  |                                                             Description                                                              |
  +==========================+========+======================================================================================================================================+
  | ``active``               |  bool  | true if active, false if inactive.                                                                                                   |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``cacheurl``             | string | Cache URL rule to apply to this delivery service.                                                                                    |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``ccrDnsTtl``            | string | The TTL of the DNS response for A or AAAA queries requesting the IP address of the tr. host.                                         |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``cdnId``                | string | Id of the CDN to which the delivery service belongs to.                                                                              |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``cdnName``              | string | Name of the CDN to which the delivery service belongs to.                                                                            |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``checkPath``            | string | The path portion of the URL to check this deliveryservice for health.                                                                |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``displayName``          | string | The display name of the delivery service.                                                                                            |
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
  | ``exampleURLs``          |  array | Entry points into the CDN for this deliveryservice.                                                                                  |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``geoLimitRedirectUrl``  | string |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``geoLimit``             | string | - 0: None - no limitations                                                                                                           |
  |                          |        | - 1: Only route on CZF file hit                                                                                                      |
  |                          |        | - 2: Only route on CZF hit or when from USA                                                                                          |
  |                          |        |                                                                                                                                      |
  |                          |        | Note that this does not prevent access to content or makes content secure; it just prevents                                          |
  |                          |        | routing to the content by Traffic Router.                                                                                            |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``geoLimitCountries``    | string |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``geoProvider``          | string |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``globalMaxMbps``        | string | The maximum global bandwidth allowed on this deliveryservice. If exceeded, the traffic routes to the                                 |
  |                          |        | dnsByPassIp* for DNS deliveryservices and to the httpBypassFqdn for HTTP deliveryservices.                                           |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``globalMaxTps``         | string | The maximum global transactions per second allowed on this deliveryservice. When this is exceeded                                    |
  |                          |        | traffic will be sent to the dnsByPassIp* for DNS deliveryservices and to the httpBypassFqdn for                                      |
  |                          |        | HTTP deliveryservices                                                                                                                |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``httpBypassFqdn``       | string | The HTTP destination to use for bypass on an HTTP deliveryservice - bypass starts when serving more than the                         |
  |                          |        | globalMaxMbps traffic on this deliveryservice.                                                                                       |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``id``                   | string | The deliveryservice id (database row number).                                                                                        |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``infoUrl``              | string | Use this to add a URL that points to more information about that deliveryservice.                                                    |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``initialDispersion``    | string |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``ipv6RoutingEnabled``   |  bool  | false: send IPv4 address of Traffic Router to client on HTTP type del.                                                               |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``lastUpdated``          | string |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``logsEnabled``          |  bool  |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``longDesc``             | string | Description field 1.                                                                                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``longDesc1``            | string | Description field 2.                                                                                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``longDesc2``            | string | Description field 2.                                                                                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``matchList``            | array  | Array of matchList hashes.                                                                                                           |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``>>type``               | string | The type of MatchList (one of :ref:to-api-v11-types use_in_table='regex').                                                           |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``>>setNumber``          | string | The set Number of the matchList.                                                                                                     |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``>>pattern``            | string | The regexp for the matchList.                                                                                                        |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``maxDnsAnswers``        | string | The maximum number of IPs to put in a A/AAAA response for a DNS deliveryservice (0 means all                                         |
  |                          |        | available).                                                                                                                          |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``midHeaderRewrite``     | string | The MID header rewrite actions to perform.                                                                                           |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``missLat``              | string | The latitude to use when the client cannot be found in the CZF or the Geo lookup.                                                    |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``missLong``             | string | The longitude to use when the client cannot be found in the CZF or the Geo lookup.                                                   |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``multiSiteOrigin``      |  bool  | Is the Multi Site Origin feature enabled for this delivery service (0=false, 1=true). See :ref:`rl-multi-site-origin`                |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``multiSiteOriginAlgor`` |  bool  | Is the Multi Site Origin feature enabled for this delivery service (0=false, 1=true). See :ref:`rl-multi-site-origin`                |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``orgServerFqdn``        | string | The origin server base URL (FQDN when used in this instance, includes the                                                            |
  |                          |        | protocol (http:// or https://) for use in retrieving content from the origin server.                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``originShield``         | string |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``profileDescription``   | string | The description of the Traffic Router Profile with which this deliveryservice is associated.                                         |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``profileId``            | string | The id of the Traffic Router Profile with which this deliveryservice is associated.                                                  |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``profileName``          | string | The name of the Traffic Router Profile with which this deliveryservice is associated.                                                |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``protocol``             | string | - 0: serve with http:// at EDGE                                                                                                      |
  |                          |        | - 1: serve with https:// at EDGE                                                                                                     |
  |                          |        | - 2: serve with both http:// and https:// at EDGE                                                                                    |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``qstringIgnore``        | string | - 0: no special query string handling; it is for use in the cache-key and pass up to origin.                                         |
  |                          |        | - 1: ignore query string in cache-key, but pass it up to parent and or origin.                                                       |
  |                          |        | - 2: drop query string at edge, and do not use it in the cache-key.                                                                  |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``rangeRequestHandling`` | string | How to treat range requests:                                                                                                         |
  |                          |        |                                                                                                                                      |
  |                          |        | - 0 Do not cache (ranges requested from files taht are already cached due to a non range request will be a HIT)                      |
  |                          |        | - 1 Use the `background_fetch <https://docs.trafficserver.apache.org/en/latest/reference/plugins/background_fetch.en.html>`_ plugin. |
  |                          |        | - 2 Use the cache_range_requests plugin.                                                                                             |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``regexRemap``           | string | Regex Remap rule to apply to this delivery service at the Edge tier.                                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``regionalGeoBlocking``  |  bool  | Regex Remap rule to apply to this delivery service at the Edge tier.                                                                 |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``remapText``            | string | Additional raw remap line text.                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``signed``               |  bool  | - false: token based auth (see :ref:token-based-auth) is not enabled for this deliveryservice.                                       |
  |                          |        | - true: token based auth is enabled for this deliveryservice.                                                                        |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``sslKeyVersion``        | string |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``trRequestHeaders``     | string |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``trResponseHeaders``    | string |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``type``                 | string | The type of this deliveryservice (one of :ref:to-api-v11-types use_in_table='deliveryservice').                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``typeId``               | string | The type of this deliveryservice (one of :ref:to-api-v11-types use_in_table='deliveryservice').                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``xmlId``                | string | Unique string that describes this deliveryservice.                                                                                   |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+

  **Response Example** ::

    {
      "response": [
        {
            "active": true,
            "cacheurl": null,
            "ccrDnsTtl": "3600",
            "cdnId": "2",
            "cdnName": "over-the-top",
            "checkPath": "",
            "displayName": "My Cool Delivery Service",
            "dnsBypassCname": "",
            "dnsBypassIp": "",
            "dnsBypassIp6": "",
            "dnsBypassTtl": "30",
            "dscp": "40",
            "edgeHeaderRewrite": null,
            "exampleURLs": [
                "http://edge.foo-ds.foo.bar.net"
            ],
            "geoLimit": "0",
            "geoLimitCountries": null,
            "geoLimitRedirectURL": null,
            "geoProvider": "0",
            "globalMaxMbps": null,
            "globalMaxTps": "0",
            "httpBypassFqdn": "",
            "id": "442",
            "infoUrl": "",
            "initialDispersion": "1",
            "ipv6RoutingEnabled": true,
            "lastUpdated": "2016-01-26 08:49:35",
            "logsEnabled": false,
            "longDesc": "",
            "longDesc1": "",
            "longDesc2": "",
            "matchList": [
                {
                    "pattern": ".*\\.foo-ds\\..*",
                    "setNumber": "0",
                    "type": "HOST_REGEXP"
                }
            ],
            "maxDnsAnswers": "0",
            "midHeaderRewrite": null,
            "missLat": "41.881944",
            "missLong": "-87.627778",
            "multiSiteOrigin": false,
            "multiSiteOriginAlgorithm": null,
            "orgServerFqdn": "http://baz.boo.net",
            "originShield": null,
            "profileDescription": "Content Router for over-the-top",
            "profileId": "5",
            "profileName": "ROUTER_TOP",
            "protocol": "0",
            "qstringIgnore": "1",
            "rangeRequestHandling": "0",
            "regexRemap": null,
            "regionalGeoBlocking": false,
            "remapText": null,
            "signed": false,
            "sslKeyVersion": "0",
            "trRequestHeaders": null,
            "trResponseHeaders": "Access-Control-Allow-Origin: *",
            "type": "HTTP",
            "typeId": "8",
            "xmlId": "foo-ds"
        }
      ]
    }

|


.. _to-api-v12-ds-health:

Health
++++++

**GET /api/1.2/deliveryservices/:id/state**

  Retrieves the failover state for a delivery service.

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +------------------+---------+-------------------------------------------------+
  |    Parameter     |  Type   |                   Description                   |
  +==================+=========+=================================================+
  | ``failover``     |  hash   |                                                 |
  +------------------+---------+-------------------------------------------------+
  | ``>locations``   |  array  |                                                 |
  +------------------+---------+-------------------------------------------------+
  | ``>destination`` |  hash   |                                                 |
  +------------------+---------+-------------------------------------------------+
  | ``>>location``   |  string |                                                 |
  +------------------+---------+-------------------------------------------------+
  | ``>>type``       |  string |                                                 |
  +------------------+---------+-------------------------------------------------+
  | ``>configured``  | boolean |                                                 |
  +------------------+---------+-------------------------------------------------+
  | ``>enabled``     | boolean |                                                 |
  +------------------+---------+-------------------------------------------------+
  | ``enabled``      | boolean |                                                 |
  +------------------+---------+-------------------------------------------------+

  **Response Example** ::

    {
        "response": {
            "failover": {
                "locations": [ ],
                "destination": {
                    "location": null,
                    "type": "DNS",
                },
                "configured": false,
                "enabled": false
            },
            "enabled": true
        }
    }

|

**GET /api/1.2/deliveryservices/:id/health**

  Retrieves the health of all locations (cache groups) for a delivery service.

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +------------------+--------+-------------------------------------------------+
  |    Parameter     |  Type  |                   Description                   |
  +==================+========+=================================================+
  | ``totalOnline``  | int    | Total number of online caches across all CDNs.  |
  +------------------+--------+-------------------------------------------------+
  | ``totalOffline`` | int    | Total number of offline caches across all CDNs. |
  +------------------+--------+-------------------------------------------------+
  | ``cachegroups``  | array  | A collection of cache groups.                   |
  +------------------+--------+-------------------------------------------------+
  | ``>online``      | int    | The number of online caches for the cache group |
  +------------------+--------+-------------------------------------------------+
  | ``>offline``     | int    | The number of offline caches for the cache      |
  |                  |        | group.                                          |
  +------------------+--------+-------------------------------------------------+
  | ``>name``        | string | Cache group name.                               |
  +------------------+--------+-------------------------------------------------+

  **Response Example** ::

    {
     "response": {
        "totalOnline": 148,
        "totalOffline": 0,
        "cachegroups": [
           {
              "online": 8,
              "offline": 0,
              "name": "us-co-denver"
           },
           {
              "online": 7,
              "offline": 0,
              "name": "us-de-newcastle"
           }
        ]
     }
    }

|

**GET /api/1.2/deliveryservices/:id/capacity**

  Retrieves the capacity percentages of a delivery service.

  Authentication Required: Yes

  Role(s) Required: None

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
    }


|

**GET /api/1.2/deliveryservices/:id/routing**

  Retrieves the routing method percentages of a delivery service.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +-----------------+----------+---------------------------------------------------+
  | Name            | Required | Description                                       |
  +=================+==========+===================================================+
  |id               | yes      | delivery service id.                              |
  +-----------------+----------+---------------------------------------------------+

  **Response Properties**

  +--------------------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
  |    Parameter             |  Type  |                                                         Description                                                         |
  +==========================+========+=============================================================================================================================+
  | ``staticRoute``          | number | The percentage of Traffic Router responses for this deliveryservice satisfied with pre-configured DNS entries.              |
  +--------------------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
  | ``miss``                 | number | The percentage of Traffic Router responses for this deliveryservice that were a miss (no location available for client IP). |
  +--------------------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
  | ``geo``                  | number | The percentage of Traffic Router responses for this deliveryservice satisfied using 3rd party geo-IP mapping.               |
  +--------------------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
  | ``err``                  | number | The percentage of Traffic Router requests for this deliveryservice resulting in an error.                                   |
  +--------------------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
  | ``cz``                   | number | The percentage of Traffic Router requests for this deliveryservice satisfied by a CZF (coverage zone file) hit.             |
  +--------------------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
  | ``dsr``                  | number | The percentage of Traffic Router requests for this deliveryservice satisfied by sending the                                 |
  |                          |        | client to the overflow CDN.                                                                                                 |
  +--------------------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
  | ``fed``                  | number | The percentage of Traffic Router requests for this deliveryservice satisfied by sending the client to a federated CDN.      |
  +--------------------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
  | ``regionalAlternate``    | number | The percentage of Traffic Router requests for this deliveryservice satisfied by sending the client to the alternate         |
  |                          |        | regional geoblocking URL.                                                                                                   |
  +--------------------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
  | ``regionalDenied``       | number | The percent of Traffic Router requests for this deliveryservice denied due to geolocation policy.                           |
  +--------------------------+--------+-----------------------------------------------------------------------------------------------------------------------------+

  **Response Example** ::

    {
     "response": {
        "staticRoute": 0,
        "miss": 0,
        "geo": 37.8855391018869,
        "err": 0,
        "cz": 62.1144608981131,
        "dsr": 0,
        "fed": 0,
        "regionalAlternate": 0,
        "regionalDenied": 0
     },
    }


.. _to-api-v12-ds-server:

Server
++++++

**GET /api/1.2/deliveryserviceserver**

  Authentication Required: Yes

  Role(s) Required: None

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

**GET /api/1.2/deliveryservices/xmlId/:xmlid/sslkeys**

  Authentication Required: Yes

  Role(s) Required: Admin

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

**GET /api/1.2/deliveryservices/hostname/:hostname/sslkeys**

  Authentication Required: Yes

  Role(s) Required: Admin

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

**GET /api/1.2/deliveryservices/xmlId/:xmlid/sslkeys/delete**

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

  Role(s) Required: Admin

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

|

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

  Allows user to add SSL crt, csr, and private key for a delivery service.

  Authentication Required: Yes

  Role(s) Required:  Admin

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

|

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

**POST /api/1.2/deliveryservices/request**

  Allows a user to send delivery service request details to a specified email address.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Properties**

  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  |  Parameter                             |  Type  | Required |           Description                                                                       |
  +========================================+========+==========+=============================================================================================+
  | ``emailTo``                            | string | yes      | The email to which the delivery service request will be sent.                               |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``details``                            | hash   | yes      | Parameters for the delivery service request.                                                |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>customer``                          | string | yes      | Name of the customer to associated with the delivery service.                               |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>deliveryProtocol``                  | string | yes      | Eg. http or http/https                                                                      |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>routingType``                       | string | yes      | Eg. DNS or HTTP Redirect                                                                    |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>serviceDesc``                       | string | yes      | A description of the delivery service.                                                      |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>peakBPSEstimate``                   | string | yes      | Used to manage cache efficiency and plan for capacity.                                      |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>peakTPSEstimate``                   | string | yes      | Used to manage cache efficiency and plan for capacity.                                      |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>maxLibrarySizeEstimate``            | string | yes      | Used to manage cache efficiency and plan for capacity.                                      |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>originURL``                         | string | yes      | The URL path to the origin server.                                                          |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>hasOriginDynamicRemap``             | bool   | yes      | This is a feature which allows services to use multiple origin URLs for the same service.   |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>originTestFile``                    | string | yes      | A URL path to a test file available on the origin server.                                   |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>hasOriginACLWhitelist``             | bool   | yes      | Is access to your origin restricted using an access control list (ACL or whitelist) of Ips? |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>originHeaders``                     | string | no       | Header values that must be passed to requests to your origin.                               |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>otherOriginSecurity``               | string | no       | Other origin security measures that need to be considered for access.                       |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>queryStringHandling``               | string | yes      | How to handle query strings that come with the request.                                     |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>rangeRequestHandling``              | string | yes      | How to handle range requests.                                                               |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>hasSignedURLs``                     | bool   | yes      | Are Urls signed?                                                                            |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>hasNegativeCachingCustomization``   | bool   | yes      | Any customization required for negative caching?                                            |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>negativeCachingCustomizationNote``  | string | yes      | Negative caching customization instructions.                                                |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>serviceAliases``                    | array  | no       | Service aliases which will be used for this service.                                        |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>rateLimitingGBPS``                  | int    | no       | Rate Limiting - Bandwidth (Gbps)                                                            |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>rateLimitingTPS``                   | int    | no       | Rate Limiting - Transactions/Second                                                         |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>overflowService``                   | string | no       | An overflow point (URL or IP address) used if rate limits are met.                          |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>headerRewriteEdge``                 | string | no       | Headers can be added or altered at each layer of the CDN.                                   |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>headerRewriteMid``                  | string | no       | Headers can be added or altered at each layer of the CDN.                                   |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>headerRewriteRedirectRouter``       | string | no       | Headers can be added or altered at each layer of the CDN.                                   |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
  | ``>notes``                             | string | no       | Additional instructions to provide the delivery service provisioning team.                  |
  +----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+

  **Request Example** ::

    {
       "emailTo": "foo@bar.com",
       "details": {
          "customer": "XYZ Corporation",
          "contentType": "video-on-demand",
          "deliveryProtocol": "http",
          "routingType": "dns",
          "serviceDesc": "service description goes here",
          "peakBPSEstimate": "less-than-5-Gbps",
          "peakTPSEstimate": "less-than-1000-TPS",
          "maxLibrarySizeEstimate": "less-than-200-GB",
          "originURL": "http://myorigin.com",
          "hasOriginDynamicRemap": false,
          "originTestFile": "http://myorigin.com/crossdomain.xml",
          "hasOriginACLWhitelist": true,
          "originHeaders": "",
          "otherOriginSecurity": "",
          "queryStringHandling": "ignore-in-cache-key-and-pass-up",
          "rangeRequestHandling": "range-requests-not-used",
          "hasSignedURLs": true,
          "hasNegativeCachingCustomization": true,
          "negativeCachingCustomizationNote": "negative caching instructions",
          "serviceAliases": [
             "http://alias1.com",
             "http://alias2.com"
          ],
          "rateLimitingGBPS": 50,
          "rateLimitingTPS": 5000,
          "overflowService": "http://overflowcdn.com",
          "headerRewriteEdge": "",
          "headerRewriteMid": "",
          "headerRewriteRedirectRouter": "",
          "notes": ""
       }
    }

|

  **Response Properties**

  +-------------+--------+----------------------------------+
  |  Parameter  |  Type  |           Description            |
  +=============+========+==================================+
  | ``alerts``  | array  | A collection of alert messages.  |
  +-------------+--------+----------------------------------+
  | ``>level``  | string | Success, info, warning or error. |
  +-------------+--------+----------------------------------+
  | ``>text``   | string | Alert message.                   |
  +-------------+--------+----------------------------------+
  | ``version`` | string |                                  |
  +-------------+--------+----------------------------------+

  **Response Example** ::

    {
      "alerts": [
            {
                "level": "success",
                "text": "Delivery Service request sent to foo@bar.com."
            }
        ]
    }

|

**POST /api/1.2/deliveryservices**

  Allows user to create a delivery service.

  Authentication Required: Yes

  Role(s) Required:  admin or oper

  **Request Properties**

  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | Parameter              | Required | Description                                                                                             |
  +========================+==========+=========================================================================================================+
  | xmlId                  | yes      | Unique string that describes this deliveryservice.                                                      |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | active                 | yes      | true if active, false if inactive.                                                                      |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | cacheurl               | no       | Cache URL rule to apply to this delivery service.                                                       |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | protocol               | yes      | - 0: serve with http:// at EDGE                                                                         |
  |                        |          | - 1: serve with https:// at EDGE                                                                        |
  |                        |          | - 2: serve with both http:// and https:// at EDGE                                                       |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | ccrDnsTtl              | no       | The TTL of the DNS response for A or AAAA queries requesting the IP address of the tr.host.             |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | checkPath              | no       | The path portion of the URL to check this deliveryservice for health.                                   |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | dnsBypassIp            | no       | The IPv4 IP to use for bypass on a DNS deliveryservice - bypass starts when serving more than the       |
  |                        |          | globalMaxMbps traffic on this deliveryservice.                                                          |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | dnsBypassIp6           | no       | The IPv6 IP to use for bypass on a DNS deliveryservice - bypass starts when serving more than the       |
  |                        |          | globalMaxMbps traffic on this deliveryservice.                                                          |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | dnsBypassTtl           | no       | The TTL of the DNS bypass response.                                                                     |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | dscp                   | no       | The Differentiated Services Code Point (DSCP) with which to mark downstream (EDGE -> customer) traffic. |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | edgeHeaderRewrite      | no       | The EDGE header rewrite actions to perform.                                                             |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | geoLimit               | no       | - 0: None - no limitations                                                                              |
  |                        |          | - 1: Only route on CZF file hit                                                                         |
  |                        |          | - 2: Only route on CZF hit or when from geo limit countries                                             |
  |                        |          |                                                                                                         |
  |                        |          | Note that this does not prevent access to content or makes content secure; it just prevents             |
  |                        |          | routing to the content by Traffic Router.                                                               |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | geoLimitCountries      | no       | The geo limit countries.                                                                                |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | geoLimitRedirectURL    | no       | This is the URL Traffic Router will redirect to when Geo Limit Failure.                                 |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | geoProvider            | no       | - 0: Maxmind(default)                                                                                   |
  |                        |          | - 1: Neustar                                                                                            |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | globalMaxMbps          | no       | The maximum global bandwidth allowed on this deliveryservice. If exceeded, the traffic routes to the    |
  |                        |          | dnsByPassIp* for DNS deliveryservices and to the httpBypassFqdn for HTTP deliveryservices.              |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | globalMaxTps           | no       | The maximum global transactions per second allowed on this deliveryservice. When this is exceeded       |
  |                        |          | traffic will be sent to the dnsByPassIp* for DNS deliveryservices and to the httpBypassFqdn for         |
  |                        |          | HTTP deliveryservices                                                                                   |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | httpBypassFqdn         | no       | The HTTP destination to use for bypass on an HTTP deliveryservice - bypass starts when serving more     |
  |                        |          | than the globalMaxMbps traffic on this deliveryservice.                                                 |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | infoUrl                | no       | Use this to add a URL that points to more information about that deliveryservice.                       |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | ipv6RoutingEnabled     | no       | false: send IPv4 address of Traffic Router to client on HTTP type del.                                  |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | longDesc               | no       | Description field.                                                                                      |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | longDesc1              | no       | Description field 1.                                                                                    |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | longDesc2              | no       | Description field 2.                                                                                    |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | matchList              | yes      | Array of matchList hashes.                                                                              |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | >type                  | yes      | The type of MatchList (one of :ref:to-api-v12-types use_in_table='regex').                              |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | >setNumber             | yes      | The set Number of the matchList.                                                                        |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | >pattern               | yes      | The regexp for the matchList.                                                                           |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | maxDnsAnswers          | no       | The maximum number of IPs to put in a A/AAAA response for a DNS deliveryservice (0 means all            |
  |                        |          | available).                                                                                             |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | missLat                | no       | The latitude to use when the client cannot be found in the CZF or the Geo lookup.                       |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | missLong               | no       | The longitude to use when the client cannot be found in the CZF or the Geo lookup.                      |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | midHeaderRewrite       | no       | The MID header rewrite actions to perform.                                                              |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | multiSiteOrigin        | yes      | 1 if enabled, 0 if disabled.                                                                            |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | orgServerFqdn          | yes      | The origin server base URL (FQDN when used in this instance, includes the                               |
  |                        |          | protocol (http:// or https://) for use in retrieving content from the origin server.                    |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | profileName            | yes      | Traffic router profile name, for example "CCR_CDN"                                                      |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | qstringIgnore          | no       | - 0: no special query string handling; it is for use in the cache-key and pass up to origin.            |
  |                        |          | - 1: ignore query string in cache-key, but pass it up to parent and or origin.                          |
  |                        |          | - 2: drop query string at edge, and do not use it in the cache-key.                                     |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | regexRemap             | no       | Regex Remap rule to apply to this delivery service at the Edge tier.                                    |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | remapText              | no       | Additional raw remap line text.                                                                         |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | signed                 | no       | - false: token based auth (see :ref:token-based-auth) is not enabled for this deliveryservice.          |
  |                        |          | - true: token based auth is enabled for this deliveryservice.                                           |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | rangeRequestHandling   | no       | How to treat range requests:                                                                            |
  |                        |          |                                                                                                         |
  |                        |          | - 0 Do not cache (ranges requested from files taht are already cached due to a non range request will   |
  |                        |          |   be a HIT)                                                                                             |
  |                        |          | - 1 Use the background_fetch plugin.                                                                    |
  |                        |          | - 2 Use the cache_range_requests plugin.                                                                |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | type                   | yes      | The type of this deliveryservice (one of :ref:to-api-v12-types use_in_table='deliveryservice').         |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | displayName            | yes      | Display name                                                                                            |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | cdnName                | yes      | cdn name                                                                                                |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | dnsBypassCname         | no       | Bypass CNAME                                                                                            |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | trResponseHeaders      | no       | Traffic router additional response headers                                                              |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | initialDispersion      | no       | Initial dispersion                                                                                      |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | regionalGeoBlocking    | no       | Is the Regional Geo Blocking feature enabled for this delivery service.                                 |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | sslKeyVersion          | no       | SSL key version                                                                                         |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | originShield           | no       | Origin shield                                                                                           |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | trRequestHeaders       | no       | Traffic router log request headers                                                                      |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | logsEnabled            | no       | - false: No                                                                                             |
  |                        |          | - true: Yes                                                                                             |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+


  **Request Example** ::

    {
        "xmlId": "my_ds_1",
        "displayName": "my_ds_displayname_1",
        "protocol": "1",
        "orgServerFqdn": "http://10.75.168.91",
        "cdnName": "cdn_number_1",
        "profileName": "CCR_CDN1",
        "type": "HTTP",
        "multiSiteOrigin": "0",
        "active": "false",
        "matchList": [
            {
                "type":  "HOST_REGEXP",
                "pattern": ".*\\.ds_1\\..*"
                "setNumber": "0"
            },
            {
                "type":  "HOST_REGEXP",
                "pattern": ".*\\.my_vod1\\..*"
                "setNumber": "1"
            }
        ]
    }


  **Response Example** ::

    {
        "response":{
            "xmlId":"my_ds_1",
            "active":"false",
            "protocol":"0",
            "missLong":null,
            "maxDnsAnswers":"0",
            "profileName": "CCR_CDN1",
            "multiSiteOrigin":"0",
            "dnsBypassIp6":null,
            "globalMaxTps":"0",
            "orgServerFqdn":"http:\/\/10.75.168.91",
            "infoUrl":null,
            "rangeRequestHandling":null,
            "id":"311",
            "trResponseHeaders":null,
            "ipv6RoutingEnabled":null,
            "midHeaderRewrite":null,
            "longDesc":null,
            "httpBypassFqdn":null,
            "cdnName":"cdn_number_1",
            "protocol":"1",
            "missLat":null,
            "globalMaxMbps":"0",
            "initialDispersion":null,
            "type":"HTTP",
            "geoLimit":null,
            "dnsBypassTtl":null,
            "dnsBypassCname":null,
            "ccrDnsTtl":null,
            "longDesc2":null,
            "remapText":null,
            "dnsBypassIp":null,
            "longDesc1":null,
            "checkPath":null,
            "qstringIgnore":null,
            "dscp":"1",
            "regexRemap":null,
            "edgeHeaderRewrite":null,
            "sslKeyVersion":"0",
            "displayName":"my_ds_displayname_1",
            "cacheurl":null,
            "signed":"0",
            "matchList":[
                {
                    "type":"HOST_REGEXP",
                    "setNumber":"0",
                    "pattern":".*\\.ds_1\\..*"
                },
                {
                    "type":"HOST_REGEXP",
                    "setNumber":"1",
                    "pattern":".*\\.my_vod1\\..*"
                }
            ],
            "regionalGeoBlocking":0,
            "originShield":null,
            "trRequestHeaders":null,
            "geoProvider":"0",
            "logsEnabled":"false",
        }
        "alerts":[
            {
                "level": "success",
                "text": "Delivery service was created: 312"
            }
        ]
    }

|

**PUT /api/1.2/deliveryservices/{:id}**

  Allows user to edit a delivery service.

  Authentication Required: Yes

  Role(s) Required:  admin or oper

  **Request Route Parameters**

  +-----------------+----------+---------------------------------------------------+
  | Name            | Required | Description                                       |
  +=================+==========+===================================================+
  |id               | yes      | delivery service id.                              |
  +-----------------+----------+---------------------------------------------------+

  **Request Properties**

  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | Parameter              | Required | Description                                                                                             |
  +========================+==========+=========================================================================================================+
  | xmlId                  | yes      | Unique string that describes this deliveryservice.                                                      |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | active                 | yes      | true if active, false if inactive.                                                                      |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | cacheurl               | no       | Cache URL rule to apply to this delivery service.                                                       |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | protocol               | yes      | - 0: serve with http:// at EDGE                                                                         |
  |                        |          | - 1: serve with https:// at EDGE                                                                        |
  |                        |          | - 2: serve with both http:// and https:// at EDGE                                                       |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | ccrDnsTtl              | no       | The TTL of the DNS response for A or AAAA queries requesting the IP address of the tr.host.             |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | checkPath              | no       | The path portion of the URL to check this deliveryservice for health.                                   |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | dnsBypassIp            | no       | The IPv4 IP to use for bypass on a DNS deliveryservice - bypass starts when serving more than the       |
  |                        |          | globalMaxMbps traffic on this deliveryservice.                                                          |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | dnsBypassIp6           | no       | The IPv6 IP to use for bypass on a DNS deliveryservice - bypass starts when serving more than the       |
  |                        |          | globalMaxMbps traffic on this deliveryservice.                                                          |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | dnsBypassTtl           | no       | The TTL of the DNS bypass response.                                                                     |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | dscp                   | no       | The Differentiated Services Code Point (DSCP) with which to mark downstream (EDGE -> customer) traffic. |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | edgeHeaderRewrite      | no       | The EDGE header rewrite actions to perform.                                                             |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | geoLimit               | no       | - 0: None - no limitations                                                                              |
  |                        |          | - 1: Only route on CZF file hit                                                                         |
  |                        |          | - 2: Only route on CZF hit or when from geo limit countries                                             |
  |                        |          |                                                                                                         |
  |                        |          | Note that this does not prevent access to content or makes content secure; it just prevents             |
  |                        |          | routing to the content by Traffic Router.                                                               |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | geoLimitCountries      | no       | The geo limit countries.                                                                                |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | geoLimitRedirectURL    | no       | This is the URL Traffic Router will redirect to when Geo Limit Failure.                                 |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | geoProvider            | no       | - 0: Maxmind(default)                                                                                   |
  |                        |          | - 1: Neustar                                                                                            |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | globalMaxMbps          | no       | The maximum global bandwidth allowed on this deliveryservice. If exceeded, the traffic routes to the    |
  |                        |          | dnsByPassIp* for DNS deliveryservices and to the httpBypassFqdn for HTTP deliveryservices.              |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | globalMaxTps           | no       | The maximum global transactions per second allowed on this deliveryservice. When this is exceeded       |
  |                        |          | traffic will be sent to the dnsByPassIp* for DNS deliveryservices and to the httpBypassFqdn for         |
  |                        |          | HTTP deliveryservices                                                                                   |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | httpBypassFqdn         | no       | The HTTP destination to use for bypass on an HTTP deliveryservice - bypass starts when serving more     |
  |                        |          | than the globalMaxMbps traffic on this deliveryservice.                                                 |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | infoUrl                | no       | Use this to add a URL that points to more information about that deliveryservice.                       |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | ipv6RoutingEnabled     | no       | false: send IPv4 address of Traffic Router to client on HTTP type del.                                  |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | longDesc               | no       | Description field.                                                                                      |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | longDesc1              | no       | Description field 1.                                                                                    |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | longDesc2              | no       | Description field 2.                                                                                    |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | matchList              | yes      | Array of matchList hashes.                                                                              |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | >type                  | yes      | The type of MatchList (one of :ref:to-api-v12-types use_in_table='regex').                              |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | >setNumber             | yes      | The set Number of the matchList.                                                                        |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | >pattern               | yes      | The regexp for the matchList.                                                                           |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | maxDnsAnswers          | no       | The maximum number of IPs to put in a A/AAAA response for a DNS deliveryservice (0 means all            |
  |                        |          | available).                                                                                             |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | missLat                | no       | The latitude to use when the client cannot be found in the CZF or the Geo lookup.                       |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | missLong               | no       | The longitude to use when the client cannot be found in the CZF or the Geo lookup.                      |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | midHeaderRewrite       | no       | The MID header rewrite actions to perform.                                                              |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | multiSiteOrigin        | yes      | 1 if enabled, 0 if disabled.                                                                            |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | orgServerFqdn          | yes      | The origin server base URL (FQDN when used in this instance, includes the                               |
  |                        |          | protocol (http:// or https://) for use in retrieving content from the origin server.                    |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | profileName            | yes      | Traffic router profile name, for example "CCR_CDN"                                                      |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | qstringIgnore          | no       | - 0: no special query string handling; it is for use in the cache-key and pass up to origin.            |
  |                        |          | - 1: ignore query string in cache-key, but pass it up to parent and or origin.                          |
  |                        |          | - 2: drop query string at edge, and do not use it in the cache-key.                                     |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | regexRemap             | no       | Regex Remap rule to apply to this delivery service at the Edge tier.                                    |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | remapText              | no       | Additional raw remap line text.                                                                         |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | signed                 | no       | - false: token based auth (see :ref:token-based-auth) is not enabled for this deliveryservice.          |
  |                        |          | - true: token based auth is enabled for this deliveryservice.                                           |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | rangeRequestHandling   | no       | How to treat range requests:                                                                            |
  |                        |          |                                                                                                         |
  |                        |          | - 0 Do not cache (ranges requested from files taht are already cached due to a non range request will   |
  |                        |          |   be a HIT)                                                                                             |
  |                        |          | - 1 Use the background_fetch plugin.                                                                    |
  |                        |          | - 2 Use the cache_range_requests plugin.                                                                |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | type                   | yes      | The type of this deliveryservice (one of :ref:to-api-v12-types use_in_table='deliveryservice').         |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | displayName            | yes      | Display name                                                                                            |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | cdnName                | yes      | cdn name                                                                                                |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | dnsBypassCname         | no       | Bypass CNAME                                                                                            |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | trResponseHeaders      | no       | Traffic router additional response headers                                                              |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | initialDispersion      | no       | Initial dispersion                                                                                      |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | regionalGeoBlocking    | no       | Is the Regional Geo Blocking feature enabled for this delivery service.                                 |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | sslKeyVersion          | no       | SSL key version                                                                                         |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | originShield           | no       | Origin shield                                                                                           |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | trRequestHeaders       | no       | Traffic router log request headers                                                                      |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+
  | logsEnabled            | no       | - false: No                                                                                             |
  |                        |          | - true: Yes                                                                                             |
  +------------------------+----------+---------------------------------------------------------------------------------------------------------+


  **Request Example** ::

    {
        "xmlId": "my_ds_2",
        "displayName": "my_ds_displayname_2",
        "protocol": "1",
        "orgServerFqdn": "http://10.75.168.91",
        "cdnName": "cdn_number_1",
        "profileName": "CCR_CDN1",
        "type": "HTTP",
        "multiSiteOrigin": "0",
        "active": "true",
        "matchList": [
            {
                "type":  "HOST_REGEXP",
                "pattern": ".*\\.ds_1\\..*"
                "setNumber": "0"
            },
            {
                "type":  "HOST_REGEXP",
                "pattern": ".*\\.my_vod1\\..*"
                "setNumber": "1"
            }
        ]
    }


  **Response Example** ::

    {
        "response":{
            "xmlId":"my_ds_2",
            "active":"true",
            "protocol":"0",
            "missLong":null,
            "maxDnsAnswers":"0",
            "profileName": "CCR_CDN1",
            "multiSiteOrigin":"0",
            "dnsBypassIp6":null,
            "globalMaxTps":"0",
            "orgServerFqdn":"http:\/\/10.75.168.91",
            "infoUrl":null,
            "rangeRequestHandling":null,
            "id":"311",
            "trResponseHeaders":null,
            "ipv6RoutingEnabled":null,
            "midHeaderRewrite":null,
            "longDesc":null,
            "httpBypassFqdn":null,
            "cdnName":"cdn_number_1",
            "protocol":"1",
            "missLat":null,
            "globalMaxMbps":"0",
            "initialDispersion":null,
            "type":"HTTP",
            "geoLimit":null,
            "dnsBypassTtl":null,
            "dnsBypassCname":null,
            "ccrDnsTtl":null,
            "longDesc2":null,
            "remapText":null,
            "dnsBypassIp":null,
            "longDesc1":null,
            "checkPath":null,
            "qstringIgnore":null,
            "dscp":"1",
            "regexRemap":null,
            "edgeHeaderRewrite":null,
            "sslKeyVersion":"0",
            "displayName":"my_ds_displayname_2",
            "cacheurl":null,
            "signed":"0",
            "matchList":[
                {
                    "type":"HOST_REGEXP",
                    "setNumber":"0",
                    "pattern":".*\\.ds_1\\..*"
                },
                {
                    "type":"HOST_REGEXP",
                    "setNumber":"1",
                    "pattern":".*\\.my_vod1\\..*"
                }
            ],
            "regionalGeoBlocking":0,
            "originShield":null,
            "trRequestHeaders":null,
            "geoProvider":"0",
            "logsEnabled":"false",
        }
        "alerts":[
            {
                "level": "success",
                "text": "Delivery service was updated: 312"
            }
        ]
    }

|

**DELETE /api/1.2/deliveryservices/{:id}**

  Allows user to delete a delivery service.

  Authentication Required: Yes

  Role(s) Required:  admin or oper

  **Request Route Parameters**

  +-----------------+----------+---------------------------------------------------+
  | Name            | Required | Description                                       |
  +=================+==========+===================================================+
  | id              | yes      | delivery service id.                              |
  +-----------------+----------+---------------------------------------------------+

   **Response Example** ::

    {
           "alerts": [
                     {
                             "level": "success",
                             "text": "Delivery service was deleted."
                     }
             ],
    }

|

**POST /api/1.2/deliveryservices/:xml_id/servers**

  Assign caches to a delivery service.

  Authentication Required: Yes

  Role(s) Required:  admin or oper

  **Request Route Parameters**

  +--------+----------+-----------------------------------+
  | Name   | Required | Description                       |
  +========+==========+===================================+
  | xml_id | yes      | the xml_id of the deliveryservice |
  +--------+----------+-----------------------------------+

  **Request Properties**

  +--------------+----------+-------------------------------------------------------------------------------------------------------------+
  | Parameter    | Required | Description                                                                                                 |
  +==============+==========+=============================================================================================================+
  | serverNames  | yes      | array of hostname of cache servers to assign to this deliveryservice, for example: [ "server1", "server2" ] |
  +--------------+----------+-------------------------------------------------------------------------------------------------------------+

  **Request Example** ::

    {
        "serverNames": [
            "tc1_ats1"
        ]
    }

  **Response Properties**

  +--------------+--------+-------------------------------------------------------------------------------------------------------------+
  | Parameter    | Type   | Description                                                                                                 |
  +==============+========+=============================================================================================================+
  | xml_id       | string | Unique string that describes this delivery service.                                                         |
  +--------------+--------+-------------------------------------------------------------------------------------------------------------+
  | serverNames  | string | array of hostname of cache servers to assign to this deliveryservice, for example: [ "server1", "server2" ] |
  +--------------+--------+-------------------------------------------------------------------------------------------------------------+


   **Response Example** ::

    {
        "response":{
            "serverNames":[
                "tc1_ats1"
            ],
            "xmlId":"my_ds_1"
        }
    }

|

