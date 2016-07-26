FROM centos:7

MAINTAINER Dan Kirkwood

# all ENV vars can be controlled by, e.g. `docker run -e BRANCH=1.6.x <image>`
ENV GITREPO=https://github.com/Comcast/traffic_control
ENV BRANCH=master

RUN	yum -y update && \
	yum -y clean all
RUN	yum -y install \
		epel-release \
		git \
		rpm-build

RUN	yum -y install \
		golang && \
	yum -y clean all

RUN mkdir -p /traffic_control
VOLUME	/traffic_control
ADD run.sh /
ENTRYPOINT /run.sh
