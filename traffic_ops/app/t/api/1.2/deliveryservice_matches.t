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

BEGIN { $ENV{MOJO_MODE} = "test" }

my $schema = Schema->connect_to_database;
my $t      = Test::Mojo->new('TrafficOps');

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);

#load core test data
Test::TestHelper->load_core_data($schema);

#NEGATIVE TESTING -- No Privs
ok $t->post_ok( '/api/1.1/user/login', json => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(200),
	'Log into the admin user?';

ok $t->get_ok("/api/1.2/deliveryservice_matches.json")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_has( '/response', 'has a response' )->json_is( '/response/0/dsName', 'steering_ds1' )->json_has( '/response/0/patterns', 'has a first match' )
	->json_is( '/response/1/dsName', 'steering_ds2' )->json_has( '/response/1/patterns', 'has a second match' ), 'Query matches';

ok $t->get_ok("/api/1.2/deliveryservice_matches.json?format=file")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_hasnt( '/response', 'should not have a response' )->json_is( '/0/dsName', 'steering_ds1' )->json_has( '/0/patterns', 'has a first match' )
	->json_is( '/1/dsName', 'steering_ds2' )->json_has( '/1/patterns', 'has a second match' ), 'Query matches';

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

#NEGATIVE TESTING -- No Privs
ok $t->post_ok( '/api/1.1/user/login', json => { u => Test::TestHelper::PORTAL_USER, p => Test::TestHelper::PORTAL_USER_PASSWORD } )->status_is(200),
	'Log into the portal user?';

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

done_testing();
