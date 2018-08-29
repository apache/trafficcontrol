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
FMT = "%(levelname)s: %(filename)s line %(lineno)d in %(funcName)s: %(message)s"

HOSTNAME = platform.node().split('.')[0] # Not strictly accurate, but generally good enough

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

# Current Run Mode
MODE = None


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
		globals()[DISTRO] = distro.LinuxDistribution().id()

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

def getConfigFiles() -> typing.List[str]:
	"""
	Gets the list of configuration files used by this server's profile
	"""
	global HOSTNAME

	logging.info("Retrieving configuration files.")

	files = getJSONResponse("/api/1.3/servers/%s/configfiles/ats" % (HOSTNAME,))

	if not files:
		logging.critical("Could not retrieve configuration files.")
		return 1

	logging.debug("Config Files raw response: %s", files)

	return 0

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

	if not pid:
		# This is the parent
		return True

	# De-couple from parent environment
	os.chdir('/')
	os.setsid()
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

	if not pid:
		# This is the parent
		exit(0)

	# Now actually exec the program
	sub = subprocess.Popen(args, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
	out, err = sub.communicate()

	if sub.returncode:
		logging.debug("stdout: %s", out.decode())
		logging.debug("stderr: %s", err.decode())

	exit(sub.returncode)

def setATSStatus(status:bool) -> bool:
	"""
	Sets the status of the system's ATS process to on if `status` is True, else off.

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

			if status and ATSAlreadyRunning:
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


###############################################################################
#####                                                                     #####
#####                         MAIN MODE ROUTINES                          #####
#####                                                                     #####
###############################################################################
def syncDSState() -> int:
	"""
	Queries Traffic Ops for the Delivery Service's sync state
	"""
	global HOSTNAME

	logging.info("starting syncDS State fetch")

	try:
		updateStatus = getJSONResponse("/api/1.3/servers/%s/update_status" % HOSTNAME)[0]
	except IndexError:
		logging.critical("Server not found in Traffic Ops config")
		return 1

	try:
		if not updateStatus['upd_pending']:
			logging.info("No update pending.")
			return 0

		statusDir = os.path.join(os.path.abspath(os.path.dirname(__file__)), "status")
		setStatusFile(statusDir, updateStatus['status'], create=True)
	except (KeyError, AttributeError) as e:
		logging.critical("Unsupported Traffic Ops version.")
		logging.error("%s", e)
		logging.warning("%s", e, exc_info=True)
		logging.debug("%s", e, stack_info=True)
		return 1

	return 0

def revalidate() -> int:
	"""
	Performs revalidation.
	"""
	global HOSTNAME

	logging.info("starting revalidation")

	try:
		updateStatus = getJSONResponse("/api/1.3/servers/%s/update_status" % HOSTNAME)[0]
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
### Platform Mapping ###
########################
packageIsInstalled = {'centos': RedHatInstalled,
                      'fedora': RedHatInstalled,
                      'rhel': RedHatInstalled}
installPackage = {'centos': RedHatInstall,
                  'fedora': RedHatInstall,
                  'rhel': RedHatInstall}
uninstallPackage = {'centos': RedHatUninstall,
                    'fedora': RedHatUninstall,
                    'rhel': RedHatUninstall}
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
	packages = getJSONResponse("/ort/%s/packages" % (HOSTNAME,))
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
		if install and not installPackage[DISTRO](list(install)):
			logging.critical("Failed to install packages, possibly permission denied?")
			return False

		logging.info("Done.")
		logging.info("Uninstalling packages...")
		if uninstall and not uninstallPackage[DISTRO](list(uninstall)):
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
	chkconfig = getJSONResponse("/ort/%s/chkconfig" % HOSTNAME)

	if chkconfig is None:
		raise ConnectionError("Server or server chkconfig not found on server!")

	logging.info("chkconfig response: %r", chkconfig)

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
#####                             MAIN FLOW                               #####
#####                                                                     #####
###############################################################################
def doMain() -> int:
	"""
	Performs operations based on the run mode.
	This can be thought of as the "true" main function.
	"""
	global MODE

	headerComment = getHeaderComment()
	myFiles = getConfigFiles()

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
	except ConnectionError:
		logging.critical("Couldn't reach /ort/%s/* endpoints - "\
		                 "ensure %s points to Traffic OPS not Traffic PORTAL",
		                 HOSTNAME,
		                 TO_URL)
		return 1

	myFiles = getConfigFiles()
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

	logging.info("Hostname detected as '%s'", HOSTNAME)

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
