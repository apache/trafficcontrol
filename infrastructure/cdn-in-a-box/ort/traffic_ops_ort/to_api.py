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

import json
import logging
import re
import subprocess
import typing
from os import path

from munch import Munch
from requests.compat import urljoin
from requests.exceptions import RequestException

from trafficops.tosession import TOSession
from trafficops.restapi import LoginError, OperationError, InvalidJSONError

from . import packaging, utils
from .config_files import ConfigFile
from .configuration import Configuration

class API(TOSession):
	"""
	This class extends :class:`trafficops.tosession.TOSession` to provide some ease-of-use
	functionality for getting things needed by :term:`ORT`.
	"""

	#: This should always be the latest API version supported - note this breaks compatibility with
	#: older ATC versions. Go figure.
	VERSION = "3.0"

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

		self.t3c_generate_cmd = [
			"t3c-generate",
			"--log-location-error=stderr",
			"--log-location-warning=stderr",
			"--log-location-info=stderr",
			"--dir={}".format(path.join(conf.tsroot, "etc/trafficserver"))
		]

		self.t3c_request_cmd = [
			"t3c-request",
			"--traffic-ops-url=http{}://{}:{}".format("s" if conf.useSSL else "", conf.toHost, conf.toPort),
			"--cache-host-name={}".format(self.hostname),
			"--traffic-ops-user={}".format(conf.username),
			"--traffic-ops-password={}".format(conf.password),
			"--log-location-error=stderr",
			"--log-location-debug=stderr",
			"--log-location-info=stderr"
		]

		self.t3c_update_cmd = [
			"t3c-update",
			"--traffic-ops-url=http{}://{}:{}".format("s" if conf.useSSL else "", conf.toHost, conf.toPort),
			"--cache-host-name={}".format(self.hostname),
			"--traffic-ops-user={}".format(conf.username),
			"--traffic-ops-password={}".format(conf.password),
			"--log-location-error=stderr",
			"--log-location-debug=stderr",
			"--log-location-info=stderr",
		]

		if conf.timeout is not None and conf.timeout >= 0:
			self.t3c_update_cmd.append("--traffic-ops-timeout-milliseconds={}".format(conf.timeout))
			self.t3c_request_cmd.append("--traffic-ops-timeout-milliseconds={}".format(conf.timeout))
		if conf.rev_proxy_disable:
			self.t3c_update_cmd.append("--traffic-ops-disable-proxy")
			self.t3c_request_cmd.append("--traffic-ops-disable-proxy")
		if not conf.verify:
			self.t3c_update_cmd.append("--traffic-ops-insecure")
			self.t3c_request_cmd.append("--traffic-ops-insecure")

		if conf.via_string_release > 0:
			self.t3c_generate_cmd.append("--via-string-release")
		if conf.disable_parent_config_comments > 0:
			self.t3c_generate_cmd.append("--disable-parent-config-comments")

		if conf.timeout is not None and conf.timeout >= 0:
			self.t3c_request_cmd.append("--traffic-ops-timeout-milliseconds={}".format(conf.timeout))
		if conf.rev_proxy_disable:
			self.t3c_request_cmd.append("--traffic-ops-disable-proxy")
		if not conf.verify:
			self.t3c_request_cmd.append("--traffic-ops-insecure")

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

	def get_statuses(self) -> typing.List[dict]:
		"""
		Retrieves all statuses from the Traffic Ops instance - using t3c-request.

		:returns: Representations of status objects
		:raises: ConnectionError if fetching the statuses fails for any reason
		"""

		t3c_request_cmd = self.t3c_request_cmd + ["--get-data=statuses"]

		for _ in range(self.retries):
			try:
				proc = subprocess.run(t3c_request_cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
				logging.debug("Raw response: %s", proc.stdout.decode())
				if proc.stderr.decode():
					logging.error(proc.stderr.decode())
				if proc.returncode == 0:
					return json.loads(proc.stdout.decode())
			except (subprocess.SubprocessError, OSError, json.JSONDecodeError) as e:
				logging.error("status fetch failure: %s", e)
		raise ConnectionError("Failed to fetch statuses from t3c-request")

	def getMyPackages(self) -> typing.List[packaging.Package]:
		"""
		Fetches a list of the packages specified by Traffic Ops that should exist on this server.

		:returns: all of the packages which this system must have, according to Traffic Ops.
		:raises ConnectionError: if fetching the package list fails
		"""
		logging.info("Fetching this server's package list from Traffic Ops")

		t3c_request_cmd = self.t3c_request_cmd + ["--get-data=packages"]

		for _ in range(self.retries):
			try:
				proc = subprocess.run(t3c_request_cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
				logging.debug("Raw output: %s", proc.stdout.decode())
				if proc.stderr.decode():
					logging.error("proc.stderr.decode()")
				if proc.returncode == 0:
					return [packaging.Package(p) for p in json.loads(proc.stdout.decode())]
			except (ValueError, IndexError, json.JSONDecodeError, OSError, subprocess.SubprocessError) as e:
				logging.debug("package fetch failure: %r", e, stack_info=True, exc_info=True)

		raise ConnectionError("Failed to get a response for packages")

	def getMyConfigFiles(self, conf:Configuration) -> typing.List[ConfigFile]:
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

		t3c_request_cmd = self.t3c_request_cmd + ["--get-data=config"]

		if conf.mode is Configuration.Modes.REVALIDATE:
			t3c_request_cmd = self.t3c_request_cmd + ["--reval-only"]

#		for _ in range(self.retries):
#			try:
#				proc = subprocess.run(t3c_request_cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
#				logging.debug("Raw output: %s", proc.stdout.decode())
#				if proc.stderr.decode():
#					logging.error("proc.stderr.decode()")
#				if proc.returncode == 0:
#					return [packaging.Package(p) for p in json.loads(proc.stdout.decode())]
#			except (ValueError, IndexError, json.JSONDecodeError, OSError, subprocess.SubprocessError) as e:
#				logging.debug("package fetch failure: %r", e, stack_info=True, exc_info=True)
#		raise ConnectionError("Failed to get a response for packages")

		t3c_generate_cmd = self.t3c_generate_cmd

		if conf.mode is Configuration.Modes.REVALIDATE:
			t3c_generate_cmd = self.t3c_generate_cmd + ["--revalidate-only"]

		req_out = ''
		for _ in range(self.retries):
			try:
				proc = subprocess.run(t3c_request_cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
				logging.debug("request config Raw response: %s", len(proc.stdout.decode()))
				if proc.stderr.decode():
					logging.error(proc.stderr.decode())
				logging.debug("request config code " + str(proc.returncode))
				if proc.returncode == 0:
					req_out = proc.stdout
					logging.debug("request config set req_out")
					break
			except (subprocess.SubprocessError, ValueError, OSError) as e:
				logging.debug("config file fetch failure: %r", e, exc_info=True, stack_info=True)

		logging.debug("request config output len " + str(len(req_out)))
		if len(req_out) > 0:
			try:
				logging.debug("calling t3c_generate_cmd")
				proc = subprocess.run(t3c_generate_cmd, input=req_out, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
				logging.debug("generate raw response len: %s", len(proc.stdout))
				if proc.stderr.decode():
					logging.error(proc.stderr.decode())
				logging.debug("generate config code " + str(proc.returncode))
				if proc.returncode == 0:
					return [ConfigFile(tsroot=conf.tsroot, contents=p.get("text"), path=(p.get("path") + '/' + p.get("name"))) for p in json.loads(proc.stdout.decode())]
			except (subprocess.SubprocessError, ValueError, OSError) as e:
				logging.debug("config file generate failure: %r", e, exc_info=True, stack_info=True)

		raise ConnectionError("Failed to fetch configuration files from Traffic Ops")

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

		t3c_update_cmd = self.t3c_update_cmd + ["--set-update-status=false", "--set-reval-status=false"]

		for _ in range(self.retries):
			try:
				proc = subprocess.run(t3c_update_cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
				logging.info(proc.stdout.decode())
				logging.error(proc.stderr.decode())
				if proc.returncode == 0:
					break
			except (LoginError, InvalidJSONError, OperationError, RequestException) as e:
				logging.error("TO update failure: %r", e, exc_info=True, stack_info=True)
		else:
			raise ConnectionError("Failed to update Traffic Ops - connection was lost")

	def getMyChkconfig(self) -> typing.List[dict]:
		"""
		Fetches the 'chkconfig' for this server

		:returns: An iterable list of 'chkconfig' entries
		:raises ConnectionError: when something goes wrong communicating with Traffic Ops
		"""

		logging.info("Fetching chkconfig")

		t3c_request_cmd = self.t3c_request_cmd + ["--get-data=chkconfig"]

		for _ in range(self.retries):
			try:
				proc = subprocess.run(t3c_request_cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
				logging.debug("Raw response: %s", proc.stdout.decode())
				logging.error(proc.stderr.decode())
				if proc.returncode == 0:
					return json.loads(proc.stdout.decode())
			except (json.JSONDecodeError, OSError, subprocess.SubprocessError) as e:
				logging.debug("chkconfig fetch failure: %r", e, exc_info=True, stack_info=True)

		raise ConnectionError("Failed to fetch 'chkconfig' from Traffic Ops - connection lost")

	def getMyUpdateStatus(self) -> dict:
		"""
		Gets the update status of a server.

		:raises ConnectionError: if something goes wrong communicating with the server
		:returns: An object representing the API's response
		"""
		logging.info("Fetching update status from Traffic Ops")

		t3c_request_cmd = self.t3c_request_cmd + ["--get-data=update-status"]

		for _ in range(self.retries):
			try:
				proc = subprocess.run(t3c_request_cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
				logging.debug("Raw response: %s", proc.stdout.decode())
				logging.error(proc.stderr.decode())
				if proc.returncode == 0:
					return json.loads(proc.stdout.decode())
			except (subprocess.SubprocessError, OSError, json.JSONDecodeError) as e:
				logging.debug("update status fetch failure: %r", e, exc_info=True, stack_info=True)

		raise ConnectionError("Failed to fetch update status - connection was lost")

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
