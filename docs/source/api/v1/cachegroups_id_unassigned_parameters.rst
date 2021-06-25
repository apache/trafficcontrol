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

.. _to-api-v1-cachegroups-id-unassigned_parameters:

********************************************
``cachegroups/{{id}}/unassigned_parameters``
********************************************
.. deprecated:: ATCv4
	This endpoint and all of its functionality is deprecated. All of the information it can return can be obtained using :ref:`to-api-v1-cachegroupparameters` & :ref:`to-api-v1-parameters`.

``GET``
=======
Gets all the :term:`Parameters` that are *not* a specific :ref:`Cache Group's parameters <cache-group-parameters>`.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| Name        | Required | Description                                                                                                   |
	+=============+==========+===============================================================================================================+
	| parameterId | no       | Show only the :term:`Parameter` with the given ID                                                             |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| orderby     | no       | Choose the ordering of the results - must be the name of one of the fields of the objects in the ``response`` |
	|             |          | array                                                                                                         |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| sortOrder   | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")                      |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| limit       | no       | Choose the maximum number of results to return                                                                |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| offset      | no       | The number of results to skip before beginning to return results. Must use in conjunction with limit          |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| page        | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long   |
	|             |          | and the first page is 1. If ``offset`` was defined, this query parameter has no effect. ``limit`` must be     |
	|             |          | defined to make use of ``page``.                                                                              |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| newerThan   | no       | Return only :term:`Parameters` that were most recently updated no earlier than this date/time, which may be   |
	|             |          | given as an :rfc:`3339`-formatted string or as number of nanoseconds since the Unix Epoch (midnight on        |
	|             |          | January 1\ :sup:`st` 1970 UTC).                                                                               |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| olderThan   | no       | Return only :term:`Parameters` that were most recently updated no later than this date/time, which may be     |
	|             |          | given as an :rfc:`3339`-formatted string or as number of nanoseconds since the Unix Epoch (midnight on        |
	|             |          | January 1\ :sup:`st` 1970 UTC).                                                                               |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+

.. versionadded:: ATCv6
	The ``newerThan`` and ``olderThan`` query string parameters were in :abbr:`ATC (Apache Traffic Control)` version 6.0.

.. table:: Request Path Parameters

	+--------+----------+----------------------------------------------------+
	| Name   | Required | Description                                        |
	+========+==========+====================================================+
	| ``id`` | yes      | The :ref:`cache-group-id` of a :term:`Cache Group` |
	+--------+----------+----------------------------------------------------+


Response Structure
------------------
:configFile:  The :term:`Parameter`'s :ref:`parameter-config-file`
:id:          The :term:`Parameter`'s :ref:`parameter-id`
:lastUpdated: The date and time at which this :term:`Parameter` was last updated, in :ref:`non-rfc-datetime`
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
	],
	"alerts": [
		{
			"text": "This endpoint is deprecated, please use GET /cachegroupparameters & GET /parameters instead",
			"level": "warning"
		}
	]}
