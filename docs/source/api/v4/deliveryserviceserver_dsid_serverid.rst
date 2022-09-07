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

.. _to-api-v4-deliveryserviceserver-dsid-serverid:

***********************************************
``deliveryserviceserver/{{DSID}}/{{serverID}}``
***********************************************

``DELETE``
==========
Removes a :term:`cache server` from a :term:`Delivery Service`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [1]_
:Permissions Required: DELIVERY-SERVICE:READ, DELIVERY-SERVICE:UPDATE, SERVER:READ, SERVER:UPDATE
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+----------+----------+---------------------------------------------------------------+
	| Name     | Required | Description                                                   |
	+==========+==========+===============================================================+
	| dsId     | yes      | An integral, unique identifier for a :term:`Delivery Service` |
	+----------+----------+---------------------------------------------------------------+
	| serverID | yes      | An integral, unique identifier for a server                   |
	+----------+----------+---------------------------------------------------------------+

.. note:: The server identified by ``serverID`` must be a :term:`cache server`, or the assignment will fail.

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
	Whole-Content-Sha512: rGD2sOMHYF0sga1zuTytyLHCUkkc3ZwQRKvZ/HuPzObOP4WztKTOVXB4uhs3iJqBg9zRB2TucMxONHN+3/yShQ==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 01 Nov 2018 14:24:34 GMT
	Content-Length: 80

	{ "alerts": [
		{
			"text": "Server unlinked from delivery service.",
			"level": "success"
		}
	]}

.. [1] Users with the "admin" or "operations" roles will be able to delete *any*:term:`Delivery Service`, whereas other users will only be able to delete :term:`Delivery Services` that their tenant has permissions to delete.
