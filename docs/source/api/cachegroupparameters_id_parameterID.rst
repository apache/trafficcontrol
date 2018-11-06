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

.. _to-api-cachegroupparameters-id-parameterID:

********************************************
cachegroupparameters/{{ID}}/{{parameter ID}}
********************************************

``DELETE``
==========
Delete a Cache Group parameter association.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+------------------+----------+-----------------------------------------------------------------------------------------+
	| Name             | Required | Description                                                                             |
	+==================+==========+=========================================================================================+
	| ``cachegroup_id``| yes      | Unique identifier for the Cache Group which will have the parameter association deleted |
	+------------------+----------+-----------------------------------------------------------------------------------------+
	| ``parameter_id`` | yes      | Unique identifier for the parameter which will be removed from a Cache Group            |
	+------------------+----------+-----------------------------------------------------------------------------------------+

Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{ "alerts":[
		{
			"level": "success",
			"text": "Cache group parameter association was deleted."
		}
	]}
