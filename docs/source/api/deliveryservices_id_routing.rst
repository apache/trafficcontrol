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

.. _to-api-deliveryservices-id-routing:

*****************************************
``/api/1.2/deliveryservices/:id/routing``
*****************************************

``GET``
=======
Retrieves routing method statistics for a particular Delivery Service

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [1]_
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------------+----------+----------------------------------------------------------------------+
	| Name            | Required | Description                                                          |
	+=================+==========+======================================================================+
	| id              | yes      | The integral, unique identifier for the Delivery Service of interest |
	+-----------------+----------+----------------------------------------------------------------------+

Response Structure
------------------
:cz:                The percent of requests to the Traffic Router for this Delivery Service that were satisfied by a coverage zone file (CZF)
:dsr:               The percent of requests to the Traffic Router for this Delivery Service that were satisfied by sending the client to an overflow Delivery Service
:err:               The percent of requests to the Traffic Router for this Delivery Service that resulted in an error
:fed:               The percent of requests to the Traffic Router for this Delivery Service that were satisfied by sending the client to a federated CDN
:geo:               The percent of requests to the Traffic Router for this Delivery Service that were satisfied using 3rd party geographic IP mapping
:miss:              The percent of requests to the Traffic Router for this Delivery Service that could not be satisfied
:regionalAlternate: The percent of requests to the Traffic Router for this Delivery Service that were satisfied by sending the client to the alternate, Regional Geo-blocking URL
:regionalDenied:    The percent of Traffic Router requests for this Delivery Service that were denied due to geographic location policy
:staticRoute:       The percent of requests to the Traffic Router for this Delivery Service that were satisfied with pre-configured DNS entries

.. code-block:: json
	:caption: Response Example

	{ "response": {
		"staticRoute": 0,
		"geo": 100,
		"err": 0,
		"fed": 0,
		"cz": 0,
		"dsr": 0,
		"regionalAlternate": 0,
		"deepCz": 0,
		"regionalDenied": 0,
		"miss": 0
	}}

.. [1] Users with the roles "admin" and/or "operations" will be able to see details for *all* Delivery Services, whereas any other user will only see details for the Delivery Services their Tenant is allowed to see.
