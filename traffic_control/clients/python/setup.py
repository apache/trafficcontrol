#!/usr/bin/env python3
# -*- coding: utf-8 -*-

#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

"""
Setuptools build/install script for the Python Traffic Control client
"""

import sys
import os
from setuptools import setup, find_packages
from trafficops import __version__ as version

HERE = os.path.abspath(os.path.dirname(__file__))

with open(os.path.join(HERE, "README.rst")) as fd:
	setup(
	       name="Apache-TrafficControl",
	       version=version,
	       author="Apache Software Foundation",
	       author_email="dev@trafficcontrol.apache.org",
	       description="Python API Client for Traffic Control",
	       long_description=fd.read(),
	       url="http://trafficcontrol.apache.org/",
	       license="http://www.apache.org/licenses/LICENSE-2.0",
	       classifiers=[
	           'Development Status :: 4 - Beta',
	           'Intended Audience :: Developers',
	           'Intended Audience :: Information Technology',
	           'License :: OSI Approved :: Apache Software License',
	           'Operating System :: OS Independent',
	           'Programming Language :: Python :: 2.7',
	           'Programming Language :: Python :: 3',
	           'Programming Language :: Python :: 3.6',
	           'Programming Language :: Python :: 3.7',
	           'Programming Language :: Python :: 3.8',
	           'Topic :: Database :: Front-Ends',
	           'Topic :: Internet :: WWW/HTTP',
	           'Topic :: Software Development :: Libraries :: Python Modules'
	       ],
	       keywords="Apache TrafficControl Traffic Control Client TrafficOps Ops",
	       packages=find_packages(exclude=["contrib", "docs", "tests"]),
	       install_requires=[
	           "future>=0.16.0",
	           "requests>=2.13.0",
	           "munch>=2.1.1",
	       ],
	       extras_require={
	           "dev": [
	               "pylint"
	           ]
	       }
	)

if ((sys.version_info[0] == 2 and sys.version_info < (2, 7))
   or (sys.version_info[0] == 3 and sys.version_info < (3, 6))):
	MSG = ('WARNING: This library may not work properly with Python {0}, '
		   'as it is untested for this version.')
	print(MSG.format(sys.version.split(' ', 1)[0]))
