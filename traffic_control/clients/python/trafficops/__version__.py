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
Versioning
==========
The :mod:`trafficops.__version__` module contains only the :data:`__version__` "constant" which
gives the version of this *Apache-TrafficControl package* and **not** the version of
*Apache Traffic Control* for which it was made. The two are versioned separately, to allow the
client to grow in a version-controlled manner without being tied to the release cadence of Apache
Traffic Control as a whole.

Version 3.0 is supported for use with Apache Traffic Control versions 7.0 and 6.1. New functionality
will be added as the :ref:`to-api` evolves, but changes to this client will remain non-breaking for
existing code using it until the next major version is released.
"""

__version__ = '3.1.0'
