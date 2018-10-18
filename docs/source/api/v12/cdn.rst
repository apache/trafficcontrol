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

**POST /api/1.2/cdns/{:id}/queue_update**

	Queue or dequeue updates for all servers assigned to a specific CDN.

	Authentication Required: Yes

	Role(s) Required: admin or oper

	**Request Route Parameters**

	+-----------------+----------+----------------------+
	| Name            | Required | Description          |
	+=================+==========+======================+
	| id              | yes      | the cdn id.          |
	+-----------------+----------+----------------------+

	**Request Properties**

	+--------------+---------+-----------------------------------------------+
	| Name         | Type    | Description                                   |
	+==============+=========+===============================================+
	| action       | string  | queue or dequeue                              |
	+--------------+---------+-----------------------------------------------+

	**Request Example** ::

		{
				"action": "queue"
		}


	**Response Properties**

	+-----------------+---------+----------------------------------------------------+
	| Name            | Type    | Description                                        |
	+=================+=========+====================================================+
	| action          | string  | The action processed, queue or dequeue.            |
	+-----------------+---------+----------------------------------------------------+
	| cdnId           | integer | cdn id                                             |
	+-----------------+---------+----------------------------------------------------+

	**Response Example** ::

		{
			"response": {
						"action": "queue",
						"cdn": 1
				}
		}


.. _to-api-v12-cdn-dnsseckeys:

DNSSEC Keys
+++++++++++

**GET /api/1.2/cdns/name/:name/dnsseckeys**

	Gets a list of dnsseckeys for a CDN and all associated Delivery Services.

	Authentication Required: Yes

	Role(s) Required: Admin

	**Request Route Parameters**

	+----------+----------+-------------+
	|   Name   | Required | Description |
	+==========+==========+=============+
	| ``name`` | yes      |             |
	+----------+----------+-------------+

	**Response Properties**

	+-------------------------------+--------+---------------------------------------------------------------+
	|           Parameter           |  Type  |                          Description                          |
	+===============================+========+===============================================================+
	| ``cdn name/ds xml_id``        | string | identifier for ds or cdn                                      |
	+-------------------------------+--------+---------------------------------------------------------------+
	| ``>zsk/ksk``                  | array  | collection of zsk/ksk data                                    |
	+-------------------------------+--------+---------------------------------------------------------------+
	| ``>>ttl``                     | string | time-to-live for dnssec requests                              |
	+-------------------------------+--------+---------------------------------------------------------------+
	| ``>>inceptionDate``           | string | epoch timestamp for when the keys were created                |
	+-------------------------------+--------+---------------------------------------------------------------+
	| ``>>expirationDate``          | string | epoch timestamp representing the expiration of the keys       |
	+-------------------------------+--------+---------------------------------------------------------------+
	| ``>>private``                 | string | encoded private key                                           |
	+-------------------------------+--------+---------------------------------------------------------------+
	| ``>>public``                  | string | encoded public key                                            |
	+-------------------------------+--------+---------------------------------------------------------------+
	| ``>>name``                    | string | domain name                                                   |
	+-------------------------------+--------+---------------------------------------------------------------+
	| ``version``                   | string | API version                                                   |
	+-------------------------------+--------+---------------------------------------------------------------+
	| ``ksk>>dsRecord>>algorithm``  | string | The algorithm of the referenced DNSKEY-recor.                 |
	+-------------------------------+--------+---------------------------------------------------------------+
	| ``ksk>>dsRecord>>digestType`` | string | Cryptographic hash algorithm used to create the Digest value. |
	+-------------------------------+--------+---------------------------------------------------------------+
	| ``ksk>>dsRecord>>digest``     | string | A cryptographic hash value of the referenced DNSKEY-record.   |
	+-------------------------------+--------+---------------------------------------------------------------+

	**Response Example** ::

		{
			"response": {
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
						dsRecord: {
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
				},
				... repeated for each ds in the cdn
			},
		}


|

**GET /api/1.2/cdns/name/:name/dnsseckeys/delete**

	Delete dnssec keys for a cdn and all associated delivery services.

	Authentication Required: Yes

	Role(s) Required: Admin

	**Request Route Parameters**

	+----------+----------+----------------------------------------------------------+
	|   Name   | Required |                       Description                        |
	+==========+==========+==========================================================+
	| ``name`` | yes      | name of the CDN for which you want to delete dnssec keys |
	+----------+----------+----------------------------------------------------------+

	**Response Properties**

	+--------------+--------+------------------+
	|  Parameter   |  Type  |   Description    |
	+==============+========+==================+
	| ``response`` | string | success response |
	+--------------+--------+------------------+

	**Response Example**
	::

		{
			"response": "Successfully deleted dnssec keys for <cdn>"
		}

|

**POST /api/1.2/deliveryservices/dnsseckeys/generate**

	Generates ZSK and KSK keypairs for a CDN and all associated Delivery Services.

	Authentication Required: Yes

	Role(s) Required:  Admin

	**Request Properties**

	+-----------------------+---------+------------------------------------------------+
	|       Parameter       |   Type  |                  Description                   |
	+=======================+=========+================================================+
	| ``key``               | string  | name of the cdn                                |
	+-----------------------+---------+------------------------------------------------+
	| ``name``              | string  | domain name of the cdn                         |
	+-----------------------+---------+------------------------------------------------+
	| ``ttl``               | string  | time to live                                   |
	+-----------------------+---------+------------------------------------------------+
	| ``kskExpirationDays`` | string  | Expiration (in days) for the key signing keys  |
	+-----------------------+---------+------------------------------------------------+
	| ``zskExpirationDays`` | string  | Expiration (in days) for the zone signing keys |
	+-----------------------+---------+------------------------------------------------+
	| ``effectiveDate``     | int     | UNIX epoch start date for the signing keys     |
	+-----------------------+---------+------------------------------------------------+

	**Request Example** ::

		{
			"key": "cdn1",
			"name" "ott.kabletown.com",
			"ttl": "60",
			"kskExpirationDays": "365",
			"zskExpirationDays": "90",
			"effectiveDate": 1012615322
		}

	**Response Properties**

	+--------------+--------+-----------------+
	|  Parameter   |  Type  |   Description   |
	+==============+========+=================+
	| ``response`` | string | response string |
	+--------------+--------+-----------------+
	| ``version``  | string | API version     |
	+--------------+--------+-----------------+

	**Response Example** ::


		{
			"response": "Successfully created dnssec keys for cdn1"
		}

.. _to-api-v12-cdn-sslkeys:

SSL Keys
+++++++++++

**GET /api/1.2/cdns/name/:name/sslkeys**

	Returns ssl certificates for all Delivery Services that are a part of the CDN.

	Authentication Required: Yes

	Role(s) Required: Admin

	**Request Route Parameters**

	+----------+----------+-------------+
	|   Name   | Required | Description |
	+==========+==========+=============+
	| ``name`` | yes      |             |
	+----------+----------+-------------+

	**Response Properties**

	+-------------------------------+--------+---------------------------------------------------------------+
	|           Parameter           |  Type  |                          Description                          |
	+===============================+========+===============================================================+
	| ``deliveryservice``           | string | identifier for deliveryservice xml_id                         |
	+-------------------------------+--------+---------------------------------------------------------------+
	| ``certificate``               | array  | collection of certificate                                     |
	+-------------------------------+--------+---------------------------------------------------------------+
	| ``>>key``                     | string | base64 encoded private key for ssl certificate                |
	+-------------------------------+--------+---------------------------------------------------------------+
	| ``>>crt``                     | string | base64 encoded ssl certificate                                |
	+-------------------------------+--------+---------------------------------------------------------------+


	**Response Example** ::

		{
			"response": [
				{
					"deliveryservice": "ds1",
					"certificate": {
						"crt": "base64encodedcrt1",
						"key": "base64encodedkey1"
					}
				},
				{
					"deliveryservice": "ds2",
					"certificate": {
						"crt": "base64encodedcrt2",
						"key": "base64encodedkey2"
					}
				}
			]
		}
