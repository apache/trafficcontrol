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
This module deals with the management of configuration files,
presumably for a cache server
"""

import logging
import os
import re
import typing

from base64 import b64decode

from trafficops.restapi import OperationError, InvalidJSONError, LoginError

from .configuration import Configuration
from .utils import getYesNoResponse as getYN


#: Holds a set of service names that need reloaded configs, mapped to a boolean which indicates
#: whether (:const:`True`) or not (:const:`False`) a full restart is required.
RELOADS_REQUIRED = set()

#: A constant that holds the absolute path to the backup directory for configuration files
BACKUP_DIR = "/opt/ort/backups"

#: a pre-compiled regular expression to use in parsing
SSL_KEY_REGEX = re.compile(r'^\s*ssl_cert_name\=(.*)\s+ssl_key_name\=(.*)\s*$')

class ConfigurationError(Exception):
	"""
	Represents an error updating configuration files
	"""
	pass

class ConfigFile():
	"""
	Represents a configuration file on a host system.
	"""

	fname = ""    #: The base name of the file
	location = "" #: An absolute path to the directory containing the file
	URI = ""      #: A URI where the actual file contents can be found
	contents = "" #: The full contents of the file - as configured in TO, not the on-disk contents
	sanitizedContents = "" #: Will store the contents after sanitization

	def __init__(self, raw:dict = None, toURL:str = "", tsroot:str = "/", *unused_args, contents: str = None, path: str = None):
		"""
		Constructs a :class:`ConfigFile` object from a raw API response

		:param raw: A raw config file from an API response
		:param toURL: The URL of a valid Traffic Ops host
		:param tsroot: The absolute path to the root of an Apache Traffic Server installation
		:param contents: Directly constructs a ConfigFile from the passed contents. Must be used with path, and causes raw to be ignored.
		:param path: Sets the full path to the file when constructing ConfigFiles directly from contents.
		:raises ValueError: if ``raw`` does not faithfully represent a configuration file

		>>> a = ConfigFile({"fnameOnDisk": "test",
		...                 "location": "/path/to",
		...                 "apiUri":"/test",
		...                 "scope": "servers"}, "http://example.com/")
		>>> a
		ConfigFile(path='/path/to/test', URI='http://example.com/test', scope='servers')
		>>> a.SSLdir
		'/etc/trafficserver/ssl'
		>>> ConfigFile(contents='testquest', path='/path/to/test')
		ConfigFile(path='/path/to/test', URI=None, scope=None)
		"""
		self.SSLdir = os.path.join(tsroot, "etc", "trafficserver", "ssl")

		if contents is not None:
			if path is None:
				raise ValueError("cannot construct from direct contents without setting path")
			self.location, self.fname = os.path.split(path)
			self.contents = contents
			self.scope = None
			return
		if raw is not None:
			try:
				self.fname = raw["fnameOnDisk"]
				self.location = raw["location"]
				if "apiUri" in raw:
					self.URI = toURL + raw["apiUri"].lstrip('/')
				else:
					self.URI = raw["url"]
				self.scope = raw["scope"]
			except (KeyError, TypeError, IndexError) as e:
				raise ValueError from e

	def __repr__(self) -> str:
		"""
		Implements ``repr(self)``

		>>> repr(ConfigFile({"fnameOnDisk": "test",
		...                  "location": "/path/to",
		...                  "apiUri": "test",
		...                  "scope": "servers"}, "http://example.com/"))
		"ConfigFile(path='/path/to/test', URI='http://example.com/test', scope='servers')"
		"""
		return "ConfigFile(path=%r, URI=%r, scope=%r)" %\
		          (self.path, self.URI if self.URI else None, self.scope)

	@property
	def path(self) -> str:
		"""
		The full path to the file on disk

		:returns: The system's path separator concatenating :attr:`location` and :attr`fname`
		"""
		return os.path.join(self.location, self.fname)

	def fetchContents(self, api:'to_api.API'):
		"""
		Fetches the file contents from :attr:`URI` if possible. Sets :attr:`contents` when
		successful.

		:param api: A valid, authenticated API session for use when interacting with Traffic Ops
		:raises ConnectionError: if something goes wrong fetching the file contents from Traffic
			Ops
		"""
		logging.info("Fetching file %s", self.fname)

		try:
			self.contents = api.getRaw(self.URI)
		except ValueError as e:
			raise ConnectionError from e

		logging.info("fetched")

	def backup(self, contents:str, mode:Configuration.Modes):
		"""
		Creates a backup of this file under the :data:`BACKUP_DIR` directory

		:param contents: The actual, on-disk contents from the original file
		:param mode: The current run-mode of :program:`traffic_ops_ort`
		:raises OSError: if the backup directory does not exist, or a backup of this file
			could not be written into it.
		"""
		backupfile = os.path.join(BACKUP_DIR, self.fname)
		willClobber = False
		if os.path.isfile(backupfile):
			willClobber = True

		if mode is Configuration.Modes.INTERACTIVE:
			prmpt = ("Write backup file %s%%s?" % backupfile)
			prmpt %= " - will clobber existing file by the same name - " if willClobber else ''
			if not getYN(prmpt, default='Y'):
				return

		elif willClobber:
			logging.warning("Clobbering existing backup file '%s'!", backupfile)

		if mode is not Configuration.Modes.REPORT:
			with open(backupfile, 'w') as fp:
				fp.write(contents)

		logging.info("Backup File written")


	def update(self, conf:Configuration) -> bool:
		"""
		Updates the file if required, backing up as necessary

		:param conf: An object that represents the configuration of :program:`traffic_ops_ort`
		:returns: whether or not the file on disk actually changed
		:raises OSError: when reading/writing files fails for some reason
		"""
		from .services import NEEDED_RELOADS, FILES_THAT_REQUIRE_RELOADS

		if not self.contents:
			self.fetchContents(conf.api)
			finalContents = sanitizeContents(self.contents, conf)
		elif self.URI:
			finalContents = self.contents
		else:
			finalContents = sanitizeContents(self.contents, conf)

		# Ensure POSIX-compliant files
		if not finalContents.endswith('\n'):
			finalContents += '\n'
		logging.info("Sanitized output: \n%s", finalContents)
		self.sanitizedContents = finalContents

		if not os.path.isdir(self.location):
			if (conf.mode is Configuration.Modes.INTERACTIVE and
			    not getYN("Create configuration directory %s?" % self.path, 'Y')):
				logging.warning("%s will not be created - some services may not work properly!",
				                self.path)
				return False

			logging.info("Directory %s will be created", self.location)
			logging.info("File %s will be created", self.path)

			if conf.mode is not Configuration.Modes.REPORT:
				os.makedirs(self.location)
				with open(self.path, 'x') as fp:
					fp.write(finalContents)
			return True

		if not os.path.isfile(self.path):
			if (conf.mode is Configuration.Modes.INTERACTIVE and\
			    not getYN("Create configuration file %s?"%self.path, default='Y')):
				logging.warning("%s will not be created - some services may not work properly!",
				                self.path)
				return False

			logging.info("File %s will be created", self.path)

			if conf.mode is not Configuration.Modes.REPORT:
				with open(self.path, 'x') as fp:
					fp.write(finalContents)

			if self.fname == "ssl_multicert.config":
				return self.advancedSSLProcessing(conf)
			return True

		written = False
		with open(self.path, 'r+') as fp:
			onDiskContents = fp.readlines()
			if filesDiffer(finalContents.splitlines(), onDiskContents):
				self.backup(''.join(onDiskContents), conf.mode)
				if conf.mode is not Configuration.Modes.REPORT:
					fp.seek(0)
					fp.truncate()


					fp.write(finalContents)

				written = True
				logging.info("File written to %s", self.path)
			else:
				logging.info("File doesn't differ from disk; nothing to do")

		# Now we need to do some advanced processing to a couple specific filenames... unfortunately
		# But ONLY if the object wasn't directly constructed.
		if self.fname == "ssl_multicert.config" and self.URI:
			return self.advancedSSLProcessing(conf) or written

		return written

	def advancedSSLProcessing(self, conf:Configuration):
		"""
		Does advanced processing on ssl_multicert.config files

		:param conf: An object that represents the configuration of :program:`traffic_ops_ort`
		:raises OSError: when reading/writing files fails for some reason
		"""
		global SSL_KEY_REGEX

		logging.info("Doing advanced SSL key processing for CDN '%s'", conf.ServerInfo.cdnName)

		try:
			r = conf.api.get_cdn_ssl_keys(cdn_name=conf.ServerInfo.cdnName)

			if r[1].status_code != 200 and r[1].status_code != 204:
				raise OSError("Bad response code: %d - raw response: %s" %
				                               (r[1].status_code,    r[1].text))
		except (OperationError, LoginError, InvalidJSONError, ValueError) as e:
			raise OSError("Invalid values encountered when communicating with Traffic Ops!") from e

		logging.debug("Raw response from Traffic Ops: %s", r[1].text)

		written = False
		for l in self.sanitizedContents.splitlines()[1:]:
			logging.debug("advanced processing for line: %s", l)

			# for some reason, pylint is detecting this regular expression as a string
			#pylint: disable=E1101
			m = SSL_KEY_REGEX.search(l)
			#pylint: enable=E1101

			if m is None:
				continue

			full = m.group(2)
			if full.endswith(".key"):
				full = full[:-4]

			wildcard = full.find('.')
			if wildcard >= 0:
				wildcard = '*'+full[wildcard:]
			else:
				# Not sure if this is a reasonable default - if there's no '.' in the key name,
				# then there's probably something wrong...
				wildcard = "*." + full

			logging.debug("Searching for '%s' or '%s' matches", full, wildcard)

			for cert in r[0]:
				if cert.hostname == full or cert.hostname == wildcard:
					key = ConfigFile()
					key.location = self.SSLdir
					key.fname = m.group(2)
					key.contents = b64decode(cert.certificate.key).decode()

					logging.info("Processing private SSL key %s ...", key.fname)
					written = key.update(conf)
					logging.info("Done.")

					crt = ConfigFile()
					crt.location = self.SSLdir
					crt.fname = m.group(1)
					crt.contents = b64decode(cert.certificate.crt).decode()

					logging.info("Processing SSL certificate %s ...", crt.fname)
					written = crt.update(conf)
					logging.info("Done.")
					break
			else:
				logging.critical("Failed to find SSL key in %s for '%s' or by wildcard '%s'!",
				                         conf.ServerInfo.cdnName,  full,            wildcard)
				raise OSError("No cert/key pair for ssl_multicert.config line '%s'" % l)

		# If even one key was written, we need to make ATS aware of the configuration changes
		return written

def filesDiffer(a:typing.List[str], b:typing.List[str]) -> bool:
	"""
	Compares two files for meaningingful differences. Traffic Ops Headers are
	stripped out of the file contents before comparison. Trailing whitespace
	is ignored

	:param a: The contents of the first file, as a list of its lines
	:param b: The contents of the second file, as a list of its lines
	:returns: :const:`True` if the files have any differences, :const:`False`
	"""
	a = [l.rstrip() for l in a if l.rstrip() and not l.startswith("# DO NOT EDIT") and\
	                                             not l.startswith("# TRAFFIC OPS NOTE:")]
	b = [l.rstrip() for l in b if l.rstrip() and not l.startswith("# DO NOT EDIT") and\
	                                             not l.startswith("# TRAFFIC OPS NOTE:")]

	if len(a) != len(b):
		return True

	for i, l in enumerate(a):
		if l != b[i]:
			return True

	return False

def sanitizeContents(raw:str, conf:Configuration) -> str:
	"""
	Performs pre-processing on a raw configuration file

	:param raw: The raw contents of the file as returned by a request to its URL
	:param conf: An object that represents the configuration of :program:`traffic_ops_ort`
	:returns: The same contents, but with special replacement strings parsed out and HTML-encoded
		symbols decoded to their literal values
	"""
	out = []

	lines = (conf.ServerInfo.sanitize(raw, conf.hostname) if conf.ServerInfo else raw).splitlines()
	for line in lines:
		tmp=(" ".join(line.split())).strip() #squeezes spaces and trims leading and trailing spaces
		tmp=tmp.replace("&amp;", '&') #decodes HTML-encoded ampersands
		tmp=tmp.replace("&gt;", '>') #decodes HTML-encoded greater-than symbols
		tmp=tmp.replace("&lt;", '<') #decodes HTML-encoded less-than symbols
		out.append(tmp)

	return '\n'.join(out)

def initBackupDir(mode:Configuration.Modes):
	"""
	Initializes a backup directory as a subdirectory of the directory containing
	this ORT script.

	:param mode: The current run-mode of :program:`traffic_ops_ort`
	:raises OSError: if the backup directory initialization fails
	"""
	global BACKUP_DIR

	logging.info("Initializing backup dir %s", BACKUP_DIR)

	if not os.path.isdir(BACKUP_DIR):
		if mode is not Configuration.Modes.REPORT:
			os.mkdir(BACKUP_DIR)
		else:
			logging.error("Cannot create non-existent backup dir in REPORT mode!")
	else:
		logging.info("Backup dir already exists - nothing to do")
