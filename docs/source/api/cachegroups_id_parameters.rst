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

.. _to-api-cachegroup_id_parameters:

*********************************
``cachegroups/{{ID}}/parameters``
*********************************
Gets all the parameters associated with a Cache Group

.. seealso:: :ref:`param-prof`

``GET``
=======
:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------------------+----------+-----------------------+
	|       Name       | Required | Description           |
	+==================+==========+=======================+
	| ``id``           | yes      | Cache Group ID        |
	+------------------+----------+-----------------------+


Response Structure
------------------
:configFile:  Configuration file associated with the parameter
:id:          A numeric, unique identifier for this parameter
:lastUpdated: The Time / Date this entry was last updated
:name:        Name of the parameter
:secure:      Is the parameter value only visible to admin users
:value:       Value of the parameter

.. code-block:: json
	:caption: Response Example

	{ "response": [
		{
			"lastUpdated": "2015-08-27 21:11:49+00",
			"value": "http://localhost:8088/teakCatalog/GOID_15101A.csv",
			"secure": false,
			"name": "teakcluster.TeakCatalogUrl",
			"id": 1409,
			"configFile": "to_ext_teak.config"
		},
		{
			"lastUpdated": "2015-09-29 18:59:47+00",
			"value": "http://localhost:8088/teakCatalog/GOID_15101A.csv",
			"secure": false,
			"name": "teakcluster.TeakCatalogUrl",
			"id": 2351,
			"configFile": "teak.config"
		},
		{
			"lastUpdated": "2018-09-18 20:33:34.182367+00",
			"value": "172.24.74.70",
			"secure": false,
			"name": "teakcluster.LBVIP",
			"id": 6571,
			"configFile": "to_ext_teak.config"
		},
		{
			"lastUpdated": "2018-09-18 20:36:13.227886+00",
			"value": "172.24.74.70",
			"secure": false,
			"name": "teakcluster.LBVIP",
			"id": 6572,
			"configFile": "teak.config"
		},
		{
			"lastUpdated": "2018-09-27 21:09:10.036367+00",
			"value": "ccdn-ats-tk-15101-01:ccdn-ats-tk-15101-02:ccdn-ats-tk-15101-03:ccdn-ats-tk-15101-04:ccdn-ats-tk-15101-05:ccdn-ats-tk-15101-06:ccdn-ats-tk-15101-07",
			"secure": false,
			"name": "node_order",
			"id": 6574,
			"configFile": "to_ext_teak.config"
		},
		{
			"lastUpdated": "2015-08-27 21:11:49+00",
			"value": "10.42.18.196",
			"secure": false,
			"name": "cgw.cgwServer",
			"id": 1148,
			"configFile": "to_ext_teak.config"
		},
		{
			"lastUpdated": "2015-08-27 21:11:49+00",
			"value": "80",
			"secure": false,
			"name": "cgw.cgwPort",
			"id": 1101,
			"configFile": "to_ext_teak.config"
		}
	]}
