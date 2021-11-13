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

.. _to-api-v3-cachegroups-id-parameters:

*********************************
``cachegroups/{{ID}}/parameters``
*********************************

.. deprecated:: ATCv6

``GET``
=======
Gets all of a :ref:`Cache Group's parameters <cache-group-parameters>`.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| Name        | Required | Description                                                                                                   |
	+=============+==========+===============================================================================================================+
	| parameterId | no       | Show only the :term:`Parameter` with the given :ref:`parameter-id`                                            |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| orderby     | no       | Choose the ordering of the results - must be the name of one of the fields of the objects in the ``response`` |
	|             |          | array                                                                                                         |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| sortOrder   | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")                      |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| limit       | no       | Choose the maximum number of results to return                                                                |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| offset      | no       | The number of results to skip before beginning to return results. Must use in conjunction with limit          |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| page        | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long   |
	|             |          | and the first page is 1. If ``offset`` was defined, this query parameter has no effect. ``limit`` must be     |
	|             |          | defined to make use of ``page``.                                                                              |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+

.. table:: Request Path Parameters

	+-----------+----------------------------------------------------------+
	| Parameter | Description                                              |
	+===========+==========================================================+
	| ID        | The :ref:`cache-group-id` of a :term:`Cache Group`       |
	+-----------+----------------------------------------------------------+


Response Structure
------------------
:configFile:  The :term:`Parameter`'s :ref:`parameter-config-file`
:id:          The :term:`Parameter`'s :ref:`parameter-id`
:lastUpdated: The date and time at which this :term:`Parameter` was last updated, in :ref:`non-rfc-datetime`
:name:        :ref:`parameter-name` of the :term:`Parameter`
:secure:      A boolean value describing whether or not the :term:`Parameter` is :ref:`parameter-secure`
:value:       The :term:`Parameter`'s :ref:`parameter-value`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 14 Nov 2018 19:56:23 GMT
	X-Server-Name: traffic_ops_golang/
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: DfqPtySzVMpnBYqVt/45sSRG/1pRTlQdIcYuQZ0CQt79QSHLzU5e4TbDqht6ntvNP041LimKsj5RzPlPX1n6tg==
	Content-Length: 135

	{ "response": [
		{
			"lastUpdated": "2018-11-14 18:22:43.754786+00",
			"value": "foobar",
			"secure": false,
			"name": "foo",
			"id": 124,
			"configFile": "bar"
		}
	]}
