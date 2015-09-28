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

****************************
Traffic Stats Administration
****************************
Installation
========================
Traffic Stats consists of three components:  Traffic Stats, InfluxDb, and Grafana.  Below are instructions on installing each component.

**Installing Traffic Stats:**

	- Download the Traffic Stats RPM from the traffic control `downloads <http://traffic-control-cdn.net/downloads/index.html>`_ page.
	- Copy the Traffic Stats RPM to your server
	- sudo rpm -ivh <traffic_stats rpm>

      Note:  This installation actually creates two separate services:  write_traffic_stats and ts_daily_summary.  More information on these services can be found in the overview section.     

**Installing InfluxDb:**

	In order to store traffic stats data you will need to install InfluxDb.  It is recommended InfluxDb be installed in a 3 server cluster; VMs are acceptable. The documentation for installing InfluxDb can be found on the InfluxDb `website <https://influxdb.com/docs/v0.9/introduction/installation.html>`_.

**Installing Grafana:**

	Grafana is used to display Traffic Stats/InfluxDb data in Traffic Ops.  Grafana is typically run on the same server as Traffic Stats but this is not a requirement.  Grafana can be installed on any server that can access InfluxDb and be accessed by Traffic Ops.  Documentation on installing Grafana can be found `here <http://docs.grafana.org/installation/>`_.

Configuration
=========================

**Configuring Traffic Stats:**

	Traffic Stats’ configuration file can be found in /opt/traffic_stats/conf/traffic_stats.cfg.
	The following values need to be configured: 

	     - *toUser:* The user used to connect to Traffic Ops
	     - *toPasswd:*  The password to use when connecting to Traffic Ops
	     - *toUrl:*  The URL of the Traffic Ops server used by Traffic Stats
	     - *influxUser:*  The user to use when connecting to InfluxDb (if configured on InfluxDb, else leave default)
	     - *influxPassword:*  That password to use when connecting to InfluxDb (if configured, else leave blank)
	     - *polling interval:*  The interval at which Traffic Monitor is polled and stats are stored in InfluxDb
	     - *statusToMon:*  The status of Traffic Monitor to poll (poll ONLINE or OFFLINE traffic monitors)
	     - *seelogConfig:*  The absolute path of the seelong config file
	     - *dailySummaryPollingInterval:* The interval, in seconds, at which Traffic Stats checks to see if daily stats need to be computed and stored.

**Configuring InfluxDb:**

	At a minimum, InfluxDb will need to be configured for clutstering.  Documentation on clustering configuration can be found on the clustering page of the `InfluxDb Website <https://influxdb.com/docs/v0.9/concepts/clustering.html>`_.

**Configuring Grafana:**

	In order for grafana to integrate with Traffic Ops, it will need to be configured to allow anonymous authorization.  More information can be found on the configuration page of the `Grafana Website  <http://docs.grafana.org/installation/configuration/#authanonymous>`_. 

	Traffic Ops uses Grafana to display stats data to users.  In order for the graphs to display correctly, you will need to install the ``traffic_ops_scripted.js`` file from ``/traffic_control/traffic_stats/grafana`` to the ``/usr/share/grafana/public/dashboards/`` on the grafana server.  

	More information can be found in the `scripted dashboards <http://docs.grafana.org/reference/scripting/>`_ section of the Grafana documentation.

**Configuring Traffic Ops for Traffic Stats:**

	- The influxDb servers need to be added to Traffic Ops with profile = InfluxDb.  Make sure to use port 8086 in the configuration.
	- The traffic stats server should be added to Traffic Ops with profile = Traffic Stats.
	- Parameters for which stats will be collected are added with the release, but any changes can be made via parameters that are assigned to the Traffic Stats profile.
