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

.. _to-api-deliveryservices-xmlid-urisignkeys:

*******************************************
``deliveryservices/{{xml_id}}/urisignkeys``
*******************************************

DELETE

	Deletes URISigning objects for a delivery service.

	Authentication Required: Yes

	Role(s) Required: admin

	**Request Route Parameters**

	+-----------+----------+----------------------------------------+
	|    Name   | Required |              Description               |
	+===========+==========+========================================+
	| xml_id    | yes      | xml_id of the desired delivery service |
	+-----------+----------+----------------------------------------+

**GET deliveryservices/:xml_id/urisignkeys**

	Retrieves one or more URISigning objects for a delivery service.

	Authentication Required: Yes

	Role(s) Required: admin

	**Request Route Parameters**

	+-----------+----------+----------------------------------------+
	|    Name   | Required |              Description               |
	+===========+==========+========================================+
	| xml_id    | yes      | xml_id of the desired delivery service |
	+-----------+----------+----------------------------------------+

	**Response Properties**

	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	|    Parameter        |  Type  |                                                               Description                                                               |
	+=====================+========+=========================================================================================================================================+
	| ``Issuer``          | string | a string describing the issuer of the URI signing object. Multiple URISigning objects may be returned in a response, see example        |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``renewal_kid``     | string | a string naming the jwt key used for renewals.                                                                                          |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``keys``            | string | json array of jwt symmetric keys                                                             .                                          |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``alg``             | string | this parameter repeats for each jwt key in the array and specifies the jwa encryption algorithm to use with this key, :rfc:`7518`       |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``kid``             | string | this parameter repeats for each jwt key in the array and specifies the unique id for the key as defined in :rfc:`7516`                  |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``kty``             | string | this parameter repeats for each jwt key in the array and specifies the key type as defined in :rfc:`7516`                               |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``k``               | string | this parameter repeats for each jwt key in the array and specifies the base64 encoded symmetric key see :rfc:`7516`                     |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+

	**Response Example** ::

		{
			"Kabletown URI Authority": {
				"renewal_kid": "Second Key",
				"keys": [
					{
						"alg": "HS256",
						"kid": "First Key",
						"kty": "oct",
						"k": "Kh_RkUMj-fzbD37qBnDf_3e_RvQ3RP9PaSmVEpE24AM"
					},
					{
						"alg": "HS256",
						"kid": "Second Key",
						"kty": "oct",
						"k": "fZBpDBNbk2GqhwoB_DGBAsBxqQZVix04rIoLJ7p_RlE"
					}
				]
			}
		}


**POST deliveryservices/:xml_id/urisignkeys**

	Assigns URISigning objects to a delivery service.

	Authentication Required: Yes

	Role(s) Required: admin

	**Request Route Parameters**

	+-----------+----------+----------------------------------------+
	|    Name   | Required |              Description               |
	+===========+==========+========================================+
	|   xml_id  | yes      | xml_id of the desired delivery service |
	+-----------+----------+----------------------------------------+

	**Request Properties**

	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	|    Parameter        |  Type  |                                                               Description                                                               |
	+=====================+========+=========================================================================================================================================+
	| ``Issuer``          | string | a string describing the issuer of the URI signing object. Multiple URISigning objects may be returned in a response, see example        |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``renewal_kid``     | string | a string naming the jwt key used for renewals.                                                                                          |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``keys``            | string | json array of jwt symmetric keys                                                             .                                          |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``alg``             | string | this parameter repeats for each jwt key in the array and specifies the jwa encryption algorithm to use with this key, :rfc:`7518`       |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``kid``             | string | this parameter repeats for each jwt key in the array and specifies the unique id for the key as defined in :rfc:`7516`                  |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``kty``             | string | this parameter repeats for each jwt key in the array and specifies the key type as defined in :rfc:`7516`                               |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``k``               | string | this parameter repeats for each jwt key in the array and specifies the base64 encoded symmetric key see :rfc:`7516`                     |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+

	**Request Example** ::

		{
			"Kabletown URI Authority": {
				"renewal_kid": "Second Key",
				"keys": [
					{
						"alg": "HS256",
						"kid": "First Key",
						"kty": "oct",
						"k": "Kh_RkUMj-fzbD37qBnDf_3e_RvQ3RP9PaSmVEpE24AM"
					},
					{
						"alg": "HS256",
						"kid": "Second Key",
						"kty": "oct",
						"k": "fZBpDBNbk2GqhwoB_DGBAsBxqQZVix04rIoLJ7p_RlE"
					}
				]
			}
		}

**PUT deliveryservices/:xml_id/urisignkeys**

	updates URISigning objects on a delivery service.

	Authentication Required: Yes

	Role(s) Required: admin

	**Request Route Parameters**

	+-----------+----------+----------------------------------------+
	|    Name   | Required |              Description               |
	+===========+==========+========================================+
	|  xml_id   | yes      | xml_id of the desired delivery service |
	+-----------+----------+----------------------------------------+

	**Request Properties**

	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	|    Parameter        |  Type  |                                                               Description                                                               |
	+=====================+========+=========================================================================================================================================+
	| ``Issuer``          | string | a string describing the issuer of the URI signing object. Multiple URISigning objects may be returned in a response, see example        |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``renewal_kid``     | string | a string naming the jwt key used for renewals.                                                                                          |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``keys``            | string | json array of jwt symmetric keys                                                             .                                          |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``alg``             | string | this parameter repeats for each jwt key in the array and specifies the jwa encryption algorithm to use with this key, :rfc:`7518`       |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``kid``             | string | this parameter repeats for each jwt key in the array and specifies the unique id for the key as defined in :rfc:`7516`                  |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``kty``             | string | this parameter repeats for each jwt key in the array and specifies the key type as defined in :rfc:`7516`                               |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``k``               | string | this parameter repeats for each jwt key in the array and specifies the base64 encoded symmetric key see :rfc:`7516`                     |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+

	**Request Example** ::

		{
			"Kabletown URI Authority": {
				"renewal_kid": "Second Key",
				"keys": [
					{
						"alg": "HS256",
						"kid": "First Key",
						"kty": "oct",
						"k": "Kh_RkUMj-fzbD37qBnDf_3e_RvQ3RP9PaSmVEpE24AM"
					},
					{
						"alg": "HS256",
						"kid": "Second Key",
						"kty": "oct",
						"k": "fZBpDBNbk2GqhwoB_DGBAsBxqQZVix04rIoLJ7p_RlE"
					}
				]
			}
		}

|
