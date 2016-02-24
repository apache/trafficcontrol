############################################################
# Dockerfile to build Traffic Monitor 1.4 container images
# Based on CentOS 6.6
############################################################

# Example Build and Run:
# docker build --file Traffic_Monitor_Dockerfile --tag traffic_monitor:1.4 .
# docker run --name my-traffic-monitor-0 --hostname my-traffic-monitor-0 --net=cdnet --env TRAFFIC_OPS_URI=http://my-traffic-ops:3000 --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --detach traffic_monitor:1.4

FROM centos:6.6
MAINTAINER Robert Butts

RUN curl -O http://traffic-control-cdn.net/downloads/1.4.0/RELEASE-1.4.0-RC10/traffic_monitor-1.4.0-2426.2876438a.el6.x86_64.rpm
RUN yum install -y wget tar unzip perl-JSON perl-WWW-Curl
RUN wget --no-cookies --no-check-certificate --header "Cookie: gpw_e24=http%3A%2F%2Fwww.oracle.com%2F; oraclelicense=accept-securebackup-cookie" "http://download.oracle.com/otn-pub/java/jdk/7u79-b15/jdk-7u79-linux-x64.rpm"
RUN yum install -y jdk-7u79-linux-x64.rpm
RUN yum install -y traffic_monitor-1.4.0-2426.2876438a.el6.x86_64.rpm

EXPOSE 80
ADD Traffic_Monitor_Dockerfile_run.sh /
ENTRYPOINT /Traffic_Monitor_Dockerfile_run.sh
