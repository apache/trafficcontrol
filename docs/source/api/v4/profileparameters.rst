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

.. _to-api-v4-profileparameters:

*********************
``profileparameters``
*********************

``GET``
=======

Retrieves all :term:`Parameter`/:term:`Profile` assignments.

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: PROFILE:READ, PARAMETER:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                                                   |
	+===========+==========+===============================================================================================================+
	| orderby   | no       | Choose the ordering of the results - must be the name of one of the fields of the objects in the ``response`` |
	|           |          | array                                                                                                         |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| sortOrder | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")                      |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| limit     | no       | Choose the maximum number of results to return                                                                |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| offset    | no       | The number of results to skip before beginning to return results. Must use in conjunction with limit          |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| page      | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long   |
	|           |          | and the first page is 1. If ``offset`` was defined, this query parameter has no effect. ``limit`` must be     |
	|           |          | defined to make use of ``page``.                                                                              |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+

Response Structure
------------------
:lastUpdated: The date and time at which this :term:`Profile`/:term:`Parameter` association was last modified, in :ref:`non-rfc-datetime`
:parameter:   The :ref:`parameter-id` of a :term:`Parameter` assigned to ``profile``
:profile:     The :ref:`profile-name` of the :term:`Profile` to which the :term:`Parameter` identified by ``parameter`` is assigned

.. code-block:: http
	:caption: Response Structure

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: +bnMkRgdx4bJoGGlr3mZl539obj3aQAP8e65FAXgywdRAUfXZCFM6VNDn7wScXBmvF2SFXo9F+MhuSwrtB9mPg==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 10 Dec 2018 15:09:13 GMT
	Transfer-Encoding: chunked

	{ "response": [
		{
			"lastUpdated": "2018-12-05 17:50:49+00",
			"profile": "GLOBAL",
			"parameter": 4
		},
		{
			"lastUpdated": "2018-12-05 17:50:49+00",
			"profile": "GLOBAL",
			"parameter": 5
		}
	]}

.. note:: The response example for this endpoint has been truncated to only the first two elements of the resulting array, as the output was hundreds of lines long.

``POST``
========
Associate a :term:`Parameter` to a :term:`Profile`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type: Object
:Permissions Required: PROFILE:READ, PARAMETER:READ, PROFILE:UPDATE

Request Structure
-----------------
This endpoint accepts two formats for the request payload:

Single Object Format
	For assigning a single :term:`Parameter` to a single :term:`Profile`
Array Format
	For making multiple assignments of :term:`Parameters` to :term:`Profiles` simultaneously

Single Object Format
""""""""""""""""""""
:parameterId: The :ref:`parameter-id` of a :term:`Parameter` to assign to some :term:`Profile`
:profileId:   The :ref:`profile-id` of the :term:`Profile` to which the :term:`Parameter` identified by ``parameterId`` will be assigned

.. code-block:: http
	:caption: Request Example - Single Object Format

	POST /api/4.0/profileparameters HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 36
	Content-Type: application/json

	{
		"profileId": 18,
		"parameterId": 1
	}

Array Format
""""""""""""
:parameterId: The :ref:`parameter-id` of a :term:`Parameter` to assign to some :term:`Profile`
:profileId:   The :ref:`profile-id` of the :term:`Profile` to which the :term:`Parameter` identified by ``parameterId`` will be assigned

.. code-block:: http
	:caption: Request Example - Array Format

	POST /api/4.0/profileparameters HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 88
	Content-Type: application/json

	[{
		"profileId": 18,
		"parameterId": 2
	},
	{
		"profileId": 18,
		"parameterId": 3
	}]

Response Structure
------------------
:lastUpdated: The date and time at which the :term:`Profile`/:term:`Parameter` assignment was last modified, in :ref:`non-rfc-datetime`
:parameter:   :ref:`parameter-name` of the :term:`Parameter` which is assigned to ``profile``
:parameterId: The :ref:`parameter-id` of the assigned :term:`Parameter`
:profile:     :ref:`profile-name` of the :term:`Profile` to which the :term:`Parameter` is assigned
:profileId:   The :ref:`profile-id` of the :term:`Profile` to which the :term:`Parameter` identified by ``parameterId`` is assigned

.. code-block:: http
	:caption: Response Example - Single Object Format

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: eDmIwlzX44fZdxLRPHMNa8aoGAK5fQv9Y70A2eeQHfEkliU4evwcsQ4WeHcH0l3/wPTGlpyC0gwLo8LQQpUxWQ==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 10 Dec 2018 13:50:11 GMT
	Content-Length: 166

	{ "alerts": [
		{
			"text": "profileParameter was created.",
			"level": "success"
		}
	],
	"response": {
		"lastUpdated": null,
		"profile": null,
		"profileId": 18,
		"parameter": null,
		"parameterId": 1
	}}
