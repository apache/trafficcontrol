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

WORKDIR /vol/traffic_control/traffic_portal/build

# bower will not run as root by default
RUN	echo '{ "allow_root": true }' > /root/.bowerrc
CMD	./build_rpm.sh
