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
		which && \
	yum -y clean all

### java/mvn
WORKDIR /opt
RUN curl -L -b "oraclelicense=a" http://download.oracle.com/otn-pub/java/jdk/8u101-b13/jdk-8u101-linux-x64.tar.gz | tar xzf - && \
    curl -L http://mirror.nexcess.net/apache/maven/maven-3/3.3.9/binaries/apache-maven-3.3.9-bin.tar.gz | tar xzf -
RUN ln -s jdk* java && ln -s apache-maven* maven

RUN alternatives --install /usr/bin/java java /opt/java/bin/java 2 && \
    alternatives --install /usr/bin/jar jar /opt/java/bin/jar 2 && \
    alternatives --install /usr/bin/javac javac /opt/java/bin/javac 2 && \
    alternatives --set java /opt/java/bin/java && \
    alternatives --set jar /opt/java/bin/jar && \
    alternatives --set javac /opt/java/bin/javac && \
    alternatives --install /usr/bin/mvn mvn /opt/maven/bin/mvn 1
###

RUN	mkdir -p /dist
VOLUME	/dist

WORKDIR /
CMD	git clone --branch $BRANCH $GITREPO && cd traffic_control && \
	./traffic_router/build/build_rpm.sh && \
	cp -f dist/* /dist/.
