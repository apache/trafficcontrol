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

    # Change BASE_IMAGE to centos when RHEL_VERSION=7
ARG BASE_IMAGE=rockylinux \
    RHEL_VERSION=8
FROM ${BASE_IMAGE}:${RHEL_VERSION}
ARG RHEL_VERSION=8

RUN if [[ "${RHEL_VERSION%%.*}" -eq 7 ]]; then \
        yum -y install dnf || exit 1; \
    fi

ARG VNC_BUILD_USER
ENV VNC_USER=$VNC_BUILD_USER

RUN dnf -y install "https://download1.rpmfusion.org/free/el/rpmfusion-free-release-${RHEL_VERSION%%.*}.noarch.rpm" epel-release && \
    dnf -y install xterm firefox git tigervnc-server sudo bind-utils net-tools which passwd which \
                   fluxbox 'google-noto*' terminus-fonts tigervnc vlc wget openssl curl nc && \
    dnf -y clean all && rm -rf /var/cache/dnf

RUN if [[ "${RHEL_VERSION%%.*}" -ge 8 ]]; then \
        curl -Lo /usr/bin/vncserver https://git.centos.org/rpms/tigervnc/raw/9e6ab1bc80/f/SOURCES/vncserver; \
    fi

RUN useradd -m $VNC_USER && \
    echo "$VNC_USER ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

USER $VNC_USER
WORKDIR /home/$VNC_USER

RUN rm -rf /home/$VNC_USER/.vnc && \
    mkdir /home/$VNC_USER/.vnc && \
    mkdir /home/$VNC_USER/.fluxbox && \
    fluxbox-generate_menu -k -g -B -su -t xterm -b firefox && \
    sed 's/{xterm}/{xterm -bg black -fg white +sb}/g' -i /home/$VNC_USER/.fluxbox/menu && \
    dd if=/dev/urandom of=/dev/stdout count=12 bs=1 | vncpasswd -f > /home/$VNC_USER/.vnc/passwd && \
    chmod 600 /home/$VNC_USER/.vnc/passwd

ADD optional/vnc/vnc_startup.sh /home/$VNC_USER/.vnc/xstartup 

USER root

RUN systemd-machine-id-setup && \
    chmod +x /home/$VNC_USER/.vnc/xstartup

ADD optional/vnc/run.sh /

COPY dns/set-dns.sh \
     dns/insert-self-into-dns.sh \
     /usr/local/sbin/

EXPOSE 5909/tcp

CMD /run.sh
