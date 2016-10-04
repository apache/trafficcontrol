############################################################
# Dockerfile to build Traffic Stats 1.4 container images
# Based on CentOS 6.6
############################################################

# Example Build and Run:
# docker build --rm --tag traffic_stats:1.6 Traffic_Stats
#
# docker run --name my-traffic-stats --hostname my-traffic-stats --net cdnet --env TRAFFIC_OPS_URI=http://my-traffic-ops:3000 --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --env CERT_COUNTRY=US --env CERT_STATE=Colorado --env CERT_CITY=Denver --env CERT_COMPANY=NotComcast --detach traffic_stats:1.4

FROM centos:6.6
MAINTAINER Dan Kirkwood

# Default value for RPM -- override with `docker build --build-arg RPM=...'
ARG INFLUXDB=http://influxdb.s3.amazonaws.com/influxdb-0.11.1-1.x86_64.rpm
ARG GRAFANA=https://grafanarel.s3.amazonaws.com/builds/grafana-3.1.1-1470047149.x86_64.rpm
ARG RPM=http://traffic-control-cdn.net/downloads/1.6.0/RELEASE-1.6.0/traffic_stats-1.6.0-3503.4899d302.x86_64.rpm

RUN yum install -y tar unzip perl-JSON perl-WWW-Curl

ADD $INFLUXDB /
ADD $GRAFANA /
ADD $RPM /

RUN yum install -y /$(basename $INFLUXDB) /$(basename $RPM) /$(basename $GRAFANA)
RUN rm /$(basename $INFLUXDB) /$(basename $RPM) /$(basename $GRAFANA)

EXPOSE 80 8086 8083
ADD run.sh /
ENTRYPOINT /run.sh
