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
FROM docker.io/centos:7

ARG BUILD_VNC_USER=demo
ARG BUILD_VNC_PASSWD=demo

ENV VNC_USER=$BUILD_VNC_USER
ENV VNC_PASSWD=$BUILD_VNC_PASSWD

RUN yum -y install https://download1.rpmfusion.org/free/el/rpmfusion-free-release-7.noarch.rpm epel-release && \
    yum -y install xterm firefox git tigervnc-server sudo bind-utils net-tools which passwd which \
                   fluxbox webcore-fonts terminus-fonts vnc mplayer wget openssl curl supervisor \
                   tinyproxy openssh-clients openssh-server 

# Disabled Xfce (too large)
#RUN yum -y groupinstall Xfce
RUN yum clean all

# Create testclient/vnc/proxy user(s)
RUN useradd -m $VNC_USER && \
    echo "$VNC_USER ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

USER $VNC_USER
WORKDIR /home/$VNC_USER

RUN rm -rf /home/$VNC_USER/.vnc && \
    mkdir /home/$VNC_USER/.vnc && \
    mkdir /home/$VNC_USER/.fluxbox && \
    fluxbox-generate_menu -k -g -B -su -t xterm -b firefox && \
    sed 's/{xterm}/{xterm -bg black -fg white +sb}/g' -i /home/$VNC_USER/.fluxbox/menu && \
    echo "$VNC_USER" | vncpasswd -f > /home/$VNC_USER/.vnc/passwd && \
    chmod 600 /home/$VNC_USER/.vnc/passwd

ADD testclient/vnc_startup.sh /home/$VNC_USER/.vnc/xstartup 

USER root

ADD testclient/run.sh \
    traffic_ops/to-access.sh \
    /

RUN chmod +x /home/$VNC_USER/.vnc/xstartup

VOLUME /ssl

##########################
# Ports 
##########################
# - 2200 - SSH 
# - 5900 - VNC Port
# - 7070 - Tinyproxy HTTP 
# - 9090 - Socks5 Proxy
##########################
EXPOSE 2200/tcp 5900/tcp 7070/tcp 9090/tcp

RUN yum --enablerepo=extras install -y docker-client-latest

CMD /run.sh
