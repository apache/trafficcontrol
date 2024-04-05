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

FROM alpine:latest AS build-trafficserver
ARG ATS_VERSION
ADD https://archive.apache.org/dist/trafficserver/trafficserver-${ATS_VERSION}.tar.bz2 /tmp/
RUN set -o errexit -o nounset -o xtrace; \
    cd tmp; \
    dirname=trafficserver-${ATS_VERSION}; \
    tar xf ${dirname}.tar.bz2; \
    rm ${dirname}.tar.bz2; \
    apk add --no-cache \
        # configure dependencies
        g++ \
        perl \
        openssl-dev \
        pcre-dev \
        make \
        # build dependencies
        libexecinfo-dev \
        fortify-headers \
        linux-headers \
        zlib-dev; \
    cd $dirname; \
    ./configure \
        --disable-tests \
        --enable-experimental-plugins \
        --prefix=/ \
        --with-user=ats \
        --with-group=ats; \
    # Not limiting threads makes the GitHub Actions Hosted Agent lose
    # communication with the server due to memory exhaustion, making the
    # Action fail.
    threads=$(( $(nproc) )); \
    make -j $threads; \
    adduser -D ats; \
    make install DESTDIR=/tmp/built; \
    cd ..; \
    rm -r $dirname

FROM alpine:latest
COPY --from=build-trafficserver /tmp/built/ /
RUN apk add --no-cache \
        # runtime dependencies
        libexecinfo \
        libstdc++ \
        pcre && \
    adduser -D ats
USER ats
CMD /bin/traffic_server
