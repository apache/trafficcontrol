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

.. _to-api-v4-deliveryservices-xmlid-xmlid-sslkeys-renew:

**************************************************
``deliveryservices/xmlId/{{XMLID}}/sslkeys/renew``
**************************************************

``POST``
========
Uses :abbr:`ACME (Automatic Certificate Management Environment)` protocol to renew SSL keys for a :term:`Delivery Service`.

:Auth. Required: Yes
:Roles Required: "admin"
:Permissions Required: ACME:READ, DS-SECURITY-KEY:DELETE, DS-SECURITY-KEY:READ, DS-SECURITY-KEY:CREATE, DS-SECURITY-KEY:UPDATE, DELIVERY-SERVICE:READ, DELIVERY-SERVICE:UPDATE
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+-------+------------------------------------------------------+
	|  Name |              Description                             |
	+=======+======================================================+
	| XMLID | The 'xml_id' of the desired :term:`Delivery Service` |
	+-------+------------------------------------------------------+


Request Structure
-----------------
No parameters available


Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{ "alerts": [{
		"level": "success",
		"text": "Certificate for test-xml-id successfully renewed."
	}]}
