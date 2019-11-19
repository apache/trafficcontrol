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

class Base(object):
	"""
	Base adapter class.
	This class implements the API required for storing and retriving the content kept in the vault.
	
	Implemented methods 
	:meth:`init_cfg` given a config-parser object, read the parameters required for 
	adapter's operation. Return "success" boolean value
	:meth:`init_ping` prepare the adapter for beign able to provide info to reply for 
	"ping" requests. Return "success" boolean value
	:meth:`ping` test the 'ping' status with the adapter. 
	Return a tuple: "success" boolean & "value" kept as ping variable
	:meth:`getParameter` given a parameter key (in url-path format) retrieve the parameter value.
	Return a tuple: "success" boolean & "value" kept in the parameter
	:meth:`searchParameter` given a parameters key prefix (in url-path format) and, a dict holding
	variable key filters, and a key holding filters on the values as well.
	Return "success" boolean indicating a sucessful write, and a key->value dictionary 
	for the relevant parameters
	:meth:`setParameter` given a parameter key (in url-path format) and a value string,
	keep the parameter value. Return "success" boolean indicating a sucessful write
	:meth:`deleteParameter` given a parameter key (in url-path format), delete the parameter
	from the DB. Return "success" boolean indicating a sucessful deletion		

	Methods to be implemented at derived classes:
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
		
		self.logger = logger

	def init_cfg(self, fullConfig):
		"""
		Initialize the class basic parameters. Part of Adapter required API.
		:param fullConfig: configuration to operate upon.
		:type fullConfig: configparser.ConfigParser class
		:return: 'True' for successful initialization
		:rtype: Boolean
		"""
		return self._init_cfg(fullConfig)

	def init_ping(self):
		"""
		Initialize the class ability to answer for ping. Part of Adapter required API.
		:return: 'True' for successful initialization
		:rtype: Boolean
		"""
		return self._init_ping()

	def ping(self):
		"""
		get value for the ping request. Part of Adapter required API.
		:return: A tuple - 'True' for successful retrival and the retrieved value
		:rtype: Tuple[Boolean, str]
		"""
		return self._ping()

	def getParameter(self, parameterUrlPath):
		"""
		Get value for the specified parameter. Part of Adapter required API.
		:param parameterUrlPath: the key of the parameter as presented as url path 
		(tokens seperated by "/", with "/" as a prefix)
		:type parameterUrlPath: str
		:return: A tuple - 'True' for successful retrival and the retrieved value
		:rtype: Tuple[Boolean, str]
		"""
		parameterStoragePath = self._get_parameter_storage_path(parameterUrlPath)		
		success, value = self._read_parameter_by_storage_path(parameterStoragePath)
		if not success:
			self.logger.error("Failed to bring parameter %s", parameterUrlPath)
			return False, ""
		if value is None:
			self.logger.error("Could not find parameter %s", parameterUrlPath)
			return True, None
		self.logger.debug("Parameter get response for %s ready", parameterUrlPath)
		return True, value

	def searchParameters(self, parameterKeyPrefixUrlPath, keyFilters, filters):		
		"""
		Get key/value dict of parameters by key pprefix
		:param parameterKeyPrefixUrlPath: the key prefix for the parameter as presented as url path 
		(tokens seperated by "/", with "/" as a prefix)
		:type parameterKeyPrefixUrlPath: str
		:param keyFilters: a dictionary of filter-name/filter-functions - 
		each function gets a single variable - a parameter key - 
		and return "True" if the variable should be included in the response.
		:type keyFilters: Dict[str, function[str]]
		:param filters: a dictionary of filter-name/filter-functions - 
		each function gets 2 variable - a parameter key and val - 
		and return "True" if the variable should be included in the response.
		:type filters: Dict[str, function[str,str]]		
		:return: 'True' for successful retrival and the retrieved values dict (parameter-key/value)
		:rtype: Tuple[boolean, Dict[str, str]]		
		"""
		self.logger.debug("Get parameters under path %s", parameterKeyPrefixUrlPath)
		parameterStoragePathPrefix = self._get_parameter_storage_path(parameterKeyPrefixUrlPath)
		success, items = self._read_parameters_by_storage_path(parameterStoragePathPrefix, keyFilters=keyFilters)
		if not success:
			self.logger.error("Failed to bring parameters by prefix %s", parameterStoragePathPrefix)
			return False, ""

		filtered = {}
		for key, val in items.items():#items() - supporting python 2&3
			skip = False
			for filterName, filterfunc in filters.items():#items() - supporting python 2&3
				if not filterfunc(key, val):
					self.logger.debug("Parameter %s dropped, not matching filter %s", key, filterName)
					skip = True
					break
			if skip:
				continue
			filtered[key] = val

		return True, filtered

	def setParameter(self, parameterUrlPath, value):
		"""
		Set value for the specified parameter. Part of Adapter required API.
		:param parameterUrlPath: the key of the parameter as presented as url path 
		(tokens seperated by "/", with "/" as a prefix)
		:type parameterUrlPath: str
		:param value: the value to be kept
		:type value: str
		:return: 'True' for successful settings
		:rtype: Boolean
		"""
		self.logger.debug("Set parameter %s", parameterUrlPath)
		parameterStoragePath = self._get_parameter_storage_path(parameterUrlPath)
		success = self._write_parameter_by_storage_path(parameterStoragePath, value)
		if not success:
			self.logger.error("Failed to set parameter %s", parameterUrlPath)
			return False
		return True

	def deleteParameter(self, parameterUrlPath):
		"""
		Delete the specified parameter. Part of Adapter required API.
		:param parameterUrlPath: the key of the parameter as presented as url path 
		(tokens seperated by "/", with "/" as a prefix)
		:type parameterUrlPath: str
		:return: 'True' for successful deletion
		:rtype: Boolean
		"""
		self.logger.debug("Delete parameter %s", parameterUrlPath)
		parameterStoragePath = self._get_parameter_storage_path(parameterUrlPath)
		success = self._remove_parameter_by_storage_path(parameterStoragePath)
		if not success:
			self.logger.error("Failed to delete parameter %s", parameterUrlPath)
			return False
		return True

	def _init_cfg(self, fullConfig):
		"""
		Method to be implemented at derived classes
		Initialize the class basic parameters. Part of Adapter required API.
		:param fullConfig: configuration to operate upon.
		:type fullConfig: configparser.ConfigParser class
		:return: 'True' for successful initialization
		:rtype: Boolean
		"""
		raise NotImplementedError()

	def _get_parameter_storage_path(self, parameterUrlPath):
		"""
		Method to be implemented at derived classes
		Conversion function - taking a key's path and translate to a file path on the file system
		:param parameterUrlPath: the "url-path" like key of the variable
		:type parameterUrlPath: str
		:return: file path of where the value is be kept
		:rtype: str
		"""
		raise NotImplementedError()

	def _get_parameter_storage_path_from_url_path(self, parameterStoragePath):
		"""
		Method to be implemented at derived classes
		Conversion function - taking file path on the file system and translate to key's path
		:param parameterStoragePath: the file name holding a value
		:type parameterUrlPath: str
		:return: the matching variable url-path like key
		:rtype: str
		"""
		raise NotImplementedError()

	def _init_ping(self):
		"""
		Method to be implemented at derived classes
		Initialize the class ability to answer for ping. Part of Adapter required API.
		:return: 'True' for successful initialization
		:rtype: Boolean
		"""
		raise NotImplementedError()

	def _ping(self):
		"""
		Method to be implemented at derived classes
		get value for the ping request. Part of Adapter required API.
		:return: A tuple - 'True' for successful retrival and the retrieved value
		:rtype: Tuple[Boolean, str]
		"""
		raise NotImplementedError()

	def _read_parameter_by_storage_path(self, parameterStoragePath):
		"""
		Method to be implemented at derived classes
		Reading the value from the provided file name.
		:param parameterStoragePath: the file name 
		:type parameterUrlPath: str
		:return: 'True' for successful retrivaland the retrieved value
		:rtype: Tuple[boolean, str]
		"""
		raise NotImplementedError()

	def _read_parameters_by_storage_path(self, parameterStoragePathPrefix, keyFilters):
		"""
		Method to be implemented at derived classes
		Reading the values of the parameters the provided directory.
		:param parameterStoragePathPrefix: the directory to look into
		:type parameterStoragePathPrefix: str
		:param keyFilters: filter-name/filter-func dict, holding functions that get a key as 
		input and retunn "true" if key should be included in the result
		:type keyFilters: Dict[str,function[str]]
		:return: 'True' for successful retrival and a dict for key-name/value 
		:rtype: Tuple[boolean, Dict[str, str]]
		"""
		raise NotImplementedError()


	def _write_parameter_by_storage_path(self, parameterStoragePath, value):
		"""
		Method to be implemented at derived classes
		Writing the value to the provided file name.
		:param parameterStoragePath: the file name 
		:type parameterUrlPath: str
		:param value: value to be writen
		:type parameterUrlPath: str
		:return: 'True' for successful writing
		:rtype: Boolean		
	   """
		raise NotImplementedError()

	def _remove_parameter_by_storage_path(self, parameterStoragePath):
		"""
		Method to be implemented at derived classes
		Deleting the the provided file.
		:param parameterStoragePath: the file name 
		:type parameterUrlPath: str
		:return: 'True' for successful deletion
		:rtype: Boolean		
		"""
		raise NotImplementedError()