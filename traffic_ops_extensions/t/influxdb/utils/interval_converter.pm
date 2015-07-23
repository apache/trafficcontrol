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
use strict;
use warnings;
no warnings 'once';
use warnings 'all';
use Extensions::InfluxDB::Utils::IntervalConverter;
use Data::Dumper;

my $t;
my $ic = new Extensions::InfluxDB::Utils::IntervalConverter();

$t = $ic->to_seconds("1000ms");
is( $t, 1 );

$t = $ic->to_seconds("60s");
is( $t, 60 );

$t = $ic->to_seconds("30s");
is( $t, 30 );

$t = $ic->to_seconds("1m");
is( $t, 60 );

$t = $ic->to_seconds("30m");
is( $t, 1800 );

$t = $ic->to_seconds("60m");
is( $t, 3600 );

$t = $ic->to_seconds("10m");
is( $t, 600 );

$t = $ic->to_seconds("1d");
is( $t, 86400 );

$t = $ic->to_seconds("3h");
is( $t, 10800 );

$t = $ic->to_seconds("2w");
is( $t, 1209600 );

$t = $ic->to_seconds("1mo");
is( $t, 2629744 );

$t = $ic->to_milliseconds("1ms");
is( $t, 1 );

$t = $ic->to_nanoseconds("1s");
is( $t, 1000000000 );

$t = $ic->to_seconds("1ns");
is( $t, 1e-10 );

done_testing();
