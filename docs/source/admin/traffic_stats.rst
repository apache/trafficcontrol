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

Traffic Stats actually consists of three seperate components:  Traffic Stats, InfluxDB, and Grafana.  See below for information on installing and configuring each component as well as configuring the integration between the three and Traffic Ops.

Installation
========================

**Installing Traffic Stats:**

	- Download the Traffic Stats RPM from the traffic control `downloads <http://traffic-control-cdn.net/downloads/index.html>`_ page.
	- Copy the Traffic Stats RPM to your server
	- sudo rpm -ivh <traffic_stats rpm>

**Installing InfluxDB:**

	In order to store traffic stats data you will need to install InfluxDB.  It is recommended InfluxDB be installed in a 3 server cluster; VMs are acceptable. The documentation for installing InfluxDB can be found on the InfluxDB `website <https://influxdb.com/docs/v0.9/introduction/installation.html>`_.

**Installing Grafana:**

	Grafana is used to display Traffic Stats/InfluxDB data in Traffic Ops.  Grafana is typically run on the same server as Traffic Stats but this is not a requirement.  Grafana can be installed on any server that can access InfluxDB and be accessed by Traffic Ops.  Documentation on installing Grafana can be found `here <http://docs.grafana.org/installation/>`_.

Configuration
=========================

**Configuring Traffic Stats:**

	Traffic Stats' configuration file can be found in /opt/traffic_stats/conf/traffic_stats.cfg.
	The following values need to be configured:

	     - *toUser:* The user used to connect to Traffic Ops
	     - *toPasswd:*  The password to use when connecting to Traffic Ops
	     - *toUrl:*  The URL of the Traffic Ops server used by Traffic Stats
	     - *influxUser:*  The user to use when connecting to InfluxDB (if configured on InfluxDB, else leave default)
	     - *influxPassword:*  That password to use when connecting to InfluxDB (if configured, else leave blank)
	     - *polling interval:*  The interval at which Traffic Monitor is polled and stats are stored in InfluxDB
	     - *statusToMon:*  The status of Traffic Monitor to poll (poll ONLINE or OFFLINE traffic monitors)
	     - *seelogConfig:*  The absolute path of the seelong config file
	     - *dailySummaryPollingInterval:* The interval, in seconds, at which Traffic Stats checks to see if daily stats need to be computed and stored.
	     - *cacheRetentionPolicy:* The default retention policy for cache stats
	     - *dsRetentionPolicy:* The default retention policy for deliveryservice stats
	     - *dailySummaryRetentionPolicy:* The retention policy to be used for the daily stats

**Configuring InfluxDB:**

	It is HIGHLY recommended that InfluxDB be configured for clustering.  Documentation on clustering configuration can be found on the clustering page of the `InfluxDB Website <https://influxdb.com/docs/v0.9/concepts/clustering.html>`_.

	Once InfluxDB is installed and clustering is configured, Databases and Retention Policies need to be created.  Traffic Stats writes to three different databases: cache_stats, deliveryservice_stats, and daily_stats.  More information about the databases and what data is stored in each can be found on the `overview <../overview/traffic_stats.html>`_ page.

	By default the cache_stats and deliveryservice_stats databases are set to store raw data for 26 hours with a retention policy called "daily"; the daily_stats database is set to store data infinitely with a retention policy called "daily_stats".The following commands can be completed via the InfluxDB `client <https://influxdb.com/download/index.html>`_ or via the admin user interface (http://influxdb_url:8083).

	*Creating Databases:*
		``create database cache_stats``

		``create database deliveryservice_stats``

		``create database daily_stats``

	*Creating Retention Policies:*
		``create retention policy daily on cache_stats duration 26h replication 3 DEFAULT``

		``create retention policy daily on deliveryservice_stats duration 26h replication 3 DEFAULT``

		``create retention policy daily_stats on daily_stats duration INF replication 3 DEFAULT``


**Configuring Grafana:**

		In Traffic Ops the Health -> Graph View tab can be configured to display grafana graphs using influxDb data.  In order for this to work correctly, you will need two things 1) a parameter added to traffic ops with the graph URL (we will discuss later) and 2) the graphs created in grafana.  See below for how to create some simple graphs in grafana.  These instructions assume that InfluxDB has been installed and conifugred and that data has been written to it.  If this is not true, you will not see any graphs.

		- Login to grafana as an admin user http://grafana_url:3000/login
		- Choose Data Sources and then Add New
		- Name your data source (we name our data sources to match the database name, cache_stats and delivery_service stats)
		- Change the type to InfluxDB 0.9.x
		- For URL use https://grafana_url (see below on setting up the httpd proxy)
		- For Access choose 'direct'
		- Under the InfluxDB Details section enter the name of your database and enter a username and password for InfluxDB if you created one. If you did not create a username and password for influxdb just enter anything.
		- Click the 'Add' button to save the Data Source
		- Click on the 'Home' dropdown at the top of the screen and choose New at the bottom
		- Click on the green menu bar (with 3 lines) at the top and choose Add Panel -> Graph
		- Where it says 'No Title (click here)' click and choose edit
		- Choose your data source at the bottom
		- You can have grafana help you create a query, or you can create your own.  Here is a sample query:

			``SELECT sum(value)*1000/6 FROM "bandwidth" WHERE $timeFilter and time < now() - 60s GROUP BY time(60s), cdn``
		- Once you have the graph the way you want it, click the 'Save Dashboard' button at the top
		- You should now have a new saved graph

	In order for Traffic Ops users to see Grafana graphs, Grafana will need to allow anonymous access.  Information on how to configure anonymous access can be found on the configuration page of the `Grafana Website  <http://docs.grafana.org/installation/configuration/#authanonymous>`_.

	Traffic Ops uses custom dashboards to display information about individual delivery services or cache groups.  In order for the custom graphs to display correctly, the `traffic_ops_*.js <https://github.com/Comcast/traffic_control/blob/master/traffic_stats/grafana/>`_ files need to be in the ``/usr/share/grafana/public/dashboards/`` directory on the grafana server.  If your Grafana server is the same as your Traffic Stats server the RPM install process will take care of putting the files in place.  If your grafana server is different from your Traffic Stats server, you will need to manually copy the files to the correct directory.  

	More information on custom scripted graphs can be found in the `scripted dashboards <http://docs.grafana.org/reference/scripting/>`_ section of the Grafana documentation.

**Configuring httpd proxying for SSL**

	Currently InfluxDB does not support HTTPS for queries (should be implemented very soon).  Since Traffic Ops is HTTPS, we need to be able to make HTTPS requests to grafana and influxdb.  We can accomplish the need to use HTTPS by installing httpd with the mod_ssl plugin and then configuring proxying of grafana and influxdb HTTPS calls to HTTP. Below are the steps for setting up the HTTPS to HTTP proxy.  This should be performed on the same server that is running grafana. This is also useful if you are running InfluxDB with Private IP addresses.

	1. Download and install httpd  `from here <http://httpd.apache.org/download.cgi>`_
	2. Create SSL certs
	3. Install and configure mod_ssl per `this link <http://dev.antoinesolutions.com/apache-server/mod_ssl>`_
	4. Create a file called grafana_proxy.conf in the /etc/httpd/conf.d directory
	5. Add the following information to grafana_proxy.conf:

	::

				ProxyPass /dashboard http://localhost:3000/dashboard
				ProxyPass /css http://localhost:3000/css
				ProxyPass /app http://localhost:3000/app
				ProxyPass /api http://localhost:3000/api
				ProxyPass /img http://localhost:3000/img
				ProxyPass /fonts http://localhost:3000/fonts
				ProxyPass /public http://localhost:3000/public
				ProxyPass /login http://localhost:3000/login
				ProxyPass /logout http://localhost:3000/logout
				
				# The following ProxyPassReverse doesn't work for some.
				ProxyPassReverse / http://localhost:3000/

				<Proxy balancer://influxDb>
				BalancerMember http://<influxDb1>:8086
				BalancerMember http://<influxDb2>:8086
				BalancerMember http://<influxDb3>:8086
				</Proxy>
				ProxyPass /query balancer://influxDb/query
				
				# This works better for some
				ProxyPass / http://localhost:3000/

	6. Restart httpd ``service httpd restart``
	7. Test grafana works by connect to grafana via https ``https://grafanaUrl``


**Configuring Traffic Ops for Traffic Stats:**

	- The influxDb servers need to be added to Traffic Ops with profile = InfluxDB.  Make sure to use port 8086 in the configuration.
	- The traffic stats server should be added to Traffic Ops with profile = Traffic Stats.
	- Parameters for which stats will be collected are added with the release, but any changes can be made via parameters that are assigned to the Traffic Stats profile.

**Configuring Traffic Ops to use Grafana Dashboards**

	To configure Traffic Ops to use Grafana Dashboards, you need to enter the following parameters and assign them to the GLOBAL profile.  This assumes you followed the above instructions to install and configure InfluxDB and Grafana.  You will need to place 'cdn-stats' and 'deliveryservice-stats' with the name of your dashboards.

	+---------------------------+------------------------------------------------------------------------------------------------+
	|       parameter name      |                                        parameter value                                         |
	+===========================+================================================================================================+
	| all_graph_url             | https://<grafana_url>/dashboard/db/deliveryservice-stats                                       |
	+---------------------------+------------------------------------------------------------------------------------------------+
	| cachegroup_graph_url      | https://<grafanaHost>/dashboard/script/traffic_ops_cachegroup.js?which=                        |
	+---------------------------+------------------------------------------------------------------------------------------------+
	| deliveryservice_graph_url | https://<grafanaHost>/dashboard/script/traffic_ops_devliveryservice.js?which=                  |
	+---------------------------+------------------------------------------------------------------------------------------------+
	| server_graph_url          | https://<grafanaHost>/dashboard/script/traffic_ops_server.js?which=                            |
	+---------------------------+------------------------------------------------------------------------------------------------+
	| visual_status_panel_1     | https://<grafanaHost>/dashboard/solo/db/cdn-stats?panelId=2&fullscreen&from=now-24h&to=now-60s |
	+---------------------------+------------------------------------------------------------------------------------------------+
	| visual_status_panel_2     | https://<grafanaHost>/dashboard/solo/db/cdn-stats?panelId=1&fullscreen&from=now-24h&to=now-60s |
	+---------------------------+------------------------------------------------------------------------------------------------+
