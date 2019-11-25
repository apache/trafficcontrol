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

import sys
if sys.version_info >= (3, 0):
	from abc import ABC, abstractmethod	
else:	
	ABC = object
	abstractmethod = lambda f: f
	


class AdapterBase(ABC):
	"""
	Base adapter class.
	This class implements the API required for storing and retriving the content kept in the vault.
	
	Methods to be implemented at derived classes:
	:meth:`get_parameter_storage_path` given a url-key of a parameter return the storage path
	:meth:`get_parameter_storage_path_from_url_path` given a storage path of a parameter
	retrun the url-key
	:meth:`init_cfg` given a config-parser object, read the parameters required for 
	adapter's operation. Return "success" bool value
	:meth:`init` prepare the adapter - connecting to storgae / setting the conntent for 
	"ping" requests etc.. Return "success" bool value
	:meth:`ping` test the 'ping' status with the adapter. 
	Return a tuple: "success" bool & "value" kept as ping variable
	:meth:`read_parameter_by_storage_path` given a storage path retrieve the parameter value.
	Return a tuple: "success" bool & "value" kept in the parameter
	:meth:`read_parameters_by_storage_path` given a storage and and a key holding 
	filters on the key and values.
	Return "success" bool indicating a sucessful write, and a key->value dictionary 
	for the relevant parameters
	:meth:`write_parameter_by_storage_path` given a storage path and a value string,
	keep the parameter value. Return "success" bool indicating a sucessful write
	:meth:`remove_parameter_by_storage_path` given a storage path, delete the parameter
	from the DB. Return "success" bool indicating a sucessful deletion
	"""

	@abstractmethod
	def init_cfg(self, fullConfig):# -> bool:
		"""
		Method to be implemented at derived classes
		Initialize the class basic parameters. Part of Adapter required API.
		:param fullConfig: configuration to operate upon.
		:type fullConfig: configparser.ConfigParser class
		:return: 'True' for successful initialization
		:rtype: bool
		"""
		raise NotImplementedError()#...

	@abstractmethod	
	def init(self):# -> bool:
		"""
		Method to be implemented at derived classes
		Initialize the class - e.g. connection to storage & ability to answer for ping. 
		:return: 'True' for successful initialization
		:rtype: bool
		"""
		raise NotImplementedError()#...

	@abstractmethod	
	def get_parameter_storage_path(self, parameterUrlPath):# -> (bool, str):
		"""
		Method to be implemented at derived classes
		Conversion function - taking a key's path and translate to a file path on the file system
		:param parameterUrlPath: the "url-path" like key of the variable
		:type parameterUrlPath: str
		:return: "success" bool and a file path of where the value is be kept
		:rtype: Tuple[bool, str]
		"""
		raise NotImplementedError()#...

	@abstractmethod	
	def get_parameter_url_path_from_storage_path(self, parameterStoragePath):# -> (bool, str):
		"""
		Method to be implemented at derived classes
		Conversion function - taking file path on the file system and translate to key's path
		:param parameterStoragePath: the file name holding a value
		:type parameterUrlPath: str
		:return: "success" bool and the matching variable url-path like key
		:rtype: Tuple[bool, str]
		"""
		raise NotImplementedError()#...

	@abstractmethod	
	def ping(self):# -> bool:
		"""
		Method to be implemented at derived classes
		Check ping connection
		:return: 'True' for successful connection with the storage layer
		:rtype: bool
		"""
		raise NotImplementedError()#...

	@abstractmethod	
	def read_parameter_by_storage_path(self, parameterStoragePath):# -> (bool, str):
		"""
		Method to be implemented at derived classes
		Reading the value from the provided file name.
		:param parameterStoragePath: the file name 
		:type parameterUrlPath: str
		:return: 'True' for successful retrivaland the retrieved value
		:rtype: Tuple[bool, str]
		"""
		raise NotImplementedError()#...

	@abstractmethod	
	def read_parameters_by_storage_path(self, parameterStoragePathPrefix, keyFilters):# -> (bool, dict(str,str)):
		"""
		Method to be implemented at derived classes
		Reading the values of the parameters the provided directory.
		:param parameterStoragePathPrefix: the directory to look into
		:type parameterStoragePathPrefix: str
		:param keyFilters: filter-name/filter-func dict, holding functions that get a key as 
		input and retunn "true" if key should be included in the result
		:type keyFilters: Dict[str,function[str]]
		:return: 'True' for successful retrival and a dict for key-name/value 
		:rtype: Tuple[bool, Dict[str, str]]
		"""
		raise NotImplementedError()#...


	@abstractmethod	
	def write_parameter_by_storage_path(self, parameterStoragePath, value):# -> bool:
		"""
		Method to be implemented at derived classes
		Writing the value to the provided file name.
		:param parameterStoragePath: the file name 
		:type parameterUrlPath: str
		:param value: value to be writen
		:type parameterUrlPath: str
		:return: 'True' for successful writing
		:rtype: bool		
	   """
		raise NotImplementedError()#...

	@abstractmethod	
	def remove_parameter_by_storage_path(self, parameterStoragePath):# -> bool:
		"""
		Method to be implemented at derived classes
		Deleting the the provided file.
		:param parameterStoragePath: the file name 
		:type parameterUrlPath: str
		:return: 'True' for successful deletion
		:rtype: bool		
		"""
		raise NotImplementedError()#...

