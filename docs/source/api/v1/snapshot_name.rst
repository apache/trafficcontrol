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

.. _to-api-v1-snapshot-name:

*********************
``snapshot/{{name}}``
*********************
.. deprecated:: ATCv4
	Use the ``PUT`` method of :ref:`to-api-snapshot` with the query parameter ``cdn`` instead.

``PUT``
=======
Performs a CDN :term:`Snapshot`. Effectively, this propagates the new *configuration* of the CDN to its *operating state*, which replaces the output of the :ref:`to-api-v1-cdns-name-snapshot` endpoint with the output of the :ref:`to-api-v1-cdns-name-snapshot-new` endpoint.
This also changes the output of the :ref:`to-api-v1-cdns-name-configs-monitoring` endpoint since that endpoint returns the latest monitoring information from the *operating state*.

.. Note:: Snapshotting the CDN also deletes all HTTPS certificates for every :term:`Delivery Service` which has been deleted since the last :term:`Snapshot`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-----------------------------------------------------------------+
	| Name | Description                                                     |
	+======+=================================================================+
	| name | The name of the CDN for which a :term:`Snapshot` shall be taken |
	+------+-----------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	PUT /api/1.4/snapshot/CDN-in-a-Box HTTP/1.1
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
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: gmaWI0tKgNFPYO0zMrLCGDosBJkPbeIvW4BH9tEh96VjBqyWqzjgPySoMa3ViM1BQXA6VAUOGmc76VyhBsaTzA==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 18 Mar 2020 15:51:48 GMT
	Content-Length: 160

	{
		"alerts": [
			{
				"text": "This endpoint is deprecated, please use PUT /snapshots with either the query parameter cdn or cdnID instead",
				"level": "warning"
			}
		],
		"response": "SUCCESS"
	}