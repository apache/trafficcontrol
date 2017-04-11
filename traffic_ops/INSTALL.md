# From new CentOS 7 install:

## Disable selinux:
### change `/etc/selinux/config` to `SELINUX=disabled`

## Install Postgreql 9.6 client and libraries

Instructions are here: https://yum.postgresql.org/

- grab the link for CentOS 7 and install:

    $ sudo yum install https://download.postgresql.org/pub/repos/yum/9.6/redhat/rhel-7-x86_64/pgdg-centos96-9.6-xxxx.noarch.rpm
    
  NOTE: get a valid link from https://yum.postgresql.org/ with the correct version number.
  
- install `postgresql96` (for psql commands) and `postgresql96-devel` (for includes/libraries needed to install `DBD::Pg` perl library)

    $ sudo yum install postgresql96 postgresql96-devel

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

Add yourself to `docker` group

    $ sudo usermod -G docker $USER

Remember to logout and login again..   You should see `docker` in your list of groups:

    $ id
    uid=9876(myuser) gid=9876(myuser) groups=9876(myuser),990(docker) ...

Edit `mysql-to-postgres.env` to suit your needs.
* modify `POSTGRES_*` vars to apply to new postgres container that will house your database, e.g.
  * POSTGRES_USER=postgres
  * POSTGRES_PASSWORD=itSas3cre4
  
If migrating from an existing pre-2.0 traffic_ops server (mysql):
* `TO_*` vars for admin access to existing mysql-based `traffic_ops` (to get a db dump), e.g.
  * TO_SERVER=https://trafficops.example.com
  * TO_USER=dennisr
  
* `MYSQL_*` vars to apply to temporary mysql container -- really no need to change..

Start a docker container to run postgres

    $ cd incubator-trafficcontrol/traffic_ops/app/db/pg-migration
    $ ./start_postgres.sh

Run migration from existing mysql-based `traffic_ops`

    $ ./migrate.sh

## Install `traffic_ops`

    $ sudo yum install traffic_ops

## Install `go` and `git` (required for `goose` and some `Perl` modules)

    $ sudo yum install git go
    
## Install Perl modules

    $ sudo cpanm Carton

IMPORTANT!!: We're using a later version of Postgresql,  so it's not installed in the default place.
We need to tell carton where it is so the `DBD::Pg` module is installed correctly.

    $ sudo su -
    # cd /opt/traffic_ops/app
    # POSTGRES_HOME=/usr/pgsql-9.6 /usr/local/bin/carton


## Install goose

    $ sudo GOPATH=/usr/local go get bitbucket.org/liamstask/goose/cmd/goose


## Install any extensions needed

   - install in /opt/traffic_ops_extensions
   
## Install `openssl` certs (or use this to generate them)

   $ sudo /opt/traffic_ops/install/bin/generateCert

## as the root user run postinstall
    $ sudo su -
    # export POSTGRES_HOME=/usr/pgsql-9.6
    # export GOPATH=/usr/local
    # /opt/traffic_ops/install/bin/postinstall
