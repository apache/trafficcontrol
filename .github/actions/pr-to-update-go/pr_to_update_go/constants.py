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
from typing import Final

ENV_GIT_AUTHOR_NAME: Final = 'GIT_AUTHOR_NAME'
ENV_GITHUB_REPOSITORY: Final = 'GITHUB_REPOSITORY'
ENV_GITHUB_REPOSITORY_OWNER: Final = 'GITHUB_REPOSITORY_OWNER'
ENV_GITHUB_TOKEN: Final = 'GITHUB_TOKEN'
ENV_GO_VERSION_FILE: Final = 'GO_VERSION_FILE'
GIT_AUTHOR_EMAIL_TEMPLATE: Final = '{git_author_name}@users.noreply.github.com'
GO_REPO_NAME: Final = 'golang/go'
GO_VERSION_URL: Final = 'https://golang.org/dl/?mode=json'
RELEASE_PAGE_URL: Final = 'https://golang.org/doc/devel/release.html'
