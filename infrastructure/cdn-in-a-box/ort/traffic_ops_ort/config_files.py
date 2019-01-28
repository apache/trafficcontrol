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

from trafficops.restapi import OperationError, InvalidJSONError

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

	def __init__(self, raw:dict = None):
		"""
		Constructs a :class:`ConfigFile` object from a raw API response

		:param raw: A raw config file from an API response
		:raises ValueError: if ``raw`` does not faithfully represent a configuration file

		>>> ConfigFile({"fnameOnDisk": "test",
		...             "location": "/path/to",
		...             "apiURI": "http://test",
		...             "scope": "servers"}))
		ConfigFile(path='/path/to/test', URI='http://test', scope='servers')
		"""
		# TODO: pass these in as parameters? Configuration object?
		from .configuration import TO_HOST, TO_PORT, TO_USE_SSL, TS_ROOT

		if raw is not None:
			try:
				self.fname = raw["fnameOnDisk"]
				self.location = raw["location"]
				if "apiUri" in raw:
					self.URI = "https://" if TO_USE_SSL else "http://"
					self.URI = "%s%s:%d/%s" % (self.URI, TO_HOST, TO_PORT, raw["apiUri"].lstrip('/'))
				else:
					self.URI = raw["url"]
				self.scope = raw["scope"]
			except (KeyError, TypeError, IndexError) as e:
				raise ValueError from e

		self.SSLdir = os.path.join(TS_ROOT, "etc", "trafficserver", "ssl")

	def __repr__(self) -> str:
		"""
		Implements ``repr(self)``

		>>> repr(ConfigFile({"fnameOnDisk": "test",
		...                  "location": "/path/to",
		...                  "apiURI": "http://test",
		...                  "scope": "servers"}))
		"ConfigFile(path='/path/to/test', URI='http://test', scope='servers')"
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

	def backup(self, contents:str):
		"""
		Creates a backup of this file under the :data:`BACKUP_DIR` directory

		:param contents: The actual, on-disk contents from the original file
		:raises OSError: if the backup directory does not exist, or a backup of this file
			could not be written into it.
		"""
		from .configuration import MODE, Modes
		from .utils import getYesNoResponse

		backupfile = os.path.join(BACKUP_DIR, self.fname)
		willClobber = False
		if os.path.isfile(backupfile):
			willClobber = True

		if MODE is Modes.INTERACTIVE:
			prmpt = ("Write backup file %s%%s?" % backupfile)
			prmpt %= " - will clobber existing file by the same name - " if willClobber else ''
			if not getYesNoResponse(prmpt, default='Y'):
				return

		elif willClobber:
			logging.warning("Clobbering existing backup file '%s'!", backupfile)

		if MODE is not Modes.REPORT:
			with open(backupfile, 'w') as fp:
				fp.write(contents)

		logging.info("Backup File written")


	def update(self, api:'to_api.API', cdn:str):
		"""
		Updates the file if required, backing up as necessary

		:param api: A valid, authenticated API session for use when interacting with Traffic Ops
		:param cdn: The name of the CDN to which this server belongs (needed for SSL keys)
		:raises OSError: when reading/writing files fails for some reason
		"""
		from . import utils
		from .configuration import MODE, Modes, SERVER_INFO
		from .services import NEEDED_RELOADS, FILES_THAT_REQUIRE_RELOADS

		if not self.contents:
			self.fetchContents(api)
			finalContents = sanitizeContents(self.contents)
		else:
			finalContents = self.contents

		# Ensure POSIX-compliant files
		if not finalContents.endswith('\n'):
			finalContents += '\n'
		logging.info("Sanitized output: \n%s", finalContents)
		self.sanitizedContents = finalContents

		if not os.path.isdir(self.location):
			if MODE is Modes.INTERACTIVE and\
			   not utils.getYesNoResponse("Create configuration directory %s?" % self.path, 'Y'):
				logging.warning("%s will not be created - some services may not work properly!",
				                self.path)
				return

			logging.info("Directory %s will be created", self.location)
			logging.info("File %s will be created", self.path)

			if MODE is not Modes.REPORT:
				os.makedirs(self.location)
				with open(self.path, 'x') as fp:
					fp.write(finalContents)
				return

		if not os.path.isfile(self.path):
			if MODE is Modes.INTERACTIVE and\
			   not utils.getYesNoResponse("Create configuration file %s?"%self.path, default='Y'):
				logging.warning("%s will not be created - some services may not work properly!",
				                self.path)
				return

			logging.info("File %s will be created", self.path)

			if MODE is not Modes.REPORT:
				with open(self.path, 'x') as fp:
					fp.write(finalContents)
				return

		with open(self.path, 'r+') as fp:
			onDiskContents = fp.readlines()
			if filesDiffer(finalContents.splitlines(), onDiskContents):
				self.backup(''.join(onDiskContents))
				if MODE is not Modes.REPORT:
					fp.seek(0)
					fp.truncate()


					fp.write(finalContents)
					if self.fname in FILES_THAT_REQUIRE_RELOADS:
						NEEDED_RELOADS.add(FILES_THAT_REQUIRE_RELOADS[self.fname])
				logging.info("File written to %s", self.path)
			else:
				logging.info("File doesn't differ from disk; nothing to do")

		# Now we need to do some advanced processing to a couple specific filenames... unfortunately
		if self.fname == "ssl_multicert.config":
			self.advancedSSLProcessing(api, cdn)

	def advancedSSLProcessing(self, api:'to_api.API', cdn:str):
		"""
		Does advanced processing on ssl_multicert.config files

		:param api: A valid, authenticated API session for use when interacting with Traffic Ops
		:param cdn: The name of the CDN to which this server belongs (needed for SSL keys)
		:raises OSError: when reading/writing files fails for some reason
		"""
		global SSL_KEY_REGEX

		logging.info("Doing advanced SSL key processing for CDN '%s'", cdn)

		try:
			r = api.get_cdn_ssl_keys(cdn_name=cdn)

			if r[1].status_code != 200 and r[1].status_code != 204:
				raise ValueError("Bad response code: %d - raw response: %s" %
				                               (r[1].status_code,    r[1].text))
		except (OperationError, InvalidJSONError, ValueError) as e:
			logging.error("Invalid values encountered when communicating with Traffic Ops!")
			logging.debug("%r", e, stack_info=True, exc_info=True)
			raise ValueError from e

		logging.debug("Raw response from Traffic Ops: %s", r[1].text)

		for l in self.sanitizedContents.splitlines()[1:]:
			logging.debug("advanced processing for line: %s", l)
			m = SSL_KEY_REGEX.search(l)

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
					key = type(self)()
					key.location = self.SSLdir
					key.fname = m.group(2)
					key.contents = b64decode(cert.certificate.key).decode()

					logging.info("Processing private SSL key %s ...", key.fname)
					key.update(api, cdn)
					logging.info("Done.")

					crt = type(self)()
					crt.location = self.SSLdir
					crt.fname = m.group(1)
					crt.contents = b64decode(cert.certificate.crt).decode()

					logging.info("Processing SSL certificate %s ...", crt.fname)
					crt.update(api, cdn)
					logging.info("Done.")
					break
			else:
				logging.critical("Failed to find SSL key in %s for '%s' or by wildcard '%s'!",
				                                           cdn,    full,            wildcard)
				raise ValueError("No cert/key pair for ssl_multicert.config line '%s'" % l)

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

def sanitizeContents(raw:str) -> str:
	"""
	Performs pre-processing on a raw configuration file

	:param raw: The raw contents of the file as returned by a request to its URL
	:returns: The same contents, but with special replacement strings parsed out and HTML-encoded
		symbols decoded to their literal values
	"""
	from .configuration import SERVER_INFO
	out = []

	# These double curly braces escape the behaviour of Python's `str.format` method to attempt
	# to use curly brace-enclosed text as a key into a dictonary of its arguments. They'll be
	# rendered into single braces in the output of `.format`, leaving the string ultimately
	# unchanged in that respect.
	for line in SERVER_INFO.sanitize(raw).splitlines():
		tmp=(" ".join(line.split())).strip() #squeezes spaces and trims leading and trailing spaces
		tmp=tmp.replace("&amp;", '&') #decodes HTML-encoded ampersands
		tmp=tmp.replace("&gt;", '>') #decodes HTML-encoded greater-than symbols
		tmp=tmp.replace("&lt;", '<') #decodes HTML-encoded less-than symbols
		out.append(tmp)

	return '\n'.join(out)

def initBackupDir():
	"""
	Initializes a backup directory as a subdirectory of the directory containing
	this ORT script.

	:raises OSError: if the backup directory initialization fails
	"""
	global BACKUP_DIR
	from . import configuration as conf

	logging.info("Initializing backup dir %s", BACKUP_DIR)

	if not os.path.isdir(BACKUP_DIR):
		if conf.MODE != conf.Modes.REPORT:
			os.mkdir(BACKUP_DIR)
		else:
			logging.error("Cannot create non-existent backup dir in REPORT mode!")
	else:
		logging.info("Backup dir already exists - nothing to do")
