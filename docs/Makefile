# Makefile for Sphinx documentation
#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

ifeq ($(OS),Windows_NT)
    T3C_DIR_SEP := \\
	REM := del /s /q
	PY := python
else
    T3C_DIR_SEP := /
	REM := rm -rf
	PY := python3
endif

VERSION ?= $(file < ../VERSION)

# You can set these variables from the command line.
SPHINXOPTS    = -j auto
SPHINXBUILD   = $(PY) -m sphinx
PAPER         = letter
BUILDDIR      = build

# Internal variables.
ALLSPHINXOPTS   = -d $(BUILDDIR)/doctrees $(PAPEROPT_$(PAPER)) $(SPHINXOPTS) source
# the i18n builder cannot share the environment and doctrees with the others
I18NSPHINXOPTS  = $(PAPEROPT_$(PAPER)) $(SPHINXOPTS) source

.PHONY: all help clean html dirhtml singlehtml pickle json htmlhelp qthelp devhelp epub latex latexpdf text man changes linkcheck doctest gettext apichanges

T3C_SOURCE_DIRS := $(patsubst %$(T3C_DIR_SEP),%,$(dir $(wildcard ..$(T3C_DIR_SEP)cache-config$(T3C_DIR_SEP)**$(T3C_DIR_SEP)README.md)))
T3C_RST_SRC := $(join $(T3C_SOURCE_DIRS), $(T3C_SOURCE_DIRS:..$(T3C_DIR_SEP)cache-config$(T3C_DIR_SEP)%=$(T3C_DIR_SEP)%.rst))
T3C_RST_OUT := $(patsubst %,source/admin/t3c/$(T3C_DIR_SEP)%,$(notdir $(T3C_RST_SRC)))



# WARNING!!!!!!
# 1) do not use underscores, _, in source/conf.py
# 2) do not use incomplete tables in source; tables that contain blank areas between row cell borders cause latex to blow up.
# 3) do not use spaces in file names
# Failure to heed the above will break the latexpdf build target.
all: html

$(T3C_RST_SRC):
	@echo $(T3C_SOURCE_DIRS)
	@echo $(T3C_RST_SRC)
	@echo $(T3C_RST_OUT)
	$(MAKE) -C ../cache-config/ rst

$(T3C_RST_OUT): $(T3C_RST_SRC)
	cp -f $? source/admin/t3c/

help:
	@echo "Please use \`make <target>' where <target> is one of"
	@echo "  html       to make standalone HTML files"
	@echo "  dirhtml    to make HTML files named index.html in directories"
	@echo "  singlehtml to make a single large HTML file"
	@echo "  pickle     to make pickle files"
	@echo "  json       to make JSON files"
	@echo "  htmlhelp   to make HTML files and a HTML help project"
	@echo "  qthelp     to make HTML files and a qthelp project"
	@echo "  devhelp    to make HTML files and a Devhelp project"
	@echo "  epub       to make an epub"
	@echo "  latex      to make LaTeX files, you can set PAPER=a4 or PAPER=letter"
	@echo "  latexpdf   to make LaTeX files and run them through pdflatex"
	@echo "  text       to make text files"
	@echo "  man        to make manual pages"
	@echo "  texinfo    to make Texinfo files"
	@echo "  info       to make Texinfo files and run them through makeinfo"
	@echo "  gettext    to make PO message catalogs"
	@echo "  changes    to make an overview of all changed/added/deprecated items"
	@echo "  linkcheck  to check all external links for integrity"
	@echo "  doctest    to run all doctests embedded in the documentation (if enabled)"
	@echo "  coverage   to run coverage check of the documentation (if enabled)"

clean:
	-$(REM) $(BUILDDIR) $(T3C_RST_OUT)

$(BUILDDIR):
	-mkdir $(BUILDDIR)

html: $(BUILDDIR) $(T3C_RST_OUT)
	$(SPHINXBUILD) -b html $(ALLSPHINXOPTS) $(BUILDDIR)/html
	@echo
	@echo "Build finished. The HTML pages are in $(BUILDDIR)/html."

dirhtml: $(BUILDDIR) $(T3C_RST_OUT)
	$(SPHINXBUILD) -b dirhtml $(ALLSPHINXOPTS) $(BUILDDIR)/dirhtml
	@echo
	@echo "Build finished. The HTML pages are in $(BUILDDIR)/dirhtml."

singlehtml: $(BUILDDIR) $(T3C_RST_OUT)
	$(SPHINXBUILD) -b singlehtml $(ALLSPHINXOPTS) $(BUILDDIR)/singlehtml
	@echo
	@echo "Build finished. The HTML page is in $(BUILDDIR)/singlehtml."

pickle: $(BUILDDIR) $(T3C_RST_OUT)
	$(SPHINXBUILD) -b pickle $(ALLSPHINXOPTS) $(BUILDDIR)/pickle
	@echo
	@echo "Build finished; now you can process the pickle files."

json: $(BUILDDIR) $(T3C_RST_OUT)
	$(SPHINXBUILD) -b json $(ALLSPHINXOPTS) $(BUILDDIR)/json
	@echo
	@echo "Build finished; now you can process the JSON files."

htmlhelp: $(BUILDDIR) $(T3C_RST_OUT)
	$(SPHINXBUILD) -b htmlhelp $(ALLSPHINXOPTS) $(BUILDDIR)/htmlhelp
	@echo
	@echo "Build finished; now you can run HTML Help Workshop with the" \
	      ".hhp project file in $(BUILDDIR)/htmlhelp."

qthelp: $(BUILDDIR) $(T3C_RST_OUT)
	$(SPHINXBUILD) -b qthelp $(ALLSPHINXOPTS) $(BUILDDIR)/qthelp
	@echo
	@echo "Build finished; now you can run "qcollectiongenerator" with the" \
	      ".qhcp project file in $(BUILDDIR)/qthelp, like this:"
	@echo "# qcollectiongenerator $(BUILDDIR)/qthelp/IPCDN_Architecture.qhcp"
	@echo "To view the help file:"
	@echo "# assistant -collectionFile $(BUILDDIR)/qthelp/IPCDN_Architecture.qhc"

devhelp: $(BUILDDIR) $(T3C_RST_OUT)
	$(SPHINXBUILD) -b devhelp $(ALLSPHINXOPTS) $(BUILDDIR)/devhelp
	@echo
	@echo "Build finished."
	@echo "To view the help file:"
	@echo "# mkdir -p $$HOME/.local/share/devhelp/IPCDN_Architecture"
	@echo "# ln -s $(BUILDDIR)/devhelp $$HOME/.local/share/devhelp/IPCDN_Architecture"
	@echo "# devhelp"

epub: $(BUILDDIR) $(T3C_RST_OUT)
	$(SPHINXBUILD) -b epub $(ALLSPHINXOPTS) $(BUILDDIR)/epub
	@echo
	@echo "Build finished. The epub file is in $(BUILDDIR)/epub."

latex: $(BUILDDIR) $(T3C_RST_OUT)
	$(SPHINXBUILD) -b latex $(ALLSPHINXOPTS) $(BUILDDIR)/latex
	@echo
	@echo "Build finished; the LaTeX files are in $(BUILDDIR)/latex."
	@echo "Run \`make' in that directory to run these through (pdf)latex" \
	      "(use \`make latexpdf' here to do that automatically)."

latexpdf: $(BUILDDIR) latex $(T3C_RST_OUT)
	@echo "Running LaTeX files through pdflatex..."
	$(MAKE) -C $(BUILDDIR)/latex all-pdf
	@echo "pdflatex finished; the PDF files are in $(BUILDDIR)/latex."

text: $(BUILDDIR) $(T3C_RST_OUT)
	$(SPHINXBUILD) -b text $(ALLSPHINXOPTS) $(BUILDDIR)/text
	@echo
	@echo "Build finished. The text files are in $(BUILDDIR)/text."

man: $(BUILDDIR) $(T3C_RST_OUT)
	$(SPHINXBUILD) -b man $(ALLSPHINXOPTS) $(BUILDDIR)/man
	@echo
	@echo "Build finished. The manual pages are in $(BUILDDIR)/man."

texinfo: $(BUILDDIR) $(T3C_RST_OUT)
	$(SPHINXBUILD) -b texinfo $(ALLSPHINXOPTS) $(BUILDDIR)/texinfo
	@echo
	@echo "Build finished. The Texinfo files are in $(BUILDDIR)/texinfo."
	@echo "Run \`make' in that directory to run these through makeinfo" \
	      "(use \`make info' here to do that automatically)."

info: $(BUILDDIR) $(T3C_RST_OUT)
	$(SPHINXBUILD) -b texinfo $(ALLSPHINXOPTS) $(BUILDDIR)/texinfo
	@echo "Running Texinfo files through makeinfo..."
	make -C $(BUILDDIR)/texinfo info
	@echo "makeinfo finished; the Info files are in $(BUILDDIR)/texinfo."

gettext: $(BUILDDIR) $(T3C_RST_OUT)
	$(SPHINXBUILD) -b gettext $(I18NSPHINXOPTS) $(BUILDDIR)/locale
	@echo
	@echo "Build finished. The message catalogs are in $(BUILDDIR)/locale."

changes: $(BUILDDIR) $(T3C_RST_OUT)
	$(SPHINXBUILD) -b changes -D version=$(VERSION) $(ALLSPHINXOPTS) $(BUILDDIR)/changes
	@echo
	@echo "The overview file is in $(BUILDDIR)/changes."

apichanges: $(BUILDDIR)
	$(SPHINXBUILD) -b changes -d $(BUILDDIR)/doctrees $(PAPEROPT_$(PAPER)) $(SPHINXOPTS) -c source/ -D version=$(VERSION) source/api/ $(BUILDDIR)/apichanges

linkcheck: $(BUILDDIR) $(T3C_RST_OUT)
	$(SPHINXBUILD) -b linkcheck $(ALLSPHINXOPTS) $(BUILDDIR)/linkcheck
	@echo
	@echo "Link check complete; look for any errors in the above output " \
	      "or in $(BUILDDIR)/linkcheck/output.txt."

doctest: $(BUILDDIR) $(T3C_RST_OUT)
	$(SPHINXBUILD) -b doctest $(ALLSPHINXOPTS) $(BUILDDIR)/doctest
	@echo "Testing of doctests in the sources finished, look at the " \
	      "results in $(BUILDDIR)/doctest/output.txt."

coverage: $(BUILDDIR) $(T3C_RST_OUT)
	$(SPHINXBUILD) -b coverage $(ALLSPHINXOPTS) $(BUILDDIR)/coverage
	@echo "Testing of coverage in the sources finished, look at the " \
	      "results in $(BUILDDIR)/coverage/python.txt."
