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
"""

import argparse
import logging
import os
import random
import time
import typing
import sys

random.seed(time.time())

#: The repository root directory
ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "..", ".."))

#: An absolute path to the Traffic Ops python packages (This assumes that the script is run from
#: within the repository's normal directory structure)
TO_LIBS_PATH = os.path.join(ROOT, "traffic_control", "clients", "python", "trafficops")

sys.path.insert(0, TO_LIBS_PATH)
sys.path.insert(0, os.path.join(TO_LIBS_PATH, "trafficops"))
from trafficops.tosession import TOSession
from common.restapi import LoginError, OperationError, InvalidJSONError

#: A format specifier for logging output. Propagates to all imported modules.
LOG_FMT = "%(levelname)s: %(asctime)s line %(lineno)d in %(module)s.%(funcName)s: %(message)s"

__version__ = "1.0.0"

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
		for file in inst.getServerConfigFiles(servername=server.hostName)[0].configFiles:
			if "apiUri" in file:
				yield file.apiUri
			else:
				logging.info("config file %s for server %s has non-API URI - skipping",
				                    file.location, server.hostName)

def getCRConfigs(A:TOSession, B:TOSession) -> typing.Generator[str, None, None]:
	"""
	Generates a list of routes to CRConfig files for all CDNs present in both A and B

	:param A: The first Traffic Ops instance
	:param B: The second Traffic Ops instance
	:returns: A list of routes to CRConfig files
	"""
	cdns = {c.name for c in A.get_cdns()[0]}.intersection({c.name for c in B.get_cdns()[0]})

	if not cdns:
		logging.error("The two instances have NO CDNs in common! This almost certainly means that "\
		              "you're not doing what you want to do")
		yield from []
	else:
		yield from ["CRConfig-Snapshots/%s/CRConfig.json" % cdn for cdn in cdns]



def genRoutes(A:TOSession, B:TOSession) -> typing.List[str]:
	"""
	Generates routes to check for ATS config files from two valid Traffic Ops sessions

	:param A: The first Traffic Ops instance
	:param B: The second Traffic Ops instance
	:returns: A list of routes representative of the configuration files for a bunch of servers
	"""
	profiles = ({p.id: p for p in A.get_profiles()[0]}, {p.id: p for p in B.get_profiles()[0]})
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
			servers = A.get_servers(query_params={"profileId": profile[0]})[0]
			try:
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
			else:
				sampleServers.append(sampleServer)

	generatedRoutes = set()
	for route in getConfigRoutesForServers(sampleServers, A):
		if route not in generatedRoutes:
			yield route
			generatedRoutes.add(route)

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

	instanceA = kwargs.InstanceA
	instanceB = kwargs.InstanceB
	loginA = kwargs.LoginA.split(':')
	loginB = kwargs.LoginB.split(':') if kwargs.LoginB else loginA
	verify = not kwargs.insecure

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
				return 1
			instanceA = {"host": instanceA[0], "port": int(instanceA[1])}
		else:
			instanceA = {"host": instanceA, "port": 443}
	except TypeError as e:
		logging.critical("'%s' is not a valid port number!", instanceA[1])
		logging.debug("%s", e, exc_info=True, stack_info=True)
		return 1

	try:
		if ':' in instanceB:
			instanceB = instanceB.split(':')
			if len(instanceB) != 2:
				logging.critical("'%s' is not a valid Traffic Ops URL!", kwargs.InstanceB)
			instanceB = {"host": instanceB[0], "port": int(instanceB[1])}
		else:
			instanceB = {"host": instanceB, "port": 443}
	except TypeError as e:
		logging.critical("'%s' is not a valid port number!", instanceB[1])
		logging.debug("%s", e, exc_info=True, stack_info=True)
		return 1

	# Instantiate connections and login
	with TOSession(host_ip=instanceA["host"], host_port=instanceA["port"], verify_cert=verify) as A,\
	TOSession(host_ip=instanceB["host"], host_port=instanceB["port"], verify_cert=verify) as B:


		try:
			A.login(loginA[0], loginA[1])
			B.login(loginB[0], loginB[1])
		except OSError as e:
			logging.debug("%s", e, exc_info=True, stack_info=True)
			logging.critical("Failed to connect to Traffic Ops")
			return 2
		except (OperationError, LoginError) as e:
			logging.debug("%s", e, exc_info=True, stack_info=True)
			logging.critical("Failed to log in to Traffic Ops")
			logging.error("Error was '%s' - are you sure your URLs and credentials are correct?", e)
		for route in genRoutes(A, B):
			print(route)

	return 0


if __name__ == '__main__':
	parser = argparse.ArgumentParser(description="A simple script to generate API routes to server"\
	                                 " configuration files for a given pair of Traffic Ops "\
	                                 "instances. This, for the purpose of using the 'compare' tool",
	                                 formatter_class=argparse.ArgumentDefaultsHelpFormatter)

	parser.add_argument("InstanceA",
	                    help="The full URL of the first Traffic Ops instance",
	                    type=str)
	parser.add_argument("InstanceB",
	                    help="The full URL of the second Traffic Ops instance",
	                    type=str)
	parser.add_argument("LoginA",
	                    help="A login string containing credentials for logging into the first "\
	                         "Traffic Ops instance (InstanceA) in the format 'username:password'.",
	                    type=str)
	parser.add_argument("LoginB",
	                    help="A login string containing credentials for logging into the second "\
	                         "Traffic Ops instance (InstanceB) in the format 'username:password'. "\
	                         "If not given, LoginA will be re-used for the second connection",
	                    type=str,
	                    nargs='?')
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
	                         "OR 'CRITICAL'",
	                    type=str,
	                    default="INFO")
	parser.add_argument("-q", "--quiet",
	                    help="Suppresses all logging output - even for critical errors",
	                    action="store_true")
	args = parser.parse_args()
	exit(main(args))
