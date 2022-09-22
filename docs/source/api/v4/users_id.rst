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

.. _to-api-v4-users-id:

****************
``users/{{ID}}``
****************

``GET``
=======
Retrieves a specific user.

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: USER:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-------------------------------------------------------------+
	| Name |                       Description                           |
	+======+=============================================================+
	|  ID  | The integral, unique identifier of the user to be retrieved |
	+------+-------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/4.0/users/3 HTTP/1.1
	User-Agent: python-requests/2.25.1
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...


Response Structure
------------------
:addressLine1:   The user's address - including street name and number
:addressLine2:   An additional address field for e.g. apartment number
:changeLogCount: The number of change log entries created by the user
:city:           The name of the city wherein the user resides
:company:        The name of the company for which the user works
:country:        The name of the country wherein the user resides
:email:          The user's email address
:fullName:       The user's full name, e.g. "John Quincy Adams"
:gid:            A deprecated field only kept for legacy compatibility reasons that used to contain the UNIX group ID of the user - now it is always ``null``

	.. deprecated:: 4.0
		This field is serves no known purpose, and shouldn't be used for anything so it can be removed in the future.

:id:                An integral, unique identifier for this user
:lastAuthenticated: The date and time at which the user was last authenticated, in :rfc:`3339`
:lastUpdated:       The date and time at which the user was last modified, in :ref:`non-rfc-datetime`
:newUser:           A meta field with no apparent purpose that is usually ``null`` unless explicitly set during creation or modification of a user via some API endpoint
:phoneNumber:       The user's phone number
:postalCode:        The postal code of the area in which the user resides
:publicSshKey:      The user's public key used for the SSH protocol
:registrationSent:  If the user was created using the :ref:`to-api-v4-users-register` endpoint, this will be the date and time at which the registration email was sent - otherwise it will be ``null``
:role:              The name of the role assigned to this user
:stateOrProvince:   The name of the state or province where this user resides
:tenant:            The name of the tenant to which this user belongs
:tenantId:          The integral, unique identifier of the tenant to which this user belongs
:ucdn:              The name of the :abbr:`uCDN (Upstream Content Delivery Network)` to which the user belongs

	.. versionadded:: 4.0

:uid: A deprecated field only kept for legacy compatibility reasons that used to contain the UNIX user ID of the user - now it is always ``null``

	.. deprecated:: 4.0
		This field is serves no known purpose, and shouldn't be used for anything so it can be removed in the future.

:username: The user's username

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Encoding: gzip
	Content-Type: application/json
	Permissions-Policy: interest-cohort=()
	Set-Cookie: mojolicious=...; Path=/; Expires=Fri, 13 May 2022 23:48:14 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 13 May 2022 22:48:14 GMT
	Content-Length: 350

	{ "response": [
		{
			"addressLine1": "22 Mike Wazowski You've Got Your Life Back Lane",
			"addressLine2": null,
			"changeLogCount": 0,
			"city": "Monstropolis",
			"company": null,
			"country": null,
			"email": "mwazowski@minc.biz",
			"fullName": "Mike Wazowski",
			"gid": null,
			"id": 3,
			"lastAuthenticated": null,
			"lastUpdated": "2022-05-13T22:13:54.605052Z",
			"newUser": true,
			"phoneNumber": null,
			"postalCode": null,
			"publicSshKey": null,
			"registrationSent": null,
			"role": "admin",
			"stateOrProvince": null,
			"tenant": "root",
			"tenantId": 1,
			"ucdn": "",
			"uid": null,
			"username": "mike"
		}
	]}

``PUT``
=======

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: USER:UPDATE, USER:READ
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------+
	| Name |                       Description                          |
	+======+============================================================+
	|  ID  | The integral, unique identifier of the user to be modified |
	+------+------------------------------------------------------------+

:addressLine1:       An optional field which should contain the user's address - including street name and number
:addressLine2:       An optional field which should contain an additional address field for e.g. apartment number
:city:               An optional field which should contain the name of the city wherein the user resides
:company:            An optional field which should contain the name of the company for which the user works
:country:            An optional field which should contain the name of the country wherein the user resides
:email:              The user's email address The given email is validated (circuitously) by `GitHub user asaskevich's regular expression <https://github.com/asaskevich/govalidator/blob/9a090521c4893a35ca9a228628abf8ba93f63108/patterns.go#L7>`_ . Note that it can't actually distinguish a valid, deliverable, email address but merely ensure the email is in a commonly-found format.
:fullName:           The user's full name, e.g. "John Quincy Adams"
:gid:            A deprecated field only kept for legacy compatibility reasons that used to contain the UNIX group ID of the user - now it is always ``null``

	.. deprecated:: 4.0
		This field is serves no known purpose, and shouldn't be used for anything so it can be removed in the future.

:id:              This field *may* optionally be given, but **must** match the user's existing ID as IDs are immutable
:localPasswd:     The user's password
:newUser:         An optional meta field with no apparent purpose - don't use this
:phoneNumber:     An optional field which should contain the user's phone number
:postalCode:      An optional field which should contain the user's postal code
:publicSshKey:    An optional field which should contain the user's public encryption key used for the SSH protocol
:role:            The name of the Role which will be granted to the user
:stateOrProvince: An optional field which should contain the name of the state or province in which the user resides
:tenantId:        The integral, unique identifier of the tenant to which the new user shall belong
:ucdn:            The name of the :abbr:`uCDN (Upstream Content Delivery Network)` to which the user belongs

	.. versionadded:: 4.0

:uid: A deprecated field only kept for legacy compatibility reasons that used to contain the UNIX user ID of the user - now it is always ``null``

	.. deprecated:: 4.0
		This field is serves no known purpose, and shouldn't be used for anything so it can be removed in the future.

:username: The user's username

.. code-block:: http
	:caption: Request Structure

	PUT /api/4.0/users/3 HTTP/1.1
	User-Agent: python-requests/2.25.1
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 476

	{
		"addressLine1": "not a real address",
		"addressLine2": "not a real address either",
		"city": "not a real city",
		"company": "not a real company",
		"country": "not a real country",
		"email": "mwazowski@minc.biz",
		"fullName": "Mike Wazowski",
		"phoneNumber": "not a real phone number",
		"postalCode": "not a real postal code",
		"publicSshKey": "not a real ssh key",
		"stateOrProvince": "not a real state or province",
		"tenantId": 1,
		"role": "admin",
		"username": "mike"
	}


Response Structure
------------------
:addressLine1:   The user's address - including street name and number
:addressLine2:   An additional address field for e.g. apartment number
:changeLogCount: The number of change log entries created by the user
:city:           The name of the city wherein the user resides
:company:        The name of the company for which the user works
:country:        The name of the country wherein the user resides
:email:          The user's email address
:fullName:       The user's full name, e.g. "John Quincy Adams"
:gid:            A deprecated field only kept for legacy compatibility reasons that used to contain the UNIX group ID of the user - now it is always ``null``

	.. deprecated:: 4.0
		This field is serves no known purpose, and shouldn't be used for anything so it can be removed in the future.

:id:               An integral, unique identifier for this user
:lastAuthenticated: The date and time at which the user was last authenticated, in :rfc:`3339`
:lastUpdated:      The date and time at which the user was last modified, in :ref:`non-rfc-datetime`
:newUser:          A meta field with no apparent purpose that is usually ``null`` unless explicitly set during creation or modification of a user via some API endpoint
:phoneNumber:      The user's phone number
:postalCode:       The postal code of the area in which the user resides
:publicSshKey:     The user's public key used for the SSH protocol
:registrationSent: If the user was created using the :ref:`to-api-v4-users-register` endpoint, this will be the date and time at which the registration email was sent - otherwise it will be ``null``
:role:             The name of the role assigned to this user
:stateOrProvince:  The name of the state or province where this user resides
:tenant:           The name of the tenant to which this user belongs
:tenantId:         The integral, unique identifier of the tenant to which this user belongs
:ucdn:              The name of the :abbr:`uCDN (Upstream Content Delivery Network)` to which the user belongs

	.. versionadded:: 4.0

:uid: A deprecated field only kept for legacy compatibility reasons that used to contain the UNIX user ID of the user - now it is always ``null``

	.. deprecated:: 4.0
		This field is serves no known purpose, and shouldn't be used for anything so it can be removed in the future.

:username: The user's username

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Encoding: gzip
	Content-Type: application/json
	Permissions-Policy: interest-cohort=()
	Set-Cookie: mojolicious=...; Path=/; Expires=Fri, 13 May 2022 23:50:25 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 13 May 2022 22:50:25 GMT
	Content-Length: 399

	{ "alerts": [
		{
			"text": "user was updated.",
			"level": "success"
		}
	],
	"response": {
		"addressLine1": "not a real address",
		"addressLine2": "not a real address either",
		"changeLogCount": 0,
		"city": "not a real city",
		"company": "not a real company",
		"country": "not a real country",
		"email": "mwazowski@minc.biz",
		"fullName": "Mike Wazowski",
		"gid": null,
		"id": 3,
		"lastAuthenticated": null,
		"lastUpdated": "2022-05-13T22:50:25.965004Z",
		"newUser": false,
		"phoneNumber": "not a real phone number",
		"postalCode": "not a real postal code",
		"publicSshKey": "not a real ssh key",
		"registrationSent": null,
		"role": "admin",
		"stateOrProvince": "not a real state or province",
		"tenant": "root",
		"tenantId": 1,
		"ucdn": "",
		"uid": null,
		"username": "mike"
	}}
