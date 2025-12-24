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

FROM rockylinux:8
RUN dnf -y install epel-release \
  && dnf -y install \
  ansible \
  git \
  python3-pip \
  python3-devel \
  libxml2-devel \
  libxslt-devel \
  libffi-devel \
  openssl-devel \
  gcc \
  && dnf clean all \
  && pip3 install --upgrade pip \
  && pip3 install --upgrade setuptools \
  && pip3 install --upgrade pyOpenSSL python-gilt paramiko Jinja2
RUN mkdir -p /opt/atc/ && mkdir ~/.ssh && echo -e "Host *\n   StrictHostKeyChecking no\n   UserKnownHostsFile=/dev/null" > ~/.ssh/config
COPY . /opt/atc

ENTRYPOINT ["/bin/bash"]
