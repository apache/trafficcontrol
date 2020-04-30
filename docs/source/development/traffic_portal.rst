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

**************
Traffic Portal
**************

Introduction
============
Traffic Portal is an `AngularJS <https://angularjs.org/>`_ client served from a `Node.js <https://nodejs.org/en/>`_ web server designed to consume the :ref:`to-api`. Traffic Portal is the official replacement for the legacy Traffic Ops UI.

Software Requirements
=====================
To work on Traffic Portal you need a \*nix (MacOS and Linux are most commonly used) environment that has the following installed:

	* `Ruby Devel 2.0.x or above <https://www.rpmfind.net/linux/rpm2html/search.php?query=ruby-devel>`_
	* `Compass 1.0.x or above <http://compass-style.org/>`_
	* `Node.js 12.0.x or above <https://nodejs.org/en/>`_
	* `Bower 1.7.9 or above <https://www.npmjs.com/package/bower>`_
	* `Grunt CLI 1.2.0 or above <https://github.com/gruntjs/grunt-cli>`_
	* Access to a working instance of Traffic Ops

.. note:: The Traffic Portal consumes the Traffic Ops API. Modify traffic_portal/conf/config.js to specify the location of Traffic Ops.

Traffic Portal Project Tree Overview
=====================================
* **traffic_control/traffic_portal/app/src** - contains HTML, JavaScript and :abbr:`SCSS (Sassy CSS)` source files.

Installing The Traffic Portal Developer Environment
===================================================
#. Clone the `Traffic Control Repository <https://github.com/apache/trafficcontrol>`_
#. Navigate to the ``traffic_portal`` subdirectory of your cloned repository.
#. Run ``npm install`` to install application dependencies into ``traffic_portal/node_modules``. Only needs to be done the first time unless ``traffic_portal/package.json`` changes.
#. Run ``bower install`` to install client-side dependencies into ``traffic_portal/app/bower_components``. Only needs to be done the first time unless ``traffic_portal/bower.json`` changes.
#. Run ``grunt`` to package the application into ``traffic_portal/app/dist``, start a local HTTPS server (Express), and start a file watcher.
#. Modify ``traffic_portal/conf/config.js``:

	#. Valid SSL certificates and keys are needed for Traffic Portal to run. Generate these (e.g. using `this SuperUser answer <https://superuser.com/questions/226192/avoid-password-prompt-for-keys-and-prompts-for-dn-information#answer-226229>`_) and update ``ssl``.
	#. Modify ``api.base_url`` to point to your Traffic Ops API endpoint.
	#. Modify ``files.static`` to be ``./app/dist/public``.
	#. Modify ``log.stream`` to be ``./server/log/access.log``.

#. Navigate to http(s)://localhost:[port|sslPort defined in ``traffic_portal/conf/config.js``]
