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

.. _to-api-v4-async_status:

***********************
``async_status/{{id}}``
***********************

``GET``
=======
Returns a status update for an asynchronous task.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: ASYNC-STATUS:READ
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| Name | Required | Description                                                                                                                          |
	+======+==========+======================================================================================================================================+
	| id   | yes      | The integral, unique identifier for the desired asynchronous job status. This will be provided when the asynchronous job is started. |
	+------+----------+--------------------------------------------------------------------------------------------------------------------------------------+


Response Structure
------------------
:id:         The integral, unique identifier for the asynchronous job status.
:status:     The status of the asynchronous job. This will be `PENDING`, `SUCCEEDED`, or `FAILED`.
:start_time: The time the asynchronous job was started.
:end_time:   The time the asynchronous job completed. This will be `null` if it has not completed yet.
:message:    A message about the job status.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Type: application/json

	{ "response":
		{
			"id":1,
			"status":"PENDING",
			"start_time":"2021-02-18T17:13:56.352261Z",
			"end_time":null,
			"message":"Async job has started."
		}
	}
