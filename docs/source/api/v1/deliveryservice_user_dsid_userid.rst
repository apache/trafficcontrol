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

.. _to-api-v1-deliveryservice_user-dsid-userid:

********************************************
``deliveryservice_user/{{dsID}}/{{userID}}``
********************************************
.. deprecated:: ATCv4

``DELETE``
==========
Removes a :term:`Delivery Service` from a user.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| Name   | Description                                                                                                                             |
	+========+=========================================================================================================================================+
	| dsId   | An integral, unique identifier for the :term:`Delivery Service` which should no longer be assigned to the user identified by ``userID`` |
	+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| userId | An integral, unique identifier for the user to whom the :term:`Delivery Service` identified by ``dsID`` should no longer be assigned    |
	+--------+-----------------------------------------------------------------------------------------------------------------------------------------+

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 14 Nov 2018 21:40:06 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: /eNE3LhFABGukcczjxJOYiwmfVTUUKII9RRuZi14AbF65BLhHdXZ5lAVEi4Hc65+ojNaijBgI9jTmgO4XCcP/A==
	Content-Length: 100

	{ "alerts": [
		{
			"level": "success",
			"text": "User [ test ] unlinked from deliveryservice [ 1 | demo1 ]."
		},
		{
			"level": "warning",
			"text": "This endpoint and its functionality is deprecated, and will be removed in the future"
		}
	]}
