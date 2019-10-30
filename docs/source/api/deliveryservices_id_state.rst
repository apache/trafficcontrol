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

.. _to-api-deliveryservices-id-state:

*********************************
``deliveryservices/{{ID}}/state``
*********************************

.. deprecated:: 1.1
	Use :ref:`to-api-cachegroup_fallbacks` instead to configure Cache Group fallbacks

``GET``
=======
Retrieves the fail-over state for a :term:`Delivery Service`.

:Auth. Required: Yes
:Roles Required: None\ [1]_
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------------------------------------------+
	| Name | Description                                                                     |
	+======+=================================================================================+
	|  ID  | The integral, unique identifier of the :term:`Delivery Service` being inspected |
	+------+---------------------------------------------------------------------------------+

Response Structure
------------------
:enabled:  ``true`` if failover has been enabled for this :term:`Delivery Service`, ``false`` otherwise
:failover: An object describing the failover configuration for this :term:`Delivery Service`

	:configured:  ``true`` if this failover configuration has been updated by some Traffic Ops user, ``false`` otherwise
	:destination: An object describing the Cache Group within this :term:`Delivery Service` which will utilize this failover configuration

		:location: The integral, unique identifier of a Cache Group within this :term:`Delivery Service` which will utilize this failover configuration
		:type:     The 'type' of the Cache Group identified by ``location``

	:enabled:   ``true`` if failover has been enabled for this :term:`Delivery Service`, ``false`` otherwise
	:locations: An array of integral, unique identifiers for Cache Groups to use for failover

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Thu, 15 Nov 2018 14:54:17 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Thu, 15 Nov 2018 18:54:17 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: 6dswLNVRAYBxXAQjXu8MfnLpQ94b9HyrL7ROzhF2pw+tBotgU98zhQRQoQrPEwrTVranTxTUyxP2icFfv5vh7g==
	Content-Length: 112

		{ "response": {
			"failover": {
				"locations": [],
				"destination": null,
				"configured": false,
				"enabled": false
			},
			"enabled": false
		}}



.. [1] If a user does not have either the "admin" nor "operations" role, then only :term:`Delivery Services` assigned to the user's Tenant will be able to be queried with this endpoint
