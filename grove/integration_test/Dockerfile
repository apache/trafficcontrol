FROM centos:7
MAINTAINER The CentOS Project <cloud-ops@centos.org>
LABEL Vendor="CentOS" \
      License=GPLv2 \
      Version=2.4.6-40


RUN yum -y --setopt=tsflags=nodocs update && \
    yum -y --setopt=tsflags=nodocs install httpd && \
    yum -y --setopt=tsflags=nodocs install perl && \
    yum -y --setopt=tsflags=nodocs install git && \
    yum -y --setopt=tsflags=nodocs install golang && \
    yum -y --setopt=tsflags=nodocs install openssl && \
    yum clean all

#EXPOSE 80

# Simple startup script to avoid some issues observed with container restart
ADD setup-and-run.sh setup-and-run.sh /
RUN chmod -v +x /setup-and-run.sh
ADD remap-base-test.json /remap-base-test.json
ADD grove.cfg /grove.cfg
ADD tests /tests
ADD compare_gets.go /compare_gets.go

CMD ["/setup-and-run.sh"]
