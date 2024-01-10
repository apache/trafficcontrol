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

.. _dev-traffic-ops:

***********
Traffic Ops
***********
At one point, Traffic Ops was a collection of Perl scripts, and while the current program is written in Go, many of its tools and utilities are still written in Perl.

Introduction
============
Traffic Ops at its core is mainly a PostgreSQL database used to store configuration information for :abbr:`ATC (Apache Traffic Control)`, and a set of RESTful API endpoints for interacting with and manipulating that information. It also serves as the single point of authentication for :abbr:`ATC (Apache Traffic Control)` components (that is, when one hears "user" in an :abbr:`ATC (Apache Traffic Control)` context it nearly always means a "user" as configured in Traffic Ops) and provides interfaces to other :abbr:`ATC (Apache Traffic Control)` components by proxy. Additionally, there is some miscellaneous, at times obscure functionality to Traffic Ops, such as generating arbitrary Linux system images.

Software Requirements
=====================
Traffic Ops is only supported on CentOS 7+ systems (although many developers do use Mac OS with some success). Here are the requirements:

- `PostgreSQL 13.2 <https://www.postgresql.org/download/>`_ - the machine where Traffic Ops is running must have the client tool set (e.g. :manpage:`psql(1)`), but the actual database can be run anywhere so long as it is accessible.

	.. note:: Prior to version 13.2, Traffic Ops used version 9.6. For upgrading an existing Mac OS Homebrew-based PostgreSQL instance, you can use `Homebrew <https://brew.sh/>`_ to easily upgrade from 9.6 to 13.2:

		.. code-block:: shell

			brew services stop postgresql
			brew upgrade postgresql
			brew postgresql-upgrade-database
			brew cleanup postgresql@9.6
			brew services start postgresql

- :manpage:`openssl(1SSL)` is recommended to generate server certificates, though not strictly required if certificates can be obtained by other means.
- Some kind of SMTP server is required for certain :ref:`to-api` endpoints to work, but for purposes unrelated to them an SMTP server is not required. :ref:`ciab` comes with a relayless SMTP server for testing (you can view the emails that Traffic Ops sends, but they aren't sent anywhere outside CDN-in-a-Box).

.. tip:: Alternatively, development and testing can be done using :ref:`ciab` - albeit somewhat more slowly.

Perl Script Requirements
------------------------
Not much code is still in Perl, but for the scripts the following are needed:

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


Go Program Requirements
-----------------------
- |install-go-link|_
- If the system's Go compiler doesn't provide it implicitly, also note that all Go code in the :abbr:`ATC (Apache Traffic Control)` repository should be formatted using `gofmt <https://golang.org/cmd/gofmt/>`_

.. |install-go-link| replace:: Go :atc-go-version:`_` or later
.. _install-go-link: http://golang.org/doc/install

All Go code dependencies are managed through the :atc-file:`go.mod`, :atc-file:`go.sum`, and :atc-file:`vendor/modules.txt` files. With the exception of ``golang.org/x`` packages (see :ref:`below <dev-traffic-ops-golang-x>`), module dependencies in :atc-file:`vendor/` are tracked in Git and should thus be available without any extra work - and any new dependencies should be properly "vendored" into that same, top-level directory. No other :atc-file:`vendor/` directories exist, as Go modules only supports a single vendor directory.

.. _dev-traffic-ops-golang-x:

Per the Go language standard's authoritative source's recommendation, all sub-packages of ``golang.org/x`` are treated as a part of the compiler, and so need not ever be "vendored" as though they were an external dependency. These dependencies are not listed explicitly here, so it is strongly advised that they be fetched using :manpage:`go-get(1)` rather than downloaded by hand.

.. tip:: All new dependencies need to be subject to community review to ensure necessity (because it will be added in its entirety to the repository, after all) and license compliance via `the developer mailing list <mailto:dev@trafficcontrol.apache.org>`_.

Traffic Ops Project Tree Overview
=================================
- :atc-file:`traffic_ops/` - The root of the Traffic Ops project

	- app/ - Holds most of the Perl code base, though many of the files contained herein are also used by the Go implementation

		.. note:: This directory is home to many things that no longer work as intended or have been superseded by other things.

		- bin/ - Directory for scripts and tools, :manpage:`cron(8)` jobs, etc.

			- checks/ - Contains the :ref:`to-check-ext` scripts that are provided by default
			- db/ - Contains scripts that manipulate the database beyond the scope of setup, migration, and seeding
			- tests/ - Integration and unit test scripts for automation purposes - in general this has been superseded by :atc-file:`traffic_ops/testing/api/`

		- conf/ - Aggregated configuration for Traffic Ops. For convenience, different environments for the :ref:`database-management` tool are already set up

			- development/ - Configuration files for the "development" environment
			- integration/ - Configuration files for the "integration" environment
			- misc/ - Miscellaneous configuration files.
			- production/ - Configuration files for the "production" environment
			- test/ - Configuration files for the "test" environment

		- db/ - Database setup, seeding, and upgrade/downgrade helpers

			- migrations/ - Database migration files
			- tools/ - Contains helper scripts for easing upgrade transitions when selective data manipulation must be done to achieve a desirable state

		- script/ - Mojolicious bootstrap/startup scripts.
		- templates/ - Mojolicious Embedded Perl (:file:`{template name}.ep`) files for the now-removed Traffic Ops UI

	- build/ - Contains files that are responsible for packaging Traffic Ops into an RPM file - and also for doing the same with :term:`ORT`
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

		- api/ - Integration testing for the Traffic Ops Go client (:atc-godoc:`traffic_ops/v4-client`) and Traffic Ops

	- traffic_ops_golang/ - The root of the Go implementation's code-base

		.. note:: The vast majority of subdirectories of :atc-file:`traffic_ops/traffic_ops_golang/` contain handlers for the :ref:`to-api`, and are named according to the endpoint they handle. What follows is a list of subdirectories of interest that have a special role (i.e. don't handle a :ref:`to-api` endpoint).

		.. seealso:: The GoDoc documentation for this package: :atc-godoc:`/`

		- api/ - A library for use by :ref:`to-api` handlers that provides helpful utilities for common tasks like obtaining a database transaction handle or accessing Traffic Ops configuration
		- auth/ - Contains definitions of privilege levels and access control code used in routing and provides a library for dealing with password and token-based authentication
		- config/ - Defines configuration structures and methods for reading them in from files
		- dbhelpers/ - Assorted utilities that provide functionality for common database tasks, e.g. "Get a user by email"
		- plugin/ - The Traffic Ops plugin system, with examples
		- trafficvault/ - This package provides the Traffic Vault interface and associated backend implementations for other handlers to interact with Traffic Vault.
		- routing/ - Contains logic for mapping all of the :ref:`to-api` endpoints to their handlers, as well as proxying requests back to the Perl implementation and managing plugins, and also provides some wrappers around registered handlers that set common HTTP headers and connection options
		- swaggerdocs/ A currently abandoned attempt at defining the :ref:`to-api` using `Swagger <https://swagger.io/>`_ - it may be picked up again at some point in the (distant) future
		- tenant/ - Contains utilities for dealing with :term:`Tenantable <Tenant>` resources, particularly for checking for permissions
		- tocookie/ - Defines the method of generating the ``mojolicious`` cookie used by Traffic Ops for authentication
		- vendor/ - contains "vendored" Go packages from third party sources

	- v3-client - The official Traffic Ops Go client package for working with the version 3 :ref:`to-api`.
	- v4-client - The official Traffic Ops Go client package for working with the version 4 :ref:`to-api`.
	- v5-client - The official Traffic Ops Go client package for working with the version 5 :ref:`to-api`.
	- vendor/ - contains "vendored" Go packages from third party sources

.. _database-management:

.. program:: admin

app/db/admin
============
The :program:`app/db/admin` binary is for use in managing the Traffic Ops (and Traffic Vault PostgreSQL backend) database tables. This essentially serves as a front-end for `Migrate <https://github.com/golang-migrate/migrate>`_, though  ``dbconf.yml`` comes from `Goose <https://github.com/kevinburke/goose/blob/1.15/db-sample/dbconf.yml>`_, which Traffic Ops used to use before switching to Migrate.

.. note:: For proper resolution of configuration and SOL statement files, it's recommended that this binary be run from the ``app`` directory

Usage
-----
``db/admin [options] command``

Options and Arguments
---------------------
.. option:: --env ENVIRONMENT

	An optional environment specification that causes the database configuration to be read out of the corresponding section of the :atc-file:`app/db/dbconf.yml` configuration file. One of:

	- development
	- integration
	- production
	- test

	:program:`admin` sets :envvar:`MOJO_MODE` to the value of the environment as specified by this option. (Default: ``development``)

.. option:: --trafficvault

	When used, commands will be run against the Traffic Vault PostgreSQL backend database as specified in the :atc-file:`app/db/trafficvault/dbconf.yml` configuration file.

.. option:: command

	The :option:`command` specifies the operation to be performed on the database. It must be one of:

	createdb
		Creates the database for the current environment.
	create_migration
		Creates a pair of timestamped up/down migrations titled NAME.
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
		Sets up the database for the current environment according to the SQL statements in :atc-file:`traffic_ops/app/db/create_tables.sql` or :atc-file:`traffic_ops/app/db/trafficvault/create_tables.sql` if the ``--trafficvault`` option is used
	migrate
		Runs a migration on the database for the current environment
	patch
		Patches the database for the current environment using the SQL statements from the ``app/db/patches.sql``. This command is not supported when using the ``--trafficvault`` option
	redo
		Rolls back the most recently applied migration, then run it again
	reset
		Creates the user defined for the current environment, drops the database for the current environment, creates a new one, loads the schema into it, and runs a single migration on it
	seed
		Executes the SQL statements from the ``app/db/seeds.sql`` file for loading static data. This command is not supported when using the ``--trafficvault`` option. The seed data is constructed under the assumption all migrations for the release have been run, so ``migrate``/\ ``upgrade`` *must* be run first.
	show_users
		Displays a list of all users registered with the PostgreSQL server
	status
		Deprecated, ``status`` is now an alias for ``dbversion`` and will be removed in a future Traffic Control release.
	upgrade
		Performs a migration on the database for the current environment, then patches it using the SQL statements from the :atc-file:`traffic_ops/app/db/patches.sql` file.

.. code-block:: bash
	:caption: Example Usage

	db/admin --env=test reset

The environments are defined in the :atc-file:`traffic_ops/app/db/dbconf.yml` file, and the name of the database generated will be the name of the environment for which it was created. If the ``--trafficvault`` option is used, the :file:`app/db/trafficvault/dbconf.yml` file defines this information.

Resolving Migration Failures
----------------------------

If you encounter an error running a migration, you will see a message like

.. code-block:: bash
	:caption: db/admin error example

	[root@trafficops app]# db/admin -env production migrate
	Error running migrate up: migration failed: syntax error at or near "This_is_a_syntax_error" (column 1) in line 18: /*

That means that the migration timestamp in the ``version`` column of the ``schema_migrations`` table has been updated to the version of the migration that failed, but the ``dirty`` column is also set, and if you try to run another migration (either up or down), you will see

.. code-block:: bash
	:caption: db/admin error migrating when the database version is dirty

	[root@trafficops app]# db/admin -env production migrate
	Error running migrate up: Dirty database version 2021070800000000. Fix and force version.

You will need to manually fix the database. When you are sure that it is fixed, you can unset the ``dirty`` column manually using an SQL client.

Installing The Developer Environment
====================================
To install the Traffic Ops Developer environment:

#. Clone the `Traffic Control repository <https://github.com/apache/trafficcontrol>`_ from GitHub. In most cases it is best to clone this directly into :file:`{GOPATH}/src/github.com/apache/trafficcontrol`, as otherwise the Go implementation will not function properly.

#. Install any required Go dependencies - the suggested method is using :manpage:`go-get(1)`.

	.. code-block:: bash
		:caption: Install Go Development Dependencies

		# assuming current working directory is the repository root
		go mod vendor -v

#. Set up a role (user) in PostgreSQL

	.. seealso:: `PostgreSQL instructions on setting up a database <https://wiki.postgresql.org/wiki/First_steps>`_.


#. Use the ``reset`` and ``upgrade`` :option:`command`\ s of :program:`admin` (see :ref:`database-management` for usage) to set up the ``traffic_ops`` database(s) and optionally with the ``--trafficvault`` option to set up the ``traffic_vault`` database(s).
#. Run the :atc-file:`traffic_ops/install/bin/postinstall` script, it will prompt for information like the default user login credentials.
#. To run Traffic Ops, follow the instructions in :ref:`to-running`.

.. program:: traffic_vault_migrate

app/db/traffic_vault_migrate
==============================
The ``traffic_vault_migrate`` tool - located at :file:`traffic_ops/app/db/traffic_vault_migrate/traffic_vault_migrate.go` in the `Apache Traffic Control repository <https://github.com/apache/trafficcontrol>`_ -
is used to transfer TV keys between database servers. It interfaces directly with each backend so Traffic Ops/Vault being available is not a requirement.
The tool assumes that the schema for each backend is already setup as according to the :ref:`admin setup <traffic_vault_admin>`.

.. program:: traffic_vault_migrate

Usage
-----------
``traffic_vault_migrate [-cdhmr] [-e value] [-f value] [-g value] [-i value] [-l value] [-o value] [-t value]``

.. option:: -c, --compare

		Compare 'to' and 'from' backend keys. Will fetch keys from the dbs of both 'to' and 'from', sorts them by cdn/ds/version and does a deep comparison.

		.. note:: Mutually exclusive with :option:`-r`/:option:`--dry`

.. option:: -d, --dump

		Write keys (from 'from' server) to disk in the folder 'dump' with the unix permissions 0640.

		.. warning:: This can write potentially sensitive information to disk, use with care.

		.. note:: Mutually exclusive with :option:`-l`/:option:`--fill`

.. option:: -e LEVEL, --logLevel=LEVEL

		Print everything at above specified log level (error|warning|info|debug|event) [info]

		.. note:: Mutually exclusive with :option:`-l`/:option:`--logCfg`

.. option:: -f CFG, --fromCfgPath=CFG

		From server config file ["riak.json"]

.. option:: -g CFG, --toCfgPath=CFG

		To server config file ["pg.json"]

.. option:: -h, --help

		Displays usage information

.. option:: -i DIR, --fill DIR

		Insert data into `to` server with data in this directory

		.. note:: Mutually exclusive with :option:`-d`/:option:`--dump`

.. option:: -l CFG, --logCfg CFG

		Log configuration file

		.. note:: Mutually exclusive with :option:`-e`/:option:`--logLevel`

.. option:: -o TYPE, --toType=TYPE

		From server types (Riak|PG) [PG]

.. option:: -m, --noConfirm

		Don't require confirmation before inserting records

.. option:: -r, --dry

		Do not perform writes. Will do a basic output of the keys on the 'from' backend.

		.. note:: Mutually exclusive with :option:`-c`/:option:`--compare`

.. option:: -t TYPE, --fromType=TYPE

		From server types (Riak|PG) [Riak]


Riak
----------

riak.json
""""""""""

 :user: The username used to log into the Riak server.

 :password: The password used to log into the Riak server.

 :host: The hostname for the Riak server.

 :port: The port for which the Riak server is listening for protobuf connections.

 :timeout: The number of seconds to wait for each operation.

 :insecure: (Optional) Determines whether to verify insecure certificates.

 :tlsVersion: (Optional) Max TLS version supported. Valid values are  "10", "11", "12", "13".


Postgres
---------
:program:`traffic_vault_migrate` will properly handle both encryption and decryption of postgres data as that is done on the client side.

pg.json
"""""""""

 :user: The username used to log into the PG server.

 :password: The password for the user to log into the PG server.

 :database: The database to connect to.

 :port: The port on which the PG server is listening.

 :host: The hostname of the PG server.

 :sslmode: The ssl settings for the client connection, `explanation here <https://www.postgresql.org/docs/13/libpq-ssl.html#LIBPQ-SSL-SSLMODE-STATEMENTS>`_. Options are 'disable', 'allow', 'prefer', 'require', 'verify-ca' and 'verify-full'

 :aesKey: The base64 encoding of a 16, 24, or 32 bit AES key.


Logging
----------

The log configuration file has the structure:

 :error_log: Where to output error messages (stderr|stdout|null)

 :warning_log: Where to output warning messages (stderr|stdout|null)

 :info_log: Where to output info messages (stderr|stdout|null)

 :debug_log: Where to output error messages (stderr|stdout|null)

 :event_log: Where to output error messages (stderr|stdout|null)

Adding a Migration Backend
-----------------------------
To add a plugin, implement the traffic_vault_migrate.go:TVBackend interface and add the backend to the returned values in :atc-godoc:`traffic_ops/app/db/traffic_vault_migrate.supportBackends`.

Test Cases
==========

.. _to-go-tests:

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

	Specify the path to the `Test Configuration File`_. If not specified, it will attempt to read a file named ``traffic-ops-test.conf`` in the working directory.

	.. seealso:: `Configuring the Integration Tests`_ for a detailed explanation of the format of this configuration file.

.. _dev-traffic-ops-fixtures:

.. option:: --fixtures FIXTURES

	Specify the path to a file containing static data for the tests to use. This should almost never be used, because many of the tests depend on the data having a certain content and structure. If not specified, it will attempt to read a file named ``tc-fixtures.json`` in the working directory.

.. option:: --includeSystemTests {no|yes}

	Specify whether to run tests that depend on additional components like an SMTP server or a Traffic Vault server. Default: ``no``

Configuring the Integration Tests
"""""""""""""""""""""""""""""""""
Configuration is mainly done through the configuration file passed as :option:`--cfg`, but is also available through the following environment variables.

In addition to the variables described here, the integration tests support identifying the network location of the Traffic Ops instance via :envvar:`TO_URL`.

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

.. envvar:: TO_USER_ADMIN

	If set, will define the name of a user with the "admin" :term:`Role` that will be created by the tests\ [#existinguser]_.

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

Test Configuration File
'''''''''''''''''''''''
The configuration file for the tests (defined by :option:`--cfg`) is a JSON-encoded object with the following properties.

.. warning:: Many of these configuration options are overridden by variables in the execution environment. Where this is a problem, there is an associated warning. In general, this issue is tracked by :issue:`3975`.

:default: An object containing sub-objects relating to default configuration settings for connecting to external resources during testing

	:logLocations: An object containing key/value pairs where the keys are log levels and each associated value is the file location to which logs of that level will be written. The allowed values respect the reserved special names used by the :atc-godoc:`lib/go-log` package. Omitted keys are treated as though their values were ``null``, in which case that level is written to `/dev/null`. The allowed keys are:

		- debug
		- error
		- event
		- info
		- warning

	:session: An object containing key/value pairs that define the default settings used by Traffic Ops "session" connections

		:timeoutInSecs: At the time of this writing this is the only meaningful configuration option that may be present under ``session``. It specifies the timeouts used by client connections during testing as an integer number of seconds. The default if not specified (or overridden) is 0, meaning no limit.

			.. warning:: This configuration is overridden by :envvar:`SESSION_TIMEOUT_IN_SECS`.

:trafficOps: An object containing information that defines the running Traffic Ops instance to use in testing.

	:password: This password will be used for all created users used by the test suite - it does not need to be the password of any pre-existing user. The default if not specified (or overridden) is an empty string, which may or may not cause problems.

		.. warning:: This is overridden by :envvar:`TO_USER_PASSWORD`.

	:URL: The network location of the running Traffic Ops server, including schema, hostname and optionally port number e.g. ``https://localhost:6443``.

		.. warning:: This is overridden by :envvar:`TO_URL`.


	:users: An object containing key-value pairs where the keys are the names of :term:`Roles` and the values are the usernames of users that will be created with the associated :term:`Role` for testing purposes. *There are very few good reasons why the values should not just be the same as the keys*. The default for any missing (and not overridden) key is the empty string which is *won't* work so please don't leave any undefined. The allowed keys are:

		- admin

			.. warning:: The value of this key is overridden by :envvar:`TO_USER_ADMIN`.

		- disallowed

			.. warning:: The value of this key is overridden by :envvar:`TO_USER_DISALLOWED`.

		- extension

			.. warning:: The value of this key is overridden by :envvar:`TO_USER_EXTENSION`.

		- federation

			.. warning:: The value of this key is overridden by :envvar:`TO_USER_FEDERATION`.

		- operations

			.. warning:: The value of this key is overridden by :envvar:`TO_USER_OPERATIONS`.

		- portal

			.. warning:: The value of this key is overridden by :envvar:`TO_USER_PORTAL`.

		- readOnly

			.. warning:: The value of this key is overridden by :envvar:`TO_USER_READ_ONLY`.

:trafficOpsDB: An object containing information that defines the database to use in testing\ [#integrationdb]_.

	:dbname: The name of the database to which the tests will connect\ [#integrationdb]_.

		.. warning:: This is overridden by :envvar:`TODB_NAME`.

	:description: An utterly cosmetic option that need not exist at all which, if set, gives a description of the database to which the tests will connect. This has no effect except possibly changing one line of debug output.

		.. warning:: This is overridden by :envvar:`TODB_DESCRIPTION`

	:hostname: The :abbr:`FQDN (Fully Qualified Domain Name)` of the server on which the database is running\ [#integrationdb]_

		.. warning:: This is overridden by :envvar:`TODB_HOSTNAME`.

	:password: The password to use when authenticating with the database

		.. warning:: This is overridden by :envvar:`TODB_PASSWORD`.

	:port: The port on which the database listens for connections\ [#integrationdb]_ - as a **string**

		.. warning:: This is overridden by :envvar:`TODB_PORT`.

	:type: The "type" of database being used\ [#integrationdb]_. This should **never** be set to anything besides ``"Pg"``, anything else results in undefined behavior (although it's equally possible that it simply won't have any effect).

		.. warning:: This is overridden by :envvar:`TODB_TYPE`.

	:ssl: An optional boolean value that defines whether or not the database uses SSL encryption for its connections - default if not specified (or overridden) is ``false``

		.. warning:: This is overridden by :envvar:`TODB_SSL`.

	:user: The name of the user as whom to authenticate with the database

		.. warning:: This is overridden by :envvar:`TODB_USER`.

Writing New Endpoints
=====================
.. note:: Most new endpoints are accompanied by database schema changes which necessitate a new migration under :atc-file:`traffic_ops/app/db/migrations` and database best-practices are not discussed in this section.

.. seealso:: This section contains a quick overview of API endpoint development; for the full guidelines for API endpoints, consult :ref:`api-guidelines`.

The first thing to consider when writing a new endpoint is what the requests it will serve will look like. It's recommended that new endpoints avoid using "path parameters" when possible, and instead try to utilize request bodies and/or query string parameters. For example, instead of ``/foos/{{ID}}`` consider simply ``/foos`` with a supported ``id`` query parameter. The request *methods* should be restricted to the following, and respect each method's associated meaning.

DELETE
	Removes a resource or one or more of its representations from the server. This should **always** be the method used when deleting objects.
GET
	Retrieves a representation of some resource. This should *always* be used for read-only operations and note that the requesting client **never** expects the state of the server to change as a result of a request using the GET method.
POST
	Requests that the server process some passed data. This is used most commonly to create new objects on the server, but can also be used more generally e.g. with a request for regenerating encryption keys. Although this isn't strictly creating new API resources, it does change the state of the server and so this is more appropriate than GET.
PUT
	Places a new representation of some resource on the server. This is typically used for updating existing objects. For creating *new* representations/objects, use POST instead. When using PUT note that clients expect it to be :dfn:`idempotent`, meaning that subsequent identical PUT requests should result in the same server state. What this means is that it's standard to require that *all* of the information defining a resource be provided for each request even if the vast majority of it isn't changing.

The HEAD and OPTIONS request methods have default implementations for any properly defined :ref:`to-api` route, and so should almost never be defined explicitly. Other request methods (e.g. CONNECT) are currently unused and ought to stay that way for the time being.

.. note:: Utilizing the PATCH method is unfeasible at the time of this writing but progress toward supporting it is being made, albeit slowly in the face of other priorities.

.. seealso:: The :abbr:`MDN (Mozilla Developer Network)`'s `documentation on the various HTTP request methods <https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods>`_.

The final step of creating any :ref:`to-api` endpoint is to write documentation for it. When doing so, be sure to follow *all* of the guidelines laid out in :ref:`docs-guide`. *If documentation doesn't exist for new functionality then it has accomplished* **nothing** *because no one using Traffic Control will know it exists*. Omitted documentation is how a project winds up with a dozen different API endpoints that all do essentially the same thing.

Framework Options
-----------------
The Traffic Ops code base offers two basic frameworks for defining a new endpoint. Either one may be used at the author's discretion (or even neither if desired and appropriate - though that seems unlikely).

Generic "CRUDer"
""""""""""""""""
The "Generic 'CRUDer'", as it's known, is a pattern of API endpoint development that principally involves defining a ``type`` that implements the :to-godoc:`api.CRUDer` interface. A description of what that entails is best left to the actual GoDoc documentation.

.. seealso:: The :to-godoc:`api.GenericCreate`, :to-godoc:`api.GenericDelete`, :to-godoc:`api.GenericRead`, and :to-godoc:`api.GenericUpdate` helpers are often used to provide the default operations of creating, deleting, reading, and updating objects, respectively. When the API endpoint being written is only meant to perform these basic operations on an object or objects stored in the database, these should be totally sufficient.

This method offers a lot of functionality "out-of-the-box" as compared to the `API Info`_ method, but because of that is also restrictive. For example, it is not possible to write an endpoint that returns data not encoded as JSON using this method. That's an uncommon use-case, but not unheard-of.

This method is best used for basic creation, reading, update, and deletion operations performed on simple objects with no structural differences across API versions.

API Info
""""""""
Endpoint handlers can also be defined by simply implementing the :godoc:`net/http.HandlerFunc` interface. The :godoc:`net/http.Request` reference passed into such handlers provides identifying information for the authenticated user (where applicable) in its context.

To easily obtain the information needed to identify a user and their associated permissions, as well as server configuration information and a database transaction handle, authors should use the :to-godoc:`api.NewInfo` function which will return all of that information in a single structure as well as any errors encountered during the process and an appropriate HTTP response code in case of such errors.

This method offers fine control over the endpoint's logic, but tends to be much more verbose than the endpoints written using the `Generic "CRUDer"`_ method. For example, a handler for retrieving an object from the database and returning it to the requesting client encoded as JSON can be twenty or more lines of code, whereas a single call to :to-godoc:`api.GenericCreate` provides equivalent functionality.

This method is best used when requests are meant to have extensive side-effects, are performed on unusually structured objects, need fine control of the HTTP headers/options, or operate on objects that have different structures or meanings across API versions.

Extensions
==========
What's typically meant by "extension" in the context of Traffic Ops is a :ref:`to-check-ext` which provides data for server "checks" which can be viewed in Traffic Portal under :menuselection:`Monitor --> Cache Checks`. This type of extension need not know nor even care which implementation it is being used with, as it interacts with Traffic Ops through the :ref:`to-api`.

Traffic Ops supports overrides or new definitions for non-standard :ref:`to-api` routes. This type of "extension" is typically reffered to as a "plugin," and they are described in `Go Plugins`_.

.. _to_go_plugins:

Go Plugins
----------
A plugin is defined by a Go source file in the :atc-file:`traffic_ops/traffic_ops_golang/plugin` directory, which is expected to be named :file:`{plugin name}.go`. A plugin is registered to Traffic Ops by a call to :to-godoc:`plugin.AddPlugin` in the source file's special ``init`` function.

A plugin is only enabled at runtime if its name is present in the :ref:`cdn.conf` file's ``traffic_ops_golang.plugins`` array.

Each plugin may also define any, all, or none of the lifecycle hooks provided: ``load``, ``startup``, and ``onRequest``

load
	The ``load`` function of a plugin, if defined, needs to implement the :to-godoc:`plugin.LoadFunc` interface, and will be run when the server starts and after configuration has been loaded. It will be passed the plugins own configuration as it was defined in the :ref:`cdn.conf` file's ``traffic_ops_golang.plugin_config`` map.
onRequest
	The ``onRequest`` function of a plugin, if defined, needs to implement the :to-godoc:`plugin.OnRequestFunc` interface, and will be called on **every** request made to the :ref:`to-api`. Because of this, it's imperative that the function exit as soon as possible. Note that once one plugin reports that it has served the request, no others will be tried. The order in which plugins are tried is defined by their order in the ``traffic_ops_golang.plugins`` array of the :ref:`cdn.conf` configuration file.

		.. seealso:: It's very common for this function to behave like a :ref:`to-api` endpoint, so when writing a plugin it may be useful to review `Writing New Endpoints`_.
startup
	Like ``load``, the ``startup`` function of a plugin, if defined, will be called when the server starts and after configuration has been loaded. *Unlike* ``load``, however, this function should implement the :to-godoc:`plugin.StartupFunc` interface and will be passed in the entirety of the server's configuration, including its own configuration and any shared plugin configuration data as defined in the :ref:`cdn.conf` file's ``traffic_ops_golang.plugin_shared_config`` map.

Example
"""""""
An example "Hello World" plugin that serves the ``/_hello`` request path by just writing "Hello World" in the body of a 200 OK response back to the client is provided in :atc-file:`traffic_ops/traffic_ops_golang/plugin/hello_world.go`:

.. literalinclude:: ../../../traffic_ops/traffic_ops_golang/plugin/hello_world.go
	:language: go
	:linenos:
	:tab-width: 4

Check Extensions
----------------
:ref:`to-check-ext` allow you to add custom checks to the :menuselection:`Monitor --> Cache Checks` view.

Extensions are managed using the ``$TO_HOME/bin/extensions`` command line script

.. seealso:: For more information see :ref:`admin-to-ext-script`.

.. [#integrationdb] The Traffic Ops instance *must* be using the same PostgreSQL database that the tests will use.
.. [#existinguser] This does not need to match the name of any pre-existing user.
