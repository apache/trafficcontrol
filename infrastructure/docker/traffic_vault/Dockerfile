############################################################
# Dockerfile to build Riak container images
#   as Traffic Vault for Traffic Control 1.6.0
# Based on CentOS 6.6
############################################################

# Example Build and Run:
# docker build --rm --tag traffic_vault:1.6.0 traffic_vault
# docker run --name my-traffic-vault --hostname my-traffic-vault --net cdnet --env ADMIN_PASS=riakadminsecret --env USER_PASS=marginallylesssecret --env CERT_COUNTRY=US --env CERT_STATE=Colorado --env CERT_CITY=Denver --env CERT_COMPANY=NotComcast --env TRAFFIC_OPS_URI=http://my-traffic-ops:3000 --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --env DOMAIN=cdnet --detach traffic_vault:1.6.0

FROM centos:6.6
MAINTAINER Dan Kirkwood

ARG ADMIN_PASS
ARG USER_PASS

# download and install Riak
RUN curl -s https://packagecloud.io/install/repositories/basho/riak/script.rpm.sh | bash
RUN yum install -y riak-2.1.1-1.el6.x86_64

# Set the Riak certs in the config (this cert+key will be created in the run.sh script
RUN sed -i -- 's/## ssl.certfile = $(platform_etc_dir)\/cert.pem/ssl.certfile = \/etc\/riak\/certs\/server.crt/g' /etc/riak/riak.conf
RUN sed -i -- 's/## ssl.keyfile = $(platform_etc_dir)\/key.pem/ssl.keyfile = \/etc\/riak\/certs\/server.key/g' /etc/riak/riak.conf
RUN sed -i -- 's/## ssl.cacertfile = $(platform_etc_dir)\/cacertfile.pem/ssl.cacertfile = \/etc\/riak\/certs\/ca-bundle.crt/g' /etc/riak/riak.conf

RUN sed -i -- "s/nodename = riak@127.0.0.1/nodename = riak@0.0.0.0/g" /etc/riak/riak.conf
RUN sed -i -- "s/listener.http.internal = 127.0.0.1:8098/listener.http.internal = 0.0.0.0:8098/g" /etc/riak/riak.conf
RUN sed -i -- "s/listener.protobuf.internal = 127.0.0.1:8087/listener.protobuf.internal = 0.0.0.0:8087/g" /etc/riak/riak.conf
RUN sed -i -- "s/## listener.https.internal = 127.0.0.1:8098/listener.https.internal = 0.0.0.0:8088/g" /etc/riak/riak.conf

RUN mkdir /etc/riak/certs

RUN echo "tls_protocols.tlsv1 = on" >> /etc/riak/riak.conf 

EXPOSE 8098 8087 8088
ADD run.sh /
ENTRYPOINT /run.sh
