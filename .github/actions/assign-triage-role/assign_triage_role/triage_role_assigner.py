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
import os
import re
import sys
from datetime import date, timedelta
from typing import Optional, Final

import yaml
from github import TimelineEvent
from github.Branch import Branch
from github.Commit import Commit
from github.ContentFile import ContentFile
from github.GithubException import BadCredentialsException, GithubException, UnknownObjectException
from github.InputGitAuthor import InputGitAuthor
from github.Issue import Issue
from github.Label import Label
from github.MainClass import Github
from github.NamedUser import NamedUser
from github.PaginatedList import PaginatedList
from github.PullRequest import PullRequest
from github.Repository import Repository
from yaml import YAMLError

from assign_triage_role.constants import GH_TIMELINE_EVENT_TYPE_CROSS_REFERENCE, ENV_GITHUB_TOKEN, ENV_GITHUB_REPOSITORY, ENV_SINCE_DAYS_AGO, ENV_MINIMUM_COMMITS, ASF_YAML_FILE, APACHE_LICENSE_YAML, ENV_GITHUB_REF_NAME, GIT_AUTHOR_EMAIL_TEMPLATE, ENV_GIT_AUTHOR_NAME, SINGLE_PR_TEMPLATE_FILE, SINGLE_CONTRIBUTOR_TEMPLATE_FILE, PR_TEMPLATE_FILE, EMPTY_CONTRIB_LIST_LIST, EMPTY_LIST_OF_CONTRIBUTORS, CONGRATS, EXPIRE, ENV_GITHUB_REPOSITORY_OWNER


class TriageRoleAssigner:
	gh: Github
	repo: Repository
	minimum_commits: int
	since_days_ago: int

	def __init__(self, gh: Github) -> None:
		self.gh = gh
		repo_name: str = self.get_repo_name()
		self.repo = self.get_repo(repo_name)

		self.minimum_commits = int(self.getenv(ENV_MINIMUM_COMMITS))
		self.since_days_ago = int(self.getenv(ENV_SINCE_DAYS_AGO))

	@staticmethod
	def getenv(env_name: str) -> str:
		return os.environ[env_name]

	def get_repo_name(self) -> str:
		repo_name: str = self.getenv(ENV_GITHUB_REPOSITORY)
		return repo_name

	def get_repo(self, repo_name: str) -> Repository:
		try:
			repo: Repository = self.gh.get_repo(repo_name)
		except BadCredentialsException:
			print(f'Credentials from {ENV_GITHUB_TOKEN} were bad.')
			sys.exit(1)
		return repo

	def get_committers(self) -> dict[str, None]:
		committers: list[str] = sorted([user.login for user in self.repo.get_collaborators() if user.permissions.push])
		committers_dict: dict[str, None] = {committer: None for committer in committers}
		return committers_dict

	def prs_by_contributor(self, since_day: date, today: date, committers: dict[str, None]) -> dict[NamedUser, list[(Issue, Issue)]]:
		# Search for PRs and Issues on the parent repo if running on a fork
		repo_name = self.repo.full_name if self.repo.parent is None else self.repo.parent.full_name

		query: str = f'repo:{repo_name} is:issue linked:pr is:closed closed:{since_day}..{today}'
		linked_issues: PaginatedList[Issue] = self.gh.search_issues(query=query)
		prs_by_contributor: dict[NamedUser, list[(Issue, Issue)]] = dict[NamedUser, list[(Issue, Issue)]]()
		for linked_issue in linked_issues:
			timeline: PaginatedList[TimelineEvent] = linked_issue.get_timeline()
			pull_request: Optional[Issue] = None
			for event in timeline:
				if event.id is not None:
					continue
				if event.event != GH_TIMELINE_EVENT_TYPE_CROSS_REFERENCE:
					continue
				pr_text = event.raw_data['source']['issue']
				if 'pull_request' not in pr_text:
					continue
				pull_request = Issue(self.repo._requester, event.raw_headers, pr_text, completed=True)
			# Skip unmerged PRs
			if 'merged_at' not in pull_request.pull_request.raw_data:
				continue
			# Do not break, in case the Issue has ever been linked to more than 1 PR in the past
			if pull_request is None:
				raise Exception(f'Unable to find a linked Pull Request for Issue {self.repo.full_name}#{linked_issue.number}')
			author: NamedUser = pull_request.user
			if author.login in committers:
				continue
			if author not in prs_by_contributor:
				prs_by_contributor[author] = list[(Issue, Issue)]()
			prs_by_contributor[author].append((pull_request, linked_issue))
		return prs_by_contributor

	def ones_who_meet_threshold(self, prs_by_contributor: dict[NamedUser, list[(Issue, Issue)]]) -> dict[str, list[(Issue, Issue)]]:
		prs_by_contributor: dict[str, list[(Issue, Issue)]] = {
			# use only the username as the dict key
			contributor.login: pull_requests
			# sort contributors by commit count
			for contributor, pull_requests in sorted(
				prs_by_contributor.items(),
				key=lambda item: len(item[1]),
				# highest commit count first
				reverse=True)
			# only include contributors who had at least self.minimum_commits Issue-closing Pull Requests merged in the past self.since_days_ago days
			if len(pull_requests) >= self.minimum_commits
		}
		return prs_by_contributor

	def set_collaborators_in_asf_yaml(self, prs_by_contributor: dict[str, list[(Issue, Issue)]], description: str):
		collaborators: list[str] = [contributor for contributor in prs_by_contributor]
		with open(ASF_YAML_FILE) as stream:
			github_key: Final[str] = 'github'
			collaborators_key: Final[str] = 'collaborators'
			try:
				asf_yaml: dict[str, dict] = yaml.safe_load(stream)
			except YAMLError as e:
				print(f'Could not load YAML file {ASF_YAML_FILE}: {e}')
				sys.exit(1)
		if github_key not in asf_yaml:
			asf_yaml[github_key] = dict[str, dict]()
		asf_yaml[github_key][collaborators_key] = collaborators

		with open(os.path.join(os.path.dirname(__file__), APACHE_LICENSE_YAML)) as stream:
			apache_license = stream.read().format(DESCRIPTION=description, ISSUE_THRESHOLD=self.minimum_commits, SINCE_DAYS_AGO=self.since_days_ago)

		with open(ASF_YAML_FILE, 'w') as stream:
			stream.write(apache_license)
			yaml.dump(asf_yaml, stream)

	def push_changes(self, target_branch_name: str, source_branch_name: str, commit_message: str) -> Commit:
		target_branch: Branch = self.repo.get_branch(target_branch_name)
		sha: str = target_branch.commit.sha
		source_branch_ref: str = f'refs/heads/{source_branch_name}'
		self.repo.create_git_ref(source_branch_ref, sha)
		print(f'Created branch {source_branch_name}')

		with open(ASF_YAML_FILE) as stream:
			asf_yaml = stream.read()

		asf_yaml_contentfile: ContentFile = self.repo.get_contents(ASF_YAML_FILE, source_branch_ref)
		kwargs = {'path': ASF_YAML_FILE,
			'message': commit_message,
			'content': asf_yaml,
			'sha': asf_yaml_contentfile.sha,
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
		print(f'Updated {ASF_YAML_FILE} on {self.repo.name} branch {source_branch_name}')
		return commit

	def get_repo_file_contents(self, branch: str) -> str:
		return self.repo.get_contents(ASF_YAML_FILE,
			f'refs/heads/{branch}').decoded_content.rstrip().decode()

	def branch_exists(self, branch: str) -> bool:
		try:
			self.get_repo_file_contents(branch)
			return True
		except GithubException as e:
			message = e.data.get('message')
			if not re.match(r'No commit found for the ref', message):
				raise e
		return False

	@staticmethod
	def list_of_contributors(prs_by_contributor: dict[str, list[(Issue, Issue)]], today: date) -> tuple[str, str, str]:
		if len(prs_by_contributor) > 0:
			joiner: str = ', ' if len(prs_by_contributor) > 2 else ' '
			list_of_contributors: list[str] = [f'@{contributor}' for contributor in prs_by_contributor.keys()]
			if len(list_of_contributors) > 1:
				list_of_contributors[-1] = f'and {list_of_contributors[-1]}'
			list_of_contributors: str = joiner.join(list_of_contributors)
			congrats: str = CONGRATS
			expire: str = EXPIRE.format(MONTH=today.strftime('%B'))
		else:
			list_of_contributors: str = EMPTY_LIST_OF_CONTRIBUTORS
			congrats: str = ''
			expire: str = ''

		return list_of_contributors, congrats, expire

	def get_pr_body(self, prs_by_contributor: dict[str, list[(Issue, Issue)]], since_day: date, today: date) -> str:
		with open(os.path.join(os.path.dirname(__file__), SINGLE_PR_TEMPLATE_FILE)) as stream:
			pr_line_template = stream.read()
		with open(os.path.join(os.path.dirname(__file__), SINGLE_CONTRIBUTOR_TEMPLATE_FILE)) as stream:
			contrib_list_template = stream.read()
		with open(os.path.join(os.path.dirname(__file__), PR_TEMPLATE_FILE)) as stream:
			pr_template = stream.read()

		contrib_list_list: str = str()
		for contributor, pr_tuples in prs_by_contributor.items():
			pr_list: str = ''
			for pr, linked_issue in pr_tuples:
				pr_line = pr_line_template.format(ISSUE_NUMBER=linked_issue.number, PR_NUMBER=pr.number)
				pr_list += pr_line + '\n'
			contrib_list: str = contrib_list_template.format(CONTRIBUTOR_USERNAME=contributor, CONTRIBUTION_COUNT=len(pr_tuples), PR_LIST=pr_list)
			contrib_list_list += contrib_list + '\n'
		if contrib_list_list == '':
			contrib_list_list = EMPTY_CONTRIB_LIST_LIST

		list_of_contributors, congrats, expire = self.list_of_contributors(prs_by_contributor, today)

		pr_body: str = pr_template.format(CONTRIB_LIST_LIST=contrib_list_list, MONTH=today.strftime('%B'), CONGRATS=congrats, LIST_OF_CONTRIBUTORS=list_of_contributors, EXPIRE=expire, ISSUE_THRESHOLD=self.minimum_commits, SINCE_DAYS_AGO=self.since_days_ago, SINCE_DAY=since_day, TODAY=today)
		# If on a fork, do not ping users or reference Issues or Pull Requests
		if self.repo.parent is not None:
			pr_body = re.sub(r'@(?!trafficcontrol)([A-Za-z0-9]+)', r'＠\1', pr_body)
			pr_body = re.sub(r'#([0-9])', r'⌗\1', pr_body)
		print('Templated PR body')
		return pr_body

	def create_pr(self, prs_by_contributor: dict[str, list[(Issue, Issue)]], commit_message: str, owner: str,
			source_branch_name: str, target_branch: str, since_day: date, today: date) -> None:
		prs: PaginatedList = self.gh.search_issues(f'repo:{self.repo.full_name} is:pr is:open head:{source_branch_name}')
		for list_item in prs:
			pr: PullRequest = self.repo.get_pull(list_item.number)
			if pr.head.ref != source_branch_name:
				continue
			print(f'Pull request for branch {source_branch_name} already exists:\n{pr.html_url}')
			return

		pr_body: str = self.get_pr_body(prs_by_contributor, since_day, today)
		pr: PullRequest = self.repo.create_pull(
			title=commit_message,
			body=pr_body,
			head=f'{owner}:{source_branch_name}',
			base=target_branch,
			maintainer_can_modify=True,
		)
		try:
			collaborators_label: Label = self.repo.get_label('collaborators')
			process_label: Label = self.repo.get_label('process')
			pr.add_to_labels(collaborators_label, process_label)
		except UnknownObjectException:
			print('Unable to find a label named "collaborators".')
		print(f'Created pull request {pr.html_url}')

	def get_repo_owner(self) -> str:
		repo_name: str = self.getenv(ENV_GITHUB_REPOSITORY_OWNER)
		return repo_name

	def run(self) -> None:
		committers: dict[str, None] = self.get_committers()
		today: date = date.today()
		since_day: date = today - timedelta(days=self.since_days_ago)
		prs_by_contributor: dict[NamedUser, list[(Issue, Issue)]] = self.prs_by_contributor(since_day, today, committers)
		prs_by_contributor: dict[str, list[(Issue, Issue)]] = self.ones_who_meet_threshold(prs_by_contributor)
		description: str = f'ATC Collaborators for {today.strftime("%B %Y")}'
		self.set_collaborators_in_asf_yaml(prs_by_contributor, description)

		source_branch_name: Final[str] = f'collaborators-{today.strftime("%Y-%m")}'
		commit_message: str = description
		target_branch_name: str = self.getenv(ENV_GITHUB_REF_NAME)
		if not self.branch_exists(source_branch_name):
			self.push_changes(target_branch_name, source_branch_name, commit_message)
		self.repo.get_git_ref(f'heads/{source_branch_name}')
		owner: str = self.get_repo_owner()
		self.create_pr(prs_by_contributor, commit_message, owner, source_branch_name, target_branch_name, since_day, today)
