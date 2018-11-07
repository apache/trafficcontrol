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

.. _to-cdns-id-queue_update:

****************************
``cdns/{{id}}/queue_update``
****************************

``POST``
========
Queue or dequeue updates for all servers assigned to a specific CDN.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------------+----------+---------------------------------------------------------------------------+
	| Name            | Required | Description                                                               |
	+=================+==========+===========================================================================+
	| id              | yes      | The integral, unique identifier for the CDN on which to (de)queue updates |
	+-----------------+----------+---------------------------------------------------------------------------+

.. table:: Request Data Parameters

	+--------------+---------+----------+-----------------------------------------------+
	| Name         | Type    | Required | Description                                   |
	+==============+=========+==========+===============================================+
	| action       | string  | yes      | One of "queue" or "dequeue" as appropriate    |
	+--------------+---------+----------+-----------------------------------------------+

Response Structure
------------------
:action: The action processed, either ``"queue"`` or ``"dequeue"``
:cdnId:  The integral, unique identifier for the CDN on which updates were (de)queued

.. code-block:: json
	:caption: Response Example

	{ "response": {
			"action": "queue",
			"cdn": 1
		}
	}
