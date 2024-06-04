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

set -o nounset
old_log_dir=/opt/traffic_router/var/log
new_log_dir=/var/log/traffic_router
if [[ -d "$old_log_dir" ]]; then
	if [[ -d "$new_log_dir" ]]; then
		(
		# Include files starting with . in the * glob
		shopt -s dotglob
		mv "$old_log_dir"/* "$new_log_dir" || true
		)
		rmdir "$old_log_dir"
	else
		mv "$old_log_dir" "$new_log_dir"
	fi
	sync
fi

# figure out which version of traffic_router is currently running
# and then shut it down. Running both test just in case.
set +e

# delete the expanded war files from the previous version
if [[ -e /opt/traffic_router/webapps/core ]]; then
  echo "Deleting previous version of Traffic Router webapp"
  rm -rf /opt/traffic_router/webapps/core
fi

rm -rf /opt/traffic_router/webapps/*

