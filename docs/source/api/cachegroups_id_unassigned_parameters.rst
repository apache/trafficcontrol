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

.. _to-api-cachegroups-id-unassigned_parameters:

********************************************
``cachegroups/{{id}}/unassigned_parameters``
********************************************

``GET``
=======
Gets all the :term:`Parameters` *not* associated with a specific :term:`Cache Group`

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------------------+----------+---------------------------------------------------------+
	|       Name       | Required | Description                                             |
	+==================+==========+=========================================================+
	| ``id``           | yes      | An integral, unique identifier of a :term:`Cache Group` |
	+------------------+----------+---------------------------------------------------------+


Response Structure
------------------
:configFile:  The :term:`Parameter`'s :ref:`parameter-config-file`
:id:          The :term:`Parameter`'s :ref:`parameter-id`
:lastUpdated: The date and time at which this :term:`Parameter` was last updated, in an ISO-like format
:name:        :ref:`parameter-name` of the :term:`Parameter`
:secure:      A boolean value that describes whether or not the :term:`Parameter` is :ref:`parameter-secure`
:value:       The :term:`Parameter`'s :ref:`parameter-value`

.. code-block:: json
	:caption: Response Example

	{ "response": [
		{
			"lastUpdated": "2018-10-09 11:14:33.862905+00",
			"value": "/opt/trafficserver/etc/trafficserver",
			"secure": false,
			"name": "location",
			"id": 6836,
			"configFile": "hdr_rw_bamtech-nhl-live.config"
		},
		{
			"lastUpdated": "2018-10-09 11:14:33.862905+00",
			"value": "/opt/trafficserver/etc/trafficserver",
			"secure": false,
			"name": "location",
			"id": 6837,
			"configFile": "hdr_rw_mid_bamtech-nhl-live.config"
		},
		{
			"lastUpdated": "2018-10-09 11:55:46.014844+00",
			"value": "/opt/trafficserver/etc/trafficserver",
			"secure": false,
			"name": "location",
			"id": 6842,
			"configFile": "hdr_rw_bamtech-nhl-live-t.config"
		},
		{
			"lastUpdated": "2018-10-09 11:55:46.014844+00",
			"value": "/opt/trafficserver/etc/trafficserver",
			"secure": false,
			"name": "location",
			"id": 6843,
			"configFile": "hdr_rw_mid_bamtech-nhl-live-t.config"
		}
	]}
