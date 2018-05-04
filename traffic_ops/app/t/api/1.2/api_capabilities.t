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

#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

BEGIN { $ENV{MOJO_MODE} = "test" }

my $schema = Schema->connect_to_database;
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');

Test::TestHelper->unload_core_data($schema);
Test::TestHelper->load_core_data($schema);

ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';

#clear the tables
my $th = $dbh->prepare("TRUNCATE TABLE capability CASCADE");
$th->execute();

#add capabilities required for the tests (basic-read & cdn-write)
my $description = "Basic read operations";
my $cap_name = "basic-read";
$t->post_ok("/api/1.2/capabilities" => {Accept => 'application/json'} => json => {
		"name" => $cap_name, "description" => $description
	})->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content};} )
	->json_is( "/response/description" => $description )
	->json_is( "/response/name" => $cap_name )
	, 'Does capability details return?';

$description = "CDN write operations";
$cap_name = "cdn-write";
$t->post_ok("/api/1.2/capabilities" => {Accept => 'application/json'} => json => {
		"name" => $cap_name, "description" => $description
	})->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content};} )
	->json_is( "/response/description" => $description )
	->json_is( "/response/name" => $cap_name )
	, 'Does capability details return?';

$t->get_ok("/api/1.2/api_capabilities")->status_is(200)->json_is( "/response", [] )->or( sub { diag $t->tx->res->content->asset->{content}; } );;

# adding valid entry
my $http_method = "GET";
my $http_route = "sample/route";
my $cap_name = "basic-read";
$t->post_ok("/api/1.2/api_capabilities" => {Accept => 'application/json'} => json => {
			"httpMethod" => $http_method, "httpRoute" => $http_route, "capability" => $cap_name
		})->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content};} )
	->json_is( "/response/id" => 1 )
	->json_is( "/response/httpMethod" => $http_method )
	->json_is( "/response/httpRoute" => $http_route )
	->json_is( "/response/capability" => $cap_name )
	, 'Does mapping details return?';

#verifying the create worked
$t->get_ok("/api/1.2/api_capabilities")->status_is(200)
	->json_is( "/response/0/id" => 1 )
	->json_is( "/response/0/httpMethod" => $http_method )
	->json_is( "/response/0/httpRoute" => $http_route )
	->json_is( "/response/0/capability" => $cap_name )
	->or( sub { diag $t->tx->res->content->asset->{content}; } );;

#verifying get single
$t->get_ok("/api/1.2/api_capabilities/1")->status_is(200)
	->json_is( "/response/0/id" => 1 )
	->json_is( "/response/0/httpMethod" => $http_method )
	->json_is( "/response/0/httpRoute" => $http_route )
	->json_is( "/response/0/capability" => $cap_name )
	->or( sub { diag $t->tx->res->content->asset->{content}; } );;

#insert the same mapping twice - fails
$t->post_ok("/api/1.2/api_capabilities" => {Accept => 'application/json'} => json => {
		"httpMethod" => $http_method, "httpRoute" => $http_route, "capability" => $cap_name
	})->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content};} )
	->json_is( "/alerts/0/text" => "HTTP method \'$http_method\', route \'$http_route\' are already mapped to capability: $cap_name" )
	, 'Is same entry twice?';

#edit a mapping
my $cap_name_updated = "cdn-write";
$t->put_ok("/api/1.2/api_capabilities/1" => {Accept => 'application/json'} => json => {
		"httpMethod" => $http_method, "httpRoute" => $http_route, "capability" => $cap_name_updated
	})->status_is(200)
	->json_is( "/response/id" => 1 )
	->json_is( "/response/httpMethod" => $http_method )
	->json_is( "/response/httpRoute" => $http_route )
	->json_is( "/response/capability" => $cap_name_updated )
	, 'Did update succeed?';

#get after update
$t->get_ok("/api/1.2/api_capabilities/1" => {Accept => 'application/json'} )->status_is(200)
	->json_is( "/response/0/id" => 1 )
	->json_is( "/response/0/httpMethod" => $http_method )
	->json_is( "/response/0/httpRoute" => $http_route )
	->json_is( "/response/0/capability" => $cap_name_updated )
	, 'Did get after update succeed?';

#edit the mapping back
$t->put_ok("/api/1.2/api_capabilities/1" => {Accept => 'application/json'} => json => {
		"httpMethod" => $http_method, "httpRoute" => $http_route, "capability" => $cap_name
	})->status_is(200)
	->json_is( "/response/id" => 1 )
	->json_is( "/response/httpMethod" => $http_method )
	->json_is( "/response/httpRoute" => $http_route )
	->json_is( "/response/capability" => $cap_name )
	, 'Did update succeed?';

#get after update
$t->get_ok("/api/1.2/api_capabilities/1" => {Accept => 'application/json'} )->status_is(200)
	->json_is( "/response/0/id" => 1 )
	->json_is( "/response/0/httpMethod" => $http_method )
	->json_is( "/response/0/httpRoute" => $http_route )
	->json_is( "/response/0/capability" => $cap_name )
	, 'Did get after update back succeed?';

#insert another mapping
my $http_method_post = "POST";
my $route_sample2 = "sample/route2";
$t->post_ok("/api/1.2/api_capabilities" => {Accept => 'application/json'} => json => {
		"httpMethod" => $http_method_post, "httpRoute" => $route_sample2, "capability" => $cap_name
	})->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content};} )
	->json_is( "/response/id" => 2 )
	->json_is( "/response/httpMethod" => $http_method_post )
	->json_is( "/response/httpRoute" => $route_sample2 )
	->json_is( "/response/capability" => $cap_name )
	, 'Does mapping details return?';

#get by cap name
$t->get_ok("/api/1.2/api_capabilities?capability=$cap_name")->status_is(200)
	->json_is( "/response/0/id" => 1 )
	->json_is( "/response/0/httpMethod" => $http_method )
	->json_is( "/response/0/httpRoute" => $http_route )
	->json_is( "/response/0/capability" => $cap_name )
	->json_is( "/response/1/id" => 2 )
	->json_is( "/response/1/httpMethod" => $http_method_post )
	->json_is( "/response/1/httpRoute" => $route_sample2 )
	->json_is( "/response/1/capability" => $cap_name )
	->or( sub { diag $t->tx->res->content->asset->{content}; } );;

#test delete
$t->delete_ok("/api/1.2/api_capabilities/2")->status_is(200)
	->json_is( "/alerts/0/text" => "API-capability mapping deleted." )
	, 'Did delete succeed?';

#make sure mapping was deleted
$t->get_ok("/api/1.2/api_capabilities/2")->status_is(200)->json_is( "/response", [] )->or( sub { diag $t->tx->res->content->asset->{content}; } );;

#negative tests

# adding invalid entry - no httpMethod
$t->post_ok("/api/1.2/api_capabilities" => {Accept => 'application/json'} => json => {
		"httpRoute" => $http_route, "capability" => $cap_name
	})->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content};} )
	->json_is( "/alerts/0/text" => "HTTP method is required." )
	, 'Was invalid insert (no httpMethod) reject correctly?';

# adding invalid entry - no route
$t->post_ok("/api/1.2/api_capabilities" => {Accept => 'application/json'} => json => {
		"httpMethod" => $http_method, "capability" => $cap_name
	})->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content};} )
	->json_is( "/alerts/0/text" => "Route is required." )
	, 'Was invalid insert (no route) reject correctly?';

# adding invalid entry - empty route
$t->post_ok("/api/1.2/api_capabilities" => {Accept => 'application/json'} => json => {
		"httpMethod" => $http_method, "capability" => $cap_name, "httpRoute" => ""
	})->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content};} )
	->json_is( "/alerts/0/text" => "Route is required." )
	, 'Was invalid insert (no route) reject correctly?';

# adding invalid entry - no capability
$t->post_ok("/api/1.2/api_capabilities" => {Accept => 'application/json'} => json => {
		"httpMethod" => $http_method, "httpRoute" => $http_route
	})->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content};} )
	->json_is( "/alerts/0/text" => "Capability name is required." )
	, 'Was invalid insert (no capability) reject correctly?';

# adding invalid entry - empty capability
$t->post_ok("/api/1.2/api_capabilities" => {Accept => 'application/json'} => json => {
		"httpMethod" => $http_method, "httpRoute" => $http_route, "capability" => ""
	})->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content};} )
	->json_is( "/alerts/0/text" => "Capability name is required." )
	, 'Was invalid insert (no capability) reject correctly?';

# adding invalid entry - invalid httpMethod
my $invalid_http_method = 'BAD';
$t->post_ok("/api/1.2/api_capabilities" => {Accept => 'application/json'} => json => {
		"httpMethod" => $invalid_http_method, "httpRoute" => $http_route, "capability" => $cap_name
	})->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content};} )
	->json_is( "/alerts/0/text" => "HTTP method \'$invalid_http_method\' is invalid. Valid values are: DELETE, GET, PATCH, POST, PUT" )
	, 'Was invalid insert (no capability) reject correctly?';

# adding invalid entry - non-existing capability
my $non_existing_cap = "non-existing";
$t->post_ok("/api/1.2/api_capabilities" => {Accept => 'application/json'} => json => {
		"httpMethod" => $http_method, "httpRoute" => $http_route, "capability" => $non_existing_cap
	})->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content};} )
	->json_is( "/alerts/0/text" => "Capability \'$non_existing_cap\' does not exist." )
	, 'Was invalid insert (no capability) reject correctly?';


ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

$dbh->disconnect();
done_testing();

