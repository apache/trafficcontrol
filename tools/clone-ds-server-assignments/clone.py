#!/usr/bin/env python3

"""
.. _clone.py:

.. program:: clone.py

This module is a script that can be used to copy server assignments from one
Delivery Service to another.

Usage
=====
``clone.py [-h] [-k] [-t TO_URL] [-u USERNAME] [-p PASSWORD] [-a API_VERSION] [-r] FROM TO``

Arguments and Flags
===================
.. option:: FROM

	The XMLID of the Delivery Service from which to copy server assignments.

.. option:: TO

	The XMLID of the Delivery Service to which to copy server assignments.

.. option:: -h, --help

	Print usage information and exit.

.. option:: -a API_VERSION, --api-version API_VERSION

	Specifies the version of the Traffic Ops API that will be used for all
	requests. (Default: 2.0)

.. option:: -k, --insecure

	Do not verify SSL certificates - typically useful for making requests to
	development or testing servers as they frequently have self-signed
	certificates. (Default: false)

.. option:: -r, --replace

	Replace any and all existing server assignments on the Delivery Service
	specified by TO, rather than simply adding the ones from FROM to what it
	already has. (Default: false)

.. option:: --to-url URL

	The :abbr:`FQDN (Fully Qualified Domain Name)` and optionally the port and
	scheme of the Traffic Ops server. This will override :envvar:`TO_URL`. The
	format is the same as for :envvar:`TO_URL`. (Default: uses the value of
	:envvar:`TO_URL`)

	.. note:: All requests will use HTTPS - even if this URL is given with an
		``http://`` scheme.

.. option:: --to-password PASSWORD

	The password to use when authenticating to Traffic Ops. Overrides
	:envvar:`TO_PASSWORD`. (Default: uses the value of :envvar:`TO_PASSWORD`)

.. option:: --to-user USERNAME

	The username to use when connecting to Traffic Ops. Overrides
	:envvar:`TO_USER`. (Default: uses the value of :envvar:`TO_USER`)

Environment Variables
=====================
If defined, :program:`clone.py` will use the :envvar:`TO_URL`,
:envvar:`TO_USER`, and :envvar`TO_PASSWORD` environment variables to define
its connection to and authentication with the Traffic Ops server. Typically,
setting these is easier than using the long options :option:`--to-url`,
:option:`--to-user`, and :option:`--to-password` on every invocation.

Module Reference
================

"""

from argparse import ArgumentParser, ArgumentDefaultsHelpFormatter
import logging
import os
import re
import sys
import typing

from trafficops import TOSession, restapi

# Prevent logging from external libraries from bubbling up to stderr
logging.disable(logging.CRITICAL)

class Config(typing.NamedTuple):
	"""
	Holds configuration options for running the script.
	"""
	api_version: str
	from_DS: str
	to_DS: str
	insecure: bool
	url: str
	port: int
	username: str
	password: str
	replace: bool

def main(cfg: Config) -> int:
	"""
	The main routine, entrypoint for the script.

	:param cfg: A configuration for the script to use.
	:returns: An exit code for the script.
	"""
	with TOSession(cfg.url, api_version=cfg.api_version, verify_cert=not cfg.insecure, host_port=cfg.port) as to:
		to.login(cfg.username, cfg.password)
		f = to.get_deliveryservices(query_params={"xmlId": cfg.from_DS})[0]
		if isinstance(f, list):
			f = f[0]
		t = to.get_deliveryservices(query_params={"xmlId": cfg.to_DS})[0]
		if isinstance(t, list):
			t = t[0]

		from_assignments = to.get_deliveryservice_servers(delivery_service_id=f["id"])[0]
		servers = [s["id"] for s in from_assignments]
		payload = {
			"dsId": t["id"],
			"replace": cfg.replace,
			"servers": servers
		}
		to.assign_deliveryservice_servers_by_ids(data=payload)

	return 0

if __name__ == "__main__":
	parser = ArgumentParser(
		description="Clones server assignments from one Delivery Service to another",
		formatter_class=ArgumentDefaultsHelpFormatter
	)
	parser.add_argument("FROM", type=str, help="The DS from which to copy server assignments")
	parser.add_argument("TO", type=str, help="The DS to which to copy server assignments")
	parser.add_argument("-k", "--insecure", action="store_true", help="Skip TO server certificate verification")
	parser.add_argument(
		"-t", "--to-url",
		type=str,
		help="The URL of the TO instance to use - overrides the TO_URL environment variable"
	)
	parser.add_argument(
		"-u", "--username",
		type=str,
		help="The username to use when authenticating with TO - overrides the TO_USER environment variable"
	)
	parser.add_argument(
		"-p", "--password",
		type=str,
		help="The password to use when authenticating with TO - overrides the TO_PASSWORD environment variable"
	)
	parser.add_argument(
		"-a", "--api-version",
		type=str,
		help="Specify an API version to use",
		default="2.0"
	)
	parser.add_argument(
		"-r", "--replace",
		action="store_true",
		help="replace existing server assignments on the 'TO' Delivery Service"
	)

	args = parser.parse_args()

	url = os.environ["TO_URL"]
	if args.to_url:
		url = args.to_url
	if not url:
		print("A Traffic Ops URL must be specified", file=sys.stderr)
		sys.exit(1)

	url = url.lower()
	url_pattern = re.compile(r'''^(?:https?://)?([a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*)(:[0-9]+)?''')
	match = url_pattern.match(url)
	port = None
	if match:
		groups = match.groups()
		url = groups[0]
		if groups[1]:
			port = int(groups[1][1:])

	user = os.environ["TO_USER"]
	if args.username:
		user = args.username
	if not user:
		print("A Traffic Ops username must be specified", file=sys.stderr)
		sys.exit(1)

	passwd = os.environ["TO_PASSWORD"]
	if args.password:
		passwd = args.password
	if not passwd:
		print("A Traffic Ops password must be specified", file=sys.stderr)
		sys.exit(1)

	try:
		cfg = Config(
			args.api_version,
			args.FROM,
			args.TO,
			args.insecure,
			url,
			port or 443,
			user,
			passwd,
			args.replace
		)
		sys.exit(main(cfg))
	except KeyboardInterrupt:
		sys.exit(2)
	except restapi.OperationError as e:
		print("Copying assignments failed:", e, file=sys.stderr)
		sys.exit(3)
