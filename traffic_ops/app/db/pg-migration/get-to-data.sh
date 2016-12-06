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
#!/bin/bash -x

output=$1
[[ -n $output ]] && output="-o $output"


cookiejar=/tmp/cookiejar
cred=/tmp/cred.json

cat >$cred <<-CREDS
	{ "u" : "$TO_USER", "p" : "$TO_PASSWORD" }
CREDS

curl -k -H "Accept: application/json" --cookie "$cookiejar" --cookie-jar "$cookiejar" -X POST --data @"$cred" "$TO_SERVER/api/1.2/user/login"
curl $output -k -s --cookie "$cookiejar" -X GET "$TO_SERVER/dbdump"
