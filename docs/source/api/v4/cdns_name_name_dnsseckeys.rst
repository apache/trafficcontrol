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

.. _to-api-v4-cdns-name-name-dnsseckeys:

*********************************
``cdns/name/{{name}}/dnsseckeys``
*********************************

``GET``
=======
Gets a list of DNSSEC keys for CDN and all associated :term:`Delivery Services`. Before returning response to user, this will make sure DNSSEC keys for all :term:`Delivery Services` exist and are not expired. If they don't exist or are expired, they will be (re-)generated.

:Auth. Required: Yes
:Roles Required: "admin"
:Permissions Required: DNS-SEC:READ, CDN:READ, DELIVERY-SERVICE:READ
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+----------------------------------------------------+
	| Name | Description                                        |
	+======+====================================================+
	| name | The name of the CDN for which keys will be fetched |
	+------+----------------------------------------------------+

Response Structure
------------------
:name: The name of the CDN or :term:`Delivery Service` to which the enclosed keys belong

	:zsk: The short-term :abbr:`ZSK (Zone-Signing Key)`

		:expirationDate: A Unix epoch timestamp (in seconds) representing the date and time whereupon the key will expire
		:inceptionDate:  A Unix epoch timestamp (in seconds) representing the date and time when the key was created
		:name:           The name of the domain for which this key will be used
		:private:        Encoded private key
		:public:         Encoded public key
		:ttl:            The time for which the key should be trusted by the client

	:ksk: The long-term :abbr:`KSK (Key-Signing Key)`

		:dsRecord: An optionally present object containing information about the algorithm used to generate the key

			:algorithm: The name of the algorithm used to generate the key
			:digest: A hash of the DNSKEY record
			:digestType: The type of hash algorithm used to create the value of ``digest``

		:expirationDate: A Unix epoch timestamp (in seconds) representing the date and time whereupon the key will expire
		:inceptionDate:  A Unix epoch timestamp (in seconds) representing the date and time when the key was created
		:name:           The name of the domain for which this key will be used
		:private:        Encoded private key
		:public:         Encoded public key
		:ttl:            The time for which the key should be trusted by the client

.. code-block:: json
	:caption: Response Example

	{ "response": {
		"cdn1": {
			"zsk": {
				"ttl": "60",
				"inceptionDate": "1426196750",
				"private": "zsk private key",
				"public": "zsk public key",
				"expirationDate": "1428788750",
				"name": "foo.kabletown.com."
			},
			"ksk": {
				"name": "foo.kabletown.com.",
				"expirationDate": "1457732750",
				"public": "ksk public key",
				"private": "ksk private key",
				"inceptionDate": "1426196750",
				"ttl": "60",
				"dsRecord": {
					"algorithm": "5",
					"digestType": "2",
					"digest": "abc123def456"
				}
			}
		},
		"ds-01": {
			"zsk": {
				"ttl": "60",
				"inceptionDate": "1426196750",
				"private": "zsk private key",
				"public": "zsk public key",
				"expirationDate": "1428788750",
				"name": "ds-01.foo.kabletown.com."
			},
			"ksk": {
				"name": "ds-01.foo.kabletown.com.",
				"expirationDate": "1457732750",
				"public": "ksk public key",
				"private": "ksk private key",
				"inceptionDate": "1426196750"
			}
		}
	}}

``DELETE``
==========
Delete DNSSEC keys for a CDN and all associated :term:`Delivery Services`.

:Auth. Required: Yes
:Roles Required: "admin"
:Permissions Required: DNS-SEC:DELETE, CDN:UPDATE, DELIVERY-SERVICE:UPDATE, CDN:READ
:Response Type:  Object (string)

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-----------------------------------------------------------+
	| Name |                       Description                         |
	+======+===========================================================+
	| name | The name of the CDN for which DNSSEC keys will be deleted |
	+------+-----------------------------------------------------------+

Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{
		"response": "Successfully deleted dnssec keys for test"
	}
