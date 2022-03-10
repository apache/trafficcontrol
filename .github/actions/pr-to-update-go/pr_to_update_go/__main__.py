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
A utility for automatically updating Go. For full usage information, refer to
the action's README.rst file.

Options
	--update-version-only   Exit after updating GO_VERSION file.
"""
import os
import sys
from argparse import ArgumentParser, Namespace

from github.MainClass import Github

from pr_to_update_go.go_pr_maker import GoPRMaker
from pr_to_update_go.constants import ENV_GITHUB_TOKEN


def main() -> None:
	"""
	The entrypoint for running the PR-maker.
	"""
	parser = ArgumentParser()
	parser.add_argument(
		'--update-version-only',
		action="store_true",
		help='Exit after updating the GO_VERSION file'
	)
	args: Namespace = parser.parse_args()

	try:
		github_token = os.environ[ENV_GITHUB_TOKEN]
	except KeyError:
		print(f'Environment variable {ENV_GITHUB_TOKEN} must be defined.')
		sys.exit(1)
	gh_api = Github(login_or_token=github_token)
	GoPRMaker(gh_api).run(args.update_version_only)

main()
