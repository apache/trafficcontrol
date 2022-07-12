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
The unit testing module for the pr_to_update_go package.

Note that this module merely reiterates the docstring tests - and only covers a
proper subset of the test cases covered by said docstring tests.
"""

from unittest import TestCase

from pr_to_update_go.go_pr_maker import get_major_version, parse_release_notes


class TestGoPRMaker(TestCase):
	"""
	Tests the go_pr_maker module.

	Note that this reiterates the docstring tests - and only covers a proper
	subset of those test cases.
	"""

	def test_get_major_version(self) -> None:
		"""Tests the get_major_version function."""
		version: str = '1.2.3'
		expected_major_version: str = '1.2'
		actual_major_version: str = get_major_version(version)
		self.assertEqual(expected_major_version, actual_major_version)
		return

	def test_get_release_notes(self) -> None:
		"""Tests the get_release_notes function."""
		go_version = '4.15.6'
		expected_release_notes = '<p> go4.15.6 The expected release notes </p>'
		release_notes_with_whitespace = f"""<p>
                go{go_version} The expected release notes
            </p>"""
		content = f"""go4.15.5 text before
        {release_notes_with_whitespace}
        text <p>after</p> 4.15.7"""
		actual_release_notes = parse_release_notes(go_version, content)
		self.assertEqual(expected_release_notes, actual_release_notes)

		go_version = '7.18.9'
		expected_release_notes = '<p id="7.18.9"> go7.18.9 The expected release notes </p>'
		release_notes_with_whitespace = f"""<p id="7.18.9">
                go{go_version} The expected release notes
            </p>"""
		content = f"""go7.18.9 text before
        {release_notes_with_whitespace}
        text <p>after</p> 7.18.9"""
		actual_release_notes = parse_release_notes(go_version, content)
		self.assertEqual(expected_release_notes, actual_release_notes)
