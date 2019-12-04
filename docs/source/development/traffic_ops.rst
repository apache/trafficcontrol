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
At its current stage in development, "Traffic Ops" actually refers to a concept with two implementations. The original Traffic Ops was written as a collection of Perl scripts based on the `Mojolicious framework <https://mojolicious.org/>`_ framework. At some point, the relatively poor performance and lack of knowledgeable developers as the project expanded became serious issues, and so for the past few years Traffic Ops has been undergoing a steady rewrite to Go.

Introduction
============
Traffic Ops at its core is mainly a PostgreSQL database used to store configuration information for :abbr:`ATC (Apache Traffic Control)`, and a set of RESTful API endpoints for interacting with and manipulating that information. It also serves as the single point of authentication for :abbr:`ATC (Apache Traffic Control)` components (that is, when one hears "user" in an :abbr:`ATC (Apache Traffic Control)` context it nearly always means a "user" as configured in Traffic Ops) and provides interfaces to other :abbr:`ATC (Apache Traffic Control)` components by proxy. Additionally, there is some miscellaneous, at times obscure functionality to Traffic Ops, such as generating arbitrary Linux system images.

Software Requirements
=====================
Traffic Ops is only supported on CentOS 7+ systems (although many developers do use Mac OS with some success).

The two different implementations have different requirements, but they do share a few:

- `Goose <https://bitbucket.org/liamstask/goose/>`_ (although the ``postinstall`` Perl script will install this if desired)
- `PostgreSQL 9.6.6 <https://www.postgresql.org/download/>`_ - the machine where (either implementation of) Traffic Ops is running must have the client tool set (e.g. :manpage:`psql(1)`), but the actual database can be run anywhere so long as it is accessible.
- :manpage:`openssl(1SSL)` is recommended to generate server certificates, though not strictly required if certificates can be obtained by other means.
- Some kind of SMTP server is required for certain :ref:`to-api` endpoints to work, but for purposes unrelated to them an SMTP server is not required.

.. tip:: Alternatively, development and testing can be done using :ref:`ciab` - albeit somewhat more slowly.

Perl Implementation Requirements
--------------------------------
Most dependencies are managed by `Carton 1.0.12 <http://search.cpan.org/~miyagawa/Carton-v1.0.12/lib/Carton.pm>`_, but there are some - outside of those shared with the Go implementation - that are not managed by that system.

- `Carton itself <http://search.cpan.org/~miyagawa/Carton-v1.0.12/lib/Carton.pm>`_
- Perl 5.10.1
- libpcap and libpcap development library - usually ``libpcap-dev`` or ``libpcap-devel`` in your system's native package manager.
- libpq and libpq development library - usually ``libpq-dev`` or ``libpq-devel`` in your system's native package manager.
- The `JSON Perl pod from CPAN <https://metacpan.org/pod/JSON>`_
- The `JSON::PP Perl pod from CPAN <https://metacpan.org/pod/JSON::PP>`_
- Developers should use `Perltidy <http://perltidy.sourceforge.net/>`_ to format their Perl code.

	.. code-block:: text
		:caption: Example Perltidy Configuration (usually in :file:`{HOME}/.perltidyrc`)

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


Go Implementation Requirements
------------------------------
- `Go 1.11 <http://golang.org/doc/install>`_
- If the system's Go compiler doesn't provide it implicitly, also note that all Go code in the :abbr:`ATC (Apache Traffic Control)` repository should be formatted using `gofmt <https://golang.org/cmd/gofmt/>`_

Traffic Ops Project Tree Overview
=================================
- :atc-file:`traffic_ops/` - The root of the Traffic Ops project

	- app/ - Holds most of the Perl code base, though many of the files contained herein are also used by the Go implementation

		.. note:: This directory is home to many things that no longer work as intended or have been superseded by other things - most notably code for the now-removed Traffic Ops UI. That does *not*, however, mean that they are safe to remove. The API code that is still relied upon today is deeply entangled with the UI code, and in a dynamic language like Perl it can be very dangerous to remove things, because it may not be apparent that something is broken until it's already in production. So please don't remove anything in here until we're ready to excise the Perl implementation entirely.

		- bin/ - Directory for scripts and tools, :manpage:`cron(8)` jobs, etc.

			- checks/ - Contains the :ref:`to-check-ext` scripts that are provided by default
			- db/ - Contains scripts that manipulate the database beyond the scope of setup, migration, and seeding
			- tests/ - Integration and unit test scripts for automation purposes - in general this has been superseded by :atc-file:`traffic_ops/testing/api/`\ [#perltest]_

		- conf/ - Aggregated configuration for Traffic Ops. For convenience, different environments for the :ref:`database-management` tool are already set up

			- development/ - Configuration files for the "development" environment
			- integration/ - Configuration files for the "integration" environment
			- misc/ - Miscellaneous configuration files.
			- production/ - Configuration files for the "production" environment
			- test/ - Configuration files for the "test" environment

		- db/ - Database setup, seeding, and upgrade/downgrade helpers

			- migrations/ - Database migration files
			- tools/ - Contains helper scripts for easing upgrade transitions when selective data manipulation must be done to achieve a desirable state

		- lib/ - Contains the main handling logic for the Perl implementation

			- API/ - Mojolicious Controllers for the :ref:`to-api`
			- Common/ - Code that is shared between both the :ref:`to-api` and the now-removed Traffic Ops UI
			- Connection/ - Adapter definitions for connecting to external services
			- Extensions/ - Contains :ref:`to-datasource-ext`
			- Fixtures/ - Test-case fixture data for the "testing" environment\ [#perltest]_

				- Integration/ - Integration tests\ [#perltest]_

			- Helpers/ - Contains route handlers for the Traffic Stats-related endpoints
			- MojoPlugins/ - Mojolicious Plugins for common controller code
			- Schema/Result/ - Contains schema definitions generated from a constructed database for use with the `DBIx Perl pod suite <https://metacpan.org/search?q=DBIx>`_. These were machine-generated in 2016 and *absolutely* **no one** *should be touching them ever again*.
			- /Test - Common helpers for testing
			- UI/ - Mojolicious controllers for the now-removed Traffic Ops UI
			- Utils/ - Contains helpful utilities for certain objects and tasks

				- Helper/ - Common utilities for the Traffic Ops application

		- public/ - A directory from which files are served statically over HTTP by the Perl implementation. One common use-case is for hosting a :term:`Coverage Zone File` for Traffic Router.
		- script/ - Mojolicious bootstrap/startup scripts.
		- t/ - Unit tests for both the API (in ``api/``) and the now-removed Traffic Ops UI\ [#perltest]_

			- api/ - Unit tests for the API\ [#perltest]_

		- t_integration/ - High-level integration tests\ [#perltest]_
		- templates/ - Mojolicious Embedded Perl (:file:`{template name}.ep`) files for the now-removed Traffic Ops UI

	- build/ - Contains files that are responsible for packaging Traffic Ops into an RPM file - and also for doing the same with :term:`ORT`
	- client/ - The Go client library for Traffic Ops
	- etc/ - Configuration files for various systems associated with running production instances of Traffic Ops, which are installed under ``/etc`` by the Traffic Ops RPM

		- cron.d/ - Holds specifications for :manpage:`cron(8)` jobs that need to be run periodically on Traffic Ops servers

			.. note:: At least one of these jobs expects itself to be run on a server that has the Perl implementation of Traffic Ops installed under ``/opt/traffic_ops/``. Nothing terrible will happen if that's not true, just that it/they won't work. Installation using the RPM will set up all of these kinds of things up automatically.

		- init.d/ - Contains the old, initscripts-based job control for Traffic Ops
		- logrotate.d/ - Specifications for the Linux :manpage:`logrotate(8)` utility for Traffic Ops log files
		- profile.d/traffic_ops.sh - Sets up common environment variables for working with Traffic Ops

	- install/ - Contains all of the resources necessary for a full install of Traffic Ops

		- bin/ - Binaries related to installing Traffic Ops, as well as installing its prerequisites, certificates, and database
		- data/ - Contains things that need to be accessible by the running server for certain functionality - typically installed to ``/var/www/data`` by the RPM (hence the name).
		- etc/ - This directory left empty; it's used to contain post-installation extensions and resources
		- lib/ - Contains libraries used by the various installation binaries

	- ort/ - Contains :term:`ORT` and :abbr:`ATS (Apache Traffic Server)` configuration file-generation logic and tooling
	- testing/ - Holds utilities for testing the :ref:`to-api`

		- api/ - Integration testing for the `Traffic Ops Go client <https://godoc.org/github.com/apache/trafficcontrol/traffic_ops/client>`_ and Traffic Ops
		- compare/ - Contains :ref:`compare-tool`

	- traffic_ops_golang/ - The root of the Go implementation's code-base

		.. note:: The vast majority of subdirectories of :atc-file:`traffic_ops/traffic_ops_golang/` contain handlers for the :ref:`to-api`, and are named according to the endpoint they handle. What follows is a list of subdirectories of interest that have a special role (i.e. don't handle a :ref:`to-api` endpoint).

		.. seealso:: `The GoDoc documentation for this package <https://godoc.org/apache/trafficcontrol/traffic_ops/traffic_ops_golang>`_

		- api/ - A library for use by :ref:`to-api` handlers that provides helpful utilities for common tasks like obtaining a database transaction handle or accessing Traffic Ops configuration
		- auth/ - Contains definitions of privilege levels and access control code used in routing and provides a library for dealing with password and token-based authentication
		- config/ - Defines configuration structures and methods for reading them in from files
		- dbhelpers/ - Assorted utilities that provide functionality for common database tasks, e.g. "Get a user by email"
		- plugin/ - The Traffic Ops plugin system, with examples
		- riaksvc/ - In addition to handling routes that deal with storing secrets in or retrieving secrets from Traffic Vault, this package provides a library of functions for interacting with Traffic Vault for other handlers to use.
		- routing/ - Contains logic for mapping all of the :ref:`to-api` endpoints to their handlers, as well as proxying requests back to the Perl implementation and managing plugins, and also provides some wrappers around registered handlers that set common HTTP headers and connection options
		- swaggerdocs/ A currently abandoned attempt at defining the :ref:`to-api` using `Swagger <https://swagger.io/>`_ - it may be picked up again at some point in the (distant) future
		- tenant/ - Contains utilities for dealing with :term:`Tenantable <Tenant>` resources, particularly for checking for permissions
		- tocookie/ - Defines the method of generating the ``mojolicious`` cookie used by Traffic Ops for authentication
		- vendor/ - contains "vendored" Go packages from third party sources

	- vendor/ - contains "vendored" Go packages from third party sources

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

The environments are defined in the :atc-file:`traffic_ops/app/db/dbconf.yml` file, and the name of the database generated will be the name of the environment for which it was created.

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

#. Run the :atc-file:`traffic_ops/install/bin/postinstall` script, it will prompt for information like the default user login credentials
#. To run Traffic Ops, follow the instructions in :ref:`to-running`.

Test Cases
==========

Perl Tests
----------
Use `prove <http://perldoc.perl.org/prove.html>`_ (should be installed with Perl) to execute test cases\ [#perltest]_. Execute after a ``carton install`` of all required dependencies:

- To run the Unit Tests: ``prove -qrp  app/t/``
- To run the Integration Tests: ``prove -qrp app/t_integration/``

Go Tests
--------
Many (but not all) endpoint handlers and utility packages in the Go code-base define Go unit tests that can be run with :manpage:`go-test(1)`. There are integration tests for the Traffic Ops Go client in :atc-file:`traffic_ops/testing/api/`.

.. code-block:: bash
	:caption: Sample Run of Go Unit Tests

	cd traffic_ops/traffic_ops_golang

	# run just one test
	go test ./about

	# run all of the tests
	go test ./...

There are a few prerequisites to running the Go client integration tests:

- A PostgreSQL server must be accessible and have a Traffic Ops database schema set up (though not necessarily populated with anything).
- A running Traffic Ops Go implementation instance must be accessible - it shouldn't be necessary to also be running the Perl implementation.

	.. note:: For testing purposes, SSL certificates are not verified, so self-signed certificates will work fine.

	.. note:: It is *highly* recommended that the Traffic Ops instance be run on the same machine as the integration tests, as otherwise network latency can cause the tests to exceed their threshold time limit of ten minutes.

The integration tests are run using :manpage:`go-test(1)`, with two configuration options available.

.. note:: It should be noted that the integration tests will output thousands of lines of highly repetitive text not directly related to the tests its running at the time - even if the ``-v`` flag is not passed to :manpage:`go-test(1)`. This problem is tracked by :issue:`4017`.

.. warning:: Running the tests will wipe the connected database clean, so do not **ever** run it on an instance of Traffic Ops that holds meaningful data.

.. option:: --cfg CONFIG

	Specify the path to a configuration file for the tests. If not specified, it will attempt to read a file named ``traffic-ops-test.config`` in the working directory.

	.. seealso:: `Configuring the Integration Tests`_ for a detailed explanation of the format of this configuration file.
.. option:: --fixtures FIXTURES

	Specify the path to a file containing static data for the tests to use. This should almost never be used, because many of the tests depend on the data having a certain content and structure. If not specified, it will attempt to read a file named ``tc-fixtures.json`` in the working directory.

Configuring the Integration Tests
"""""""""""""""""""""""""""""""""
Configuration is mainly done through the configuration file passed as :option:`--cfg`, but is also available through the following environment variables.

.. envvar:: SESSION_TIMEOUT_IN_SECS

	Sets the timeout of requests made to the Traffic Ops instance, in seconds.

.. envvar:: TODB_DESCRIPTION

	An utterly cosmetic variable which, if set, gives a description of the PostgreSQL database to which the tests will connect. This has no effect except possibly changing one line of debug output.

.. envvar:: TODB_HOSTNAME

	If set, will define the :abbr:`FQDN (Fully Qualified Domain Name)` at which the PostgreSQL server to be used by the tests resides\ [#integrationdb]_.

.. envvar:: TODB_NAME

	If set, will define the name of the database to which the tests will connect\ [#integrationdb]_.

.. envvar:: TODB_PASSWORD

	If set, defines the password to use when authenticating with the PostgreSQL server.

.. envvar:: TODB_PORT

	If set, defines the port on which the PostgreSQL server listens\ [#integrationdb]_.

.. envvar:: TODB_SSL

	If set, must be one of the following values:

	true
		The PostgreSQL server to which the tests will connect uses SSL on the port on which it will be contacted.
	false
		The PostgreSQL server to which the tests will connect does not use SSL on the port on which it will be contacted.

.. envvar:: TODB_TYPE

	If set, tells the database driver used by the tests the kind of SQL database to which they are connecting\ [#integrationdb]_. This author has no idea what will happen if this is set to something other than ``Pg``, but it's possible the tests will fail to run. Certainly never do it.

.. envvar:: TODB_USER

	If set, defines the user as whom to authenticate with the PostgreSQL server.

.. envvar:: TO_URL

	If set, will define the URL at which the Traffic Ops instance is running - including port number.

.. envvar:: TO_USER_DISALLOWED

	If set, will define the name of a user with the "disallowed" :term:`Role` that will be created by the tests\ [#existinguser]_.

.. envvar:: TO_USER_EXTENSION

	If set, will define the name of a user with the "extension" :term:`Role` that will be created by the tests\ [#existinguser]_.

	.. caution:: Due to legacy constraints, the only truly safe value for this is ``extension`` - anything else could cause the tests to fail.

.. envvar:: TO_USER_FEDERATION

	If set, will define the name of a user with the "federation" :term:`Role` that will be created by the tests\ [#existinguser]_.

.. envvar:: TO_USER_OPERATIONS

	If set, will define the name of a user with the "operations" :term:`Role` that will be created by the tests\ [#existinguser]_.

.. envvar:: TO_USER_PASSWORD

	If set, will define the password used by all users created by the tests. This does not need to be the password of any pre-existing user.

.. envvar:: TO_USER_PORTAL

	If set, will define the name of a user with the "portal" :term:`Role` that will be created by the tests\ [#existinguser]_.

.. envvar:: TO_USER_READ_ONLY

	If set, will define the name of a user with the "read-only" :term:`Role` that will be created by the tests\ [#existinguser]_.


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

.. [#perltest] As progress continues on moving Traffic Ops to run entirely in Go, the number of passing tests has steadily decreased. This means that the tests are not a reliable way to test Traffic Ops, as they are expected to fail more and more as functionality is stripped from the Perl codebase.
.. [#integrationdb] The Traffic Ops instance *must* be using the same PostgreSQL database that the tests will use.
.. [#existinguser] This does not need to match the name of any pre-existing user.
