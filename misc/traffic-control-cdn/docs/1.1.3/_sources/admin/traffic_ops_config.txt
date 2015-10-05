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

Configuring Traffic Ops
%%%%%%%%%%%%%%%%%%%%%%%

Follow the steps below to configure the newly installed Traffic Ops Instance.

Installing the SSL Cert
=======================
By default, Traffic Ops runs as an SSL web server, and a certificate needs to be installed.  TBD.

Content Delivery Networks
=========================

.. _rl-param-prof:

Parameters an profiles
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

..Note:: The Traffic Server profiles contain some information that is specific to the hardware being used (most notably the disk configuration), so some parameters will have to be changed to reflect your configuration. Future releases of Traffic Control will separate the hardware and software profiles so it is easier to "mix-and-match" different hardware configurations. 

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
| purge_allow_ip           | ip_allow.config   | The IP address that is allowed to "purge" content on the CDN through regex_revalidate                                   |
+--------------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------+
| health.threshold.loadavg | rascal.properties | The Unix load average at which Traffic Router will stop sending traffic to this cache                                   |
+--------------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------+
| health.threshold.\\      | rascal.properties | The amount of bandwidth that Traffic Router will try to keep available on the cache.                                    |
| availableBandwidthInKbps |                   | For example: "">1500000" means stop sending new traffic to this cache when traffic is at 8.5Gbps on a 10Gbps interface. |
+--------------------------+-------------------+-------------------------------------------------------------------------------------------------------------------------+


Regions, Locations and Cache Groups
===================================
All servers have to have a `location`, which is their physical location. Each location is part of a `region`, and each region is part of a `division`. For Example, ``Denver`` could be a location in the ``Mile High`` region and that region could be part of the ``West`` division. Enter your divisions first in  `Misc->Divisions`, then enter the regions in `Misc->Regions`, referencing the divisions entered, and finally, enter the physical locations in `Misc->Locations`, referencing the regions entered. 

All servers also have to be part of a `cache group`. A cache group is a logical grouping of caches, that don't have to be in the same physical location (in fact, usually a cache group is spread across minimally 2 physical locations for redundancy purposes), but share geo coordinates for content routing purposes. JvD to add more.






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




