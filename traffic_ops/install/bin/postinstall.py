#!/usr/bin/python3
"""
This script is meant as a drop-in replacement for the old _postinstall Perl script.

It does, however, offer several more command-line flags not present in the original, to aid in
testing.

-a, --automatic               If there are questions in the config file which do not have answers,
                              the script will look to the defaults for the answer. If the answer is
                              not in the defaults the script will exit.
--cfile [FILE]                An input config file used to ask and answer questions.
--debug                       Enables verbose logging output.
--defaults [FILE]             Writes out a configuration file with defaults which can be used as
                              input. If no FILE is given, writes to stdout.
-n, --no-root                 Enable running as a non-root user (may cause failure).
-r DIR, --root-directory DIR  Set the directory to be treated as the system's root directory (e.g.
                              for testing). Default: /
-u USER, --ops-user USER      Specify a username to own Traffic Ops files and processes.
                              Default: trafops
-g GROUP, --ops-group GROUP   Specify the group to own Traffic Ops files and processes.
                              Default: trafops
--no-restart-to               Skip restarting Traffic Ops after configuration and database changes
                              are applied.
--no-database                 Skip all database operations.

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
import random
import re
import shutil
import stat
import string
import subprocess
import sys
import typing

# Paths for output configuration files
DATABASE_CONF_FILE = "/opt/traffic_ops/app/conf/production/database.conf"
DB_CONF_FILE = "/opt/traffic_ops/app/db/dbconf.yml"
CDN_CONF_FILE = "/opt/traffic_ops/app/conf/cdn.conf"
LDAP_CONF_FILE = "/opt/traffic_ops/app/conf/ldap.conf"
USERS_CONF_FILE = "/opt/traffic_ops/install/data/json/users.json"
PROFILES_CONF_FILE = "/opt/traffic_ops/install/data/profiles/"
OPENSSL_CONF_FILE = "/opt/traffic_ops/install/data/json/openssl_configuration.json"
PARAM_CONF_FILE = "/opt/traffic_ops/install/data/json/profiles.json"


POST_INSTALL_CFG = "/opt/traffic_ops/install/data/json/post_install.json"

# Log file for the installer
# TODO: determine if logging to a file should be directly supported.
# LOG_FILE = "/var/log/traffic_ops/postinstall.log"

# Log file for CPAN output
# TODO: The Perl used to "rotate" this file on every run, for some reason. Should we?
# CPAN_LOG_FILE = "/var/log/traffic_ops/cpan.log"

# Configuration file output with answers which can be used as input to postinstall
# TODO: Perl used to always write its defaults out to this file when requested.
# Python, instead, outputs to stdout. This is breaking, but more flexible. Change it?
# OUTPUT_CONFIG_FILE = "/opt/traffic_ops/install/bin/configuration_file.json"

class Question:
	"""
	Question represents a single question to be asked of the user, to determine a configuration
	value.

	>>> Question("question", "answer", "var")
	Question(question='question', default='answer', config_var='var', hidden=False)
	"""

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
		qstn = self.question
		ans = self.default
		cfgvr = self.config_var
		hddn = self.hidden
		return f"Question(question='{qstn}', default='{ans}', config_var='{cfgvr}', hidden={hddn})"

	def ask(self) -> str:
		"""
		Asks the user the Question interactively.

		If 'hidden' is true, output will not be echoed.
		"""
		if self.hidden:
			while True:
				passwd = getpass.getpass(self)
				if not passwd:
					continue
				if passwd == getpass.getpass(f"Re-Enter {self.question}: "):
					return passwd
				print("Error: passwords do not match, try again")
		ipt = input(self)
		return ipt if ipt else self.default

	def to_json(self) -> str:
		"""
		Converts a question to JSON encoding.

		>>> Question("Do the thing?", "yes", "cfg_var", True).to_json()
		'{"Do the thing?": "yes", "config_var": "cfg_var", "hidden": true}'
		>>> Question("Do the other thing?", "no", "other cfg_var").to_json()
		'{"Do the other thing?": "no", "config_var": "other cfg_var"}'
		"""
		qstn = self.question
		ans = self.default
		cfgvr = self.config_var
		if self.hidden:
			return f'{{"{qstn}": "{ans}", "config_var": "{cfgvr}", "hidden": true}}'
		return f'{{"{qstn}": "{ans}", "config_var": "{cfgvr}"}}'

	def serialize(self) -> object:
		"""Returns a serializable dictionary, suitable for converting to JSON."""
		return {self.question: self.default, "config_var": self.config_var, "hidden": self.hidden}

class User(typing.NamedTuple):
	"""Users represents a user that will be inserted into the Traffic Ops database."""

	#: The user's username.
	username: str
	#: The user's password - IN PLAINTEXT.
	password: str

class SSLConfig:
	"""SSLConfig bundles the options for generating new (self-signed) SSL certificates"""

	def __init__(self, gen_cert: bool, cfg_map: typing.Dict[str, str]):

		self.gen_cert = gen_cert
		self.rsa_password = cfg_map["rsaPassword"]
		self.params = "/C={country}/ST={state}/L={locality}/O={company}/OU={org_unit}/CN={common_name}/"
		self.params = self.params.format(**cfg_map)

class CDNConfig(typing.NamedTuple):
	"""CDNConfig holds all of the options needed to format a cdn.conf file."""
	gen_secret: bool
	num_secrets: int
	port: int
	num_workers: int
	url: str
	ldap_conf_location: str

	def generate_secret(self, conf):
		"""
		Generates new secrets - if configured to do so - and adds them to the passed cdn.conf
		configuration.
		"""
		if not self.gen_secret:
			return

		if isinstance(conf, dict) and "secrets" in conf and isinstance(conf["secrets"], list):
			logging.debug("Secrets found in cdn.conf file")
		else:
			conf["secrets"] = []
			logging.debug("No secrets found in cdn.conf file")

		conf["secrets"].insert(0, random_word())

		if self.num_secrets and len(conf["secrets"]) > self.num_secrets:
			conf["secrets"] = conf["secrets"][:self.num_secrets - 1]

	def insert_url(self, conf):
		"""
		Inserts the configured URL - if it is not an empty string - into the passed cdn.conf
		configuration, in to.base_url.
		"""
		if not self.url:
			return

		if "to" not in conf or not isinstance(conf["to"], dict):
			conf["to"] = {}
		conf["to"]["base_url"] = self.url

#pylint:disable=bad-continuation
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
		Question(
			"DNS sub-domain for which your CDN is authoritative",
			"cdn1.kabletown.net",
			"dns_subdomain"
		)
	]
}
#pylint:enable=bad-continuation

class ConfigEncoder(json.JSONEncoder):
	"""
	ConfigEncoder encodes a dictionary of filenames to configuration question lists as JSON.

	>>> ConfigEncoder().encode({'/test/file':[Question('question', 'default', 'cfg_var', True)]})
	'{"/test/file": [{"question": "default", "config_var": "cfg_var", "hidden": true}]}'
	"""

	# The linter is just wrong about this
	#pylint:disable=method-hidden
	def default(self, o) -> object:
		"""
		Returns a serializable representation of 'o'.

		Specifically, it does this by attempting to convert a dictionary of filenames to Question
		lists to a dictionary of filenames to lists of dictionaries of strings to strings, falling
		back on default encoding if the proper typing is not found.
		"""
		if isinstance(o, Question):
			return o.serialize()

		return json.JSONEncoder.default(self, o)
	#pylint:enable=method-hidden

def get_config(questions: typing.List[Question], fname: str, automatic: bool = False) -> dict:
	"""Asks all provided questions, or uses their defaults in automatic mode"""

	logging.info("===========%s===========", fname)

	config = {}

	for question in questions:
		answer = question.default if automatic else question.ask()

		config[question.config_var] = answer

	return config

def generate_db_conf(qstns: typing.List[Question], fname: str, automatic: bool, root: str) -> dict:
	"""
	Generates the database.conf file and returns a map of its configuration.

	Also writes the configuration file to the file 'fname' under the directory 'root'.
	"""
	db_conf = get_config(qstns, fname, automatic)
	typ = db_conf.get("type", "UNKNOWN")
	hostname = db_conf.get("hostname", "UNKNOWN")
	port = db_conf.get("port", "UNKNOWN")

	db_conf["description"] = f"{typ} database on {hostname}:{port}"

	path = os.path.join(root, fname.lstrip('/'))
	with open(path, 'w+') as conf_file:
		json.dump(db_conf, conf_file, indent="\t")
		print(file=conf_file)

	logging.info("Database configuration has been saved")

	return db_conf

def generate_todb_conf(qstns: list, fname: str, auto: bool, root: str, conf: dict) -> dict:
	"""
	Generates the dbconf.yml file and returns a map of its configuration.

	Also writes the configuration file to the file 'fname' under the directory 'root'.
	"""
	todbconf = get_config(qstns, fname, auto)

	driver = "postgres"
	if "type" not in conf:
		logging.warning("Driver type not found in todb config; using 'postgres'")
	else:
		driver = "postgres" if conf["type"] == "Pg" else conf["type"]

	path = os.path.join(root, fname.lstrip('/'))
	hostname = conf.get('hostname', 'UNKNOWN')
	port = conf.get('port', 'UNKNOWN')
	user = conf.get('user', 'UNKNOWN')
	password = conf.get('password', 'UNKNOWN')
	dbname = conf.get('dbname', 'UNKNOWN')

	open_line = f"host={hostname} port={port} user={user} password={password} dbname={dbname}"
	with open(path, 'w+') as conf_file:
		print("production:", file=conf_file)
		print("    driver:", driver, file=conf_file)
		print(f"    open: {open_line} sslmode=disable", file=conf_file)

	return todbconf

def generate_ldap_conf(questions: typing.List[Question], fname: str, automatic: bool, root: str):
	"""
	Generates the ldap.conf file by asking the questions or using default answers in auto mode.

	Also writes the configuration to the file 'fname' under the directory 'root'
	"""
	use_ldap_question = [q for q in questions if q.question == "Do you want to set up LDAP?"]
	if not use_ldap_question:
		logging.warning("Couldn't find question asking if LDAP should be set up, using default: no")
		return
	use_ldap = use_ldap_question[0].default if automatic else use_ldap_question[0].ask()

	if use_ldap.casefold() not in {'y', 'yes'}:
		logging.info("Not setting up ldap")
		return

	ldap_conf = get_config([q for q in questions if q is not use_ldap_question[0]], fname, automatic)
	keys = (
		'host',
		'admin_dn',
		'admin_pass',
		'search_base',
		'search_query',
		'insecure',
		'ldap_timeout_secs'
	)

	for key in keys:
		if key not in ldap_conf:
			raise ValueError(f"{key} is a required key in {fname}")

	if not re.fullmatch(r"\S+:\d+", ldap_conf["host"]):
		raise ValueError(f"host in {fname} must be of form 'hostname:port'")

	path = os.path.join(root, fname.lstrip('/'))
	os.makedirs(os.path.dirname(path), exist_ok=True)
	with open(path, 'w+') as conf_file:
		json.dump(ldap_conf, conf_file, indent="\t")
		print(file=conf_file)

def hash_pass(passwd: str) -> str:
	"""
	Generates a Scrypt-based hash of the given password in a Perl-compatible format.
	It's hard-coded - like the Perl - to use 64 random bytes for the salt, n=16384,
	r=8, p=1 and dklen=64.
	"""
	salt = os.urandom(64)
	n = 16384
	r = 8
	p = 1
	hashed = hashlib.scrypt(passwd.encode(), salt=salt, n=n, r=r, p=p, dklen=64)

	hashed_b64 = base64.standard_b64encode(hashed).decode()
	salt_b64 = base64.standard_b64encode(salt).decode()

	return f"SCRYPT:{n}:{r}:{p}:{salt_b64}:{hashed_b64}"

def generate_users_conf(qstns: typing.List[Question], fname: str, auto: bool, root: str) -> User:
	"""
	Generates a users.json file from the given questions and returns a User containing the same
	information.
	"""
	config = get_config(qstns, fname, auto)

	if "tmAdminUser" not in config or "tmAdminPw" not in config:
		raise ValueError(f"{fname} must include 'tmAdminUser' and 'tmAdminPw'")

	hashed_pass = hash_pass(config["tmAdminPw"])

	path = os.path.join(root, fname.lstrip('/'))
	with open(path, 'w+') as conf_file:
		json.dump({"username": config["tmAdminUser"], "password": hashed_pass}, conf_file, indent="\t")
		print(file=conf_file)

	return User(config["tmAdminUser"], config["tmAdminPw"])

def generate_profiles_dir(questions: typing.List[Question]):
	"""
	I truly have no idea what's going on here. This is what the Perl did, so I
	copied it. It does nothing. Literally nothing.
	"""
	user_in = questions

def generate_openssl_conf(questions: typing.List[Question], fname: str, auto: bool) -> SSLConfig:
	"""
	Constructs an SSLConfig by asking the passed questions, or using their default answers if in
	auto mode.
	"""
	cfg_map = get_config(questions, fname, auto)
	if "genCert" not in cfg_map:
		raise ValueError("missing 'genCert' key")

	gen_cert = cfg_map["genCert"].casefold() in {"y", "yes"}

	return SSLConfig(gen_cert, cfg_map)

def generate_param_conf(qstns: typing.List[Question], fname: str, auto: bool, root: str) -> dict:
	"""
	Generates a profiles.json by asking the passed questions, or using their default answers in auto
	mode.

	Also writes the file to 'fname' in the directory 'root'.
	"""
	conf = get_config(qstns, fname, auto)

	path = os.path.join(root, fname.lstrip('/'))
	with open(path, 'w+') as conf_file:
		json.dump(conf, conf_file, indent="\t")
		print(file=conf_file)

	return conf

def sanity_check_config(cfg: typing.Dict[str, typing.List[Question]], automatic: bool) -> int:
	"""
	Checks a user-input configuration file, and outputs the number of files in the
	default question set that did not appear in the input.

	:param cfg: The user's parsed input questions.
	:param automatic: If :keyword:`True` all missing questions will use their default answers.
	Otherwise, the user will be prompted for answers.
	"""
	diffs = 0

	for fname, file in DEFAULTS.items():
		if fname not in cfg:
			logging.warning("File '%s' found in defaults but not config file", fname)
			cfg[fname] = []

		for default_value in file:
			for config_value in cfg[fname]:
				if default_value.config_var == config_value.config_var:
					break
			else:
				question = default_value.question
				answer = default_value.default

				if not automatic:
					logging.info("Prompting user for answer")
					if default_value.hidden:
						answer = default_value.ask()
				elif default_value.hidden:
					logging.info("Adding question '%s' with default answer", question)
				else:
					logging.info("Adding question '%s' with default answer %s", question, answer)

				# The Perl here would ask questions, but those would just get asked later
				# anyway, so I'm not sure why.
				cfg[fname].append(Question(question, answer, default_value.config_var, default_value.hidden))
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
		if not isinstance(questions, list):
			raise ValueError(f"file '{file}' has malformed questions")

		qstns = []
		for qstn in questions:
			if not isinstance(qstn, dict):
				raise ValueError(f"file '{file}' has a malformed question ({qstn})")
			try:
				question = next(key for key in qstn.keys() if qstn not in ("hidden", "config_var"))
			except StopIteration:
				raise ValueError(f"question in '{file}' has no question/answer properties ({qstn})")

			answer = qstn[question]
			if not isinstance(question, str) or not isinstance(answer, str):
				errstr = f"question in '{file}' has malformed question/answer property ({question}: {answer})"
				raise ValueError(errstr)

			del qstn[question]
			hidden = False
			if "hidden" in qstn:
				hidden = bool(qstn["hidden"])
				del qstn["hidden"]

			if "config_var" not in qstn:
				raise ValueError(f"question in '{file}' has no 'config_var' property")
			cfg_var = qstn["config_var"]
			if not isinstance(cfg_var, str):
				raise ValueError(f"question in '{file}' has malformed 'config_var' property ({cfg_var})")
			del qstn["config_var"]

			if qstn:
				logging.warning("Found unknown extra properties in question in '%s' (%r)", file, qstn.keys())

			qstns.append(Question(question, answer, cfg_var, hidden=hidden))
		ret[file] = qstns

	return ret

def setup_maxmind(maxmind_answer: str, root: str):
	"""
	If 'maxmind_answer' is a truthy response ('y' or 'yes' (case-insensitive), sets up a Maxmind
	database using `wget`.
	"""
	if maxmind_answer.casefold() not in {'y', 'yes'}:
		logging.info("Not downloading Maxmind data")
		return

	os.chdir(os.path.join(root, 'opt/traffic_ops/app/public/routing'))

	wget = "/usr/bin/wget"
	cmd = [wget, "https://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz"]
	# Perl ignored errors downloading the databases, so we do too
	try:
		subprocess.run(cmd, capture_output=True, check=True, universal_newlines=True)
	except subprocess.SubprocessError as e:
		logging.error("Failed to download MaxMind data")
		logging.debug("(ipv4) Exception: %s", e)

	cmd[1] = "https://geolite.maxmind.com/download/geoip/database/GeoLiteCityv6-beta/GeoLiteCityv6.dat.gz"
	try:
		subprocess.run(cmd, capture_output=True, check=True, universal_newlines=True)
	except subprocess.SubprocessError as e:
		logging.error("Failed to download MaxMind data")
		logging.debug("(ipv6) Exception: %s", e)

def exec_openssl(description: str, *cmd_args) -> bool:
	"""
	Executes openssl with the supplied command-line arguments.

	:param description: Describes the operation taking place for logging purposes.
	:returns: Whether or not the execution succeeded, success being defined by an exit code of zero
	"""
	logging.info(description)

	cmd = ["/usr/bin/openssl", *cmd_args]

	while True:
		proc = subprocess.run(cmd, capture_output=True, universal_newlines=True, check=False)
		if proc.returncode == 0:
			return True

		logging.debug("openssl exec failed with code %s; stderr: %s", proc.returncode, proc.stderr)
		while True:
			ans = input(f"{description} failed. Try again (y/n) [y]: ")
			if not ans or ans.casefold().startswith('n'):
				return False
			if ans.casefold().startswith('y'):
				break

def setup_certificates(conf: SSLConfig, root: str, ops_user: str, ops_group: str) -> int:
	"""
	Generates self-signed SSL certificates from the given configuration.
	:returns: For whatever reason this subroutine needs to dictate the return code of the script, so that's what it returns.
	"""
	if not conf.gen_cert:
		logging.info("Not generating openssl certification")
		return 0

	if not os.path.isfile('/usr/bin/openssl') or not os.access('/usr/bin/openssl', os.X_OK):
		logging.error("Unable to install SSL certificates as openssl is not installed")
		cmd = os.path.join(root, "opt/traffic_ops/install/bin/generateCert")
		logging.error("Install openssl and then run %s to install SSL certificates", cmd)
		return 4

	logging.info("Installing SSL Certificates")
	logging.info("\n\tWe're now running a script to generate a self signed X509 SSL certificate")
	logging.info("Postinstall SSL Certificate Creation")

	# Perl logs this before actually generating a key. So we do too.
	logging.info("The server key has been generated")

	args = (
		"genrsa",
		"-des3",
		"-out",
		"server.key",
		"-passout",
		f"pass:{conf.rsa_password}",
		"1024"
	)
	if not exec_openssl("Generating an RSA Private Server Key", *args):
		return 1

	args = (
		"req",
		"-new",
		"-key",
		"server.key",
		"-out",
		"server.csr",
		"-passin",
		f"pass:{conf.rsa_password}",
		"-subj",
		conf.params
	)
	if not exec_openssl("Creating a Certificate Signing Request (CSR)", *args):
		return 1

	logging.info("The Certificate Signing Request has been generated")
	os.rename("server.key", "server.key.orig")

	args = (
		"rsa",
		"-in",
		"server.key.orig",
		"-out",
		"server.key",
		"-passin",
		f"pass:{conf.rsa_password}"
	)
	if not exec_openssl("Removing the pass phrase from the server key", *args):
		return 1

	logging.info("The pass phrase has been removed from the server key")

	args = (
		"x509",
		"-req",
		"-days",
		"365",
		"-in",
		"server.csr",
		"-signkey",
		"server.key",
		"-out",
		"server.crt"
	)
	if not exec_openssl("Generating a Self-signed certificate", *args):
		return 1

	logging.info("A server key and self signed certificate has been generated")
	logging.info("Installing a server key and certificate")

	keypath = os.path.join(root, 'etc/pki/tls/private/localhost.key')
	shutil.copy("server.key", keypath)
	os.chmod(keypath, stat.S_IRUSR | stat.S_IWUSR)
	shutil.chown(keypath, user=ops_user, group=ops_group)

	logging.info("The private key has been installed")
	logging.info("Installing self signed certificate")

	certpath = os.path.join(root, 'etc/pki/tls/certs/localhost.crt')
	shutil.copy("server.crt", certpath)
	os.chmod(certpath, stat.S_IRUSR | stat.S_IWUSR)
	shutil.chown(certpath, user=ops_user, group=ops_group)

	logging.info("Saving the self signed csr")

	csrpath = os.path.join(root, 'etc/pki/tls/certs/localhost.csr')
	shutil.copy("server.csr", csrpath)
	os.chmod(csrpath, stat.S_IRUSR | stat.S_IWUSR | stat.S_IRGRP | stat.S_IWGRP | stat.S_IROTH)
	shutil.chown(csrpath, user=ops_user, group=ops_group)

	log_msg = """
        The self signed certificate has now been installed.

        You may obtain a certificate signed by a Certificate Authority using the
        server.csr file saved in the current directory.  Once you have obtained
        a signed certificate, copy it to %s and
        restart Traffic Ops."""
	logging.info(log_msg, certpath)

	cdn_conf_path = os.path.join(root, "opt/traffic_ops/app/conf/cdn.conf")

	try:
		with open(cdn_conf_path) as conf_file:
			cdn_conf = json.load(conf_file)
	except (OSError, json.JSONDecodeError) as e:
		raise OSError(f"reading {cdn_conf_path}: {e}") from e

	if (
		not isinstance(cdn_conf, dict) or
		"hypnotoad" not in cdn_conf or
		not isinstance(cdn_conf["hypnotoad"], dict)
	):
		logging.critical("Malformed %s; improper object and/or missing 'hypnotoad' key", cdn_conf_path)
		return 1

	hypnotoad = cdn_conf["hypnotoad"]
	if (
		"listen" not in hypnotoad or
		not isinstance(hypnotoad["listen"], list) or
		not hypnotoad["listen"] or
		not isinstance(hypnotoad["listen"][0], str)
	):
		log_msg = """	The "listen" portion of %s is missing from %s
	Please ensure it contains the same structure as the one originally installed"""
		logging.error(log_msg, cdn_conf_path, cdn_conf_path)
		return 1

	listen = hypnotoad["listen"][0]

	if f"cert={certpath}" not in listen or f"key={keypath}" not in listen:
		log_msg = """	The "listen" portion of %s is:
	%s
	and does not reference the same "cert=" and "key=" values as are created here.
	Please modify %s to add the following as parameters:
	?cert=%s&key=%s"""
		logging.error(log_msg, cdn_conf_path, listen, cdn_conf_path, certpath, keypath)
		return 1

	return 0

def random_word(length: int = 12) -> str:
	"""
	Returns a randomly generated string 'length' characters long containing only word
	characters ([a-zA-Z0-9_]).
	"""
	word_chars = string.ascii_letters + string.digits + '_'
	return ''.join(random.choice(word_chars) for _ in range(length))

def generate_cdn_conf(questions: typing.List[Question], fname: str, automatic: bool, root: str):
	"""
	Generates some properties of a cdn.conf file based on the passed questions.

	This modifies or writes the file 'fname' under the directory 'root'.
	"""
	cdn_conf = get_config(questions, fname, automatic)

	if "genSecret" not in cdn_conf:
		raise ValueError("missing 'genSecret' config_var")

	gen_secret = cdn_conf["genSecret"].casefold() in {'y', 'yes'}

	try:
		num_secrets = int(cdn_conf["keepSecrets"])
	except KeyError as e:
		raise ValueError("missing 'keepSecrets' config_var") from e
	except ValueError as e:
		raise ValueError(f"invalid 'keepSecrets' config_var value: {e}") from e

	try:
		port = cdn_conf.get("port")
	except KeyError as e:
		raise ValueError("missing 'port' config_var") from e
	except ValueError as e:
		raise ValueError(f"invalid 'port' config_var value: {e}") from e

	try:
		workers = int(cdn_conf["workers"])
	except KeyError as e:
		raise ValueError("missing 'workers' config_var") from e
	except ValueError as e:
		raise ValueError(f"invalid 'workers' config_var value: {e}")

	try:
		url = cdn_conf["base_url"]
	except KeyError as e:
		raise ValueError("missing 'base_url' config_var") from e

	try:
		ldap_loc = cdn_conf["ldap_conf_location"]
	except KeyError as e:
		raise ValueError("missing 'ldap_conf_location' config_var") from e

	conf = CDNConfig(gen_secret, num_secrets, port, workers, url, ldap_loc)

	path = os.path.join(root, fname.lstrip('/'))
	existing_conf = {}
	if os.path.isfile(path):
		with open(path) as conf_file:
			try:
				existing_conf = json.load(conf_file)
			except json.JSONDecodeError as e:
				raise ValueError(f"invalid existing cdn.config at {path}: {e}") from e

	if not isinstance(existing_conf, dict):
		logging.warning("Existing cdn.conf (at '%s') is not an object - overwriting", path)
		existing_conf = {}

	conf.generate_secret(existing_conf)
	conf.insert_url(existing_conf)

	if (
		"traffic_ops_golang" not in existing_conf or
		not isinstance(existing_conf["traffic_ops_golang"], dict)
	):
		existing_conf["traffic_ops_golang"] = {}

	existing_conf["traffic_ops_golang"]["port"] = conf.port
	err_log = os.path.join(root, "var/log/traffic_ops/error.log")
	existing_conf["traffic_ops_golang"]["log_location_error"] = err_log
	access_log = os.path.join(root, "var/log/traffic_ops/access.log")
	existing_conf["traffic_ops_golang"]["log_location_event"] = access_log

	if "hypnotoad" not in existing_conf or not isinstance(existing_conf["hypnotoad"], dict):
		existing_conf["hypnotoad"]["workers"] = conf.num_workers

	with open(path, "w+") as conf_file:
		json.dump(existing_conf, conf_file, indent="\t")
		print(file=conf_file)
	logging.info("CDN configuration has been saved")

def db_connection_string(dbconf: dict) -> str:
	"""
	Constructs a database connection string from the passed configuration object.
	"""
	user = dbconf["user"]
	password = dbconf["password"]
	db_name = "traffic_ops" if dbconf["type"] == "Pg" else dbconf["type"]
	hostname = dbconf["hostname"]
	port = dbconf["port"]
	return f"postgresql://{user}:{password}@{hostname}:{port}/{db_name}"

def exec_psql(conn_str: str, query: str) -> str:
	"""
	Executes SQL queries by forking and exec-ing '/usr/bin/psql'.

	:param conn_str: A "connection string" that defines the postgresql resource in the format
	{schema}://{user}:{password}@{host or IP}:{port}/{database}
	:param query: The query to be run. It can actually be a script containing multiple queries.
	:returns: The comma-separated columns of each line-delimited row of the results of the query.
	"""
	cmd = ["/usr/bin/psql", "--tuples-only", "-d", conn_str, "-c", query]
	proc = subprocess.run(cmd, capture_output=True, universal_newlines=True, check=False)
	if proc.returncode != 0:
		logging.debug("psql exec failed; stderr: %s\n\tstdout: %s", proc.stderr, proc.stdout)
		raise OSError("failed to execute database query")
	return proc.stdout.strip()

def invoke_db_admin_pl(action: str, root: str):
	"""
	Exectues admin with the given action, and looks for it from the given root directory.
	"""
	path = os.path.join(root, "opt/traffic_ops/app")
	# This is a workaround for admin using hard-coded relative paths. That
	# should be fixed at some point, IMO, but for now this works.
	os.chdir(path)
	cmd = [os.path.join(path, "db/admin"), "--env=production", action]
	proc = subprocess.run(cmd, capture_output=True, universal_newlines=True, check=False)
	if proc.returncode != 0:
		logging.debug("admin exec failed; stderr: %s\n\tstdout: %s", proc.stderr, proc.stdout)
		raise OSError(f"Database {action} failed")
	logging.info("Database %s succeeded", action)

def setup_database_data(conn_str: str, user: User, param_conf: dict, root: str):
	"""
	Sets up all necessary initial database data using `/usr/bin/sql`
	"""
	logging.info("paramconf %s", param_conf)
	logging.info("Setting up the database data")

	tables_found_query = '''
		SELECT EXISTS(
			SELECT 1
			FROM pg_tables
			WHERE schemaname = 'public'
				AND tablename = 'tm_user'
		);'''
	if exec_psql(conn_str, tables_found_query) == "t":
		logging.info("Found existing tables skipping table creation")
	else:
		invoke_db_admin_pl("load_schema", root)

	invoke_db_admin_pl("migrate", root)
	invoke_db_admin_pl("seed", root)
	invoke_db_admin_pl("patch", root)

	hashed_pass = hash_pass(user.password)
	insert_admin_query = '''
		INSERT INTO tm_user (username, tenant_id, role, local_passwd, confirm_local_passwd)
		VALUES (
			'{}',
			(SELECT id FROM tenant WHERE name = 'root'),
			(SELECT id FROM role WHERE name = 'admin'),
			'{hashed_pass}',
			'{hashed_pass}'
		)
		ON CONFLICT (username) DO NOTHING;
	'''.format(user.username, hashed_pass=hashed_pass)
	_ = exec_psql(conn_str, insert_admin_query)

	logging.info("=========== Setting up cdn")
	insert_cdn_query = "\n\t-- global parameters" + '''
		INSERT INTO cdn (name, domain_name, dnssec_enabled)
		VALUES ('{cdn_name}', '{dns_subdomain}', false)
		ON CONFLICT DO NOTHING;
	'''.format(**param_conf)
	logging.info("\n%s", insert_cdn_query)
	_ = exec_psql(conn_str, insert_cdn_query)

	tm_url = param_conf["tm.url"]

	logging.info("=========== Setting up parameters")
	insert_parameters_query = "\n\t-- global parameters" + '''
		INSERT INTO parameter (name, config_file, value)
		VALUES ('tm.url', 'global', '{tm_url}'),
			('tm.infourl', 'global', '{tm_url}/doc'),
		-- CRConfic.json parameters
			('geolocation.polling.url', 'CRConfig.json', '{tm_url}/routing/GeoLite2-City.mmdb.gz'),
			('geolocation6.polling.url', 'CRConfig.json', '{tm_url}/routing/GeoLiteCityv6.dat.jz')
		ON CONFLICT (name, config_file, value) DO NOTHING;
	'''.format(tm_url=tm_url)
	logging.info("\n%s", insert_parameters_query)
	_ = exec_psql(conn_str, insert_parameters_query)

	logging.info("\n=========== Setting up profiles")
	insert_profiles_query = "\n\t-- global parameters" + '''
		INSERT INTO profile (name, description, type, cdn)
		VALUES ('GLOBAL' 'Global Traffic Ops profile, DO NOT DELETE', 'UNK_PROFILE', (SELECT id FROM cdn WHERE name='ALL'))
		ON CONFLICT DO NOTHING;

		INSERT INTO profile_parameter (profile, parameter)
		VALUES
			(
				(SELECT id FROM profile WHERE name = 'GLOBAL'),
				(
					SELECT id
					FROM parameter
					WHERE name = 'tm.url'
						AND config_file = 'global'
						AND value = '{tm_url}'
				)
			),
			(
				(SELECT id FROM profile WHERE name = 'GLOBAL'),
				(
					SELECT id
					FROM parameter
					WHERE name = 'tm.infourl'
						AND config_file = 'global'
						AND value = '{tm_url}/doc'
				)
			),
			(
				(SELECT id FROM profile WHERE name = 'GLOBAL'),
				(
					SELECT id
					FROM parameter
					WHERE name = 'geolocation.polling.url'
						AND config_file = 'CRConfig.json'
						AND value = '{tm_url}/routing/GeoLite2-City.mmdb.gz'
				)
			),
			(
				(SELECT id FROM profile WHERE name = 'GLOBAL'),
				(
					SELECT id
					FROM parameter
					WHERE name = 'geolocation6.polling.url'
						AND config_file = 'CRConfig.json'
						AND value = '{tm_url}/routing/GeoLiteCityv6.mmdb.gz'
				)
			)
		ON CONFLICT (profile, parameter) DO NOTHING;
	'''.format(tm_url=tm_url)
	logging.info("\n%s", insert_profiles_query)
	_ = exec_psql(conn_str, insert_cdn_query)

def main(
automatic: bool,
debug: bool,
defaults: typing.Optional[str],
cfile: typing.Optional[str],
root_dir: str,
ops_user: str,
ops_group: str,
no_restart_to: bool,
no_database: bool
) -> int:
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
					with open(defaults, "w") as dump_file:
						json.dump(DEFAULTS, dump_file, indent="\t")
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

	if not cfile:
		logging.info("No input file given - using defaults")
		user_input = DEFAULTS
	else:
		logging.info("Using input file %s", cfile)
		try:
			with open(cfile) as conf_file:
				user_input = unmarshal_config(json.load(conf_file))
			diffs = sanity_check_config(user_input, automatic)
			logging.info(
			"File sanity check complete - found %s difference%s",
			diffs,
			'' if diffs == 1 else 's'
			)
		except (OSError, ValueError, json.JSONDecodeError) as e:
			logging.critical("Reading in input file '%s': %s", cfile, e)
			return 1

	try:
		path = os.path.join(root_dir, "opt/traffic_ops/install/bin")
		# os.chdir(path)
	except OSError as e:
		logging.critical("Attempting to change directory to '%s': %s", path, e)
		return 1

	try:
		dbconf = generate_db_conf(user_input[DATABASE_CONF_FILE], DATABASE_CONF_FILE, automatic, root_dir)
		todbconf = generate_todb_conf(user_input[DB_CONF_FILE], DB_CONF_FILE, automatic, root_dir, dbconf)
		generate_ldap_conf(user_input[LDAP_CONF_FILE], LDAP_CONF_FILE, automatic, root_dir)
		admin_conf = generate_users_conf(
		user_input[USERS_CONF_FILE],
		USERS_CONF_FILE,
		automatic,
		root_dir
		)
		generate_profiles_dir(user_input[PROFILES_CONF_FILE])
		opensslconf = generate_openssl_conf(user_input[OPENSSL_CONF_FILE], OPENSSL_CONF_FILE, automatic)
		paramconf = generate_param_conf(user_input[PARAM_CONF_FILE], PARAM_CONF_FILE, automatic, root_dir)
		postinstall_cfg = os.path.join(root_dir, POST_INSTALL_CFG.lstrip('/'))
		if not os.path.isfile(postinstall_cfg):
			with open(postinstall_cfg, 'w+') as conf_file:
				print("{}", file=conf_file)
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

	try:
		cert_code = setup_certificates(opensslconf, root_dir, ops_user, ops_group)
		if cert_code:
			return cert_code
	except OSError as e:
		logging.critical("Setting up SSL Certificates: %s", e)
		return 1

	try:
		generate_cdn_conf(user_input[CDN_CONF_FILE], CDN_CONF_FILE, automatic, root_dir)
	except OSError as e:
		logging.critical("Generating cdn.conf: %s", e)
		return 1

	if not no_database:
		try:
			conn_str = db_connection_string(dbconf)
		except KeyError as e:
			logging.error("Missing database connection variable: %s", e)
			logging.error("Can't connect to the database.  Use the script `/opt/traffic_ops/install/bin/todb_bootstrap.sh` on the db server to create it and run `postinstall` again.")
			return -1

		if not os.path.isfile("/usr/bin/psql") or not os.access("/usr/bin/psql", os.X_OK):
			logging.critical("psql is not installed, please install it to continue with database setup")
			return 1

		try:
			setup_database_data(conn_str, admin_conf, paramconf, root_dir)
		except (OSError, subprocess.SubprocessError)as e:
			logging.error("Failed to set up database: %s", e)
			logging.error("Can't connect to the database.  Use the script `/opt/traffic_ops/install/bin/todb_bootstrap.sh` on the db server to create it and run `postinstall` again.")
			return -1


	if not no_restart_to:
		logging.info("Starting Traffic Ops")
		try:
			cmd = ["/sbin/service", "traffic_ops", "restart"]
			proc = subprocess.run(cmd, capture_output=True, universal_newlines=True, check=False)
		except (OSError, subprocess.SubprocessError) as e:
			logging.critical("Failed to restart Traffic Ops, return code %s: %s", proc.returncode, e)
			logging.debug("stderr: %s\n\tstdout: %s", proc.stderr, proc.stdout)
			return 1
		# Perl didn't actually do any "waiting" before reporting success, so
		# neither do we
		logging.info("Waiting for Traffic Ops to restart")
	else:
		logging.info("Skipping Traffic Ops restart")
	logging.info("Success! Postinstall complete.")

	return 0

if __name__ == '__main__':
	PARSER = argparse.ArgumentParser()
	PARSER.add_argument(
		"-a",
		"--automatic",
		help="If there are questions in the config file which do not have answers, the script " +
		"will look to the defaults for the answer. If the answer is not in the defaults the " +
		"script will exit",
		action="store_true"
	)
	PARSER.add_argument(
		"--cfile",
		help="An input config file used to ask and answer questions",
		type=str,
		default=None
	)
	PARSER.add_argument("--debug", help="Enables verbose output", action="store_true")
	PARSER.add_argument(
		"--defaults",
		help="Writes out a configuration file with defaults which can be used as input",
		type=str,
		nargs="?",
		default=None,
		const=""
	)
	PARSER.add_argument(
		"-n",
		"--no-root",
		help="Enable running as a non-root user (may cause failure)",
		action="store_true"
	)
	PARSER.add_argument(
		"-r",
		"--root-directory",
		help="Set the directory to be treated as the system's root directory (e.g. for testing)",
		type=str,
		default="/"
	)
	PARSER.add_argument(
		"-u",
		"--ops-user",
		help="Specify a username to own Traffic Ops files and processes",
		type=str,
		default="trafops"
	)
	PARSER.add_argument(
		"-g",
		"--ops-group",
		help="Specify the group to own Traffic Ops files and processes",
		type=str,
		default="trafops"
	)
	PARSER.add_argument(
		"--no-restart-to",
		help="Skip restarting Traffic Ops after configuration and database changes are applied",
		action="store_true"
	)
	PARSER.add_argument("--no-database", help="Skip all database operations", action="store_true")

	ARGS = PARSER.parse_args()

	if not ARGS.no_root and os.getuid() != 0:
		logging.error("You must run this script as the root user")
		logging.shutdown()
		sys.exit(1)

	try:
		EXIT_CODE = main(
		ARGS.automatic,
		ARGS.debug,
		ARGS.defaults,
		ARGS.cfile,
		os.path.abspath(ARGS.root_directory),
		ARGS.ops_user,
		ARGS.ops_group,
		ARGS.no_restart_to,
		ARGS.no_database
		)
		sys.exit(EXIT_CODE)
	except KeyboardInterrupt:
		sys.exit(1)
	finally:
		logging.shutdown()
