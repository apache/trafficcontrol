#!/usr/bin/env bash

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

TEST_DIR=$1

if [[ -z  $TEST_DIR ]]; then
	echo "Usage: ./run-docker-integration-test.sh integration/test/directory"
	echo "Example: ./run-docker-integration-test.sh ./docker-integration-tests"
	exit 2
fi

TEST_FILES=
for filename in ${TEST_DIR}/*; do
	if [[ -x "$filename" ]]; then
		TEST_FILES="$TEST_FILES:$filename"
	fi
done

EXIT_CODE=0

ORIG_IFS=$IFS
IFS=:
for filepath in $TEST_FILES; do
		if [[ ! -z  $filepath ]]; then
				filename=$(basename $filepath)
				docker cp $filepath traffic_ops:/
				docker exec -it traffic_ops "/$filename"
				code=$?
				if [ "$code" -ne 0 ]; then
						EXIT_CODE=1
				fi
				docker exec -it traffic_ops rm "/$filename"
		fi
done
IFS=$ORIG_IFS

exit $EXIT_CODE
