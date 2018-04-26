#!/usr/bin/env python
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

"""
Generate an ansible inventory compatable dataset from TO

Before using, make sure you have the TrafficOps Python Native Client library available inside your env
# Usage: python TO.py -to <TO_USER> <TO_PASS> -url to.kabletown.invalid --list
export TO_USERNAME=<TO_USERname>
export TO_PASSWORD=<TO_PASSword>
export TO_URL=<TO_URL>
"""

import json
import argparse
import logging
import os
import collections
from trafficops.tosession import TOSession

# Disable rest api logging, allowing invalid certs causes warnings
logging.getLogger('common.restapi').disabled = True
#logging.getLogger('common.restapi').logging_level = logging.DEBUG

def empty_inventory():
    """Generate a valid empty inventory"""
    return {'_meta': {'hostvars': {}}}


class AnsibleInventory(object):
    """Wrapper class for needed methods"""

    def __init__(self, user, password, url):
        """Init base members"""
        self.to_user = user
        self.to_pass = password
        self.to_url = url

    @classmethod
    def populate_server_profile_vars(cls, api, profile_id):
        """Generate the server profile variables once as we see it"""
        server_vars = {}
        server_vars['hosts'] = []
        server_vars['vars'] = {}
        profile = api.get_profile_by_id(
            profile_id=profile_id)[0]
        server_vars['vars']['server_profile_description'] = profile[0]['description']
        server_vars['vars']['server_profile_type'] = profile[0]['type']
        server_vars['vars']['server_profile_routingDisabled'] = profile[0]['routingDisabled']
        server_vars['vars']['server_profile_parameters'] = []
        params = api.get_parameters_by_profile_id(
            id=profile_id)[0]
        for param in params:
            tmp_param = {
                'name': param['name'],
                'value': param['value'],
                'configFile': param['configFile']}
            server_vars['vars']['server_profile_parameters'].append(tmp_param)
        return server_vars

    @classmethod
    def populate_cachegroups(cls, api, cachegroup_id):
        """Generate the values for cachegroups once on first sight"""
        var_data = {}
        cgdata = collections.namedtuple('Cgdata', ['cgvars',
                                                   'primary_parent_group_name',
                                                   'secondary_parent_group_name'])
        var_data['hosts'] = []
        var_data['vars'] = {}
        cachegroup = api.get_cachegroup_by_id(
            cache_group_id=cachegroup_id)[0]
        var_data['vars']['cachegroup_name'] = cachegroup[0]['name']
        var_data['vars']['cachegroup_shortName'] = cachegroup[0]['shortName']
        var_data['vars']['cachegroup_parentCachegroupName'] = \
            cachegroup[0]['parentCachegroupName']
        var_data['vars']['cachegroup_secondaryParentCachegroupName'] = \
            cachegroup[0]['secondaryParentCachegroupName']
        var_data['vars']['cachegroup_typeName'] = cachegroup[0]['typeName']
        if cachegroup[0]['parentCachegroupName'] is None:
            flat_parent_cg = "parentCachegroup|None"
        else:
            flat_parent_cg = "parentCachegroup|" + \
                cachegroup[0]['parentCachegroupName']

        if cachegroup[0]['secondaryParentCachegroupName'] is None:
            flat_second_parent_cg = "secondaryParentCachegroup|None"
        else:
            flat_second_parent_cg = "secondaryParentCachegroup|" + \
                cachegroup[0]['secondaryParentCachegroupName']
        out = cgdata(cgvars=var_data,
                     primary_parent_group_name=flat_parent_cg,
                     secondary_parent_group_name=flat_second_parent_cg)
        return out

    def generate_inventory_list(self, target_environment):  # pylint: disable=too-many-statements
        """Generate the inventory list for the specified environment"""
        traffic_ops_api = TOSession(
            self.to_url, verify_cert=False)
        traffic_ops_api.login(self.to_user, self.to_pass)
        servers = traffic_ops_api.get_servers()[0]
        out = {}
        out['_meta'] = {}
        out['_meta']['hostvars'] = {}
        out[target_environment] = {}
        out[target_environment]['hosts'] = []
        out["ungrouped"] = {}
        out['ungrouped']['hosts'] = []
        out['cachegroup'] = {}
        out['cachegroup']['children'] = []
        out['server_type'] = {}
        out['server_type']['children'] = []
        out['server_cdnName'] = {}
        out['server_cdnName']['children'] = []
        out['server_profile'] = {}
        out['server_profile']['children'] = []
        out['server_status'] = {}
        out['server_status']['children'] = []
        for server in servers:
            fqdn = server['hostName'] + '.' + server['domainName']
            out["ungrouped"]['hosts'].append(fqdn)
            out[target_environment]['hosts'].append(fqdn)
            out['_meta']['hostvars'][fqdn] = {}
            out['_meta']['hostvars'][fqdn]['server_toEnvironment'] = target_environment
            out['_meta']['hostvars'][fqdn]['server_cachegroup'] = server['cachegroup']
            out['_meta']['hostvars'][fqdn]['server_cdnName'] = server['cdnName']
            out['_meta']['hostvars'][fqdn]['server_id'] = server['id']
            out['_meta']['hostvars'][fqdn]['server_ipAddress'] = server['ipAddress']
            out['_meta']['hostvars'][fqdn]['server_ip6Address'] = server['ip6Address']
            out['_meta']['hostvars'][fqdn]['server_offlineReason'] = server['offlineReason']
            out['_meta']['hostvars'][fqdn]['server_physLocation'] = server['physLocation']
            out['_meta']['hostvars'][fqdn]['server_profile'] = server['profile']
            out['_meta']['hostvars'][fqdn]['server_profileDesc'] = server['profileDesc']
            out['_meta']['hostvars'][fqdn]['server_status'] = server['status']
            out['_meta']['hostvars'][fqdn]['server_type'] = server['type']
            flat_server_profile = "server_profile|" + server['profile']
            flat_cachegroup = "cachegroup|" + server['cachegroup']
            flat_server_type = "server_type|" + server['type']
            flat_server_cdn_name = "server_cdnName|" + server['cdnName']
            flat_server_status = "server_status|" + server['status']
            if flat_server_profile not in out.keys():
                out['server_profile']['children'].append(flat_server_profile)
                out[flat_server_profile] = self.populate_server_profile_vars(
                    traffic_ops_api,
                    server['profileId'])
            out[flat_server_profile]['hosts'].append(fqdn)
            if flat_cachegroup not in out.keys():
                out['cachegroup']['children'].append(flat_cachegroup)
                # out[flat_cachegroup] = self.populate_cachegroups(
                cgdata = self.populate_cachegroups(
                    traffic_ops_api,
                    server['cachegroupId'])
                out[flat_cachegroup] = cgdata.cgvars
                flat_parent_cg = cgdata.primary_parent_group_name
                flat_second_parent_cg = cgdata.secondary_parent_group_name
                if flat_parent_cg not in out.keys():
                    out[flat_parent_cg] = {}
                    out[flat_parent_cg]['children'] = []
                if flat_second_parent_cg not in out.keys():
                    out[flat_second_parent_cg] = {}
                    out[flat_second_parent_cg]['children'] = []
                out[flat_parent_cg]['children'].append(flat_cachegroup)
                out[flat_second_parent_cg]['children'].append(flat_cachegroup)
            out[flat_cachegroup]['hosts'].append(fqdn)
            if flat_server_type not in out.keys():
                out['server_type']['children'].append(flat_server_type)
                out[flat_server_type] = {}
                out[flat_server_type]['hosts'] = []
            out[flat_server_type]['hosts'].append(fqdn)
            if flat_server_cdn_name not in out.keys():
                out['server_cdnName']['children'].append(flat_server_cdn_name)
                out[flat_server_cdn_name] = {}
                out[flat_server_cdn_name]['hosts'] = []
            out[flat_server_cdn_name]['hosts'].append(fqdn)
            if flat_server_status not in out.keys():
                out['server_status']['children'].append(flat_server_status)
                out[flat_server_status] = {}
                out[flat_server_status]['hosts'] = []
            out[flat_server_status]['hosts'].append(fqdn)
        return out

    def to_inventory(self):
        return self.generate_inventory_list(self.to_url)


if __name__ == "__main__":
    PARSER = argparse.ArgumentParser(
        description='Generate ansible inventory from TrafficOps')

    PARSER.add_argument(
        '-to',
        nargs='+',
        help='Please pass username and password for TrafficOps')
    PARSER.add_argument(
        '-url',
        type=str,
        help='TO URL')
    PARSER.add_argument('--list', action='store_true')
    PARSER.add_argument('--host', type=str, help='Target TO Server')

    ARGS = PARSER.parse_args()

    TO_USER = os.getenv("TO_USERNAME")
    TO_PASS = os.getenv("TO_PASSWORD")
    TO_URL = os.getenv("TO_URL")

    if ARGS.to:
        if TO_USER is None:
            TO_USER = ARGS.to[0]

        if TO_PASS is None:
            TO_PASS = ARGS.to[1]

    if TO_URL is None:
        if ARGS.env is None:
            TO_URL = "https://localhost:8080"
        else:
            TO_URL = ARGS.env

    if TO_USER and TO_PASS is not None:
        if ARGS.list:
            TMP_ANSIBLEINVENTORY = AnsibleInventory(TO_USER, TO_PASS, TO_URL)
            INVENTORY = TMP_ANSIBLEINVENTORY.to_inventory()
        # Since we're supplying hostvar metadata, --host support isn't required
        else:
            INVENTORY = empty_inventory()
    else:
        INVENTORY = empty_inventory()

    print json.dumps(INVENTORY)
