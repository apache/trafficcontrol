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

Cache Polling URL
-----------------------------------

The :term:`cache servers` are polled at the URL specified in the ``health.polling.url`` :term:`parameter`, on the :term:`cache server`'s :term:`profile`.

This :term:`parameter` must have the config file ``rascal.properties``.

The value is a template with the text ``${hostname}`` being replaced with the :term:`cache server`'s Network IP (IPv4), and ``${interface_name}`` being replaced with the :term:`cache server`'s network Interface Name. For example, ``http://${hostname}/_astats?application=&inf.name=${interface_name}``.

If the template contains a port, that port will be used, and the :term:`cache server`'s HTTPS and TCP Ports will not be added.

If the template does not contain a port, then if the template starts with ``https`` the :term:`cache server`'s HTTPS Port will be added, and if the template doesn't start with ``https`` the :term:`cache server`'s TCP Port will be added.

Examples:

Template ``http://${hostname}/_astats?application=&inf.name=${interface_name}`` Server IP ``192.0.2.42`` Server TCP Port ``8080`` HTTPS Port ``8443`` becomes ``http://192.0.2.42:8080/_astats?application=&inf.name=${interface_name}``.
Template ``https://${hostname}/_astats?application=&inf.name=${interface_name}`` Server IP ``192.0.2.42`` Server TCP Port ``8080`` HTTPS Port ``8443`` becomes ``https://192.0.2.42:8443/_astats?application=&inf.name=${interface_name}``.
Template ``http://${hostname}:1234/_astats?application=&inf.name=${interface_name}`` Server IP ``192.0.2.42`` Server TCP Port ``8080`` HTTPS Port ``8443`` becomes ``http://192.0.2.42:1234/_astats?application=&inf.name=${interface_name}``.
Template ``https://${hostname}:1234/_astats?application=&inf.name=${interface_name}`` Server IP ``192.0.2.42`` Server TCP Port ``8080`` HTTPS Port ``8443`` becomes ``https://192.0.2.42:1234/_astats?application=&inf.name=${interface_name}``.

Stat and Health Flush Configuration
-----------------------------------
The Monitor has a health flush interval, a stat flush interval, and a stat buffer interval. Recall that the monitor polls both stats and health. The health poll is so small and fast, a buffer is largely unnecessary. However, in a large CDN, the stat poll may involve thousands of :term:`cache server` s with thousands of stats each, or more, and CPU may be a bottleneck.

The flush intervals, ``health_flush_interval_ms`` and ``stat_flush_interval_ms``, indicate how often to flush stats or health, if results are continuously coming in with no break. This prevents starvation. Ideally, if there is enough CPU, the flushes should never occur. The default flush times are 200 milliseconds, which is suggested as a reasonable starting point; operators may adjust them higher or lower depending on the need to get health data and stop directing client traffic to unhealthy :term:`cache server` s as quickly as possible, balanced by the need to reduce CPU usage.

The stat buffer interval, ``stat_buffer_interval_ms``, also provides a temporal buffer for stat processing. Stats will not be processed except after this interval, whereupon all pending stats will be processed, unless the flush interval occurs as a starvation safety. The stat buffer and flush intervals may be thought of as a state machine with two states: the "buffer state" accepts results until the buffer interval has elapsed, whereupon the "flush state" is entered, and results are accepted while outstanding, and processed either when no results are outstanding or the flush interval has elapsed.

Troubleshooting and Log Files
=============================
Traffic Monitor log files are in ``/opt/traffic_monitor/var/log/``.
