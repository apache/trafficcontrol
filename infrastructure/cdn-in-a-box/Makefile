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
# Based on CentOS 8
############################################################

ifeq ($(RHEL_VERSION),)
	export RHEL_VERSION := 8
else
	RHEL_VERSION := $(shell echo $(RHEL_VERSION) | cut -f1 -d.)
endif

CIAB_DIR := $(realpath $(dir $(MAKEFILE_LIST)))
TC_DIR := $(shell cd "$(CIAB_DIR)/../.." && pwd -P)

PKG_COMMAND := $(TC_DIR)/pkg
PKG_FLAGS := -v -$(RHEL_VERSION)
BUILD_SUFFIX := _build
BUILD_NUMBER := $(shell git rev-list HEAD 2>/dev/null | wc -l | tr -d '[[:space:]]').$(shell git rev-parse --short=8 HEAD)
TC_VERSION := $(shell cat "$(TC_DIR)/VERSION")
TOMCAT_VERSION := $(shell grep '^\s*TOMCAT_VERSION=' "$(TC_DIR)/traffic_router/build/build_rpm.sh"  | cut -d= -f2)
TOMCAT_RELEASE := $(shell grep '^\s*TOMCAT_RELEASE=' "$(TC_DIR)/traffic_router/build/build_rpm.sh"  | cut -d= -f2)

SPECIAL_SAUCE := $(TC_VERSION)-$(BUILD_NUMBER).el$(RHEL_VERSION).x86_64.rpm
SPECIAL_SEASONING := $(TOMCAT_VERSION).$(TOMCAT_RELEASE)-$(BUILD_NUMBER).el$(RHEL_VERSION).x86_64.rpm

TO_SOURCE := $(wildcard $(TC_DIR)/traffic_ops/**)
TO_SOURCE += $(wildcard $(TC_DIR)/traffic_ops_db/**)
ORT_SOURCE:= $(wildcard $(TC_DIR)/traffic_ops_ort/atstccfg/**.go)
ORT_SOURCE+= $(wildcard $(TC_DIR)/lib/**.go)
TM_SOURCE := $(wildcard $(TC_DIR)/traffic_monitor/**)
TP_SOURCE := $(wildcard $(TC_DIR)/traffic_portal/**)
TR_SOURCE := $(wildcard $(TC_DIR)/traffic_router/**)
TS_SOURCE := $(wildcard $(TC_DIR)/traffic_stats/**)

TM_RPM := $(CIAB_DIR)/traffic_monitor/traffic_monitor.rpm
TO_RPM := $(CIAB_DIR)/traffic_ops/traffic_ops.rpm
TP_RPM := $(CIAB_DIR)/traffic_portal/traffic_portal.rpm
TR_RPM := $(CIAB_DIR)/traffic_router/traffic_router.rpm
TOMCAT_RPM := $(CIAB_DIR)/traffic_router/tomcat.rpm
TS_RPM := $(CIAB_DIR)/traffic_stats/traffic_stats.rpm
ORT_RPM := $(CIAB_DIR)/cache/traffic_ops_ort.rpm

TM_DIST_RPM := $(TC_DIR)/dist/traffic_monitor-$(SPECIAL_SAUCE)
TO_DIST_RPM := $(TC_DIR)/dist/traffic_ops-$(SPECIAL_SAUCE)
TP_DIST_RPM := $(TC_DIR)/dist/traffic_portal-$(SPECIAL_SAUCE)
TR_DIST_RPM := $(TC_DIR)/dist/traffic_router-$(SPECIAL_SAUCE)
TOMCAT_DIST_RPM := $(TC_DIR)/dist/tomcat-$(SPECIAL_SEASONING)
TS_DIST_RPM := $(TC_DIR)/dist/traffic_stats-$(SPECIAL_SAUCE)
ORT_DIST_RPM := $(TC_DIR)/dist/traffic_ops_ort-$(SPECIAL_SAUCE)

.PHONY: all build-builders clean debug native nearly-all pull-builders very-clean

# Default target; builds all pre-requisite rpms from source trees
all: $(ORT_RPM) $(TM_RPM) $(TP_RPM) $(TO_RPM) $(TR_RPM) $(TOMCAT_RPM) $(TS_RPM)

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

# Actual output rpm recipies
$(TM_RPM): $(TM_DIST_RPM)
	cp -f "$?" "$@"
$(TO_RPM): $(TO_DIST_RPM)
	cp -f "$?" "$@"
$(TP_RPM): $(TP_DIST_RPM)
	cp -f "$?" "$@"
$(TR_RPM): $(TR_DIST_RPM)
	cp -f "$?" "$@"
$(TOMCAT_RPM): $(TOMCAT_DIST_RPM)
	cp -f "$?" "$@"
$(TS_RPM): $(TS_DIST_RPM)
	cp -f "$?" "$@"
$(ORT_RPM): $(ORT_DIST_RPM)
	cp -f "$?" "$@"

# RPM target aliases
traffic_monitor/traffic_monitor.rpm: $(TM_RPM)
traffic_ops/traffic_ops.rpm: $(TO_RPM)
traffic_portal/traffic_portal.rpm: $(TP_RPM)
traffic_router/traffic_router.rpm: $(TR_RPM)
traffic_router/tomcat.rpm: $(TOMCAT_RPM)
traffic_stats/traffic_stats.rpm: $(TS_RPM)
cache/traffic_ops_ort.rpm: $(ORT_RPM)

# Dist rpms
$(TM_DIST_RPM): $(TM_SOURCE)
	"$(PKG_COMMAND)" $(PKG_FLAGS) traffic_monitor$(BUILD_SUFFIX)

$(TO_DIST_RPM): $(TO_SOURCE)
	"$(PKG_COMMAND)" $(PKG_FLAGS) traffic_ops$(BUILD_SUFFIX)

$(TP_DIST_RPM): $(TP_SOURCE)
	"$(PKG_COMMAND)" $(PKG_FLAGS) traffic_portal$(BUILD_SUFFIX)

$(TR_DIST_RPM) $(TOMCAT_DIST_RPM): $(TR_SOURCE)
	"$(PKG_COMMAND)" $(PKG_FLAGS) traffic_router$(BUILD_SUFFIX)

$(TS_DIST_RPM): $(TS_SOURCE)
	"$(PKG_COMMAND)" $(PKG_FLAGS) traffic_stats$(BUILD_SUFFIX)

$(ORT_DIST_RPM): $(ORT_SOURCE)
	"$(PKG_COMMAND)" $(PKG_FLAGS) traffic_ops_ort$(BUILD_SUFFIX)

clean:
	cd "$(CIAB_DIR)"
	$(RM) $(TM_RPM) $(TO_RPM) $(TP_RPM) $(TR_RPM) $(TOMCAT_RPM) $(ORT_RPM) $(TS_RPM)

very-clean: clean
	$(warning This will destroy ALL OUTPUT RPMS IN 'dist'. Please be sure this is what you want)
	sleep 2 # Give users a second to cancel
	$(RM) -r "$(TC_DIR)/dist"/*
