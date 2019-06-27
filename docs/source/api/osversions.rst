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

.. _to-api-osversions:

**************
``osversions``
**************
.. seealso:: :ref:`tp-tools-generate-iso`

``GET``
=======
Gets all available Operating System (OS) versions for ISO generation, as well as the name of the directory where the "kickstarter" files are found.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
No parameters available.

Response Structure
------------------
This endpoint has no constant keys in its ``response``. Instead, each key in the ``response`` object is the name of an OS, and the value is a string that names the directory where the ISO source can be found. These directories sit under `/var/www/files/` on the Traffic Ops host machine by default, or at the location defined by the ``kickstart.files.location`` :term:`Parameter` of the Traffic Ops server's :term:`Profile`, if it is defined.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Fri, 30 Nov 2018 19:14:36 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Fri, 30 Nov 2018 23:14:36 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: RxbRY2DZ+lYOdTzzUETEZ3wtLBiD2BwXMVuaZjhe4a4cwgcZKRBWxZ6Qy5YYujFe1+UBiTG4sML/Amn27F4AVg==
	Content-Length: 38

	{ "response": {
		"CentOS 7.2": "centos72"
	}}
