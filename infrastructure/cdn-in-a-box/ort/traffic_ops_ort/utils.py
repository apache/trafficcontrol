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
This module contains miscellaneous utilities, typically dealing with string
manipulation or user input/output
"""
import logging
from sys import stderr
import requests
import typing

def getYesNoResponse(prmpt:str, default:str = None) -> bool:
	"""
	Utility function to get an interactive yes/no response to the prompt `prmpt`

	:param prmpt: The prompt to display to users
	:param default: The default response; should be one of ``'y'``, ``"yes"``, ``'n'`` or ``"no"``
		(case insensitive)
	:raises AttributeError: if 'prmpt' and/or 'default' is/are not strings
	:returns: the parsed response as a boolean
	"""
	if default:
		prmpt = prmpt.rstrip().rstrip(':') + '['+default+"]:"
	while True:
		choice = input(prmpt).lower()

		if choice in {'y', 'yes'}:
			return True
		if choice in {'n', 'no'}:
			return False
		if not choice and default is not None:
			return default.lower() in {'y', 'yes'}

		print("Please enter a yes/no response.", file=stderr)

def getTextResponse(uri:str, cookies:dict = None, verify:bool = True) -> str:
	"""
	Gets the plaintext response body of an HTTP ``GET`` request

	:param uri: The full path to a resource for the request
	:param cookies: An optional dictionary of cookie names mapped to values
	:param verify: If :const:`True`, the SSL keys used to communicate with the full URI will be
		verified

	:raises ConnectionError: when an error occurs trying to communicate with the server
	:raises ValueError: if the server's response cannot be interpreted as a UTF-8 string - e.g.
		when the response body is raw binary data but the response headers claim it's UTF-16
	"""
	logging.info("Getting plaintext response via 'HTTP GET %s'", uri)

	response = requests.get(uri, cookies=cookies, verify=verify)

	if response.status_code not in range(200, 300):
		logging.warning("Status code (%d) seems to indicate failure!", response.status_code)
		logging.debug("Response: %r\n%r", response.headers, response.content)

	return response.text

def getJSONResponse(uri:str, cookies:dict = None, verify:bool = True) -> dict:
	"""
	Retrieves a JSON object from some HTTP API

	:param uri: The URI to fetch
	:param cookies: A dictionary of cookie names mapped to values
	:param verify: If this is :const:`True`, the SSL keys will be verified during handshakes with
		'https' URIs
	:returns: The decoded JSON object
	:raises ConnectionError: when an error occurs trying to communicate with the server
	:raises ValueError: when the request completes successfully, but the response body
		does not represent a JSON-encoded object.
	"""

	logging.info("Getting JSON response via 'HTTP GET %s", uri)

	try:
		response = requests.get(uri, cookies=cookies, verify=verify)
	except (ValueError, ConnectionError, requests.exceptions.RequestException) as e:
		raise ConnectionError from e

	if response.status_code not in range(200, 300):
		logging.warning("Status code (%d) seems to indicate failure!", response.status_code)
		logging.debug("Response: %r\n%r", response.headers, response.content)

	return response.json()

def parse_multipart(raw: str) -> typing.List[typing.Tuple[str, str]]:
	"""
	Parses a multipart/mixed-type payload and returns each contiguous chunk.

	:param raw: The raw payload - without any HTTP status line.
	:returns: A list where each element is a tuple where the first element is a chunk of the message. All headers are discarded except 'Path', which is the second element of each tuple if it was found in the chunk.
	:raises: ValueError if the raw payload cannot be parsed as a multipart/mixed-type message.

	>>> testdata = '''MIME-Version: 1.0\\r
	... Content-Type: multipart/mixed; boundary="test"\\r
	... \\r
	... --test\\r
	... Content-Type: text/plain; charset=us-ascii\\r
	... Path: /path/to/ats/root/directory/etc/trafficserver/fname\\r
	... \\r
	... # A fake testing file that wasn't generated at all on some date
	... CONFIG proxy.config.way.too.many.period.separated.words INT 1
	...
	... --test\\r
	... Content-Type: text/plain; charset=utf8\\r
	... Path: /path/to/ats/root/directory/etc/trafficserver/othername\\r
	... \\r
	... # The same header again
	... CONFIG proxy.config.the.same.insane.chain.of.words.again.but.the.last.one.is.different INT 0
	...
	... --test--\\r
	... '''
	>>> output = parse_multipart(testdata)
	>>> print(output[0][0])
	# A fake testing file that wasn't generated at all on some date
	CONFIG proxy.config.way.too.many.period.separated.words INT 1

	>>> output[0][1]
	'/path/to/ats/root/directory/etc/trafficserver/fname'
	>>> print(output[1][0])
	# The same header again
	CONFIG proxy.config.the.same.insane.chain.of.words.again.but.the.last.one.is.different INT 0

	>>> output[1][1]
	'/path/to/ats/root/directory/etc/trafficserver/othername'
	"""
	try:
		hdr_index = raw.index("\r\n\r\n")
		headers = {line.split(':')[0].casefold(): line.split(':')[1] for line in raw[:hdr_index].splitlines()}
	except (IndexError, ValueError) as e:
		raise ValueError("Invalid or corrupt multipart header") from e

	ctype = headers.get("content-type")
	if not ctype:
		raise ValueError("Message is missing 'Content-Type' header")

	try:
		param_index = ctype.index(";")
		params = {param.split('=')[0].strip(): param.split('=')[1].strip() for param in ctype[param_index+1:].split(';')}
	except (IndexError, ValueError) as e:
		raise ValueError("Invalid or corrupt 'Content-Type' header") from e

	boundary = params.get("boundary", "").strip('"\'')
	if not boundary:
		raise ValueError("'Content-Type' header missing 'boundary' parameter")

	chunks = raw.split(f"--{boundary}")[1:] # ignore prologue
	if chunks[-1].strip() != "--":
		logging.warning("Final chunk appears invalid - possible bad message payload")
	else:
		chunks = chunks[:-1]

	ret = []
	for i, chunk in enumerate(chunks):
		try:
			hdr_index = chunk.index("\r\n\r\n")
			headers = {line.split(':')[0].casefold(): line.split(':')[1] for line in chunk[:hdr_index].splitlines() if line}
		except (IndexError, ValueError) as e:
			logging.debug("chunk: %s", chunk)
			raise ValueError(f"Chunk #{i} poorly formed") from e

		ret.append((chunk[hdr_index+4:].replace("\r","").strip(), headers.get("path").strip()))

	return ret
