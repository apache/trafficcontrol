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

# Check for proper invocation
PWD := $(strip $(lastword $(patsubst %/,%,$(notdir $(shell pwd)))))
makefile_dir := $(strip $(notdir $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))))

ifneq ($(PWD),$(makefile_dir))
$(error This makefile MUST be run from within its directory)
endif

PKG_COMMAND := ../../pkg
PKG_FLAGS := -v
BUILD_SUFFIX := _build
BUILD_NUMBER := $(shell git rev-list HEAD 2>/dev/null | wc -l | tr -d '[[:space:]]').$(shell git rev-parse --short=8 HEAD)
TC_VERSION := $(shell cat "../../VERSION")
TOMCAT_VERSION := $(shell grep '^\s*TOMCAT_VERSION=' ../../traffic_router/build/build_rpm.sh  | cut -d= -f2)
TOMCAT_RELEASE := $(shell grep '^\s*TOMCAT_RELEASE=' ../../traffic_router/build/build_rpm.sh  | cut -d= -f2)

RHEL_VERSION := el7

SPECIAL_SAUCE := $(TC_VERSION)-$(BUILD_NUMBER).$(RHEL_VERSION).x86_64.rpm
SPECIAL_SEASONING := $(TOMCAT_VERSION).$(TOMCAT_RELEASE)-$(BUILD_NUMBER).$(RHEL_VERSION).x86_64.rpm

TO_SOURCE := $(wildcard ../../traffic_ops/**/*)
TO_SOURCE += $(wildcard ../../traffic_ops_db/**/*)
ORT_SOURCE:= $(wildcard ../../traffic_ops_ort/atstccfg/**/*.go)
ORT_SOURCE+= $(wildcard ../../lib/**/*.go)
TM_SOURCE := $(wildcard ../../traffic_monitor/**/*)
TP_SOURCE := $(wildcard ../../traffic_portal/**/*)
TR_SOURCE := $(wildcard ../../traffic_router/**/*)
TS_SOURCE := $(wildcard ../../traffic_stats/**/*)

.PHONY: clean very-clean all nearly-all debug native

# Default target; builds all pre-requisite rpms from source trees
all: cache/traffic_ops_ort.rpm traffic_monitor/traffic_monitor.rpm traffic_portal/traffic_portal.rpm traffic_ops/traffic_ops.rpm traffic_router/traffic_router.rpm traffic_router/tomcat.rpm traffic_stats/traffic_stats.rpm

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
PKG_COMMAND := ../../build/clean_build.sh
PKG_FLAGS :=
BUILD_SUFFIX :=
ifneq ($(MAKECMDGOALS),native)
MAKECMDGOALS := $(filter-out native,$(MAKECMDGOALS))
native: $(MAKECMDGOALS)
else
native: all
endif
endif

#.PHONY: native

# Actual output rpm recipies
traffic_monitor/traffic_monitor.rpm: ../../dist/traffic_monitor-$(SPECIAL_SAUCE)
	cp -f $? $@
traffic_ops/traffic_ops.rpm: ../../dist/traffic_ops-$(SPECIAL_SAUCE)
	cp -f $? $@
traffic_portal/traffic_portal.rpm: ../../dist/traffic_portal-$(SPECIAL_SAUCE)
	cp -f $? $@
traffic_router/traffic_router.rpm: ../../dist/traffic_router-$(SPECIAL_SAUCE)
	cp -f $? $@
traffic_router/tomcat.rpm: ../../dist/tomcat-$(SPECIAL_SEASONING)
	cp -f $? $@
traffic_stats/traffic_stats.rpm: ../../dist/traffic_stats-$(SPECIAL_SAUCE)
	cp -f $? $@
cache/traffic_ops_ort.rpm: ../../dist/traffic_ops_ort-$(SPECIAL_SAUCE)
	cp -f $? $@

# Dist rpms
../../dist/traffic_monitor-$(SPECIAL_SAUCE): $(TM_SOURCE)
	$(PKG_COMMAND) $(PKG_FLAGS) traffic_monitor$(BUILD_SUFFIX)

../../dist/traffic_ops-$(SPECIAL_SAUCE): $(TO_SOURCE)
	$(PKG_COMMAND) $(PKG_FLAGS) traffic_ops$(BUILD_SUFFIX)

../../dist/traffic_portal-$(SPECIAL_SAUCE): $(TP_SOURCE)
	$(PKG_COMMAND) $(PKG_FLAGS) traffic_portal$(BUILD_SUFFIX)

../../dist/traffic_rou%er-$(SPECIAL_SAUCE) ../../dist/tomca%-$(SPECIAL_SEASONING): $(TR_SOURCE)
	$(PKG_COMMAND) $(PKG_FLAGS) traffic_router$(BUILD_SUFFIX)

../../dist/traffic_stats-$(SPECIAL_SAUCE): $(TS_SOURCE)
	$(PKG_COMMAND) $(PKG_FLAGS) traffic_stats$(BUILD_SUFFIX)

../../dist/traffic_ops_ort-$(SPECIAL_SAUCE): $(ORT_SOURCE)
	$(PKG_COMMAND) $(PKG_FLAGS) traffic_ops_ort$(BUILD_SUFFIX)

clean:
	$(RM) traffic_monitor/traffic_monitor.rpm traffic_ops/traffic_ops.rpm traffic_portal/traffic_portal.rpm traffic_router/traffic_router.rpm traffic_router/tomcat.rpm cache/traffic_ops_ort.rpm traffic_stats/traffic_stats.rpm

very-clean: clean
	$(warning This will destroy ALL OUTPUT RPMS IN 'dist'. Please be sure this is what you want)
	sleep 2 # Give users a sceond to cancel
	$(RM) ../../dist/*
