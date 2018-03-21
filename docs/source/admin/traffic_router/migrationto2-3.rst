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
Traffic Router - Migrating to 2.3
*********************************
.. contents::
  :depth: 2
  :backlinks: top

Release Notes v2.3
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
* tomcat >= 8.5-28 installed or available in Yum repository (This package is created automatically by the Traffic Router build process)

Upgrade Procedure
==========================
* upload tomcat.rpm to a Yum repository
* uninstall previous version of Traffic Router
* install Traffic Router package
* restore property files

Upload tomcat.rpm
-----------------
The 'tomcat' package gets created when you build Traffic Router. You must either add it to the yum repo where you keep all of the Traffic Control packages, or manually copy it to the servers where you will be installing Traffic Router and run ```yum install [path to package]```

Uninstall Previous Traffic Router
---------------------------------
* ``` yum remove traffic_router ```
* Copy the traffic_router properties files you want to keep to a safe place: ``` cp /opt/traffic_router/conf/*.properties ~ ```
* delete whatever remains: ``` rm -rf /opt/traffic_router; rm -rf /opt/tomcat; rm -f /etc/init.d/tomcat ```

Install Traffic Router Package
------------------------------
If openssl, apr, tomcat-native, jdk and tomcat_tr packages are all in an available repository then you just need to run: ``` yum install traffic_router ```.


Restore Property Files
------------------------------
Replace the Traffic Router properties files with the correct ones  you saved earlier: ``` cp ~/.properties /opt/traffic_router/conf ```

Development Environment Upgrade
===============================

If you already have a development environment set up for the previous version of Traffic Router, then you will need to need to get and install these libraries on your workstation: openssl, apr and tomcat-native.
Also, whenever you run either 'mvn clean verify' or 'TrafficRouterStart' you will need to pass a command line parameter telling Java where to look for the 'tomcat-native' libraries:
``` mvn clean verify -Djava.library.path=[tomcat native library path on your box] ```
``` java -Djava.library.path=[tomcat native library path on your box] TrafficRouterStart ```

