# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

"""
This package is meant to fully implement the Traffic Ops Operational
Readiness Test - which was originally written in a single, chickenscratch
Perl script. When the :func:`main` function is run, it acts (more or less)
exactly like that legacy script, with the ability to set system configuration
files and start, stop, and restart HTTP cache servers etc.

.. program:: traffic_ops_ort

.. seealso:: Contributions to :program:`traffic_ops_ort` should follow the :ref:`ATC Python contribution guidelines <py-contributing>`

This package provides an executable script named :program:`traffic_ops_ort`

Usage
=====
There are two main ways to invoke :program:`traffic_ops_ort`. The first method uses what's referred
to as the "legacy call signature" and is meant to match the Perl command line arguments.

.. code-block:: text
	:caption: Legacy Call Signature

	traffic_ops_ort [-k] [-h] [-v] [--dispersion DISP] [--login_dispersion DISP]
	         [--retries RETRIES] [--wait_for_parents INT] [--rev_proxy_disable]
	         [--ts_root PATH] MODE LOG_LEVEL TO_URL LOGIN``

The second method - called the "new call signature" - aims to reduce the complexity of the
:term:`ORT` command line. Rather than require a URL and "login string" for connecting and
authenticating with the Traffic Ops server, these pieces of information are optional and may be
provided by the :option:`--to_url`, :option:`-u`/:option:`--to_user`, and
:option:`-p`/:option:`--to_password` options, respectively. If they are NOT provided, then their values
will be obtained from the :envvar:`TO_URL`, :envvar:`TO_USER`, and :envvar:`TO_PASSWORD` environment
variables, respectively. Note that :program:`traffic_ops_ort` cannot be run using the new call
signature without providing a definition for each of these, either on the command line or in the
execution environment.

.. code-block:: text
	:caption: New call signature

	traffic_ops_ort [-k] [-h] [-v] [--dispersion DISP] [--login_dispersion DISP]
	     [--retries RETRIES] [--wait_for_parents INT] [--rev_proxy_disable]
	     [--ts_root PATH] [-l LOG_LEVEL] [-u USER] [-p PASSWORD] [--to_url URL] MODE

These two call signatures should not be mixed, and :program:`traffic_ops_ort` will exit with an
error if they are.

Arguments and Flags
-------------------
.. option:: -h, --help

	Print usage information and exit

.. option:: -v, --version

	Print version information and exit

.. option:: -t, --timeout

	Sets the timeout in milliseconds for connections to Traffic Ops.

.. option:: -k, --insecure

	An optional flag which, when used, disables the checking of SSL certificates for validity

.. option:: --dispersion DISP

	Wait a random number between 0 and ``DISP`` seconds before starting. This option *only* has any
	effect if :option:`MODE` is ``SYNCDS``. (Default: 300)

.. option:: --login_dispersion DISP

	Wait a random number between 0 and ``DISP`` seconds before authenticating with Traffic Ops.
	(Default: 0)

.. option:: --retries RETRIES

	If connection to Traffic Ops fails, retry ``RETRIES`` times before giving up (Default: 3).

.. option:: --wait_for_parents INT

	If ``INT`` is anything but 0, do not apply updates if parents of this server have pending
	updates. This option requires an integer argument for legacy compatibility reasons; 0 is
	considered ``False``, anything else is ``True``. (Default: 1)

.. option:: --rev_prox_disable

	Make requests directly to the Traffic Ops server, bypassing a reverse proxy if one exists.

.. option:: --ts_root PATH

	An optional flag which, if present, specifies the absolute path to the install directory of
	Apache Traffic Server. A common alternative to the default is ``/opt/trafficserver``.
	(Default: ``/``)

.. option:: MODE

	Specifies :program:`traffic_ops_ort`'s mode of operation. Must be one of:

	REPORT
		Runs as though the mode was BADASS, but doesn't actually change anything on the system.

		.. tip:: This is normally useful with a verbose :option:`LOG_LEVEL` to check the state of
			the system

	INTERACTIVE
		Runs as though the mode was BADASS, but asks the user for confirmation before making changes
	REVALIDATE
		Will not restart Apache Traffic Server, install packages, or enable/disable system services
		and will exit immediately if this server does not have revalidations pending. Also, the only
		configuration file that will be updated is `regex_revalidate.config`.
	SYNCDS
		Will not restart Apache Traffic Server, and will exit immediately if this server does not
		have updates pending. Otherwise, the same as BADASS
	BADASS
		Applies all pending configuration in Traffic Ops, and attempts to solve encountered problems
		when possible. This will install packages, enable/disable system services, and will start or
		restart Apache Traffic Server as necessary.

.. option:: LOG_LEVEL, -l LOG_LEVEL, --log_level LOG_LEVEL

	Sets the verbosity of output provided by :program:`traffic_ops_ort`. This argument is positional
	in the legacy call signature, but optional in the new call signature, wherein it has a default
	value of "WARN". Must be one of (case-insensitive):

	NONE
		Will output nothing, not even fatal errors.
	CRITICAL
		Will only output error messages that indicate an unrecoverable error.
	FATAL
		A synonym for "CRITICAL"
	ERROR
		Will output more general errors about conditions that are causing problems of some kind.
	WARN
		In addition to error information, will output warnings about conditions that may cause
		problems, or possible misconfiguration.
	INFO
		Outputs informational messages about what :program:`traffic_ops_ort` is doing as it
		progresses.
	DEBUG
		Outputs detailed debug information, including stack traces.

		.. note:: Not all stack traces indicate problems with :program:`traffic_ops_ort`. Stack
			traces are printed whenever an exception is encountered, whether or not it could be
			handled.

	TRACE
		A synonym for "DEBUG"
	ALL
		A synonym for "DEBUG"

	.. note:: All logging is sent to STDERR. INTERACTIVE :option:`MODE` prompts are printed to STDOUT

.. option:: TO_URL, --to_url TO_URL

	This must be at minimum an :abbr:`FQDN (Fully Qualified Domain Name)` that resolves to the
	Traffic Ops server, but may optionally include the schema and/or port number. E.g.
	``https://trafficops.infra.ciab.test:443``, ``https://trafficops.infra.ciab.test``,
	``trafficops.infra.ciab.test:443``, and ``trafficops.infra.ciab.test`` are all acceptable, and
	in fact are all equivalent. When given a value without a schema, HTTPS will be the assumed
	protocol, and when a port number is not present, 443 will be assumed except in the case that
	the schema *is* provided and is ``http://`` (case-insensitive) in which case 80 will be assumed.

	This argument is positional in the legacy call signature, but is optional in the new call
	signature. When the new call signature is used and this option is not present on the command
	line, its value will be obtained from :envvar:`TO_URL`. Note that :program:`traffic_ops_ort`
	cannot be run using the new call signature unless this value is defined, either on the command
	line or in the execution environment.

.. option:: LOGIN

	The information used to authenticate with Traffic Ops. This must consist of a username and a
	password, delimited by a colon (``:``). E.g. ``admin:twelve``. This argument is not used in the
	new call signature, instead :option:`-u`/:option:`--to_user` and
	:option:`-p`/:option:`--to_password` are used to separately set the authentication user and
	password, respectively.

	.. warning:: The first colon found in this string is considered the delimiter. There is no way
		to escape the delimeter. This effectively means that usernames containing colons cannot be
		used to authenticate with Traffic Ops, though passwords containing colons should be fine.

.. option:: -u USER, --to_user USER

	Specifies the username of the user as whom to authenticate when connecting to Traffic Ops. This
	option is only available using the new call signature. If not provided when using said new call
	signature, the value will be obtained from the :envvar:`TO_USER` environment variable. Note that
	:program:`traffic_ops_ort` cannot be run using the new call signature unless this value is
	defined, either on the command line or in the execution environment.

.. option:: -p PASSWORD, --to_password PASSWORD

	Specifies the password of the user identified by :envvar:`TO_USER` (or
	:option:`-u`/:option:`--to_user` if overridden) to use when authenticating to Traffic Ops. This
	option is only available using the new call signature. If not provided when using said new call
	signature, the value will be obtained from the :envvar:`TO_PASSWORD`  environment variable. Note
	that :program:`traffic_ops_ort` cannot be run using the new call signature unless this value is
	defined, either on the command line or in the execution environment.

.. option:: --hostname HOSTNAME

	Causes ORT to request configuration information for the server named ``HOSTNAME`` instead of
	detecting the server's actual hostname. This is primarily useful for testing purposes.

Environment Variables
---------------------
:program:`traffic_ops_ort` supports authentication with a Traffic Ops instance using the environment
variables :envvar:`TO_URL`, :envvar:`TO_USER` and :envvar:`TO_PASSWORD`.

.. _ort-special-strings:

Strings with Special Meaning to ORT
===================================
When processing configuration files, if :program:`traffic_ops_ort` encounters any of the strings in
the :ref:`Replacement Strings <ort-replacement-strings>` table it will perform the indicated
replacement. This means that these strings can be used to create templates in :term:`Profile`
:term:`Parameters` and certain :term:`Delivery Service` configuration fields.

.. _ort-replacement-strings:

.. table:: Replacement Strings

	+-------------------------+--------------------------------------------------------------------+
	| String                  | Replaced With                                                      |
	+=========================+====================================================================+
	| ``__CACHE_IPV4__``      | The IPv4 address of the :term:`cache server` on which              |
	|                         | :program:`traffic_ops_ort` is running.                             |
	+-------------------------+--------------------------------------------------------------------+
	| ``__FULL_HOSTNAME__``   | The full hostname (i.e. including the full domain to which it      |
	|                         | belongs) of the :term:`cache server` on which                      |
	|                         | :program:`traffic_ops_ort` is running.                             |
	+-------------------------+--------------------------------------------------------------------+
	| ``__HOSTNAME__``        | The (short) hostname of the :term:`cache server` on which          |
	|                         | :program:`traffic_ops_ort` is running.                             |
	+-------------------------+--------------------------------------------------------------------+
	| ``__RETURN__``          | A newline character (``\\n``).                                      |
	+-------------------------+--------------------------------------------------------------------+
	| ``__SERVER_TCP_PORT__`` | If the :term:`cache server` on which :program:`traffic_ops_ort` is |
	|                         | being run has a TCP port configured to something besides ``80``,   |
	|                         | this will be replaced with that TCP port value. *If it* **is**     |
	|                         | *set to ``80``, this string will simply be removed,* **NOT**       |
	|                         | *replaced with* **ANYTHING**.                                      |
	+-------------------------+--------------------------------------------------------------------+
	| ``##OVERRIDE##``        | This string is only valid in the content of files named            |
	|                         | "remap.config". It is further described in `Remap Override`_       |
	+-------------------------+--------------------------------------------------------------------+

.. deprecated:: ATCv4.0
	The use of ``__RETURN__`` in lieu of a true newline character is (finally) no longer necessary,
	and the ability to do so will be removed in the future.

.. note:: There is currently no way to indicate that a server's IPv6 address should be inserted -
	only IPv4 is supported.

.. _ort-remap-override:

Remap Override
--------------
.. warning:: The ANY_MAP ``##OVERRIDE##`` special string is a temporary solution and will be
	deprecated once Delivery Service Versioning is implemented. For this reason it is suggested that
	it not be used unless absolutely necessary.

The ``##OVERRIDE##`` template string allows the :term:`Delivery Service` :ref:`ds-raw-remap` field
to override to fully override the :term:`Delivery Service`'s line in the
`remap.config ATS configuration file <https://docs.trafficserver.apache.org/en/7.1.x/admin-guide/files/remap.config.en.html>`_,
generated by Traffic Ops. The end result is the original, generated line commented out, prepended
with ``##OVERRIDDEN##`` and the ``##OVERRIDE##`` rule is activated in its place. This behavior is
used to incrementally deploy plugins used in this configuration file. Normally, this entails cloning
the :term:`Delivery Service` that will have the plugin, ensuring it is assigned to a subset of the
:term:`cache servers` that serve the :term:`Delivery Service` content, then using this
``##OVERRIDE##`` rule to create a ``remap.config`` rule that will use the plugin, overriding the
normal rule. Simply grow the subset over time at the desired rate to slowly deploy the plugin. When
it encompasses all :term:`cache servers` that serve the original :term:`Delivery Service`'s content,
the "override :term:`Delivery Service`" can be deleted and the original can use a
non-``##OVERRIDE##`` :ref:`ds-raw-remap` to add the plugin.

.. code-block:: text
	:caption: Example of Remap Override

	# This is the original line as generated by Traffic Ops
	map http://from.example.com/ http://to.example.com/

	# This is the raw remap text as configured on the delivery service
	##OVERRIDE## map http://from.example.com/ http://to.example.com/ some_plugin.so

	# The resulting content is what actually winds up in the remap.config file:
	##OVERRIDE##
	map http://from.example.com/ http://to.example.com/ some_plugin.so
	##OVERRIDDEN## map http://from.example.com/ http://to.example.com/

.. warning:: The "from" URL must exactly match for this to properly work (e.g. including trailing
	URL '/'), otherwise :abbr:`ATS (Apache Traffic Server)` may fail to initialize or reload while
	processing :file:`remap.config`.

.. tip:: To assist in troubleshooting, it is strongly recommended that any ``##OVERRIDE##`` rules in
	use should be documented on the original :term:`Delivery Service`.

Module Contents
===============
"""

__version__ = "0.11.0"
__author__  = "Brennan Fieck"

import argparse
import datetime
from distutils.spawn import find_executable
import logging
import os
import random
import sys
import time

from requests.exceptions import RequestException
from trafficops.restapi import LoginError, OperationError, InvalidJSONError

def doMain(args:argparse.Namespace) -> int:
	"""
	Runs the main routine based on the parsed arguments to the script

	:param args: A parsed argument list as returned from :meth:`argparse.ArgumentParser.parse_args`
	:returns: an exit code for the script.
	:raises AttributeError: when the namespace is missing required arguments
	"""
	from . import configuration, main_routines, to_api
	random.seed(time.time())

	try:
		conf = configuration.Configuration(args)
	except ValueError as e:
		logging.critical(e)
		logging.debug("%r", e, exc_info=True, stack_info=True)
		return 1

	if conf.login_dispersion:
		disp = random.randint(0, conf.login_dispersion)
		logging.info("Login dispersion is active - sleeping for %d seconds before continuing", disp)
		time.sleep(disp)

	try:
		with to_api.API(conf) as api:
			conf.api = api
			return main_routines.run(conf)
	except (LoginError, OperationError, InvalidJSONError, RequestException) as e:
		logging.critical("Failed to connect and authenticate with the Traffic Ops server")
		logging.error(e)
		logging.debug("%r", e, exc_info=True, stack_info=True)
		return 1

_EPILOG = """traffic_ops_ort supports two calling conventions, one is intended to be fully
compatible with the Perl implementation, while the other is intended to be an improvement over the
former. Essentially this means that either the `Log_Level`, `Traffic_Ops_URL` and
`Traffic_Ops_Login`` must all be given, or none of them. If none of them are given, the log level
will be determined by the `-l`/`--log_level` option, the Traffic Ops server URL will be constructed
from the information available in the `TO_URL` environment and/or the `--to_url` option, the
Traffic Ops user will be determined from the information available in the `TO_USER` environment
variable and/or the `-u`/`--to_user` option, and the Traffic Ops user's password will be determined
from the information available in the `TO_PASSWORD` environment variable and/or the
`-p`/`--to_password` option.
""".replace('\n', ' ') + "\n\n" + ("Note that passing a negative integer to options that expect "
                                   "integers will instead set them zero.")

def main() -> int:
	"""
	The ORT entrypoint, parses argv before handing it off to :func:`doMain`.

	:returns: An exit code for :program:`traffic_ops_ort`
	"""
	global _EPILOG

	# I have no idea why, but the old ORT script does this on every run.
	print(datetime.datetime.utcnow().strftime("%a %b %d %H:%M:%S UTC %Y"))

	from .configuration import LogLevels, Configuration

	runModesAllowed = {str(x) for x in Configuration.Modes}.union(
	                  {str(x).lower() for x in Configuration.Modes})
	logLevelsAllowed = {str(x) for x in LogLevels}.union({str(x).lower() for x in LogLevels})

	parser = argparse.ArgumentParser(description="A Python-based TO_ORT implementation",
	                                 epilog=_EPILOG,
	                                 formatter_class=argparse.ArgumentDefaultsHelpFormatter)

	parser.add_argument("Mode",
	                    help="REPORT: Do nothing, but print what would be done\n"\
	                         "REPORT, INTERACTIVE, REVALIDATE, SYNCDS, BADASS",
	                    choices=runModesAllowed,
	                    type=str)
	parser.add_argument("legacy",
	                    help="ALL/TRACE, DEBUG, INFO, WARN, ERROR, FATAL/CRITICAL, NONE",
	                    metavar="Log_Level",
	                    nargs="?",
	                    action="append",
	                    choices=logLevelsAllowed,
	                    type=str)
	parser.add_argument("-l", "--log_level",
	                    help="Sets the logging level. (Default: WARN)",
	                    type=str,
	                    choices=logLevelsAllowed)
	parser.add_argument("legacy",
	                    help="URL to Traffic Ops host. Example: https://trafficops.company.net",
	                    metavar="Traffic_Ops_URL",
	                    nargs="?",
	                    action="append",
	                    type=str)
	parser.add_argument("--to_url",
	                    help=("A URL or hostname - optionally with a port specification - that "
	                          "points to a Traffic Ops server. e.g. `trafficops.infra.ciab.test` or"
	                          " `https://trafficops.infra.ciab.test:443`. This overrides the TO_URL"
	                          " environment variable"),
	                    type=str)
	parser.add_argument("legacy",
	                    help="Example: 'username:password'",
	                    metavar="Traffic_Ops_Login",
	                    nargs="?",
	                    action="append",
	                    type=str)
	parser.add_argument("-u", "--to_user",
	                    help=("The username to use when authenticating to the Traffic Ops server. "
	                          "This overrides the TO_USER environment variable."),
	                    type=str)
	parser.add_argument("-p", "--to_password",
	                    help=("The password to use when authenticating to the Traffic Ops server. "
	                          "This overrides the TO_PASSWORD environment variable."),
	                    type=str)
	parser.add_argument("--dispersion",
	                    help="wait a random number between 0 and <dispersion> before starting.",
	                    type=int,
	                    default=300)
	parser.add_argument("--hostname",
	                    help="Pretend to be a server with the provided hostname instead of using "
	                         "this server's actual hostname in communications with Traffic Ops")
	parser.add_argument("--login_dispersion",
	                    help="wait a random number between 0 and <login_dispersion> before login.",
	                    type=int,
	                    default=0)
	parser.add_argument("--retries",
	                    help="retry connection to Traffic Ops URL <retries> times.",
	                    type=int,
	                    default=3)
	parser.add_argument("--wait_for_parents",
	                    help="do not update if parent_pending = 1 in the update json.",
	                    type=int,
	                    default=1)
	parser.add_argument("--rev_proxy_disable",
	                    help="bypass the reverse proxy even if one has been configured.",
	                    action="store_true")
	parser.add_argument("--ts_root",
	                    help="Specify the root directory at which Apache Traffic Server is installed"\
	                         " (e.g. '/opt/trafficserver')",
	                    type=str,
	                    default="/")
	parser.add_argument("-k", "--insecure",
	                    help="Skip verification of SSL certificates for Traffic Ops connections. "\
	                         "DON'T use this in production!",
	                    action="store_true")
	parser.add_argument("-t", "--timeout",
	                    help="Sets the timeout in milliseconds for requests made to Traffic Ops.",
	                    type=int,
	                    default=None)
	parser.add_argument("--via-string-release",
			    		help="set the ATS via string to the package release instead of version",
			    		type=int,
			    		default=0)
	parser.add_argument("--disable-parent-config-comments",
						help="Do not insert comments in parent.config files",
						type=int,
						default=0)
	parser.add_argument("-v", "--version",
	                    action="version",
	                    version="%(prog)s v"+__version__,
	                    help="Print version information and exit.")

	args = parser.parse_args()

	# New call signature
	if None in args.legacy:
		if any(args.legacy):
			print("Legacy mode call signature cannot be partial!", file=sys.stderr)
			print("(Hint: use -h/--help for usage)", file=sys.stderr)
			return 1

		try:
			args.to_url = args.to_url if args.to_url else os.environ["TO_URL"]
			args.to_password = args.to_password if args.to_password else os.environ["TO_PASSWORD"]
			args.to_user = args.to_user if args.to_user else os.environ["TO_USER"]
		except KeyError as e:
			print("Neither option nor environment variable defined for %s!" % e.args[0],
			      file=sys.stderr)
			print("(Hint: use -h/--help for usage)", file=sys.stderr)
			return 1

		args.log_level = args.log_level if args.log_level else "WARN"

	# Illegal mixed call signature
	elif (args.to_url is not None or
	      args.to_user is not None or
	      args.log_level is not None or
	      args.to_password is not None):

		print("Do not mix legacy call signature with new-style call signature!", file=sys.stderr)
		print("(Hint: use -h/--help for usage)", file=sys.stderr)
		return 1

	# Legacy call signature
	else:
		args.log_level, args.to_url, login = args.legacy
		try:
			args.to_user, args.to_password = login.split(':')
		except ValueError:
			print("Invalid Traffic_Ops_Login format! Should be 'username:password'",file=sys.stderr)
			print("(Hint: use -h/--help for usage)", file=sys.stderr)
			return 1

	if not find_executable("t3c"):
		print("Could not find t3c executable - this is required to run ORT!", file=sys.stderr)
		return 1

	return doMain(args)
