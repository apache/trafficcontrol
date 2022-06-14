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

files_changed_in_pr() {
	if [[ "$GITHUB_REF" == refs/pull/*/merge ]]; then
		pr_number="$(<<<"$GITHUB_REF" grep -o '[0-9]\+')"
		files_changed="$(curl "${GITHUB_API_URL}/repos/${GITHUB_REPOSITORY}/pulls/${pr_number}/files" | jq -r .[].filename)"
	else
		files_changed="$(git diff-tree --no-commit-id --name-only -r "$GITHUB_SHA")"
	fi
}

CODE=0;

migration_dirs=(traffic_ops/app/db/migrations traffic_ops/app/db/trafficvault/migrations);

for migration_dir in ${migration_dirs[@]}; do
	if ! [[ -d "$migration_dir" ]]; then
		echo "No migrations exist at $migration_dir - skipping!" >&2;
		continue;
	fi
	pushd $migration_dir;

	# Ensure proper order
	SORTED="$(mktemp)";
	SORTEDDASHN="$(mktemp)";

	ls | sort > "$SORTED";
	ls | sort -n > "$SORTEDDASHN";

	if [ ! -z "$(diff $SORTED $SORTEDDASHN)" ]; then
		echo "ERROR: expected sort -n and sort to give the same migration order:" >&2;
		diff "$SORTED" "$SORTEDDASHN" >&2;
		CODE=1;
	fi

	rm "$SORTED" "$SORTEDDASHN";

	# No two migrations may share a timestamp
	for direction in up down; do
		while read -r file; do
			echo "ERROR: more than one file uses timestamp $file - timestamps must be unique" >&2;
			CODE=1;
		done < <(ls *".${direction}.sql" | cut -d_ -f1 | uniq -d)
	done

	# Files must be named like {{timestamp}}_{{migration name}}.up.sql
	pattern='^[0-9]+_[^.]+\.(up|down)\.sql$'
	for file in *; do
		if ! <<<"$file" grep -qE "$pattern"; then
			echo "ERROR: ${migration_dir}/${file}: wrong filename, must match ${pattern}" >&2;
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

	if [[ $LATEST_FILE_TIME != ${mtime_array[$mtime_length-1]} ]] && <<<"$(files_changed_in_pr)" grep -q "^${LATEST_FILE}$"; then
		echo "ERROR: latest added/modified file: $LATEST_FILE is not in the right order" >&2;
		CODE=1;
	fi

	set +e;
	# All new migrations must use 16-digit timestamps.
	VIOLATING_FILES="$(ls | sort | cut -d _ -f 1 | grep -vE '^[0-9]{16}$' | grep -vE '^00000000000000$')";
	set -e;

	if [[ ! -z "$VIOLATING_FILES" ]]; then
		for file in "$VIOLATING_FILES"; do
			echo "ERROR: $migration_dir/$file(name).sql: wrong filename, all migrations after 2020-06-16 must use 16-digit timestamps in their filenames" >&2;
		done
		CODE=1;
	fi
	popd;
done

exit "$CODE";
