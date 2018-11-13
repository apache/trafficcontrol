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

.. _to-api-deliveryservice_user-dsid-userid:

********************************************
``deliveryservice_user/{{dsID}}/{{userID}}``
********************************************

``DELETE``
==========
Removes a Delivery Service from a user.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+--------+----------+---------------------------------------------------------------------------------------------------------------------------------+
	| Name   | Required | Description                                                                                                                     |
	+========+==========+=================================================================================================================================+
	| dsId   | yes      | An integral, unique identifier for the Delivery Service which should no longer be assigned to the user identified by ``userID`` |
	+--------+----------+---------------------------------------------------------------------------------------------------------------------------------+
	| userId | yes      | An integral, unique identifier for the user to whom the Delivery Service identified by ``dsID`` should no longer be assigned    |
	+--------+----------+---------------------------------------------------------------------------------------------------------------------------------+

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Type: application/json

	{ "alerts": [
		{
			"level": "success",
			"text": "User and delivery service were unlinked."
		}
	]}
