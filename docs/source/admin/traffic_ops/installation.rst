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

.. _to-install:

************************
Traffic Ops - Installing
************************

System Requirements
-------------------
The user must have the following for a successful minimal install:

- CentOS 7+
- 2 machines or Virtual Machines (VMs), each with at least 2 (v)CPUs, 4GB of RAM, and 20 GB of disk space.
- Access to CentOS Base and EPEL repositories
- Access to `The Comprehensive Perl Archive Network (CPAN) <http://www.cpan.org/>`_

As of version 2.0 only PostgreSQL is supported as the database. This documentation assumes CentOS 7.2 and PostgreSQL 9.6.3 for a production install.

.. highlight:: none

Installation
------------

#. Install PostgreSQL Database

	.. note:: For more information on installing PostgreSQL, see `their documentation <https://www.postgresql.org/docs/>`_.

	For a production install it is best to install PostgreSQL on its own server/VM. To install PostgreSQL, on the PostgreSQL host (hostname ``pg`` in example),
	run the following commands as the root user (or with ``sudo``):

	.. code-block:: shell

		yum update -y
		yum install -y https://download.postgresql.org/pub/repos/yum/9.6/redhat/rhel-7-x86_64/pgdg-centos96-9.6-3.noarch.rpm
		yum install -y postgresql96-server
		su - postgres -c /usr/pgsql-9.6/bin/initdb -A md5 -W #-W forces the user to provide a superuser (postgres) password


	Edit ``/var/lib/pgsql/9.6/data/pg_hba.conf`` to allow your Traffic Ops instance to access the PostgreSQL server. For example if you are going to install Traffic Ops on ``99.33.99.1`` add::

		host  all   all     99.33.99.1/32 md5

	to the appropriate section of this file. Edit the ``/var/lib/pgsql/9.6/data/postgresql.conf`` file to add the appropriate listen_addresses or ``listen_addresses = '*'``, set ``timezone = 'UTC'``, and start the database:

	.. code-block:: shell

		systemctl enable postgresql-9.6
		systemctl start postgresql-9.6
		systemctl status postgresql-9.6


#. Build a Traffic Ops ``.rpm`` file using the instructions under the :ref:`dev-building` page.

#. Install PostgreSQL. To install the PostgreSQL 9.6 yum repository, run this command as the root user (or with ``sudo``):

	.. code-block:: shell

		yum install -y https://download.postgresql.org/pub/repos/yum/9.6/redhat/rhel-7-x86_64/pgdg-centos96-9.6-3.noarch.rpm

#. Install the Traffic Ops RPM. The Traffic Ops RPM file should have been built in an earlier step. To install it, simply run the following command as the root user (or with ``sudo``):

	.. code-block:: shell

		yum install -y ./dist/traffic_ops-2.0.0-xxxx.yyyyyyy.el7.x86_64.rpm

	.. note:: This will install the PostgreSQL client, ``psql`` as a dependency.

#. Install Additional Packages. Some packages on which Traffic Ops depends not have been installed as direct dependencies of the ``traffic_ops-<version stuff>.rpm``. To explicitly install these, run the following commands as the root user (or with ``sudo``):

	.. code-block:: shell

		yum install -y git
		wget -q https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz -O go.tar.gz
		tar -C /usr/local -xzf go.tar.gz
		PATH=$PATH:/usr/local/go/bin                    # go binaries are needed in the path for the 'postinstall' script
		go get bitbucket.org/liamstask/goose/cmd/goose

	.. note:: These are for the Traffic Control version 2.0.0 install, this may change, but the explicit installs won't hurt.

#. Login to the Database from the Traffic Ops machine. At this point you should be able to login from the Traffic Ops (hostname ``to`` in the example) host to the PostgreSQL (hostname ``pg`` in the example) host like so:

	.. code-block:: psql

		to-# psql -h 99.33.99.1 -U postgres
		Password for user postgres:
		psql (9.6.3)
		Type "help" for help.

		postgres=#


#. Create the User and Database. In this example, we use user: ``traffic_ops``, password: ``tcr0cks``, database: ``traffic_ops``:

	.. code-block:: psql

		to-# psql -U postgres -h 99.33.99.1 -c "CREATE USER traffic_ops WITH ENCRYPTED PASSWORD 'tcr0cks';"
		Password for user postgres:
		CREATE ROLE
		to-# createdb traffic_ops --owner traffic_ops -U postgres -h 99.33.99.1
		Password:
		to-#

#. Run the ``postinstall`` Script. Now, run the following command as the root user (or with ``sudo``): ``/opt/traffic_ops/install/bin/postinstall``. The ``postinstall`` script will first get all packages needed from CPAN. This may take a while, expect up to 30 minutes on the first install. If there are any prompts in this phase, please just answer with the defaults (some CPAN installs can prompt for install questions). When this phase is complete, you will see ``Complete! Modules were installed into /opt/traffic_ops/app/local``. Some additional files will be installed, and then it will proceed with the next phase of the install, where it will ask you about the local environment for your CDN. Please make sure you remember all your answers and verify that the database answers match the information previously used to create the database. Example output:

	.. code-block:: none

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



	.. table:: Explanation of the information that needs to be provided:

		+----------------------------------------------------+----------------------------------------------------------------------------------------------+
		| Field                                              | Description                                                                                  |
		+====================================================+==============================================================================================+
		| Database type                                      | This requests the type of database to be used. Answer with the default - 'Pg' to indicate a  |
		|                                                    | PostgreSQL database.                                                                         |
		+----------------------------------------------------+----------------------------------------------------------------------------------------------+
		| Database name                                      | The name of the database Traffic Ops uses to store the configuration information.            |
		+----------------------------------------------------+----------------------------------------------------------------------------------------------+
		| Database server hostname IP or FQDN                | The hostname of the database server (``pg`` in the example).                                 |
		+----------------------------------------------------+----------------------------------------------------------------------------------------------+
		| Database port number                               | The database port number. The default value, 5432, should be correct unless you changed it   |
		|                                                    | during the setup.                                                                            |
		+----------------------------------------------------+----------------------------------------------------------------------------------------------+
		| Traffic Ops database user                          | The username Traffic Ops will use to read/write from the database.                           |
		+----------------------------------------------------+----------------------------------------------------------------------------------------------+
		| Password for Traffic Ops                           | The password for the database user that Traffic Ops uses.                                    |
		+----------------------------------------------------+----------------------------------------------------------------------------------------------+
		| Database server root (admin) user name             | Privileged database user that has permission to create the database and user for Traffic Ops.|
		+----------------------------------------------------+----------------------------------------------------------------------------------------------+
		| Database server root (admin) user password         | The password for the privileged database user.                                               |
		+----------------------------------------------------+----------------------------------------------------------------------------------------------+
		| Traffic Ops URL                                    | The URL to connect to this instance of Traffic Ops, usually https://<Traffic Ops host FQDN>/ |
		+----------------------------------------------------+----------------------------------------------------------------------------------------------+
		| Human-readable CDN Name                            | The name of the first CDN which Traffic Ops will be manage.                                  |
		+----------------------------------------------------+----------------------------------------------------------------------------------------------+
		| DNS sub-domain for which your CDN is authoritative | The DNS domain that will be delegated to this Traffic Control CDN.                           |
		+----------------------------------------------------+----------------------------------------------------------------------------------------------+
		| Administration username for Traffic Ops            | The Administration (highest privilege) Traffic Ops user to create. Use this user to login    |
		|                                                    | for the first time and create other users.                                                   |
		+----------------------------------------------------+----------------------------------------------------------------------------------------------+
		| Password for the admin user                        | The password for the administrative Traffic Ops user.                                        |
		+----------------------------------------------------+----------------------------------------------------------------------------------------------+


Traffic Ops is now installed!


**To complete the Traffic Ops Setup See:** :ref:`default-profiles`


Upgrading Traffic Ops
=====================
To upgrade from older Traffic Ops versions, run the following commands as the root user (or with ``sudo``):

	.. code-block:: shell

		systemctl stop traffic_ops
		yum upgrade traffic_ops
		pushd /opt/traffic_ops/app/
		PERL5LIB=/opt/traffic_ops/app/lib:/opt/traffic_ops/app/local/lib/perl5 ./db/admin.pl --env production upgrade

After this completes, see :ref:`to-install` to run the ``postinstall`` script.
Once the ``postinstall`` script, has finished, run the following command as the root user (or with ``sudo``):
``systemctl start traffic_ops``
