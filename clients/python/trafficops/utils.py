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
utils
=====
Useful utility methods
"""

# Core Modules
import os
import inspect
import logging

# Python 2 to Python 3 Compatibility
from builtins import str

LOGGER = logging.getLogger(__name__)

def log_with_debug_info(logging_level=logging.INFO, msg=u'', parent=False, separator=u':'):
	"""
	Uses inspect module(reflection) to gather debugging information for the source file name,
	function name, and line number of the calling function/method.

	:param logging_level: The logging level from the logging module constants
						  E.g. logging.INFO, logging.DEBUG, etc.
	:type logging_level: int
	:param msg: The message to log.
	:type msg: Text
	:param parent: If True, use the caller's parent information instead of the caller's information
		in the message.
	:type parent: bool
	:param separator: The string to use for the component separator
	:type separator: Text
	:return: '<file name>:<function name>:<line number>: <msg>'
		e.g. 'tosession.py:_build_endpoint:199: This is a message to log.'
	:rtype: Text
	"""

	frame,\
	file_path,\
	line_number,\
	function_name,\
	_,\
	_ = inspect.stack()[1 if not parent else 2]
	file_name = os.path.split(file_path)[-1]
	calling_module = inspect.getmodule(frame).__name__
	debug_msg = separator.join(map(str, (file_name, function_name, line_number, u' '))) + str(msg)

	# Log to the calling module logger.  If calling_module is '__main__', use the root logger.
	logger = logging.getLogger(calling_module if calling_module != u'__main__' else '')
	logger.log(logging_level, debug_msg)
