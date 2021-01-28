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

.. _to-api-v3-deliveryservices-xmlid-urisignkeys:

*******************************************
``deliveryservices/{{xml_id}}/urisignkeys``
*******************************************

``DELETE``
==========
Deletes URISigning objects for a :term:`Delivery Service`.

:Auth. Required: Yes
:Roles Required: admin\ [#tenancy]_
:Response Type:  ``undefined``

Request Structure
-----------------

.. table:: Request Path Parameters

	+-----------+----------+----------------------------------------+
	|    Name   | Required |              Description               |
	+===========+==========+========================================+
	| xml_id    | yes      | xml_id of the desired delivery service |
	+-----------+----------+----------------------------------------+

Response Structure
------------------
TBD

``GET``
=======
Retrieves one or more URISigning objects for a delivery service.

:Auth. Required: Yes
:Roles Required: admin\ [#tenancy]_
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Route Parameters

	+-----------+----------+----------------------------------------+
	|    Name   | Required |              Description               |
	+===========+==========+========================================+
	| xml_id    | yes      | xml_id of the desired delivery service |
	+-----------+----------+----------------------------------------+

Response Structure
------------------

:Issuer:      a string describing the issuer of the URI signing object. Multiple URISigning objects may be returned in a response, see example
:renewal_kid: a string naming the jwt key used for renewals
:keys:        json array of jwt symmetric keys
:alg:         this parameter repeats for each jwt key in the array and specifies the jwa encryption algorithm to use with this key, :rfc:`7518`
:kid:         this parameter repeats for each jwt key in the array and specifies the unique id for the key as defined in :rfc:`7516`
:kty:         this parameter repeats for each jwt key in the array and specifies the key type as defined in :rfc:`7516`
:k:           this parameter repeats for each jwt key in the array and specifies the base64 encoded symmetric key see :rfc:`7516`

.. code-block:: json
	:caption: Response Example

	{ "Kabletown URI Authority": {
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
	}}


``POST``
========
Assigns URISigning objects to a delivery service.

:Auth. Required: Yes
:Roles Required: admin\ [#tenancy]_
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------+----------+----------------------------------------+
	|    Name   | Required |              Description               |
	+===========+==========+========================================+
	|   xml_id  | yes      | xml_id of the desired delivery service |
	+-----------+----------+----------------------------------------+

Request Structure
-----------------
:Issuer:      a string describing the issuer of the URI signing object. Multiple URISigning objects may be returned in a response, see example
:renewal_kid: a string naming the jwt key used for renewals
:keys:        json array of jwt symmetric keys
:alg:         this parameter repeats for each jwt key in the array and specifies the jwa encryption algorithm to use with this key, :rfc:`7518`
:kid:         this parameter repeats for each jwt key in the array and specifies the unique id for the key as defined in :rfc:`7516`
:kty:         this parameter repeats for each jwt key in the array and specifies the key type as defined in :rfc:`7516`
:k:           this parameter repeats for each jwt key in the array and specifies the base64 encoded symmetric key see :rfc:`7516`

.. code-block:: json
	:caption: Request Example

	{ "Kabletown URI Authority": {
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
	}}

``PUT``
=======
updates URISigning objects on a delivery service.

:Auth. Required: Yes
:Roles Required: admin\ [#tenancy]_
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------+----------+----------------------------------------+
	|    Name   | Required |              Description               |
	+===========+==========+========================================+
	|  xml_id   | yes      | xml_id of the desired delivery service |
	+-----------+----------+----------------------------------------+

Request Structure
-----------------
:Issuer:      a string describing the issuer of the URI signing object. Multiple URISigning objects may be returned in a response, see example
:renewal_kid: a string naming the jwt key used for renewals
:keys:        json array of jwt symmetric keys
:alg:         this parameter repeats for each jwt key in the array and specifies the jwa encryption algorithm to use with this key, :rfc:`7518`
:kid:         this parameter repeats for each jwt key in the array and specifies the unique id for the key as defined in :rfc:`7516`
:kty:         this parameter repeats for each jwt key in the array and specifies the key type as defined in :rfc:`7516`
:k:           this parameter repeats for each jwt key in the array and specifies the base64 encoded symmetric key see :rfc:`7516`

.. code-block:: json
	:caption: Request Example

	{ "Kabletown URI Authority": {
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
	}}

.. [#tenancy] URI Signing Keys can only be created, viewed, deleted, or modified on :term:`Delivery Services` that either match the requesting user's :term:`Tenant` or are descendants thereof.
