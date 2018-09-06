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
# Based on CentOS 7.2
############################################################

# Check for proper invocation
PWD := $(strip $(lastword $(patsubst %/,%,$(notdir $(shell pwd)))))
makefile_dir := $(strip $(notdir $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))))

ifneq ($(PWD),$(makefile_dir))
$(error This makefile MUST be run from within its directory)
endif

BUILD_NUMBER := $(shell git rev-list HEAD 2>/dev/null | wc -l).$(shell git rev-parse --short=8 HEAD)
TC_VERSION := $(shell cat "../../VERSION")

ifeq ($(shell which rpm 2>/dev/null),)
RHEL_VERSION :=el7
$(info Couldn't find 'rpm' binary, building for RHEL version $(RHEL_VERSION))
else
RHEL_VERSION :=el$(shell rpm -q --qf "%{VERSION}" $(shell rpm -q --whatprovides redhat-release))
endif

SPECIAL_SAUCE := $(TC_VERSION)-$(BUILD_NUMBER).$(RHEL_VERSION).x86_64.rpm

TP_SOURCE := $(wildcard ../../traffic_portal/**/*)
ORT_SOURCE:= $(wildcard ../../traffic_ops/bin/*)
TO_SOURCE := $(wildcard ../../traffic_ops/**/*)
TO_SOURCE += $(wildcard ../../traffic_ops_db/**/*)
TM_SOURCE := $(wildcard ../../traffic_monitor/**/*)


.PHONY: clean very-clean all nearly-all

# Default target; builds all pre-requisite rpms from source trees
all: traffic_monitor/traffic_monitor.rpm traffic_portal/traffic_portal.rpm traffic_ops/traffic_ops.rpm mid/traffic_ops_ort.rpm edge/traffic_ops_ort.rpm

# Actual output rpm recipies
traffic_monitor/traffic_monitor.rpm: ../../dist/traffic_monitor-$(SPECIAL_SAUCE)
	cp -f $? $@
traffic_ops/traffic_ops.rpm: ../../dist/traffic_ops-$(SPECIAL_SAUCE)
	cp -f $? $@
mid/traffic_ops_ort.rpm edge/traffic_ops_ort.rpm: ../../dist/traffic_ops_ort-$(SPECIAL_SAUCE)
	cp -f $? $@
traffic_portal/traffic_portal.rpm: ../../dist/traffic_portal-$(SPECIAL_SAUCE)
	cp -f $? $@

# Dist rpms
../../dist/traffic_monitor-$(SPECIAL_SAUCE): $(TM_SOURCE)
	sudo ../../pkg -v traffic_monitor_build

../../dist/traffic_ops-$(SPECIAL_SAUCE): $(TO_SOURCE)
	sudo ../../pkg -v traffic_ops_build

../../dist/traffic_ops_ort-$(SPECIAL_SAUCE): $(ORT_SOURCE)
	sudo ../../pkg -v traffic_ops_build

../../dist/traffic_portal-$(SPECIAL_SAUCE): $(TM_SOURCE)
	sudo ../../pkg -v traffic_portal_build

clean:
	$(RM) traffic_ops/traffic_ops.rpm traffic_monitor/traffic_monitor.rpm traffic_portal/traffic_portal.rpm edge/traffic_ops_ort.rpm mid/traffic_ops_ort.rpm

very-clean: clean
	$(warning This will destroy ALL OUTPUT RPMS IN 'dist'. Please be sure this is what you want)
	sleep 2 # Give users a sceond to cancel
	$(RM) ../../dist/*
