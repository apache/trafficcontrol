#
# Copyright 2011-2014, Comcast Corporation. This software and its contents are
# Comcast confidential and proprietary. It cannot be used, disclosed, or
# distributed without Comcast's prior written permission. Modification of this
# software is only allowed at the direction of Comcast Corporation. All allowed
# modifications must be provided to Comcast Corporation.
#
#!/bin/sh

BCK_FILE="/opt/tmredis/backup/tm_redis_`date +%m%d20%y`.json"; export BCK_FILE
LOG_FILE="/opt/tmredis/var/log/tmredis/backup_redis_daily.out"; export LOG_FILE
/opt/tmredis/bin/backup_redis_daily -file=$BCK_FILE -redis=localhost:6379 > $LOG_FILE 2>&1 

if [ -f $BCK_FILE ]; then
	/bin/gzip $BCK_FILE
else
	echo "ERROR: $BCK_FILE does not exist"
fi
