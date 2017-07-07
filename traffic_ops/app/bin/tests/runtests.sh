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
testdir=$1
dbenv=$2
dbhost=$3
dbport=${4:-5432}
export PERL5LIB=/opt/traffic_ops/app/lib:/opt/traffic_ops/app/local/lib/perl5

usage() {
        echo "Usage: $(basename $0) <test dir> <dbenv> <host> <port>"
        echo "  e.g. $(basename $0) ./t test db 5432"
}

finish() {
        local st=$?
        [[ $st -ne 0 ]] && echo "Exiting with status $st"
        [[ -n $msg ]] && echo $msg
}

trap finish EXIT

dbconf=/opt/traffic_ops/app/conf/$dbenv/database.conf
if [[ ! -f $dbconf ]]; then
        usage
        msg="$dbconf should be a file"
        exit 1
fi

if [[ ! -d $testdir ]]; then
        usage
        msg="$testdir should be a directory"
        exit 1
fi

while ! nc $dbhost $dbport </dev/null; do # &>/dev/null; do
        dig $dbhost +short
        echo "waiting for $dbhost:$dbport"
        sleep 3
done

# get dbname, user, password from database.conf;  update it with hostname of db container
read dbname user pw <<<$(python -c "import json; d=json.load(open('$dbconf')); d['hostname']='$dbhost'; d['port']='$dbport'; json.dump(d, open('$dbconf', 'w')); print d['dbname'],d['user'],d['password']") || exit $?

# update dbconf.yml
sed -E -i "s/host=[^ ]+/host=$dbhost/" /opt/traffic_ops/app/db/dbconf.yml

# create user if doesn't exist
st=$(psql -h$dbhost -p$dbport -Upostgres -tAc "SELECT 1 from pg_roles WHERE rolname='$user'") || exit $?
if [[ $st != 1 ]]; then
        psql -h$dbhost -Upostgres -etAc "CREATE USER $user with LOGIN ENCRYPTED PASSWORD '$pw'" || exit $?
fi


st=$(psql -h$dbhost -p$dbport -Upostgres -tAc "SELECT 1 FROM pg_database WHERE datname='$dbname'") || exit $?
if [[ $st != 1 ]]; then
        createdb -h$dbhost -p$dbport -Upostgres -e --owner $user $dbname || exit $?
fi

cd /opt/traffic_ops/app
export USER=root

echo "/opt/traffic_ops/app/db/dbconf.yml"
cat /opt/traffic_ops/app/db/dbconf.yml

echo "/opt/traffic_ops/app/conf/$dbenv/database.conf"
cat "/opt/traffic_ops/app/conf/$dbenv/database.conf"

export GOROOT=/usr/local/go
export GOPATH=/opt/traffic_ops/go
PATH=$PATH:$GOPATH/bin:$GOROOT/bin

export PGOPTIONS='--client-min-messages=warning'
./db/admin.pl --env=test reset

prove -qrp $testdir
