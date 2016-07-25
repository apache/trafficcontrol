FROM centos:6.8

# all ENV vars can be controlled by, e.g. `docker run -e BRANCH=1.6.x <image>`
ENV GITREPO=https://github.com/Comcast/traffic_control
ENV BRANCH=master

ENV GITVERSION=2.9.1
ENV GOVERSION=1.6.2

RUN yum -y update

# These are needed to build up-to-date go and git
RUN yum -y install \
	expat-devel \
	gcc \
	gettext \
	libcurl-devel \
	openssl-devel \
	perl-ExtUtils-MakeMaker

RUN yum -y install \
	rpm-build \
	rsync

# Install up-to-date go and git
RUN	curl https://storage.googleapis.com/golang/go$GOVERSION.linux-amd64.tar.gz | \
	tar -C /usr/local -xzf - && \
	ln -s /usr/local/go/bin/* /usr/local/bin/.
RUN 	curl https://www.kernel.org/pub/software/scm/git/git-$GITVERSION.tar.gz | \
	tar -C / -xzf -
RUN	(cd /git-$GITVERSION && make prefix=/usr/local install && rm -r /git-$GITVERSION)

RUN mkdir -p /traffic_control
VOLUME	/traffic_control

RUN	git clone ${GITREPO} && \
	cd /traffic_control && \
	git checkout ${BRANCH} && \
	cd /traffic_control/traffic_stats && \
	./build/build_rpm.sh
