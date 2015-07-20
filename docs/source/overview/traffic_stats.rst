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

.. _reference-label-tc-ts:
.. |arrow| image:: fwda.png


Traffic Stats
=============
Traffic Stats is a collection of utilities written in `Go <http.golang.org>`_ that are used to acquire and store statistics about CDNs controlled by Traffic Control.  Traffic Stats mines metrics from Traffic Monitor's JSON APIs and stores the data in `InfluxDb <http://influxdb.com>`_.  Data is typically stored in InfluxDb on a short-term basis (30 days or less) and is used to drive graphs created by `Grafana <http://grafana.org>`_ which are linked from Traffic Ops.  Traffic Stats contains two seperate services: write traffic stats and ts daily summary.  See below for more information. 

|arrow| Write Traffic Stats
------------------------
Write Traffic Stats gathers stat data for Edge Caches and Delivery Services at a configurable interval from the Traffic Monitor API's and stores the data in InfluxDb. 

|arrow| TS Daily Summary
------------------------
TS (Traffic Stats) Daily Summary is a process that runs once a day, gathers summary data for the previous day from InfluxDb, and stores it in the Traffic Ops Database.  The stats that are currently summarized are Max Bandwidth and Bytes Served.


Any number of Traffic Stats instances may run on a given CDN to collect metrics from Traffic Monitor, however, integration with a long term metrics storage system is implementation dependent. 

Traffic Stats does not influence overall CDN operation, but is required in order to display charts in Traffic Operations.