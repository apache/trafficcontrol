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

.. _to-api-user-current:

****************
``user/current``
****************

``GET``
=======
.. caution:: As a username is needed to log in, any administrator or application must necessarily know the current username at any given time. Thus it's generally better to use the ``username`` query parameter of a ``GET`` request to :ref:`to-api-users` instead.

Retrieves the details of the authenticated user.

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: None
:Response Type:  Object

Request Structure
-----------------
No parameters available.

.. code-block:: http
	:caption: Request Example

	GET /api/5.0/user/current HTTP/1.1
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
:gid:            A deprecated field only kept for legacy compatibility reasons that used to contain the UNIX group ID of the user

	.. deprecated:: 4.0
		This field is serves no known purpose, and shouldn't be used for anything so it can be removed in the future.

:id:                An integral, unique identifier for this user
:lastAuthenticated: The date and time at which the user was last authenticated, in :rfc:`3339` format
:lastUpdated:       The date and time at which the user was last modified, in :rfc:`3339` format
:newUser:           A meta field with no apparent purpose that is usually ``null`` unless explicitly set during creation or modification of a user via some API endpoint
:phoneNumber:       The user's phone number
:postalCode:        The postal code of the area in which the user resides
:publicSshKey:      The user's public key used for the SSH protocol
:registrationSent:  If the user was created using the :ref:`to-api-users-register` endpoint, this will be the date and time at which the registration email was sent - otherwise it will be ``null``
:role:              The name of the :term:`Role` assigned to this user
:stateOrProvince:   The name of the state or province where this user resides
:tenant:            The name of the :term:`Tenant` to which this user belongs
:tenantId:          The integral, unique identifier of the :term:`Tenant` to which this user belongs
:ucdn:              The name of the :abbr:`uCDN (Upstream Content Delivery Network)` to which the user belongs
:uid:               A deprecated field only kept for legacy compatibility reasons that used to contain the UNIX user ID of the user

	.. deprecated:: 4.0
		This field is serves no known purpose, and shouldn't be used for anything so it can be removed in the future.

:username: The user's username

.. code-block:: http
	:caption: Response Example


	HTTP/1.1 200 OK
	Content-Encoding: gzip
	Content-Type: application/json
	Permissions-Policy: interest-cohort=()
	Set-Cookie: mojolicious=...; Path=/; Expires=Fri, 13 May 2022 23:42:05 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 13 May 2022 22:42:05 GMT
	Content-Length: 311

	{ "response": {
		"addressLine1": null,
		"addressLine2": null,
		"changeLogCount": 1,
		"city": null,
		"company": null,
		"country": null,
		"email": "admin@no-reply.atc.test",
		"fullName": "Development Admin User",
		"gid": null,
		"id": 2,
		"lastAuthenticated": "2022-05-13T22:42:05.495439Z",
		"lastUpdated": "2022-05-13T22:42:05.495439Z",
		"newUser": false,
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
		"username": "admin"
	}}

``PUT``
=======
.. warning:: Assuming the current user's integral, unique identifier is known, it's generally better to use the ``PUT`` method of the :ref:`to-api-users` instead.

.. warning:: Users that login via LDAP pass-back cannot be modified

Updates the date for the authenticated user.

:Auth. Required: Yes
:Roles Required: None
:Permissions Required:  None
:Response Type:  Object

Request Structure
-----------------
:addressLine1: The user's address - including street name and number
:addressLine2: An additional address field for e.g. apartment number
:city:         The name of the city wherein the user resides
:company:      The name of the company for which the user works
:country:      The name of the country wherein the user resides
:email:        The user's email address - cannot be an empty string\ [#notnull]_. The given email is validated (circuitously) by `GitHub user asaskevich's regular expression <https://github.com/asaskevich/govalidator/blob/9a090521c4893a35ca9a228628abf8ba93f63108/patterns.go#L7>`_ . Note that it can't actually distinguish a valid, deliverable, email address but merely ensure the email is in a commonly-found format.
:fullName:     The user's full name, e.g. "John Quincy Adams"
:gid:          A legacy field only kept for legacy compatibility reasons that used to contain the UNIX group ID of the user - please don't use this

	.. deprecated:: 4.0
		This field is serves no known purpose, and shouldn't be used for anything so it can be removed in the future.

:id:              The user's integral, unique, identifier - this cannot be changed\ [#notnull]_
:localPasswd:     Optionally, the user's password. This should never be given if it will not be changed. An empty string or ``null`` can be used to explicitly specify no change.
:phoneNumber:     The user's phone number
:postalCode:      The user's postal code
:publicSshKey:    The user's public encryption key used for the SSH protocol
:role:            The integral, unique identifier of the highest permission :term:`Role` which will be permitted to the user - this cannot be altered from the user's current :term:`Role`\ [#notnull]_
:stateOrProvince: The state or province in which the user resides
:tenantId:        The integral, unique identifier of the :term:`Tenant` to which the new user shall belong\ [#tenancy]_\ [#notnull]_
:ucdn:            The name of the :abbr:`uCDN (Upstream Content Delivery Network)` to which the user belongs
:uid:             A legacy field only kept for legacy compatibility reasons that used to contain the UNIX user ID of the user - please don't use this

	.. deprecated:: 4.0
		This field is serves no known purpose, and shouldn't be used for anything so it can be removed in the future.

:username: The user's new username\ [#notnull]_

.. code-block:: http
	:caption: Request Example

	PUT /api/5.0/user/current HTTP/1.1
	User-Agent: python-requests/2.25.1
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 562

	{
		"addressLine1": null,
		"addressLine2": null,
		"changeLogCount": 1,
		"city": null,
		"company": null,
		"country": null,
		"email": "admin@no-reply.atc.test",
		"fullName": "Development Admin User",
		"gid": null,
		"id": 2,
		"lastAuthenticated": "2022-05-13T22:42:05.495439Z",
		"lastUpdated": "2022-05-13T22:42:05.495439Z",
		"newUser": false,
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
		"username": "admin"
	}

Response Structure
------------------
:addressLine1:   The user's address - including street name and number
:addressLine2:   An additional address field for e.g. apartment number
:changeLogCount: The number of change log entries created by the user
:city:           The name of the city wherein the user resides
:company:        The name of the company for which the user works
:country:        The name of the country wherein the user resides
:email:          The user's email address validated (circuitously) by `GitHub user asaskevich's regular expression <https://github.com/asaskevich/govalidator/blob/9a090521c4893a35ca9a228628abf8ba93f63108/patterns.go#L7>`_ . Note that it can't actually distinguish a valid, deliverable, email address but merely ensure the email is in a commonly-found format.
:fullName:       The user's full name, e.g. "John Quincy Adams"
:gid:            A legacy field only kept for legacy compatibility reasons that used to contain the UNIX group ID of the user

	.. deprecated:: 4.0
		This field is serves no known purpose, and shouldn't be used for anything so it can be removed in the future.

:id:               An integral, unique identifier for this user
:lastAuthenticated: The date and time at which the user was last authenticated, in :rfc:`3339`
:lastUpdated:      The date and time at which the user was last modified, in :ref:`non-rfc-datetime`
:newUser:          A meta field with no apparent purpose
:phoneNumber:      The user's phone number
:postalCode:       The postal code of the area in which the user resides
:publicSshKey:     The user's public key used for the SSH protocol
:registrationSent: If the user was created using the :ref:`to-api-users-register` endpoint, this will be the date and time at which the registration email was sent - otherwise it will be ``null``
:role:             The name of the :term:`Role` assigned to this user
:stateOrProvince:  The name of the state or province where this user resides
:tenant:           The name of the :term:`Tenant` to which this user belongs
:tenantId:         The integral, unique identifier of the :term:`Tenant` to which this user belongs
:ucdn:             The name of the :abbr:`uCDN (Upstream Content Delivery Network)` to which the user belongs
:uid:              A legacy field only kept for legacy compatibility reasons that used to contain the UNIX user ID of the user

	.. deprecated:: 4.0
		This field is serves no known purpose, and shouldn't be used for anything so it can be removed in the future.

:username: The user's username

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Encoding: gzip
	Content-Type: application/json
	Permissions-Policy: interest-cohort=()
	Set-Cookie: mojolicious=...; Path=/; Expires=Fri, 13 May 2022 23:45:22 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 13 May 2022 22:45:22 GMT
	Content-Length: 370

	{ "alerts": [
		{
			"text": "User profile was successfully updated",
			"level": "success"
		}
	],
	"response": {
		"addressLine1": null,
		"addressLine2": null,
		"changeLogCount": 1,
		"city": null,
		"company": null,
		"country": null,
		"email": "admin@no-reply.atc.test",
		"fullName": "Development Admin User",
		"gid": null,
		"id": 2,
		"lastAuthenticated": "2022-05-13T22:44:55.973452Z",
		"lastUpdated": "2022-05-13T22:45:22.505401Z",
		"newUser": false,
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
		"username": "admin"
	}}

.. [#notnull] This field cannot be ``null``.
.. [#tenancy] This endpoint respects tenancy; a user cannot assign itself to a :term:`Tenant` that is not the same :term:`Tenant` to which it was previously assigned or a descendant thereof.
