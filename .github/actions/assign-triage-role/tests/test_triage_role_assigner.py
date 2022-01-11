"""
Test the Triage Role Assigner
"""
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
from datetime import date
from unittest import TestCase

from assign_triage_role.triage_role_assigner import TriageRoleAssigner


class TestTriageRoleAssigner(TestCase):
	"""
	Tests for the Triage Role Assigner
	"""

	def test_list_of_contributors(self) -> None:
		"""
		Test TriageRoleAssigner.list_of_contributors()
		"""
		today = date.fromisoformat("2021-12-10")

		empty_prs_by_contributor = {}
		expected_list_of_contributors = "no one"
		expected_congrats = ""
		expected_expire = ""
		list_of_contributors, congrats, expire = TriageRoleAssigner.list_of_contributors(
			empty_prs_by_contributor, today)
		self.assertEqual(expected_list_of_contributors, list_of_contributors)
		self.assertEqual(expected_congrats, congrats)
		self.assertEqual(expected_expire, expire)

		prs_by_contributor = {"Namey Name": []}
		expected_list_of_contributors = "@Namey Name"
		expected_congrats = "Congrats! "
		expected_expire = "These privileges will expire at the end of December."
		list_of_contributors, congrats, expire = TriageRoleAssigner.list_of_contributors(
			prs_by_contributor, today)
		self.assertEqual(expected_list_of_contributors, list_of_contributors)
		self.assertEqual(expected_congrats, congrats)
		self.assertEqual(expected_expire, expire)

		prs_by_contributor = {"Namey Name": [], "A Contributor": []}
		expected_list_of_contributors = "@Namey Name and @A Contributor"
		expected_congrats = "Congrats! "
		expected_expire = "These privileges will expire at the end of December."
		list_of_contributors, congrats, expire = TriageRoleAssigner.list_of_contributors(
			prs_by_contributor, today)
		self.assertEqual(expected_list_of_contributors, list_of_contributors)
		self.assertEqual(expected_congrats, congrats)
		self.assertEqual(expected_expire, expire)

		prs_by_contributor = {"Namey Name": [], "A Contributor": [], "Someone Else": []}
		expected_list_of_contributors = "@Namey Name, @A Contributor, and @Someone Else"
		expected_congrats = "Congrats! "
		expected_expire = "These privileges will expire at the end of December."
		list_of_contributors, congrats, expire = TriageRoleAssigner.list_of_contributors(
			prs_by_contributor, today)
		self.assertEqual(expected_list_of_contributors, list_of_contributors)
		self.assertEqual(expected_congrats, congrats)
		self.assertEqual(expected_expire, expire)
