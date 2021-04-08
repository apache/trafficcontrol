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

// this is the config that is consumed by /server.js on traffic portal startup (sudo service traffic_portal start)
module.exports = {
    timeout: '120s',
    useSSL: true, // set to true if you plan to use https (self-signed or trusted certs).
    port: 8880, // set to http port
    sslPort: 8443, // set to https port
    // if useSSL is true, generate ssl certs and provide the proper locations.
    ssl: {
        key:    '../traffic_ops/traffic_ops_golang/localhost.key',
        cert:   '../traffic_ops/traffic_ops_golang/localhost.crt',
        ca:     [ '../traffic_ops/traffic_ops_golang/localhost.crt' ]
    },
    // set api 'base_url' to the traffic ops api url (all api calls made from the traffic portal will be proxied to the api base_url)
    api: {
        base_url: 'https://localhost:6443/api/'
    },
    // default static files location (this is where the traffic portal html, css and javascript was installed. rpm installs these files at /opt/traffic_portal/public
    // change this to ./app/dist/public/ if you are running locally for development
    files: {
        static: './app/dist/public/'
    },
    // default log location (this is where traffic_portal logs are written)
    // change this to ./server/log/access.log if you are running traffic portal locally for development
    log: {
        stream: './access.log'
    },
    reject_unauthorized: 0 // 0 if using self-signed certs, 1 if trusted certs
};

