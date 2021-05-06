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
        bash \
        bind-tools \
        curl \
        # gettext contains envsubst for to-enroll
        gettext \
        jq \
        lighttpd \
        net-tools && \
    rm -rf  /sbin/route \
            /etc/lighttpd/lighttpd.conf \
            /var/www/localhost/

COPY traffic_router/core/src/test/resources/czmap.json \
     traffic_router/core/src/test/resources/geo/GeoLite2-City.mmdb.gz \
    /var/www/html/

COPY infrastructure/cdn-in-a-box/static/lighttpd.conf /etc/lighttpd/
COPY infrastructure/cdn-in-a-box/static/run.sh \
     infrastructure/cdn-in-a-box/traffic_ops/to-access.sh \
     infrastructure/cdn-in-a-box/enroller/server_template.json \
     /

COPY infrastructure/cdn-in-a-box/dns/set-dns.sh \
     infrastructure/cdn-in-a-box/dns/insert-self-into-dns.sh \
     /usr/local/sbin/

EXPOSE 80

CMD /run.sh
