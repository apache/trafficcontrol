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

.. _to-api-servers-hostname-update_status:

**************************************
``servers/{{hostname}}/update_status``
**************************************

.. note:: This endpoint only truly has meaning for :term:`cache servers`, though it will return a valid response for any server configured in Traffic Ops.

``GET``
=======
Retrieves information regarding pending updates and :term:`Content Invalidation Jobs` for a given server

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: SERVER:READ
:Response Type: Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+----------+----------------------------------------------------+
	| Name     | Description                                        |
	+==========+====================================================+
	| hostname | The (short) hostname of the server being inspected |
	+----------+----------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/5.0/servers/edge/update_status HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
Each object in the returned array\ [#uniqueness]_ will contain the following fields:

:configUpdateTime:     The last time an update was requested for this server. This field defaults to standard epoch
:configApplyTime:      The last time an update was applied for this server. This field defaults to standard epoch
:host_id:              The integral, unique identifier for the server for which the other fields in this object represent the pending updates and revalidation status
:host_name:            The (short) hostname of the server for which the other fields in this object represent the pending updates and revalidation status
:parent_pending:       A boolean telling whether or not any :term:`Topology` ancestor or :term:`parent` of this server has pending updates
:parent_reval_pending: A boolean telling whether or not any :term:`Topology` ancestor or :term:`parent` of this server has pending :term:`Content Invalidation Jobs`
:reval_pending:        ``true`` if the server has pending :term:`Content Invalidation Jobs`, ``false`` otherwise
:revalUpdateTime:      The last time a content invalidation/revalidation request was submitted for this server. This field defaults to standard epoch
:revalApplyTime:       The last time a content invalidation/revalidation request was applied by this server. This field defaults to standard epoch
:status:               The name of the status of this server

	.. seealso:: :ref:`health-proto` gives more information on how these statuses are used, and the ``GET`` method of the :ref:`to-api-statuses` endpoint can be used to retrieve information about all server statuses configured in Traffic Ops.

:upd_pending:       ``true`` if the server has pending updates, ``false`` otherwise
:use_reval_pending: A boolean which tells :term:`ORT` whether or not this version of Traffic Ops should use pending :term:`Content Invalidation Jobs`

	.. note:: This field was introduced to give :term:`ORT` the ability to work with Traffic Control versions 1.x and 2.x seamlessly - as of Traffic Control v3.0 there is no reason for this field to ever be ``false``.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: R6BjNVrcecHGn3eGDqQ1yDiBnEDGQe7QtOMIsRwlpck9SZR8chRQznrkTF3YdROAZ1l8BxR3fXTIvKHIzK2/dA==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 04 Feb 2019 16:24:01 GMT
	Content-Length: 174

	{ "response": [{
		"host_name": "edge",
		"upd_pending": false,
		"reval_pending": false,
		"use_reval_pending": true,
		"host_id": 10,
		"status": "REPORTED",
		"parent_pending": false,
		"parent_reval_pending": false,
		"config_update_time": "2022-02-18T13:52:47.129174-07:00",
		"config_apply_time": "2022-02-18T13:52:47.129174-07:00",
		"revalidate_update_time": "2022-02-28T15:44:15.895145-07:00",
		"revalidate_apply_time": "2022-02-18T13:52:47.129174-07:00"
	}]}

.. [#uniqueness] The returned object is an array, and there is no guarantee that one server exists for a given hostname. However, for each server in the array, that server's update status will be accurate for the server with that particular server ID.
