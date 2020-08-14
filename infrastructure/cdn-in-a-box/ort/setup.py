#!/usr/bin/env python3
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
The setuptools-based install script for Traffic Ops ORT
"""

import os
import sys

from setuptools import setup, find_packages

here = os.path.abspath(os.path.dirname(__file__))
sys.path.append(here)
import traffic_ops_ort

with open(os.path.join(here, 'README.rst')) as f:
	long_description = f.read()

setup(
	name='traffic_ops_ort',
	version=traffic_ops_ort.__version__,
	description="The Traffic Ops Operational Readiness Test, checks for and applies changes from "\
	"Traffic Ops to a server",
	long_description=long_description,
	url='https://github.com/ocket8888/trafficcontrol/tree/ort-without-systemd/infrastructure/cdn-in-a-box/ort',
	author='Brennan Fieck',
	author_email='Brennan_WilliamFieck@comcast.com',
	classifiers=[
		'Development Status :: 3 - Alpha',
		'Environment :: Console',
		'Intended Audience :: Telecommunications Industry',
		'Intended Audience :: Developers',
		'Intended Audience :: System Administrators',
		'License :: OSI Approved :: Apache Software License',
		'Operating System :: POSIX :: Linux',
		'Programming Language :: Python :: Implementation :: CPython',
		'Programming Language :: Python :: Implementation :: PyPy',
		'Programming Language :: Python :: 3 :: Only',
		'Programming Language :: Python :: 3.4',
		'Programming Language :: Python :: 3.5',
		'Programming Language :: Python :: 3.6',
		'Programming Language :: Python :: 3.7',
		'Topic :: Internet',
		'Topic :: Internet :: WWW/HTTP',
		'Topic :: Internet :: Log Analysis',
		'Topic :: Utilities'
	],
	keywords='network connection configuration TrafficControl',
	packages=find_packages(exclude=['contrib', 'docs', 'tests']),
	install_requires=['setuptools', 'typing', 'requests', 'urllib3', 'distro', 'psutil', 'Apache-TrafficControl'],
	# data_files=[('etc/crontab', ['traffic_ops_ort.crontab'])],
	entry_points={
		'console_scripts': [
			'traffic_ops_ort=traffic_ops_ort:main',
		],
	},
	python_requires='~=3.4'
)
