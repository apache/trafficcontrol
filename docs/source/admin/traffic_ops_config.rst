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

Configuring Traffic Ops
%%%%%%%%%%%%%%%%%%%%%%%

Follow the steps below to configure the newly installed Traffic Ops Instance.

Installing the SSL Cert
=======================
By default, Traffic Ops runs as an SSL web server, and a certificate needs to be installed.  TBD.

Content Delivery Networks
=========================

.. _rl-param-prof:

Profile Parameters
======================
Many of the settings for the different servers in a Traffic Control CDN are controlled by parameters in the parameter view of Traffic Ops. Parameters are grouped in profiles and profiles are assigned to a server. For a typical cache there are hundreds of configuration settings to apply. The Traffic Ops parameter view contains the defined settings. To make life easier, Traffic Ops allows for duplication, comparison, import and export of Profiles. Traffic Ops also has a "Global profile" - the parameters in this profile are going to be applied to all servers in the Traffic Ops instance, or apply to Traffic Ops themselves. These parameters are:


.. index::
  Global Profile

+--------------------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------+
|           Name           |  Config file  |                                                                 Value                                                                 |
+==========================+===============+=======================================================================================================================================+
| tm.url                   | global        | The URL where this Traffic Ops instance is being served from.                                                                         |
+--------------------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------+
| tm.toolname              | global        | The name of the Traffic Ops tool. Usually "Traffic Ops". Used in the About screen and in the comments headers of the files generated. |
+--------------------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------+
| tm.infourl               | global        | This is the "for more information go here" URL, which is visible in the About page.                                                   |
+--------------------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------+
| tm.logourl               | global        | This is the URL of the logo for Traffic Ops and can be relative if the logo is under traffic_ops/app/public.                          |
+--------------------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------+
| tm.instance_name         | global        | The name of the Traffic Ops instance. Can be used when multiple instances are active. Visible in the About page.                      |
+--------------------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------+
| tm.traffic_mon_fwd_proxy | global        | When collecting stats from Traffic Monitor, Traffic Ops uses this forward proxy to pull the stats through.                            |
|                          |               | This can be any of the MID tier caches, or a forward cache specifically deployed for this purpose. Setting                            |
|                          |               | this variable can significantly lighten the load on the Traffic Monitor system and it is recommended to                               |
|                          |               | set this parameter on a production system.                                                                                            |
+--------------------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------+
| geolocation.polling.url  | CRConfig.json | The location to get the GeoLiteCity database from.                                                                                    |
+--------------------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------+
| geolocation6.polling.url | CRConfig.json | The location to get the IPv6 GeoLiteCity database from.                                                                               |
+--------------------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------+

These parameters should be set to reflect the local environment.


After running the postinstall script, Traffic Ops has the following profiles pre-loaded:

+----------+-------------------------------------------------------------------------------------------------+
|   Name   |                                           Description                                           |
+==========+=================================================================================================+
| EDGE1    | The profile to be applied to the latest supported version of ATS, when running as an EDGE cache |
+----------+-------------------------------------------------------------------------------------------------+
| TR1      | The profile to be applied to the latest version of Traffic Router                               |
+----------+-------------------------------------------------------------------------------------------------+
| TM1      | The profile to be applied to the latest version of Traffic Monitor                              |
+----------+-------------------------------------------------------------------------------------------------+
| MID1     | The profile to be applied to the latest supported version of ATS, when running as an MID cache  |
+----------+-------------------------------------------------------------------------------------------------+
| RIAK_ALL | Riak profile for all CDNs to be applied to the Traffic Vault servers                            |
+----------+-------------------------------------------------------------------------------------------------+

.. Note:: The Traffic Server profiles contain some information that is specific to the hardware being used (most notably the disk configuration), so some parameters will have to be changed to reflect your configuration. Future releases of Traffic Control will separate the hardware and software profiles so it is easier to "mix-and-match" different hardware configurations.

Below is a list of cache parameters that are likely to need changes from the default profiles shipped with Traffic Ops:

+--------------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------+
|           Name           |    Config file    |                                                       Description                                                       |
+==========================+===================+=========================================================================================================================+
| allow_ip                 | astats.config     | This is a comma separated  list of IPv4 CIDR blocks that will have access to the astats statistics on the caches.       |
|                          |                   | The Traffic Monitor IP addresses have to be included in this, if they are using IPv4 to monitor the caches.             |
+--------------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------+
| allow_ip6                | astats.config     | This is a comma separated  list of IPv6 CIDR blocks that will have access to the astats statistics on the caches.       |
|                          |                   | The Traffic Monitor IP addresses have to be included in this, if they are using IPv6 to monitor the caches.             |
+--------------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------+
| Drive_Prefix             | storage.config    | JvD/Jeff to supply blurb                                                                                                |
+--------------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------+
| Drive_Letters            | storage.config    | JvD/Jeff to supply blurb                                                                                                |
+--------------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------+
| purge_allow_ip           | ip_allow.config   | The IP address range that is allowed to execute the PURGE method on the caches (not related to :ref:`rl-purge`)         |
+--------------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------+
| health.threshold.loadavg | rascal.properties | The Unix load average at which Traffic Router will stop sending traffic to this cache                                   |
+--------------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------+
| health.threshold.\\      | rascal.properties | The amount of bandwidth that Traffic Router will try to keep available on the cache.                                    |
| availableBandwidthInKbps |                   | For example: "">1500000" means stop sending new traffic to this cache when traffic is at 8.5Gbps on a 10Gbps interface. |
+--------------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------+

Below is a list of Traffic Server plugins that need to be configured in the parameter table:

+------------------+---------------+------------------------------------------------------+------------------------------------------------------------------------------------------------------------+
|       Name       |  Config file  |                     Description                      |                                                  Details                                                   |
+==================+===============+======================================================+============================================================================================================+
| astats_over_http | package       | The package version for the astats_over_http plugin. | `astats_over_http <http://trafficcontrol.apache.org/downloads/index.html>`_                                  |
+------------------+---------------+------------------------------------------------------+------------------------------------------------------------------------------------------------------------+
| trafficserver    | package       | The package version for the trafficserver plugin.    | `trafficserver <http://trafficcontrol.apache.org/downloads/index.html>`_                                     |
+------------------+---------------+------------------------------------------------------+------------------------------------------------------------------------------------------------------------+
| regex_revalidate | plugin.config | The config to be used for regex_revalidate.          | `regex_revalidate <https://docs.trafficserver.apache.org/en/5.3.x/reference/plugins/regex_remap.en.html>`_ |
|                  |               | For example: --config regex_revalidate.config        |                                                                                                            |
+------------------+---------------+------------------------------------------------------+------------------------------------------------------------------------------------------------------------+
| remap_stats      | plugin.config | The config to be used for remap_stats.               | `remap_stats <https://github.com/apache/trafficserver/tree/master/plugins/experimental/remap_stats>`_      |
|                  |               | Value is left blank.                                 |                                                                                                            |
+------------------+---------------+------------------------------------------------------+------------------------------------------------------------------------------------------------------------+


Regions, Locations and Cache Groups
===================================
All servers have to have a `location`, which is their physical location. Each location is part of a `region`, and each region is part of a `division`. For Example, ``Denver`` could be a location in the ``Mile High`` region and that region could be part of the ``West`` division. Enter your divisions first in  `Misc->Divisions`, then enter the regions in `Misc->Regions`, referencing the divisions entered, and finally, enter the physical locations in `Misc->Locations`, referencing the regions entered. 

All servers also have to be part of a `cache group`. A cache group is a logical grouping of caches, that don't have to be in the same physical location (in fact, usually a cache group is spread across minimally 2 physical Locations for redundancy purposes), but share geo coordinates for content routing purposes. JvD to add more.



Configuring Content Purge
=========================
Content purge using ATS is not simple; there is no file system to delete files/directories from, and in large caches it can be hard to delete a simple regular expression from the cache. This is why Traffic Control uses the `Regex Revalidate Plugin <https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/regex_revalidate.en.html>`_ to purge content from the system. We don't actually remove the content, we have a check that gets run before each request on each cache to see if this request matches a list of regular expressions, and if it does, we force a revalidation to the origin, making the original content inaccessible. The regex_revalidate plugin will monitor it's config file, and will pick up changes to it without a `traffic_line -x` signal to ATS. Changes to this file need to be distributed to the highest tier (MID) caches in the CDN before they are distributed to the lower tiers, to prevent filling the lower tiers with the content that should be purged from the higher tiers without hitting the origin. This is why the ort script (see :ref:`reference-traffic-ops-ort`) will by default push out config changes to MID first, confirm that they have all been updated, and then push out the changes to the lower tiers. In large CDNs, this can make the distribution and time to activation of the purge too long, and because of that there is the option to not distribute the `regex_revalidate.config` file using the ort script, but to do this using other means. By default, Traffic Ops will use ort to distribute the `regex_revalidate.config` file. 

Content Purge is controlled by the following parameters in the profile of the cache:

+----------------------+-------------------------+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------+
|         Name         |       Config file       |                   Description                    |                                                                         Details                                                                         |
+======================+=========================+==================================================+=========================================================================================================================================================+
| location             | regex_revalidate.config | What location the file should be in on the cache | The presence of this parameter tells ort to distribute this file; delete this parameter from the profile if this file is distributed using other means. |
+----------------------+-------------------------+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------+
| maxRevalDurationDays | regex_revalidate.config | The maximum time a purge can be active           | To prevent a build up of many checks before each request, this is longest time the system will allow                                                    |
+----------------------+-------------------------+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------+
| regex_revalidate     | plugin.config           | The config to be used for regex_revalidate.      | `regex_revalidate <https://docs.trafficserver.apache.org/en/5.3.x/reference/plugins/regex_remap.en.html>`_                                              |
|                      |                         | For example: --config regex_revalidate.config    |                                                                                                                                                         |
+----------------------+-------------------------+--------------------------------------------------+---------------------------------------------------------------------------------------------------------------------------------------------------------+

Note that the TTL the adminstrator enters in the purge request should be longer than the TTL of the content to ensure the bad content will not be used. If the CDN is serving content of unknown, or unlimited TTL, the administrator should consider using `proxy-config-http-cache-guaranteed-min-lifetime <https://docs.trafficserver.apache.org/en/latest/admin-guide/files/records.config.en.html#proxy-config-http-cache-guaranteed-min-lifetime>`_ to limit the maximum time an object can be in the cache before it is considered stale, and set that to the same value as `maxRevalDurationDays` (Note that the former is in seconds and the latter is in days, so convert appropriately).



.. _Creating-CentOS-Kickstart:

Creating the CentOS Kickstart File
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
The kickstart file is a text file, containing a list of items, each identified by a keyword. You can create it by using the Kickstart Configurator application, or writing it from scratch. The Red Hat Enterprise Linux installation program also creates a sample kickstart file based on the options that you selected during installation. It is written to the file ``/root/anaconda-ks.cfg``. This file is editable using most text editors that can save files as ASCII text.

To generate ISO, the CentOS Kickstart is necessary:

1. Create a kickstart file.
2. Create a boot media with the kickstart file or make the kickstart file available on the network.
3. Make the installation tree available.
4. Start the kickstart installation.

Create a ks.src file in the root of the selection location. See the example below: 

::


 mkdir newdir
 cd newdir/
 cp -r ../centos65/* .
 vim ks.src
 vim isolinux/isolinux.cfg
 cd vim osversions.cfg
 vim osversions.cfg


This is a standard kickstart formatted file that the generate ISO process uses to create the kickstart (ks.cfg) file for the install. The generate ISO process uses the ks.src, overwriting any information set in the Generate ISO tab in Traffic Ops, creating ks.cfg.

.. Note:: Streamline your install folder for under 1GB, which assists in creating a CD.   

.. seealso:: For in-depth instructions, please see `Kickstart Installation <https://access.redhat.com/documentation/en-US/Red_Hat_Enterprise_Linux/6/html/Installation_Guide/s1-kickstart2-howuse.html>`_




