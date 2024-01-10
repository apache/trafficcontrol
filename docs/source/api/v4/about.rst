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

.. _to-api-v4-about:

***********
``about``
***********

``GET``
=======

Returns info about the Traffic Ops build that is currently running, generated at startup. The output will be the same until the Traffic Ops :ref:`version changes <to-upgrading>`.

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: None
:Response Type:  Object

Request Structure
-----------------
No parameters available.

.. code-block:: http
	:caption: Request Example

	GET /api/4.0/about HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...

Response Structure
------------------
:commitHash:    The `Git <https://git-scm.com/>`_ commit hash that Traffic Ops was built at.
:commits:       The number of commits in the branch of the commit that Traffic Ops was built at, including that commit. Calculated by extracting the commit count from running ``git describe --tags --long``.
:goVersion:     The version of `Go <https://golang.org/>`_ that was used to build Traffic Ops.
:release:       The major version of CentOS or Red Hat Enterprise Linux that the build environment was running.
:name:          The human-readable name of the `RPM <https://rpm-packaging-guide.github.io/#packaging-software>`_ file.
:RPMVersion:    The entire name of the RPM file, excluding the file extension.
:Version:       The version of :abbr:`ATC (Apache Traffic Control)` that this version of Traffic Control belongs to.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 24 Feb 2020 19:35:28 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: 7SVQsddCUVRs+sineziRGR6OyMli7XLZbjxyMQgW6E506bh5thMOuttPFT7aJckDcgT45PlhexycwlApOHI4Vw==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 24 Feb 2020 18:35:28 GMT
	Content-Length: 145

	{
		"commitHash": "1c9a2e9c",
		"commits": "10555",
		"goVersion": "go1.11.13",
		"release": "el7",
		"name": "traffic_ops",
		"RPMVersion": "traffic_ops-4.0.0-10555.1c9a2e9c.el7",
		"Version": "4.0.0"
	}
