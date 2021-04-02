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

// API Utility
import { Agent } from "https";

import axios from 'axios';
import randomIpv6 from "random-ipv6";

import { config, randomize } from '../config';

export class API {

    private cookie = "";
    private readonly config = config;

    constructor() {
        axios.defaults.headers.common['Accept'] = 'application/json'
        axios.defaults.headers.common['Authorization'] = 'No-Auth'
        axios.defaults.headers.common['Content-Type'] = 'application/json'
        axios.defaults.httpsAgent = new Agent({ rejectUnauthorized: false })
    }

    Login = async function () {
        try {
            const response = await axios({
                method: 'post',
                url: this.config.params.apiUrl + '/user/login',
                data: {
                    u: this.config.params.login.username,
                    p: this.config.params.login.password
                }
            });
            this.cookie = await response.headers["set-cookie"][0];
            return response
        } catch (error) {
            return error;
        }
    }

    SendRequest = async function (route, method, data) {
        try {
            let response
            this.Randomize(data)

            if(data.hasOwnProperty('getRequest')){
                let response = await this.GetId(data);
                if (response != null) {
                    throw new Error('Failed to get id:\nResponse Status: ' + response.statusText + '\nResponse Data: ' + response.data)
                }
            }

            switch (method) {
                case "post":
                    response = await axios({
                        method: method,
                        url: this.config.params.apiUrl + route,
                        headers: { Cookie: this.cookie},
                        data: data
                    });
                    break;
                case "get":
                    response = await axios({
                        method: method,
                        url: this.config.params.apiUrl + route,
                        headers: { Cookie: this.cookie},
                    });
                    break;
                case "delete":
                    if ((data.route).includes('?name')){
                        data.route = data.route + randomize
                    }
                    if ((data.route).includes('?id')){
                        data.route = data.route + data.id;
                    }
                    if((data.route).includes('/service_categories/')){
                        data.route = data.route + randomize
                    }
                    response = await axios({
                        method: method,
                        url: this.config.params.apiUrl + data.route,
                        headers: { Cookie: this.cookie},
                    });
                    break;
            }
            if (response.status == 200 || response.status == 201) {
                return null
            } else {
                console.log("Reponse Data: " , response.data);
                console.log("Response: " , response);
                throw new Error('Request Failed:\nResponse Status: ' + response.statusText + '\nResponse Data: ' + response.data);
            }
        } catch (error) {
            return error;
        }
    }

    GetId = async function (data) {
        for(var i = 0; i < data.getRequest.length; i++) {
            var query = '?' + data.getRequest[i].queryKey  + '=' + data.getRequest[i].queryValue + randomize;
            try {
                const response = await axios({
                    method: 'get',
                    url: this.config.params.apiUrl + data.getRequest[i].route + query,
                    headers: { Cookie: this.cookie},
               });

               if (response.status == 200) {
                    if(data.getRequest[i].hasOwnProperty('isArray')){
                        data[data.getRequest[i].replace] = [await response.data.response[0].id];
                    } else if (data.getRequest[i].replace == "route") {
                        data[data.getRequest[i].replace] = data.route + response.data.response[0].id;
                    } else {
                        data[data.getRequest[i].replace] = await response.data.response[0].id;
                    }
                } else {
                    return response
                }
            } catch (error) {
                return error;
            }
        }
        return null
    }

   Randomize = function(data) {
        if(data.hasOwnProperty('email')) {
            data['email'] = data.fullName + randomize + data.email;
        }
        if(data.hasOwnProperty('fullName')) {
            data['fullName'] = data.fullName + randomize;
        }
        if(data.hasOwnProperty('hostName')) {
            data['hostName'] = data.hostName + randomize;
        }
        if(data.hasOwnProperty('ipAddress')) {
            data['ipAddress'] = (Math.floor(Math.random() * 255) + 1)+"."+(Math.floor(Math.random() * 255))+"."+(Math.floor(Math.random() * 255))+"."+(Math.floor(Math.random() * 255));
        }
        if(data.hasOwnProperty('name')) {
            data['name'] = data.name + randomize;
        }
        if(data.hasOwnProperty('requiredCapability')) {
            data['requiredCapability'] = data.requiredCapability + randomize;
        }
        if(data.hasOwnProperty('serverCapability')) {
            data['serverCapability'] = data.serverCapability + randomize;
        }
        if(data.hasOwnProperty('username')) {
            data['username'] = data.username + randomize;
        }
        if(data.hasOwnProperty('xmlId')) {
            data['xmlId'] = data.xmlId + randomize;
        }
        if(data.hasOwnProperty('shortName')) {
            data['shortName'] = data.shortName + randomize;
        }
        if(data.hasOwnProperty('divisionName')) {
            data['divisionName'] = data.divisionName + randomize;
        }
        if(data.hasOwnProperty('domainName')) {
            data['domainName'] = data.domainName + randomize;
        }
        if(data.hasOwnProperty('nodes')){
           for(var i in  data['nodes']){
               data['nodes'][i].cachegroup = data['nodes'][i].cachegroup + randomize;
           }
        }
        if(data.hasOwnProperty('interfaces')){
            let ipv6 = randomIpv6();
            for(var i in data['interfaces']){
                for(var j in data['interfaces'][i].ipAddresses){
                   data['interfaces'][i].ipAddresses[j].address = ipv6.toString();
                }
            }
        }
    }

    UseAPI = async function(data) {
        try {
            let response = await this.Login();
            if (response.status == 200) {
                for(var i = 0; i < data.Prerequisites.length; i++){
                    for(var j = 0; j < data.Prerequisites[i].Data.length; j++){
                        let output = await this.SendRequest(data.Prerequisites[i].Route, data.Prerequisites[i].Method, data.Prerequisites[i].Data[j]);
                        if (output != null) {
                            console.error(`UseAPI failed on Action ${data.Prerequisites[i].Action} with index ${i}, and Data index ${j}`);
                            throw new Error(output);
                        }
                    }
                }
                return null
            } else if (response.status == undefined) {
                throw new Error(`Error requesting ${this.config.params.apiUrl}: ${response}`);
            } else {
                throw new Error('Login failed:\nResponse Status: ' + response.statusText + '\nResponse Data: ' + response.data)
            }
        } catch (error) {
            return error;
        }
    }
}
