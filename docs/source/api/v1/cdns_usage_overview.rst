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

.. _to-api-v1-cdns-usage-overview:

***********************
``cdns/usage/overview``
***********************

.. versionadded:: 1.2

.. deprecated:: ATCv4
	This endpoint and its functionality is deprecated, and will be removed in the future.

``GET``
=======
Retrieves the high-level CDN usage metrics from Traffic Stats

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
No parameters available.

Response Structure
------------------
:currentGbps: The current throughput of all CDNs, in Gigabits per second
:maxGbps:     The all-time maximum throughput of all CDNs, in Gigabits per second
:source:      The name of the service providing the statistics. This will almost always be "TrafficStats"
:tps:         The number of transactions being performed per second
:version:     The version of the service providing the statistics (named in ``"source"``)

.. warning:: The ``"tps"`` field is currently broken, and will return ``0`` every time. See `GitHub issue #1020 <https://github.com/apache/trafficcontrol/issues/1020>`_ for more information.

.. code-block:: json
	:caption: Response Example

	{
		"alerts": [
			{
				"level": "warning",
				"text": "This endpoint and its functionality is deprecated, and will be removed in the future"
			}
		],
		"response": {
			"currentGbps": 975.920621333333,
			"source": "TrafficStats",
			"tps": 0,
			"version": "1.2",
			"maxGbps": 12085
		}
	}
