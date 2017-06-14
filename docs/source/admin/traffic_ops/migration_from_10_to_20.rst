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
  
.. _rl-ps:

Traffic Ops - Migrating from 1.x to 2.x
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

In Traffic Ops 2.x MySQL was removed and Postgres was replaced as the database of choice for the unforeseen future.  A Docker-based migration tool was developed to
help with that conversion using an open source Postgres tool called `pgloader <http://pgloader.io/>`_.  The following instructions will help configuring the Migration tool

System Requirements
-------------------
The user must have the following for a successful minimal install:

* CentOS 7.2+
* Docker installed (this migration was tested against version **docker-engine-selinux-17.05.0.ce-1.el7.centos.noarch.rpm**)
* Postgres has been installed according to :ref:`rl-to-install`

Setup the traffic_ops_db directory
----------------------------------

   Modify /opt dir permission to make it writable and owned by postgres:postgres

   ::

   $ sudo chmod 755 /opt 
   
   Download the Traffic Control tarball for 2.0.0

   :: 

     $ cd /opt
     $ wget https://dist.apache.org/repos/dist/release/incubator/trafficcontrol/<tarball_version>

   Extract the **traffic_ops_db** dir to **/opt/traffic_ops_db**

   :: 

   $ tar -zxvf trafficcontrol-incubating-<version>.tar.gz --strip=1 trafficcontrol-incubating-<version>/traffic_ops_db
   $ sudo chown -R postgres:postgres /opt/traffic_ops_db

.. highlight:: none

Migration Preparation
---------------------
Be sure there is connectivity between your MySQL server's IP address/port and your Postgres server's IP address/port.

Navigating the Database Migration
---------------------------------
Begin the database migration after settings up the **/opt/traffic_ops_db** directory

   Switch to the postgres user so permissions stay intact.
   :: 

   $ su - postgres
   $ cd /opt/traffic_ops_db/

1. Configure the **/opt/traffic_ops_db/pg-migration/mysql-to-postgres.env** migration for your source MySQL and target Postgres settings


2. Run the migration, watch the console output for any errors (it may take some time)
   :: 

   $ ./migrate.sh


  Your MySQL data should now be ported into your new instance of Postgres!





