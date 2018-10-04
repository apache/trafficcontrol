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
This module contains functionality for dealing with the Traffic Ops ReST API
"""

import datetime
import typing
import logging
import requests

from . import packaging

#: Caches update statuses mapped by hostnames
CACHED_UPDATE_STATUS = {}

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

	def __format__(self, fmt:str) -> str:
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

def TOPost(uri:str, data:dict) -> str:
	"""
	POSTs the passed data in a request to the specified API endpoint

	:param uri: The Traffic Ops URL-relative path to an API endpoint, e.g. if the intention is
		to post to ``https://TO_URL:TO_PORT/api/1.3/users``, this should just be
		``'api/1.3/users'``

			.. note:: This function will ensure the proper concatenation of the Traffic Ops URL
				to the request path; callers need not worry about whether the ``uri`` ought to
				begin with a slash.

	:returns: The Traffic Ops server's response to the POST request - possibly empty - as a UTF-8
		string
	:raises ConnectionError: when an error occurs trying to communicate with Traffic Ops
	"""
	from . import configuration as conf

	uri = '/'.join((conf.TO_URL, uri.lstrip('/')))
	logging.info("POSTing %r to %s", data, uri)

	try:
		resp = requests.post(uri, cookies=conf.getTOCookie(), verify=conf.VERIFY, data=data)
	except (PermissionError, requests.exceptions.RequestException) as e:
		raise ConnectionError from e

	logging.debug("Raw response from Traffic Ops: %s\n%s\n%s", resp, resp.headers, resp.content)

	return resp.text

def getTOJSONResponse(uri:str) -> dict:
	"""
	A wrapper around :func:`traffic_ops_ort.utils.getJSONResponse` that handles cookies and
	tacks on the top-level Traffic Ops URL.

	:param uri: The Traffic Ops URL-relative path to a JSON API endpoint, e.g. if the intention
		is to get ``https://TO_URL:TO_PORT/api/1.3/ping``, this should just be ``'api/1.3/ping'``

			.. note:: This function will ensure the proper concatenation of the Traffic Ops URL
				to the request path; callers need not worry about whether the ``uri`` ought to
				begin with a slash.

	:returns: The decoded JSON response as an object

			.. note:: If the API response containes a 'response' object, this function will
				only return that object. Also, if the API response contains an 'alerts' object,
				they will be logged appropriately

	:raises ConnectionError: when an error occurs trying to communicate with Traffic Ops
	:raises ValueError: when the request completes successfully, but the response body
		does not represent a JSON-encoded object.
	"""
	global API_LOGGERS
	from . import configuration as conf, utils

	uri = '/'.join((conf.TO_URL, uri.lstrip('/')))
	logging.info("Fetching Traffic Ops API response: %s", uri)

	if datetime.datetime.now().timestamp() >= conf.TO_COOKIE.expires:
		try:
			conf.getNewTOCookie()
		except PermissionError as e:
			raise ConnectionError from e

	resp = utils.getJSONResponse(uri,
	                             cookies = {conf.TO_COOKIE.name:conf.TO_COOKIE.value},
	                             verify = conf.VERIFY)

	if "response" in resp:
		if "alerts" in resp:
			for alert in resp["alerts"]:
				if "level" in alert:
					msg = alert["text"] if "text" in alert else "Unkown"
					API_LOGGERS[alert["level"]](msg)
				elif "text" in alert:
					logging.warning("Traffic Ops API alert: %s", alert["text"])
					logging.debug("Weird alert encountered: %r", alert)


		return resp["response"]

	return resp

def getUpdateStatus(host:str) -> dict:
	"""
	Gets the update status of a server.

	.. note:: If the global :data:`CACHED_UPDATE_STATUS` cached response is set, this function will
		default to that object. If it is *not* set, then this function will set it.

	:param host: The (short) hostname of the server to query
	:raises ValueError: if ``host`` is not a :const:`str`
	:raises PermissionError: if a new cookie is required, but fails to be aquired
	:returns: An object representing the API's response
	"""
	global CACHED_UPDATE_STATUS

	logging.info("Fetching update status for %s from Traffic Ops", host)
	if not isinstance(host, str):
		raise ValueError("First argument ('host') must be 'str', not '%s'" % type(host))

	if host in CACHED_UPDATE_STATUS:
		return CACHED_UPDATE_STATUS[host]

	CACHED_UPDATE_STATUS[host] = getTOJSONResponse("api/1.3/servers/%s/update_status" % host)

	return CACHED_UPDATE_STATUS[host]

def getMyStatus() -> str:
	"""
	Fetches the status of this server as set in Traffic Ops

	:raises ConnectionError: if fetching the status fails
	:raises ValueError: if the :data:`traffic_ops_ort.configuration.HOSTNAME` is not properly set,
		or a weird value is stored in the global :data:`CACHED_UPDATE_STATUS` response cache.
	:returns: the name of the status to which this server is set in the Traffic Ops configuration

	.. note:: If the global :data:`CACHED_UPDATE_STATUS` cached response is set, this function will
		default to the status provided by that object.
	"""
	global CACHED_UPDATE_STATUS
	from .configuration import HOSTNAME

	try:
		if HOSTNAME[0] in CACHED_UPDATE_STATUS:
			myStatus = CACHED_UPDATE_STATUS[HOSTNAME[0]]
			if "status" in CACHED_UPDATE_STATUS[HOSTNAME[0]]:
				return CACHED_UPDATE_STATUS[HOSTNAME[0]]["status"]

			logging.warning("CACHED_UPDATE_STATUS possibly set improperly")
			logging.warning("clearing this server's cached entry!")
			logging.debug("value was %r", myStatus)
			del CACHED_UPDATE_STATUS[HOSTNAME[0]]

	except (IndexError, KeyError) as e:
		raise ValueError from e

	myStatus = getUpdateStatus(HOSTNAME[0])

	try:
		return myStatus[0]["status"]
	except (IndexError, KeyError) as e:
		logging.error("Malformed response from Traffic Ops to update status request!")
		raise ConnectionError from e

def getStatuses() -> typing.Generator[str, None, None]:
	"""
	Yields a successive list of statuses supported by Traffic Ops.

	.. note:: This is implemented by iterating the :data:`CACHED_STATUSES` global cache -
		first populating it if it is empty - and so the validity of its outputs
		depends on the validity of the data stored therein

	:raises ValueError: if a response from the TO API is successful, but cannot be parsed as
		JSON
	:raises TypeError: if :data:`CACHED_STATUSES` is not iterable
	:raises ConnectionError: if something goes wrong contacting the Traffic Ops API
	:returns: an iterable generator that yields status names as strings
	"""
	global CACHED_STATUSES

	logging.info("Retrieving statuses from Traffic Ops")

	if CACHED_STATUSES:
		logging.debug("Using cached statuses: %r", CACHED_STATUSES)
		yield from CACHED_STATUSES
	else:
		statuses = getTOJSONResponse("api/1.3/statuses")
		yield from statuses

def getMyPackages() -> typing.List[packaging.Package]:
	"""
	Fetches a list of the packages specified by Traffic Ops that should exist on this server.

	:returns: all of the packages which this system must have, according to Traffic Ops.
	:raises ConnectionError: if fetching the package list fails
	:raises ValueError: if the API endpoint returns a malformed response that can't be parsed
	"""
	from .configuration import HOSTNAME

	logging.info("Fetching this server's package list from Traffic Ops")

	myPackages=getTOJSONResponse('/'.join(("ort", HOSTNAME[0], "packages")))

	logging.debug("Raw package response: %r", myPackages)

	return [packaging.Package(p) for p in myPackages]

def updateTrafficOps():
	"""
	Updates Traffic Ops's knowledge of this server's update status.
	"""
	from .configuration import MODE, Modes, HOSTNAME
	from .utils import getYesNoResponse as getYN

	if MODE is Modes.INTERACTIVE and not getYN("Update Traffic Ops?", default='Y'):
		logging.warning("Update will not be performed; you should do this manually")
		return

	logging.info("Updating Traffic Ops")

	if MODE is Modes.REPORT:
		return

	payload = {"updated": False, "reval_updated": False}
	response = TOPost("/update/%s" % HOSTNAME[0], payload)

	if response:
		logging.info("Traffic Ops response: %s", response)

def getMyConfigFiles() -> typing.List[dict]:
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

	uri = "/api/1.3/servers/%s/configfiles/ats" % configuration.HOSTNAME[0]

	myFiles = getTOJSONResponse(uri)

	try:
		configuration.SERVER_INFO = ServerInfo(myFiles["info"])
		return myFiles["configFiles"]
	except KeyError as e:
		raise ValueError from e

def getMyChkconfig() -> typing.List[dict]:
	"""
	Fetches the 'chkconfig' for this server

	:returns: An iterable list of 'chkconfig' entries
	:raises ConnectionError: when something goes wrong communicating with Traffic Ops
	:raises ValueError: when a response was successfully obtained from the Traffic Ops API, but the
		response could not successfully be parsed as JSON, or was missing information
	"""
	from . import configuration

	uri = "/ort/%s/chkconfig" % configuration.HOSTNAME[0]
	logging.info("Fetching chkconfig from %s", uri)

	return getTOJSONResponse(uri)
