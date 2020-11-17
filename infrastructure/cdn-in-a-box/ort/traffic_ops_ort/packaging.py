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
This module deals with managing packages installed on the host system.
It attempts to do so in a distribution-agnostic way, but actually only
support a strict set of distributions with well-known package managers.
"""

import logging
import subprocess

from .configuration import Configuration
from .utils import getYesNoResponse as getYN

class _MetaPackage(type):
	"""
	This factory is responsible for constructing a :class:`Package` class which properly
	utilizes the host system's package management platform.
	"""

	def __new__(mcs, name, bases, dct) -> type:
		"""
		Customization point for distribution-dependent functionality
		"""
		from .configuration import DISTRO
		pack = super().__new__(mcs, name, bases, dct)

		if DISTRO in {'fedora', 'centos', 'rhel'}:
			concat = '-'
			pack.checkInstallList = lambda x: [ p
			                                    for p in
			                                    subprocess.Popen(["/bin/rpm", "-q", x.name],
			                                                     stdout=subprocess.PIPE)
			                                              .communicate()[0].decode().splitlines()
			                                    if not p.endswith("is not installed")]
			pack.installArgs = ["/bin/dnf", "install", "-y"]
			pack.uninstallArgs = ["/bin/dnf", "remove", "-y"]

		elif DISTRO in {'ubuntu', 'linuxmint', 'debian'}:
			concat = '='
			pack.checkInstallList = lambda x: [ "{1}={2}".format(*p.split())
			                                    for p in
			                                    subprocess.Popen(["/usr/bin/dpkg", "-l", x.name],
			                                                     stdout=subprocess.PIPE)
			                                    .communicate()[0].decode().splitlines()[5:]
			                                    if p ]

			pack.installArgs = ["/usr/bin/apt-get", "install", "-y"]
			pack.uninstallArgs = ["/usr/bin/apt-get", "purge", "-y"]

		# TODO - is this reasonable? I mean, I KNOW it's unreasonable because this whole module
		# shouldn't exist, but is this a good fallback since it DOES exist?
		else:
			concat = " v"
			pack.checkInstallList = lambda x: []
			pack.installArgs = ["/bin/true"]
			pack.uninstallArgs = ["/bin/true"]
		pack.__str__ = lambda x: concat.join((x.name, x.version)) if x.version else x.name

		return pack




class Package(metaclass=_MetaPackage):
	"""
	Represents a package installed (or about to be installed) on the host system
	"""

	#: The package's name
	name = None

	#: Optionally, a specific version of the package
	version = None

	def __init__(self, pkg:dict):
		"""
		Constructs a :class:`Package` object from a raw JSON response

		:param pkg: a parsed JSON response, expected to contain a ``"name"`` key and optionally a
			``"version"`` key.
		:raises ValueError: if ``pkg`` does not faithfully represent a package
		"""
		if "name" not in pkg:
			raise ValueError("%r does not represent a package!" % pkg)

		self.name = pkg["name"]
		self.version = pkg["version"] if "version" in pkg else ""

		# These are defined in the metaclass based on the host system's Linux distribution, but are
		# specified here for the benefit of static analysis tools
		self.checkInstallList = getattr(self, 'checkInstallList', lambda: ())
		self.installArgs = getattr(self, 'installArgs', None)
		self.uninstallArgs = getattr(self, 'uninstallArgs', None)

	def __repr__(self) -> str:
		"""
		Implements ``repr(self)``
		"""
		if self.version:
			return "Package(name=%r, version=%r)" % (self.name, self.version)
		return "Package(name=%r)" % (self.name,)



	def isInstalled(self) -> bool:
		"""
		Checks if this package is already present on the system

		:returns: whether or not the package could be found
		"""
		for pkg in self.checkInstallList():
			if (self.version and pkg == str(self)) or pkg.startswith(self.name):
				return True
		return False

	def install(self, conf:Configuration) -> int:
		"""
		Installs this package.

		:param conf: An object containing the configuration for :program:`traffic_ops_ort`

		:returns: the exit code of the install process
		"""
		if self.isInstalled():
			logging.info("%s is already installed - nothing to do", self)
			return 0

		if conf.mode is Configuration.Modes.INTERACTIVE and not getYN("Install %s?" % self, 'Y'):
			logging.warning("%s will not be installed, dependencies may be unsatisfied!", self)
			return 0

		logging.info("Installing %s", self)

		if conf.mode is Configuration.Modes.REPORT:
			return 0

		try:
			sub = subprocess.Popen(self.installArgs + [str(self)],
			                       stdout=subprocess.PIPE,
			                       stderr=subprocess.PIPE)
			out, err = sub.communicate(timeout=60)
		except subprocess.CalledProcessError as e:
			logging.debug("%r", e, stack_info=True, exc_info=True)
			logging.error("%s could not be installed!", self)
			return -1
		except TimeoutError as e:
			logging.debug("%r", e, stack_info=True, exc_info=True)
			logging.error("Package install timed out!")
			return -1

		logging.debug("STDOUT: %s", out.decode())
		logging.debug("STDERR: %s", err.decode())

		return sub.returncode

	def uninstall(self, conf:Configuration) -> int:
		"""
		Uninstalls this package. I have no idea how one would make use of this from within ATC...

		:returns: the exit code of the uninstall process
		"""
		if not self.isInstalled():
			logging.info("%s is not installed - nothing to do", self)
			return 0

		if conf.mode is Configuration.Modes.INTERACTIVE and not getYN("Uninstall %s?" % self, 'Y'):
			logging.warning("%s will not be installed, dependencies may be out of date!", self)
			return 0

		logging.info("Uninstalling %s", self)

		if conf.mode is Configuration.Modes.REPORT:
			return 0

		try:
			sub = subprocess.Popen(self.uninstallArgs + [str(self)],
			                       stdout=subprocess.PIPE,
			                       stderr=subprocess.PIPE)
			out, err = sub.communicate(timeout=60)
		except subprocess.CalledProcessError as e:
			logging.debug("%r", e, stack_info=True, exc_info=True)
			logging.error("%s could not be uninstalled!", self)
			return -1
		except TimeoutError as e:
			logging.debug("%r", e, stack_info=True, exc_info=True)
			logging.error("Package uninstall timed out!")
			return -1

		logging.debug("STDOUT: %s", out.decode())
		logging.debug("STDERR: %s", err.decode())
		return sub.returncode
