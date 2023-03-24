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
Utility for automatically updating chromedriver. See README.rst for more details.
"""
import os
import sys
from typing import Optional, NamedTuple

from github.MainClass import Github
from github.Repository import Repository
from github.GithubException import UnknownObjectException
from github.GitRef import GitRef
from github.PullRequest import PullRequest
from github.InputGitTreeElement import InputGitTreeElement

try:
	from chromedriver_updater.constants import PR_GITHUB_TOKEN, GITHUB_REPO, GITHUB_REF_NAME, \
		BRANCH_NAME, GIT_AUTHOR_NAME, TRAFFIC_PORTAL, TRAFFIC_PORTAL_V2, GITHUB_REPOSITORY_OWNER
except ModuleNotFoundError:
	from constants import PR_GITHUB_TOKEN, GITHUB_REPO, GITHUB_REF_NAME, BRANCH_NAME, \
		GIT_AUTHOR_NAME, TRAFFIC_PORTAL, TRAFFIC_PORTAL_V2, GITHUB_REPOSITORY_OWNER


class UpdateEntry(NamedTuple):
	"""
	Named tuple surrounding each line in the updates file
	"""
	path: str
	old_version: str
	new_version: str

	def __str__(self) -> str:
		return f"{self.path}: {self.old_version} -> {self.new_version}\n"


def parse_update_entry(entry: str) -> UpdateEntry:
	"""
	Parses string to UpdateEntry
	:param entry: String to parse
	:return: Parsed string as UpdateEntry
	"""
	if ":" not in entry:
		print(f"Invalid update entry '{entry}', expected format {{path}}:{{old}},{{new}}",
			  file=sys.stderr)
		sys.exit(1)
	parts = entry.split(":")
	if len(parts) != 2 or "," not in parts[1]:
		print(f"Invalid update entry '{entry}', expected format {{path}}:{{old}},{{new}}",
			  file=sys.stderr)
		sys.exit(1)
	if parts[0].endswith("/"):
		project = parts[0][:-1]
	else:
		project = parts[0]
	parts = parts[1].split(",")
	if len(parts) != 2:
		print(f"Invalid update entry '{entry}', expected format {{path}}:{{old}},{{new}}",
			  file=sys.stderr)
		sys.exit(1)

	return UpdateEntry(path=project, old_version=parts[0], new_version=parts[1])


def path_to_project(path: str) -> Optional[str]:
	"""
	Converts paths to projects
	:param path: Path to check
	:return: The project string, or None
	"""
	if "traffic_portal" in path:
		return TRAFFIC_PORTAL
	if "traffic-portal" in path:
		return TRAFFIC_PORTAL_V2
	return None


class PRCreator(Github):
	"""
	Creates the PR to update chromedriver
	"""
	repo: Repository

	def __init__(self, *args, **kwargs):
		super().__init__(*args, **kwargs)
		self.repo = self.get_repo(GITHUB_REPO)

	def get_pr(self) -> Optional[PullRequest]:
		"""
		Retrieve the PR opened by this script if available
		:return:
		"""
		for issue in self.search_issues("repo:%s is:pr is:open head:%s" %
									 (self.repo.full_name, BRANCH_NAME)):
			return issue.as_pull_request()
		return None

	def get_branch(self) -> Optional[GitRef]:
		"""
		Retrieve the git reference to this scripts branch if available
		:return:
		"""
		try:
			return self.repo.get_git_ref(f"heads/{BRANCH_NAME}")
		except UnknownObjectException:
			return None

	def create_git_tree_from_entry(self, entry: UpdateEntry, file: str) -> InputGitTreeElement:
		"""
		Creates a git element for the purposes of committing
		:param entry:
		:param file:
		:return:
		"""
		file_path = os.path.join(entry.path, file)
		if not os.path.exists(file_path):
			print(f"Could not find '{file_path}' to commit", file=sys.stderr)
			sys.exit(1)
		with open(file_path, "r", encoding="utf-8") as changed_file:
			change_file = changed_file.read()
		blob = self.repo.create_git_blob(change_file, "utf-8")
		return InputGitTreeElement(path=file_path, mode='100644', type='blob', sha=blob.sha)

	def create_pull_request(self, update_entries: [UpdateEntry]) -> PullRequest:
		"""
		Create a pull request based on update entries
		:param update_entries:
		:return:
		"""
		with open(os.path.join(os.path.dirname(__file__), "templates", "pr.md"),
				  encoding="utf-8") as template_file:
			pr_template = template_file.read()
		tree_elements = []
		projects = ""
		updated = ""
		for update in update_entries:
			tree_elements.append(self.create_git_tree_from_entry(update, "package.json"))
			tree_elements.append(self.create_git_tree_from_entry(update, "package-lock.json"))
			updated += str(update)
			project = path_to_project(update.path)
			if project is None:
				print(f"Unknown project from path {update.path}")
			elif project not in projects:
				projects += f"* {project}\n"
		ref = self.repo.create_git_ref(ref=f"refs/heads/{BRANCH_NAME}",
									   sha=self.repo.get_branch(GITHUB_REF_NAME).commit.sha)
		branch = self.repo.get_branch(BRANCH_NAME)
		base_tree = self.repo.get_git_tree(sha=branch.commit.sha)
		commit = self.repo.create_git_commit("Update chromedriver",
											 self.repo.create_git_tree(tree_elements, base_tree),
											 [self.repo.get_git_commit(sha=branch.commit.sha)])
		ref.edit(sha=commit.sha)

		pull = self.repo.create_pull(title="Update Chromedriver Versions", maintainer_can_modify=True,
									 body=pr_template.format(UPDATES=updated, PROJECTS=projects),
									 head=f"{GITHUB_REPOSITORY_OWNER}:{BRANCH_NAME}", base=GITHUB_REF_NAME)

		try:
			labels = [self.repo.get_label("tests"), self.repo.get_label("dependencies")]
			if TRAFFIC_PORTAL in projects:
				labels.append(self.repo.get_label(TRAFFIC_PORTAL))
			if TRAFFIC_PORTAL_V2 in projects:
				labels.append(self.repo.get_label(TRAFFIC_PORTAL_V2))
			pull.add_to_labels(*labels)
		except UnknownObjectException:
			print("Could not find labels", file=sys.stderr)

		return pull


if __name__ == "__main__":
	print("Logging into github")
	gh = PRCreator(login_or_token=PR_GITHUB_TOKEN)

	if len(sys.argv) > 2:
		print("chromedriver_updater [updates file]", file=sys.stderr)
		sys.exit(1)

	pull_request = gh.get_pr()
	if pull_request is not None:
		print(f"PR already exists; number: {pull_request.number}, url: {pull_request.html_url}")
		sys.exit(0)

	if len(sys.argv) == 1:
		sys.exit(0)

	if not os.path.exists(sys.argv[1]):
		print(f"File '{sys.argv[1]}' does not exist", file=sys.stderr)
		sys.exit(1)
	with open(sys.argv[1], "r", encoding="utf-8") as f:
		updates = [line.rstrip() for line in f.readlines()]
	if len(updates) == 0:
		print("Nothing to update")
		sys.exit(0)
	branch_ref = gh.get_branch()
	if branch_ref is not None:
		print("Branch already exists (but not a pull request), recreating")
		branch_ref.delete()

	updates = [parse_update_entry(update) for update in updates if update != ""]
	gh.create_pull_request(updates)
