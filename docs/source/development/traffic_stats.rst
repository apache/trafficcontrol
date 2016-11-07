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

Traffic Stats
*************

Introduction
============
Traffic Stats is a utility written in `Go <http.golang.org>`_ that is used to acquire and store statistics about CDNs controlled by Traffic Control.  Traffic Stats mines metrics from Traffic Monitor's JSON APIs and stores the data in `InfluxDb <http://influxdb.com>`_.  Data is typically stored in InfluxDb on a short-term basis (30 days or less) and is used to drive graphs created by `Grafana <http://grafana.org>`_ which are linked from Traffic Ops.  Traffic Stats also calculates daily statistics from InfluxDb and stores them in the Traffic Ops database.

Software Requirements
=====================
To work on Traffic Stats you need a \*nix (MacOS and Linux are most commonly used) environment that has the following installed:

	* `Go 1.7.x or above <https://golang.org/doc/install>`_
	* Access to a working instance of Traffic Ops
	* Access to a working instance of Traffic Monitor
	* `InfluxDb version 1.0.0 or greater <https://influxdb.com/download/index.html>`_

Traffic Stats Project Tree Overview
=====================================
	* **traffic_control/traffic_stats** - contains Go source files and Files used to create the Traffic Stats rpm.
	* **traffic_control/traffic_stats/grafana/** - contains a javascript file which is installed on the grafana server.  This allows Traffic Ops to create custom dashboards for Delivery Services, Caches, etc.
	* **traffic_control/traffic_stats/influxdb_tools/** - contains one tool to create the databases and retention policies needed by Traffic Stats as well as continuous queries to downsample data; contains another tool to sync downsampled data between influxdb instances.  This is is helpful if you have multiple instances and they get out of sync with data.


Go Formatting Conventions
============================
In general `Go fmt <https://golang.org/cmd/gofmt/>`_ is the standard for formatting go code.  It is also recommended to use `Go lint <https://github.com/golang/lint>`_.

Installing The Developer Environment
====================================
To install the Traffic Ops Developer environment:

	- Clone the traffic_control repository using Git into a location accessible by your $GOPATH
	- Navigate to the traffic_ops/client directory of your cloned repository. (This is the directory containing Traffic Ops client code used by Traffic Stats)
	- From the traffic_ops/client directory run ``go test`` to test the client code.  This will run all unit tests for the client and return the results.  If there are missing dependencies you will need to run ``go get <dependency name>`` to get the dependency
	- Once the tests pass, run ``go install`` to build and install the Traffic Ops client package.  This makes it accessible to Traffic Stats.
	- Navigate to your cloned repo under Traffic Stats
	- Run ``go build traffic_stats.go`` to build traffic_stats.  You will need to run ``go get`` for any missing dependencies.


Test Cases
==========
	Currently there are no automated tests for Traffic Stats :( but pull requests are always welcome to fix this problem!

