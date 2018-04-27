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

.. _to-api-v11-users:

Users
=====

.. _to-api-v11-users-route:

/api/1.1/users
++++++++++++++

**GET /api/1.1/users**

  Retrieves all users.

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``addressLine1``      | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``addressLine2``      | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``city``              | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``company``           | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``country``           | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``email``             | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``fullName``          | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``gid``               | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``id``                | hash   |                                                |
  +----------------------+--------+------------------------------------------------+
  |``lastUpdated``       | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``newUser``           | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``phoneNumber``       | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``postalCode``        | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``publicSshKey``      | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``registrationSent``  | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``role``              | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``roleName``          | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``stateOrProvince``   | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``uid``               | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``username``          | string |                                                |
  +----------------------+--------+------------------------------------------------+

  **Response Example** ::

   {
      "response": [
		 {
			"addressLine1": "",
			"addressLine2": "",
			"city": "",
			"company": "",
			"country": "",
			"email": "email1@email.com",
			"fullName": "Tom Simpson",
			"gid": "0",
			"id": "53",
			"lastUpdated": "2016-01-26 10:22:07",
			"newUser": true,
			"phoneNumber": "",
			"postalCode": "",
			"publicSshKey": "xxx",
			"registrationSent": true,
			"role": "6",
			"rolename": "admin",
			"stateOrProvince": "",
			"uid": "0",
			"username": "tsimpson"
		 },
		 {
		 	... more users
		 },
        ]
    }

|


**GET /api/1.1/users/:id**

  Retrieves user by ID.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +-----------+----------+---------------------------------------------+
  |   Name    | Required |                Description                  |
  +===========+==========+=============================================+
  |   ``id``  |   yes    | User id.                                    |
  +-----------+----------+---------------------------------------------+

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``addressLine1``      | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``addressLine2``      | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``city``              | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``company``           | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``country``           | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``email``             | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``fullName``          | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``gid``               | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``id``                | hash   |                                                |
  +----------------------+--------+------------------------------------------------+
  |``lastUpdated``       | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``newUser``           | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``phoneNumber``       | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``postalCode``        | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``publicSshKey``      | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``registrationSent``  | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``role``              | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``roleName``          | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``stateOrProvince``   | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``uid``               | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``username``          | string |                                                |
  +----------------------+--------+------------------------------------------------+

  **Response Example** ::

   {
      "response": [
		 {
			"addressLine1": "",
			"addressLine2": "",
			"city": "",
			"company": "",
			"country": "",
			"email": "email1@email.com",
			"fullName": "Tom Simpson",
			"gid": "0",
			"id": "53",
			"lastUpdated": "2016-01-26 10:22:07",
			"newUser": true,
			"phoneNumber": "",
			"postalCode": "",
			"publicSshKey": "xxx",
			"registrationSent": true,
			"role": "6",
			"rolename": "admin",
			"stateOrProvince": "",
			"uid": "0",
			"username": "tsimpson"
		 }
        ]
    }

|


**GET /api/1.1/user/current**

  Retrieves the profile for the authenticated user.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``email``             | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``city``              | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``id``                | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``phoneNumber``       | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``company``           | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``country``           | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``fullName``          | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``localUser``         | boolean|                                                |
  +----------------------+--------+------------------------------------------------+
  |``uid``               | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``stateOrProvince``   | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``username``          | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``newUser``           | boolean|                                                |
  +----------------------+--------+------------------------------------------------+
  |``addressLine2``      | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``role``              | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``addressLine1``      | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``gid``               | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``postalCode``        | string |                                                |
  +----------------------+--------+------------------------------------------------+

  **Response Example** ::

    {
           "response": {
                            "email": "email@email.com",
                            "city": "",
                            "id": "50",
                            "phoneNumber": "",
                            "company": "",
                            "country": "",
                            "fullName": "Tom Callahan",
                            "localUser": true,
                            "uid": "0",
                            "stateOrProvince": "",
                            "username": "tommyboy",
                            "newUser": false,
                            "addressLine2": "",
                            "role": "6",
                            "addressLine1": "",
                            "gid": "0",
                            "postalCode": ""
           },
    }

|
  
**POST /api/1.1/user/current/update**

  Updates the date for the authenticated user.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``email``             | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``city``              | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``id``                | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``phoneNumber``       | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``company``           | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``country``           | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``fullName``          | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``localUser``         | boolean|                                                |
  +----------------------+--------+------------------------------------------------+
  |``uid``               | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``stateOrProvince``   | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``username``          | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``newUser``           | boolean|                                                |
  +----------------------+--------+------------------------------------------------+
  |``addressLine2``      | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``role``              | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``addressLine1``      | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``gid``               | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``postalCode``        | string |                                                |
  +----------------------+--------+------------------------------------------------+

  **Request Example** ::

    {
     "user": {
        "email": "",
        "city": "",
        "id": "",
        "phoneNumber": "",
        "company": "",
        "country": "",
        "fullName": "",
        "localUser": true,
        "uid": "0",
        "stateOrProvince": "",
        "username": "tommyboy",
        "newUser": false,
        "addressLine2": "",
        "role": "6",
        "addressLine1": "",
        "gid": "0",
        "postalCode": ""
     }
    }

  **Response Properties**

  +-------------+--------+----------------------------------+
  |  Parameter  |  Type  |           Description            |
  +=============+========+==================================+
  | ``alerts``  | array  | A collection of alert messages.  |
  +-------------+--------+----------------------------------+
  | ``>level``  | string | Success, info, warning or error. |
  +-------------+--------+----------------------------------+
  | ``>text``   | string | Alert message.                   |
  +-------------+--------+----------------------------------+
  | ``version`` | string |                                  |
  +-------------+--------+----------------------------------+

  **Response Example** ::

    {
          "alerts": [
                    {
                            "level": "success",
                            "text": "UserProfile was successfully updated."
                    }
            ],
    }

|

**GET /api/1.1/user/current/jobs.json**

  Retrieves the user's list of jobs.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Query Parameters**

  +--------------+----------+----------------------------------------+
  |    Name      | Required |              Description               |
  +==============+==========+========================================+
  | ``keyword``  | no       | PURGE                                  |
  +--------------+----------+----------------------------------------+

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``keyword``           | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``objectName``        | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``assetUrl``          | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``assetType``         | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``status``            | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``dsId``              | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``dsXmlId``           | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``username``          | boolean|                                                |
  +----------------------+--------+------------------------------------------------+
  |``parameters``        | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``enteredTime``       | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``objectType``        | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``agent``             | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``id``                | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``startTime``         | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``version``           | string |                                                |
  +----------------------+--------+------------------------------------------------+

  **Response Example**
  ::

    {
     "response": [
        {
           "id": "1",
           "keyword": "PURGE",
           "objectName": null,
           "assetUrl": "",
           "assetType": "file",
           "status": "PENDING",
           "dsId": "9999",
           "dsXmlId": "ds-xml-id",
           "username": "peewee",
           "parameters": "TTL:56h",
           "enteredTime": "2015-01-21 18:00:16",
           "objectType": null,
           "agent": "",
           "startTime": "2015-01-21 10:45:38"
        }
     ],
    }


|

**POST/api/1.1/user/current/jobs**

Invalidating content on the CDN is sometimes necessary when the origin was mis-configured and something is cached in the CDN that needs to be removed. Given the size of a typical Traffic Control CDN and the amount of content that can be cached in it, removing the content from all the caches may take a long time. To speed up content invalidation, Traffic Ops will not try to remove the content from the caches, but it makes the content inaccessible using the *regex_revalidate* ATS plugin. This forces a *revalidation* of the content, rather than a new get.

.. Note:: This method forces a HTTP *revalidation* of the content, and not a new *GET* - the origin needs to support revalidation according to the HTTP/1.1 specification, and send a ``200 OK`` or ``304 Not Modified`` as applicable.

Authentication Required: Yes

Role(s) Required: Yes

  **Request Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``dsId``              | string | Unique Delivery Service ID                     |
  +----------------------+--------+------------------------------------------------+
  |``regex``             | string | Path Regex this should be a                    |
  |                      |        | `PCRE <http://www.pcre.org/>`_ compatible      |
  |                      |        | regular expression for the path to match for   |
  |                      |        | forcing the revalidation. Be careful to only   |
  |                      |        | match on the content you need to remove -      |
  |                      |        | revalidation is an expensive operation for     |
  |                      |        | many origins, and a simple ``/.*`` can cause   |
  |                      |        | an overload condition of the origin.           |
  +----------------------+--------+------------------------------------------------+
  |``startTime``         | string | Start Time is the time when the revalidation   |
  |                      |        | rule will be made active. Populate             |
  |                      |        | with the current time to schedule ASAP.        |
  +----------------------+--------+------------------------------------------------+
  |``ttl``               | int    | Time To Live is how long the revalidation rule |
  |                      |        | will be active for in hours. It usually makes  |
  |                      |        | sense to make this the same as the             |
  |                      |        | ``Cache-Control`` header from the origin which |
  |                      |        | sets the object time to live in cache          |
  |                      |        | (by ``max-age`` or ``Expires``). Entering a    |
  |                      |        | longer TTL here will make the caches do        |
  |                      |        | unnecessary work.                              |
  +----------------------+--------+------------------------------------------------+

  **Request Example** ::

    {
           "dsId": "9999",
           "regex": "/path/to/content.jpg",
           "startTime": "2015-01-27 11:08:37",
           "ttl": 54
    }

|

  **Response Properties**

  +-------------+--------+----------------------------------+
  |  Parameter  |  Type  |           Description            |
  +=============+========+==================================+
  | ``alerts``  | array  | A collection of alert messages.  |
  +-------------+--------+----------------------------------+
  | ``>level``  | string | Success, info, warning or error. |
  +-------------+--------+----------------------------------+
  | ``>text``   | string | Alert message.                   |
  +-------------+--------+----------------------------------+
  | ``version`` | string |                                  |
  +-------------+--------+----------------------------------+

  **Response Example** ::

    {
          "alerts":
                  [
                      { 
                            "level": "success",
                            "text": "Successfully created purge job for: ."
                      }
                  ],
    }


|

**POST /api/1.1/user/login**

  Authentication of a user using username and password. Traffic Ops will send back a session cookie.

  Authentication Required: No

  Role(s) Required: None

  **Request Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``u``                 | string | username                                       |
  +----------------------+--------+------------------------------------------------+
  |``p``                 | string | password                                       |
  +----------------------+--------+------------------------------------------------+

  **Request Example** ::

    {
       "u": "username",
       "p": "password"
    }

|

  **Response Properties**

  +-------------+--------+----------------------------------+
  |  Parameter  |  Type  |           Description            |
  +=============+========+==================================+
  | ``alerts``  | array  | A collection of alert messages.  |
  +-------------+--------+----------------------------------+
  | ``>level``  | string | Success, info, warning or error. |
  +-------------+--------+----------------------------------+
  | ``>text``   | string | Alert message.                   |
  +-------------+--------+----------------------------------+
  | ``version`` | string |                                  |
  +-------------+--------+----------------------------------+

  **Response Example** ::

   {
     "alerts": [
        {
           "level": "success",
           "text": "Successfully logged in."
        }
     ],
    }

|

**GET /api/1.1/user/:id/deliveryservices/available.json**

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +-----------------+----------+---------------------------------------------------+
  | Name            | Required | Description                                       |
  +=================+==========+===================================================+
  |id               | yes      |                                                   |
  +-----------------+----------+---------------------------------------------------+

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``xmlId``             | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``id``                | string |                                                |
  +----------------------+--------+------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "xmlId": "ns-img",
           "id": "90"
        },
        {
           "xmlId": "ns-img-secure",
           "id": "280"
        }
     ],
    }

|

**POST /api/1.1/user/login/token**

  Authentication of a user using a token.

  Authentication Required: No

  Role(s) Required: None

  **Request Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``t``                 | string | token-value                                    |
  +----------------------+--------+------------------------------------------------+

  **Request Example** ::

    {
       "t": "token-value"
    }

|

  **Response Properties**

  +-------------+--------+-------------+
  |  Parameter  |  Type  | Description |
  +=============+========+=============+
  | ``alerts``  | array  |             |
  +-------------+--------+-------------+
  | ``>level``  | string |             |
  +-------------+--------+-------------+
  | ``>text``   | string |             |
  +-------------+--------+-------------+
  | ``version`` | string |             |
  +-------------+--------+-------------+

  **Response Example** ::

    {
     "alerts": [
        {
           "level": "error",
           "text": "Unauthorized, please log in."
        }
     ],
    }

|

**POST /api/1.1/user/logout**

  User logout. Invalidates the session cookie.

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``alerts``            | array  |                                                |
  +----------------------+--------+------------------------------------------------+
  |* ``level``           | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |* ``text``            | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``version``           | string |                                                |
  +----------------------+--------+------------------------------------------------+

  **Response Example**

  ::

    {
     "alerts": [
        {
           "level": "success",
           "text": "You are logged out."
        }
     ],
    }


|

**POST /api/1.1/user/reset_password**

  Reset user password.

  Authentication Required: No

  Role(s) Required: None

  **Request Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``email``             | string | The email address of the user to initiate      |
  |                      |        | password reset.                                |
  +----------------------+--------+------------------------------------------------+

  **Request Example**
  ::

    {
     "email": "email@email.com"
    }

|

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``alerts``            | array  | A collection of alert messages.                |
  +----------------------+--------+------------------------------------------------+
  |* ``level``           | string | Success, info, warning or error.               |
  +----------------------+--------+------------------------------------------------+
  |* ``text``            | string | Alert message.                                 |
  +----------------------+--------+------------------------------------------------+
  |``version``           | string |                                                |
  +----------------------+--------+------------------------------------------------+

  **Response Example** ::

    

    {
     "alerts": [
        {
           "level": "success",
           "text": "Successfully sent password reset to email 'email@email.com'"
        }
     ],
    }

  
