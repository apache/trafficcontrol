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
	
	Test::TestHelper->unload_core_data($schema);
	Test::TestHelper->teardown( $schema, 'Log' );
	Test::TestHelper->teardown( $schema, 'Role' );
	Test::TestHelper->teardown( $schema, 'TmUser' );

	Test::TestHelper->load_core_data($schema);

	# Verify the user
	ok my $user = $schema->resultset('TmUser')->find( { username => $login_user } ), 'Does the portal user exist?';
	
	ok $t->post_ok( '/login', => form => { u => $login_user, p => $login_password} )->status_is(302);
		
	#adding a user
	my $addedUserName = "user1";
	my $addedUserEmail = "abc\@z.com";

	ok $t->post_ok('/api/1.2/users' => {Accept => 'application/json'} => json => {
        	"username" => $addedUserName, "fullName"=>"full name", "email" => $addedUserEmail, "localPassword" => "pass", "confirmLocalPassword"=> "pass", "role" => 4 })
        	->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/username" =>  $addedUserName )
		->json_is( "/response/email" =>  $addedUserEmail)
        	    , 'Failed adding user?';

	#same name again - fail
	ok $t->post_ok('/api/1.2/users' => {Accept => 'application/json'} => json => {
        	"username" => $addedUserName, "fullName"=>"full name1", "email" => "xy\@z.com", "localPassword" => "pass", "confirmLocalPassword"=> "pass", "role" => 4 })
        	->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		       , 'Success same user...';

	#bad email - fail
	ok $t->post_ok('/api/1.2/users' => {Accept => 'application/json'} => json => {
        	"username" => "user2", "fullName"=>"full name2", "email" => "xy", "localPassword" => "pass", "confirmLocalPassword"=> "pass", "role" => 4 })
        	->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	       , 'Success bad email...';

	#adding same email again - fail
	ok $t->post_ok('/api/1.2/users' => {Accept => 'application/json'} => json => {
        	"username" => "new-user", "fullName"=>"full name3", "email" => $addedUserEmail, "localPassword" => "pass", "confirmLocalPassword"=> "pass", "role`" => 4 })
        	->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	       , 'Success same email...';
	       
	my $userid = $schema->resultset('TmUser')->find( { username => $addedUserName } )->id, 'Does the portal user exist?';
	       
	ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
}

my $schema = Schema->connect_to_database;
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');

run_ut($t, $schema, Test::TestHelper::ADMIN_USER,  Test::TestHelper::ADMIN_USER_PASSWORD);

$dbh->disconnect();
done_testing();


