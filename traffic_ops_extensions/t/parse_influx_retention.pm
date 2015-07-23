package main;
#
# Copyright 2015 Comcast Cable Communications Management, LLC
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
use Test::More;
use Test::Mojo;
use DBI;
use Data::Dumper;
use strict;
use warnings;
use JSON;
use Utils::Helper::Extensions;
Utils::Helper::Extensions->use;

BEGIN { $ENV{MOJO_MODE} = "test" }

my $ir                          = new Extensions::InfluxDB::Helper::InfluxResponse();
my $retention_period_in_seconds = $ir->parse_retention_period_in_seconds("120h0m0s");
is( $retention_period_in_seconds, 432000 );

$retention_period_in_seconds = $ir->parse_retention_period_in_seconds("120h2m0s");
is( $retention_period_in_seconds, 432120 );

# Test hours+minutes
$retention_period_in_seconds = $ir->parse_retention_period_in_seconds("120h2m45s");
is( $retention_period_in_seconds, 432165 );
done_testing();
