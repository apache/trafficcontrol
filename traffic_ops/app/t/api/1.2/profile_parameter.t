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

ok $t->post_ok('/api/1.2/profileparameters' => {Accept => 'application/json'} => json => {
	"profileId" => 3, "parameterId" => 4 })->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/0/profileId" => "3" )
	->json_is( "/response/0/parameterId" => "4" )
		, 'Does the profile parameter details return?';

ok $t->post_ok('/api/1.2/profileparameters' => {Accept => 'application/json'} => json => [
	{ "profileId" => 3, "parameterId" => 5 },
	{ "profileId" => 3, "parameterId" => 6 }
	])->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/0/profileId" => "3" )
	->json_is( "/response/0/parameterId" => "5" )
	->json_is( "/response/1/profileId" => "3" )
	->json_is( "/response/1/parameterId" => "6" )
		, 'Does the profile parameter details return?';

ok $t->post_ok('/api/1.2/profileparameters' => {Accept => 'application/json'} => json => [])->status_is(400);

ok $t->post_ok('/api/1.2/profileparameters' => {Accept => 'application/json'} => json => {
	"profileId" => 3, "parameterId" => 4 })->status_is(400);

ok $t->post_ok('/api/1.2/profileparameters' => {Accept => 'application/json'} => json => {
	"profileId" => 3, "parameterId" => 2 })->status_is(400);

ok $t->delete_ok('/api/1.2/profileparameters/3/5' => {Accept => 'application/json'})->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/alerts/0/level", "success" )
	->json_is( "/alerts/0/text", "Profile parameter association was deleted." );

my @associated_params = &get_parameter_ids(3);
my @expected = (3,4,6);
ok( @associated_params ~~ @expected );

ok $t->delete_ok('/api/1.2/profileparameters/3/5' => {Accept => 'application/json'})->status_is(400);

ok $t->post_ok('/api/1.2/profiles/name/CCR1/parameters' => {Accept => 'application/json'} => json => 
        [
            {
                "name"          => "param1",
                "configFile"    => "configFile1",
                "value"         => "value1"
            },
            {
                "name"          => "param2",
                "configFile"    => "configFile2",
                "value"         => "value2"
            },
        ] ) ->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/profileName" => "CCR1" )
	->json_is( "/response/parameters/0/name" => "param1" )
	->json_is( "/response/parameters/0/configFile" => "configFile1" )
	->json_is( "/response/parameters/0/value" => "value1" )
	->json_is( "/response/parameters/0/secure" => "0" )
	->json_is( "/response/parameters/1/name" => "param2" )
	->json_is( "/response/parameters/1/configFile" => "configFile2" )
	->json_is( "/response/parameters/1/value" => "value2" )
	->json_is( "/response/parameters/1/secure" => "0" )
		, 'Does the profile_parameters create details return?';

ok $t->post_ok('/api/1.2/profiles/name/CCR1/parameters' => {Accept => 'application/json'} => json => 
        [
            {
                "name"          => "param1",
                "configFile"    => "configFile1",
                "value"         => "value1",
                "secure"        => "0"
            },
            {
                "name"          => "param3",
                "configFile"    => "configFile3",
                "value"         => "value3",
                "secure"        => "0"
            },
        ] ) ->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
		, 'Does the profile_parameters create details return?';

my $prof_id = &get_profile_id("CCR1");

ok $t->post_ok('/api/1.2/profiles/'. $prof_id .'/parameters' => {Accept => 'application/json'} => json =>
        [
            {
                "name"          => "param11",
                "configFile"    => "configFile11",
                "value"         => "value11"
            },
            {
                "name"          => "param21",
                "configFile"    => "configFile21",
                "value"         => "value21"
            },
        ] ) ->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/profileName" => "CCR1" )
	->json_is( "/response/parameters/0/name" => "param11" )
	->json_is( "/response/parameters/0/configFile" => "configFile11" )
	->json_is( "/response/parameters/0/value" => "value11" )
	->json_is( "/response/parameters/0/secure" => "0" )
	->json_is( "/response/parameters/1/name" => "param21" )
	->json_is( "/response/parameters/1/configFile" => "configFile21" )
	->json_is( "/response/parameters/1/value" => "value21" )
	->json_is( "/response/parameters/1/secure" => "0" )
		, 'Does the profile_parameters create details return?';

ok $t->post_ok('/api/1.2/profiles/'. $prof_id . '/parameters' => {Accept => 'application/json'} => json =>
        [
            {
                "name"          => "param11",
                "configFile"    => "configFile11",
                "value"         => "value11",
                "secure"        => "0"
            },
            {
                "name"          => "param31",
                "configFile"    => "configFile31",
                "value"         => "value31",
                "secure"        => "0"
            },
        ] ) ->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
		, 'Does the profile_parameters create details return?';

ok $t->post_ok('/api/1.2/profiles/name/CCR1/parameters' => {Accept => 'application/json'} => json => 
        [
            {
                "configFile"    => "configFile1",
                "value"         => "value1",
                "secure"        => "0"
            },
        ] ) ->status_is(400)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_like( "/alerts/0/text" => qr/^there is parameter name does not provide/ )
		, 'Does the profile_parameters create details return?';

ok $t->post_ok('/api/1.2/profiles/name/CCR1/parameters' => {Accept => 'application/json'} => json => {
                "name"          => "param1",
                "value"         => "value1",
                "secure"        => "0"
        }) ->status_is(400)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_like( "/alerts/0/text" => qr/^there is parameter configFile does not provide/ )
		, 'Does the profile_parameters create details return?';

ok $t->post_ok('/api/1.2/profiles/'. $prof_id . '/parameters' => {Accept => 'application/json'} => json => {
                "name"          => "param1",
                "value"         => "value1",
                "secure"        => "0"
        }) ->status_is(400)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_like( "/alerts/0/text" => qr/^there is parameter configFile does not provide/ )
		, 'Does the profile_parameters create details return?';

ok $t->post_ok('/api/1.2/profiles/name/CCR1/parameters' => {Accept => 'application/json'} => json => {
                "name"          => "param1",
                "configFile"    => "configFile1",
                "secure"        => "0"
        }) ->status_is(400)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_like( "/alerts/0/text" => qr/^there is parameter value does not provide/ )
		, 'Does the profile_parameters create details return?';
ok $t->post_ok('/api/1.2/profiles/name/CCR1/parameters' => {Accept => 'application/json'} => json => {
                "name"          => "param1",
                "configFile"    => "configFile1",
                "value"         => "value1",
                "secure"        => "abc"
        }) ->status_is(400)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_like( "/alerts/0/text" => qr/^secure must 0 or 1/ )
		, 'Does the profile_parameters create details return?';

ok $t->post_ok('/api/1.2/profiles/name/CCR11/parameters' => {Accept => 'application/json'} => json => {
                "name"          => "param1",
                "configFile"    => "configFile1",
                "value"         => "value1"
        }) ->status_is(404)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
		, 'Does the profile_parameters create details return?';

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();

sub get_profile_id {
    my $profile_name = shift;
    my $q      = "select * from profile where name = \'$profile_name\'";
    my $get_svr = $dbh->prepare($q);
    $get_svr->execute();
    my $p = $get_svr->fetchall_arrayref( {} );
    $get_svr->finish();
    my $id = $p->[0]->{id};
    return $id;
}

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
