#!/bin/bash
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
#
# ------------------------------------------------------

echo "Dumping Data from Traffic Ops Instance: $TO_SERVER" 
output=/tmp/trafficops_init.sql
[[ -n $output ]] && output="-o $output"

cookiejar=/tmp/cookiejar
cred=/tmp/cred.json

echo "mig:MYSQL_HOST: $MYSQL_HOST"
echo "mig:MYSQL_USER: $MYSQL_USER"
echo "mig:MYSQL_PASSWORD: $MYSQL_PASSWORD"
echo "mig:MYSQL_DATABASE: $MYSQL_DATABASE"

echo "mig:POSTGRES_HOST: $POSTGRES_HOST"
echo "mig:POSTGRES_USER: $POSTGRES_USER"
echo "mig:POSTGRES_DATABASE: $POSTGRES_DATABASE"
echo "mig:POSTGRES_PASSWORD: $POSTGRES_PASSWORD"

cat >$cred <<-CREDS
	{ "u" : "$TO_USER", "p" : "$TO_PASSWORD" }
CREDS

curl -f -k -H "Accept: application/json" --cookie "$cookiejar" --cookie-jar "$cookiejar" -X POST --data @"$cred" "$TO_SERVER/api/1.2/user/login"  || exit 1
curl $output -f -k -s --cookie "$cookiejar" -X GET "$TO_SERVER/dbdump"  || exit 1

echo  "[client]" > /root/.my.cnf
echo  "user=$MYSQL_USER" >> /root/.my.cnf 
echo  "password=$MYSQL_PASSWORD" >> /root/.my.cnf 
chmod 0600 /root/.my.cnf
mysql -h $MYSQL_HOST $MYSQL_DATABASE < /tmp/trafficops_init.sql

pgloader -v \
	--cast 'type tinyint to smallint drop typemod' \
	--cast 'type varchar to text drop typemod' \
	--cast 'type double to numeric drop typemod' \
	mysql://$MYSQL_USER:$MYSQL_PASSWORD@$MYSQL_HOST/$MYSQL_DATABASE \
	postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST/$POSTGRES_DATABASE

# For debugging
#while true; do
#    echo "Waiting.."
#    sleep 3
#done

