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

############################################################
# Dockerfile to build Edge-Tier Cache container images for
# Apache Traffic Control
# Based on Rocky Linux 8
############################################################

ifeq ($(RHEL_VERSION),)
	export RHEL_VERSION := 8
else
	RHEL_VERSION := $(shell echo $(RHEL_VERSION) | cut -f1 -d.)
endif

ifeq ($(shell echo $(RHEL_VERSION) | grep -Eo '^[0-9]+'),7)
	export BASE_IMAGE := centos
else
	BASE_IMAGE := rockylinux
endif

CIAB_DIR_RELATIVE := $(dir $(MAKEFILE_LIST))
CIAB_DIR_ABSOLUTE := $(shell cd $(CIAB_DIR_RELATIVE) && pwd)
TC_DIR := $(CIAB_DIR_RELATIVE)../..

PKG_COMMAND := $(TC_DIR)/pkg
PKG_FLAGS := -v -$(RHEL_VERSION)
BUILD_SUFFIX := _build
BUILD_NUMBER := $(shell bash -c ' \
	set -o errexit -o pipefail -o xtrace; \
	git rev-parse --is-inside-work-tree >/dev/null; \
	output="$$(git describe \
		--long \
		--tags \
		--match="RELEASE-[0-9].[0-9].[0-9]" \
		--match="RELEASE-[0-9][0-9].[0-9][0-9].[0-9][0-9]" \
		--match="v[0-9].[0-9].[0-9]" \
		--match="v[0-9][0-9].[0-9][0-9].[0-9][0-9]" | \
		awk -F- "{print \$$(NF-1)}" || \
		echo 0)"; \
	output+=".$$(git rev-parse --short=8 HEAD)"; \
	echo "$$output"' || \
	cat $(TC_DIR)/BUILD_NUMBER)
BUILD_ARCH   := $(shell docker create --name=ciab-get-$(BASE_IMAGE)-$(RHEL_VERSION)-arch --rm $(BASE_IMAGE):$(RHEL_VERSION) rpm --eval %_arch >/dev/null 2>&1 || true && docker start --attach ciab-get-$(BASE_IMAGE)-$(RHEL_VERSION)-arch )
TC_VERSION := $(shell cat "$(TC_DIR)/VERSION")
TOMCAT_VERSION := $(shell grep '^\s*TOMCAT_VERSION=' "$(TC_DIR)/.env"  | cut -d= -f2)
ATS_VERSION := $(shell $(CIAB_DIR_RELATIVE)bin/ats-version.sh)
ifeq ($(RHEL_VERSION),7)
# The commit hash in `git describe` is only 7 characters in CentOS 7's git 1.8
ATS_VERSION := $(shell echo "$(ATS_VERSION)" | sed -E 's/^(.*-[0-9]+\.[0-9a-f]{7})[0-9a-f]*/\1/')
endif

SPECIAL_SAUCE := $(TC_VERSION)-$(BUILD_NUMBER).el$(RHEL_VERSION).$(BUILD_ARCH).rpm
SPECIAL_SAUCE_NOARCH := $(TC_VERSION)-$(BUILD_NUMBER).el$(RHEL_VERSION).noarch.rpm
SPECIAL_SEASONING := $(TOMCAT_VERSION)-1.el$(RHEL_VERSION).noarch.rpm
SPECIAL_SYRUP := $(ATS_VERSION).el$(RHEL_VERSION).$(BUILD_ARCH).rpm

ATS_SOURCE := $(wildcard $(TC_DIR)/cache-config/testing/docker/trafficserver/**)
TO_SOURCE := $(wildcard $(TC_DIR)/traffic_ops/**)
TO_SOURCE += $(wildcard $(TC_DIR)/traffic_ops_db/**)
ORT_SOURCE := $(wildcard $(TC_DIR)/cache-config/**)
ORT_SOURCE += $(wildcard $(TC_DIR)/lib/**.go)
TM_SOURCE := $(wildcard $(TC_DIR)/traffic_monitor/**)
TP_SOURCE := $(wildcard $(TC_DIR)/traffic_portal/**)
TP2_SOURCE := $(wildcard $(TC_DIR)/experimental/traffic-portal/**)
TR_SOURCE := $(wildcard $(TC_DIR)/traffic_router/**)
TS_SOURCE := $(wildcard $(TC_DIR)/traffic_stats/**)
TCH_SOURCE := $(wildcard $(TC_DIR)/tc-health-client/**)

ATS_RPM := cache/trafficserver.rpm
TM_RPM := traffic_monitor/traffic_monitor.rpm
TO_RPM := traffic_ops/traffic_ops.rpm
TP_RPM := traffic_portal/traffic_portal.rpm
TP2_RPM := optional/traffic_portal_v2/traffic_portal_v2.rpm
TR_RPM := traffic_router/traffic_router.rpm
TOMCAT_RPM := traffic_router/tomcat.rpm
TS_RPM := traffic_stats/traffic_stats.rpm
ORT_RPM := cache/trafficcontrol-cache-config.rpm
TCH_RPM := health/trafficcontrol-health-client.rpm

ATS_RPM_RELATIVE := $(CIAB_DIR_RELATIVE)$(ATS_RPM)
TM_RPM_RELATIVE := $(CIAB_DIR_RELATIVE)$(TM_RPM)
TO_RPM_RELATIVE := $(CIAB_DIR_RELATIVE)$(TO_RPM)
TP_RPM_RELATIVE := $(CIAB_DIR_RELATIVE)$(TP_RPM)
TP2_RPM_RELATIVE := $(CIAB_DIR_RELATIVE)$(TP2_RPM)
TR_RPM_RELATIVE := $(CIAB_DIR_RELATIVE)$(TR_RPM)
TOMCAT_RPM_RELATIVE := $(CIAB_DIR_RELATIVE)$(TOMCAT_RPM)
TS_RPM_RELATIVE := $(CIAB_DIR_RELATIVE)$(TS_RPM)
ORT_RPM_RELATIVE := $(CIAB_DIR_RELATIVE)$(ORT_RPM)
TCH_RPM_RELATIVE := $(CIAB_DIR_RELATIVE)$(TCH_RPM)

ATS_RPM_ABSOLUTE := $(CIAB_DIR_ABSOLUTE)/$(ATS_RPM)
TM_RPM_ABSOLUTE := $(CIAB_DIR_ABSOLUTE)/$(TM_RPM)
TO_RPM_ABSOLUTE := $(CIAB_DIR_ABSOLUTE)/$(TO_RPM)
TP_RPM_ABSOLUTE := $(CIAB_DIR_ABSOLUTE)/$(TP_RPM)
TP2_RPM_ABSOLUTE := $(CIAB_DIR_ABSOLUTE)/$(TP2_RPM)
TR_RPM_ABSOLUTE := $(CIAB_DIR_ABSOLUTE)/$(TR_RPM)
TOMCAT_RPM_ABSOLUTE := $(CIAB_DIR_ABSOLUTE)/$(TOMCAT_RPM)
TS_RPM_ABSOLUTE := $(CIAB_DIR_ABSOLUTE)/$(TS_RPM)
ORT_RPM_ABSOLUTE := $(CIAB_DIR_ABSOLUTE)/$(ORT_RPM)
TCH_RPM_ABSOLUTE := $(CIAB_DIR_ABSOLUTE)/$(TCH_RPM)

ATS_DIST_RPM := $(TC_DIR)/dist/trafficserver-$(SPECIAL_SYRUP)
TM_DIST_RPM := $(TC_DIR)/dist/traffic_monitor-$(SPECIAL_SAUCE)
TO_DIST_RPM := $(TC_DIR)/dist/traffic_ops-$(SPECIAL_SAUCE)
TP_DIST_RPM := $(TC_DIR)/dist/traffic_portal-$(SPECIAL_SAUCE)
TP2_DIST_RPM := $(TC_DIR)/dist/traffic_portal_v2-$(SPECIAL_SAUCE)
TR_DIST_RPM := $(TC_DIR)/dist/traffic_router-$(SPECIAL_SAUCE_NOARCH)
TOMCAT_DIST_RPM := $(TC_DIR)/dist/tomcat-$(SPECIAL_SEASONING)
TS_DIST_RPM := $(TC_DIR)/dist/traffic_stats-$(SPECIAL_SAUCE)
ORT_DIST_RPM := $(TC_DIR)/dist/trafficcontrol-cache-config-$(SPECIAL_SAUCE)
TCH_DIST_RPM := $(TC_DIR)/dist/trafficcontrol-health-client-$(SPECIAL_SAUCE)

.PHONY: all build-builders clean debug native nearly-all pull-builders very-clean

# Default target; builds all pre-requisite rpms from source trees
all: $(ATS_RPM_RELATIVE) $(ORT_RPM_RELATIVE) $(TM_RPM_RELATIVE) $(TP_RPM_RELATIVE) $(TO_RPM_RELATIVE) $(TR_RPM_RELATIVE) $(TOMCAT_RPM_RELATIVE) $(TS_RPM_RELATIVE) $(TCH_RPM_RELATIVE)

ifneq ($(filter build-builders,$(MAKECMDGOALS)),)
PKG_FLAGS += -b
ifneq ($(MAKECMDGOALS),build-builders)
MAKECMDGOALS := $(filter-out build-builders,$(MAKECMDGOALS))
build-builders: $(MAKECMDGOALS)
else
build-builders: all
endif
endif

ifneq ($(filter debug,$(MAKECMDGOALS)),)
PKG_FLAGS += -d
export DEBUG_BUILD = true
ifneq ($(MAKECMDGOALS),debug)
MAKECMDGOALS := $(filter-out debug,$(MAKECMDGOALS))
debug: $(MAKECMDGOALS)
else
debug: all
endif
endif

ifneq ($(filter native,$(MAKECMDGOALS)),)
PKG_COMMAND := $(TC_DIR)/build/clean_build.sh
PKG_FLAGS :=
BUILD_SUFFIX :=
ifneq ($(MAKECMDGOALS),native)
MAKECMDGOALS := $(filter-out native,$(MAKECMDGOALS))
native: $(MAKECMDGOALS)
else
native: all
endif
endif

ifneq ($(filter pull-builders,$(MAKECMDGOALS)),)
PKG_FLAGS += -p
ifneq ($(MAKECMDGOALS),pull-builders)
MAKECMDGOALS := $(filter-out pull-builders,$(MAKECMDGOALS))
pull-builders: $(MAKECMDGOALS)
else
pull-builders: all
endif
endif

# Relative path output RPM recipies
$(ATS_RPM_RELATIVE): $(ATS_DIST_RPM)
	cp -f "$?" "$@" || (rm "$(CIAB_DIR_RELATIVE)/cache/ATS_VERSION"; false)
$(TM_RPM_RELATIVE): $(TM_DIST_RPM)
	cp -f "$?" "$@"
$(TO_RPM_RELATIVE): $(TO_DIST_RPM)
	cp -f "$?" "$@"
$(TP_RPM_RELATIVE): $(TP_DIST_RPM)
	cp -f "$?" "$@"
$(TP2_RPM_RELATIVE): $(TP2_DIST_RPM)
	cp -f "$?" "$@"
$(TR_RPM_RELATIVE): $(TR_DIST_RPM)
	cp -f "$?" "$@"
$(TOMCAT_RPM_RELATIVE): $(TOMCAT_DIST_RPM)
	cp -f "$?" "$@"
$(TS_RPM_RELATIVE): $(TS_DIST_RPM)
	cp -f "$?" "$@"
$(ORT_RPM_RELATIVE): $(ORT_DIST_RPM)
	cp -f "$?" "$@"
$(TCH_RPM_RELATIVE): $(TCH_DIST_RPM)
	cp -f "$?" "$@"

# Absolute path output RPM recipies
$(ATS_RPM_ABSOLUTE): $(ATS_DIST_RPM)
	cp -f "$?" "$@" || (rm "$(CIAB_DIR_ABSOLUTE)/cache/ATS_VERSION"; false)
$(TM_RPM_ABSOLUTE): $(TM_DIST_RPM)
	cp -f "$?" "$@"
$(TO_RPM_ABSOLUTE): $(TO_DIST_RPM)
	cp -f "$?" "$@"
$(TP_RPM_ABSOLUTE): $(TP_DIST_RPM)
	cp -f "$?" "$@"
$(TP2_RPM_ABSOLUTE): $(TP2_DIST_RPM)
	cp -f "$?" "$@"
$(TR_RPM_ABSOLUTE): $(TR_DIST_RPM)
	cp -f "$?" "$@"
$(TOMCAT_RPM_ABSOLUTE): $(TOMCAT_DIST_RPM)
	cp -f "$?" "$@"
$(TS_RPM_ABSOLUTE): $(TS_DIST_RPM)
	cp -f "$?" "$@"
$(ORT_RPM_ABSOLUTE): $(ORT_DIST_RPM)
	cp -f "$?" "$@"
$(TCH_RPM_ABSOLUTE): $(TCH_DIST_RPM)
	cp -f "$?" "$@"

# Dist rpms
$(ATS_DIST_RPM): $(ATS_SOURCE)
	docker compose -f $(TC_DIR)/cache-config/testing/docker/docker-compose-ats-build.yml build --parallel trafficserver_build && docker compose -f $(TC_DIR)/cache-config/testing/docker/docker-compose-ats-build.yml run --rm trafficserver_build

$(TM_DIST_RPM): $(TM_SOURCE)
	"$(PKG_COMMAND)" $(PKG_FLAGS) traffic_monitor$(BUILD_SUFFIX)

$(TO_DIST_RPM): $(TO_SOURCE)
	"$(PKG_COMMAND)" $(PKG_FLAGS) traffic_ops$(BUILD_SUFFIX)

$(TP_DIST_RPM): $(TP_SOURCE)
	"$(PKG_COMMAND)" $(PKG_FLAGS) traffic_portal$(BUILD_SUFFIX)

$(TP2_DIST_RPM): $(TP2_SOURCE)
	"$(PKG_COMMAND)" -o $(PKG_FLAGS) traffic_portal_v2$(BUILD_SUFFIX)

$(TR_DIST_RPM) $(TOMCAT_DIST_RPM): $(TR_SOURCE)
	"$(PKG_COMMAND)" $(PKG_FLAGS) traffic_router$(BUILD_SUFFIX)

$(TS_DIST_RPM): $(TS_SOURCE)
	"$(PKG_COMMAND)" $(PKG_FLAGS) traffic_stats$(BUILD_SUFFIX)

$(ORT_DIST_RPM): $(ORT_SOURCE)
	"$(PKG_COMMAND)" $(PKG_FLAGS) cache-config$(BUILD_SUFFIX)

$(TCH_DIST_RPM): $(TCH_SOURCE)
	"$(PKG_COMMAND)" $(PKG_FLAGS) tc-health-client$(BUILD_SUFFIX)

clean:
	cd "$(CIAB_DIR_RELATIVE)"
	$(RM) $(TM_RPM_RELATIVE) $(TO_RPM_RELATIVE) $(TP_RPM_RELATIVE) $(TP2_RPM_RELATIVE) $(TR_RPM_RELATIVE) $(TOMCAT_RPM_RELATIVE) $(ORT_RPM_RELATIVE) $(TS_RPM_RELATIVE) $(TCH_RPM_RELATIVE)

very-clean: clean
	$(warning This will destroy ALL OUTPUT RPMS IN 'dist'. Please be sure this is what you want)
	sleep 2 # Give users a second to cancel
	$(RM) -r "$(TC_DIR)/dist"/*
