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

Traffic Portal
**************

Introduction
============
Traffic Portal is an `AngularJS 1.x <https://angularjs.org/>`_ client served from a `Node.js <https://nodejs.org/en/>`_ web server designed to consume the Traffic Ops 1.x API. Traffic Portal provides a set of functionality restricted to the delivery service(s) of the authenticated user. Functionality primarily includes a set of charts / graphs designed to provide insight into the performance of a user's delivery service(s).

Software Requirements
=====================
To work on Traffic Portal you need a \*nix (MacOS and Linux are most commonly used) environment that has the following installed:

	* `Ruby 2.0.x or above <https://www.ruby-lang.org/en/>`_
	* `Compass 1.0.x or above <http://compass-style.org/>`_
	* `Node.js 6.0.x or above <https://nodejs.org/en/>`_
	* `Bower 1.7.9 or above <https://nodejs.org/en/>`_
	* `Grunt CLI 1.2.0 or above <https://github.com/gruntjs/grunt-cli>`_
	* Access to a working instance of Traffic Ops

Traffic Portal Project Tree Overview
=====================================
	* **traffic_control/traffic_portal/app/src** - contains HTML, JavaScript and Sass source files.

Installing The Traffic Portal Developer Environment
===================================================

	- Clone the traffic_control repository
	- Navigate to the traffic_control/traffic_portal of your cloned repository.
	- Run ``npm install`` to install application dependencies. Only needs to be done the first time unless traffic_portal/package.json changes.
	- Run ``bower install`` to install client-side dependencies. Only needs to be done the first time unless traffic_portal/bower.json changes.
	- Run ``grunt`` to package the application into traffic_portal/app/dist, start a local http server (Express), and start a file watcher.
	- Navigate to http://localhost:8080

Notes
=====

- The Traffic Portal consumes the Traffic Ops API. By default, Traffic Portal assumes Traffic Ops is running on http://localhost:3000. Temporarily modify conf/config.js if you need to change the location of Traffic Ops.



