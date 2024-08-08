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

.. tip:: For the purposes of development, it may be easier to use :ref:`dev`.

Some CDN-in-a-Box components can be used with a debugger to step through lines of code, set breakpoints, see the state of all variables in each scope, etc. at runtime. Components that support debugging:

* `Enroller`_
* `t3c on Caches`_
	- `t3c on Edge Cache`_
	- `t3c on Mid 01 Cache`_
	- `t3c on Mid 02 Cache`_
* `Traffic Monitor`_
* `Traffic Ops`_
* `Traffic Router`_
* `Traffic Stats`_

Enroller
========

* In ``infrastructure/cdn-in-a-box``, open ``variables.env`` and set ``ENROLLER_DEBUG_ENABLE`` to ``true``.

* Stop CDN-in-a-Box if it is running and remove any existing volumes. Build/rebuild the ``enroller-debug`` image each time you have changed :atc-file:`infrastructure/cdn-in-a-box/enroller/enroller.go`. Then, start CDN-in-a-Box.

	.. code-block:: shell
		:caption: docker compose command for debugging the CDN in a Box Enroller

		alias mydc='docker compose -f docker-compose.yml -f docker-compose.expose-ports.yml -f optional/docker-compose.debugging.yml'
		mydc down -v
		mydc build enroller
		mydc up

* Install `an IDE that supports delve <https://github.com/go-delve/delve/blob/master/Documentation/EditorIntegration.md>`_ and create a debugging configuration over port 2343. If you are using VS Code, the configuration should look like this:

	.. code-block:: json
		:caption: VS Code launch.json for debugging the CDN in a Box Enroller

		{
			"version": "0.2.0",
			"configurations": [
				{
					"name": "Enroller",
					"type": "go",
					"request": "attach",
					"mode": "remote",
					"port": 2343,
					"cwd": "${workspaceRoot}/",
					"remotePath": "/go/src/github.com/apache/trafficcontrol/",
				}
			]
		}

* Use the debugging configuration you created to start debugging the Enroller. It should connect without first breaking at any line.

For an example of usage, set a breakpoint at `the toSession.CreateDeliveryServiceV30() call in enrollDeliveryService() <https://github.com/apache/trafficcontrol/blob/RELEASE-5.1.1/infrastructure/cdn-in-a-box/enroller/enroller.go#L209>`_, then wait for the Enroller to process a file from ``/shared/enroller/deliveryservices/`` (only exists within the Docker container).

t3c on Caches
=============

t3c on Edge Cache
-----------------
* Navigate to the ``infrastructure/cdn-in-a-box`` directory. Remove ``infrastructure/cdn-in-a-box/cache/trafficcontrol-cache-config.rpm`` because it contains release Go binaries that do not include useful debugging information. Rebuild the RPM with no optimization, for debugging:

	.. code-block:: shell
		:caption: Remove release RPMs, then build debug RPMs

		make very-clean
		make debug cache/trafficcontrol-cache-config.rpm

	.. tip:: If you have gone through the steps to :ref:`dev-building-natively`, you can run ``make debug native cache/trafficcontrol-cache-config.rpm`` instead of ``make debug cache/trafficcontrol-cache-config.rpm`` to build the RPM quickly.

* Still in ``infrastructure/cdn-in-a-box``, open ``variables.env`` and set ``T3C_DEBUG_COMPONENT_EDGE`` to ``t3c-apply`` (used for this example). A list of valid values for ``T3C_DEBUG_COMPONENT_EDGE``:
	- t3c-apply
	- t3c-check
	- t3c-check-refs
	- t3c-check-reload
	- t3c-diff
	- t3c-generate
	- t3c-request
	- t3c-update

* Stop CDN-in-a-Box if it is running and remove any existing volumes. Build the ``edge-debug`` image to make sure it uses our fresh ``trafficcontrol-cache-config.rpm``. Then, start CDN-in-a-Box:

	.. code-block:: shell
		:caption: docker compose command for debugging ``t3c`` running on the Edge Cache

		alias mydc='docker compose -f docker-compose.yml -f docker-compose.expose-ports.yml -f optional/docker-compose.debugging.yml'
		mydc down -v
		mydc build edge
		mydc up -d
		mydc logs -f trafficmonitor

* Install `an IDE that supports delve <https://github.com/go-delve/delve/blob/master/Documentation/EditorIntegration.md>`_ and create a debugging configuration over port 2347. If you are using VS Code, the configuration should look like this:

	.. code-block:: json
		:caption: VS Code launch.json for debugging ``t3c`` on the Edge Cache

		{
			"version": "0.2.0",
			"configurations": [
				{
					"name": "t3c on Edge",
					"type": "go",
					"request": "attach",
					"mode": "remote",
					"port": 2347,
					"cwd": "${workspaceRoot}",
					"remotePath": "/tmp/go/src/github.com/apache/trafficcontrol",
				}
			]
		}

Wait for Traffic Monitor to start, which will indicate that the SSL keys have been generated. Because ``T3C_DEBUG_COMPONENT_EDGE`` is set to the name of one of the ``t3c`` binaries, ``t3c`` will *not* run automatically every minute. Start it it manually:

.. code-block:: shell
	:caption: Run ``t3c-apply`` with debugging enabled

	[user@computer cdn-in-a-box]$ mydc exec edge t3c apply --run-mode=badass --traffic-ops-url=https://trafficops.infra.ciab.test --traffic-ops-user=admin --traffic-ops-password=twelve12 --git=yes --dispersion=0 --log-location-error=stdout --log-location-warning=stdout --log-location-info=stdout all
	API server listening at: [::]:2347

The *API server listening* message is from ``dlv``, indicating it is ready to accept a connection from your IDE. Note that, unlike the other components, execution of ``t3c`` does not begin until your IDE connects to ``dlv``.

For this example, set a breakpoint at `the assignment of "##OVERRIDDEN## " + str to newstr in torequest.processRemapOverrides() <https://github.com/apache/trafficcontrol/blob/dde7f69d49/cache-config/t3c-apply/torequest/torequest.go#L336>`_.

Use the debugging configuration you created to connect to ``dlv`` and start debugging ``t3c``.

t3c on Mid 01 Cache
-------------------
* Navigate to the ``infrastructure/cdn-in-a-box`` directory. Remove ``infrastructure/cdn-in-a-box/cache/trafficcontrol-cache-config.rpm`` because it contains release Go binaries that do not include useful debugging information. Rebuild the RPM with no optimization, for debugging:

	.. code-block:: shell
		:caption: Remove release RPMs, then build debug RPMs

		make very-clean
		make debug cache/trafficcontrol-cache-config.rpm

	.. tip:: If you have gone through the steps to :ref:`dev-building-natively`, you can run ``make debug native cache/trafficcontrol-cache-config.rpm`` instead of ``make debug cache/trafficcontrol-cache-config.rpm`` to build the RPM quickly.

* Still in ``infrastructure/cdn-in-a-box``, open ``variables.env`` and set ``T3C_DEBUG_COMPONENT_MID_01`` to ``t3c-apply`` (used for this example). A list of valid values for ``T3C_DEBUG_COMPONENT_MID_01``:
	- t3c-apply
	- t3c-check
	- t3c-check-refs
	- t3c-check-reload
	- t3c-diff
	- t3c-generate
	- t3c-request
	- t3c-update

* Stop CDN-in-a-Box if it is running and remove any existing volumes. Build the ``mid-debug`` image to make sure it uses our fresh ``trafficcontrol-cache-config.rpm``. Then, start CDN-in-a-Box:

	.. code-block:: shell
		:caption: docker compose command for debugging ``t3c`` running on the Mid 01 Cache

		alias mydc='docker compose -f docker-compose.yml -f docker-compose.expose-ports.yml -f optional/docker-compose.debugging.yml'
		mydc down -v
		mydc build mid-01
		mydc up -d
		mydc logs -f trafficmonitor

* Install `an IDE that supports delve <https://github.com/go-delve/delve/blob/master/Documentation/EditorIntegration.md>`_ and create a debugging configuration over port 2348. If you are using VS Code, the configuration should look like this:

	.. code-block:: json
		:caption: VS Code launch.json for debugging ``t3c`` on the Mid 01 Cache

		{
			"version": "0.2.0",
			"configurations": [
				{
					"name": "t3c on Mid 01",
					"type": "go",
					"request": "attach",
					"mode": "remote",
					"port": 2348,
					"cwd": "${workspaceRoot}",
					"remotePath": "/tmp/go/src/github.com/apache/trafficcontrol",
				}
			]
		}

Wait for Traffic Monitor to start, which will indicate that the SSL keys have been generated. Because ``T3C_DEBUG_COMPONENT_MID_01`` is set to the name of one of the ``t3c`` binaries, ``t3c`` will *not* run automatically every minute. Start it it manually:

.. code-block:: shell
	:caption: Run ``t3c-apply`` with debugging enabled

	[user@computer cdn-in-a-box]$ mydc exec mid-01 t3c apply --run-mode=badass --traffic-ops-url=https://trafficops.infra.ciab.test --traffic-ops-user=admin --traffic-ops-password=twelve12 --git=yes --dispersion=0 --log-location-error=stdout --log-location-warning=stdout --log-location-info=stdout all
	API server listening at: [::]:2348

The *API server listening* message is from ``dlv``, indicating it is ready to accept a connection from your IDE. Note that, unlike the other components, execution of ``t3c`` does not begin until your IDE connects to ``dlv``.

For this example, set a breakpoint at `the assignment of "##OVERRIDDEN## " + str to newstr in torequest.processRemapOverrides() <https://github.com/apache/trafficcontrol/blob/dde7f69d49/cache-config/t3c-apply/torequest/torequest.go#L336>`_.

Use the debugging configuration you created to connect to ``dlv`` and start debugging ``t3c``.

t3c on Mid 02 Cache
-------------------
* Navigate to the ``infrastructure/cdn-in-a-box`` directory. Remove ``infrastructure/cdn-in-a-box/cache/trafficcontrol-cache-config.rpm`` because it contains release Go binaries that do not include useful debugging information. Rebuild the RPM with no optimization, for debugging:

	.. code-block:: shell
		:caption: Remove release RPMs, then build debug RPMs

		make very-clean
		make debug cache/trafficcontrol-cache-config.rpm

	.. tip:: If you have gone through the steps to :ref:`dev-building-natively`, you can run ``make debug native cache/trafficcontrol-cache-config.rpm`` instead of ``make debug cache/trafficcontrol-cache-config.rpm`` to build the RPM quickly.

* Still in ``infrastructure/cdn-in-a-box``, open ``variables.env`` and set ``T3C_DEBUG_COMPONENT_MID_02`` to ``t3c-apply`` (used for this example). A list of valid values for ``T3C_DEBUG_COMPONENT_MID_02``:
	- t3c-apply
	- t3c-check
	- t3c-check-refs
	- t3c-check-reload
	- t3c-diff
	- t3c-generate
	- t3c-request
	- t3c-update

* Stop CDN-in-a-Box if it is running and remove any existing volumes. Build the ``mid-debug`` image to make sure it uses our fresh ``trafficcontrol-cache-config.rpm``. Then, start CDN-in-a-Box:

	.. code-block:: shell
		:caption: docker compose command for debugging ``t3c`` running on the Mid 02 Cache

		alias mydc='docker compose -f docker-compose.yml -f docker-compose.expose-ports.yml -f optional/docker-compose.debugging.yml'
		mydc down -v
		mydc build mid-02
		mydc up -d
		mydc logs -f trafficmonitor

* Install `an IDE that supports delve <https://github.com/go-delve/delve/blob/master/Documentation/EditorIntegration.md>`_ and create a debugging configuration over port 2349. If you are using VS Code, the configuration should look like this:

	.. code-block:: json
		:caption: VS Code launch.json for debugging ``t3c`` on the Mid 02 Cache

		{
			"version": "0.2.0",
			"configurations": [
				{
					"name": "t3c on Mid 02",
					"type": "go",
					"request": "attach",
					"mode": "remote",
					"port": 2349,
					"cwd": "${workspaceRoot}",
					"remotePath": "/tmp/go/src/github.com/apache/trafficcontrol",
				}
			]
		}

Wait for Traffic Monitor to start, which will indicate that the SSL keys have been generated. Because ``T3C_DEBUG_COMPONENT_MID_02`` is set to the name of one of the ``t3c`` binaries, ``t3c`` will *not* run automatically every minute. Start it it manually:

.. code-block:: shell
	:caption: Run ``t3c-apply`` with debugging enabled

	[user@computer cdn-in-a-box]$ mydc exec mid-02 t3c apply --run-mode=badass --traffic-ops-url=https://trafficops.infra.ciab.test --traffic-ops-user=admin --traffic-ops-password=twelve12 --git=yes --dispersion=0 --log-location-error=stdout --log-location-warning=stdout --log-location-info=stdout all
	API server listening at: [::]:2349

The *API server listening* message is from ``dlv``, indicating it is ready to accept a connection from your IDE. Note that, unlike the other components, execution of ``t3c`` does not begin until your IDE connects to ``dlv``.

For this example, set a breakpoint at `the assignment of "##OVERRIDDEN## " + str to newstr in torequest.processRemapOverrides() <https://github.com/apache/trafficcontrol/blob/dde7f69d49/cache-config/t3c-apply/torequest/torequest.go#L336>`_.

Use the debugging configuration you created to connect to ``dlv`` and start debugging ``t3c``.

Traffic Monitor
===============

* Navigate to the ``infrastructure/cdn-in-a-box`` directory. Remove the existing RPMs because they contain release Go binaries do not include useful debugging information. Rebuild the RPMs with no optimization, for debugging:

	.. code-block:: shell
		:caption: Remove release RPMs, then build debug RPMs

		make very-clean
		make debug traffic_monitor/traffic_monitor.rpm

	.. tip:: If you have gone through the steps to :ref:`dev-building-natively`, you can run ``make debug native traffic_monitor/traffic_monitor.rpm`` instead of ``make debug traffic_monitor/traffic_monitor.rpm`` to build the RPM quickly.

* Still in ``infrastructure/cdn-in-a-box``, open ``variables.env`` and set ``TM_DEBUG_ENABLE`` to ``true``.

* Stop CDN-in-a-Box if it is running and remove any existing volumes. Build the ``trafficmonitor-debug`` image to make sure it uses our fresh ``traffic_monitor.rpm``. Then, start CDN-in-a-Box:

	.. code-block:: shell
		:caption: docker compose command for debugging Traffic Monitor

		alias mydc='docker compose -f docker-compose.yml -f docker-compose.expose-ports.yml -f optional/docker-compose.debugging.yml'
		mydc down -v
		mydc build trafficmonitor
		mydc up

* Install `an IDE that supports delve <https://github.com/go-delve/delve/blob/master/Documentation/EditorIntegration.md>`_ and create a debugging configuration over port 2344. If you are using VS Code, the configuration should look like this:

	.. code-block:: json
		:caption: VS Code launch.json for debugging Traffic Monitor

		{
			"version": "0.2.0",
			"configurations": [
				{
					"name": "Traffic Monitor",
					"type": "go",
					"request": "attach",
					"mode": "remote",
					"port": 2344,
					"cwd": "${workspaceRoot}",
					"remotePath": "/tmp/go/src/github.com/apache/trafficcontrol",
				}
			]
		}

* Use the debugging configuration you created to start debugging Traffic Monitor. It should connect without first breaking at any line.

For an example of usage, set a breakpoint at `the o.m.RLock() call in ThreadsafeEvents.Get() <https://github.com/apache/trafficcontrol/blob/RELEASE-5.1.1/traffic_monitor/health/event.go#L71>`_, then visit http://trafficmonitor.infra.ciab.test/publish/EventLog (see :ref:`Traffic Monitor APIs: /publish/EventLog <tm-publish-EventLog>`).

Traffic Ops
===========

* Navigate to the ``infrastructure/cdn-in-a-box`` directory. Remove the existing RPMs because they contain release Go binaries do not include useful debugging information. Rebuild the RPMs with no optimization, for debugging:

	.. code-block:: shell
		:caption: Remove release RPMs, then build debug RPMs

		make very-clean
		make debug traffic_ops/traffic_ops.rpm

	.. tip:: If you have gone through the steps to :ref:`dev-building-natively`, you can run ``make debug native traffic_ops/traffic_ops.rpm`` instead of ``make debug traffic_ops/traffic_ops.rpm`` to build the RPM quickly.

* Still in ``infrastructure/cdn-in-a-box``, open ``variables.env`` and set ``TO_DEBUG_ENABLE`` to ``true``.

* Stop CDN-in-a-Box if it is running and remove any existing volumes. Build the ``trafficops-debug`` image to make sure it uses our fresh ``traffic_ops.rpm``. Then, start CDN-in-a-Box:

	.. code-block:: shell
		:caption: docker compose command for debugging Traffic Ops

		alias mydc='docker compose -f docker-compose.yml -f docker-compose.expose-ports.yml -f optional/docker-compose.debugging.yml'
		mydc down -v
		mydc build trafficops
		mydc up

* Install `an IDE that supports delve <https://github.com/go-delve/delve/blob/master/Documentation/EditorIntegration.md>`_ and create a debugging configuration over port 2345. If you are using VS Code, the configuration should look like this:

	.. code-block:: json
		:caption: VS Code launch.json for debugging Traffic Ops

		{
			"version": "0.2.0",
			"configurations": [
				{
					"name": "Traffic Ops",
					"type": "go",
					"request": "attach",
					"mode": "remote",
					"port": 2345,
					"cwd": "${workspaceRoot}",
					"remotePath": "/tmp/go/src/github.com/apache/trafficcontrol",
				}
			]
		}

* Use the debugging configuration you created to start debugging Traffic Ops. It should connect without first breaking at any line.

For an example of usage, set a breakpoint at `the log.Debugln() call in TOProfile.Read() <https://github.com/apache/trafficcontrol/blob/RELEASE-5.1.1/traffic_ops/traffic_ops_golang/profile/profiles.go#L148>`_, then visit https://trafficportal.infra.ciab.test/api/4.0/profiles (after logging into :ref:`tp-overview`).

Traffic Router
==============

* Navigate to the ``infrastructure/cdn-in-a-box`` directory.

* In ``variables.env``, set ``TR_DEBUG_ENABLE`` to ``true``.

* Install a debugging-capabe Java IDE or text editor of your choice. If unsure, install IntelliJ IDEA Community Edition.

* At the base of the repository (not in the ``cdn-in-a-box`` directory), open the ``traffic_router`` directory in your IDE.

* Add a new "Remote" (Java) debug configuration. Use port 5005.

* Start CDN-in-a-Box, including the "expose ports" "debugging" compose files:

	.. code-block:: shell
		:caption: docker compose command for debugging Traffic Router

		alias mydc='docker compose -f docker-compose.yml -f docker-compose.expose-ports.yml -f optional/docker-compose.debugging.yml'
		mydc down -v
		mydc build trafficrouter
		mydc up -d
		mydc logs --follow trafficrouter

* Watch the ``trafficrouter`` container's log. After DNS and certificate operations, the enroller, and Traffic Monitor, Traffic Router will start. Look for ``Listening for transport dt_socket at address: 5005`` in the example log below:

	.. code-block:: shell
		:caption: Log of the Docker container for Traffic Router

		        Warning:
		        The JKS keystore uses a proprietary format. It is recommended to migrate to PKCS12 which is an industry standard format using "keytool -importkeystore -srckeystore /opt/traffic_router/conf/keyStore.jks -destkeystore /opt/traffic_router/conf/keyStore.jks -deststoretype pkcs12".
		        Certificate stored in file <trafficrouter.infra.ciab.test.crt>

		        Warning:
		        The JKS keystore uses a proprietary format. It is recommended to migrate to PKCS12 which is an industry standard format using "keytool -importkeystore -srckeystore /opt/traffic_router/conf/keyStore.jks -destkeystore /opt/traffic_router/conf/keyStore.jks -deststoretype pkcs12".
		        Waiting for enroller initial data load to complete....
		        Waiting for enroller initial data load to complete....
		        Waiting for enroller initial data load to complete....
		        Waiting for enroller initial data load to complete....
		        Waiting for enroller initial data load to complete....
		        Waiting for enroller initial data load to complete....
		        Waiting for enroller initial data load to complete....
		        Waiting for enroller initial data load to complete....
		        Waiting for enroller initial data load to complete....
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        Waiting for Traffic Monitor to start...
		        tail: cannot open '/var/log/tomcat/catalina.log' for reading: No such file or directory
		        tail: cannot open '/var/log/tomcat/catalina.2020-02-21.log' for reading: No such file or directory
		        ==> /var/log/traffic_router/traffic_router.log <==

		        ==> /var/log/traffic_routr/access.log <==
		        Tomcat started.
		        tail: '/var/log/tomcat/catalina.log' has appeared;  following end of new file
		        tail: '/var/log/tomcat/catalina.2020-02-21.log' has appeared;  following end of new file

		        ==> /var/log/traffic_router/traffic_router.log <==
		        INFO  2020-02-21T05:16:07.557 [Thread-3] org.apache.traffic_control.traffic_router.protocol.LanguidPoller - Waiting for state from mbean path traffic-router:name=languidState
		        INFO  2020-02-21T05:16:07.557 [Thread-4] org.apache.traffic_control.traffic_router.protocol.LanguidPoller - Waiting for state from mbean path traffic-router:name=languidState
		        INFO  2020-02-21T05:16:07.558 [Thread-5] org.apache.traffic_control.traffic_router.protocol.LanguidPoller - Waiting for state from mbean path traffic-router:name=languidState
		        INFO  2020-02-21T05:16:07.559 [Thread-6] org.apache.traffic_control.traffic_router.protocol.LanguidPoller - Waiting for state from mbean path traffic-router:name=languidState

		        ==> /var/log/tomcat/catalina.log <==
		        Listening for transport dt_socket at address: 5005

		Watch for the line that mentions port 5005 -----------^^^^

		        ==> /var/log/tomcat/catalina.2020-02-21.log <==
		        21-Feb-2020 05:16:07.359 WARNING [main] org.apache.traffic_control.traffic_router.protocol.LanguidNioProtocol.<clinit> Adding BouncyCastle provider
		        21-Feb-2020 05:16:07.452 WARNING [main] org.apache.traffic_control.traffic_router.protocol.LanguidNioProtocol.<init> Serving wildcard certs for multiple domains
		        21-Feb-2020 05:16:07.459 WARNING [main] org.apache.traffic_control.traffic_router.protocol.LanguidNioProtocol.<init> Serving wildcard certs for multiple domains
		        21-Feb-2020 05:16:07.459 WARNING [main] org.apache.traffic_control.traffic_router.protocol.LanguidNioProtocol.<init> Serving wildcard certs for multiple domains
		        21-Feb-2020 05:16:07.461 INFO [main] org.apache.traffic_control.traffic_router.protocol.LanguidNioProtocol.setSslImplementationName setSslImplementation: org.apache.traffic_control.traffic_router.protocol.RouterSslImplementation

* When you see that Tomcat is listening for debugger connections on port 5005, start debugging using the debug configuration that you created.

Traffic Stats
===============

* Navigate to the ``infrastructure/cdn-in-a-box`` directory. Remove the existing RPMs because they contain release Go binaries do not include useful debugging information. Rebuild the RPMs with no optimization, for debugging:

	.. code-block:: shell
		:caption: Remove release RPMs, then build debug RPMs

		make very-clean
		make debug traffic_stats/traffic_stats.rpm

	.. tip:: If you have gone through the steps to :ref:`dev-building-natively`, you can run ``make debug native traffic_stats/traffic_stats.rpm`` instead of ``make debug traffic_stats/traffic_stats.rpm`` to build the RPMs quickly.

* Still in ``infrastructure/cdn-in-a-box``, open ``variables.env`` and set ``TS_DEBUG_ENABLE`` to ``true``.

* Stop CDN-in-a-Box if it is running and remove any existing volumes. Build the ``trafficstats-debug`` image to make sure it uses our fresh ``traffic_stats.rpm``. Then, start CDN-in-a-Box:

	.. code-block:: shell
		:caption: docker compose command for debugging Traffic Stats

		alias mydc='docker compose -f docker-compose.yml -f docker-compose.expose-ports.yml -f optional/docker-compose.debugging.yml'
		mydc down -v
		mydc build trafficstats
		mydc up

* Install `an IDE that supports delve <https://github.com/go-delve/delve/blob/master/Documentation/EditorIntegration.md>`_ and create a debugging configuration over port 2346. If you are using VS Code, the configuration should look like this:

	.. code-block:: json
		:caption: VS Code launch.json for debugging Traffic Stats

		{
			"version": "0.2.0",
			"configurations": [
				{
					"name": "Traffic Stats",
					"type": "go",
					"request": "attach",
					"mode": "remote",
					"port": 2346,
					"cwd": "${workspaceRoot}",
					"remotePath": "/tmp/go/src/github.com/apache/trafficcontrol",
				}
			]
		}

* Use the debugging configuration you created to start debugging Traffic Stats. It should connect without first breaking at any line.

For an example of usage, set a breakpoint at `the http.Get() call in main.getURL() <https://github.com/apache/trafficcontrol/blob/RELEASE-5.1.1/traffic_stats/traffic_stats.go#L706>`_, then wait 10 seconds for the breakpoint to be hit.

Troubleshooting
===============

* If you are debugging a Golang project and you don't see the values of all variables, or stepping to the next line puts you several lines ahead, rebuild the Docker image with an RPM built using :bash:`make debug`.
