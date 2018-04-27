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

.. _to-api-v12-server:

Server
======

.. _to-api-v12-servers-route:

/api/1.2/servers
++++++++++++++++

**GET /api/1.2/servers**

  Retrieves properties of CDN servers.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Query Parameters**

  +--------------------+----------+---------------------------------------------+
  |   Name             | Required |                Description                  |
  +====================+==========+=============================================+
  | ``dsId``           | no       | Used to filter servers by delivery service. |
  +--------------------+----------+---------------------------------------------+
  | ``status``         | no       | Used to filter servers by status.           |
  +--------------------+----------+---------------------------------------------+
  | ``type``           | no       | Used to filter servers by type.             |
  +--------------------+----------+---------------------------------------------+
  | ``profileId``      | no       | Used to filter servers by profile ID.       |
  +--------------------+----------+---------------------------------------------+
  | ``cdn``            | no       | Used to filter servers by CDN ID.           |
  +--------------------+----------+---------------------------------------------+
  | ``cachegroup``     | no       | Used to filter servers by cache group ID.   |
  +--------------------+----------+---------------------------------------------+
  | ``physLocation``   | no       | Used to filter servers by phys location ID. |
  +--------------------+----------+---------------------------------------------+

  **Response Properties**

  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  |     Parameter      |  Type  |                                                Description                                                 |
  +====================+========+============================================================================================================+
  | ``cachegroup``     | string | The cache group name (see :ref:`to-api-v11-cachegroup`).                                                   |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``cachegroupId``   | string | The cache group id.                                                                                        |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``cdnId``          | string | Id of the CDN to which the server belongs to.                                                              |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``cdnName``        | string | Name of the CDN to which the server belongs to.                                                            |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``domainName``     | string | The domain name part of the FQDN of the cache.                                                             |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``guid``           | string | An identifier used to uniquely identify the server.                                                        |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``hostName``       | string | The host name part of the cache.                                                                           |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``httpsPort``      | string | The HTTPS port on which the main application listens (443 in most cases).                                  |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``id``             | string | The server id (database row number).                                                                       |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``iloIpAddress``   | string | The IPv4 address of the lights-out-management port.                                                        |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``iloIpGateway``   | string | The IPv4 gateway address of the lights-out-management port.                                                |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``iloIpNetmask``   | string | The IPv4 netmask of the lights-out-management port.                                                        |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``iloPassword``    | string | The password of the of the lights-out-management user (displays as ****** unless you are an 'admin' user). |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``iloUsername``    | string | The user name for lights-out-management.                                                                   |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``interfaceMtu``   | string | The Maximum Transmission Unit (MTU) to configure for ``interfaceName``.                                    |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``interfaceName``  | string | The network interface name used for serving traffic.                                                       |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``ip6Address``     | string | The IPv6 address/netmask for ``interfaceName``.                                                            |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``ip6Gateway``     | string | The IPv6 gateway for ``interfaceName``.                                                                    |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``ipAddress``      | string | The IPv4 address for ``interfaceName``.                                                                    |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``ipGateway``      | string | The IPv4 gateway for ``interfaceName``.                                                                    |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``ipNetmask``      | string | The IPv4 netmask for ``interfaceName``.                                                                    |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``lastUpdated``    | string | The Time and Date for the last update for this server.                                                     |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``mgmtIpAddress``  | string | The IPv4 address of the management port (optional).                                                        |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``mgmtIpGateway``  | string | The IPv4 gateway of the management port (optional).                                                        |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``mgmtIpNetmask``  | string | The IPv4 netmask of the management port (optional).                                                        |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``offlineReason``  | string | A user-entered reason why the server is in ADMIN_DOWN or OFFLINE status.                                   |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``physLocation``   | string | The physical location name (see :ref:`to-api-v11-phys-loc`).                                               |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``physLocationId`` | string | The physical location id (see :ref:`to-api-v11-phys-loc`).                                                 |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``profile``        | string | The assigned profile name (see :ref:`to-api-v11-profile`).                                                 |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``profileDesc``    | string | The assigned profile description (see :ref:`to-api-v11-profile`).                                          |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``profileId``      | string | The assigned profile Id (see :ref:`to-api-v11-profile`).                                                   |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``rack``           | string | A string indicating rack location.                                                                         |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``routerHostName`` | string | The human readable name of the router.                                                                     |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``routerPortName`` | string | The human readable name of the router port.                                                                |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``status``         | string | The Status string (See :ref:`to-api-v11-status`).                                                          |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``statusId``       | string | The Status id (See :ref:`to-api-v11-status`).                                                              |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``tcpPort``        | string | The default TCP port on which the main application listens (80 for a cache in most cases).                 |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``type``           | string | The name of the type of this server (see :ref:`to-api-v11-type`).                                          |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``typeId``         | string | The id of the type of this server (see :ref:`to-api-v11-type`).                                            |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``updPending``     |  bool  |                                                                                                            |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+

  **Response Example** ::

   {
      "response": [
          {
              "cachegroup": "us-il-chicago",
              "cachegroupId": "3",
              "cdnId": "3",
              "cdnName": "CDN-1",
              "domainName": "chi.kabletown.net",
              "guid": null,
              "hostName": "atsec-chi-00",
              "id": "19",
              "iloIpAddress": "172.16.2.6",
              "iloIpGateway": "172.16.2.1",
              "iloIpNetmask": "255.255.255.0",
              "iloPassword": "********",
              "iloUsername": "",
              "interfaceMtu": "9000",
              "interfaceName": "bond0",
              "ip6Address": "2033:D0D0:3300::2:2/64",
              "ip6Gateway": "2033:D0D0:3300::2:1",
              "ipAddress": "10.10.2.2",
              "ipGateway": "10.10.2.1",
              "ipNetmask": "255.255.255.0",
              "lastUpdated": "2015-03-08 15:57:32",
              "mgmtIpAddress": "",
              "mgmtIpGateway": "",
              "mgmtIpNetmask": "",
              "offlineReason": "N/A",
              "physLocation": "plocation-chi-1",
              "physLocationId": "9",
              "profile": "EDGE1_CDN1_421_SSL",
              "profileDesc": "EDGE1_CDN1_421_SSL profile",
              "profileId": "12",
              "rack": "RR 119.02",
              "routerHostName": "rtr-chi.kabletown.net",
              "routerPortName": "2",
              "status": "ONLINE",
              "statusId": "6",
              "tcpPort": "80",
              "httpsPort": "443",
              "type": "EDGE",
              "typeId": "3",
              "updPending": false
          },
          {
          ... more server data
          }
        ]
    }

|

**GET /api/1.2/servers/:id**

  Retrieves properties of a CDN server by server ID.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +-----------+----------+---------------------------------------------+
  |   Name    | Required |                Description                  |
  +===========+==========+=============================================+
  |   ``id``  |   yes    | Server id.                                  |
  +-----------+----------+---------------------------------------------+

  **Response Properties**

  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  |     Parameter      |  Type  |                                                Description                                                 |
  +====================+========+============================================================================================================+
  | ``cachegroup``     | string | The cache group name (see :ref:`to-api-v11-cachegroup`).                                                   |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``cachegroupId``   | string | The cache group id.                                                                                        |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``cdnId``          | string | Id of the CDN to which the server belongs to.                                                              |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``cdnName``        | string | Name of the CDN to which the server belongs to.                                                            |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``domainName``     | string | The domain name part of the FQDN of the cache.                                                             |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``guid``           | string | An identifier used to uniquely identify the server.                                                        |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``hostName``       | string | The host name part of the cache.                                                                           |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``httpsPort``      | string | The HTTPS port on which the main application listens (443 in most cases).                                  |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``id``             | string | The server id (database row number).                                                                       |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``iloIpAddress``   | string | The IPv4 address of the lights-out-management port.                                                        |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``iloIpGateway``   | string | The IPv4 gateway address of the lights-out-management port.                                                |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``iloIpNetmask``   | string | The IPv4 netmask of the lights-out-management port.                                                        |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``iloPassword``    | string | The password of the of the lights-out-management user (displays as ****** unless you are an 'admin' user). |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``iloUsername``    | string | The user name for lights-out-management.                                                                   |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``interfaceMtu``   | string | The Maximum Transmission Unit (MTU) to configure for ``interfaceName``.                                    |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``interfaceName``  | string | The network interface name used for serving traffic.                                                       |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``ip6Address``     | string | The IPv6 address/netmask for ``interfaceName``.                                                            |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``ip6Gateway``     | string | The IPv6 gateway for ``interfaceName``.                                                                    |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``ipAddress``      | string | The IPv4 address for ``interfaceName``.                                                                    |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``ipGateway``      | string | The IPv4 gateway for ``interfaceName``.                                                                    |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``ipNetmask``      | string | The IPv4 netmask for ``interfaceName``.                                                                    |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``lastUpdated``    | string | The Time and Date for the last update for this server.                                                     |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``mgmtIpAddress``  | string | The IPv4 address of the management port (optional).                                                        |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``mgmtIpGateway``  | string | The IPv4 gateway of the management port (optional).                                                        |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``mgmtIpNetmask``  | string | The IPv4 netmask of the management port (optional).                                                        |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``offlineReason``  | string | A user-entered reason why the server is in ADMIN_DOWN or OFFLINE status.                                   |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``physLocation``   | string | The physical location name (see :ref:`to-api-v11-phys-loc`).                                               |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``physLocationId`` | string | The physical location id (see :ref:`to-api-v11-phys-loc`).                                                 |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``profile``        | string | The assigned profile name (see :ref:`to-api-v11-profile`).                                                 |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``profileDesc``    | string | The assigned profile description (see :ref:`to-api-v11-profile`).                                          |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``profileId``      | string | The assigned profile Id (see :ref:`to-api-v11-profile`).                                                   |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``rack``           | string | A string indicating rack location.                                                                         |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``routerHostName`` | string | The human readable name of the router.                                                                     |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``routerPortName`` | string | The human readable name of the router port.                                                                |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``status``         | string | The Status string (See :ref:`to-api-v11-status`).                                                          |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``statusId``       | string | The Status id (See :ref:`to-api-v11-status`).                                                              |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``tcpPort``        | string | The default TCP port on which the main application listens (80 for a cache in most cases).                 |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``type``           | string | The name of the type of this server (see :ref:`to-api-v11-type`).                                          |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``typeId``         | string | The id of the type of this server (see :ref:`to-api-v11-type`).                                            |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``updPending``     |  bool  |                                                                                                            |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+

  **Response Example** ::

   {
      "response": [
          {
              "cachegroup": "us-il-chicago",
              "cachegroupId": "3",
              "cdnId": "3",
              "cdnName": "CDN-1",
              "domainName": "chi.kabletown.net",
              "guid": null,
              "hostName": "atsec-chi-00",
              "id": "19",
              "iloIpAddress": "172.16.2.6",
              "iloIpGateway": "172.16.2.1",
              "iloIpNetmask": "255.255.255.0",
              "iloPassword": "********",
              "iloUsername": "",
              "interfaceMtu": "9000",
              "interfaceName": "bond0",
              "ip6Address": "2033:D0D0:3300::2:2/64",
              "ip6Gateway": "2033:D0D0:3300::2:1",
              "ipAddress": "10.10.2.2",
              "ipGateway": "10.10.2.1",
              "ipNetmask": "255.255.255.0",
              "lastUpdated": "2015-03-08 15:57:32",
              "mgmtIpAddress": "",
              "mgmtIpGateway": "",
              "mgmtIpNetmask": "",
              "offlineReason": "N/A",
              "physLocation": "plocation-chi-1",
              "physLocationId": "9",
              "profile": "EDGE1_CDN1_421_SSL",
              "profileDesc": "EDGE1_CDN1_421_SSL profile",
              "profileId": "12",
              "rack": "RR 119.02",
              "routerHostName": "rtr-chi.kabletown.net",
              "routerPortName": "2",
              "status": "ONLINE",
              "statusId": "6",
              "tcpPort": "80",
              "httpsPort": "443",
              "type": "EDGE",
              "typeId": "3",
              "updPending": false
          }
        ]
    }

|


**GET /api/1.2/servers/:id/deliveryservices**

  Retrieves all delivery services assigned to the server. See also `Using Traffic Ops - Delivery Service <http://trafficcontrol.apache.org/docs/latest/admin/traffic_ops_using.html#delivery-service>`_.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +-----------------+----------+---------------------------------------------------+
  | Name            | Required | Description                                       |
  +=================+==========+===================================================+
  | ``id``          | yes      | Server ID.                                        |
  +-----------------+----------+---------------------------------------------------+

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
  | ``deepCachingType``      | string | When to do Deep Caching for this Delivery Service:                                                                                   |
  |                          |        |                                                                                                                                      |
  |                          |        | - NEVER (default)                                                                                                                    |
  |                          |        | - ALWAYS                                                                                                                             |
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
  | ``routingName``          | string | The routing name of this deliveryservice.                                                                                            |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``signed``               |  bool  | - false: token based auth (see :ref:token-based-auth) is not enabled for this deliveryservice.                                       |
  |                          |        | - true: token based auth is enabled for this deliveryservice.                                                                        |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``sslKeyVersion``        | string |                                                                                                                                      |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``tenant``               | string | Owning tenant name                                                                                                                   |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``tenantId``             | int    | Owning tenant ID.                                                                                                                    |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``trRequestHeaders``     | string | List of header keys separated by ``__RETURN__``. Listed headers will be included in TR access log entries under the "rh=" token.     |
  +--------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``trResponseHeaders``    | string | List of header ``name:value`` pairs separated by ``__RETURN__``. Listed pairs will be included in all TR HTTP responses.             |
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
            "deepCachingType": "NEVER",
            "displayName": "My Cool Delivery Service",
            "dnsBypassCname": "",
            "dnsBypassIp": "",
            "dnsBypassIp6": "",
            "dnsBypassTtl": "30",
            "dscp": "40",
            "edgeHeaderRewrite": null,
            "exampleURLs": [
                "http://foo.foo-ds.foo.bar.net"
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
            "routingName": "foo",
            "signed": false,
            "sslKeyVersion": "0",
            "tenant": "root",
            "tenantId": 1,
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


**GET /api/1.2/servers/totals**

  Retrieves a count of CDN servers by type.

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +-----------+--------+------------------------------------------------------------------------+
  | Parameter |  Type  |                             Description                                |
  +===========+========+========================================================================+
  | ``count`` | int    | The number of servers of this type in this instance of Traffic Ops.    |
  +-----------+--------+------------------------------------------------------------------------+
  | ``type``  | string | The name of the type of the server count (see :ref:`to-api-v12-type`). |
  +-----------+--------+------------------------------------------------------------------------+

  **Response Example** ::

    {
      "response": [
        {
          "count": 4,
          "type": "CCR"
        },
        {
          "count": 55,
          "type": "EDGE"
        },
        {
          "type": "MID",
          "count": 18
        },
        {
          "count": 0,
          "type": "INFLUXDB"
        },
        {
          "count": 4,
          "type": "RASCAL"
        }
    }

|

**GET /api/1.2/servers/status**

  Retrieves a count of CDN servers by status.

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +-----------------+--------+-----------------------------------------------------------------------------------------------------------------------+
  | Parameter       |  Type  |                             Description                                                                               |
  +=================+========+=======================================================================================================================+
  | ``ONLINE``      | int    | The number of ONLINE servers. Traffic Monitor will not monitor the state of ONLINE servers. True health is unknown.   |
  +-----------------+--------+-----------------------------------------------------------------------------------------------------------------------+
  | ``REPORTED``    | int    | The number of REPORTED servers. Traffic Monitor monitors the state of REPORTED servers and removes them if unhealthy. |
  +-----------------+--------+-----------------------------------------------------------------------------------------------------------------------+
  | ``OFFLINE``     | int    | The number of OFFLINE servers. Used for longer-term maintenance. These servers are excluded from CRConfig.json.       |
  +-----------------+--------+-----------------------------------------------------------------------------------------------------------------------+
  | ``ADMIN_DOWN``  | int    | The number of ADMIN_DOWN servers. Used for short-term maintenance. These servers are included in CRConfig.json.       |
  +-----------------+--------+-----------------------------------------------------------------------------------------------------------------------+
  | ``CCR_IGNORE``  | int    | The number of CCR_IGNORE servers. These servers are excluded from CRConfig.json.                                      |
  +-----------------+--------+-----------------------------------------------------------------------------------------------------------------------+
  | ``PRE_PROD``    | int    | The number of PRE_PROD servers. Used for servers to be deployed. These servers are excluded from CRConfig.json.       |
  +-----------------+--------+-----------------------------------------------------------------------------------------------------------------------+

  **Response Example** ::

    {
      "response":
        {
          "ONLINE": 100,
          "OFFLINE": 23,
          "REPORTED": 45,
          "ADMIN_DOWN": 4,
          "CCR_IGNORE": 1,
          "PRE_PROD": 0,
        }
    }

|


**GET /api/1.2/servers/hostname/:name/details**

  Retrieves the details of a server.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +----------+----------+----------------------------------+
  |   Name   | Required |           Description            |
  +==========+==========+==================================+
  | ``name`` | yes      | The host name part of the cache. |
  +----------+----------+----------------------------------+

  **Response Properties**

  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  |      Parameter       |  Type  |                                                 Description                                                 |
  +======================+========+=============================================================================================================+
  | ``cachegroup``       | string | The cache group name (see :ref:`to-api-v12-cachegroup`).                                                    |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``deliveryservices`` | array  | Array of strings with the delivery service ids assigned (see :ref:`to-api-v12-ds`).                         |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``domainName``       | string | The domain name part of the FQDN of the cache.                                                              |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``hardwareInfo``     | hash   | Hwinfo struct (see :ref:`to-api-v12-hwinfo`).                                                               |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``hostName``         | string | The host name part of the cache.                                                                            |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``id``               | string | The server id (database row number).                                                                        |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``iloIpAddress``     | string | The IPv4 address of the lights-out-management port.                                                         |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``iloIpGateway``     | string | The IPv4 gateway address of the lights-out-management port.                                                 |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``iloIpNetmask``     | string | The IPv4 netmask of the lights-out-management port.                                                         |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``iloPassword``      | string | The password of the of the lights-out-management user  (displays as ****** unless you are an 'admin' user). |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``iloUsername``      | string | The user name for lights-out-management.                                                                    |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``interfaceMtu``     | string | The Maximum Transmission Unit (MTU) to configure for ``interfaceName``.                                     |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``interfaceName``    | string | The network interface name used for serving traffic.                                                        |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``ip6Address``       | string | The IPv6 address/netmask for ``interfaceName``.                                                             |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``ip6Gateway``       | string | The IPv6 gateway for ``interfaceName``.                                                                     |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``ipAddress``        | string | The IPv4 address for ``interfaceName``.                                                                     |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``ipGateway``        | string | The IPv4 gateway for ``interfaceName``.                                                                     |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``ipNetmask``        | string | The IPv4 netmask for ``interfaceName``.                                                                     |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``lastUpdated``      | string | The Time/Date of the last update for this server.                                                           |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``mgmtIpAddress``    | string | The IPv4 address of the management port (optional).                                                         |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``mgmtIpGateway``    | string | The IPv4 gateway of the management port (optional).                                                         |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``mgmtIpNetmask``    | string | The IPv4 netmask of the management port (optional).                                                         |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``physLocation``     | string | The physical location name (see :ref:`to-api-v12-phys-loc`).                                                |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``profile``          | string | The assigned profile name (see :ref:`to-api-v12-profile`).                                                  |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``rack``             | string | A string indicating rack location.                                                                          |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``routerHostName``   | string | The human readable name of the router.                                                                      |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``routerPortName``   | string | The human readable name of the router port.                                                                 |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``status``           | string | The Status string (See :ref:`to-api-v12-status`).                                                           |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``tcpPort``          | string | The default TCP port on which the main application listens (80 for a cache in most cases).                  |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``httpsPort``        | string | The default HTTPS port on which the main application listens (443 for a cache in most cases).               |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``type``             | string | The name of the type of this server (see :ref:`to-api-v12-type`).                                           |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``xmppId``           | string | Deprecated.                                                                                                 |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``xmppPasswd``       | string | Deprecated.                                                                                                 |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+

  **Response Example** ::

    {
      "response": {
        "cachegroup": "us-il-chicago",
        "deliveryservices": [
          "1",
          "2",
          "3",
          "4"
        ],
        "domainName": "chi.kabletown.net",
        "hardwareInfo": {
          "Physical Disk 0:1:3": "D1S2",
          "Physical Disk 0:1:2": "D1S2",
          "Physical Disk 0:1:15": "D1S2",
          "Power Supply.Slot.2": "04.07.15",
          "Physical Disk 0:1:24": "YS08",
          "Physical Disk 0:1:1": "D1S2",
          "Model": "PowerEdge R720xd",
          "Physical Disk 0:1:22": "D1S2",
          "Physical Disk 0:1:18": "D1S2",
          "Enterprise UEFI Diagnostics": "4217A5",
          "Lifecycle Controller": "1.0.8.42",
          "Physical Disk 0:1:8": "D1S2",
          "Manufacturer": "Dell Inc.",
          "Physical Disk 0:1:6": "D1S2",
          "SysMemTotalSize": "196608",
          "PopulatedDIMMSlots": "24",
          "Physical Disk 0:1:20": "D1S2",
          "Intel(R) Ethernet 10G 2P X520 Adapter": "13.5.7",
          "Physical Disk 0:1:14": "D1S2",
          "BACKPLANE FIRMWARE": "1.00",
          "Dell OS Drivers Pack, 7.0.0.29, A00": "7.0.0.29",
          "Integrated Dell Remote Access Controller": "1.57.57",
          "Physical Disk 0:1:5": "D1S2",
          "ServiceTag": "D6XPDV1",
          "PowerState": "2",
          "Physical Disk 0:1:23": "D1S2",
          "Physical Disk 0:1:25": "D903",
          "BIOS": "1.3.6",
          "Physical Disk 0:1:12": "D1S2",
          "System CPLD": "1.0.3",
          "Physical Disk 0:1:4": "D1S2",
          "Physical Disk 0:1:0": "D1S2",
          "Power Supply.Slot.1": "04.07.15",
          "PERC H710P Mini": "21.0.2-0001",
          "PowerCap": "689",
          "Physical Disk 0:1:16": "D1S2",
          "Physical Disk 0:1:10": "D1S2",
          "Physical Disk 0:1:11": "D1S2",
          "Lifecycle Controller 2": "1.0.8.42",
          "BP12G+EXP 0:1": "1.07",
          "Physical Disk 0:1:9": "D1S2",
          "Physical Disk 0:1:17": "D1S2",
          "Broadcom Gigabit Ethernet BCM5720": "7.2.20",
          "Physical Disk 0:1:21": "D1S2",
          "Physical Disk 0:1:13": "D1S2",
          "Physical Disk 0:1:7": "D1S2",
          "Physical Disk 0:1:19": "D1S2"
        },
        "hostName": "atsec-chi-00",
        "id": "19",
        "iloIpAddress": "172.16.2.6",
        "iloIpGateway": "172.16.2.1",
        "iloIpNetmask": "255.255.255.0",
        "iloPassword": "********",
        "iloUsername": "",
        "interfaceMtu": "9000",
        "interfaceName": "bond0",
        "ip6Address": "2033:D0D0:3300::2:2/64",
        "ip6Gateway": "2033:D0D0:3300::2:1",
        "ipAddress": "10.10.2.2",
        "ipGateway": "10.10.2.1",
        "ipNetmask": "255.255.255.0",
        "mgmtIpAddress": "",
        "mgmtIpGateway": "",
        "mgmtIpNetmask": "",
        "physLocation": "plocation-chi-1",
        "profile": "EDGE1_CDN1_421_SSL",
        "rack": "RR 119.02",
        "routerHostName": "rtr-chi.kabletown.net",
        "routerPortName": "2",
        "status": "ONLINE",
        "tcpPort": "80",
        "httpsPort": "443",
        "type": "EDGE",
        "xmppId": "atsec-chi-00-dummyxmpp",
        "xmppPasswd": "X"

      }
    }

|

**POST /api/1.2/servercheck**

  Post a server check result to the serverchecks table.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +----------------------------+----------+-------------+
  |            Name            | Required | Description |
  +============================+==========+=============+
  | ``id``                     | yes      |             |
  +----------------------------+----------+-------------+
  | ``host_name``              | yes      |             |
  +----------------------------+----------+-------------+
  | ``servercheck_short_name`` | yes      |             |
  +----------------------------+----------+-------------+
  | ``value``                  | yes      |             |
  +----------------------------+----------+-------------+

  **Request Example** ::

    {
     "id": "",
     "host_name": "",
     "servercheck_short_name": "",
     "value": ""
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

    Response Example:

    {
      "alerts":
        [
          {
            "level": "success",
            "text": "Server Check was successfully updated."
          }
        ],
    }

|

**POST /api/1.2/servers**

  Allow user to create a server.

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Properties**

  +----------------+----------+-------------------------------------------------------------+
  |      Name      | Required |                  Description                                |
  +================+==========+=============================================================+
  | hostName       | yes      | The host name part of the server.                           |
  +----------------+----------+-------------------------------------------------------------+
  | domainName     | yes      | The domain name part of the FQDN of the cache.              |
  +----------------+----------+-------------------------------------------------------------+
  | cachegroupId   | yes      | Cache Group ID                                              |
  +----------------+----------+-------------------------------------------------------------+
  | interfaceName  | yes      | The interface name (e.g. eth0, p2p1).                       |
  +----------------+----------+-------------------------------------------------------------+
  | ipAddress      | yes      | Must be unique per server profile.                          |
  +----------------+----------+-------------------------------------------------------------+
  | ipNetmask      | yes      | The IPv4 Netmask                                            |
  +----------------+----------+-------------------------------------------------------------+
  | ipGateway      | yes      | The IPv4 Gateway.                                           |
  +----------------+----------+-------------------------------------------------------------+
  | interfaceMtu   | yes      | 1500 or 9000                                                |
  +----------------+----------+-------------------------------------------------------------+
  | physLocationId | yes      | The ID of the Physical Location.                            |
  +----------------+----------+-------------------------------------------------------------+
  | typeId         | yes      | The ID of the Server Type                                   |
  +----------------+----------+-------------------------------------------------------------+
  | profileId      | yes      | Profile ID - Profile's CDN must match server's.             |
  +----------------+----------+-------------------------------------------------------------+
  | cdnId          | yes      | CDN ID the server belongs to                                |
  +----------------+----------+-------------------------------------------------------------+
  | updPending     | yes      | Is there an update pending for this server. (true or false) |
  +----------------+----------+-------------------------------------------------------------+
  | statusId       | yes      | The Status ID of the server.                                |
  +----------------+----------+-------------------------------------------------------------+
  | tcpPort        | no       | Must be a valid TCP port if specified.                      |
  +----------------+----------+-------------------------------------------------------------+
  | httpsPort      | no       | Must be a valid TCP port if specified.                      |
  +----------------+----------+-------------------------------------------------------------+
  | xmppId         | no       |                                                             |
  +----------------+----------+-------------------------------------------------------------+
  | xmppPasswd     | no       |                                                             |
  +----------------+----------+-------------------------------------------------------------+
  | ip6Address     | no       | IPv6 address and prefix. Must be unique per server profile. |
  +----------------+----------+-------------------------------------------------------------+
  | ip6Gateway     | no       | IPv6 Gateway                                                |
  +----------------+----------+-------------------------------------------------------------+
  | rack           | no       | The rack location in the Data Center.                       |
  +----------------+----------+-------------------------------------------------------------+
  | mgmtIpAddress  | no       | The IPv4 management address.                                |
  +----------------+----------+-------------------------------------------------------------+
  | mgmtIpNetmask  | no       | The IPv4 management netmask.                                |
  +----------------+----------+-------------------------------------------------------------+
  | mgmtIpGateway  | no       | The IPv4 management gateway.                                |
  +----------------+----------+-------------------------------------------------------------+
  | iloIpAddress   | no       | The IPv4 ILO address.                                       |
  +----------------+----------+-------------------------------------------------------------+
  | iloIpNetmask   | no       | The IPv4 ILO netmask.                                       |
  +----------------+----------+-------------------------------------------------------------+
  | iloIpGateway   | no       | The IPv4 ILO gateway.                                       |
  +----------------+----------+-------------------------------------------------------------+
  | iloUsername    | no       | The ILO username.                                           |
  +----------------+----------+-------------------------------------------------------------+
  | iloPassword    | no       | The ILO password.                                           |
  +----------------+----------+-------------------------------------------------------------+
  | routerHostName | no       | The hostname of the router the server is connected to.      |
  +----------------+----------+-------------------------------------------------------------+
  | routerPortName | no       | The portname in the router.                                 |
  +----------------+----------+-------------------------------------------------------------+

  **Request Example** ::

    {
        "hostName": "tc1_ats1",
        "domainName": "cdn1.kabletown.test",
        "cachegroupId": 1,
        "cdnId": 1,
        "interfaceName": "eth0",
        "ipAddress": "10.74.27.188",
        "ipNetmask": "255.255.255.0",
        "ipGateway": "10.74.27.1",
        "interfaceMtu": 1500,
        "physLocationId": 1,
        "typeId": 1,
        "profileId": 1,
	"updPending": true,
	"statusId": 1
    }

|

  **Response Properties**

  +----------------+--------+------------------------------------------------+
  |      Name      |  Type  |                  Description                   |
  +================+========+================================================+
  | hostName       | string | The host name part of the server.              |
  +----------------+--------+------------------------------------------------+
  | Name           | string | Description                                    |
  +----------------+--------+------------------------------------------------+
  | domainName     | string | The domain name part of the FQDN of the cache. |
  +----------------+--------+------------------------------------------------+
  | cachegroup     | string | cache group name                               |
  +----------------+--------+------------------------------------------------+
  | interfaceName  | string |                                                |
  +----------------+--------+------------------------------------------------+
  | ipAddress      | string |                                                |
  +----------------+--------+------------------------------------------------+
  | ipNetmask      | string |                                                |
  +----------------+--------+------------------------------------------------+
  | ipGateway      | string |                                                |
  +----------------+--------+------------------------------------------------+
  | interfaceMtu   | string | 1500 or 9000                                   |
  +----------------+--------+------------------------------------------------+
  | physLocation   | string |                                                |
  +----------------+--------+------------------------------------------------+
  | type           | string | server type                                    |
  +----------------+--------+------------------------------------------------+
  | profile        | string |                                                |
  +----------------+--------+------------------------------------------------+
  | cdnName        | string | cdn name the server belongs to                 |
  +----------------+--------+------------------------------------------------+
  | tcpPort        | string |                                                |
  +----------------+--------+------------------------------------------------+
  | httpsPort      | string |                                                |
  +----------------+--------+------------------------------------------------+
  | xmppId         | string |                                                |
  +----------------+--------+------------------------------------------------+
  | xmppPasswd     | string |                                                |
  +----------------+--------+------------------------------------------------+
  | ip6Address     | string |                                                |
  +----------------+--------+------------------------------------------------+
  | ip6Gateway     | string |                                                |
  +----------------+--------+------------------------------------------------+
  | rack           | string |                                                |
  +----------------+--------+------------------------------------------------+
  | mgmtIpAddress  | string |                                                |
  +----------------+--------+------------------------------------------------+
  | mgmtIpNetmask  | string |                                                |
  +----------------+--------+------------------------------------------------+
  | mgmtIpGateway  | string |                                                |
  +----------------+--------+------------------------------------------------+
  | iloIpAddress   | string |                                                |
  +----------------+--------+------------------------------------------------+
  | iloIpNetmask   | string |                                                |
  +----------------+--------+------------------------------------------------+
  | iloIpGateway   | string |                                                |
  +----------------+--------+------------------------------------------------+
  | iloUsername    | string |                                                |
  +----------------+--------+------------------------------------------------+
  | iloPassword    | string |                                                |
  +----------------+--------+------------------------------------------------+
  | routerHostName | string |                                                |
  +----------------+--------+------------------------------------------------+
  | routerPortName | string |                                                |
  +----------------+--------+------------------------------------------------+

  **Response Example** ::

    {
        'response' : {
	    'profileId' : 1,
            'xmppPasswd' : '**********',
            'profile' : 'EDGE1_CDN1_421',
            'iloUsername' : 'username',
	    'statusId' : 1,
            'status' : 'REPORTED',
            'ipAddress' : '10.74.27.188',
            'cdnId' : 1,
            'physLocation' : 'plocation-chi-1',
            'cachegroup' : 'cache_group_edge',
            'interfaceName' : 'eth0',
            'ip6Gateway' : null,
            'iloPassword' : null,
            'id' : 1003,
            'routerPortName' : null,
            'lastUpdated' : '2016-01-25 14:16:16',
            'ipNetmask' : '255.255.255.0',
            'ipGateway' : '10.74.27.1',
            'tcpPort' : 80,
            'httpsPort' : 443,
            'mgmtIpAddress' : null,
            'ip6Address' : null,
            'interfaceMtu' : 1500,
            'iloIpGateway' : null,
            'hostName' : 'tc1_ats1',
            'xmppId' : 'tc1_ats1',
            'rack' : null,
            'mgmtIpNetmask' : null,
            'iloIpAddress' : null,
            'mgmtIpGateway' : null,
            'type' : 'EDGE',
            'domainName' : 'cdn1.kabletown.test',
            'iloIpNetmask' : null,
            'routerHostName' : null,
	    'updPending' : false,
	    'guid' : null,
	    'physLocationId' : 1,
	    'offlineReason' : 'N\/A',
	    'cachegroupId' : 1,
	    'typeId' : 1,
	    'cdnName' : 'cdn1',
	    'profileDesc' : 'The profile description'
        }
    }

|

**PUT /api/1.2/servers/{:id}**

  Allow user to edit server through api.

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Route Parameters**

  +------+----------+-------------------------------+
  | Name | Required | Description                   |
  +======+==========+===============================+
  | id   | yes      | The id of the server to edit. |
  +------+----------+-------------------------------+

  **Request Properties**

  +----------------+----------+-------------------------------------------------+
  |      Name      | Required |                  Description                    |
  +================+==========+=================================================+
  | hostName       | yes      | The host name part of the server.               |
  +----------------+----------+-------------------------------------------------+
  | domainName     | yes      | The domain name part of the FQDN of the cache.  |
  +----------------+----------+-------------------------------------------------+
  | cachegroup     | yes      | cache group name                                |
  +----------------+----------+-------------------------------------------------+
  | interfaceName  | yes      |                                                 |
  +----------------+----------+-------------------------------------------------+
  | ipAddress      | yes      | Must be unique per server profile.              |
  +----------------+----------+-------------------------------------------------+
  | ipNetmask      | yes      |                                                 |
  +----------------+----------+-------------------------------------------------+
  | ipGateway      | yes      |                                                 |
  +----------------+----------+-------------------------------------------------+
  | interfaceMtu   | no       | 1500 or 9000                                    |
  +----------------+----------+-------------------------------------------------+
  | physLocation   | yes      |                                                 |
  +----------------+----------+-------------------------------------------------+
  | type           | yes      | server type                                     |
  +----------------+----------+-------------------------------------------------+
  | profile        | yes      | Profile ID - Profile's CDN must match server's. |
  +----------------+----------+-------------------------------------------------+
  | cdnName        | yes      | cdn name the server belongs to                  |
  +----------------+----------+-------------------------------------------------+
  | tcpPort        | no       |                                                 |
  +----------------+----------+-------------------------------------------------+
  | httpsPort      | no       |                                                 |
  +----------------+----------+-------------------------------------------------+
  | xmppId         | no       |                                                 |
  +----------------+----------+-------------------------------------------------+
  | xmppPasswd     | no       |                                                 |
  +----------------+----------+-------------------------------------------------+
  | ip6Address     | no       | Must be unique per server profile.              |
  +----------------+----------+-------------------------------------------------+
  | ip6Gateway     | no       |                                                 |
  +----------------+----------+-------------------------------------------------+
  | rack           | no       |                                                 |
  +----------------+----------+-------------------------------------------------+
  | mgmtIpAddress  | no       |                                                 |
  +----------------+----------+-------------------------------------------------+
  | mgmtIpNetmask  | no       |                                                 |
  +----------------+----------+-------------------------------------------------+
  | mgmtIpGateway  | no       |                                                 |
  +----------------+----------+-------------------------------------------------+
  | iloIpAddress   | no       |                                                 |
  +----------------+----------+-------------------------------------------------+
  | iloIpNetmask   | no       |                                                 |
  +----------------+----------+-------------------------------------------------+
  | iloIpGateway   | no       |                                                 |
  +----------------+----------+-------------------------------------------------+
  | iloUsername    | no       |                                                 |
  +----------------+----------+-------------------------------------------------+
  | iloPassword    | no       |                                                 |
  +----------------+----------+-------------------------------------------------+
  | routerHostName | no       |                                                 |
  +----------------+----------+-------------------------------------------------+
  | routerPortName | no       |                                                 |
  +----------------+----------+-------------------------------------------------+

  **Request Example** ::

    {
        "hostName": "tc1_ats2",
        "domainName": "my.test.com",
        "cachegroup": "cache_group_edge",
        "cdnName": "cdn_number_1",
        "interfaceName": "eth0",
        "ipAddress": "10.74.27.188",
        "ipNetmask": "255.255.255.0",
        "ipGateway": "10.74.27.1",
        "interfaceMtu": "1500",
        "physLocation": "plocation-chi-1",
        "type": "EDGE",
        "profile": "EDGE1_CDN1_421"
    }

|

  **Response Properties**

  +----------------+--------+------------------------------------------------+
  |      Name      |  Type  |                  Description                   |
  +================+========+================================================+
  | hostName       | string | The host name part of the server.              |
  +----------------+--------+------------------------------------------------+
  | Name           | string | Description                                    |
  +----------------+--------+------------------------------------------------+
  | domainName     | string | The domain name part of the FQDN of the cache. |
  +----------------+--------+------------------------------------------------+
  | cachegroup     | string | cache group name                               |
  +----------------+--------+------------------------------------------------+
  | interfaceName  | string |                                                |
  +----------------+--------+------------------------------------------------+
  | ipAddress      | string |                                                |
  +----------------+--------+------------------------------------------------+
  | ipNetmask      | string |                                                |
  +----------------+--------+------------------------------------------------+
  | ipGateway      | string |                                                |
  +----------------+--------+------------------------------------------------+
  | interfaceMtu   | string | 1500 or 9000                                   |
  +----------------+--------+------------------------------------------------+
  | physLocation   | string |                                                |
  +----------------+--------+------------------------------------------------+
  | type           | string | server type                                    |
  +----------------+--------+------------------------------------------------+
  | profile        | string |                                                |
  +----------------+--------+------------------------------------------------+
  | cdnName        | string | cdn name the server belongs to                 |
  +----------------+--------+------------------------------------------------+
  | tcpPort        | string |                                                |
  +----------------+--------+------------------------------------------------+
  | httpsPort      | string |                                                |
  +----------------+--------+------------------------------------------------+
  | xmppId         | string |                                                |
  +----------------+--------+------------------------------------------------+
  | xmppPasswd     | string |                                                |
  +----------------+--------+------------------------------------------------+
  | ip6Address     | string |                                                |
  +----------------+--------+------------------------------------------------+
  | ip6Gateway     | string |                                                |
  +----------------+--------+------------------------------------------------+
  | rack           | string |                                                |
  +----------------+--------+------------------------------------------------+
  | mgmtIpAddress  | string |                                                |
  +----------------+--------+------------------------------------------------+
  | mgmtIpNetmask  | string |                                                |
  +----------------+--------+------------------------------------------------+
  | mgmtIpGateway  | string |                                                |
  +----------------+--------+------------------------------------------------+
  | iloIpAddress   | string |                                                |
  +----------------+--------+------------------------------------------------+
  | iloIpNetmask   | string |                                                |
  +----------------+--------+------------------------------------------------+
  | iloIpGateway   | string |                                                |
  +----------------+--------+------------------------------------------------+
  | iloUsername    | string |                                                |
  +----------------+--------+------------------------------------------------+
  | iloPassword    | string |                                                |
  +----------------+--------+------------------------------------------------+
  | routerHostName | string |                                                |
  +----------------+--------+------------------------------------------------+
  | routerPortName | string |                                                |
  +----------------+--------+------------------------------------------------+

  **Response Example** ::

    {
        'response' : {
            'xmppPasswd' : '**********',
            'profile' : 'EDGE1_CDN1_421',
            'iloUsername' : null,
            'status' : 'REPORTED',
            'ipAddress' : '10.74.27.188',
            'cdnId' : '1',
            'physLocation' : 'plocation-chi-1',
            'cachegroup' : 'cache_group_edge',
            'interfaceName' : 'eth0',
            'ip6Gateway' : null,
            'iloPassword' : null,
            'id' : '1003',
            'routerPortName' : null,
            'lastUpdated' : '2016-01-25 14:16:16',
            'ipNetmask' : '255.255.255.0',
            'ipGateway' : '10.74.27.1',
            'tcpPort' : '80',
            'httpsPort' : '443',
            'mgmtIpAddress' : null,
            'ip6Address' : null,
            'interfaceMtu' : '1500',
            'iloIpGateway' : null,
            'hostName' : 'tc1_ats2',
            'xmppId' : 'tc1_ats1',
            'rack' : null,
            'mgmtIpNetmask' : null,
            'iloIpAddress' : null,
            'mgmtIpGateway' : null,
            'type' : 'EDGE',
            'domainName' : 'my.test.com',
            'iloIpNetmask' : null,
            'routerHostName' : null
        }
    }

|

**PUT /api/1.2/servers/{:id}/status**

  Updates server status and queues updates on all child caches if server type is EDGE or MID. Also, captures offline reason if status is set to ADMIN_DOWN or OFFLINE and prepends offline reason with the user that initiated the status change.

  Authentication Required: Yes

  Role(s) Required: Admin or Operations

  **Request Route Parameters**

  +------+----------+-------------------------------+
  | Name | Required | Description                   |
  +======+==========+===============================+
  | id   | yes      | The id of the server.         |
  +------+----------+-------------------------------+

  **Request Properties**

  +----------------+----------+-------------------------------------------------+
  |      Name      | Required |                  Description                    |
  +================+==========+=================================================+
  | status         | yes      | Status ID or name.                              |
  +----------------+----------+-------------------------------------------------+
  | offlineReason  | yes|no   | Required if status is ADMIN_DOWN or OFFLINE.    |
  +----------------+----------+-------------------------------------------------+

  **Request Example** ::

    {
        "status": "ADMIN_DOWN",
        "offlineReason": "Bad drives"
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

  **Response Example** ::

    {
          "alerts": [
                    {
                            "level": "success",
                            "text": "Updated status [ ADMIN_DOWN ] for foo.bar.net [ user23: bad drives ] and queued updates on all child caches"
                    }
            ],
    }

|

**DELETE /api/1.2/servers/{:id}**

  Allow user to delete server through api.

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Route Parameters**

  +------+----------+---------------------------------+
  | Name | Required | Description                     |
  +======+==========+=================================+
  | id   | yes      | The id of the server to delete. |
  +------+----------+---------------------------------+

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
                            "text": "Server was deleted."
                    }
            ],
    }

|

**POST /api/1.2/servers/{:id}/queue_update**

  Queue or dequeue updates for a specific server.

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Route Parameters**

  +-----------+----------+------------------+
  | Name      | Required | Description      |
  +===========+==========+==================+
  | id        | yes      | the server id.   |
  +-----------+----------+------------------+

  **Request Properties**

  +--------------+---------+-----------------------------------------------+
  | Name         | Type    | Description                                   |
  +==============+=========+===============================================+
  | action       | string  | queue or dequeue                              |
  +--------------+---------+-----------------------------------------------+

  **Response Properties**

  +--------------+---------+-----------------------------------------------+
  | Name         | Type    | Description                                   |
  +==============+=========+===============================================+
  | action       | string  | The action processed, queue or dequeue.       |
  +--------------+---------+-----------------------------------------------+
  | serverId     | integer | server id                                     |
  +--------------+---------+-----------------------------------------------+

  **Response Example** ::

    {
      "response": {
          "serverId": "1",
          "action": "queue"
      }
    }

|

