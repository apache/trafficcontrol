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
ARG GO_VERSION
FROM alpine:latest AS certbuilder
RUN apk add --no-cache openssl
RUN openssl genrsa -passout pass:x -out server.pass.key 2048 && \
	openssl rsa -passin pass:x -in server.pass.key -out server.key && \
	rm server.pass.key && \
	openssl req -new -key server.key -out server.csr \
		-subj "/C=US/ST=CO/L=Denver/O=Apache/OU=Traffic Control/CN=trafficops.dev.ciab.test" && \
	openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt && \
	openssl rand 32 | base64 > /aes.key && \
	chmod 644 /server.key /aes.key

FROM golang:${GO_VERSION}-alpine AS trafficops-dev

ENV TC="/go/src/github.com/apache/trafficcontrol" GOFLAGS="--buildvcs=false"
VOLUME $TC
ENV ADMIN="$TC/traffic_ops/app/db/admin"
EXPOSE 443 6444

COPY --from=certbuilder /server.key /server.crt /aes.key /
RUN apk add --no-cache make inotify-tools postgresql-client gcc libc-dev && \
	go install github.com/go-delve/delve/cmd/dlv@latest && \
	rm -rf $GOPATH/pkg/*
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY .pgpass /root/.pgpass
RUN chmod 0600 /root/.pgpass

# These versions don't actually do anything - because the dev env can't make ISO
# system images - but having them defined anyway helps when testing related
# functionality, to a point.
RUN mkdir -p /var/www/files/ && echo '{"CentOS 7": "centos7", "Rocky Linux 8": "rocky8"}' > /var/www/files/osversions.json

CMD $TC/dev/traffic_ops/run.sh
