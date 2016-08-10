FROM centos:7

MAINTAINER Dan Kirkwood

RUN	yum -y install \
		epel-release \
		git \
		rpm-build

# all ENV vars can be controlled by, e.g. `docker run -e BRANCH=1.6.x <image>`
ENV GITREPO=https://github.com/Comcast/traffic_control
ENV BRANCH=master

#############################################################################
# Above are common for all sub-projects.
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

RUN	mkdir -p /dist
VOLUME	/dist

WORKDIR /
CMD	git clone --branch $BRANCH $GITREPO && cd traffic_control && \
	./traffic_portal/build/build_rpm.sh && \
	cp -f dist/* /dist/.
