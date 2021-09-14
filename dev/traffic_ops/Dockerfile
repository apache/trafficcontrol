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
FROM alpine:latest AS certbuilder
RUN apk add --no-cache openssl
RUN openssl genrsa -passout pass:x -out server.pass.key 2048 && \
	openssl rsa -passin pass:x -in server.pass.key -out server.key && \
	rm server.pass.key && \
	openssl req -new -key server.key -out server.csr \
		-subj "/C=US/ST=CO/L=Denver/O=Apache/OU=Traffic Control/CN=trafficops.dev.ciab.test" && \
	openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt && \
	openssl rand 32 | base64 > /aes.key

FROM alpine:latest AS trafficops-dev

ENV TC="/root/go/src/github.com/apache/trafficcontrol/"
VOLUME /root/go/src/github.com/apache/trafficcontrol
ENV ADMIN="$TC/traffic_ops/app/db/admin"
EXPOSE 443 6444

COPY --from=certbuilder /server.key /server.crt /aes.key /
RUN apk add --no-cache make inotify-tools go postgresql-client && go install github.com/go-delve/delve/cmd/dlv@latest && ln -s /root/go/bin/dlv /usr/bin/dlv
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY .pgpass /root/.pgpass
RUN chmod 0600 /root/.pgpass

ENTRYPOINT /root/go/src/github.com/apache/trafficcontrol/dev/traffic_ops/run.sh
