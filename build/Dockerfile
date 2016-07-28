FROM centos:7

MAINTAINER Dan Kirkwood

RUN	yum -y install git

# all ENV vars can be controlled by, e.g. `docker run -e BRANCH=1.6.x <image>`
ENV GITREPO=https://github.com/Comcast/traffic_control
ENV BRANCH=master

# clone/checkout using above env vars
ADD clone-for-docker.sh /
ENTRYPOINT /clone-for-docker.sh
VOLUME /repo
