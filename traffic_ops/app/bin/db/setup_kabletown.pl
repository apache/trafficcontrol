#!/usr/bin/perl

package main;
#
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
use Mojo::Base -strict;
use Schema;
use Test::IntegrationTestHelper;
use strict;
use warnings;

# If MOJO_MODE not already defined, default to development
BEGIN { $ENV{MOJO_MODE} //= "development" }
my $schema = Schema->connect_to_database;
my $t      = Test::Mojo->new('TrafficOps');

Test::IntegrationTestHelper->unload_core_data($schema);
Test::IntegrationTestHelper->load_core_data($schema);
