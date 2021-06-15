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
This module handles the reloading of service configuration files when they are changed.
in :mod:`config_files`, when a configuration file is changed, its base name is looked up
in this module's :data:`FILES_THAT_REQUIRE_RELOADS` map, and - if found - the function to
which it is mapped will be added to the :data:`NEEDED_RELOADS` set. Thus, to implement this
for a new service, all that need be done is to write a function to perform configuration
reloads, then add the names of configuration files to the map, pointing at the new function.
"""

import logging
import os
import subprocess
import typing

from functools import partial

import psutil

from .configuration import Configuration
from .utils import getYesNoResponse as getYN

#: Holds the list of reloads needed due to configuration file changes
NEEDED_RELOADS = set()

#: True if the host system has systemd D-Bus - actual value set at runtime
HAS_SYSTEMD = False

try:
	output = subprocess.check_output(["systemctl", "--no-pager"], stderr=subprocess.STDOUT)
except subprocess.CalledProcessError:
	logging.debug("Host system does NOT have systemd - stack trace:", exc_info=True,stack_info=True)
else:
	HAS_SYSTEMD = True

def reloadATSConfigs(conf:Configuration) -> bool:
	"""
	This function will reload configuration files for the Apache Trafficserver caching HTTP
	proxy. It does this by calling ``traffic_ctl config reload`

	:param conf: An object representing the configuration of :program:`traffic_ops_ort`
	:returns: whether or not the reload succeeded (as indicated by the exit code of
		``traffic_ctl``)
	:raises OSError: when something goes wrong executing the child process
	"""
	# First of all, ATS must be running for this to work
	if not setATSStatus(True, conf):
		logging.error("Cannot reload configs, ATS not running!")
		return False

	cmd = [os.path.join(conf.tsroot, "bin", "traffic_ctl"), "config", "reload"]
	cmdStr = ' '.join(cmd)

	if ( conf.mode is Configuration.Modes.INTERACTIVE and
	     not getYN("Run command '%s' to reload configuration?" % cmdStr, default='Y')):
		logging.warning("Configuration will not be reloaded for Apache Trafficserver!")
		logging.warning("Changes will NOT be applied!")
		return True

	logging.info("Apache Trafficserver configuration reload will be done via: %s", cmdStr)

	if conf.mode is Configuration.Modes.REPORT:
		return True

	sub = subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
	out, err = sub.communicate()

	if sub.returncode:
		logging.debug("========== PROCESS STDOUT ==========")
		logging.debug("%s", out.decode())
		logging.debug("========== PROCESS STDERR ==========")
		logging.debug("%s", err.decode())
		logging.debug("====================================")
		return False
	return True

def restartATS(conf:Configuration) -> bool:
	"""
	A convenience function for calling :func:`setATSStatus` for restarts.

	:param conf: An object representing the configuration of :program:`traffic_ops_ort`
	:returns: whether or not the restart was successful (or unnecessary)
	"""

	doRestart = ( conf.mode is Configuration.Modes.BADASS or
	              conf.mode is Configuration.Modes.REPORT or
	              ( conf.mode is Configuration.Modes.INTERACTIVE and
	                getYN("Restart ATS?", default='Y')))

	return setATSStatus(True, conf, restart=doRestart)


def restartService(service:str, conf:Configuration) -> bool:
	"""
	Restarts a generic systemd service

	:param service: The name of the service to be restarted
	:param conf: An object representing the configuration of :program:`traffic_ops_ort`
	:returns: Whether or not the restart was successful
	"""
	global HAS_SYSTEMD

	if not HAS_SYSTEMD:
		logging.warning("This system doesn't have systemd, services cannot be restarted")
		return True

	if conf.mode is not Configuration.Modes.REPORT and (
	   conf.mode is not Configuration.Modes.INTERACTIVE or getYN("Restart %s?" % service, 'Y')):
		logging.info("Restarting %s", service)
		try:
			sub = subprocess.Popen(["systemctl", "restart", service],
			                       stdout=subprocess.PIPE,
			                       stderr=subprocess.PIPE)
			out, err = sub.communicate()
			logging.debug("stdout: %s\nstderr: %s", out, err)
		except (OSError, subprocess.CalledProcessError) as e:
			logging.error("An error occurred when restarting %s: %s", service, e)
			logging.debug("%r", e, exc_info=True, stack_info=True)
			return False
	return True

#: A big ol' map of filenames to the services which require reloads when said files change
FILES_THAT_REQUIRE_RELOADS = {"records.config":          reloadATSConfigs,
                              "remap.config":            reloadATSConfigs,
                              "parent.config":           reloadATSConfigs,
                              "cache.config":            reloadATSConfigs,
                              "hosting.config":          reloadATSConfigs,
                              "astats.config":           reloadATSConfigs,
                              "logs_xml.config":         reloadATSConfigs,
                              "ssl_multicert.config":    reloadATSConfigs,
                              "regex_revalidate.config": reloadATSConfigs,
                              "plugin.config":           restartATS,
                              "ntpd.conf":               partial(restartService, "ntpd"),
                              "50-ats.rules":            restartATS}

def doReloads(conf:Configuration) -> bool:
	"""
	Performs all necessary service restarts/configuration reloads

	:param conf: An object representing the configuration of :program:`traffic_ops_ort`
	:returns: whether or not the reloads/restarts went successfully
	"""
	global NEEDED_RELOADS

	# If ATS is being restarted, configuration reloads will be implicit
	if restartATS in NEEDED_RELOADS and reloadATSConfigs in NEEDED_RELOADS:
		NEEDED_RELOADS.discard(reloadATSConfigs)

	for reload in NEEDED_RELOADS:
		try:
			if not reload(conf):
				return False
		except OSError as e:
			logging.error("An error occurred when reloading service configuration files: %s",e)
			logging.debug("%s", e, exc_info=True, stack_info=True)
			return False

	return True

def getProcessesIfRunning(name:str) -> typing.Optional[psutil.Process]:
	"""
	Retrieves a process by name, if it exists.

	.. warning:: Process names don't have to be unique, this will return the process with the
		lowest PID that matches ``name``. This can also only return processes visible to the
		user running the Python interpreter.

	:param name: the name for which to search
	:returns: a process if one is found that matches ``name``, else :const:`None`
	:raises OSError: if the process table cannot be iterated
	"""
	logging.debug("Iterating process list - looking for %s", name)
	for process in psutil.process_iter():

		# Found
		if process.name() == name:
			logging.debug("Running process found (pid: %d)", process.pid)
			return process

	logging.debug("No process named '%s' was found", name)

	return None

def setATSStatus(status:bool, conf:Configuration, restart:bool = False) -> bool:
	"""
	Sets the status of the system's ATS process.

	:param status: Specifies whether ATS should be running (:const:`True`) or not (:const:`False`)
	:param restart: If this is :const:`True`, then ATS will be restarted if it is already running

		.. note:: ``restart`` has no effect if ``status`` is :const:`False`

	:returns: whether or not the status setting was successful (or unnecessary)
	:raises OSError: when there is a problem executing the subprocess
	"""
	existingProcess = getProcessesIfRunning("[TS_MAIN]")

	# ATS is not running
	if existingProcess is None:
		if not status:
			logging.info("ATS already stopped - nothing to do")
			return True
		logging.info("ATS not running, will be started")
		arg = "start"


	# ATS is running, but has a bad status
	elif status and existingProcess.status() not in {psutil.STATUS_RUNNING, psutil.STATUS_SLEEPING}:
		logging.warning("ATS already running, but status is %s - restarting",
		                                       existingProcess.status())
		arg = "restart"

	# ATS is running and should be stopped
	elif not status:
		logging.info("ATS is running, will be stopped")
		arg = "stop"

	# ATS is running fine, but we want to restart it
	elif restart:
		logging.info("ATS process found - restarting")
		arg = "restart"

	# ATS is running fine already
	else:
		logging.info("ATS already running - nothing to do")
		return True

	tsexe = os.path.join(conf.tsroot, "bin", "trafficserver")
	if ( conf.mode is Configuration.Modes.INTERACTIVE and
	     not getYN("Run command '%s %s'?" % (tsexe, arg))):
		logging.warning("ATS status will not be set - Traffic Ops may not expect this!")
		return True

	logging.info("ATS status will be set using: %s %s", tsexe, arg)

	if conf.mode is not Configuration.Modes.REPORT:

		sub = subprocess.Popen([tsexe, arg], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
		out, err = sub.communicate()

		if sub.returncode:
			logging.error("Failed to start/stop/restart trafficserver!")
			logging.warning("Is the 'trafficserver' script located at %s?", tsexe)
			logging.debug(out.decode())
			logging.debug(err.decode())
			return False
	return True

def setServiceStatus(chkconfig:dict, mode:Configuration.Modes) -> bool:
	"""
	Sets the status of a service based on its 'chkconfig'.
	A 'chkconfig' consists of a list of run-levels with either 'on' or 'off' as values.
	This allowed specifying what run-levels needed a service. It's now totally deprecated,
	but the Traffic Ops back-end doesn't know that yet...

	:param chkconfig: A single chkconfig
	:param mode: The current run-mode
	:returns: whether or not the service's status was set successfully
	"""
	global HAS_SYSTEMD

	try:
		status = "enable" if "on" in chkconfig["value"] else "disable"
		service = chkconfig['name']
	except KeyError as e:
		logging.error("'%r' could not be parsed as a chkconfig object!", chkconfig)
		logging.debug("%s", e, exc_info=True, stack_info=True)
		return False

	if (mode is Configuration.Modes.INTERACTIVE and
	    not getYN("%s %s?" % (service, status), default='Y')):
		logging.warning("%s will not be %sd - some things may break!", service, status)
		return True

	logging.info("%s will be %sd", service, status)

	if mode is not Configuration.Modes.REPORT:
		try:
			sub = subprocess.Popen(["systemctl", status, service],
			                       stdout=subprocess.PIPE,
			                       stderr=subprocess.PIPE)
			out, err = sub.communicate()
			logging.debug("output")
		except (OSError, subprocess.CalledProcessError) as e:
			logging.error("An error occurred when %sing %s: %s", status[:-1], service, e)
			logging.debug("%r", e, exc_info=True, stack_info=True)
			return False

	return True
