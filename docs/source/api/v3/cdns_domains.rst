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

.. _to-api-v3-cdns-domains:

****************
``cdns/domains``
****************

``GET``
=======
Gets a list of domains and their related Traffic Router :term:`Profiles` for all CDNs.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
No parameters available.

Response Structure
------------------
:domainName:         The :abbr:`TLD (Top-Level Domain)` assigned to this CDN
:parameterId:        The :ref:`parameter-id` for the :term:`Parameter` that sets this :abbr:`TLD (Top-Level Domain)` on the Traffic Router
:profileDescription: A short, human-readable description of the Traffic Router's profile
:profileId:          The :ref:`profile-id` of the :term:`Profile` assigned to the Traffic Router responsible for serving ``domainName``
:profileName:        The :ref:`profile-name` of the :term:`Profile` assigned to the Traffic Router responsible for serving ``domainName``

.. code-block:: json
	:caption: Response Example

	{ "response": [
		{
			"profileId": 12,
			"parameterId": -1,
			"profileName": "CCR_CIAB",
			"profileDescription": "Traffic Router for CDN-In-A-Box",
			"domainName": "mycdn.ciab.test"
		}
	]}
