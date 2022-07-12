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

.. _dev-traffic-stats:

*************
Traffic Stats
*************

Introduction
============
Traffic Stats is a utility written in `Go <http.golang.org>`_ that is used to acquire and store statistics about CDNs controlled by Traffic Control. Traffic Stats mines metrics from the :ref:`tm-api` and stores the data in `InfluxDB <http://influxdb.com>`_.  Data is typically stored in InfluxDB on a short-term basis (30 days or less) and is used to drive graphs created by `Grafana <http://grafana.org>`_ which are linked from Traffic Ops. Traffic Stats also calculates daily statistics from InfluxDb and stores them in the Traffic Ops database.

Software Requirements
=====================
* A \*nix (MacOS and Linux are most commonly used) environment
* `Go 1.7.x or above <https://golang.org/doc/install>`_
* Access to a working instance of Traffic Ops
* Access to a working instance of Traffic Monitor
* `InfluxDB version 1.0.0 or greater <https://influxdata.com/downloads>`_

Traffic Stats Project Tree Overview
=====================================
* :file:`traffic_stats/` - contains Go source files and files used to create the Traffic Stats RPM.

	* :file:`grafana/` - contains a javascript file which is installed on the Grafana server. This allows Traffic Ops to create custom dashboards for :term:`Delivery Services`, :term:cache server`\ s, etc.
	* :file:`influxdb_tools/` - contains :ref:`sync_ts_databases` and :ref:`create_ts_databases` which are helpful if you have multiple instances and they get out of sync with data.


Go Formatting Conventions
=========================
In general `Go fmt <https://golang.org/cmd/gofmt/>`_ is the standard for formatting Go code. It is also recommended to use `Go lint <https://github.com/golang/lint>`_.

Installing The Developer Environment
====================================
#. Clone the traffic_control repository using Git into a location accessible by your $GOPATH
#. Navigate to the :atc-file:`traffic_ops/v4-client` directory of your cloned repository. (This is the directory containing Traffic Ops client code used by Traffic Stats)
#. From the :atc-file:`traffic_ops/v4-client` directory, run ``go test`` to test the client code. This will run all unit tests for the client and return the results. If there are missing dependencies you will need to run ``go mod vendor -v`` to get the dependencies
#. Once the tests pass, run ``go install`` to build and install the Traffic Ops client package. This makes it accessible to Traffic Stats.
#. Navigate to your cloned repository under Traffic Stats
#. Run ``go build traffic_stats.go`` to build traffic_stats.  You will need to run ``go mod vendor -v`` for any missing dependencies.
