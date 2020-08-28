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
from . adapter_base import AdapterBase

class FsAdapter(AdapterBase):
	"""
	Fs (file system) adapter class.
	This class implements the API required for storing and retriving the content kept in the vault.
	This specific Adapter works using files upon a file-system. 
	Inputs configuration (held under "fs-adapter" section in the config file) include:
	:param db-base-os-path: The path in which the DB files are stored
	:param ping-os-path: The path of a variable mimicing the RIAK ping functionality
	"""
	
	def __init__ (self, logger):
		"""
		The class constructor.
		:param logger: logger to send log messages to
		:type logger: a python logging logger class
		"""
		self.logger = logger

	def init_cfg(self, fullConfig):# -> bool:
		"""
		Initialize the class basic parameters. Part of Adapter required API.
		Read the relevant section in the configuration
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

	def init(self):# -> bool:
		"""
		Initialize the class ability to answer for ping. Part of Adapter required API.
		"""
		value = ":)"
		success = self.write_parameter_by_storage_path(self.pingStoragePath, value)
		if not success:
			self.logger.error("Failed to set parameter %s", self.pingStoragePath)
			return False
		return True

	def get_parameter_storage_path(self, parameterUrlPath):# -> bool, str:
		"""
		Conversion function - taking a key's path and translate to a file path on the file system
		"""
		#verifying the existance of "/" at the beging
		if not parameterUrlPath.startswith("/"):
			self.logger.error("Failed to translate url path '%s': no leading '/'", parameterUrlPath)
			return False, ""
		# avoiding the use of os.path.join in purpuse. 
		# we do not want the "/" at the begining push us to the root path		
		return True, self.basePath + parameterUrlPath.replace("/", os.path.sep)

	def get_parameter_url_path_from_storage_path(self, parameterStoragePath):# -> str:
		"""
		Conversion function - taking file path on the file system and translate to key's path
		"""
		return True, "/"+os.path.relpath(parameterStoragePath, self.basePath).replace(os.path.sep, "/")

	def ping(self):# -> bool:
		"""
		Get value for the ping request by its path. Part of Adapter required API.
		"""
		success, value = self.read_parameter_by_storage_path(self.pingStoragePath)
		if not success or value is None:
			self.logger.error("no ping response")
			return False
		self.logger.debug("ping response: %s", value)
		return True

	def read_parameter_by_storage_path(self, parameterStoragePath):# -> (bool, str):
		"""
		Reading the value from the provided file name.
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

	def read_parameters_by_storage_path(self, parameterStoragePathPrefix, keyFilters):# -> (bool, dict(str, str)):
		"""
		Reading the values of the parameters the provided directory.
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
					self.logger.debug("Parameter os path %s dropped, not matching filter '%s'", self.get_parameter_url_path_from_storage_path(fileName)[1], filterName)
					filteredOut = True
					break
			if filteredOut:
				continue

			success_path, parameterUrlPath = self.get_parameter_url_path_from_storage_path(fileName)
			if not success_path:
				self.logger.error("%s parameter os path is invalid.", fileName)
				return False, None				
			success, value = self.read_parameter_by_storage_path(fileName)
			if not success:
				self.logger.error("%s parameter os path not found.", fileName)
				return False, None
				
			parameters[parameterUrlPath] = value

		return True, parameters


	def write_parameter_by_storage_path(self, parameterStoragePath, value):# -> bool:
		"""
		Writing the value to the provided file name.
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

	def remove_parameter_by_storage_path(self, parameterStoragePath):# -> bool:
		"""
		Deleting the the provided file.
		"""
		self.logger.debug("Delete parameter os path %s", parameterStoragePath)
		try:
			os.remove(_parameterStoragePath)
		except Exception as e:
			self.logger.exception("could not delete parameter os path %s", parameterStoragePath)
			return False
		self.logger.debug("Delete parameter os path %s done", parameterStoragePath)
		return True

