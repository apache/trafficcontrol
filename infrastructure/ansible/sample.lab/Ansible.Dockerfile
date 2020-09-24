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

FROM centos:7.4.1708
MAINTAINER Jonathan Gray
RUN yum -y install epel-release \
  && yum -y install \
  ansible \
  git \
  python-pip \
  python-devel \
  libxml2-devel \
  libxslt-devel \
  libffi-devel \
  openssl-devel \
  gcc \
  && yum clean all \
  && pip install --upgrade pip \
  && pip install --upgrade setuptools \
  && pip install --upgrade pyOpenSSL python-gilt paramiko Jinja2
RUN mkdir -p /opt/atc/ && mkdir ~/.ssh && echo -e "Host *\n   StrictHostKeyChecking no\n   UserKnownHostsFile=/dev/null" > ~/.ssh/config
COPY . /opt/atc

ENTRYPOINT ["/bin/bash"]
