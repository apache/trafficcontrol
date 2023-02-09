#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
#
ARG BASE_IMAGE=rockylinux \
    RHEL_VERSION=8
FROM ${BASE_IMAGE}:${RHEL_VERSION}
ARG RHEL_VERSION=8

RUN dnf install -y openssl curl
RUN openssl genrsa -passout pass:x -out server.pass.key 2048 && \
	openssl rsa -passin pass:x -in server.pass.key -out server.key && \
	rm server.pass.key && \
	openssl req -new -key server.key -out server.csr \
		-subj "/C=US/ST=CO/L=Denver/O=Apache/OU=Traffic Control/CN=trafficops.dev.ciab.test" && \
	openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt && \
	openssl rand 32 | base64 > /aes.key && \
    chmod 644 /server.key /aes.key

RUN	curl -sL https://rpm.nodesource.com/setup_16.x | bash - && \
	dnf -y install nodejs git

ENV TC="/trafficcontrol"
VOLUME $TC

EXPOSE 443

ADD traffic-portal/build/run.sh /
ADD traffic-portal/build/config.json /
RUN chmod +x /run.sh
CMD /run.sh

