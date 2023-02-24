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
import os
from typing import Optional, NamedTuple

from github.MainClass import Github
from github.Repository import Repository
from github.GithubException import GithubException, UnknownObjectException
from github.GitRef import GitRef
from github.PullRequest import PullRequest
from github.InputGitTreeElement import InputGitTreeElement
from chromedriver_updater.constants import PR_GITHUB_TOKEN, GITHUB_REPO, GITHUB_REF_NAME, BRANCH_NAME, GIT_AUTHOR_NAME
import sys


class UpdateEntry(NamedTuple):
	path: str
	old_version: str
	new_version: str


def parse_update_entry(entry: str) -> UpdateEntry:
	parts = entry.split(":")
	project = parts[0]
	parts = parts[1].split(",")

	return UpdateEntry(project, parts[0], parts[1])


def path_to_project(path: str) -> str:
	if "traffic_portal" in path:
		return "Traffic Portal"
	elif "traffic-portal" in path:
		return "Traffic Portal v2"
	return ""


class PRCreator(Github):
	repo: Repository

	def __init__(self, *args, **kwargs):
		super().__init__(*args, **kwargs)
		self.repo = self.get_repo(GITHUB_REPO)

	def pr_exists(self) -> bool:
		for pr in self.search_issues("repo:%s is:pr is:open head:%s" % (self.repo.full_name, BRANCH_NAME)):
			return True
		return False

	def branch_exists(self) -> Optional[GitRef]:
		try:
			b = self.repo.get_git_ref("heads/%s" % BRANCH_NAME)
			return b
		except UnknownObjectException as e:
			return None

	def create_git_tree_from_entry(self, entry: UpdateEntry, file: str) -> InputGitTreeElement:
		file_path = os.path.join(entry.path, file)
		with open(file_path, "r") as f:
			change_file = f.read()
		blob = self.repo.create_git_blob(change_file, "utf-8")
		return InputGitTreeElement(path=file_path, mode='100644', type='blob', sha=blob.sha)

	def create_commit(self, updates: [UpdateEntry]) -> PullRequest:
		with open(os.path.join(os.path.dirname(__file__), "templates", "pr.md"), encoding="utf-8") as f:
			pr_template = f.read()
		tree_elements = []
		projects = ""
		updated = ""
		for update in updates:
			tree_elements.append(self.create_git_tree_from_entry(update, "package.json"))
			tree_elements.append(self.create_git_tree_from_entry(update, "package-lock.json"))
			projects += "* %s\n" % path_to_project(update.path)
			updated += "%s: %s -> %s\n" % (update.path, update.old_version, update.new_version)
		head_branch = self.repo.get_branch(GITHUB_REF_NAME)
		branch_ref = self.repo.create_git_ref(ref="refs/heads/%s" % BRANCH_NAME, sha=head_branch.commit.sha)
		branch = self.repo.get_branch(BRANCH_NAME)
		base_tree = self.repo.get_git_tree(sha=branch.commit.sha)
		tree = self.repo.create_git_tree(tree_elements, base_tree)
		parent = self.repo.get_git_commit(sha=branch.commit.sha)
		commit = self.repo.create_git_commit("Update chromedriver", tree, [parent])
		branch_ref.edit(sha=commit.sha)

		return self.repo.create_pull(title="Update Chromedriver Versions", maintainer_can_modify=True,
									 body=pr_template.format(UPDATES=updated, PROJECTS=projects),
									 head="%s:%s" % (GIT_AUTHOR_NAME, BRANCH_NAME), base=GITHUB_REF_NAME)


if __name__ == "__main__":
	print("Logging into github")
	gh = PRCreator(login_or_token=PR_GITHUB_TOKEN)

	if len(sys.argv) == 1:
		if gh.pr_exists():
			print("PR %s already exists" % BRANCH_NAME)
		sys.exit(0)
	elif len(sys.argv) != 2:
		print("chromedriver_updater [updates file]")
		sys.exit(1)

	updatesFile = sys.argv[1]
	if not os.path.exists(updatesFile):
		print("No updates")
		sys.exit(0)
	with open(updatesFile, "r") as f:
		updates = [line.rstrip() for line in f.readlines()]
	branch_ref = gh.branch_exists()
	if branch_ref is not None:
		print("Branch already exists (but not the pull request), recreating")
		branch_ref.delete()

	updates = [parse_update_entry(update) for update in updates if update != ""]
	gh.create_commit(updates)
