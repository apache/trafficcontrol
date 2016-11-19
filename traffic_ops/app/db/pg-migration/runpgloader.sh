#!/bin/bash -x

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
