FROM centos:7

MAINTAINER Dan Kirkwood

RUN	yum -y update && \
	yum -y clean all
RUN	yum -y install \
		epel-release \
		git \
		rpm-build

RUN	yum -y install \
		golang && \
	yum -y clean all

WORKDIR /repo/traffic_control/traffic_stats/build
ENTRYPOINT ./build_rpm.sh
