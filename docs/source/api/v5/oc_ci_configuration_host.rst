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

.. _to-api-oc-fci-configuration-host:

********************************
``OC/CI/configuration/{{host}}``
********************************

``PUT``
=======
Triggers an asynchronous task to update the configuration for the :abbr:`uCDN (Upstream Content Delivery Network)` and the specified host by adding the request to a queue to be reviewed later. This returns a 202 Accepted status and an endpoint to be used for status updates.

.. note:: Users with the ``ICDN:UCDN-OVERRIDE`` permission will need to provide a "ucdn" query parameter to bypass the need for :abbr:`uCDN (Upstream Content Delivery Network)` information in the :abbr:`JWT (JSON Web Token)` and allow them to view all :abbr:`CDNi (Content Delivery Network Interconnect)` information.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: CDNI:UPDATE
:Response Type:  Object

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

.. table:: Request Path Parameters

	+-------+-----------------------------------------------------------------------------------+
	| Name  |                 Description                                                       |
	+=======+===================================================================================+
	|  host | The text identifier for the host domain to be updated with the new configuration. |
	+-------+-----------------------------------------------------------------------------------+

:type: A string of the type of metadata to follow. See :rfc:`8006` for possible values. Only a selection of these are supported.
:host-metadata: An array of generic metadata objects that conform to :rfc:`8006`.
:generic-metadata-type: A string of the type of metadata to follow conforming to :rfc:`8006`.
:generic-metadata-value: An array of generic metadata value objects conforming to :rfc:`8006` and :abbr:`SVA (Streaming Video Alliance)` specifications.

.. code-block:: http
	:caption: Example /OC/CI/configuration Request

	PUT /api/5.0/oc/ci/configuration/example.com HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 181
	Content-Type: application/json

	{
		"type": "MI.HostMetadata",
		"host-metadata": [
			{
				"generic-metadata-type": "MI.RequestedCapacityLimits",
				"generic-metadata-value": {
					"requested-limits": [
						{
							"limit-type": "egress",
							"limit-value": 20000,
							"footprints": [
								{
									"footprint-type": "ipv4cidr",
									"footprint-value": [
										"127.0.0.1",
										"127.0.0.2"
									]
								}
							]
						}
					]
				}
			}
		]
	}

Response Structure
------------------

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 202 Accepted
	Content-Type: application/json

	{ "alerts": [
		{
			"text": "CDNi configuration update request received. Status updates can be found here: /api/5.0/async_status/1",
			"level": "success"
		}
	]}
