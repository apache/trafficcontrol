#!/usr/bin/python3

#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

import getpass
import logging
import sys
import typing

# Paths for output configuration files
DATABASE_CONF_FILE = "/opt/traffic_ops/app/conf/production/database.conf"
DB_CONF_FILE       = "/opt/traffic_ops/app/db/dbconf.yml"
CDN_CONF_FILE      = "/opt/traffic_ops/app/conf/cdn.conf"
LDAP_CONF_FILE     = "/opt/traffic_ops/app/conf/ldap.conf"
USERS_CONF_FILE    = "/opt/traffic_ops/install/data/json/users.json"
PROFILES_CONF_FILE = "/opt/traffic_ops/install/data/profiles/"
OPENSSL_CONF_FILE  = "/opt/traffic_ops/install/data/json/openssl_configuration.json"
PARAM_CONF_FILE    = "/opt/traffic_ops/install/data/json/profiles.json"

CUSTOM_PROFILE_DIR = PROFILES_CONF_FILE + "custom"

# Location of Traffic Ops profiles
PROFILE_DIR = "/opt/traffic_ops/install/data/profiles/"
POST_INSTALL_CFG = "/opt/traffic_ops/install/data/json/post_install.json"

# Log file for the installer
LOG_FILE = "/var/log/traffic_ops/postinstall.log"

# Log file for CPAN output
CPAN_LOG_FILE = "/var/log/traffic_ops/cpan.log"

# Configuration file output with answers which can be used as input to postinstall
OUTPUT_CONFIG_FILE = "/opt/traffic_ops/install/bin/configuration_file.json"

def get_field(question: str, config_answer: str, hidden: bool = False) -> str:
	if hidden:
		while True:
			pw = getpass.getpass(question)
			if pw == getpass.getpass(f"Re-Enter {question}"):
				return pw
			print("Error: passwords do not match, try again", file=sys.stderr)

	if config_answer:
		ipt = input(f"{question} [{config_answer}]: ")
		return ipt if ipt else config_answer
	return input(question + ": ")

def get_config(user_input: typing.Dict[str, typing.List[typing.Dict[str, str]]], fname: str, automatic: bool = False) -> dict:

	if fname not in user_input:
		raise ValueError(f"No {fname} found in config")

	logging.info(f"==========={fname}===========")

	config = {}

	for var in user_input[fname]:
		question = next(key for key in var.keys() if key != "hidden" and key != "config_var")
		hidden = var.get("hidden")
		answer = var.get(question) if automatic else get_field(question, var.get(question), hidden)

		config[var.get(config_var)] = answer

	return config

def generate_db_conf(user_input: dict, fname: str, todb_fname: str, automatic: bool):
	"""
	"""
	db_conf = get_config(user_input, fname, automatic)
