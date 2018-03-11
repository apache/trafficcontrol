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

*****************************
Traffic Portal Administration
*****************************
The following are requirements to ensure an accurate set up:

* CentOS 6.7 or 7
* Node.js 6.0.x or above

**Installing Traffic Portal**

	- Download the Traffic Portal RPM from `Apache Jenkins <https://builds.apache.org/view/S-Z/view/TrafficControl/job/incubator-trafficcontrol-master-build/>`_ or build the Traffic Portal RPM from source (./pkg -v traffic_portal_build).
	- Copy the Traffic Portal RPM to your server
	- curl --silent --location https://rpm.nodesource.com/setup_6.x | sudo bash -
	- sudo yum install -y nodejs
	- sudo yum install -y <traffic_portal rpm>

**Configuring Traffic Portal**

	- update /etc/traffic_portal/conf/config.js (if upgrade, reconcile config.js with config.js.rpmnew and then delete config.js.rpmnew)
	- update /opt/traffic_portal/public/traffic_portal_properties.json (if upgrade, reconcile traffic_portal_properties.json with traffic_portal_properties.json.rpmnew and then delete traffic_portal_properties.json.rpmnew)
	- [OPTIONAL] update /opt/traffic_portal/public/resources/assets/css/custom.css (to customize traffic portal skin)

**Starting Traffic Portal**

	- sudo service traffic_portal start

**Stopping Traffic Portal**

	- sudo service traffic_portal stop







