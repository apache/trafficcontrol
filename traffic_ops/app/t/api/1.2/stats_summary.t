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
use Schema;
use Test::TestHelper;
use Test::MockModule;
use Test::MockObject;
use strict;
use warnings;
use JSON;
use Fixtures::StatsSummary;

BEGIN { $ENV{MOJO_MODE} = "test" }

my $schema = Schema->connect_to_database;
my $t      = Test::Mojo->new('TrafficOps');

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);
Test::TestHelper->teardown( $schema, "StatsSummary" );

#load core test data
Test::TestHelper->load_core_data($schema);

my $schema_values = { schema => $schema, no_transactions => 1 };
my $stats_summary = Fixtures::StatsSummary->new($schema_values);
Test::TestHelper->load_all_fixtures($stats_summary);


my $cdn      = "test-cdn1";
my $deliveryservice  = "test-ds1";
my $stat_name  = "test_stat";
my $stat_value    = "3.1415";

my ($sec,$min,$hour,$mday,$mon,$year,$wday,$yday,$isdst)=localtime(time);
my $summary_time = sprintf ( "%04d-%02d-%02d %02d:%02d:%02d",
                                   $year+1900,$mon+1,$mday,$hour,$min,$sec);
my $stat_date = sprintf ( "%04d-%02d-%02d",
                                   $year+1900,$mon+1,$mday);

#login
ok $t->post_ok( '/api/1.1/user/login', json => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(200),
	'Log into the portal user?';

#create
ok $t->post_ok(
	'/api/1.2/stats_summary/create',
	json => {
		cdnName => $cdn,
		deliveryServiceName => $deliveryservice,
		statName => $stat_name,
		statValue => $stat_value,
		summaryTime => $summary_time,
		statDate => $stat_date,
	}
	)->status_is(200)->json_has("Successfully added stats summary record")
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

#get_object
ok $t->get_ok("/api/1.2/stats_summary.json")->status_is(200)
	->json_has($stat_name)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# logout
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

done_testing();
