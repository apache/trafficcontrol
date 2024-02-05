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

.. _ts-admin:

****************************
Traffic Stats Administration
****************************
Traffic Stats consists of three separate components: Traffic Stats, InfluxDB, and Grafana. See below for information on installing and configuring each component as well as configuring the integration between the three and Traffic Ops.

Installation
============

Installing Traffic Stats
------------------------
- See the `downloads <https://trafficcontrol.apache.org/downloads/index.html>`_ page for Traffic Control to get the latest release.
- Follow the instructions in :ref:`dev-building` to generate an RPM.
- Copy the RPM to your server
- Install the generated Traffic Stats RPM with :manpage:`yum(8)` or :manpage:`rpm(8)`

Installing InfluxDB
-------------------
..  Note::As of Traffic Stats 1.8.0, InfluxDB 1.0.0 or higher is required. For InfluxDB versions less than 1.0.0 use Traffic Stats 1.7.x

In order to store Traffic Stats data you will need to install `InfluxDB <https://docs.influxdata.com/influxdb/latest/introduction/installation/>`_. While not required, it is recommended to use some sort of high availability option like `Influx enterprise <https://portal.influxdata.com/>`_, `InfluxDB Relay <https://github.com/influxdata/influxdb-relay>`_, or another `high availability option <https://www.influxdata.com/high-availability/>`_.

Installing Grafana
------------------
Grafana is used to display Traffic Stats/InfluxDB data in Traffic Ops. Grafana is typically run on the same server as Traffic Stats but this is not a requirement. Grafana can be installed on any server that can access InfluxDB and can be accessed by Traffic Ops. Documentation on installing Grafana can be found on the `Grafana website <http://docs.grafana.org/installation/>`__.

Configuration
=============

Configuring Traffic Stats
-------------------------
Traffic Stats's configuration file can be found in :file:`/opt/traffic_stats/conf/traffic_stats.cfg`. It is a JSON-format file with the following properties.

:toUser: The username of the user as whom to connect to Traffic Ops
:toPasswd: The password to use when authenticating with Traffic Ops
:toUrl: The URL of the Traffic Ops server used by Traffic Stats
:toRequestTimeout: The time, in seconds, before a client request to Traffic Ops is canceled. Defaults to 10 if no value provided
:influxUser: Optionally specify the user to use when connecting to InfluxDB (if InfluxDB is not configured to require authentication, has no effect)
:influxPassword: Optionally specify the password to use when connecting to InfluxDB (if InfluxDB is not configured to require authentication, has no effect)
:pollingInterval: The interval in seconds for which Traffic Monitor will wait between polling for stats and storing them in InfluxDB
:statusToMon: The name of the :term:`Status` of the Traffic Monitors to poll (e.g. ``ONLINE`` or ``OFFLINE``)
:seelogConfig: Optionally specify the absolute path to a `seelog <https://github.com/cihub/seelog>`_ configuration file. Has no effect if the ``logs`` property is present.

	.. deprecated:: ATCv6.1
		This will be removed in the future, configurations should be migrated to the new ``logs`` property.

:logs: This property is an object containing keys that specify locations for different log streams.

	.. versionadded:: ATCv6.1

	:error: The location of error-level logs. If omitted, ``null``, an empty string (``""``), or the special path ``/dev/null`` (even on Windows), error-level logs will not be emitted. The special values "stderr" and "stdout" cause logging to use STDOUT or STDERR, respectively\ [#logfiles]_.
	:warning: The location of warning-level logs. If omitted, ``null``, an empty string (``""``), or the special path ``/dev/null`` (even on Windows), warning-level logs will not be emitted. The special values "stderr" and "stdout" cause logging to use STDOUT or STDERR, respectively\ [#logfiles]_.
	:info: The location of info-level logs. If omitted, ``null``, an empty string (``""``), or the special path ``/dev/null`` (even on Windows), info-level logs will not be emitted. The special values "stderr" and "stdout" cause logging to use STDOUT or STDERR, respectively\ [#logfiles]_.
	:debug: The location of debug-level logs. If omitted, ``null``, an empty string (``""``), or the special path ``/dev/null`` (even on Windows), debug-level logs will not be emitted. The special values "stderr" and "stdout" cause logging to use STDOUT or STDERR, respectively\ [#logfiles]_.
	:event: The location of event-level logs. If omitted, ``null``, an empty string (``""``), or the special path ``/dev/null`` (even on Windows), event-level logs will not be emitted. The special values "stderr" and "stdout" cause logging to use STDOUT or STDERR, respectively\ [#logfiles]_.

		.. note:: At the time of this writing, Traffic Stats does not make use of the "event" log level.


:dailySummaryPollingInterval: The interval, in seconds, on which Traffic Stats checks to see if daily stats need to be computed and stored.
:cacheRetentionPolicy: The default retention policy for cache stats
:dsRetentionPolicy: The default retention policy for :term:`Delivery Service` statistics
:dailySummaryRetentionPolicy: The retention policy to be used for the daily statistics
:influxUrls: An array of InfluxDB hosts for Traffic Stats to write stats to.

Configuring InfluxDB
--------------------
As mentioned above, it is recommended that InfluxDB be running in some sort of high availability configuration. There are several ways to achieve high availability so it is best to consult the high availability options on the `InfuxDB website <https://www.influxdata.com/high-availability/>`_.

Once InfluxDB is installed and configured, databases and retention policies need to be created. Traffic Stats writes to three different databases: cache_stats, deliveryservice_stats, and daily_stats. More information about the databases and what data is stored in each can be found in the `Traffic Stats Overview <tc-ts>`_.

To easily create databases, retention policies, and continuous queries, run :program:`create_ts_databases` from the :file:`/opt/traffic_stats/influxdb_tools` directory on your Traffic Stats server. See the `InfluxDB Tools`_ section for more information.

.. _grafana-config:

Configuring Grafana
-------------------
Grafana can be configured to display graphs using InfluxDB data.
See below for how to create some simple graphs in Grafana. These instructions assume that InfluxDB has been configured and that data has been written to it. If this is not true, you will not see any graphs.

To create a graph in Grafana, you can follow these basic steps:

#. Login to Grafana as an administrative user
#. Click on :menuselection:`Data Sources --> Add New`
#. Enter the necessary information to configure your data source
#. Click on :menuselection:`Home --> New` at the bottom
#. Click on :menuselection:`"Collapsed Menu Icon" Button --> Add Panel --> Graph`
#. Where it says :guilabel:`No Title (click here)` click and choose edit
#. Choose your data source at the bottom
#. You can have Grafana help you create a query, or you can create your own.

	.. code-block:: postgresql
		:caption: Sample Query

		SELECT sum(value)*1000 FROM "monthly"."bandwidth.cdn.1min" GROUP BY time(60s), cdn;

#. Once you have the graph the way you want it, click the :guilabel:`Save Dashboard` button at the top
#. You should now have a new saved graph

Grafana uses Grafana Scenes to display information about individual :term:`Delivery Services` or :term:`Cache Groups`. In order for the custom graphs to display correctly, the built files of :atc-file:`traffic_stats/trafficcontrol-scenes/`  need to be placed in the :file:`/var/lib/grafana/plugins/trafficcontrol-scenes-app` directory on the Grafana server. If your Grafana server is the same as your Traffic Stats server the RPM install process will take care of putting the files in place. If your Grafana server is different from your Traffic Stats server, you will need to manually copy the files to the correct directory.

To view dynamic dashboards from Grafana Scenes, visit: ``https://grafanaHost/a/trafficcontrol-scenes-app``

.. seealso:: More information on Grafana Scenes can be found in the `blog post <https://grafana.com/blog/2023/09/12/grafana-scenes-is-generally-available-start-building-highly-interactive-apps-today/>`_ of Grafana.

Configuring Traffic Portal for Traffic Stats
--------------------------------------------
- The InfluxDB servers need to be added to Traffic Portal with a :term:`Profile` that has the :ref:`profile-type` InfluxDB. Make sure to use port 8086 in the configuration.
- The traffic stats server should be added to Traffic Ops with a :term:`Profile` that has the :ref:`profile-type` TRAFFIC_STATS.
- :term:`Parameters` for which stats will be collected are added with the release, but any changes can be made via :term:`Parameters` that are assigned to the Traffic Stats :term:`Profile`.

Configuring Traffic Portal to use Grafana Dashboards
----------------------------------------------------
To configure Traffic Portal to use Grafana Dashboards, you need to enter the following :term:`Parameters` and assign them to the special GLOBAL :term:`Profile`. This assumes you followed instructions in the Installation_, `Configuring Traffic Stats`_, `Configuring InfluxDB`_, and `Configuring Grafana`_ sections.

.. table:: Traffic Stats Parameters

	+---------------------------+--------------------------------------------------------------------------------------------------------------------+
	|       parameter name      |                                        parameter value                                                             |
	+===========================+====================================================================================================================+
	| all_graph_url             | :file:`https://{grafanaHost}/dashboard/db/{deliveryservice-stats-dashboard}`                                       |
	+---------------------------+--------------------------------------------------------------------------------------------------------------------+
	| cachegroup_graph_url      | :file:`https://{grafanaHost}/dashboard/script/traffic_ops_cachegroup.js?which=`                                    |
	+---------------------------+--------------------------------------------------------------------------------------------------------------------+
	| deliveryservice_graph_url | :file:`https://{grafanaHost}/dashboard/script/traffic_ops_deliveryservice.js?which=`                               |
	+---------------------------+--------------------------------------------------------------------------------------------------------------------+
	| server_graph_url          | :file:`https://{grafanaHost}/dashboard/script/traffic_ops_server.js?which=`                                        |
	+---------------------------+--------------------------------------------------------------------------------------------------------------------+
	| visual_status_panel_1     | :file:`https://{grafanaHost}/dashboard-solo/db/{cdn-stats-dashboard}?panelId=2&fullscreen&from=now-24h&to=now-60s` |
	+---------------------------+--------------------------------------------------------------------------------------------------------------------+
	| visual_status_panel_2     | :file:`https://{grafanaHost}/dashboard-solo/db/{cdn-stats-dashboard}?panelId=1&fullscreen&from=now-24h&to=now-60s` |
	+---------------------------+--------------------------------------------------------------------------------------------------------------------+
	| daily_bw_url              | :file:`https://{grafanaHost}/dashboard-solo/db/{daily-summary-dashboard}?panelId=1&fullscreen&from=now-3y&to=now`  |
	+---------------------------+--------------------------------------------------------------------------------------------------------------------+
	| daily_served_url          | :file:`https://{grafanaHost}/dashboard-solo/db/{daily-summary-dashboard}?panelId=2&fullscreen&from=now-3y&to=now`  |
	+---------------------------+--------------------------------------------------------------------------------------------------------------------+

where

grafanaHost
	is the :abbr:`FQDN (Fully Qualified Domain Name)` of the Grafana server (again, usually the same as the Traffic Stats server),
cdn-stats-dashboard
	is the name of the Dashboard providing CDN-level statistics,
deliveryservice-stats-dashboard
	is the name of the Dashboard providing :term:`Delivery Service`-level statistics, and
daily-summary-dashboard
	is the name of the Dashboard providing a daily summary of general statistics that would be of interest to administrators using Traffic Portal

InfluxDB Tools
==============
Under the Traffic Stats source directory there is a directory called ``influxdb_tools``. These tools are meant to be used as one-off scripts to help a user quickly get new databases and continuous queries setup in InfluxDB. They are specific for Traffic Stats and are not meant to be generic to InfluxDB. Below is an brief description of each script along with how to use it.

.. _create_ts_databases:

.. program:: create_ts_databases

create/create_ts_databases.go
-----------------------------
This program creates all `databases <https://docs.influxdata.com/influxdb/latest/concepts/key_concepts/#database>`_, `retention policies <https://docs.influxdata.com/influxdb/latest/concepts/key_concepts/#retention-policy>`_, and `continuous queries <https://docs.influxdata.com/influxdb/v0.11/query_language/continuous_queries/>`_ required by Traffic Stats.

Pre-Requisites
""""""""""""""
* Go 1.7 or later
* Configured ``$GOPATH`` environment variable

Options and Arguments
"""""""""""""""""""""
.. option:: --help

	(Optional) Print usage information and exit (with a failure exit code for some reason)

.. option:: --password password

	The password that will be used by the user defined by :option:`--user` to authenticate.

.. option:: --replication N

	(Optional) The number of nodes in the cluster (default: 3)

.. option:: --url URL

	The InfluxDB server's root URL - including port number, if required (default: ``http://localhost:8086``)

.. option:: --user username

	The name of the user to use when connecting to InfluxDB

.. _sync_ts_databases:

.. program:: sync_ts_databases

sync/sync_ts_databases.go
-------------------------
This program is used to sync one InfluxDB environment to another. Only data from continuous queries is synced as it is down-sampled data and much smaller in size than syncing raw data. Possible use cases are syncing from production to development or syncing a new cluster once brought online.

Pre-Requisites
""""""""""""""
* Go 1.7 or later
* Configured ``$GOPATH`` environment variable

Options and Arguments
"""""""""""""""""""""
.. option:: --database database_name

	(Optional) Specify the name of a specific database to sync (default: all databases)

.. option:: --days N

	The number of days in the past to sync. ``0`` means 'all'

.. option:: --help

	(Optional) Print usage information and exit

.. option:: --source-password password

	The password of the user named by :option:`--source-user`

.. option:: --source-url URL

	(Optional) The URL of the InfluxDB instance _from_ which data will be copied (default: ``http://localhost:8086``)

.. option:: --source-user username

	The name of the user as whom the utility will connect to the source InfluxDB instance

.. option:: --target-password password

	The password of the user named by :option:`--target-user`

.. option:: --target-url URL

	(Optional) The URL of the InfluxDB instance _to_ which data will be copied (default: ``http://localhost:8086``)

.. option:: --target-user username

	The name of the user as whom the utility will connect to the target InfluxDB instance

.. [#logfiles] To log to files named literally "stdout" or "stderr", use an absolute or relative file path e.g. "./stdout" or "/path/to/stderr".
