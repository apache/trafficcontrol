#!/usr/bin/env python3
"""
This script is meant as a drop-in replacement for the old _postinstall.pl Perl script.

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
# There's a bug in asteroid with Python 3.9's NamedTuple being
# recognized for the dynamically generated class that it is. Should be fixed
# with the next release, but 'til then...

import argparse
import base64
import errno
import getpass
import grp
import hashlib
import json
import logging
import os
import pwd
import random
import re
import shutil
import stat
import string
import subprocess
import sys

from struct import unpack, pack
from typing import Any, Dict, List, Optional, Union, NamedTuple

# Paths for output configuration files
DATABASE_CONF_FILE = "/opt/traffic_ops/app/conf/production/database.conf"
DB_CONF_FILE = "/opt/traffic_ops/app/db/dbconf.yml"
TV_DATABASE_CONF_FILE = "/opt/traffic_ops/app/conf/production/tv.conf"
TV_DB_CONF_FILE = "/opt/traffic_ops/app/db/trafficvault/dbconf.yml"
CDN_CONF_FILE = "/opt/traffic_ops/app/conf/cdn.conf"
LDAP_CONF_FILE = "/opt/traffic_ops/app/conf/ldap.conf"
USERS_CONF_FILE = "/opt/traffic_ops/install/data/json/users.json"
PROFILES_CONF_FILE = "/opt/traffic_ops/install/data/profiles/"
OPENSSL_CONF_FILE = "/opt/traffic_ops/install/data/json/openssl_configuration.json"
PARAM_CONF_FILE = "/opt/traffic_ops/install/data/json/profiles.json"
TRAFFIC_VAULT_AES_KEY_FILE = "/opt/traffic_ops/app/conf/aes.key"


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

INDENT = "\t"

class Question:
	"""
	Question represents a single question to be asked of the user, to determine a configuration
	value.

	>>> Question("question", "answer", "var")
	Question(question='question', default='answer', config_var='var', hidden=False)
	"""

	def __init__(self, question: str, default: str, config_var: str, hidden: bool = False) -> None:
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
				passwd = getpass.getpass(str(self))
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

	def serialize(self) -> Dict[str, Union[str, bool]]:
		"""Returns a serializable dictionary, suitable for converting to JSON."""
		return {self.question: self.default, "config_var": self.config_var, "hidden": self.hidden}

class User(NamedTuple):
	"""
	A User represents a user that will be inserted into the Traffic Ops database.

	Attributes
	----------
	self.username: str
		The user's username.
	self.password: str
		The user's password - IN PLAINTEXT.
	"""
	username: str
	password: str

class SSLConfig:
	"""SSLConfig bundles the options for generating new (self-signed) SSL certificates"""

	def __init__(self, gen_cert: bool, cfg_map: Dict[str, str]) -> None:
		self.gen_cert = gen_cert
		self.rsa_password = cfg_map["rsaPassword"]
		self.params = "/C={country}/ST={state}/L={locality}/O={company}/OU={org_unit}/CN={common_name}/"
		self.params = self.params.format(**cfg_map)

class CDNConfig(NamedTuple):
	"""CDNConfig holds all of the options needed to format a cdn.conf file."""

	gen_secret: bool
	num_secrets: int
	port: str
	num_workers: int
	url: str
	ldap_conf_location: str

	def generate_secret(self, conf: Dict[Any, Any]):
		"""
		Generates new secrets - if configured to do so - and adds them to the passed cdn.conf
		configuration.
		"""
		if not self.gen_secret:
			return

		if "secrets" in conf and isinstance(conf["secrets"], list):
			logging.debug("Secrets found in cdn.conf file")
		else:
			conf["secrets"] = []
			logging.debug("No secrets found in cdn.conf file")

		secrets: List[str] = conf["secrets"]
		secrets.insert(0, random_word())

		if self.num_secrets and len(secrets) > self.num_secrets:
			conf["secrets"] = secrets[:self.num_secrets - 1]

	def insert_url(self, conf: Dict[Any, Any]):
		"""
		Inserts the configured URL - if it is not an empty string - into the passed cdn.conf
		configuration, in to.base_url.
		"""
		if not self.url:
			return

		if "to" not in conf or not isinstance(conf["to"], dict):
			conf["to"] = {}
		conf["to"]["base_url"] = self.url

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
	TV_DATABASE_CONF_FILE: [
		Question("Traffic Vault Database type", "Pg", "type"),
		Question("Traffic Vault Database name", "traffic_vault", "dbname"),
		Question("Traffic Vault Database server hostname IP or FQDN", "localhost", "hostname"),
		Question("Traffic Vault Database port number", "5432", "port"),
		Question("Traffic Vault database user", "traffic_vault", "user"),
		Question("Password for Traffic Vault database user", "", "password", hidden=True)
	],
	CDN_CONF_FILE: [
		Question("Generate a new secret?", "yes", "genSecret"),
		Question("Number of secrets to keep?", "1", "keepSecrets"),
		Question("Port to serve on?", "443", "port"),
		Question("Number of workers?", "12", "workers"),
		Question("Traffic Ops url?", "http://localhost:3000", "base_url"),
		Question("ldap.conf location?", "/opt/traffic_ops/app/conf/ldap.conf", "ldap_conf_location"),
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

class ConfigEncoder(json.JSONEncoder):
	"""
	ConfigEncoder encodes a dictionary of filenames to configuration question lists as JSON.

	>>> ConfigEncoder().encode({'/test/file':[Question('question', 'default', 'cfg_var', True)]})
	'{"/test/file": [{"question": "default", "config_var": "cfg_var", "hidden": true}]}'
	"""

	def default(self, o: Any) -> Any:
		"""
		Returns a serializable representation of 'o'.

		Specifically, it does this by attempting to convert a dictionary of filenames to Question
		lists to a dictionary of filenames to lists of dictionaries of strings to strings, falling
		back on default encoding if the proper typing is not found.
		"""
		if isinstance(o, Question):
			return o.serialize()

		return json.JSONEncoder.default(self, o)

def get_config(questions: List[Question], fname: str, automatic: bool = False) -> Dict[str, str]:
	"""Asks all provided questions, or uses their defaults in automatic mode"""

	logging.info("===========%s===========", fname)

	config: Dict[str, str] = {}

	for question in questions:
		answer = question.default if automatic else question.ask()

		config[question.config_var] = answer

	return config

def generate_db_conf(qstns: List[Question], fname: str, automatic: bool, root: str) -> Dict[str, str]:
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
	with open(path, 'w+', encoding="utf-8") as conf_file:
		json.dump(db_conf, conf_file, indent=INDENT)
		print(file=conf_file)

	logging.info("Database configuration has been saved")

	return db_conf

def generate_todb_conf(fname: str, root: str, conf: Dict[str, str]):
	"""
	Generates the dbconf.yml file.

	Also writes the configuration file to the file 'fname' under the directory 'root'.
	"""

	driver = "postgres"
	if "type" not in conf:
		logging.warning("Driver type not found in todb config; using 'postgres'")
	else:
		driver = "postgres" if conf["type"] == "Pg" else conf["type"]

	path = os.path.join(root, fname.lstrip("/"))
	hostname = conf.get("hostname", "UNKNOWN")
	port = conf.get("port", "UNKNOWN")
	user = conf.get("user", "UNKNOWN")
	password = conf.get("password", "UNKNOWN")
	dbname = conf.get("dbname", "UNKNOWN")

	open_line = f"host={hostname} port={port} user={user} password={password} dbname={dbname}"
	with open(path, "w+", encoding="utf-8") as conf_file:
		print("production:", file=conf_file)
		print("    driver:", driver, file=conf_file)
		print("    open:", open_line, "sslmode=disable", file=conf_file)

def generate_ldap_conf(questions: List[Question], fname: str, automatic: bool, root: str):
	"""
	Generates the ldap.conf file by asking the questions or using default answers in auto mode.

	Also writes the configuration to the file 'fname' under the directory 'root'
	"""
	use_ldap_question = [q for q in questions if q.question == "Do you want to set up LDAP?"]
	if not use_ldap_question:
		logging.warning("Couldn't find question asking if LDAP should be set up, using default: no")
		return
	use_ldap = use_ldap_question[0].default if automatic else use_ldap_question[0].ask()

	if use_ldap.lower() not in {"y", "yes"}:
		logging.info("Not setting up ldap")
		return

	ldap_conf = get_config([q for q in questions if q is not use_ldap_question[0]], fname, automatic)
	keys = (
		"host",
		"admin_dn",
		"admin_pass",
		"search_base",
		"search_query",
		"insecure",
		"ldap_timeout_secs"
	)

	for key in keys:
		if key not in ldap_conf:
			raise ValueError(f"{key} is a required key in {fname}")

	keys_converted = {"password": "admin_pass", "hostname": "host"}
	for deprecated, key in keys_converted.items():
		if deprecated in ldap_conf and ldap_conf[key] == '':
			ldap_conf[key] = ldap_conf[deprecated]

	if not re.match(r"^\S+:\d+$", ldap_conf["host"]):
		raise ValueError(f"host in {fname} must be of form 'hostname:port'")

	path = os.path.join(root, fname.lstrip('/'))
	try:
		os.makedirs(os.path.dirname(path))
	except OSError as e:
		if e.errno == errno.EEXIST:
			pass
	with open(path, "w+", encoding="utf-8") as conf_file:
		json.dump(ldap_conf, conf_file, indent=INDENT)
		print(file=conf_file)

def hash_pass(passwd: str) -> str:
	"""
	Generates a Scrypt-based hash of the given password in a Perl-compatible format.
	It's hard-coded - like the Perl - to use 64 random bytes for the salt, n=16384,
	r=8, p=1 and dklen=64.
	"""
	n = 2 ** 14
	r_val = 8
	p_val = 1
	dklen = 64
	salt = os.urandom(dklen)
	if sys.version_info.major >= 3 and hasattr(hashlib, 'scrypt'): # Python 2.7 and CentOS 7's Python 3.6 do not include hashlib.scrypt()
		hashed = hashlib.scrypt(passwd.encode(), salt=salt, n=n, r=r_val, p=p_val, dklen=dklen)
	else:
		hashed = Scrypt(password=passwd.encode(), salt=salt, cost_factor=n, block_size_factor=r_val, parallelization_factor=p_val, key_length=dklen).derive()
	hashed_b64 = base64.standard_b64encode(hashed).decode()
	salt_b64 = base64.standard_b64encode(salt).decode()

	return f"SCRYPT:{n}:{r_val}:{p_val}:{salt_b64}:{hashed_b64}"

class Scrypt:
	"""
	Implements SCRYPT encryption based on the configuration given at object
	construction.
	"""
	def __init__(self, password: bytes, salt: bytes, cost_factor: int, block_size_factor: int, parallelization_factor: int, key_length: int):
		self.password = password
		self.salt = salt
		self.cost_factor = cost_factor
		self.block_size_factor = block_size_factor
		self.parallelization_factor = parallelization_factor
		self.key_length = key_length
		self.block_unit = 32 * self.block_size_factor

	def derive(self) -> bytes:
		"""
		Derives an encrypted bytestring representative of the password given at
		initialization.
		"""
		salt_length = 2 ** 7 * self.block_size_factor * self.parallelization_factor
		pack_format = f"<{'L' * int(salt_length / 4)}"  # `<` means `little-endian` and `L` means `unsigned long`
		salt = hashlib.pbkdf2_hmac('sha256', password=self.password, salt=self.salt, iterations=1, dklen=salt_length)
		block = list(unpack(pack_format, salt))
		block = self.ROMix(block)
		salt = pack(pack_format, *block)
		key = hashlib.pbkdf2_hmac('sha256', password=self.password, salt=salt, iterations=1, dklen=self.key_length)
		return key

	def ROMix(self, block: List[int]) -> List[int]:
		xored_block = [0] * len(block)
		variations: List[List[int]] = [[]] * self.cost_factor
		variations[0] = block
		index = 1
		while index < self.cost_factor:
			variations[index] = self.block_mix(variations[index - 1])
			index += 1
		block = self.block_mix(variations[-1])
		for _ in variations:
			variation_index = block[self.block_unit - 16] % self.cost_factor
			variation = variations[variation_index]
			for index, _unused in enumerate(xored_block):
				xored_block[index] = block[index] ^ variation[index]
			block = self.block_mix(xored_block)
		return block

	def block_mix(self, previous_block: List[int]) -> List[int]:
		block = previous_block.copy()
		x_length = 16  # x is the list of numbers within `block` that we mix
		copy_index = self.block_unit - x_length
		x = previous_block[copy_index:copy_index + x_length]
		octet_index = 0
		block_xor_index = 0
		while octet_index < 2 * self.block_size_factor:
			for index, _ in enumerate(x):
				x[index] ^= previous_block[block_xor_index + index]
			block_xor_index += x_length
			self.salsa20(x)
			block_offset = (int(octet_index / 2) + octet_index % 2 * self.block_size_factor) * x_length
			block[block_offset:block_offset + x_length] = x
			octet_index += 1
		return block

	def salsa20(self, block: List[int]):
		x = block.copy()
		for _ in range(4):
			# These bit shifting operations could be condensed into a single line of list comprehensions,
			# but there is a >3x performance benefit from writing it out explicitly.
			bits = x[0] + x[12] & 0xffffffff
			x[4] ^= bits << 7 | bits >> 32 - 7
			bits = x[4] + x[0] & 0xffffffff
			x[8] ^= bits << 9 | bits >> 32 - 9
			bits = x[8] + x[4] & 0xffffffff
			x[12] ^= bits << 13 | bits >> 32 - 13
			bits = x[12] + x[8] & 0xffffffff
			x[0] ^= bits << 18 | bits >> 32 - 18
			bits = x[5] + x[1] & 0xffffffff
			x[9] ^= bits << 7 | bits >> 32 - 7
			bits = x[9] + x[5] & 0xffffffff
			x[13] ^= bits << 9 | bits >> 32 - 9
			bits = x[13] + x[9] & 0xffffffff
			x[1] ^= bits << 13 | bits >> 32 - 13
			bits = x[1] + x[13] & 0xffffffff
			x[5] ^= bits << 18 | bits >> 32 - 18
			bits = x[10] + x[6] & 0xffffffff
			x[14] ^= bits << 7 | bits >> 32 - 7
			bits = x[14] + x[10] & 0xffffffff
			x[2] ^= bits << 9 | bits >> 32 - 9
			bits = x[2] + x[14] & 0xffffffff
			x[6] ^= bits << 13 | bits >> 32 - 13
			bits = x[6] + x[2] & 0xffffffff
			x[10] ^= bits << 18 | bits >> 32 - 18
			bits = x[15] + x[11] & 0xffffffff
			x[3] ^= bits << 7 | bits >> 32 - 7
			bits = x[3] + x[15] & 0xffffffff
			x[7] ^= bits << 9 | bits >> 32 - 9
			bits = x[7] + x[3] & 0xffffffff
			x[11] ^= bits << 13 | bits >> 32 - 13
			bits = x[11] + x[7] & 0xffffffff
			x[15] ^= bits << 18 | bits >> 32 - 18
			bits = x[0] + x[3] & 0xffffffff
			x[1] ^= bits << 7 | bits >> 32 - 7
			bits = x[1] + x[0] & 0xffffffff
			x[2] ^= bits << 9 | bits >> 32 - 9
			bits = x[2] + x[1] & 0xffffffff
			x[3] ^= bits << 13 | bits >> 32 - 13
			bits = x[3] + x[2] & 0xffffffff
			x[0] ^= bits << 18 | bits >> 32 - 18
			bits = x[5] + x[4] & 0xffffffff
			x[6] ^= bits << 7 | bits >> 32 - 7
			bits = x[6] + x[5] & 0xffffffff
			x[7] ^= bits << 9 | bits >> 32 - 9
			bits = x[7] + x[6] & 0xffffffff
			x[4] ^= bits << 13 | bits >> 32 - 13
			bits = x[4] + x[7] & 0xffffffff
			x[5] ^= bits << 18 | bits >> 32 - 18
			bits = x[10] + x[9] & 0xffffffff
			x[11] ^= bits << 7 | bits >> 32 - 7
			bits = x[11] + x[10] & 0xffffffff
			x[8] ^= bits << 9 | bits >> 32 - 9
			bits = x[8] + x[11] & 0xffffffff
			x[9] ^= bits << 13 | bits >> 32 - 13
			bits = x[9] + x[8] & 0xffffffff
			x[10] ^= bits << 18 | bits >> 32 - 18
			bits = x[15] + x[14] & 0xffffffff
			x[12] ^= bits << 7 | bits >> 32 - 7
			bits = x[12] + x[15] & 0xffffffff
			x[13] ^= bits << 9 | bits >> 32 - 9
			bits = x[13] + x[12] & 0xffffffff
			x[14] ^= bits << 13 | bits >> 32 - 13
			bits = x[14] + x[13] & 0xffffffff
			x[15] ^= bits << 18 | bits >> 32 - 18

		for index in range(16):
			block[index] = block[index] + x[index] & 0xffffffff

def generate_users_conf(qstns: List[Question], fname: str, auto: bool, root: str) -> User:
	"""
	Generates a users.json file from the given questions and returns a User containing the same
	information.
	"""
	config = get_config(qstns, fname, auto)

	if "tmAdminUser" not in config or "tmAdminPw" not in config:
		raise ValueError(f"{fname} must include 'tmAdminUser' and 'tmAdminPw'")

	hashed_pass = hash_pass(config["tmAdminPw"])

	path = os.path.join(root, fname.lstrip('/'))
	with open(path, "w+", encoding="utf-8") as conf_file:
		json.dump({"username": config["tmAdminUser"], "password": hashed_pass}, conf_file, indent=INDENT)
		print(file=conf_file)

	return User(config["tmAdminUser"], config["tmAdminPw"])

def generate_profiles_dir(questions: List[Question]):
	"""
	I truly have no idea what's going on here. This is what the Perl did, so I
	copied it. It does nothing. Literally nothing.
	"""
	#pylint:disable=unused-variable
	user_in = questions
	#pylint:enable=unused-variable

def generate_openssl_conf(questions: List[Question], fname: str, auto: bool) -> SSLConfig:
	"""
	Constructs an SSLConfig by asking the passed questions, or using their default answers if in
	auto mode.
	"""
	cfg_map = get_config(questions, fname, auto)
	if "genCert" not in cfg_map:
		raise ValueError("missing 'genCert' key")

	gen_cert = cfg_map["genCert"].lower() in {"y", "yes"}

	return SSLConfig(gen_cert, cfg_map)

def generate_param_conf(qstns: List[Question], fname: str, auto: bool, root: str) -> Dict[str, str]:
	"""
	Generates a profiles.json by asking the passed questions, or using their default answers in auto
	mode.

	Also writes the file to 'fname' in the directory 'root'.
	"""
	conf = get_config(qstns, fname, auto)

	path = os.path.join(root, fname.lstrip("/"))
	with open(path, "w+", encoding="utf-8") as conf_file:
		json.dump(conf, conf_file, indent=INDENT)
		print(file=conf_file)

	return conf

def sanity_check_config(cfg: Dict[str, List[Question]], automatic: bool) -> int:
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

def unmarshal_config(dct: Dict[str, Any]) -> Dict[str, List[Question]]:
	"""
	Reads in a raw parsed configuration file and returns the resulting configuration.

	>>> unmarshal_config({"test": [{"Do the thing?": "yes", "config_var": "thing"}]})
	{'test': [Question(question='Do the thing?', default='yes', config_var='thing', hidden=False)]}
	>>> unmarshal_config({"test": [{"foo": "", "config_var": "bar", "hidden": True}]})
	{'test': [Question(question='foo', default='', config_var='bar', hidden=True)]}
	"""
	ret: Dict[str, List[Question]] = {}
	for file, questions in dct.items():
		if not isinstance(questions, list):
			raise ValueError(f"file '{file}' has malformed questions")

		qstns: List[Question] = []
		qstn: Any
		for qstn in questions:
			if not isinstance(qstn, dict):
				raise ValueError(f"file '{file}' has a malformed question ({qstn})")
			question: Any
			try:
				question = next(key for key in qstn.keys() if key not in {"hidden", "config_var"})
			except StopIteration:
				raise ValueError(f"question in '{file}' has no question/answer properties ({qstn})")

			answer: Any = qstn[question]
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
			cfg_var: Any = qstn["config_var"]
			if not isinstance(cfg_var, str):
				raise ValueError(f"question in '{file}' has malformed 'config_var' property ({cfg_var})")
			del qstn["config_var"]

			if qstn:
				logging.warning("Found unknown extra properties in question in '%s' (%r)", file, qstn.keys())

			qstns.append(Question(question, answer, cfg_var, hidden=hidden))
		ret[file] = qstns

	return ret

def write_encryption_key(aes_key_location: str):
	"""
	Creates an AES encryption key for the postgres traffic vault backend

	:param aes_key_location: Denotes the location of the aes encryption key file
	:returns: None
	"""

	args = (
		"rand",
		"-out",
		aes_key_location,
		"-base64",
		"32"
	)
	if not exec_openssl(f"Generating an AES encryption key to {aes_key_location}", *args):
		logging.debug("AES key generation failed")
		raise OSError("failed to generate AES key")

def exec_openssl(description: str, *cmd_args: str) -> bool:
	"""
	Executes openssl with the supplied command-line arguments.

	:param description: Describes the operation taking place for logging purposes.
	:returns: Whether or not the execution succeeded, success being defined by an exit code of zero
	"""
	logging.info(description)

	cmd = ("/usr/bin/openssl",) + cmd_args

	while True:
		proc = subprocess.Popen(
			cmd,
			stderr=subprocess.PIPE,
			stdout=subprocess.PIPE,
			universal_newlines=True,
		)
		proc.wait()
		if proc.returncode == 0:
			return True

		output = proc.communicate()
		logging.debug("openssl exec failed with code %s; stderr: %s", proc.returncode, output[1])
		while True:
			ans = input(f"{description} failed. Try again (y/n) [y]: ")
			if not ans or ans.lower().startswith("n"):
				return False
			if ans.lower().startswith("y"):
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
	os.chown(keypath, pwd.getpwnam(ops_user).pw_uid, grp.getgrnam(ops_group).gr_gid)

	logging.info("The private key has been installed")
	logging.info("Installing self signed certificate")

	certpath = os.path.join(root, 'etc/pki/tls/certs/localhost.crt')
	shutil.copy("server.crt", certpath)
	os.chmod(certpath, stat.S_IRUSR | stat.S_IWUSR)
	os.chown(certpath, pwd.getpwnam(ops_user).pw_uid, grp.getgrnam(ops_group).gr_gid)

	logging.info("Saving the self signed csr")

	csrpath = os.path.join(root, 'etc/pki/tls/certs/localhost.csr')
	shutil.copy("server.csr", csrpath)
	os.chmod(csrpath, stat.S_IRUSR | stat.S_IWUSR | stat.S_IRGRP | stat.S_IWGRP | stat.S_IROTH)
	os.chown(csrpath, pwd.getpwnam(ops_user).pw_uid, grp.getgrnam(ops_group).gr_gid)

	log_msg = """
        The self signed certificate has now been installed.

        You may obtain a certificate signed by a Certificate Authority using the
        server.csr file saved in the current directory.  Once you have obtained
        a signed certificate, copy it to %s and
        restart Traffic Ops."""
	logging.info(log_msg, certpath)

	cdn_conf_path = os.path.join(root, "opt/traffic_ops/app/conf/cdn.conf")

	try:
		with open(cdn_conf_path, encoding="utf-8") as conf_file:
			cdn_conf = json.load(conf_file)
	except (OSError, ValueError) as e:
		raise OSError(f"reading {cdn_conf_path}: {e}") from e

	if (
		not isinstance(cdn_conf, dict) or
		"traffic_ops_golang" not in cdn_conf or
		not isinstance(cdn_conf["traffic_ops_golang"], dict)
	):
		logging.critical("Malformed %s; improper object and/or missing 'traffic_ops_golang' key", cdn_conf_path)
		return 1

	to_golang: Any = cdn_conf["traffic_ops_golang"]
	if (
		"cert" not in to_golang or
		not isinstance(to_golang["cert"], str)
	):
		log_msg = """	The "cert" portion of %s is missing from %s
	Please ensure it contains the same structure as the one originally installed"""
		logging.error(log_msg, cdn_conf_path, cdn_conf_path)
		return 1

	if (
		"key" not in to_golang or
		not isinstance(to_golang["key"], str)
	):
		log_msg = """	The "key" portion of %s is missing from %s
	Please ensure it contains the same structure as the one originally installed"""
		logging.error(log_msg, cdn_conf_path, cdn_conf_path)
		return 1

	return 0

def random_word(length: int = 12) -> str:
	"""
	Returns a randomly generated string 'length' characters long containing only word
	characters ([a-zA-Z0-9_]).
	"""
	word_chars = string.ascii_letters + string.digits + "_"
	return "".join(random.choice(word_chars) for _ in range(length))

def generate_cdn_conf(questions: List[Question], fname: str, automatic: bool, root: str) -> bool:
	"""
	Generates some properties of a cdn.conf file based on the passed questions.

	This modifies or writes the file 'fname' under the directory 'root'.
	:returns: A boolean value denoting whether or not a postgres traffic vault backend is configured.
	"""
	cdn_conf = get_config(questions, fname, automatic)

	if "genSecret" not in cdn_conf:
		raise ValueError("missing 'genSecret' config_var")

	gen_secret = cdn_conf["genSecret"].lower() in {'y', 'yes'}

	try:
		num_secrets = int(cdn_conf["keepSecrets"])
	except KeyError as e:
		raise ValueError("missing 'keepSecrets' config_var") from e
	except ValueError as e:
		raise ValueError(f"invalid 'keepSecrets' config_var value: {e}") from e

	try:
		port = cdn_conf["port"]
		int(port)
	except KeyError as e:
		raise ValueError("missing 'port' config_var") from e
	except ValueError as e:
		raise ValueError(f"invalid 'port' config_var value: {e}") from e

	try:
		workers = int(cdn_conf["workers"])
	except KeyError as e:
		raise ValueError("missing 'workers' config_var") from e
	except ValueError as e:
		raise ValueError(f"invalid 'workers' config_var value: {e}") from e

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
	existing_conf: Dict[Any, Any] = {}
	if os.path.isfile(path):
		with open(path, encoding="utf-8") as conf_file:
			try:
				existing_conf = json.load(conf_file)
			except ValueError as e:
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
	traffic_vault_backend = "postgres"
	tv_aes_key_location = os.path.join(root, TRAFFIC_VAULT_AES_KEY_FILE.lstrip('/'))

	with open(path, "w+", encoding="utf-8") as conf_file:
		json.dump(existing_conf, conf_file, indent=INDENT)
		print(file=conf_file)
	logging.info("CDN configuration has been saved")
	try:
		traffic_vault_backend: Any = existing_conf["traffic_ops_golang"]["traffic_vault_backend"]
	except KeyError as e:
		logging.warning("no traffic vault backend configured, using default postgres")

	if traffic_vault_backend == "postgres":
		try:
			tv_aes_key_location: Any = existing_conf["traffic_ops_golang"]["traffic_vault_config"]["aes_key_location"]
		except KeyError as e:
			logging.warning("no traffic vault aes encryption key location specified, using default %s", TRAFFIC_VAULT_AES_KEY_FILE)
		write_encryption_key(tv_aes_key_location)

	return traffic_vault_backend == "postgres"

def db_connection_string(dbconf: Dict[str, Any]) -> str:
	"""
	Constructs a database connection string from the passed configuration object.
	"""
	user = dbconf["user"]
	password = dbconf["password"]
	db_name = "traffic_ops" if dbconf["type"] == "Pg" else dbconf["type"]
	hostname = dbconf["hostname"]
	port = dbconf["port"]
	return f"postgresql://{user}:{password}@{hostname}:{port}/{db_name}"

def exec_psql(conn_str: str, query: str, **args: str) -> str:
	"""
	Executes SQL queries by forking and exec-ing '/usr/bin/psql'.

	:param conn_str: A "connection string" that defines the postgresql resource in the format
	{schema}://{user}:{password}@{host or IP}:{port}/{database}
	:param query: The query to be run. It can actually be a script containing multiple queries.
	:returns: The comma-separated columns of each line-delimited row of the results of the query.
	"""
	cmd = ["/usr/bin/psql", "--tuples-only", "-d", conn_str, "-c", query] + list(args.values())
	proc = subprocess.Popen(
		cmd,
		stderr=subprocess.PIPE,
		stdout=subprocess.PIPE,
		universal_newlines=True,
	)
	proc.wait()
	output = proc.communicate()
	if proc.returncode != 0:
		logging.debug("psql exec failed; stderr: %s\n\tstdout: %s", output[1], output[0])
		raise OSError("failed to execute database query")
	return output[0].strip()

def invoke_db_admin_pl(action: str, root: str, tv: bool):
	"""
	Executes admin with the given action, and looks for it from the given root directory.
	"""
	path = os.path.join(root, "opt/traffic_ops/app")
	# This is a workaround for admin using hard-coded relative paths. That
	# should be fixed at some point, IMO, but for now this works.
	os.chdir(path)
	cmd = [os.path.join(path, "db/admin"), "--env=production", action]
	if tv:
		cmd = [os.path.join(path, "db/admin"), "--trafficvault","--env=production", action]
	proc = subprocess.Popen(
		cmd,
		stderr=subprocess.PIPE,
		stdout=subprocess.PIPE,
		universal_newlines=True,
	)
	output = proc.communicate()
	if proc.returncode != 0:
		logging.debug("admin exec failed; stderr: %s\n\tstdout: %s", output[1], output[0])
		raise OSError(f"Database {action} failed")
	logging.info("Database %s succeeded", action)

def setup_database_data(
	conn_str: str,
	user: User,
	param_conf: Dict[str, str],
	root: str,
	postgresTV: bool
):
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
		invoke_db_admin_pl("load_schema", root, False)

	invoke_db_admin_pl("migrate", root, False)
	invoke_db_admin_pl("seed", root, False)
	invoke_db_admin_pl("patch", root, False)

	if postgresTV:
		invoke_db_admin_pl("create_user", root, True)
		invoke_db_admin_pl("createdb", root, True)
		invoke_db_admin_pl("load_schema", root, True)
		invoke_db_admin_pl("migrate", root, True)

	hashed_pass = hash_pass(user.password)
	insert_admin_query = f'''
		INSERT INTO tm_user (username, tenant_id, role, local_passwd, confirm_local_passwd)
		VALUES (
			'{user.username}',
			(SELECT id FROM tenant WHERE name = 'root'),
			(SELECT id FROM role WHERE name = 'admin'),
			'{hashed_pass}',
			'{hashed_pass}'
		)
		ON CONFLICT (username) DO NOTHING;
	'''
	_ = exec_psql(conn_str, insert_admin_query)

	logging.info("=========== Setting up cdn")
	insert_cdn_query = "\n\t-- global parameters" + f'''
		INSERT INTO cdn (name, domain_name, dnssec_enabled)
		VALUES ('{param_conf["cdn_name"]}', '{param_conf["dns_subdomain"]}', false)
		ON CONFLICT DO NOTHING;
	'''
	logging.info("\n%s", insert_cdn_query)
	_ = exec_psql(conn_str, insert_cdn_query)

	tm_url = param_conf["tm.url"]

	logging.info("=========== Setting up parameters")
	insert_parameters_query = "\n\t-- global parameters" + f'''
		INSERT INTO parameter (name, config_file, value)
		VALUES ('tm.url', 'global', '{tm_url}'),
			('tm.infourl', 'global', '{tm_url}/doc'),
		-- CRConfic.json parameters
			('geolocation.polling.url', 'CRConfig.json', '{tm_url}/routing/GeoLite2-City.mmdb.gz'),
			('geolocation6.polling.url', 'CRConfig.json', '{tm_url}/routing/GeoLiteCityv6.dat.jz')
		ON CONFLICT (name, config_file, value) DO NOTHING;
	'''
	logging.info("\n%s", insert_parameters_query)
	_ = exec_psql(conn_str, insert_parameters_query)

	logging.info("\n=========== Setting up profiles")
	insert_profiles_query = "\n\t-- global parameters" + f'''
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
	'''
	logging.info("\n%s", insert_profiles_query)
	_ = exec_psql(conn_str, insert_cdn_query)

def main(
	automatic: bool,
	debug: bool,
	defaults: Optional[str],
	cfile: str,
	root_dir: str,
	ops_user: str,
	ops_group: str,
	no_restart_to: bool,
	no_database: bool,
) -> int:
	"""
	Runs the main routine given the parsed arguments as input.
	"""
	postgres_tv = False
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
					with open(defaults, "w", encoding="utf-8") as dump_file:
						json.dump(DEFAULTS, dump_file, indent=INDENT)
				except OSError as e:
					logging.critical("Writing output: %s", e)
					return 1
			else:
				json.dump(DEFAULTS, sys.stdout, cls=ConfigEncoder, indent=INDENT)
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
			with open(cfile, encoding="utf-8") as conf_file:
				user_input = unmarshal_config(json.load(conf_file))
			diffs = sanity_check_config(user_input, automatic)
			logging.info(
			"File sanity check complete - found %s difference%s",
			diffs,
			'' if diffs == 1 else 's'
			)
		except (OSError, ValueError) as e:
			logging.critical("Reading in input file '%s': %s", cfile, e)
			return 1

	try:
		dbconf = generate_db_conf(user_input[DATABASE_CONF_FILE], DATABASE_CONF_FILE, automatic, root_dir)
		generate_todb_conf(DB_CONF_FILE, root_dir, dbconf)
		# the new "/opt/traffic_ops/app/conf/production/tv.conf" section for Traffic Vault PostgreSQL backend is optional
		if TV_DATABASE_CONF_FILE in user_input:
			tv_dbconf = generate_db_conf(user_input[TV_DATABASE_CONF_FILE], TV_DATABASE_CONF_FILE, automatic, root_dir)
			generate_todb_conf(TV_DB_CONF_FILE, root_dir, tv_dbconf)
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
			with open(postinstall_cfg, "w+", encoding="utf-8") as conf_file:
				print("{}", file=conf_file)
	except OSError as e:
		logging.critical("Writing configuration: %s", e)
		return 1
	except ValueError as e:
		logging.critical("Generating configuration: %s", e)
		return 1

	try:
		cert_code = setup_certificates(opensslconf, root_dir, ops_user, ops_group)
		if cert_code:
			return cert_code
	except OSError as e:
		logging.critical("Setting up SSL Certificates: %s", e)
		return 1

	try:
		postgres_tv = generate_cdn_conf(user_input[CDN_CONF_FILE], CDN_CONF_FILE, automatic, root_dir)
	except OSError as e:
		logging.critical("Generating cdn.conf: %s", e)
		return 1

	if not no_database:
		try:
			conn_str = db_connection_string(dbconf)
		except KeyError as e:
			logging.error("Missing database connection variable: %s", e)
			logging.error(
				"Can't connect to the database.  " \
				"Make sure traffic_ops database and username exist in postgres."
			)
			return 1

		if not os.path.isfile("/usr/bin/psql") or not os.access("/usr/bin/psql", os.X_OK):
			logging.critical("psql is not installed, please install it to continue with database setup")
			return 1

		def db_connect_failed():
			logging.error("Failed to set up database: %s", e)
			logging.error(
				"Can't connect to the database.  "
				"Make sure traffic_ops database and username exist in postgres."
			)

		try:
			setup_database_data(conn_str, admin_conf, paramconf, root_dir, postgres_tv)
		except (subprocess.CalledProcessError, OSError) as e:
			db_connect_failed()
			return 1
		except subprocess.SubprocessError as e:
			db_connect_failed()
			return 1


	if not no_restart_to:
		logging.info("Starting Traffic Ops")
		try:
			cmd = ["/sbin/service", "traffic_ops", "restart"]
			with subprocess.Popen(
				cmd,
				stderr=subprocess.PIPE,
				stdout=subprocess.PIPE,
				universal_newlines=True,
			) as proc:
				try:
					if proc.wait():
						raise subprocess.CalledProcessError(proc.returncode, cmd)
				except (subprocess.CalledProcessError) as e:
					output = proc.communicate()
					logging.critical("Failed to restart Traffic Ops, return code %s: %s", e.returncode, e)
					logging.debug("stderr: %s\n\tstdout: %s", output[1], output[0])
					return 1
		except OSError as e:
			logging.critical("Failed to restart Traffic Ops: unknown error occurred: %s", e)
			return 1
		# Perl didn't actually do any "waiting" before reporting success, so
		# neither do we
		logging.info("Waiting for Traffic Ops to restart")
	else:
		logging.info("Skipping Traffic Ops restart")
	logging.info("Success! Postinstall complete.")

	return 0

if __name__ == '__main__':
	logging.basicConfig(stream=sys.stdout)

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
	PARSER.add_argument(
		"-cfile",
		help=argparse.SUPPRESS,
		type=str,
		default=None,
		dest="legacy_cfile"
	)
	PARSER.add_argument("--debug", help="Enables verbose output", action="store_true")
	PARSER.add_argument("-debug", help=argparse.SUPPRESS, dest="legacy_debug", action="store_true")
	PARSER.add_argument(
		"--defaults",
		help="Writes out a configuration file with defaults which can be used as input",
		type=str,
		nargs="?",
		default=None,
		const=""
	)
	PARSER.add_argument(
		"-defaults",
		help=argparse.SUPPRESS,
		type=str,
		nargs="?",
		default=None,
		const="",
		dest="legacy_defaults"
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

	USED_LEGACY_ARGS = False
	DEFAULTS_ARG = None
	if ARGS.legacy_defaults:
		if ARGS.defaults:
			logging.error("cannot specify both '--defaults' and '-defaults'")
			sys.exit(1)
		USED_LEGACY_ARGS = True
		DEFAULTS_ARG = ARGS.legacy_defaults
	else:
		DEFAULTS_ARG = ARGS.defaults

	DEBUG = False
	if ARGS.legacy_debug:
		if ARGS.debug:
			logging.error("cannot specify both '--debug' and '-debug'")
			sys.exit(1)
		USED_LEGACY_ARGS = True
		DEBUG = ARGS.legacy_debug
	else:
		DEBUG = ARGS.debug

	CFILE = None
	if ARGS.legacy_cfile:
		if ARGS.cfile:
			logging.error("cannot specify both '--cfile' and '-cfile'")
			sys.exit(1)
		USED_LEGACY_ARGS = True
		CFILE = ARGS.legacy_cfile
	else:
		CFILE = ARGS.cfile

	if not ARGS.no_root and os.getuid() != 0:
		logging.error("You must run this script as the root user")
		logging.shutdown()
		sys.exit(1)

	if USED_LEGACY_ARGS:
		logging.warning(
			"passing long options with a single '-' is deprecated, please use '--' in the future"
		)

	try:
		EXIT_CODE = main(
			ARGS.automatic,
			DEBUG,
			DEFAULTS_ARG,
			CFILE,
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
