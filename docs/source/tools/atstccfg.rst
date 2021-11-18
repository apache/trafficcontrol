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
:program:`atstccfg` is a tool for generating configuration files server-side on :abbr:`ATC (Apache Traffic Control)` cache servers.

.. warning:: :program:`atstccfg` does not have a stable command-line interface, it can and will change without warning. Scripts should avoid calling it for the time being, as its only intended caller is :term:`ORT`.

The source code for :program:`atstccfg` may be found in :atc-file:`traffic_ops/ort/atstccfg`, and the Go module reference is :atc-godoc:`traffic_ops/ort/atstccfg`.

Usage
=====
- ``atstccfg -h``
- ``atstccfg -v``
- ``atstccfg -l``
- ``atstccfg [-e ERROR_LOCATION] [-i INFO_LOCATION] [-p] [-P TO_PASSWORD] [-r N] [-s] [-t TIMEOUT] [-u TO_URL] [-U TO_USER] [-w WARNING_LOCATION] [-y] [--dir TSROOT] -n CACHE_NAME``
- ``atstccfg [-e ERROR_LOCATION] [-i INFO_LOCATION] [-p] [-P TO_PASSWORD] [-r N] [-s] [-t TIMEOUT] [-u TO_URL] [-U TO_USER] [-w WARNING_LOCATION] [--dir TSROOT] -n CACHE_NAME -d DATA``
- ``atstccfg [-e ERROR_LOCATION] [-i INFO_LOCATION] [-p] [-P TO_PASSWORD] [-r N] [-s] [-t TIMEOUT] [-u TO_URL] [-U TO_USER] [-w WARNING_LOCATION] [--dir TSROOT] -n CACHE_NAME -a REVAL_STATUS -q QUEUE_STATUS``

When called using the fourth form, :program:`atstccfg` provides its normal output. This is the entirety of all configuration files necessary for the server, all provided at once. The output is in :mimetype:`mixed/multipart` format, defined by :rfc:`1521`. Each chunk of the message comes with the proprietary ``Path`` header, which specifies the exact location on disk of the file whose contents are contained in that chunk.

Options
-------
.. option:: -a REVAL_STATUS, --set-reval-status REVAL_STATUS

	Sets the ``reval_pending`` property of the server in Traffic Ops. Must be 'true' or 'false'. Requires :option:`--set-queue-status` also be set. This disables normal output.

.. option:: -e ERROR_LOCATION, --log-location-error ERROR_LOCATION

	The file location to which to log errors. Respects the special string constants of :atc-godoc:`lib/go-log`. Default: 'stderr'

.. option:: -d DATA, --get-data DATA

	Specifies non-configuration-file data to retrieve from Traffic Ops. This disables normal output. Valid values are:

	chkconfig
		Retrieves information about the services which should be running on the :term:`cache server`. The output will be in JSON-encoded format as an array of objects with these fields:

		:name:  The name of the service. This should correspond to an existing systemd service unit file.
		:value: A "chkconfig" line describing on which "run-levels" the services should be running. See the :manpage:`chkconfig(8)` manual pages for details on what this field means.

	packages
		Retrieves information about the packages which should exist on the :term:`cache server`. The output will be in JSON-encoded format as an array of
		objects with these fields:

		:name:    The name of the package. This should hopefully be a meaningful package name for the :term:`cache server`'s package management system.
		:version: The version of the package which should be installed. This might also be an empty string which means "any version will do".

	statuses
		Retrieves all statuses from Traffic Ops. This is defined to be exactly the ``response`` object from the response to a GET request made to the :ref:`to-api-statuses` Traffic Ops API endpoint.
	system-info
		Retrieves generic information about the Traffic Control system from the :ref:`to-api-system-info` API endpoint. The output is the ``parameters`` object of the responses from GET requests to that endpoint (still JSON-encoded).
	update-status
		Retrieves information about the current update status using :ref:`to-api-servers-hostname-update_status`. The response is in the same format as the responses for that endpoint's GET method handler - except that that endpoint returns an array and this :program:`atstccfg` call signature returns a single one of those elements. Which one is chosen is arbitrary (hence undefined behavior when more than one server with the same hostname exists).

.. option:: --dir TSROOT

	Specifies a directory path in which to place Traffic Server configuration
	files, in the event that "location" :term:`Parameters` are not found for
	them. If this is not given and location :term:`Parameters` are not found for
	required files, then :program:`atstccfg` will exit with an error.

	The files that :program:`atstccfg` considers "required" for these purposes
	are:

	- cache.config
	- hosting.config
	- parent.config
	- plugin.config
	- records.config
	- remap.config
	- storage.config
	- volume.config

.. option:: -h, --help

	Print usage information and exit.

.. option:: -i INFO_LOCATION, --log-location-info INFO_LOCATION

	The file location to which to log information messages. Respects the special string constants of :atc-godoc:`lib/go-log`. Default: 'stderr'

.. option:: -l, --list-plugins

	List the loaded plug-ins and then exit.

.. option:: -n NAME, --cache-host-name NAME

	Required. Specifies the (short) hostname of the :term:`cache server` for which output will be generated. Must be the server hostname in Traffic Ops, not a URL, or :abbr:`FQDN (Fully Qualified Domain Name)`. Behavior when more than one server exists with the passed hostname is undefined.

.. option:: -p, --traffic-ops-disable-proxy

	Bypass the Traffic Ops caching proxy and make requests directly to Traffic Ops. Has no effect if no such proxy exists.

.. option:: -P TO_PASSWORD, --traffic-ops-password TO_PASSWORD

	Authenticate using this password - if not given, :program:`atstccfg` will attempt to use the value of the :envvar:`TO_PASSWORD` environment variable.

.. option:: -q QUEUE_STATUS, --set-queue-status QUEUE_STATUS

	Sets the ``upd_pending`` property of the server identified by :option:`--cache-host-name` to the specified value, which must be 'true' or 'false'. Requires :option:`--set-reval-status` to also be set.

.. option:: -r N, --num-retries N

	The number of times to retry getting a file if it fails. Default: 5

.. option:: -s, --traffic-ops-insecure

	If given, SSL certificate errors will be ignored when communicating with Traffic Ops.

	.. caution:: The use of this option in production environments is discouraged.

.. option:: -t TIMEOUT, --traffic-ops-timeout-milliseconds TIMEOUT

	Sets the timeout - in milliseconds - for requests made to Traffic Ops. Default: 30000

.. option:: -u TO_URL, --traffic-ops-url TO_URL

	Request this URL, e.g. ``https://trafficops.infra.ciab.test/servers/edge/configfiles/ats``. If not given, :program:`atstccfg` will attempt to use the value of the :envvar:`TO_URL` environment variable.

.. option:: -U TO_USER, --traffic-ops-user TO_USER

	Authenticate as the user ``TO_USER`` - if not given, :program:`atstccfg` will attempt to use the value of the :envvar:`TO_USER` environment variable.

.. option:: -v, --version

	Print version information and exit.

.. option:: -w WARNING_LOCATION, --log-location-warning WARNING_LOCATION

	The file location to which to log warnings. Respects the special string constants of :atc-godoc:`lib/go-log`. Default: 'stderr'

.. option:: -y, --revalidate-only

	When given, :program:`atstccfg` will only emit files relevant for updating :term:`Content Invalidation Jobs`. for Apache Traffic Server implementations, this limits the output to be only files named ``regex_revalidate.config``. Has no effect if :option:`--get-data` or :option:`--set-queue-status`/:option:`--set-reval-status` is/are used.

Environment Variables
---------------------
:program:`atstccfg` supports authentication with a Traffic Ops instance using the environment variables :envvar:`TO_URL` (if :option:`-u`/:option:`--traffic-ops-url` is not given), :envvar:`TO_USER` (if :option:`-U`/:option:`--traffic-ops-user` is not given), and :envvar:`TO_PASSWORD` (if :option:`-P`/:option:`--traffic-ops-password` is not given).
