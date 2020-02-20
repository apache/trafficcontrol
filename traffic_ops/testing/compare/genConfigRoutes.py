#!/usr/bin/env python3
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

"""
This script is meant to generate a list of Traffic Ops API routes that point to configuration files
for cache servers. It verifies that servers of the same name both exist and have the same routes.

Note that the Python Apache-TrafficControl client package must be installed to use this script.
"""

import argparse
import logging
import os
import random
import time
import typing
import sys

from trafficops.tosession import TOSession
from trafficops.restapi import LoginError, OperationError, InvalidJSONError

random.seed(time.time())

#: A format specifier for logging output. Propagates to all imported modules.
LOG_FMT = "%(levelname)s: %(asctime)s line %(lineno)d in %(module)s.%(funcName)s: %(message)s"

__version__ = "2.0.0-1"

def getConfigRoutesForServers(servers:typing.List[dict], inst:TOSession) \
                                                               -> typing.Generator[str, None, None]:
	"""
	Generates a list of routes to the config files for a given set of servers and a given traffic
	ops instance

	:param servers: a list of server objects
	:param inst: A valid, authenticated, and connected Traffic Ops instance
	:returns: A list of routes to config files for the ``servers``. These will be relative to the
		url of the ``inst``
	"""
	for server in servers:
		try:
			yield "/api/1.3/servers/%s/configfiles/ats" % server.hostName
			for file in inst.get_server_config_files(host_name=server.hostName)[0].configFiles:
				if "apiUri" in file:
					yield file.apiUri
				else:
					logging.info("config file %s for server %s has non-API URI - skipping",
				                       file.location, server.hostName)
		except (AttributeError, UnicodeError, IndexError, KeyError, InvalidJSONError, OperationError) as e:
			logging.debug("%r", e, exc_info=True, stack_info=True)
			logging.error("Invalid API response for server %s config files: %s", server.hostName, e)

def getCRConfigs(A:TOSession, B:TOSession) -> typing.Generator[str, None, None]:
	"""
	Generates a list of routes to CRConfig files for all CDNs present in both A and B

	:param A: The first Traffic Ops instance
	:param B: The second Traffic Ops instance
	:returns: A list of routes to CRConfig files
	"""
	try:
		Acdns = A.get_cdns()[0]
		Bcdns = B.get_cdns()[0]
	except (UnicodeError, IndexError, KeyError, InvalidJSONError, OperationError) as e:
		logging.debug("%r", e, exc_info=True, stack_info=True)
		logging.critical("Unable to get CDN lists: %s", e)
		return

	cdns = {c.name for c in Acdns}.intersection({c.name for c in Bcdns})

	if not cdns:
		logging.error("The two instances have NO CDNs in common! This almost certainly means that "\
		              "you're not doing what you want to do")
		return

	for cdn in cdns:
		yield "/CRConfig-Snapshots/%s/CRConfig.json" % cdn
		yield "/api/2.0/cdns/%s/snapshot" % cdn
		yield "/api/2.0/cdns/%s/snapshot/new" % cdn


def consolidateVariables(kwargs:argparse.Namespace) -> typing.Tuple[str, str,
                                                         typing.Tuple[str, str], typing.Tuple[str]]:
	"""
	Consolidates the arguments passed on the command line with the ones in the environment

	:param kwargs: The arguments passed on the command line
	:returns: In order: the reference Traffic Ops URL, the testing Traffic Ops URL, the login
		information for the reference instance, and the login information for the testing instance
	:raises ValueError: if a required variable is not defined
	"""
	instanceA = kwargs.refURL if kwargs.refURL else os.environ.get("TO_URL", None)
	if instanceA is None:
		logging.critical("Must specify the URL of the reference instance!")
		raise ValueError()

	instanceB = kwargs.testURL if kwargs.testURL else os.environ.get("TEST_URL", None)
	if instanceB is None:
		logging.critical("Must specify the URL of the testing instance!")
		raise ValueError()

	refUser = kwargs.refUser if kwargs.refUser else os.environ.get("TO_USER", None)
	if refUser is None:
		logging.critical("Must specify reference instance username!")
		raise ValueError()

	refPasswd = kwargs.refPasswd if kwargs.refPasswd else os.environ.get("TO_PASSWORD", None)
	if refPasswd is None:
		logging.critical("Must specify reference instance password!")
		raise ValueError()

	testUser = kwargs.testUser if kwargs.testUser else os.environ.get("TEST_USER", refUser)
	testPasswd = kwargs.testPasswd if kwargs.testPasswd else os.environ.get("TEST_PASSWORD", refPasswd)

	# Peel off all schemas
	if instanceA.startswith("https://"):
		instanceA = instanceA[8:]
	elif instanceA.startswith("http://"):
		instanceA = instanceA[7:]

	if instanceB.startswith("https://"):
		instanceB = instanceB[8:]
	elif instanceB.startswith("http://"):
		instanceB = instanceB[7:]

	# Parse out port numbers, if specified
	try:
		if ':' in instanceA:
			instanceA = instanceA.split(':')
			if len(instanceA) != 2:
				logging.critical("'%s' is not a valid Traffic Ops URL!", kwargs.InstanceA)
				raise ValueError()
			instanceA = {"host": instanceA[0], "port": int(instanceA[1])}
		else:
			instanceA = {"host": instanceA, "port": 443}
	except TypeError as e:
		logging.critical("'%s' is not a valid port number!", instanceA[1])
		raise ValueError from e

	try:
		if ':' in instanceB:
			instanceB = instanceB.split(':')
			if len(instanceB) != 2:
				logging.critical("'%s' is not a valid Traffic Ops URL!", kwargs.InstanceB)
				raise ValueError()
			instanceB = {"host": instanceB[0], "port": int(instanceB[1])}
		else:
			instanceB = {"host": instanceB, "port": 443}
	except TypeError as e:
		logging.critical("'%s' is not a valid port number!", instanceB[1])
		raise ValueError from e

	return (instanceA, instanceB, (refUser, refPasswd), (testUser, testPasswd))

def genRoutes(A:TOSession, B:TOSession, snapshots:bool, skip_servers:bool) ->\
                                                                  typing.Generator[str, None, None]:
	"""
	Generates routes to check for ATS config files from two valid Traffic Ops sessions

	:param A: The first Traffic Ops instance
	:param B: The second Traffic Ops instance
	:param snapshots: If ``true``, generate CDN snapshot routes, otherwise don't
	:param skip_servers: If ``true``, generation of server config files will be skipped
	:returns: A list of routes representative of the configuration files for a bunch of servers
	"""
	generatedRoutes = set()
	if not skip_servers:
		try:
			profiles = ({p.id: p for p in A.get_profiles()[0]}, {p.id: p for p in B.get_profiles()[0]})
		except (UnicodeError, InvalidJSONError, OperationError) as e:
			logging.critical("Could not fetch server profiles: %s", e)
			logging.debug("%r", e, exc_info=True, stack_info=True)
			return

		profileIds = (set(profiles[0].keys()), set(profiles[1].keys()))

		# Differences and intersections:
		for key in profileIds[0].difference(profileIds[1]):
			del profiles[0][key]
			logging.warning("profile %s found in %s but not in %s!", key, A.to_url, B.to_url)
		for key in profileIds[1].difference(profileIds[0]):
			del profiles[1][key]
			logging.warning("profile %s found in %s but not in %s!", key, B.to_url, A.to_url)

		# Now only check for identical profiles - we wouldn't expect the config files generated from
		# different profiles to be the same.
		commonProfiles = set()
		for profileId, profile in profiles[0].items():
			if profiles[1][profileId].name == profile.name:
				commonProfiles.add((profileId, profile.name, profile.type))
			else:
				logging.error("profile %s is not the same profile in both instances!", profileId)

		sampleServers = []
		for profile in commonProfiles:
			if profile[2] == "ATS_PROFILE":
				try:
					servers = A.get_servers(query_params={"profileId": profile[0]})[0]
					serverIndex = random.randint(0, len(servers)-1)
					sampleServer = servers[serverIndex]
					del servers[serverIndex]
					while not B.get_servers(query_params={"id": sampleServer.id})[0]:
						logging.warning("Server %s found in %s but not in %s!", sampleServer.id,
						                                  A.to_url, B.to_url)
						serverIndex = random.randint(0, len(servers)-1)
						sampleServer = servers[serverIndex]
						del servers[serverIndex]
				except (IndexError, ValueError):
					logging.error("Server list for profile %s exhausted without finding a sample!",
					                                  profile[1])
				except (UnicodeError, InvalidJSONError, OperationError) as e:
					logging.error("Invalid JSON response fetching server list for %s: %s", profile[2],e)
					logging.debug("%r", e, exc_info=True, stack_info=True)
				else:
					sampleServers.append(sampleServer)

		for route in getConfigRoutesForServers(sampleServers, A):
			if route not in generatedRoutes:
				yield route
				generatedRoutes.add(route)

	if snapshots:
		for route in getCRConfigs(A, B):
			if route not in generatedRoutes:
				yield route
				generatedRoutes.add(route)

def main(kwargs:argparse.Namespace) -> int:
	"""
	Runs the commandline specified by ``kwargs``.

	:param kwargs: An object that provides the attribute namespace representing this script's
		options. See ``genConfigRoutes.py --help`` for more information.
	:returns: an exit code for the program
	:raises KeyError: when ``kwargs`` does not faithfully represent a valid command line
	"""
	global LOG_FMT

	if kwargs.quiet:
		level = logging.CRITICAL + 1
	else:
		level = logging.getLevelName(kwargs.log_level)

	try:
		logging.basicConfig(level=level, format=LOG_FMT)
		logging.getLogger().setLevel(level)
	except ValueError:
		print("Unrecognized log level:", kwargs.log_level, file=sys.stderr)
		return 1

	try:
		instanceA, instanceB, loginA, loginB = consolidateVariables(kwargs)
	except ValueError as e:
		logging.debug("%s", e, exc_info=True, stack_info=True)
		logging.critical("(hint: try '-h'/'--help')")
		return 1

	verify = not kwargs.insecure

	# Instantiate connections and login
	with TOSession(host_ip=instanceA["host"], host_port=instanceA["port"], verify_cert=verify) as A,\
	TOSession(host_ip=instanceB["host"], host_port=instanceB["port"], verify_cert=verify) as B:


		try:
			A.login(loginA[0], loginA[1])
			B.login(loginB[0], loginB[1])
		except (OSError, LoginError) as e:
			logging.debug("%s", e, exc_info=True, stack_info=True)
			logging.critical("Failed to connect to Traffic Ops")
			return 2
		except (OperationError, InvalidJSONError) as e:
			logging.debug("%s", e, exc_info=True, stack_info=True)
			logging.critical("Failed to log in to Traffic Ops")
			logging.error("Error was '%s' - are you sure your URLs and credentials are correct?", e)
			return 2
		for route in genRoutes(A, B, kwargs.snapshot, kwargs.no_server_configs):
			print(route)

	return 0


if __name__ == '__main__':
	parser = argparse.ArgumentParser(description="A simple script to generate API routes to server"\
	                                 " configuration files for a given pair of Traffic Ops "\
	                                 "instances. This, for the purpose of using the 'compare' tool",
	                                 formatter_class=argparse.ArgumentDefaultsHelpFormatter)

	parser.add_argument("--refURL",
	                    help="The full URL of the reference Traffic Ops instance",
	                    type=str)
	parser.add_argument("--testURL",
	                    help="The full URL of the testing Traffic Ops instance",
	                    type=str)
	parser.add_argument("--refUser",
	                    help="A username for logging into the reference Traffic Ops instance.",
	                    type=str)
	parser.add_argument("--refPasswd",
	                    help="A password for logging into the reference Traffic Ops instance",
	                    type=str)
	parser.add_argument("--testUser",
	                    help="A username for logging into the testing Traffic Ops instance. If "\
	                         "not given, the value for the reference instance will be used.",
	                    type=str)
	parser.add_argument("--testPasswd",
	                    help="A password for logging into the testing Traffic Ops instance. If "\
	                         "not given, the value for the reference instance will be used.",
	                    type=str)
	parser.add_argument("-k", "--insecure",
	                    help="Do not verify SSL certificate signatures against *either* Traffic "\
	                         "Ops instance",
	                    action="store_true")
	parser.add_argument("-v", "--version",
	                    help="Print version information and exit",
	                    action="version",
	                    version="%(prog)s v"+__version__)
	parser.add_argument("-l", "--log_level",
	                    help="Sets the Python log level, one of 'DEBUG', 'INFO', 'WARN', 'ERROR', "\
	                         "or 'CRITICAL'",
	                    type=str,
	                    default="INFO")
	parser.add_argument("-q", "--quiet",
	                    help="Suppresses all logging output - even for critical errors",
	                    action="store_true")
	parser.add_argument("-s", "--snapshot",
	                    help="Produce snapshot routes in the output (CRConfig.json, snapshot/new "\
	                         "etc.)",
	                    action="store_true")
	parser.add_argument("-C", "--no-server-configs",
	                    help="Do not generate routes for server config files",
	                    action="store_true")
	args = parser.parse_args()
	exit(main(args))
