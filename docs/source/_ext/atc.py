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

REPO_URI = "https://github.com/apache/trafficcontrol/"

# -- Implementation detail directive -----------------------------------------
from docutils import nodes
from sphinx.util.docutils import SphinxDirective
from sphinx.locale import translators, _

class impl(nodes.Admonition, nodes.Element):
	pass

def visit_impl_node(self, node):
	self.visit_admonition(node)

def depart_impl_node(self, node):
	self.depart_admonition(node)

class ImplementationDetail(SphinxDirective):

	has_content = True
	required_arguments = 0
	optional_arguments = 1
	final_argument_whitespace = True

	label_text = 'Implementation Detail'

	def run(self):
		impl_node = impl('\n'.join(self.content))
		impl_node += nodes.title(_(self.label_text), _(self.label_text))
		self.state.nested_parse(self.content, self.content_offset, impl_node)
		if self.arguments:
			n, m = self.state.inline_text(self.arguments[0], self.lineno)
			impl_node.append(nodes.paragraph('', '', *(n + m)))
		return [impl_node]

# -- Go Version role --------------------------------------------------
# Returns the value of the Go version stored in GO_VERSION to minor version
# precision.
import os
import re

from docutils.nodes import strong
from docutils.parsers.rst.states import Inliner
from typing import Tuple, List

def atc_go_version(unused_typ: str,
                   unused_rawtext: str,
                   unused_text: str,
                   unused_lineno: int,
                   unused_inliner: Inliner,
                   options=None,
                   content=None) -> Tuple[List[strong], list]:
	go_version_file = os.path.join(os.path.dirname(__file__), '../../../GO_VERSION')
	with open(file=go_version_file) as go_version_file:
		go_version = go_version_file.read()
	major_minor_version = re.match(pattern=r'\d+\.\d+', string=go_version).group()
	strong_node = nodes.strong(major_minor_version, major_minor_version)
	return [strong_node], []

# -- Issue role --------------------------------------------------------------
from docutils import utils

ISSUE_URI = REPO_URI + "issues/%s"

def issue_role(unused_typ,
               unused_rawtext,
               text,
               unused_lineno,
               unused_inliner,
               options=None,
               content=None):
	if options is None:
		options = {}
	if content is None:
		content = []

	issue = utils.unescape(text)
	text = 'Issue #' + issue
	refnode = nodes.reference(text, text, refuri=ISSUE_URI % issue)
	return [refnode], []

# -- Pull Request Role -------------------------------------------------------
PR_URI = REPO_URI + "pull/%s"

def pr_role(unused_typ,
            unused_rawtext,
            text,
            unused_lineno,
            unused_inliner,
            options=None,
            content=None):
	if options is None:
		options = {}
	if content is None:
		content = []

	pr = utils.unescape(text)
	text = 'Pull Request ' + pr
	refnode = nodes.reference(text, text, refuri=PR_URI % pr)
	return [refnode], []

# -- ATC file role -----------------------------------------------------------
FILE_URI = REPO_URI + "tree/master/%s"
def atc_file_role(unused_typ,
                  unused_rawtext,
                  text,
                  unused_lineno,
                  unused_inliner,
                  options=None,
                  content=None):
	if options is None:
		options = {}
	if content is None:
		content = []

	text = utils.unescape(text)
	litnode = nodes.literal(text, text)
	refnode = nodes.reference(text, '', litnode, refuri=FILE_URI % text)
	return [refnode], []

# -- GoDoc role (absolute) ---------------------------------------------------
GODOC_URI = "https://pkg.go.dev/"
def godoc_role(unused_typ,
               unused_rawtext,
               text,
               unused_lineno,
               unused_inliner,
               options=None,
               content=None):
	if options is None:
		options = {}
	if content is None:
		content = []

	text = utils.unescape(text)
	litnode = nodes.literal(text, text)
	last_period_index = text.rfind('.')
	if last_period_index != -1 and text.rfind('/') < last_period_index:
		text = text[:last_period_index] + '#' + text[last_period_index + 1:]
	refnode = nodes.reference(text, '', litnode, refuri=GODOC_URI + text)
	return [refnode], []

# -- GoDoc role (atc-relative) ----------------------------------------------
ATC_GODOC_PREFIX = "github.com/apache/trafficcontrol/v7/"
ATC_GODOC_URI = GODOC_URI + ATC_GODOC_PREFIX
def atc_godoc_role(unused_typ,
                   unused_rawtext,
                   text,
                   unused_lineno,
                   unused_inliner,
                   options=None,
                   content=None):
	if options is None:
		options = {}
	if content is None:
		content = []

	text = utils.unescape(text)
	literaltext = ATC_GODOC_PREFIX + text
	litnode = nodes.literal(literaltext, literaltext)
	refnode = nodes.reference(text, '', litnode, refuri=ATC_GODOC_URI + text.replace('.', '#', 1))
	return [refnode], []

# -- GoDoc role (to-relative) -----------------------------------------------
TO_GODOC_PREFIX = ATC_GODOC_PREFIX + "traffic_ops/traffic_ops_golang/"
TO_GODOC_URI = GODOC_URI + TO_GODOC_PREFIX
def to_godoc_role(unused_typ,
                  unused_rawtext,
                  text,
                  unused_lineno,
                  unused_inliner,
                  options=None,
                  content=None):
	if options is None:
		options = {}
	if content is None:
		content = []

	text = utils.unescape(text)
	literaltext = TO_GODOC_PREFIX + text
	litnode = nodes.literal(literaltext, literaltext)
	refnode = nodes.reference(text, '', litnode, refuri=TO_GODOC_URI + text.replace('.', '#', 1))
	return [refnode], []


def setup(app: object) -> dict:
	app.add_node(impl,
	             html=(visit_impl_node, depart_impl_node),
	             latex=(visit_impl_node, depart_impl_node),
	             text=(visit_impl_node, depart_impl_node))
	app.add_directive("impl-detail", ImplementationDetail)
	app.add_role("atc-go-version", atc_go_version)
	app.add_role("issue", issue_role)
	app.add_role("pr", pr_role)
	app.add_role("pull-request", pr_role)
	app.add_role("atc-file", atc_file_role)
	app.add_role("godoc", godoc_role)
	app.add_role("atc-godoc", atc_godoc_role)
	app.add_role("to-godoc", to_godoc_role)

	return {
		'version': '0.2',
		'parallel_read_safe': True,
		'parallel_write_safe': True,
	}
