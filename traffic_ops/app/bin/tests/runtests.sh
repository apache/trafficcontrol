#!/usr/bin/env bash
#
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# set up db based on given conf file

set -x
export PERL5LIB=/opt/traffic_ops/app/lib:/opt/traffic_ops/app/local/lib/perl5

usage() {
        echo "Usage: $(basename $0) <test dir> <test env> <host> <port>"
        echo "  e.g. $(basename $0) ./t test db 5432"
}

finish() {
        local st=$?
        [[ $st -ne 0 ]] && echo "Exiting with status $st"
        [[ -n $msg ]] && echo $msg
}

trap finish EXIT

dbconf=/opt/traffic_ops/app/conf/$TESTENV/database.conf
if [[ ! -f $dbconf ]]; then
        usage
        msg="$dbconf should be a file"
        exit 1
fi

if [[ ! -d $TESTDIR ]]; then
        usage
        msg="$TESTDIR should be a directory"
        exit 1
fi

while ! nc $DBHOST $DBPORT </dev/null; do # &>/dev/null; do
        echo "waiting for $DBHOST:$DBPORT"
        sleep 3
done

# get dbname, user, password from database.conf;  update it with hostname of db container
read dbname user pw <<<$(python -c "import json; d=json.load(open('$dbconf')); d['hostname']='$DBHOST'; d['port']='$DBPORT'; json.dump(d, open('$dbconf', 'w')); print d['dbname'],d['user'],d['password']") || exit $?

# update dbconf.yml
sed -E -i "s/host=[^ ]+/host=$DBHOST/" /opt/traffic_ops/app/db/dbconf.yml

# create user if doesn't exist
st=$(psql -h$DBHOST -p$DBPORT -Upostgres -tAc "SELECT 1 from pg_roles WHERE rolname='$user'") || exit $?
if [[ $st != 1 ]]; then
        psql -h$DBHOST -Upostgres -etAc "CREATE USER $user with LOGIN ENCRYPTED PASSWORD '$pw'" || exit $?
fi


st=$(psql -h$DBHOST -p$DBPORT -Upostgres -tAc "SELECT 1 FROM pg_database WHERE datname='$dbname'") || exit $?
if [[ $st != 1 ]]; then
        createdb -h$DBHOST -p$DBPORT -Upostgres -e --owner $user $dbname || exit $?
fi

cd /opt/traffic_ops/app
export USER=root

echo "/opt/traffic_ops/app/db/dbconf.yml"
cat /opt/traffic_ops/app/db/dbconf.yml

echo "/opt/traffic_ops/app/conf/$TESTENV/database.conf"
cat "/opt/traffic_ops/app/conf/$TESTENV/database.conf"

export GOROOT=/usr/local/go
export GOPATH=/opt/traffic_ops/go
PATH=$PATH:$GOPATH/bin:$GOROOT/bin

export PGOPTIONS='--client-min-messages=warning'
./db/admin.pl --env=$TESTENV reset

prove -qrp $TESTDIR
