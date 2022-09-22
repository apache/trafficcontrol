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
	:file:`/api/{api-version}/` where ``api-version`` is :option:`--api-version`. This behavior can
	be disabled by using :option:`--raw-path`.

.. option:: DATA

	An optional positional argument that is a data payload to pass to the Traffic Ops server in the
	request body. If this is the absolute or relative path to a file, the contents of the file will
	instead be read and used as the request payload.

.. option:: -h, --help

	Print usage information and exit

.. option:: -a API_VERSION, --api-version API_VERSION

	Specifies the version of the Traffic Ops API that will be used for the request. Has no effect if
	:option:`--raw-path` is used. (Default: 4.1)

.. option:: -f, --full

	Output the full HTTP exchange including request method line, request headers, request body (if
	any), response status line, and response headers (as well as the response body, if any). This is
	equivalent to using :option:`--request-headers`, :option:`--request-payload`, and
	:option:`--response-headers` at the same time, and those options will have no effect if given.
	(Default: false)

.. option:: -k, --insecure

	Do not verify SSL certificates - typically useful for making requests to development or testing
	servers as they frequently have self-signed certificates. (Default: false)

.. option:: -p, --pretty

	Pretty-print any payloads that are output as formatted JSON. Has no effect on plaintext payloads.
	Uses tab characters for indentation. (Default: false)

.. option:: -r, --raw-path

	Request exactly :option:`PATH`; do not preface the request path with :file:`/api/{api_version}`.
	This effectively means that :option:`--api-version` will have no effect. (Default: false)

.. option:: -v, --version

	Print version information and exit.

.. option:: --request-headers

	Output the request method line and any and all request headers. (Default: false)

.. option:: --request-payload

	Output the request body if any was sent. Will attempt to pretty-print the body as JSON if
	:option:`--pretty` is used. (Default: false)

.. option:: --response-headers

	Output the response status line and any and all response headers. (Default: false)

.. option:: --to-url URL

	The :abbr:`FQDN (Fully Qualified Domain Name)` and optionally the port and scheme of the Traffic
	Ops server. This will override :envvar:`TO_URL`. The format is the same as for :envvar:`TO_URL`.
	(Default: uses the value of :envvar:`TO_URL`)

.. option:: --to-password PASSWORD

	The password to use when authenticating to Traffic Ops. Overrides :envvar:`TO_PASSWORD`.
	(Default: uses the value of :envvar:`TO_PASSWORD`)

.. option:: --to-user USERNAME

	The username to use when connecting to Traffic Ops. Overrides :envvar:`TO_USER`. (Default: uses
	the value of :envvar:`TO_USER`)

Environment Variables
---------------------
If defined, :program:`toaccess` scripts will use the :envvar:`TO_URL`, :envvar:`TO_USER`, and
:envvar`TO_PASSWORD` environment variables to define their connection to and authentication with the
Traffic Ops server. Typically, setting these is easier than using the long options :option:`--to-url`,
:option:`--to-user`, and :option:`--to-password` on every invocation.

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
	passing ``https://test:`` to :option:`--to-url` will cause an exit code of ``2``, not ``32``.
4
	An HTTP client error occurred. The HTTP stack will be printed to stdout as indicated by other
	options - meaning by default it will only print the response payload if one was given, but will
	respect options like e.g. :option:`--request-payload` as well as
	:option:`-p`/:option:`--pretty`.
5
	An HTTP server error occurred. The HTTP stack will be printed to stdout as indicated by other
	options - meaning by default it will only print the response payload if one was given, but will
	respect options like e.g. :option:`--request-payload` as well as
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
import json
import logging
import os
import sys
from urllib.parse import urlparse

from trafficops.restapi import LoginError, OperationError, InvalidJSONError
from trafficops.tosession import TOSession
from trafficops.__version__ import __version__

from requests.exceptions import RequestException

l = logging.getLogger()
l.disabled = True
logging.basicConfig(level=logging.CRITICAL+1)

def output(r, pretty, request_header, response_header, request_payload, indent = '\t'):
	"""
	Prints the passed response object in a format consistent with the other parameters.

	:param r: The :mod:`requests` response object being printed
	:param pretty: If :const:`True`, attempt to pretty-print payloads as JSON
	:param request_header: If :const:`True`, print request line and request headers
	:param response_header: If :const:`True`, print response line and response headers
	:param request_payload: If :const:`True`, print the request payload
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
	                               "rather than (respectively) the '--to-url', '--to-user' and "
	                               "'--to-password' command-line flags. Those flags are only "
	                               "required when said environment variables are not defined.\n"
	                               "%(prog)s will exit with a success provided a response was "
	                               "received and the status code of said response was less than 400. "
	                               "The exit code will be 1 if command line arguments cannot be "
	                               "parsed or authentication with the Traffic Ops server fails. "
	                               "In the event of some unknown error occurring when waiting for a "
	                               "response, the exit code will be 2. If the server responds with "
	                               "a status code indicating a client or server error, that status "
	                               "code will be used as the exit code."))

	parser.add_argument("--to-url",
	                    type=str,
	                    help=("The fully qualified domain name of the Traffic Ops server. Overrides "
	                         "'$TO_URL'. The format for both the environment variable and the flag "
	                         "is '[scheme]hostname[:port]'. That is, ports should be specified here, "
	                         "and they need not start with 'http://' or 'https://'. HTTPS is the "
	                         "assumed protocol unless the scheme _is_ provided and is 'http://'."))
	parser.add_argument("--to-user",
	                    type=str,
	                    help="The username to use when connecting to Traffic Ops. Overrides '$TO_USER")
	parser.add_argument("--to-password",
	                    type=str,
	                    help="The password to use when authenticating to Traffic Ops. Overrides '$TO_PASSWORD'")
	parser.add_argument("-k", "--insecure", action="store_true", help="Do not verify SSL certificates")
	parser.add_argument("-f", "--full",
	                    action="store_true",
	                    help=("Also output HTTP request/response lines and headers, and request payload. "
	                         "This is equivalent to using '--request-headers', '--response-headers' "
	                         "and '--request-payload' at the same time."))
	parser.add_argument("--request-headers",
	                    action="store_true",
	                    help="Output request method line and headers")
	parser.add_argument("--response-headers",
	                    action="store_true",
	                    help="Output response status line and headers")
	parser.add_argument("--request-payload",
	                    action="store_true",
	                    help="Output request payload (will try to pretty-print if '--pretty' is given)")
	parser.add_argument("-r", "--raw-path",
	                    action="store_true",
	                    help="Request exactly PATH; it won't be prefaced with '/api/{{api-version}}/")
	parser.add_argument("-a", "--api-version",
	                    type=float,
	                    default=4.1,
	                    help="Specify the API version to request against")
	parser.add_argument("-p", "--pretty",
	                    action="store_true",
	                    help=("Pretty-print payloads as JSON. "
	                         "Note that this will make Content-Type headers \"wrong\", in general"))
	parser.add_argument("-v", "--version",
	                    action="version",
	                    help="Print version information and exit",
	                    version="%(prog)s v"+__version__)
	parser.add_argument("PATH", help="The path to the resource being requested - omit '/api/2.x'")
	parser.add_argument("DATA",
	                    help=("An optional data string to pass with the request. If this is a "
	                         "filename, the contents of the file will be sent instead."),
	                    nargs='?')


	args = parser.parse_args()

	try:
		to_host = args.to_url if args.to_url else os.environ["TO_URL"]
	except KeyError as e:
		raise KeyError("Traffic Ops hostname not set! Set the TO_URL environment variable or use "\
		               "'--to-url'.") from e

	original_to_host = to_host
	to_host = urlparse(to_host, scheme="https")
	useSSL = to_host.scheme.lower() == "https"
	to_port = to_host.port
	if to_port is None:
		if useSSL:
			to_port = 443
		else:
			to_port = 80

	to_host = to_host.hostname
	if not to_host:
		raise KeyError(f"Invalid URL/host for Traffic Ops: '{original_to_host}'")

	s = TOSession(to_host,
	              host_port=to_port,
	              ssl=useSSL,
	              api_version=f"{args.api_version:.1f}",
	              verify_cert=not args.insecure)

	data = args.DATA
	if data and os.path.isfile(data):
		with open(data) as f:
			data = f.read()

	if isinstance(data, str):
		data = data.encode()

	try:
		to_user = args.to_user if args.to_user else os.environ["TO_USER"]
	except KeyError as e:
		raise KeyError("Traffic Ops user not set! Set the TO_USER environment variable or use "\
		               "'--to-user'.") from e

	try:
		to_passwd = args.to_password if args.to_password else os.environ["TO_PASSWORD"]
	except KeyError as e:
		raise KeyError("Traffic Ops password not set! Set the TO_PASSWORD environment variable or "\
		               "use '--to-password'") from e

	# TOSession objects return LoginError when certs are invalid, OperationError when
	# login actually fails
	try:
		s.login(to_user, to_passwd)
	except LoginError as e:
		raise PermissionError(
			"certificate verification failed, the system may have a self-signed certificate - try using -k/--insecure"
		) from e
	except (OperationError, InvalidJSONError) as e:
		raise PermissionError(e) from e
	except RequestException as e:
		raise ConnectionError("Traffic Ops host not found: Name or service not known") from e

	return (s,
	       args.PATH,
	       data,
	       (
	         args.request_headers or args.full,
	         args.response_headers or args.full,
	         args.request_payload or args.full
	       ),
	       args.raw_path,
	       args.pretty)

def request(method):
	"""
	All of the scripts wind up calling this function to handle their common functionality.

	:param method: The name of the request method to use (case-insensitive)
	:returns: The program's exit code
	"""
	try:
		s, path, data, full, raw, pretty = parse_arguments("to%s" % method)
	except (PermissionError, KeyError, ConnectionError) as e:
		print(e, file=sys.stderr)
		return 1

	if raw:
		path = '/'.join((s.to_url.rstrip('/'), path.lstrip('/')))
	else:
		path = '/'.join((s.base_url.rstrip('/'), path.lstrip('/')))

	try:
		if data is not None:
			r = s._session.request(method, path, data=data)
		else:
			r = s._session.request(method, path)
	except (RequestException, ValueError) as e:
		print("Error occurred: ", e, file=sys.stderr)
		return 2

	output(r, pretty, *full)
	return 0 if r.status_code < 400 else r.status_code // 100

def get():
	"""
	Entry point for :program:`toget`

	:returns: The program's exit code
	"""
	return request("get")

def put():
	"""
	Entry point for :program:`toput`

	:returns: The program's exit code
	"""
	return request("put")

def post():
	"""
	Entry point for :program:`topost`

	:returns: The program's exit code
	"""
	return request("post")

def delete():
	"""
	Entry point for :program:`todelete`

	:returns: The program's exit code
	"""
	return request("delete")

def options():
	"""
	Entry point for :program:`tooptions`

	:returns: The program's exit code
	"""
	return request("options")

def head():
	"""
	Entry point for :program:`tohead`

	:returns: The program's exit code
	"""
	return request("head")

def patch():
	"""
	Entry point for :program:`topatch`

	:returns: The program's exit code
	"""
	return request("patch")
