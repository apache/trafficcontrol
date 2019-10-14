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

.. _to-api-user-current-update:

***********************
``user/current/update``
***********************
.. deprecated:: 1.4
	Use the ``PUT`` method of :ref:`to-api-users` instead.

``POST``
========
Updates the date for the authenticated user.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  ``undefined``

Request Structure
-----------------
:addressLine1:       An optional field which should contain the user's address - including street name and number
:addressLine2:       An optional field which should contain an additional address field for e.g. apartment number
:city:               An optional field which should contain the name of the city wherein the user resides
:company:            An optional field which should contain the name of the company for which the user works
:confirmLocalPasswd: The 'confirm' field in a new user's password specification - must match ``localPasswd``
:country:            An optional field which should contain the name of the country wherein the user resides
:email:              The user's email address

	.. versionchanged:: 1.4
		Prior to version 1.4, the email was validated using the `Email::Valid Perl package <https://metacpan.org/pod/Email::Valid>`_ but is now validated (circuitously) by `GitHub user asaskevich's regular expression <https://github.com/asaskevich/govalidator/blob/9a090521c4893a35ca9a228628abf8ba93f63108/patterns.go#L7>`_ . Note that neither method can actually distinguish a valid, deliverable, email address but merely ensure the email is in a commonly-found format.

:fullName:        The user's full name, e.g. "John Quincy Adams"
:localPasswd:     The user's password
:newUser:         An optional meta field with no apparent purpose - don't use this
:phoneNumber:     An optional field which should contain the user's phone number
:postalCode:      An optional field which should contain the user's postal code
:publicSshKey:    An optional field which should contain the user's public encryption key used for the SSH protocol
:role:            The number that corresponds to the highest permission role which will be permitted to the user
:stateOrProvince: An optional field which should contain the name of the state or province in which the user resides
:tenantId:        The integral, unique identifier of the tenant to which the new user shall belong

	.. note:: This field is optional if and only if tenancy is not enabled in Traffic Control

:username: The user's new username

.. code-block:: http
	:caption: Request Example

	POST /api/1.4/user/current/update HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 483
	Content-Type: application/json

	{ "user": {
		"addressLine1": "not a real address",
		"addressLine2": "not a real address either",
		"city": "not a real city",
		"company": "not a real company",
		"country": "not a real country",
		"email": "not@real.email",
		"fullName": "Not a real fullName",
		"phoneNumber": "not a real phone number",
		"postalCode": "not a real postal code",
		"publicSshKey": "not a real ssh key",
		"stateOrProvince": "not a real state or province",
		"tenantId": 1,
		"role": 1,
		"username": "admin"
	}}

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Thu, 13 Dec 2018 21:04:36 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Fri, 14 Dec 2018 01:04:36 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: sHFqZQ4Cv7IIWaIejoAvM2Fr/HSupcX3D16KU/etjw+4jcK9EME3Bq5ohLC+eQ52BDCKW2Ra+AC3TfFtworJww==
	Content-Length: 79

	{ "alerts": [
		{
			"level": "success",
			"text": "User profile was successfully updated"
		},
		{
			"level": "warning",
			"text": "This endpoint is deprecated, please use 'PUT /api/1.4/user/current' instead"
		}
	]}
