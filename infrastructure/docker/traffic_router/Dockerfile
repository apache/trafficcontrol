############################################################
# Dockerfile to build Traffic Router 1.4 container images
# Based on CentOS 6.6
############################################################

# Example Build and Run:
# docker build --rm --tag traffic_router:1.4 Traffic_Router
# docker run --name my-traffic-router --hostname my-traffic-router --net cdnet --env TRAFFIC_OPS_URI=http://my-traffic-ops:3000 --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --env TRAFFIC_MONITORS="my-traffic-monitor-0:80;my-traffic-monitor-1:80" --env ORIGIN_URI="http://my-hotair.cdnet" --detach traffic_router:1.4

FROM centos:6.6
MAINTAINER Robert Butts

RUN curl -O http://traffic-control-cdn.net/downloads/1.4.0/RELEASE-1.4.0-RC10/traffic_router-1.4.0-2426.2876438a.el6.x86_64.rpm
RUN yum install -y wget tar unzip perl-JSON perl-WWW-Curl
RUN wget --no-cookies --no-check-certificate --header "Cookie: gpw_e24=http%3A%2F%2Fwww.oracle.com%2F; oraclelicense=accept-securebackup-cookie" "http://download.oracle.com/otn-pub/java/jdk/7u79-b15/jdk-7u79-linux-x64.rpm"
RUN yum install -y jdk-7u79-linux-x64.rpm
RUN yum install -y traffic_router-1.4.0-2426.2876438a.el6.x86_64.rpm

EXPOSE 53 80 3333
ADD run.sh /
ENTRYPOINT /run.sh
