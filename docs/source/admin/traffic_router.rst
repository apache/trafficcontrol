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

* Successful install of Traffic Ops
* Successful install of Traffic Monitor
* Administrative access to the Traffic Ops
* Physical address of the site
* perl-JSON
* perl-WWW-Curl

1. Install a traffic router: ``sudo yum install traffic_router``.
2. Edit ``/opt/traffic_router/conf/traffic_monitor.properties`` and put in the correct online Traffic Monitor(s) for your CDN.

 Example: ::


   

3. Start Tomcat and test lookups with dig and curl against that server.
4. Verify that the correct firewall rules are in place.
5. Use a security scanner, such as `nmap <http://nmap.org/>`_, to test that the Traffic Router has ports 53/tcp, 80/tcp and 8080/tcp open.

 Below is a listing of the old Traffic Router open ports: ::


 


6. Set the servers to ONLINE.

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

