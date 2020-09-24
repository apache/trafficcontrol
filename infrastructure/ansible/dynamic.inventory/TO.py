#!/usr/bin/env python3
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
"""A python script that converts TrafficOps information into an Ansible Inventory"""

import json
import argparse
import os
import collections
from trafficops.tosession import TOSession


def empty_inventory():
	"""Generate a valid empty inventory"""
	return {'_meta': {'hostvars': {}}}


class AnsibleInventory():
	"""Wrapper class for needed methods"""

	def __init__(self, user, password, url, verify_cert):
		"""Init base members"""
		self.to_user = user
		self.to_pass = password
		self.to_url = url
		self.verify_cert = verify_cert

	@staticmethod
	def populate_server_profile_vars(api, profile_id):
		"""Generate the server profile variables once as we see it"""
		server_vars = {}
		server_vars['hosts'] = []
		server_vars['vars'] = {}
		profile = api.get_profiles(id=profile_id)[0]
		server_vars['vars']['server_profile_description'] = profile[0]['description']
		server_vars['vars']['server_profile_type'] = profile[0]['type']
		server_vars['vars']['server_profile_routingDisabled'] = profile[0]['routingDisabled']
		server_vars['vars']['server_profile_parameters'] = []
		params = api.get_parameters_by_profile_id(id=profile_id)[0]
		for param in params:
			tmp_param = {
				'name': param['name'],
				'value': param['value'],
				'configFile': param['configFile']}
			server_vars['vars']['server_profile_parameters'].append(tmp_param)
		return server_vars

	@staticmethod
	def populate_cachegroups(api, cachegroup_id):
		"""Generate the values for cachegroups once on first sight"""
		var_data = {}
		cgdata = collections.namedtuple(
			'Cgdata', [
				'cgvars', 'primary_parent_group_name', 'secondary_parent_group_name'])
		var_data['hosts'] = []
		var_data['vars'] = {}
		cachegroup = api.get_cachegroups(id=cachegroup_id)[0]
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

	def generate_inventory_list(self, target_to):
		"""Generate the inventory list for the specified TrafficOps instance"""
		with TOSession(self.to_url, verify_cert=self.verify_cert) as traffic_ops_api:
			traffic_ops_api.login(self.to_user, self.to_pass)
			servers = traffic_ops_api.get_servers()[0]
			out = {}
			out['_meta'] = {}
			out['_meta']['hostvars'] = {}
			out[target_to] = {}
			out[target_to]['hosts'] = []
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
				out[target_to]['hosts'].append(fqdn)
				out['_meta']['hostvars'][fqdn] = {}
				out['_meta']['hostvars'][fqdn]['server_toFQDN'] = target_to
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
				if flat_server_profile not in out:
					out['server_profile']['children'].append(
						flat_server_profile)
					out[flat_server_profile] = self.populate_server_profile_vars(
						traffic_ops_api, server['profileId'])
				out[flat_server_profile]['hosts'].append(fqdn)
				if flat_cachegroup not in out:
					out['cachegroup']['children'].append(flat_cachegroup)
					cgdata = self.populate_cachegroups(
						traffic_ops_api,
						server['cachegroupId'])
					out[flat_cachegroup] = cgdata.cgvars
					flat_parent_cg = cgdata.primary_parent_group_name
					flat_second_parent_cg = cgdata.secondary_parent_group_name
					if flat_parent_cg not in out:
						out[flat_parent_cg] = {}
						out[flat_parent_cg]['children'] = []
					if flat_second_parent_cg not in out:
						out[flat_second_parent_cg] = {}
						out[flat_second_parent_cg]['children'] = []
					out[flat_parent_cg]['children'].append(flat_cachegroup)
					out[flat_second_parent_cg]['children'].append(
						flat_cachegroup)
				out[flat_cachegroup]['hosts'].append(fqdn)
				if flat_server_type not in out:
					out['server_type']['children'].append(flat_server_type)
					out[flat_server_type] = {}
					out[flat_server_type]['hosts'] = []
				out[flat_server_type]['hosts'].append(fqdn)
				if flat_server_cdn_name not in out:
					out['server_cdnName']['children'].append(
						flat_server_cdn_name)
					out[flat_server_cdn_name] = {}
					out[flat_server_cdn_name]['hosts'] = []
				out[flat_server_cdn_name]['hosts'].append(fqdn)
				if flat_server_status not in out:
					out['server_status']['children'].append(flat_server_status)
					out[flat_server_status] = {}
					out[flat_server_status]['hosts'] = []
				out[flat_server_status]['hosts'].append(fqdn)
		return out

	def to_inventory(self):
		"""A wrapper function blending the target url in"""
		return self.generate_inventory_list(self.to_url)

#
# Thanks to Maxim for the snipit on handling bool parameters.
# https://stackoverflow.com/questions/15008758/parsing-boolean-values-with-argparse
#


def str2bool(x):
	"""A helper function to help with truthiness"""
	if isinstance(x, bool):
		return x
	if x.lower() in ('yes', 'true', 't', 'y', '1'):
		return True
	if x.lower() in ('no', 'false', 'f', 'n', '0'):
		return False
	raise argparse.ArgumentTypeError('Boolean value expected.')


if __name__ == "__main__":
	PARSER = argparse.ArgumentParser(
		description='Generate an Ansible inventory from TrafficOps')

	PARSER.add_argument(
		'--username',
		type=str,
		metavar='username',
		default=os.environ.get('TO_USER', None),
		help='TrafficOps Username. Environment Var: TO_USER Default: None')
	PARSER.add_argument(
		'--password',
		type=str,
		metavar='password',
		default=os.environ.get('TO_PASSWORD', None),
		help='TrafficOps Password. Environment Var: TO_PASSWORD Default: None')
	PARSER.add_argument(
		'--url',
		type=str,
		metavar='to.kabletown.invalid:8443',
		default=os.environ.get('TO_URL', None),
		help='TrafficOps FQDN and optional HTTPS Port. Environment Var: TO_URL Default: None')
	PARSER.add_argument(
		'--verify_cert',
		type=str2bool,
		default=os.environ.get('TO_VERIFY_CERT', "true"),
		metavar="(true, false, yes, no, t, f, y, n, 0, or 1)",
		help='Perform SSL Certificate Verification. Environment Var: TO_VERIFY_CERT Default: true')
	PARSER.add_argument(
		'--list',
		action='store_true',
		help='Primary argument to enable retrieval of TO data.  Required per calling convention.')
	PARSER.add_argument(
		'--host',
		type=str,
		metavar='do_not_use',
		default=None,
		help='Ignored parameter that must be present due to calling convention.')
	ARGS = PARSER.parse_args()

	if ARGS.username and ARGS.password and ARGS.url:
		if ARGS.list:
			INVENTORY = AnsibleInventory(
				ARGS.username,
				ARGS.password,
				ARGS.url,
				ARGS.verify_cert).to_inventory()
		# Since we're supplying hostvar metadata, --host support isn't required
		else:
			INVENTORY = empty_inventory()
	else:
		INVENTORY = empty_inventory()

	print(json.dumps(INVENTORY))
