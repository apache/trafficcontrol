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

.. _to-api-oc-ci-configuration_requests:

********************************
``OC/CI/configuration/requests``
********************************

``GET``
=======
Returns the requested updates for :abbr:`CDNi (Content Delivery Network Interconnect)` configurations. An optional ``id`` parameter will return only information for a specific request.

:Auth. Required: Yes
:Roles Required: "admin"
:Permissions Required: CDNI-ADMIN:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                                                   |
	+===========+==========+===============================================================================================================+
	| id        | no       | Return only the configuration requests identified by this integral, unique identifier                         |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+

Response Structure
------------------
:id:                     An integral, unique identifier for the requested configuration updates.
:ucdn:                   The name of the :abbr:`uCDN (Upstream Content Delivery Network)` to which the requested changes apply.
:data:                   An array of generic :abbr:`FCI (Footprint and Capabilities Advertisement Interface)` base objects.
:host:                   The domain to which the requested changes apply.
:requestType:            A string of the type of configuration update request.
:asyncStatusId:          An integral, unique identifier for the associated asynchronous status.
:generic-metadata-type:  A string of the type of metadata to follow conforming to :rfc:`8006`.
:generic-metadata-value: An array of generic metadata value objects conforming to :rfc:`8006` and :abbr:`SVA (Streaming Video Alliance)` specifications.
:footprints:             An array of footprints impacted by this generic base object.

.. note:: These are meant to be generic and therefore there is not much information in these documents. For further information please see :rfc:`8006`, :rfc:`8007`, :rfc:`8008`, and the :abbr:`SVA (Streaming Video Alliance)` documents titled `Footprint and Capabilities Interface: Open Caching API`, `Open Caching API Implementation Guidelines`, `Configuration Interface: Part 1 Specification - Overview & Architecture`, `Configuration Interface: Part 2 Specification – CDNi Metadata Model Extensions`, and `Configuration Interface: Part 3 Specification – Publishing Layer APIs`.

.. code-block:: json
	:caption: Example /OC/CI/configuration/requests Response

	{
		"response": [
			{
				"id": 1,
				"ucdn": "ucdn1",
				"data": [
					{
						"generic-metadata-type": "MI.RequestedCapacityLimits",
						"generic-metadata-value": {
							"requested-limits": [
								{
									"limit-type": "egress",
									"limit-value": 232323,
									"footprints": [
										{
											"footprint-type": "ipv4cidr",
											"footprint-value": [
												"127.0.0.1",
												"127.0.0.2"
											]
										},
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
					}
				],
				"host": "example.com",
				"requestType": "hostConfigUpdate",
				"asyncStatusId": 0
			}
		]
	}
