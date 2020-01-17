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

.. role:: bash(code)
	:language: bash

.. _dev-debugging-ciab:

*****************************
Debugging inside CDN-in-a-Box
*****************************

Some CDN-in-a-Box components can be used with a debugger to step through lines of code, set breakpoints, see the state of all variables in each scope, etc. at runtime. Components that support debugging:

* `Traffic Monitor`_
* `Traffic Ops (Go)`_

Traffic Monitor
===============

* Navigate to the ``infrastructure/cdn-in-a-box`` directory. Remove the existing RPMs because they contain release Go binaries do not include useful debugging information. Rebuild the RPMs with no optimization, for debugging:

.. code-block:: shell
	:caption: Remove release RPMs, then build debug RPMs

	make very-clean
	make debug

* Still in ``infrastructure/cdn-in-a-box``, open ``variables.env`` and set ``TM_DEBUG_ENABLE`` to ``true``.

* Stop CDN-in-a-Box if it is running and remove any existing volumes. Rebuild the ``trafficmonitor`` image without reusing any cached layers to make sure it uses our fresh ``traffic_monitor.rpm``. Then, start CDN-in-a-Box.

.. code-block:: shell
	:caption: docker-compose command for debugging Traffic Monitor

	alias mydc='docker-compose -f docker-compose.yml -f docker-compose.expose-ports.yml optional/docker-compose.debugging.yml'
	mydc down -v
	mydc build --no-cache trafficmonitor
	mydc up

* Install `an IDE that supports delve <https://github.com/Microsoft/vscode-go/wiki/Debugging-Go-code-using-VS-Code>`_ and create a debugging configuration over port 2344.

* Use the debugging configuration you created to start debugging Traffic Monitor. It should connect without first breaking at any line.

For an example of usage, set a breakpoint at `the o.m.RLock() call in ThreadsafeEvents.Get() <https://github.com/apache/trafficcontrol/blob/RELEASE-4.0.0-RC3/traffic_monitor/health/event.go#L69>`_, then visit http://trafficmonitor.infra.ciab.test/publish/EventLog (see :ref:`Traffic Monitor APIs: /publish/EventLog <tm-publish-EventLog>`).

Traffic Ops (Go)
================

* Navigate to the ``infrastructure/cdn-in-a-box`` directory. Remove the existing RPMs because they contain release Go binaries do not include useful debugging information. Rebuild the RPMs with no optimization, for debugging:

.. code-block:: shell
	:caption: Remove release RPMs, then build debug RPMs

	make very-clean
	make debug

* Still in ``infrastructure/cdn-in-a-box``, open ``variables.env`` and set ``TO_DEBUG_ENABLE`` to ``true``.

* Stop CDN-in-a-Box if it is running and remove any existing volumes. Rebuild the ``trafficops-go`` image without reusing any cached layers to make sure it uses our fresh ``traffic_ops.rpm``. Then, start CDN-in-a-Box.

.. code-block:: shell
	:caption: docker-compose command for debugging Traffic Ops

	alias mydc='docker-compose -f docker-compose.yml -f docker-compose.expose-ports.yml optional/docker-compose.debugging.yml'
	mydc down -v
	mydc build --no-cache trafficops
	mydc up

* Install `an IDE that supports delve <https://github.com/Microsoft/vscode-go/wiki/Debugging-Go-code-using-VS-Code>`_ and create a debugging configuration over port 2345.

* Use the debugging configuration you created to start debugging Traffic Ops. It should connect without first breaking at any line.

For an example of usage, set a breakpoint at `the log.Debugln() call in TOProfile.Read() <https://github.com/apache/trafficcontrol/blob/RELEASE-4.0.0-RC3/traffic_ops/traffic_ops_golang/profile/profiles.go#L129>`_, then visit https://trafficportal.infra.ciab.test/api/1.5/profiles (after logging into :ref:`tp-overview`).

Troubleshooting
===============

* If you are debugging a Golang project and you don't see the values of all variables, or stepping to the next line puts you several lines ahead, rebuild the Docker image with an RPM built using :bash:`make debug`.
