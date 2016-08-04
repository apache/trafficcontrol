FROM centos:7

MAINTAINER Dan Kirkwood

RUN	yum -y update && \
	yum -y clean all
RUN	yum -y install \
		epel-release \
		git \
		rpm-build

# all ENV vars can be controlled by, e.g. `docker run -e BRANCH=1.6.x <image>`
ENV GITREPO=https://github.com/Comcast/traffic_control
ENV BRANCH=master

# clone/checkout using above env vars
ADD clone-for-docker.sh /
CMD /clone-for-docker.sh
