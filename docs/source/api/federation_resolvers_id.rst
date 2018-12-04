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

.. _to-api-federation_resolvers_id:

*******************************
``federation_resolvers/{{ID}}``
*******************************

``DELETE``
==========
Deletes a federation resolver.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-----------------------------------------------------------------------+
	| Name | Description                                                           |
	+======+=======================================================================+
	|  ID  | Integral, unique identifier for the federation resolver to be deleted |
	+------+-----------------------------------------------------------------------+

Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"level": "success",
			"text": "Federation resolver deleted [ IP = 2.2.2.2/32 ] with id: 27"
		}
	]}
