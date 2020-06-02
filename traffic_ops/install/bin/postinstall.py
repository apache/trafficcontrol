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
import base64
import getpass
import hashlib
import json
import logging
import os
import re
import subprocess
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
				if pw == getpass.getpass(f"Re-Enter {self.question}: "):
					return pw
				print("Error: passwords do not match, try again")
		ipt = input(self)
		return ipt if ipt else self.default

	def toJSON(self) -> str:
		"""
		Converts a question to JSON encoding.
		>>> Question("Do the thing?", "yes", "cfg_var", True).toJSON()
		'{"Do the thing?": "yes", "config_var": "cfg_var", "hidden": true}'
		>>> Question("Do the other thing?", "no", "other cfg_var").toJSON()
		'{"Do the other thing?": "no", "config_var": "other cfg_var"}'
		"""
		if self.hidden:
			return '{{"{}": "{}", "config_var": "{}", "hidden": true}}'.format(self.question, self.default, self.config_var)
		return '{{"{}": "{}", "config_var": "{}"}}'.format(self.question, self.default, self.config_var)

	def serialize(self) -> object:
		return {self.question: self.default, "config_var": self.config_var, "hidden": self.hidden}

class User(typing.NamedTuple):
	username: str
	password: str

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

class ConfigEncoder(json.JSONEncoder):
	"""
	ConfigEncoder encodes a dictionary of filenames to configuration question lists as JSON
	>>> ConfigEncoder().encode({'/test/file':[Question('question', 'default', 'cfg_var', True)]})
	'{"/test/file": [{"question": "default", "config_var": "cfg_var", "hidden": true}]}'
	"""
	def default(self, o) -> object:
		"""
		Returns a serializable representation of 'o' - specifically by attempting
		to convert a dictionary of filenames to Question lists to a dictionary of
		filenames to lists of dictionaries of strings to strings, falling back on
		default encoding if the proper typing is not found.
		"""
		if isinstance(o, Question):
			return o.serialize()

		return json.JSONEncoder.default(self, o)

def get_config(questions: typing.List[Question], fname: str, automatic: bool = False) -> dict:

	logging.info(f"==========={fname}===========")

	config = {}

	for q in questions:
		answer = q.default if automatic else q.ask()

		config[q.config_var] = answer

	return config

def generate_db_conf(questions: typing.List[Question], fname: str, automatic: bool, root: str) -> dict:
	"""
	"""
	db_conf = get_config(questions, fname, automatic)
	db_conf["description"] = f"{db_conf.get('type', 'UNKNOWN')} database on {db_conf.get('hostname','UNKOWN')}:{db_conf.get('port', 'UNKNOWN')}"

	path = os.path.join(root, fname.lstrip('/'))
	with open(path, 'w+') as fd:
		json.dump(db_conf, fd, indent="\t")
		print(file=fd)

	logging.info("Database configuration has been saved")

	return db_conf

def generate_todb_conf(questions: typing.List[Question], fname: str, automatic: bool, root: str, dbconf: dict) -> dict:
	todbconf = get_config(questions, fname, automatic)

	driver = "postgres"
	if "type" not in dbconf:
		logging.warning("Driver type not found in todb config; using 'postgres'")
	else:
		driver = "postgres" if dbconf["type"] == "Pg" else dbconf["type"]

	path = os.path.join(root, fname.lstrip('/'))
	hostname = dbconf.get('hostname', 'UNKNOWN')
	port = dbconf.get('port', 'UNKNOWN')
	user = dbconf.get('user', 'UNKNOWN')
	password = dbconf.get('password', 'UNKNOWN')
	dbname = dbconf.get('dbname', 'UNKNOWN')
	with open(path, 'w+') as fd:
		print("production", file=fd)
		print("    driver:", driver, file=fd)
		print(f"    open: host={hostname} port={port} user={user} password={password} dbname={dbname} sslmode=disable", file=fd)

	return todbconf

def generate_ldap_conf(questions: typing.List[Question], fname: str, automatic: bool, root: str):
	use_ldap_question = [q for q in questions if q.question == "Do you want to set up LDAP?"]
	if not use_ldap_question:
		logging.warning("Couldn't find question asking if LDAP should be set up, using default: no")
		return
	use_ldap = use_ldap_question[0].default if automatic else use_ldap_question[0].ask()

	if use_ldap.casefold() not in {'y', 'yes'}:
		logging.info("Not setting up ldap")
		return

	ldapConf = get_config([q for q in questions if q is not use_ldap_question[0]], fname, automatic)
	for key in ('host', 'admin_dn', 'admin_pass', 'search_base', 'search_query', 'insecure', 'ldap_timeout_secs'):
		if key not in ldapConf:
			raise ValueError(f"{key} is a required key in {fname}")

	if not re.fullmatch(r"\S+:\d+", ldapConf["host"]):
		raise ValueError(f"host in {fname} must be of form 'hostname:port'")

	path = os.path.join(root, fname.lstrip('/'))
	os.makedirs(os.path.dirname(path), exist_ok=True)
	with open(path, 'w+') as fd:
		json.dump(ldapConf, fd, indent="\t")
		print(file=fd)

def hash_pass(passwd: str) -> str:
	"""
	Generates a Scrypt-based hash of the given password in a Perl-compatible format.
	It's hard-coded - like the Perl - to use 64 random bytes for the salt, n=16384,
	r=8, p=1 and dklen=64.
	"""
	salt=os.urandom(64)
	n=16384
	r=8
	p=1
	hashed = hashlib.scrypt(passwd.encode(), salt=salt, n=n, r=r, p=p, dklen=64)

	hashed_b64 = base64.standard_b64encode(hashed).decode()
	salt_b64 = base64.standard_b64encode(salt).decode()

	return f"SCRYPT:{n}:{r}:{p}:{salt_b64}:{hashed_b64}"

def generate_users_conf(questions: typing.List[Question], fname: str, automatic: bool, root: str) -> User:
	config = get_config(questions, fname, automatic)

	if "tmAdminUser" not in config or "tmAdminPw" not in config:
		raise ValueError(f"{fname} must include 'tmAdminUser' and 'tmAdminPw'")

	hashedPass = hash_pass(config["tmAdminPw"])

	path = os.path.join(root, fname.lstrip('/'))
	with open(path, 'w+') as fd:
		json.dump({"username": config["tmAdminUser"], "password": hashedPass}, fd, indent="\t")
		print(file=fd)

	return User(config["tmAdminUser"], config["tmAdminPw"])

def generate_profiles_dir(questions: typing.List[Question], fname: str):
	"""
	I truly have no idea what's going on here. This is what the Perl did, so I
	copied it. It does nothing. Literally nothing.
	"""
	user_in = questions

def generate_openssl_conf(questions: typing.List[Question], fname: str, automatic: bool) -> dict:
	return get_config(questions, fname, automatic)

def generate_param_conf(questions: typing.List[Question], fname: str, automatic: bool, root: str) -> dict:
	conf = get_config(questions, fname, automatic)

	path = os.path.join(root, fname.lstrip('/'))
	with open(path, 'w+') as fd:
		json.dump(conf, fd, indent="\t")
		print(file=fd)

	return conf

def sanity_check_config(cfg: typing.Dict[str, typing.List[Question]], automatic: bool) -> int:
	"""
	Checks a user-input configuration file, and outputs the number of files in the
	default question set that did not appear in the input.

	:param cfg: The user's parsed input questions.
	:param automatic: If :keyword:`True` all missing questions will use their default answers. Otherwise, the user will be prompted for answers.
	"""
	diffs = 0

	for fname, file in DEFAULTS.items():
		if fname not in cfg:
			logging.warning("File '%s' found in defaults but not config file", fname)
			cfg[fname] = []

		for defaultValue in file:
			for configValue in cfg[fname]:
				if defaultValue.config_var == configValue.config_var:
					break
			else:
				continue

			question = defaultValue.question
			answer = defaultValue.answer

			if not automatic:
				logging.info("Prompting user for answer")
				if defaultValue.hidden:
					answer = defaultValue.ask()
			else:
				logging.info("Adding question '%s' with default answer%s", question, f" {answer}" if not defaultValue.hidden else "")

			# The Perl here would ask questions, but those would just get asked later
			# anyway, so I'm not sure why.
			cfg[fname].append(Question(question, answer, defaultValue.config_var, defaultValue.hidden))
			diffs += 1

	return diffs

def unmarshal_config(dct: dict) -> typing.Dict[str, typing.List[Question]]:
	"""
	Reads in a raw parsed configuration file and returns the resulting configuration.
	>>> unmarshal_config({"test": [{"Do the thing?": "yes", "config_var": "thing"}]})
	{'test': [Question(question='Do the thing?', default='yes', config_var='thing', hidden=False)]}
	>>> unmarshal_config({"test": [{"foo": "", "config_var": "bar", "hidden": True}]})
	{'test': [Question(question='foo', default='', config_var='bar', hidden=True)]}
	"""
	ret = {}
	for file, questions in dct.items():
		if type(questions) is not list:
			raise ValueError(f"file '{file}' has malformed questions")

		qs = []
		for q in questions:
			if type(q) is not dict:
				raise ValueError(f"file '{file}' has a malformed question ({q})")
			try:
				question = next(key for key in q.keys() if q != "hidden" and q != "config_var")
			except StopIteration:
				raise ValueError(f"question in '{file}' has no question/answer properties ({q})")

			answer = q[question]
			if type(question) is not str or type(answer) is not str:
				raise ValueError(f"question in '{file}' has malformed question/answer property ({question}: {answer})")

			del q[question]
			hidden = False
			if "hidden" in q:
				hidden = bool(q["hidden"])
				del q["hidden"]

			if "config_var" not in q:
				raise ValueError(f"question in '{file}' has no 'config_var' property")
			cfg_var = q["config_var"]
			if type(cfg_var) is not str:
				raise ValueError(f"question in '{file}' has malformed 'config_var' property ({cfg_var})")
			del q["config_var"]

			if q:
				logging.warning("Found unknown extra properties in question in '%s' (%r)", file, q.keys())

			qs.append(Question(question, answer, cfg_var, hidden=hidden))
		ret[file] = qs

	return ret

def setup_maxmind(mm: str, root: str):
	"""
	If 'mm' is a truthy response ('y' or 'yes' (case-insensitive), sets up a Maxmind database using `wget`.
	"""
	if mm.casefold() not in {'y', 'yes'}:
		logging.info("Not downloading Maxmind data")

	os.chdir(os.path.join(root, 'opt/traffic_ops/app/public/routing'))

	# Perl ignored errors downloading the databases, so we do too
	try:
		subprocess.run(["/usr/bin/wget", "https://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz"], capture_output=True, check=True, universal_newlines=True)
	except subprocess.SubprocessError as e:
		logging.error("Failed to download MaxMind data")
		logging.debug("(ipv4) Exception: %s", e)

	try:
		subprocess.run(["/usr/bin/wget", "https://geolite.maxmind.com/download/geoip/database/GeoLiteCityv6-beta/GeoLiteCityv6.dat.gz"], capture_output=True, check=True, universal_newlines=True)
	except subprocess.SubprocessError as e:
		logging.error("Failed to download MaxMind data")
		logging.debug("(ipv6) Exception: %s", e)

def main(automatic: bool, debug: bool, defaults: str = None, cfile: str = None, root_dir: str = "/") -> int:
	"""
	Runs the main routine given the parsed arguments as input.
	"""
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
				json.dump(DEFAULTS, sys.stdout, cls=ConfigEncoder, indent="\t")
				print()
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
				userInput = json.load(fd, object_hook=unmarshal_config)
			diffs = sanity_check_config(userInput, automatic)
			logging.info(f"File sanity check complete - found {diffs} difference{'' if diffs == 1 else 's'}")
		except (OSError, ValueError, json.JSONDecodeError) as e:
			logging.critical("Reading in input file '%s': %s", cfile, e)
			return 1

	try:
		path = os.path.join(root_dir, "opt/traffic_ops/install/bin")
		# os.chdir(path)
	except OSError as e:
		logging.critical(f"Attempting to change directory to '{path}': {e}")
		return 1

	try:
		dbconf = generate_db_conf(userInput[DATABASE_CONF_FILE], DATABASE_CONF_FILE, automatic, root_dir)
		todbconf = generate_todb_conf(userInput[DB_CONF_FILE], DB_CONF_FILE, automatic, root_dir, dbconf)
		generate_ldap_conf(userInput[LDAP_CONF_FILE], LDAP_CONF_FILE, automatic, root_dir)
		admin_conf = generate_users_conf(userInput[USERS_CONF_FILE], USERS_CONF_FILE, automatic, root_dir)
		custom_profile = generate_profiles_dir(userInput[PROFILES_CONF_FILE], PROFILES_CONF_FILE)
		opensslconf = generate_openssl_conf(userInput[OPENSSL_CONF_FILE], OPENSSL_CONF_FILE, automatic)
		paramconf = generate_param_conf(userInput[PARAM_CONF_FILE], PARAM_CONF_FILE, automatic, root_dir)
		postinstall_cfg = os.path.join(root_dir, POST_INSTALL_CFG.lstrip('/'))
		if not os.path.isfile(postinstall_cfg):
			with open(postinstall_cfg, 'w+') as fd:
				print("{}", file=fd)
	except OSError as e:
		logging.critical("Writing configuration: %s", e)
		return 1
	except ValueError as e:
		logging.critical("Generating configuration: %s", e)
		return 1

	try:
		setup_maxmind(todbconf.get("maxmind", "no"), root_dir)
	except OSError as e:
		logging.critical("Setting up MaxMind: %s", e)
		return 1

	return 0

if __name__ == '__main__':
	parser = argparse.ArgumentParser()
	parser.add_argument("-a", "--automatic", help="If there are questions in the config file which do not have answers, the script will look to the defaults for the answer. If the answer is not in the defaults the script will exit", action="store_true")
	parser.add_argument("--cfile", help="An input config file used to ask and answer questions", type=str, default=None)
	parser.add_argument("--debug", help="Enables verbose output", action="store_true")
	parser.add_argument("--defaults", help="Writes out a configuration file with defaults which can be used as input", type=str, nargs="?", default=None, const="")
	parser.add_argument("-n", "--no-root", help="Enable running as a non-root user (may cause failure)", action="store_true")
	parser.add_argument("-r", "--root-directory", help="Set the directory to be treated as the system's root directory (e.g. for testing)", type=str, default="/")

	args = parser.parse_args()

	if not args.no_root and os.getuid() != 0:
		logging.error("You must run this script as the root user")
		sys.exit(1)
	sys.exit(main(args.automatic, args.debug, args.defaults, args.cfile, args.root_directory))
