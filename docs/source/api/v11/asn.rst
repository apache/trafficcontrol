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


.. _to-api-v11-asn:

****
ASNS
****

.. _to-api-v11-asns-route:

``/api/1.1/asns``
=================
List all ASNS

``GET``
-------
:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

.. table:: Response Properties

	+------------------+---------------------------+----------------------+
	|        Key       |             Type          |     Description      |
	+==================+===========================+======================+
	| ``asns``         | array of ``asn`` objects  | A collection of asns |
	+------------------+---------------------------+----------------------+

.. table:: ``asn`` Properties

	+------------------+---------+-------------------------------------------------------------------------+
	|       Key        |  Type   |                           Description                                   |
	+==================+=========+=========================================================================+
	| ``lastUpdated``  | string  | The Time / Date this server entry was last updated                      |
	+------------------+---------+-------------------------------------------------------------------------+
	| ``id``           | integer | Local unique identifier for the ASN                                     |
	+------------------+---------+-------------------------------------------------------------------------+
	| ``asn``          | integer | Autonomous System Numbers per APNIC for identifying a service provider  |
	+------------------+---------+-------------------------------------------------------------------------+
	| ``cachegroup``   | string  | Related cachegroup name                                                 |
	+------------------+---------+-------------------------------------------------------------------------+
	| ``cachegroupId`` | integer | Related cachegroup name                                                 |
	+------------------+---------+-------------------------------------------------------------------------+

.. code-block:: json
	:caption: Response Example

	{
	 "response": {
			"asns": [
				 {
						"lastUpdated": "2012-09-17 21:41:22",
						"id": 27,
						"asn": 7015,
						"cachegroup": "us-ma-woburn",
						"cachegroupId": 2
				 },
				 {
						"lastUpdated": "2012-09-17 21:41:22",
						"id": 28,
						"asn": 7016,
						"cachegroup": "us-pa-pittsburgh",
						"cachegroupID": 3
				 }
			]
	 },
	}

