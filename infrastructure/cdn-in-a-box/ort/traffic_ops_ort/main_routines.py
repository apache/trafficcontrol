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
This module is meant to run the main routines of the script,
and performs a variety of operations based on the run mode.
"""

import os
import logging
import random
import time

from .configuration import Configuration
from .utils import getYesNoResponse as getYN

#: A constant that holds the absolute path to the status file directory
STATUS_FILE_DIR = "/opt/ort/status"

class ORTException(Exception):
	"""Signifies an ORT related error"""
	pass

def syncDSState(conf:Configuration) -> bool:
	"""
	Queries Traffic Ops for the :term:`Delivery Service`'s sync state

	:param conf: The script's configuration

	:returns: whether or not an update is needed
	:raises ConnectionError: when something goes wrong communicating with Traffic Ops
	"""
	logging.info("starting syncDS state fetch")

	updateStatus = conf.api.getMyUpdateStatus()

	logging.debug("Retrieved raw update status: %r", updateStatus)

	if 'upd_pending' not in updateStatus:
		raise ConnectionError("Malformed API response doesn't indicate if updates are pending!")

	if not updateStatus['upd_pending']:
		return False

	if conf.wait_for_parents and 'parent_pending' in updateStatus and updateStatus["parent_pending"]:
		logging.warning("One or more parents still have updates pending, waiting for parents.")
		return False

	if conf.mode is Configuration.Modes.SYNCDS and conf.dispersion:
		disp = random.randint(0, conf.dispersion)
		logging.info("Dispersion is set. Will sleep for %d seconds before continuing", disp)
		time.sleep(disp)

	return True

def revalidateState(conf:Configuration) -> bool:
	"""
	Checks the revalidation status of this server in Traffic Ops

	:param conf: The script's configuration

	:returns: whether or not this server has a revalidation pending
	:raises ConnectionError: when something goes wrong communicating with Traffic Ops
	"""
	logging.info("starting revalidation state fetch")

	updateStatus = conf.api.getMyUpdateStatus()

	logging.debug("Retrieved raw revalidation status: %r", updateStatus)
	if (conf.wait_for_parents and
		"parent_reval_pending" in updateStatus and
	    updateStatus["parent_reval_pending"]):
		logging.info("Parent revalidation is pending - waiting for parent")
		return False

	return "reval_pending" in updateStatus and updateStatus["reval_pending"]

def deleteOldStatusFiles(myStatus:str, conf:Configuration):
	"""
	Attempts to delete any and all old status files

	:param myStatus: the current status - files by this name will not be deleted
	:param conf: An object containing the configuration of :program:`traffic_ops_ort`
	:raises ConnectionError: if there's an issue retrieving a list of statuses from Traffic Ops
	:raises OSError: if a file cannot be deleted for any reason
	"""
	logging.info("Deleting old status files (those that are not %s)", myStatus)

	doDeleteFiles = conf.mode is not Configuration.Modes.REPORT

	for status in conf.api.get_statuses():

		# Only the status name matters
		try:
			status = status["name"]
		except KeyError as e:
			logging.debug("Bad status object: %r", status)
			raise ConnectionError from e

		if doDeleteFiles and status != myStatus:
			fname = os.path.join("/opt/ORTstatus", status)
			if not os.path.isfile(fname):
				continue
			logging.info("File '%s' to be deleted", fname)

			# check for user confirmation before deleting files in 'INTERACTIVE' mode
			if conf.mode is not Configuration.Modes.INTERACTIVE or getYN("Delete file %s?" % fname):
				logging.warning("Deleting file '%s'!", fname)
				os.remove(fname)

def setStatusFile(conf:Configuration) -> bool:
	"""
	Attempts to set the status file according to this server's reported status in Traffic Ops.

	.. warning:: This will create the directory '/opt/ORTstatus' if it does not exist, and may
		delete files there without warning!

	:param conf: An object that contains the configuration for :program:`traffic_ops_ort`
	:returns: whether or not the status file could be set properly
	"""
	global STATUS_FILE_DIR
	logging.info("Setting status file")

	try:
		myStatus = conf.api.getMyStatus()
	except ConnectionError as e:
		logging.error("Failed to set status file - Traffic Ops connection failed")
		return False

	if not os.path.isdir(STATUS_FILE_DIR):
		logging.warning("status directory does not exist, creating...")
		doMakeDir = conf.mode is not Configuration.Modes.REPORT

		# Check for user confirmation if in 'INTERACTIVE' mode
		if doMakeDir and (conf.mode is not Configuration.Modes.INTERACTIVE or
		                  getYN("Create status directory '%s'?" % STATUS_FILE_DIR, default='Y')):
			try:
				os.makedirs(STATUS_FILE_DIR)
			except OSError as e:
				logging.error("Failed to create status directory '%s' - %s", STATUS_FILE_DIR, e)
				logging.debug("%s", e, exc_info=True, stack_info=True)
				return False
	else:
		try:
			deleteOldStatusFiles(myStatus, conf)
		except ConnectionError as e:
			logging.error("Failed to delete old status files - Traffic Ops connection failed.")
			logging.debug("%s", e, exc_info=True, stack_info=True)
			return False
		except OSError as e:
			logging.error("Failed to delete old status files - %s", e)
			logging.debug("%s", e, exc_info=True, stack_info=True)
			return False

	fname = os.path.join(STATUS_FILE_DIR, myStatus)
	if not os.path.isfile(fname):
		logging.info("File '%s' to be created", fname)
		if conf.mode is not Configuration.Modes.REPORT and (
		   conf.mode is not Configuration.Modes.INTERACTIVE or getYN("Create file '%s'?", 'y')):

			try:
				with open(fname, 'x'):
					pass
			except OSError as e:
				logging.error("Failed to create status file - %s", e)
				logging.debug("%s", e, exc_info=True, stack_info=True)
				return False

	return True

def processPackages(conf:Configuration) -> bool:
	"""
	Manages the packages that Traffic Ops reports are required for this server.

	:param conf: An object containing the configuration of :program:`traffic_ops_ort`
	:returns: whether or not the package processing was successfully completed
	"""
	try:
		myPackages = conf.api.getMyPackages()
	except ConnectionError as e:
		logging.error("Packages not found or API response malformed! - %s", e)
		logging.debug("%s", e, exc_info=True, stack_info=True)
		return False

	for package in myPackages:
		if package.install(conf):
			if conf.mode is not Configuration.Modes.BADASS:
				return False
			logging.warning("Failed to install %s, but we're BADASS, so moving on!", package)

	return True

def processServices(conf:Configuration) -> bool:
	"""
	Manages the running processes of the server, according to an ancient system known as 'chkconfig'

	:param conf: An object containing the configuration for :program:`traffic_ops_ort`
	:returns: whether or not the service processing was completed successfully
	"""
	from . import services

	if not services.HAS_SYSTEMD:
		logging.warning("This system doesn't have systemd, services cannot be enabled/disabled")
		return True


	try:
		chkconfig = conf.api.getMyChkconfig()
	except ConnectionError as e:
		logging.error("Failed to fetch 'chkconfig' from Traffic Ops! (%s)", e)
		logging.debug("%r", e, exc_info=True, stack_info=True)
		return False

	for item in chkconfig:
		logging.debug("Processing item %r", item)

		if not services.setServiceStatus(item, conf.mode):
			return False

	return True

def processConfigurationFiles(conf:Configuration) -> bool:
	"""
	Updates and backs up all of a server's configuration files.

	:param conf: An object containing the configuration for :program:`traffic_ops_ort`
	:returns: whether or not the configuration changes were successful
	"""
	from . import config_files, services

	try:
		config_files.initBackupDir(conf.mode)
	except OSError as e:
		logging.error("Couldn't create backup directory!")
		logging.warning("%s", e)
		logging.debug("", exc_info=True, stack_info=True)
		return False

	try:
		myFiles = conf.api.getMyConfigFiles(conf)
	except ConnectionError as e:
		logging.critical("Failed to fetch configuration files; Traffic Ops connection failed! %s",e)
		logging.debug("%s", e, exc_info=True, stack_info=True)
		return False

	for file in myFiles:
		try:
			logging.info("\n============ Processing File: %s ============", file.fname)
			if file.update(conf) and file.fname in services.FILES_THAT_REQUIRE_RELOADS:
				services.NEEDED_RELOADS.add(services.FILES_THAT_REQUIRE_RELOADS[file.fname])
			logging.info("\n============================================\n")

		# A bad object could just reflect an inconsistent reply structure from the API, so BADASSes
		# will attempt to continue. However, an issue updating a valid configuration is not
		# recoverable, even for BADASSes
		except OSError as e:
			logging.error("An error occurred while trying to update %s", file.fname)
			logging.debug("%s", e, exc_info=True, stack_info=True)
			return False

	return True

def run(conf:Configuration) -> int:
	"""
	This function is the entrypoint into the script's main flow from :func:`traffic_ops_ort.doMain`
	It runs the appropriate actions depending on the run mode.

	:param conf: An object that holds the script's configuration

	:returns: an exit code for the script
	"""
	from . import services

	# If this is just a revalidation, then we can exit if there's no revalidation pending
	if conf.mode is Configuration.Modes.REVALIDATE:
		try:
			updateRequired = revalidateState(conf)
		except ConnectionError as e:
			logging.critical("Server configuration unreachable, or not found in Traffic Ops!")
			logging.error(e)
			logging.debug("%r", e, exc_info=True, stack_info=True)
			return 2

		if not updateRequired:
			logging.info("No revalidation pending")
			return 0

		logging.info("in REVALIDATE mode; skipping package/service processing")

	# In all other cases, we check for an update to the Delivery Service and apply any found
	# changes
	else:
		try:
			updateRequired = syncDSState(conf)
		except ConnectionError as e:
			logging.critical("Server configuration unreachable, or not found in Traffic Ops!")
			logging.error(e)
			logging.debug("%r", e, exc_info=True, stack_info=True)
			return 2

		# Bail on failures - unless this script is BADASS!
		if not setStatusFile(conf):
			if conf.mode is not Configuration.Modes.BADASS:
				logging.critical("Failed to set status as specified by Traffic Ops")
				return 2
			logging.warning("Failed to set status but we're BADASS, so moving on.")

		logging.info("\nProcessing Packages...")
		if not processPackages(conf):
			logging.critical("Failed to process packages")
			if conf.mode is not Configuration.Modes.BADASS:
				return 2
			logging.warning("Package processing failed but we're BADASS, so attempting to move on")
		logging.info("Done.\n")

		logging.info("\nProcessing Services...")
		if not processServices(conf):
			logging.critical("Failed to process services.")
			if conf.mode is not Configuration.Modes.BADASS:
				return 2
			logging.warning("Service processing failed but we're BADASS, so attempting to move on")
		logging.info("Done.\n")


	# All modes process configuration files
	logging.info("\nProcessing Configuration Files...")
	if not processConfigurationFiles(conf):
		logging.critical("Failed to process configuration files.")
		return 2
	logging.info("Done.\n")

	if updateRequired:
		if (conf.mode is not Configuration.Modes.INTERACTIVE or
		    getYN("Update Traffic Ops?", default='Y')):

			logging.info("\nUpdating Traffic Ops...")
			conf.api.updateTrafficOps(conf.mode)
			logging.info("Done.\n")
		else:
			logging.warning("Traffic Ops was not notified of changes. You should do this manually.")

	else:
		logging.info("Traffic Ops update not necessary")

	if services.NEEDED_RELOADS and not services.doReloads(conf):
		logging.critical("Failed to reload all configuration changes")
		return 2

	return 0
