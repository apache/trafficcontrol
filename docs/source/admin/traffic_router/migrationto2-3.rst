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

*********************************
Traffic Router - Migrating to 3.0
*********************************
.. contents::
	:depth: 2
	:backlinks: top

Release Notes v3.0
==================
* Replaced custom Java :abbr:`SNI (Server Name Indication)` implementation with a native implementation using tomcat-native, :abbr:`APR (Apache Portable Runtime)` and OpenSSL. This should significantly improve the performance of routing HTTPS :term:`Delivery Services`.

	.. seealso:: `The Server Name Indication Wikipedia page <https://en.wikipedia.org/wiki/Server_Name_Indication>`_, `The Apache Portable Runtime project site <https://apr.apache.org/>`_ and/or `the OpenSSL project site <https://www.openssl.org/>`_

* Upgraded to Tomcat 8.5.30
* Separated the Traffic Router installation from the Tomcat deployment and created a new 'tomcat' package for installing Tomcat. Traffic Router and Tomcat can now be upgraded independently
* Converted Traffic Router to a :manpage:`systemd(1)` service
* Modified the development test and deployment processes to be more consistent with production

System Requirements
===================
* Centos 7.9 or CentOS 8.2
* OpenSSL >= 1.0.2 installed
* JDK >= 8.0 installed or available in an accessible :manpage:`yum(8)` repository
* :abbr:`APR (Apache Portable Runtime)` >= 1.4.8-3 installed or available in an accessible :manpage:`yum(8)` repository
* Tomcat Native >= 1.2.16 installed or available in an accessible :manpage:`yum(8)` repository
* tomcat >= 8.5-30 installed or available in an accessible :manpage:`yum(8)` repository (This package is created automatically by the Traffic Router build process)

Upgrade Procedure
=================
* upload the :file:`dist/tomcat-{version string}.rpm` file generated as a part of the build instructions outlined in :ref:`dev-building` to an accessible :manpage:`yum(8)` repository
* update the ``traffic_router`` package with :manpage:`yum(8)`
* restore property files

Upload tomcat.rpm
-----------------
The :file:`term-{version string}.rpm` package should have been created when Traffic Router was built according to the instructions in :ref:`dev-building`. It must must either be added to an accessible :manpage:`yum(8)` repository, or manually copied to the servers where Traffic Router will be installed. It is generally better that it be added to a :manpage:`yum(8)` repository because then it will be installed automatically when Traffic Router is updated.

Update the traffic_router Package
---------------------------------
If ``openssl``, ``apr``, ``tomcat-native``, ``java-1.8.0-openjdk``, ``java-1.8.0-openjdk-devel`` and ``tomcat_tr`` packages are all in an available :manpage:`yum(8)` repository then an upgrade can be performed by running ``yum update traffic_router`` as the root user or with :manpage:`sudo(8)`. This will first cause the ``apr``, ``tomcat-native``, ``java-1.8.0-openjdk``, ``java-1.8.0-openjdk-devel`` and ``tomcat`` packages to be installed. When the ``tomcat`` package runs, it will cause any older versions of ``traffic_router`` or ``tomcat`` to be uninstalled. This is because the previous versions of the ``traffic_router`` package included an untracked installation of ``tomcat``.

Restore Property Files
----------------------
The install process does not override or replace any of the files in the :file:`/opt/traffic_router/conf` directory. Previous versions of the :file:`traffic_ops.properties`, :file:`traffic_monitor.properties` and :file:`startup.properties` should still be good. On a new install replace the Traffic Router properties files with the correct ones for the CDN.

Development Environment Upgrade
===============================
If a development environment is already set up for the previous version of Traffic Router, then ``openssl``, ``apr`` and ``tomcat-native`` will need to be manually installed with :manpage:`yum(8)` or :manpage:`rpm(8)`. Also, whenever either ``mvn clean verify`` or ``TrafficRouterStart`` is/are run, the location of the ``tomcat-native`` libraries will need to be made known to the :abbr:`JVM (Java Virtual Machine)` via command line arguments.

.. code-block:: shell
	:caption: Example Commands Specifying a Path to the tomcat-native Library

	mvn clean verify -Djava.library.path=[tomcat native library path on your box]
	java -Djava.library.path=[tomcat native library path on your box] TrafficRouterStart
