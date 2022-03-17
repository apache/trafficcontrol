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

.. _tm-api:

********************
Traffic Monitor APIs
********************
The Traffic Monitor URLs below allow certain query parameters for use in controlling the data returned.

.. note:: Unlike :ref:`Traffic Ops API endpoints <to-api>`\ , no authentication is required for any of these, and as such there can be no special role requirements for a user.

.. _tm-publish-EventLog:

``/publish/EventLog``
=====================
Gets a log of recent changes in the availability of polled caches.

``GET``
-------
:Response Type: Array (key 'events' contains an array of all data)

Response Structure
""""""""""""""""""
:event: an entry in the top-level ``events`` array

	:description: A string containing short description of the event
	:hostname:    A string containing the server's full hostname
	:index:       A serial integer that is incremented for each sequential  event
	:isAvailable: A boolean value indicating whether the server is available following this event
	:name:        The server's short hostname as a string
	:time:        A UNIX timestamp as an integer
	:type:        The type of the server as a string

.. code-block:: json
	:caption: Example Response

	{ "events": [
		{
			"time": 1538417713,
			"index": 67848,
			"description": "REPORTED - loadavg too high (36.37 \u003e 25.00) (health)",
			"name": "edge",
			"hostname": "edge",
			"type":"EDGE",
			"isAvailable":false
		}
	]}

``/publish/CacheStats``
=======================
Statistics gathered for each cache.

``GET``
-------
:Response Type: Object

Request Structure
"""""""""""""""""
.. table:: Request Query Parameters

	+--------------+---------+------------------------------------------------+
	|  Parameter   | Type    |                  Description                   |
	+==============+=========+================================================+
	| ``hc``       | integer | The history count, number of items to display. |
	+--------------+---------+------------------------------------------------+
	| ``stats``    | string  | A comma separated list of stats to display.    |
	+--------------+---------+------------------------------------------------+
	| ``wildcard`` | boolean | Controls whether specified stats should be     |
	|              |         | treated as partial strings.                    |
	+--------------+---------+------------------------------------------------+

.. code-block:: http
	:caption: Example Request

	GET /publish/CacheStats HTTP/1.1
	Accept: */*
	Content-Type: application/json

Response Structure
""""""""""""""""""
:pp: Stores any provided request parameters provided as a string
:date: A ``ctime``-like string representation of the time at which the response was served
:caches: An object with keys that are the names of monitored :term:`cache servers`

	:<server name>: Each server's object is a collection of keys that are the names of statistics

		:<interface name>: The name of the network interface under the same sever

			:<statistic name>: The name of the statistic which this array represents. Each value in the array is one (and usually only one) object with the following structure:

				:value: The statistic's value. This is *always* a string, even if that string only contains a number.
				:time: An integer UNIX timestamp indicating the start time for this value of this statistic
				:span: The span of time - in milliseconds - for which this value is valid. This is determined by the polling interval for the statistic

.. code-block:: http
	:caption: Example Response

	HTTP/1.1 200 OK
	Content-Type: application/json
	Date: Thu, 14 May 2020 15:48:55 GMT
	Transfer-Encoding: chunked

	{
		"pp": "",
		"date": "Thu May 14 15:48:55 UTC 2020",
		"caches": {
			"mid": {
				"eth0": {
					"ats.proxy.process.ssl.cipher.user_agent.PSK-AES256-GCM-SHA384": [
						{
							"value": "0",
							"time": 1589471325624,
							"span": 99
						}
					]
				},
				"aggregate": {
					"ats.proxy.process.http.milestone.server_begin_write": [
						{
							"value": "174",
							"time": 1589471325624,
							"span": 1
						}
					]
				},
				"lo": {
					"ats.proxy.node.http.transaction_counts_avg_10s.miss_changed": [
						{
							"value": "0",
							"time": 1589471325624,
							"span": 99
						}
					]
				}
			},
			"edge": {
				"eth0": {
					"ats.proxy.process.ssl.cipher.user_agent.PSK-AES256-GCM-SHA384": [
						{
							"value": "0",
							"time": 1589471325624,
							"span": 99
						}
					]
				},
				"aggregate": {
					"ats.proxy.process.http.milestone.server_begin_write": [
						{
							"value": "174",
							"time": 1589471325624,
							"span": 1
						}
					]
				},
				"lo": {
					"ats.proxy.node.http.transaction_counts_avg_10s.miss_changed": [
						{
							"value": "0",
							"time": 1589471325624,
							"span": 99
						}
					]
				}
			}
		}
	}


``publish/CacheStats/{{cache}}``
================================
Statistics gathered for only a single cache.

``GET``
-------
:Response Type: Object

Request Structure
"""""""""""""""""
.. table:: Request Path Parameters

	+-----------+--------+----------------------------------+
	| Parameter | Type   |           Description            |
	+===========+========+==================================+
	| ``cache`` | string | The name of the cache to inspect |
	+-----------+--------+----------------------------------+

.. table:: Request Query Parameters

	+--------------+---------+------------------------------------------------+
	|  Parameter   | Type    |                  Description                   |
	+==============+=========+================================================+
	| ``hc``       | integer | The history count, number of items to display. |
	+--------------+---------+------------------------------------------------+
	| ``stats``    | string  | A comma separated list of stats to display.    |
	+--------------+---------+------------------------------------------------+
	| ``wildcard`` | boolean | Controls whether specified stats should be     |
	|              |         | treated as partial strings.                    |
	+--------------+---------+------------------------------------------------+

.. code-block:: http
	:caption: Example Request

	GET /api/CacheStats/mid HTTP/1.1
	Accept: */*
	Content-Type: application/json

Response Structure
""""""""""""""""""
:pp: Stores any provided request parameters provided as a string
:date: A ``ctime``-like string representation of the time at which the response was served
:caches: An object with keys that are the names of monitored :term:`cache servers` - only the cache named by the ``cache`` request path parameter will be shown

	:<server name>: The requested server's object is a collection of keys that are the names of statistics

		:<interface name>: The name of the network interface under the same sever

			:<statistic name>: The name of the statistic which this array represents. Each value in the array is one (and usually only one) object with the following structure:

				:value: The statistic's value. This is *always* a string, even if that string only contains a number.
				:time: An integer UNIX timestamp indicating the start time for this value of this statistic
				:span: The span of time - in milliseconds - for which this value is valid. This is determined by the polling interval for the statistic

.. code-block:: http
	:caption: Example Response

	HTTP/1.1 200 OK
	Content-Type: application/json
	Date: Thu, 14 May 2020 15:54:35 GMT
	Transfer-Encoding: chunked

	{
		"pp": "",
		"date": "Thu May 14 15:48:55 UTC 2020",
		"caches": {
			"mid": {
				"eth0": {
					"ats.proxy.process.ssl.cipher.user_agent.PSK-AES256-GCM-SHA384": [
						{
							"value": "0",
							"time": 1589471325624,
							"span": 99
						}
					]
				},
				"aggregate": {
					"ats.proxy.process.http.milestone.server_begin_write": [
						{
							"value": "174",
							"time": 1589471325624,
							"span": 1
						}
					]
				},
				"lo": {
					"ats.proxy.node.http.transaction_counts_avg_10s.miss_changed": [
						{
							"value": "0",
							"time": 1589471325624,
							"span": 99
						}
					]
				}
			}
		}
	}

``/publish/DsStats``
====================
Statistics gathered for :term:`Delivery Services`

``GET``
-------
:Response Type: Object

Request Structure
"""""""""""""""""
.. table:: Request Query Parameters

	+--------------+---------+------------------------------------------------+
	|  Parameter   | Type    |                  Description                   |
	+==============+=========+================================================+
	| ``hc``       | int     | The history count, number of items to display. |
	+--------------+---------+------------------------------------------------+
	| ``stats``    | string  | A comma separated list of stats to display.    |
	+--------------+---------+------------------------------------------------+
	| ``wildcard`` | boolean | Controls whether specified stats should be     |
	|              |         | treated as partial strings.                    |
	+--------------+---------+------------------------------------------------+

Response Structure
""""""""""""""""""

TODO

``/publish/DsStats/{{deliveryService}}``
========================================
Statistics gathered for this :term:`Delivery Service` only.

``GET``
-------
:Response Type: ?

Request Structure
"""""""""""""""""
.. table:: Request Path Parameters

	+---------------------+--------+-----------------------------------------------------+
	| Parameter           | Type   | Description                                         |
	+=====================+========+=====================================================+
	| ``deliveryService`` | string | The name of the :term:`Delivery Service` to inspect |
	+---------------------+--------+-----------------------------------------------------+


.. table:: Request Query Parameters

	+--------------+---------+------------------------------------------------+
	|  Parameter   | Type    |                  Description                   |
	+==============+=========+================================================+
	| ``hc``       | integer | The history count, number of items to display. |
	+--------------+---------+------------------------------------------------+
	| ``stats``    | string  | A comma separated list of stats to display.    |
	+--------------+---------+------------------------------------------------+
	| ``wildcard`` | boolean | Controls whether specified stats should be     |
	|              |         | treated as partial strings.                    |
	+--------------+---------+------------------------------------------------+

Response Structure
""""""""""""""""""

TODO

``/publish/CrStates``
=====================
The current state of this CDN per the :ref:`health-proto`.

``GET``
-------
:Response Type: Object

.. code-block:: http
	:caption: Example Request

	GET /publish/CrStates HTTP/1.1
	Accept: */*

Response Structure
""""""""""""""""""
:caches: An object with keys that are the names of monitored :term:`cache servers`.

	:isAvailable: Whether or not this :term:`cache server` is available for routing overall
	:ipv4Available: Whether or not an IPv4 interface on this :term:`cache server` is available for routing.
	:ipv6Available: Whether or not an IPv6 interface on this :term:`cache server` is available for routing.
	:status: The status of this server, along with any additional reason for it to be marked as such
	:lastPoll: The last time the health data for this server was polled by a traffic monitor

:deliveryServices: An object with keys that are the :ref:`XMLIDs <ds-xmlid>` of monitored :term:`Delivery Services`.

	:disabledLocations: An array of the names of disabled "locations" (i.e. :term:`Cache Groups`) for this :term:`Delivery Service`.
	:isAvailable: Whether or not this :term:`Delivery Service` is available for routing

.. code-block:: http
	:caption: Example Response

	HTTP/1.1 200 OK
	Content-Type: application/json
	Date: Thu, 14 May 2020 15:54:35 GMT
	Transfer-Encoding: chunked

	{
		"caches": {
			"edge": {
				"isAvailable": true,
				"ipv4Available": true,
				"ipv6Available": false,
				"status": "REPORTED - available",
				"lastPoll": "2022-03-15T17:54:03.821178179Z"
			}
		},
		"deliveryServices": {
			"dev-ds": {
				"disabledLocations": [],
				"isAvailable": true
			}
		}
	}


``/publish/CrConfig``
=====================
The CDN :term:`Snapshot` (historically named a "CRConfig") served to and consumed by Traffic Router.

``GET``
-------
:Response Type: ?

Response Structure
""""""""""""""""""

TODO

``/publish/PeerStates``
=======================
The health state information from all peer Traffic Monitors.

``GET``
-------
:Response Type: ?

Request Structure
"""""""""""""""""
.. table:: Request Query Parameters

	+--------------+---------+------------------------------------------------+
	|  Parameter   | Type    |                  Description                   |
	+==============+=========+================================================+
	| ``hc``       | integer | The history count, number of items to display. |
	+--------------+---------+------------------------------------------------+
	| ``stats``    | string  | A comma separated list of stats to display.    |
	+--------------+---------+------------------------------------------------+
	| ``wildcard`` | boolean | Controls whether specified stats should be     |
	|              |         | treated as partial strings.                    |
	+--------------+---------+------------------------------------------------+

Response Structure
""""""""""""""""""

TODO


``/publish/DistributedPeerStates``
==================================
The health state information from all distributed peer Traffic Monitors.

``GET``
-------
:Response Type: ?

Request Structure
"""""""""""""""""
.. table:: Request Query Parameters

	+--------------+---------+------------------------------------------------+
	|  Parameter   | Type    |                  Description                   |
	+==============+=========+================================================+
	| ``hc``       | integer | The history count, number of items to display. |
	+--------------+---------+------------------------------------------------+
	| ``stats``    | string  | A comma separated list of stats to display.    |
	+--------------+---------+------------------------------------------------+
	| ``wildcard`` | boolean | Controls whether specified stats should be     |
	|              |         | treated as partial strings.                    |
	+--------------+---------+------------------------------------------------+

Response Structure
""""""""""""""""""

TODO


``/publish/Stats``
==================
The general statistics about Traffic Monitor.

``GET``
-------
:Response Type: ?

Response Structure
""""""""""""""""""

TODO

``/publish/StatSummary``
========================
The summary of :term:`cache server` statistics.

``GET``
-------
:Response Type: ?

Request Structure
"""""""""""""""""
.. table:: Request Query Parameters

	+---------------+---------+-----------------------------------------------------------+
	|   Parameter   |   Type  |                        Description                        |
	+===============+=========+===========================================================+
	| ``startTime`` | number  | Window start. The number of milliseconds since the epoch. |
	+---------------+---------+-----------------------------------------------------------+
	| ``endTime``   | number  | Window end. The number of milliseconds since the epoch.   |
	+---------------+---------+-----------------------------------------------------------+
	| ``hc``        | integer | The history count, number of items to display.            |
	+---------------+---------+-----------------------------------------------------------+
	| ``stats``     | string  | A comma separated list of stats to display.               |
	+---------------+---------+-----------------------------------------------------------+
	| ``wildcard``  | boolean | Controls whether specified stats should be                |
	|               |         | treated as partial strings.                               |
	+---------------+---------+-----------------------------------------------------------+
	| ``cache``     | string  | Summary statistics for just this cache.                   |
	+---------------+---------+-----------------------------------------------------------+

Response Structure
""""""""""""""""""

TODO

``/publish/ConfigDoc``
======================
The overview of configuration options.

``GET``
-------
:Response Type: ?

Response Structure
""""""""""""""""""

TODO
