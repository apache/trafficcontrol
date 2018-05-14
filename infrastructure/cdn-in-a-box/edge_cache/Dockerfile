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
# Dockerfile to build Edge Cache instance (as a reverse proxy)
############################################################

FROM ats:latest

# Modify configs as necessary (keeping backups just in case)
WORKDIR /opt/traffic_server/etc/trafficserver
COPY infrastructure/cdn-in-a-box/edge_cache/config-edit.sh .
RUN chmod +x config-edit.sh
RUN ./config-edit.sh

# If all's well, serve content forever
WORKDIR /opt/traffic_server/
RUN cat etc/trafficserver/remap.config
EXPOSE 8080
CMD /opt/traffic_server/bin/traffic_server -o3
