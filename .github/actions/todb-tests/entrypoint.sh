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

set -ex

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
for file in "$(ls | uniq -d)"; do
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

# All new migrations must use 16-digit timestamps.
VIOLATING_FILES="$(ls | cut -d _ -f 1 | sed -n -e '/2020061622101648/,$p' | grep -vE '^\d{16}$')";

if [[ ! -z "$VIOLATING_FILES" ]]; then
	for file in "$VIOLATING_FILES"; do
		echo "ERROR: traffic_ops/app/db/migrations/$file(name).sql: wrong filename, all migrations after 2020-06-16 must use 16-digit timestamps in their filenames" >&2;
	done
	CODE=1;
fi

exit "$CODE";
