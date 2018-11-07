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

.. _to-api-v12-cachegroupfallbacks-route:

************************
``cachegroup_fallbacks``
************************

``GET``
=======
Retrieve fallback-related configurations for a Cache Group.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters**

	+--------------+----------+---------------------------------------------------------------------------------------------------+
	| Name         | Required | Description                                                                                       |
	+==============+==========+===================================================================================================+
	| cacheGroupId |yes\ [1]_ | The integral, unique identifier of a Cache Group whose fallback configurations shall be retrieved |
	+--------------+----------+---------------------------------------------------------------------------------------------------+
	| fallbackId   |yes\ [1]_ | The integral, unique identifier of a fallback Cache Group                                         |
	+--------------+----------+---------------------------------------------------------------------------------------------------+

.. [1] At least one of these must be provided, not necessarily both (though both is perfectly valid)

Response Structure
------------------
:cacheGroupId:   The integral, unique identifier of the Cache Group described by this entry
:cacheGroupName: The name of the Cache Group described by this entry
:fallbackId:     The integral, unique identifier of the Cache Group on which the Cache Group described by this entry will fall back
:fallbackName:   The name of the Cache Group on which the Cache Group described by this entry will fall back
:fallbackOrder:  The order of the fallback described by "fallbackId" and "fallbackName" in the list of fallbacks for the Cache Group described by this entry

.. code-block:: json
	:caption: Response Example

	{ "response": [
		{
			"cacheGroupId":1,
			"cacheGroupName":"GROUP1",
			"fallbackId":2,
			"fallbackOrder":10,
			"fallbackName":"GROUP2"
		}
	]}

``POST``
========
Creates fallback configuration for a Cache Group.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Data Parameters

	+----------------------------------+---------+----------+---------------------------------------------------------------------------------------------------------------------+
	| Parameter                        | Type    | Required | Description                                                                                                         |
	+==================================+=========+==========+=====================================================================================================================+
	| ``cacheGroupId``                 | integer | yes      | Integral, unique identifier of a Cache Group to which to assign a fallback                                          |
	+----------------------------------+---------+----------+---------------------------------------------------------------------------------------------------------------------+
	| ``fallbackId``                   | integer | yes      | Integral, unique identifier of a Cache Group on which the Cache Group identified by ``cacheGroupId`` will fall back |
	+----------------------------------+---------+----------+---------------------------------------------------------------------------------------------------------------------+
	| ``fallbackOrder``                | integer | yes      | The order of this fallback for the Cache Group identified by ``cacheGroupId``                                       |
	+----------------------------------+---------+----------+---------------------------------------------------------------------------------------------------------------------+

.. note:: The request data should be an array of these objects (and any number can be submitted per request), see the example

.. code-block:: json
	:caption: Request Example

	[
		{
			"cacheGroupId": 1,
			"fallbackId": 3,
			"fallbackOrder": 10
		 }
	]

Response Structure
------------------
:cacheGroupId:   The integral, unique identifier of the Cache Group to which this fallback was assigned
:cacheGroupName: The name of the Cache Group to which this fallback was assigned
:fallbackId:     The integral, unique identifier of the Cache Group on which this entries Cache Group will fall back
:fallbackName:   The name of the Cache Group on which this entries Cache Group will fall back
:fallbackOrder:  The order of the fallback described by "fallbackId" and "fallbackName" in the list of fallbacks for the Cache Group described by this entry


.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"level":"success",
			"text":"Backup configuration CREATE for cache group 1 successful."
		}
	],
	"response": [
		{
			"cacheGroupId":1,
			"cacheGroupName":"GROUP1",
			"fallbackId":3,
			"fallbackName":"GROUP2",
			"fallbackorder":10,
		}
	]}

``PUT``
=======
Updates an existing fallback configuration for one or more Cache Groups.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Data Parameters

	+----------------------------------+---------+----------+---------------------------------------------------------------------------------------------------------------------+
	| Parameter                        | Type    | Required | Description                                                                                                         |
	+==================================+=========+==========+=====================================================================================================================+
	| ``cacheGroupId``                 | integer | yes      | Integral, unique identifier of a Cache Group to which to assign a fallback                                          |
	+----------------------------------+---------+----------+---------------------------------------------------------------------------------------------------------------------+
	| ``fallbackId``                   | integer | yes      | Integral, unique identifier of a Cache Group on which the Cache Group identified by ``cacheGroupId`` will fall back |
	+----------------------------------+---------+----------+---------------------------------------------------------------------------------------------------------------------+
	| ``fallbackOrder``                | integer | yes      | The order of this fallback for the Cache Group identified by ``cacheGroupId``                                       |
	+----------------------------------+---------+----------+---------------------------------------------------------------------------------------------------------------------+

.. note:: The request data should be an array of these objects (and any number can be submitted per request), see the example

.. code-block:: json
	:caption: Request Example

		[
			 {
					"cacheGroupId": 1,
					"fallbackId": 3,
					"fallbackOrder": 10
			 }
		]

Response Structure
------------------
:cacheGroupId:   The integral, unique identifier of the Cache Group to which this fallback was assigned
:cacheGroupName: The name of the Cache Group to which this fallback was assigned
:fallbackId:     The integral, unique identifier of the Cache Group on which this entries Cache Group will fall back
:fallbackName:   The name of the Cache Group on which this entries Cache Group will fall back
:fallbackOrder:  The order of the fallback described by "fallbackId" and "fallbackName" in the list of fallbacks for the Cache Group described by this entry

.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"level":"success",
			"text":"Backup configuration UPDATE for cache group 1 successful."
		}
	],
	"response": [
		{
			"cacheGroupId":1,
			"cacheGroupName":"GROUP1",
			"fallbackId":3,
			"fallbackName":"GROUP2",
			"fallbackorder":10,
		}
	]}

``DELETE``
==========
Delete fallback list assigned to a Cache Group

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Query Parameters**

	+--------------+----------+---------------------------------------------------------------------------------------------------+
	| Name         | Required | Description                                                                                       |
	+==============+==========+===================================================================================================+
	| cacheGroupId |yes\ [2]_ | The integral, unique identifier of a Cache Group whose fallback configurations shall be retrieved |
	+--------------+----------+---------------------------------------------------------------------------------------------------+
	| fallbackId   |yes\ [2]_ | The integral, unique identifier of a fallback Cache Group                                         |
	+--------------+----------+---------------------------------------------------------------------------------------------------+

.. [2] At least one of "cacheGroupId" or "fallbackId" must be sent with the request. If both are sent, a single fallback relationship is deleted, whereas using only "cacheGroupId" will result in all fallbacks being removed from the Cache Group identified by that integral, unique identifier, and using only "fallbackId" will remove the Cache Group identified by *that* integral, unique identifier from all other Cache Groups' fallback lists.

Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"level": "success",
			"text": "Backup configuration DELETED"
		}
	]}

