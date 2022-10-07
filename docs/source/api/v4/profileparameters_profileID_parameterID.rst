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

.. _to-api-v4-profileparameters-profileID-parameterID:

***************************************************
``profileparameters/{{profileID}}/{{parameterID}}``
***************************************************

``DELETE``
==========
Deletes a :term:`Profile`/:term:`Parameter` association.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: PROFILE:UPDATE, PROFILE:READ, PARAMETER:READ
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+-------------+------------------------------------------------------------------------------------------------------------------------------+
	|    Name     | Description                                                                                                                  |
	+=============+==============================================================================================================================+
	|  profileID  | The :ref:`profile-id` of the :term:`Profile` from which a :term:`Parameter` shall be removed                                 |
	+-------------+------------------------------------------------------------------------------------------------------------------------------+
	| parameterID | The :ref:`parameter-id` of the :term:`Parameter` which shall be removed from the :term:`Profile` identified by ``profileID`` |
	+-------------+------------------------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/4.0/profileparameters/18/129 HTTP/1.1
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
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: JQuBqHyT9MnNwO9NSIDVQhkRtXdeAJc95W1pF2dwQeoBFmf0Y8knXm3/O/rbJDEoUC7DhUQN1aoYIsqqmz4qQQ==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 10 Dec 2018 15:00:15 GMT
	Content-Length: 71

	{ "alerts": [
		{
			"text": "profileParameter was deleted.",
			"level": "success"
		}
	]}
