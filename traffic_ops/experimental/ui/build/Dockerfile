FROM centos:7

MAINTAINER Dan Kirkwood

RUN	yum -y install \
		epel-release \
		git \
		rpm-build && \
	yum -y clean all

ENV GITREPO https://github.com/apache/incubator-trafficcontrol.git
ENV BRANCH master

WORKDIR /repo

# traffic_ops_v2 specific
RUN	yum -y install \
		gcc \
		libffi-devel \
		make \
		nodejs \
		npm \
		ruby-devel \
		rubygems

RUN	gem install compass
RUN	npm -g install bower grunt-cli

# bower will not run as root by default
RUN	echo '{ "allow_root": true }' > /root/.bowerrc

###

VOLUME ./artifacts:/artifacts

ADD run-build.sh .

CMD ./run-build.sh traffic_ops_v2 2>&1 | tee /artifacts/build-traffic_ops_v2.log

# vi:syntax=Dockerfile
