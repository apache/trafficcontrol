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

.. _tr-api:

******************
Traffic Router API
******************
By default, Traffic Router serves its API via HTTP (not HTTPS) on port 3333. This can be configured in :file:`/opt/traffic_router/conf/server.xml` or by setting a :term:`Parameter` with the :ref:`parameter-name` "api.port", and the :ref:`parameter-config-file` "server.xml" on the Traffic Router's :term:`Profile`.

The API can be configured via HTTPS on port 3443 in :file:`/opt/traffic_router/conf/server.xml` or by setting a :term:`Parameter` named ``secure.api.port`` with ``configFile`` ``server.xml`` on the Traffic Router's :term:`Profile`.  The post install script will generate self signed certificates at ``/opt/traffic_router/conf/``, create a new Java Keystore named :file:`/opt/traffic_router/conf/keyStore.jks`, and add the new certificate to the Keystore.  The password for the Java Keystore and the Keystore location are stored in :file:`/opt/traffic_router/conf/https.properties`.
To override the self signed certificates with new ones from a certificate authority, update the properties for the Keystore location and password at :file:`/opt/traffic_router/conf/https.properties`.


The API can be configured via HTTPS on port 3443 in :file:`/opt/traffic_router/conf/server.xml` or by setting a :term:`Parameter` named ``secure.api.port`` with ``configFile`` ``server.xml`` on the Traffic Router's :term:`Profile`.  When ``systemctl start traffic_router`` is run, it will generate self signed certificates at ``/opt/traffic_router/conf/``, create a new Java Keystore named :file:`/opt/traffic_router/conf/keyStore.jks`, and add the new certificate to the Keystore.  The password for the Java Keystore and the Keystore location are stored in :file:`/opt/traffic_router/conf/https.properties`.
To override the self signed certificates with new ones from a certificate authority, either replace the Java Keystore in the default location or update the properties for the new Keystore location and password at :file:`/opt/traffic_router/conf/https.properties` and then restart the Traffic Router using ``systemctl``.

Other attributes of the default certificate can also be customized by specifying appropriate values for the following properties in :file:`/opt/traffic_router/conf/https.properties`. These properties are listed below:

.. table:: HTTPS Certificate Attributes

	+------------------------------------------+------------------------------------------------------------------------+---------------------------------------------------------+
	| Name                                     | Description                                                            | Default                                                 |
	+==========================================+========================================================================+=========================================================+
	|  https.certificate.location              | The location of the certificate key store                              |                                                         |
	+------------------------------------------+------------------------------------------------------------------------+---------------------------------------------------------+
	|  https.password                          | The password for the certificate key store                             |                                                         |
	+------------------------------------------+------------------------------------------------------------------------+---------------------------------------------------------+
	|  https.key.size                          | The size for the HTTPS keys                                            | 2048                                                    |
	+------------------------------------------+------------------------------------------------------------------------+---------------------------------------------------------+
	|  https.signature.algorithm               | The HTTPS signing algorithm to be used                                 | SHA1WithRSA                                             |
	+------------------------------------------+------------------------------------------------------------------------+---------------------------------------------------------+
	|  https.validity.years                    | The amount of time (in years) for which the cert is valid              | 3                                                       |
	+------------------------------------------+------------------------------------------------------------------------+---------------------------------------------------------+
	|  https.certificate.country               | The country of the certificate                                         | US                                                      |
	+------------------------------------------+------------------------------------------------------------------------+---------------------------------------------------------+
	|  https.certificate.state                 | The state of the certificate                                           | CO                                                      |
	+------------------------------------------+------------------------------------------------------------------------+---------------------------------------------------------+
	|  https.certificate.locality              | The locality of the certificate                                        | Denver                                                  |
	+------------------------------------------+------------------------------------------------------------------------+---------------------------------------------------------+
	|  https.certificate.organization          | The organization of the certificate                                    | Apache Traffic Control                                  |
	+------------------------------------------+------------------------------------------------------------------------+---------------------------------------------------------+
	|  https.certificate.organizational.unit   | The organizational unit of the certificate                             | Apache Foundation, Hosted by Traffic Control, CDNDefault|
	+------------------------------------------+------------------------------------------------------------------------+---------------------------------------------------------+

Traffic Router API endpoints only respond to ``GET`` requests.

.. _tr-api-crs-stats:

``/crs/stats``
==============
General stats.

Request Structure
-----------------
.. code-block:: http
	:caption: Request Example

	GET /crs/stats HTTP/1.1
	Host: trafficrouter.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Type: application/json;charset=UTF-8
	Transfer-Encoding: chunked
	Date: Tue, 15 Jan 2019 21:02:09 GMT

	{ "app": {
		"buildTimestamp": "2019-01-10",
		"name": "traffic_router",
		"deploy-dir": "/opt/traffic_router",
		"git-revision": "437e9df81",
		"version": "3.0.0"
	},
	"stats": {
		"dnsMap": {},
		"httpMap": {
			"video.demo1.mycdn.ciab.test": {
				"czCount": 0,
				"geoCount": 0,
				"deepCzCount": 0,
				"missCount": 0,
				"dsrCount": 0,
				"errCount": 0,
				"staticRouteCount": 0,
				"fedCount": 0,
				"regionalDeniedCount": 0,
				"regionalAlternateCount": 0
			}
		},
		"totalDnsCount": 0,
		"totalHttpCount": 1,
		"totalDsMissCount": 0,
		"appStartTime": 1547584831677,
		"averageDnsTime": 0,
		"averageHttpTime": 1547584863270,
		"updateTracker": {
			"lastHttpsCertificatesCheck": 1547586068932,
			"lastGeolocationDatabaseUpdaterUpdate": 1547584858917,
			"lastCacheStateCheck": 1547586128932,
			"lastCacheStateChange": 1547584867102,
			"lastNetworkUpdaterUpdate": 1547584857484,
			"lastHttpsCertificatesUpdate": 1547586071079,
			"lastSteeringWatcherUpdate": 1547584923514,
			"lastConfigCheck": 1547586127344,
			"lastConfigChange": 1547584863406,
			"lastNetworkUpdaterCheck": 1547584857465,
			"lastGeolocationDatabaseUpdaterCheck": 1547584858906,
			"lastFederationsWatcherUpdate": 1547584863433,
			"lastHttpsCertificatesFetchSuccess": 1547586069070,
			"lastSteeringWatcherCheck": 1547586124630,
			"lastFederationsWatcherCheck": 1547586124584,
			"lastHttpsCertificatesFetchAttempt": 1547586068932
		}
	}}

.. _tr-api-crs-stats-ip-ip:

``/crs/stats/ip/{{IP}}``
================================
Geolocation information for an IPv4 or IPv6 address.

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------------------+
	| Name | Description                                                            |
	+======+========================================================================+
	|  IP  | The IP address for which statics will be returned. May be IPv4 or IPv6 |
	+------+------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /crs/stats/ip/255.255.255.255 HTTP/1.1
	Host: trafficrouter.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Disposition: inline;filename=f.txt
	Content-Type: application/json;charset=UTF-8
	Transfer-Encoding: chunked
	Date: Tue, 15 Jan 2019 21:06:09 GMT

	{ "locationByGeo": {
		"city": "Woodridge",
		"countryCode": "US",
		"latitude": "41.7518",
		"postalCode": "60517",
		"countryName": "United States",
		"longitude": "-88.0489"
	},
	"locationByFederation": "not found",
	"requestIp": "69.241.118.34",
	"locationByCoverageZone": "not found",
	"locationByDeepCoverageZone": "not found"
	}

.. _tr-api-crs-locations:

``/crs/locations``
==================
A list of configured :term:`Cache Groups` to which the Traffic Router is capable of routing client traffic.

Request Structure
-----------------
.. code-block:: http
	:caption: Request Example

	GET /crs/locations HTTP/1.1
	Host: trafficrouter.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*

Response Structure
------------------
:locations: An array of strings that are the :ref:`Names of Cache Groups <cache-group-name>` to which this Traffic Router is capable of routing client traffic

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Type: application/json;charset=UTF-8
	Transfer-Encoding: chunked
	Date: Tue, 15 Jan 2019 21:12:17 GMT

	{ "locations": [
		"CDN_in_a_Box_Edge"
	]}

.. _tr-api-crs-locations-caches:

``/crs/locations/caches``
=========================
A mapping of caches to cache groups and their current health state.

Request Structure
-----------------
.. code-block:: http
	:caption: Request Example

	GET /crs/locations/caches HTTP/1.1
	Host: trafficrouter.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Type: application/json;charset=UTF-8
	Transfer-Encoding: chunked
	Date: Tue, 15 Jan 2019 21:15:53 GMT

	{ "locations": {
		"CDN_in_a_Box_Edge": [
			{
				"cacheId": "edge",
				"fqdn": "edge.infra.ciab.test",
				"ipAddresses": [
					"172.16.239.100",
					"fc01:9400:1000:8:0:0:0:100"
				],
				"port": 0,
				"adminStatus": null,
				"lastUpdateHealthy": false,
				"lastUpdateTime": 0,
				"connections": 0,
				"currentBW": 0,
				"availBW": 0,
				"cacheOnline": true
			}
		]
	}}

.. _tr-api-crs-locations-cachegroup-caches:

``/crs/locations/{{cachegroup}}/caches``
========================================
A list of :term:`cache servers` for this :term:`Cache Group` only.

Request Structure
-----------------
.. table:: Request Path Parameters

	+------------+------------------------------------------------------------------------------------------------------------------------------+
	| Name       | Description                                                                                                                  |
	+============+==============================================================================================================================+
	| cachegroup | The :ref:`Name of a Cache Group <cache-group-name>` of which a list of constituent :term:`cache servers` will be retrieved   |
	+------------+------------------------------------------------------------------------------------------------------------------------------+


.. code-block:: http
	:caption: Request Example

	GET /crs/locations/CDN_in_a_Box_Edge/caches HTTP/1.1
	Host: trafficrouter.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Type: application/json;charset=UTF-8
	Transfer-Encoding: chunked
	Date: Tue, 15 Jan 2019 21:18:25 GMT

	{ "caches": [
		{
			"cacheId": "edge",
			"fqdn": "edge.infra.ciab.test",
			"ipAddresses": [
				"172.16.239.100",
				"fc01:9400:1000:8:0:0:0:100"
			],
			"port": 0,
			"adminStatus": null,
			"lastUpdateHealthy": false,
			"lastUpdateTime": 0,
			"connections": 0,
			"currentBW": 0,
			"availBW": 0,
			"cacheOnline": true
		}
	]}

.. _tr-api-crs-consistenthash-cache-coveragezone:

``/crs/consistenthash/cache/coveragezone``
===========================================
The resulting cache of the consistent hash using coverage zone file for a given client IP, :term:`Delivery Service`, and request path.

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| Name              | Required | Description                                                                                                  |
	+===================+==========+==============================================================================================================+
	| ip                | yes      | The IP address of a potential client                                                                         |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| deliveryServiceId | yes      | The integral, unique identifier?/'xml_id'?/name? of a :term:`Delivery Service` served by this Traffic Router |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| requestPath       | yes      | The... request path?                                                                                         |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+

Response Structure
------------------
TBD

.. _tr-api-crs-consistenthash-cache-deep-coveragezone:

``/crs/consistenthash/cache/deep/coveragezone``
===============================================
The resulting cache of the consistent hash using deep coverage zone file (deep caching) for a given client IP, :term:`Delivery Service`, and request path.

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| Name              | Required | Description                                                                                                  |
	+===================+==========+==============================================================================================================+
	| ip                | yes      | The IP address of a potential client                                                                         |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| deliveryServiceId | yes      | The integral, unique identifier?/'xml_id'?/name? of a :term:`Delivery Service` served by this Traffic Router |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| requestPath       | yes      | The... request path?                                                                                         |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+

Response Structure
------------------
TBD

.. _tr-api-crs-consistenthash-cache-geolocation:

``/crs/consistenthash/cache/geolocation``
=========================================
The resulting cache of the consistent hash using geographic location for a given client IP, :term:`Delivery Service`, and request path.

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| Name              | Required | Description                                                                                                  |
	+===================+==========+==============================================================================================================+
	| ip                | yes      | The IP address of a potential client                                                                         |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| deliveryServiceId | yes      | The integral, unique identifier?/'xml_id'?/name? of a :term:`Delivery Service` served by this Traffic Router |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| requestPath       | yes      | The... request path?                                                                                         |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+

Response Structure
------------------
TBD

.. _tr-api-crs-consistenthash-deliveryservice:

``/crs/consistenthash/deliveryservice/``
========================================
The resulting :term:`Delivery Service` of the consistent hash for a given :term:`Delivery Service` and request path -- used to test STEERING :term:`Delivery Services`.

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| Name              | Required | Description                                                                                                  |
	+===================+==========+==============================================================================================================+
	| deliveryServiceId | yes      | The integral, unique identifier?/'xml_id'?/name? of a :term:`Delivery Service` served by this Traffic Router |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| requestPath       | yes      | The... request path?                                                                                         |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /crs/consistenthash/deliveryservice?deliveryServiceId=demo1&requestPath=/ HTTP/1.1
	Host: trafficrouter.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Type: application/json;charset=UTF-8
	Transfer-Encoding: chunked
	Date: Tue, 15 Jan 2019 21:40:51 GMT

	{ "id": "demo1",
	"coverageZoneOnly": false,
	"geoRedirectUrl": null,
	"geoRedirectFile": null,
	"geoRedirectUrlType": "INVALID_URL",
	"routingName": "video",
	"missLocation": {
		"latitude": 42,
		"longitude": -88,
		"postalCode": null,
		"city": null,
		"countryCode": null,
		"countryName": null,
		"defaultLocation": false,
		"properties": {
			"city": null,
			"countryCode": null,
			"latitude": "42.0",
			"postalCode": null,
			"countryName": null,
			"longitude": "-88.0"
		}
	},
	"dispersion": {
		"limit": 1,
		"shuffled": true
	},
	"ip6RoutingEnabled": true,
	"responseHeaders": {},
	"requestHeaders": [],
	"regionalGeoEnabled": false,
	"geolocationProvider": "maxmindGeolocationService",
	"anonymousIpEnabled": false,
	"sslEnabled": true,
	"acceptHttp": true,
	"deepCache": "NEVER",
	"dns": false,
	"locationLimit": 0,
	"maxDnsIps": 0,
	"sslReady": true,
	"available": true
	}

.. _tr-api-crs-coveragezone-caches:

``/crs/coveragezone/caches``
============================
A list of caches for a given :term:`Delivery Service` and :term:`Cache Group`.

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------------------+----------+--------------------------------------------------------------------------------------------------------------------------------+
	| Name              | Required | Description                                                                                                                    |
	+===================+==========+================================================================================================================================+
	| deliveryServiceId | yes      | The integral, unique identifier?/'xml_id'?/name? of a :term:`Delivery Service` served by this Traffic Router                   |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------------------------+
	| cacheLocationId   | yes      | The :ref:`Name of a Cache Group <cache-group-name>` to which this Traffic Router is capable of routing client traffic          |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------------------------+

Response Structure
------------------
TBD

``/crs/coveragezone/cachelocation``
===================================
The resulting :term:`Cache Group` for a given client IP and :term:`Delivery Service`.

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| Name              | Required | Description                                                                                                  |
	+===================+==========+==============================================================================================================+
	| ip                | yes      | The IP address of a potential client                                                                         |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| deliveryServiceId | yes      | The integral, unique identifier?/'xml_id'?/name? of a :term:`Delivery Service` served by this Traffic Router |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+

Response Structure
------------------
TBD

.. _tr-api-crs-deepcoveragezone-cachelocation:

``/crs/deepcoveragezone/cachelocation``
=======================================
The resulting :term:`Cache Group` using the :term:`Deep Coverage Zone File` (deep caching) for a given client IP and :term:`Delivery Service`.

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| Name              | Required | Description                                                                                                  |
	+===================+==========+==============================================================================================================+
	| ip                | yes      | The IP address of a potential client                                                                         |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| deliveryServiceId | yes      | The integral, unique identifier?/'xml_id'?/name? of a :term:`Delivery Service` served by this Traffic Router |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+

Response Structure
------------------
TBD

.. _tr-api-crs-consistenthash-patternbased-regex:

``/crs/consistenthash/patternbased/regex``
==========================================
The resulting path that will be used for consistent hashing when the given regex is applied to the given request path.

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| Name              | Required | Description                                                                                                  |
	+===================+==========+==============================================================================================================+
	| regex             | yes      | The (URI encoded) regular expression to be used to test pattern based consistent hashing                     |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| requestPath       | yes      | The (URI encoded) request path to use to test pattern based consistent hashing                               |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /crs/consistenthash/patternbased/regex?regex=%2F.*%3F%28%2F.*%3F%2F%29.*%3F%28%5C.m3u8%29&requestPath=%2Ftext1234%2Fname%2Fasset.m3u8 HTTP/1.1
	Host: localhost:3333
	User-Agent: curl/7.54.0
	Accept: */*

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Type: application/json;charset=UTF-8
	Transfer-Encoding: chunked
	Date: Fri, 15 Feb 2019 22:06:53 GMT

	{
	"resultingPathToConsistentHash":"/name/.m3u8",
	"consistentHashRegex":"/.*?(/.*?/).*?(\\.m3u8)",
	"requestPath":"/text1234/name/asset.m3u8"
	}

.. _tr-api-crs-consistenthash-patternbased-deliveryservice:

``/crs/consistenthash/patternbased/deliveryservice``
====================================================
The resulting path that will be used for consistent hashing for the given delivery service and the given request path.

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| Name              | Required | Description                                                                                                  |
	+===================+==========+==============================================================================================================+
	| requestPath       | yes      | The (URI encoded) request path to use to test pattern based consistent hashing                               |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| deliveryServiceId | yes      | The integral, unique identifier?/'xml_id'?/name? of a :term:`Delivery Service` served by this Traffic Router |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /crs/consistenthash/patternbased/deliveryservice?deliveryServiceId=asdf&requestPath=%2Fsometext1234%2Fstream_name%2Fasset_name.m3u8 HTTP/1.1
	Host: localhost:3333
	User-Agent: curl/7.54.0
	Accept: */*

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Type: application/json;charset=UTF-8
	Transfer-Encoding: chunked
	Date: Fri, 15 Feb 2019 22:12:38 GMT

	{
	"resultingPathToConsistentHash":"/sometext1234/stream_name/asset_name.m3u8",
	"deliveryServiceId":"asdf",
	"requestPath":"/sometext1234/stream_name/asset_name.m3u8"
	}

.. _tr-api-crs-consistenthash-cache-coveragezone-steering:

``/crs/consistenthash/cache/coveragezone/steering``
===================================================
The resulting cache of the consistent hash using coverage zone for a given client IP, delivery service and, request path -- used to test cache selection for steering delivery services.

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| Name              | Required | Description                                                                                                  |
	+===================+==========+==============================================================================================================+
	| requestPath       | yes      | The (URI encoded) request path to use to test pattern based consistent hashing                               |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| deliveryServiceId | yes      | The integral, unique identifier?/'xml_id'?/name? of a :term:`Delivery Service` served by this Traffic Router |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| ip                | yes      | The IP address of a potential client                                                                         |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------+

Response Structure
------------------
TBD
