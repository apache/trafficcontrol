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

.. index::
  Traffic Ops - Default Profiles
  
.. _rl-to-default-profiles:

Traffic Ops - Default Profiles
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

Traffic Ops has the concept of :ref:`rl-working-with-profiles`, which are an integral function within Traffic Ops.  To get started, a set of default Traffic Ops profiles need to be imported into Traffic Ops
to get started to support Traffic Control components Traffic Router, Traffic Monitor, and Apache Traffic Server.

`Download Default Profiles from here <http://trafficcontrol.incubator.apache.org/downloads/profiles/>`_ 

.. _rl-to-profiles-min-needed:

Minimum Traffic Ops Profiles needed
-----------------------------------
   * EDGE_ATS_<version>_<platform>_PROFILE.traffic_ops
   * MID_ATS_<version>_<platform>_PROFILE.traffic_ops
   * TRAFFIC_MONITOR_PROFILE.traffic_ops
   * TRAFFIC_ROUTER_PROFILE.traffic_ops
   * TRAFFIC_STATS_PROFILE.traffic_ops
   


Steps to Import a Profile
-------------------------
1. Sign into Traffic Ops

2. Navigate to 'Parameters->Select Profile'

3. Click the "Import Profile" button at the bottom

4. Choose the specific profile you want to import from your download directory

5. Click 'Submit'

6. Continue these steps for each :ref:`rl-to-profiles-min-needed` above
