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
use Test::More;
use Test::Mojo;
use DBI;
use Data::Dumper;
use strict;
use warnings;
use Schema;
use Test::TestHelper;
use Test::MockModule;
use Extensions::TrafficStats::Connection::InfluxDBAdapter;
use JSON;

BEGIN { $ENV{MOJO_MODE} = "test" }

# NOTE:
#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

my $dbh    = Schema->database_handle;
my $schema = Schema->connect_to_database;
my $t      = Test::Mojo->new('TrafficOps');

Test::TestHelper->unload_core_data($schema);
Test::TestHelper->load_core_data($schema);

ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

my $fake_lwp = new Test::MockModule( 'Extensions::TrafficStats::Connection::InfluxDBAdapter', no_auto => 1 );
my $fake_header = HTTP::Headers->new;
$fake_header->header( 'Content-Type' => 'application/json' );    # set

my $fake_answer = {
	response => {
		series => {
			count   => 2,
			columns => [ "time", "mean" ],
			values  => [ [ "2015-05-07T02:00:00Z", 3309856.31666667 ], [ "2015-05-07T02:00:00Z", 3309856.31666667 ], ],
		},
		summary => {
			average               => 1140.232,
			fifthPercentile       => 0,
			ninetyFifthPercentile => 1561.47,
			ninetyFifthPercentile => 1561.47,
			min                   => 619.22,
			max                   => 1561.47,
			total                 => 6841.39,
		},
		parameters => {
			interval              => "60s",
			fifthPercentile       => 0,
			ninetyFifthPercentile => 1561.47,
			ninetyFifthPercentile => 1561.47,
			min                   => 619.22,
			max                   => 1561.47,
			total                 => 6841.39,
		},
	},
};

my $json_response = encode_json($fake_answer);
my $api_version   = "1.2";

my $fake_response = HTTP::Response->new( 200, undef, $fake_header, $json_response );
$fake_lwp->mock( 'query', sub { return $fake_response } );
ok $t->get_ok("/api/$api_version/cdns/usage/overview.json")->status_is(200)->json_has( '/response', 'has a response' ), 'Query1';

#ok $t->get_ok(
#'/api/1.2/cache_stats.json?cdnName=cdn1&=test-ds1&cacheGroupName=us-co-denver&metricType=kbps&startDate=2015-05-06T20:00:00-06:00&endDate=2015-05-06T20:00:00-06:00&interval=60s'
#)->status_is(200)->json_has( '/response', 'has a response' ), 'Query1';

#ok $t->get_ok(
#'/api/1.2/cache_stats.json?cdnName=cdn1&cacheGroupName=us-co-denver&metricType=kbps&startDate=2015-05-06T20:00:00-06:00&endDate=2015-05-06T20:00:00-06:00&interval=60s'
#)->status_is(200)->json_has( '/response', 'has a response' ), 'Query2';

$dbh->disconnect();
done_testing();
