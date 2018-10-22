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

.. _to-api-cdns-routing:

*************************
``/api/1.x/cdns/routing``
*************************

``GET``
=======
Retrieves the aggregate routing percentages of Cache Groups assigned to any CDN.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
No parameters available

Response Structure
------------------
:cz:          Used Coverage Zone geographic IP mapping
:dsr:         Overflow traffic sent to secondary CDN
:err:         Error localizing client IP
:geo:         Used 3rd party geographic IP mapping
:miss:        No location available for client IP
:staticRoute: Used pre-configured DNS entries

.. code-block:: json
	:caption: Response Example

	{ "response": {
		"staticRoute": 0,
		"geo": 20.6251834458468,
		"err": 0,
		"fed": 0.287643087760493,
		"cz": 79.0607572644555,
		"regionalAlternate": 0,
		"dsr": 0,
		"miss": 0.0264162019371881,
		"regionalDenied": 0
	}}

