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

.. _to-api-cdns-name-name:

*******************************
``/api/1.x/cdns/name/{{name}}``
*******************************

``GET``
=======
Extract information about a CDN, identified by name.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------+----------+---------------------------------------------+
	|   Name    | Required |                Description                  |
	+===========+==========+=============================================+
	|  ``name`` |   yes    | The name of the CDN to be inspected         |
	+-----------+----------+---------------------------------------------+

Response Structure
------------------
:dnssecEnabled: ``true`` if DNSSEC is enabled on this CDN, otherwise ``false``
:domainName:    Top Level Domain name within which this CDN operates
:id:            The integral, unique identifier for the CDN
:lastUpdated:   Date and time when the CDN was last modified in ISO format
:name:          The name of the CDN

.. code-block:: json
	:caption: Response Example

	{ "response": [
		{
			"dnssecEnabled": false,
			"domainName": "mycdn.ciab.test",
			"id": 2,
			"lastUpdated": "2018-10-16 20:10:49+00",
			"name": "CDN-in-a-Box"
		}
	]}


``DELETE``
==========
Allows a user to delete a CDN by name

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------+----------+---------------------------------------------+
	|   Name    | Required |                Description                  |
	+===========+==========+=============================================+
	|  ``name`` |   yes    | The name of the CDN to be inspected         |
	+-----------+----------+---------------------------------------------+

Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"text": "cdn was deleted.",
			"level": "success"
		}
	]}
