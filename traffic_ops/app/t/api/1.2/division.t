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

$t->get_ok("/api/1.2/divisions")->status_is(200)->json_is( "/response/0/id", 100 )
	->json_is( "/response/0/name", "mountain" )->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok("/api/1.2/divisions/100")->status_is(200)->json_is( "/response/0/id", 100 )
	->json_is( "/response/0/name", "mountain" )->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok('/api/1.2/divisions' => {Accept => 'application/json'} => json => {
        "name" => "division1" })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "division1" )
            , 'Does the division details return?';

ok $t->post_ok('/api/1.2/divisions' => {Accept => 'application/json'} => json => {
        "name" => "division1" })->status_is(400);

my $division_id = &get_division_id('division1');
ok $t->put_ok('/api/1.2/divisions/' . $division_id  => {Accept => 'application/json'} => json => {
			"name" => "division2"
		})
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/name" => "division2" )
		->json_is( "/alerts/0/level" => "success" )
	, 'Does the division2 details return?';

ok $t->delete_ok('/api/1.2/divisions/' . $division_id)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();

sub get_division_id {
	my $name = shift;
	my $q    = "select id from division where name = \'$name\'";
	my $get_svr = $dbh->prepare($q);
	$get_svr->execute();
	my $p = $get_svr->fetchall_arrayref( {} );
	$get_svr->finish();
	my $id = $p->[0]->{id};
	return $id;
}

