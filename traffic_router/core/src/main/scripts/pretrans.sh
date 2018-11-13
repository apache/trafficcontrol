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
# and then shut it down. Running both test just in case.
set +e

if [[ -e "/etc/init.d/tomcat" ]]; then
  echo "Stopping tomcat service..."
  /sbin/service tomcat stop
  chkconfig tomcat off
fi

echo "Stopping traffic_router services"
/usr/bin/systemctl list-unit-files traffic_router.service > /dev/null
[ $? -eq 0 ] && /usr/bin/systemctl stop traffic_router
echo "Done stopping traffic_router services"

