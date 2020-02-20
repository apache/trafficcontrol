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

.. _to-api-v1-deliveryservices-xmlid-xmlid-sslkeys-delete:

***************************************************
``deliveryservices/xmlId/{{xmlid}}/sslkeys/delete``
***************************************************

``GET``
=======
:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object (string)

Request Structure
-----------------
.. table:: Request Path Parameters

	+-------+----------+-------------------------------------------------------------+
	| Name  | Required | Description                                                 |
	+=======+==========+=============================================================+
	| xmlId | yes      | The :ref:`ds-xmlid` of the desired :term:`Delivery Service` |
	+-------+----------+-------------------------------------------------------------+

.. table:: Request Query Parameters

	+---------+----------+------------------------------------------------------------+
	|   Name  | Required |          Description                                       |
	+=========+==========+============================================================+
	| version | no       | The version number of the SSL keys that shall be retrieved |
	+---------+----------+------------------------------------------------------------+

Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{ "response": "Successfully deleted ssl keys for <xml_id>" }
