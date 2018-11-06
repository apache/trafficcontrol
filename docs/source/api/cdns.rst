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

.. _to-api-cdns:

****
cdns
****
Extract information about all CDNs

``GET``
=======
:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
No parameters available

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
			"id": 1,
			"name": "over-the-top",
			"dnssecEnabled": false,
			"lastUpdated": "2014-10-02 08:22:43",
			"domainName": "top.comcast.net"
		},
		{
			"id": 2,
			"name": "cdn2",
			"dnssecEnabled": true,
			"lastUpdated": "2014-10-02 08:22:43",
			"domainName": "2.comcast.net"
		}
	]}

``POST``
========
Allows user to create a CDN

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
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
			"text": "cdn was created.",
			"level": "success"
		}
	],
	"response": {
		"dnssecEnabled": false,
		"domainName": "test.ciab.test",
		"id": 4,
		"lastUpdated": "2018-10-17 15:54:52+00",
		"name": "test"
	}}

