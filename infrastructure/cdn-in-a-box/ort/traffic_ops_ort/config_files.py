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

import os
import logging
import typing

#: Holds a set of service names that need reloaded configs, mapped to a boolean which indicates
#: whether (:const:`True`) or not (:const:`False`) a full restart is required.
RELOADS_REQUIRED = set()

#: A constant that holds the absolute path to the backup directory for configuration files
BACKUP_DIR = "/opt/ort/backups"

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

	def __init__(self, raw:dict):
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
		from .configuration import TO_URL

		try:
			self.fname = raw["fnameOnDisk"]
			self.location = raw["location"]
			if "apiUri" in raw:
				self.URI = '/'.join((TO_URL, raw["apiUri"].lstrip('/')))
			else:
				self.URI = raw["url"]
			self.scope = raw["scope"]
		except (KeyError, TypeError, IndexError) as e:
			raise ValueError from e

		self.path = os.path.join(self.location, self.fname)

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

	def __str__(self) -> str:
		"""
		Implements ``str(self)``

		:returns: The contents of the file, fetched using :meth:`fetchContents` if necessary
		"""
		if self.contents:
			return self.contents

		try:
			self.fetchContents()
		except ConnectionError as e:
			logging.debug("%s", e, exc_info=True, stack_info=True)

		return self.contents

	def fetchContents(self):
		"""
		Fetches the file contents from :attr:`URI` if possible. Sets :attr:`contents` when
		successful.

		:raises ConnectionError: if something goes wrong fetching the file contents from Traffic
			Ops
		"""
		from . import configuration as conf, utils
		logging.info("Fetching file %s", self.fname)

		try:
			# You can't really check this any earlier, since even a 'url' field could, potentially,
			# point at Traffic Ops
			if self.URI.startswith(conf.TO_URL):
				cookies = {conf.TO_COOKIE.name: conf.TO_COOKIE.value}
			else:
				cookies = None

			self.contents = utils.getTextResponse(self.URI, cookies=cookies, verify=conf.VERIFY)
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


	def update(self):
		"""
		Updates the file if required, backing up as necessary

		:raises OSError: when reading/writing files fails for some reason
		"""
		from . import utils
		from .configuration import MODE, Modes, SERVER_INFO
		from .services import NEEDED_RELOADS, FILES_THAT_REQUIRE_RELOADS

		finalContents = sanitizeContents(str(self))
		# Ensure POSIX-compliant files
		if not finalContents.endswith('\n'):
			finalContents += '\n'
		logging.info("Sanitized output: \n%s", finalContents)

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
	for line in raw.replace('{', "{{").replace('}', "}}").format(SERVER_INFO).splitlines():
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
