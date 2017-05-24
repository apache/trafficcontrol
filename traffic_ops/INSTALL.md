# From new CentOS 7 install:

## Disable selinux:
### change `/etc/selinux/config` to `SELINUX=disabled`

## Add access to Postgreql 9.6 yum repository 

Instructions are here: https://yum.postgresql.org/

- From this page,  copy the link for CentOS 7 and install:

    $ sudo yum install https://download.postgresql.org/pub/repos/yum/9.6/redhat/rhel-7-x86_64/pgdg-centos96-9.6-xxxx.noarch.rpm
    
## Install Postgres 9.6 server (in a container or on the host)

### on the host:

    $ sudo su -
    # yum install postgresql96-server
    $ su - postgres
    $ /usr/pgsql-9.6/bin/initdb -A md5 -W #-W forces the user to provide a superuser (postgres) password
    $ exit
    # systemctl start postgresql-9.6
    # systemctl status postgresql-9.6

### -or- in a container

NOTE: you do *not* need postgresql96-server if running postgres within a `docker` container.

Install `docker` and `docker-compose` using instructions here:

    https://docs.docker.com/engine/installation/linux/centos/
    
    https://docs.docker.com/compose/install/


## Install `traffic_ops`

    $ sudo yum install traffic_ops

## Install any extensions needed

   - install in /opt/traffic_ops_extensions
   
## Install `openssl` certs (or use this to generate them)

   $ sudo /opt/traffic_ops/install/bin/generateCert

## as the root user run postinstall
    $ sudo su -
    # /opt/traffic_ops/install/bin/postinstall
