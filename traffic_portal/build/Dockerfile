FROM traffic_control_builder:latest

MAINTAINER Dan Kirkwood

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

# bower will not run as root by default
RUN	echo '{ "allow_root": true }' > /root/.bowerrc

ENV	GITREPO=https://github.com/Comcast/traffic_control
ENV	BRANCH=master

CMD	/clone-for-docker.sh $GITREPO $BRANCH && \
		cd /vol/traffic_control/traffic_portal/build && \
		./build_rpm.sh
