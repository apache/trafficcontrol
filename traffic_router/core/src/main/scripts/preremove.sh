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

# stop traffic router depending on which OS we are on CentOS6 or CentOS7
$VERSION_ID = %{centos_ver}

if [[ -n "$VERSION_ID" && "$VERSION_ID" == "7" ]]; then
	/usr/bin/sudo /usr/bin/systemctl stop traffic_router
else
	/sbin/service traffic_router stop
fi