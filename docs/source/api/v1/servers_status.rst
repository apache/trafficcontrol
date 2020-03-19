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

.. _to-api-v1-servers-status:

******************
``servers/status``
******************
.. deprecated:: ATCv4
	Use the ``GET`` method of :ref:`to-api-v1-servers` and filter client side instead.

``GET``
=======
Retrieves an aggregated view of all server statuses across all CDNs

:Auth. Required: Yes
:Roles Required: None
:Response Type: Object

Request Structure
-----------------
.. table:: Request Query Parameters

	+------------+----------+-------------------------------------------------------------------------------------------------------------------+
	| Name       | Required | Description                                                                                                       |
	+============+==========+===================================================================================================================+
	| type       | no       | Return status counts for only servers of this :term:`Type`                                                        |
	+------------+----------+-------------------------------------------------------------------------------------------------------------------+

Response Structure
------------------
:status: Every key in the ``response`` object will be the name of a valid server status, with a value that is the number of servers with that status. If there are no servers with a given status, that status will not appear as a key.

	.. seealso:: :ref:`to-api-v1-statuses` can be queried to retrieve all possible server statuses, as well as to create new statuses or modify existing statuses.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Mon, 04 Feb 2019 16:22:14 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: M072YRXvtNwjnCfntv/W3AsSpOhCl7Cpm0UDznOcXxwwgRYSGXx2MoeovXSNzYim62FJJoQJom1ccRSAW9ZMcA==
	Content-Length: 38

	{ "response": {
		"REPORTED": 2,
		"ONLINE": 9
	},
	"alerts": [
		{
			"text": "This endpoint is deprecated, please use GET /servers instead",
			"level": "warning"
		}
	]}
