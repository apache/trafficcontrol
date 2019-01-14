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
This module is responsible for holding information related to
the configuration of the ORT script; it has constants that
hold and set up the log level, run modes, Traffic Ops login
credentials etc.
"""

import enum
import logging
import os
import platform

import distro
import requests

#: Contains the host's hostname as a tuple of ``(short_hostname, full_hostname)``
HOSTNAME = (platform.node().split('.')[0], platform.node())

#: contains identifying information about the host system's Linux distribution
DISTRO = distro.LinuxDistribution().id()

#: Holds information about the host system, required for processing configuration files,
#: and also possibly useful in other situations
SERVER_INFO = None

#: This sets whether or not to verify SSL certificates when communicated with Traffic Ops.
#: Does not affect non-Traffic Ops servers
VERIFY = True

#: If set to :const:`True`, this script will not apply updates until all of its parents have
#: finished applying their updates
WAIT_FOR_PARENTS = False


class Modes(enum.IntEnum):
	"""
	Enumerated run modes
	"""
	REPORT = 0      #: Do nothing, only report what would be done
	INTERACTIVE = 1 #: Ask for user confirmation before modifying the system
	REVALIDATE = 2  #: Only check for configuration file changes and content revalidations
	SYNCDS = 3      #: Check for and apply Delivery Service changes
	BADASS = 4      #: Apply all settings specified in Traffic Ops, and attempt to solve all problems

	def __str__(self) -> str:
		"""
		Implements ``str(self)`` by returning enum member's name
		"""
		return self.name

#: Holds the current run mode
MODE = None

def setMode(mode:str) -> bool:
	"""
	Sets the script's run mode in the global variable :data:`MODE`

	:param mode: Expected to be the name of a :class:`Modes` constant.
	:returns: whether or not the run mode could be set successfully
	:raises ValueError: when ``mode`` is not a :const:`str`
	"""
	try:
		mode = Modes[mode.upper()]
	except KeyError:
		return False
	except (AttributeError, ValueError) as e:
		raise ValueError from e

	global MODE
	MODE = mode

	return True


class LogLevels(enum.IntEnum):
	"""
	Enumerated Log levels
	"""
	ALL      = logging.NOTSET     #: Outputs all logging information
	TRACE    = logging.NOTSET     #: Synonym for :attr:`ALL`
	DEBUG    = logging.DEBUG      #: Outputs debugging information as well as all higher log levels
	INFO     = logging.INFO       #: Outputs informational messages as well as all higher log levels
	WARN     = logging.WARNING    #: Outputs warnings, errors and fatal messages
	ERROR    = logging.ERROR      #: Errors - but not more verbose warnings - will be output
	FATAL    = logging.CRITICAL   #: Outputs only reasons why the script exited prematurely
	CRITICAL = logging.CRITICAL   #: Synonym for :attr:`FATAL`
	NONE     = logging.CRITICAL+1 #: Silent mode - no output at all

	def __str__(self) -> str:
		"""
		Implements ``str(self)`` by returning the enum member's name.
		(Coalesces synonyms to their legacy names)
		"""
		return self.name if self != logging.CRITICAL else "FATAL"

#: A format specifier for logging output. Propagates to all imported modules.
LOG_FORMAT = "%(levelname)s: %(asctime)s line %(lineno)d in %(module)s.%(funcName)s: %(message)s"

def setLogLevel(level:str) -> bool:
	"""
	Sets the global logger's log level to the desired name.

	:param level: Expected to be the name of a :class:`LogLevels` constant.
	:returns: whether or not the log level could be set successfully
	:raises ValueError: when the type of ``level`` is not :const:`str`
	"""
	try:
		level = LogLevels[level.upper()]
	except KeyError:
		return False
	except (AttributeError, ValueError) as e:
		raise ValueError from e

	logging.basicConfig(level=level, format=LOG_FORMAT)
	logging.getLogger().setLevel(level)

	return True


#: An absolute path to the root installation directory of the Apache Trafficserver installation
TS_ROOT = None

def setTSRoot(tsroot:str) -> bool:
	"""
	Sets the global variable :data:`TS_ROOT`.

	:param tsroot: Should be an absolute path to the directory containing the system's Apache
		Trafficserver installation.
	:returns: whether or not the installation path could be set successfully
	:raises ValueError: if ``tsroot`` is not a :const:`str`

	"""
	try:
		tsroot = tsroot.strip()

		if tsroot != '/' and tsroot.endswith('/'):
			tsroot = tsroot.rstrip('/')

		if not os.path.isdir(tsroot) or\
		   not os.path.isfile(os.path.join(tsroot, 'bin', 'trafficserver')):

			return False
	except (OSError, AttributeError, ValueError) as e:
		raise ValueError from e

	global TS_ROOT
	TS_ROOT = tsroot
	return True


#: :const:`True` if Traffic Ops communicates using SSL, :const:`False` otherwise
TO_USE_SSL = False

#: Holds only the :abbr:`FQDN (Fully Quallified Domain Name)` of the Traffic Ops server
TO_HOST = None

#: Holds the port number on which the Traffic Ops server listens for incoming HTTP(S) requests
TO_PORT = None

def setTOURL(url:str) -> bool:
	"""
	Sets the :data:`TO_USE_SSL`, :data:`TO_PORT` and :data:`TO_HOST` global variables and verifies,
	them.

	:param url: A full URL (including schema - and port when necessary) specifying the location of
		a running Traffic Ops server
	:returns: whether or not the URL could be set successfully
	:raises ValueError: when ``url`` is not a :const:`str`
	"""
	global VERIFY
	try:
		url = url.rstrip('/')
		_ = requests.head(url, verify=VERIFY)
	except requests.exceptions.RequestException as e:
		logging.error("%s", e)
		logging.debug("%s", e, exc_info=True, stack_info=True)
		return False
	except (AttributeError, ValueError) as e:
		raise ValueError from e

	global TO_HOST, TO_PORT, TO_USE_SSL

	port = None

	if url.lower().startswith("http://"):
		url = url[7:]
		port = 80
		TO_USE_SSL = False
	elif url.lower().startswith("https://"):
		url = url[8:]
		port = 443
		TO_USE_SSL = True

	# I'm assuming here that a valid FQDN won't include ':' - and it shouldn't.
	portpoint = url.find(':')
	if portpoint > 0:
		TO_HOST = url[:portpoint]
		port = int(url[portpoint+1:])
	else:
		TO_HOST = url

	if port is None:
		raise ValueError("Couldn't determine port number from URL '%s'!" % url)

	TO_PORT = port

	return True

#: Holds the username used for logging in to Traffic Ops
USERNAME = None

#: Holds the password used to authenticate :data:`USERNAME` with Traffic Ops
PASSWORD = None

def setTOCredentials(login:str) -> bool:
	"""
	Parses and returns a JSON-encoded login string for the Traffic Ops API.
	This will set :data:`USERNAME` and :data:`PASSWORD` if login is successful.

	:param login: The raw login info as passed on the command line (e.g. 'username:password')
	:raises ValueError: if ``login`` is not a :const:`str`.
	:returns: whether or not the login could be set and validated successfully
	"""
	try:
		u, p = login.split(':')
		login = '{{"u": "{0}", "p": "{1}"}}'.format(u, p)
	except IndexError:
		logging.critical("Bad Traffic_Ops_Login: '%s' - should be like 'username:password'", login)
		return False
	except (AttributeError, ValueError) as e:
		raise ValueError from e

	global USERNAME, PASSWORD
	USERNAME = u
	PASSWORD = p

	return True
