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

class Db(object):
	"""
	DB class representing the DB layer.
	
	Implemented methods 
	:meth:`ping` test the 'ping' status with the adapter. 
	Return a tuple: "success" bool & "value" kept as ping variable
	:meth:`getParameter` given a parameter key (in url-path format) retrieve the parameter value.
	Return a tuple: "success" bool & "value" kept in the parameter
	:meth:`searchParameter` given a parameters key prefix (in url-path format) and, a dict holding
	variable key filters, and a key holding filters on the values as well.
	Return "success" bool indicating a sucessful write, and a key->value dictionary 
	for the relevant parameters
	:meth:`setParameter` given a parameter key (in url-path format) and a value string,
	keep the parameter value. Return "success" bool indicating a sucessful write
	:meth:`deleteParameter` given a parameter key (in url-path format), delete the parameter
	from the DB. Return "success" bool indicating a sucessful deletion		
	"""

	def __init__ (self, logger, storage_adaper):
		"""
		The class constructor.
		:param logger: logger to send log messages to
		:type logger: a python logging logger class
		:param storage_adaper: an initalized storage adapter
		:type storage_adaper: a storage.adapter_base.AdapterBase class 
		"""
		
		self.logger = logger
		self.storage_adaper = storage_adaper

	def ping(self):
		"""
		get value for the ping request. Part of Adapter required API.
		:return: A tuple - 'True' for successful retrival and the retrieved value
		:rtype: Tuple[bool, str]
		"""
		if self.storage_adaper.ping():
			return (True, "OK")
		return False, None

	def getParameter(self, parameterUrlPath):
		"""
		Get value for the specified parameter. Part of Adapter required API.
		:param parameterUrlPath: the key of the parameter as presented as url path 
		(tokens seperated by "/", with "/" as a prefix)
		:type parameterUrlPath: str
		:return: A tuple - 'True' for successful retrival and the retrieved value
		:rtype: Tuple[bool, str]
		"""
		success_path, parameterStoragePath = self.storage_adaper.get_parameter_storage_path(parameterUrlPath)
		if not success_path:
			self.logger.error("Invalid parameter url path: %s", parameterUrlPath)
			return False, None
		success, value = self.storage_adaper.read_parameter_by_storage_path(parameterStoragePath)
		if not success:
			self.logger.error("Failed to bring parameter %s", parameterUrlPath)
			return False, None
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
		:rtype: Tuple[bool, Dict[str, str]]		
		"""
		self.logger.debug("Get parameters under path %s", parameterKeyPrefixUrlPath)
		success_path, parameterStoragePathPrefix = self.storage_adaper.get_parameter_storage_path(parameterKeyPrefixUrlPath)
		if not success_path:
			self.logger.error("Invalid parameter url path prefix: %s", parameterKeyPrefixUrlPath)
			return False, None
		success, items = self.storage_adaper.read_parameters_by_storage_path(parameterStoragePathPrefix, keyFilters=keyFilters)
		if not success:
			self.logger.error("Failed to bring parameters by prefix %s", parameterStoragePathPrefix)
			return False, None

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
		:rtype: bool
		"""
		self.logger.debug("Set parameter %s", parameterUrlPath)
		success_path, parameterStoragePath = self.storage_adaper.get_parameter_storage_path(parameterUrlPath)
		if not success_path:
			self.logger.error("Invalid parameter url path: %s", parameterUrlPath)
			return False, None
		success = self.storage_adaper.write_parameter_by_storage_path(parameterStoragePath, value)
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
		:rtype: bool
		"""
		self.logger.debug("Delete parameter %s", parameterUrlPath)
		success_path, parameterStoragePath = self.storage_adaper.get_parameter_storage_path(parameterUrlPath)
		if not success_path:
			self.logger.error("Invalid parameter url path: %s", parameterUrlPath)
			return False, None
		success = self.storage_adaper.remove_parameter_by_storage_path(parameterStoragePath)
		if not success:
			self.logger.error("Failed to delete parameter %s", parameterUrlPath)
			return False
		return True

