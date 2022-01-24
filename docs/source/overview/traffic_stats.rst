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

.. _ts-overview:

*************
Traffic Stats
*************
:dfn:`Traffic Stats` is a program written in `Go <http://golang.org>`_ that is used to acquire and store statistics about CDNs controlled by Traffic Control. Traffic Stats mines metrics from the :ref:`tm-api` and stores the data in `InfluxDB <http://influxdb.com>`_ or `Kafka <https://kafka.apache.org/>`_. Data is typically stored in InfluxDB on a short-term basis (30 days or less), and Kafka data is available for consumption based on retention period. The data from InfluxDB is then used to drive graphs created by `Grafana <http://grafana.org>`_ - which are linked to from :ref:`tp-overview` - as well as provide data exposed through the :ref:`to-api`. Traffic Stats performs two functions:

- Gathers statistics for Edge-tier :term:`cache servers` and :term:`Delivery Services` at a configurable interval (10 second default) from the :ref:`tm-api` and stores the data in InfluxDB or Kafka
- Summarizes all of the statistics once a day (around midnight UTC) and creates a daily report containing the Max :abbr:`Gbps (Gigabits per second)` Served and the Total Bytes Served.

Statistics are stored in three different databases:

- ``cache_stats``: Stores data gathered from edge-tier :term:`cache servers`. The `measurements <https://influxdb.com/docs/v0.9/concepts/glossary.html#measurement>`_ stored by ``cache_stats`` are:

	- ``bandwidth``
	- ``maxKbps``
	- ``client_connections`` (``ats.proxy.process.http.current_client_connections``).

Cache Data is stored with `tags <https://influxdb.com/docs/v0.9/concepts/glossary.html#tag>`_ for hostname, :term:`Cache Group`, and CDN. Data can be queried using tags.

- ``deliveryservice_stats``: Stores data for :term:`Delivery Services`. The measurements stored by ``deliveryservice_stats`` are:

	- ``kbps``
	- ``status_4xx``
	- ``status_5xx``
	- ``tps_2xx``
	- ``tps_3xx``
	- ``tps_4xx``
	- ``tps_5xx``
	- ``tps_total``

:term:`Delivery Service` statistics are stored with tags for :term:`Cache Group`, CDN, and :term:`Delivery Service` :ref:`ds-xmlid`.

- ``daily_stats``: Stores summary data for daily activities. The statistics that are currently summarized are:

	- Max Bandwidth
	- Bytes Served

Daily stats are stored by CDN.

When Kafka is enabled, Cache and Delivery Service statistics are sent through JSON format with optional TLS authentication.

Traffic Stats does not influence overall CDN operation, but is required with InfluxDB enabled in order to display charts in :ref:`tp-overview`.