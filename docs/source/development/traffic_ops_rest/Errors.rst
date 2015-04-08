.. 
.. Copyright 2015 Comcast Cable Communications Management, LLC
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

.. _to-api-error:

Error Responses
===============
  

**Response Properties**

+----------------------+--------+------------------------------------------------+
| Parameter            | Type   | Description                                    |
+======================+========+================================================+
|``alerts``            | array  | A collection of alert messages.                |
+----------------------+--------+------------------------------------------------+
|> ``level``           | string | Success, info, warning or error.               |
+----------------------+--------+------------------------------------------------+
|> ``text``            | string | Alert message.                                 |
+----------------------+--------+------------------------------------------------+
|``version``           | string |                                                |
+----------------------+--------+------------------------------------------------+

.. _reference-label-400:

HTTP Status Code: 400
---------------------

**Response Message** 

These errors may happen for POST /api/1.1/myuser/purge only.

::


  HTTP Status Code: 400
  Reason: Unauthorized

**Examples**

::


  {
   "alerts": [
      {
         "level": "error",
         "text": "Field [ X ] required."
      }
   ],
   "version": "1.1"
  }
  
  {
   "alerts": [
      {
         "level": "error",
         "text": "Username already taken."
      }
   ],
   "version": "1.1"
  }

**Response Message** 

This error may happen for POST /api/1.1/user/password/reset only.

::


  HTTP Status Code: 400
  Reason: Unauthorized

**Example**


::

  {
   "alerts": [
      {
         "level": "error",
         "text": "Email not found [email]."
      }
   ],
   "version": "1.1"
  }


**Response Message** 

This error may happen for POST /api/1.1/user/jobs/purge only.

::


  HTTP Status Code: 400
  Reason: Bad Request

**Example**

::


  {
   "alerts": [
      {
         "level": "error",
         "text": "[ validation error message ]"
      }
   ],
   "version": "1.1"
  }

**Response Message** 

These errors may happen for POST /api/1.1/to_extensions only.

::


  HTTP Status Code: 400
  Reason: Bad Request

**Examples**

::


  {
   "alerts": [
      {
         "level": "error",
         "text": "ToExtension update not supported; delete and re-add."
      }
   ],
   "version": "1.1"
  }

  {
   "alerts": [
      {
         "level": "error",
         "text": "Invalid Extension type: [ type ]"
      }
   ],
   "version": "1.1"
  }

  {
   "alerts": [
      {
         "level": "error",
         "text": "A Check extension is already loaded with name = [ name ]"
      }
   ],
   "version": "1.1"
  }
  
  {
   "alerts": [
      {
         "level": "error",
         "text": "No open slots left for checks, delete one first."
      }
   ],
   "version": "1.1"
  }

.. _reference-label-401:

HTTP Status Code: 401
---------------------

**Response Message** 

General error.

::


  HTTP Status Code: 401
  Reason: Unauthorized

**Example**

::


  {
   "alerts": [
      {
         "level": "error",
         "text": "Unauthorized, please log in."
      }
   ],
   "version": "1.1"
  }

.. _reference-label-403:

HTTP Status Code: 403
---------------------

**Response Message** 

General error.

::


  HTTP Status Code: 403
  Reason: Delivery service not assigned to user.

**Example**


::


  {
   "alerts": [
      {
         "level": "error",
         "text": "Forbidden"
      }
   ],
   "version": "1.1"
  }

**Response Message** 

This error may happen for POST /api/1.1/servercheck, POST /api/1.1/to_extensions/:id/delete, and POST /api/1.1/to_extensions only.


::


  HTTP Status Code: 403
  Reason: Forbidden

**Example**


::


  {
   "alerts": [
      {
         "level": "error",
         "text": "Invalid user for this API. Only the \"extension\" user can use this."
      }
   ],
   "version": "1.1"
  }
  
.. _reference-label-404:

HTTP Status Code: 404
---------------------

**Response Message** 

This error may happen for POST /api/1.1/servercheck only.

::


  HTTP Status Code: 404
  Reason: Resource not found


**Example**

::


  {
   "alerts": [
      {
         "level": "error",
         "text": "Server not found"
      }
   ],
   "version": "1.1"
  }

**Response Message** 

This error may happen for POST /api/1.1/to_extensions/:id/delete only.

::


  HTTP Status Code: 404
  Reason: Resource not found


**Example**

::


  {
   "alerts": [
      {
         "level": "error",
         "text": "ToExtension with id [ id ] not found."
      }
   ],
   "version": "1.1"
  }

**Response Message** 

This error may happen for POST /api/1.1/servercheck only.


::


  HTTP Status Code: 404
  Reason: Resource not found


**Example**

::


  {
   "alerts": [
      {
         "level": "error",
         "text": "Server Check Extension [ server check short name ] not found - Do you need to install it?"
      }
   ],
   "version": "1.1"
  } 
  
