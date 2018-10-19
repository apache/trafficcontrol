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
#

# figure out which version of traffic_router is currently running
# and then shut it down
set +e
chkconfig --list tomcat >/dev/null

if [ $? -eq 0 ]; then
  /sbin/service tomcat stop
else
  /usr/bin/systemctl list-unit-files traffic_router.service > /dev/null

  [ $? -eq 0 ] && /usr/bin/systemctl stop traffic_router
fi

# delete the expanded war files from the previous version
if [[ -e /opt/traffic_router/webapps/core ]]; then
  echo "Deleting previous version of Traffic Router webapp"
  rm -rf /opt/traffic_router/webapps/core
fi

rm -rf /opt/traffic_router/webapps/*

