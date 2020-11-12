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
At its current stage in development, "Traffic Ops" actually refers to a concept with two implementations. The original Traffic Ops was written as a collection of Perl scripts based on the `Mojolicious framework <https://mojolicious.org/>`_ framework. At some point, the relatively poor performance and lack of knowledgeable developers as the project expanded became serious issues, and so for the past few years Traffic Ops has undergone a rewrite to Go.

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
- Some kind of SMTP server is required for certain :ref:`to-api` endpoints to work, but for purposes unrelated to them an SMTP server is not required. :ref:`ciab` comes with a relayless SMTP server for testing (you can view the emails that Traffic Ops sends, but they aren't sent anywhere outside CDN-in-a-Box).

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
- |install-go-link|_
- If the system's Go compiler doesn't provide it implicitly, also note that all Go code in the :abbr:`ATC (Apache Traffic Control)` repository should be formatted using `gofmt <https://golang.org/cmd/gofmt/>`_

.. |install-go-link| replace:: Go :atc-go-version:`_` or later
.. _install-go-link: http://golang.org/doc/install

All Go code dependencies are managed through the :atc-file:`vendor/` directory and should thus be available without any extra work - and any new dependencies should be properly "vendored" into that same, top-level directory. Some dependencies have been "vendored" into :atc-file:`traffic_ops/vendor` and :atc-file:`traffic_ops/traffic_ops_golang/vendor` but the preferred location for new dependencies is under that top-level :atc-file:`vendor/` directory.

Per the Go language standard's authoritative source's recommendation, all sub-packages of ``golang.org/x`` are treated as a part of the compiler, and so need not ever be "vendored" as though they were an external dependency. These dependencies are not listed explicitly here, so it is strongly advised that they be fetched using :manpage:`go-get(1)` rather than downloaded by hand.

.. tip:: All new dependencies need to be subject to community review to ensure necessity (because it will be added in its entirety to the repository, after all) and license compliance via `the developer mailing list <mailto:dev@trafficcontrol.apache.org>`_.

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

	:program:`admin` sets :envvar:`MOJO_MODE` to the value of the environment as specified by this option. (Default: ``development``)

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

#. Clone the `Traffic Control repository <https://github.com/apache/trafficcontrol>`_ from GitHub. In most cases it is best to clone this directly into :file:`{GOPATH}/src/github.com/apache/trafficcontrol`, as otherwise the Go implementation will not function properly.
#. Install the local dependencies using `Carton <https://metacpan.org/release/Carton>`_.

	.. code-block:: bash
		:caption: Install Perl Dependencies

		# assuming current working directory is the repository root
		cd traffic_ops/app
		carton

#. Install any required Go dependencies - the suggested method is using :manpage:`go-get(1)`.

	.. code-block:: bash
		:caption: Install Go Development Dependencies

		# assuming current working directory is the repository root
		go get -v ./lib/... ./traffic_ops/traffic_ops_golang/...

#. Set up a role (user) in PostgreSQL

	.. seealso:: `PostgreSQL instructions on setting up a database <https://wiki.postgresql.org/wiki/First_steps>`_.


#. Use the ``reset`` and ``upgrade`` :option:`command`\ s of :program:`admin` (see :ref:`database-management` for usage) to set up the ``traffic_ops`` database(s).
#. Run the :atc-file:`traffic_ops/install/bin/postinstall` script, it will prompt for information like the default user login credentials
#. To run Traffic Ops, follow the instructions in :ref:`to-running`.

Test Cases
==========

Perl Tests
----------
Use `prove <http://perldoc.perl.org/prove.html>`_ (should be installed with Perl) to execute test cases\ [#perltest]_. Execute after a ``carton install`` of all required dependencies:

- To run the Unit Tests: ``prove -qrp  app/t/``
- To run the Integration Tests: ``prove -qrp app/t_integration/``

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
.. option:: --fixtures FIXTURES

	Specify the path to a file containing static data for the tests to use. This should almost never be used, because many of the tests depend on the data having a certain content and structure. If not specified, it will attempt to read a file named ``tc-fixtures.json`` in the working directory.

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

	:logLocations: An object containing key/value pairs where the keys are log levels and each associated value is the file location to which logs of that level will be written. The allowed values respect the `reserved special names used by the github.com/apache/trafficcontrol/lib/go-log package <https://godoc.org/github.com/apache/trafficcontrol/lib/go-log#pkg-constants>`_. Omitted keys are treated as though their values were ``null``, in which case that level is written to `/dev/null`. The allowed keys are:

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
All new :ref:`to-api` endpoints should be written in Go, so writing endpoints for the Perl implementation is not discussed here. Furthermore, most new endpoints are accompanied by database schema changes which necessitate a new migration under :atc-file:`traffic_ops/app/db/migrations` and database best-practices are not discussed in this section.

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

This method offers a lot of functionality "out-of-the-box" as compared to the `APIInfo`_ method, but because of that is also restrictive. For example, it is not possible to write an endpoint that returns data not encoded as JSON using this method. That's an uncommon use-case, but not unheard-of.

This method is best used for basic creation, reading, update, and deletion operations performed on simple objects with no structural differences across API versions.

APIInfo
"""""""
Endpoint handlers can also be defined by simply implementing the :godoc:`net/http.HandlerFunc` interface. The :godoc:`net/http.Request` reference passed into such handlers provides identifying information for the authenticated user (where applicable) in its context.

To easily obtain the information needed to identify a user and their associated permissions, as well as server configuration information and a database transaction handle, authors should use the :to-godoc:`api.NewInfo` function which will return all of that information in a single structure as well as any errors encountered during the process and an appropriate HTTP response code in case of such errors.

This method offers fine control over the endpoint's logic, but tends to be much more verbose than the endpoints written using the `Generic "CRUDer"`_ method. For example, a handler for retrieving an object from the database and returning it to the requesting client encoded as JSON can be twenty or more lines of code, whereas a single call to :to-godoc:`api.GenericCreate` provides equivalent functionality.

This method is best used when requests are meant to have extensive side-effects, are performed on unusually structured objects, need fine control of the HTTP headers/options, or operate on objects that have different structures or meanings across API versions.

Rewriting a Perl Endpoint
-------------------------
When rewriting endpoints from Perl, some special considerations must be taken.

- Any rules and guidelines herein outlined that are broken by the Perl handler *must* also be broken in the rewritten Go handler to maintain compatibility within the API. New features can be added in the latest unreleased version of the API so long as they are appropriately documented, but avoid the temptation to fix things that seem broken. Such changes are best left to re-implementation of the API in a subsequent major version. The exceptions to this rule are if the broken behavior constitutes a security vulnerability (in which case be sure to follow the instructions on `the Apache Software Foundation security page <https://www.apache.org/security/>`_) or if it happens in the event of a server or client error. For example, many Perl handlers will spit out an HTML page in the event of a server-side error while the standard behavior of the :ref:`to-api` in such cases is to return the appropriate HTTP response code and a response body containing a JSON-encoded ``alerts`` object describing the nature of the error.
- Mark newly rewritten endpoints in their :atc-file:`traffic_ops/traffic_ops_golang/routing/routes.go` definition with ``perlBypass`` to ensure that upon upgrading it is possible to configure the server to fall back on the Perl implementation. That way, any erroneous rewrites that wind up in production environments can be quickly bypassed in favor of the old, known-to-be-working version.
- The Perl handlers support any combination of optional trailing ``/`` and ``.json`` on endpoint routes, and rewritten route definitions ought to support that. For example, the endpoint ``/foo`` can with equal validity from the Perl implementation's perspective as ``/foo.json``, ``/foo/``, ``/foo.json/`` (for some reason), and even (horrendously) as ``/foo/.json``.
- It's possible that a route definition for the newly rewritten route already exists, explicitly defining a proxy to the Perl implementation using ``handlerToFunc(proxyHandler)`` to avoid collisions with later-defined routes. These will need to be deleted in order for the route to be properly handled.

Extensions
==========
Both the Perl and Go implementation support different kinds of extensions.

What's typically meant by "extension" in the context of Traffic Ops is a :ref:`to-check-ext` which provides data for server "checks" which can be viewed in Traffic Portal under :menuselection:`Monitor --> Cache Checks`. This type of extension need not know nor even care which implementation it is being used with, as it interacts with Traffic Ops through the :ref:`to-api`. These are described in `Legacy Perl Extensions`_ as their description remains rather Perl-centric, but in principle their operation is not limited to the context of the Perl Implementation.

Both Perl and Go also support overrides or new definitions for non-standard :ref:`to-api` routes. It is strongly recommended that no Perl-based extensions of this type be written, but for posterity they are described in `Legacy Perl Extensions`_. The Go implementation refers to this type of "extension" as a "plugin," and they are described in `Go Plugins`_.

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

Legacy Perl Extensions
----------------------
Traffic Ops Extensions are a way to enhance the basic functionality of Traffic Ops in a customizable manner. There are two types of extensions:

:ref:`to-check-ext`
	These allow you to add custom checks to the :menuselection:`Monitor --> Cache Checks` view.

:ref:`to-datasource-ext`
	These allow you to add statistic sources for the graph views and APIs.

Extensions are managed using the ``$TO_HOME/bin/extensions`` command line script

.. seealso:: For more information see :ref:`admin-to-ext-script`.


Extensions at Runtime
"""""""""""""""""""""
The search path for :ref:`to-datasource-ext` depends on the configuration of the ``PERL5LIB`` environment variable, which is pre-configured in the Traffic Ops start scripts. All :ref:`to-check-ext` must be located in ``$TO_HOME/bin/checks``

	.. code-block:: bash
		:caption: Example ``PERL5LIB`` Configuration

		export PERL5LIB=/opt/traffic_ops_extensions/private/lib/Extensions:/opt/traffic_ops/app/lib/Extensions/TrafficStats

To prevent :ref:`to-datasource-ext` namespace collisions within Traffic Ops all :ref:`to-datasource-ext` should follow the package naming convention '``Extensions::<ExtensionName>``'

``TrafficOpsRoutes.pm``
"""""""""""""""""""""""
Traffic Ops accesses each extension through the addition of a URL route as a custom hook. These routes will be defined in a file called ``TrafficOpsRoutes.pm`` that should be present in the top directory of your Extension. The routes that are defined should follow the `Mojolicious route conventions <https://mojolicious.org/perldoc/Mojolicious/Guides/Routing#Routes>`_.


Development Configuration
"""""""""""""""""""""""""
To incorporate any custom :ref:`to-datasource-ext` during development set your ``PERL5LIB`` environment variable with any number of colon-separated directories with the understanding that the ``PERL5LIB`` search order is from left to right through this list. Once Perl locates your custom route or Perl package/class it 'pins' on that class or Mojolicious Route and doesn't look any further, which allows for the developer to override Traffic Ops functionality.

.. [#perltest] As progress continues on moving Traffic Ops to run entirely in Go, the number of passing tests has steadily decreased. This means that the tests are not a reliable way to test Traffic Ops, as they are expected to fail more and more as functionality is stripped from the Perl codebase.
.. [#integrationdb] The Traffic Ops instance *must* be using the same PostgreSQL database that the tests will use.
.. [#existinguser] This does not need to match the name of any pre-existing user.
