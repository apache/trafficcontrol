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

.. _to-api-profileparameters:

*********************
``profileparameters``
*********************

``GET``
=======
.. deprecated:: 1.1
	To get the profiles associated with a particular parameter, use the ``param`` query parameter of :ref:`to-api-profiles` instead. To see the parameters associated with a particular profile, refer to the ``params`` key in the response of a ``GET`` request to :ref:`to-api-profiles-id` instead.

Retrieves all parameter/profile assignments.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
No parameters available

Response Structure
------------------
:lastUpdated: The date and time at which this profile/parameter association was last modified
:parameter:   An integral, unique identifier for a parameter assigned to ``profile``
:profile:     The name of the profile to which the parameter identified by ``parameter`` is assigned

.. code-block:: http
	:caption: Response Structure

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
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
Associate parameter to profile.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
This endpoint accepts two formats for the request payload:

Single Object Format
	For assigning a single parameter to a single profile
Array Format
	For making multiple assignments of parameters to profiles simultaneously

Single Object Format
""""""""""""""""""""
:parameterId: The integral, unique identifier of a parameter to assign to some profile
:profileId:   The integral, unique identifier of the profile to which the parameter identified by ``parameterId`` will be assigned

.. code-block:: http
	:caption: Request Example - Single Object Format

	POST /api/1.4/profileparameters HTTP/1.1
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
.. caution:: Array format is broken as of the time of this writing. Follow `GitHub Issue #3103 <https://github.com/apache/trafficcontrol/issues/3103>`_ for further developments.

:parameterId: The integral, unique identifier of a parameter to assign to some profile
:profileId:   The integral, unique identifier of the profile to which the parameter identified by ``parameterId`` will be assigned

.. code-block:: http
	:caption: Request Example - Array Format

	POST /api/1.4/profileparameters HTTP/1.1
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
:lastUpdated: The date and time at which the profile/parameter assignment was last modified, in ISO format
:parameter:   Name of the parameter which is assigned to ``profile``
:parameterId: The integral, unique identifier of the assigned parameter
:profile:     Name of the profile to which the parameter is assigned
:profileId:   The integral, unique identifier of the profile to which the parameter identified by ``parameterId`` is assigned

.. code-block:: http
	:caption: Response Example - Single Object Format

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
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
