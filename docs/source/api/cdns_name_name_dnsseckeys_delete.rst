|

**GET /api/1.1/cdns/name/:name/dnsseckeys/delete**

	Delete dnssec keys for a cdn and all associated delivery services.

	Authentication Required: Yes

	Role(s) Required: Admin

	**Request Path Parameters**

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

