#!/bin/bash
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

set -o errexit -o nounset
trap 'echo "Error on line ${LINENO} of ${0}"; exit 1' ERR

load() {
	if [[ -e "$archive_filename" ]]; then
		docker image load -i "$archive_filename"
	else
		echo "No tarred image found named ${archive_filename}"
		docker pull "alpine@${alpine_digest}"
	fi
}

save() {
	if [[ -e "$archive_filename" ]]; then
		echo "Docker image archive ${archive_filename} already exists. Skipping save..."
		return
	fi
	mkdir -p docker-images
	docker image save "alpine@${alpine_digest}" -o "$archive_filename"
	echo "Saved tarred image ${archive_filename}"
}

if [[ $# -ge 1 ]]; then
	action="$1"
	shift
else
	echo 'Argument `load-or-save` is required but was not found.'
	exit 1
fi
if [[ $# -ge 1 ]]; then
	alpine_digest="$1"
	shift
else
	echo 'Input `digest` is required but was not found.'
	exit 1
fi
archive_filename="docker-images/alpine@${alpine_digest}.tar.gz"

if [[ "$action" != load && "$action" != save ]]; then
	export action
	<<-'MESSAGE' envsubst
	Invalid value `${action}` was found for input load-or-save. Valid values:
	* load
	* save
MESSAGE
	exit 1
fi

"$action"
