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

.. _to-api-v4-cdn-locks:

*****************
``cdn_locks``
*****************

.. versionadded:: 4.0

``GET``
=======
Gets information for all CDN locks.

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: CDN:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+---------------+----------+-----------------------------------------------------------------------------------+
	| Parameter     | Required | Description                                                                       |
	+===============+==========+===================================================================================+
	| username      | no       | Return only the CDN lock that the user with ``username`` possesses                |
	+---------------+----------+-----------------------------------------------------------------------------------+
	| cdn           | no       | Return only the CDN lock for the CDN that has the name ``cdn``                    |
	+---------------+----------+-----------------------------------------------------------------------------------+

Response Structure
------------------
:userName:         The username for which the lock exists.
:cdn:              The name of the CDN for which the lock exists.
:message:          The message or reason that the user specified while acquiring the lock.
:soft:             Whether or not this is a soft(shared) lock.
:sharedUserNames:  An array of the usernames that the creator of the lock has shared their lock with.
:lastUpdated:      Time that this lock was last updated(created).

.. code-block:: http
	:caption: Response Example

	HTTP/2 200
	Content-Type: application/json

	{ "response": [
		{
			"userName": "foo",
			"cdn": "bar",
			"message": "acquiring lock to snap CDN",
			"soft": true,
			"sharedUserNames": [
				"user1"
			],
			"lastUpdated": "2021-05-26T09:31:57-06"
		}
	]}

``POST``
========
Allows user to acquire a lock on a CDN.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: CDN-LOCK:CREATE, CDN:READ
:Response Type:  Object

Request Structure
-----------------
The request body must be a single ``CDN Lock`` object with the following keys:

:cdn:             The name of the CDN for which the user wants to acquire a lock.
:message:         The message or reason for the user to acquire the lock. This is an optional field.
:sharedUserNames: An array of the usernames that the creator of the lock wants to share their lock with. This is an optional field.
:soft:            Whether or not this is a soft(shared) lock. This is an optional field; ``soft`` will be set to ``true`` by default.

.. code-block:: http
	:caption: Request Example

	POST /api/4.0/cdn_locks HTTP/2
	Host: localhost:8443
	User-Agent: curl/7.64.2
	Accept: */*
	Cookie: mojolicious=...
	Content-Type: application/json
	Content-Length: 81

	{
		"cdn": "bar",
		"message": "acquiring lock to snap CDN",
		"sharedUserNames": [
				"user1"
		],
		"soft": true
	}

Response Structure
------------------
:userName:         The username for which the lock was created.
:cdn:              The name of the CDN for which the lock was created.
:message:          The message or reason that the user specified while acquiring the lock.
:soft:             Whether or not this is a soft(shared) lock.
:sharedUserNames:  An array of the usernames that the creator of the lock has shared their lock with.
:lastUpdated:      Time that this lock was last updated(created).

.. code-block:: http
	:caption: Response Example

	HTTP/2 201
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Wed, 26 May 2021 17:59:10 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: IWjt4zhg4OlPDTfOebjMTS1uHsZ8LycEaHgSS3KHnmc6Vvmw5/S6q70CCnbAePV2x1bxKkVEifTIxfft8vq3sg==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 26 May 2021 16:59:10 GMT
	Content-Length: 204

	{ "alerts": [
		{
			"text": "soft CDN lock acquired!",
			"level":"success"
		}
	],
	"response": {
		"userName": "foo",
		"cdn": "bar",
		"message": "acquiring lock to snap CDN",
		"soft": true,
		"sharedUserNames": [
			"user1"
		],
		"lastUpdated": "2021-05-26T10:59:10-06"
	}}

``DELETE``
----------
Deletes an existing ``CDN Lock``.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: CDN-LOCK:DELETE, CDN:READ
:Response Type: Object

Request Structure
-----------------
.. table:: Request Query Parameters

	+---------------+----------+-----------------------------------------------------------------------------------+
	| Parameter     | Required | Description                                                                       |
	+===============+==========+===================================================================================+
	| cdn           | yes      | Delete the CDN lock for the CDN that has the name ``cdn``                         |
	+---------------+----------+-----------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/4.0/cdn_locks?cdn=bar HTTP/2
	Host: localhost:8443
	User-Agent: curl/7.64.1
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 0
	Content-Type: application/json

Response Structure
------------------

.. code-block:: http
	:caption: Response Example

	HTTP/2 200
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Wed, 26 May 2021 22:20:10 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: p/M2OEmhaws6QLhzzoSBvpC5UnIM+/84RI1wO42PYXiyUKWnxoQQEtm4lkN+K5NOKIH+OkyUlI2ovQZP6lGOcg==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 26 May 2021 21:20:10 GMT
	Content-Length: 202

	{ "alerts": [
		{
			"text": "cdn lock deleted",
			"level":"success"
		}
	],
	"response": {
		"userName": "foo",
		"cdn": "bar",
		"message": "acquiring lock to snap CDN",
		"soft": true,
		"sharedUserNames": [
			"user1"
		],
		"lastUpdated": "2021-05-26T10:59:10-06"
	}}
	
Response Structure
-----------------
:HTTPstatuscode: This would typically be a 204 No Content for a successful deletion, or a 4xx status code (such as 400 Bad Request, 401 Unauthorized, or 403 Forbidden) if the deletion was not successful.   	
:Message:   This would provide additional information about the outcome of the request. For example, it may include an error message if the request failed, or a success message if the lock was deleted. 
:Additionaldata: Depending on the specific implementation of the API, the response object may include additional data such as the ID of the deleted lock, the timestamp of when the lock was deleted, or other information relevant to the request.




