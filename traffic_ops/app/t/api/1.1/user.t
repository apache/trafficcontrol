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
use Data::Dumper;
use strict;
use warnings;
use Schema;
use Test::TestHelper;
use Fixtures::TmUser;
use Fixtures::Deliveryservice;
use Digest::SHA1 qw(sha1_hex);

#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

BEGIN { $ENV{MOJO_MODE} = "test" }

my $dbh    = Schema->database_handle;
my $schema = Schema->connect_to_database;
my $t      = Test::Mojo->new('TrafficOps');

Test::TestHelper->unload_core_data($schema);
Test::TestHelper->teardown( $schema, 'Log' );
Test::TestHelper->teardown( $schema, 'Role' );
Test::TestHelper->teardown( $schema, 'TmUser' );

Test::TestHelper->load_core_data($schema);

ok my $admin_user  = $schema->resultset('TmUser')->find( { username => Test::TestHelper::ADMIN_USER } ),  'Does the admin user exist?';
ok my $portal_user = $schema->resultset('TmUser')->find( { username => Test::TestHelper::PORTAL_USER } ), 'Does the portal user exist?';

# Verify the Admin user
ok $t->post_ok( '/api/1.1/user/login', json => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(200),
	"Can the admin login?";

ok $t->get_ok( '/api/1.1/user/2/deliveryservices/available.json',
	json => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/response/0/id", "2" )->json_is( "/response/0/xmlId", "test-ds2" ),
	"Can the admin get available delivery services";

ok $t->get_ok('/api/1.1/user/current.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/username", Test::TestHelper::ADMIN_USER ), "Verify the admin can reach the current user";

ok $t->post_ok( '/api/1.1/user/current/update', json => { user => { addressLine1 => 'newaddress', email => 'testportal@kabletown.com', role => 4 } } )
	->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/alerts/0/level", "success" ),
	"Verify that an admin can update the user and role";

# Verify changes.
ok $t->get_ok('/api/1.1/user/current.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/username", Test::TestHelper::ADMIN_USER )->json_is( "/response/addressLine1", 'newaddress' )
	->json_is( "/response/email", 'testportal@kabletown.com' )->json_is( "/response/role", 4 ), "Verify the admin updated user and role";

#Test::TestHelper->teardown( $schema, 'TmUser' );
ok $t->post_ok('/api/1.1/user/logout')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# Verify the Portal user
$t->post_ok( '/api/1.1/user/login', json => { u => Test::TestHelper::PORTAL_USER, p => Test::TestHelper::PORTAL_USER_PASSWORD } )->status_is(200);
$t->get_ok('/api/1.1/user/current.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/username", Test::TestHelper::PORTAL_USER );

ok $t->post_ok( '/api/1.1/user/current/update', json => { user => { addressLine1 => 'newaddress', email => 'testportal1@kabletown.com', role => 4 } } )
	->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/alerts/0/level", "success" ),
	"Verify that we can update the user but ignore the role unless an admin";

# Verify changes.
ok $t->get_ok('/api/1.1/user/current.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/username", Test::TestHelper::PORTAL_USER )->json_is( "/response/addressLine1", 'newaddress' )
	->json_is( "/response/email", 'testportal1@kabletown.com' )->json_is( "/response/role", 6 ), "Verify the update happened and the role didn't change";

# Test required fields
$t->post_ok( '/api/1.1/user/current/update', json => { address_line1 => 'newaddress' } )->status_is(400)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/alerts/0/text", "email is required" );

ok $t->post_ok('/api/1.1/user/logout')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();

done_testing();
