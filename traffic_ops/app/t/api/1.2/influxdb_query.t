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
use strict;
use warnings;
no warnings 'once';
use warnings 'all';
use Test::TestHelper;
use Test::MockModule;
use Connection::InfluxDBAdapter;
use Data::Dumper;
use Builder::DeliveryServiceStatsQuery;

BEGIN {
	use_ok('Test::Exception');
}

my $iq = Builder::DeliveryServiceStatsQuery->new(
	{
		ds_name     => "ds_stats",
		series_name => "kbps",
		start_date  => "2015-01-01T00:00:00-07:00",
		end_date    => "2015-01-30T00:00:00-07:00",
		interval    => "60s",
		limit       => 10
	}
);

undef $\;
my $summary_q = $iq->summary_query();

=cut
SELECT mean(value), percentile(value, 5), percentile(value, 95), percentile(value, 98), min(value), max(value), sum(value), count(value) FROM kbps WHERE time > '2015-01-01T00:00:00-07:00' AND time < '2015-01-30T00:00:00-07:00' AND deliveryservice = 'ds_stats' GROUP BY time(), deliveryservice
=cut

my $expected_q =
	"SELECT mean(value), percentile(value, 5), percentile(value, 95), percentile(value, 98), min(value), max(value), sum(value), count(value) FROM kbps WHERE time > '2015-01-01T00:00:00-07:00' AND
                                          time < '2015-01-30T00:00:00-07:00' AND deliveryservice = 'ds_stats' GROUP BY time(60s), deliveryservice";

$summary_q =~ s/\\n/ /g;
$summary_q =~ s/\s+/ /g;

$expected_q =~ s/\\n//g;
$expected_q =~ s/\s+/ /g;

is( $summary_q, $expected_q, 'Compare Summary queries' );

my $series_q = $iq->series_query();
$series_q =~ s/\\n/ /g;
$series_q =~ s/\s+/ /g;

$expected_q =
	"SELECT mean(value) FROM kbps WHERE time > '2015-01-01T00:00:00-07:00' AND time < '2015-01-30T00:00:00-07:00' AND deliveryservice = 'ds_stats' GROUP BY time(60s) ORDER BY asc";
$expected_q =~ s/\\n/ /g;
$expected_q =~ s/\s+/ /g;

is( $series_q, $expected_q, 'Compare Series queries' );

$iq = Builder::DeliveryServiceStatsQuery->new( { XXX => 'XXX' } );
throws_ok {
	$iq->summary_query()
}
qr/'XXX' is not a valid key constructor key./, 'Check invalid parameter key';

done_testing();
