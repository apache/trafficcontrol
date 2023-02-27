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

try:
	from chromedriver_updater.constants import PR_GITHUB_TOKEN, GITHUB_REPO, GITHUB_REF_NAME, BRANCH_NAME, \
		GIT_AUTHOR_NAME, TRAFFIC_PORTAL, TRAFFIC_PORTAL_V2
except ModuleNotFoundError:
	from constants import PR_GITHUB_TOKEN, GITHUB_REPO, GITHUB_REF_NAME, BRANCH_NAME, GIT_AUTHOR_NAME, TRAFFIC_PORTAL, \
		TRAFFIC_PORTAL_V2
import sys


class UpdateEntry(NamedTuple):
	path: str
	old_version: str
	new_version: str

	def __str__(self) -> str:
		return "%s: %s -> %s\n" % (self.path, self.old_version, self.new_version)


def parse_update_entry(entry: str) -> UpdateEntry:
	if ":" not in entry:
		print("Invalid update entry '%s', expected format {path}:{old},{new}" % entry, file=sys.stderr)
		sys.exit(1)
	parts = entry.split(":")
	if len(parts) != 2 or "," not in parts[1]:
		print("Invalid update entry '%s', expected format {path}:{old},{new}" % entry, file=sys.stderr)
		sys.exit(1)
	if parts[0].endswith("/"):
		project = parts[0][:-1]
	else:
		project = parts[0]
	parts = parts[1].split(",")
	if len(parts) != 2:
		print("Invalid update entry '%s', expected format {path}:{old},{new}" % entry, file=sys.stderr)
		sys.exit(1)

	return UpdateEntry(path=project, old_version=parts[0], new_version=parts[1])


def path_to_project(path: str) -> Optional[str]:
	if "traffic_portal" in path:
		return TRAFFIC_PORTAL
	elif "traffic-portal" in path:
		return TRAFFIC_PORTAL_V2
	return None


class PRCreator(Github):
	repo: Repository

	def __init__(self, *args, **kwargs):
		super().__init__(*args, **kwargs)
		self.repo = self.get_repo(GITHUB_REPO)

	def get_pr(self) -> Optional[PullRequest]:
		for pr in self.search_issues("repo:%s is:pr is:open head:%s" % (self.repo.full_name, BRANCH_NAME)):
			return pr.as_pull_request()
		return None

	def get_branch(self) -> Optional[GitRef]:
		try:
			return self.repo.get_git_ref("heads/%s" % BRANCH_NAME)
		except UnknownObjectException:
			return None

	def create_git_tree_from_entry(self, entry: UpdateEntry, file: str) -> InputGitTreeElement:
		file_path = os.path.join(entry.path, file)
		if not os.path.exists(file_path):
			print("Could not find '%s' to commit" % file_path, file=sys.stderr)
			sys.exit(1)
		with open(file_path, "r") as f:
			change_file = f.read()
		blob = self.repo.create_git_blob(change_file, "utf-8")
		return InputGitTreeElement(path=file_path, mode='100644', type='blob', sha=blob.sha)

	def create_pull_request(self, update_entries: [UpdateEntry]) -> PullRequest:
		with open(os.path.join(os.path.dirname(__file__), "templates", "pr.md"), encoding="utf-8") as f:
			pr_template = f.read()
		tree_elements = []
		projects = ""
		updated = ""
		for update in update_entries:
			tree_elements.append(self.create_git_tree_from_entry(update, "package.json"))
			tree_elements.append(self.create_git_tree_from_entry(update, "package-lock.json"))
			updated += str(update)
			project = path_to_project(update.path)
			if project is None:
				print("Unknown project from path %s" % update.path)
			elif project not in projects:
				projects += "* %s\n" % project
		head_branch = self.repo.get_branch(GITHUB_REF_NAME)
		ref = self.repo.create_git_ref(ref="refs/heads/%s" % BRANCH_NAME, sha=head_branch.commit.sha)
		branch = self.repo.get_branch(BRANCH_NAME)
		base_tree = self.repo.get_git_tree(sha=branch.commit.sha)
		tree = self.repo.create_git_tree(tree_elements, base_tree)
		parent = self.repo.get_git_commit(sha=branch.commit.sha)
		commit = self.repo.create_git_commit("Update chromedriver", tree, [parent])
		ref.edit(sha=commit.sha)

		pr = self.repo.create_pull(title="Update Chromedriver Versions", maintainer_can_modify=True,
								   body=pr_template.format(UPDATES=updated, PROJECTS=projects),
								   head="%s:%s" % (GIT_AUTHOR_NAME, BRANCH_NAME), base=GITHUB_REF_NAME)

		try:
			test_label = self.repo.get_label("tests")
			deps_label = self.repo.get_label("dependencies")
			labels = [test_label, deps_label]
			if TRAFFIC_PORTAL in projects:
				tp_label = self.repo.get_label(TRAFFIC_PORTAL)
				labels.append(tp_label)
			if TRAFFIC_PORTAL_V2 in projects:
				tpv2_label = self.repo.get_label(TRAFFIC_PORTAL_V2)
				labels.append(tpv2_label)
			pr.add_to_labels(*labels)
		except UnknownObjectException:
			print("Could not find labels", file=sys.stderr)

		return pr


if __name__ == "__main__":
	print("Logging into github")
	gh = PRCreator(login_or_token=PR_GITHUB_TOKEN)

	if len(sys.argv) > 2:
		print("chromedriver_updater [updates file]", file=sys.stderr)
		sys.exit(1)

	pr = gh.get_pr()
	if pr is not None:
		print("PR already exists; number: %s, url: %s" % (pr.number, pr.html_url))
		sys.exit(0)

	if len(sys.argv) == 1:
		sys.exit(0)

	updatesFile = sys.argv[1]
	if not os.path.exists(updatesFile):
		print("File '%s' does not exist" % updatesFile, file=sys.stderr)
		sys.exit(1)
	with open(updatesFile, "r") as f:
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
