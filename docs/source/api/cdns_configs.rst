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
.. _to-api-cdns-config:

****************
``cdns/configs``
****************
.. deprecated:: 1.0
	Use one of :ref:`to-api-cdns-name-configs-monitoring`, :ref:`to-api-cdns-name-configs-routing`, or :ref:`to-api-servers-server-configfiles-ats` instead.

.. caution:: This endpoint doesn't appear to work as of Traffic Control version 3.0.0 - it is strongly advised that its used be avoided.

``GET``
=======
Retrieves CDN configuration information.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
No parameters available.

Response Properties
-------------------
:config_file: Presumably the name of some configuration file\ [1]_
:id:          The integral, unique identifier for this CDN
:name:        The CDN's name
:value:       Presumably the content of some configuration file\ [1]_

.. [1] These values are currently missing from this endpoint's output. **DO NOT count on this endpoint to provide this information**.
