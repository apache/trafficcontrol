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

"""Required packages to run api contract test cases."""
from setuptools import setup
setup(
    name='trafficcontrol',
    version='1.0.0',
    install_requires=[
    'Apache-TrafficControl @ file:///os.path.abspath(os.getcwd())/traffic_control/clients/python',
    'pytest==7.2.1'
    ],
)
