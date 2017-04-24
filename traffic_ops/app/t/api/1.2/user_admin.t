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
use Data::Dumper;

#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

BEGIN { $ENV{MOJO_MODE} = "test" }

sub run_ut {
	my $t = shift;
	my $schema = shift;
	my $login_user = shift;
	my $login_password = shift;
	my $use_tenancy = shift;

	Test::TestHelper->unload_core_data($schema);
	Test::TestHelper->teardown( $schema, 'Log' );
	Test::TestHelper->teardown( $schema, 'Role' );
	Test::TestHelper->teardown( $schema, 'TmUser' );

	Test::TestHelper->load_core_data($schema);

	my $tenant_id = $schema->resultset('TmUser')->find( { username => $login_user } )->get_column('tenant_id');
	my $tenant_name = defined ($tenant_id) ? $schema->resultset('Tenant')->find( { id => $tenant_id } )->get_column('name') : "null";

	# Verify the user
	ok my $user = $schema->resultset('TmUser')->find( { username => $login_user } ), 'Does the user exist?';
	
	ok $t->post_ok( '/login', => form => { u => $login_user, p => $login_password} )->status_is(302);

	my $useTenancyParamId = &get_param_id('use_tenancy');
	ok $t->put_ok('/api/1.2/parameters/' . $useTenancyParamId => {Accept => 'application/json'} => json => {
				'value'      => $use_tenancy,
			})->status_is(200)
			->or( sub { diag $t->tx->res->content->asset->{content}; } )
			->json_is( "/response/name" => "use_tenancy" )
			->json_is( "/response/configFile" => "global" )
			->json_is( "/response/value" => $use_tenancy )
		, 'Was the disabling paramter set?';

	#adding a user
	my $addedUserName = "user1";
	my $addedUserEmail = "abc\@z.com";

	ok $t->post_ok('/api/1.2/users' => {Accept => 'application/json'} => json => {
        	"username" => $addedUserName, "fullName"=>"full name", "email" => $addedUserEmail, "localPasswd" => "longerpass", "confirmLocalPasswd"=> "longerpass", "role" => 4 , "tenantId" => $tenant_id})
        	->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/username" =>  $addedUserName )
		->json_is( "/response/email" =>  $addedUserEmail)
		->json_is( "/response/tenantId" =>  $tenant_id)
        	    , 'Success added user?';

	#same name again - fail
	ok $t->post_ok('/api/1.2/users' => {Accept => 'application/json'} => json => {
        	"username" => $addedUserName, "fullName"=>"full name1", "email" => "xy\@z.com", "localPasswd" => "longerpass", "confirmLocalPasswd"=> "longerpass", "role" => 4 })
        	->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		       , 'Success same user...';

	#bad email - fail
	ok $t->post_ok('/api/1.2/users' => {Accept => 'application/json'} => json => {
        	"username" => "user2", "fullName"=>"full name2", "email" => "xy", "localPassword" => "longerpass", "confirmLocalPasswd"=> "longerpass", "role" => 4 })
        	->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	       , 'Success bad email...';

	#adding same email again - fail
	ok $t->post_ok('/api/1.2/users' => {Accept => 'application/json'} => json => {
        	"username" => "new-user", "fullName"=>"full name3", "email" => $addedUserEmail, "localPasswd" => "longerpass", "confirmLocalPasswd"=> "longerpass", "role" => 4 })
        	->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	       , 'Success same email...';
	       
	my $userid = $schema->resultset('TmUser')->find( { username => $addedUserName } )->id, 'Does the user exist?';

	#login as the user, and do something, to verify the user can log in with the given password
	ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
	ok $t->post_ok( '/login', => form => { u => $addedUserName, p => "longerpass"} )->status_is(302);
	ok $t->get_ok('/api/1.2/users/'.$userid)->status_is(200);
	#back to the standard user
	ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
	ok $t->post_ok( '/login', => form => { u => $login_user, p => $login_password} )->status_is(302);
	       
	if (defined($tenant_id) and !$use_tenancy){
		#verify the update with no "tenant" removed the tenant
		$t->put_ok( '/api/1.2/users/'.$userid,
			json => { "username" => $addedUserName."1", "fullName"=>"full name", "email" => $addedUserEmail."1", "localPasswd" => "longerpass", "confirmLocalPasswd"=> "longerpass", "role" => 4} )
			->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
			->json_is( "/alerts/0/text", "User update was successful." )
			->json_is( "/response/tenantId", undef);
			
		#putting the tenant back the tenant
		$t->put_ok( '/api/1.2/users/'.$userid,
			json => { "username" => $addedUserName."2", "tenantId" => $tenant_id, "fullName"=>"full name", "email" => $addedUserEmail."2", "localPasswd" => "longerpass", "confirmLocalPasswd"=> "longerpass", "role" => 4} )
			->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
			->json_is( "/response/tenantId", $tenant_id)
			->json_is( "/alerts/0/text", "User update was successful." );
		
	
		#removed the tenant explicitly
		$t->put_ok( '/api/1.2/users/'.$userid,
	 	json => {  "username" => $addedUserName."3", "tenantId" => undef, "fullName"=>"full name", "email" => $addedUserEmail."3", "localPasswd" => "longerpass", "confirmLocalPasswd"=> "longerpass", "role" => 4} )
			->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
			->json_is( "/alerts/0/text", "User update was successful." )
			->json_is( "/response/tenantId", undef);
			
	}
	       
	ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
}

my $schema = Schema->connect_to_database;
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');

run_ut($t, $schema, Test::TestHelper::ADMIN_USER,  Test::TestHelper::ADMIN_USER_PASSWORD, 0);
run_ut($t, $schema, Test::TestHelper::ADMIN_ROOT_USER,  Test::TestHelper::ADMIN_ROOT_USER_PASSWORD, 0);
run_ut($t, $schema, Test::TestHelper::ADMIN_ROOT_USER,  Test::TestHelper::ADMIN_ROOT_USER_PASSWORD, 1);

$dbh->disconnect();
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

