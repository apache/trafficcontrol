/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

var ApplicationService = function($rootScope, $anchorScroll, $http, messageModel, propertiesModel) {

    let storageAvailable = function(type) {
        try {
            let storage = window[type],
                x = '__storage_test__';

            storage.setItem(x, x);
            storage.removeItem(x);
            return true;
        }
        catch(e) {
            return e;
        }
    };

    let init = function() {
        $http.defaults.withCredentials = true;

        // jquery DataTables default overrides
        $.extend(true, $.fn.dataTable.defaults, {
            "dom": "lfBrtip",
            "buttons": [
                { "extend": "csv", "text": "Export as CSV", "titleAttr": "Export as CSV", "className": "btn-link" }
            ],
            "colReorder": {
                realtime: false
            },
            "stateSave": true,
            "scrollX": true
        });

        if (!storageAvailable('localStorage')) {
            messageModel.setMessages([ { level: 'warning', text: 'A browser that supports local storage is required to use ' + propertiesModel.properties.name } ], false);
        }
    };
    init();

    $rootScope.$on("$viewContentLoaded", function() {
        $anchorScroll(); // scrolls window to top
    });

};

ApplicationService.$inject = ['$rootScope', '$anchorScroll', '$http', 'messageModel', 'propertiesModel'];
module.exports = ApplicationService;
