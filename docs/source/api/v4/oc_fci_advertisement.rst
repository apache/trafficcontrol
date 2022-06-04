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

.. _to-api-oc-fci-advertisement:

************************
``OC/FCI/advertisement``
************************

``GET``
=======
Returns the complete footprint and capabilities information structure the :abbr:`dCDN (Downstream Content Delivery Network)` wants to expose to a given :abbr:`uCDN (Upstream Content Delivery Network)`.

.. note:: Users with the ``ICDN:UCDN-OVERRIDE`` permission will need to provide a "ucdn" query parameter to bypass the need for :abbr:`uCDN (Upstream Content Delivery Network)` information in the :abbr:`JWT (JSON Web Token)` and allow them to view all :abbr:`CDNi (Content Delivery Network Interconnect)` information.

:Auth. Required: No
:Roles Required: "admin" or "operations"
:Permissions Required: CDNI:READ
:Response Type:  Array

Request Structure
-----------------
This requires authorization using a :abbr:`JWT (JSON Web Token)` provided by the :abbr:`dCDN (Downstream Content Delivery Network)` to identify the :abbr:`uCDN (Upstream Content Delivery Network)`. This token must include the following claims:

.. table:: Required JWT claims

	+-----------------+--------------------------------------------------------------------------------------------------------------------+
	|    Name         | Description                                                                                                        |
	+=================+====================================================================================================================+
	|      iss        | Issuer claim as a string key for the :abbr:`uCDN (Upstream Content Delivery Network)`                              |
	+-----------------+--------------------------------------------------------------------------------------------------------------------+
	|      aud        | Audience claim as a string key for the :abbr:`dCDN (Downstream Content Delivery Network)`                          |
	+-----------------+--------------------------------------------------------------------------------------------------------------------+
	|      exp        | Expiration claim as the expiration date as a Unix epoch timestamp (in seconds)                                     |
	+-----------------+--------------------------------------------------------------------------------------------------------------------+

Response Structure
------------------
:capabilities:     An array of generic :abbr:`FCI (Footprint and Capabilities Advertisement Interface)` base objects.
:capability-type:  A string of the type of base object.
:capability-value: An array of the value for the base object.
:footprints:       An array of footprints impacted by this generic base object.

	.. note:: These are meant to be generic and therefore there is not much information in these documents. For further information please see :rfc:`8006`, :rfc:`8007`, :rfc:`8008`, and the :abbr:`SVA (Streaming Video Alliance)` documents titled `Footprint and Capabilities Interface: Open Caching API`, `Open Caching API Implementation Guidelines`, `Configuration Interface: Part 1 Specification - Overview & Architecture`, `Configuration Interface: Part 2 Specification – CDNi Metadata Model Extensions`, and `Configuration Interface: Part 3 Specification – Publishing Layer APIs`.

.. code-block:: json
	:caption: Example /OC/FCI/advertisement Response

	{
		"capabilities": [
			{
				"capability-type": "FCI.CapacityLimits",
				"capability-value": [
					{
						"limits": [
							{
								"id": "host_limit_requests_requests",
								"scope": {
									"type": "testScope",
									"value": [
										"test.com"
									]
								},
								"limit-type": "requests",
								"maximum-hard": 20,
								"maximum-soft": 15,
								"telemetry-source": {
									"id": "request_metrics",
									"metric": "requests"
								}
							},
							{
								"id": "total_limit_egress_capacity",
								"limit-type": "egress",
								"maximum-hard": 202020,
								"maximum-soft": 500,
								"telemetry-source": {
									"id": "capacity_metrics",
									"metric": "capacity"
								}
							}
						]
					}
				],
				"footprints": [
					{
						"footprint-type": "countrycode",
						"footprint-value": [
							"us"
						]
					}
				]
			},
			{
				"capability-type": "FCI.Telemetry",
				"capability-value": {
					"sources": [
						{
							"id": "capacity_metrics",
							"type": "generic",
							"metrics": [
								{
									"name": "capacity",
									"time-granularity": 0,
									"data-percentile": 50,
									"latency": 0
								}
							],
							"configuration": {
								"url": "example.com/telemetry1"
							}
						}
					]
				},
				"footprints": [
					{
						"footprint-type": "countrycode",
						"footprint-value": [
							"us"
						]
					}
				]
			}
		]
	}

