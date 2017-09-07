import sys

from setuptools import setup

version = ''
for line in open('trafficops/__init__.py').readlines():
    if '__version__' in line:
        version = line.split('=')[-1].strip().strip("'")
        break

setup(
    name='TrafficOps',
    version=version,
    author='Robert Scrimo Jr.',
    author_email='robert_scrimo@comcast.com',
    packages=['trafficops'],
    url='http://trafficcontrol.apache.org/',
    license='http://www.apache.org/licenses/LICENSE-2.0',
    description='Python API Client for Traffic Ops',
    long_description=open('README.txt').read(),
    install_requires=[
        "future>=0.16.0",
        "requests>=2.13.0",
        "munch>=2.1.1",
    ],
)

if ((sys.version_info[0] == 2 and sys.version_info < (2, 7))
   or (sys.version_info[0] == 3 and sys.version_info < (3, 6))):
    msg = ('WARNING: This library may not work properly with Python {0}, '
           'as it is untested for this version.')
    print(msg.format(sys.version.split(' ', 1)[0]))
