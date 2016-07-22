############################################################
# Dockerfile to build Traffic Stats 1.4 container images
# Based on CentOS 6.6
############################################################

# Example Build and Run:
# docker build --file Traffic_Stats_Dockerfile --rm --tag traffic_stats:1.4 .
# docker run --name my-traffic-stats --hostname my-traffic-stats --net cdnet --env TRAFFIC_OPS_URI=http://my-traffic-ops:3000 --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --env CERT_COUNTRY=US --env CERT_STATE=Colorado --env CERT_CITY=Denver --env CERT_COMPANY=NotComcast --detach traffic_stats:1.4

FROM centos:6.6
MAINTAINER Robert Butts

RUN curl -O http://traffic-control-cdn.net/downloads/1.4.0/RELEASE-1.4.0-RC10/traffic_stats-1.4.0-2426.2876438a.x86_64.rpm
RUN yum install -y wget tar unzip perl-JSON perl-WWW-Curl
RUN yum install -y traffic_stats-1.4.0-2426.2876438a.x86_64.rpm
RUN curl -O http://influxdb.s3.amazonaws.com/influxdb-0.9.5.1-1.x86_64.rpm
RUN yum install -y influxdb-0.9.5.1-1.x86_64.rpm

RUN yum localinstall -y https://grafanarel.s3.amazonaws.com/builds/grafana-2.6.0-1.x86_64.rpm

EXPOSE 80 8086 8083
ADD Traffic_Stats_Dockerfile_run.sh /
ENTRYPOINT /Traffic_Stats_Dockerfile_run.sh
