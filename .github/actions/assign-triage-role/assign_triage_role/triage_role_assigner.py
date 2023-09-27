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
import re
import sys
from datetime import date, timedelta
from http.client import NOT_FOUND
from typing import Any, Hashable, NotRequired, Optional, Final, TypedDict
from xml.dom import minidom
from xml.dom.minidom import Node
import yaml

from github.Commit import Commit
from github.GithubException import GithubException, UnknownObjectException
from github.InputGitAuthor import InputGitAuthor
from github.Issue import Issue
from github.MainClass import Github
from github.NamedUser import NamedUser
from github.Repository import Repository

from assign_triage_role.constants import (
	GH_TIMELINE_EVENT_TYPE_CROSS_REFERENCE,
	ASF_YAML_FILE,
	APACHE_LICENSE_YAML,
	GIT_AUTHOR_EMAIL_TEMPLATE,
	GITHUB_SERVER_URL,
	SINGLE_PR_TEMPLATE_FILE,
	SINGLE_CONTRIBUTOR_TEMPLATE_FILE,
	PR_TEMPLATE_FILE,
	EMPTY_CONTRIB_LIST_LIST,
	EMPTY_LIST_OF_CONTRIBUTORS,
	CONGRATS,
	EXPIRE,
	GITHUB_REPOSITORY,
	GIT_AUTHOR_NAME,
	GITHUB_REPOSITORY_OWNER,
	MINIMUM_COMMITS,
	SINCE_DAYS_AGO,
	GITHUB_REF_NAME,
	PR_GITHUB_TOKEN
)

class UpdateFileArgs(TypedDict):
	path: str
	message: str
	content: str
	sha: str
	branch: str

	author: NotRequired[InputGitAuthor]
	committer: NotRequired[InputGitAuthor]


class TriageRoleAssigner(Github):
	"""
	Triage Role Assigner
	"""
	repo: Repository
	minimum_commits: int
	since_days_ago: int
	today: date
	target_branch_name: str
	owner: str

	def __init__( self, *args: Any, **kwargs: Any) -> None:
		super().__init__(*args, **kwargs)
		repo_name = GITHUB_REPOSITORY
		self.repo = self.get_repo(repo_name)

		self.minimum_commits = int(MINIMUM_COMMITS)
		self.since_days_ago = int(SINCE_DAYS_AGO)
		self.today = date.today()

		self.target_branch_name = GITHUB_REF_NAME
		self.owner = GITHUB_REPOSITORY_OWNER

	def since_day(self) -> date:
		"""
		Gets a date :var self.since_days_ago: before :var self.day:
		"""
		return self.today - timedelta(days=self.since_days_ago)

	def get_committers(self) -> set[str]:
		"""
		Gets a set of committer usernames
		"""
		return {user.login for user in self.repo.get_collaborators() if user.permissions.push}

	def prs_by_contributor(self, committers: set[str]) -> dict[NamedUser, list[tuple[Issue, Issue]]]:
		"""
		Returns a dict of Pull Requests, associated by committer, within the last
		:var self.since_day: days of :var self.today:.
		"""
		# Search for PRs and Issues on the parent repo if running on a fork
		repo_name = self.repo.full_name if self.repo.parent is None else self.repo.parent.full_name

		query = (f"repo:{repo_name} is:issue linked:pr is:closed closed:"
		         f"{self.since_day()}..{self.today}")
		linked_issues = self.search_issues(query=query)
		prs_by_contributor: dict[NamedUser, list[tuple[Issue, Issue]]] = {}
		for linked_issue in linked_issues:
			timeline = linked_issue.get_timeline()
			pull_request: Optional[Issue] = None
			for event in timeline:
				if event.id is not None:
					continue
				if event.event != GH_TIMELINE_EVENT_TYPE_CROSS_REFERENCE:
					continue
				pr_text: dict[str, Any] = event.raw_data["source"]["issue"]
				if "pull_request" not in pr_text:
					continue
				pull_request = Issue(self.repo.__getattribute__("_requester"), event.raw_headers,
					pr_text, completed=True)
			# Do not break, in case the Issue has ever been linked to more than 1 PR in the past
			if pull_request is None:
				raise LookupError(
					f"Unable to find a linked Pull Request for Issue {self.repo.full_name}#{linked_issue.number}")
			# Skip unmerged PRs
			if "merged_at" not in pull_request.pull_request.raw_data:
				continue
			author = pull_request.user
			if author.login in committers:
				continue
			if author not in prs_by_contributor:
				prs_by_contributor[author] = []
			prs_by_contributor[author].append((pull_request, linked_issue))
		return prs_by_contributor

	def ones_who_meet_threshold(self, prs_by_contributor: dict[NamedUser,
	list[tuple[Issue, Issue]]]) -> dict[str, list[tuple[Issue, Issue]]]:
		"""
		Returns a dict of contributors who had at least self.minimum_commits Issue-closing Pull
		Requests merged in the past self.since_days_ago days
		"""
		return {
			# use only the username as the dict key
			contributor.login: pull_requests
			# sort contributors by commit count
			for contributor, pull_requests in sorted(
				prs_by_contributor.items(),
				key=lambda item: len(item[1]),
				# highest commit count first
				reverse=True)
			# only include contributors who had at least self.minimum_commits Issue-closing Pull
			# Requests merged in the past self.since_days_ago days
			if len(pull_requests) >= self.minimum_commits
		}

	def set_collaborators_in_asf_yaml(
		self,
		prs_by_contributor: dict[str, list[tuple[Issue, Issue]]],
		description: str,
		repo_url: str
	) -> None:
		"""
		Writes the list of collaborators to .asf.yaml
		"""
		collaborators: list[str] = list(prs_by_contributor)
		with open(ASF_YAML_FILE, encoding="utf-8") as stream:
			github_key: Final[str] = "github"
			collaborators_key: Final[str] = "collaborators"
			asf_yaml: dict[str, dict[Hashable, Any]] = yaml.safe_load(stream)
		if github_key not in asf_yaml:
			asf_yaml[github_key] = {}
		asf_yaml[github_key][collaborators_key] = collaborators

		with open(
			os.path.join(
				os.path.dirname(__file__),
				APACHE_LICENSE_YAML
			),
			encoding="utf-8"
		) as stream:
			apache_license = stream.read().format(
				DESCRIPTION=description,
				ISSUE_THRESHOLD=self.minimum_commits,
				SINCE_DAYS_AGO=self.since_days_ago,
				REPO_URL=repo_url
			)

		with open(ASF_YAML_FILE, "w", encoding="utf-8") as stream:
			stream.write(apache_license)
			yaml.dump(asf_yaml, stream)

	def push_changes(self, source_branch_name: str, commit_message: str) -> Commit:
		"""
		Commits the changes to the remote
		"""
		target_branch = self.repo.get_branch(self.target_branch_name)
		sha = target_branch.commit.sha
		source_branch_ref = f"refs/heads/{source_branch_name}"
		self.repo.create_git_ref(source_branch_ref, sha)
		print(f"Created branch {source_branch_name}")

		with open(ASF_YAML_FILE, encoding="utf-8") as stream:
			asf_yaml = stream.read()

		asf_yaml_contentfile = self.repo.get_contents(ASF_YAML_FILE, source_branch_ref)
		if isinstance(asf_yaml_contentfile, list):
			asf_yaml_contentfile = asf_yaml_contentfile[0]
		kwargs: UpdateFileArgs = {"path": ASF_YAML_FILE,
			"message": commit_message,
			"content": asf_yaml,
			"sha": asf_yaml_contentfile.sha,
			"branch": source_branch_name,
		}
		try:
			git_author_email = GIT_AUTHOR_EMAIL_TEMPLATE.format(git_author_name=GIT_AUTHOR_NAME)
			author = InputGitAuthor(name=GIT_AUTHOR_NAME, email=git_author_email)
			kwargs["author"] = author
			kwargs["committer"] = author
		except KeyError:
			print("Committing using the default author")

		commit = self.repo.update_file(**kwargs).get("commit")
		if not isinstance(commit, Commit):
			raise TypeError(f"expected a commit, but got: {type(commit)}")
		print(f"Updated {ASF_YAML_FILE} on {self.repo.name} branch {source_branch_name}")
		return commit

	def get_repo_file_contents(self, branch: str) -> str:
		"""
		Uses the GitHub API to get the contents of .asf.yaml
		"""
		asf_file = self.repo.get_contents(ASF_YAML_FILE, f"refs/heads/{branch}")
		if isinstance(asf_file, list):
			asf_file = asf_file[0]
		return asf_file.decoded_content.rstrip().decode()

	def branch_exists(self, branch: str) -> bool:
		"""
		Checks if a remote branch already exists
		"""
		try:
			self.get_repo_file_contents(branch)
			return True
		except GithubException as e:
			if e.status != NOT_FOUND:
				raise e
		return False

	@staticmethod
	def list_of_contributors(prs_by_contributor: dict[str, list[tuple[Issue, Issue]]],
			today: date) -> tuple[str, str, str]:
		"""
		Returns a list of contributors in a tuple, along with :var congrats: and :var expire:,
		whose values depend on the length of that list.
		"""
		if len(prs_by_contributor) > 0:
			joiner = ", " if len(prs_by_contributor) > 2 else " "
			list_of_contributors = [f"@{contributor}" for contributor in
				prs_by_contributor.keys()]
			if len(list_of_contributors) > 1:
				list_of_contributors[-1] = f"and {list_of_contributors[-1]}"
			contributors = joiner.join(list_of_contributors)
			congrats = CONGRATS
			expire = EXPIRE.format(MONTH=today.strftime("%B"))
		else:
			contributors = EMPTY_LIST_OF_CONTRIBUTORS
			congrats = ""
			expire = ""

		return contributors, congrats, expire

	@staticmethod
	def remove_comments(pr_body: str) -> str:
		"""
		Removes comments from the Pull Request body
		"""
		body = minidom.parseString(f"<body>{pr_body}</body>").firstChild
		if body is None:
			raise ValueError("failed to parse PR body")
		return "".join(node.toxml()
			for node in body.childNodes if node.nodeType != Node.COMMENT_NODE)

	def get_pr_body(self, prs_by_contributor: dict[str, list[tuple[Issue, Issue]]]) -> str:
		"""
		Renders the Pull Request template
		"""
		with open(os.path.join(os.path.dirname(__file__), SINGLE_PR_TEMPLATE_FILE),
				encoding="utf-8") as stream:
			pr_line_template = stream.read()
		with open(os.path.join(os.path.dirname(__file__),
				SINGLE_CONTRIBUTOR_TEMPLATE_FILE), encoding="utf-8") as stream:
			contrib_list_template = stream.read()
		with open(os.path.join(os.path.dirname(__file__), PR_TEMPLATE_FILE),
				encoding="utf-8") as stream:
			pr_template = stream.read()

		def contrib_list(contributor: str, pr_tuples: list[tuple[Issue, Issue]]) -> str:
			pr_list = "\n".join(
				pr_line_template.format(ISSUE_NUMBER=linked_issue.number, PR_NUMBER=pr.number
				) for pr, linked_issue in pr_tuples)
			return contrib_list_template.format(CONTRIBUTOR_USERNAME=contributor,
				CONTRIBUTION_COUNT=len(pr_tuples), PR_LIST=pr_list)

		contrib_list_list = "\n".join(
			contrib_list(contributor, pr_tuples
			) for contributor, pr_tuples in prs_by_contributor.items()
		) if len(prs_by_contributor) > 0 else EMPTY_CONTRIB_LIST_LIST

		list_of_contributors, congrats, expire = self.list_of_contributors(prs_by_contributor,
			self.today)

		repo_url = "/".join((GITHUB_SERVER_URL, GITHUB_REPOSITORY))
		pr_body = pr_template.format(
			CONTRIB_LIST_LIST=contrib_list_list,
			MONTH=self.today.strftime("%B"),
			CONGRATS=congrats,
			LIST_OF_CONTRIBUTORS=list_of_contributors,
			EXPIRE=expire,
			ISSUE_THRESHOLD=self.minimum_commits,
			SINCE_DAYS_AGO=self.since_days_ago,
			SINCE_DAY=self.since_day(),
			TODAY=self.today,
			REPO_URL=repo_url
		)
		# If on a fork, do not ping users or reference Issues or Pull Requests
		if self.repo.parent is not None:
			pr_body = re.sub(r"@(?!trafficcontrol)([A-Za-z0-9]+)", r"＠\1", pr_body)
			pr_body = re.sub(r"#([0-9])", r"⌗\1", pr_body)
		pr_body = self.remove_comments(pr_body)
		print("Templated PR body")
		return pr_body

	def create_pr(self, prs_by_contributor: dict[str, list[tuple[Issue, Issue]]], commit_message: str,
			source_branch_name: str) -> None:
		"""
		Submits a Pull Request
		"""
		pull_requests = self.search_issues(
			f"repo:{self.repo.full_name} is:pr is:open head:{source_branch_name}")
		for list_item in pull_requests:
			pull_request = self.repo.get_pull(list_item.number)
			if pull_request.head.ref != source_branch_name:
				continue
			print(f"Pull request for branch {source_branch_name} already exists:\n"
			      f"{pull_request.html_url}", file=sys.stderr)
			return

		pr_body = self.get_pr_body(prs_by_contributor)
		pull_request = self.repo.create_pull(
			title=commit_message,
			body=pr_body,
			head=f"{self.owner}:{source_branch_name}",
			base=self.target_branch_name,
			maintainer_can_modify=True,
		)
		try:
			collaborators_label = self.repo.get_label("collaborators")
			process_label = self.repo.get_label("process")
			pull_request.add_to_labels(collaborators_label, process_label)
		except UnknownObjectException:
			print("Unable to find a label named \"collaborators\".", file=sys.stderr)
		print(f"Created pull request {pull_request.html_url}")

	def run(self) -> None:
		"""
		Runs ths Triage Role Assigner
		"""
		committers = self.get_committers()
		prs_by_contributor = self.ones_who_meet_threshold(self.prs_by_contributor(committers))
		description = f"ATC Collaborators for {self.today.strftime('%B %Y')}"
		repo_url = "/".join((GITHUB_SERVER_URL, GITHUB_REPOSITORY))
		self.set_collaborators_in_asf_yaml(prs_by_contributor, description, repo_url)

		source_branch_name: Final[str] = f"collaborators-{self.today.strftime('%Y-%m')}"
		commit_message = description
		if not self.branch_exists(source_branch_name):
			self.push_changes(source_branch_name, commit_message)
		self.repo.get_git_ref(f"heads/{source_branch_name}")

		pr_github_token = PR_GITHUB_TOKEN
		self.__init__(login_or_token=pr_github_token)

		self.create_pr(prs_by_contributor, commit_message, source_branch_name)
