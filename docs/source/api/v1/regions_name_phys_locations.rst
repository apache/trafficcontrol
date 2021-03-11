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

.. _to-api-v1-regions-name-phys_locations:

***************************************
``regions/:region_name/phys_locations``
***************************************
.. deprecated:: ATCv4

``POST``
========
Creates a new :term:`Physical Location` within the specified :term:`Region`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------------------------------------+
	| Name |                Description                                                |
	+======+===========================================================================+
	| name | The name of the region in which the new physical location will be created |
	+------+---------------------------------------------------------------------------+

:address:   The physical location's street address
:city:      The name of the city in which the physical location lies
:comments:  An optional string for containing any and all human-readable comments
:email:     An optional string containing email address of the physical location's ``poc``
:name:      An optional name of the physical location
:phone:     An optional string containing the phone number where the the physical location's ``poc`` might be reached
:poc:       The name of a "point of contact" for the physical location
:shortName: An abbreviation of the ``name``
:state:     An abbreviation of the name of the state or province within which this physical location lies
:zip:       The zip code of the physical location

.. code-block:: http
	:caption: Request Structure

	POST /api/1.4/regions/Greater_London/phys_locations HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 326
	Content-Type: application/json

	{
		"address": "Buckingham Palace",
		"city": "London",
		"comments": "Buckingham Palace",
		"email": "steve.kingstone@royal.gsx.gov.uk",
		"name": "Great_Britain",
		"phone": "0-843-816-6276",
		"poc": "Her Majesty The Queen Elizabeth Alexandra Mary Windsor II",
		"regionId": 3,
		"shortName": "uk",
		"state": "Westminster",
		"zip": "SW1A 1AA"
	}


Response Structure
------------------
:address:     The physical location's street address
:city:        The name of the city in which the physical location lies
:comments:    Any and all human-readable comments
:email:       The email address of the physical location's ``poc``
:id:          An integral, unique identifier for the physical location
:name:        The name of the physical location
:phone:       A phone number where the the physical location's ``poc`` might be reached
:poc:         The name of a "point of contact" for the physical location
:regionId:    An integral, unique identifier for the region within which the physical location lies
:regionName:  The name of the region within which the physical location lies
:shortName:   An abbreviation of the ``name``
:state:       An abbreviation of the name of the state or province within which this physical location lies
:zip:         The zip code of the physical location

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Thu, 06 Dec 2018 00:44:58 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: Pjhw/e8+2I4gQiazMv0aGzTAebmZD5yBsI5iyT7MzGbQbkuIlH4k7qlYa9JiiN9ExT69p+P8NgOQyKKsvOnmmg==
	Content-Length: 354

	{ "response": {
		"regionName": "Greater_London",
		"poc": "Her Majesty The Queen Elizabeth Alexandra Mary Windsor II",
		"name": "Great_Britain",
		"comments": "Buckingham Palace",
		"phone": "0-843-816-6276",
		"state": "Westminster",
		"regionId": 3,
		"email": "steve.kingstone@royal.gsx.gov.uk",
		"zip": "SW1A 1AA",
		"city": "London",
		"id": 4,
		"address": "Buckingham Palace",
		"shortName": "uk"
	},
	"alerts": [
		{
			"level": "warning",
			"text": "This endpoint is deprecated, please use 'POST /phys_locations' instead"
		}
	]}
