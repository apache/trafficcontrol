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

#verifying the basic cfg
$t->get_ok("/api/1.2/tenants")->status_is(200)->json_is( "/response/0/name", "root" )->or( sub { diag $t->tx->res->content->asset->{content}; } );;

my $root_tenant_id = &get_tenant_id('root');

ok $t->post_ok('/api/1.2/tenants' => {Accept => 'application/json'} => json => {
        "name" => "tenantA", "parentId" => $root_tenant_id })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "tenantA" )
	->json_is( "/response/parentId" =>  $root_tenant_id)
            , 'Does the tenant details return?';

#same name - would not accept
ok $t->post_ok('/api/1.2/tenants' => {Accept => 'application/json'} => json => {
        "name" => "tenantA", "parentId" => $root_tenant_id })->status_is(400);

#no name - would not accept
ok $t->post_ok('/api/1.2/tenants' => {Accept => 'application/json'} => json => {
        "parentId" => $root_tenant_id })->status_is(400);

#no parent - would not accept
ok $t->post_ok('/api/1.2/tenants' => {Accept => 'application/json'} => json => {
        "name" => "tenantB" })->status_is(400);

#rename
my $tenantA_id = &get_tenant_id('tenantA');
ok $t->put_ok('/api/1.2/tenants/' . $tenantA_id  => {Accept => 'application/json'} => json => {
			"name" => "tenantA2", "parentId" => $root_tenant_id 
		})
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/name" => "tenantA2" )
		->json_is( "/response/id" => $tenantA_id )
		->json_is( "/response/parentId" => $root_tenant_id )
		->json_is( "/alerts/0/level" => "success" )
	, 'Does the tenant2 details return?';

#cannot change tenant parent to undef
ok $t->put_ok('/api/1.2/tenants/' . $tenantA_id  => {Accept => 'application/json'} => json => {
			"name" => "tenantC", 
		})->status_is(400);

#adding a child tenant
ok $t->post_ok('/api/1.2/tenants' => {Accept => 'application/json'} => json => {
        "name" => "tenantD", "parentId" => $tenantA_id })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "tenantD" )
	->json_is( "/response/parentId" =>  $tenantA_id)
            , 'Does the tenant details return?';

#cannot delete a tenant that have children
ok $t->delete_ok('/api/1.2/tenants/' . $tenantA_id)->status_is(500);

my $tenantD_id = &get_tenant_id('tenantD');

ok $t->delete_ok('/api/1.2/tenants/' . $tenantD_id)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->delete_ok('/api/1.2/tenants/' . $tenantA_id)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();

sub get_tenant_id {
	my $name = shift;
	my $q    = "select id from tenant where name = \'$name\'";
	my $get_svr = $dbh->prepare($q);
	$get_svr->execute();
	my $p = $get_svr->fetchall_arrayref( {} );
	$get_svr->finish();
	my $id = $p->[0]->{id};
	return $id;
}

