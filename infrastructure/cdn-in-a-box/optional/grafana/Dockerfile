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

FROM grafana/grafana:10.2.3-ubuntu

USER root
ARG TRAFFIC_TS_RPM=traffic_stats/traffic_stats.rpm

RUN apt-get update && \
    apt-get install -y dnsutils net-tools gettext-base p7zip-full netcat && \
    rm -rf /var/lib/apt/lists/*

ADD enroller/server_template.json \
    traffic_ops/to-access.sh \
    optional/grafana/run-grafana.sh \
    optional/grafana/datasources.yml.template \
    optional/grafana/app.yaml.template \
    $TRAFFIC_TS_RPM \
    /

RUN cd ~ \
    && 7z x /traffic_stats.rpm \
    && 7z e traffic_stats-*.cpio module.js plugin.json module.js.map -r \
    && mkdir -p /var/lib/grafana/plugins/trafficcontrol-scenes-app/ \
    && mv plugin.json module.js.map module.js /var/lib/grafana/plugins/trafficcontrol-scenes-app/ \
    && rm /traffic_stats.rpm \
    && rm ~/traffic_stats-*.cpio

COPY dns/set-dns.sh \
     dns/insert-self-into-dns.sh \
     /usr/local/sbin/

ENTRYPOINT [ "/run-grafana.sh" ]
EXPOSE 443
