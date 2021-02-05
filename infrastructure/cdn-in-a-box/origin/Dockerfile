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
############################################################

FROM alpine:3.13

RUN apk add --no-cache \
    net-tools \
    lighttpd \
    bash \
    curl \
    bind-tools \
    # contains envsubst for to-enroll
    gettext \
    jq

RUN rm /sbin/route

RUN rm /etc/lighttpd/lighttpd.conf
RUN rm -rf /var/www/localhost/

ADD origin/content /var/www/html/

ADD origin/lighttpd.conf /etc/lighttpd/lighttpd.conf
ADD origin/run.sh \
    traffic_ops/to-access.sh \
    enroller/server_template.json \
    /

COPY dns/set-dns.sh \
     dns/insert-self-into-dns.sh \
     /usr/local/sbin/

EXPOSE 80

CMD /run.sh
