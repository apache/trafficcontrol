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

var DeliveryServiceSslKeysService = function($http, messageModel, ENV) {
    this.successMessage = 'SSL Keys generated and updated for ';
    this.acmeSuccessMessage = 'ACME call has been made successfully. This may take a few minutes. Please watch for a notification in the Change Log. Delivery Service = ';

    this.generateSslKeys = function(deliveryService, sslKeys, generateSslKeyForm) {
        return this.generateSslKeysBase(deliveryService, sslKeys, generateSslKeyForm, 'deliveryservices/sslkeys/generate', this.successMessage);
    };

    this.generateSslKeysWithAcme = function(deliveryService, sslKeys, generateSslKeyForm) {
        return this.generateSslKeysBase(deliveryService, sslKeys, generateSslKeyForm, 'deliveryservices/sslkeys/generate/acme', null);
    };

    this.generateSslKeysBase = function(deliveryService, sslKeys, generateSslKeyForm, endpoint, message) {
        if (sslKeys.hasOwnProperty('version')){
            generateSslKeyForm.version = parseInt(sslKeys.version, 10) + 1;
        } else {
            generateSslKeyForm.version = 1;
        }

        generateSslKeyForm.cdn = deliveryService.cdnName;
        generateSslKeyForm.deliveryservice = deliveryService.xmlId;
        generateSslKeyForm.key = deliveryService.xmlId;
        generateSslKeyForm.authType = sslKeys.authType;

        return $http.post(ENV.api.unstable + endpoint, generateSslKeyForm).then(
            function(result) {
                if (message === null) {
                    messageModel.setMessages(result.data.alerts, true);
                } else {
                    messageModel.setMessages([{level: 'success', text: message + deliveryService.xmlId}], true);
                }
                return result.data.response;
            },
            function(err) {
                if (err.data && err.data.alerts) {
                    messageModel.setMessages(err.data.alerts, false);
                }
                throw err;
            }
        );
    };

    this.renewCert = function(deliveryService) {
        return $http.post(ENV.api.unstable + "deliveryservices/xmlId/" + deliveryService.xmlId + "/sslkeys/renew").then(
            function(result) {
                messageModel.setMessages(result.data.alerts, false);
                return result.data.response;
            },
            function(err) {
                if (err.data && err.data.alerts) {
                    messageModel.setMessages(err.data.alerts, false);
                }
                throw err;
            }
        );
    };

    this.addSslKeys = function(sslKeys, deliveryService) {

        sslKeys.key = deliveryService.xmlId;
        if (sslKeys.hasOwnProperty('version')){
            sslKeys.version = parseInt(sslKeys.version, 10) + 1;
        } else {
            sslKeys.version = 1;
        }

        sslKeys.cdn = deliveryService.cdnName;
        sslKeys.deliveryservice = deliveryService.xmlId;

        return $http.post(ENV.api.unstable + "deliveryservices/sslkeys/add", sslKeys).then(
            function(result) {
                messageModel.setMessages(result.data.alerts, false);
                return result.data.response;
            },
            function(err) {
                if (err.data && err.data.alerts) {
                    messageModel.setMessages(err.data.alerts, false);
                }
                throw err;
            }
        );
    };

    this.deleteCert = function(deliveryService) {
        return $http.delete(ENV.api.unstable + "deliveryservices/xmlId/" + deliveryService.xmlId + "/sslkeys").then(
            function(result) {
                messageModel.setMessages(result.data.alerts, false);
                return result;
            },
            function(err) {
                if (err.data && err.data.alerts) {
                    messageModel.setMessages(err.data.alerts, false);
                    throw err;
                }
            }
        );
    };

    this.getSslKeys = function(deliveryService) {
        return $http.get(ENV.api.unstable + "deliveryservices/xmlId/" + deliveryService.xmlId + "/sslkeys", {params: {decode: "true"}}).then(
            function(result) {
                return result.data.response;
            },
            function(err) {
                if (err.data && err.data.alerts) {
                    messageModel.setMessages(err.data.alerts, false);
                }
                throw err;
            }
        );
    };

    this.getAcmeProviders = function() {
        return $http.get(ENV.api.unstable + 'acme_accounts/providers').then(
            function (result) {
                return result.data.response;
            },
            function (err) {
                throw err;
            }
        );
    };


    this.revertCert = function(deliveryService) {
        return $http.get(ENV.api.unstable + "deliveryservices/xmlId/" + deliveryService.xmlId + "/sslkeys", {params: {decode: "true"}}).then(
            function(result) {
                let prevVersion = parseInt(result.data.response.version, 10) - 1;
                $http.get(ENV.api.unstable + "deliveryservices/xmlId/" + deliveryService.xmlId + "/sslkeys", {params: {decode: "true", version: prevVersion}}).then(
                    function(result) {
                        let prevKeys = result.data.response;
                        prevKeys.cdn = deliveryService.cdnName;
                        prevKeys.deliveryservice = deliveryService.xmlId;

                        return $http.post(ENV.api.unstable + "deliveryservices/sslkeys/add", prevKeys).then(
                            function(result) {
                                messageModel.setMessages(result.data.alerts, false);
                                return result.data.response;
                            },
                            function(err) {
                                if (err.data && err.data.alerts) {
                                    messageModel.setMessages(err.data.alerts, false);
                                }
                                throw err;
                            }
                        );
                    },
                    function(err) {
                        if (err.data && err.data.alerts) {
                            messageModel.setMessages(err.data.alerts, false);
                        }
                        throw err;
                    }
                );
            },
            function(err) {
                if (err.data && err.data.alerts) {
                    messageModel.setMessages(err.data.alerts, false);
                }
                throw err;
            }
        );
    };
};

DeliveryServiceSslKeysService.$inject = ['$http', 'messageModel', 'ENV'];
module.exports = DeliveryServiceSslKeysService;
