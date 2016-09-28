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

ok $t->post_ok('/api/1.2/profileparameters/3' => {Accept => 'application/json'} => json => {
	"parametersId" => [4,5] })->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/id" => "3" )
	->json_is( "/response/parametersId/0" => "3" )
	->json_is( "/response/parametersId/1" => "4" )
	->json_is( "/response/parametersId/2" => "5" )
		, 'Does the profile parameter details return?';

ok $t->post_ok('/api/1.2/profileparameters/3' => {Accept => 'application/json'} => json => {
	"parametersId" => [] })->status_is(400);

ok $t->post_ok('/api/1.2/profileparameters/3' => {Accept => 'application/json'} => json => {
	"parametersId" => [2] })->status_is(400);

ok $t->delete_ok('/api/1.2/profileparameters/3/5' => {Accept => 'application/json'})->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/alerts/0/level", "success" )
	->json_is( "/alerts/0/text", "Profile parameter association was deleted." );

my @associated_params = &get_parameter_ids(3);
my @expected = (3,4);
ok( @associated_params ~~ @expected );

ok $t->delete_ok('/api/1.2/profileparameters/3/5' => {Accept => 'application/json'})->status_is(400);

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();

sub get_parameter_ids {
    my $profile_id = shift;
    my $q      = "select * from profile_parameter where profile = \'$profile_id\'";
    my $get_svr = $dbh->prepare($q);
    $get_svr->execute();
    my $p = $get_svr->fetchall_arrayref( {} );
    $get_svr->finish();
    my @ids;
    foreach my $id (@$p) {
        push(@ids, $id->{parameter});
    }
    return @ids;
}
