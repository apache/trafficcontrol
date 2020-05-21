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
#. Make sure that compass is installed and functioning correctly by running ``compass version``. If ``compass`` is not running, you can install it on macOS as follows:
    #. macOS comes with its own version of ruby built into it. In order to install compass, if you run a command like ``sudo gem install compass``, it will take the default ruby available in mac. In order to properly install compass, here is what you should be doing:

    #. ``brew install ruby``

    #. ``gem update --system``

    #. ``gem install compass``

        .. tip:: You need to install ``sass`` before you install ``compass``. You can do that by running:
                 ``gem install sass``

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

#. Run ``grunt`` to package the application into ``traffic_portal/app/dist``, start a local HTTPS server (Express), and start a file watcher.
#. Modify ``traffic_portal/conf/config.js``:

	#. Valid SSL certificates and keys are needed for Traffic Portal to run. Generate these (e.g. using `this SuperUser answer <https://superuser.com/questions/226192/avoid-password-prompt-for-keys-and-prompts-for-dn-information#answer-226229>`_) and update ``ssl``.
	#. Modify ``api.base_url`` to point to your Traffic Ops API endpoint.
	#. Modify ``files.static`` to be ``./app/dist/public``.
	#. Modify ``log.stream`` to be ``./server/log/access.log``.

#. Navigate to http(s)://localhost:[port|sslPort defined in ``traffic_portal/conf/config.js``]
