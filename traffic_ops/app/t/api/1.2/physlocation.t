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

ok $t->post_ok('/api/1.2/regions/Denver Region/phys_locations' => {Accept => 'application/json'} => json => {
        "name" => "physical location1" ,
        "shortName" => "physloc1"})->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "physical location1" )
	->json_is( "/response/shortName" => "physloc1" )
    ->json_is( "/response/regionName" => "Denver Region" )
            , 'Does the physical location details return?';

ok $t->post_ok('/api/1.2/regions/non_region/phys_locations' => {Accept => 'application/json'} => json => {
        "name" => "physical location1",
        "shortName" => "mountain"})->status_is(400);

$t->get_ok("/api/1.2/phys_locations")->status_is(200)->json_is( "/response/0/id", 200 )
	->json_is( "/response/0/name", "Boulder" )->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok("/api/1.2/phys_locations/100")->status_is(200)->json_is( "/response/0/id", 100 )
	->json_is( "/response/0/name", "Denver" )->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok('/api/1.2/phys_locations' => {Accept => 'application/json'} => json => {
			"name" => "phys1",
			"shortName" => "phys1",
			"address" => "address",
			"city" => "city",
			"state" => "state",
			"zip" => "zip",
			"regionId" => "string",
		})
		->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/alerts/0/level" => "error" )
		->json_is( "/alerts/0/text" => "regionId must be a positive integer" )
	, 'Is phys location NOT created?';

ok $t->post_ok('/api/1.2/phys_locations' => {Accept => 'application/json'} => json => {
			"name" => "phys1",
			"shortName" => "phys1",
			"address" => "address",
			"city" => "city",
			"state" => "state",
			"zip" => "zip",
			"regionId" => 100,
		})
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/name" => "phys1" )
		->json_is( "/alerts/0/level" => "success" )
		->json_is( "/alerts/0/text" => "Phys location creation was successful." )
	, 'Is phys location created?';

my $phys_loc_id = &get_phys_location_id('phys1');

ok $t->put_ok('/api/1.2/phys_locations/' . $phys_loc_id  => {Accept => 'application/json'} => json => {
			"name" => "phys2",
			"shortName" => "phys2",
			"address" => "address",
			"city" => "city",
			"state" => "state",
			"zip" => "zip",
			"regionId" => 100,
		})
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/name" => "phys2" )
		->json_is( "/alerts/0/level" => "success" )
	, 'Is the phys location updated?';

ok $t->delete_ok('/api/1.2/phys_locations/' . $phys_loc_id)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();

sub get_phys_location_id {
	my $name = shift;
	my $q    = "select id from phys_location where name = \'$name\'";
	my $get_svr = $dbh->prepare($q);
	$get_svr->execute();
	my $p = $get_svr->fetchall_arrayref( {} );
	$get_svr->finish();
	my $id = $p->[0]->{id};
	return $id;
}

