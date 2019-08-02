#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
FROM golang:1.11
MAINTAINER dev@trafficcontrol.apache.org

# ADD traffic_ops /go/src/github.com/apache/trafficcontrol/traffic_ops
# ADD lib /go/src/github.com/apache/trafficcontrol/lib
# WORKDIR /go/src/github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang

ADD entrypoint.sh /

# CMD bash -c 'go get -v && go test -cover -v ./... ../../lib/go-tc/...'
ENTRYPOINT ["/entrypoint.sh"]
