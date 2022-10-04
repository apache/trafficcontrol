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

.. _to-api-v4-cdns-capacity:

*****************
``cdns/capacity``
*****************

``GET``
=======
Retrieves the aggregate capacity percentages of all :term:`Cache Groups` for a given CDN.

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: CDN:READ
:Response Type:  Object

Request Structure
-----------------
No parameters available.

Response Structure
------------------
:availablePercent:   The percent of available (unused) bandwidth to 64 bits of precision\ [1]_
:unavailablePercent: The percent of unavailable (used) bandwidth to 64 bits of precision\ [1]_
:utilizedPercent:    The percent of bandwidth currently in use to 64 bits of precision\ [1]_
:maintenancePercent: The percent of bandwidth being used for administrative or analytical processes internal to the CDN to 64 bits of precision\ [1]_

.. code-block:: json
	:caption: Response Example

	{ "response": {
		"availablePercent": 89.0939840205533,
		"unavailablePercent": 0,
		"utilizedPercent": 10.9060020300395,
		"maintenancePercent": 0.0000139494071146245
	}}

.. [1] Following `IEEE 754 <https://ieeexplore.ieee.org/document/4610935>`_
