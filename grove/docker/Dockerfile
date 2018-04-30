# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Example Build and Run:
# docker build --rm --tag grove:0.1 --build-arg=RPM=grove-0.1.127-1.x86_64.rpm .
#
# docker run --name my-grove-cache --hostname my-grove-cache --net cdnet --env REMAP_PATH=/config/remap.json -v config:/config --detach grove:0.1
#


FROM centos/systemd

RUN yum install -y initscripts epel-release openssl

ARG RPM=grove.rpm
ADD $RPM /
RUN yum install -y /$(basename $RPM)

RUN setcap 'cap_net_bind_service=+ep' /usr/sbin/grove

EXPOSE 80 443
ADD docker-entrypoint.sh /
ENTRYPOINT /docker-entrypoint.sh
