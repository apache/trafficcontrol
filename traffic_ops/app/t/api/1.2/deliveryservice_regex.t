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
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);

#load core test data
Test::TestHelper->load_core_data($schema);

#NEGATIVE TESTING -- No Privs
ok $t->post_ok( '/api/1.1/user/login', json => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(200),
	'Log into the admin user?';

ok $t->get_ok("/api/1.2/deliveryservices_regexes")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_has( '/response', 'has a response' )->json_is( '/response/0/dsName', 'steering-ds1' )->json_has( '/response/0/regexes/0/type', 'has a regex type' )
	->json_is( '/response/1/dsName', 'steering-ds2' )
	->json_has( '/response/1/regexes', 'has a second regex' )->json_has( '/response/7/regexes/0/type', 'has a second regex type' ), 'Query regexes';

$t->get_ok("/api/1.2/deliveryservices/100/regexes")->status_is(200)->json_is( "/response/0/id", 200 )
	->json_is( "/response/0/pattern" => '.*\.foo\..*' )
	->json_is( "/response/0/type" => 19 )
	->json_is( "/response/0/typeName" => 'HOST_REGEXP' )
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok("/api/1.2/deliveryservices/100/regexes")->status_is(200)->json_is( "/response/1/id", 800 )
	->json_is( "/response/1/pattern" => '.*\.steering-ds1\..*' )
	->json_is( "/response/1/type" => 19 )
	->json_is( "/response/1/typeName" => 'HOST_REGEXP' )
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok("/api/1.2/deliveryservices/100/regexes/200")->status_is(200)->json_is( "/response/0/id", 200 )
	->json_is( "/response/0/pattern" => '.*\.foo\..*' )
	->json_is( "/response/0/type" => 19 )
	->json_is( "/response/0/typeName" => 'HOST_REGEXP' )
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->put_ok('/api/1.2/deliveryservices/100/regexes/200' => {Accept => 'application/json'} => json => {
			"pattern" => '.*\.foo-bar\..*',
			"type" => 20,
			"setNumber" => 22,
		})
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/pattern" => '.*\.foo-bar\..*' )
		->json_is( "/response/type" => 20 )
		->json_is( "/response/typeName" => 'PATH_REGEXP' )
		->json_is( "/response/setNumber" => 22 )
		->json_is( "/alerts/0/level" => "success" )
	, 'Did the delivery service regex update?';

ok $t->delete_ok('/api/1.2/deliveryservices/100/regexes/200')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->delete_ok('/api/1.2/deliveryservices/100/regexes/800')->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/alerts/0/level" => "error" )
		->json_is( "/alerts/0/text" => "A delivery service must have at least one regex." )
	, 'Does the delivery service regex delete fail because each ds must have at least one regex?';;

ok $t->post_ok('/api/1.2/deliveryservices/100/regexes' => {Accept => 'application/json'} => json => {
			"pattern" => "foo.bar.com",
			"type" => 19,
			"setNumber" => 2,
		})
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/pattern" => "foo.bar.com" )
		->json_is( "/response/type" => 19 )
		->json_is( "/response/typeName" => "HOST_REGEXP" )
		->json_is( "/response/setNumber" => 2 )
		->json_is( "/alerts/0/level" => "success" )
		->json_is( "/alerts/0/text" => "Delivery service regex creation was successful." )
	, 'Is the delivery service regex created?';

ok $t->post_ok('/api/1.2/deliveryservices/100/regexes' => {Accept => 'application/json'} => json => {
			"pattern" => "foo2.bar.com",
			"type" => 12,
			"setNumber" => 2,
		})
		->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/alerts/0/level" => "error" )
		->json_is( "/alerts/0/text" => "Invalid regex type" )
	, 'Does the delivery service regex create fail due to bad regex type?';

#prepare for negative test - enable the ds-user tablet
my $useTenancyParamId = &get_param_id('use_tenancy');
ok $t->put_ok('/api/1.2/parameters/' . $useTenancyParamId => {Accept => 'application/json'} => json => {
			'value'      => '0',
		})->status_is(200)
		->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/name" => "use_tenancy" )
		->json_is( "/response/configFile" => "global" )
		->json_is( "/response/value" => "0" )
	, 'Was the disabling paramter set?';

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

#NEGATIVE TESTING -- No Privs
ok $t->post_ok( '/api/1.1/user/login', json => { u => Test::TestHelper::PORTAL_USER, p => Test::TestHelper::PORTAL_USER_PASSWORD } )->status_is(200),
	'Log into the portal user?';

# Verify Permissions
ok $t->get_ok("/api/1.2/deliveryservices_regexes")->status_is(403)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

done_testing();

sub get_param_id {
	my $name = shift;
	my $q      = "select id from parameter where name = \'$name\'";
	my $get_svr = $dbh->prepare($q);
	$get_svr->execute();
	my $p = $get_svr->fetchall_arrayref( {} );
	$get_svr->finish();
	my $id = $p->[0]->{id};
	return $id;
}
