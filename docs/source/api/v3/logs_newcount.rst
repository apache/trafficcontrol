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


.. _to-api-v3-logs-newcount:

*****************
``logs/newcount``
*****************

``GET``
=======
Gets the number of new changes made to the Traffic Control system - "new" being defined as the last time the client requested either :ref:`to-api-v3-logs`

.. note:: This endpoint's functionality is implemented by the :ref:`to-api-v3-logs` endpoint's response setting cookies for the client to use when requesting _this_ endpoint. Take care that your client respects cookies!

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
No parameters available

Response Structure
------------------
:newLogcount: The integer number of new changes

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Thu, 15 Nov 2018 15:17:35 GMT
	X-Server-Name: traffic_ops_golang/
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: Ugdqe8GXKSOExphwbDX/Gli+2vBpubttbpfYMbJaCP7adox3MzmVRi2RxTDL5kwPewrcL1CO88zGITskhOsc9g==
	Content-Length: 30

	{ "response": {
		"newLogcount": 4
	}}
