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

.. _to-api-deliveryservices-dnsseckeys-generate:

*************************************************
``/api/1.x/deliveryservices/dnsseckeys/generate``
*************************************************

``POST``
========
Generates Zone-Signing Key (ZSK) and Key-Signing Key (KSK) keypairs for a CDN and all associated Delivery Services.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object (string)

Request Structure
-----------------
.. table:: Request Properties

	+-----------------------+---------+----------+------------------------------------------------+
	|       Parameter       |   Type  | Required |                  Description                   |
	+=======================+=========+==========+================================================+
	| ``key``               | string  | yes      | Name of the CDN                                |
	+-----------------------+---------+----------+------------------------------------------------+
	| ``name``              | string  | yes      | Domain name used by the CDN                    |
	+-----------------------+---------+----------+------------------------------------------------+
	| ``ttl``               | string  | yes      | Time for which the keypairs shall remain valid |
	+-----------------------+---------+----------+------------------------------------------------+
	| ``kskExpirationDays`` | string  | yes      | Expiration (in days) for the KSKs              |
	+-----------------------+---------+----------+------------------------------------------------+
	| ``zskExpirationDays`` | string  | yes      | Expiration (in days) for the ZSKs              |
	+-----------------------+---------+----------+------------------------------------------------+
	| ``effectiveDate``     | integer | yes      | UNIX epoch start date for the signing keys     |
	+-----------------------+---------+----------+------------------------------------------------+

.. versionchanged:: 1.2
	Added required 'effectiveDate' field to request

Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{
		"response": "Successfully created dnssec keys for cdn1"
	}

