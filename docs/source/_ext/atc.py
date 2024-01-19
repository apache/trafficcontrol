"""
This file is a Sphinx extension used by the ATC project to provide some custom
behavior to the Sphinx runtime.
"""
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
import os
import re
from typing import Any, Final, Optional, Dict, List, Tuple

from docutils import nodes, utils
from sphinx.application import Sphinx
from sphinx.util.docutils import SphinxDirective
from sphinx.locale import _
from sphinx.domains import changeset


REPO_URI: Final = "https://github.com/apache/trafficcontrol/"

class impl(nodes.Admonition, nodes.Element):
	"""
	This class is just used to provide translatability for the documentation.
	Which is currently unimplemented, but the translators will choke on
	directives that don't have visit and depart methods defined, so if that's
	ever added in the future the sort of "incomplete" definitions of directives
	here would be a problem.
	"""

def visit_impl_node(self: nodes.Admonition, node: nodes.Node):
	"""
	Implementation of Node entry for custom admonitions for translation
	purposes.
	"""
	self.visit_admonition(node)

def depart_impl_node(self: nodes.Admonition, node: nodes.Node):
	"""
	Implementation of Node entry for custom admonitions for translation
	purposes.
	"""
	self.depart_admonition(node)

# -- Implementation detail directive -----------------------------------------
class ImplementationDetail(SphinxDirective):
	"""
	An admonition that describes something that is possibly interesting or
	informative context for our particular implementation of something.

	Example:

		.. impl-detail::
			This is because the underlying number is 32 bits wide.
	"""
	has_content = True
	required_arguments = 0
	optional_arguments = 1
	final_argument_whitespace = True

	label_text = "Implementation Detail"

	def run(self) -> List[nodes.Node]:
		"""
		Converts the content of the directive into a proper Node.
		"""
		impl_node = impl("\n".join(self.content))
		impl_node += nodes.title(_(self.label_text), _(self.label_text))
		self.state.nested_parse(self.content, self.content_offset, impl_node)
		if self.arguments:
			n, m = self.state.inline_text(self.arguments[0], self.lineno)
			impl_node.append(nodes.paragraph("", "", *(n + m)))
		return [impl_node]

# -- Version Removed directive ----------------------------------------
changeset.versionlabels["versionremoved"] = "Removed in version %s"
changeset.versionlabel_classes["versionremoved"] = "removed"
class VersionRemoved(changeset.VersionChange):
	"""
	A directive that annotates the version at which something has been removed.

	Example:

		******************************************
		``deliveryservices_required_capabilities``
		******************************************
		.. versionremoved:: 5.0
			This is now a property of Delivery Services themselves instead.
	"""

# -- Go Version role --------------------------------------------------
# Returns the value of the Go version stored in GO_VERSION to minor version
# precision.

def atc_go_version(
	unused_typ: None,
	unused_rawtext: None,
	unused_text: None,
	lineno: int,
	unused_inliner: None,
	unused_options: None=None,
	unused_content: None=None
) -> Tuple[List[nodes.Node], List[nodes.Node]]:
	"""
	A role that inserts the Go version used/required by this version of ATC.

	Example:

		:atc-go-version:_
	"""
	go_version_file = os.path.join(os.path.dirname(__file__), "../../../GO_VERSION")
	with open(file=go_version_file, encoding="utf-8") as go_version_file:
		go_version = go_version_file.read()

	matches = re.match(pattern=r"\d+\.\d+", string=go_version)
	if matches is None:
		raise ValueError(f"Go version found that could not be parsed: '{go_version}' (from line {lineno})")
	major_minor_version = matches.group()
	strong_node = nodes.strong(major_minor_version, major_minor_version)
	return [strong_node], []

# -- Issue role --------------------------------------------------------------

ISSUE_URI: Final = REPO_URI + "issues/%s"

def issue_role(
	unused_typ: None,
	unused_rawtext: None,
	text: str,
	unused_lineno: None,
	unused_inliner: None,
	options: Optional[Dict[str, Any]]=None,
	content: Optional[List[Any]]=None
) -> Tuple[List[nodes.Node], List[nodes.Node]]:
	"""
	A role that can be used to link to an Issue by number.

	Example:

		This is tracked by :issue:1234
	"""
	if options is None:
		options = {}
	if content is None:
		content = []

	issue = utils.unescape(text)
	text = f"Issue #{issue}"
	refnode = nodes.reference(text, text, refuri=ISSUE_URI % issue)
	return [refnode], []

# -- Pull Request Role -------------------------------------------------------
PR_URI: Final = REPO_URI + "pull/%s"

def pr_role(
	unused_typ: None,
	unused_rawtext: None,
	text: str,
	unused_lineno: None,
	unused_inliner: None,
	options: Optional[Dict[str, Any]]=None,
	content: Optional[List[Any]]=None
) -> Tuple[List[nodes.Node], List[nodes.Node]]:
	"""
	A role that can be used to link to a Pull Request by number.

	Example:

		This is tracked by :pr:1234
	"""

	if options is None:
		options = {}
	if content is None:
		content = []

	pr = utils.unescape(text)
	text = f"Pull Request {pr}"
	refnode = nodes.reference(text, text, refuri=PR_URI % pr)
	return [refnode], []

# -- ATC file role -----------------------------------------------------------
FILE_URI = REPO_URI + "tree/master/%s"
def atc_file_role(
	unused_typ: None,
	unused_rawtext: None,
	text: str,
	unused_lineno: None,
	unused_inliner: None,
	options: Optional[Dict[str, Any]]=None,
	content: Optional[List[Any]]=None
) -> Tuple[List[nodes.Node], List[nodes.Node]]:
	"""
	A role that can be used to link to a file within the ATC repository.

	Example:

		The version of the ATC project is maintained in :atc-file:VERSION
	"""
	if options is None:
		options = {}
	if content is None:
		content = []

	text = utils.unescape(text)
	litnode = nodes.literal(text, text)
	refnode = nodes.reference(text, "", litnode, refuri=FILE_URI % text)
	return [refnode], []

# -- GoDoc role (absolute) ---------------------------------------------------
GODOC_URI: Final = "https://pkg.go.dev/"
def godoc_role(
	unused_typ: None,
	unused_rawtext: None,
	text: str,
	unused_lineno: None,
	unused_inliner: None,
	options: Optional[Dict[str, Any]]=None,
	content: Optional[List[Any]]=None
) -> Tuple[List[nodes.Node], List[nodes.Node]]:
	"""
	Links to the GoDoc documentation for some package.

	Example:

		We use the :godoc:`github.com/jmoiron/sqlx` package for SQL stuff.
	"""
	if options is None:
		options = {}
	if content is None:
		content = []

	text = utils.unescape(text)
	litnode = nodes.literal(text, text)
	last_period_index = text.rfind(".")
	if last_period_index != -1 and text.rfind("/") < last_period_index:
		text = text[:last_period_index] + "#" + text[last_period_index + 1:]
	refnode = nodes.reference(text, "", litnode, refuri=GODOC_URI + text)
	return [refnode], []

# -- GoDoc role (atc-relative) ----------------------------------------------
ATC_GODOC_PREFIX: Final = "github.com/apache/trafficcontrol/"
ATC_GODOC_URI: Final = GODOC_URI + ATC_GODOC_PREFIX
def atc_godoc_role(
	unused_typ: None,
	unused_rawtext: None,
	text: str,
	unused_lineno: None,
	unused_inliner: None,
	options: Optional[Dict[str, Any]]=None,
	content: Optional[List[Any]]=None
) -> Tuple[List[nodes.Node], List[nodes.Node]]:
	"""
	A role that can be used to link to the GoDoc documentation for a symbol from
	the ATC project. This is equivalent to using the ``godoc`` role, but allows
	one to omit "github.com/apache/trafficcontrol", which can be tedious and
	consume a lot of space.

	Example:

		Servers are modeled by :atc-godoc:`lib/go-tc#ServerV5`
	"""
	if options is None:
		options = {}
	if content is None:
		content = []

	text = utils.unescape(text)
	literaltext = ATC_GODOC_PREFIX + text
	litnode = nodes.literal(literaltext, literaltext)
	refnode = nodes.reference(text, "", litnode, refuri=ATC_GODOC_URI + text.replace(".", "#", 1))
	return [refnode], []

# -- GoDoc role (to-relative) -----------------------------------------------
TO_GODOC_PREFIX: Final = ATC_GODOC_PREFIX + "traffic_ops/traffic_ops_golang/"
TO_GODOC_URI: Final = GODOC_URI + TO_GODOC_PREFIX
def to_godoc_role(
	unused_typ: None,
	unused_rawtext: None,
	text: str,
	unused_lineno: None,
	unused_inliner: None,
	options: Optional[Dict[Any, Any]]=None,
	content: Optional[List[Any]]=None
) -> Tuple[List[nodes.Node], List[nodes.Node]]:
	"""
	A role that can be used to link to the GoDoc documentation for a symbol from
	within Traffic Ops's package. This is equivalent to using the ``godoc``
	role, but allows one to omit
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/", which
	can be tedious and consume a lot of space.

	Example:

		API endpoints make heavy use of :to-godoc:`api#Info` objects.
	"""
	if options is None:
		options = {}
	if content is None:
		content = []

	text = utils.unescape(text)
	literaltext = TO_GODOC_PREFIX + text
	litnode = nodes.literal(literaltext, literaltext)
	refnode = nodes.reference(text, "", litnode, refuri=TO_GODOC_URI + text.replace(".", "#", 1))
	return [refnode], []


def setup(app: Sphinx) -> Dict[str, Any]:
	"""
	This is the function that adds all of the above declared directives and
	roles to the Sphinx environment. In order to make use of a new directive or
	role, you **must** add it to the app within this setup function.
	"""
	app.add_node(
		impl,
		html=(visit_impl_node, depart_impl_node),
		latex=(visit_impl_node, depart_impl_node),
		text=(visit_impl_node, depart_impl_node)
	)
	app.add_directive("impl-detail", ImplementationDetail)
	app.add_directive("versionremoved", VersionRemoved)
	app.add_role("atc-go-version", atc_go_version)
	app.add_role("issue", issue_role)
	app.add_role("pr", pr_role)
	app.add_role("pull-request", pr_role)
	app.add_role("atc-file", atc_file_role)
	app.add_role("godoc", godoc_role)
	app.add_role("atc-godoc", atc_godoc_role)
	app.add_role("to-godoc", to_godoc_role)

	return {
		"version": "1.0",
		"parallel_read_safe": True,
		"parallel_write_safe": True,
	}
