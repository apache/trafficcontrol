#
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
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

class Fs(object):
    """
    Fs (file system) adapter class.
    This class implements the API required for storing and retriving the content kept in the vault.
    This specific Adapter works using files upon a file-system. 
    Inputs configuration (held under "fs-adapter" section in the config file) include:
    :param db-base-os-path: The path in which the DB files are stored
    :param ping-os-path: The path of a variable mimicing the RIAK ping functionality

    Interface methods 
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
        myCfgData = dict(fullConfig.items("fs-adapter")) if fullConfig.has_section("fs-adapter") else {}
        self.basePath = myCfgData.get("db-base-os-path")
        if self.basePath is None:
            self.logger.error("Missing %s/%s configuration", "fs-adapter", "db-base-os-path")
            return False 
        self.pingOsPath = myCfgData.get("ping-os-path")
        if self.pingOsPath is None:
            self.logger.error("Missing %s/%s configuration", "fs-adapter", "ping-os-path")
            return False 
        return True

    def init_ping(self):
        """
        Initialize the class ability to answer for ping. Part of Adapter required API.
        :return: 'True' for successful initialization
        :rtype: Boolean
        """
        value = "OK"
        success = self._write_parameter_os_path(self.pingOsPath, value)
        if not success:
            self.logger.error("Failed to set parameter %s", self.pingOsPath)
            return False
        return True

    def ping(self):
        """
        get value for the ping request. Part of Adapter required API.
        :return: A tuple - 'True' for successful retrival and the retrieved value
        :rtype: Tuple[Boolean, str]
        """
        success, value = self._read_parameter_os_path(self.pingOsPath)
        if not success or value is None:
            self.logger.error("no ping response")
            return (False, "")
        self.logger.debug("ping response: %s", value)
        return (True, value)

    def getParameter(self, parameterKeyUrlPath):
        """
        Get value for the specified parameter. Part of Adapter required API.
        :param parameterKeyUrlPath: the key of the parameter as presented as url path 
        (tokens seperated by "/", with "/" as a prefix)
        :type parameterKeyUrlPath: str
        :return: A tuple - 'True' for successful retrival and the retrieved value
        :rtype: Tuple[Boolean, str]
        """
        parameterOsPath = self._get_parameter_os_path(parameterKeyUrlPath)        
        success, value = self._read_parameter_os_path(parameterOsPath)
        if not success:
            self.logger.error("Failed to bring parameter %s", parameterKeyUrlPath)
            return False, ""
        if value is None:
            self.logger.error("Could not find parameter %s", parameterKeyUrlPath)
            return True, None
        self.logger.debug("Parameter get response for %s ready", parameterKeyUrlPath)
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
        parameterOsPathPrefix = self._get_parameter_os_path(parameterKeyPrefixUrlPath)
        success, items = self._read_parameter_os_paths(parameterOsPathPrefix, keyFilters=keyFilters)
        if not success:
            self.logger.error("Failed to bring parameters by prefix %s", parameterOsPathPrefix)
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

    def setParameter(self, parameterKeyUrlPath, value):
        """
        Set value for the specified parameter. Part of Adapter required API.
        :param parameterKeyUrlPath: the key of the parameter as presented as url path 
        (tokens seperated by "/", with "/" as a prefix)
        :type parameterKeyUrlPath: str
        :param value: the value to be kept
        :type value: str
        :return: 'True' for successful settings
        :rtype: Boolean
        """
        self.logger.debug("Set parameter %s", parameterKeyUrlPath)
        parameterOsPath = self._get_parameter_os_path(parameterKeyUrlPath)
        success = self._write_parameter_os_path(parameterOsPath, value)
        if not success:
            self.logger.error("Failed to set parameter %s", parameterKeyUrlPath)
            return False
        return True

    def deleteParameter(self, parameterKeyUrlPath):
        """
        Delete the specified parameter. Part of Adapter required API.
        :param parameterKeyUrlPath: the key of the parameter as presented as url path 
        (tokens seperated by "/", with "/" as a prefix)
        :type parameterKeyUrlPath: str
        :return: 'True' for successful deletion
        :rtype: Boolean
        """
        self.logger.debug("Delete parameter %s", parameterKeyUrlPath)
        parameterOsPath = self._get_parameter_os_path(parameterKeyUrlPath)
        success = self._remove_parameter_os_path(parameterOsPath)
        if not success:
            self.logger.error("Failed to delete parameter %s", parameterKeyUrlPath)
            return False
        return True

    def _get_parameter_os_path(self, parameterKeyUrlPath):
        """
        Conversion function - taking a key's path and translate to a file path on the file system
        :param parameterKeyUrlPath: the "url-path" like key of the variable
        :type parameterKeyUrlPath: str
        :return: file path of where the value is be kept
        :rtype: str
        """
        return os.path.join(self.basePath, parameterKeyUrlPath.lstrip("/").replace("/", os.path.sep))

    def _get_parameter_os_path_key_url_path(self, parameterOsPath):
        """
        Conversion function - taking file path on the file system and translate to key's path
        :param parameterOsPath: the file name holding a value
        :type parameterKeyUrlPath: str
        :return: the matching variable url-path like key
        :rtype: str
        """
        return "/"+os.path.relpath(parameterOsPath, self.basePath).replace(os.path.sep, "/")
        

    def _read_parameter_os_path(self, parameterOsPath):
        """
        Reading the value from the provided file name.
        :param parameterOsPath: the file name 
        :type parameterKeyUrlPath: str
        :return: 'True' for successful retrivaland the retrieved value
        :rtype: Tuple[boolean, str]
        """
        self.logger.debug("Get parameter by os path: %s", parameterOsPath)
        try:
            with open(parameterOsPath) as fd:
                value = fd.read()
        except OSError as e:
            self.logger.exception("%s parameter os path not found.", parameterOsPath)
            return False, None

        self.logger.debug("Get parameter by os path %s succeed", parameterOsPath)
        return True, value

    def _read_parameter_os_paths(self, parameterOsPathPrefix, keyFilters):
        """
        Reading the values of the parameters the provided directory.
        :param parameterOsPathPrefix: the directory to look into
        :type parameterOsPathPrefix: str
        :param keyFilters: filter-name/filter-func dict, holding functions that get a key as 
        input and retunn "true" if key should be included in the result
        :type keyFilters: Dict[str,function[str]]
        :return: 'True' for successful retrival and a dict for key-name/value 
        :rtype: Tuple[boolean, Dict[str, str]]
        """
        self.logger.debug("Get parameters under os path %s", parameterOsPathPrefix)
        parameters = {}

        fileNames = []
        for (dirpath, _, filenames) in os.walk(parameterOsPathPrefix):
            fileNames += [os.path.join(dirpath, file) for file in filenames]

        for fileName in fileNames:
            filteredOut = False
            for filterName, filterfunc in keyFilters.items():#items() - supporting python 2&3
                if not filterfunc(fileName):
                    self.logger.debug("Parameter os path %s dropped, not matching filter %s", self._get_parameter_os_path_key_url_path(fileName), filterName)
                    filteredOut = True
                    break
            if filteredOut:
                continue

            parameterKeyUrlPath = self._get_parameter_os_path_key_url_path(fileName)
            success, value = self._read_parameter_os_path(fileName)
            if not success:
                self.logger.error("%s parameter os path not found.", fileName)
                return False, None
                
            parameters[parameterKeyUrlPath] = value

        return True, parameters


    def _write_parameter_os_path(self, parameterOsPath, value):
        """
        Writing the value to the provided file name.
        :param parameterOsPath: the file name 
        :type parameterKeyUrlPath: str
        :param value: value to be writen
        :type parameterKeyUrlPath: str
        :return: 'True' for successful writing
        :rtype: Boolean        
       """
        self.logger.debug("Set parameter by os path %s", parameterOsPath)
        try:
            dirname = os.path.dirname(parameterOsPath)
            if dirname and not os.path.exists(dirname):
                os.makedirs(dirname)
            with open(parameterOsPath, "w") as fd:
                fd.write(value)            
        except OSError as e:
            self.logger.exception("could not post parameter os path %s", parameterOsPath)
            return False
        self.logger.debug("Set parameter os path %s done", parameterOsPath)
        return True

    def _remove_parameter_os_path(self, parameterOsPath):
        """
        Deleting the the provided file.
        :param parameterOsPath: the file name 
        :type parameterKeyUrlPath: str
        :return: 'True' for successful deletion
        :rtype: Boolean        
       """
        self.logger.debug("Delete parameter os path %s", parameterOsPath)
        try:
            os.remove(_parameterOsPath)
        except OSError as e:
            self.logger.exception("could not delete parameter os path %s", parameterOsPath)
            return False
        self.logger.debug("Delete parameter os path %s done", parameterOsPath)
        return True

