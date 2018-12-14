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

.. _to-api-parameters-id-profiles:

******************************
``parameters/{{ID}}/profiles``
******************************
.. deprecated:: 1.1
	Use the ``param`` query parameter of :ref:`to-api-profiles` instead.

``GET``
=======
Retrieves all profiles assigned to the parameter.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+--------------------------------------------------------------------------------------------+
	| Name |                    Description                                                             |
	+======+============================================================================================+
	|  ID  | An integral, unique identifier that specifies for which parameter shall profiles be listed |
	+------+--------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Structure

	GET /api/1.4/parameters/4/profiles HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:description:     A description of profile
:id:              An integral, unique identifier for this profile
:lastUpdated:     The date and time at which this profile was last updated
:name:            Profile name
:routingDisabled: An integer that defines whether or not Traffic Routers will route to servers using these profiles - can only be one of:

	0
		Traffic Routers will route traffic to these servers normally
	1
		Traffic Routers will ignore these servers, and not route traffic to them

:type: The profile's type

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 05 Dec 2018 20:51:23 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Thu, 06 Dec 2018 00:51:23 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: y5fA9q1VogDGxL66ka+ofTtLo3JiTj+Bdrvc4DnfrjFyzqll+537WySFj1nE0C29Twx5l/C8JEHy3Byaz/wbfA==
	Content-Length: 184

	{ "response": [
		{
			"routingDisabled": 0,
			"lastUpdated": "2018-12-05 17:50:49.007102+00",
			"name": "GLOBAL",
			"type": "UNK_PROFILE",
			"id": 1,
			"description": "Global Traffic Ops profile, DO NOT DELETE"
		}
	]}
