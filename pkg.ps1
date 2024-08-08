#!/usr/bin/env pwsh
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

<#
.SYNOPSIS
    Builds Traffic Control components.
.DESCRIPTION
	A convenience script for generating RPM packages using Docker.
	If no projects are listed, all projects will be built.
.PARAMETER q
    Quiet mode. Supresses output.
.PARAMETER v
	Verbose mode. Lists all build output.
.PARAMETER l
	List available projects.
.EXAMPLE
    C:\Users\admin\go\src\github.com\apache\trafficcontrol> .\pkg.ps1 -l
    source
	traffic_monitor_build
	traffic_ops_build
	traffic_portal_build
	traffic_router_build
	traffic_stats_build
	grove_build
	grovetccfg_build
	weasel
	docs
.NOTES
    Author: Apache Software Foundation
    Date:   July 30, 2019
#>
param(
	[switch]$q,
	[switch]$v,
	[switch]$l
)

$DOCKER = (Get-Command docker).Source;
if ($DOCKER -eq "") {
	Write-Error "docker is required for a docker build." -Category ObjectNotFound;
	exit 1;
}

Push-Location $PSScriptRoot;
$COMPOSE_FILE = "./infrastructure/docker/build/docker-compose.yml";
$COMPOSE = (Get-Command docker compose).Source;
$COMPOSE_ARGS = ""
if ($COMPOSE -eq "") {
	& $DOCKER "inspect docker-compose:latest";
	if ($? -eq $false) {
		& $DOCKER "pull docker-compose:latest";
		if ($? -eq $false) {
			Write-Error "Couldn't pull docker compose - please connect to the internet or install docker-compose." -Category NotInstalled;
			exit 1;
		}
	}
	$USERHOME = $home;
	$COMPOSE = $DOCKER;
	$COMPOSE_ARGS = "run --rm '$env:COMPOSE_OPTIONS' -v '${PSScriptRoot}:$PSScriptRoot' -v '${USERHOME}:/root'";
}

$PROJECTS = & $COMPOSE $COMPOSE_ARGS "-f" $COMPOSE_FILE "config" "--services" | Out-String;

if ($l) {
	Write-Host $PROJECTS;
	exit 0;
}

if ($args.Length -eq 0) {
	$args = $PROJECTS.Split([Environment]::Newline);
}

$code=$true;
for ($i=0; $i -lt $args.Length; $i++) {
	$PROJECT = $args[$i];
	if ($PROJECT -eq "") {
		continue;
	}
	Write-Host "Building $PROJECT";

	if ($v) {
		& $COMPOSE $COMPOSE_ARGS "-f" $COMPOSE_FILE "build" $PROJECT | Write-Verbose;
		if ($?) {
			& $COMPOSE $COMPOSE_ARGS "-f" $COMPOSE_FILE "run" "--rm" $PROJECT | Write-Verbose;
			if ($? -eq $false) {
				Write-Error "Failed to execute $PROJECT";
				$code=$false;
			}
		} else {
			Write-Error "Failed to build $PROJECT";
			$code=$false;
		}
	} else {
		& $COMPOSE $COMPOSE_ARGS "-f" $COMPOSE_FILE "build" $PROJECT | Out-Null;
		if ($?) {
			& $COMPOSE $COMPOSE_ARGS "-f" $COMPOSE_FILE "run" "--rm" $PROJECT | Out-Null;
			if ($? -eq $false) {
				Write-Error "Failed to execute $PROJECT";
				$code=$false;
			}
		} else {
			Write-Error "Failed to build $PROJECT";
			$code=$false;
		}
	}
}

exit $code;
