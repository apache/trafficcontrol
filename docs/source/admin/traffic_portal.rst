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
Traffic Portal Administration
*****************************
The following are requirements to ensure an accurate set up:

* CentOS 6.7 or 7
* Node.js 6.0.x or above

**Installing Traffic Portal**

	- Download the Traffic Portal RPM from the traffic control `downloads <http://traffic-control-cdn.net/downloads/index.html>`_ page or build from `source <https://github.com/Comcast/traffic_control/tree/master/traffic_portal/build>`_.
	- Copy the Traffic Portal RPM to your server
	- curl --silent --location https://rpm.nodesource.com/setup_6.x | sudo bash -
	- sudo yum install -y nodejs
	- sudo yum install -y <traffic_portal rpm>

**Configuring Traffic Portal**

	- cd /etc/traffic_portal/conf
	- sudo cp config-template.js config.js
	- sudo vi config.js (read the inline comments)
	- [OPTIONAL] sudo vi /opt/traffic_portal/public/traffic_portal_properties.json (to customize traffic portal content)
	- [OPTIONAL] sudo vi /opt/traffic_portal/public/resources/assets/css/custom.css (to customize traffic portal skin)

**Starting Traffic Portal**

	- sudo service traffic_portal start

**Stopping Traffic Portal**

	- sudo service traffic_portal stop







