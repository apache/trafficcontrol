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

*****************************
Traffic Router Administration
*****************************
Installing Traffic Router
==========================
The following are requirements to ensure an accurate set up:

* CentOS 6
* 4 vCPUs
* 8GB RAM
* Successful install of Traffic Ops
* Successful install of Traffic Monitor
* Administrative access to the Traffic Ops
* Physical address of the site
* perl-JSON
* perl-WWW-Curl

1. Enter the Traffic Router server into Traffic Ops.
2. Make sure the FQDN of the Traffic Monitor is resolvable in DNS.
3. Install a traffic router: ``sudo yum install traffic_router``.
4. Edit ``/opt/traffic_router/conf/traffic_monitor.properties`` and put in the correct online Traffic Monitor(s) for your CDN.

 Example: ::

	# list of ips or hostnames delimited by semicolon (;)
	traffic_monitor.bootstrap.hosts=traffic-mon-01.cdn.kabletown.net:80;

	# Instead of using the traffic_monitor.bootstrap.hosts property as a bootstrap
	# source before switching to ONLINE Monitors in the TrConfig, always
	# use the hosts listed for TrConfig and TrStates. Defaults to false.
	traffic_monitor.bootstrap.local = false

	# traffic_monitor.properties: url that should normally point to this file
	traffic_monitor.properties=file:/opt/traffic_router/conf/traffic_monitor.properties

	# Frequency for reloading this file
	# traffic_monitor.properties.reload.period=60000
   

5. Start Tomcat: ``sudo service tomcat start``, and test lookups with dig and curl against that server.
6. Snapshot CRConfig 

* This instantly associates production traffic on the servers. They need to be online when you change the DNS records.

7. Add the correct DNS entries to the SOA records for the CDN on which you are working.

8. Add the servers to the NS and SOA records for your domain.

Configuring Traffic Router
==========================
1. From **Misc > Profiles**, verify the following:
 
 * The Traffic Router information.
 * The profile is set correctly.
 * The Status is set to OFFLINE.

2. Verify the functionality of the DNS entry for the Traffic Router.
3. Click **Tools > Generate ISO**.
4. Complete the necessary fields.
5. Click **Download ISO**.

Troubleshooting and log files
=============================
Traffic Router log files are in ``/opt/traffic_router/var/log/``, and tomcat log files are in ``/opt/tomcat/logs/``.
