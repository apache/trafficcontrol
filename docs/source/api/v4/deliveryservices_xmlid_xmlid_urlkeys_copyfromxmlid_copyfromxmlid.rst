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

.. _to-api-v4-deliveryservices-xmlid-xml_id-urlkeys-copyFrom_xml_id:

*******************************************************************************
``deliveryservices/xmlId/{{xml_id}}/urlkeys/copyFromXmlId/{{copyFrom_xml_id}}``
*******************************************************************************

``POST``
========
Allows a user to copy URL signing keys from a specified :term:`Delivery Service` to another :term:`Delivery Service`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: DS-SECURITY-KEY:READ, DS-SECURITY-KEY:CREATE, DELIVERY-SERVICE:READ, DELIVERY-SERVICE:UPDATE
:Response Type:  Object (string)

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------------+--------------------------------------------------------------------------------------+
	| Name            | Description                                                                          |
	+=================+======================================================================================+
	| xml_id          | The :ref:`ds-xmlid` of the :term:`Delivery Service` *to* which keys will be copied   |
	+-----------------+--------------------------------------------------------------------------------------+
	| copyFrom_xml_id | The :ref:`ds-xmlid` of the :term:`Delivery Service` *from* which keys will be copied |
	+-----------------+--------------------------------------------------------------------------------------+

Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{
		"response": "Successfully copied and stored keys"
	}
