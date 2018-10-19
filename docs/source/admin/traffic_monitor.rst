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

.. _tm-golang:

Installing Traffic Monitor
==========================

The following are hard requirements requirements for Traffic Monitor to operate:

* CentOS 6+
* Successful install of Traffic Ops (usually on a separate machine)
* Administrative access to the Traffic Ops (usually on a separate machine)


These are the recommended hardware specifications for a production deployment of Traffic Monitor:

* 8 CPUs
* 16GB of RAM
* It is also recommended that you know the physical address of the site where the Traffic Monitor machine lives for optimal performance

#. Enter the Traffic Monitor server into Traffic Portal

	.. note:: For legacy compatibility reasons, the 'Type' field of a new Traffic Monitor server must be 'RASCAL'.

#. Make sure the Fully Qualified Domain Name (FQDN) of the Traffic Monitor is resolvable in DNS.
#. Install Traffic Monitor, either from source or by running the command ``yum install traffic_monitor`` as the root user, or with ``sudo``.
#. Configure Traffic Monitor. See :ref:`here <tm-configure>`
#. Start the service, usually by running the command ``systemctl start traffic_monitor`` as the root user, or with ``sudo``
#. Verify Traffic Monitor is running by e.g. opening your preferred web browser to port 80 on the Traffic Monitor host.

Configuring Traffic Monitor
===========================

Configuration Overview
----------------------

.. _tm-configure:

Traffic Monitor is configured via two JSON configuration files, ``traffic_ops.cfg`` and ``traffic_monitor.cfg``, by default located in the ``conf`` directory in the install location. ``traffic_ops.cfg`` contains Traffic Ops connection information. Specify the URL, username, and password for the instance of Traffic Ops of which this Traffic Monitor is a member. ``traffic_monitor.cfg`` contains log file locations, as well as detailed application configuration variables such as processing flush times and initial poll intervals. Once started with the correct configuration, Traffic Monitor downloads its configuration from Traffic Ops and begins polling caches. Once every cache has been polled, :ref:`health-proto` state is available via RESTful JSON endpoints.


Troubleshooting and Log Files
=============================
Traffic Monitor log files are in ``/opt/traffic_monitor/var/log/``.
