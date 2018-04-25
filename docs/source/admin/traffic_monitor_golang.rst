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

******************************
Traffic Monitor Administration
******************************

.. _rl-tm-golang:

* These instructions are for the Golang Traffic Monitor, for the legacy Java version, see :ref:`here <rl-tm-java>`.

Installing Traffic Monitor
==========================

The following are requirements to ensure an accurate set up:

* CentOS 6+
* 8 vCPUs
* 16GB RAM
* Successful install of Traffic Ops
* Administrative access to the Traffic Ops
* Physical address of the site

1. Enter the Traffic Monitor server into Traffic Ops
2. Make sure the FQDN of the Traffic Monitor is resolvable in DNS.
3. Install Traffic Monitor: ``sudo yum -y install traffic_monitor``
4. Configure Traffic Monitor. See :ref:`here <rl-tm-configure>`
5. Start the service: ``sudo service traffic_monitor start`` ::

    Starting traffic_monitor:

6. Verify Traffic Monitor is running by pointing your browser to port 80 on the Traffic Monitor host.

Configuring Traffic Monitor
===========================

Configuration Overview
----------------------

.. _rl-tm-configure:

Traffic Monitor is configured via two JSON configuration files, ``traffic_ops.cfg`` and ``traffic_monitor.cfg``, by default located in the ``conf`` directory in the install location.

The ``traffic_ops.cfg`` config contains Traffic Ops connection information. Specify the URL, username, and password for the instance of Traffic Ops for which this Traffic Monitor is a member.

The ``traffic_monitor.cfg`` config contains log file locations, as well as detailed application configuration variables, such as processing flush times and initial poll intervals.

Once started with the correct configuration, Traffic Monitor downloads its configuration from Traffic Ops and begins polling caches. Once every cache has been polled, health protocol state is available via RESTful JSON endpoints.


Troubleshooting and log files
=============================
Traffic Monitor log files are in ``/opt/traffic_monitor/var/log/``.
