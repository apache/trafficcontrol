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

.. _to-api-v4-deliveryservices_required_capabilities:

******************************************
``deliveryservices_required_capabilities``
******************************************

.. deprecated:: 4.1
	This endpoint will be removed in a future release, in favor of :ref:`ds-required-capabilities` being a part of :term:`Delivery Services`.

.. versionremoved:: 5.0

``GET``
=======
Gets all associations of :term:`Server Capability` to :term:`Delivery Services`.

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: DELIVERY-SERVICE:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+--------------------+----------+---------------------------------------------------------------------------------------------------------------+
	| Name               | Required | Description                                                                                                   |
	+====================+==========+===============================================================================================================+
	| deliveryServiceID  | no       | Filter :term:`Server Capability` associations by :term:`Delivery Service` integral, unique identifier         |
	+--------------------+----------+---------------------------------------------------------------------------------------------------------------+
	| xmlID              | no       | Filter :term:`Server Capability` associations by :term:`Delivery Service` :ref:`ds-xmlid`                     |
	+--------------------+----------+---------------------------------------------------------------------------------------------------------------+
	| requiredCapability | no       | Filter :term:`Server Capability` associations by :term:`Server Capability` name                               |
	+--------------------+----------+---------------------------------------------------------------------------------------------------------------+
	| orderby            | no       | Choose the ordering of the results - must be the name of one of the fields of the objects in the ``response`` |
	|                    |          | array                                                                                                         |
	+--------------------+----------+---------------------------------------------------------------------------------------------------------------+
	| sortOrder          | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")                      |
	+--------------------+----------+---------------------------------------------------------------------------------------------------------------+
	| limit              | no       | Choose the maximum number of results to return                                                                |
	+--------------------+----------+---------------------------------------------------------------------------------------------------------------+
	| offset             | no       | The number of results to skip before beginning to return results. Must use in conjunction with limit.         |
	+--------------------+----------+---------------------------------------------------------------------------------------------------------------+
	| page               | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long   |
	|                    |          | and the first page is 1. If ``offset`` was defined, this query parameter has no effect. ``limit`` must be     |
	|                    |          | defined to make use of ``page``.                                                                              |
	+--------------------+----------+---------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/4.0/deliveryservices_required_capabilities HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:deliveryServiceID:   The associated :term:`Delivery Service`'s integral, unique identifier
:xmlID:               The associated :term:`Delivery Service`'s :ref:`ds-xmlid`
:lastUpdated:         The date and time at which this association between the :term:`Delivery Service` and the :term:`Server Capability` was last updated, in :ref:`non-rfc-datetime`
:requiredCapability:  The :term:`Server Capability`'s name

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: UFO3/jcBFmFZM7CsrsIwTfPc5v8gUiXqJm6BNp1boPb4EQBnWNXZh/DbBwhMAOJoeqDImoDlrLnrVjQGO4AooA==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 07 Oct 2019 22:15:11 GMT
	Content-Length: 396

	{
		"response": [
			{
				"deliveryServiceID": 1,
				"lastUpdated": "2019-10-07 22:05:31+00",
				"requiredCapability": "ram",
				"xmlId": "example_ds-1"
			},
			{
				"deliveryServiceID": 2,
				"lastUpdated": "2019-10-07 22:05:31+00",
				"requiredCapability": "disk",
				"xmlId": "example_ds-2"
			}
		]
	}

``POST``
========
Associates a :term:`Server Capability` with a :term:`Delivery Service`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: DELIVERY-SERVICE:READ, DELIVERY-SERVICE:UPDATE
:Response Type:  Object

.. note:: A :term:`Server Capability` can only be made required on a :term:`Delivery Service` if its associated Servers already have that :term:`Server Capability` assigned.

Request Structure
-----------------
:deliveryServiceID:   The integral, unique identifier of the :term:`Delivery Service` to be associated
:requiredCapability:  The name of the :term:`Server Capability` to be associated

.. code-block:: http
	:caption: Request Example

	POST /api/4.0/deliveryservices_required_capabilities HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 56
	Content-Type: application/json

	{
		"deliveryServiceID": 1,
		"requiredCapability": "disk"
	}

Response Structure
------------------
:deliveryServiceID:   The newly associated :term:`Delivery Service`'s integral, unique identifier
:lastUpdated:         The date and time at which this association between the :term:`Delivery Service` and the :term:`Server Capability` was last updated, in :ref:`non-rfc-datetime`
:requiredCapability:  The newly associated :term:`Server Capability`'s name

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: eQrl48zWids0kDpfCYmmtYMpegjnFxfOVvlBYxxLSfp7P7p6oWX4uiC+/Cfh2X9i3G+MQ36eH95gukJqOBOGbQ==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 07 Oct 2019 22:15:11 GMT
	Content-Length: 287

	{
		"alerts": [
			{
				"level": "success",
				"text": "deliveryservice.RequiredCapability was created."
			}
		],
		"response": {
			"deliveryServiceID": 1,
			"lastUpdated": "2019-10-07 22:15:11+00",
			"requiredCapability": "disk"
		}
	}

``DELETE``
==========
Dissociate a :term:`Server Capability` from a :term:`Delivery Service`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: DELIVERY-SERVICE:READ, DELIVERY-SERVICE:UPDATE
:Response Type:  ``undefined``

Request Structure
-----------------
:deliveryServiceID:   The integral, unique identifier of the :term:`Delivery Service` from which a :term:`Server Capability` will be dissociated
:requiredCapability:  The name of the :term:`Server Capability` to dissociate

.. code-block:: http
	:caption: Request Example

	POST /api/4.0/deliveryservices_required_capabilities HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 56
	Content-Type: application/json

	{
		"deliveryServiceID": 1,
		"requiredCapability": "disk"
	}

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: eQrl48zWids0kDpfCYmmtYMpegjnFxfOVvlBYxxLSfp7P7p6oWX4uiC+/Cfh2X9i3G+MQ36eH95gukJqOBOGbQ==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 07 Oct 2019 22:15:11 GMT
	Content-Length: 127

	{
		"alerts": [
			{
				"level": "success",
				"text": "deliveryservice.RequiredCapability was deleted."
			}
		]
	}
