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

.. _to-api-cachegroupparameters:

*********************************
``/api/1.x/cachegroupparameters``
*********************************
Extract information about parameters associated with Cache Groups

``GET``
=======
:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Response Structure
------------------
No available parameters

Response Structure
------------------
:cachegroupParameters: An array of identifying information for parameters assigned to Cache Group profiles

	:parameter:    Numeric ID of the parameter
	:last_updated: Date and time of last modification in ISO format
	:cachegroup:   Name of the Cache Group

.. code-block:: json
	:caption: Response Example

	{ "response": {
		"cachegroupParameters": [
			{
				"parameter": "379",
				"last_updated": "2013-08-05 18:49:37",
				"cachegroup": "us-ca-sanjose"
			},
			{
				"parameter": "380",
				"last_updated": "2013-08-05 18:49:37",
				"cachegroup": "us-ca-sanjose"
			},
			{
				"parameter": "379",
				"last_updated": "2013-08-05 18:49:37",
				"cachegroup": "us-ma-woburn"
			}
		]
	}}
