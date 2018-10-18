DNSSEC Keys
+++++++++++

**GET /api/1.1/cdns/name/:name/dnsseckeys**

	Gets a list of dnsseckeys for CDN and all associated Delivery Services.
	Before returning response to user, check to make sure keys aren't expired.  If they are expired, generate new ones.
	Before returning response to user, make sure dnssec keys for all delivery services exist.  If they don't exist, create them.

	Authentication Required: Yes

	Role(s) Required: Admin

	**Request Path Parameters**

	+----------+----------+-------------+
	|   Name   | Required | Description |
	+==========+==========+=============+
	| ``name`` | yes      |             |
	+----------+----------+-------------+

	**Response Properties**

	+------------------------+--------+---------------------------------------------------------+
	|       Parameter        |  Type  |                       Description                       |
	+========================+========+=========================================================+
	| ``cdn name/ds xml_id`` | string | identifier for ds or cdn                                |
	+------------------------+--------+---------------------------------------------------------+
	| ``>zsk/ksk``           | array  | collection of zsk/ksk data                              |
	+------------------------+--------+---------------------------------------------------------+
	| ``>>ttl``              | string | time-to-live for dnssec requests                        |
	+------------------------+--------+---------------------------------------------------------+
	| ``>>inceptionDate``    | string | epoch timestamp for when the keys were created          |
	+------------------------+--------+---------------------------------------------------------+
	| ``>>expirationDate``   | string | epoch timestamp representing the expiration of the keys |
	+------------------------+--------+---------------------------------------------------------+
	| ``>>private``          | string | encoded private key                                     |
	+------------------------+--------+---------------------------------------------------------+
	| ``>>public``           | string | encoded public key                                      |
	+------------------------+--------+---------------------------------------------------------+
	| ``>>name``             | string | domain name                                             |
	+------------------------+--------+---------------------------------------------------------+
	| ``version``            | string | API version                                             |
	+------------------------+--------+---------------------------------------------------------+


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
						"ttl": "60"
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

