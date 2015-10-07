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
  
.. _rl-ps:

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

1. Install Traffic Ops: ``sudo yum install traffic_ops``





2. After installation of Traffic Ops rpm enter the following command: ``sudo /opt/traffic_ops/install/bin/postinstall``

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

    Downloading MaxMind data.
    --2015-04-14 02:14:32--  http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz
    Resolving geolite.maxmind.com... 141.101.115.190, 141.101.114.190, 2400:cb00:2048:1::8d65:73be, ...
    Connecting to geolite.maxmind.com|141.101.115.190|:80... connected.
    HTTP request sent, awaiting response... 200 OK
    Length: 17633433 (17M) [application/octet-stream]
    Saving to: “GeoLite2-City.mmdb.gz”

    100%[==================================================================================================================================================================>] 17,633,433  7.03M/s   in 2.4s

    2015-04-14 02:14:35 (7.03 MB/s) - “GeoLite2-City.mmdb.gz” saved [17633433/17633433]

    Copying coverage zone file to public dir.

    Installing SSL Certificates.

      We're now running a script to generate a self signed X509 SSL certificate.
      When prompted to enter a pass phrase, just enter 'pass' each time.  The
      pass phrase will be stripped from the private key before installation.

      When prompted to enter a 'challenge password', just hit the ENTER key.

      The remaining enformation Country, State, Locality, etc... are required to
      generate a properly formatted SSL certificate.

    Hit Enter when you are ready to continue:
    Postinstall SSL Certificate Creation.

    Generating an RSA Private Server Key.

    Generating RSA private key, 1024 bit long modulus
    ..........................++++++
    .....................++++++
    e is 65537 (0x10001)
    Enter pass phrase for server.key:
    Verifying - Enter pass phrase for server.key:

    The server key has been generated.

    Creating a Certificate Signing Request (CSR)

    Enter pass phrase for server.key:
    You are about to be asked to enter information that will be incorporated
    into your certificate request.
    What you are about to enter is what is called a Distinguished Name or a DN.
    There are quite a few fields but you can leave some blank
    For some fields there will be a default value,
    If you enter '.', the field will be left blank.
    -----
    Country Name (2 letter code) [XX]:US
    State or Province Name (full name) []:CO
    Locality Name (eg, city) [Default City]:Denver
    Organization Name (eg, company) [Default Company Ltd]:
    Organizational Unit Name (eg, section) []:
    Common Name (eg, your name or your server's hostname) []:
    Email Address []:

    Please enter the following 'extra' attributes
    to be sent with your certificate request
    A challenge password []:pass
    An optional company name []:

    The Certificate Signing Request has been generated.
    Removing the pass phrase from the server key.
    Enter pass phrase for server.key.orig:
    writing RSA key

    The pass phrase has been removed from the server key.

    Generating a Self-signed certificate.
    Signature ok
    subject=/C=US/ST=CO/L=Denver/O=Default Company Ltd
    Getting Private key

    A server key and self signed certificate has been generated.

    Installing the server key and server certificate.

    The private key has been installed.

    Installing the self signed certificate.

    Saving the self signed csr.

      The self signed certificate has now been installed.

      You may obtain a certificate signed by a Certificate Authority using the
      server.csr file saved in the current directory.  Once you have obtained
      a signed certificate, copy it to /etc/pki/tls/certs/localhost.crt and
      restart Traffic Ops.



    SSL Certificates have been installed.

    Starting Traffic Ops.

    Starting Traffic Ops

    Subroutine TrafficOps::has redefined at /opt/traffic_ops/app/local/lib/perl5/Mojo/Base.pm line 38.
    Subroutine TrafficOps::has redefined at /opt/traffic_ops/app/local/lib/perl5/Mojo/Base.pm line 38.
    Loading config from /opt/traffic_ops/app/conf/cdn.conf
    Reading log4perl config from /opt/traffic_ops/app/conf/production/log4perl.conf
    Starting hot deployment for Hypnotoad server 32192.

    Waiting for Traffic Ops to start.


    Shutdown Traffic Ops [y/n] [n]:  n

    To start Traffic Ops:  service traffic_ops start
    To stop Traffic Ops:   service traffic_ops stop

    traffic_ops #

Traffic Ops is now installed!

3. Download the web dependencies (this will be added to the installer in the future): ::

    traffic_ops # pwd
    /opt/traffic_ops/install/bin
    traffic_ops # ./download_web_deps
    Finished curling https://cdn.datatables.net/1.10.4/js/jquery.dataTables.min.js | size is: 78746
    Finished curling https://github.com/fancyapps/fancyBox/zipball/v2.1.5 | size is: 541026
    Finished curling http://www.flotcharts.org/downloads/flot-0.8.3.zip | size is: 649913
    Finished curling https://github.com/krzysu/flot.tooltip/releases/download/0.8.4/jquery.flot.tooltip-0.8.4.zip | size is: 7669
    Finished curling https://gflot.googlecode.com/svn-history/r154/trunk/flot/jquery.flot.axislabels.js | size is: 17321
    Finished curling https://github.com/alpixel/jMenu/archive/master.zip | size is: 41836
    Finished curling https://code.jquery.com/jquery-1.11.2.min.js | size is: 95931
    Finished curling https://code.jquery.com/ui/1.11.4/jquery-ui.min.js | size is: 240427
    Finished curling https://code.jquery.com/ui/1.7.3/themes/dark-hive/jquery-ui.css | size is: 27499
    Finished curling http://jquery-ui.googlecode.com/svn/tags/1.7.3/themes/dark-hive/images/ui-bg_flat_30_cccccc_40x100.png | size is: 180
    Finished curling http://jquery-ui.googlecode.com/svn/tags/1.7.3/themes/dark-hive/images/ui-bg_flat_50_5c5c5c_40x100.png | size is: 180
    Finished curling http://jquery-ui.googlecode.com/svn/tags/1.7.3/themes/dark-hive/images/ui-bg_glass_40_ffc73d_1x400.png | size is: 131
    Finished curling http://jquery-ui.googlecode.com/svn/tags/1.7.3/themes/dark-hive/images/ui-bg_highlight-hard_20_0972a5_1x100.png | size is: 114
    Finished curling http://jquery-ui.googlecode.com/svn/tags/1.7.3/themes/dark-hive/images/ui-bg_highlight-soft_33_003147_1x100.png | size is: 127
    Finished curling http://jquery-ui.googlecode.com/svn/tags/1.7.3/themes/dark-hive/images/ui-bg_highlight-soft_35_222222_1x100.png | size is: 113
    Finished curling http://jquery-ui.googlecode.com/svn/tags/1.7.3/themes/dark-hive/images/ui-bg_highlight-soft_44_444444_1x100.png | size is: 117
    Finished curling http://jquery-ui.googlecode.com/svn/tags/1.7.3/themes/dark-hive/images/ui-bg_highlight-soft_80_eeeeee_1x100.png | size is: 95
    Finished curling http://jquery-ui.googlecode.com/svn/tags/1.7.3/themes/dark-hive/images/ui-bg_loop_25_000000_21x21.png | size is: 235
    Finished curling http://jquery-ui.googlecode.com/svn/tags/1.7.3/themes/dark-hive/images/ui-icons_222222_256x240.png | size is: 4369
    Finished curling http://jquery-ui.googlecode.com/svn/tags/1.7.3/themes/dark-hive/images/ui-icons_4b8e0b_256x240.png | size is: 4369
    Finished curling http://jquery-ui.googlecode.com/svn/tags/1.7.3/themes/dark-hive/images/ui-icons_a83300_256x240.png | size is: 4369
    Finished curling http://jquery-ui.googlecode.com/svn/tags/1.7.3/themes/dark-hive/images/ui-icons_cccccc_256x240.png | size is: 4369
    Finished curling http://jquery-ui.googlecode.com/svn/tags/1.7.3/themes/dark-hive/images/ui-icons_ffffff_256x240.png | size is: 4369
    Finished curling https://maxcdn.bootstrapcdn.com/bootstrap/3.3.4/js/bootstrap.min.js | size is: 35951
    Output file: ../../app/public/js/jquery.dataTables.min.js does not exist, putting into place.
    Making dir: ../../app/public/js/fancybox/
    Output file: ../../app/public/js/fancybox//jquery.fancybox-buttons.js does not exist. Putting file from zip into place.
    Output file: ../../app/public/js/fancybox//fancybox_loading@2x.gif does not exist. Putting file from zip into place.
    Output file: ../../app/public/js/fancybox//fancybox_loading.gif does not exist. Putting file from zip into place.
    Output file: ../../app/public/js/fancybox//fancybox_buttons.png does not exist. Putting file from zip into place.
    Output file: ../../app/public/js/fancybox//jquery.fancybox-thumbs.js does not exist. Putting file from zip into place.
    Output file: ../../app/public/js/fancybox//jquery.fancybox-buttons.css does not exist. Putting file from zip into place.
    Output file: ../../app/public/js/fancybox//jquery.fancybox-thumbs.css does not exist. Putting file from zip into place.
    Output file: ../../app/public/js/fancybox//fancybox_sprite@2x.png does not exist. Putting file from zip into place.
    Output file: ../../app/public/js/fancybox//jquery.fancybox.css does not exist. Putting file from zip into place.
    Output file: ../../app/public/js/fancybox//jquery.fancybox-media.js does not exist. Putting file from zip into place.
    Output file: ../../app/public/js/fancybox//fancybox_overlay.png does not exist. Putting file from zip into place.
    Output file: ../../app/public/js/fancybox//fancybox_sprite.png does not exist. Putting file from zip into place.
    Output file: ../../app/public/js/fancybox//jquery.fancybox.js does not exist. Putting file from zip into place.
    Making dir: ../../app/public/js/flot/
    Output file: ../../app/public/js/flot//jquery.flot.min.js does not exist. Putting file from zip into place.
    Output file: ../../app/public/js/flot//jquery.flot.selection.js does not exist. Putting file from zip into place.
    Output file: ../../app/public/js/flot//jquery.flot.stack.js does not exist. Putting file from zip into place.
    Output file: ../../app/public/js/flot//jquery.flot.time.js does not exist. Putting file from zip into place.
    Output file: ../../app/public/js/flot//jquery.flot.tooltip.js does not exist. Putting file from zip into place.
    Output file: ../../app/public/js/flot/jquery.flot.axislabels.js does not exist, putting into place.
    Output file: ../../app/public/js//jMenu.jquery.min.js does not exist. Putting file from zip into place.
    Output file: ../../app/public/css//jmenu.css does not exist. Putting file from zip into place.
    Output file: ../../app/public/js/jquery-1.11.2.min.js does not exist, putting into place.
    Output file: ../../app/public/js/jquery-ui.min.js does not exist, putting into place.
    Output file: ../../app/public/css/jquery-ui.css does not exist, putting into place.
    Making dir: ../../app/public/css/images/
    Output file: ../../app/public/css/images/ui-bg_flat_30_cccccc_40x100.png does not exist, putting into place.
    Output file: ../../app/public/css/images/ui-bg_flat_50_5c5c5c_40x100.png does not exist, putting into place.
    Output file: ../../app/public/css/images/ui-bg_glass_40_ffc73d_1x400.png does not exist, putting into place.
    Output file: ../../app/public/css/images/ui-bg_highlight-hard_20_0972a5_1x100.png does not exist, putting into place.
    Output file: ../../app/public/css/images/ui-bg_highlight-soft_33_003147_1x100.png does not exist, putting into place.
    Output file: ../../app/public/css/images/ui-bg_highlight-soft_35_222222_1x100.png does not exist, putting into place.
    Output file: ../../app/public/css/images/ui-bg_highlight-soft_44_444444_1x100.png does not exist, putting into place.
    Output file: ../../app/public/css/images/ui-bg_highlight-soft_80_eeeeee_1x100.png does not exist, putting into place.
    Output file: ../../app/public/css/images/ui-bg_loop_25_000000_21x21.png does not exist, putting into place.
    Output file: ../../app/public/css/images/ui-icons_222222_256x240.png does not exist, putting into place.
    Output file: ../../app/public/css/images/ui-icons_4b8e0b_256x240.png does not exist, putting into place.
    Output file: ../../app/public/css/images/ui-icons_a83300_256x240.png does not exist, putting into place.
    Output file: ../../app/public/css/images/ui-icons_cccccc_256x240.png does not exist, putting into place.
    Output file: ../../app/public/css/images/ui-icons_ffffff_256x240.png does not exist, putting into place.
    Output file: ../../app/public/js/bootstrap.min.js does not exist, putting into place.
    traffic_ops #

Upgrading Traffic Ops
=====================
To upgrade:

1. Enter the following command:``service traffic_ops stop``
2. Enter the following command:``yum upgrade traffic_ops``
3. See :ref:`rl-ps` to run the post install.
4. Enter the following command:``service traffic_ops start``
