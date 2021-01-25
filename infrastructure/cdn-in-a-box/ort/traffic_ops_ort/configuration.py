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

import argparse
import enum
import logging
import os
import platform
import typing

import distro
import requests


#: A format specifier for logging output. Propagates to all imported modules.
LOG_FORMAT = "%(levelname)s: %(asctime)s line %(lineno)d in %(module)s.%(funcName)s: %(message)s"


#: contains identifying information about the host system's Linux distribution
DISTRO = distro.LinuxDistribution().id()


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


class Configuration():
	"""
	Represents a configured state for :program:`traffic_ops_ort`.
	"""

	class Modes(enum.IntEnum):
		"""
		Enumerated representations for run modes for valid configurations.
		"""
		REPORT = 0      #: Do nothing, only report what would be done
		INTERACTIVE = 1 #: Ask for user confirmation before modifying the system
		REVALIDATE = 2  #: Only check for configuration file changes and content revalidations
		SYNCDS = 3      #: Check for and apply Delivery Service changes
		BADASS = 4      #: Apply all settings specified in Traffic Ops, with no restrictions

		def __str__(self) -> str:
			"""
			Implements ``str(self)``

			:returns: the enum member's name
			"""
			return self.name


	#: Holds a reference to a :class:`to_api.API` object used by this configuration - must be set
	#: manually.
	api = None


	#: Holds a reference to a :class:`to_api.ServerInfo` object used by this configuration - must be
	#: set manually.
	ServerInfo = None


	def __init__(self, args:argparse.Namespace):
		"""
		Constructs the configuration object.

		:param args: Should be the result of parsing command-line arguments to :program:`traffic_ops_ort`
		:raises ValueError: if an error occurred setting up the configuration
		"""
		global DISTRO

		self.dispersion = args.dispersion if args.dispersion > 0 else 0
		self.login_dispersion = args.login_dispersion if args.login_dispersion > 0 else 0
		self.wait_for_parents = bool(args.wait_for_parents)
		self.retries = args.retries if args.retries > 0 else 0
		self.rev_proxy_disable = args.rev_proxy_disable
		self.verify = not args.insecure
		self.timeout = args.timeout
		self.via_string_release = args.via_string_release
		self.disable_parent_config_comments = args.disable_parent_config_comments

		setLogLevel(args.log_level)

		logging.info("Distribution detected as: '%s'", DISTRO)

		if not args.hostname:
			self.hostname = (platform.node().split('.')[0], platform.node())
			logging.info("Hostname detected as: '%s'", self.fullHostname)
		else:
			self.hostname = (args.hostname, args.hostname)
			logging.info("Hostname set to: '%s'", self.fullHostname)

		try:
			self.mode = Configuration.Modes[args.Mode.upper()]
		except KeyError as e:
			raise ValueError("Unrecognized Mode: '%s'" % args.Mode)

		self.tsroot = parseTSRoot(args.ts_root)
		logging.info("ATS root installation directory set to '%s'", self.tsroot)

		self.useSSL, self.toHost, self.toPort = parseTOURL(args.to_url, self.verify)
		self.username, self.password = args.to_user, args.to_password


	@property
	def shortHostname(self) -> str:
		"""
		Convenience accessor for the short hostname of this server

		:returns: The (short) hostname of this server as detected by :func:`platform.node`
		"""
		return self.hostname[0]

	@property
	def fullHostname(self) -> str:
		"""
		Convenience accessor for the full hostname of this server

		:returns: The hostname of this server as detected by :func:`platform.node`
		"""
		return self.hostname[1]

	@property
	def TOURL(self) -> str:
		"""
		Convenience function to construct a full URL out of whatever information was given at runtime

		:returns: The configuration's URL which points to its Traffic Ops server instance

		.. note:: This is totally constructed from information given on the command line; the
			resulting URL may actually point to a reverse proxy for the Traffic Ops server and not
			the server itself.
		"""
		return "%s://%s:%d/" % ("https" if self.useSSL else "https", self.toHost, self.toPort)


def setLogLevel(level:str):
	"""
	Parses a string to return the requested :class:`LogLevels` member, to which it will then set
	the global logging level.

	:param level: the name of a LogLevels enum constant
	:raises ValueError: if ``level`` cannot be parsed to an actual LogLevel
	"""
	global LOG_FORMAT

	try:
		level = LogLevels[level.upper()]
	except KeyError as e:
		raise ValueError("Unrecognized log level: '%s'" % level) from e

	logging.basicConfig(level=level, format=LOG_FORMAT)
	logging.getLogger().setLevel(level)


def parseTSRoot(tsroot:str) -> str:
	"""
	Parses and validates a given path as a path to the root of an Apache Traffic Server installation

	:param tsroot: The relative or absolute path to the root of this server's ATS installation
	:raises ValueError: if ``tsroot`` is not an existing path, or does not contain the ATS binary
	"""
	tsroot = tsroot.strip()
	if tsroot != '/' and tsroot.endswith('/'):
		tsroot = tsroot.rstrip('/')

	try:
		if not os.path.isdir(tsroot):
			raise ValueError("'%s' is not a directory!" % tsroot)

		binpath = os.path.join(tsroot, 'bin', 'trafficserver')
		if not os.path.isfile(binpath):
			raise ValueError("'%s' does not exist! '%s' is not the root of a Traffic Server"
			                 "installation" % (binpath, tsroot))
	except OSError as e:
		raise ValueError("Couldn't set the ATS root install directory: %s" % e) from e

	return tsroot


def parseTOURL(url:str, verify:bool) -> typing.Tuple[bool, str, int]:
	"""
	Parses and verifies the passed URL and breaks it into parts for the caller

	:param url: At minimum an FQDN for a Traffic Ops server, but can include schema and port number
	:param verify: Whether or not to verify the server's SSL certificate
	:returns: Whether or not the Traffic Ops server uses SSL (http vs https), the server's FQDN, and the port on which it listens
	:raises ValueError: if ``url`` does not point at a valid HTTP server or is incorrectly formatted
	"""
	url = url.rstrip('/')

	useSSL, host, port = True, None, 443

	try:
		_ = requests.head(url, verify=verify)
	except requests.exceptions.RequestException as e:
		raise ValueError("Cannot contact any server at '%s' (%s)" % (url, e)) from e

	if url.lower().startswith("http://"):
		port = 80
		useSSL = False
		url = url[7:]
	elif url.lower().startswith("https://"):
		url = url[8:]

	# I'm assuming here that a valid FQDN won't include ':' - and it shouldn't
	portpoint = url.find(':')
	if portpoint > 0:
		host = url[:portpoint]
		port = int(url[portpoint+1:])
	else:
		host = url

	return useSSL, host, port
