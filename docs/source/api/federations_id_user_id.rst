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

.. _to-api-federations-id-users-id:

***************************************
``federations/{{ID}}/users/{{userID}}``
***************************************

``DELETE``
==========
Removes a user from a federation.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+--------+----------------------------------------------------------------------------------------------------------------+
	|  Name  | Description                                                                                                    |
	+========+================================================================================================================+
	|   ID   | An integral, unique identifier for the federation from which the user identified by ``userID`` will be removed |
	+--------+----------------------------------------------------------------------------------------------------------------+
	| userID | An integral, unique identifier for the user who will be removed from the federation identified by ``ID``       |
	+--------+----------------------------------------------------------------------------------------------------------------+

Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"level": "success",
			"text": "Removed user [ bobmack ] from federation [ cname1. ]"
		}
	]}
