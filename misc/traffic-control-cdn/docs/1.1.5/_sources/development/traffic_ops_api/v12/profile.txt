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

.. _to-api-v12-profile:


Profiles
========

.. _to-api-v12-profiles-route:

/api/1.2/profiles
+++++++++++++++++

**GET /api/1.2/profiles**

	Authentication Required: Yes

	**Response Properties**

	+-----------------+--------+----------------------------------------------------+
	|    Parameter    |  Type  |                    Description                     |
	+=================+========+====================================================+
	| ``lastUpdated`` | array  | The Time / Date this server entry was last updated |
	+-----------------+--------+----------------------------------------------------+
	| ``name``        | string | The name for the profile                           |
	+-----------------+--------+----------------------------------------------------+
	| ``id``          | string | Primary key                                        |
	+-----------------+--------+----------------------------------------------------+
	| ``description`` | string | The description for the profile                    |
	+-----------------+--------+----------------------------------------------------+


	**Response Example** ::

		TBD
  

| 


**GET /api/1.2/profiles/trimmed.json**

	Authentication Required: Yes

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

	 	TBD 

