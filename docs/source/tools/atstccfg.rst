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

.. program: atstccfg

.. _atstccfg:

********
atstccfg
********
:program:`atstccfg` is a tool for generating configuration files server-side on :abbr:`ATC (Apache Traffic Control)` cache servers. It stores its generated/cached files in ``/tmp/atstccfg_cache/`` for re-use.

.. warning:: :program:`atstccfg` does not have a stable command-line interface, it can and will change without warning. Scripts should avoid calling it for the time being, as its only intended caller is :term:`ORT`.

The source code for :program:`atstccfg` may be found in :atc-file:`traffic_ops/ort/atstccfg`, and the Go module reference is :atc-godoc:`traffic_ops/ort/atstccfg`.

Usage
=====
``atstccfg [-u TO_URL] [-U TO_USER] [-P TO_PASSWORD] [-n] [-r N] [-e ERROR_LOCATION] [-w WARNING_LOCATION] [-i INFO_LOCATION] [-g] [-s] [-t TIMEOUT] [-a MAX_AGE] [-l] [-h] [-v]``

Options
-------
.. option:: -a AGE, --cache-file-max-age-seconds AGE

	Sets the maximum age - in seconds - a cached response can be in order to be considered "fresh" - older files will be re-generated and cached. Default: 60

.. option:: -e ERROR_LOCATION, --log-location-error ERROR_LOCATION

	The file location to which to log errors. Respects the special string constants of :atc-godoc:`lib/go-log`. Default: 'stderr'

.. option:: -g, --print-generated-files

	If given, the names of files generated (and not proxied to Traffic Ops) will be printed to stdout, then :program:`atstccfg` will exit.

.. option:: -i INFO_LOCATION, --log-location-info INFO_LOCATION

	The file location to which to log information messages. Respects the special string constants of :atc-godoc:`lib/go-log`. Default: 'stderr'

.. option:: -h, --help

	Print usage information and exit.

.. option:: -l, --list-plugins

	List the loaded plugins and then exit.

.. option:: -n, --no-cache

	If given, existing cache files will not be used. Cache files will still be created, existing ones just won't be used.

.. option:: -P TO_PASSWORD, --traffic-ops-password TO_PASSWORD

	Authenticate using this password - if not given, atstccfg will attempt to use the value of the :envvar:`TO_PASS` environment variable.

.. option:: -r N, --num-retries N

	The number of times to retry getting a file if it fails. Default: 5

.. option:: -s, --traffic-ops-insecure

	If given, SSL certificate errors will be ignored when communicating with Traffic Ops.

	.. caution:: For (hopefully) obvious reasons, the use of this option in production environments is discouraged.

.. option:: -t TIMEOUT, --traffic-ops-timeout-milliseconds TIMEOUT

	Sets the timeout - in milliseconds - for requests made to Traffic Ops. Default: 10000

.. option:: -u TO_URL, --traffic-ops-url TO_URL

	Request this URL, e.g. ``https://trafficops.infra.ciab.test/servers/edge/configfiles/ats``. If not given, :program:`atstccfg` will attempt to use the value of the :envvar:`TO_URL` environment variable.

.. option:: -U TO_USER, --traffic-ops-user TO_USER

	Authenticate as the user ``TO_USER`` - if not given, :program:`atstccfg` will attempt to use the value of the :envvar:`TO_USER` environment variable.

.. option:: -v, --version

	Print version information and exit.

.. option:: -w WARNING_LOCATION, --log-location-warning WARNING_LOCATION

	The file location to which to log warnings. Respects the special string constants of :atc-godoc:`lib/go-log`. Default: 'stderr'

Environment Variables
---------------------

.. envvar:: TO_USER

	Defines the user as whom to authenticate with Traffic Ops. This is only used if :option:`-U`/:option:`--traffic-ops-user` is not given.

.. envvar:: TO_PASS

	Defines the password to use when authenticating with Traffic Ops. This is only used if :option:`-P`/:option:`--traffic-ops-password` is not given.

.. envvar:: TO_URL

	Defines the *full* URL to be requested. This is only used if :option:`-u`/:option:`--traffic-ops-url` is not given.
