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

.. _to-api-v4-users:

*********
``users``
*********

``GET``
=======
Retrieves all requested users.

:Auth. Required: Yes
:Roles Required: None\ [#tenancy]_
:Permissions Required: USER:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------+----------+------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                              |
	+===========+==========+==========================================================================================+
	| id        | no       | Return only the user identified by this integral, unique identifier                      |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| tenant    | no       | Return only users belonging to the :term:`Tenant` identified by tenant name              |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| role      | no       | Return only users belonging to the :term:`Role` identified by role name                  |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| username  | no       | Return only the user with this username                                                  |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| orderby   | no       | Choose the ordering of the results - must be the name of one of the fields of the        |
	|           |          | objects in the ``response`` array                                                        |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| sortOrder | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc") |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| limit     | no       | Choose the maximum number of results to return                                           |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| offset    | no       | The number of results to skip before beginning to return results. Must use in            |
	|           |          | conjunction with limit                                                                   |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| page      | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages |
	|           |          | are ``limit`` long and the first page is 1. If ``offset`` was defined, this query        |
	|           |          | parameter has no effect. ``limit`` must be defined to make use of ``page``.              |
	+-----------+----------+------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/4.0/users?username=mike HTTP/1.1
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
	Set-Cookie: mojolicious=...; Path=/; Expires=Fri, 13 May 2022 23:16:14 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 13 May 2022 22:16:14 GMT
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

``POST``
========
Creates a new user.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [#tenancy]_
:Permissions Required: USER:CREATE, USER:READ
:Response Type:  Object

Request Structure
-----------------
:addressLine1:       An optional field which should contain the user's address - including street name and number
:addressLine2:       An optional field which should contain an additional address field for e.g. apartment number
:city:               An optional field which should contain the name of the city wherein the user resides
:company:            An optional field which should contain the name of the company for which the user works
:country:            An optional field which should contain the name of the country wherein the user resides
:email:              The user's email address The given email is validated (circuitously) by `GitHub user asaskevich's regular expression <https://github.com/asaskevich/govalidator/blob/9a090521c4893a35ca9a228628abf8ba93f63108/patterns.go#L7>`_ . Note that it can't actually distinguish a valid, deliverable, email address but merely ensure the email is in a commonly-found format.
:fullName:           The user's full name, e.g. "John Quincy Adams"
:gid:                A deprecated field only kept for legacy compatibility reasons that used to contain the UNIX group ID of the user

	.. deprecated:: 4.0
		This field is serves no known purpose, and shouldn't be used for anything so it can be removed in the future.

:localPasswd:     The user's password
:newUser:         An optional meta field with no apparent purpose - don't use this
:phoneNumber:     An optional field which should contain the user's phone number
:postalCode:      An optional field which should contain the user's postal code
:publicSshKey:    An optional field which should contain the user's public encryption key used for the SSH protocol
:role:            The name that corresponds to the highest permission role which will be permitted to the user
:stateOrProvince: An optional field which should contain the name of the state or province in which the user resides
:tenantId:        The integral, unique identifier of the tenant to which the new user shall belong
:ucdn:            The name of the :abbr:`uCDN (Upstream Content Delivery Network)` to which the user belongs

	.. versionadded:: 4.0

:uid: A deprecated field only kept for legacy compatibility reasons that used to contain the UNIX user ID of the user

	.. deprecated:: 4.0
		This field is serves no known purpose, and shouldn't be used for anything so it can be removed in the future.

:username: The new user's username

.. code-block:: http
	:caption: Request Example

	POST /api/4.0/users HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 304
	Content-Type: application/json

	{
		"username": "mike",
		"addressLine1": "22 Mike Wazowski You've Got Your Life Back Lane",
		"city": "Monstropolis",
		"compary": "Monsters Inc.",
		"email": "mwazowski@minc.biz",
		"fullName": "Mike Wazowski",
		"localPasswd": "BFFsully",
		"confirmLocalPasswd": "BFFsully",
		"newUser": true,
		"role": "admin",
		"tenantId": 1
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

	HTTP/1.1 201 Created
	Content-Encoding: gzip
	Content-Type: application/json
	Location: /api/4.0/users?id=3
	Permissions-Policy: interest-cohort=()
	Set-Cookie: mojolicious=...; Path=/; Expires=Fri, 13 May 2022 23:13:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 13 May 2022 22:13:54 GMT
	Content-Length: 382

	{ "alerts": [
		{
			"text": "user was created.",
			"level": "success"
		}
	],
	"response": {
		"addressLine1": "22 Mike Wazowski You've Got Your Life Back Lane",
		"addressLine2": null,
		"changeLogCount": null,
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
	}}

.. [#tenancy] While no roles are required, this endpoint does respect tenancy. A user will only be able to see, create, delete or modify other users belonging to the same tenant, or its descendants.
