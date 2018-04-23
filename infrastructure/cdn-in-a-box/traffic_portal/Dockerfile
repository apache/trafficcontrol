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

FROM node:latest AS trafficportal-build

# Getting the whole source.. should eventually find a way to not include everything (since we don't need it all)
COPY traffic_portal/ /

# Installs dependencies
RUN apt-get update -y && apt-get install -y libffi-dev ruby-dev rubygems vim

# Installs rubygem dependencies
RUN gem update --system
RUN gem install --no-rdoc --no-ri compass
RUN gem install --no-rdoc --no-ri sass -v 3.4.22

# Global node dependencies (*shudders involuntarily*)
RUN npm install -g bower grunt-cli

# Sure, why not run javascript as root?
RUN echo '{ "allow_root": true }' > /root/.bowerrc

# The build step, using yum.js (even though we already have apt-get.js) and make.js
WORKDIR /traffic_portal
RUN npm install
RUN bower install
RUN grunt dist


# Now the actual deployment container
FROM node:latest as trafficportal

# This controls the place from which Traffic Portal is deployed
ENV TP_HOME /opt/traffic_portal

# The source location in the previous build stage
ENV SRC /traffic_portal

# This really ought to be taken care of by docker-compose, but idk; the other pieces are doing this
ENV TO_PORT 443
ENV TO_HOST trafficops-go

# Setting up directory structure
RUN mkdir -p /etc/traffic_portal
RUN mkdir -p /etc/logrotate.d
RUN mkdir -p /etc/init.d
RUN mkdir -p /var/log/traffic_portal
RUN mkdir -p ${TP_HOME}/server

# Copy the final product into its proper place(s)
WORKDIR ${TP_HOME}
COPY --from=trafficportal-build ${SRC}/server.js .
# COPY --from=trafficportal-build ${SRC}/conf /etc/traffic_portal/
COPY --from=trafficportal-build ${SRC}/build/etc/init.d/traffic_portal /etc/init.d/
COPY --from=trafficportal-build ${SRC}/build/etc/logrotate.d/traffic_portal /etc/logrotate.d/
COPY --from=trafficportal-build ${SRC}/build/etc/logrotate.d/traffic_portal-access /etc/logrotate.d/
COPY --from=trafficportal-build ${SRC}/app/dist/ .
RUN mv dist/* .
RUN rmdir dist

# Now to build the release info
RUN echo -e "{\n\"Version\":\"2.1\",\n\"Build Date\":\"$(dat +'%Y-%m-%d %H:%M:%S')\"\n}" > public/traffic_portal_release.json

# Make some things executable
RUN chmod +x ${TP_HOME}/node_modules/forever/bin/forever
RUN chmod +x /etc/init.d/traffic_portal

EXPOSE 61443

COPY infrastructure/cdn-in-a-box/traffic_portal/run.sh /

CMD /run.sh
