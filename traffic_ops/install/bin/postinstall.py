#!/usr/bin/python3
"""
>>> [c for c in [[a for a in b if not a.config_var] for b in DEFAULTS.values()] if c]
[]
"""

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

import argparse
import getpass
import logging
import os
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

class Question():

	def __init__(self, question: str, default: str, config_var: str, hidden: bool = False):
		self.question = question
		self.default = default
		self.config_var = config_var
		self.hidden = hidden

	def __str__(self) -> str:
		if self.default:
			return f"{self.question} [{self.default}]: "
		return f"{self.question}: "

	def __repr__(self) -> str:
		return f"Question(question='{self.question}', default='{self.default}', config_var='{self.config_var}', hidden={self.hidden})"

	def ask(self) -> str:
		if self.hidden:
			while True:
				pw = getpass.getpass(self)
				if pw == getpass.getpass(f"Re-Enter {self.question}"):
					return pw
				print("Error: passwords do not match, try again")
		ipt = input(self)
		return ipt if ipt else self.default


# The default question/answer set
DEFAULTS = {
	DATABASE_CONF_FILE: [
		Question("Database type", "Pg", "type"),
		Question("Database name", "traffic_ops", "dbname"),
		Question("Database server hostname IP or FQDN", "localhost", "hostname"),
		Question("Database port number", "5432", "port"),
		Question("Traffic Ops database user", "traffic_ops", "user"),
		Question("Password for Traffic Ops database user", "", "password", hidden=True)
	],
	DB_CONF_FILE: [
		Question("Database server root (admin) user", "postgres", "pgUser"),
		Question("Password for database server admin", "", "pgPassword", hidden=True),
		Question("Download Maxmind Database?", "yes", "maxmind")
	],
	CDN_CONF_FILE: [
		Question("Generate a new secret?", "yes", "genSecret"),
		Question("Number of secrets to keep?", "1", "keepSecrets"),
		Question("Port to serve on?", "443", "port"),
		Question("Number of workers?", "12", "workers"),
		Question("Traffic Ops url?", "http://localhost:3000", "base_url"),
		Question("ldap.conf location?", "/opt/traffic_ops/app/conf/ldap.conf", "ldap_conf_location")
	],
	LDAP_CONF_FILE:[
		Question("Do you want to set up LDAP?", "no", "setupLdap"),
		Question("LDAP server hostname", "", "host"),
		Question("LDAP Admin DN", "", "admin_dn"),
		Question("LDAP Admin Password", "", "admin_pass", hidden=True),
		Question("LDAP Search Base", "", "search_base"),
		Question("LDAP Search Query", "", "search_query"),
		Question("LDAP Skip TLS verify", "", "insecure"),
		Question("LDAP Timeout Seconds", "", "ldap_timeout_secs")
	],
	USERS_CONF_FILE: [
		Question("Administration username for Traffic Ops", "admin", "tmAdminUser"),
		Question("Password for the admin user", "", "tmAdminPw", hidden=True)
	],
	PROFILES_CONF_FILE: [
		Question("Add custom profiles?", "no", "custom_profiles")
	],
	OPENSSL_CONF_FILE: [
		Question("Do you want to generate a certificate?", "yes", "genCert"),
		Question("Country Name (2 letter code)", "", "country"),
		Question("State or Province Name (full name)", "", "state"),
		Question("Locality Name (eg, city)", "", "locality"),
		Question("Organization Name (eg, company)", "", "company"),
		Question("Organizational Unit Name (eg, section)", "", "org_unit"),
		Question("Common Name (eg, your name or your server's hostname)", "", "common_name"),
		Question("RSA Passphrase", "CHANGEME!!", "rsaPassword", hidden=True)
	],
	PARAM_CONF_FILE: [
		Question("Traffic Ops url", "https://localhost", "tm.url"),
		Question("Human-readable CDN Name. (No whitespace, please)", "kabletown_cdn", "cdn_name"),
		Question("DNS sub-domain for which your CDN is authoritative", "cdn1.kabletown.net", "dns_subdomain")
	]
}

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

def sanity_check_config(cfg: dict):
	diffs = 0

	for fname, file in DEFAULTS:
		if fname not in cfg:
			logging.warning("File '%s' found in defaults but not config file", fname)
			cfg[fname] = []

		for defaultValue in file:
			for configValue in cfg[fname]:
				if defaultValue["config_var"] == configValue.get("config_var"):
					break
			else:
				continue



def main(automatic: bool, debug: bool, defaults: str = None, cfile: str = None) -> int:
	if debug:
		logging.getLogger().setLevel(logging.DEBUG)
	else:
		logging.getLogger().setLevel(logging.INFO)

	# At this point, the Perl script... unzipped its own logfile?

	logging.info("Starting postinstall")
	# The Perl printed this whether or not the logger was actually at the debug level
	# so we do too
	logging.info("Debug is on")

	if automatic:
		logging.info("Running in automatic mode")

	if defaults is not None:
		try:
			if defaults:
				try:
					with open(defaults, "w") as fd:
						json.dump(DEFAULTS, fd, indent="\t")
				except OSError as e:
					logging.critical("Writing output: %s", e)
					return 1
			else:
				json.dump(DEFAULTS, sys.stdout, indent="\t")
		except ValueError as e:
			logging.critical("Converting defaults to JSON: %s", e)
			return 1
		return 0

	userInput = None
	if not cfile:
		logging.info("No input file given - using defaults")
		userInput = DEFAULTS
	else:
		logging.info("Using input file %s", cfile)
		try:
			with open(cfile) as fd:
				userInput = json.load(fd)
			sanity_check_config(userInput)
		except (OSError, ValueError) as e:
			logging.critical("Reading in input file '%s': %s", cfile, e)
			return 1
	return 0

if __name__ == '__main__':
	parser = argparse.ArgumentParser()
	parser.add_argument("-a", "--automatic", help="If there are questions in the config file which do not have answers, the script will look to the defaults for the answer. If the answer is not in the defaults the script will exit", action="store_true")
	parser.add_argument("--cfile", help="An input config file used to ask and answer questions", type=str, default=None)
	parser.add_argument("--debug", help="Enables verbose output", action="store_true")
	parser.add_argument("--defaults", help="Writes out a configuration file with defaults which can be used as input", type=str, nargs="?", default=None, const="")
	parser.add_argument("-n", "--no-root", help="Enable running as a non-root user (may cause failure)", action="store_true")

	args = parser.parse_args()

	if not args.no_root and os.getuid() != 0:
		logging.error("You must run this script as the root user")
		sys.exit(1)
	sys.exit(main(args.automatic, args.debug, args.defaults, args.cfile))
