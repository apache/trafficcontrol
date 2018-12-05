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

.. _to-api-parameters-id-unassigned_profiles:

*****************************************
``parameters/{{ID}}/unassigned_profiles``
*****************************************
.. warning:: There are **very** few good reasons to use this endpoint - be sure not limit said use.

``GET``
=======
Retrieves all profiles to which the specified parameter is NOT assigned to the parameter.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-------------------------------------------------------------------------------------------------------+
	| Name |                    Description                                                                        |
	+======+=======================================================================================================+
	|  ID  | An integral, unique identifier that specifies for which parameter unassigned profiles shall be listed |
	+------+-------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/parameters/43/unassigned_profiles HTTP/1.1
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
	Date: Wed, 05 Dec 2018 21:47:48 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Thu, 06 Dec 2018 01:47:48 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: c669pfp2I0FB2xQ1v6RKDbuld5cXvAtGnr7wEzW4ku+7EimNvEyBBPPW4G7FLRQvwO0y/0hWoJcm4/ZYGBR39g==
	Transfer-Encoding: chunked

	{ "response": [
		{
			"cdn": 1,
			"lastUpdated": "2018-12-05 17:50:49.007102+00",
			"name": "GLOBAL",
			"description": "Global Traffic Ops profile, DO NOT DELETE",
			"cdnName": "ALL",
			"routingDisabled": false,
			"id": 1,
			"type": "UNK_PROFILE"
		},
		{
			"cdn": 1,
			"lastUpdated": "2018-12-05 17:50:49.024653+00",
			"name": "TRAFFIC_ANALYTICS",
			"description": "Traffic Analytics profile",
			"cdnName": "ALL",
			"routingDisabled": false,
			"id": 2,
			"type": "UNK_PROFILE"
		}
	]}

.. note:: The Response Example above has been truncated to only its first two array elements, as the true output was very long.
