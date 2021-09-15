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

	* `Ruby Devel 2.0.x or above <https://www.rpmfind.net/linux/rpm2html/search.php?query=ruby-devel>`_
	* `Compass 1.0.x or above <http://compass-style.org/>`_
	* `Node.js 12.0.x or above <https://nodejs.org/en/>`_
	* `Grunt CLI 1.2.0 or above <https://github.com/gruntjs/grunt-cli>`_
	* Access to a working instance of Traffic Ops

.. _dev-tp-global-npm:

Install Global NPM Packages
---------------------------

Grunt CLI can be installed using NPM.

.. code-block:: shell
	:caption: Install Grunt CLI

	npm -g install grunt-cli

.. _dev-tp-compass:

Install Compass
---------------

Compass can be installed using ``gem`` manually, or by using ``bundle``

.. tip:: Bundle will automatically install the correct version of the gems.

#. ``brew install ruby``/``apt-get install ruby``/``yum install ruby``

#. ``gem update --system``

#. At this point, you can either manually install the gems or use bundler

	#. For manually: ``gem install sass compass``

	#. For automatically: ``gem install bundle && bundle install``

	.. note:: Bundle requires ruby versions > 2.3.0, so if you're using a version of ruby < 2.3.0 then this will not work.

#. Make sure that ``compass`` and ``sass`` are part of your ``PATH`` environment variable.

#. If not, you can see where gem installs ``compass`` and ``sass`` by running:
	``gem environment``

#. In there, you can see where ruby is installing all the gems. Add that path to your ``PATH`` environment variable.
	For example, it is ``/usr/local/lib/ruby/gems/2.7.0/gems/compass-1.0.3/bin/`` for this test setup.

#. Once you have installed ``compass`` successfully, make sure you can reach it by typing:
	``compass version``
	This should give a valid output. For example, for the test setup, the output is:

.. code-block:: text
	:caption: Compass version output

	Compass 1.0.3 (Polaris)
	Copyright (c) 2008-2020 Chris Eppstein
	Released under the MIT License.
	Compass is charityware.
	Please make a tax deductable donation for a worthy cause: http://umdf.org/compass


Traffic Portal Project Tree Overview
=====================================
* **traffic_control/traffic_portal/app/src** - contains HTML, JavaScript and :abbr:`SCSS (Sassy CSS)` source files.

Installing The Traffic Portal Developer Environment
===================================================
#. Clone the `Traffic Control Repository <https://github.com/apache/trafficcontrol>`_
#. Navigate to the ``traffic_portal`` subdirectory of your cloned repository.
#. Run ``npm install`` to install application dependencies into ``traffic_portal/node_modules``. Only needs to be done the first time unless ``traffic_portal/package.json`` changes.
#. Make sure that compass is installed and functioning correctly by running ``compass version``. If compass is not available, then it can be installed following the instructions under :ref:`dev-tp-compass`.

#. Modify ``traffic_portal/conf/configDev.js``:
	#. Valid SSL certificates and keys are needed for Traffic Portal to run. Generate these (e.g. using `this SuperUser answer <https://superuser.com/questions/226192/avoid-password-prompt-for-keys-and-prompts-for-dn-information#answer-226229>`_) and update ``ssl``.
	#. Modify ``api.base_url`` to point to your Traffic Ops API endpoint.
#. Run ``grunt`` to package the application into ``traffic_portal/app/dist``, start a local HTTPS server (Express), and start a file watcher.
#. Navigate to http(s)://localhost:[port|sslPort defined in ``traffic_portal/conf/configDev.js``]

.. note:: The Traffic Portal consumes the Traffic Ops API. Modify traffic_portal/conf/configDev.js to specify the location of Traffic Ops.
