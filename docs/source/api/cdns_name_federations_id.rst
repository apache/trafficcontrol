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

.. _to-api-cdns-name-federations-id:

************************************
``cdns/{{name}}/federations/{{ID}}``
************************************

``GET``
=======
Retrieves a specific federation used within a specific CDN.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-------------------------------------------------------------------------------------+
	| Name | Description                                                                         |
	+======+=====================================================================================+
	| name | The name of the CDN for which the federation identified by ``ID`` will be inspected |
	+------+-------------------------------------------------------------------------------------+
	|  ID  | An integral, unique identifier for the federation to be inspected                   |
	+------+-------------------------------------------------------------------------------------+

Response Structure
------------------
:cname:           The Canonical Name (CNAME) used by the federation
:deliveryService: An object with keys that provide identifying information for the Delivery Service using this federation

	:id:    The integral, unique identifer for the Delivery Service
	:xmlId: The Delivery Service's uniquely identifying 'xml_id'

:description: An optionally-present field containing a description of the field

	.. note:: This key will only be present if the description was provided when the federation was created. Refer to the ``POST`` method of the :ref:`to-api-cdns-name-federations` endpoint to see how federations can be created.

:lastUpdated: The date and time at which this federation was last modified, in ISO format
:ttl:         Time to Live (TTL) for the ``cname``, in hours

.. code-block:: json
	:caption: Response Example

	{ "response": [
		{
			"id": 41
			"cname": "booya.com.",
			"ttl": 34,
			"description": "fooya",
			"deliveryService": {
				"id": 61,
				"xmlId": "the-xml-id"
			},
			"lastUpdated": "2018-08-01 14:41:48+00"
		}
	]}

``PUT``
=======
Updates a federation.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-------------------------------------------------------------------------------------+
	| Name | Description                                                                         |
	+======+=====================================================================================+
	| name | The name of the CDN for which the federation identified by ``ID`` will be inspected |
	+------+-------------------------------------------------------------------------------------+
	|  ID  | An integral, unique identifier for the federation to be inspected                   |
	+------+-------------------------------------------------------------------------------------+

:cname: The Canonical Name (CNAME) used by the federation

	.. note:: The CNAME must end with a "``.``"

:description: An optional description of the federation
:ttl:         Time to Live (TTL) for the name record used for ``cname``

.. code-block:: json
	:caption: Request Example

	{
		"cname": "the.cname.com.",
		"ttl": 48,
		"description": "the description"
	}

Response Structure
------------------
:cname:       The Canonical Name (CNAME) used by the federation
:description: An optionally-present field containing a description of the field

	.. note:: This key will only be present if the description was provided when the federation was created

:lastUpdated: The date and time at which this federation was last modified, in ISO format
:ttl:         Time to Live (TTL) for the ``cname``, in hours


.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"level": "success",
			"text": "Federation updated [ cname = the.cname. ] with id: 26."
		}
	],
	"response": {
		"id": 26,
		"cname": "the.cname.com.",
		"ttl": 48,
		"description": "the description",
		"lastUpdated": "2018-08-01 14:41:48+00"
	}}

``DELETE``
==========
Deletes a specific federation.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-------------------------------------------------------------------------------------+
	| Name | Description                                                                         |
	+======+=====================================================================================+
	| name | The name of the CDN for which the federation identified by ``ID`` will be inspected |
	+------+-------------------------------------------------------------------------------------+
	|  ID  | An integral, unique identifier for the federation to be inspected                   |
	+------+-------------------------------------------------------------------------------------+

.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"level": "success",
			"text": "Federation deleted [ cname = the.cname. ] with id: 26."
		}
	]}
