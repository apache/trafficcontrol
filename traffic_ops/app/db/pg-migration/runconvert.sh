#!/bin/bash -x
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
#

set -x

waiting=/sync/waiting-for-pgloader
touch $waiting

# Wait for pgloader to finish
while [[ -f $waiting ]]; do
    ls -l $waiting
    sleep 3
done

echo "Looks like pgloader is finished..  Converting.."

# Load required conversion of booleans
psql postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST/$POSTGRES_DB < ./convert_bools.sql
