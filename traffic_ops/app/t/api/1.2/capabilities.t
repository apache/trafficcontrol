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

#clear the table
my $th = $dbh->prepare("TRUNCATE TABLE capability CASCADE");
$th->execute();

$t->get_ok("/api/1.2/capabilities")->status_is(200)->json_is( "/response", [] )->or( sub { diag $t->tx->res->content->asset->{content}; } );;

# adding valid entry
my $description = "Basic read operations";
my $cap_name = "basic-read";
$t->post_ok("/api/1.2/capabilities" => {Accept => 'application/json'} => json => {
			"name" => $cap_name, "description" => $description
		})->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content};} )
	->json_is( "/response/description" => $description )
	->json_is( "/response/name" => $cap_name )
	, 'Does capability details return?';

#verifying the create worked
$t->get_ok("/api/1.2/capabilities")->status_is(200)
	->json_is( "/response/0/name" => $cap_name )
	->json_is( "/response/0/description" => $description )
	->or( sub { diag $t->tx->res->content->asset->{content}; } );;

#verifying get by capapbility
$t->get_ok("/api/1.2/capabilities/$cap_name")->status_is(200)
	->json_is( "/response/0/name" => $cap_name )
	->json_is( "/response/0/description" => $description )
	->or( sub { diag $t->tx->res->content->asset->{content}; } );;

#insert the same capability twice - fails
$t->post_ok("/api/1.2/capabilities" => {Accept => 'application/json'} => json => {
		"name" => $cap_name, "description" => $description
	})->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content};} )
	->json_is( "/alerts/0/text" => "Capability \'$cap_name\' already exists." )
	, 'Is same entry twice?';

#edit a capability
my $description_updated = "edited desctiption";
$t->put_ok("/api/1.2/capabilities/$cap_name" => {Accept => 'application/json'} => json => {
		"description" => $description_updated
	})->status_is(200)
	->json_is( "/response/name" => $cap_name )
	->json_is( "/response/description" => $description_updated )
	, 'Did update succeed?';

#get after update
$t->get_ok("/api/1.2/capabilities/$cap_name" => {Accept => 'application/json'} )->status_is(200)
	->json_is( "/response/0/name" => $cap_name )
	->json_is( "/response/0/description" => $description_updated )
	, 'Did get after update succeed?';

#edit the mapping back
$t->put_ok("/api/1.2/capabilities/$cap_name" => {Accept => 'application/json'} => json => {
		"name" => $cap_name, "description" => $description
	})->status_is(200)
	->json_is( "/response/name" => $cap_name )
	->json_is( "/response/description" => $description )
	, 'Did update succeed?';

#get after update
$t->get_ok("/api/1.2/capabilities/$cap_name" => {Accept => 'application/json'} )->status_is(200)
	->json_is( "/response/0/name" => $cap_name )
	->json_is( "/response/0/description" => $description )
	, 'Did get after update back succeed?';

#insert another capability
my $cap_name_basic_write = "basic-write";
my $description_basic_write = "Basic write operations";
$t->post_ok("/api/1.2/capabilities" => {Accept => 'application/json'} => json => {
		"name" => $cap_name_basic_write, "description" => $description_basic_write
	})->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content};} )
	->json_is( "/response/name" => $cap_name_basic_write )
	->json_is( "/response/description" => $description_basic_write )
	, 'Does capability details return?';

#get by cap name
$t->get_ok("/api/1.2/capabilities/$cap_name_basic_write")->status_is(200)
	->json_is( "/response/0/name" => $cap_name_basic_write )
	->json_is( "/response/0/description" => $description_basic_write )
	->or( sub { diag $t->tx->res->content->asset->{content}; } );;

#test delete
$t->delete_ok("/api/1.2/capabilities/$cap_name_basic_write")->status_is(200)
	->json_is( "/alerts/0/text" => "Capability deleted." )
	, 'Did delete succeed?';

#make sure mapping was deleted
$t->get_ok("/api/1.2/capabilities/$cap_name_basic_write")->status_is(200)->json_is( "/response", [] )->or( sub { diag $t->tx->res->content->asset->{content}; } );;

#negative tests
# adding invalid entry - no description
$t->post_ok("/api/1.2/capabilities" => {Accept => 'application/json'} => json => {
		"name" => $cap_name
	})->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content};} )
	->json_is( "/alerts/0/text" => "Description is required." )
	, 'Was invalid insert (no description) reject correctly?';

# adding invalid entry - empty description
$t->post_ok("/api/1.2/capabilities" => {Accept => 'application/json'} => json => {
		"name" => $cap_name, "description" => ""
	})->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content};} )
	->json_is( "/alerts/0/text" => "Description is required." )
	, 'Was invalid insert (no description) reject correctly?';

# adding invalid entry - no name
$t->post_ok("/api/1.2/capabilities" => {Accept => 'application/json'} => json => {
		"description" => $description
	})->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content};} )
	->json_is( "/alerts/0/text" => "Name is required." )
	, 'Was invalid insert (no route) reject correctly?';

# adding invalid entry - empty name
$t->post_ok("/api/1.2/capabilities" => {Accept => 'application/json'} => json => {
		"description" => $description, "name" => ""
	})->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content};} )
	->json_is( "/alerts/0/text" => "Name is required." )
	, 'Was invalid insert (no route) reject correctly?';

# trying to delete a referenced capability. first add a mapping to it.
my $http_method = "GET";
my $http_route = "sample/route";
$t->post_ok("/api/1.2/api_capabilities" => {Accept => 'application/json'} => json => {
		"httpMethod" => $http_method, "httpRoute" => $http_route, "capability" => $cap_name
	})->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content};} )
	->json_is( "/response/id" => 1 )
	->json_is( "/response/httpMethod" => $http_method )
	->json_is( "/response/httpRoute" => $http_route )
	->json_is( "/response/capability" => $cap_name )
	, 'Does mapping details return?';

#test delete -  should fail
$t->delete_ok("/api/1.2/capabilities/$cap_name")->status_is(400)
	->json_is( "/alerts/0/text" => "Capability \'$cap_name\' is refered by an api_capability mapping: 1. Deletion failed." )
	, 'Did delete succeed?';


ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

$dbh->disconnect();
done_testing();

