#!/usr/bin/env python3

"""
This script aims to be a drop-in replacement for the aged
`traffic_ops_ort.pl` script. Its primary purpose is for
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

# Holds the info needed to query Traffic Ops
TO_URL, TO_LOGIN, TO_COOKIE, TS_ROOT = (None,)*4

# Logging info
LOG_LEVELS = {"ALL":   logging.NOTSET,
              "DEBUG": logging.DEBUG,
              "INFO":  logging.INFO,
              "WARN":  logging.WARNING,
              "ERROR": logging.ERROR,
              "FATAL": logging.CRITICAL}
FMT = "%(levelname)s: line %(lineno)d in %(module)s.%(funcName)s: %(message)s"

# Not strictly accurate, but generally good enough
HOSTNAME = (platform.node().split('.')[0], platform.node())

class Modes(enum.IntEnum):
	"""
	Enumerated run modes
	"""
	REPORT = 0
	INTERACTIVE = 1
	REVALIDATE = 2
	SYNCDS = 3
	BADASS = 4

	def __str__(self) -> str:
		"""
		Implements `str(self)` by returning enum member's name
		"""
		return self.name

class ORTException(Exception):
	"""Represents an error while processing ORT API responses, etc."""
	pass


# Current Run Mode
MODE = None


#This is the set of files which will require an ATS restart upon changes
ATS_FILES = {"records.config",
             "remap.config",
             "parent.config",
             "cache.config",
             "hosting.config",
             "astats.config",
             "logs_xml.config",
             "ssl_multicert.config"}
ATS_NEEDS_RESTART = False

###############################################################################
#####                                                                     #####
#####                     PYTHON DEPENDENCY HANDLING                      #####
#####                                                                     #####
###############################################################################
def installPythonPackages(packages:typing.List[str]) -> bool:
	"""
	Attempts to install the packages listed in `packages` and
	add them to the global scope.

	Returns a truthy value indicating success.
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
	Handles the case of missing packages.

	Installs packages by calling `installPythonPackages`.
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
def getJSONResponse(uri:str, expectedStatus:int = 200) -> object:
	"""
	Returns the JSON-encoded contents (as a `dict`) of the response for a GET request to `uri`
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

def getRawResponse(uri:str, expectedStatus:int=200, TOrelative:bool=True, verify:bool=False) ->str:
	"""
	Returns the raw body of a GET request for the specified URI.
	(actually encodes to utf-8 string)

	Note that the behaviour of treating the uri as relative to TO_URL may be overridden, unlike
	`getJSONResponse`.
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
	Removes all files in `statusDir` that aren't `status`, and creates `status`
	if it doesn't exist and `create` is True.

	Raises an OSError if that fails.
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
			fname = os.path.join(statusDir, stat)
			if stat != status and os.path.isfile(fname):
				logging.info("Removing %s", fname)
				if not MODE == Modes.REPORT:
					os.remove(os.path.join(statusDir, stat))

	fname = os.path.join(statusDir, status)
	if create and not os.path.isfile(fname):
		logging.info("creating %s", fname)
		if MODE:
			with open(os.path.join(statusDir, status), "x"):
				pass
#pylint: disable=R1710
def startDaemon(args:typing.List[str], stdout:str='/dev/null', stderr:str='/dev/null') -> bool:
	"""
	Starts a daemon process to execute the command line given by 'args'
	and returns a boolean indicating success.

	Note that this can only indicate the success of the first fork.

	The first fork will exit successfully as long as the second fork doesn't
	raise an OSError. The second fork will exit with the same returncode as
	the exec'd process.
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
	Sets the status of the system's ATS process to on if `status` is True, else off.
	If `restart` is True, then ATS will be restarted if already running.
	(`restart` has no effect if `status` is False)

	Returns a boolean indicator of success.
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
		return startDaemon([os.path.join(TS_ROOT, "bin", "traffic_server"), arg],
		                   stdout=os.path.join(TS_ROOT, "var", "log", "trafficserver", "traffic.out"),
		                   stderr=os.path.join(TS_ROOT, "var", "log", "trafficserver", "error.log"))

	return True

def getHeaderComment() -> str:
	"""
	Gets the header for the Traffic Ops system
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
	Parses the passed login and returns a
	JSON string used for login.

	Will test the login before returning, and
	raise a `PermissionError` if credentials
	are refused. Also sets the TO_COOKIE
	global variable.
	"""
	global TO_COOKIE

	login = '{{"u": "{0}", "p": "{1}"}}'.format(*login.split(':'))

	logging.debug("TO_LOGIN: %s", login)

	# Obtain login cookie
	cookie = requests.post(TO_URL + '/api/1.3/user/login', data=login, verify=False)

	if not cookie.cookies or 'mojolicious' not in cookie.cookies:
		logging.error("Response code: %d", cookie.status_code)
		logging.warning("Response Headers: %s", cookie.headers)
		logging.debug("Response: %s", cookie.content)
		raise PermissionError("Login credentials rejected by Traffic Ops")

	TO_COOKIE = {"mojolicious": cookie.cookies["mojolicious"]}

	return login

def getYesNoResponse(prmpt:str, default:str = None) -> bool:
	"""
	Utility function to get an interactive yes/no response to the prompt `prmpt`
	"""
	if default:
		prmpt = prmpt.rstrip().rstrip(':') + '['+default+"]:"
	while True:
		choice = input(prmpt).lower()

		if choice in {'y', 'yes'}:
			return True
		elif choice in {'n', 'no'}:
			return False

		print("Please enter a yes/no response.", file=sys.stderr)

###############################################################################
#####                                                                     #####
#####                         MAIN MODE ROUTINES                          #####
#####                                                                     #####
###############################################################################
def syncDSState() -> bool:
	"""
	Queries Traffic Ops for the Delivery Service's sync state.

	Returns True if an update is needed, False if it isn't.
	If something goes wrong, it'll raise an ORTException containing the error
	message.
	"""
	global HOSTNAME

	logging.info("starting syncDS State fetch")

	try:
		updateStatus = getJSONResponse("/api/1.3/servers/%s/update_status" % HOSTNAME[0])[0]
	except IndexError:
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
	"""
	global HOSTNAME

	logging.info("starting revalidation")

	try:
		updateStatus = getJSONResponse("/api/1.3/servers/%s/update_status" % HOSTNAME[0])[0]
	except IndexError:
		logging.critical("Server not found in Traffic Ops config")
		return 1

	logging.debug("updateStatus raw response: %s", updateStatus)

	try:
		if not updateStatus['reval_pending']:
			logging.info("No revalidation pending.")
			return 0

		if updateStatus['parent_reval_pending']:
			logging.critical("Parent revalidation is pending.")
			return 1

		statusDir = os.path.join(os.path.abspath(os.path.dirname(__file__)), "status")
		setStatusFile(statusDir, updateStatus['status'])
	except (KeyError, AttributeError):
		logging.critical("Unsupported Traffic Ops version")
		logging.warning("%s", e, exc_info=True)
		logging.debug("%s", e, stack_info=True)
		return 1

	return 0

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
	Returns the list of packages installed by the name 'package',
	optionally with a specific version.
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
	Installs the packages in the `packages` list and
	returns a boolean success indicator.
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
	Removes the packages in the `packages` list, and
	returns a boolean indicator of success.
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
	Returns the list of packages installed by the name 'package',
	optionally with a specific version.
	"""
	logging.debug("Checking for Ubuntu-like package %s", package)

	sub = subprocess.Popen(["/usr/bin/dpkg", "-l", package], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
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
	Installs the packages in the `packages` list and
	returns a boolean success indicator
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
	Uninstalls the packages in the `packages` list and returns a
	boolean indicator of success.
	"""
	sub = subprocess.Popen(["/usr/bin/apt-get", "purge", "-y"] + packages,
	                       stdout=subprocess.PIPE,
	                       stderr=subprocess.PIPE)
	out, err = subprocess.communicate()

	if sub.returncode:
		logging.debug("apt-get stdout: %s", out.decode())
		logging.debug("apt-get stderr: %s", err.decode())
		return False

	logging.info("Successfully uninstalled packages: %s", ", ".join(packages))
	return True

########################
### Platform Mapping ###
########################
packageIsInstalled = {'centos': RedHatInstalled,
                      'fedora': RedHatInstalled,
                      'rhel': RedHatInstalled,
                      'ubuntu': UbuntuInstalled,
                      'linuxmint': UbuntuInstalled}
installPackage = {'centos': RedHatInstall,
                  'fedora': RedHatInstall,
                  'rhel': RedHatInstall,
                  'ubuntu': UbuntuInstall,
                  'linuxmint': UbuntuInstall}
uninstallPackage = {'centos': RedHatUninstall,
                    'fedora': RedHatUninstall,
                    'rhel': RedHatUninstall,
                    'ubuntu': UbuntuUninstall,
                    'linuxmint': UbuntuUninstall}
packageConcat = {'centos': '-',
                 'fedora': '-',
                 'rhel': '-',
                 'ubuntu': '=',
                 'linuxmint': '='}

def processPackages() -> bool:
	"""
	Manages the packages that Traffic Ops reports are required for this server.

	Returns a boolean indication of success.
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

		if MODE == Modes.INTERACTIVE
		logging.info("Installing packages...")
		if install:
			if MODE == INTERACTIVE and not getYesNoResponse("Would you like to install the following packages: %s ?"%", ".join(install)):
				logging.critical("User chose not to install packages - cannot continue!")
				return False

			if not installPackage[DISTRO](list(install)):
				logging.critical("Failed to install packages, possibly permission denied?")
				return False

		logging.info("Done.")


		logging.info("Uninstalling packages...")
		if uninstall:
			if MODE == INTERACTIVE and not getYesNoResponse("Would you like to remove the following packages: %s ?"%", ".join(uninstall)):
				logging.critical("User chose not to remove packages - cannot continue!")
				return False

			if not uninstallPackage[DISTRO](list(uninstall)):
				logging.critical("Failed to uninstall packages, possibly permission denied?")
				return False

		logging.info("Done.")

	return True

def processChkconfig() -> bool:
	"""
	'Processes chkconfig' - whatever that means/is

	Returns a boolean indicator of success.
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
	Creates a backup file named 'fname' with the contents `contents` and returns
	True if the operation succeeded, else False.

	Note: will always return True in REPORT mode.
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

def updateConfig(directory:str, fname:str, contents:str) -> bool:
	"""
	Updates the config file specified by `file` to contain `contents`.

	Returns a boolean indicator of success.
	"Success" is defined as being able to write the file contents and create any necessary backups
	if the mode is not BADASS - in which case failure to back the file up is 'acceptable'.

	This will make a backup in the `backup` subdirectory of the directory
	containing this script if the file on disk differs from `contents`.
	"""
	global MODE, ATS_FILES, ATS_NEEDS_RESTART

	file = os.path.join(directory, fname)

	logging.info("Udpating config file '%s'", file)

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
			diskContents = fd.read()
	except OSError:
		logging.warning("Couldn't read on-disk file: %s", e)
		logging.debug("", exc_info=True, stack_info=True)
		if MODE != Modes.BADASS:
			return False

	if diskContents.strip() == contents.strip():
		logging.info("on-disk contents match Traffic Ops - nothing to do")
		return True

	if not mkbackup(os.path.basename(file), diskContents) and MODE != Modes.BADASS:
		logging.error("Failed to create backup.")
		return False

	try:
		with open(file, 'w') as fd:
			fd.write(contents)
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
	Process the passed file and value to apply a configuration.

	Returns a boolean indicator of success
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
	processes the passed JSON object containing config file definitions

	Returns a boolean indicator of success
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


	if MODE == Modes.REVALIDATE:
		logging.info("======== Revalidating, no package processing needed ========")
		return revalidate()

	logging.info("======== Start processing packages ========")
	try:
		if not processPackages():
			logging.critical("Unrecoverable error occured when processing packages.")
			return 1
		logging.info("======== Start processing services ========")
		if not processChkconfig():
			logging.critical("Unrecoverable error occured when processing chkconfig.")
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

	return 0

def main() -> int:
	"""
	Runs the program, returns an exit code
	"""
	global TO_COOKIE, TO_LOGIN, TO_URL, LOG_LEVELS, MODE, needInstall, DISTRO, TS_ROOT, FMT

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
	parser.add_argument("-dispersion",
	                    help="wait a random number between 0 and <dispersion> before starting.",
	                    type=int,
	                    default=300)
	parser.add_argument("-login_dispersion",
	                    help="wait a random number between 0 and <login_dispersion> before login.",
	                    type=int,
	                    default=0)
	parser.add_argument("-retries",
	                    help="retry connection to Traffic Ops URL <retries> times.",
	                    type=int,
	                    default=3)
	parser.add_argument("-wait_for_parents",
	                    help="do not update if parent_pending = 1 in the update json.",
	                    type=int,
	                    default=1)
	parser.add_argument("-rev_proxy_disabled",
	                    help="bypass the reverse proxy even if one has been configured.",
	                    type=int,
	                    default=0)
	parser.add_argument("-ts_root",
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

	logging.debug("test")
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
	logging.info("Traffic Server root is at %s", TS_ROOT)

	logging.info("Distro detected as '%s'", DISTRO)

	logging.info("Hostname detected as '%s'", HOSTNAME[1])

	TO_URL = args.Traffic_Ops_URL.rstrip('/')

	logging.info("Traffic_Ops_URL: %s", TO_URL)

	if not TO_URL.startswith("http"):
		logging.critical("Malformed Traffic_Ops_URL: '%s' - "\
		                 "must include scheme (e.g. http://traffic.ops)", TO_URL)
		return 1

	# Litmus test to make sure the server exists and can be reached
	_ = requests.head(TO_URL, verify=False)

	try:
		TO_LOGIN = setTO_LOGIN(args.Traffic_Ops_Login)
	except IndexError:
		logging.critical("Bad Traffic_Ops_Login: '%s' - should be like 'username:password'")
		return 1
	except PermissionError:
		logging.critical("Failed to obtain cookied from Traffic Ops")
		return 1

	return doMain()


if __name__ == '__main__':
	try:
		exit(main())
	except requests.exceptions.ConnectionError as e:
		print("Failed to connect to Traffic Ops:", e, file=sys.stderr)
		exit(1)
