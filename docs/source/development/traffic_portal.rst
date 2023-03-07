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

.. _dev-traffic-portal:

**************
Traffic Portal
**************

Introduction
============
Traffic Portal is an `AngularJS <https://angularjs.org/>`_ client served from a `Node.js <https://nodejs.org/en/>`_ web server designed to consume the :ref:`to-api`. Traffic Portal is the official replacement for the legacy Traffic Ops UI.

.. _dev-tp-software-requirements:

Software Requirements
=====================
To work on Traffic Portal you need a \*nix (MacOS and Linux are most commonly used) environment that has the following installed:

* `Node.js 16.0.x or above <https://nodejs.org/en/>`_
* `Grunt CLI 1.2.0 or above <https://github.com/gruntjs/grunt-cli>`_
* Access to a working instance of Traffic Ops

.. _dev-tp-global-npm:

Install Global NPM Packages
---------------------------

Grunt CLI can be installed using NPM.

.. code-block:: shell
	:caption: Install Grunt CLI

	npm -g install grunt-cli


Traffic Portal Project Tree Overview
=====================================
* **traffic_control/traffic_portal/app/src** - contains HTML, JavaScript and :abbr:`SCSS (Sassy CSS)` source files.

Installing The Traffic Portal Developer Environment
===================================================
#. Clone the `Traffic Control Repository <https://github.com/apache/trafficcontrol>`_
#. Navigate to the ``traffic_portal`` subdirectory of your cloned repository.
#. Run ``npm install`` to install application dependencies into ``traffic_portal/node_modules``. Only needs to be done the first time unless ``traffic_portal/package.json`` changes.

#. Modify :atc-file:`traffic_portal/conf/configDev.js`:

	#. Valid SSL certificates and keys are needed for Traffic Portal to run. Generate these (e.g. using `this SuperUser answer <https://superuser.com/questions/226192/avoid-password-prompt-for-keys-and-prompts-for-dn-information#answer-226229>`_) and update ``ssl``.
	#. Modify ``api.base_url`` to point to your Traffic Ops API endpoint.

#. Run ``grunt`` to package the application into ``traffic_portal/app/dist``, start a local HTTPS server (Express), and start a file watcher. To use a custom configuration file (not just :atc-file:`traffic_portal/conf/config.js` or :atc-file:`traffic_portal/conf/configDev.js`), set the `TP_SERVER_CONFIG_FILE` environment variable to the location of the desired file.
#. Navigate to http(s)://localhost:[port|sslPort defined in the configuration file used (default: :atc-file:`traffic_portal/conf/configDev.js`)]
