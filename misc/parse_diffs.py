#!/usr/bin/env python3

import re
import typing

from enum import Enum

class Level(Enum):
	ERROR = "error"
	NOTICE = "notice" # notice not implemented (yet?)
	WARNING = "warning" # warning unused - diffs are assumed errors

class Annotation(typing.NamedTuple):
	level: Level
	file: str
	line: int
	message: str

	def __str__(self) -> str:
		return f"::{self.level.value} file={self.file},line={self.line}::{self.message}"

	def __repr__(self) -> str:
		return f"Annotation(level={self.level}, file='{self.file}', line={self.line})"

CHUNK_HEADER_PATTERN = re.compile(r"^@@ -\d+,\d+ \+(\d+),(\d+) @@")

def parseChunk(chunk: str, file: str) -> Annotation:
	"""
	parseChunk parses a single diff chunk and produces a corresponding annotation.

	>>> chunk = '''@@ -1,3 +1,3 @@
	...  {
	... +    "test": "quest"
	... -    "foo": "bar"
	...  }'''
	>>> ann = parseChunk(chunk, "test")
	>>> ann
	Annotation(level=Level.ERROR, file='test', line=2)
	>>> print(ann)
	::error file=test,line=2::Format Error
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

	content = "\n".join(["Format Error", "```diff", chunk, "```"])

	return Annotation(Level.ERROR, file, line, content)


FILE_HEADER_PATTERN = re.compile(r"^diff --git a/(.+) b/(.+)$")

def parseFile(contents: str) -> typing.List[Annotation]:
	"""
	parseFile parses the diff for a single file and returns the corresponding annotations.

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
	>>> anns = parseFile(file)
	>>> anns
	[Annotation(level=Level.ERROR, file='test', line=27), Annotation(level=Level.ERROR, file='test', line=88)]
	"""
	lines = contents.splitlines()
	if len(lines) < 5:
		raise ValueError(f"'{contents}' does not represent a file diff in git format")

	match = FILE_HEADER_PATTERN.match(lines[0])
	if not match or len(match.groups()) != 2:
		raise ValueError(f"invalid git diff file header: '{lines[0]}'")

	fname = match.groups()[1]

	lines = lines[4:]
	contents = "\n".join(lines)

	chunk = [lines[0]]
	annotations = []
	for line in lines[1:]:
		if CHUNK_HEADER_PATTERN.match(line):
			annotations.append(parseChunk("\n".join(chunk), fname))
			chunk = []

		chunk.append(line)

	if chunk:
		annotations.append(parseChunk("\n".join(chunk), fname))

	return annotations

