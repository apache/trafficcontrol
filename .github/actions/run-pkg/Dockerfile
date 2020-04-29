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

FROM docker:stable

VOLUME /var/run/docker.sock:/var/run/docker.sock
RUN apk add curl bash py-pip python-dev libffi-dev openssl-dev gcc libc-dev make docker-compose

# RUN apt-get update && \
# 	apt-get install -y \
# 		apt-transport-https \
# 		ca-certificates \
# 		curl \
# 		gnupg-agent \
# 		software-properties-common \
# 		lsb-release &&\
# 	curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - &&\
# 	apt-add-repository \
# 		"deb [arch=amd64] https://download.docker.com/linux/ubuntu \
# 		$(lsb_release -cs) \
# 		stable" &&\
# 	apt-get update &&\
# 	apt-get install docker-ce docker-ce-cli containerd.io

COPY entrypoint.sh /
RUN chmod a+x /entrypoint.sh

ENTRYPOINT /entrypoint.sh
