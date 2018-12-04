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

.. _to-api-federations-id-deliveryservices-id:

************************************************
``federations/{{ID}}/deliveryservices/{{dsID}}``
************************************************

``DELETE``
==========
Removes a Delivery Service from a federation.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+--------------------------------------------------------------------------------------------------------------------------+
	| Name | Description                                                                                                              |
	+======+==========================================================================================================================+
	|  ID  | The integral, unique identifier of the federation from which the Delivery Service identified by ``dsID`` will be removed |
	+------+--------------------------------------------------------------------------------------------------------------------------+
	| dsID | The integral, unique identifier of the Delivery Service which will be removed from the federation identified by ``ID``   |
	+------+--------------------------------------------------------------------------------------------------------------------------+

Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"level": "success",
			"text": "Removed delivery service [ booya-12 ] from federation [ cname1. ]"
		}
	]}
