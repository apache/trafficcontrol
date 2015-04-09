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

Traffic Ops
============
Traffic Ops uses a MySql or Postgres database to store the configuration information, and the `Mojolicious framework <http://mojolicio.us/>`_ to generate the user interface and REST APIs. 

Software Requirements
---------------------
To work on Traffic Ops you need a \*nix (MacOS and Linux are most commonly used) environment that has the following installed:

* `Carton 1.0.12 <http://search.cpan.org/~miyagawa/Carton-v1.0.12/lib/Carton.pm>`_
* `Go 1.4 <http://golang.org/doc/install>`_
* Perl 5.10.1
* Git
* MySQL 5.1.52

Traffic Ops Project Tree Overview
---------------------------------

**/opt/traffic_ops/app**

* bin/ - Directory for scripts, cronjobs, etc.

* conf/

  * /development - Development (local) specific config files.
  * /misc - Miscellaneous config files.
  * /production - Production specific config files.
  * /test - Test (unit test) specific config files.

* db/ - Database related area.

  * /migrations - Database Migration files.

* lib/

  * /API - Mojo Controllers for the /API area of the application.
  * /Common - Common Code between both the API and UI areas.
  * /Extensions      
  * Fixtures/ - Test Case fixture data for the ‘to_test’ database.
    * /Integration - Integration Tests.
  * /MojoPlugins - Mojolicious Plugins for Common Controller Code.
  * Schema/ - Database Schema area.
    * /Result - DBIx ORM related files.
  * /Test - Common Test. 
  * /UI - Mojo Controllers for the Traffic Ops UI itself.
  * Utils/           
    * /Helper - Common utilities for the Traffic Ops application.

* log/ - Log directory where the development and test files are written by the app.

* public/
             
 * css/ - Stylesheets.
 * images/ - Images.
 * js/ - Javascripts

* script/ - Mojo Bootstrap scripts.
   
* t/ - Unit Tests for the UI.

 * api/ - Unit Tests for the API.

* t_integration/ - High level tests for Integration level testing.

* templates/ - Mojo Embedded Perl (.ep) files for the UI.



Perl Formatting Conventions 
---------------------------
Perl tidy is for use in code formatting. See the following config file for formatting conventions.

::


  edit a file called $HOME/.perltidyrc

  l = 156
  et=4
  t
  ci=4
  st
  se
  vt=0
  cti=0
  pt=1
  bt=1
  sbt=1
  bbt=1
  nsfs
  nolq
  otr
  aws
  wls="= + - / * ."
  wrs=\"= + - / * .\"
  wbb =% + - * / x != == >= <= =~ < > | & **= += *= &= <<= &&= -= /= |= + >>= ||= .= %= ^= x= 


Database Management
-------------------
..  Add db naming conventions

The admin.pl script is for use in managing the Traffic Ops database tables. Below is an example of its usage. 

``$ db/admin.pl``

Usage:  db/admin.pl [--env (development|test|production)] [arguments]

Example: ``db/admin.pl --env=test reset``

Purpose:  This script is used to manage the database. The environments are defined in the dbconf.yml, as well as the database names.

+-----------+--------------------------------------------------------------------+
| Arguments | Description                                                        |
+===========+====================================================================+
| create    | Execute db 'create' the database for the current environment.      |
+-----------+--------------------------------------------------------------------+
| down      | Roll back a single migration from the current version.             |
+-----------+--------------------------------------------------------------------+
| drop      | Execute db 'drop' on the database for the current environment.     |
+-----------+--------------------------------------------------------------------+
| redo      | Roll back the most recently applied migration, then run it again.  |
+-----------+--------------------------------------------------------------------+
| reset     | Execute db drop, create, load_schema, migrate on the database for  |
|           | the current environment.                                           |
+-----------+--------------------------------------------------------------------+
| seed      | Execute SQL from db/seeds.sql for loading static data.             |
+-----------+--------------------------------------------------------------------+
| setup     | Execute db drop, create, load_schema, migrate, seed on the         |
|           | database for the current environment.                              |
+-----------+--------------------------------------------------------------------+
| status    | Print the status of all migrations.                                |
+-----------+--------------------------------------------------------------------+
| upgrade   | Execute migrate then seed on the database for the current          |
|           | environment.                                                       |
+-----------+--------------------------------------------------------------------+

Installing The Developer Environment
------------------------------------
To install the Traffic Ops Developer environment:

1. Clone the traffic_control repository using Git.
2. Install the local dependencies using Carton (cpanfile).

  ::

   $ cd trafficops/app
   $ carton

3. Enter ``db/admin.pl --env=test setup`` to set up the traffic_ops database. 

   * Unit test database: ``$ db/admin.pl --env=test setup``
   * Development database: ``$ db/admin.pl --env=development setup``
   * Integration database: ``$ db use db/admin.pl --env=integration setup``
4. (Optional) To load temporary data into the tables: ``$ perl bin/db/setup_kabletown.pl``
5. Set up a user in the database. 

 ::


  master $ db/admin.pl --env=development setup
  Using database.conf: conf/development/database.conf
  Using database.conf: conf/development/database.conf
  Using database.conf: conf/development/database.conf
  Using database.conf: conf/development/database.conf
  Using database.conf: conf/development/database.conf
  Using database.conf: conf/development/database.conf
  Executing 'drop database to_development'
  Executing 'create database to_development'
  Creating database tables...
  Warning: Using a password on the command line interface can be insecure.
  Migrating database...
  goose: migrating db environment 'development', current version: 0, target: 20150210100000
  OK    20141222103718_extension.sql
  OK    20150108100000_add_job_deliveryservice.sql
  OK    20150205100000_cg_location.sql
  OK    20150209100000_cran_to_asn.sql
  OK    20150210100000_ds_keyinfo.sql
  Seeding database...
  Warning: Using a password on the command line interface can be insecure.

5. Set up a user in MySQL.

 ::

  master $ mysql
  Welcome to the MySQL monitor.  Commands end with ; or \g.
  Your MySQL connection id is 305
  Server version: 5.6.19 Homebrew

  Copyright (c) 2000, 2014, Oracle and/or its affiliates. All rights reserved.

  Oracle is a registered trademark of Oracle Corporation and/or its
  affiliates. Other names may be trademarks of their respective
  owners.

  Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

  mysql> create user ‘to_user’@’localhost’;
  mysql> grant all on to_development.* to 'to_user'@'localhost' identified by 'twelve';
  mysql> grant all on to_test.* to 'to_user'@'localhost' identified by 'twelve';
  mysql> grant all on to_integration.* to 'to_user'@'localhost' identified by 'twelve';


6. To start Traffic Ops, enter ``$ bin/start.sh``

   The local Traffic Ops instance uses an open source framework called morbo, starting following the start command execution.

   Start up success includes the following:

  ::
   

   [2015-02-24 10:44:34,991] [INFO] Listening at "http://*:3000".
   
   Server available at http://127.0.0.1:3000.


7. Using a browser, navigate to the given address: ``http://127.0.0.1:3000``
8. For the initial log in:
  
  * User name: admin
  * Password: password

9. Change the log in information.

Test Cases
----------
Use prove to execute test cases. Execute after a carton install:

* To run the Unit Tests: ``$ local/bin/prove -qrp  t/``
* To run the Integration Tests: ``$ local/bin/prove -qrp t_integration/``

The KableTown CDN example
^^^^^^^^^^^^^^^^^^^^^^^^^

Traffic Ops Extensions
----------------------


