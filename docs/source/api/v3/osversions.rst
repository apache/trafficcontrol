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

.. _to-api-v3-osversions:

**************
``osversions``
**************
.. seealso:: :ref:`tp-tools-generate-iso`

``GET``
=======
Gets all available :abbr:`OS (Operating System)` versions for ISO generation, as well as the name of the directory where the "kickstarter" files are found.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
No parameters available.

.. _v3-response-structure:

Response Structure
------------------
This endpoint has no constant keys in its ``response``. Instead, each key in the ``response`` object is the name of an OS, and the value is a string that names the directory where the ISO source can be found. These directories sit under ``/var/www/files/`` on the Traffic Ops host machine by default, or at the location defined by the ``kickstart.files.location`` :term:`Parameter` of the Traffic Ops server's :term:`Profile`, if it is defined.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: RxbRY2DZ+lYOdTzzUETEZ3wtLBiD2BwXMVuaZjhe4a4cwgcZKRBWxZ6Qy5YYujFe1+UBiTG4sML/Amn27F4AVg==
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 30 Nov 2018 19:14:36 GMT
	Content-Length: 38

	{ "response": {
		"CentOS 7.2": "centos72"
	}}


Configuration File
------------------
The data returned from the endpoint comes directly from a configuration file. By default, the file is located at ``/var/www/files/osversions.json``.
The **directory** of the file can be changed by creating a specific :term:`Parameter` named ``kickstart.files.location`` in configuration file ``mkisofs``.

The format of the file is a JSON object as described in :ref:`v3-response-structure`.

.. code-block:: json
	:caption: Example osversions.json file

	{
		"CentOS 7.2": "centos72"
	}


The legacy Perl Traffic Ops used a Perl configuration file located by default at ``/var/www/files/osversions.cfg``. A Perl script is provided
to convert the legacy configuration file to the new JSON format. The script is located within the Traffic Control repository at ``traffic_ops/app/bin/osversions-convert.pl``.

.. code-block:: shell
	:caption: Example usage of conversion script

	./osversions-convert.pl < /var/www/files/osversions.cfg > /var/www/files/osversions.json
