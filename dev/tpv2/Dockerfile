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
	openssl rand 32 | base64 > /aes.key && \
    chmod 644 /server.key /aes.key

FROM node:20-alpine AS tpv2-dev

ENV TC="/go/src/github.com/apache/trafficcontrol"
VOLUME $TC
EXPOSE 443

COPY --from=certbuilder /server.key /server.crt /

CMD ${TC}/dev/tpv2/run.sh
