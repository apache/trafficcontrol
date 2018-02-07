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

ok my $portal_user = $schema->resultset('TmUser')->find( { username => Test::TestHelper::PORTAL_USER } ), 'Does the portal user exist?';

# Verify the Portal user
$t->post_ok( '/api/1.1/user/login', json => { u => Test::TestHelper::PORTAL_USER, p => Test::TestHelper::PORTAL_USER_PASSWORD } )->status_is(200);
$t->get_ok('/api/1.1/user/current.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/username", Test::TestHelper::PORTAL_USER );

# Test required fields
$t->post_ok( '/api/1.1/user/current/update',
	json => { user => { username => Test::TestHelper::PORTAL_USER, fullName => 'tom sawyer', email => 'testportal1@kabletown.com', address_line1 => 'newaddress', role => 6 } } )
	->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/alerts/0/text", "User profile was successfully updated" );

$t->post_ok( '/api/1.1/user/current/update',
	json => { user => { username => Test::TestHelper::PORTAL_USER, fullName => 'tom sawyer', email => 'testportal1@kabletown.com', address_line1 => 'newaddress', role => 3 } } )
	->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/alerts/0/text", "role cannot exceed current user's privilege level (15)" );

# Ensure unique emails
ok $t->post_ok( '/api/1.1/user/current/update', json => { user => { fullName => 'tom sawyer', username => Test::TestHelper::PORTAL_USER, email => 'testportal1@kabletown.com', role => 6 } } )
	->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/alerts/0/level", "success" ),
	"Verify that the emails are unique";

ok $t->post_ok( '/api/1.1/user/current/update', json => { user => { fullName => 'tom sawyer', username => Test::TestHelper::PORTAL_USER, email => '@kabletown.com', role => 6 } } )
	->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/alerts/0/level", "error" ),
	"Verify that the emails are properly formatted";

ok $t->post_ok( '/api/1.1/user/current/update', json => { user => { fullName => 'tom sawyer', username => Test::TestHelper::PORTAL_USER, email => '@kabletown.com', role => 6 } } )
	->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/alerts/0/level", "error" ),
	"Verify that the usernames are unique";

$t->post_ok( '/api/1.1/user/current/update', json => { user => { fullName => 'tom sawyer', email => 'testportal1@kabletown.com', "role" => 6 } } )->status_is(400)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/alerts/0/text", "username is required" );

$t->post_ok( '/api/1.1/user/current/update', json => { user => { fullName => 'tom sawyer', username => Test::TestHelper::PORTAL_USER, "role" => 6 } } )->status_is(400)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/alerts/0/text", "email is required" );

ok $t->post_ok('/api/1.1/user/logout')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();

done_testing();
