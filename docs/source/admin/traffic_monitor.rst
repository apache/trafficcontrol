.. 
.. Copyright 2015 Comcast Cable Communications Management, LLC
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

******************************
Traffic Monitor Administration
******************************
Installing Traffic Monitor
==========================
The following are requirements to ensure an accurate set up:

* Successful install of Traffic Ops
* Tomcat
* Administrative access to the Traffic Ops
* Physical address of the site
* perl-JSON
* perl-WWW-Curl

1. Create a FQDN that is resolvable in DNS.
2. Install Traffic Monitor and perl mods: ``yum -y install traffic_monitor perl-JSON perl-WWW-Curl``
3. Take the config from Traffic Ops: ``/opt/traffic_monitor/bin/traffic_monitor_config.pl``
4. Start Tomcat: ``service tomcat start`` ::


    Using CATALINA_BASE: /opt/tomcat
    Using CATALINA_HOME: /opt/tomcat
    Using CATALINA_TMPDIR: /opt/tomcat/temp
    Using JRE_HOME: /usr
    Using CLASSPATH:/opt/tomcat/bin/bootstrap.jar
    Using CATALINA_PID:/var/run/tomcat/tomcat.pid
    Starting tomcat [ OK ]

Configuring Traffic Monitor
===========================

Configuration Overview
----------------------
Traffic Monitor is configured using its JSON configuration file, ``traffic_monitor_config.js``. Specify the URL, username, password, and CDN name for the instance of Traffic Ops for which this Traffic Monitor is a member, and start the software.  Once started with the correct configuration, Traffic Monitor downloads its configuration from Traffic Ops and begins polling caches. Once a configurable number of polling cycles completes, health protocol state is available via RESTful JSON endpoints.

