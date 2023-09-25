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
import os
from typing import Final

ENV_GIT_AUTHOR_NAME: Final[str] = "GIT_AUTHOR_NAME"
ENV_GITHUB_REPOSITORY: Final[str] = "GITHUB_REPOSITORY"
ENV_GITHUB_REPOSITORY_OWNER: Final[str] = "GITHUB_REPOSITORY_OWNER"
ENV_GITHUB_SERVER_URL: Final[str] = "GITHUB_SERVER_URL"
ENV_GITHUB_TOKEN: Final[str] = "GITHUB_TOKEN"
ENV_PR_GITHUB_TOKEN: Final[str] = "PR_GITHUB_TOKEN"
ENV_GITHUB_REF_NAME: Final[str] = "GITHUB_REF_NAME"
ENV_MINIMUM_COMMITS: Final[str] = "MINIMUM_COMMITS"
ENV_SINCE_DAYS_AGO: Final[str] = "SINCE_DAYS_AGO"

GH_TIMELINE_EVENT_TYPE_CROSS_REFERENCE: Final[str] = "cross-referenced"
GIT_AUTHOR_EMAIL_TEMPLATE: Final[str] = "{git_author_name}@users.noreply.github.com"
ASF_YAML_FILE: Final[str] = ".asf.yaml"
SINGLE_PR_TEMPLATE_FILE: Final[str] = "templates/single_pr.md"
SINGLE_CONTRIBUTOR_TEMPLATE_FILE: Final[str] = "templates/single_contributor.md"
EMPTY_CONTRIB_LIST_LIST: Final[str] = "(None)"
CONGRATS: Final[str] = "Congrats! "
EXPIRE: Final[str] = "These privileges will expire at the end of {MONTH}."
EMPTY_LIST_OF_CONTRIBUTORS: Final[str] = "no one"
PR_TEMPLATE_FILE: Final[str] = "templates/pr_template.md"
APACHE_LICENSE_YAML: Final[str] = "templates/apache_license.yml"


def getenv(env_name: str) -> str:
	"""
	Gets environment variable :param env_name:
	"""
	env_var = os.environ.get(env_name)
	if env_var is None:
		raise NameError(f"Environment variable {env_name} is not defined")
	return env_var


GIT_AUTHOR_NAME: Final[str] = getenv(ENV_GIT_AUTHOR_NAME)
GITHUB_REPOSITORY: Final[str] = getenv(ENV_GITHUB_REPOSITORY)
GITHUB_REPOSITORY_OWNER: Final[str] = getenv(ENV_GITHUB_REPOSITORY_OWNER)
GITHUB_SERVER_URL: Final[str] = getenv(ENV_GITHUB_SERVER_URL)
GITHUB_TOKEN: Final[str] = getenv(ENV_GITHUB_TOKEN)
PR_GITHUB_TOKEN: Final[str] = getenv(ENV_PR_GITHUB_TOKEN)
GITHUB_REF_NAME: Final[str] = getenv(ENV_GITHUB_REF_NAME)
MINIMUM_COMMITS: Final[str] = getenv(ENV_MINIMUM_COMMITS)
SINCE_DAYS_AGO: Final[str] = getenv(ENV_SINCE_DAYS_AGO)
