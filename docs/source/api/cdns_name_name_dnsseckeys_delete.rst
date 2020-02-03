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

.. _to-api-cdns-name-name-dnsseckeys-delete:

****************************************
``cdns/name/{{name}}/dnsseckeys/delete``
****************************************

``GET``
=======
Delete DNSSEC keys for a CDN and all associated :term:`Delivery Services`.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object (string)

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-----------------------------------------------------------+
	| Name |                       Description                         |
	+======+===========================================================+
	| name | The name of the CDN for which DNSSEC keys will be deleted |
	+------+-----------------------------------------------------------+

Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{
		"response": "Successfully deleted dnssec keys for test"
	}

