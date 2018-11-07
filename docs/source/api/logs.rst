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
.. _to-api-logs:

********
``logs``
********

``GET``
=======
Fetches a list of changes that have been made to the Traffic Control system

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------------+---------+----------+---------------------------------------------------+
	| Name            | Type    | Required | Description                                       |
	+=================+=========+==========+===================================================+
	| ``days``        | integer | no       | The number of days of change logs to return       |
	+-----------------+---------+----------+---------------------------------------------------+
	| ``limit``       | integer | no       | The number of rows to which to limit the response |
	+-----------------+---------+----------+---------------------------------------------------+

Response Structure
------------------
:id:          Integral, unique identifier for the Log entry
:lastUpdated: Date and time at which the change was made, in ISO format
:level:       Log categories for each entry, e.g. 'UICHANGE', 'OPER', 'APICHANGE'
:message:     Log detail about what occurred
:ticketNum:   Optional field to cross reference with any bug tracking systems
:user:        Name of the user who made the change

.. code-block:: json
	:caption: Response Example

	{ "response": [
		{
			"ticketNum": null,
			"level": "APICHANGE",
			"lastUpdated": "2018-10-24 16:07:17.45204+00",
			"user": "admin",
			"id": 430,
			"message": "Server updates queued for cdn 2"
		},
		{
			"ticketNum": null,
			"level": "APICHANGE",
			"lastUpdated": "2018-10-24 16:07:16.130401+00",
			"user": "admin",
			"id": 429,
			"message": "Snapshot of CRConfig performed for CDN-in-a-Box"
		}
	]}
