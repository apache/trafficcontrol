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

.. _tc-ts:
.. |arrow| image:: fwda.png


Traffic Stats
=============
Traffic Stats is a program written in `Go <http://golang.org>`_ that is used to acquire and store statistics about CDNs controlled by Traffic Control. Traffic Stats mines metrics from Traffic Monitor's JSON APIs and stores the data in `InfluxDB <http://influxdb.com>`_. Data is typically stored in InfluxDB on a short-term basis (30 days or less). The data from InfluxDB is then used to drive graphs created by `Grafana <http://grafana.org>`_ - which are linked to from Traffic Ops - as well as provide data exposed through the Traffic Ops API. Traffic Stats performs two functions:

- Gathers stat data for Edge Caches and Delivery Services at a configurable interval (10 second default) from the Traffic Monitor API's and stores the data in InfluxDB
- Summarizes all of the stats once a day (around midnight UTC) and creates a daily report containing the Max Gbps Served and the Total Bytes Served.

Stat data is stored in three different databases:

	- ``cache_stats``: Stores data gathered from edge caches. The `measurements <https://influxdb.com/docs/v0.9/concepts/glossary.html#measurement>`_ stored by cache are: bandwidth, maxKbps, and ``client_connections`` (``ats.proxy.process.http.current_client_connections``). Cache Data is stored with `tags <https://influxdb.com/docs/v0.9/concepts/glossary.html#tag>`_ for hostname, Cache Group, and CDN. Data can be queried using tags.

	- ``deliveryservice_stats``: Stores data for Delivery Services. The measurements stored by delivery service are:

		- ``kbps``
		- ``status_4xx``
		- ``status_5xx``
		- ``tps_2xx``
		- ``tps_3xx``
		- ``tps_4xx``
		- ``tps_5xx``
		- ``tps_total``

		Delivery Service statistics are stored with tags for Cache Group, CDN, and Delivery Service ``xml_id``.

	- ``daily_stats``: Stores summary data for daily activities. The statistics that are currently summarized are Max Bandwidth and Bytes Served and they are stored by CDN.

------------

Traffic Stats does not influence overall CDN operation, but is required in order to display charts in Traffic Portal and the legacy Traffic Ops UI.
