#!/bin/bash -x
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

# make sure postgres and mysql ports are both active
echo "POSTGRES_HOST=$POSTGRES_HOST MYSQL_HOST=$MYSQL_HOST"
for c in "$POSTGRES_HOST 5432" "$MYSQL_HOST 3306"; do
	while true; do
		echo Waiting for $c
		sleep 3
		nc -z $c && break
	done
done

pgloader -v \
	--cast 'type tinyint to smallint drop typemod' \
	--cast 'type varchar to text drop typemod' \
	--cast 'type double to numeric drop typemod' \
	mysql://$MYSQL_USER:$MYSQL_PASSWORD@$MYSQL_HOST/traffic_ops_db \
	postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST/$POSTGRES_DB
exit 0
