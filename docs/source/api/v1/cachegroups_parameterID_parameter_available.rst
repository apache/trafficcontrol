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

.. _to-api-v1-cachegroups-parameterID-parameter-available:

****************************************************
``cachegroups/{{parameter ID}}/parameter/available``
****************************************************
.. deprecated:: ATCv4
.. danger:: This endpoint does not appear to work, and thus its use is strongly discouraged!

``GET``
=======
Gets a list of :term:`Cache Groups` which are available to have a specific :term:`Parameter` assigned to them.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------------------+----------+--------------------------------------------------------------+
	| Name             | Required | Description                                                  |
	+==================+==========+==============================================================+
	| ``parameter ID`` | yes      | The :ref:`parameter-id` of the :term:`Parameter` of interest |
	+------------------+----------+--------------------------------------------------------------+

Response Structure
------------------
:id:   An integer that is the :ref:`Cache Group's ID <cache-group-id>`
:name: A string that is the :ref:`Cache Group's name <cache-group-name>`

.. code-block:: json
	:caption: Response Example

	{
		"alerts": [
			{
				"level": "warning",
				"text": "This endpoint is deprecated, please use 'GET /cachegroupparameters & GET /cachegroups' instead"
			}
		],
		"response": [
			{
				"name": "dc-chicago",
				"id": "21"
			},
			{
				"name": "dc-cmc",
				"id": "22"
			}
	]}
