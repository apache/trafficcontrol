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

.. _to-api-servercheck_extensions-id:

*********************************
``servercheck/extensions/{{ID}}``
*********************************

``DELETE``
==========
Deletes a Traffic Ops server check extension definition. This does **not** delete the actual extension file.

:Auth. Required: Yes
:Roles Required: None\ [1]_
:Permissions Required: SERVER-CHECK:DELETE, SERVER-CHECK:READ, SERVER:READ
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------------------------------------+
	| Name | Description                                                               |
	+======+===========================================================================+
	|  ID  | The integral, unique identifier of the extension definition to be deleted |
	+------+---------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/5.0/servercheck/extensions/16 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 12 Dec 2018 16:33:52 GMT
	X-Server-Name: traffic_ops_golang/
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: EB0Nu85azbGzaehDTAODP3NPqWbByIza1XQhgwtsW2WTXyK/dxQtncp0YiJXyO0tH9H+n+6BBfojBOb5h0dFPA==
	Content-Length: 60

	{ "alerts": [
		{
			"level": "success",
			"text": "Extension deleted."
		}
	]}

.. [1] No roles are required to use this endpoint, however access is controlled by username. Only the reserved user ``extension`` is permitted the use of this endpoint.
