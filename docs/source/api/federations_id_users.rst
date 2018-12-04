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

.. _to-api-federations-id-users:

****************************
``federations/{{ID}}/users``
****************************

``GET``
=======
Retrieves users assigned to a federation.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-------------------------------------------------------------------------------------+
	| Name |                 Description                                                         |
	+======+=====================================================================================+
	|  ID  | The integral, unique identifier of the federation for which users will be retrieved |
	+------+-------------------------------------------------------------------------------------+

Response Structure
------------------
:company:  The company to which the user belongs
:email:    The user's email address
:fullName: The user's full name
:id:       An integral, unique identifier for the user
:role:     The user's highest role
:username: The user's short "username"

.. code-block:: json
	:caption: Response Example

	{ "response": [
		{
			"id": 41
			"username": "booya",
			"company": "XYZ Corporation",
			"role": "federation",
			"email": "booya@fooya.com",
			"fullName": "Booya Fooya"
		}
	]}

``POST``
========
Assigns one or more users to a federation.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object

Request Structure
-----------------
:userIds: An array of integral, unique identifiers for users which will be assigned to this federation
:replace: An optional boolean (default: ``false``) which, if ``true``, will cause any conflicting assignments already in place to be overridden by this request

	.. note:: If ``replace`` is not given (and/or not ``true``), then any conflicts with existing assignments will cause the entire operation to fail.

.. code-block:: json
	:caption: Request Example

	{
		"userIds": [ 2, 3, 4, 5, 6 ],
		"replace": true
	}

Response Structure
------------------
:userIds: An array of integral, unique identifiers for users which have been assigned to this federation
:replace: An optional boolean (default: ``false``) which, if ``true``, caused any conflicting assignments already in place to be overridden by this request

.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"level": "success",
			"text": "5 user(s) were assigned to the cname. federation"
		}
	],
	"response": {
		"userIds" : [ 2, 3, 4, 5, 6 ],
		"replace" : true
	}}
