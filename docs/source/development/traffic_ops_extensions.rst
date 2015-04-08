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

Traffic Ops Extensions
======================
Traffic Ops Extensions are a way to enhance the basic functionality of Traffic Ops in a custom manner. There are three types of extensions:

1. Check Extensions

  These allow you to add custom checks to the "Health->Server Checks" view.

2. Configuration Extensions

  These allow you to add custom configuration file generators.

3. Data source Extensions

  These allow you to add statistic sources for the graph views and APIs.

Extensions are managed using the $TO_HOME/bin/extensions command line script. For more information see :ref:`admin-to-ext-script`.

Check Extensions
----------------
The "Health->Server Checks" view is most commonly used to help deployment teams with their activities during initial bringup of the caches. It can however also be useful to detect problems after deployment. A check extension is a script that runs at a certain interval out of cron, and uses a column in the "Health->Server Checks" view to display the results. 

In other words, check extensions are scripts that, after registering with Traffic Ops, have a column reserved in the "Health->Server Checks" view and that usually run periodically out of cron.

.. |checkmark| image:: ../../../traffic_ops/app/public/images/good.png 

.. |X| image:: ../../../traffic_ops/app/public/images/bad.png


It is the responsibility of the check extension script to iterate over the servers it wants to check and post the results.  A check extension can have a column of |checkmark|'s and |X|'s (CHECK_EXTENSION_BOOL) or a column that shows a number (CHECK_EXTENSION_NUM). A simple example of a check extension of type CHECK_EXTENSION_NUM that will show 99.33 for all servers of type EDGE is shown below: :: 


  Script here.

Check Extension scripts are located in the $TO_HOME/bin/checks directory.

Currently, the following Check Extensions are available and installed by default:

**Cache Disk Usage Check - CDU**
  This check shows how much of the available total cache disk is in use. A "warm" cache should show 100.00.

**Cache Hit Ratio Check - CHR**
  The cache hit ratio for the cache in the last 15 minutes (the interval is determined by the cron entry). 

**DiffServe CodePoint Check - DSCP**
  Checks if the returning traffic from the cache has the correct DSCP value as assigned in the delivery service. (Some routers will overwrite DSCP)

**Maximum Transmission Check - MTU**
  Checks if the Traffic Ops host (if that is the one running the check) can send and receive 8192 size packets to the ``ip_address`` of the server in the server table.

**Operational Readiness Check - ORT**
  See :ref:`reference-traffic-ops-ort` for more information on the ort script. The ORT column shows how many changes the traffic_ops_ort.pl script would apply if it was run. The number in this column should be 0. 

**Ping Check - 10G, ILO, 10G6, FQDN**
  The bin/checks/ToPingCheck.pl is to check basic IP connectivity, and in the default setup it checks IP connectivity to the following:
  
  10G
    Is the ``ip_address`` (the main IPv4 address) from the server table pingable?
  ILO
    Is the ``ilo_ip_address`` (the lights-out-mangement IPv4 address) from the server table pingable?
  10G6
    Is the ``ip6_address`` (the main IPv6 address) from the server table pingable?
  FQDN 
    Is the Fully Qualified Domain name (the concatenation of ``host_name`` and ``.`` and ``domain_name`` from the server table) pingable?

**Traffic Router Check - RTR**
  

Configuration Extensions
------------------------
NOTE: Config Extensions are Beta at this time.


Data source Extensions
----------------------
NOTE: Data source Extensions are Beta at this time.



