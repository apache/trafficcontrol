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

.. note:: Currently, both listed methods of building Traffic Control components will produce ``*.rpm`` files, meaning that the support of these components is limited to RedHat-based distributions - and none of them are currently tested (or guaranteed to work) outside of CentOS 7 and CentOS 8, specifically.

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
- `Docker Compose <https://docs.docker.com/compose/install/>`_\ [1]_


Usage
-----
``./pkg [options] [projects]``

.. note:: The ``pkg`` script often needs to be run as ``sudo``, as certain privileges are required to run Docker containers

Options

-7    Build RPMs targeting CentOS 7 (default)
-8    Build RPMs targeting CentOS 8
-b    Build builder Docker images before building projects
-d    Disable compiler optimizations for debugging.
-l    List available projects.
-p    Pull builder Docker images, do not build them (default)
-q    Quiet mode. Supresses output. (default)
-v    Verbose mode. Lists all build output.

If present, ``projects`` should be one or more project names. When no specific project or project list is given the default projects will be built. Valid projects:

- docs
- grove_build\ [2]_
- grovetccfg_build
- source\ [2]_
- traffic_monitor_build\ [2]_
- traffic_ops_build\ [2]_
- traffic_ops_ort_build\ [2]_
- traffic_portal_build\ [2]_
- traffic_router_build\ [2]_
- traffic_stats_build\ [2]_
- weasel

Output :file:`{component}-{version}.rpm` files, build logs and source tarballs will be output to the ``dist/`` directory at the root of the Traffic Control repository directory.

.. [1] This is optional, but recommended. If a ``docker-compose`` executable is not available the ``pkg`` script will automatically download and run it using a container. This is noticeably slower than running it natively.
.. [2] This is a default project, which will be built if ``pkg`` is run with no ``projects`` argument

.. _build-with-dc:

Build Using ``docker-compose``
------------------------------
If the ``pkg`` script fails, ``docker-compose`` can still be used to build the projects directly. The compose file can be found at ``infrastructure/docker/build/docker-compose.yml`` under the repository's root directory. It can be passed directly to ``docker-compose``, either from the ``infrastructure/docker/build/`` directory or by explicitly passing a path to the ``infrastructure/docker/build/docker-compose.yml`` file via ``-f``. It is recommended that between builds ``docker-compose down -v`` is run to prevent caching of old build steps. The service names are the same as the project names described above in `Usage`_, and similar to the ``pkg`` script, the build results, logs and source tarballs may all be found in the ``dist`` directory after completion.

.. note:: Calling ``docker-compose`` in the way described above will build _all_ projects, not just the default projects.

.. seealso:: `The Docker Compose command line reference <https://docs.docker.com/compose/reference/overview/>`_

.. _dev-building-natively:

Build the RPMs Natively
=======================
A developer may end up building the RPMs several times to test or :ref:`debug <dev-debugging-ciab>` code changes, so it can be desirable to build the RPMs quickly for this purpose. Natively building the RPMs has the lowest build time of any building method.

Install the Dependencies
------------------------

.. table:: Build dependencies for Traffic Control

	+------------------------------------+---------------------+----------------------------+------------------------+---------------------------+---------------------------+--------------------------+----------+------------------------------+--------------------------+
	|                                    | Common dependencies | :ref:`dev-traffic-monitor` | :ref:`dev-traffic-ops` | :ref:`dev-traffic-portal` | :ref:`dev-traffic-router` | :ref:`dev-traffic-stats` | Grove    | Grove TC Config (grovetccfg) | :ref:`Docs <docs-guide>` |
	+====================================+=====================+============================+========================+===========================+===========================+==========================+==========+==============================+==========================+
	| macOS (homebrew_)\ [3]_            | - rpm               | - go                       | - go                   | - npm                     | - maven                   | - go                     | - go     | - go                         | - python3                |
	|                                    |                     |                            |                        | - bower                   |                           |                          |          |                              |                          |
	|                                    |                     |                            |                        | - grunt-cli               |                           |                          |          |                              |                          |
	+------------------------------------+---------------------+----------------------------+------------------------+---------------------------+---------------------------+--------------------------+----------+------------------------------+--------------------------+
	| CentOS/Red Hat/Fedora (yum_)\ [4]_ | - git               |                            |                        | - epel-release            | - java-1.8.0-openjdk      |                          |          |                              | - python3-devel          |
	|                                    | - rpm-build         |                            |                        | - npm                     | - maven                   |                          |          |                              | - gcc                    |
	|                                    | - rsync             |                            |                        | - nodejs-grunt-cli        |                           |                          |          |                              | - make                   |
	|                                    |                     |                            |                        | - ruby-devel              |                           |                          |          |                              |                          |
	|                                    |                     |                            |                        | - gcc                     |                           |                          |          |                              |                          |
	|                                    |                     |                            |                        | - make                    |                           |                          |          |                              |                          |
	+------------------------------------+---------------------+----------------------------+------------------------+---------------------------+---------------------------+--------------------------+----------+------------------------------+--------------------------+
	| Arch Linux (pacman_)               | - git               | - go                       | - go                   | - npm                     | - jdk8-openjdk            | - go                     | - go     | - go                         | - python-pip             |
	|                                    | - rpm-tools         |                            |                        | - bower                   | - maven                   |                          |          |                              | - python-sphinx          |
	|                                    | - diff              |                            |                        | - grunt-cli               |                           |                          |          |                              | - make                   |
	|                                    | - rsync             |                            |                        | - ruby                    |                           |                          |          |                              |                          |
	|                                    |                     |                            |                        | - gcc                     |                           |                          |          |                              |                          |
	|                                    |                     |                            |                        | - make                    |                           |                          |          |                              |                          |
	+------------------------------------+---------------------+----------------------------+------------------------+---------------------------+---------------------------+--------------------------+----------+------------------------------+--------------------------+
	| Windows (cygwin_)\ [5]_            | - git               |                            |                        | - ruby-devel              | - curl                    |                          |          |                              |                          |
	|                                    | - rpm-build         |                            |                        | - make                    |                           |                          |          |                              |                          |
	|                                    | - rsync             |                            |                        | - gcc-g++                 |                           |                          |          |                              |                          |
	+------------------------------------+---------------------+----------------------------+------------------------+---------------------------+---------------------------+--------------------------+----------+------------------------------+--------------------------+
	| Windows (chocolatey_)\ [5]_        |                     | - golang                   | - golang               | - nodejs                  | - openjdk8                | - golang                 | - golang | - golang                     | - python                 |
	|                                    |                     |                            |                        |                           | - maven                   |                          |          |                              | - pip                    |
	|                                    |                     |                            |                        |                           |                           |                          |          |                              | - make                   |
	+------------------------------------+---------------------+----------------------------+------------------------+---------------------------+---------------------------+--------------------------+----------+------------------------------+--------------------------+

.. _homebrew:   https://brew.sh/
.. _yum:        https://wiki.centos.org/PackageManagement/Yum
.. _pacman:     https://www.archlinux.org/pacman/
.. _cygwin:     https://cygwin.com/
.. _chocolatey: https://chocolatey.org/

.. [3] If you are on macOS, you additionally need to :ref:`dev-tr-mac-jdk`.

.. [4] If you are on CentOS, you need to `download Go directly <https://golang.org/dl/>`_ instead of using a package manager in order to get the latest Go version. For most users, the desired architecture is AMD64/x86_64.

.. [5] If you are on Windows, you need to install both the Cygwin packages and the Chocolatey packages in order to build the Apache Traffic Control RPMs natively.

.. |AdoptOpenJDK instructions| replace:: add the AdoptOpenJDK tap and install the ``adoptopenjdk8`` cask
.. _AdoptOpenJDK instructions: https://github.com/AdoptOpenJDK/homebrew-openjdk#other-versions

After installing the packages using your platform's package manager,

	- Install the :ref:`global NPM packages <dev-tp-global-npm>` and :ref:`Compass <dev-tp-compass>` to build Traffic Portal.

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

By default, the RPMs will be built targeting CentOS 7. CentOS 8 is also a supported build target. You can choose which CentOS version to build for (7, 8, etc.) by setting the ``RHEL_VERSION`` environment variable:

.. code-block:: shell
	:caption: Building RPMs that target CentOS 8 without the build host needing to be CentOS 8

	export RHEL_VERSION=8

.. warning:: Although there are no known issues with natively-built RPMs, the official, supported method of building the RPMs is by using :ref:`pkg <pkg>` or :ref:`docker-compose <build-with-dc>`. Use natively-built RPMs at your own risk.

Building Individual Components
==============================
Each Traffic Control component can be individually built, and the instructions for doing so may be found in their respective component's development documentation.

Building This Documentation
---------------------------
See instructions for :ref:`building the documentation <docs-build>`.
