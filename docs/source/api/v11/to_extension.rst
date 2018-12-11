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

.. _to-api-to_extensions:

*****************
``to_extensions``
*****************
.. seealso:: :ref:`admin-to-ext-script`

``GET``
=======
Retrieves the list of Traffic Ops extensions.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
No parameters available.

Response Structure
------------------
:additional_config_json: A string containing a JSON-encoded object with extra configuration options... inside a JSON object...
:description:            A short description of the extension

	.. note:: This is, unfortunately, ``null`` for all default plugins

:id:       An integral, unique identifier for this extension definition
:info_url: A URL where info about this extension may be found
:isactive: An integer describing the boolean notion of whether or not the plugin is active; one of:

	0
		disabled
	1
		enabled

:name:                  The name of the extension
:script_file:           The base filename of the script that runs for the extension
:servercheck_shortname: The name of the column in the table at 'Monitor' -> 'Cache Checks' in Traffic Portal, where "Check Extension" output is displayed

	.. note:: This field has meaning only for "Check Extensions"

:type:    The type of extension - there are a set number of allowed values which are not recorded anywhere at the time of this writing
:version: A (hopefully) semantic version number describing the version of the plugin

**Response Example** ::

	{
				 "response": [
								{
												script_file: "ping",
												version: "1.0.0",
												name: "ILO_PING",
												description: null,
												info_url: "http://foo.com/bar.html",
												additional_config_json: "{ "path": "/api/1.1/servers.json", "match": { "type": "EDGE"}, "select": "ilo_ip_address", "cron": "9 * * * *" }",
												isactive: "1",
												id: "1",
												type: "CHECK_EXTENSION_BOOL",
												servercheck_short_name: "ILO"
								},
								{
												script_file: "ping",
												version: "1.0.0",
												name: "10G_PING",
												description: null,
												info_url: "http://foo.com/bar.html",
												additional_config_json: "{ "path": "/api/1.1/servers.json", "match": { "type": "EDGE"}, "select": "ip_address", "cron": "18 * * * *" }",
												isactive: "1",
												id: "2",
												type: "CHECK_EXTENSION_BOOL",
												servercheck_short_name: "10G"
								}
				 ],
	}


|

**POST /api/1.1/to_extensions**

	Creates a Traffic Ops extension.

	Authentication Required: Yes

	Role(s) Required: None

	**Request Parameters**

	+--------------------------+--------+--------------------------------------------+
	| Parameter                | Type   | Description                                |
	+==========================+========+============================================+
	|``name``                  | string |                                            |
	+--------------------------+--------+--------------------------------------------+
	|``version``               | string |                                            |
	+--------------------------+--------+--------------------------------------------+
	|``info_url``              | string |                                            |
	+--------------------------+--------+--------------------------------------------+
	|``script_file``           | string |                                            |
	+--------------------------+--------+--------------------------------------------+
	|``isactive``              | string |                                            |
	+--------------------------+--------+--------------------------------------------+
	|``additional_config_json``| string |                                            |
	+--------------------------+--------+--------------------------------------------+
	|``description``           | string |                                            |
	+--------------------------+--------+--------------------------------------------+
	|``servercheck_short_name``| string |                                            |
	+--------------------------+--------+--------------------------------------------+
	|``type``                  | string |                                            |
	+--------------------------+--------+--------------------------------------------+

	**Request Example** ::


		{
					"name": "ILO_PING",
					"version": "1.0.0",
					"info_url": "http://foo.com/bar.html",
					"script_file": "ping",
					"isactive": "1",
					"additional_config_json": "{ "path": "/api/1.1/servers.json", "match": { "type": "EDGE"}",
					"description": null,
					"servercheck_short_name": "ILO"
					"type": "CHECK_EXTENSION_BOOL",
		}

|

	**Response Properties**

	+------------+--------+----------------------------------+
	| Parameter  |  Type  |           Description            |
	+============+========+==================================+
	| ``alerts`` | array  | A collection of alert messages.  |
	+------------+--------+----------------------------------+
	| ``>level`` | string | Success, info, warning or error. |
	+------------+--------+----------------------------------+
	| ``>text``  | string | Alert message.                   |
	+------------+--------+----------------------------------+

	**Response Example** ::

		{
		 "alerts": [
				{
					 "level": "success",
					 "text": "Check Extension loaded."
				}
		 ],
		}


|

**POST /api/1.1/to_extensions/:id/delete**

	Deletes a Traffic Ops extension.

	Authentication Required: Yes

	Role(s) Required: None

	**Request Route Parameters**

	+--------+----------+-----------------+
	|  Name  | Required |   Description   |
	+========+==========+=================+
	| ``id`` | yes      | TO extension id |
	+--------+----------+-----------------+

	**Response Properties**

	+------------+--------+----------------------------------+
	| Parameter  |  Type  |           Description            |
	+============+========+==================================+
	| ``alerts`` | array  | A collection of alert messages.  |
	+------------+--------+----------------------------------+
	| ``>level`` | string | Success, info, warning or error. |
	+------------+--------+----------------------------------+
	| ``>text``  | string | Alert message.                   |
	+------------+--------+----------------------------------+

	**Response Example** ::

		{
		 "alerts": [
				{
					 "level": "success",
					 "text": "Extension deleted."
				}
		 ],
		}


|

