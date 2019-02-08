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

This package provides an executable script named :program:`traffic_ops_ort`

Usage
=====
``traffic_ops_ort [-k] [--dispersion DISP] [--login_dispersion DISP] [--retries RETRIES] [--wait_for_parents INT] [--rev_proxy_disable] [--ts-root PATH] MODE LOG_LEVEL TO_URL LOGIN``

``traffic_ops_ort [-v]``

``traffic_ops_ort [-h]``

.. option:: -h, --help

	Print usage information and exit

.. option:: -v, --version

	Print version information and exit

.. option:: -k, --insecure

	An optional flag which, when used, disables the checking of SSL certificates for validity

.. option:: --dispersion DISP

	Wait a random number between 0 and ``DISP`` seconds before starting. (Default: 300)

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
		and will exit immediately if this server does not have revalidations pending. Otherwise, the
		same as BADASS.
	SYNCDS
		Will not restart Apache Traffic Server, and will exit immediately if this server does not
		have updates pending. Otherwise, the same as BADASS
	BADASS
		Applies all pending configuration in Traffic Ops, and attempts to solve encountered problems
		when possible. This will install packages, enable/disable system services, and will start or
		restart Apache Traffic Server as necessary.

.. option:: LOG_LEVEL

	Sets the verbosity of output provided by :program:`traffic_ops_ort`. Must be one of:

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

.. option:: TO_URL

	This must be the full URL that refers to the Traffic Ops server, including schema and port
	number (if needed). E.g. ``https://trafficops.infra.ciab.test:443``.

.. option:: LOGIN

	The information used to authenticate with Traffic Ops. This must consist of a username and a
	password, delimited by a colon (``:``). E.g. ``admin:twelve``.

	.. warning:: The first colon found in this string is considered the delimiter. There is no way
		to escape the delimeter. This effectively means that usernames containing colons cannot be
		used to authenticate with Traffic Ops, though passwords containing colons should be fine.

Module Contents
===============
"""

__version__ = "0.0.5"
__author__  = "Brennan Fieck"

import argparse
import datetime
import logging
import random
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

def main():
	"""
	The ORT entrypoint, parses argv before handing it off to :func:`doMain`.
	"""
	# I have no idea why, but the old ORT script does this on every run.
	print(datetime.datetime.utcnow().strftime("%a %b %d %H:%M:%S UTC %Y"))

	parser = argparse.ArgumentParser(description="A Python-based TO_ORT implementation",
	                                 epilog=("Note that passing a negative integer to options that "
	                                         "expect integers will instead set them to zero."),
	                                 formatter_class=argparse.ArgumentDefaultsHelpFormatter)

	parser.add_argument("Mode",
	                    help="REPORT: Do nothing, but print what would be done\n"\
	                         "REPORT, INTERACTIVE, REVALIDATE, SYNCDS, BADASS")
	parser.add_argument("Log_Level",
	                    help="ALL/TRACE, DEBUG, INFO, WARN, ERROR, FATAL/CRITICAL, NONE",
	                    type=str)
	parser.add_argument("Traffic_Ops_URL",
	                    help="URL to Traffic Ops host. Example: https://trafficops.company.net",
	                    type=str)
	parser.add_argument("Traffic_Ops_Login",
	                    help="Example: 'username:password'")
	parser.add_argument("--dispersion",
	                    help="wait a random number between 0 and <dispersion> before starting.",
	                    type=int,
	                    default=300)
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
	parser.add_argument("-v", "--version",
	                    action="version",
	                    version="%(prog)s v"+__version__,
	                    help="Print version information and exit.")

	exit(doMain(parser.parse_args()))
