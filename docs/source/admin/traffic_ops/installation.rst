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
  Traffic Ops - Installing

.. _rl-to-install:

Traffic Ops - Installing
%%%%%%%%%%%%%%%%%%%%%%%%

System Requirements
-------------------
The user must have the following for a successful minimal install:

* CentOS 7
* 2 VMs with at least 2 vCPUs, 4GB RAM, 20 GB disk space each
* Access to Centos Base and epel repositories
* Access to `The Comprehensive Perl Archive Network (CPAN) <http://www.cpan.org/>`_

As of version 2.0 only Postgres is supported as the database. This documentation assumes CentOS 7.2 and Postgresql 9.6.3 for a production install.

.. highlight:: none

Navigating the Install
-----------------------
To begin the install:

1. Install Postgres

  For a production install it is best to install postgres on it's own server/VM. To install postgres, on the postgres host (pg) ::

    pg-$ sudo su -
    pg-# yum -y update
    pg-# yum -y install https://download.postgresql.org/pub/repos/yum/9.6/redhat/rhel-7-x86_64/pgdg-centos96-9.6-3.noarch.rpm
    pg-# yum -y install postgresql96-server
    pg-$ su - postgres
    pg-$ /usr/pgsql-9.6/bin/initdb -A md5 -W #-W forces the user to provide a superuser (postgres) password


  Edit ``/var/lib/pgsql/9.6/data/pg_hba.conf`` to allow your traffic ops app server access. For example if you are going to install traffic ops on ``99.33.99.1`` add::

    host  all   all     99.33.99.1/32 md5

  to the appropriate section of this file. Edit the ``/var/lib/pgsql/9.6/data/postgresql.conf`` file to add the approriate listen_addresses or ``listen_addresses = '*'``, set ``timezone = 'UTC'``,  and start the database: ::

    pg-$ exit
    pg-# systemctl enable postgresql-9.6
    pg-# systemctl start postgresql-9.6
    pg-# systemctl status postgresql-9.6


2. Build Traffic Ops

   Build a Traffic Ops rpm using the instructions under the :ref:`dev-building` page.


3. Install Postgresql

  Install the postgresql 9.6 yum repository access. ::

    to-$ sudo su -
    to-# yum -y install https://download.postgresql.org/pub/repos/yum/9.6/redhat/rhel-7-x86_64/pgdg-centos96-9.6-3.noarch.rpm

4. Install the rpm built in step 2. ::

    to-# yum -y install ./dist/traffic_ops-2.0.0-xxxx.yyyyyyy.el7.x86_64.rpm


  Install some additional packages that it depends on that were not installed as dependecies in the previous step (these are for the 2.0.0 install, this may change, but the pre-installs won't hurt): ::

    to-# yum -y install git
    to-# wget -q https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz
    to-# tar -C /usr/local -xzf go1.8.3.linux-amd64.tar.gz
    to-# PATH=$PATH:/usr/local/go/bin             # go bins are needed in the path for postinstall
    to-# go get bitbucket.org/liamstask/goose/cmd/goose

  At this point you should be able to login to the database from the ``to`` host to the ``pg`` host like: ::

    to-# psql -h 99.33.99.1 -U postgres
    Password for user postgres:
    psql (9.6.3)
    Type "help" for help.

    postgres=#

  Use this connectivity to create the user and database. In  this example, we use user: ``traffic_ops``, password: ``tcr0cks``, database: ``traffic_ops``: ::

    to-# psql -U postgres -h 99.33.99.1 -c "CREATE USER traffic_ops  WITH ENCRYPTED PASSWORD 'tcr0cks';"
    Password for user postgres:
    CREATE ROLE
    to-# createdb traffic_ops --owner traffic_ops -U postgres -h 99.33.99.1
    Password:
    to-#


  Now, run the following command as root: ``/opt/traffic_ops/install/bin/postinstall``

  The postinstall will first get all packages needed from CPAN. This may take a while, expect up to 30 minutes on the first install.
  If there are any prompts in this phase, please just answer with the defaults (some CPAN installs can prompt for install questions).

  When this phase is complete, you will see::

      Complete! Modules were installed into /opt/traffic_ops/app/local

  Some additional files will be installed, and then it will proceed with the next phase of the install, where it will ask you about the local environment for your CDN. Please make sure you remember all your answers and the database answers match the database information previously used to create the database.


  Example output::

      ===========/opt/traffic_ops/app/conf/production/database.conf===========
      Database type [Pg]:
      Database type: Pg
      Database name [traffic_ops]:
      Database name: traffic_ops
      Database server hostname IP or FQDN [localhost]: 99.33.99.1
      Database server hostname IP or FQDN: 99.33.99.1
      Database port number [5432]:
      Database port number: 5432
      Traffic Ops database user [traffic_ops]:
      Traffic Ops database user: traffic_ops
      Password for Traffic Ops database user:
      Re-Enter Password for Traffic Ops database user:
      Writing json to /opt/traffic_ops/app/conf/production/database.conf
      Database configuration has been saved
      ===========/opt/traffic_ops/app/db/dbconf.yml===========
      Database server root (admin) user [postgres]:
      Database server root (admin) user: postgres
      Password for database server admin:
      Re-Enter Password for database server admin:
      Download Maxmind Database? [yes]:
      Download Maxmind Database?: yes
      ===========/opt/traffic_ops/app/conf/cdn.conf===========
      Generate a new secret? [yes]:
      Generate a new secret?: yes
      Number of secrets to keep? [10]:
      Number of secrets to keep?: 10
      Not setting up ldap
      ===========/opt/traffic_ops/install/data/json/users.json===========
      Administration username for Traffic Ops [admin]:
      Administration username for Traffic Ops: admin
      Password for the admin user:
      Re-Enter Password for the admin user:
      Writing json to /opt/traffic_ops/install/data/json/users.json
      ===========/opt/traffic_ops/install/data/json/openssl_configuration.json===========
      Do you want to generate a certificate? [yes]:
      Country Name (2 letter code): US
      State or Province Name (full name): CO
      Locality Name (eg, city): Denver
      Organization Name (eg, company): Super CDN, Inc
      Organizational Unit Name (eg, section):
      Common Name (eg, your name or your server's hostname):
      RSA Passphrase:
      Re-Enter RSA Passphrase:
      ===========/opt/traffic_ops/install/data/json/profiles.json===========
      Traffic Ops url [https://localhost]:
      Traffic Ops url: https://localhost
      Human-readable CDN Name.  (No whitespace, please) [kabletown_cdn]: blue_cdn
      Human-readable CDN Name.  (No whitespace, please): blue_cdn
      DNS sub-domain for which your CDN is authoritative [cdn1.kabletown.net]: blue-cdn.supercdn.net
      DNS sub-domain for which your CDN is authoritative: blue-cdn.supercdn.net
      Writing json to /opt/traffic_ops/install/data/json/profiles.json
      Downloading Maxmind data
      --2017-06-11 15:32:41--  http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz
      Resolving geolite.maxmind.com (geolite.maxmind.com)... 2400:cb00:2048:1::6810:262f, 2400:cb00:2048:1::6810:252f, 104.16.38.47, ...
      Connecting to geolite.maxmind.com (geolite.maxmind.com)|2400:cb00:2048:1::6810:262f|:80... connected.

      ... much SQL output skipped

      Starting Traffic Ops
      Restarting traffic_ops (via systemctl):                    [  OK  ]
      Waiting for Traffic Ops to restart
      Success! Postinstall complete.

      to-# ifconfig


  Explanation of the information that needs to be provided:

    +----------------------------------------------------+----------------------------------------------------------------------------------------------+
    | Field                                              | Description                                                                                  |
    +====================================================+==============================================================================================+
    | Database type                                      | Pg                                                                                           |
    +----------------------------------------------------+----------------------------------------------------------------------------------------------+
    | Database name                                      | The name of the database Traffic Ops uses to store the configuration information             |
    +----------------------------------------------------+----------------------------------------------------------------------------------------------+
    | Database server hostname IP or FQDN                | The hostname of the database server                                                          |
    +----------------------------------------------------+----------------------------------------------------------------------------------------------+
    | Database port number                               | The database port number                                                                     |
    +----------------------------------------------------+----------------------------------------------------------------------------------------------+
    | Traffic Ops database user                          | The username Traffic Ops will use to read/write from the database                            |
    +----------------------------------------------------+----------------------------------------------------------------------------------------------+
    | Password for traffic ops                           | The password for the above database user                                                     |
    +----------------------------------------------------+----------------------------------------------------------------------------------------------+
    | Database server root (admin) user name             | Privileged database user that has permission to create the database and user for Traffic Ops |
    +----------------------------------------------------+----------------------------------------------------------------------------------------------+
    | Database server root (admin) user password         | The password for the above privileged database user                                          |
    +----------------------------------------------------+----------------------------------------------------------------------------------------------+
    | Traffic Ops url                                    | The URL to connect to this instance of Traffic Ops, usually https://<traffic ops host FQDN>/ |
    +----------------------------------------------------+----------------------------------------------------------------------------------------------+
    | Human-readable CDN Name                            | The name of the first CDN traffic Ops will be managing                                       |
    +----------------------------------------------------+----------------------------------------------------------------------------------------------+
    | DNS sub-domain for which your CDN is authoritative | The DNS domain that will be delegated to this Traffic Control CDN                            |
    +----------------------------------------------------+----------------------------------------------------------------------------------------------+
    | Administration username for Traffic Ops            | The Administration (highest privilege) Traffic Ops user to create;                           |
    |                                                    | use this user to login for the first time and create other users                             |
    +----------------------------------------------------+----------------------------------------------------------------------------------------------+
    | Password for the admin user                        | The password for the above user                                                              |
    +----------------------------------------------------+----------------------------------------------------------------------------------------------+


Traffic Ops is now installed!


**To complete the Traffic Ops Setup See:** :ref:`rl-to-default-profiles`


Upgrading Traffic Ops
=====================
To upgrade:

#. Enter the following command:``service traffic_ops stop``
#. Enter the following command:``yum upgrade traffic_ops``
#. Enter the following command from the /opt/traffic_ops/app directory:
   ``PERL5LIB=/opt/traffic_ops/app/lib:/opt/traffic_ops/app/local/lib/perl5 ./db/admin.pl --env production upgrade``
#. See :ref:`rl-to-install` to run postinstall.
#. Enter the following command:``service traffic_ops start``




