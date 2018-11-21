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
The ``compare`` tool is used to compare the output of a set of :ref:`Traffic Ops API <to-api>` endpoints between two running instances of Traffic Ops. The idea is that two different versions of Traffic Ops with the same data will have differences in the output of their API endpoints *if and only if* either the change was intentional, or a new bug was introduced in the newer version. Typically, this isn't really true, due to rapidly changing data structures like timestamps in the API outputs, but this should offer a good starting point for identifying bugs in changes made to the Traffic Ops API.

Location and Dependencies
=========================
The ``compare`` tool is written in Go, and can be found from within the Traffic Control repository at ``traffic_ops/testing/compare/``. The main file of interest is ``compare.go``, which contains the main routine and logic for checking endpoints. To build the executable, run ``go build .`` from within the ``traffic_ops/testing/compare/`` directory. Alternatively, run the file without storing a built binary by using ``go run <path to compare.go>``. In order to build/run the ``compare`` tool, the following dependencies should be satisfied, probably via ``go get``:

* github.com/apache/trafficcontrol/lib/go-tc\ [1]_
* github.com/kelseyhightower/envconfig
* golang.org/x/net/publicsuffix

The directory containing the ``compare`` tool also contains an executable Python 3 script named ``genConfigRoutes.py``. This script can be used to scrape the two Traffic Ops instances for API routes that resolve to generated configuration files for mid-tier and edge-tier cache servers, which can then be fed directly into the ``compare`` tool via a file or pipe. While the script itself has no actual dependencies, it *must* be run from within the full Traffic Control repository, as it imports the Python client for Traffic Ops (located in ``traffic_control/clients/python/trafficops`` inside the repository). The client itself has its own documented dependencies

.. TODO: ^ make that last statement not a dirty lie ^

.. [1] Theoretically, if you downloaded the Traffic Control repository properly (into ``$GOPATH/src/github.com/apache/trafficcontrol``), this will already be satisfied.

Usage
=====

``compare``
-----------
Usage: compare [-hsV] [-f value] [--ref_passwd value] [--ref_url value] [--ref_user value] [-r value] [--test_passwd value] [--test_url value] [--test_user value] [parameters ...]

--ref_passwd=value        The password for logging into the reference Traffic Ops instance (overrides TO_PASSWORD environment variable)
--ref_url=value           The URL for the reference Traffic Ops instance (overrides TO_URL environment variable)
--ref_user=value          The username for logging into the reference Traffic Ops instance (overrides TO_USER environment variable)
--test_passwd=value       The password for logging into the testing Traffic Ops instance (overrides TEST_PASSWORD environment variable)
--test_url=value          The URL for the testing Traffic Ops instance (overrides TEST_URL environment variable)
--test_user=value         The username for logging into the testing Traffic Ops instance (overrides TEST_USER environment variable)
-f, --file=value          File listing routes to test (will read from stdin if not given)
-h, --help                Print usage information and exit
-r, --results_path=value  Directory where results will be written
-V, --version             Print version information and exit

.. versionchanged:: 3.0.0
	Removed the ``-s`` command line switch to compare CDN snapshots - this is now the responsibility of the genConfigRoutes.py script.

The typical way to use ``compare`` is to first specify some environment variables:

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

These can be overridden by command line switches as described above. If a username and/or password is not given for the testing instance (either via environment variables or on the command line), it/they will be assumed to be the same as the one/those specified for the reference instance.

genConfigRoutes.py
------------------

``usage: genConfigRoutes.py [-h] [--refURL REFURL] [--testURL TESTURL]``
                          ``[--refUser REFUSER] [--refPasswd REFPASSWD]``
                          ``[--testUser TESTUSER] [--testPasswd TESTPASSWD] [-k]``
                          ``[-v] [-l LOG_LEVEL] [-q]``

A simple script to generate API routes to server configuration files for a
given pair of Traffic Ops instances. This, for the purpose of using the
'compare' tool


-h, --help                           show this help message and exit
--refURL REFURL                      The full URL of the reference Traffic Ops instance (default: None)
--testURL TESTURL                    The full URL of the testing Traffic Ops instance (default: None)
--refUser REFUSER                    A username for logging into the reference Traffic Ops instance. (default: None)
--refPasswd REFPASSWD                A password for logging into the reference Traffic Ops instance (default: None)
--testUser TESTUSER                  A username for logging into the testing Traffic Ops instance. If not given, the value for the reference instance will be used. (default: None)
--testPasswd TESTPASSWD              A password for logging into the testing Traffic Ops instance. If not given, the value for the reference instance will be used. (default: None)
-k, --insecure                       Do not verify SSL certificate signatures against *either* Traffic Ops instance (default: False)
-v, --version                        Print version information and exit
-l LOG_LEVEL, --log_level LOG_LEVEL  Sets the Python log level, one of 'DEBUG', 'INFO', 'WARN', 'ERROR', or 'CRITICAL' (default: INFO)
-q, --quiet                          Suppresses all logging output - even for critical errors (default: False)
-s, --snapshot                       Produce snapshot routes in the output (CRConfig.json, snapshot/new etc.) (default: False)
-C, --no-server-configs              Do not generate routes for server config files (default: False)


.. note:: If you're using a CDN-in-a-Box environment for testing, it's likely that you'll need the ``-k``/``--insecure`` option if you're outside the Docker network

.. note:: This script will use the same environment variables as `compare`, which can be overridden by the above  command line parameters

The genConfigRoutes.py script will output list of unique API routes (relative to the desired Traffic Ops URL) that point to generated configuration files for a sample set of servers common to both  Traffic Ops instances. The results are printed to stdout, making the output perfect for piping directly into ``compare`` like so:

.. code-block:: shell

	./genConfigRoutes.py https://trafficopsA.example.test https://trafficopsB.example.test username:password | ./compare

\... assuming the proper environment variables have been set for ``compare``.

Usage with Docker
=================
A Dockerfile is provided to run tests on a pair of instances given the configuration environment variables necessary. This will generate configuration file routes using ``genConfigRoutes.py``, and add them to whatever is already contained in ``traffic_ops/testing/compare/testroutes.txt``, then run the ``compare`` tool on the final API route list. Build artifacts (i.e. anything output files created by the `compare` tool) are placed in the `/artifacts/` directory on the container. To retrieve these results, the use of a volume is recommended. The build context *must* be at the root of the Traffic Control repository, as the tools have dependencies on the Traffic Control clients.

Arguments can be passed to the genConfigRoutes.py script by defining the build-time argument ``MODE``. By default it expands to ``-s`` to allow the generation of CDN snapshot routes. It is not necessary to pass ``-k``/``--insecure``, as the Dockerfile will do that implicitly.

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

.. note:: Unlike using the ``genRoutesConfig.py`` script and/or the ``compare`` on their own, *all* of the variables must be defined, even if they are duplicates.
