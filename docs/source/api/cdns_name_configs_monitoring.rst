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

.. _to-api-cdns-name-configs-monitoring:

************************************
``cdns/{{name}}/configs/monitoring``
************************************

.. seealso:: :ref:`health-proto`

``GET``
=======
Retrieves information concerning the monitoring configuration for a specific CDN.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------------------+
	| Name | Description                                                            |
	+======+========================================================================+
	| name | The name of the CDN for which monitoring configuration will be fetched |
	+------+------------------------------------------------------------------------+

Response Structure
------------------
:cacheGroups: An array of objects representing each of the Cache Groups being monitored within this CDN

	:coordinates: An object representing the physical location of this Cache Group

		:latitude:  The geographic latitude of this Cache Group
		:longitude: The geographic longitude of this Cache Group

	:name: The name of this Cache Group

:config: A collection of parameters used to configure the monitoring behaviour of Traffic Monitor

	:hack.ttl:                    Unknown
	:health.event-count:          The total number of health events to store
	:health.polling.interval:     An interval in milliseconds on which to poll for health statistics
	:health.threadPool:           The number of threads to be used for health polling
	:health.timepad:              A 'padding time' to add to requests to spread them out for Traffic Control systems that use a large number of Traffic Monitors
	:tm.crConfig.polling.url:     The URL from which a :term:`Snapshot` can be obtained
	:tm.dataServer.polling.url:   The URL from which a list of data servers can be obtained
	:tm.healthParams.polling.url: The URL from which a list of health-polling parameters can be obtained
	:tm.polling.interval:         The interval at which to poll for configuration updates

:deliveryServices: An array of objects representing each :term:`Delivery Service` provided by this CDN

	:status:             The :term:`Delivery Service`'s status
	:totalKbpsThreshold: A threshold rate of data transfer this :term:`Delivery Service` is configured to handle, in Kilobits per second
	:totalTpsThreshold:  A threshold amount of transactions per second that this :term:`Delivery Service` is configured to handle
	:xmlId:              An integral, unique identifier for this Deliver Service (named "xmlId" for legacy reasons)

:profiles: An array of the profiles in use by the :term:`cache server` s and :term:`Delivery Service`\ s belonging to this CDN

	:name:       The profile's name
	:parameters: An array of the parameters in this profile that relate to monitoring configuration. This can be ``null`` if the servers using this profile cannot be monitored (e.g. Traffic Routers)

		:health.connection.timeout:                 A timeout value, in milliseconds, to wait before giving up on a health check request
		:health.polling.url:                        A URL to request for polling health. Substitutions can be made in a shell-like syntax using the properties of an object from the ``"trafficServers"`` array
		:health.threshold.availableBandwidthInKbps: The total amount of bandwidth that servers using this profile are allowed, in Kilobits per second. This is a string and using comparison operators to specify ranges, e.g. ">10" means "more than 10 kbps"
		:health.threshold.loadavg:                  The UNIX loadavg at which the server should be marked "unhealthy" - see ``man uptime``
		:health.threshold.queryTime:                The highest allowed length of time for completing health queries (after connection has been established) in milliseconds
		:history.count:                             The number of past events to store; once this number is reached, the oldest event will be forgotten before a new one can be added

	:type: The type of the profile

:trafficMonitors: An array of objects representing each Traffic Monitor that monitors this CDN (this is used by Traffic Monitor's "peer polling" function)

	:fqdn:     AN FQDN that resolves to the IP (and/or IPv6) address of the server running this Traffic Monitor instance
	:hostname: The hostname of the server running this Traffic Monitor instance
	:ip6:      The IPv6 address of this Traffic Monitor - when applicable
	:ip:       The IP address of this Traffic Monitor
	:port:     The port on which this Traffic Monitor listens for incoming connections
	:profile:  The name of the profile assigned to this Traffic Monitor
	:status:   The status of the server running this Traffic Monitor instance

:trafficServers: An array of objects that represent the caches being monitored within this CDN

	:cacheGroup:    The Cache Group to which this cache belongs
	:fqdn:          A Fully Qualified Domain Name (FQDN) that resolves to the :term:`cache server`'s IP (or IPv6) address
	:hashId:        The short name for the :term:`cache server` - named "hashId" for legacy reasons
	:hostName:      The (short) hostname of the :term:`cache server`
	:interfacename: The name of the network interface device being used by the cache's HTTP proxy
	:ip6:           The cache's IPv6 address - when applicable
	:ip:            The cache's IP address
	:port:          The port on which the cache listens for incoming connections
	:profile:       The name of the profile assigned to this cache
	:status:        The status of the Cache
	:type:          The type of the cache - should be either ``EDGE`` or ``MID``

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: uLR+tRoqR8SYO38j3DV9wQ+IkJ7Kf+MCoFkcWZtsgbpLJ+0S6f+IiI8laNVeDgrM/P23MAQ6BSepm+EJRl1AXQ==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 14 Nov 2018 21:09:31 GMT
	Transfer-Encoding: chunked

	{ "response": {
		"trafficServers": [
			{
				"profile": "ATS_EDGE_TIER_CACHE",
				"status": "REPORTED",
				"ip": "172.16.239.100",
				"ip6": "fc01:9400:1000:8::100",
				"port": 80,
				"cachegroup": "CDN_in_a_Box_Edge",
				"hostname": "edge",
				"fqdn": "edge.infra.ciab.test",
				"interfacename": "eth0",
				"type": "EDGE",
				"hashid": "edge"
			},
			{
				"profile": "ATS_MID_TIER_CACHE",
				"status": "REPORTED",
				"ip": "172.16.239.120",
				"ip6": "fc01:9400:1000:8::120",
				"port": 80,
				"cachegroup": "CDN_in_a_Box_Mid",
				"hostname": "mid",
				"fqdn": "mid.infra.ciab.test",
				"interfacename": "eth0",
				"type": "MID",
				"hashid": "mid"
			}
		],
		"trafficMonitors": [
			{
				"profile": "RASCAL-Traffic_Monitor",
				"status": "ONLINE",
				"ip": "172.16.239.40",
				"ip6": "fc01:9400:1000:8::40",
				"port": 80,
				"cachegroup": "CDN_in_a_Box_Edge",
				"hostname": "trafficmonitor",
				"fqdn": "trafficmonitor.infra.ciab.test"
			}
		],
		"cacheGroups": [
			{
				"name": "CDN_in_a_Box_Mid",
				"coordinates": {
					"latitude": 38.897663,
					"longitude": -77.036574
				}
			},
			{
				"name": "CDN_in_a_Box_Edge",
				"coordinates": {
					"latitude": 38.897663,
					"longitude": -77.036574
				}
			}
		],
		"profiles": [
			{
				"name": "CCR_CIAB",
				"type": "CCR",
				"parameters": null
			},
			{
				"name": "ATS_EDGE_TIER_CACHE",
				"type": "EDGE",
				"parameters": {
					"health.connection.timeout": 2000,
					"health.polling.url": "http://${hostname}/_astats?application=&inf.name=${interface_name}",
					"health.threshold.availableBandwidthInKbps": ">1750000",
					"health.threshold.loadavg": "25.0",
					"health.threshold.queryTime": 1000,
					"history.count": 30
				}
			},
			{
				"name": "ATS_MID_TIER_CACHE",
				"type": "MID",
				"parameters": {
					"health.connection.timeout": 2000,
					"health.polling.url": "http://${hostname}/_astats?application=&inf.name=${interface_name}",
					"health.threshold.availableBandwidthInKbps": ">1750000",
					"health.threshold.loadavg": "25.0",
					"health.threshold.queryTime": 1000,
					"history.count": 30
				}
			}
		],
		"deliveryServices": [],
		"config": {
			"hack.ttl": 30,
			"health.event-count": 200,
			"health.polling.interval": 6000,
			"health.threadPool": 4,
			"health.timepad": 0,
			"heartbeat.polling.interval": 3000,
			"location": "/opt/traffic_monitor/conf",
			"peers.polling.interval": 3000,
			"tm.crConfig.polling.url": "https://${tmHostname}/CRConfig-Snapshots/${cdnName}/CRConfig.xml",
			"tm.dataServer.polling.url": "https://${tmHostname}/dataserver/orderby/id",
			"tm.healthParams.polling.url": "https://${tmHostname}/health/${cdnName}",
			"tm.polling.interval": 2000
		}
	}}
