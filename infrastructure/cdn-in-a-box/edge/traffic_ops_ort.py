#!/usr/bin/env python3
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
This script aims to be a drop-in replacement for the aged
``traffic_ops_ort.pl`` script. Its primary purpose is for
management of a cache server via configuration from
Traffic Ops. This script will install/upgrade packages
as necessary, will ensure services that ought to be running
are running, and sets up ATS and ATS Plugin configuration files.
"""


import argparse
import datetime
import sys
import os
import platform
import typing
import logging
import enum
import subprocess
import time
import re

needInstall = []

try:
	import requests
except ImportError:
	logging.error("You must have the 'requests' package installed to use this script.")
	logging.warning("(Hint: try `pip3 install requests`)")
	needInstall.append("requests")

try:
	import urllib3
	urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
except ImportError:
	logging.error("You must have the `urllib3` library installed to use this script.")
	logging.warning("(Needed by the `requests` package)")
	logging.warning("(Hint: try `pip3 install urllib3`)")
	needInstall.append("urllib3")

try:
	import distro
	DISTRO = distro.LinuxDistribution().id()
except ImportError:
	logging.error("You must have the `distro` package installed to use this script.")
	logging.warning("Hint: try `pip3 install distro`)")
	needInstall.append("distro")

try:
	import psutil
except ImportError:
	logging.error("You must have the `psutil` package installed to use this script.")
	logging.warning("(Hint: try `pip3 install psutil`)")
	needInstall.append("psutil")

__version__ = "0.1"

#: Holds the full URL including schema (e.g. 'http') and port that points at Traffic Ops
TO_URL = None

#: Holds the Traffic Ops login information as a JSON-encoded string
TO_LOGIN = None

#: Holds the Mojolicious cookie returned by Traffic Ops on successful login, and needed for
#: subsequent requests.
TO_COOKIE = None

#: An absolute path to the root installation directory of Apache Trafficserver
TS_ROOT = None

#: A map of command-line log level options to :mod:`logging` log level constants
LOG_LEVELS = {"ALL":   logging.NOTSET,
              "DEBUG": logging.DEBUG,
              "INFO":  logging.INFO,
              "WARN":  logging.WARNING,
              "ERROR": logging.ERROR,
              "FATAL": logging.CRITICAL}

#: A format specifier for logging output
FMT = "%(levelname)s: line %(lineno)d in %(module)s.%(funcName)s: %(message)s"

#: Contains the host's hostname as a tuple of ``(short_hostname, full_hostname)``
HOSTNAME = (platform.node().split('.')[0], platform.node())

class Modes(enum.IntEnum):
	"""
	Enumerated run modes
	"""
	REPORT = 0      #: Do nothing, only report what would be done
	INTERACTIVE = 1 #: Ask for user confirmation before modifying the system
	REVALIDATE = 2  #: Only check for config file changes and content revalidations
	SYNCDS = 3      #: Check for and apply Delivery Service changes
	BADASS = 4      #: Apply all settings specified in Traffic Ops, and attempt to solve all problems

	def __str__(self) -> str:
		"""
		Implements ``str(self)`` by returning enum member's name
		"""
		return self.name

#: Stores the current Run Mode
MODE = None

class ORTException(Exception):
	"""Represents an error while processing ORT API responses, etc."""
	pass

#: This is the set of files which will require an ATS restart when changed
ATS_FILES = {"records.config",
             "remap.config",
             "parent.config",
             "cache.config",
             "hosting.config",
             "astats.config",
             "logs_xml.config",
             "ssl_multicert.config"}

#: Global state variable that tracks whether or not ATS should be restarted
ATS_NEEDS_RESTART = False

###############################################################################
#####                                                                     #####
#####                     PYTHON DEPENDENCY HANDLING                      #####
#####                                                                     #####
###############################################################################
def installPythonPackages(packages:typing.List[str]) -> bool:
	"""
	Attempts to install the packages listed in ``packages`` and
	add them to the global scope.

	:param packages: A list of missing Python dependencies

	:return: a truthy value indicating success.
	"""
	logging.info("Attempting install of %s", ','.join(packages))

	# Ensure `pip` is installed
	try:
		import pip
	except ImportError as e:
		logging.info("`pip` package not installed. Attempting install with `ensurepip` module")
		logging.debug("%s", e, exc_info=True, stack_info=True)
		import ensurepip
		try:
			ensurepip.bootstrap(altinstall=True)
		except (EnvironmentError, PermissionError) as err:
			logging.info("Permission Denied, attempting user-level install")
			logging.debug("%s", err, exc_info=True, stack_info=True)
			ensurepip.bootstrap(altinstall=True, user=True)
		import pip

	# Get main pip function
	try:
		pipmain = pip.main
	except AttributeError as e: # This happens when using a version of pip >= 10.0
		from pip._internal import main as pipmain


	# Attempt package install
	ret = pipmain(["install"] + packages)
	if ret == 1:
		logging.info("Possible 'Permission Denied' error, attempting user-level install")
		ret = pipmain(["install", "--user"] + packages)

	logging.debug("Pip return code was %d", ret)
	if not ret:
		import importlib
		for package in packages:
			try:
				globals()[package] = importlib.import_module(package)
			except ImportError:
				logging.error("Failed to import %s", package)
				logging.warning("Install appeared succesful - subsequent run may succeed")
				logging.debug("%s", e, exc_info=True, stack_info=True)
				return False

		urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
		globals()["DISTRO"] = distro.LinuxDistribution().id()

	return not ret

def handleMissingPythonPackages(packages:typing.List[str]) -> bool:
	"""
	Handles the case of missing packages, by installing them in :attr:`Modes.BADASS` mode,
	asking the user for confirmation in :attr:`Modes.INTERACTIVE` mode, and failing in all
	other run modes.

	Installs packages by calling :func:`installPythonPackages`.

	:param packages: A list of missing Python dependencies
	:return: A boolean indicator of success regarding the installation of missing dependencies
	"""
	logging.info("Packages needed for this script but not installed: %s", ','.join(packages))

	if MODE == Modes.BADASS:
		logging.warning("Mode is BADASS - attempting to install missing packages")
		doInstall = True
	elif MODE == Modes.INTERACTIVE:
		print("The following packages are needed to run this script, but are not installed:")
		print(','.join(packages))
		choice = input("Would you like the script to attempt to install them now? [y/N]: ")

		while choice and choice.lower() not in {'y', 'n', 'yes', 'no'}:
			print("Invalid choice:", choice, file=sys.stderr)
			choice = input("Would you like the script to attempt to install them now? [y/N]: ")

		doInstall = choice and choice.lower() in {'y', 'yes'}
	else:
		logging.error("Cannot install packages in current run mode.")
		return False

	if doInstall:
		logging.warning("Installing python packages: %s", ','.join(packages))
		if installPythonPackages(packages):
			logging.info("Python packages installed successfully.")
		else:
			logging.critical("Missing Python packages could not be installed.")
			return False
	else:
		logging.critical("Packages %s missing and will not be installed - cannot continue",
		                 ','.join(packages))
		return False

	return True


###############################################################################
#####                                                                     #####
#####                              UTILITIES                              #####
#####                                                                     #####
###############################################################################
def getJSONResponse(uri:str, expectedStatus:int = 200) -> typing.Optional[object]:
	"""
	Gets a JSON-encoded response to a Traffic Ops API ``GET`` request

	:param uri: The request URI, relative to :data:`TO_URL`
	:param expectedStatus: The expected status code of the response. If the actual response doesn't
		match, this function will return :const:`None`.

	:returns: On success, a parsed JSON object - or :const:`None` on failure.
	"""
	global TO_URL, TO_COOKIE

	logging.info("Getting json response via 'HTTP GET %s%s", TO_URL, uri)

	response = requests.get(TO_URL + uri, cookies=TO_COOKIE, verify=False)

	if response.status_code != expectedStatus:
		logging.error("Failed to get a response from '%s': server returned status code %d",
		              uri,
		              response.status_code)
		logging.debug("Response: %s\n%r\n%r", response, response.headers, response.content)
		# For some applications, empty responses may be valid, so returning None rather than exiting.
		return None

	return response.json()

def getRawResponse(uri:str,
                   expectedStatus:int = 200,
                   TOrelative:bool = True,
                   verify:bool = False) -> typing.Optional[str]:
	"""
	Gets the raw response body of an HTTP ``GET`` request - optionally one that is relative to
	:data:`TO_URL`.

	.. note:: Actually encodes to utf-8 string

	:param uri: The full (unless ``Torelative`` is :const:`True`) path to a resource for the request
	:param expectedStatus: The expected status code of the response. If the actual response doesn't
		match, this function will return :const:`None`
	:param TOrelative: If true, the ``uri`` will be treated as relative to :data:`TO_URL`
	:param verify: If true, the SSL keys used to communicate with the full URI will be verified

	:return: The body of the response on success, :const:`None` on failure
	"""
	global TO_URL, TO_COOKIE

	if TOrelative:
		uri = TO_URL + uri

	response = requests.get(uri, cookies=TO_COOKIE, verify=verify)
	if response.status_code != expectedStatus:
		logging.error("Failed to get a response from '%s': server returned status code %d",
		              uri,
		              response.status_code)
		logging.debug("Response: %s\n%r\n%r", response, response.headers, response.content)
		return None

	return response.text

def setStatusFile(statusDir:str, status:str, create:bool = False):
	"""
	Removes all files in ``statusDir`` that aren't ``status``, and creates ``status``
	if it doesn't exist and `create` is True.

	:param statusDir: The absolute path to the directory for Traffic Ops status files
	:param status: The name of the status to be set
	:param create: If :const:`True`, a non-existant ``statusDir/status`` file will be created
	:raises OSError: when reading or writing to the status file fails
	"""
	global MODE

	logging.info("Setting status file")

	if not os.path.isdir(statusDir):
		logging.info("Creating directory %s", statusDir)
		if MODE:
			os.mkdir(statusDir)
	else:
		try:
			statuses = getJSONResponse("/api/1.3/statuses")['response']
		except (KeyError, AttributeError) as e:
			logging.error("Bad API response from /api/1.3/statuses")
			logging.warning("Terminating status file creation/cleanup prematurely.")
			logging.debug("%s", e, exc_info=True, stack_info=True)
			return

		for stat in statuses:
			fname = os.path.join(statusDir, stat["name"])
			if stat != status and os.path.isfile(fname):
				logging.info("Removing %s", fname)
				if not MODE == Modes.REPORT:
					os.remove(fname)

	fname = os.path.join(statusDir, status)
	if create and not os.path.isfile(fname):
		logging.info("creating %s", fname)
		if MODE:
			with open(os.path.join(statusDir, status), "x"):
				pass

#pylint: disable=R1710
def startDaemon(args:typing.List[str], stdout:str='/dev/null', stderr:str='/dev/null') -> bool:
	"""
	Starts a daemon process to execute an external command

	.. note:: this can only indicate the success of the first fork.

	The first fork will exit successfully as long as the second fork doesn't
	raise an OSError. The second fork will exit with the same returncode as
	the exec'd process.

	:param args: The command line to be executed. ``args[0]`` should be the name of the executable
	:param stdout: A filename to which the command's stdout will be re-directed
	:param stderr: A filename to which the command's stderr will be re-directed

	:returns: whether or not the fork succeeded
	"""

	logging.debug("Forking a process - argv: '%s'", ' '.join(args))

	try:
		pid = os.fork()
	except OSError as e:
		logging.error("First fork failed - aborting")
		logging.debug("%s", e, exc_info=True, stack_info=True)
		return False

	if pid:
		# This is the parent
		return True

	# De-couple from parent environment
	os.chdir('/')
	try:
		os.setsid()
	except PermissionError:
		logging.debug("Failure to `setsid`: pid=%d, pgid=%d", os.getpid(), os.getpgid(os.getpid()))
		logging.debug("", exc_info=True, stack_info=True)
	os.umask(0)
	sys.stdin = open('/dev/null')
	sys.stdout = open(stdout, 'w')
	sys.stdout = open(stderr, 'w')

	# Do the second fork magic

	try:
		pid = os.fork()
	except OSError as e:
		logging.error("Error in forked process: %s", e)
		logging.debug("", exc_info=True, stack_info=True)
		exit(1)

	if pid:
		# This is the parent
		exit(0)

	# Now actually exec the program
	try:
		os.execvp(args[0], args[1:])
	except OSError:
		logging.critical("Failure to start %s", args[0])
		logging.debug("", exc_info=True, stack_info=True)

	# If somehow we get down here, it's time to bail
	exit(1)
#pylint: enable=R1710

def setATSStatus(status:bool, restart:bool = False) -> bool:
	"""
	Sets the status of the system's ATS process.

	:param status: Specifies whether ATS should be running (:const:`True`) or not (:const:`False`)
	:param restart: If this is :const:`True`, then ATS will be restarted if it is already running

		.. note:: ``restart`` has no effect if ``status`` is :const:`False`

	:returns: whether or not the status setting was successful (or unnecessary)
	"""
	global MODE, TS_ROOT

	logging.debug("Iterating process list")

	arg = None
	for process in psutil.process_iter():

		# Found an ATS process
		if process.name() == "[TS_MAIN]":
			logging.debug("ATS process found (pid: %d)", process.pid)
			ATSAlreadyRunning = process.status() in {psutil.STATUS_RUNNING, psutil.STATUS_SLEEPING}

			if status and ATSAlreadyRunning and restart:
				logging.info("ATS process found; restarting")
				arg = "restart"
			elif status and ATSAlreadyRunning:
				logging.info("ATS already running; nothing to do.")
			elif status:
				logging.warning("ATS process is running, but status is '%s' - restarting", process.status())
				arg = "restart"
			else:
				logging.warning("ATS is running; stopping ATS")
				arg = "stop"

			break
	else:
		if status:
			logging.warning("ATS not already running; starting ATS.")
			arg = "start"
		else:
			logging.info("ATS already not running; nothing to do.")

	if arg and MODE != "REPORT":
		tsexe = os.path.join(TS_ROOT, "bin", "trafficserver")
		sub = subprocess.Popen([tsexe, arg], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
		out, err = sub.communicate()

		if sub.returncode:
			logging.error("Failed to start trafficserver!")
			logging.warning("Is the 'trafficserver' script located at %s?", tsexe)
			logging.debug(out.decode())
			logging.debug(err.decode())
			return False

	return True

def getHeaderComment() -> str:
	"""
	Gets the header for the Traffic Ops system

	:returns: The ``toolname`` field of the Traffic Ops header
	"""
	response = getJSONResponse("/api/1.3/system/info.json")
	logging.debug("system/info.json response: %s", response)

	if response is None or \
	   "response" not in response or \
	   "parameters" not in response["response"] or \
	   "tm.toolname" not in response["response"]["parameters"]:
		logging.error("Did not find tm.toolname!")
		return ''

	tmToolname = response["response"]["parameters"]["tm.toolname"]
	logging.info("Found tm.toolname: %s", tmToolname)

	return tmToolname

def setTO_LOGIN(login:str) -> str:
	"""
	Parses and returns a JSON-encoded login string for the Traffic Ops API.
	This will set :data:`TO_COOKIE` if login is successful.

	:param login: The raw login info as passed on the command line (e.g. 'username:password')
	:raises PermissionError: if the provided credentials are refused
	:returns: a JSON-encoded login object suitable for passing to the Traffic Ops API's login endpoint
	"""
	global TO_COOKIE

	try:
		login = '{{"u": "{0}", "p": "{1}"}}'.format(*login.split(':'))
	except IndexError as e:
		logging.critical("Bad Traffic_Ops_Login: '%s' - should be like 'username:password'", login)
		raise ORTException()

	logging.debug("TO_LOGIN: %s", login)

	# Obtain login cookie
	cookie = requests.post(TO_URL + '/api/1.3/user/login', data=login, verify=False)

	if not cookie.cookies or 'mojolicious' not in cookie.cookies:
		logging.critical("Login credentials rejected by Traffic Ops")
		logging.error("Response code: %d", cookie.status_code)
		logging.warning("Response Headers: %s", cookie.headers)
		logging.debug("Response: %s", cookie.content)
		raise ORTException()

	TO_COOKIE = {"mojolicious": cookie.cookies["mojolicious"]}

	return login

def getYesNoResponse(prmpt:str, default:str = None) -> bool:
	"""
	Utility function to get an interactive yes/no response to the prompt `prmpt`

	:param prmpt: The prompt to display to users
	:param default: The default response; should be one of ``'y'``, ``"yes"``, ``'n'`` or ``"no"``
		(case insensitive)
	:returns: the parsed response as a boolean
	"""
	if default:
		prmpt = prmpt.rstrip().rstrip(':') + '['+default+"]:"
	while True:
		choice = input(prmpt).lower()

		if choice in {'y', 'yes'}:
			return True
		elif choice in {'n', 'no'}:
			return False
		elif not choice and default is not None:
			return default.lower() in {'y', 'yes'}

		print("Please enter a yes/no response.", file=sys.stderr)

###############################################################################
#####                                                                     #####
#####                         MAIN MODE ROUTINES                          #####
#####                                                                     #####
###############################################################################
def syncDSState() -> bool:
	"""
	Queries Traffic Ops for the Delivery Service's sync state.

	:raises ORTException: All possible errors are coalesced to this with a suitable error message
	:returns: :const:`True` if an update is needed, :const:`False` if it isn't.
	"""
	global HOSTNAME

	logging.info("starting syncDS State fetch")

	try:
		updateStatus = getJSONResponse("/api/1.3/servers/%s/update_status" % HOSTNAME[0])[0]
	except IndexError as e:
		logging.critical("Server not found in Traffic Ops config")
		logging.debug("%s", e, exc_info=True, stack_info=True)
		raise ORTException("Failed to contact API endpoint.")

	try:
		if not updateStatus['upd_pending']:
			logging.info("No update pending.")
			return False

		statusDir = os.path.join(os.path.abspath(os.path.dirname(__file__)), "status")
		setStatusFile(statusDir, updateStatus['status'], create=True)
	except (KeyError, AttributeError) as e:
		logging.critical("Unsupported Traffic Ops version.")
		logging.warning("%s", e, exc_info=True)
		logging.debug("", stack_info=True)
		raise ORTException("Traffic Ops version not supported")

	return True

def revalidate() -> int:
	"""
	Performs revalidation.

	:returns: ``0`` indicates success, ``1`` indicates no revalidation was pending and ``2``
		indicates failure
	"""
	global HOSTNAME

	logging.info("starting revalidation")

	try:
		updateStatus = getJSONResponse("/api/1.3/servers/%s/update_status" % HOSTNAME[0])[0]
	except IndexError as e:
		logging.critical("Server not found in Traffic Ops config")
		logging.debug("%s", e, exc_info=True, stack_info=True)
		return 2

	logging.debug("updateStatus raw response: %s", updateStatus)

	try:
		if not updateStatus['reval_pending']:
			logging.info("No revalidation pending.")
			return 1

		if updateStatus['parent_reval_pending']:
			logging.critical("Parent revalidation is pending.")
			return 1

		statusDir = os.path.join(os.path.abspath(os.path.dirname(__file__)), "status")
		setStatusFile(statusDir, updateStatus['status'])
	except (KeyError, AttributeError):
		logging.critical("Unsupported Traffic Ops version")
		logging.warning("%s", e, exc_info=True)
		logging.debug("%s", e, stack_info=True)
		return 2

	return 0

def updateOps() -> int:
	"""
	Updates Traffic Ops as needed

	:returns: An exit code for the main routine
	"""
	global MODE, HOSTNAME, TO_URL, TO_COOKIE

	revalPending = revalidate() == 0

	if (MODE==Modes.INTERACTIVE and getYesNoResponse("Update Traffic Ops?", 'Y'))\
	   or MODE!=Modes.REVALIDATE:

		logging.info("starting Traffic Ops update for upd_pending")
		payload = {"updated": False, "reval_updated": False}
		response = requests.post(TO_URL+"/update/%s" % HOSTNAME[0],
		                         cookies=TO_COOKIE,
		                         verify=False,
		                         data=payload)

		logging.debug("Raw response from Traffic Ops: %s\n%s\n%s",
		              response,
		              response.headers,
		              response.content)

	elif MODE == Modes.REVALIDATE:
		logging.info("starting Traffic Ops update for reval_pending")
		payload = {"updated": False, "reval_updated": True}
		response = requests.post(TO_URL+"/update/%s" % HOSTNAME[0],
		                         cookies=TO_COOKIE,
		                         verify=False,
		                         data=payload)

		logging.debug("Raw response from Traffic Ops: %s\n%s\n%s",
		              response,
		              response.headers,
		              response.content)

	else:
		logging.warning("Update will not be performed at this time; you should do this manually")
	return 0

#: run modes mapped to handler functions
HANDLERS = {"revalidate": revalidate,
            "syncds": lambda: None}


###############################################################################
#####                                                                     #####
#####                         PACKAGE MANAGEMENT                          #####
#####                                                                     #####
###############################################################################

########################
###      RedHat      ###
########################
def RedHatInstalled(package:str, version:str = None) -> typing.List[str]:
	"""
	RedHat-specific function to check for the existence of a package on the system

	:param package: the package name for which to check
	:param version: an optional version specification

	:returns: the list of packages installed by the name ``package``
	"""
	logging.debug("Checking for RedHat-like package %s", package)
	arg = package if not version else package + '-' + version

	# This should never be done, but CentOS insists on being behind the rest of the
	# world, so done it must be.
	sub = subprocess.Popen(["/bin/rpm", "-q", arg], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
	out, err = sub.communicate()

	if sub.returncode:
		# logging.error("Failed to query rpm database for %s", package)
		logging.debug(out.decode())
		logging.debug(err.decode())
		return []

	return out.decode().split()

def RedHatInstall(packages:typing.List[str]) -> bool:
	"""
	RedHat-specific function to install a list of packages.

	:param packages: A list of package names (optionally concatenated with versions) to be installed
	:returns: whether or not all packages could be successfully installed
	"""
	sub = subprocess.Popen(["/bin/yum", "install", "-y"] + packages,
	                       stdout=subprocess.PIPE,
	                       stderr=subprocess.PIPE)
	out, err = sub.communicate()

	if sub.returncode:
		logging.debug("yum stdout: %s", out.decode())
		logging.debug("yum stderr: %s", err.decode())
		return False

	logging.info("Successfully installed packages: %s", ", ".join(packages))
	return True

def RedHatUninstall(packages:typing.List[str]) -> bool:
	"""
	RedHat-specific function to uninstall a list of packages.

	:param packages: A list of packages to be uninstalled
	:returns: whether or not all packages could be successfully installed
	"""
	sub = subprocess.Popen(["/bin/yum", "remove", "-y"] + packages,
	                       stdout=subprocess.PIPE,
	                       stderr=subprocess.PIPE)
	out, err = sub.communicate()

	if sub.returncode:
		logging.debug("yum stdout: %s", out.decode())
		logging.debug("yum stderr: %s", err.decode())
		return False

	logging.info("Successfully uninstalled packages: %s", ", ".join(packages))
	return True

########################
###      Ubuntu      ###
########################
def UbuntuInstalled(package:str, version:str = None) -> typing.List[str]:
	"""
	Ubuntu-specific function to check for the existence of a package on the system

	:param package: the package name for which to check
	:param version: an optional version specification

	:returns: the list of packages installed by the name ``package``
	"""
	logging.debug("Checking for Ubuntu-like package %s", package)

	sub = subprocess.Popen(["/usr/bin/dpkg", "-l", package],
	                       stdout=subprocess.PIPE,
	                       stderr=subprocess.PIPE)
	out, err = sub.communicate()

	if sub.returncode:
		logging.debug("dpkg stdout: %s", out.decode())
		logging.debug("dpkg stderr: %s", err.decode())
		return []

	try:
		out = [(p.split()[1], p.split()[2]) for p in out[5:].decode().splitlines() if p]

		if version is not None:
			# TODO - better version checking
			out = [p for p in out if version in p[1]]

		return [p[0] for p in out]
	except IndexError:
		logging.warning("Error encountered while processing installed Debian packages")
		logging.debug("Package: %s", package, exc_info=True, stack_info=True)

	return []

def UbuntuInstall(packages:typing.List[str]) -> bool:
	"""
	Ubuntu-specific function to install a list of packages.

	:param packages: A list of package names (optionally concatenated with versions) to be installed
	:returns: whether or not all packages could be successfully installed
	"""
	sub = subprocess.Popen(["/usr/bin/apt-get", "install", "-y"] + packages,
	                       stdout=subprocess.PIPE,
	                       stderr=subprocess.PIPE)
	out, err = sub.communicate()

	if sub.returncode:
		logging.debug("apt-get stdout: %s", out.decode())
		logging.debug("apt-get stderr: %s", err.decode())
		return False

	logging.info("Successfully installed packages: %s", ", ".join(packages))
	return True

def UbuntuUninstall(packages:typing.List[str]) -> bool:
	"""
	Ubuntu-specific function to uninstall a list of packages.

	:param packages: A list of packages to be uninstalled
	:returns: whether or not all packages could be successfully installed
	"""
	sub = subprocess.Popen(["/usr/bin/apt-get", "purge", "-y"] + packages,
	                       stdout=subprocess.PIPE,
	                       stderr=subprocess.PIPE)
	out, err = sub.communicate()

	if sub.returncode:
		logging.debug("apt-get stdout: %s", out.decode())
		logging.debug("apt-get stderr: %s", err.decode())
		return False

	logging.info("Successfully uninstalled packages: %s", ", ".join(packages))
	return True

########################
### Platform Mapping ###
########################
#: Maps distro names to functions that check for the existence of packages on those distros
packageIsInstalled = {'centos': RedHatInstalled,
                      'fedora': RedHatInstalled,
                      'rhel': RedHatInstalled,
                      'ubuntu': UbuntuInstalled,
                      'linuxmint': UbuntuInstalled}

#: Maps distro names to functions that install packages on those distros
installPackage = {'centos': RedHatInstall,
                  'fedora': RedHatInstall,
                  'rhel': RedHatInstall,
                  'ubuntu': UbuntuInstall,
                  'linuxmint': UbuntuInstall}

#: Maps distro names to functions that uninstall packages on those distros
uninstallPackage = {'centos': RedHatUninstall,
                    'fedora': RedHatUninstall,
                    'rhel': RedHatUninstall,
                    'ubuntu': UbuntuUninstall,
                    'linuxmint': UbuntuUninstall}

#: Maps distro names to the symbols that concatenate package names with version semantics for said
#: distros
packageConcat = {'centos': '-',
                 'fedora': '-',
                 'rhel': '-',
                 'ubuntu': '=',
                 'linuxmint': '='}

def processPackages() -> bool:
	"""
	Manages the packages that Traffic Ops reports are required for this server.

	:returns: whether or not the package processing was successfully completed
	"""
	global HOSTNAME, DISTRO, MODE, packageIsInstalled, installPackage, packageConcat

	logging.info("Fetching packages from Traffic Ops")
	packages = getJSONResponse("/ort/%s/packages" % (HOSTNAME[0],))
	logging.info("Response: %s", packages)

	if packages is None:
		raise ConnectionError("Server or server packages not found on server!")

	install, uninstall = {p['name'] + packageConcat[DISTRO] + p['version'] for p in packages}, set()
	for package in packages:
		similarPackages = packageIsInstalled[DISTRO](package['name'])

		# An error occurred, and the package query failed (different than empty response)
		if similarPackages is None:
			logging.critical("Failed to check packages against system.")
			return False

		logging.debug("List of packages similar to %r: %r", package, similarPackages)

		# We check for the pre-existence of packages in two ways; one for
		# RedHat-based distros, and one for Debian/Ubuntu-based distros

		if any(p.startswith(
		          package['name']+packageConcat[DISTRO]+package['version'])
		        for p in similarPackages):
			logging.info("package %s is installed", package['name'])
			p = package['name'] + packageConcat[DISTRO] + package['version']
			logging.info("%s is no longer marked for install", p)
			install.remove(p)
		else:
			logging.info("%r all marked for uninstall", similarPackages)
			uninstall.update(similarPackages)

	logging.info("Marked %d packages for install and %d packages for uninstall.",
	               len(install),               len(uninstall))

	if not MODE == Modes.REPORT:

		logging.info("Installing packages...")
		if install:
			prompt = "Would you like to install the following packages: %s ?" % (", ".join(install))
			if MODE == Modes.INTERACTIVE and not getYesNoResponse(prompt):
				logging.critical("User chose not to install packages - cannot continue!")
				return False

			if not installPackage[DISTRO](list(install)):
				logging.critical("Failed to install packages, possibly permission denied?")
				return False

		logging.info("Done.")


		logging.info("Uninstalling packages...")
		if uninstall:
			prompt = "Would you like to remove the following packages: %s ?" % (", ".join(uninstall))
			if MODE == Modes.INTERACTIVE and not getYesNoResponse(prompt):
				logging.critical("User chose not to remove packages - cannot continue!")
				return False

			if not uninstallPackage[DISTRO](list(uninstall)):
				logging.critical("Failed to uninstall packages, possibly permission denied?")
				return False

		logging.info("Done.")

	return True

def processChkconfig() -> bool:
	"""
	Process the list of services Traffic Ops reports ought to be running on this system

	:returns: whether or not all services could be processed

	.. warning:: This actually only handles Apache Traffic Server at the moment
	"""
	global HOSTNAME, MODE, DISTRO

	logging.info("Processesing Chkconfig...")
	chkconfig = getJSONResponse("/ort/%s/chkconfig" % HOSTNAME[0])

	if chkconfig is None:
		raise ConnectionError("Server or server chkconfig not found on server!")

	logging.debug("chkconfig response: %r", chkconfig)

	for item in chkconfig:
		logging.debug("Processing item: %r", item)

		# A special catch for ATS so we don't need to deal with systemd.
		if item['name'] == "trafficserver":
			if not setATSStatus("on" in item['value']):
				logging.critical("Failed to set ATS Status")
				return False
		else:
			logging.info("checking on service: %s", item['name'])


	return True


###############################################################################
#####                                                                     #####
#####                           CONFIGURATION                             #####
#####                                                                     #####
###############################################################################
def getConfigFiles() -> typing.List[str]:
	"""
	Gets the list of configuration files used by this server's profile

	:returns: The list of configuration files used by this server as reported by Traffic Ops
	"""
	global HOSTNAME

	logging.info("Retrieving configuration files.")

	files = getJSONResponse("/api/1.3/servers/%s/configfiles/ats" % (HOSTNAME[0],))

	if not files:
		logging.critical("Could not retrieve configuration files.")
		return None

	logging.debug("Config Files raw response: %s", files)

	return files

def initBackup() -> bool:
	"""
	Initializes a backup directory as a subdirectory of the directory containing
	this ORT script.

	:returns: whether or not the backup directory was successfully initialized
	"""
	global MODE

	here = os.path.abspath(os.path.dirname(__file__))
	backupdir = os.path.join(here, "backup")

	logging.info("Initializing backup dir %s", backupdir)

	if not os.path.isdir(backupdir):
		if MODE != Modes.REPORT:
			try:
				os.mkdir(backupdir)
			except OSError:
				logging.error("Couldn't create backup dir")
				logging.warning("%s", e)
				logging.debug("", exc_info=True, stack_info=True)
				return False
		else:
			logging.error("Cannot create non-existent backup dir in REPORT mode!")
			return True

	logging.info("Backup dir already exists - nothing to do")
	return True

def mkbackup(fname:str, contents:str) -> bool:
	"""
	Creates a backup of a specific file with its original contents.

	:param fname: The name of the file being backed up (basename only)
	:param contents: The file's original contents
	:returns: whether or not the backup succeeded

		.. note: will always return :const:`True` in :attr:`Modes.REPORT` mode.
	"""
	global MODE

	if MODE == Modes.REPORT:
		logging.info("REPORT mode - nothing to do")
		return True

	backupfile = os.path.join(os.path.abspath(os.path.dirname(__file__)), "backup", fname)
	if os.path.isfile(backupfile):
		logging.warning("Clobbering existing backup file '%s'!", backupfile)

	try:
		with open(backupfile, 'w') as fd:
			fd.write(contents)
	except OSError as e:
		logging.warning("Failed to write backup file: %s", e)
		logging.debug("", exc_info=True, stack_info=True)
		return False

	logging.info("Backup of %s written to %s", fname, backupfile)
	return True

def stripComments(s) -> str:
        """
        Strips comments from a string
        """
        return re.sub(r'(?m)^\s*^ *<!--.*\n?', '', s)

def updateConfig(directory:str, fname:str, contents:str) -> bool:
	"""
	Updates a single configuration file

	This will make a backup in the `backup` subdirectory of the directory
	containing this script if the file on disk differs from ``contents``.

	:param directory: The directory which contains the configuration file
	:param fname: The basename of the configuration file
	:param contents: The contents which this file will hold
	:returns: whether or not the update was successful

		.. note:: "Success" is defined as being able to write the file contents and create any
			necessary backups if the run mode is not :attr:`Modes.BADASS` - in which case failure
			to back the file up is 'acceptable'.

	"""
	global MODE, ATS_FILES, ATS_NEEDS_RESTART

	file = os.path.join(directory, fname)

	logging.info("Updating config file '%s'", file)

	if not os.path.isfile(file):
		if MODE != Modes.REPORT:
			logging.info("File does not exist - creating")
			try:
				with open(file, 'w') as fd:
					fd.write(contents)
			except OSError as e:
				logging.error("Couldn't write to file")
				logging.warning("%s", e)
				logging.debug("", exc_info=True, stack_info=True)
				return False

			logging.info("File written.")

		return True

	logging.debug("Reading in file on disk")
	try:
		with open(file) as fd:
			diskContents = ''.join([line for line in fd.readlines()\
			                if line and not line.startswith("# DO NOT EDIT")\
			                and not line.startswith("# TRAFFIC OPS NOTE:")])
	except OSError:
		logging.warning("Couldn't read on-disk file: %s", e)
		logging.debug("", exc_info=True, stack_info=True)
		if MODE != Modes.BADASS:
			return False

	# Certain headers from Traffic Ops shouldn't be considered when checking for equality
	importantContents = '\n'.join([l for l in contents.split('\n')\
	                     if l and not l.startswith("# DO NOT EDIT")\
	                     and not l.startswith("# TRAFFIC OPS NOTE:")])

	if stripComments(diskContents.strip()) == stripComments(contents.strip()):
		logging.info("on-disk contents match Traffic Ops - nothing to do")
		return True

	if not mkbackup(os.path.basename(file), diskContents) and MODE != Modes.BADASS:
		logging.error("Failed to create backup.")
		return False

	try:
		with open(file, 'w') as fd:
			fd.write(contents + '\n' if not contents.endswith('\n') else contents)
	except OSError as e:
		logging.error("Failed to update config file '%s'", file)
		logging.warning("%s", e)
		logging.debug("", exc_info=True, stack_info=True)
		return False

	# If update was needed and successful, then check if ats should be restarted
	if fname in ATS_FILES:
		ATS_NEEDS_RESTART = True

	logging.info("%s has been updated", file)
	return True

def sanitizeContents(contents:str) -> str:
	"""
	Sanitizes the input `contents` string to be a well-behaved config file.
	"""
	out = []
	for line in contents.splitlines():
		tmp=(" ".join(line.split())).strip() #squeezes spaces and trims leading and trailing spaces
		tmp=tmp.replace("&amp;", '&') #decodes HTML-encoded ampersands
		tmp=tmp.replace("&gt;", '>') #decodes HTML-encoded greater-than symbols
		tmp=tmp.replace("&lt;", '<') #decodes HTML-encoded less-than symbols
		out.append(tmp)

	return "\n".join(out)

def processConfigFile(file:dict, port:int, ip:str) -> bool:
	"""
	Process a given configuration file object to produce the specified contents

	:param file: An object representing a configuration file
	:param port: This server's TCP port
	:param ip: This server's IPv4 address
	:returns: whether or not the update was successful
	"""
	global MODE, HOSTNAME, TO_URL

	try:
		fname = file['fnameOnDisk']
		scope = file['scope']
		location = file['location']
		uri, contents = None, None
		if 'apiUri' in file:
			uri = TO_URL + file['apiUri']
		elif 'url' in file:
			uri = file['url']
		else:
			contents = file['contents']
	except KeyError as e:
		logging.error("Malformed config file")
		logging.warning("%s", e)
		logging.debug("", exc_info=True, stack_info=True)
		return False

	logging.info("======== Start processing config file: %s ========", fname)

	if not os.path.isdir(location):
		if MODE != Modes.REPORT:
			logging.debug("location '%s' doesn't exist; creating.")

			try:
				os.makedirs(location)
			except OSError as e:
				logging.error("Couldn't create directory %s", location)
				logging.warning("%s", e)
				logging.debug("", exc_info=True, stack_info=True)
				return False
		else:
			# Even though nothing was done and an error gets reported, we return a success here
			# because presumably everything would go fine were this not REPORT mode.
			logging.error("Cannot create dirs in REPORT mode!")
			return True

	# `contents` should only be `None` if `uri` is defined
	if contents is None:
		contents = getRawResponse(uri, TOrelative=False)

	if contents is None: #still...
		return False

	contents = contents.replace("__HOSTNAME__", HOSTNAME[0])
	contents = contents.replace("__FULL_HOSTNAME__", HOSTNAME[1])
	contents = contents.replace("__RETURN__", '\n')
	contents = contents.replace("__CACHE_IPV4__", ip)

	# Don't ask me why, but the reference ORT implementation just strips these ones out
	# if the tcp port is 80.
	contents = contents.replace("__SERVER_TCP_PORT__", str(port) if port != 80 else "")

	contents = sanitizeContents(contents)

	logging.debug("Sanitized file contents from Traffic Ops database:\n%s\n" % contents)

	logging.warning("Skipping pre-requisite checks - plugins may not be installed!!")
	logging.info("Dependent packages should be specified as profile parameters.")

	updateConfig(location, fname, contents)

	logging.info("======== End processing config file: %s ========", fname)

	return True

def processConfigFiles(files:list, port:int, ip:str) -> bool:
	"""
	Processes the passed JSON object containing configuration file definitions

	:param files: A list of configuration file objects
	:param port: This server's TCP port
	:param ip: This server's IPv4 address
	:returns: whether or not all configuration files could be processed
	"""
	global MODE

	if not initBackup():
		return False

	for file in files:
		if not processConfigFile(file, port, ip):
			logging.error("Failed to process config file '%s'", file)

			if MODE != Modes.BADASS:
				return False

			logging.warning("We're in BADASS mode, attempting to continue")

	return True


###############################################################################
#####                                                                     #####
#####                             MAIN FLOW                               #####
#####                                                                     #####
###############################################################################
def doMain() -> int:
	"""
	Performs operations based on the run mode.
	This can be thought of as the "true" main function.

	:returns: an exit code for the script
	"""
	global MODE, ATS_NEEDS_RESTART

	headerComment = getHeaderComment()

	try:
		DSUpdateNeeded = syncDSState()
	except ORTException as e:
		logging.critical("Failed to get Update Pending from Traffic Ops: %s", e)
		return 1

	myFiles = getConfigFiles()
	if myFiles is None:
		return 1
	if "configFiles" not in myFiles or "info" not in myFiles:
		logging.critical("Malformed response from configfiles/ats endpoint - unable to continue")
		return 1

	try:
		# This is needed for templated responses from the config file API endpoints
		tcpPort = myFiles["info"]["serverTcpPort"]
		serverIpv4 = myFiles["info"]["serverIpv4"]
	except KeyError:
		logging.critical("Malformed response from configfiles/ats endpoint - unable to continue")
		logging.debug("", exc_info=True, stack_info=True)
		return 1


	try:
		if MODE == Modes.REVALIDATE:
			logging.info("======== Revalidating, no package processing needed ========")
			reval = revalidate()
			if reval:
				# Bail, possibly with an exit code
				return reval - 1

		else:
			logging.info("======== Start processing packages ========")
			if not processPackages():
				logging.critical("Unrecoverable error occurred when processing packages.")
				return 1
			logging.info("======== Start processing services ========")
			if not processChkconfig():
				logging.critical("Unrecoverable error occurred when processing services.")
				return 1

		if not processConfigFiles(myFiles["configFiles"], tcpPort, serverIpv4):
			logging.critical("Unrecoverable error occurred when processing config files")
			return 1
	except ConnectionError:
		logging.critical("Couldn't reach /ort/%s/* endpoints - "\
		                 "ensure %s points to Traffic OPS not Traffic PORTAL",
		                 HOSTNAME[0],
		                 TO_URL)
		return 1

	if ATS_NEEDS_RESTART:
		logging.info("Restarting ATS")
		if not setATSStatus(True, restart=True):
			logging.critical("Failed to restart ATS")
			return 1
	else:
		logging.info("ATS Restart unnecessary")

	if MODE != Modes.REPORT and DSUpdateNeeded:
		return updateOps()

	logging.info("Ops update unnecessary; exiting.")

	return 0

def main() -> int:
	"""
	Runs the program

	:returns: an exit code for the script
	"""
	global TO_COOKIE, TO_LOGIN, TO_URL, LOG_LEVELS, MODE, needInstall, DISTRO, TS_ROOT, FMT

	# I have no idea why, but the old ORT script does this on every run.
	print(datetime.datetime.utcnow().strftime("%a %b %d %H:%M:%S UTC %Y"))

	parser = argparse.ArgumentParser(description="A Python-based TO_ORT implementation",
	                                 epilog="Doesn't support the 'TRACE' or 'NONE' log levels.",
	                                 formatter_class=argparse.ArgumentDefaultsHelpFormatter)

	parser.add_argument("Mode",
	                    help="REPORT: Do nothing, but print what would be done\n"\
	                         "")
	parser.add_argument("Log_Level",
	                    help="ALL, DEBUG, INFO, WARN, ERROR, FATAL",
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
	parser.add_argument("--rev_proxy_disabled",
	                    help="bypass the reverse proxy even if one has been configured.",
	                    type=int,
	                    default=0)
	parser.add_argument("--ts_root",
	                    help="Specify the root directory at which Apache Traffic Server is installed"\
	                         " (e.g. '/opt/trafficserver')",
	                    type=str,
	                    default="/")
	parser.add_argument("-v", "--version",
	                    action="version",
	                    version="%(prog)s v"+__version__,
	                    help="Print version information and exit.")

	args = parser.parse_args()

	logLevel = args.Log_Level.upper()
	if logLevel not in LOG_LEVELS:
		print("Unrecognized/Unsupported log level:", args.Log_Level, file=sys.stderr)
		return 1

	logging.basicConfig(level=LOG_LEVELS[logLevel], format=FMT)
	logging.getLogger().setLevel(LOG_LEVELS[logLevel])

	try:
		MODE = Modes[args.Mode.upper()]
	except KeyError as e:
		logging.critical("Unknown mode '%s'", args.Mode)
		logging.debug("%s", e, exc_info=True, stack_info=True)
		return 1

	logging.info("Running in %s mode", MODE)

	if needInstall:
		if not handleMissingPythonPackages(needInstall):
			return 1

	tsroot = args.ts_root.strip()
	tsbin = os.path.join(tsroot, "bin", "traffic_server")
	if not os.path.isdir(tsroot) or not os.path.isfile(tsbin):
		logging.critical("Unable to find root Apache Traffic Server installation")
		logging.error("No such directory: '%s' or no such file: '%s", tsroot, tsbin)
		return 1

	TS_ROOT = args.ts_root
	logging.info(\
	        "Traffic Server root is at %s\n\tDistro detected as '%s'\n\tHostname detected as '%s'",
	        TS_ROOT,
	        DISTRO,
	        HOSTNAME[1])

	TO_URL = args.Traffic_Ops_URL.rstrip('/')

	logging.info("Traffic_Ops_URL: %s", TO_URL)

	if not TO_URL.startswith("http"):
		logging.critical("Malformed Traffic_Ops_URL: '%s' - "\
		                 "must include scheme (e.g. http://traffic.ops)", TO_URL)
		return 1

	# Litmus test to make sure the server exists and can be reached
	try:
		_ = requests.head(TO_URL, verify=False)
	except requests.exceptions.RequestException as e:
		logging.critical("Malformed or Invalid Traffic_Ops_URL")
		logging.error("%s", e)
		logging.debug("%s", e, exc_info=True, stack_info=True)
		return 1

	try:
		TO_LOGIN = setTO_LOGIN(args.Traffic_Ops_Login)
	except ORTException:
		return 1

	return doMain()


if __name__ == '__main__':
	try:
		exit(main())
	except requests.exceptions.ConnectionError as e:
		print("Failed to connect to Traffic Ops:", e, file=sys.stderr)
		exit(1)
