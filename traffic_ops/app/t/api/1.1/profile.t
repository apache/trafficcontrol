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
use strict;
use warnings;
use Test::TestHelper;

#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

BEGIN { $ENV{MOJO_MODE} = "test" }

my $t      = Test::Mojo->new('TrafficOps');
my $dbh    = Schema->database_handle;
my $schema = Schema->connect_to_database;

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);

#load core test data
Test::TestHelper->load_core_data($schema);

ok $t->post_ok( '/login', => form => { u => Test::TestHelper::PORTAL_USER, p => Test::TestHelper::PORTAL_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok("/api/1.1/profiles.json?orderby=name")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/0/name", "CCR1" )->json_is( "/response/0/description", "ccr description" )
	->json_is( "/response/1/name", "CCR2" )->json_is( "/response/1/description", "ccr description" )
	->json_is( "/response/2/name", "EDGE1" )->json_is( "/response/2/description", "edge description" )
	->json_is( "/response/3/name", "EDGE2" )->json_is( "/response/3/description", "edge description" )
	->json_is( "/response/4/name", "MID1" )->json_is( "/response/4/description", "mid description" )
	->json_is( "/response/5/name", "MISC" )->json_is( "/response/5/description", "misc profile description" );

$t->get_ok("/api/1.1/profiles/trimmed.json")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();
