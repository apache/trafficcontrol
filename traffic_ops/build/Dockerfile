FROM centos:6.8

MAINTAINER Dan Kirkwood

RUN	yum -y install \
		epel-release \
		git \
		rpm-build

# install latest git
RUN yum install -y \
	expat-devel \
	gcc \
	gettext \
	libcurl-devel \
	openssl-devel \
	perl-ExtUtils-MakeMaker \
	tar \
	tcl && \
	git clone https://github.com/git/git.git && \
	cd git && \
	make prefix=/usr/local all && \
	make prefix=/usr/local install


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
