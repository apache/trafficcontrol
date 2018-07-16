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

############################################################
# Dockerfile to build Mid-Tier Cache container images for
# Apache Traffic Control
# Based on CentOS 7.2
############################################################

FROM alpine:latest

RUN apk add --no-cache lighttpd

RUN rm /etc/lighttpd/lighttpd.conf
RUN rm -rf /var/www/localhost/
RUN mkdir -p /var/www/html/

ADD lighttpd.conf /etc/lighttpd/lighttpd.conf
ADD index.html /var/www/html/index.html
ADD tc_logo.svg /var/www/html/tc_logo.svg

EXPOSE 80

CMD lighttpd -t -f /etc/lighttpd/lighttpd.conf && lighttpd -D -f /etc/lighttpd/lighttpd.conf
