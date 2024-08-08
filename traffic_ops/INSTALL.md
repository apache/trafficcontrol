<!--
    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.
-->

# From new CentOS 7 install:

## Disable selinux:
### change `/etc/selinux/config` to `SELINUX=disabled`

## Add access to PostgreSQL yum repository 

Instructions are here: https://yum.postgresql.org/

- From this page,  copy the link for CentOS 7 and install:

    `$ sudo yum install https://download.postgresql.org/pub/repos/yum/reporpms/EL-7-x86_64/pgdg-redhat-repo-latest.noarch.rpm`
    
## Install Postgres 13.2 server (in a container or on the host)

### on the host:

    $ sudo su -
    # yum install postgresql13-server
    $ su - postgres
    $ /usr/pgsql-13/bin/initdb -A md5 -W #-W forces the user to provide a superuser (postgres) password
    $ exit
    # systemctl start postgresql-13
    # systemctl status postgresql-13

### -or- in a container

NOTE: you do *not* need postgresql13-server if running postgres within a `docker` container.

Install `docker` and `docker compose` using instructions here:

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
