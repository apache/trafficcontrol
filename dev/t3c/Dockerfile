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
FROM alpine:latest AS traffic-server-builder

RUN apk add --no-cache build-base perl libexecinfo-dev pcre-dev libressl-dev libtool linux-headers openssl-dev zlib-dev
ADD https://downloads.apache.org/trafficserver/trafficserver-9.1.1.tar.bz2 /
RUN tar -xf trafficserver-9.1.1.tar.bz2 && cd trafficserver-9.1.1 && mkdir /ats && ./configure --prefix / --enable-experimental-plugins && make -j && make install && ls -R /ats

FROM alpine:latest AS t3c-dev

ENV TC="/root/go/src/github.com/apache/trafficcontrol/"
VOLUME /root/go/src/github.com/apache/trafficcontrol
EXPOSE 80 8081

COPY --from=traffic-server-builder /bin/traffic_cache_tool /bin/traffic_crashlog /bin/traffic_ctl /bin/traffic_layout /bin/traffic_logcat /bin/traffic_logstats /bin/traffic_manager /bin/traffic_server /bin/traffic_via /bin/trafficserver /bin/tspush /bin/tsxs /bin/
COPY --from=traffic-server-builder /etc/trafficserver /etc/trafficserver
COPY --from=traffic-server-builder /include/ts /include/ts
COPY --from=traffic-server-builder /include/tscpp /include/tscpp
COPY --from=traffic-server-builder /lib/perl5 /lib/perl5
COPY --from=traffic-server-builder /lib/pkgconfig/trafficserver.pc /lib/pkgconfig/trafficserver.pc
COPY --from=traffic-server-builder /lib/libtscore.la /lib/libtscore.so /lib/libtscore.so.9 /lib/libtscore.so.9.1.1 /lib/libtscppapi.la /lib/libtscppapi.so /lib/libtscppapi.so.9 /lib/libtscppapi.so.9.1.1 /lib/libtscpputil.la /lib/libtscpputil.so /lib/libtscpputil.so.9 /lib/libtscpputil.so.9.1.1 /lib/libtsmgmt.la /lib/libtsmgmt.so /lib/libtsmgmt.so.9 /lib/libtsmgmt.so.9.1.1 /lib/plugin_init_fail.la /lib/plugin_init_fail.so /lib/plugin_instinit_fail.la /lib/plugin_instinit_fail.so /lib/plugin_missing_deleteinstance.la /lib/plugin_missing_deleteinstance.so /lib/plugin_missing_doremap.la /lib/plugin_missing_doremap.so /lib/plugin_missing_init.la /lib/plugin_missing_init.so /lib/plugin_missing_newinstance.la /lib/plugin_missing_newinstance.so /lib/plugin_required_cb.la /lib/plugin_required_cb.so /lib/plugin_testing_calls.la /lib/plugin_testing_calls.so /lib/plugin_v1.la /lib/plugin_v1.so /lib/plugin_v2.la /lib/plugin_v2.so /lib/
COPY --from=traffic-server-builder /libexec/trafficserver /libexec/trafficserver
COPY --from=traffic-server-builder /share/man /share/man
RUN mkdir /share/trafficserver && mkdir -p /var/log/trafficserver && mkdir /var/trafficserver

RUN apk add --no-cache make inotify-tools go pcre libexecinfo && \
	go install github.com/go-delve/delve/cmd/dlv@latest && \
	ln -s /root/go/bin/dlv /usr/bin/dlv && \
	mkdir /lib64 && \
	ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN addgroup ats && adduser -S -s /bin/ash -G ats ats && chown -R ats:ats /etc/trafficserver /var/trafficserver /share/trafficserver /var/log/trafficserver
RUN echo "stats_over_http.so" >> /etc/trafficserver/plugin.config && echo "system_stats.so" >> /etc/trafficserver/plugin.config

CMD /root/go/src/github.com/apache/trafficcontrol/dev/t3c/run.sh
