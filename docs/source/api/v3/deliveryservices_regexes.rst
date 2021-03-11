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


.. _to-api-v3-deliveryservices_regexes:

****************************
``deliveryservices_regexes``
****************************

``GET``
=======
Retrieves routing regular expressions for all :term:`Delivery Services`.

:Auth. Required: Yes
:Roles Required: None\ [1]_
:Response Type:  Array

Request Structure
-----------------
No parameters available

Response Structure
------------------
:dsName:  The name of the :term:`Delivery Service` represented by this object
:regexes: An array of objects that represent various routing regular expressions used by ``dsName``

	:pattern:   The actual regular expression - ``\``\ s are escaped
	:setNumber: The order in which the regular expression is evaluated against requests
	:type:      The type of regular expression - determines that against which it will be evaluated

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: +2MI+Q/NJqTizlMR/MhPAL+yu6/z/Yqvo5fDO8F593RMOmK6dX/Al4wARbEG+HQaJNgSCRPsiLVATusrmnnCMA==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 27 Nov 2018 19:22:59 GMT
	Content-Length: 110

	{ "response": [
		{
			"regexes": [
				{
					"type": "HOST_REGEXP",
					"setNumber": 0,
					"pattern": ".*\\.demo1\\..*"
				}
			],
			"dsName": "demo1"
		}
	]}

.. [1] If tenancy is used, then users (regardless of role) will only be able to see the routing regular expressions used by :term:`Delivery Services` their tenant has permissions to see.
