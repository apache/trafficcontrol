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

.. _to-api-cachegroups-id-queue_update:

*******************************
cachegroups/{{ID}}/queue_update
*******************************

``POST``
========
Queue or dequeue updates for all servers assigned to a cache group limited to a specific CDN.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------------+----------+-------------------------------------------------------------------------------------------------+
	| Name            | Required | Description                                                                                     |
	+=================+==========+=================================================================================================+
	| ID              | yes      | The integral, unique identifier for the Cache Group for which updates are being queued/dequeued |
	+-----------------+----------+-------------------------------------------------------------------------------------------------+

.. table:: Request Data Parameters

	+--------------+---------+----------+-----------------------------------------------------------------------------+
	| Name         | Type    | Required | Description                                                                 |
	+==============+=========+==========+=============================================================================+
	| action       | string  | yes      | The action to perform; one of "queue" or "dequeue"                          |
	+--------------+---------+----------+-----------------------------------------------------------------------------+
	| cdn          | string  | no\ [1]_ | The full name of the CDN in need of update queue/dequeue                    |
	+--------------+---------+----------+-----------------------------------------------------------------------------+
	| cdnId        | string  | no\ [1]_ | The integral, unique identifier for the CDN in need of update queue/dequeue |
	+--------------+---------+----------+-----------------------------------------------------------------------------+

.. [1] Either 'cdn' or 'cdnID' *must* be in the request data (but not both).

Response Structure
------------------
:action:         The action processed, one of "queue" or "dequeue"
:cachegroupId:   The integral, unique identifier of the Cache Group for which updates were queued/dequeued
:cachegroupName: The name of the Cache Group for which updates were queued/dequeued
:cdn:            The name of the CDN to which the queue/dequeue operation was restricted
:serverNames:    An array of the (short) hostnames of the servers within the Cache Group which are also assigned to the CDN specified in the ``"cdn"`` field

.. code-block:: json
	:caption: Response Example

	{ "response": {
		"cachegroupName": "CDN_in_a_Box_Edge",
		"action": "dequeue",
		"serverNames": [
			"edge",
			"trafficmonitor",
			"trafficrouter"
		],
		"cdn": "CDN-in-a-Box",
		"cachegroupID": 7
	}}

