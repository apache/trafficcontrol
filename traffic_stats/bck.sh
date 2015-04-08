#
#
#!/bin/sh
#
# Copyright 2015 Comcast Cable Communications Management, LLC
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
#

BCK_FILE="/opt/tmredis/backup/tm_redis_`date +%m%d20%y`.json"; export BCK_FILE
LOG_FILE="/opt/tmredis/var/log/tmredis/backup_redis_daily.out"; export LOG_FILE
/opt/tmredis/bin/backup_redis_daily -file=$BCK_FILE -redis=localhost:6379 > $LOG_FILE 2>&1 

if [ -f $BCK_FILE ]; then
	/bin/gzip $BCK_FILE
else
	echo "ERROR: $BCK_FILE does not exist"
fi
