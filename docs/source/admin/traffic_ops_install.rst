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

.. index::
  Traffic Ops - Installing 
  
Installing Traffic Ops
%%%%%%%%%%%%%%%%%%%%%%

System Requirements
-------------------
The user must have the following for a successful install:

* CentOS 6
* 4 vCPUs
* 32GB RAM
* 20 GB disk space
* YUM repository with minimally the following dependecies avaliable

  * apr 1.3.9-5 
  * apr-util 1.3.9-3 
  * apr-util-ldap 1.3.9-3   
  * expat-devel 2.0.1-11 
  * genisoimage 1.1.9-12  
  * httpd 2.2.15
  * httpd-tools 2.2.15  
  * libpcap-devel 14:1.4
  * mod_ssl  1:2.2.15-29
  * mysql 5.1.71 
  * autoconf 2.63-5.1.
  * automake 1.11.1-4
  * gcc 4.4.7-4
  * gettext 0.17-16
  * libcurl-devel 7.19.7-37
  * libtool 2.2.6-15.5
  * mysql-devel 5.1.73-3
  * perl-CPAN 1.9402-136
  * libcurl 7.19.7-37
  * openssl 1.0.1e-30
  * cloog-ppl 0.15.7-1.2
  * cpp 4.4.7-4
  * cvs 1.11.23-16
  * libgomp 4.4.7-4
  * libidn-devel 1.18-2
  * m4 1.4.13-5
  * mpfr 2.4.1-6
  * perl-Digest-SHA 1:5.47-136
  * ppl 0.10.2-11
  * curl 7.19.7-37
  * openssl-devel 1.0.1e-30
 
* Access to `The Comprehensive Perl Archive Network (CPAN) <http://www.cpan.org/>`_

.. Note:: The above versions are known to work on CentOS 6.5. Higher versions may work.

.. Note:: Although Traffic Ops supports both MySQL and Postgres as a database, support for MySQL is more mature and better tested. It is best to use MySQL when first getting started, and the rest of this quide assumes MySQL as the database.

Navigating the Install
-----------------------
To begin the install:

1. Install Traffipc Ops: ``sudo yum install traffic_ops``

.. Example output ::


..     trafficops-vm # yum install traffic_ops
..     Loaded plugins: fastestmirror, security
..     Loading mirror speeds from cached hostfile
..     Setting up Install Process
..     Resolving Dependencies
..     --> Running transaction check
..     ---> Package traffic_ops.x86_64 0:1.28-1505 will be installed
..     --> Processing Dependency: perl-Digest-SHA1 for package: traffic_ops-1.28-1505.x86_64
..     --> Processing Dependency: perl-DBI for package: traffic_ops-1.28-1505.x86_64
..     --> Processing Dependency: perl-DBD-MySQL for package: traffic_ops-1.28-1505.x86_64
..     --> Processing Dependency: mysql-server for package: traffic_ops-1.28-1505.x86_64
..     --> Processing Dependency: mysql for package: traffic_ops-1.28-1505.x86_64
..     --> Processing Dependency: mod_ssl for package: traffic_ops-1.28-1505.x86_64
..     --> Processing Dependency: mkisofs for package: traffic_ops-1.28-1505.x86_64
..     --> Processing Dependency: libpcap-devel for package: traffic_ops-1.28-1505.x86_64
..     --> Processing Dependency: expat-devel for package: traffic_ops-1.28-1505.x86_64
..     --> Running transaction check
..     ---> Package expat-devel.x86_64 0:2.0.1-11.el6_2 will be installed
..     ---> Package genisoimage.x86_64 0:1.1.9-12.el6 will be installed
..     ---> Package libpcap-devel.x86_64 14:1.4.0-1.20130826git2dbcaa1.el6 will be installed
..     ---> Package mod_ssl.x86_64 1:2.2.15-30.el6.centos will be installed
..     --> Processing Dependency: httpd-mmn = 20051115 for package: 1:mod_ssl-2.2.15-30.el6.centos.x86_64
..     --> Processing Dependency: httpd = 2.2.15-30.el6.centos for package: 1:mod_ssl-2.2.15-30.el6.centos.x86_64
..     --> Processing Dependency: httpd for package: 1:mod_ssl-2.2.15-30.el6.centos.x86_64
..     ---> Package mysql.x86_64 0:5.1.73-3.el6_5 will be installed
..     --> Processing Dependency: mysql-libs = 5.1.73-3.el6_5 for package: mysql-5.1.73-3.el6_5.x86_64
..     ---> Package mysql-server.x86_64 0:5.1.73-3.el6_5 will be installed
..     ---> Package perl-DBD-MySQL.x86_64 0:4.013-3.el6 will be installed
..     ---> Package perl-DBI.x86_64 0:1.609-4.el6 will be installed
..     ---> Package perl-Digest-SHA1.x86_64 0:2.12-2.el6 will be installed
..     --> Running transaction check
..     ---> Package httpd.x86_64 0:2.2.15-30.el6.centos will be installed
..     --> Processing Dependency: httpd-tools = 2.2.15-30.el6.centos for package: httpd-2.2.15-30.el6.centos.x86_64
..     --> Processing Dependency: apr-util-ldap for package: httpd-2.2.15-30.el6.centos.x86_64
..     --> Processing Dependency: libaprutil-1.so.0()(64bit) for package: httpd-2.2.15-30.el6.centos.x86_64
..     --> Processing Dependency: libapr-1.so.0()(64bit) for package: httpd-2.2.15-30.el6.centos.x86_64
..     ---> Package mysql-libs.x86_64 0:5.1.71-1.el6 will be updated
..     ---> Package mysql-libs.x86_64 0:5.1.73-3.el6_5 will be an update
..     --> Running transaction check
..     ---> Package apr.x86_64 0:1.3.9-5.el6_2 will be installed
..     ---> Package apr-util.x86_64 0:1.3.9-3.el6_0.1 will be installed
..     ---> Package apr-util-ldap.x86_64 0:1.3.9-3.el6_0.1 will be installed
..     ---> Package httpd-tools.x86_64 0:2.2.15-30.el6.centos will be installed
..     --> Finished Dependency Resolution

..     Dependencies Resolved

..     ====================================================================================================================================================================================
..      Package                                Arch                         Version                                                   Repository                                      Size
..     ====================================================================================================================================================================================
..     Installing:
..      traffic_ops                            x86_64                       1.28-1505                                                 local-copy-of-yum_NOARCH                        33 M
..     Installing for dependencies:
..      apr                                    x86_64                       1.3.9-5.el6_2                                             local-copy-of-yum_REPO                         123 k
..      apr-util                               x86_64                       1.3.9-3.el6_0.1                                           local-copy-of-yum_REPO                          87 k
..      apr-util-ldap                          x86_64                       1.3.9-3.el6_0.1                                           local-copy-of-yum_REPO                          15 k
..      expat-devel                            x86_64                       2.0.1-11.el6_2                                            local-copy-of-yum_REPO                         120 k
..      genisoimage                            x86_64                       1.1.9-12.el6                                              local-copy-of-yum_REPO                         348 k
..      httpd                                  x86_64                       2.2.15-30.el6.centos                                      local-copy-of-yum_REPO                         821 k
..      httpd-tools                            x86_64                       2.2.15-30.el6.centos                                      local-copy-of-yum_REPO                          73 k
..      libpcap-devel                          x86_64                       14:1.4.0-1.20130826git2dbcaa1.el6                         local-copy-of-yum_REPO                         114 k
..      mod_ssl                                x86_64                       1:2.2.15-30.el6.centos                                    local-copy-of-yum_REPO                          91 k
..      mysql                                  x86_64                       5.1.73-3.el6_5                                            local-copy-of-yum_REPO                         894 k
..      mysql-server                           x86_64                       5.1.73-3.el6_5                                            local-copy-of-yum_REPO                         8.6 M
..      perl-DBD-MySQL                         x86_64                       4.013-3.el6                                               local-copy-of-yum_REPO                         134 k
..      perl-DBI                               x86_64                       1.609-4.el6                                               local-copy-of-yum_REPO                         705 k
..      perl-Digest-SHA1                       x86_64                       2.12-2.el6                                                local-copy-of-yum_REPO                          49 k
..     Updating for dependencies:
..      mysql-libs                             x86_64                       5.1.73-3.el6_5                                            local-copy-of-yum_REPO                         1.2 M

..     Transaction Summary
..     ====================================================================================================================================================================================
..     Install      15 Package(s)
..     Upgrade       1 Package(s)

..     Total download size: 47 M
..     Is this ok [y/N]: y
..     Downloading Packages:
..     (1/16): apr-1.3.9-5.el6_2.x86_64.rpm                                                                                                                         | 123 kB     00:00
..     (2/16): apr-util-1.3.9-3.el6_0.1.x86_64.rpm                                                                                                                  |  87 kB     00:00
..     (3/16): apr-util-ldap-1.3.9-3.el6_0.1.x86_64.rpm                                                                                                             |  15 kB     00:00
..     (4/16): expat-devel-2.0.1-11.el6_2.x86_64.rpm                                                                                                                | 120 kB     00:00
..     (5/16): genisoimage-1.1.9-12.el6.x86_64.rpm                                                                                                                  | 348 kB     00:00
..     (6/16): httpd-2.2.15-30.el6.centos.x86_64.rpm                                                                                                                | 821 kB     00:00
..     (7/16): httpd-tools-2.2.15-30.el6.centos.x86_64.rpm                                                                                                          |  73 kB     00:00
..     (8/16): libpcap-devel-1.4.0-1.20130826git2dbcaa1.el6.x86_64.rpm                                                                                              | 114 kB     00:00
..     (9/16): mod_ssl-2.2.15-30.el6.centos.x86_64.rpm                                                                                                              |  91 kB     00:00
..     (10/16): mysql-5.1.73-3.el6_5.x86_64.rpm                                                                                                                     | 894 kB     00:00
..     (11/16): mysql-libs-5.1.73-3.el6_5.x86_64.rpm                                                                                                                | 1.2 MB     00:00
..     (12/16): mysql-server-5.1.73-3.el6_5.x86_64.rpm                                                                                                              | 8.6 MB     00:00
..     (13/16): perl-DBD-MySQL-4.013-3.el6.x86_64.rpm                                                                                                               | 134 kB     00:00
..     (14/16): perl-DBI-1.609-4.el6.x86_64.rpm                                                                                                                     | 705 kB     00:00
..     (15/16): perl-Digest-SHA1-2.12-2.el6.x86_64.rpm                                                                                                              |  49 kB     00:00
..     (16/16): traffic_ops-1.28-1505.x86_64.rpm                                                                                                                    |  33 MB     00:02
..     ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
..     Total                                                                                                                                                11 MB/s |  47 MB     00:04
..     Running rpm_check_debug
..     Running Transaction Test
..     Transaction Test Succeeded
..     Running Transaction
..       Installing : perl-DBI-1.609-4.el6.x86_64                                                                                                                                     1/17
..       Updating   : mysql-libs-5.1.73-3.el6_5.x86_64                                                                                                                                2/17
..       Installing : apr-1.3.9-5.el6_2.x86_64                                                                                                                                        3/17
..       Installing : apr-util-1.3.9-3.el6_0.1.x86_64                                                                                                                                 4/17
..       Installing : perl-DBD-MySQL-4.013-3.el6.x86_64                                                                                                                               5/17
..       Installing : mysql-5.1.73-3.el6_5.x86_64                                                                                                                                     6/17
..       Installing : mysql-server-5.1.73-3.el6_5.x86_64                                                                                                                              7/17
..       Installing : apr-util-ldap-1.3.9-3.el6_0.1.x86_64                                                                                                                            8/17
..       Installing : httpd-tools-2.2.15-30.el6.centos.x86_64                                                                                                                         9/17
..       Installing : httpd-2.2.15-30.el6.centos.x86_64                                                                                                                              10/17
..       Installing : 1:mod_ssl-2.2.15-30.el6.centos.x86_64                                                                                                                          11/17
..       Installing : 14:libpcap-devel-1.4.0-1.20130826git2dbcaa1.el6.x86_64                                                                                                         12/17
..       Installing : expat-devel-2.0.1-11.el6_2.x86_64                                                                                                                              13/17
..       Installing : genisoimage-1.1.9-12.el6.x86_64                                                                                                                                14/17
..       Installing : perl-Digest-SHA1-2.12-2.el6.x86_64                                                                                                                             15/17
..       Installing : traffic_ops-1.28-1505.x86_64                                                                                                                                   16/17

..     Run /opt/traffic_ops/install/bin/postinstall from the root home directory to complete the install.

..       Cleanup    : mysql-libs-5.1.71-1.el6.x86_64                                                                                                                                 17/17
..       Verifying  : 1:mod_ssl-2.2.15-30.el6.centos.x86_64                                                                                                                           1/17
..       Verifying  : apr-1.3.9-5.el6_2.x86_64                                                                                                                                        2/17
..       Verifying  : perl-DBD-MySQL-4.013-3.el6.x86_64                                                                                                                               3/17
..       Verifying  : mysql-libs-5.1.73-3.el6_5.x86_64                                                                                                                                4/17
..       Verifying  : mysql-server-5.1.73-3.el6_5.x86_64                                                                                                                              5/17
..       Verifying  : mysql-5.1.73-3.el6_5.x86_64                                                                                                                                     6/17
..       Verifying  : perl-Digest-SHA1-2.12-2.el6.x86_64                                                                                                                              7/17
..       Verifying  : apr-util-ldap-1.3.9-3.el6_0.1.x86_64                                                                                                                            8/17
..       Verifying  : perl-DBI-1.609-4.el6.x86_64                                                                                                                                     9/17
..       Verifying  : httpd-tools-2.2.15-30.el6.centos.x86_64                                                                                                                        10/17
..       Verifying  : genisoimage-1.1.9-12.el6.x86_64                                                                                                                                11/17
..       Verifying  : httpd-2.2.15-30.el6.centos.x86_64                                                                                                                              12/17
..       Verifying  : traffic_ops-1.28-1505.x86_64                                                                                                                                   13/17
..       Verifying  : expat-devel-2.0.1-11.el6_2.x86_64                                                                                                                              14/17
..       Verifying  : 14:libpcap-devel-1.4.0-1.20130826git2dbcaa1.el6.x86_64                                                                                                         15/17
..       Verifying  : apr-util-1.3.9-3.el6_0.1.x86_64                                                                                                                                16/17
..       Verifying  : mysql-libs-5.1.71-1.el6.x86_64                                                                                                                                 17/17

..     Installed:
..       traffic_ops.x86_64 0:1.28-1505

..     Dependency Installed:
..       apr.x86_64 0:1.3.9-5.el6_2             apr-util.x86_64 0:1.3.9-3.el6_0.1     apr-util-ldap.x86_64 0:1.3.9-3.el6_0.1     expat-devel.x86_64 0:2.0.1-11.el6_2
..       genisoimage.x86_64 0:1.1.9-12.el6      httpd.x86_64 0:2.2.15-30.el6.centos   httpd-tools.x86_64 0:2.2.15-30.el6.centos  libpcap-devel.x86_64 14:1.4.0-1.20130826git2dbcaa1.el6
..       mod_ssl.x86_64 1:2.2.15-30.el6.centos  mysql.x86_64 0:5.1.73-3.el6_5         mysql-server.x86_64 0:5.1.73-3.el6_5       perl-DBD-MySQL.x86_64 0:4.013-3.el6
..       perl-DBI.x86_64 0:1.609-4.el6          perl-Digest-SHA1.x86_64 0:2.12-2.el6

..     Dependency Updated:
..       mysql-libs.x86_64 0:5.1.73-3.el6_5

..     Complete!
..     trafficops-vm #

.. _rl-ps:

.. The postinstall script
.. ----------------------
2. After installation of Traffic Ops rpm enter the following command: ``/opt/traffic_ops/install/bin/postinstall``

  Example output::


      trafficops-vm # /opt/traffic_ops/install/bin/postinstall

      This script will build and package the required Traffic Ops perl modules.
      In order to complete this operation, Development tools such as the gcc
      compiler must be installed on this machine.

      Hit ENTER to continue:


  The first thing the post install will do is install additional packages needed from the yum repo.

  Ater that, it will automatically proceed to installing the required Perl packages from CPAN.

  .. Note:: Especially when installing Traffic Ops for the first time on a system this can take a long time, since many dependencies for the Mojolicous application need to be downloaded. Expect 30 minutes. 

  If there are any prompts in this phase, please just answer with the defaults (some CPAN installs can prompt for install questions). 

  When this phase is complete, you will see:: 

      ...
      Successfully installed Test-Differences-0.63
      Successfully installed DBIx-Class-Schema-Loader-0.07042
      Successfully installed Time-HiRes-1.9726 (upgraded from 1.9719)
      Successfully installed Mojolicious-Plugin-Authentication-1.26
      113 distributions installed
      Complete! Modules were installed into /opt/traffic_ops/app/local
      Linking perl libraries...
      Installing perl scripts


      This script will initialize the Traffic Ops database.
      Please enter the following information in order to completely
      configure the Traffic Ops mysql database.


      Database type [mysql]:


  The next phase of the install will ask you about the local environment for your CDN.

  .. Note:: before proceeding to this step, the database has to have at least a root password, and needs to be started. When using mysql, please type ``service mysqld start`` as root in another terminal and follow the instructions on the screen to set the root passwd.

  .. Note:: CentOS files note.

  Example output::

      Database type [mysql]:
      Database name [traffic_ops_db]:
      Database server hostname IP or FQDN [localhost]:
      Database port number [3306]:
      Traffic Ops database user [traffic_ops]:
      Password for traffic_ops:
      Re-Enter password for traffic_ops:

      Error: passwords do not match, try again.

      Password for traffic_ops:
      Re-Enter password for traffic_ops:

      Database server root (admin) user name [root]:
      Database server root password:
      Database Type: mysql
      Database Name: traffic_ops_db
      Hostname: localhost
      Port: 3306
      Database User: traffic_ops
      Is the above information correct (y/n) [n]:  y

      The database properties have been saved to /opt/traffic_ops/app/conf/production/database.conf

        The database configuration has been saved.  Now we need to set some custom
        fields that are necessary for the CDN to function correctly.


      Traffic Ops url [https://localhost]:  https://traffic-ops.kabletown.net
      Human-readable CDN Name.  (No whitespace, please) [kabletown_cdn]:
      DNS sub-domain for which your CDN is authoritative [cdn1.kabletown.net]:
      Fully qualified name of your CentOS 6.5 ISO kickstart tar file, or 'na' to skip and add files later [/var/cache/centos65.tgz]:  na
      Fully qualified location to store your ISO kickstart files [/var/www/files]:

      Traffic Ops URL: https://traffic-ops.kabletown.net
      Traffic Ops Info URL: https://traffic-ops.kabletown.net/info
      Domainname: cdn1.kabletown.net
      CDN Name: kabletown_cdn
      GeoLocation Polling URL: https://traffic-ops.kabletown.net/routing/GeoIP2-City.mmdb.gz
      CoverageZone Polling URL: https://traffic-ops.kabletown.net/routing/coverage-zone.json

      Is the above information correct (y/n) [n]:  y
      Parameter information has been saved to /opt/traffic_ops/install/data/json/parameters.json


      Adding an administration user to the Traffic Ops database.

      Administration username for Traffic Ops:  admin
      Password for the admin user admin:
      Verify the password for admin:
      Do you wish to create an ldap configuration for access to traffic ops [y/n] ? [n]:  n
      creating database
      Creating database...
      Creating user...
      Flushing privileges...
      setting up database
      Executing 'drop database traffic_ops_db'
      Executing 'create database traffic_ops_db'
      Creating database tables...
      Migrating database...
      goose: migrating db environment 'production', current version: 0, target: 20150316100000
      OK    20141222103718_extension.sql
      OK    20150108100000_add_job_deliveryservice.sql
      OK    20150205100000_cg_location.sql
      OK    20150209100000_cran_to_asn.sql
      OK    20150210100000_ds_keyinfo.sql
      OK    20150304100000_add_ip6_ds_routing.sql
      OK    20150310100000_add_bg_fetch.sql
      OK    20150316100000_move_hdr_rw.sql
      Seeding database...
      Database initialization succeeded.
      seeding profile data...
      name EDGE1 description Edge 1
      name TR1 description Traffic Router 1
      name TM1 description Traffic Monitor 1
      name MID1 description Mid 1
      seeding parameter data...

  Explanation of the information that needs to be provided:

    +----------------------------------------------------+-----------------------------------------------------------------------------------------------+
    |                       Field                        |                                          Description                                          |
    +====================================================+===============================================================================================+
    | Database type                                      | mysql or postgres                                                                             |
    +----------------------------------------------------+-----------------------------------------------------------------------------------------------+
    | Database name                                      | The name of the database Traffic Ops uses to store the configuration information              |
    +----------------------------------------------------+-----------------------------------------------------------------------------------------------+
    | Database server hostname IP or FQDN                | The hostname of the database server                                                           |
    +----------------------------------------------------+-----------------------------------------------------------------------------------------------+
    | Database port number                               | The database port number                                                                      |
    +----------------------------------------------------+-----------------------------------------------------------------------------------------------+
    | Traffic Ops database user                          | The username Traffic Ops will use to read/write from the database                             |
    +----------------------------------------------------+-----------------------------------------------------------------------------------------------+
    | password for traffic ops                           | The passwdord for the above database user                                                     |
    +----------------------------------------------------+-----------------------------------------------------------------------------------------------+
    | Database server root (admin) user name             | Priviledged database user that has permission to create the database and user for Traffic Ops |
    +----------------------------------------------------+-----------------------------------------------------------------------------------------------+
    | Database server root (admin) user password         | The password for the above priviledged database user                                          |
    +----------------------------------------------------+-----------------------------------------------------------------------------------------------+
    | Traffic Ops url                                    | The URL to connect to this instance of Traffic Ops, usually https://<traffic ops host FQDN>/  |
    +----------------------------------------------------+-----------------------------------------------------------------------------------------------+
    | Human-readable CDN Name                            | The name of the first CDN traffic Ops will be managing                                        |
    +----------------------------------------------------+-----------------------------------------------------------------------------------------------+
    | DNS sub-domain for which your CDN is authoritative | The DNS domain that will be delegated to this Traffic Control CDN                             |
    +----------------------------------------------------+-----------------------------------------------------------------------------------------------+
    | name of your CentOS 6.5 ISO kickstart tar file     | See :ref:`Creating-CentOS-Kickstart`                                                          |
    +----------------------------------------------------+-----------------------------------------------------------------------------------------------+
    | Administration username for Traffic Ops            | The Administration (highest privilege) Traffic Ops user to create;                            |
    |                                                    | use this user to login for the first time and create other users                              |
    +----------------------------------------------------+-----------------------------------------------------------------------------------------------+
    | Password for the admin user                        | The passwd for the above user                                                                 |
    +----------------------------------------------------+-----------------------------------------------------------------------------------------------+


  The postinstall script will now seed the database with some inital configuration settings for the CDN and the servers in the CDN.

  The next phase is the download of the geo location database and configuration of information needed for SSL certificates.

  Example output::

     JvD to provide new screen scrape. 


Traffic Ops is now installed!

Upgrading Traffic Ops
=====================
To upgrade:

1. Enter the following command:``service traffic_ops stop``
2. Enter the following command:``yum upgrade traffic_ops``
3. See :ref:`rl-ps` to run the post install.
4. Enter the following command:``service traffic_ops start``
