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

.. _to-api-servercheck:

***************
``servercheck``
***************

Updates the resulting value from running a given check extension on a server.

``POST``
========
Post a server check result to the "serverchecks" table.

:Auth. Required: Yes
:Roles Required: None\ [1]_
:Response Type: Object

Request Structure
-----------------
The request only requires to have either ``host_name`` or ``id`` defined.

:host_name:              The hostname of the server to which this "servercheck" refers.
:id:                     The id of the server to which this "servercheck" refers.
:servercheck_short_name: The short name of the "servercheck".
:value:                  The value of the "servercheck"

.. code-block:: http
	:caption: Request Example

	POST /api/1.1/servercheck HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 113
	Content-Type: application/json

	{
		"id": 1,
		"host_name": "edge",
		"servercheck_short_name": "test",
		"value": 1
	}

Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"level": "success",
			"text": "Server Check was successfully updated."
		}
	]}

.. [1] No roles are required to use this endpoint, however access is controlled by username. Only the reserved user ``extension`` is permitted the use of this endpoint.

