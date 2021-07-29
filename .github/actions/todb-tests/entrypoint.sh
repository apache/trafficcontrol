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

set -e

cd traffic_ops/app/db/migrations;

# Ensure proper order
SORTED="$(mktemp)";
SORTEDDASHN="$(mktemp)";

ls | sort > "$SORTED";
ls | sort -n > "$SORTEDDASHN";

CODE=0;

if [ ! -z "$(diff $SORTED $SORTEDDASHN)" ]; then
	echo "ERROR: expected sort -n and sort to give the same migration order:" >&2;
	diff "$SORTED" "$SORTEDDASHN" >&2;
	CODE=1;
fi

rm "$SORTED" "$SORTEDDASHN";

# No two migrations may share a timestamp
ls | cut -d _ -f 1 | uniq -d | while read -r file; do
	echo "ERROR: more than one file uses timestamp $file - timestamps must be unique" >&2;
	CODE=1;
done

# Files must be named like {{timestamp}}_{{migration name}}.sql
for file in "$(ls)"; do
	if ! [[ "$file" =~ [0-9]+_[^\.]+\.sql ]]; then
		echo "ERROR: traffic_ops/app/db/migrations/$file: wrong filename, must match \d+_[^\\.]+\\.sql" >&2;
		CODE=1;
	fi
done

# Files added must have date and name later than all existing file
LATEST_FILE="$(git log -1 --name-status --diff-filter=d --format="%ct" . | tail -n 1 | awk '{print $2}' | cut -d / -f5)"
LATEST_FILE_TIME="$(git log -1 --name-status --diff-filter=d --format="%ct" . | head -n 1 )"

# Get modified times in an array
mtime_array=()
arr=($(ls))
for file in "${arr[@]}"; do
  mtime_array+=( "$(git log -1 --format=%ct  $file)" )
done
mtime_length=${#mtime_array[@]}

if [[ $LATEST_FILE_TIME != ${mtime_array[$mtime_length-1]} ]]; then
  echo "ERROR: latest added/modified file: $LATEST_FILE is not in the right order" >&2;
  CODE=1;
fi

set +e;
# All new migrations must use 16-digit timestamps.
VIOLATING_FILES="$(ls | sort | cut -d _ -f 1 | sed -n -e '/2020061622101648/,$p' | tr '[:space:]' '\n' | grep -vE '^[0-9]{16}$')";
set -e;

if [[ ! -z "$VIOLATING_FILES" ]]; then
	for file in "$VIOLATING_FILES"; do
		echo "ERROR: traffic_ops/app/db/migrations/$file(name).sql: wrong filename, all migrations after 2020-06-16 must use 16-digit timestamps in their filenames" >&2;
	done
	CODE=1;
fi

exit "$CODE";
