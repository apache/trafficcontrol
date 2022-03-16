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

*****************
Developer's Guide
*****************
Use this guide to start developing applications that consume the Traffic Control APIs, to create extensions to Traffic Ops, or work on Traffic Control itself.

.. toctree::
	:maxdepth: 1

	api_guidelines
	environment_variables
	building
	debugging
	documentation_guidelines
	godocs
	traffic_monitor
	traffic_ops
	traffic_portal
	traffic_router
	traffic_stats

.. _dev:

The Development Environment
===========================
A development environment is available in :atc-file:`dev/`. This environment only depends on `Docker <https://www.docker.com/>`_ (version 20+) and `Docker-Compose <https://docs.docker.com/compose/` (version 1.27+) and enables rapid changes to be made to components during active development. This is, in general far faster than :ref:`dev-debugging-ciab`, but covers less complex configurations for testing purposes. Continuous Integration typically makes use of CDN-in-a-Box, so developers in general are free to use the Development Environment.

.. note:: Many ports used by the development environment clash with those exposed locally by CDN-in-a-Box when the :atc-file:`infrastructure/cdn-in-a-box/docker-compose.expose-ports.yml` Compose file is included, so the two cannot be used at the same time.

.. _dev-atc:

atc
---
The command ``atc`` is made available by sourcing :atc-file:`dev/atc.dev.sh` (e.g. ``source dev/atc.dev.sh``). While at the repository root, this command can be used to manipulate the development environment - most notably stopping and starting it.

Sourcing this file also sets :envvar:`TO_URL`, :envvar:`TO_USER`, and :envvar:`TO_PASSWORD` to the values appropriate for the default setup of the development environment, such that :ref:`toaccess` may be used to access the development Traffic Ops instance without any extra steps.

.. program:: atc

.. code-block:: bash
	:caption: ``atc`` Usage

	atc [-h] OPERATION

.. option:: -h, --help

	Print usage information and exit.

Each valid ``OPERATION`` is outlined in its corresponding section.

build
"""""
Build Docker images for the environment, but do not start it.

.. code-block:: bash
	:caption: ``atc build`` Usage

	atc build [SERVICE...]

.. option:: SERVICE

	If specified, only the given services will be built. By default, all services are built.

.. code-block:: bash
	:caption: ``atc build`` Example

	# Build only Traffic Ops
	atc build trafficops

	# Build all services
	atc build

connect
"""""""
Connect to a shell session inside a development container.

.. note:: Connecting to ``trafficrouter`` will result in connecting as a non-root user, so privileged access requires a more careful, manual use of :manpage:`docker(1)`.


exec
""""
Run a command in a dev container.

.. code-block:: bash
	:caption: ``atc exec`` Usage

	atc exec SERVICE CMD...

.. option:: SERVICE

	The service within which to execute commands.

.. option:: CMD

	An argv to pass to the service as a command.

.. code-block:: bash
	:caption: ``atc exec`` Example

	# Check ping statistics for communications from Traffic Ops to Traffic Monitor.
	atc exec trafficops ping -c 4 trafficmonitor

ready
"""""
Check if the development environment is ready. If it is ready the exit code is 0, and if it isn't ready the exit code is non-zero. "Readiness" is defined by the availability of the Traffic Ops API.

.. code-block:: bash
	:caption: ``atc ready [-h] [-w]`` Usage

	atc ready [SERVICE...]

.. option:: -h, --help

	Print usage information and exit.

.. option:: -w, --wait

	Wait for ATC to be ready, instead of just checking if it is ready.

.. code-block:: bash
	:caption: ``atc ready`` Example

	# Print "ready" if ATC is ready, "not ready" if it isn't.
	if atc ready; then
		echo "ready";
	else
		echo "not ready";
	fi

	# Block until ATC is ready before getting a CDN Snapshot for the development CDN.
	atc ready -w && toget -k cdns/dev/snapshot

restart
"""""""
Restart the development environment. This is functionally equivalent to stop_ followed by start_ where the same arguments that would be passed to ``restart`` are instead given to each of those.

.. warning:: Restarting Traffic Ops also stops every service that either it depends on or that depends on it - which is all of them. However, it only *starts* the services that Traffic Ops *depends on*, which is only the database service. So ``atc restart trafficops`` stops everything and only starts Traffic Ops back up again.

.. tip:: The services automatically rebuild the ATC components they run when those components change, so usually restarting is only necessary if you're making changes to the developer environment itself.

.. code-block:: bash
	:caption: ``atc restart`` Usage

	atc restart [SERVICE...]

.. option:: SERVICE

	If specified, only the given services will be restarted. By default, all services are restarted.

.. code-block:: bash
	:caption: ``atc restart`` Example

	# Restart only Traffic Router
	atc restart trafficrouter

	# Restart all services
	atc restart

start
"""""
Start up the development environment.

.. note:: Starting Traffic Ops also starts the Traffic Ops database and Traffic Vault (which isn't its own service).

.. code-block:: bash
	:caption: ``atc start`` Usage

	atc start [SERVICE...]

.. option:: SERVICE

	If specified, only the given services will be started. By default, all services are started.

.. code-block:: bash
	:caption: ``atc start`` Example

	# Start only Traffic Portal
	atc start trafficportal

	# Start all services
	atc start

stop
""""
Stop the development environment.

.. note:: Stopping Traffic Ops also stops every service that either it depends on or that depends on it - which is all of them.

.. code-block:: bash
	:caption: ``atc stop`` Usage

	atc stop [SERVICE...]

.. option:: SERVICE

	If specified, only the given services will be built. By default, all services are built.

.. code-block:: bash
	:caption: ``atc stop`` Example

	# Stop only Traffic Router
	atc stop trafficrouter

	# Stop all services
	atc stop

t3c
---
The ``atc.dev.sh`` file also provides a way to run :ref:`t3c-t3c` commands with debugging enabled. While in most production deployments :ref:`t3c-t3c` runs on a :manpage:`cron(8)` schedule, it is never run in the development environment, normally. One must manually trigger a run in order to run it and debug it.

The usage of this provided function is exactly as if one were simply calling :ref:`t3c-t3c`. A `delve <https://github.com/go-delve/delve/tree/master/Documentation>`_ debugging session is automatically started when :ref:`t3c-t3c` is run, which listens on port 8081.
