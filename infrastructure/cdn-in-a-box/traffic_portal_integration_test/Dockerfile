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

FROM centos:7

COPY ./infrastructure/cdn-in-a-box/traffic_ops/to-access.sh /usr/local/sbin/
COPY ./traffic_portal/ /lang/traffic_portal/
# Installs the Google Chrome yum repo
COPY infrastructure/cdn-in-a-box/traffic_portal_integration_test/etc etc

WORKDIR /lang/traffic_portal/test/end_to_end
RUN yum install -y \
        bind-utils \
        git \
        google-chrome-stable \
        java-1.8.0-openjdk \
        net-tools && \
    yum -y clean all

# Download and install node
RUN curl -Lo node.tar.xz https://nodejs.org/dist/v12.16.3/node-v12.16.3-linux-x64.tar.xz && \
    tar --strip-components=1 -C /usr/local -xf node.tar.xz && \
    rm node.tar.xz && \
    npm install -g protractor@^7.0.0

RUN npm install --save-dev jasmine-reporters@^2.0.0

RUN sed -i "s/baseUrl: 'https:\/\/localhost:4443/baseUrl: 'https:\/\/trafficportal.infra.ciab.test:443/" /lang/traffic_portal/test/end_to_end/conf.js

RUN export junit_conf_jasmine="\n\
	onPrepare: function() {\n\
		var jasmineReporters = require('jasmine-reporters');\n\
		jasmine.getEnv().addReporter(\n\
			new jasmineReporters.JUnitXmlReporter({savePath:'/portaltestresults', filePrefix: 'portaltestresults', consolidateAll:true})\n\
		);\n\
	},\n\
" && perl -p -i -e 'BEGIN{ ($e=$ENV{junit_conf_jasmine}) =~ s/\\n/\n/g } if (/^\s*suites\b/) {print $e}' conf.js

RUN export junit_conf_chrome="\
		chromeOptions: {\n\
			args: ['--headless', '--disable-gpu', '--no-sandbox', '--ignore-certificate-errors', '--disable-extensions', '--whitelisted-ips=\"\"'] \n\
		},\n\
" && perl -p -i -e 'BEGIN{ ($e=$ENV{junit_conf_chrome}) =~ s/\\n/\n/g } if (/\bbrowserName\b/../}/) {print $e; $e=""}' conf.js

RUN webdriver-manager clean
RUN webdriver-manager update

COPY ./infrastructure/cdn-in-a-box/traffic_portal_integration_test/run.sh /lang/traffic_portal/test/end_to_end/
CMD ./run.sh
