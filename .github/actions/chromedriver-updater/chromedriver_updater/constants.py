#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
"""
Contains and loads constants used by the updater
"""
import os
from typing import Optional


def getenv(env_name: str) -> str:
	"""
	Gets environment variable :param env_name:
	"""
	env_var: Optional[str] = os.environ.get(env_name)
	if env_var is None:
		raise NameError(f"Environment variable {env_name} is not defined")
	return env_var


GIT_AUTHOR_NAME = getenv("GIT_AUTHOR_NAME")
GITHUB_REPO = getenv("GITHUB_REPOSITORY")
GITHUB_REPOSITORY_OWNER = getenv("GITHUB_REPOSITORY_OWNER")
GITHUB_REF_NAME = getenv("GITHUB_REF_NAME")
PR_GITHUB_TOKEN = getenv("PR_GITHUB_TOKEN")
GITHUB_TOKEN = getenv("GITHUB_TOKEN")
BRANCH_NAME = "ATC-Chromedriver-Updater"
TRAFFIC_PORTAL_V2 = "Traffic Portal v2"
TRAFFIC_PORTAL = "Traffic Portal"
