"""
Assign Triage Role
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
import sys

from github import BadCredentialsException
from yaml import YAMLError

from assign_triage_role.constants import GITHUB_TOKEN, ENV_GITHUB_TOKEN, ASF_YAML_FILE
from assign_triage_role.triage_role_assigner import TriageRoleAssigner

try:
	TriageRoleAssigner(login_or_token=GITHUB_TOKEN).run()
except BadCredentialsException as e:
	print(f"Credentials from {ENV_GITHUB_TOKEN} were bad: {e}", file=sys.stderr)
	sys.exit(1)
except YAMLError as e:
	print(f"Could not load YAML file {ASF_YAML_FILE}: {e}", file=sys.stderr)
	sys.exit(1)
