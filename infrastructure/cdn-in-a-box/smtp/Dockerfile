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

FROM maildev/maildev:2.0.0-beta3

USER root
RUN apk add --no-cache \
		bash \
		bind-tools

COPY traffic_ops/to-access.sh /
COPY smtp/run.sh /usr/bin/

COPY dns/set-dns.sh \
     dns/insert-self-into-dns.sh \
     /usr/local/sbin/

# Unset entrypoint
ENTRYPOINT []
CMD ["/usr/bin/env", "run.sh"]

HEALTHCHECK --interval=10s --timeout=1s \
	CMD bash -c 'source /to-access.sh && [[ "$(wget -qO- https://$SMTP_FQDN/healthz)" == true ]]'
