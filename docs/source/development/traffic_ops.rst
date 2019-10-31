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

***********
Traffic Ops
***********

Introduction
============
Traffic Ops uses a PostgreSQL database to store the configuration information, and the `Mojolicious framework <http://mojolicio.us/>`_ to generate the user interface and REST APIs.

Software Requirements
=====================
To work on Traffic Ops you need a CentOS 7+ environment that has the following installed:

- `Carton 1.0.12 <http://search.cpan.org/~miyagawa/Carton-v1.0.12/lib/Carton.pm>`_

	- libpcap (plus development library - usually "libpcap-dev" or "libpcap-devel")
	- libpq (plus development library - usually "libpq-dev" or "libpq-devel")
	- cpan JSON
	- cpan JSON\:\:PP

- `Go 1.8.3 <http://golang.org/doc/install>`_
- Perl 5.10.1
- Git
- PostgreSQL 9.6.6
- `Goose <https://bitbucket.org/liamstask/goose/>`_

Traffic Ops Project Tree Overview
=================================
traffic_ops/ - The root of the Traffic Ops project

	- app/ - Holds most of the Perl code base

		- bin/ - Directory for scripts, :manpage:`cron(8)` jobs, etc
		- conf/

			- development/ - Development (local) specific configuration files.
			- misc/ - Miscellaneous configuration files.
			- production/ - Production specific configuration files.
			- test/ - Test (unit test) specific configuration files.

		- db/ - Database related area.

			- migrations/ - Database Migration files.

		- lib/

			- API/ - Mojolicious Controllers for the :ref:`to-api`
			- Common/ - Common Code between both the :ref:`to-api` and the deprecated Traffic Ops UI
			- Extensions/ - Contains :ref:`to-datasource-ext`
			- Fixtures/ - Test Case fixture data for the 'to_test' database.

				- Integration/ - Integration Tests.

			- MojoPlugins/ - Mojolicious Plugins for Common Controller Code.
			- Schema/ - Database Schema area.

				- /Result - DBIx ORM related files.

			- /Test - Common Test.
			- UI/ - Mojolicious Controllers for the deprecated Traffic Ops UI.
			- Utils/

				- Helper/ - Common utilities for the Traffic Ops application.

		- log/ - Log directory where the development and test files are written
		- public/

		 - css/ - Stylesheets
		 - images/ - Images
		 - js/ - Javascripts

		- script/ - Mojolicious Bootstrap scripts.
		- t/ - Unit Tests for the UI.

		 - api/ - Unit Tests for the API.

		- t_integration/ - High level tests for Integration level testing.
		- templates/ - Mojolicious Embedded Perl (``*.ep``) files for the UI.

	- bin/ - holds executables related to Traffic Ops, but not actually a part of the Traffic Ops server's operation
	- build/ - contains files that are responsible for packaging Traffic Ops into an RPM file
	- client/ - API endpoints handled by Go
	- client_tests/ - lol
	- doc/ - contains only a :file:`coverage-zone.json` example (?) file
	- etc/ - configuration files needed for the Traffic Ops server

		- cron.d/ - holds specifications for :manpage:`cron(8)` jobs that need to be run periodically on Traffic Ops servers
		- init.d/ - contains the old initscripts-based job control for Traffic Ops
		- logrotate.d/ - specifications for the Linux :manpage:`logrotate(8)` utility for Traffic Ops log files
		- profile.d/traffic_ops.sh - sets up common environment variables for working with Traffic Ops

	- experimental/ - includes all kinds of prototype and/or abandoned tools and extensions

		- ats_config/ - an attempt to provide an easier method of obtaining and/or writing configuration files for :abbr:`ATS (Apache Traffic Server)` :term:`cache servers`
		- auth/ - a simple authentication server that mimics the authentication process of Traffic Ops, and provides a detailed view of a logged-in user's permissions and capabilities
		- goto/ - an Angular (1.x) web page backed by a Go server that provides a ReST API interface for mySQL servers
		- postgrest/ - originally probably going to be a web server that provides a ReST API for postgreSQL servers, this only contains a simple - albeit unfinished - Docker container specification for running postgreSQL client tools and/or server(s)
		- server/ - a living copy of the original attempt at re-writing Traffic Ops in Go
		- traffic_ops_auth/ - proof-of-concept for authenticating, creating and deleting users in a Traffic Ops schema.
		- url-rewriter-nginx/ - Docker container specification for a modification to the NginX web server, meant to make it suitable for use as a caching server at the Edge-tier or Mid-tier levels of the Traffic Control architecture
		- webfront/ - a simple HTTP caching server written from the ground-up, meant to be suitable as a caching server at the Edge-tier or Mid-tier levels of the Traffic Control architecture

	- install/ - contains all of the resources necessary for a full install of Traffic Ops

		- bin/ - binaries related to installing Traffic Ops, as well as installing its prerequisites, certificates, and database
		- data/ - almost nothing
		- etc/ - this directory left empty; it's used to contain post-installation extensions and resources
		- lib/ - contains libraries used by the various installation binaries

	- testing/ - holds utilities for testing the :ref:`to-api`, as well as comparing two separate API instances (for e.g. comparing a new build to a known-to-work build)
	- traffic_ops_golang/ - has all of the functionality that has been re-written from Perl into Go
	- vendor/ - contains "vendored" packages from third party sources

Perl Formatting Conventions
===========================
`Perltidy <http://perltidy.sourceforge.net/>`_ is for use in code formatting.

.. code-block:: text
	:caption: Example Perltidy Configuration (usually in :file:`~/.perltidyrc`)

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

.. _database-management:

.. program:: admin

app/db/admin
============
The :program:`app/db/admin` binary is for use in managing the Traffic Ops database tables. This essentially serves as a front-end for `Goose <https://bitbucket.org/liamstask/goose/>`_.

.. note:: For proper resolution of configuration and SOL statement files, it's recommended that this binary be run from the ``app`` directory

Usage
-----
``db/admin [options] command``

Options and Arguments
---------------------
.. option:: --env ENVIRONMENT

	An optional environment specification that causes the database configuration to be read out of the corresponding section of the :file:`app/db/dbconf.yml` configuration file. One of:

	- development
	- integration
	- production
	- test

	(Default: ``development``)

.. envvar:: MOJO_MODE

	:program:`admin` sets this to the value of the environment as specified by :option:`--env` (Default: ``development``)

.. option:: command

	The :option:`command` specifies the operation to be performed on the database. It must be one of:

	createdb
		Creates the database for the current environment
	create_user
		Creates the user defined for the current environment
	dbversion
		Displays the database version that results from the current sequence of migrations
	down
		Rolls back a single migration from the current version
	drop
		Drops the database for the current environment
	drop_user
		Drops the user defined for the current environment
	load_schema
		Sets up the database for the current environment according to the SQL statements in ``app/db/create_tables.sql``
	migrate
		Runs a migration on the database for the current environment
	patch
		Patches the database for the current environment using the SQL statements from the ``app/db/patches.sql``
	redo
		Rolls back the most recently applied migration, then run it again
	reset
		Creates the user defined for the current environment, drops the database for the current environment, creates a new one, loads the schema into it, and runs a single migration on it
	reverse_schema
		Reverse engineers the ``app/lib/Schema/Result/*`` files from the environment database
	seed
		Executes the SQL statements from the ``app/db/seeds.sql`` file for loading static data
	show_users
		Displays a list of all users registered with the PostgreSQL server
	status
		Prints the status of all migrations
	upgrade
		Performs a migration on the database for the current environment, then seeds it and patches it using the SQL statements from the ``app/db/patches.sql`` file

.. code-block:: bash
	:caption: Example Usage

	db/admin --env=test reset

The environments are defined in the :file:`app/db/dbconf.yml` file, and the name of the database generated will be the name of the environment for which it was created.

Installing The Developer Environment
====================================
To install the Traffic Ops Developer environment:

#. Clone the `Traffic Control repository <https://github.com/apache/trafficcontrol>`_ from GitHub.
#. Install the local dependencies using `Carton <https://metacpan.org/release/Carton>`_.

	.. code-block:: shell
		:caption: Install Development Dependencies

		cd traffic_ops/app
		carton

#. Set up a role (user) in PostgreSQL

	.. seealso:: `PostgreSQL instructions on setting up a database <https://wiki.postgresql.org/wiki/First_steps>`_.


#. Use the ``reset`` and ``upgrade`` :option:`command`\ s of :program:`admin` (see :ref:`database-management` for usage) to set up the ``traffic_ops`` database(s).
#. (Optional) To load the 'KableTown' example/testing data set into the tables, use the :file:`app/bin/db/setup_kabletown.pl` script.

	.. note:: To ensure proper paths to Perl libraries and resource files, ``setup_kabletown.pl`` should be run from within the ``app/`` directory.

#. Run the ``postinstall`` script, located in ``install/bin/``

#. To start Traffic Ops, use the ``start.pl`` script located in the ``app/bin`` directory. If the server starts successfully, the STDOUT of the process should contain the line ``[<date and time>] [INFO] Listening at "http://*:3000"``, followed by the line ``Server available at http://127.0.0.1:3000`` (using default settings for port number and listening address, and where ``<date and time>`` is an actual date and time in ISO format).

	.. note:: To ensure proper paths to Perl libraries and resource files, the ``start.pl`` script should be run from within the ``app/`` directory.

#. Using a web browser, navigate to the given address: ``http://127.0.0.1:3000``
#. A prompt for login credentials should appear. Assuming default settings are used, the initial login credentials will be

	:User name: ``admin``
	:Password:  ``password``

#. Change the login credentials.

Test Cases
==========
Use `prove <http://perldoc.perl.org/prove.html>`_ (should be installed with Perl) to execute test cases. Execute after a ``carton install`` of all required dependencies:

- To run the Unit Tests: ``prove -qrp  app/t/``
- To run the Integration Tests: ``prove -qrp app/t_integration/``

.. note:: As progress continues on moving Traffic Ops to run entirely in Go, the number of passing tests has steadily decreased. This means that the tests are not a reliable way to test Traffic Ops, as they are expected to fail more and more as functionality is stripped from the Perl codebase.

The KableTown CDN example
-------------------------
The integration tests will load an example CDN with most of the features of Traffic Control being used. This is mostly for testing purposes, but can also be used as an example of how to configure certain features. To load the KableTown CDN example and access it:

#. Be sure the integration tests have been run
#. Start the Traffic Ops server. The :envvar:`MOJO_MODE` environment variable should be set to the name of the environment that has been loaded.

	.. code-block:: bash
		:caption: Example Startup

		export MOJO_MODE=integration
		cd app/
		bin/start.pl

#. Using a web browser, navigate to the address Traffic Ops is serving, e.g. ``http://127.0.0.1:3000`` for default settings
#. For the initial log in:

	:User name: ``admin``
	:Password: ``password``


Extensions
==========
Traffic Ops Extensions are a way to enhance the basic functionality of Traffic Ops in a customizable manner. There are two types of extensions:

:ref:`to-check-ext`
	These allow you to add custom checks to the :menuselection:`Monitor --> Cache Checks` view.

:ref:`to-datasource-ext`
	These allow you to add statistic sources for the graph views and APIs.

Extensions are managed using the ``$TO_HOME/bin/extensions`` command line script

.. seealso:: For more information see :ref:`admin-to-ext-script`.


Extensions at Runtime
---------------------
The search path for :ref:`to-datasource-ext` depends on the configuration of the ``PERL5LIB`` environment variable, which is pre-configured in the Traffic Ops start scripts. All :ref:`to-check-ext` must be located in ``$TO_HOME/bin/checks``

	.. code-block:: bash
		:caption: Example ``PERL5LIB`` Configuration

		export PERL5LIB=/opt/traffic_ops_extensions/private/lib/Extensions:/opt/traffic_ops/app/lib/Extensions/TrafficStats

To prevent :ref:`to-datasource-ext` namespace collisions within Traffic Ops all :ref:`to-datasource-ext` should follow the package naming convention '``Extensions::<ExtensionName>``'

``TrafficOpsRoutes.pm``
-----------------------
Traffic Ops accesses each extension through the addition of a URL route as a custom hook. These routes will be defined in a file called ``TrafficOpsRoutes.pm`` that should be present in the top directory of your Extension. The routes that are defined should follow the `Mojolicious route conventions <https://mojolicious.org/perldoc/Mojolicious/Guides/Routing#Routes>`_.


Development Configuration
--------------------------
To incorporate any custom :ref:`to-datasource-ext` during development set your ``PERL5LIB`` environment variable with any number of colon-separated directories with the understanding that the ``PERL5LIB`` search order is from left to right through this list. Once Perl locates your custom route or Perl package/class it 'pins' on that class or Mojolicious Route and doesn't look any further, which allows for the developer to override Traffic Ops functionality.
