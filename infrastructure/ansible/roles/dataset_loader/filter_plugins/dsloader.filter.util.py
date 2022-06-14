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
# Make coding more python3-ish
from __future__ import (absolute_import, division, print_function, unicode_literals)
__metaclass__ = type

import re

class FilterModule(object):
    def filters(self):
        return {
            'associate_round_robin': self.associate_round_robin,
            'denormalize_association': self.denormalize_association,
            'associate_profile_to_servers': self.associate_profile_to_servers
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
        for key, vals in list(association.items()):
            for val in vals:
                out[val] = key
        return out

    def associate_profile_to_servers(self, host_list, profile_list, hardware_profile_map, hostvars, profile_whitelist):
        out = {}
        hardware_profile_suffix_regex = re.compile(r'^(.*?)(?:_hdwr\d+)?$')
        compound_key_to_profile_map = {}
        compound_key_to_host_map = {}
        for cdnName in set(hostvars[host]['cdn'] for host in hostvars):
            for componentName in set(hostvars[host]['component'] for host in hostvars):
                compound_key = ""+(cdnName)+"_"+(componentName)
                eligible_profiles = []
                for profileName in profile_list:
                    common_profileName_re = hardware_profile_suffix_regex.search(profileName) # should always match unless something is really wrong
                    if re.match(profile_whitelist[componentName],profileName) and common_profileName_re is not None:
                        # Get the name of the profile without the hardware suffix so that we ignore that factor when assigning initial placements until later
                        eligible_profiles.append(common_profileName_re.group(1).lower())
                eligible_profiles = list(set(eligible_profiles)) # deduplicate the list
                # Bail out early if there aren't any profiles that could be assigned
                if len(eligible_profiles) > 0:
                    compound_key_to_profile_map[compound_key] = eligible_profiles
                    eligible_hosts = []
                    for host in host_list:
                        if hostvars[host]['cdn'] == cdnName and hostvars[host]['component'] == componentName:
                            eligible_hosts.append(host)
                    if len(eligible_hosts) > 0:
                        compound_key_to_host_map[compound_key] = eligible_hosts

        # Ignores combinations that either have no profiles or hosts
        for key in list(compound_key_to_host_map.keys()):
            host_to_profile_commonName_map = self.denormalize_association(self.associate_round_robin(compound_key_to_profile_map[key],compound_key_to_host_map[key]))
            for host, common_profileName in list(host_to_profile_commonName_map.items()):
                if host in list(hardware_profile_map.keys()):
                    # Re-apply the hardware profile variances
                    out[host] = common_profileName+"_"+hardware_profile_map[host]['hardware_profile']
                else:
                    out[host] = common_profileName

        return out
