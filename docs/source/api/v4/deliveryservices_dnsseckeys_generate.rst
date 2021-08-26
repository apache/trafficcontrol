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

.. _to-api-v4-deliveryservices-dnsseckeys-generate:

****************************************
``deliveryservices/dnsseckeys/generate``
****************************************

``POST``
========
Generates :abbr:`ZSK (Zone-Signing Key)` and :abbr:`KSK (Key-Signing Key)` keypairs for a CDN and all associated :term:`Delivery Services`.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object (string)

Request Structure
-----------------
:effectiveDate: UNIX epoch start date for the signing keys
:key:               Name of the CDN
:kskExpirationDays: Expiration (in days) for the :abbr:`KSKs (Key-Signing Keys)`
:name:              Domain name used by the CDN
:ttl:               Time for which the keypairs shall remain valid
:zskExpirationDays: Expiration (in days) for the :abbr:`ZSKs (Zone-Signing Keys)`


Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{
		"response": "Successfully created dnssec keys for cdn1"
	}

