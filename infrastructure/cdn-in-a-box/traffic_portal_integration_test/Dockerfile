# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

FROM centos:7 as os-dependencies

# Installs the Google Chrome yum repo
COPY infrastructure/cdn-in-a-box/traffic_portal_integration_test/etc etc

RUN yum install -y \
        bind-utils \
        # jq is in EPEL
        epel-release \
        git \
        google-chrome-stable \
        java-1.8.0-openjdk \
        net-tools && \
    yum -y install jq && \
    yum -y clean all

FROM os-dependencies AS node-dependencies
# Download and install node
RUN curl -Lo node.tar.xz https://nodejs.org/dist/v12.16.3/node-v12.16.3-linux-x64.tar.xz && \
    tar --strip-components=1 -C /usr/local -xf node.tar.xz && \
    rm node.tar.xz && \
    npm install -g protractor@^7.0.0

FROM node-dependencies

COPY    traffic_portal/ /lang/traffic_portal/
WORKDIR /lang/traffic_portal/test/end_to_end
COPY    infrastructure/cdn-in-a-box/traffic_portal_integration_test/junit-jasmine-onprepare.js .

RUN npm install --save-dev jasmine-reporters@^2.0.0

# Sets the exports.config.onPrepare handler, which cannot be stored in JSON
RUN mkdir /portaltestresults && \
    cat junit-jasmine-onprepare.js >> conf.js && \
    rm junit-jasmine-onprepare.js

RUN jq ' \
        .capabilities.chromeOptions.args = [ \
            "--disable-extensions", \
            "--disable-gpu", \
            "--headless", \
            "--ignore-certificate-errors", \
            "--no-sandbox", \
            "--whitelisted-ips=" \
        ] ' \
        conf.json > conf.json.tmp && \
    mv conf.json.tmp conf.json

RUN webdriver-manager clean && \
    webdriver-manager update

COPY infrastructure/cdn-in-a-box/traffic_ops/to-access.sh \
     infrastructure/cdn-in-a-box/traffic_portal_integration_test/run.sh \
     /usr/local/sbin/

CMD run.sh
