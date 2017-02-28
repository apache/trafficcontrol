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
$t->post_ok( '/api/1.2/user/login', json => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(200);
$t->get_ok('/api/1.2/user/current.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/username", Test::TestHelper::ADMIN_USER );


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
        "username" => $addedUserName, "fullName"=>"full name", "email" => "xy\@z.com", "localPassword" => "pass", "confirmLocalPassword"=> "pass", "role" => 4 })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	       , 'Success same user...';

#bad email - fail
ok $t->post_ok('/api/1.2/users' => {Accept => 'application/json'} => json => {
        "username" => "user2", "fullName"=>"full name", "email" => "xy", "localPassword" => "pass", "confirmLocalPassword"=> "pass", "role" => 4 })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	       , 'Success bad email...';

#adding same email again - fail
ok $t->post_ok('/api/1.2/users' => {Accept => 'application/json'} => json => {
        "username" => "new-user", "fullName"=>"full name", "email" => $addedUserEmail, "localPassword" => "pass", "confirmLocalPassword"=> "pass", "role`" => 4 })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	       , 'Success same email...';


ok $t->post_ok('/api/1.2/user/logout')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();

done_testing();
