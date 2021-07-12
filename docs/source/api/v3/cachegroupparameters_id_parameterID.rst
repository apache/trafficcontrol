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

.. _to-api-v3-cachegroupparameters-id-parameterID:

***********************************************
``cachegroupparameters/{{ID}}/{{parameterID}}``
***********************************************

.. deprecated:: ATCv6

``DELETE``
==========
Dissociate a :term:`Parameter` with a :term:`Cache Group`

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+-------------+----------------------------------------------------------------------------------------------------------------+
	| Name        | Description                                                                                                    |
	+=============+================================================================================================================+
	| ID          | The :ref:`cache-group-id` of the :term:`Cache Group` which will have the :term:`Parameter` association deleted |
	+-------------+----------------------------------------------------------------------------------------------------------------+
	| parameterID | The :ref:`parameter-id` of the :term:`Parameter` which will be removed from a :term:`Cache Group`              |
	+-------------+----------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/3.0/cachegroupparameters/8/124 HTTP/1.1
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
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 14 Nov 2018 18:26:40 GMT
	X-Server-Name: traffic_ops_golang/
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: Cuj+ZPAKsDLp4FpbJDcwsWY0yVQAi1Um1CWraeTIQEMlyJSBEm17oKQWDjzTrvqqV8Prhu3gzlcHoVPzEpbQ1Q==
	Content-Length: 84

	{ "alerts": [
		{
			"level": "success",
			"text": "cachegroup parameter was deleted."
		}
	]}
