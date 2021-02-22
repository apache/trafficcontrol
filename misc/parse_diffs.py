#!/usr/bin/env python3

"""
This script parses git diffs and outputs a set of GitHub Action annotations where each section of
each file is an annotation containing the diff chunk as a message. All annotations are error-level.
"""

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

import re
import sys
import typing

from enum import Enum

class Level(Enum):
	"""
	Level encodes the level of an annotation.
	"""

	error = "error"
	notice = "notice" # notice not implemented (yet?)
	warn = "warning" # warning unused - diffs are assumed errors

class Annotation(typing.NamedTuple):
	"""
	Annotation represents a GitHub Actions Annotation. To preview the GitHub Annotation string,
	coerce it to a string value. For a fully sanitized, printable Annotation, use the sanitize
	method.

	>>> print(Annotation(level=Level.error, file='test', line=1, message='test'))
	::error file=test,line=1::test
	"""
	level: Level
	file: str
	line: int
	message: str

	def __str__(self) -> str:
		msg = self.message.replace("]", "%5D").replace(";", "%3B")
		return f"::{self.level.value} file={self.file},line={self.line}::{msg}"

	def __repr__(self) -> str:
		return f"Annotation(level={self.level}, file='{self.file}', line={self.line})"

	def sanitize(self) -> str:
		"""
		sanitize returns a sanitized string, suitable for GHA, but not for humans to read because
		all newlines have been replaced with their URL percent-encoded code points.
		"""
		return str(self).replace("\r", "%0D").replace("\n", "%0A")

CHUNK_HEADER_PATTERN = re.compile(r"^@@ -\d+,\d+ \+(\d+),(\d+) @@")

def parse_chunk(chunk: str, file: str) -> Annotation:
	"""
	parse_chunk parses a single diff chunk and produces a corresponding annotation.

	>>> chunk = '''@@ -1,3 +1,3 @@
	...  {
	... +    "test": "quest"
	... -    "foo": "bar"
	...  }'''
	>>> ann = parse_chunk(chunk, "test")
	>>> ann
	Annotation(level=Level.error, file='test', line=2)
	>>> print(ann)
	::error file=test,line=2::Format error
	```diff
	@@ -1,3 +1,3 @@
	 {
	+    "test": "quest"
	-    "foo": "bar"
	 }
	```
	"""
	lines = chunk.splitlines()
	if len(lines) < 2:
		raise ValueError(f"invalid diff chunk: {chunk}")

	match = CHUNK_HEADER_PATTERN.match(lines[0])
	if not match or len(match.groups()) != 2:
		raise ValueError(f"invalid diff chunk header: {lines[0]}")

	line = int(match.groups()[0]) + int(match.groups()[1])//2

	content = "\n".join(["Format error", "```diff", chunk, "```"])

	return Annotation(Level.error, file, line, content)


FILE_HEADER_PATTERN = re.compile(r"^diff --git a/(.+) b/(.+)$")

def parse_file(contents: str) -> typing.List[Annotation]:
	"""
	parse_file parses the diff for a single file and returns the corresponding annotations.

	>>> file = '''diff --git a/test b/test
	... index c9072dcb7..2b7686061 100644
	... --- a/test
	... +++ b/test
	... @@ -24,7 +24,7 @@ package tc
	...  // in: body
	...  type ASNsResponse struct {
	...         // in: body
	... -       Response []ASN `json:"response"`
	... +Response []ASN `json:"response"`
	...        Alerts
	... }
	...
	... @@ -85,7 +85,7 @@ type ASNNullable struct {
	...         // ID of the ASN
	...         //
	...         // required: true
	... -       ID *int `json:"id" db:"id"`
	... +ID *int `json:"id" db:"id"`
	...
	...         // LastUpdated
	...         //
	... '''
	>>> anns = parse_file(file)
	>>> len(anns)
	2
	>>> anns[0]
	Annotation(level=Level.error, file='test', line=27)
	>>> anns[1]
	Annotation(level=Level.error, file='test', line=88)
	"""
	lines = contents.splitlines()
	if len(lines) < 5:
		raise ValueError(f"'{contents}' does not represent a file diff in git format")

	match = FILE_HEADER_PATTERN.match(lines[0])
	if not match or len(match.groups()) != 2:
		raise ValueError(f"invalid git diff file header: '{lines[0]}'")

	fname = match.groups()[1]

	lines = lines[4:]

	chunk = [lines[0]]
	annotations = []
	for line in lines[1:]:
		if CHUNK_HEADER_PATTERN.match(line):
			annotations.append(parse_chunk("\n".join(chunk), fname))
			chunk = []

		chunk.append(line)

	if chunk:
		annotations.append(parse_chunk("\n".join(chunk), fname))

	return annotations

def parse_diff(diff: str) -> typing.List[Annotation]:
	"""
	parse_diff parses a git diff output and returns the corresponding annotations.

	>>> diff= '''diff --git a/test b/test
	... index c9072dcb7..2b7686061 100644
	... --- a/test
	... +++ b/test
	... @@ -24,7 +24,7 @@ package tc
	...  // in: body
	...  type ASNsResponse struct {
	...         // in: body
	... -       Response []ASN `json:"response"`
	... +Response []ASN `json:"response"`
	...         Alerts
	...  }
	...
	... @@ -85,7 +85,7 @@ type ASNNullable struct {
	...         // ID of the ASN
	...         //
	...         // required: true
	... -       ID *int `json:"id" db:"id"`
	... +ID *int `json:"id" db:"id"`
	...
	...         // LastUpdated
	...         //
	... diff --git a/quest b/quest
	... index 283901f14..0c1e2b0c1 100644
	... --- a/quest
	... +++ b/quest
	... @@ -1,7 +1,7 @@
	...  package tc
	...
	...  import (
	... -       "database/sql"
	... +"database/sql"
	...  )
	...
	...  /*
	... '''
	>>> anns = parse_diff(diff)
	>>> len(anns)
	3
	>>> anns[0]
	Annotation(level=Level.error, file='test', line=27)
	>>> anns[1]
	Annotation(level=Level.error, file='test', line=88)
	>>> anns[2]
	Annotation(level=Level.error, file='quest', line=4)
	"""

	lines = diff.splitlines()
	if len(lines) < 5:
		raise ValueError(f"'{diff}' does not represent a git diff")

	match = FILE_HEADER_PATTERN.match(lines[0])
	if not match or len(match.groups()) != 2:
		raise ValueError(f"invalid git diff file header: '{lines[0]}''")

	file = lines[:4]
	lines = lines[4:]
	annotations = []
	for line in lines:
		if FILE_HEADER_PATTERN.match(line):
			annotations += parse_file("\n".join(file))
			file = []
		file.append(line)

	if file:
		annotations += parse_file("\n".join(file))

	return annotations

def main() -> int:
	"""
	Runs the main program, based on the passed-in arguments.

	Returns an exit code based on success or failure of the script - NOT based
	on the presence of any error-level annotations.
	"""
	try:
		print(*(x.sanitize() for x in parse_diff(sys.stdin.read())), sep="\n")
		return 0
	except ValueError as e:
		print("error:", e, file=sys.stderr)
		return 1
	except OSError as e:
		print("error reading input:", e, file=sys.stderr)
		return 2

if __name__ == "__main__":
	sys.exit(main())
