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

.. _to-api-v1-types-trimmed:

*****************
``types/trimmed``
*****************
.. deprecated:: ATCv4
	This endpoint and all of its functionality is deprecated. All of the information it can return can be more completely obtained with :ref:`to-api-v1-types`.

``GET``
=======
Retrieves only the names of all of the :term:`Types` of things configured in Traffic Ops. Yes, that is as specific as a description of a 'type' can be.

.. warning:: This endpoint is of limited use because it doesn't tell you what the type of each :term:`Type` is, which describes the types of objects that it can describe. No, I did not just have a stroke while writing this.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
No parameters available

Response Structure
------------------
:name: The name of the type

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Encoding: gzip
	Content-Length: 389
	Content-Type: application/json
	Date: Fri, 31 Jan 2020 18:09:29 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Fri, 31 Jan 2020 22:09:29 GMT; path=/; HttpOnly
	Vary: Accept-Encoding

	{ "alerts": [
		{
			"level": "warning",
			"text": "This endpoint is deprecated, please use '/types' instead"
		}
	],
	"response": [
		{
			"name": "AAAA_RECORD"
		},
		{
			"name": "ANY_MAP"
		},
		{
			"name": "A_RECORD"
		},
		{
			"name": "BIND"
		},
		{
			"name": "CCR"
		},
		{
			"name": "CHECK_EXTENSION_BOOL"
		},
		{
			"name": "CHECK_EXTENSION_NUM"
		},
		{
			"name": "CHECK_EXTENSION_OPEN_SLOT"
		},
		{
			"name": "CLIENT_STEERING"
		},
		{
			"name": "CNAME_RECORD"
		},
		{
			"name": "CONFIG_EXTENSION"
		},
		{
			"name": "DNS"
		},
		{
			"name": "DNS_LIVE"
		},
		{
			"name": "DNS_LIVE_NATNL"
		},
		{
			"name": "EDGE"
		},
		{
			"name": "EDGE_LOC"
		},
		{
			"name": "ENROLLER"
		},
		{
			"name": "GRAFANA"
		},
		{
			"name": "HEADER_REGEXP"
		},
		{
			"name": "HOST_REGEXP"
		},
		{
			"name": "HTTP"
		},
		{
			"name": "HTTP_LIVE"
		},
		{
			"name": "HTTP_LIVE_NATNL"
		},
		{
			"name": "HTTP_NO_CACHE"
		},
		{
			"name": "INFLUXDB"
		},
		{
			"name": "MID"
		},
		{
			"name": "MID_LOC"
		},
		{
			"name": "ORG"
		},
		{
			"name": "ORG_LOC"
		},
		{
			"name": "PATH_REGEXP"
		},
		{
			"name": "RASCAL"
		},
		{
			"name": "RESOLVE4"
		},
		{
			"name": "RESOLVE6"
		},
		{
			"name": "RIAK"
		},
		{
			"name": "STATISTIC_EXTENSION"
		},
		{
			"name": "STEERING"
		},
		{
			"name": "STEERING_GEO_ORDER"
		},
		{
			"name": "STEERING_GEO_WEIGHT"
		},
		{
			"name": "STEERING_ORDER"
		},
		{
			"name": "STEERING_REGEXP"
		},
		{
			"name": "STEERING_WEIGHT"
		},
		{
			"name": "TC_LOC"
		},
		{
			"name": "TRAFFIC_ANALYTICS"
		},
		{
			"name": "TRAFFIC_OPS"
		},
		{
			"name": "TRAFFIC_OPS_DB"
		},
		{
			"name": "TRAFFIC_PORTAL"
		},
		{
			"name": "TRAFFIC_STATS"
		},
		{
			"name": "TR_LOC"
		},
		{
			"name": "TXT_RECORD"
		}
	]}
