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


sub run_ut {
	my $t = shift;
	my $schema = shift;
	my $login_user = shift;
	my $login_password = shift;

	Test::TestHelper->unload_core_data($schema);
	Test::TestHelper->teardown( $schema, 'Log' );
	Test::TestHelper->teardown( $schema, 'Role' );
	Test::TestHelper->teardown( $schema, 'TmUser' );

	Test::TestHelper->load_core_data($schema);
	
	my $tenant_id = $schema->resultset('TmUser')->find( { username => $login_user } )->get_column('tenant_id');
	my $tenant_name = defined ($tenant_id) ? $schema->resultset('Tenant')->find( { id => $tenant_id } )->get_column('name') : undef;
	my $tenant = defined ($tenant_name) ? $tenant_name : "null";

	ok my $portal_user = $schema->resultset('TmUser')->find( { username => $login_user } ), 'Tenant $tenant_name: Does the portal user exist?';


	# Verify the Portal user
	$t->post_ok( '/api/1.2/user/login', json => { u => $login_user, p => $login_password} )->status_is(200);
	$t->get_ok('/api/1.2/user/current.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/username", $login_user )
		->json_is( "/response/tenantId", $tenant_id)
		->json_is( "/response/tenant",   $tenant_name);

	# Test required fields
	$t->post_ok( '/api/1.2/user/current/update',
		json => { user => { username => $login_user, fullName => 'tom sawyer', email => 'testportal1@kabletown.com', address_line1 => 'newaddress', tenantId => $tenant_id, role => 6 } } )
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/alerts/0/text", "User profile was successfully updated" );

	$t->post_ok( '/api/1.2/user/current/update',
		json => { user => { username => $login_user, fullName => 'tom sawyer', email => 'testportal1@kabletown.com', address_line1 => 'newaddress', tenantId => $tenant_id, role => 3 } } )
		->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/alerts/0/text", "role cannot exceed current user's privilege level (15)" );

	#verify tenancy	
	$t->get_ok('/api/1.2/user/current.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/username", $login_user )
		->json_is( "/response/tenantId", $tenant_id)
		->json_is( "/response/tenant",   $tenant_name);

	# Test required fields
	if (defined($tenant_id)){
		#verify the update with no "tenant" do not removed the tenant
		$t->post_ok( '/api/1.2/user/current/update',
			json => { user => { username => $login_user, fullName => 'tom sawyer', email => 'testportal1@kabletown.com', address_line1 => 'newaddress', role => 6 } } )
			->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
			->json_is( "/alerts/0/text", "User profile was successfully updated" );
		#verify tenancy	
		$t->get_ok('/api/1.2/user/current.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
			->json_is( "/response/username", $login_user )
			->json_is( "/response/tenantId", $tenant_id)
			->json_is( "/response/tenant",   $tenant_name);

		#cannot removed the tenant on current user
		$t->post_ok( '/api/1.2/user/current/update',
			json => { user => { username => $login_user, fullName => 'tom sawyer', email => 'testportal1@kabletown.com', address_line1 => 'newaddress', tenantId => undef, role => 6 } } )
			->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
			->json_is( "/alerts/0/text", "User profile was successfully updated" );
		#verify tenancy	
		$t->get_ok('/api/1.2/user/current.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
			->json_is( "/response/username", $login_user )
			->json_is( "/response/tenantId", $tenant_id)
			->json_is( "/response/tenant",   $tenant_name);
	
		#putting the tenant back the update with no "tenant" removed the tenant
		$t->post_ok( '/api/1.2/user/current/update',
			json => { user => { username => $login_user, fullName => 'tom sawyer', email => 'testportal1@kabletown.com', address_line1 => 'newaddress', tenantId => $tenant_id, role => 6 } } )
			->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
			->json_is( "/alerts/0/text", "User profile was successfully updated" );
		#verify tenancy	
		$t->get_ok('/api/1.2/user/current.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
			->json_is( "/response/username", $login_user )
			->json_is( "/response/tenantId", $tenant_id);
	}
	
	# Ensure unique emails
	ok $t->post_ok( '/api/1.2/user/current/update', json => { user => { username => $login_user, fullName => 'tom sawyer', email => 'testportal1@kabletown.com', tenantId => $tenant_id, role => 6 } } )
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/alerts/0/level", "success" ),
		"Tenant $tenant: Verify that the emails are unique";

	ok $t->post_ok( '/api/1.2/user/current/update', json => { user => { username => $login_user, fullName => 'tom sawyer', email => '@kabletown.com', tenantId => $tenant_id, role => 6 } } )
		->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/alerts/0/level", "error" ),
		"Tenant $tenant: Verify that the emails are properly formatted";

	ok $t->post_ok( '/api/1.2/user/current/update', json => { user => { username => $login_user, fullName => 'tom sawyer', email => '@kabletown.com', tenantId => $tenant_id, role => 6 } } )
		->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/alerts/0/level", "error" ),
		"Tenant $tenant: Verify that the usernames are unique";

	$t->post_ok( '/api/1.2/user/current/update', json => { user => { fullName => 'tom sawyer', email => 'testportal1@kabletown.com', tenantId => $tenant_id, role => 6 } } )->status_is(400)
		->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/alerts/0/text", "username is required" );

	$t->post_ok( '/api/1.2/user/current/update', json => { user => { fullName => 'tom sawyer', username => $login_user, tenantId => $tenant_id, role => 6 } } )->status_is(400)
		->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/alerts/0/text", "email is required" );

	$t->post_ok( '/api/1.2/user/current/update', json => { user => { email => 'testportal1@kabletown.com', username => $login_user, tenantId => $tenant_id, role => 6 } } )->status_is(400)
		->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/alerts/0/text", "fullName is required" );

	ok $t->post_ok('/api/1.2/user/logout')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
}

my $dbh    = Schema->database_handle;
my $schema = Schema->connect_to_database;
my $t      = Test::Mojo->new('TrafficOps');

run_ut($t, $schema, Test::TestHelper::PORTAL_USER,  Test::TestHelper::PORTAL_USER_PASSWORD);
run_ut($t, $schema, Test::TestHelper::PORTAL_ROOT_USER,  Test::TestHelper::PORTAL_ROOT_USER_PASSWORD);

$dbh->disconnect();

done_testing();
