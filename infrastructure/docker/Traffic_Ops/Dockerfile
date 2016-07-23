############################################################
# Dockerfile to build Traffic Ops 1.4 container images
# Based on CentOS 6.6
############################################################

# Example Build and Run:
# docker network create cdnet
# docker build --file Traffic_Ops_Dockerfile --rm --tag traffic_ops:1.4 .
# docker run --name my-traffic-ops-mysql --hostname my-traffic-ops-mysql --net cdnet --env MYSQL_ROOT_PASSWORD=secretrootpass --detach mysql:5.5
# docker run --name my-traffic-ops --hostname my-traffic-ops --net cdnet --publish 443:443 --env MYSQL_IP=my-traffic-ops-mysql --env MYSQL_PORT=3306 --env MYSQL_ROOT_PASS=secretrootpass --env MYSQL_TRAFFIC_OPS_PASS=supersecretpassword --env ADMIN_USER=superroot --env ADMIN_PASS=supersecreterpassward --env CERT_COUNTRY=US --env CERT_STATE=Colorado --env CERT_CITY=Denver --env CERT_COMPANY=NotComcast --env TRAFFIC_VAULT_PASS=marginallylesssecret --env DOMAIN=cdnet --detach traffic_ops:1.4

FROM centos:6.6
MAINTAINER Robert Butts

RUN yum install -y tar gcc-c++ perl cpan openssl-devel mysql-devel expect

RUN (echo y;echo o conf prerequisites_policy follow;echo o conf urllist push http://cpan.yimg.com/;echo o conf urllist push ftp://cpan.cs.utah.edu/CPAN/;echo o conf urllist push ftp://cpan.cs.utah.edu/CPAN/;echo o conf urllist push ftp://mirrors.syringanetworks.net/CPAN/;echo o conf commit;echo install YAML;echo install IPC::Cmd;echo install HTTP::Tiny;echo install CPAN;echo reload cpan)|cpan

RUN cpan Carton

RUN curl -O http://traffic-control-cdn.net/downloads/1.4.0/RELEASE-1.4.0-RC10/traffic_ops-1.4.0-2426.2876438a.x86_64.rpm
RUN yum install -y traffic_ops-1.4.0-2426.2876438a.x86_64.rpm

RUN cd /opt/traffic_ops/app && carton

# \todo Workaround for PR #1153. Remove when fixed.
RUN sed -i -- 's/system("mysql $db_name/system("mysql -h $host_ip -P $host_port $db_name/g' /opt/traffic_ops/app/db/admin.pl
RUN sed -i -- 's/my $host_ip/$host_ip/g' /opt/traffic_ops/app/db/admin.pl
RUN sed -i -- 's/my $host_port/$host_port/g' /opt/traffic_ops/app/db/admin.pl
RUN sed -i -- 's/my \$db_password = '\'''\'';/my \$db_password = '\'''\'';\nmy \$host_ip = '\'''\'';\nmy \$host_port = '\'''\'';/g' /opt/traffic_ops/app/db/admin.pl

RUN export PERL5LIB=/opt/traffic_ops/app/local/lib/perl5/:/opt/traffic_ops/install/lib/ \
&& export TERM=xterm \
&& export USER=root \
&& yes | /opt/traffic_ops/install/bin/build_trafficops_perl_library -i

RUN export PERL5LIB=/opt/traffic_ops/app/local/lib/perl5/:/opt/traffic_ops/install/lib/ \
&& export TERM=xterm \
&& export USER=root \
&& /opt/traffic_ops/install/bin/download_web_deps -i

# \todo Workaround for issue #1016. Remove when fixed.
RUN sed -i -- "s#/opt/rascal/conf#/opt/traffic_monitor/conf#g" /opt/traffic_ops/install/data/profiles/profile.traffic_monitor.traffic_ops

RUN sed -i -- 's#],#,{"config_file":"CRConfig.json","name":"api.auth.url","value":"http://${toHostname}/api/1.1/user/login"}],#g' /opt/traffic_ops/install/data/profiles/profile.traffic_router.traffic_ops

# \todo add Drive Letters to postinstall input
RUN sed -i -- 's/"value": "b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y"/"value": "0"/g' /opt/traffic_ops/install/data/profiles/profile.trafficserver_edge.traffic_ops
RUN sed -i -- 's/"value": "b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y"/"value": "0"/g' /opt/traffic_ops/install/data/profiles/profile.trafficserver_mid.traffic_ops

RUN echo "{\"user\": \"riakuser\",\"password\": \"$TRAFFIC_VAULT_PASS\"}" > /opt/traffic_ops/app/conf/production/riak.conf

EXPOSE 443
WORKDIR /opt/traffic_ops/app
ENV MOJO_MODE production
ADD profile.origin.traffic_ops /
ADD Traffic_Ops_Dockerfile_run.sh /
ENTRYPOINT /Traffic_Ops_Dockerfile_run.sh
