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

.. _to-api-cacheassignmentgroups:

***************
``cacheassignmentgroups``
***************
.. versionadded:: 1.4

``GET``
=======
Gets a list of all cache assignment groups in the Traffic Ops database

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+---------+----------+--------------------------------------------------------------------------------+
	| Name    | Required | Description                                                                    |
	+=========+==========+================================================================================+
	| id      | no       | Return only cache assignment groups that have this integral, unique identifier |
	+---------+----------+--------------------------------------------------------------------------------+
	| name    | no       | Return only cache assignment groups with this name                             |
	+---------+----------+--------------------------------------------------------------------------------+
	| server  | no       | Return only cache assignment groups with this server assigned                  |
	+---------+----------+--------------------------------------------------------------------------------+
	| cdnId   | no       | Return only cache assignment groups with this CDN ID                           |
	+---------+----------+--------------------------------------------------------------------------------+
	| orderby | no       | Order results by this field, may be one of: id, name, cdnId                    |
	+---------+----------+--------------------------------------------------------------------------------+


Response Structure
------------------
:id:          Integral, unique, identifier for this cache assignment group
:name:        The name of the cache assignment group
:description: Description of the cache assignment group
:cdnId:       ID of the CDN to which this cache assignment group belongs
:servers:     List of server IDs assigned to this cache assignment group
:lastUpdated: The time and date at which this entry was last updated, in a ``ctime``-like format

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK

	{ "response": [
		{
            "id": 5,		
            "name": "Live Caches",            
            "description": "All Live Caches",
	        "cdnId": 1,
            "servers": [1, 2, 3],
  			"lastUpdated": "2018-10-24 16:07:05+00"
		}
	]}

``POST``
========
Creates a new cache assignment group pair

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
:name:        The name of the cache assignment group
:description: Description of the cache assignment group
:cdnId:       ID of the CDN to which this cache assignment group belongs
:servers:     List of server IDs assigned to this cache assignment group


.. code-block:: http
	:caption: Request Example

	POST /api/1.4/cacheassignmentgroups/ HTTP/1.1
	Cookie: mojolicious=...
	Content-Type: application/json

    {
       "name": "East Coast Edge Caches",
       "description": "All Edge Caches on East Coast of US",
       "cdnId": 1,
       "servers": [10, 12, 14]
    }

Response Structure
------------------
:id:          Integral, unique, identifier for this cache assignment group
:name:        The name of the cache assignment group
:description: Description of the cache assignment group
:cdnId:       ID of the CDN to which this cache assignment group belongs
:servers:     List of server IDs assigned to this cache assignment group
:lastUpdated: The time and date at which this entry was last updated, in a ``ctime``-like format

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK

	{ "response": [
		{
            "id": 5,		
            "name": "Live Caches",            
            "description": "All Live Caches",
	        "cdnId": 1,
            "servers": [1, 2, 3],
  			"lastUpdated": "2018-10-24 16:07:05+00"
		}
	  ], 
	 "alerts": [
		{
			"text": "cacheassignmentgroup was created.",
			"level": "success"
		}
	 ]
	}


``PUT``
=======
Updates a cache assignment group

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Query Parameters

	+------+----------+-----------------------------------------------------------------------+
	| Name | Required | Description                                                           |
	+======+==========+=======================================================================+
	| id   | yes      | The integral, unique identifier of the cache assignment group to edit |
	+------+----------+-----------------------------------------------------------------------+

:name:        The name of the cache assignment group
:description: Description of the cache assignment group
:cdnId:       ID of the CDN to which this cache assignment group belongs
:servers:     List of server IDs assigned to this cache assignment group


.. code-block:: http
	:caption: Request Example

	POST /api/1.4/cacheassignmentgroups/?id=4 HTTP/1.1
	Cookie: mojolicious=...
	Content-Type: application/json

    {
       "name": "East Coast Edge Caches",
       "description": "All Edge Caches on East Coast of US",
       "cdnId": 1,
       "servers": [10, 12, 14]
    }

Response Structure
------------------
:id:          Integral, unique, identifier for this cache assignment group
:name:        The name of the cache assignment group
:description: Description of the cache assignment group
:cdnId:       ID of the CDN to which this cache assignment group belongs
:servers:     List of server IDs assigned to this cache assignment group
:lastUpdated: The time and date at which this entry was last updated, in a ``ctime``-like format

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK

	{ "response": [
		{
            "id": 5,		
            "name": "Live Caches",            
            "description": "All Live Caches",
	        "cdnId": 1,
            "servers": [1, 2, 3],
  			"lastUpdated": "2018-10-24 16:07:05+00"
		}
	  ], 
	 "alerts": [
		{
			"text": "cacheassignmentgroup was created.",
			"level": "success"
		}
	 ]
	}

``DELETE``
==========
Deletes a cache assignment group

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Query Parameters

	+------+----------+-------------------------------------------------------------------------+
	| Name | Required | Description                                                             |
	+======+==========+=========================================================================+
	| id   | yes      | The integral, unique identifier of the cache assignment group to delete |
	+------+----------+-------------------------------------------------------------------------+

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	

	{ "alerts": [
		{
			"text": "cacheassignmentgroup was deleted.",
			"level": "success"
		}
	]}
