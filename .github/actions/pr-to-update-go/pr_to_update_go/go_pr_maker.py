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
from typing import Optional

import requests
from github import GitRef
from github.Branch import Branch
from github.Commit import Commit
from github.GitCommit import GitCommit
from github.GitTree import GitTree
from github.GithubObject import NotSet
from github.InputGitTreeElement import InputGitTreeElement
from github.Requester import Requester

from requests import Response

from github.GithubException import BadCredentialsException, GithubException, UnknownObjectException
from github.InputGitAuthor import InputGitAuthor
from github.Label import Label
from github.MainClass import Github
from github.Milestone import Milestone
from github.PaginatedList import PaginatedList
from github.PullRequest import PullRequest
from github.Repository import Repository

from pr_to_update_go.constants import ENV_GITHUB_TOKEN, GO_VERSION_URL, ENV_GITHUB_REPOSITORY, \
	ENV_GITHUB_REPOSITORY_OWNER, GO_REPO_NAME, RELEASE_PAGE_URL, ENV_GO_VERSION_FILE, \
	ENV_GIT_AUTHOR_NAME, GIT_AUTHOR_EMAIL_TEMPLATE


class GoPRMaker:
	"""
	A class to generate pull requests for the purpose of updating the Go version in a repository.
	"""
	gh: Github
	latest_go_version: str
	repo: Repository
	author: InputGitAuthor

	def __init__(self, gh: Github) -> None:
		"""
		:param gh: Github
		:rtype: None
		"""
		self.gh = gh
		repo_name: str = self.get_repo_name()
		self.repo = self.get_repo(repo_name)

		try:
			git_author_name = self.getenv(ENV_GIT_AUTHOR_NAME)
			git_author_email = GIT_AUTHOR_EMAIL_TEMPLATE.format(git_author_name=git_author_name)
			self.author = InputGitAuthor(git_author_name, git_author_email)
		except KeyError:
			self.author = NotSet
			print('Will commit using the default author')

	def branch_exists(self, branch: str) -> bool:
		"""
		:param branch:
		:type branch:
		:return:
		:rtype: bool
		"""
		try:
			repo_go_version = self.get_repo_go_version(branch)
			if self.latest_go_version == repo_go_version:
				print(f'Branch {branch} already exists')
				return True
		except GithubException as e:
			message = e.data.get('message')
			if not re.match(r'No commit found for the ref', message):
				raise e
		return False

	def update_branch(self, branch_name: str, sha: str) -> None:
		"""
		:param branch_name:
		:type branch_name:
		:param sha:
		:type sha:
		:return:
		:rtype: None
		"""
		requester: Requester = self.repo._requester
		patch_parameters = {
			'sha': sha,
		}
		requester.requestJsonAndCheck(
			'PATCH', self.repo.url + f'/git/refs/heads/{branch_name}', input=patch_parameters
		)
		return

	def run(self, update_version_only: bool = False) -> None:
		"""
		:return:
		:rtype: None
		"""
		repo_go_version = self.get_repo_go_version()
		self.latest_go_version = self.get_latest_major_upgrade(repo_go_version)
		commit_message: str = f'Update Go version to {self.latest_go_version}'

		source_branch_name: str = f'go-{self.latest_go_version}'
		target_branch: str = 'master'
		if repo_go_version == self.latest_go_version:
			print(f'Go version is up-to-date on {target_branch}, nothing to do.')
			return

		commit: Optional[Commit] = None
		if not self.branch_exists(source_branch_name):
			commit = self.set_go_version(self.latest_go_version, commit_message,
				source_branch_name)
		if commit is None:
			source_branch_ref: GitRef = self.repo.get_git_ref('heads/go-1.15.15')
			commit = self.repo.get_commit(source_branch_ref.object.sha)
		subprocess.run(['git', 'fetch', 'origin'], check=True)
		subprocess.run(['git', 'checkout', commit.sha], check=True)
		if update_version_only:
			print(f'Branch {source_branch_name} has been created, exiting...')
			return

		update_golang_org_x_commit: Optional[GitCommit] = self.update_golang_org_x(commit)
		if isinstance(update_golang_org_x_commit, GitCommit):
			sha: str = update_golang_org_x_commit.sha
			self.update_branch(source_branch_name, sha)

		owner: str = self.get_repo_owner()
		self.create_pr(self.latest_go_version, commit_message, owner, source_branch_name,
			target_branch)

	@staticmethod
	def getenv(env_name: str) -> str:
		"""
		:param env_name: str
		:return:
		:rtype: str
		"""
		return os.environ[env_name]

	def get_repo(self, repo_name: str) -> Repository:
		"""
		:param repo_name: str
		:return:
		:rtype: Repository
		"""
		try:
			repo: Repository = self.gh.get_repo(repo_name)
		except BadCredentialsException:
			print(f'Credentials from {ENV_GITHUB_TOKEN} were bad.')
			sys.exit(1)
		return repo

	@staticmethod
	def get_major_version(from_go_version: str) -> str:
		"""
		:param from_go_version: str
		:return:
		:rtype: str
		"""
		return re.search(pattern=r'^\d+\.\d+', string=from_go_version).group(0)

	def get_latest_major_upgrade(self, from_go_version: str) -> str:
		"""
		:param from_go_version: str
		:return:
		:rtype: str
		"""
		major_version = self.get_major_version(from_go_version)
		go_version_response: Response = requests.get(GO_VERSION_URL)
		go_version_response.raise_for_status()
		go_version_content: list = json.loads(go_version_response.content)
		index = 0
		fetched_go_version: str = ''
		while True:
			if not go_version_content[index]['stable']:
				continue
			go_version_name: str = go_version_content[index]['version']
			fetched_go_version = re.search(pattern=r'[\d.]+', string=go_version_name).group(0)
			if major_version == self.get_major_version(fetched_go_version):
				break
			index += 1
		if major_version != self.get_major_version(fetched_go_version):
			raise Exception(f'No supported {major_version} Go versions exist.')
		print(f'Latest version of Go {major_version} is {fetched_go_version}')
		return fetched_go_version

	def get_repo_name(self) -> str:
		"""
		:return:
		:rtype: str
		"""
		repo_name: str = self.getenv(ENV_GITHUB_REPOSITORY)
		return repo_name

	def get_repo_owner(self) -> str:
		"""
		:return:
		:rtype: str
		"""
		repo_name: str = self.getenv(ENV_GITHUB_REPOSITORY_OWNER)
		return repo_name

	def get_go_milestone(self, go_version: str) -> str:
		"""
		:param go_version: str
		:return:
		"""
		go_repo: Repository = self.get_repo(GO_REPO_NAME)
		milestones: PaginatedList[Milestone] = go_repo.get_milestones(state='all', sort='due_on',
			direction='desc')
		milestone_title = f'Go{go_version}'
		for milestone in milestones:  # type: Milestone
			if milestone.title == milestone_title:
				print(f'Found Go milestone {milestone.title}')
				return milestone.raw_data.get('html_url')
		raise Exception(f'Could not find a milestone named {milestone_title}.')

	@staticmethod
	def get_release_notes_page() -> str:
		"""
		:return:
		:rtype: str
		"""
		release_history_response: Response = requests.get(RELEASE_PAGE_URL)
		release_history_response.raise_for_status()
		return release_history_response.content.decode()

	@staticmethod
	def get_release_notes(go_version: str, release_notes_content: str) -> str:
		"""
		:param go_version: str
		:param release_notes_content: str
		:return:
		:rtype: str
		"""
		go_version_pattern = go_version.replace('.', '\\.')
		release_notes_pattern: str = f'<p>\\s*\\n\\s*go{go_version_pattern}.*?</p>'
		release_notes_matches = re.search(release_notes_pattern, release_notes_content,
			re.MULTILINE | re.DOTALL)
		if release_notes_matches is None:
			raise Exception(f'Could not find release notes on {RELEASE_PAGE_URL}')
		release_notes = re.sub(r'[\s\t]+', ' ', release_notes_matches.group(0))
		return release_notes

	def get_pr_body(self, go_version: str, milestone_url: str) -> str:
		"""
		:param go_version: str
		:param milestone_url: str
		:return:
		:rtype: str
		"""
		with open(os.path.dirname(__file__) + '/pr_template.md') as file:
			pr_template = file.read()
		go_major_version = self.get_major_version(go_version)

		release_notes = self.get_release_notes(go_version, self.get_release_notes_page())
		pr_body: str = pr_template.format(GO_VERSION=go_version, GO_MAJOR_VERSION=go_major_version,
			RELEASE_NOTES=release_notes, MILESTONE_URL=milestone_url)
		print('Templated PR body')
		return pr_body

	def get_repo_go_version(self, branch: str = 'master') -> str:
		"""
		:param branch: str
		:return:
		:rtype: str
		"""
		return self.repo.get_contents(self.getenv(ENV_GO_VERSION_FILE),
			f'refs/heads/{branch}').decoded_content.rstrip().decode()

	def set_go_version(self, go_version: str, commit_message: str,
			source_branch_name: str) -> Commit:
		"""
		:param go_version: str
		:param commit_message: str
		:param source_branch_name: str
		:return:
		:rtype: str
		"""
		master: Branch = self.repo.get_branch('master')
		sha: str = master.commit.sha
		ref: str = f'refs/heads/{source_branch_name}'
		self.repo.create_git_ref(ref, sha)

		print(f'Created branch {source_branch_name}')
		go_version_file: str = self.getenv(ENV_GO_VERSION_FILE)
		go_file_contents = self.repo.get_contents(go_version_file, ref)
		kwargs = {'path': go_version_file,
			'message': commit_message,
			'content': (go_version + '\n'),
			'sha': go_file_contents.sha,
			'branch': source_branch_name,
		}
		try:
			git_author_name = self.getenv(ENV_GIT_AUTHOR_NAME)
			git_author_email = GIT_AUTHOR_EMAIL_TEMPLATE.format(git_author_name=git_author_name)
			author: InputGitAuthor = InputGitAuthor(name=git_author_name, email=git_author_email)
			kwargs['author'] = author
			kwargs['committer'] = author
		except KeyError:
			print('Committing using the default author')

		commit: Commit = self.repo.update_file(**kwargs).get('commit')
		print(f'Updated {go_version_file} on {self.repo.name}')
		return commit

	def update_golang_org_x(self, previous_commit: Commit) -> Optional[GitCommit]:
		"""
		:param previous_commit:
		:type previous_commit:
		:return:
		:rtype: Optional[GitCommit]
		"""
		subprocess.run(['git', 'fetch', 'origin'], check=True)
		subprocess.run(['git', 'checkout', previous_commit.sha], check=True)
		script_path: str = os.path.join(os.path.dirname(__file__), 'update_golang_org_x.sh')
		subprocess.run([script_path], check=True)
		files_to_check: list[str] = ['go.mod', 'go.sum', os.path.join('vendor', 'modules.txt')]
		tree_elements: list[InputGitTreeElement] = []
		for file in files_to_check:
			diff_process = subprocess.run(['git', 'diff', '--exit-code', '--', file])
			if diff_process.returncode == 0:
				continue
			with open(file) as stream:
				content: str = stream.read()
			tree_element: InputGitTreeElement = InputGitTreeElement(path=file, mode='100644',
				type='blob', content=content)
			tree_elements.append(tree_element)
		if len(tree_elements) == 0:
			print('No golang.org/x/ dependencies need to be updated.')
			return
		tree_hash = subprocess.check_output(
			['git', 'log', '-1', '--pretty=%T', previous_commit.sha]).decode().strip()
		base_tree: GitTree = self.repo.get_git_tree(sha=tree_hash)
		tree: GitTree = self.repo.create_git_tree(tree_elements, base_tree)
		commit_message: str = f'Update golang.org/x/ dependencies for go{self.latest_go_version}'
		previous_git_commit: GitCommit = self.repo.get_git_commit(previous_commit.sha)
		git_commit: GitCommit = self.repo.create_git_commit(message=commit_message, tree=tree,
			parents=[previous_git_commit],
			author=self.author, committer=self.author)
		print('Updated golang.org/x/ dependencies')
		return git_commit

	def create_pr(self, latest_go_version: str, commit_message: str, owner: str,
			source_branch_name: str, target_branch: str) -> None:
		"""
		:param latest_go_version: str
		:param commit_message: str
		:param owner: str
		:param source_branch_name: str
		:param target_branch: str
		:return:
		:rtype: None
		"""
		prs: PaginatedList = self.gh.search_issues(
			f'repo:{self.repo.full_name} is:pr is:open head:{source_branch_name}')
		for list_item in prs:
			pr: PullRequest = self.repo.get_pull(list_item.number)
			if pr.head.ref != source_branch_name:
				continue
			print(f'Pull request for branch {source_branch_name} already exists:\n{pr.html_url}')
			return

		milestone_url: str = self.get_go_milestone(latest_go_version)
		pr_body: str = self.get_pr_body(latest_go_version, milestone_url)
		pr: PullRequest = self.repo.create_pull(
			title=commit_message,
			body=pr_body,
			head=f'{owner}:{source_branch_name}',
			base=target_branch,
			maintainer_can_modify=True,
		)
		try:
			go_version_label: Label = self.repo.get_label('go version')
			pr.add_to_labels(go_version_label)
		except UnknownObjectException:
			print('Unable to find a label named "go version".')
		print(f'Created pull request {pr.html_url}')
