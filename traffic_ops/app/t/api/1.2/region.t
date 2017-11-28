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

ok $t->post_ok('/api/1.2/divisions/mountain/regions' => {Accept => 'application/json'} => json => {
        "name" => "region1"})->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "region1" )
	->json_is( "/response/divisionName" => "mountain" )
            , 'Does the region details return?';

ok $t->post_ok('/api/1.2/divisions/mountain/regions' => {Accept => 'application/json'} => json => {
        "name" => "region1"})->status_is(400);

$t->get_ok("/api/1.2/regions")->status_is(200)->json_is( "/response/0/id", 200 )
	->json_is( "/response/0/name", "Boulder Region" )->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok("/api/1.2/regions/100")->status_is(200)->json_is( "/response/0/id", 100 )
	->json_is( "/response/0/name", "Denver Region" )->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok('/api/1.2/regions' => {Accept => 'application/json'} => json => {
			"name" => "reg1",
			"division" => "string"
		})
		->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/alerts/0/level" => "error" )
		->json_is( "/alerts/0/text" => "division must be a positive integer" )
	, 'Did the region create fail?';

ok $t->post_ok('/api/1.2/regions' => {Accept => 'application/json'} => json => {
			"name" => "reg1",
			"division" => 100
		})
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/name" => "reg1" )
		->json_is( "/response/division" => 100 )
		->json_is( "/alerts/0/level" => "success" )
		->json_is( "/alerts/0/text" => "Region create was successful." )
	, 'Do the region details return?';

my $region_id = &get_reg_id('reg1');

ok $t->put_ok('/api/1.2/regions/' . $region_id  => {Accept => 'application/json'} => json => {
			"name" => "reg2",
			"division" => 100
		})
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/name" => "reg2" )
		->json_is( "/response/division" => 100 )
		->json_is( "/alerts/0/level" => "success" )
	, 'Do the regions details return?';

ok $t->delete_ok('/api/1.2/regions/' . $region_id)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
done_testing();

sub get_reg_id {
	my $name = shift;
	my $q    = "select id from region where name = \'$name\'";
	my $get_svr = $dbh->prepare($q);
	$get_svr->execute();
	my $p = $get_svr->fetchall_arrayref( {} );
	$get_svr->finish();
	my $id = $p->[0]->{id};
	return $id;
}

