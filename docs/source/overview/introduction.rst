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

Introduction
============
Traffic Control is a control plane for a CDN, which includes all of the components mentioned in the CDN Basics section, except for the Log File Analysis System. The caching software chosen for Traffic Control is `Apache Traffic Server <http://trafficserver.apache.org/>`_ (ATS). Although the current release only supports ATS as a cache, this may change with future releases. 

Traffic Control was first developed at Comcast for internal use and released to Open Source in April of 2015.

Traffic Control implements the blue boxes in the architecture diagram below. 


.. image:: traffic_control_overview_3.png
	:align: center

In the next sections each of these components will be explained further.
