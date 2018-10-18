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

.. _to-api-cdns-id:

************************
``/api/1.x/cdns/{{ID}}``
************************

``GET``
=======
Extract information about a specific CDN.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------+----------+---------------------------------------------+
	|   Name    | Required |                Description                  |
	+===========+==========+=============================================+
	|   ``id``  |   yes    | ID of the CDN to inspect                    |
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

``PUT``
=======
Allows a user to edit a specific CDN

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------+----------+---------------------------------------------+
	|   Name    | Required |                Description                  |
	+===========+==========+=============================================+
	|   ``id``  |   yes    | ID of the CDN to inspect                    |
	+-----------+----------+---------------------------------------------+

.. table:: Request Data Parameters

	+-------------------+---------+----------+-----------------------------------------------------------+
	|    Parameter      |  Type   | Required |        Description                                        |
	+===================+=========+==========+===========================================================+
	| ``name``          | string  | yes      | Name of the new CDN                                       |
	+-------------------+---------+----------+-----------------------------------------------------------+
	| ``domainName``    | string  | yes      | The top-level domain (TLD) belonging to the new CDN       |
	+-------------------+---------+----------+-----------------------------------------------------------+
	| ``dnssecEnabled`` | boolean | yes      | ``true`` if this CDN will use DNSSEC, ``false`` otherwise |
	+-------------------+---------+----------+-----------------------------------------------------------+

Response Structure
------------------
:dnssecEnabled: ``true`` if the CDN uses DNSSEC, ``false`` otherwise
:domainName:    The top-level domain (TLD) assigned to the newly created CDN
:id:            An integral, unique identifier for the newly created CDN
:name:          The newly created CDN's name


.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"text": "cdn was updated.",
			"level": "success"
		}
	],
	"response": {
		"dnssecEnabled": false,
		"domainName": "foo.ciab.test",
		"id": 7,
		"lastUpdated": "2018-10-17 17:37:34+00",
		"name": "Foo"
	}}

``DELETE``
==========
Allows a user to delete a specific CDN

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------+----------+---------------------------------------------+
	|   Name    | Required |                Description                  |
	+===========+==========+=============================================+
	|   ``id``  |   yes    | ID of the CDN to inspect                    |
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
