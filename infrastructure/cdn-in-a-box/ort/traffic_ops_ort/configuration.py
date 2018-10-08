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

import datetime
import enum
import logging
import os
import platform
import typing

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


#: Holds the full URL including schema (e.g. 'http') and port that points at Traffic Ops
TO_URL = None

def setTOURL(url:str) -> bool:
	"""
	Sets the :data:`TO_URL` global variable and verifies it

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

	global TO_URL
	TO_URL = url

	return True


#: Holds a Mojolicious cookie for validating connections to Traffic Ops
TO_COOKIE = None

#: Holds the login information for re-obtaining a cookie when the one in :data:`TO_COOKIE` expires
TO_LOGIN = None

def setTOCredentials(login:str) -> bool:
	"""
	Parses and returns a JSON-encoded login string for the Traffic Ops API.
	This will set :data:`TO_COOKIE` and :data:`TO_LOGIN` if login is successful.

	:param login: The raw login info as passed on the command line (e.g. 'username:password')
	:raises ValueError: if ``login`` is not a :const:`str`.
	:returns: whether or not the login could be set and validated successfully
	"""
	try:
		login = '{{"u": "{0}", "p": "{1}"}}'.format(*login.split(':'))
	except IndexError:
		logging.critical("Bad Traffic_Ops_Login: '%s' - should be like 'username:password'", login)
		return False
	except (AttributeError, ValueError) as e:
		raise ValueError from e

	global TO_LOGIN
	TO_LOGIN = login

	try:
		getNewTOCookie()
	except PermissionError:
		return False

	return True

def getNewTOCookie():
	"""
	Re-obtains a cookie from Traffic Ops based on the login credentials in :data:`TO_LOGIN`

	:raises PermissionError: if :data:`TO_LOGIN` or :data:`TO_URL` are unset, invalid,
		or the wrong type
	"""
	global TO_URL, TO_LOGIN, VERIFY, TO_COOKIE
	if TO_URL is None or not isinstance(TO_URL, str) or\
	   TO_LOGIN is None or not isinstance(TO_LOGIN, str):
		raise PermissionError("TO_URL and TO_LOGIN must be set prior to calling this function!")

	try:
		# Obtain login cookie
		cookie = requests.post(TO_URL + '/api/1.3/user/login', data=TO_LOGIN, verify=VERIFY)
	except requests.exceptions.RequestException as e:
		logging.critical("Login credentials rejected by Traffic Ops")
		raise PermissionError from e

	if not cookie.cookies or 'mojolicious' not in cookie.cookies:
		logging.error("Response code: %d", cookie.status_code)
		logging.warning("Response Headers: %s", cookie.headers)
		logging.debug("Response: %s", cookie.content)
		raise PermissionError("Login credentials rejected by Traffic Ops")

	TO_COOKIE = [c for c in cookie.cookies if c.name == "mojolicious"][0]

def getTOCookie() -> typing.Dict[str, str]:
	"""
	A small, convenience wrapper for getting a current, valid Traffic Ops authentication cookie. If
	:data:`TO_COOKIE` is expired, this function requests a new one from Traffic Ops

	:returns: A cookie dataset that may be passed directly to :mod:`requests` functions
	:raises PermissionError: if :data:`TO_LOGIN` and/or :data:`TO_URL` are unset, invalid, or the
		wrong type
	"""
	global TO_COOKIE

	if datetime.datetime.now().timestamp() >= TO_COOKIE.expires:
		getNewTOCookie()

	return {TO_COOKIE.name: TO_COOKIE.value}
