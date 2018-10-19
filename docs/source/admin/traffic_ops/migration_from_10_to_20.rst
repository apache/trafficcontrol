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
	Traffic Ops - Migrating from Traffic Ops 1.x to Traffic Ops 2.x

.. _ps:

***************************************
Traffic Ops - Migrating from 1.x to 2.x
***************************************
In Traffic Ops 2.x the database used to store CDN information was changed from MySQL to PostgreSQL. PostgreSQL will remain the Traffic Ops database for the foreseeable future.
A Docker-based migration tool was developed to help with the conversion process using an open-source PostgreSQL tool called `pgloader <http://pgloader.io/>`_.
The following instructions will help configuring the Migration tool

System Requirements
-------------------
The user must have the following for a successful minimal install:

* CentOS 7.2+
* Docker installed [1]_
* PostgreSQL installed according to :ref:`to-install`

Setup the ``traffic_ops_db`` Directory
--------------------------------------
#. Modify the permissions of the ``/opt`` directory to make it writable by and owned by the ``postgres`` user and the ``postgres`` group. This can easily be accomplished by running the command ``chmod 755 /opt`` as the root user, or with ``sudo``.

#. Download the Traffic Control 2.0.0 tarball like so

	.. code-block:: shell

		cd /opt
		wget https://dist.apache.org/repos/dist/release/incubator/trafficcontrol/<tarball_version>

#. Extract the **only** the ``traffic_ops_db`` directory to ``/opt/traffic_ops_db``

	.. code-block:: shell

		tar -zxvf trafficcontrol-incubating-<version>.tar.gz --strip=1 trafficcontrol-incubating-<version>/traffic_ops_db
		chown -R postgres:postgres /opt/traffic_ops_db

Migration Preparation
---------------------
Be sure there is connectivity between your MySQL server's IP address/port and your PostgreSQL server's IP address/port.

Navigating the Database Migration
---------------------------------
Begin the database migration after settings up the ``/opt/traffic_ops_db`` directory
Switch to the postgres user, so that permissions remain intact.

.. code-block:: shell

	su - postgres
	cd /opt/traffic_ops_db/

#. Configure the ``/opt/traffic_ops_db/pg-migration/mysql-to-postgres.env`` file for your source MySQL and target PostgresQL settings. This part ought to be self-explanatory, given the names used in that file.

#. Run the ``migrate.sh`` script, watching the console output for any errors (this may take some time).

Your MySQL data should now be ported into your new instance of PostgreSQL!

.. [1] This migration was tested against version ``docker-engine-selinux-17.05.0.ce-1``
