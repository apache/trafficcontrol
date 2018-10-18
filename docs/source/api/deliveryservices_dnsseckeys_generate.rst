|

**POST /api/1.1/deliveryservices/dnsseckeys/generate**

	Generates zsk and ksk keypairs for a cdn and all associated delivery services.

	Authentication Required: Yes

	Role(s) Required: Admin

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

	**Request Example** ::

		{
			"key": "cdn1",
			"name" "ott.kabletown.com",
			"ttl": "60",
			"kskExpirationDays": "365",
			"zskExpirationDays": "90"
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

