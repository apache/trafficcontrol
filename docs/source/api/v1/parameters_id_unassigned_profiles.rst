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

.. _to-api-v1-parameters-id-unassigned_profiles:

*****************************************
``parameters/{{ID}}/unassigned_profiles``
*****************************************
.. deprecated:: ATCv4
.. warning:: There are **very** few good reasons to use this endpoint - be sure to limit said use.

``GET``
=======
Retrieves all :term:`Profiles` to which the specified :term:`Parameter` is *not* assigned.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+--------------------------------------------------------------------------------------------------------+
	| Name | Description                                                                                            |
	+======+========================================================================================================+
	|  ID  | The :ref:`parameter-id` of the :term:`Parameter` for which unassigned :term:`Profiles` shall be listed |
	+------+--------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/parameters/43/unassigned_profiles HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:description:     The :term:`Profile`'s :ref:`profile-description`
:id:              The :term:`Profile`'s :ref:`profile-id`
:lastUpdated:     The date and time at which this :term:`Profile` was last updated, in :ref:`non-rfc-datetime`
:name:            The :term:`Profile`'s :ref:`profile-name`
:routingDisabled: The :term:`Profile`'s :ref:`profile-routing-disabled` setting
:type:            The :term:`Profile`'s :ref:`profile-type`

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
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
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
	],
	"alerts": [
		{
			"level": "warning",
			"text": "This endpoint is deprecated, please use 'GET /profiles' instead"
		}
	]}

.. note:: The Response Example above has been truncated to only its first two array elements, as the true output was very long.
