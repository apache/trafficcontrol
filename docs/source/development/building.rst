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

.. _dev-building:

************************
Building Traffic Control
************************
The build steps for Traffic Control components are all pretty much the same, despite that they are written in a variety of different languages and frameworks. This is accomplished by using Docker.

.. note:: Currently, both listed methods of building Traffic Control components will produce ``*.rpm`` files, meaning that the support of these components is limited to RedHat-based distributions - and none of them are currently tested (or guaranteed to work) outside of Rocky Linux 8 and CentOS 7, specifically.

Downloading Traffic Control
===========================
If any local work on Traffic Monitor, Traffic Router Golang, Grove or Traffic Ops is to be done, it is **highly** recommended that `the Traffic Control repository <https://github.com/apache/trafficcontrol>`_ be downloaded inside the ``$GOPATH`` directory. Specifically, the best location is ``$GOPATH/src/github.com/apache/trafficcontrol``. Cloning the repository outside of this location will require either linking the actual directory to that point, or moving/copying it there.

.. seealso:: The Golang project's ``GOPATH`` `wiki page <https://github.com/golang/go/wiki/GOPATH>`_

.. _pkg:

Build Using ``pkg``
===================
This is the easiest way to build all the components of Traffic Control; all requirements are automatically loaded into the image used to build each component.  The ``pkg`` command can be found at the root of the Traffic Control `repository <https://github.com/apache/trafficcontrol/blob/master/pkg>`_.

Requirements
------------
- `Docker <https://docs.docker.com/engine/installation/>`_
- `Docker Compose <https://docs.docker.com/compose/install/>`_\ [#compose-optional]_


Usage
-----
``./pkg [options] [projects]``

Options
"""""""

.. option:: -7

	Build RPMs targeting CentOS 7.

	.. versionchanged:: ATCv6.0.0

		Previously, :option:`-7` was implicit if not given. As of ATC version 6.0.0, this is no longer the case, and :option:`-8` is implicit instead.

.. option:: -8

	Build RPMs targeting Rocky Linux 8 (default).

	.. versionchanged:: ATCv6.0.0

		Previously, :option:`-7` was implicit if not given. As of ATC version 6.0.0, this is no longer the case, and :option:`-8` is implicit instead.

.. option:: -a

	Build all projects, including optional ones.

.. option:: -b

	Build builder Docker images before building projects.

.. option:: -d

	Disable compiler optimizations for debugging.

.. option:: -f FILE

	Use ``FILE`` instead of the default Docker-Compose file (``./infrastructure/docker/build/docker-compose.yml``).

.. option:: -h

	Print help message and exit.

	.. versionadded:: ATCv6.1.0

.. option:: -l

	List available projects.

	.. caution:: This lists only the projects that are built by default if none are specified, not *all* projects that can be built. See :issue:`6272`.

.. option:: -L

	Don't write logs to files - respects output levels on STDERR/STDOUT as set by :option:`-q`/:option:`-v`.

.. option:: -o

	Build from the optional list. Same as passing :option:`-f` with the option-argument ``./infrastructure/docker/build/docker-compose-opt.yml``.

.. option:: -p

	Pull builder Docker images, do not build them (default).

.. option:: -q

	Quiet mode. Suppresses output (default).

.. option:: -s

	Simple output filenames - e.g. ``traffic_ops.rpm`` instead of ``traffic_ops-6.1.0-11637.ec9ff6a6.el8.x86_64.rpm``.

	.. versionadded:: ATCv6.1.0

.. option:: -S

	Skip building "source RPMs".

	.. versionadded:: ATCv6.1.0

.. option:: -v

	Verbose mode. Lists all build output.

	.. versionadded:: ATCv6.1.0


If present, ``projects`` should be one or more project names. When no specific project or project list is given the default projects will be built. Valid projects:

- ats\ [#optional-project]_
- docs\ [#default-project]_
- fakeorigin_build\ [#optional-project]_
- grove_build\ [#default-project]_
- grovetccfg_build\ [#default-project]_
- source\ [#default-project]_
- traffic_monitor_build\ [#default-project]_
- traffic_ops_build\ [#default-project]_
- cache-config_build\ [#default-project]_
- traffic_portal_build\ [#default-project]_
- traffic_router_build\ [#default-project]_
- traffic_stats_build\ [#default-project]_
- weasel\ [#default-project]_

Output :file:`{component}-{version}.rpm` files, build logs and source tarballs will be output to the ``dist/`` directory at the root of the Traffic Control repository directory.

.. _build-with-dc:

Build Using ``docker compose``
------------------------------
If the ``pkg`` script fails, ``docker compose`` can still be used to build the projects directly. The compose file can be found at ``infrastructure/docker/build/docker-compose.yml`` under the repository's root directory. It can be passed directly to ``docker compose``, either from the ``infrastructure/docker/build/`` directory or by explicitly passing a path to the ``infrastructure/docker/build/docker-compose.yml`` file via ``-f``. It is recommended that between builds ``docker compose down -v`` is run to prevent caching of old build steps. The service names are the same as the project names described above in `Usage`_, and similar to the ``pkg`` script, the build results, logs and source tarballs may all be found in the ``dist`` directory after completion.

.. note:: Calling ``docker compose`` in the way described above will build _all_ projects, not just the default projects.

.. seealso:: `The Docker Compose command line reference <https://docs.docker.com/compose/reference/overview/>`_

.. _dev-building-natively:

Build the RPMs Natively
=======================
A developer may end up building the RPMs several times to test or :ref:`debug <dev-debugging-ciab>` code changes, so it can be desirable to build the RPMs quickly for this purpose. Natively building the RPMs has the lowest build time of any building method.

Install the Dependencies
------------------------

.. table:: Build dependencies for Traffic Control

	+---------------------------------+---------------------+----------------------------+------------------------+---------------------------+---------------------------+--------------------------+----------+------------------------------+--------------------------+
	| OS/Package Manager              | Common dependencies | :ref:`dev-traffic-monitor` | :ref:`dev-traffic-ops` | :ref:`dev-traffic-portal` | :ref:`dev-traffic-router` | :ref:`dev-traffic-stats` | Grove    | Grove TC Config (grovetccfg) | :ref:`Docs <docs-guide>` |
	+=================================+=====================+============================+========================+===========================+===========================+==========================+==========+==============================+==========================+
	| macOS\ [#mac-jdk]_              | - coreutils         | - go                       | - go                   | - npm                     | - maven                   | - go                     | - go     | - go                         | - python3                |
	| (homebrew_)                     | - rpm               |                            |                        | - grunt-cli               |                           |                          |          |                              |                          |
	+---------------------------------+---------------------+----------------------------+------------------------+---------------------------+---------------------------+--------------------------+----------+------------------------------+--------------------------+
	| Rocky\ Linux\ [#rocky-go]_,     | - git               |                            |                        | - epel-release            | - java-11-openjdk         |                          |          |                              | - python3-devel          |
	| Red Hat,                        | - rpm-build         |                            |                        | - npm                     | - maven                   |                          |          |                              | - gcc                    |
	| Fedora,                         | - rsync             |                            |                        | - nodejs-grunt-cli        |                           |                          |          |                              | - make                   |
	| CentOS                          |                     |                            |                        |                           |                           |                          |          |                              |                          |
	| (yum_)                          |                     |                            |                        |                           |                           |                          |          |                              |                          |
	+---------------------------------+---------------------+----------------------------+------------------------+---------------------------+---------------------------+--------------------------+----------+------------------------------+--------------------------+
	| Arch Linux,                     | - git               | - go                       | - go                   | - npm                     | - jdk11-openjdk           | - go                     | - go     | - go                         | - python-pip             |
	| Manjaro                         | - rpm-tools         |                            |                        | - grunt-cli               | - maven                   |                          |          |                              | - python-sphinx          |
	| (pacman_)                       | - diff              |                            |                        |                           |                           |                          |          |                              | - make                   |
	+---------------------------------+---------------------+----------------------------+------------------------+---------------------------+---------------------------+--------------------------+----------+------------------------------+--------------------------+
	| Windows                         | - git               |                            |                        |                           | - curl                    |                          |          |                              |                          |
	| (cygwin_)\ [#windeps]_          | - rpm-build         |                            |                        |                           |                           |                          |          |                              |                          |
	|                                 | - rsync             |                            |                        |                           |                           |                          |          |                              |                          |
	+---------------------------------+---------------------+----------------------------+------------------------+---------------------------+---------------------------+--------------------------+----------+------------------------------+--------------------------+
	| Windows                         |                     | - golang                   | - golang               | - nodejs                  | - openjdk11               | - golang                 | - golang | - golang                     | - python                 |
	| (chocolatey_)\ [#windeps]_      |                     |                            |                        |                           | - maven                   |                          |          |                              | - pip                    |
	|                                 |                     |                            |                        |                           |                           |                          |          |                              | - make                   |
	+---------------------------------+---------------------+----------------------------+------------------------+---------------------------+---------------------------+--------------------------+----------+------------------------------+--------------------------+

.. _homebrew:   https://brew.sh/
.. _yum:        https://wiki.centos.org/PackageManagement/Yum
.. _pacman:     https://www.archlinux.org/pacman/
.. _cygwin:     https://cygwin.com/
.. _chocolatey: https://chocolatey.org/

After installing the packages using your platform's package manager,

- Install the :ref:`global NPM packages <dev-tp-global-npm>` to build Traffic Portal.
- Install the Python 3 modules used to :ref:`build the documentation <docs-build>`.

Run ``build/clean_build.sh`` directly
-------------------------------------

In a terminal, navigate to the root directory of the repository. You can run ``build/clean_build.sh`` with no arguments to build all components.

.. code-block:: shell
	:caption: ``build/clean_build.sh`` with no arguments

	build/clean_build.sh

This is the equivalent of running

.. code-block:: shell
	:caption: ``build/clean_build.sh`` with all components

	build/clean_build.sh tarball traffic_monitor traffic_ops traffic_portal traffic_router traffic_stats grove grove/grovetccfg docs

If any component fails to build, no further component builds will be attempted.

By default, the RPMs will be built targeting Rocky Linux 8. CentOS 7 is also a supported build target. You can choose which RHEL version to build for (8, 7, etc.) by setting the ``RHEL_VERSION`` environment variable:

.. code-block:: shell
	:caption: Building RPMs that target CentOS 7 without the build host needing to be CentOS 7

	export RHEL_VERSION=7

.. warning:: Although there are no known issues with natively-built RPMs, the official, supported method of building the RPMs is by using :ref:`pkg <pkg>` or :ref:`docker compose <build-with-dc>`. Use natively-built RPMs at your own risk.

Building Individual Components
==============================
Each Traffic Control component can be individually built, and the instructions for doing so may be found in their respective component's development documentation.

Building This Documentation
---------------------------
See instructions for :ref:`building the documentation <docs-build>`.

.. [#compose-optional] This is optional, but recommended. If a ``docker compose`` executable is not available the ``pkg`` script will automatically download and run it using a container. This is noticeably slower than running it natively.
.. [#optional-project] This project is "optional", which means that it cannot be built unless :option:`-o` is given.
.. [#default-project] This is a default project, which will be built if ``pkg`` is run with no ``projects`` argument
.. [#mac-jdk] If you are on macOS, you additionally need to :ref:`dev-tr-mac-jdk`.
.. [#rocky-go] If you are on Rocky Linux, you need to `download Go directly <https://golang.org/dl/>`_ instead of using a package manager in order to get the latest Go version. For most users, the desired architecture is AMD64/x86_64.
.. [#windeps] If you are on Windows, you need to install **both** the Cygwin packages and the Chocolatey packages in order to build the Apache Traffic Control RPMs natively.
