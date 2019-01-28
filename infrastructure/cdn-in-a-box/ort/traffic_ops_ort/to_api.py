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

from trafficops.tosession import TOSession

from . import packaging

class API(TOSession):
	"""
	This class extends :class:`trafficops.tosession.TOSession` to provide some ease-of-use
	functionality for getting things needed by :term:`ORT`.
	"""

	#: This should always be the latest API version supported - note this breaks compatability with
	#: older ATC versions. Go figure.
	VERSION = "1.4"

	#: Caches update statuses mapped by hostnames
	CACHED_UPDATE_STATUS = {}

	def __init__(self, username:str, password:str, toHost:str, myHostname:str, port:int = 443,
	                   verify:bool = True, useSSL:bool = True):
		"""
		This not only creates the API session, but log the user in immediately.

		:param username: The name of the user as whom :term:`ORT` will authenticate with Traffic Ops
		:param password: The password of Traffic Ops user ``username``
		:param toHost: The :abbr:`FQDN (Fully Qualified Domain Name)` of the Traffic Ops server
		:param myHostname: The (short) hostname of **this** server
		:param port: The port number on which Traffic Ops listens for incoming HTTP(S) requests
		:param verify: If :const:`True` SSL certificates will be verifed, if :const:`False` they
		               will not and warnings about unverified SSL certificates will be swallowed.
		:param useSSL: If :const:`True` :term:`ORT` will attempt to communicate with Traffic Ops
		               using SSL, if :const:`False` it will not. *This setting will be respected
		               regardless of the passed port number!*
		:raises trafficops.restapi.LoginError: when authentication with Traffic Ops fails
		:raises trafficops.restapi.OperationError: when some anonymous error occurs communicating
		                                           with the Traffic Ops server
		"""
		super(API, self).__init__(host_ip=toHost, api_version=self.VERSION, host_port=port,
		                          verify_cert=verify, ssl=useSSL)
		self.login(username, password)

		self.hostname = myHostname

	def getRaw(self, path:str) -> str:
		"""
		This gets the API response to a "raw" path, meaning it will queried directly without a
		``/api/1.x`` prefix. Because the output structure of the API response is not known, this
		returns the response body as an unprocessed string rather than a Python object via e.g.
		Munch.

		:param path: The raw path on the Traffic Ops server
		:returns: The API response payload
		"""

		r = self._session.get(path)

		if r.status_code != 200 and r.status_code != 204:
			raise ValueError("request for '%s' appears to have failed; reason: %s" % (path, r.reason))

		return r.text

	def getMyPackages(self) -> typing.List[packaging.Package]:
		"""
		Fetches a list of the packages specified by Traffic Ops that should exist on this server.

		:returns: all of the packages which this system must have, according to Traffic Ops.
		:raises ConnectionError: if fetching the package list fails
		:raises ValueError: if the API endpoint returns a malformed response that can't be parsed
		"""
		logging.info("Fetching this server's package list from Traffic Ops")

		# Ah, read-only properties that gut functionality, my favorite.
		from requests.compat import urljoin
		tmp = self.api_base_url
		self._api_base_url = urljoin(self._server_url, '/').rstrip('/') + '/'

		packagesPath = '/'.join(("ort", self.hostname, "packages"))
		myPackages = self.get(packagesPath)
		self._api_base_url = tmp

		logging.debug("Raw package response: %s", myPackages[1].text)

		return [packaging.Package(p) for p in myPackages[0]]

	def getMyConfigFiles(self) -> typing.List[dict]:
		"""
		Fetches configuration files constructed by Traffic Ops for this server

		.. note:: This function will set the :data:`traffic_ops_ort.configuration.SERVER_INFO`
			object to an instance of :class:`ServerInfo` with the provided information.

		:returns: A list of constructed config file objects
		:raises ConnectionError: when something goes wrong communicating with Traffic Ops
		:raises ValueError: when a response was successfully obtained from the Traffic Ops API, but the
			response could not successfully be parsed as JSON, or was missing information
		"""
		from . import configuration

		logging.info("Fetching list of configuration files from Traffic Ops")

		# The API function decorator confuses pylint into thinking this doesn't return
		#pylint: disable=E1111
		myFiles = self.get_server_config_files(host_name=self.hostname)
		#pylint: enable=E1111

		logging.debug("Raw response from Traffic Ops: %s", myFiles[1].text)
		myFiles = myFiles[0]

		try:
			configuration.SERVER_INFO = ServerInfo(myFiles.info)
			return myFiles.configFiles
		except (KeyError, AttributeError) as e:
			raise ValueError from e

	def updateTrafficOps(self):
		"""
		Updates Traffic Ops's knowledge of this server's update status.
		"""
		from .configuration import MODE, Modes
		from .utils import getYesNoResponse as getYN

		if MODE is Modes.INTERACTIVE and not getYN("Update Traffic Ops?", default='Y'):
			logging.warning("Update will not be performed; you should clear updates manually")
			return

		logging.info("Updating Traffic Ops")

		if MODE is Modes.REPORT:
			return

		payload = {"updated": False, "reval_updated": False}
		response = self._session.post('/'.join((self._server_url.rstrip('/'),
		                                        "update",
		                                        self.hostname)
		                             ), data=payload)

		if response.text:
			logging.info("Traffic Ops response: %s", response.text)

	def getMyChkconfig(self) -> typing.List[dict]:
		"""
		Fetches the 'chkconfig' for this server

		:returns: An iterable list of 'chkconfig' entries
		:raises ConnectionError: when something goes wrong communicating with Traffic Ops
		:raises ValueError: when a response was successfully obtained from the Traffic Ops API, but the
			response could not successfully be parsed as JSON, or was missing information
		"""


		# Ah, read-only properties that gut functionality, my favorite.
		from requests.compat import urljoin
		tmp = self.api_base_url
		self._api_base_url = urljoin(self._server_url, '/').rstrip('/') + '/'

		uri = "ort/%s/chkconfig" % self.hostname
		logging.info("Fetching chkconfig from %s", uri)

		r = self.get(uri)
		self._api_base_url = tmp
		logging.debug("Raw response from Traffic Ops: %s", r[1].text)

		return r[0]

	def getUpdateStatus(self, host:str) -> dict:
		"""
		Gets the update status of a server.

		.. note:: If the :data:`self.CACHED_UPDATE_STATUS` cached response is set, this function will
			default to that object. If it is *not* set, then this function will set it.

		:param host: The (short) hostname of the server to query
		:raises PermissionError: if a new cookie is required, but fails to be aquired
		:returns: An object representing the API's response
		"""

		logging.info("Fetching update status for %s from Traffic Ops", host)

		if host in self.CACHED_UPDATE_STATUS:
			return self.CACHED_UPDATE_STATUS[host]

		# The API function decorator confuses pylint into thinking this doesn't return
		#pylint: disable=E1111
		r = self.get_server_update_status(server_name=host)
		#pylint: enable=E1111
		logging.debug("Raw response from Traffic Ops: %s", r[1].text)

		self.CACHED_UPDATE_STATUS[host] = r[0]

		return r[0]

	def getMyStatus(self) -> str:
		"""
		Fetches the status of this server as set in Traffic Ops

		:raises ConnectionError: if fetching the status fails
		:raises ValueError: if the :data:`traffic_ops_ort.configuration.HOSTNAME` is not properly set,
			or a weird value is stored in the global :data:`CACHED_UPDATE_STATUS` response cache.
		:returns: the name of the status to which this server is set in the Traffic Ops configuration

		.. note:: If the global :data:`CACHED_UPDATE_STATUS` cached response is set, this function will
			default to the status provided by that object.
		"""


		logging.info("Fetching server status from Traffic Ops")

		# The API function decorator confuses pylint into thinking this doesn't return
		#pylint: disable=E1111
		r = self.get_servers(query_params={"hostName": self.hostname})
		#pylint: enable=E1111

		logging.debug("Raw response from Traffic Ops: %s", r[1].text)

		r = r[0][0]

		try:
			return r.status
		except (IndexError, KeyError, AttributeError) as e:
			logging.error("Malformed response from Traffic Ops to update status request!")
			raise ConnectionError from e

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

	def sanitize(self, fmt:str) -> str:
		"""
		Implements ``str.format(self)``
		"""
		from .configuration import HOSTNAME
		fmt = fmt.replace("__HOSTNAME__", HOSTNAME[0])
		fmt = fmt.replace("__FULL_HOSTNAME__", HOSTNAME[1])
		fmt = fmt.replace("__RETURN__", '\n')
		fmt = fmt.replace("__CACHE_IPV4__", self.serverIpv4)

		# Don't ask me why, but the reference ORT implementation just strips these ones out
		# if the tcp port is 80.
		return fmt.replace("__SERVER_TCP_PORT__", str(self.serverTcpPort)\
		                                          if self.serverTcpPort != 80 else "")
