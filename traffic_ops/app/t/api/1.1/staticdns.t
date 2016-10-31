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
no warnings 'once';
use warnings 'all';
use Test::TestHelper;
use Fixtures::Staticdnsentry;

#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

BEGIN { $ENV{MOJO_MODE} = "test" }

my $schema = Schema->connect_to_database;
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');

Test::TestHelper->unload_core_data($schema);
Test::TestHelper->load_core_data($schema);
Test::TestHelper->load_all_fixtures( Fixtures::Staticdnsentry->new( { schema => $schema, no_transactions => 1 } ) );

ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok('/api/1.1/staticdnsentries.json?orderby=host')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/0/ttl", "3600" )->json_is( "/response/0/host", "AAAA_RECORD_HOST" )->json_is( "/response/0/cachegroup", "mid-northeast-group" )
	->json_is( "/response/0/deliveryservice", "test-ds1" )->json_is( "/response/0/address", "127.0.0.1" )->json_is( "/response/0/type", "AAAA_RECORD" )

	->json_is( "/response/2/ttl", "3600" )->json_is( "/response/2/host", "CNAME_HOST" )->json_is( "/response/2/cachegroup", "mid-northwest-group" )
	->json_is( "/response/2/deliveryservice", "test-ds2" )->json_is( "/response/2/address", "127.0.0.1" )->json_is( "/response/2/type", "CNAME_RECORD" );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();
