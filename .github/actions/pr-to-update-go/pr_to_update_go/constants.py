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
The constants module holds some constants used by the PR maker.

Environment variable names, and the meanings of the values of those variables:

	ENV_GIT_AUTHOR_NAME         - The username of the author for commits
	ENV_GITHUB_REPOSITORY       - The "name" of the repository set by GHA (e.g. octocat/Hello-World)
	ENV_GITHUB_REPOSITORY_OWNER - The repository owner's name set by GHA (e.g. octocat)
	ENV_GITHUB_TOKEN            - The token used to access the GitHub API - set by GHA
	PR_GITHUB_TOKEN             - The token used to access the GitHub API for creating the Pull
	                              Request- set by GHA
	ENV_GO_VERSION_FILE         - The repository-relative path to the file containing the Go version
	ENV_ENV_FILE                - The repository-relative path to an environment file containing
	                              a line setting the variable GO_VERSION to the Go version
	                              (e.g. GO_VERSION=3.2.1)

Miscellaneous:

	GIT_AUTHOR_EMAIL_TEMPLATE - Template used to construct the Git Author's email address.
	GO_REPO_NAME              - The name of the official Go repository.
	GO_VERSION_URL            - A URL from which information about Go releases is fetched.
	RELEASE_PAGE_URL          - The URL of a webpage containing changelog notes about Go releases.
"""
from typing import Final

ENV_GIT_AUTHOR_NAME: Final = 'GIT_AUTHOR_NAME'
"""
The name of the environment variable set to username of the author for commits.
"""

ENV_GITHUB_REPOSITORY: Final = 'GITHUB_REPOSITORY'
"""
The name of the environment variable set to "name" of the repository set by GHA
(e.g. octocat/Hello-World).
"""

ENV_GITHUB_REPOSITORY_OWNER: Final = 'GITHUB_REPOSITORY_OWNER'
"""
The name of the environment variable set to repository owner's name set by GHA
(e.g. octocat).
"""

ENV_GITHUB_TOKEN: Final = 'GITHUB_TOKEN'
"""
The name of the environment variable set to token used to access the GitHub
API - set by GHA.
"""

ENV_PR_GITHUB_TOKEN: Final = 'PR_GITHUB_TOKEN'
"""
The name of the environment variable set to token used to access the GitHub
API, but only for creating the Pull Request, so that Actions will run on the
generated Pull Request - set by GHA.
"""

ENV_GO_VERSION_FILE: Final = 'GO_VERSION_FILE'
"""
The name of the environment variable set to repository-relative path to the file
containing the Go version.
"""

ENV_ENV_FILE: Final = 'ENV_FILE'
"""
The name of the environment variable set to repository-relative path to an
environment file containing a line setting the variable GO_VERSION to the Go
version (e.g. GO_VERSION=3.2.1).
"""

GO_VERSION_KEY: Final = 'GO_VERSION'
"""
The key in the env file whose value corresponds to the Go version to be used by any project
using the env file
"""

GIT_AUTHOR_EMAIL_TEMPLATE: Final = '{git_author_name}@users.noreply.github.com'
"""Template used to construct the Git Author's email address."""

GO_REPO_NAME: Final = 'golang/go'
"""The name of the official Go repository."""

GO_VERSION_URL: Final = 'https://golang.org/dl/?mode=json'
"""A URL from which information about Go releases is fetched."""

RELEASE_PAGE_URL: Final = 'https://go.dev/doc/devel/release'
"""The URL of a webpage containing changelog notes about Go releases."""


__all__ = [
	"ENV_GIT_AUTHOR_NAME",
	"ENV_GITHUB_REPOSITORY",
	"ENV_GITHUB_REPOSITORY_OWNER",
	"ENV_GITHUB_TOKEN",
	"ENV_PR_GITHUB_TOKEN"
	"ENV_GO_VERSION_FILE",
	"ENV_ENV_FILE",
	"GIT_AUTHOR_EMAIL_TEMPLATE",
	"GO_REPO_NAME",
	"GO_VERSION_URL",
	"RELEASE_PAGE_URL",
	"GO_VERSION_KEY",
]
