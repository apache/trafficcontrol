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

.. _to-api-federation_resolvers:

************************
``federation_resolvers``
************************

``POST``
========
Creates a new federation resolver.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object

Request Structure
-----------------
:ipAddress: The IP address of the resolver - may be IPv4 or IPv6
:typeId:    The integral, unique identifier of the type of resolver being created - will *represent* one of:

	RESOLVE4
		Resolver is for IPv4 addresses and ``ipAddress`` is IPv4
	RESOLVE6
		Resolver is for IPv6 addresses and ``ipAddress`` is IPv6

.. code-block:: json
	:caption: Request Example

	{
		"ipAddress": "2.2.2.2/32",
		"typeId": 245
	}

Response Structure
------------------
:id:        The integral, unique identifier of the resolver
:ipAddress: The IP address of the resolver - may be IPv4 or IPv6
:type:      The type of the resolver - one of:

	RESOLVE4
		Resolver is for IPv4 addresses and ``ipAddress`` is IPv4
	RESOLVE6
		Resolver is for IPv6 addresses and ``ipAddress`` is IPv6

.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"level": "success",
			"text": "Federation resolver created [ IP = 2.2.2.2/32 ] with id: 27"
		}
	],
	"response": {
		"id" : 27,
		"ipAddress" : "2.2.2.2/32",
		"typeId" : 245,
	}}
