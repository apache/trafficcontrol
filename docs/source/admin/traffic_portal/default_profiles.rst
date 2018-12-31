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

.. _default-profiles:

****************
Default Profiles
****************
Traffic Ops has the concept of :ref:`working-with-profiles`, which are an integral component of Traffic Ops. To get started, a set of default Traffic Ops profiles are provided. These can be imported into Traffic Ops, and are required by the Traffic Control components Traffic Router, Traffic Monitor, and Apache Traffic Server (Edge-tier and Mid-tier caches). Download Default Profiles from `here <http://trafficcontrol.apache.org/downloads/profiles/>`_

.. _to-profiles-min-needed:

Minimum Traffic Ops Profiles Needed
===================================

- :file:`EDGE_ATS_{version}_{platform}_PROFILE.traffic_ops`
- :file:`MID_ATS_{version}_{platform}_PROFILE.traffic_ops`
- :file:`TRAFFIC_MONITOR_PROFILE.traffic_ops`
- :file:`TRAFFIC_ROUTER_PROFILE.traffic_ops`
- :file:`TRAFFIC_STATS_PROFILE.traffic_ops`
- :file:`EDGE_GROVE_PROFILE.traffic_ops`

.. note:: Despite that these have the ``.traffic_ops`` extension, they use JSON to store data. If your syntax highlighting doesn't work in some editor or viewer, try changing the extension to ``.json``.

.. warning:: These profiles will likely need to be modified to suit your system. Many of them contain hardware-specific parameters and parameter values.

Steps to Import a Profile
=========================
#. Sign into Traffic Portal
#. Navigate to :menuselection:`Configure --> Profiles`
#. Click on :menuselection:`More --> Import Profile`
#. Drag and drop your desired profile into the upload pane
#. Click :guilabel:`Import`
#. Continue these steps for each of the `Minimum Traffic Ops Profiles Needed`_.
