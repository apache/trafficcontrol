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
import sys
from datetime import date, timedelta
from typing import Optional

from github import TimelineEvent
from github.GithubException import BadCredentialsException
from github.Issue import Issue
from github.MainClass import Github
from github.NamedUser import NamedUser
from github.PaginatedList import PaginatedList
from github.Repository import Repository

from assign_triage_role.constants import GH_TIMELINE_EVENT_TYPE_CROSS_REFERENCE, ENV_GITHUB_TOKEN, ENV_GITHUB_REPOSITORY, ENV_SINCE_DAYS_AGO


class TriageRoleAssigner:
	gh: Github
	repo: Repository
	since_days_ago: int

	def __init__(self, gh: Github) -> None:
		self.gh = gh
		repo_name: str = self.get_repo_name()
		self.repo = self.get_repo(repo_name)

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

	def run(self) -> None:
		committers: dict[str, None] = self.get_committers()
		today: date = date.today()
		since_day: date = today - timedelta(days=self.since_days_ago)
		prs_by_contributor: dict[NamedUser, list[(Issue, Issue)]] = self.prs_by_contributor(since_day, today, committers)
