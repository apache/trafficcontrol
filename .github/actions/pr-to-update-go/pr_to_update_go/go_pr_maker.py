#!/usr/bin/env python3
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
Generate pull requests that update a repository's Go version.

Classes:

    GoPRMaker

"""
import json
import os
import re
import subprocess
import sys
from pathlib import PurePath
from typing import Optional, TypedDict, Any, Union

import requests
from dotenv import set_key

from github.Commit import Commit
from github.ContentFile import ContentFile
from github.GitCommit import GitCommit
from github.GithubException import BadCredentialsException, GithubException, UnknownObjectException
from github.GitRef import GitRef
from github.InputGitAuthor import InputGitAuthor
from github.InputGitTreeElement import InputGitTreeElement
from github.MainClass import Github
from github.Repository import Repository

from pr_to_update_go.constants import (
	ENV_ENV_FILE,
	ENV_GITHUB_TOKEN,
	ENV_PR_GITHUB_TOKEN,
	GO_VERSION_URL,
	ENV_GITHUB_REPOSITORY,
	ENV_GITHUB_REPOSITORY_OWNER,
	GO_REPO_NAME,
	RELEASE_PAGE_URL,
	ENV_GO_VERSION_FILE,
	ENV_GIT_AUTHOR_NAME,
	GIT_AUTHOR_EMAIL_TEMPLATE,
	GO_VERSION_KEY,
)


class GoVersion(TypedDict):
	"""
	A single entry in the list returned by the Go website's version listing API.
	"""
	#: The type of files is unimportant, because it's not used
	files: list[Any]
	stable: bool
	version: str


def _get_pr_body(go_version: str, milestone_url: str) -> str:
	"""
	Generates the body of a Pull Request given a Go release version and a
	URL that points to information about what changes were in said release.
	"""
	with open(os.path.join(os.path.dirname(__file__), 'pr_template.md'), encoding="UTF-8") as file:
		pr_template = file.read()
	go_major_version = get_major_version(go_version)

	release_notes = _get_release_notes(go_version)
	pr_body = pr_template.format(GO_VERSION=go_version, GO_MAJOR_VERSION=go_major_version,
		RELEASE_NOTES=release_notes, MILESTONE_URL=milestone_url)
	print('Templated PR body')
	return pr_body


def get_major_version(from_go_version: str) -> str:
	"""
	Extracts the "major" version part of a full Go release version. ("major" to
	the Go project is the part of the version that most people think of as the
	major and minor versions - refer to examples).

	>>> get_major_version("1.23.45-6rc7")
	'1.23'
	>>> get_major_version("1.2.3")
	'1.2'
	>>> get_major_version("not a release version")
	''
	"""
	match = re.search(pattern=r'^\d+\.\d+', string=from_go_version)
	if match:
		return match.group(0)
	return ""


def getenv(var: str) -> str:
	"""
	Returns the value of the environment variable with the given name.

	If ``var`` is not set in the execution environment, a KeyError is raised.

	>>> os.environ["FOO"] = "BAR"
	>>> getenv("FOO")
	'BAR'
	"""
	return os.environ[var]


def parse_release_notes(version: str, content: str) -> str:
	"""
	Parses Go version release notes.

	>>> raw = '''<html lang="en-US"><head><title>test</title></head><body>
	... <h1>Big Title</h1>
	... <script>"use strict";</script>
	... <div><section><p>
	... 	go4.15.5 text before
	... </p></section></div>
	... <div><section><p>
	... 	go4.15.6 The expected release notes
	... </p></section></div>
	... <div><section><p>
	... 	go4.15.7 text after
	... </p><section></div>
	... <style>* { display: none; }</style>
	... </body></html>'''
	>>> parse_release_notes("4.15.6", raw)
	'<p> go4.15.6 The expected release notes </p>'
	>>> raw = '''<html lang="en-US"><head><title>test</title></head><body>
	... go4.15.6 is mentioned earlier
	... go4.15.6 before on the same line as the opening tag <p>
	... go4.15.6 the actual notes.
	... </p>go4.15.6 later on the same line as the closing tag
	... go4.15.6 in a later context
	... </body></html>'''
	>>> parse_release_notes("4.15.6", raw)
	'<p> go4.15.6 the actual notes. </p>'
	"""
	go_version_pattern = version.replace('.', r"\.")
	release_notes_pattern = re.compile(
		r"<p[^>]*>\s*\n\s*go" + go_version_pattern + r".*?</p>",
		re.MULTILINE | re.DOTALL
	)
	matches = release_notes_pattern.search(content)
	if not matches:
		raise Exception(f'could not find release notes for Go {version}')
	return " ".join(matches.group(0).split())


def _get_release_notes(go_version: str) -> str:
	"""
	Gets the release notes for the given Go version.
	"""
	release_history_response = requests.get(RELEASE_PAGE_URL)
	release_history_response.raise_for_status()
	return parse_release_notes(go_version, release_history_response.content.decode())


def find_latest_major_upgrade(major_version: str, versions: list[GoVersion]) -> str:
	"""
	Finds the latest version in `versions` with the given "major" version.

	Note that this expects and relies on the ordering of the passed `versions`
	being in descending order (as returned by the Go website API).
	>>> versions=[
	... {
	... 	"stable": True,
	... 	"version": "1.3.0",
	... 	"files": []
	... },
	... {
	... 	"stable": False,
	... 	"version": "1.2.5",
	...     "files": []
	... },
	... {
	... 	"stable": True,
	...     "version": "1.2.4",
	...     "files": []
	... },
	... {
	... 	"stable": True,
	... 	"version": "one.two.three",
	... 	"files": []
	... },
	... {
	... 	"stable": True,
	... 	"version": "1.2.3",
	... 	"files": []
	... }]
	>>> find_latest_major_upgrade("1.2", versions)
	'1.2.4'
	"""
	for version in versions:
		if not version["stable"]:
			continue
		match = re.search(r"[\d.]+", version["version"])
		if not match:
			continue
		fetched_go_version = match.group(0)
		if major_version == get_major_version(fetched_go_version):
			return fetched_go_version

	raise Exception(f'no supported {major_version} Go versions exist')


def _get_latest_major_upgrade(from_go_version: str) -> str:
	"""
	Gets the version of the latest Go release that is the same "major"
	version as the passed current (or "from") Go version.

	If no stable version is found that is the same "major" version as the
	given current version, an exception is raised.
	"""
	response = requests.get(GO_VERSION_URL)
	response.raise_for_status()
	versions: list[GoVersion] = json.loads(response.content)
	major_version = get_major_version(from_go_version)
	fetched_go_version = find_latest_major_upgrade(major_version, versions)
	print(f'Latest version of Go {major_version} is {fetched_go_version}')
	return fetched_go_version


class GoPRMaker:
	"""
	A class to generate pull requests for the purpose of updating the Go version
	in a repository.
	"""
	gh_api: Github
	latest_go_version: str
	repo: Repository
	author: Optional[InputGitAuthor]

	def __init__(self, gh_api: Github):
		self.gh_api = gh_api
		self.repo = self.get_repo(getenv(ENV_GITHUB_REPOSITORY))

		try:
			git_author_name = getenv(ENV_GIT_AUTHOR_NAME)
			git_author_email = GIT_AUTHOR_EMAIL_TEMPLATE.format(git_author_name=git_author_name)
			self.author = InputGitAuthor(git_author_name, git_author_email)
		except KeyError:
			self.author = None
			print('Will commit using the default author')

	def branch_exists(self, branch: str) -> bool:
		"""
		Checks the existence of a given branch in the repository.
		"""
		try:
			repo_go_version = self.get_repo_go_version(branch)
			if self.latest_go_version == repo_go_version:
				print(f'Branch {branch} already exists')
				return True
		except GithubException as e:
			message = e.data["message"]
			if not isinstance(message, str) or not re.match("No commit found for the ref", message):
				raise e
		return False

	def update_branch(self, branch_name: str, sha: str) -> None:
		"""
		Updates the branch given by ``branch_name`` on the remote origin by
		fast-forwarding it to a commit given by its hash in ``sha``.

		Note that only fast-forward updates are possible, as this doesn't
		"force" push.
		"""
		ref = self.repo.get_git_ref(f"heads/{branch_name}")
		ref.edit(sha)

	def run(self, update_version_only: bool = False) -> None:
		"""
		This is the 'main' method of the PR maker, which does everything
		necessary to create the PR that will update the repository's Go version.
		"""
		repo_go_version = self.get_repo_go_version()
		self.latest_go_version = _get_latest_major_upgrade(repo_go_version)
		commit_message = f'Update Go version to {self.latest_go_version}'

		source_branch_name = f'go-{self.latest_go_version}'
		target_branch = 'master'
		if repo_go_version == self.latest_go_version:
			print(f'Go version is up-to-date on {target_branch}, nothing to do.')
			return

		commit: Optional[GitCommit] = None
		if not self.branch_exists(source_branch_name):
			commit = self.set_go_version(self.latest_go_version, commit_message,
				source_branch_name)
		if commit is None:
			source_branch_ref: GitRef = self.repo.get_git_ref(f'heads/{source_branch_name}')
			commit = self.repo.get_git_commit(source_branch_ref.object.sha)
		subprocess.run(['git', 'fetch', 'origin'], check=True)
		subprocess.run(['git', 'checkout', commit.sha], check=True)
		if update_version_only:
			print(f'Branch {source_branch_name} has been created, exiting...')
			return

		self.update_golang_org_x(commit, source_branch_name)

		self.create_pr(
			self.latest_go_version,
			commit_message,
			getenv(ENV_GITHUB_REPOSITORY_OWNER),
			source_branch_name,
			target_branch
		)

	def get_repo(self, repo_name: str) -> Repository:
		"""
		Fetches a PyGitHub Repository object using the passed repository name.
		"""
		try:
			repo: Repository = self.gh_api.get_repo(repo_name)
		except BadCredentialsException as e:
			raise PermissionError(f"Credentials from token '{ENV_GITHUB_TOKEN}' were bad") from e
		return repo

	def get_go_milestone(self, go_version: str) -> str:
		"""
		Gets a URL for the GitHub milestone that tracks the release of the
		passed Go version.

		If the passed version is not found to have a milestone associated with
		it, a LookupError exception is raised.
		"""
		go_repo = self.get_repo(GO_REPO_NAME)
		milestones = go_repo.get_milestones(state='all', sort='due_on', direction='desc')
		milestone_title = f'Go{go_version}'
		for milestone in milestones:
			if milestone.title == milestone_title:
				print(f'Found Go milestone {milestone.title}')
				# Technically it would probably be best to use the 'html_url'
				# returned by the GH API, but accessing that through PyGithub
				# involves using poorly-documented properties of that library,
				# as well as sacrificing type-safety.
				return f"https://github.com/{GO_REPO_NAME}/milestone/{milestone.number}"
		raise LookupError(f'could not find a milestone named {milestone_title}')

	def file_contents(self, file: str, branch: str = "master") -> ContentFile:
		"""
		Gets the contents of the given file path within the repository,
		optionally on a specific branch ("master" by default).

		All trailing whitespace (e.g. extra newlines) is stripped.

		An exception is raised if ``file`` is not a path to a regular file,
		relative to the root of the repository (on the given branch).
		"""
		contents = self.repo.get_contents(file, f"refs/heads/{branch}")
		if isinstance(contents, list):
			raise IsADirectoryError(f"cannot get file contents of '{file}': is a directory")
		return contents

	def get_repo_go_version(self, branch: str = 'master') -> str:
		"""
		Gets the current Go version used at the head of the given branch (or not
		given to use "master" by default) for the repository.
		"""
		return self.file_contents(getenv(ENV_GO_VERSION_FILE),
			branch).decoded_content.decode().strip()

	def set_go_version(self, go_version: str, commit_message: str,
			source_branch_name: str) -> Optional[GitCommit]:
		"""
		Makes the commits necessary to change the Go version used by the
		repository.

		This includes updating the GO_VERSION and .env files at the repository's
		root.
		"""
		master_tip = self.repo.get_branch('master').commit
		sha = master_tip.sha
		ref = f'refs/heads/{source_branch_name}'
		self.repo.create_git_ref(ref, sha)
		print(f'Created branch {source_branch_name}')

		go_version_file = getenv(ENV_GO_VERSION_FILE)
		with open(go_version_file, 'w') as go_version_file_stream:
			go_version_file_stream.write(f'{go_version}\n')
		env_file = getenv(ENV_ENV_FILE)
		env_path = PurePath(os.path.dirname(env_file), ".env")
		set_key(dotenv_path=env_path, key_to_set=GO_VERSION_KEY, value_to_set=go_version,
			quote_mode='never')
		return self.update_files_on_tree(head=master_tip, files_to_check=[go_version_file,
			env_file], commit_message=commit_message, source_branch_name=source_branch_name)

	def update_files_on_tree(self, head: Union[Commit, GitCommit], files_to_check: list[str],
			commit_message: str, source_branch_name:
			str) -> Optional[GitCommit]:
		"""
		Commits multiple files in a single Git commit, then reverts those changes locally.
		"""
		tree_elements: list[InputGitTreeElement] = []
		for file in files_to_check:
			diff_process = subprocess.run(['git', 'diff', '--exit-code', '--', file], check=False)
			if diff_process.returncode == 0:
				continue
			with open(file, encoding="UTF-8") as stream:
				content: str = stream.read()
			tree_element: InputGitTreeElement = InputGitTreeElement(path=file, mode='100644',
				type='blob', content=content)
			tree_elements.append(tree_element)
			subprocess.run(['git', 'checkout', '--', file], check=True)
		if len(tree_elements) == 0:
			print('No files need to be updated.')
			return None
		tree_hash = subprocess.check_output(
			['git', 'log', '-1', '--pretty=%T', head.sha]).decode().strip()
		base_tree = self.repo.get_git_tree(sha=tree_hash)
		tree = self.repo.create_git_tree(tree_elements, base_tree)
		kwargs = {}
		if self.author:
			kwargs = {'author': self.author, 'committer': self.author}
		git_commit = self.repo.create_git_commit(message=commit_message, tree=tree,
			parents=[self.repo.get_git_commit(head.sha)], **kwargs)
		self.update_branch(source_branch_name, git_commit.sha)
		return git_commit

	def update_golang_org_x(self, head: GitCommit, source_branch_name: str) -> Optional[GitCommit]:
		"""
		Updates golang.org/x/ Go dependencies as necessary for the new Go
		version.
		"""
		subprocess.run(['git', 'fetch', 'origin'], check=True)
		subprocess.run(['git', 'checkout', head.sha], check=True)
		subprocess.run([os.path.join(os.path.dirname(__file__), 'update_golang_org_x.sh')],
			check=True)

		commit_message: str = f'Update golang.org/x/ dependencies for go{self.latest_go_version}'
		git_commit = self.update_files_on_tree(head=head, files_to_check=['go.mod', 'go.sum',
			os.path.join('vendor', 'modules.txt')], commit_message=commit_message,
			source_branch_name=source_branch_name)
		print('Updated golang.org/x/ dependencies')
		return git_commit

	def create_pr(self, latest_go_version: str, commit_message: str, owner: str,
			source_branch_name: str, target_branch: str) -> None:
		"""
		Creates the pull request to update the Go version.
		"""
		prs = self.gh_api.search_issues(
			f'repo:{self.repo.full_name} is:pr is:open head:{source_branch_name}')
		for list_item in prs:
			pull_request = self.repo.get_pull(list_item.number)
			if pull_request.head.ref != source_branch_name:
				continue
			print(
				f'Pull request for branch {source_branch_name} already exists:\n{pull_request.html_url}')
			return

		milestone_url = self.get_go_milestone(latest_go_version)
		pr_body = _get_pr_body(latest_go_version, milestone_url)

		try:
			pr_github_token = getenv(ENV_PR_GITHUB_TOKEN)
			self.gh_api = Github(login_or_token=pr_github_token)
			self.repo = self.get_repo(getenv(ENV_GITHUB_REPOSITORY))
		except KeyError:
			print(f'Token in {ENV_PR_GITHUB_TOKEN} is invalid, creating the PR using the '
			      f'{ENV_GITHUB_TOKEN} token')
			pass

		pull_request = self.repo.create_pull(
			title=commit_message,
			body=pr_body,
			head=f'{owner}:{source_branch_name}',
			base=target_branch,
			maintainer_can_modify=True,
		)
		try:
			go_version_label = self.repo.get_label('go version')
			pull_request.add_to_labels(go_version_label)
		except UnknownObjectException:
			print('Unable to find a label named "go version"', file=sys.stderr)
		print(f'Created pull request {pull_request.html_url}')


if __name__ == "__main__":
	import doctest

	doctest.testmod()
