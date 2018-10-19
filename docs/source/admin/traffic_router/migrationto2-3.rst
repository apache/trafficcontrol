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
==========================
* Replaced custom Java SNI implementation with a native implementation using tomcat-native, apr (Apache Portable Runtime) and OpenSSL
  This should significantly improve the performance of routing 'https' delivery services.
* Upgraded to Tomcat 8.5.30
* Separated the Traffic Router installation from the Tomcat deployment and created a new 'tomcat' package for installing Tomcat.
  Traffic Router and Tomcat can now be upgraded independently
* Converted Traffic Router to a 'systemd' service
* Modified the development test and dev deployment processes to be more consistent with production

System Requirements
==========================
* Centos 7.2
* OpenSSL >= 1.0.2 installed
* JDK >= 8.0 installed or available in Yum repository
* APR (Apache Portable Runtime) >= 1.4.8-3 installed or available in Yum repository
* Tomcat Native >= 1.2.16 installed or available in Yum repository
* tomcat >= 8.5-30 installed or available in Yum repository (This package is created automatically by the Traffic Router build process)

Upgrade Procedure
==========================
* upload tomcat.rpm to a Yum repository
* update the traffic_router package
* restore property files

Upload tomcat.rpm
-----------------
The 'tomcat' package gets created when you build Traffic Router. You must either add it to the yum repo where you keep all of the Traffic Control packages, or manually copy it to the servers where you will be installing Traffic Router and run ``yum install [path to package]``
It is preferable that you add it to your Yum repository because then it will be installed automatically when you perform the Traffic Router update.

Update the traffic_router Package
---------------------------------
If openssl, apr, tomcat-native, jdk and tomcat_tr packages are all in an available repository then you just need to run: ``yum update traffic_router``.
This will first cause the apr, tomcat-native, jdk and tomcat packages to be installed. When the 'tomcat' package runs, it will cause any older versions of traffic_router or tomcat to be uninstalled. This is because the previous versions of the traffic_router package included an untracked installation of tomcat.


Restore Property Files
------------------------------
The install process does not override or replace any of the files in the /opt/traffic_router/conf directory. Previous versions of the traffic_ops.properties, traffic_monitor.properties and startup.properties should still be good. On a new install replace the Traffic Router properties files with the correct ones for the CDN.

Development Environment Upgrade
===============================

If you already have a development environment set up for the previous version of Traffic Router, then you will need to get and install these libraries on your workstation: openssl, apr and tomcat-native.
Also, whenever you run either 'mvn clean verify' or 'TrafficRouterStart' you will need to pass a command line parameter telling Java where to look for the 'tomcat-native' libraries:
``mvn clean verify -Djava.library.path=[tomcat native library path on your box]``
``java -Djava.library.path=[tomcat native library path on your box] TrafficRouterStart``

