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

import io
import os
import sys
from setuptools import setup, find_packages

pyversion = sys.version_info
if pyversion.major < 3:
	sys.stderr.write("Python 2 is no longer supported\n")
	sys.exit(1)

elif pyversion.major == 3 and (pyversion.minor < 6 or pyversion.minor > 8):
	MSG = ('WARNING: This library may not work properly with Python {0.major}.{0.minor}.{0.micro}, '
	       'as it is untested for this version. (3.6 <= version <=3.8 recommended)\n')
	sys.stderr.write(MSG.format(pyversion))

HERE = os.path.abspath(os.path.dirname(__file__))

about = {}
with io.open(os.path.join(HERE, 'trafficops', '__version__.py'), mode='r', encoding='utf-8') as f:
	exec(f.read(), about)


with open(os.path.join(HERE, "README.rst")) as fd:
	setup(
	       name="Apache-TrafficControl",
	       version=about["__version__"],
	       author="Apache Software Foundation",
	       author_email="dev@trafficcontrol.apache.org",
	       description="Python API Client for Traffic Control",
	       long_description='\n'.join((fd.read(), about["__doc__"])),
	       url="https://trafficcontrol.apache.org/",
	       license="http://www.apache.org/licenses/LICENSE-2.0",
	       classifiers=[
	           'Development Status :: 4 - Beta',
	           'Intended Audience :: Developers',
	           'Intended Audience :: Information Technology',
	           'License :: OSI Approved :: Apache Software License',
	           'Operating System :: OS Independent',
	           'Programming Language :: Python :: 2.7',
	           'Programming Language :: Python :: 3',
	           'Programming Language :: Python :: 3.4',
	           'Programming Language :: Python :: 3.5',
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
	           "requests>=2.26.0",
	           "munch>=2.1.1",
	       ],
	       entry_points={
	           'console_scripts': [
	       	    'toget=to_access:get',
	       	    'toput=to_access:put',
	       	    'topost=to_access:post',
	       	    'todelete=to_access:delete',
	       	    'tooptions=to_access:options',
	       	    'tohead=to_access:head',
	       	    'topatch=to_access:patch'
	            ],
	       },
	       extras_require={
	           "dev": [
	               "pylint>=2.0,<3.0"
	           ]
	       },
	       # This will only be enforced by pip versions >=9.0 (18.1 being current at the time of
	       # this writing) - i.e. not the pip installed as python34-pip from elrepo on CentOS
	       python_requires='>=3.6'
	)

