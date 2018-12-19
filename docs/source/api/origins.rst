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

.. _to-api-origins:

***********
``origins``
***********
.. versionadded:: 1.3

``GET``
=======
Gets all requested origins.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| Name            | Required | Description                                                                                                                                                                          |
	+=================+==========+======================================================================================================================================================================================+
	| cachegroup      | no       | Return only origins within the Cache Group identified by this integral, unique identifier                                                                                            |
	+-----------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| coordinate      | no       | Return only origins located at the geographic coordinates identified by this integral, unique identifier                                                                             |
	+-----------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| deliveryservice | no       | Return only origins that belong to the Delivery Service identified by this integral, unique identifier                                                                               |
	+-----------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| id              | no       | Return only the origin that has this integral, unique identifier                                                                                                                     |
	+-----------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| name            | no       | Return only origins by this name                                                                                                                                                     |
	+-----------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| profileId       | no       | Return only origins which use the profile identified by this integral, unique identifier                                                                                             |
	+-----------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| primary         | no       | If ``true``, return only origins which are the the primary origin of the Delivery Service to which they belong - if ``false`` return only origins which are *not* the primary origin |
	|                 |          | of the Delivery Service to which they belong                                                                                                                                         |
	+-----------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| tenant          | no       | Return only origins belonging to the tenant identified by this integral, unique identifier                                                                                           |
	+-----------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+

.. note:: Several fields of origin definitions which are filterable by Query Parameters are allowed to be ``null``. ``null`` values in these fields will be filtered *out* appropriately by such Query Parameters, but do note that ``null`` is not a valid value accepted by any of these Query Parameters, and attempting to pass it will result in an error.

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/origins?name=demo1 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:cachegroup:        The name of the Cache Group to which the origin belongs
:cachegroupId:      An integral, unique identifier for the Cache Group to which the origin belongs
:coordinate:        The name of a coordinate pair that defines the origin's geographic location
:coordinateID:      An integral, unique identifier for the coordinate pair that defines the origin's geographic location
:deliveryService:   The 'xml_id' of the Delivery Service to which the origin belongs
:deliveryServiceId: An integral, unique identifier for the Delivery Service to which the origin belongs
:fqdn:              The FQDN of the origin
:id:                An integral, unique identifier for this origin
:ip6Address:        The IPv6 address of the Origin
:ipAddress:         The IPv4 address of the Origin
:isPrimary:         A boolean value which, when ``true`` specifies this origin as the 'primary' origin served by ``deliveryService``
:lastUpdated:       The date and time at which this origin was last modified
:name:              The name of the origin
:port:              The TCP port on which the origin listens
:profile:           The name of the profile used by this origin
:profileId:         An integral, unique identifier for the profile used by this origin
:protocol:          The protocol used by this origin - will be one of 'http' or 'https'
:tenant:            The name of the tenant that owns this origin
:tenantId:          An integral, unique identifier for the tenant that owns this origin

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: sm8DpvdvrfdSVLtmXTdfjsZbTlbc+pI40Gy0aj00XIURTPfFXuv/4LgHb6A3r92iymbRHvFrH6qdB2g97U2sBg==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 11 Dec 2018 15:43:41 GMT
	Content-Length: 376

	{ "response": [
		{
			"cachegroup": null,
			"cachegroupId": null,
			"coordinate": null,
			"coordinateId": null,
			"deliveryService": "demo1",
			"deliveryServiceId": 1,
			"fqdn": "origin.infra.ciab.test",
			"id": 1,
			"ip6Address": null,
			"ipAddress": null,
			"isPrimary": true,
			"lastUpdated": "2018-12-10 19:11:32+00",
			"name": "demo1",
			"port": null,
			"profile": null,
			"profileId": null,
			"protocol": "http",
			"tenant": "root",
			"tenantId": 1
		}
	]}

``POST``
========
Creates a new origin definition.

.. warning:: At the time of this writing it is possible to create and/or modify origin definitions assigned to STEERING and CLIENT_STEERING Delivery Services - despite that an origin has no meaning in those contexts. In these cases, the API responses may give incorrect output - see `GitHub Issue #3107 <https://github.com/apache/trafficcontrol/issues/3107>`_ for details and updates.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
:cachegroupId:      An optional, integral, unique identifier that identifies a Cache Group to which the new origin shall belong
:coordinateID:      An optional, integral, unique identifier of a coordinate pair that shall define the origin's geographic location
:deliveryServiceId: The integral, unique identifier of the Delivery Service to which the new origin shall belong
:fqdn:              The Fully Qualified Domain Name (FQDN) of the origin
:ip6Address:        An optional string containing the IPv6 address of the origin
:ipAddress:         An optional string containing the IPv4 address of the origin
:isPrimary:         An optional boolean which, if ``true`` will set this origin as the 'primary' origin served by the Delivery Service identified by ``deliveryServiceID``

	.. note:: Though not specifying this field in this request will leave it as ``null`` in the output, Traffic Ops will silently coerce that to its default value: ``false``.

:name:      A human-friendly name of the Origin
:port:      An optional port number on which the origin listens for incoming TCP connections
:profileId: An optional, integral, unique identifier for a profile that the new origin shall use
:protocol:  The protocol used by the origin - must be one of 'http' or 'https'
:tenantId:  An optional\ [1]_, integral, unique identifier for the tenant which shall own the new origin

.. code-block:: http
	:caption: Request Example

	POST /api/1.4/origins HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 114
	Content-Type: application/json

	{
		"deliveryServiceId": 2,
		"fqdn": "example.com",
		"name": "example",
		"port": 80,
		"protocol": "http",
		"tenantId": 1
	}

.. [1] The ``tenantId`` field is required if and only if tenancy is enabled within Traffic Ops.

Response Structure
------------------
:cachegroup:        The name of the Cache Group to which the origin belongs
:cachegroupId:      An integral, unique identifier for the Cache Group to which the origin belongs
:coordinate:        The name of a coordinate pair that defines the origin's geographic location
:coordinateID:      An integral, unique identifier for the coordinate pair that defines the origin's geographic location
:deliveryService:   The 'xml_id' of the Delivery Service to which the origin belongs
:deliveryServiceId: An integral, unique identifier for the Delivery Service to which the origin belongs
:fqdn:              The FQDN of the origin
:id:                An integral, unique identifier for this origin
:ip6Address:        The IPv6 address of the Origin
:ipAddress:         The IPv4 address of the Origin
:isPrimary:         A boolean value which, when ``true`` specifies this origin as the 'primary' origin served by ``deliveryService``
:lastUpdated:       The date and time at which this origin was last modified
:name:              The name of the origin
:port:              The TCP port on which the origin listens
:profile:           The name of the profile used by this origin
:profileId:         An integral, unique identifier for the profile used by this origin
:protocol:          The protocol used by this origin - will be one of 'http' or 'https'
:tenant:            The name of the tenant that owns this origin
:tenantId:          An integral, unique identifier for the tenant that owns this origin

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: z4gp0MaqYu+gSRORhKT2eObVBuVDVx1rdteRaN5kRL9uJ3hNzUCi4dSKIt0rgNgOEDt6x/iTYrmVhr/TSHYtmA==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 11 Dec 2018 15:14:27 GMT
	Content-Length: 418

	{ "alerts": [
		{
			"text": "origin was created.",
			"level": "success"
		}
	],
	"response": {
		"cachegroup": null,
		"cachegroupId": null,
		"coordinate": null,
		"coordinateId": null,
		"deliveryService": null,
		"deliveryServiceId": 2,
		"fqdn": "example.com",
		"id": 2,
		"ip6Address": null,
		"ipAddress": null,
		"isPrimary": null,
		"lastUpdated": "2018-12-11 15:14:27+00",
		"name": "example",
		"port": 80,
		"profile": null,
		"profileId": null,
		"protocol": "http",
		"tenant": null,
		"tenantId": 1
	}}

``PUT``
=======
Updates an origin definition.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Query Parameters

	+------+----------+-----------------------------------------------------------------------+
	| Name | Required | Description                                                           |
	+======+==========+=======================================================================+
	|  id  | yes      | The integral, unique identifier of the origin definition being edited |
	+------+----------+-----------------------------------------------------------------------+

:cachegroupId:      An optional, integral, unique identifier that identifies a Cache Group to which the new origin shall belong
:coordinateID:      An optional, integral, unique identifier of a coordinate pair that shall define the origin's geographic location
:deliveryServiceId: The integral, unique identifier of the Delivery Service to which the new origin shall belong
:fqdn:              The Fully Qualified Domain Name (FQDN) of the origin
:ip6Address:        An optional string containing the IPv6 address of the origin
:ipAddress:         An optional string containing the IPv4 address of the origin
:isPrimary:         An optional boolean which, if ``true`` will set this origin as the 'primary' origin served by the Delivery Service identified by ``deliveryServiceID``
:name:      A human-friendly name of the Origin
:port:      An optional port number on which the origin listens for incoming TCP connections
:profileId: An optional, integral, unique identifier for a profile that the new origin shall use
:protocol:  The protocol used by the origin - must be one of 'http' or 'https'
:tenantId:  An optional\ [1]_, integral, unique identifier for the tenant which shall own the new origin

.. code-block:: http
	:caption: Request Example

	PUT /api/1.4/origins?id=2 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 135
	Content-Type: application/json

	{
		"deliveryServiceId": 2,
		"fqdn": "example.com",
		"isprimary": true,
		"name": "example",
		"port": 443,
		"protocol": "https",
		"tenantId": 1
	}


Response Structure
------------------
:cachegroup:        The name of the Cache Group to which the origin belongs
:cachegroupId:      An integral, unique identifier for the Cache Group to which the origin belongs
:coordinate:        The name of a coordinate pair that defines the origin's geographic location
:coordinateID:      An integral, unique identifier for the coordinate pair that defines the origin's geographic location
:deliveryService:   The 'xml_id' of the Delivery Service to which the origin belongs
:deliveryServiceId: An integral, unique identifier for the Delivery Service to which the origin belongs
:fqdn:              The FQDN of the origin
:id:                An integral, unique identifier for this origin
:ip6Address:        The IPv6 address of the Origin
:ipAddress:         The IPv4 address of the Origin
:isPrimary:         A boolean value which, when ``true`` specifies this origin as the 'primary' origin served by ``deliveryService``
:lastUpdated:       The date and time at which this origin was last modified
:name:              The name of the origin
:port:              The TCP port on which the origin listens
:profile:           The name of the profile used by this origin
:profileId:         An integral, unique identifier for the profile used by this origin
:protocol:          The protocol used by this origin - will be one of 'http' or 'https'
:tenant:            The name of the tenant that owns this origin
:tenantId:          An integral, unique identifier for the tenant that owns this origin

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: Zx7jOa7UAQxRtDenYodvGQSoooPj4m0yY0AIeUpbdelmYMiNdPYtW82BCmMesFXkmP74nV4HbTUyDHVMuJxZ7g==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 11 Dec 2018 15:40:53 GMT
	Content-Length: 420

	{ "alerts": [
		{
			"text": "origin was updated.",
			"level": "success"
		}
	],
	"response": {
		"cachegroup": null,
		"cachegroupId": null,
		"coordinate": null,
		"coordinateId": null,
		"deliveryService": null,
		"deliveryServiceId": 2,
		"fqdn": "example.com",
		"id": 2,
		"ip6Address": null,
		"ipAddress": null,
		"isPrimary": true,
		"lastUpdated": "2018-12-11 15:40:53+00",
		"name": "example",
		"port": 443,
		"profile": null,
		"profileId": null,
		"protocol": "https",
		"tenant": null,
		"tenantId": 1
	}}

``DELETE``
==========
Deletes an origin definition.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Query Parameters

	+------+----------+------------------------------------------------------------------------+
	| Name | Required | Description                                                            |
	+======+==========+========================================================================+
	|  id  | yes      | The integral, unique identifier of the origin definition being deleted |
	+------+----------+------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/1.4/origins?id=2 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: fLaY4/nh0yR38xq5weBKYg02+aQV6Z1ZroOq9UqUCHLMMrH1NMyhOHx+EphPq7JxkjmGY04WCt6VvDyjGWcgfQ==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 11 Dec 2018 17:04:14 GMT
	Content-Length: 61

	{ "alerts": [
		{
			"text": "origin was deleted.",
			"level": "success"
		}
	]}
