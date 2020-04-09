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
This module contains functionality for dealing with the Traffic Ops ReST API.
It extends the class provided by the official Apache Traffic Control Client.
"""

import typing
import logging
import re

from munch import Munch
from requests.compat import urljoin
from requests.exceptions import RequestException

from trafficops.tosession import TOSession
from trafficops.restapi import LoginError, OperationError, InvalidJSONError

from . import packaging
from .configuration import Configuration

class API(TOSession):
	"""
	This class extends :class:`trafficops.tosession.TOSession` to provide some ease-of-use
	functionality for getting things needed by :term:`ORT`.

	TODO: update to 2.0 if/when atstccfg support is integrated.
	"""

	#: This should always be the latest API version supported - note this breaks compatability with
	#: older ATC versions. Go figure.
	VERSION = "1.4"

	def __init__(self, conf:Configuration):
		"""
		This not only creates the API session, but log the user in immediately.

		:param conf: An object that represents the configuration of :program:`traffic_ops_ort`
		:raises LoginError: when authentication with Traffic Ops fails
		:raises OperationError: when some anonymous error occurs communicating with the Traffic Ops server
		"""
		super(API, self).__init__(host_ip=conf.toHost, api_version=self.VERSION,
		                          host_port=conf.toPort, verify_cert=conf.verify, ssl=conf.useSSL)

		self.retries = conf.retries

		for r in range(self.retries):
			try:
				logging.info("login attempt #%d", r)
				self.login(conf.username, conf.password)
				break
			except (LoginError, OperationError, InvalidJSONError, RequestException) as e:
				logging.debug("login failure: %r", e, stack_info=True, exc_info=True)
		else:
			raise LoginError("Failed to log in to Traffic Ops, retries exceeded.")

		self.hostname = conf.shortHostname

	def __enter__(self):
		"""
		Implements context-management for :class:`API` objects. Needs to override :class:`TOSession`
		context-management because the connection is already established during initialization.
		"""
		return self

	def getRaw(self, path:str) -> str:
		"""
		This gets the API response to a "raw" path, meaning it will queried directly without a
		``/api/1.x`` prefix. Because the output structure of the API response is not known, this
		returns the response body as an unprocessed string rather than a Python object via e.g.
		Munch.

		:param path: The raw path on the Traffic Ops server
		:returns: The API response payload
		:raises ConnectionError: When something goes wrong communicating with the Traffic Ops server
		"""
		for _ in range(self.retries):
			try:
				r = self._session.get(path)
				break
			except (LoginError, OperationError, InvalidJSONError, RequestException) as e:
				logging.debug("API failure: %r", e, stack_info=True, exc_info=True)
		else:
			raise ConnectionError("Failed to get a valid response from Traffic Ops for %s" % path)

		if r.status_code != 200 and r.status_code != 204:
			raise ConnectionError("request for '%s' appears to have failed; reason: %s" %
			                                  (path,                             r.reason))

		return r.text

	def getMyPackages(self) -> typing.List[packaging.Package]:
		"""
		Fetches a list of the packages specified by Traffic Ops that should exist on this server.

		:returns: all of the packages which this system must have, according to Traffic Ops.
		:raises ConnectionError: if fetching the package list fails
		"""
		logging.info("Fetching this server's package list from Traffic Ops")

		# Ah, read-only properties that gut functionality, my favorite.
		tmp = self.api_base_url
		self._api_base_url = urljoin(self._server_url, '/').rstrip('/') + '/'

		packagesPath = '/'.join(("ort", self.hostname, "packages"))
		for _ in range(self.retries):
			try:
				myPackages = self.get(packagesPath)
				break
			except (LoginError, OperationError, InvalidJSONError, RequestException) as e:
				logging.debug("package fetch failure: %r", e, stack_info=True, exc_info=True)
		else:
			self._api_base_url = tmp
			raise ConnectionError("Failed to get a response for packages")

		self._api_base_url = tmp

		logging.debug("Raw package response: %s", myPackages[1].text)

		try:
			return [packaging.Package(p) for p in myPackages[0]]
		except ValueError:
			raise ConnectionError

	def setConfigFileAPIVersion(self, files: Munch) -> None:
		match_api_base = re.compile(r'^(/api/)\d+\.\d+(/)')
		api_base_replacement = r'\g<1>%s\2' % API.VERSION
		for configFile in files.configFiles:
			configFile.apiUri = match_api_base.sub(api_base_replacement, configFile.apiUri)

	def getMyConfigFiles(self, conf:Configuration) -> typing.List[dict]:
		"""
		Fetches configuration files constructed by Traffic Ops for this server

		.. note:: This function will set the :attr:`serverInfo` attribute of the object passed as
			the ``conf`` argument to an instance of :class:`ServerInfo` with the provided
			information.

		:param conf: An object that represents the configuration of :program:`traffic_ops_ort`
		:returns: A list of constructed config file objects
		:raises ConnectionError: when something goes wrong communicating with Traffic Ops
		"""
		logging.info("Fetching list of configuration files from Traffic Ops")
		for _ in range(self.retries):
			try:
				# The API function decorator confuses pylint into thinking this doesn't return
				#pylint: disable=E1111
				myFiles = self.get_server_config_files(host_name=self.hostname)
				#pylint: enable=E1111
				break
			except (InvalidJSONError, LoginError, OperationError, RequestException) as e:
				logging.debug("config file fetch failure: %r", e, exc_info=True, stack_info=True)
		else:
			raise ConnectionError("Failed to fetch configuration files from Traffic Ops")

		logging.debug("Raw response from Traffic Ops: %s", myFiles[1].text)
		myFiles = myFiles[0]
		self.setConfigFileAPIVersion(myFiles)

		try:
			conf.serverInfo = ServerInfo(myFiles.info)
			# if there's a reverse proxy, switch urls.
			if conf.serverInfo.toRevProxyUrl and not conf.rev_proxy_disable:
				self._server_url = conf.serverInfo.toRevProxyUrl
				self._api_base_url = urljoin(self._server_url, '/api/%s' % self.VERSION).rstrip('/') + '/'
			return myFiles.configFiles
		except (KeyError, AttributeError, ValueError) as e:
			raise ConnectionError("Malformed response from Traffic Ops to update status request!") from e

	def updateTrafficOps(self, mode:Configuration.Modes):
		"""
		Updates Traffic Ops's knowledge of this server's update status.

		:param mode: The current run-mode of :program:`traffic_ops_ort`
		"""
		from .utils import getYesNoResponse as getYN

		if mode is Configuration.Modes.INTERACTIVE and not getYN("Update Traffic Ops?", default='Y'):
			logging.warning("Update will not be performed; you should clear updates manually")
			return

		logging.info("Updating Traffic Ops")

		if mode is Configuration.Modes.REPORT:
			return

		payload = {"updated": False, "reval_updated": False}

		for _ in range(self.retries):
			try:
				response = self._session.post('/'.join((self._server_url.rstrip('/'),
				                                        "update",
				                                        self.hostname)
				                             ), data=payload)
				break
			except (LoginError, InvalidJSONError, OperationError, RequestException) as e:
				logging.debug("TO update failure: %r", e, exc_info=True, stack_info=True)
		else:
			raise ConnectionError("Failed to update Traffic Ops - connection was lost")

		if response.text:
			logging.info("Traffic Ops response: %s", response.text)

	def getMyChkconfig(self) -> typing.List[dict]:
		"""
		Fetches the 'chkconfig' for this server

		:returns: An iterable list of 'chkconfig' entries
		:raises ConnectionError: when something goes wrong communicating with Traffic Ops
		"""


		# Ah, read-only properties that gut functionality, my favorite.
		tmp = self.api_base_url
		self._api_base_url = urljoin(self._server_url, '/').rstrip('/') + '/'

		uri = "ort/%s/chkconfig" % self.hostname
		logging.info("Fetching chkconfig from %s", uri)

		for _ in range(self.retries):
			try:
				r = self.get(uri)
				break
			except (InvalidJSONError, OperationError, LoginError, RequestException) as e:
				logging.debug("chkconfig fetch failure: %r", e, exc_info=True, stack_info=True)
		else:
			self._api_base_url = tmp
			raise ConnectionError("Failed to fetch 'chkconfig' from Traffic Ops - connection lost")

		self._api_base_url = tmp

		logging.debug("Raw response from Traffic Ops: %s", r[1].text)

		return r[0]

	def getMyUpdateStatus(self) -> dict:
		"""
		Gets the update status of a server.

		:raises ConnectionError: if something goes wrong communicating with the server
		:returns: An object representing the API's response
		"""
		logging.info("Fetching update status from Traffic Ops")
		for _ in range(self.retries):
			try:
				# The API function decorator confuses pylint into thinking this doesn't return
				#pylint: disable=E1111
				r = self.get_server_update_status(server_name=self.hostname)
				#pylint: enable=E1111
				break
			except (InvalidJSONError, LoginError, OperationError, RequestException) as e:
				logging.debug("update status fetch failure: %r", e, exc_info=True, stack_info=True)
		else:
			raise ConnectionError("Failed to fetch update status - connection was lost")

		logging.debug("Raw response from Traffic Ops: %s", r[1].text)

		return r[0]

	def getMyStatus(self) -> str:
		"""
		Fetches the status of this server as set in Traffic Ops

		:raises ConnectionError: if fetching the status fails
		:returns: the name of the status to which this server is set in the Traffic Ops configuration

		"""
		logging.info("Fetching server status from Traffic Ops")

		for _ in range(self.retries):
			try:
				# The API function decorator confuses pylint into thinking this doesn't return
				#pylint: disable=E1111
				r = self.get_servers(query_params={"hostName": self.hostname})
				#pylint: enable=E1111
				break
			except (InvalidJSONError, LoginError, OperationError,RequestException) as e:
				logging.debug("status fetch failure: %r", e, exc_info=True, stack_info=True)
		else:
			raise ConnectionError("Failed to fetch server status - connection was lost")

		logging.debug("Raw response from Traffic Ops: %s", r[1].text)

		r = r[0][0]

		try:
			return r.status
		except (IndexError, KeyError, AttributeError) as e:
			raise ConnectionError("Malformed response from Traffic Ops to update status request!") from e

#: Caches the names of statuses supported by Traffic Ops
CACHED_STATUSES = []

#: Maps Traffic Ops alert levels to logging levels
API_LOGGERS = {"error":   lambda x: logging.error("Traffic Ops API alert: %s", x),
               "warning": lambda x: logging.warning("Traffic Ops API alert: %s", x),
               "info":    lambda x: logging.info("Traffic Ops API alert: %s", x),
               "success": lambda x: logging.info("Traffic Ops API alert: %s", x)}

class ServerInfo():
	"""
	Holds information about a server, as returned by the Traffic Ops API
	``api/1.x/servers/<hostname>/configfiles/ats`` endpoint
	"""

	cdnId = -1         #: A database primary key for the CDN to which this server is assigned
	cdnName = ""       #: The name of the CDN to which this server is assigned
	profileName = ""   #: The name of the profile in use by this server
	profileId = -1     #: A database primary key for this server's profile's information
	serverId = -1      #: A database primary key for this server's information
	serverIpv4 = ""    #: This server's IPv4 address
	serverName = ""    #: This server's short hostname
	serverTcpPort = 80 #: The port on which the caching proxy of this server listens
	toUrl = ""         #: The Traffic Ops URL... not sure what that's for...

	#: This specifies the url of a reverse proxy that should be used for future requests to the
	#: Traffic Ops API - if present.
	toRevProxyUrl = ""

	def __init__(self, raw:dict):
		"""
		Constructs a server object out of some kind of raw API response

		:param raw: some kind of ungodly huge JSON object from one API endpoint or
			another. Attempts will be made to resolve inconsistent naming accross
			endpoints.
		:raises ValueError: when the passed object doesn't have all required fields
		"""
		try:
			self.cdnId =         raw["cdnId"]
			self.cdnName =       raw["cdnName"]
			self.profileName =   raw["profileName"]
			self.profileId =     raw["profileId"]
			self.serverId =      raw["serverId"]
			self.serverIpv4 =    raw["serverIpv4"]
			self.serverName =    raw["serverName"]
			self.serverTcpPort = raw["serverTcpPort"]
			self.toUrl =         raw["toUrl"]

			# This may or may not exist
			if "toRevProxyUrl" in raw:
				self.toRevProxyUrl = raw["toRevProxyUrl"]
		except (KeyError, TypeError) as e:
			raise ValueError from e


	def __repr__(self) -> str:
		"""
		Implements ``str(self)``
		"""
		out = "Server(%s)"
		return out % ', '.join(("%s=%r" % (a, self.__getattribute__(a))\
		                       for a in dir(self)\
		                       if not a.startswith('_')))

	def sanitize(self, fmt:str, hostname:typing.Tuple[str, str]) -> str:
		"""
		Sanitizes an input string with the passed hostname information

		:param fmt: The string to be sanitized
		:param hostname: A tuple containing the short and full hostnames of the server
		:returns: The string ``fmt`` after sanitization
		"""
		fmt = fmt.replace("__HOSTNAME__", hostname[0])
		fmt = fmt.replace("__FULL_HOSTNAME__", hostname[1])
		fmt = fmt.replace("__RETURN__", '\n')
		fmt = fmt.replace("__CACHE_IPV4__", self.serverIpv4)

		# Don't ask me why, but the reference ORT implementation just strips these ones out
		# if the tcp port is 80.
		return fmt.replace("__SERVER_TCP_PORT__", str(self.serverTcpPort)\
		                                          if self.serverTcpPort != 80 else "")
