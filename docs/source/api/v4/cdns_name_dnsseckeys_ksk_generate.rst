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

.. _to-api-v4-cdns-name-dnsseckeys-ksk-generate:

*****************************************
``cdns/{{name}}/dnsseckeys/ksk/generate``
*****************************************

``POST``
========
Generates a new :abbr:`KSK (Key-Signing Key)` for a specific CDN.

:Auth. Required: Yes
:Roles Required: "admin"
:Permissions Required: DNS-SEC:CREATE, CDN:UPDATE, CDN:READ
:Response Type:  Object (string)

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+----------+-----------------------------------------------------------------------------------+
	| Name | Required | Description                                                                       |
	+======+==========+===================================================================================+
	| name | yes      | The name of the CDN for which the :abbr:`KSK (Key-Signing Key)` will be generated |
	+------+----------+-----------------------------------------------------------------------------------+

:expirationDays: The integral number of days until the newly generated :abbr:`KSK (Key-Signing Key)` expires
:effectiveDate:  An optional string containing the date and time at which the newly generated :abbr:`KSK (Key-Signing Key)` becomes effective, in :RFC:`3339` format. Defaults to the current time if not specified

Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{ "response": "Successfully generated ksk dnssec keys for my-cdn-name" }
