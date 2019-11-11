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
    def __init__ (self, logger, basePath, pingRelPath):
        self.logger = logger
        self.basePath = basePath
        self.pingRelPath = pingRelPath

    def init_ping(self):
        pingPath = os.path.join(self.basePath, self.pingRelPath)
        value = "OK"
        success = self._set_param(pingPath, value)
        if not success:
            self.logger.error("Failed to set variable %s", pingPath)
            return False
        return True

    def ping(self):
        pingPath = os.path.join(self.basePath, self.pingRelPath)
        success, value = self._get_param(pingPath)
        if not success or value is None:
            self.logger.error("no ping response")
            return (False, "")
        self.logger.debug("ping response: %s", value)
        return (True, value)

    def getParameter(self, parameterPath):
        parameterPath = parameterPath.lstrip("/")
        _parameterPath = os.path.join(self.basePath, parameterPath)
        success, value = self._get_param(_parameterPath)
        if not success:
            self.logger.error("Failed to bring variable %s", parameterPath)
            return False, ""
        if value is None:
            self.logger.error("Could not find variable %s", parameterPath)
            return True, None
        self.logger.debug("Parameter get response for %s ready", parameterPath)
        #self.logger.debug("Parameter get response: %s", value)
        return True, value

    def searchParameters(self, parameterPathPrefix, keyFilters = {}, filters = {}):
        parameterPathPrefix = parameterPathPrefix.lstrip("/")
        _parametersPathPrefix = os.path.join(self.basePath, parameterPathPrefix)
        success, items = self._get_params(_parametersPathPrefix, prefixToRemove=self.basePath, keyFilters=keyFilters)
        if not success:
            self.logger.error("Failed to bring variable by prefix %s", parameterPathPrefix)
            return False, ""

        filtered = {}
        for key, val in items.iteritems():
            skip = False
            for filterName, filter in filters.iteritems():
                if not filter(key, val):
                    self.logger.debug("Parameter %s dropped, not matching filter %s", key, filterName)
                    skip = True
                    break
            if skip:
                continue
            filtered[key] = val

        return True, filtered


    def setParameter(self, parameterPath, value):
        parameterPath = parameterPath.lstrip("/")
        _parameterPath = os.path.join(self.basePath, parameterPath)
        success = self._set_param(_parameterPath, value)
        if not success:
            self.logger.error("Failed to set variable %s", parameterPath)
            return False
        return True

    def deleteParameter(self, parameterPath):
        parameterPath = parameterPath.lstrip("/")
        _parameterPath = os.path.join(self.basePath, parameterPath)
        success = self._delete_param(_parameterPath)
        if not success:
            self.logger.error("Failed to delete variable %s", parameterPath)
            return False
        return True

    def _get_param(self, paramName):
        self.logger.debug("Get parameter %s", paramName)
        try:
            with open(paramName) as fd:
                value = fd.read()
        except:
            self.logger.exception("%s parameter not found.", paramName)
            return False, None

        self.logger.debug("Get parameter %s succeed", paramName)
        #self.logger.debug("Get parameter %s=%s", paramName, value)
        return True, value

    def _get_params(self, paramNamePrefix, prefixToRemove = "", keyFilters={}):
        self.logger.debug("Get parameters under path %s", paramNamePrefix)
        prefixToRemoveLen = len(prefixToRemove)
        parameters = {}

        fileNames = []
        for (dirpath, dirnames, filenames) in os.walk(paramNamePrefix):
            fileNames += [os.path.join(dirpath, file) for file in filenames]

        for fileName in fileNames:
            filteredOut = False
            for filterName, filter in keyFilters.iteritems():
                if not filter(fileName):
                    self.logger.debug("Parameter %s dropped, not matching filter %s", fileName, filterName)
                    filteredOut = True
                    break
            if filteredOut:
                continue

            paramName = fileName[prefixToRemoveLen:]
            success, value = self._get_param(fileName)
            if not success:
                self.logger.error("%s parameter not found.", paramName)
                return False, None
                
            parameters[paramName] = value

        return True, parameters


    def _set_param(self, paramName, value):
        self.logger.debug("Set parameter %s", paramName)
        #self.logger.debug("Set parameter %s=%s", paramName, value)
        try:
            dirname = os.path.dirname(paramName)
            if not os.path.exists(dirname):
                os.makedirs(dirname)
            with open(paramName, "w") as fd:
                fd.write(value)            
        except:
            self.logger.exception("cloud not post parameter %s", paramName)
            return False
        self.logger.debug("Set parameter %s done", paramName)
        return True

    def _delete_param(self, paramName):
        self.logger.debug("Delete parameter %s", paramName)
        try:
            os.remove(paramName)
        except:
            self.logger.exception("cloud not delete parameter %s", paramName)
            return False
        self.logger.debug("Delete parameter %s done", paramName)
        return True

