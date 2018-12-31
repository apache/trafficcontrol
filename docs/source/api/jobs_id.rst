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

.. _to-api-jobs-id:

***************
``jobs/{{ID}}``
***************

``GET``
=======
Get details about a specific job.

:Auth. Required: Yes
:Roles Required: "operations" or "admin"
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------+
	| Name | Description                                                |
	+======+============================================================+
	|  ID  | An integral, unique identifier for the job to be inspected |
	+------+------------------------------------------------------------+

Response Structure
------------------
:assetUrl:        A regular expression - matching URLs will be operated upon according to ``keyword``
:createdBy:       The username of the user who initiated the job
:deliveryService: The 'xml_id' that uniquely identifies the :term:`Delivery Service` on which this job operates
:id:              An integral, unique identifier for this job
:keyword:         A keyword that represents the operation being performed by the job:

	PURGE
		This job will prevent caching of URLs matching the ``assetURL`` until it is removed (or its Time to Live expires)

:parameters: A string containing key/value pairs representing parameters associated with the job - currently only uses Time to Live e.g. ``"TTL:48h"``
:startTime:  The date and time at which the job began, in ISO format

.. code-block:: json
	:caption: Response Example

	{ "response": [
		{
			"id": 1
			"assetUrl": "http:\/\/foo-bar.domain.net\/taco.html",
			"deliveryService": "foo-bar",
			"keyword": "PURGE",
			"parameters": "TTL:48h",
			"startTime": "2015-05-14 08:56:36-06",
			"createdBy": "jdog24"
		}
	 ]}
