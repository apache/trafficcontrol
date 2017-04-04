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

   $ sudo yum install postgresql96-server

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

## Install `go` and `gcc` (required for `goose` and some `Perl` modules)

    $ sudo yum install go gcc
    
## Install Perl modules

    $ sudo cpanm Carton

IMPORTANT!!: We're using a later version of Postgresql,  so it's not installed in the default place.
We need to tell carton where it is so the `DBD::Pg` module is installed correctly.

    $ sudo su -
    # cd /opt/traffic_ops/app
    # POSTGRES_HOME=/usr/pgsql-9.6 /usr/local/bin/carton


## Install goose

    $ sudo GOPATH=/tmp GOBIN=/usr/local/bin go get bitbucket.org/liamstask/goose/cmd/goose


## Modify `traffic_ops` configuration

- `/opt/traffic_ops/app/db`
   - `dbconf.yml` 
      - modify "production" line to match user/pass from env file above
- `/opt/traffic_ops/app/conf`
   - `cdn.conf` 
      - set workers to desired value (96 is far too high for dev environment -- 15 is suggested)
      - change `to.base_url` to appropriate FQDN or IP address
   - `ldap.conf`
      - add ldap server credentials if needed
   - `production/database.conf`
      - modify to match user/pass from env file above
   - `production/riak.conf`, `production/influxdb.conf`
      - add appropriate user/password
   - `production/log4perl.conf`
      - if logging data needed,  change ERROR to DEBUG on first line

## Initialize the db

    $ cd /opt/traffic_ops/app
    $ PERL5LIB=$(pwd)/lib:$(pwd)/local/lib/perl5 db/admin.pl --env=production setup
    
## Install any extensions needed

   - install in /opt/traffic_ops_extensions
   
## Install `openssl` certs (or use this to generate them)

   - `sudo /opt/traffic_ops/install/bin/generateCert`
   
## Install web dependencies

   - `sudo /opt/traffic_ops/install/bin/download_web_deps`
