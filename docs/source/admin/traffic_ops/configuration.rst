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

.. _`regex_revalidate plugin`: https://docs.trafficserver.apache.org/en/7.1.x/admin-guide/plugins/regex_revalidate.en.html

*************************
Traffic Ops - Configuring
*************************
Follow the steps below to configure the newly installed Traffic Ops Instance.

Installing the SSL Certificate
==============================
By default, Traffic Ops runs as an SSL web server (that is, over HTTPS), and a certificate needs to be installed.

Self-signed Certificate (Development)
-------------------------------------
.. code-block:: console
	:caption: Example Procedure

	$ openssl genrsa -des3 -passout pass:x -out localhost.pass.key 2048
	Generating RSA private key, 2048 bit long modulus
	...
	$ openssl rsa -passin pass:x -in localhost.pass.key -out localhost.key
	writing RSA key
	$ rm localhost.pass.key

	$ openssl req -new -key localhost.key -out localhost.csr
	You are about to be asked to enter information that will be incorporated
	into your certificate request.
	What you are about to enter is what is called a Distinguished Name or a DN.
	There are quite a few fields but you can leave some blank
	For some fields there will be a default value,
	If you enter '.', the field will be left blank.
	-----
	Country Name (2 letter code) [XX]:US<enter>
	State or Province Name (full name) []:CO<enter>
	Locality Name (eg, city) [Default City]:Denver<enter>
	Organization Name (eg, company) [Default Company Ltd]: <enter>
	Organizational Unit Name (eg, section) []: <enter>
	Common Name (eg, your name or your server's hostname) []: <enter>
	Email Address []: <enter>

	Please enter the following 'extra' attributes
	to be sent with your certificate request
	A challenge password []: pass<enter>
	An optional company name []: <enter>
	$ openssl x509 -req -sha256 -days 365 -in localhost.csr -signkey localhost.key -out localhost.crt
	Signature ok
	subject=/C=US/ST=CO/L=Denver/O=Default Company Ltd
	Getting Private key
	$ sudo cp localhost.crt /etc/pki/tls/certs
	$ sudo cp localhost.key /etc/pki/tls/private
	$ sudo chown trafops:trafops /etc/pki/tls/certs/localhost.crt
	$ sudo chown trafops:trafops /etc/pki/tls/private/localhost.key

Certificate from Certificate Authority (Production)
---------------------------------------------------

.. Note:: You will need to know the appropriate answers when generating the certificate request file :file:`trafficopss.csr` below.

Example Procedure
"""""""""""""""""
.. code-block:: console
	:caption: Example Procedure

	$ openssl genrsa -des3 -passout pass:x -out trafficops.pass.key 2048
	Generating RSA private key, 2048 bit long modulus
	...
	$ openssl rsa -passin pass:x -in trafficops.pass.key -out trafficops.key
	writing RSA key
	$ rm localhost.pass.key

Generate the :abbr:`CSR (Certificate Signing Request)` file needed for :abbr:`CA (Certificate Authority)` request

.. code-block:: console
	:caption: Example Certificate Signing Request File Generation

	$ openssl req -new -key trafficops.key -out trafficops.csr
	You are about to be asked to enter information that will be incorporated
	into your certificate request.
	What you are about to enter is what is called a Distinguished Name or a DN.
	There are quite a few fields but you can leave some blank
	For some fields there will be a default value,
	If you enter '.', the field will be left blank.
	-----
	Country Name (2 letter code) [XX]: <enter country code>
	State or Province Name (full name) []: <enter state or province>
	Locality Name (eg, city) [Default City]: <enter locality name>
	Organization Name (eg, company) [Default Company Ltd]: <enter organization name>
	Organizational Unit Name (eg, section) []: <enter organizational unit name>
	Common Name (eg, your name or your server's hostname) []: <enter server's hostname name>
	Email Address []: <enter e-mail address>

	Please enter the following 'extra' attributes
	to be sent with your certificate request
	A challenge password []: <enter challenge password>
	An optional company name []: <enter>
	$ sudo cp trafficops.key /etc/pki/tls/private
	$ sudo chown trafops:trafops /etc/pki/tls/private/trafficops.key

You must then take the output file :file:`trafficops.csr` and submit a request to your :abbr:`CA (Certificate Authority)`. Once you get approved and receive your :file:`trafficops.crt` file

.. code-block:: shell
	:caption: Certificate Installation

	sudo cp trafficops.crt /etc/pki/tls/certs
	sudo chown trafops:trafops /etc/pki/tls/certs/trafficops.crt

If necessary, install the :abbr:`CA (Certificate Authority) certificate's ``.pem`` and ``.crt`` files in ``/etc/pki/tls/certs``.

You will need to update the file :file:`/opt/traffic_ops/app/conf/cdn.conf` with the any necessary changes.

.. code-block:: text
	:caption: Sample 'listen' Line When Path to ``trafficops.crt`` and ``trafficops.key`` are Known

	'hypnotoad' => ...
	    'listen' => 'https://[::]:443?cert=/etc/pki/tls/certs/trafficops.crt&key=/etc/pki/tls/private/trafficops.key&ca=/etc/pki/tls/certs/localhost.ca&verify=0x00&ciphers=AES128-GCM-SHA256:HIGH:!RC4:!MD5:!aNULL:!EDH:!ED'
		 ...


Content Delivery Networks
=========================

.. _param-prof:

Profile Parameters
------------------
Many of the settings for the different servers in a Traffic Control CDN are controlled by parameters in the :menuselection:`Configure --> Parameters` view of Traffic Portal. Parameters are grouped in profiles and profiles are assigned to a server or a :term:`Delivery Service`. For a typical cache there are hundreds of configuration settings to apply. The Traffic Portal :menuselection:`Parameters` view contains the defined settings. To make life easier, Traffic Portal allows for duplication, comparison, import and export of profiles. Traffic Ops also has a "Global profile" - the parameters in this profile are going to be applied to all servers in the Traffic Ops instance, or apply to Traffic Ops themselves. These parameters are explained in the :ref:`global-profile-parameters` table.

.. _global-profile-parameters:
.. table:: Global Profile Parameters

	+--------------------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------+
	|           Name           | ConfigFile    |                                                                 Value                                                                 |
	+==========================+===============+=======================================================================================================================================+
	| tm.url                   | global        | The URL at which this Traffic Ops instance services requests                                                                          |
	+--------------------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------+
	| tm.rev_proxy.url         | global        | Not required. The URL where a caching proxy for configuration files generated by Traffic Ops may be found. Requires a minimum         |
	|                          |               | :term:`ORT` version of 2.1. When configured, :term:`ORT` will request configuration files via this                                    |
	|                          |               | :abbr:`FQDN (Fully Qualified Domain Name)`, which should be set up as a reverse proxy to the Traffic Ops server(s). The suggested     |
	|                          |               | cache lifetime for these files is 3 minutes or less. This setting allows for greater scalability of a CDN maintained by Traffic Ops   |
	|                          |               | by caching configuration files of profile and CDN scope, as generating these is a very computationally expensive process              |
	+--------------------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------+
	| tm.toolname              | global        | The name of the Traffic Ops tool. Usually "Traffic Ops" - this will appear in the comment headers of the generated files              |
	+--------------------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------+
	| tm.infourl               | global        | This is the "for more information go here" URL, which used to be visible in the "About" page of the now-deprecated Traffic Ops UI     |
	+--------------------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------+
	| tm.logourl               | global        | This is the URL of the logo for Traffic Ops and can be relative if the logo is under :file:`traffic_ops/app/public`                   |
	+--------------------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------+
	| tm.instance_name         | global        | The name of the Traffic Ops instance - typically to distinguish instances when multiple are active                                    |
	+--------------------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------+
	| tm.traffic_mon_fwd_proxy | global        | When collecting stats from Traffic Monitor, Traffic Ops will use this forward proxy instead of the actual Traffic Monitor host.       |
	|                          |               | This can be any of the MID tier caches, or a forward cache specifically deployed for this purpose. Setting                            |
	|                          |               | this variable can significantly lighten the load on the Traffic Monitor system and it is recommended to                               |
	|                          |               | set this parameter on a production system.                                                                                            |
	+--------------------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------+
	| geolocation.polling.url  | CRConfig.json | The location of a geographic IP mapping database for Traffic Router instances to use                                                  |
	+--------------------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------+
	| geolocation6.polling.url | CRConfig.json | The location of a geographic IPv6 mapping database for Traffic Router instances to use                                                |
	+--------------------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------+
	| maxmind.default.override | CRConfig.json | The destination geographic coordinates to use for client location when the geographic IP mapping database returns a default location  |
	|                          |               | that matches the country code. This parameter can be specified multiple times with different values to support default overrides for  |
	|                          |               | multiple countries. The reason for the name "maxmind" is because MaxMind's GeoIP2 database is the default geographic IP mapping       |
	|                          |               | database implementation used by Comcast production servers (and the only officially supported implementation at the time of this      |
	|                          |               | writing). The format of this Parameter's value is: ``<Country Code>;<Latitude>,<Longitude>``, e.g. ``US;37.751,-97.822``              |
	+--------------------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------+

These parameters should be set to reflect the local environment.

After running the :program:`postinstall` script, Traffic Ops has the :ref:`tbl-default-profiles` pre-loaded.

.. _tbl-default-profiles:
.. table:: Default Profiles

	+----------+-------------------------------------------------------------------------------------------------------------------------------------------------+
	|   Name   | Description                                                                                                                                     |
	+==========+=================================================================================================================================================+
	| EDGE1    | The profile to be applied to the latest supported version of :abbr:`ATS (Apache Traffic Server)`, when running as an Edge-tier cache            |
	+----------+-------------------------------------------------------------------------------------------------------------------------------------------------+
	| TR1      | The profile to be applied to the latest version of Traffic Router                                                                               |
	+----------+-------------------------------------------------------------------------------------------------------------------------------------------------+
	| TM1      | The profile to be applied to the latest version of Traffic Monitor                                                                              |
	+----------+-------------------------------------------------------------------------------------------------------------------------------------------------+
	| MID1     | The profile to be applied to the latest supported version of :abbr:`ATS (Apache Traffic Server)`, when running as a Mid-tier cache              |
	+----------+-------------------------------------------------------------------------------------------------------------------------------------------------+
	| RIAK_ALL | "Riak" profile for all CDNs to be applied to the Traffic Vault servers ("Riak" being the name of the underlying database used by Traffic Vault) |
	+----------+-------------------------------------------------------------------------------------------------------------------------------------------------+

.. Note:: The "EDGE1" and "MID1" profiles contain some information that is specific to the hardware being used (most notably the disk configuration), so some parameters will have to be changed to reflect your configuration. Future releases of Traffic Control will separate the hardware and software profiles so it is easier to "mix-and-match" different hardware configurations. The :ref:`cache-server-hardware-parameters` table tabulates the cache parameters that are likely to need changes from the default profiles shipped with Traffic Ops.

.. _cache-server-hardware-parameters:
.. table:: Cache Server Hardware Parameters

	+-------------------------------------------+-------------------+----------------------------------------------------------------------------------------------------------------------------------------------+
	| Name                                      | ConfigFile        | Description                                                                                                                                  |
	+===========================================+===================+==============================================================================================================================================+
	| allow_ip                                  | astats.config     | This is a comma-separated list of IPv4 :abbr:`CIDR (Classless Inter-Domain Routing)` blocks that will have access to the 'astats' statistics |
	|                                           |                   | on the cache servers. The Traffic Monitor IP addresses have to be included in this if they are using IPv4 to monitor the cache servers       |
	+-------------------------------------------+-------------------+----------------------------------------------------------------------------------------------------------------------------------------------+
	| allow_ip6                                 | astats.config     | This is a comma-separated list of IPv6 :abbr:`CIDR (Classless Inter-Domain Routing)` blocks that will have access to the 'astats' statistics |
	|                                           |                   | on the cache servers. The Traffic Monitor IP addresses have to be included in this if they are using IPv6 to monitor the cache servers       |
	+-------------------------------------------+-------------------+----------------------------------------------------------------------------------------------------------------------------------------------+
	| Drive_Prefix                              | storage.config    | The device path start of the disks. For example, if storage devices ``/dev/sda`` through ``/dev/sdf`` are to be used for caching, this       |
	|                                           |                   | should be set to ``/dev/sd``                                                                                                                 |
	+-------------------------------------------+-------------------+----------------------------------------------------------------------------------------------------------------------------------------------+
	| Drive_Letters                             | storage.config    | A comma-separated list of the letter part of the storage devices to be used for caching. For example, if storage devices ``/dev/sda``        |
	|                                           |                   | through ``/dev/sdf`` are to be used for caching, this should be set to ``a,b,c,d,e,f``                                                       |
	+-------------------------------------------+-------------------+----------------------------------------------------------------------------------------------------------------------------------------------+
	| purge_allow_ip                            | ip_allow.config   | The IP address range that is allowed to execute the PURGE method on the caches (not related to :ref:`purge`)                                 |
	+-------------------------------------------+-------------------+----------------------------------------------------------------------------------------------------------------------------------------------+
	| coalesce_masklen_v4	                    | ip_allow.config   | The mask length to use when coalescing IPv4 networks into one line using                                                                     |
	|                                           |                   | `the NetAddr\:\:IP Perl library <http://search.cpan.org/~miker/NetAddr-IP-4.078/IP.pm>`_                                                     |
	+-------------------------------------------+-------------------+----------------------------------------------------------------------------------------------------------------------------------------------+
	| coalesce_number_v4 	                    | ip_allow.config   | The number to use when coalescing IPv4 networks into one line using                                                                          |
	|                                           |                   | `the NetAddr\:\:IP Perl library <http://search.cpan.org/~miker/NetAddr-IP-4.078/IP.pm>`_                                                     |
	+-------------------------------------------+-------------------+----------------------------------------------------------------------------------------------------------------------------------------------+
	| coalesce_masklen_v6	                    | ip_allow.config   | The mask length to use when coalescing IPv6 networks into one line using                                                                     |
	|                                           |                   | `the NetAddr\:\:IP Perl library. <http://search.cpan.org/~miker/NetAddr-IP-4.078/IP.pm>`_                                                    |
	+-------------------------------------------+-------------------+----------------------------------------------------------------------------------------------------------------------------------------------+
	| health.threshold.loadavg                  | rascal.properties | The Unix 'load average' (as given by :manpage:`uptime(1)`) at which Traffic Router will stop sending traffic to this cache                   |
	+-------------------------------------------+-------------------+----------------------------------------------------------------------------------------------------------------------------------------------+
	| health.threshold.availableBandwidthInKbps | rascal.properties | The amount of bandwidth (in kilobits per second) that Traffic Router will try to keep available on the cache. For example ">1500000" means   |
	|                                           |                   | "stop sending new traffic to this cache server when traffic is at 8.5Gbps on a 10Gbps interface"                                             |
	+-------------------------------------------+-------------------+----------------------------------------------------------------------------------------------------------------------------------------------+

The :ref:`plugin-parameters` table contains all Traffic Server plug-ins that must be configured as global parameters.

.. _plugin-parameters:
.. table:: Plugin Parameters

	+------------------+---------------+-------------------------------------------------------------------------------------------------------------------------------------------+
	|       Name       | ConfigFile    | Description                                                                                                                               |
	+==================+===============+===========================================================================================================================================+
	| astats_over_http | package       | The package version for the :abbr:`ATS (Apache Traffic Server)`                                                                           |
	|                  |               | `astats_over_http plugin <https://github.com/apache/trafficcontrol/tree/master/traffic_server/plugins/astats_over_http>`_                 |
	+------------------+---------------+----------------------------------------------------------+--------------------------------------------------------------------------------+
	| trafficserver    | package       | The package version of :abbr:`ATS (Apache Traffic Server)`                                                                                |
	+------------------+---------------+----------------------------------------------------------+--------------------------------------------------------------------------------+
	| regex_revalidate | plugin.config | The configuration to be used for the :abbr:`ATS (Apache Traffic Server)` `regex_revalidate plugin`_                                       |
	+------------------+---------------+----------------------------------------------------------+--------------------------------------------------------------------------------+
	| remap_stats      | plugin.config | The configuration to be used for the :abbr:`ATS (Apache Traffic Server)`                                                                  |
	|                  |               | `remap_stats plugin <https://github.com/apache/trafficserver/tree/master/plugins/experimental/remap_stats>`_ - value should be left blank |
	+------------------+---------------+-------------------------------------------------------------------------------------------------------------------------------------------+

Cache server parameters for special configurations, which are unlikely to need changes but may be useful in particular circumstances, may be found in the :ref:`special-parameters` table.

.. _special-parameters:
.. table:: Special Parameters

	+--------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| Name         | ConfigFile    | Description                                                                                                                                                                     |
	+==============+===================+=============================================================================================================================================================================+
	| not_a_parent | parent.config | This is a boolean flag and is considered ``true`` if it exists and has any value except ``false``. This prevents cache servers with this parameter in their profile from being  |
	|              |               | inserted into the ``parent.config`` files generated for other cache servers that have the affected cache server(s)'s :term:`Cache Group` as a parent of their own               |
	|              |               | :term:`Cache Group`. This is primarily useful for when Edge-tier cache servers are configured to have a :term:`Cache Group` of other Edge-tier cache servers as parents (a      |
	|              |               | highly unusual configuration), and it is necessary to exclude some - but not all - Edge-tier cache servers in the parent :term:`Cache Group` from the ``parent.config`` (for    |
	|              |               | example because they lack necessary capabilities), but still have all Edge-tier cache servers in the same :term:`Cache Group` in order to take traffic from ordinary            |
	|              |               | :term:`Delivery Service`\ s at that :term:`Cache Group`\ 's geographic location. Once again, this is a highly unusual scenario, and under ordinary circumstances this parameter |
	|              |               | should not exist.                                                                                                                                                               |
	+--------------+---------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+

Regions, Locations and Cache Groups
===================================
All servers have to have a :term:`Physical Location`, which defines their geographic latitude and longitude. Each :term:`Physical Location` is part of a :term:`Region`, and each :term:`Region` is part of a :term:`Division`. For example, ``Denver`` could be the name of a :term:`Physical Location` in the ``Mile High`` :term:`Region` and that :term:`Region` could be part of the ``West`` :term:`Division`. The hierarchy between these terms is illustrated graphically in :ref:`topography-hierarchy`.

.. _topography-hierarchy:
.. figure:: images/topography.*
	:align: center
	:alt: A graphic illustrating the hierarchy exhibited by topological groupings
	:figwidth: 25%

	Topography Hierarchy

To create these structures in Traffic Portal, first make at least one :term:`Division` under :menuselection:`Topology --> Divisions`. Next enter the desired :term:`Region`\ (s) in :menuselection:`Topology --> Regions`, referencing the earlier-entered :term:`Division`\ (s). Finally, enter the desired :term:`Physical Location`\ (s) in :menuselection:`Topology --> Phys Locations`, referencing the earlier-entered :term:`Region`\ (s).

All servers also have to be part of a :term:`Cache Group`. A :term:`Cache Group` is a logical grouping of cache servers, that don't have to be in the same :term:`Physical Location` (in fact, usually a :term:`Cache Group` is spread across minimally two :term:`Physical Location`\ s for redundancy purposes), but share geographical coordinates for content routing purposes.

Configuring Content Purge
=========================
Purging cached content using :abbr:`ATS (Apache Traffic Server)` is not simple; there is no file system from which to delete files and/or directories, and in large caches it can be hard to delete content matching a simple regular expression from the cache. This is why Traffic Control uses the `Regex Revalidate Plugin <https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/regex_revalidate.en.html>`_ to purge content from the cache. The cached content is not actually removed, instead a check that runs before each request on each cache server is serviced to see if this request matches a list of regular expressions. If it does, the cache server is forced to send the request upstream to its parents (possibly other caches, possibly the origin) without checking for the response in its cache. The Regex Revalidate Plugin will monitor its configuration file, and will pick up changes to it without needing to alert :abbr:`ATS (Apache Traffic Server). Changes to this file need to be distributed to the highest tier (Mid-tier) cache servers in the CDN before they are distributed to the lower tiers, to prevent filling the lower tiers with the content that should be purged from the higher tiers without hitting the origin. This is why the :term:`ORT` script will - by default - push out configuration changes to Mid-tier cache servers first, confirm that they have all been updated, and then push out the changes to the lower tiers. In large CDNs, this can make the distribution and time to activation of the purge too long, and because of that there is the option to not distribute the ``regex_revalidate.config`` file using the :term:`ORT` script, but to do this using other means. By default, Traffic Ops will use :term:`ORT` to distribute the ``regex_revalidate.config`` file. Content Purge is controlled by the parameters in the profile of the cache server specified in the :ref:`content-purge-parameters` table.

.. _content-purge-parameters:
.. table:: Content Purge Parameters

	+----------------------+-------------------------+----------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| Name                 | ConfigFile              | Description                                                                                                                                                    |
	+======================+=========================+================================================================================================================================================================+
	| location             | regex_revalidate.config | Where in the file system the ``regex_revalidate.config`` file should located on the cache server. The presence of this parameter tells ORT to distribute this  |
	|                      |                         | file; delete this parameter from the profile if this file is distributed using other means                                                                     |
	+----------------------+-------------------------+----------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| maxRevalDurationDays | regex_revalidate.config | The maximum duration for which a purge shall be active. To prevent a build-up of many checks before each request, this is longest duration (in days) for which |
	|                      |                         | the system will allow content purges to remain active                                                                                                          |
	+----------------------+-------------------------+--------------------------------------------------+-------------------------------------------------------------------------------------------------------------+
	| regex_revalidate     | plugin.config           | The configuration to be used for the `regex_revalidate plugin`_                                                                                                |
	+----------------------+-------------------------+--------------------------------------------------+-------------------------------------------------------------------------------------------------------------+
	| use_reval_pending    | global                  | Configures Traffic Ops to use a separate ``reval_pending`` flag for each cache server. When this flag is in use :term:`ORT` will check for a new               |
	|                      |                         | ``regex_revalidate.config`` every 60 seconds in "SYNCDS" mode during the dispersal timer. This will also allow :term:`ORT` to be run in "REVALIDATE" mode,     |
	|                      |                         | which will check for and clear the ``reval_pending`` flag. This can be set to run via :manpage:`cron(8)` task. Enable with a value of ``1``.                   |
	+----------------------+-------------------------+--------------------------------------------------+-------------------------------------------------------------------------------------------------------------+

.. versionadded:: 2.1
	``use_reval_pending`` was unavailable prior to Traffic Ops version 2.1.


.. Note:: The :abbr:`TTL (Time To Live)` entered by the administrator in the purge request should be longer than the :abbr:`TTL (Time To Live)` of the content to ensure the bad content will not be used. If the CDN is serving content of unknown, or unlimited :abbr:`TTL (Time To Live)`, the administrator should consider using `proxy-config-http-cache-guaranteed-min-lifetime <https://docs.trafficserver.apache.org/en/latest/admin-guide/files/records.config.en.html#proxy-config-http-cache-guaranteed-min-lifetime>`_ to limit the maximum time an object can be in the cache before it is considered stale, and set that to the same value as `maxRevalDurationDays` (Note that the former is in seconds and the latter is in days, so convert appropriately).

.. _Creating-CentOS-Kickstart:

Creating the CentOS Kickstart File
==================================
The Kickstart file is a text file, containing a list of items, each identified by a keyword. This file can be generated using the `Red Hat Kickstart Configurator application <https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/5/html/installation_guide/ch-redhat-config-kickstart>`_, or it can be written from scratch. The Red Hat Enterprise Linux installation program also creates a sample Kickstart file based on the options selected during installation. It is written to the file :file:`/root/anaconda-ks.cfg` in this case. This file is editable using most text editors.

Generating a System Image
-------------------------
#. Create a Kickstart file.
#. Create a boot media with the Kickstart file or make the Kickstart file available on the network.
#. Make the installation tree available.
#. Start the Kickstart installation.

.. code-block:: shell
	:caption: Creating a New System Image Definition Tree from an Existing One

	# Starting from the Kickstart root directory (`/var/www/files` by default)
	mkdir newdir
	cd newdir/

	# In this example, the pre-existing system image definition tree is for CentOS 7.4 located in `centos74`
	cp -r ../centos74/* .
	vim ks.src
	vim isolinux/isolinux.cfg
	cd ..
	vim osversions.cfg

:file:`ks.src` is a standard, Kickstart-formatted file that the will be used to create the Kickstart (ks.cfg) file for the install whenever a system image is generated from the source tree. :file:`ks.src` is a template - it will be overwritten by any information set in the form submitted from :menuselection:`Tools --> Generate ISO` in Traffic Portal. Ultimately, the two are combined to create the final Kickstart file (:file:`ks.cfg`).

.. Note:: It is highly recommended for ease of use that the system image source trees be kept under 1GB in size.

.. seealso:: For in-depth instructions, please see `Kickstart Installation <https://access.redhat.com/documentation/en-US/Red_Hat_Enterprise_Linux/6/html/Installation_Guide/s1-kickstart2-howuse.html>`_ in the Red Hat documentation.


Configuring the Go Application
==============================
Traffic Ops is in the process of migrating from Perl to Go, and currently runs as two applications. The Go application serves all endpoints which have been rewritten in the Go language, and transparently proxies all other requests to the old Perl application. Both applications are installed by the RPM, and both run as a single :manpage:`systemd(1)` service. When the project has fully migrated to Go, the Perl application will be removed, and the RPM and service will consist solely of the Go application.

By default, the :program:`postinstall` script configures the Go application to behave and transparently serve as the old Perl Traffic Ops did in previous versions. This includes reading the old ``cdn.conf`` and ``database.conf`` config files, and logging to the old ``access.log`` location. However, the Go Traffic Ops application may be customized by passing the command-line flag, ``-oldcfg=false``. By default, it will then look for a configuration file at :file:`/opt/traffic_ops/conf/traffic_ops_golang.config`. The new configuration file location may also be customized via the ``-cfg`` flag. A sample configuration file is installed by the RPM at :file:`/opt/traffic_ops/conf/traffic_ops_golang.config`. The new Go Traffic Ops application as a :manpage:`systemd(1)` service with a new configuration file, the ``-oldcfg=false`` and  ``-cfg`` flags may be added to the ``start`` function in the service file, located by default at :file:`/etc/init.d/traffic_ops`.
