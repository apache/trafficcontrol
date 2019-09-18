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

class FilterModule(object):
    def filters(self):
        return {
            'associate_round_robin': self.associate_round_robin,
            'denormalize_association': self.denormalize_association
        }
    def associate_round_robin(self, set1, set2):
        step_size = len(set1)
        out_list_normal = {}
        for i, val in enumerate(set1):
            # Offset the beginning of the second set by the index of where we are processing the first set
            sublist = set2[i:]
            # Snag every Nth element from the offset sublist and associte that with set1
            out_list_normal[val] = sublist[::step_size]
        return out_list_normal

    def denormalize_association(self, association):
        out = {}
        for key, vals in association.items():
            for val in vals:
                out[val] = key
        return out
