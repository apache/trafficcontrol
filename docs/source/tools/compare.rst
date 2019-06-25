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

.. _compare-tool:

****************
The Compare Tool
****************
The ``compare`` tool is used to compare the output of a set of :ref:`to-api` endpoints between two running instances of Traffic Ops. The idea is that two different versions of Traffic Ops with the same data will have differences in the output of their API endpoints *if and only if* either the change was intentional, or a new bug was introduced in the newer version. Typically, this isn't really true, due to rapidly changing data structures like timestamps in the API outputs, but this should offer a good starting point for identifying bugs in changes made to the :ref:`to-api`.

Location and Dependencies
=========================
The ``compare`` tool is written in Go, and can be found from within the Traffic Control repository at ``traffic_ops/testing/compare/``. The main file of interest is ``compare.go``, which contains the main routine and logic for checking endpoints. To build the executable, run ``go build .`` from within the ``traffic_ops/testing/compare/`` directory. Alternatively, run the file without storing a built binary by using ``go run <path to compare.go>``. In order to build/run the ``compare`` tool, the following dependencies should be satisfied, probably via ``go get``:

* github.com/apache/trafficcontrol/lib/go-tc\ [1]_
* github.com/kelseyhightower/envconfig
* golang.org/x/net/publicsuffix

The directory containing the ``compare`` tool also contains an executable Python 3 script named ``genConfigRoutes.py``. This script can be used to scrape the two Traffic Ops instances for API routes that resolve to generated configuration files for mid-tier and edge-tier :term:`cache server`\ s, which can then be fed directly into the ``compare`` tool via a file or pipe. While the script itself has no actual dependencies, it *must* be run from within the full Traffic Control repository, as it imports the Python client for Traffic Ops (located in ``traffic_control/clients/python/trafficops`` inside the repository). The client itself has its own documented dependencies

.. [1] Theoretically, if you downloaded the Traffic Control repository properly (into ``$GOPATH/src/github.com/apache/trafficcontrol``), this will already be satisfied.

Usage
=====

.. program:: compare

traffic_ops/testing/compare/compare.go
--------------------------------------
``compare [-hsV] [-f FILE] [--ref_passwd PASSWD] [--ref_url URL] [--ref_user USER] [-r PATH] [--test_passwd PASSWD] [--test_url URL] [--test_user USER] [parameters ...]``

.. option:: --ref_passwd PASSWD

	The password for logging into the reference Traffic Ops instance. This option overrides the :envvar:`TO_PASSWORD` environment variable, and is required if and only if :envvar:`TO_PASSWORD` is not set.

.. option:: --ref_url URL

	The URL that points to the reference Traffic Ops instance. This option overrides the :envvar:`TO_URL` environment variable, and is required if and only if :envvar:`TO_URL` is not set.

.. option:: --ref_user USER

	The username for logging into the reference Traffic Ops instance. This option overrides the :envvar:`TO_USER` environment variable, and is required if and only if :envvar:`TO_USER` is not set.

.. option:: --test_passwd PASSWD

	The password for logging into the testing Traffic Ops instance. This option overrides the :envvar:`TEST_PASSWORD` environment variable. Additionally, if this option is not specified *and* :envvar:`TEST_PASSWORD` is not set, the value for :envvar:`TO_PASSWORD` (or :option:`--ref_passwd` if overridden) will be used.

.. option:: --test_url URL

	The URL for the testing Traffic Ops instance. This option overrides the :envvar:`TEST_URL` environment variable. Additionally, if this option is not specified *and* :envvar:`TEST_URL` is not set, the value for :envvar:`TO_URL` (or :option:`--ref_url` if overridden) will be used.

.. option:: --test_user USER

	The username for logging into the testing Traffic Ops instance. This option overrides the :envvar:`TEST_USER` environment variable. Additionally, if this option is not specified *and* :envvar:`TEST_USER` is not set, the value for :envvar:`TO_USER` (or :option:`--ref_user` if overridden) will be used.

.. option:: -f FILE, --file FILE

	This optional flag specifies a file from which to list API paths to test. If this option is not given, :program:`compare` will read from STDIN.

.. option:: -h, --help

	Print usage information and exit

.. option:: -r PATH, --results_path PATH

	This optional flag specifies an output directory where results will be written. Default: ``./results``

.. option:: -V, --version

	Print version information and exit

.. versionchanged:: 3.0.0
	Removed the ``-s`` command line switch to compare CDN :term:`Snapshots` - this is now the responsibility of the :program:`genConfigRoutes.py` script.

.. program:: genConfigRoutes.py

traffic_ops/testing/compare/genConfigRoutes.py
----------------------------------------------
.. note:: This script uses the :ref:`py-client`, and so that must be installed to use it.

``genConfigRoutes.py [-h] [-v] [--refURL URL] [--testURL URL] [--refUser USER] [--refPasswd PASSWD] [--testUser USER] [--testPasswd PASSWD] [-k] [-l LOG_LEVEL] [-q]``

A simple script to generate API routes to server configuration files for a given pair of Traffic Ops instances. This, for the purpose of using the :program:`compare` tool

.. option:: -h, --help

	Show usage information and exit

.. option:: --refURL URL

	The full URL of the reference Traffic Ops instance. This option overrides the :envvar:`TO_URL` environment variable, and is required if and only if :envvar:`TO_URL` is not set.

.. option:: --testURL URL

	The full URL of the testing Traffic Ops instance. This option overrides the :envvar:`TEST_URL` environment variable. Additionally, if this option is not specified *and* :envvar:`TEST_URL` is not set, the value for :envvar:`TO_URL` (or :option:`--refURL` if overridden) will be used.

.. option:: --refUser USER

	A username for logging into the reference Traffic Ops instance. This option overrides the :envvar:`TO_USER` environment variable, and is required if and only if :envvar:`TO_USER` is not set.

.. option:: --refPasswd PASSWD

	A password for logging into the reference Traffic Ops instance. This option overrides the :envvar:`TO_PASSWORD` environment variable, and is required if and only if :envvar:`TO_PASSWORD` is not set.

.. option:: --testUser USER

	A username for logging into the testing Traffic Ops instance. This option overrides the :envvar:`TEST_USER` environment variable. Additionally, if this option is not specified *and* :envvar:`TEST_USER` is not set, the value for :envvar:`TO_USER` (or :option:`--refUser` if overridden) will be used.

.. option:: --testPasswd PASSWD

	A password for logging into the testing Traffic Ops instance. This option overrides the :envvar:`TEST_PASSWORD` environment variable. Additionally, if this option is not specified *and* :envvar:`TEST_PASSWORD` is not set, the value for :envvar:`TO_PASSWORD` (or :option:`--refPasswd` if overridden) will be used.

.. option:: -k, --insecure

	Do not verify SSL certificate signatures against *either* Traffic Ops instance (default: False)

.. option:: -v, --version

	Print version information and exit

.. option:: -l LOG_LEVEL, --log_level LOG_LEVEL

	Sets the Python log level, one of

	- DEBUG
	- INFO
	- WARN
	- ERROR
	- CRITICAL

	(default: INFO)

..option:: -q, --quiet

	Suppresses all logging output - even for critical errors (default: False)

.. option:: -s, --snapshot

	Produce CDN :term:`Snapshot` routes in the output (CRConfig.json, snapshot/new etc.) (default: False)

.. option:: -C, --no-server-configs

	Do not generate routes for server configuration files (default: False)

.. tip:: If you're using a CDN-in-a-Box environment for testing, it's likely that you'll need the :option:`-k`/:option:`--insecure` option if you're outside the Docker network

Environment Variables
---------------------
Both :program:`compare` and :program:`genConfigRoutes.py` require connection and authentication methods for two Traffic Ops instances. For ease of use, these can be provided by environment variables. Both programs are capable of using the same environment variables, so that they only need to be defined once each.

.. envvar:: TO_URL

	The URL of the reference Traffic Ops instance. Overridden by :option:`genConfigRoutes.py --refURL` and :option:`compare --ref_url`.

.. envvar:: TO_USER

	The username to authenticate with the reference Traffic Ops instance. Overridden by :option:`genConfigRoutes.py --refUser` and :option:`compare --ref_user`.

.. envvar:: TO_PASSWORD

	The password to authenticate with the reference Traffic Ops instance. Overridden by :option:`genConfigRoutes.py --refPasswd` and :option:`compare --ref_passwd`.

.. envvar:: TEST_URL

	The URL of the testing Traffic Ops instance. Overridden by :option:`genConfigRoutes.py --testURL` and :option:`compare --test_url`.

.. envvar:: TEST_USER

	The username to authenticate with the testing Traffic Ops instance. Overridden by :option:`genConfigRoutes.py --testUser` and :option:`compare --test_user`.

.. envvar:: TEST_PASSWORD

	The password to authenticate with the testing Traffic Ops instance. Overridden by :option:`genConfigRoutes.py --testPasswd` and :option:`compare --test_passwd`.

Usage in a Pipeline
===================
The :program:`genConfigRoutes.py` script will output list of unique API routes (relative to the desired Traffic Ops URL) that point to generated configuration files for a sample set of servers common to both  Traffic Ops instances. The results are printed to STDOUT, making the output perfect for piping directly into :program:`compare` like so:

.. code-block:: shell
	:caption: Example Pipeline from :program:`genConfigRoutes.py` into :program:`compare`

	./genConfigRoutes.py https://trafficopsA.example.test https://trafficopsB.example.test username:password | ./compare

.. note:: This is assuming the proper `Environment Variables`_ have been set for :program:`compare`.

Usage with Docker
=================
A Dockerfile is provided to run tests on a pair of instances given the configuration environment variables necessary. This will generate configuration file routes using :program:`genConfigRoutes.py`, and add them to whatever is already contained in :file:`traffic_ops/testing/compare/testroutes.txt`, then run the :program:`compare` tool on the final API route list. Build artifacts (i.e. anything output files created by the :program:`compare` tool) are placed in the :file:`/artifacts/` directory on the container. To retrieve these results, the use of a volume is recommended. The build context *must* be at the root of the Traffic Control repository, as the tools have dependencies on the Traffic Control clients.

Arguments can be passed to the :program:`genConfigRoutes.py` script by defining the build-time argument ``MODE``. By default it expands to :option:`-s` to allow the generation of CDN :term:`Snapshot` routes. It is not necessary to pass :option:`-k`/:option:`--insecure`, as the Dockerfile will do that implicitly.

In order to use the container, the following environment variables must be defined for the container at runtime:

TO_URL
	The URL of the reference Traffic Ops instance
TO_USER
	The username to authenticate with the reference Traffic Ops instance
TO_PASSWORD
	The password to authenticate with the reference Traffic Ops instance
TEST_URL
	The URL of the testing Traffic Ops instance
TEST_USER
	The username to authenticate with the testing Traffic Ops instance
TEST_PASSWORD
	The password to authenticate with the testing Traffic Ops instance

.. code-block:: shell
	:caption: Sample Script to Build and Run

	sudo docker build . -f traffic_ops/testing/compare/Dockerfile -t compare:latest
	sudo docker run -v $PWD/artifacts:/artifacts -e TO_URL="$TO_URL" -e TEST_URL="$TEST_URL" -e TO_USER="admin" -e TO_PASSWORD="twelve" -e TEST_USER="admin" -e TEST_PASSWORD="twelve" compare:latest

.. note:: The above code example assumes that the environment variables ``TO_URL`` and ``TEST_URL`` refer to the URL of the reference Traffic Ops instance and the URL of the test Traffic Ops instance, respectively (including port numbers). It also uses credentials suitable for logging into a stock :ref:`ciab` instance.

.. note:: Unlike using the :program:`genConfigRoutes.py` script and/or the :program:`compare` on their own, *all* of the variables must be defined, even if they are duplicates.
