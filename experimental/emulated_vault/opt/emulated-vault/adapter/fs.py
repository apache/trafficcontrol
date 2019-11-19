#
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#	 http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
#
#

import os
from . base import Base

class Fs(Base):
	"""
	Fs (file system) adapter class.
	This class implements the API required for storing and retriving the content kept in the vault.
	This specific Adapter works using files upon a file-system. 
	Inputs configuration (held under "fs-adapter" section in the config file) include:
	:param db-base-os-path: The path in which the DB files are stored
	:param ping-os-path: The path of a variable mimicing the RIAK ping functionality

	Interface methods 
	:meth:`_get_parameter_storage_path` given a url-key of a parameter return the storage path
	:meth:`_get_parameter_storage_path_from_url_path` given a storage path of a parameter
	retrun the url-key
	:meth:`_init_cfg` given a config-parser object, read the parameters required for 
	adapter's operation. Return "success" boolean value
	:meth:`_init_ping` prepare the adapter for beign able to provide info to reply for 
	"ping" requests. Return "success" boolean value
	:meth:`_ping` test the 'ping' status with the adapter. 
	Return a tuple: "success" boolean & "value" kept as ping variable
	:meth:`_read_parameter_by_storage_path` given a storage path retrieve the parameter value.
	Return a tuple: "success" boolean & "value" kept in the parameter
	:meth:`_read_parameters_by_storage_path` given a storage and and a key holding 
	filters on the key and values.
	Return "success" boolean indicating a sucessful write, and a key->value dictionary 
	for the relevant parameters
	:meth:`_write_parameter_by_storage_path` given a storage path and a value string,
	keep the parameter value. Return "success" boolean indicating a sucessful write
	:meth:`_remove_parameter_by_storage_path` given a storage path, delete the parameter
	from the DB. Return "success" boolean indicating a sucessful deletion		
	"""
	def __init__ (self, logger):
		"""
		The class constructor.
		:param logger: logger to send log messages to
		:type logger: a python logging logger class
		"""
		Base.__init__(self, logger)

	def _init_cfg(self, fullConfig):
		"""
		Initialize the class basic parameters. Part of Adapter required API.
		:param fullConfig: configuration to operate upon.
		:type fullConfig: configparser.ConfigParser class
		:return: 'True' for successful initialization
		:rtype: Boolean
		"""
		myCfgData = dict(fullConfig.items("fs-adapter")) if fullConfig.has_section("fs-adapter") else {}
		self.basePath = myCfgData.get("db-base-os-path")
		if self.basePath is None:
			self.logger.error("Missing %s/%s configuration", "fs-adapter", "db-base-os-path")
			return False 
		self.pingStoragePath = myCfgData.get("ping-os-path")
		if self.pingStoragePath is None:
			self.logger.error("Missing %s/%s configuration", "fs-adapter", "ping-os-path")
			return False 
		return True

	def _get_parameter_storage_path(self, parameterUrlPath):
		"""
		Conversion function - taking a key's path and translate to a file path on the file system
		:param parameterUrlPath: the "url-path" like key of the variable
		:type parameterUrlPath: str
		:return: file path of where the value is be kept
		:rtype: str
		"""
		return os.path.join(self.basePath, parameterUrlPath.lstrip("/").replace("/", os.path.sep))

	def _get_parameter_storage_path_from_url_path(self, parameterStoragePath):
		"""
		Conversion function - taking file path on the file system and translate to key's path
		:param parameterStoragePath: the file name holding a value
		:type parameterUrlPath: str
		:return: the matching variable url-path like key
		:rtype: str
		"""
		return "/"+os.path.relpath(parameterStoragePath, self.basePath).replace(os.path.sep, "/")

	def _init_ping(self):
		"""
		Initialize the class ability to answer for ping. Part of Adapter required API.
		:return: 'True' for successful initialization
		:rtype: Boolean
		"""
		value = "OK"
		success = self._write_parameter_by_storage_path(self.pingStoragePath, value)
		if not success:
			self.logger.error("Failed to set parameter %s", self.pingStoragePath)
			return False
		return True

	def _ping(self):
		"""
		get value for the ping request. Part of Adapter required API.
		:return: A tuple - 'True' for successful retrival and the retrieved value
		:rtype: Tuple[Boolean, str]
		"""
		success, value = self._read_parameter_by_storage_path(self.pingStoragePath)
		if not success or value is None:
			self.logger.error("no ping response")
			return (False, "")
		self.logger.debug("ping response: %s", value)
		return (True, value)

	def _read_parameter_by_storage_path(self, parameterStoragePath):
		"""
		Reading the value from the provided file name.
		:param parameterStoragePath: the file name 
		:type parameterStoragePath: str
		:return: 'True' for successful retrivaland the retrieved value
		:rtype: Tuple[boolean, str]
		"""
		self.logger.debug("Get parameter by os path: %s", parameterStoragePath)
		try:
			with open(parameterStoragePath) as fd:
				value = fd.read()
		except Exception as e:
			self.logger.exception("%s parameter os path not found.", parameterStoragePath)
			return False, None

		self.logger.debug("Get parameter by os path %s succeed", parameterStoragePath)
		return True, value

	def _read_parameters_by_storage_path(self, parameterStoragePathPrefix, keyFilters):
		"""
		Reading the values of the parameters the provided directory.
		:param parameterStoragePathPrefix: the directory to look into
		:type parameterStoragePathPrefix: str
		:param keyFilters: filter-name/filter-func dict, holding functions that get a key as 
		input and retunn "true" if key should be included in the result
		:type keyFilters: Dict[str,function[str]]
		:return: 'True' for successful retrival and a dict for key-name/value 
		:rtype: Tuple[boolean, Dict[str, str]]
		"""
		self.logger.debug("Get parameters under os path %s", parameterStoragePathPrefix)
		parameters = {}

		fileNames = []
		for (dirpath, _, filenames) in os.walk(parameterStoragePathPrefix):
			fileNames += [os.path.join(dirpath, file) for file in filenames]

		for fileName in fileNames:
			filteredOut = False
			for filterName, filterfunc in keyFilters.items():#items() - supporting python 2&3
				if not filterfunc(fileName):
					self.logger.debug("Parameter os path %s dropped, not matching filter %s", self._get_parameter_storage_path_from_url_path(fileName), filterName)
					filteredOut = True
					break
			if filteredOut:
				continue

			parameterUrlPath = self._get_parameter_storage_path_from_url_path(fileName)
			success, value = self._read_parameter_by_storage_path(fileName)
			if not success:
				self.logger.error("%s parameter os path not found.", fileName)
				return False, None
				
			parameters[parameterUrlPath] = value

		return True, parameters


	def _write_parameter_by_storage_path(self, parameterStoragePath, value):
		"""
		Writing the value to the provided file name.
		:param parameterStoragePath: the file name 
		:type parameterStoragePath: str
		:param value: value to be writen
		:type value: str
		:return: 'True' for successful writing
		:rtype: Boolean		
		"""
		self.logger.debug("Set parameter by os path %s", parameterStoragePath)
		try:
			dirname = os.path.dirname(parameterStoragePath)
			if dirname and not os.path.exists(dirname):
				os.makedirs(dirname)
			with open(parameterStoragePath, "w") as fd:
				fd.write(value)			
		except Exception as e:
			self.logger.exception("could not post parameter os path %s", parameterStoragePath)
			return False
		self.logger.debug("Set parameter os path %s done", parameterStoragePath)
		return True

	def _remove_parameter_by_storage_path(self, parameterStoragePath):
		"""
		Deleting the the provided file.
		:param parameterStoragePath: the file name 
		:type parameterStoragePath: str
		:return: 'True' for successful deletion
		:rtype: Boolean		
		"""
		self.logger.debug("Delete parameter os path %s", parameterStoragePath)
		try:
			os.remove(_parameterStoragePath)
		except Exception as e:
			self.logger.exception("could not delete parameter os path %s", parameterStoragePath)
			return False
		self.logger.debug("Delete parameter os path %s done", parameterStoragePath)
		return True

