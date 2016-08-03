FROM centos:7

MAINTAINER Dan Kirkwood

RUN	yum -y update && \
	yum -y clean all
RUN	yum -y install \
		epel-release \
		git \
		rpm-build

RUN	yum -y install \
		libffi-devel \
		make \
		nodejs \
		npm \
		ruby-devel \
		rubygems

RUN	yum -y clean all

RUN	gem install compass
RUN	npm -g install bower grunt-cli

WORKDIR /vol/traffic_control/traffic_portal/build

RUN	yum -y install ant.noarch

# bower will not run as root by default
RUN	echo '{ "allow_root": true }' > /root/.bowerrc
ENTRYPOINT ./build_rpm.sh
