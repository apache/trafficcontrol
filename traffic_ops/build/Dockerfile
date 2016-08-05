FROM traffic_control_builder:latest

MAINTAINER Dan Kirkwood

RUN	yum -y install \
		golang && \
	yum -y clean all


ENV	GITREPO=https://github.com/Comcast/traffic_control
ENV	BRANCH=master

CMD	/clone-for-docker.sh $GITREPO $BRANCH && \
		cd /vol/traffic_control/traffic_ops/build && \
		./build_rpm.sh
