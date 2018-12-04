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

.. _to-api-federations-id-federation_resolvers:

*******************************************
``federations/{{ID}}/federation_resolvers``
*******************************************

``GET``
=======
Retrieves federation resolvers assigned to a federation.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------------------------------------+
	| Name |                 Description                                                              |
	+======+==========================================================================================+
	|  ID  | The integral, unique identifier for the federation for which resolvers will be retrieved |
	+------+------------------------------------------------------------------------------------------+

Response Structure
------------------
:id:        The integral, unique identifier of this federation resolver
:ipAddress: The IP address of the federation resolver - may be IPv4 or IPv6
:type:      The type of resolver - one of:

	RESOLVE4
		This resolver is for IPv4 addresses (and ``ipAddress`` is IPv4)
	RESOLVE6
		This resolver is for IPv6 addresses (and ``ipAddress`` is IPv6)

.. code-block:: json
	:caption: Response Example

	{ "response": [
		{
			"id": 41
			"ipAddress": "2.2.2.2/16",
			"type": "RESOLVE4"
		}
	]}

``POST``
========
Assigns one or more resolvers to a federation.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------------------------------------+
	| Name |                 Description                                                              |
	+======+==========================================================================================+
	|  ID  | The integral, unique identifier for the federation for which resolvers will be retrieved |
	+------+------------------------------------------------------------------------------------------+

:fedResolverIds: An array of integral, unique identifiers for federation resolvers
:replace:        An optional boolean (default: ``false``) which, if ``true``, will cause any conflicting assignments already in place to be overridden by this request

	.. note:: If ``replace`` is not given (and/or not ``true``), then any conflicts with existing assignments will cause the entire operation to fail.

.. code-block:: json
	:caption: Request Example

	{
		"fedResolverIds": [ 2, 3, 4, 5, 6 ],
		"replace": true
	}

Response Structure
------------------
:fedResolverIds: An array of integral, unique identifiers for federation resolvers
:replace:        An optionally-present boolean (default: ``false``) which, if ``true``, any conflicting assignments already in place were overridden by this request

.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"level": "success",
			"text": "5 resolvers(s) were assigned to the cname. federation"
		}
	],
	"response": {
		"fedResolverIds" : [ 2, 3, 4, 5, 6 ],
		"replace" : true
	}}
