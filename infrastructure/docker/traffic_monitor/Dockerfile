############################################################
# Dockerfile to build Traffic Monitor 1.6.0 container images
# Based on CentOS 6.6
############################################################

# Example Build and Run:
# docker build --rm --build-arg JDK=http://download.oracle.com/<path to jdk rpm> --build-arg RPM=<path to traffic_monitor rpm> --tag traffic_monitor:<version> traffic_monitor
#
# docker run --name my-traffic-monitor-0 --hostname my-traffic-monitor-0 --net=cdnet --env TRAFFIC_OPS_URI=http://my-traffic-ops:3000 --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --detach traffic_monitor:1.6.0

FROM centos:6.6
MAINTAINER Dan Kirkwood

# Default values for JDK and RPM -- override with `docker build --build-arg RPM=... --build-arg JDK=...'
ARG JDK=http://download.oracle.com/otn-pub/java/jdk/7u79-b15/jdk-7u79-linux-x64.rpm
ARG RPM=http://traffic-control-cdn.net/downloads/1.6.0/RELEASE-1.6.0/traffic_monitor-1.6.0-3503.4899d302.el6.x86_64.rpm

RUN yum install -y tar unzip perl-JSON perl-WWW-Curl wget

# NOTE: Can't just use ADD because of the required cookie
RUN curl -O -L -b "oraclelicense=a" $JDK
ADD $RPM /

RUN yum install -y /$(basename $JDK) /$(basename $RPM)
RUN rm /$(basename $JDK) /$(basename $RPM)

EXPOSE 80
ADD run.sh /
ENTRYPOINT /run.sh
