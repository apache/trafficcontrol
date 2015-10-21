#!/bin/bash

#
# Copyright 2015 Comcast Cable Communications Management, LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# By default all sub-projects are built.  Supply a list of projects to build if
# only a subset is wanted.
all_projects="\
	traffic_ops \
	traffic_ops_ort \
	traffic_monitor \
	traffic_router \
	traffic_stats \
"

if [[ $# -gt 0 ]]; then
	projects="$*"
else
	projects=$all_projects
fi

for p in $projects; do
	echo "-----  Building $p ..."
	case $p in
		traffic_ops)     (cd traffic_ops/rpm && ./build_rpm.sh) ;;
		traffic_ops_ort) (cd traffic_ops/rpm && ./build_ort_rpm.sh) ;;
		traffic_monitor) (cd traffic_monitor/rpm && ./build_rpm.sh) ;;
		traffic_router)  (cd traffic_router/rpm && ./build_rpm.sh) ;;
		traffic_stats)   (cd traffic_stats/rpm && ./build_rpm.sh) ;;
		*) echo "No project named $p"; exit 1;;
	esac || (echo "$p failed: $!"; exit 1)
done
