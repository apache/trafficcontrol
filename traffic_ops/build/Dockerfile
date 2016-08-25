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
		golang && \
	yum -y clean all

RUN	mkdir -p /dist
VOLUME	/dist

WORKDIR /
CMD	echo "GITREPO=$GITREPO BRANCH=$BRANCH" && git clone --branch $BRANCH $GITREPO && cd traffic_control && git status && \
	./traffic_ops/build/build_rpm.sh && \
	cp -f dist/* /dist/.
