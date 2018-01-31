# Example Build and Run:
# docker build --rm --tag grove:0.1 --build-arg=RPM=grove-0.1.127-1.x86_64.rpm .
#
# docker run --name my-grove-cache --hostname my-grove-cache --net cdnet --env REMAP_PATH=/config/remap.json -v config:/config --detach grove:0.1
#


FROM centos/systemd

RUN yum install -y initscripts epel-release openssl

ARG RPM=grove.rpm
ADD $RPM /
RUN yum install -y /$(basename $RPM)

RUN setcap 'cap_net_bind_service=+ep' /usr/sbin/grove

EXPOSE 80 443
ADD docker-entrypoint.sh /
ENTRYPOINT /docker-entrypoint.sh
