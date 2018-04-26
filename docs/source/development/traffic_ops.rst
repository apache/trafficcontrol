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

Traffic Ops
***********

Introduction
============
Traffic Ops uses a Postgres database to store the configuration information, and the `Mojolicious framework <http://mojolicio.us/>`_ to generate the user interface and REST APIs. 

Software Requirements
=====================
To work on Traffic Ops you need a \*nix (MacOS and Linux are most commonly used) environment that has the following installed:

* `Carton 1.0.12 <http://search.cpan.org/~miyagawa/Carton-v1.0.12/lib/Carton.pm>`_

  * cpan JSON
  * cpan JSON::PP

* `Go 1.8.3 <http://golang.org/doc/install>`_
* Perl 5.10.1
* Git
* Postgres 9.6.6
* `Goose <https://bitbucket.org/liamstask/goose/>`_

Addionally, the installation of the following RPMs (or equivalent) is required:

* All RPMs listed in :ref:`rl-ps`

Traffic Ops Project Tree Overview
=======================================

**/opt/traffic_ops/app**

* bin/ - Directory for scripts, cronjobs, etc

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
  * Fixtures/ - Test Case fixture data for the 'to_test' database.
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
===========================
Perl tidy is for use in code formatting. See the following config file for formatting conventions.

::


  edit a file called $HOME/.perltidyrc

  -l=156
  -et=4
  -t
  -ci=4
  -st
  -se
  -vt=0
  -cti=0
  -pt=1
  -bt=1
  -sbt=1
  -bbt=1
  -nsfs
  -nolq
  -otr
  -aws
  -wls="= + - / * ."
  -wrs=\"= + - / * .\"
  -wbb="% + - * / x != == >= <= =~ < > | & **= += *= &= <<= &&= -= /= |= + >>= ||= .= %= ^= x="


Database Management
===================
..  Add db naming conventions

The admin.pl script is for use in managing the Traffic Ops database tables. Below is an example of its usage. 

``$ db/admin.pl``

Usage:  db/admin.pl [--env (development|test|production)] [arguments]

Example: ``db/admin.pl --env=test reset``

Purpose:  This script is used to manage the database. The environments are defined in the dbconf.yml, as well as the database names.

* To use the ``admin.pl`` script, you may need to add ``traffic_ops/lib`` and ``traffic_ops/local/lib/perl5`` to your `PERL5LIB <http://modperlbook.org/html/3-9-2-2-Using-the-PERL5LIB-environment-variable.html>`_ environment variable.

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
====================================
To install the Traffic Ops Developer environment:

1. Clone the traffic_control repository from `github.com <https://github.com/apache/incubator-trafficcontrol>`_.
2. Install the local dependencies using Carton (cpanfile).

  ::

   $ cd traffic_ops/app
   $ carton

3. Set up a role (user) in Postgres

   See Postgres instructions on initdb https://wiki.postgresql.org/wiki/First_steps 


4. Enter ``db/admin.pl --env=<enviroment name> setup`` to set up the traffic_ops database(s). 

   * Unit test database: ``$ db/admin.pl --env=test setup``
   * Development database: ``$ db/admin.pl --env=development setup``
   * Integration database: ``$ db/admin.pl --env=integration setup``

   |

   Running the the admin.pl script in setup mode should look like this: ::

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

5. (Optional) To load temporary data into the tables: ``$ perl bin/db/setup_kabletown.pl``

6. Run the postinstall script: ``traffic_ops/install/bin/postinstall``

7. To start Traffic Ops, enter ``$ bin/start.pl``

   The local Traffic Ops instance uses an open source framework called morbo, starting following the start command execution.

   Start up success includes the following:

  ::
   

   [2015-02-24 10:44:34,991] [INFO] Listening at "http://*:3000".
   
   Server available at http://127.0.0.1:3000.


8. Using a browser, navigate to the given address: ``http://127.0.0.1:3000``
9. For the initial log in:
  
  * User name: admin
  * Password: password

10. Change the log in information.

Test Cases
==========
Use prove to execute test cases. Execute after a carton install:

* To run the Unit Tests: ``$ local/bin/prove -qrp  t/``
* To run the Integration Tests: ``$ local/bin/prove -qrp t_integration/``

The KableTown CDN example
-------------------------
The integration tests will load an example CDN with most of the features of Traffic Control being used. This is mostly for testing purposes, but can also be used as an example of how to configure certain features. To load the KableTown CDN example and access it:

1. Run the integration tests 
2. Start morbo against the integration database: ``export MOJO_MODE=integration; ./bin/start.pl``
3. Using a browser, navigate to the given address: ``http://127.0.0.1:3000``
4. For the initial log in:
  
  * User name: admin
  * Password: password


Extensions
==========
Traffic Ops Extensions are a way to enhance the basic functionality of Traffic Ops in a custom manner. There are three types of extensions:

1. Check Extensions

  These allow you to add custom checks to the "Health->Server Checks" view.

2. Configuration Extensions

  These allow you to add custom configuration file generators.

3. Data source Extensions

  These allow you to add statistic sources for the graph views and APIs.

Extensions are managed using the $TO_HOME/bin/extensions command line script. For more information see :ref:`admin-to-ext-script`.

Check Extensions
----------------

In other words, check extensions are scripts that, after registering with Traffic Ops, have a column reserved in the "Health->Server Checks" view and that usually run periodically out of cron.

.. |checkmark| image:: ../../../traffic_ops/app/public/images/good.png 

.. |X| image:: ../../../traffic_ops/app/public/images/bad.png


It is the responsibility of the check extension script to iterate over the servers it wants to check and post the results.

An example script might proceed by logging into the Traffic Ops server using the HTTPS base_url provided on the command line. The script is hardcoded with an auth token that is also provisioned in the Traffic Ops User database. This token allows the script to obtain a cookie used in later communications with the Traffic Ops API. The script then obtains a list of all caches to be polled by accessing Traffic Ops' ``/api/1.1/servers.json`` REST target. This list is walked, running a command to gather the stats from that cache. For some extensions, an HTTP GET request might be made to the ATS astats plugin, while for others the cache might be pinged, or a command run over SSH. The results are then compiled into a numeric or boolean result and the script POSTs tha result back to the Traffic Ops ``/api/1.1/servercheck/`` target.

A check extension can have a column of |checkmark|'s and |X|'s (CHECK_EXTENSION_BOOL) or a column that shows a number (CHECK_EXTENSION_NUM).A simple example of a check extension of type CHECK_EXTENSION_NUM that will show 99.33 for all servers of type EDGE is shown below: :: 


  Script here.

Check Extension scripts are located in the $TO_HOME/bin/checks directory.

Currently, the following Check Extensions are available and installed by default:

**Cache Disk Usage Check - CDU**
  This check shows how much of the available total cache disk is in use. A "warm" cache should show 100.00.

**Cache Hit Ratio Check - CHR**
  The cache hit ratio for the cache in the last 15 minutes (the interval is determined by the cron entry). 

**DiffServe CodePoint Check - DSCP**
  Checks if the returning traffic from the cache has the correct DSCP value as assigned in the delivery service. (Some routers will overwrite DSCP)

**Maximum Transmission Check - MTU**
  Checks if the Traffic Ops host (if that is the one running the check) can send and receive 8192 size packets to the ``ip_address`` of the server in the server table.

**Operational Readiness Check - ORT**
  See :ref:`reference-traffic-ops-ort` for more information on the ort script. The ORT column shows how many changes the traffic_ops_ort.pl script would apply if it was run. The number in this column should be 0. 

**Ping Check - 10G, ILO, 10G6, FQDN**
  The bin/checks/ToPingCheck.pl is to check basic IP connectivity, and in the default setup it checks IP connectivity to the following:
  
  10G
    Is the ``ip_address`` (the main IPv4 address) from the server table pingable?
  ILO
    Is the ``ilo_ip_address`` (the lights-out-mangement IPv4 address) from the server table pingable?
  10G6
    Is the ``ip6_address`` (the main IPv6 address) from the server table pingable?
  FQDN 
    Is the Fully Qualified Domain name (the concatenation of ``host_name`` and ``.`` and ``domain_name`` from the server table) pingable?

**Traffic Router Check - RTR**
  Checks the state of each cache as perceived by all Traffic Monitors (via Traffic Router). This extension asks each Traffic Router for the state of the cache. A check failure is indicated if one or more monitors report an error for a cache. A cache is only marked as good if all reports are positive. (This is a pessimistic approach, opposite of how TM marks a cache as up, "the optimistic approach")
  

Configuration Extensions
------------------------
NOTE: Config Extensions are Beta at this time.


Data Source Extensions
----------------------
Traffic Ops has the ability to load custom code at runtime that allow any CDN user to build custom APIs for any requirement that Traffic Ops does not fulfill.  There are two classes of Data Source Extensions, private and public.  Private extensions are Traffic Ops extensions that are not publicly available, and should be kept in the /opt/traffic_ops_extensions/private/lib. Public extensions are Traffic Ops extensions that are Open Source in nature and free to enhance or contribute back to the Traffic Ops Open Source project and should be kept in /opt/traffic_ops/app/lib/Extensions.


Extensions at Runtime
---------------------
The search path for extensions depends on the configuration of the PERL5LIB, which is preconfigured in the Traffic Ops start scripts.  The following directory structure is where Traffic Ops will look for Extensions in this order.

Extension Directories
---------------------
PERL5LIB Example Configuration: ::

   export PERL5LIB=/opt/traffic_ops_extensions/private/lib/Extensions:/opt/traffic_ops/app/lib/Extensions/TrafficStats

Perl Package Naming Convention
------------------------------
To prevent Extension namespace collisions within Traffic Ops all Extensions should follow the package naming convention below:

Extensions::<ExtensionName>

Data Source Extension Perl package name example
Extensions::TrafficStats
Extensions::YourCustomExtension

TrafficOpsRoutes.pm
-------------------
Traffic Ops accesses each extension through the addition of a URL route as a custom hook.  These routes will be defined in a file called TrafficOpsRoutes.pm that should live in the top directory of your Extension.  The routes that are defined should follow the Mojolicious route conventions.


Development Configuration
--------------------------
To incorporate any custom Extensions during development set your PERL5LIB with any number of directories with the understanding that the PERL5LIB search order will come into play, so keep in mind that top-down is how your code will be located.  Once Perl locates your custom route or Perl package/class it 'pins' on that class or Mojo Route and doesn't look any further, which allows for the developer to *override* Traffic Ops functionality.
