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

# from . import configuration

# TODO - metaclass packages per-distribution

class Package():
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

	def __repr__(self) -> str:
		"""
		Implements ``repr(self)``
		"""
		if self.version:
			return "Package(name=%r, version=%r)" % (self.name, self.version)
		return "Package(name=%r)" % (self.name,)

	def __str__(self) -> str:
		"""
		Implements ``str(self)``
		"""
		if self.version:
			return '-'.join((self.name, self.version))
		return self.name
