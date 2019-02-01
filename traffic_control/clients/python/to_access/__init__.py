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
.. _toaccess:

.. program:: toaccess

``toaccess``
============
This module provides a set of functions meant to provide ease-of-use functionality for interacting
with the Traffic Ops API. It provides scripts named :file:`to{method}` where `method` is the name of
an HTTP method (in lowercase). Collectively they are referred to as :program:`toaccess` Implemented
methods thus far are:

- delete
- head
- get
- options
- patch
- post
- put

Arguments and Flags
-------------------
.. option:: PATH

	This is the request path. By default, whatever is passed is considered to be relative to
	:file:`/api/{api_version}/` where ``api_version`` is :option:`--api_version`. This behavior can
	be disabled by using :option:`--raw_path`.

.. option:: DATA

	An optional positional argument that is a data payload to pass to the Traffic Ops server in the
	request body. If this is the absolute or relative path to a file, the contents of the file will
	instead be read and used as the request payload.

.. option:: -h, --help

	Print usage information and exit

.. option:: -a API_VERSION, --api_version API_VERSION

	Specifies the version of the Traffic Ops API that will be used for the request. Has no effect if
	:option:`--raw_path` is used. (Default: 1.3)

.. option:: -f, --full

	Output the full HTTP exchange including request method line, request headers, request body (if
	any), response status line, and response headers (as well as the response body, if any). This is
	equivalent to using :option:`--request_headers`, :option:`--request_payload`, and
	:option:`--response_headers` at the same time, and those options will have no effect if given.
	(Default: false)

.. option:: -k, --insecure

	Do not verify SSL certificates - typically useful for making requests to development or testing
	servers as they frequently have self-signed certificates. (Default: false)

.. option:: -p, --pretty

	Pretty-print any payloads that are output as formatted JSON. Has no effect on plaintext payloads.
	Uses tab characters for indentation. (Default: false)

.. option:: -r, --raw_path

	Request exactly :option:`PATH`; do not preface the request path with :file:`/api/{api_version}`.
	This effectively means that :option:`--api_version` will have no effect. (Default: false)

.. option:: --request_headers

	Output the request method line and any and all request headers. (Default: false)

.. option:: --request_payload

	Output the request body if any was sent. Will attempt to pretty-print the body as JSON if
	:option:`--pretty` is used. (Default: false)

.. option:: --response_headers

	Output the response status line and any and all response headers. (Default: false)

.. option:: --to_url URL

	The :abbr:`FQDN (Fully Qualified Domain Name)` and optionally the port and scheme of the Traffic
	Ops server. This will override :envvar:`TO_URL`. The format is the same as for :envvar:`TO_URL`.
	(Default: uses the value of :envvar:`TO_URL`)

.. option:: --to_password PASSWORD

	The password to use when authenticating to Traffic Ops. Overrides :envvar:`TO_PASSWORD`.
	(Default: uses the value of :envvar:`TO_PASSWORD`)

.. option:: --to_user USERNAME

	The username to use when connecting to Traffic Ops. Overrides :envvar:`TO_USER`. (Default: uses
	the value of :envvar:`TO_USER`)

Environment Variables
---------------------
If defined, :program:`toaccess` scripts will use these environment variables to define their
connection to and authentication with the Traffic Ops server. Typically, setting these is easier
than using the long options :option:`--to_url`, :option:`--to_user`, and :option:`--to_password` on
every invocation.

.. envvar:: TO_PASSWORD

	Will be used to authenticate the user defined by either :option:`--to_user` or :envvar:`TO_USER`.

.. envvar:: TO_URL

	The :abbr:`FQDN (Fully Qualified Domain Name)` of the Traffic Ops server to which the script
	will connect. The format of this should be :file:`[{http or https}://]{hostname}[:{port}]`. Note
	that this may optionally start with ``http://`` or ``https://`` (case insensitive), but
	typically this is unnecessary. Also notice that the port number may be specified, though again
	this isn't usually required. All :program:`toaccess` scripts will assume that port 443 should be
	used unless otherwise specified. They will further assume that the protocol is HTTPS unless
	:envvar:`TO_URL` (or :option:`--to_url`) starts with ``http://``, in which case the default port
	will also be set to 80 unless otherwise specified in the URL.

.. envvar:: TO_USER

	The name of the user as whom to connect to the Traffic Ops server. Overriden by
	:option:`--to_user`.

Exit Codes
----------
The exit code of a :program:`toaccess` script can sometimes be used by the caller to determine what
the result of calling the script was without needing to parse the output. The exit codes used are:

0
	The command executed successfully, and the result is on STDOUT.
1
	Typically this exit code means that an error was encountered when parsing positional command
	line arguments. However, this is also the exit code used by most Python interpreters to signal
	an unhandled exception.
2
	Signifies a runtime error that caused the request to fail - this is **not** generally indicative
	of an HTTP client or server error, but rather an underlying issue connecting to or
	authenticating with Traffic Ops. This is distinct from an exit code of ``32`` in that the
	*format* of the arguments was correct, but there was some problem with the *value*. For example,
	passing ``https://test:`` to :option:`--to_url` will cause an exit code of ``2``, not ``32``.
4
	An HTTP client error occurred. The HTTP stack will be printed to stdout as indicated by other
	options - meaning by default it will only print the response payload if one was given, but will
	respect options like e.g. :option:`--request_payload` as well as
	:option:`-p`/:option:`--pretty`.
5
	An HTTP server error occurred. The HTTP stack will be printed to stdout as indicated by other
	options - meaning by default it will only print the response payload if one was given, but will
	respect options like e.g. :option:`--request_payload` as well as
	:option:`-p`/:option:`--pretty`.
32
	This is the error code emitted by Python's :mod:`argparse` module when the passed arguments
	could not be parsed successfully.

.. note:: The way exit codes ``4`` and ``5`` are implemented is by returning the status code of the
	HTTP request divided by 100 whenever it is at least 400. This means that if the Traffic Ops
	server ever started returning e.g. 700 status codes, the exit code of the script would be 7.


Module Reference
================

"""
from __future__ import print_function, raise_from

import json
import logging
import os
import sys

from trafficops.restapi import LoginError, OperationError, InvalidJSONError
from trafficops.tosession import TOSession

l = logging.getLogger()
l.disabled = True
logging.basicConfig(level=logging.CRITICAL+1)

#: The full path to a file used to store the user's Mojolicious athentication cookie (currently unused)
COOKIEFILE = "" #os.path.expanduser(os.path.join("~", ".to-auth.cookie"))

def set_cookie(cookie):
	"""
	Writes the passed cookie to the :data:`COOKIEFILE` file for later use.

	.. warning::
		Currently this function is never used, and it depends on the value of :data:`COOKIEFILE`,
		which at the moment is being set to an empty string.
	"""
	with open(COOKIEFILE, "w") as f:
		f.write(cookie)

def output(r, pretty, request_header, response_header, request_payload, indent = '\t'):
	"""
	Prints the passed response object in a format consistent with the other parameters.

	:param r: The :mod:`requests` response object being printed
	:param pretty: If :const:`True`, attempt to pretty-print payloads as JSON
	:param reqHeader: If :const:`True`, print request line and request headers
	:param respHeader: If :const:`True`, print response line and response headers
	:param reqPayload: If :const:`True`, print the request payload
	:param indent: An optional number of spaces for pretty-printing indentation (default is the tab character)
	"""
	if request_header:
		print(r.request.method, r.request.path_url, "HTTP/1.1")
		for h,v in r.request.headers.items():
			print("%s:" % h, v)
		print()

	if request_payload and r.request.body:
		try:
			result = r.request.body if not pretty else json.dumps(json.loads(r.request.body))
		except ValueError:
			result = r.request.body
		print(result, end="\n\n")

	if response_header:
		print("HTTP/1.1", r.status_code, end="")
		print(" "+r.reason if r.reason else "")
		for h,v in r.headers.items():
			print("%s:" % h, v)
		print()

	try:
		result = r.text if not pretty else json.dumps(r.json(), indent=indent)
	except ValueError:
		result = r.text
	print(result)

def parse_arguments(program):
	"""
	A common-use function that parses the command line arguments.

	:param program: The name of the program being run - used for usage informational output
	:returns: The Traffic Ops HTTP session object, the requested path, any data to be sent, an output
	          format specification, whether or not the path is raw, and whether or not output should
	          be prettified
	"""
	from argparse import ArgumentParser, ArgumentDefaultsHelpFormatter
	parser = ArgumentParser(prog=program,
	                        formatter_class=ArgumentDefaultsHelpFormatter,
	                        description="A helper program for interfacing with the Traffic Ops API",
	                        epilog=("Typically, one will want to connect and authenticate by defining "
	                               "the 'TO_URL', 'TO_USER' and 'TO_PASSWORD' environment variables "
	                               "rather than (respectively) the '--to_url', '--to_user' and "
	                               "'--to_password' command-line flags. Those flags are only "
	                               "required when said environment variables are not defined.\n"
	                               "%(prog)s will exit with a success provided a response was "
	                               "received and the status code of said response was less than 400. "
	                               "The exit code will be 1 if command line arguments cannot be "
	                               "parsed or authentication with the Traffic Ops server fails. "
	                               "In the event of some unknown error occurring when waiting for a "
	                               "response, the exit code will be 2. If the server responds with "
	                               "a status code indicating a client or server error, that status "
	                               "code will be used as the exit code."))

	parser.add_argument("--to_url",
	                    type=str,
	                    help=("The fully qualified domain name of the Traffic Ops server. Overrides "
	                         "'$TO_URL'. The format for both the environment variable and the flag "
	                         "is '[scheme]hostname[:port]'. That is, ports should be specified here, "
	                         "and they need not start with 'http://' or 'https://'. HTTPS is the "
	                         "assumed protocol unless the scheme _is_ provided and is 'http://'."))
	parser.add_argument("--to_user",
	                    type=str,
	                    help="The username to use when connecting to Traffic Ops. Overrides '$TO_USER")
	parser.add_argument("--to_password",
	                    type=str,
	                    help="The password to use when authenticating to Traffic Ops. Overrides '$TO_PASSWORD'")
	parser.add_argument("-k", "--insecure", action="store_true", help="Do not verify SSL certificates")
	parser.add_argument("-f", "--full",
	                    action="store_true",
	                    help=("Also output HTTP request/response lines and headers, and request payload. "
	                         "This is equivalent to using '--request_headers', '--response_headers' "
	                         "and '--request_payload' at the same time."))
	parser.add_argument("--request_headers",
	                    action="store_true",
	                    help="Output request method line and headers")
	parser.add_argument("--response_headers",
	                    action="store_true",
	                    help="Output response status line and headers")
	parser.add_argument("--request_payload",
	                    action="store_true",
	                    help="Output request payload (will try to pretty-print if '--pretty' is given)")
	parser.add_argument("-r", "--raw_path",
	                    action="store_true",
	                    help="Request exactly PATH; it won't be prefaced with '/api/{{api_version}}/")
	parser.add_argument("-a", "--api_version",
	                    type=float,
	                    default=1.3,
	                    help="Specify the API version to request against")
	parser.add_argument("-p", "--pretty",
	                    action="store_true",
	                    help=("Pretty-print payloads as JSON. "
	                         "Note that this will make Content-Type headers \"wrong\", in general"))
	parser.add_argument("PATH", help="The path to the resource being requested - omit '/api/1.x'")
	parser.add_argument("DATA",
	                    help=("An optional data string to pass with the request. If this is a "
	                         "filename, the contents of the file will be sent instead."),
	                    nargs='?')


	args = parser.parse_args()

	try:
		to_host = args.to_host if args.to_host else os.environ["TO_URL"]
	except KeyError as e:
		raise KeyError("Traffic Ops hostname not set! Set the TO_URL environment variable or use "\
		               "'--to_url'.") from e

	useSSL = True
	to_port = 443
	if to_host.lower().startswith("http://"):
		to_host = to_host[7:]
		useSSL = False
		to_port = 80
	elif to_host.lower().startswith("https://"):
		to_host = to_host[8:]

	portpos = to_host.find(':')
	if portpos > 0:
		to_port = int(to_host[portpos+1:])
		to_host = to_host[:portpos]

	s = TOSession(to_host,
	              host_port=to_port,
	              ssl=useSSL,
	              api_version=str(args.api_version),
	              verify_cert=not args.insecure)
	# s.create()

	# The TOSession object's methods will handle '/' stripping for us, so it's only necessary with
	# raw paths.
	path = args.PATH if not args.raw_path else '/'.join((s.to_url.rstrip('/'), args.PATH.lstrip('/')))

	data = args.DATA
	if data and os.path.isfile(data):
		with open(data) as f:
			data = f.read()

	if isinstance(data, str):
		data = data.encode()

	# if os.path.isfile(COOKIEFILE):
	# 	with open(COOKIEFILE) as f:
	# 		s._session.cookies.set("mojolicious", f.read())
	# 	return s, path, data, args.full, args.raw_path

	try:
		to_user = args.to_user if args.to_user else os.environ["TO_USER"]
	except KeyError as e:
		raise KeyError("Traffic Ops user not set! Set the TO_USER environment variable or use "\
		               "'--to_user'.") from e

	try:
		to_passwd = args.to_password if args.to_password else os.environ["TO_PASSWORD"]
	except KeyError as e:
		raise KeyError("Traffic Ops password not set! Set the TO_PASSWORD environment variable or "\
		               "use '--to_password'") from e

	try:
		s.login(to_user, to_passwd)
	except (OperationError, InvalidJSONError, LoginError) as e:
		raise PermissionError from e

	return (s,
	       path,
	       data,
	       (
	         args.request_headers or args.full,
	         args.response_headers or args.full,
	         args.request_payload or args.full
	       ),
	       args.raw_path,
	       args.pretty)


def get():
	"""
	Entry point for :program:`toget`

	:returns: The program's exit code
	"""
	try:
		s, path, data, full, raw, pretty = parse_arguments("toget")
	except (PermissionError, KeyError) as e:
		print(e, file=sys.stderr)
		return 1

	try:
		if raw:
			if data is not None:
				r = s._session.get(path, data=data)
			else:
				r = s._session.get(path)
		elif data is not None:
			r = s.get(path, data=data)[1]
		else:
			r = s.get(path)[1]
	except (OperationError, InvalidJSONError) as e:
		if e.resp is not None:
			r = e.resp
		else:
			print("Error occurred: ", e, file=sys.stderr)
			return 2

	output(r, pretty, *full)
	# set_cookie(s._session.cookies.get("mojolicious"))
	return 0 if r.status_code < 400 else r.status_code // 100

def put():
	"""
	Entry point for :program:`toput`

	:returns: The program's exit code
	"""
	try:
		s, path, data, full, raw, pretty = parse_arguments("toput")
	except (PermissionError, KeyError) as e:
		print(e, file=sys.stderr)
		return 1

	try:
		if raw:
			if data is not None:
				r = s._session.put(path, data=data)
			else:
				r = s._session.put(path)
		elif data is not None:
			r = s.put(path, data=data)[1]
		else:
			r = s.put(path)[1]
	except (OperationError, InvalidJSONError) as e:
		if e.resp is not None:
			r = e.resp
		else:
			print("Error occurred: ", e, file=sys.stderr)
			return 2

	output(r, pretty, *full)
	# set_cookie(s._session.cookies.get("mojolicious"))
	return 0 if r.status_code < 400 else r.status_code // 100

def post():
	"""
	Entry point for :program:`topost`

	:returns: The program's exit code
	"""
	try:
		s, path, data, full, raw, pretty = parse_arguments("topost")
	except (PermissionError, KeyError) as e:
		print(e, file=sys.stderr)
		return 1

	try:
		if raw:
			if data is not None:
				r = s._session.post(path, data=data)
			else:
				r = s._session.post(path)
		elif data is not None:
			r = s.post(path, data=data)[1]
		else:
			r = s.post(path)[1]
	except (OperationError, InvalidJSONError) as e:
		if e.resp is not None:
			r = e.resp
		else:
			print("Error occurred: ", e, file=sys.stderr)
			return 2

	output(r, pretty, *full)
	# set_cookie(s._session.cookies.get("mojolicious"))
	return 0 if r.status_code < 400 else r.status_code // 100

def delete():
	"""
	Entry point for :program:`todelete`

	:returns: The program's exit code
	"""
	try:
		s, path, data, full, raw, pretty = parse_arguments("todelete")
	except (PermissionError, KeyError) as e:
		print(e, file=sys.stderr)
		return 1

	try:
		if raw:
			if data is not None:
				r = s._session.delete(path, data=data)
			else:
				r = s._session.delete(path)
		elif data is not None:
			r = s.delete(path, data=data)[1]
		else:
			r = s.delete(path)[1]
	except (OperationError, InvalidJSONError) as e:
		if e.resp is not None:
			r = e.resp
		else:
			print("Error occurred: ", e, file=sys.stderr)
			return 2

	output(r, pretty, *full)
	# set_cookie(s._session.cookies.get("mojolicious"))
	return 0 if r.status_code < 400 else r.status_code // 100

def options():
	"""
	Entry point for :program:`tooptions`

	:returns: The program's exit code
	"""
	from functools import partial

	try:
		s, path, data, full, raw, pretty = parse_arguments("tooptions")
	except (PermissionError, KeyError) as e:
		print(e, file=sys.stderr)
		return 1

	try:
		if raw:
			if data is not None:
				r = s._session.options(path, data=data)
			else: r = s._session.options(path)
		elif data is not None:
			r = s.options(path, data=data)[1]
		else:
			r = s.options(path)[1]
	except (OperationError, InvalidJSONError) as e:
		if e.resp is not None:
			r = e.resp
		else:
			print("Error occurred: ", e, file=sys.stderr)
			return 2

	output(r, pretty, *full)
	# set_cookie(s._session.cookies.get("mojolicious"))
	return 0 if r.status_code < 400 else r.status_code // 100

def head():
	"""
	Entry point for :program:`tohead`

	:returns: The program's exit code
	"""
	try:
		s, path, data, full, raw, pretty = parse_arguments("tohead")
	except (PermissionError, KeyError) as e:
		print(e, file=sys.stderr)
		return 1

	try:
		if raw:
			if data is not None:
				r = s._session.head(path, data=data)
			else:
				r = s._session.head(path)
		elif data is not None:
			r = s.head(path, data=data)[1]
		else:
			r = s.head(path)[1]
	except (OperationError, InvalidJSONError) as e:
		if e.resp is not None:
			r = e.resp
		else:
			print("Error occurred: ", e, file=sys.stderr)
			return 2

	output(r, pretty, *full)
	# set_cookie(s._session.cookies.get("mojolicious"))
	return 0 if r.status_code < 400 else r.status_code // 100

def patch():
	"""
	Entry point for :program:`topatch`

	:returns: The program's exit code
	"""
	try:
		s, path, data, full, raw, pretty = parse_arguments("topatch")
	except (PermissionError, KeyError) as e:
		print(e, file=sys.stderr)
		return 1

	try:
		if raw:
			if data is not None:
				r = s._session.patch(path, data=data)
			else:
				r = s._session.patch(path)
		elif data is not None:
			r = s.patch(path, data=data)[1]
		else:
			r = s.patch(path)[1]
	except (OperationError, InvalidJSONError) as e:
		if e.resp is not None:
			r = e.resp
		else:
			print("Error occurred: ", e, file=sys.stderr)
			return 2

	output(r, pretty, *full)
	# set_cookie(s._session.cookies.get("mojolicious"))
	return 0 if r.status_code < 400 else r.status_code // 100
