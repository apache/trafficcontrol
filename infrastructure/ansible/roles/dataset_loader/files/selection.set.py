#!/usr/bin/python

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

# This script is used to provide a round-robin merging of two lists

import sys
import json

if len(sys.argv) < 3 or len(sys.argv) > 4:
	print("{}")
	sys.exit(0)

cdn_csv_list = sys.argv[1].split(',')
fqdn_csv_list = sys.argv[2].split(',')
option = ''
if len(sys.argv) == 4:
	option = sys.argv[3]
cdn_csv_list.sort()
fqdn_csv_list.sort()

step_size = len(cdn_csv_list)
out_list_normal = {}
for i, val in enumerate(cdn_csv_list):
	sublist = fqdn_csv_list[i:]
	out_list_normal[val] = ','.join(sublist[::step_size])

out_list_denormal = {}
for val, csvlist in out_list_normal.items():
	for i in csvlist.split(','):
		if i != "":
			out_list_denormal[i] = val

if option == 'denormalize':
	print(json.dumps(out_list_denormal))
else:
	print(json.dumps(out_list_normal))
