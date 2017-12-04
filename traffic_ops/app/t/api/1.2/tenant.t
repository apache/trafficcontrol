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
use JSON;
use strict;
use warnings;
no warnings 'once';
use warnings 'all';
use Test::TestHelper;
use Utils::Tenant;

#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

BEGIN { $ENV{MOJO_MODE} = "test" }

my $schema = Schema->connect_to_database;
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');

Test::TestHelper->unload_core_data($schema);
Test::TestHelper->load_core_data($schema);

ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_ROOT_USER, p => Test::TestHelper::ADMIN_ROOT_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';


# Count the 'response number'
my $responses_counter = sub {
	my $t = shift;
	my $json = decode_json( $t->tx->res->content->asset->slurp );
	my $r    = $json->{response};
	if ($r) {
		return scalar(@$r);
	}
	return 0;
};

# Count the 'response number', and compare to the give value
my $count_response_test = sub {
	my ( $t, $count ) = @_;
	return $t->success( is( $t->$responses_counter(), $count ) )->or( sub { diag $t->tx->res->content->asset->{content}; } );
};

#verifying the basic cfg
ok $t->get_ok("/api/1.2/tenants")->status_is(200)->json_is( "/response/0/name", "root" )->or( sub { diag $t->tx->res->content->asset->{content}; } );;

my $root_tenant_id = &get_tenant_id('root');

#setting with no "active" field which is optional
ok $t->post_ok('/api/1.2/tenants' => {Accept => 'application/json'} => json => {
        "name" => "tenantA", "active" => 1, "parentId" => $root_tenant_id })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "tenantA" )
	->json_is( "/response/active" =>  1)
	->json_is( "/response/parentId" =>  $root_tenant_id)
            , 'Does the tenant details return?';

#same name - would not accept
ok $t->post_ok('/api/1.2/tenants' => {Accept => 'application/json'} => json => {
        "name" => "tenantA", "active" => 1, "parentId" => $root_tenant_id })->status_is(400);

#no name - would not accept
ok $t->post_ok('/api/1.2/tenants' => {Accept => 'application/json'} => json => {
        "parentId" => $root_tenant_id })->status_is(400);

#no parent - would not accept
ok $t->post_ok('/api/1.2/tenants' => {Accept => 'application/json'} => json => {
        "name" => "tenant" })->status_is(400);

#now getting it excepted
ok $t->post_ok('/api/1.2/tenants' => {Accept => 'application/json'} => json => {
        "name" => "tenant", "active" => 1, "parentId" => $root_tenant_id })->status_is(200);

my $tenantA_id = &get_tenant_id('tenantA');
my $tenantB_id = &get_tenant_id('tenant');
#rename, and move to active
ok $t->put_ok('/api/1.2/tenants/' . $tenantA_id  => {Accept => 'application/json'} => json => {
			"name" => "tenantA2", "active" => 1, "parentId" => $root_tenant_id 
		})
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/name" => "tenantA2" )
		->json_is( "/response/id" => $tenantA_id )
		->json_is( "/response/active" => 1 )
		->json_is( "/response/parentId" => $root_tenant_id )
		->json_is( "/alerts/0/level" => "success" )
	, 'Does the tenantA2 details return?';

#change "active"
ok $t->put_ok('/api/1.2/tenants/' . $tenantA_id  => {Accept => 'application/json'} => json => {
			"name" => "tenantA2", "active" => 0, "parentId" => $root_tenant_id 
		})
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/name" => "tenantA2" )
		->json_is( "/response/id" => $tenantA_id )
		->json_is( "/response/active" => 0 )
		->json_is( "/response/parentId" => $root_tenant_id )
		->json_is( "/alerts/0/level" => "success" )
	, 'Did we moved to non active?';

#change "active" back
ok $t->put_ok('/api/1.2/tenants/' . $tenantA_id  => {Accept => 'application/json'} => json => {
			"name" => "tenantA2", "active" => 1, "parentId" => $root_tenant_id 
		})
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/name" => "tenantA2" )
		->json_is( "/response/id" => $tenantA_id )
		->json_is( "/response/active" => 1 )
		->json_is( "/response/parentId" => $root_tenant_id )
		->json_is( "/alerts/0/level" => "success" )
	, 'Did we moved back to active?';

#cannot change tenant parent to undef
ok $t->put_ok('/api/1.2/tenants/' . $tenantA_id  => {Accept => 'application/json'} => json => {
			"name" => "tenantC", "active" => 1})
			->json_is( "/alerts/0/text" => "parentId is required")
			->status_is(400);

#cannot skip "active" field on "put"
ok $t->put_ok('/api/1.2/tenants/' . $tenantA_id  => {Accept => 'application/json'} => json => {
			"name" => "tenantC", "parentId" => $root_tenant_id})
			->json_is( "/alerts/0/text" => "active is required")
			->status_is(400);

#cannot skip "name" field on "put"
ok $t->put_ok('/api/1.2/tenants/' . $tenantA_id  => {Accept => 'application/json'} => json => {
			"active" => 1, "parentId" => $root_tenant_id})
			->json_is( "/alerts/0/text" => "name is required")
			->status_is(400);

#cannot change root-tenant to inactive
ok $t->put_ok('/api/1.2/tenants/' . $root_tenant_id  => {Accept => 'application/json'} => json => {
			"name" => "root", "active" => 0, "parentId" => undef})
			->json_is( "/alerts/0/text" => "parentId is required")
			->status_is(400);

#adding a child tenant
ok $t->post_ok('/api/1.2/tenants' => {Accept => 'application/json'} => json => {
        "name" => "tenantD", "active" => 1, "parentId" => $tenantA_id })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "tenantD" )
	->json_is( "/response/active" => 1 )
	->json_is( "/response/parentId" =>  $tenantA_id)
            , 'Does the tenant details return?';

#adding a child inactive tenant
ok $t->post_ok('/api/1.2/tenants' => {Accept => 'application/json'} => json => {
        "name" => "tenantE", "active" => 0, "parentId" => $tenantA_id })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "tenantE" )
	->json_is( "/response/active" => 0 )
	->json_is( "/response/parentId" =>  $tenantA_id)
            , 'Does the tenant details return?';

my $tenantD_id = &get_tenant_id('tenantD');
my $tenantE_id = &get_tenant_id('tenantE');

#list tenants- verify heirachic order - order by id
ok $t->get_ok("/api/1.2/tenants?orderby=id")->status_is(200)
	->json_is( "/response/0/id", $root_tenant_id )
	->json_is( "/response/1/id", $tenantA_id)
	->json_is( "/response/2/id", $tenantD_id)
	->json_is( "/response/3/id", $tenantE_id)
	->json_is( "/response/4/id", $tenantB_id)->or( sub { diag $t->tx->res->content->asset->{content}; } );;

#list tenants- verify heirachic order - order by name
ok $t->get_ok("/api/1.2/tenants?orderby=name")->status_is(200)
	->json_is( "/response/0/id", $root_tenant_id )
	->json_is( "/response/2/id", $tenantA_id)
	->json_is( "/response/3/id", $tenantD_id)
	->json_is( "/response/4/id", $tenantE_id)
	->json_is( "/response/1/id", $tenantB_id)->or( sub { diag $t->tx->res->content->asset->{content}; } );;

#list tenants- verify heirachic order - order by name (default)
ok $t->get_ok("/api/1.2/tenants")->status_is(200)
	->json_is( "/response/0/id", $root_tenant_id )
	->json_is( "/response/2/id", $tenantA_id)
	->json_is( "/response/3/id", $tenantD_id)
	->json_is( "/response/4/id", $tenantE_id)
	->json_is( "/response/1/id", $tenantB_id)->or( sub { diag $t->tx->res->content->asset->{content}; } );;

#tenants heirarchy- test depth, height, root
my $tenant_utils_of_root = Utils::Tenant->new(undef, $root_tenant_id, $schema);
my $tenants_data = $tenant_utils_of_root->create_tenants_data_from_db();

ok $tenant_utils_of_root->is_root_tenant($tenants_data, $root_tenant_id) == 1; 
ok $tenant_utils_of_root->get_tenant_heirarchy_depth($tenants_data, $root_tenant_id) == 0; 
ok $tenant_utils_of_root->get_tenant_heirarchy_height($tenants_data, $root_tenant_id) == 2; 

ok $tenant_utils_of_root->is_root_tenant($tenants_data, $tenantA_id) == 0; 
ok $tenant_utils_of_root->get_tenant_heirarchy_depth($tenants_data, $tenantA_id) == 1; 
ok $tenant_utils_of_root->get_tenant_heirarchy_height($tenants_data, $tenantA_id) == 1; 

ok $tenant_utils_of_root->is_root_tenant($tenants_data, $tenantB_id) == 0; 
ok $tenant_utils_of_root->get_tenant_heirarchy_depth($tenants_data, $tenantB_id) == 1; 
ok $tenant_utils_of_root->get_tenant_heirarchy_height($tenants_data, $tenantB_id) == 0; 

ok $tenant_utils_of_root->is_root_tenant($tenants_data, $tenantD_id) == 0; 
ok $tenant_utils_of_root->get_tenant_heirarchy_depth($tenants_data, $tenantD_id) == 2; 
ok $tenant_utils_of_root->get_tenant_heirarchy_height($tenants_data, $tenantD_id) == 0; 

ok $tenant_utils_of_root->is_root_tenant($tenants_data, $tenantE_id) == 0; 
ok $tenant_utils_of_root->get_tenant_heirarchy_depth($tenants_data, $tenantE_id) == 2; 
ok $tenant_utils_of_root->get_tenant_heirarchy_height($tenants_data, $tenantE_id) == 0; 

############################
#testing tenancy checks
#root tenant - touch entire hierarchy as well as null
ok $tenant_utils_of_root->is_tenant_resource_accessible($tenants_data, $root_tenant_id) == 1; 
ok $tenant_utils_of_root->is_tenant_resource_accessible($tenants_data, undef) == 1; 
ok $tenant_utils_of_root->is_tenant_resource_accessible($tenants_data, $tenantA_id) == 1; 
ok $tenant_utils_of_root->is_tenant_resource_accessible($tenants_data, $tenantE_id) == 1; 

my $tenant_utils_of_a = Utils::Tenant->new(undef, $tenantA_id, $schema);
my $tenants_data_of_a = $tenant_utils_of_a->create_tenants_data_from_db();
#parent - no access
ok $tenant_utils_of_a->is_tenant_resource_accessible($tenants_data_of_a, $root_tenant_id) == 0; 
#undef - all have access 
ok $tenant_utils_of_a->is_tenant_resource_accessible($tenants_data_of_a, undef) == 1; 
#itself - full access
ok $tenant_utils_of_a->is_tenant_resource_accessible($tenants_data_of_a, $tenantA_id) == 1; 
# child - full access
ok $tenant_utils_of_a->is_tenant_resource_accessible($tenants_data_of_a, $tenantE_id) == 1; 
# Brother - no access
ok $tenant_utils_of_a->is_tenant_resource_accessible($tenants_data_of_a, $tenantB_id) == 0; 

#leaf test
my $tenant_utils_of_d = Utils::Tenant->new(undef, $tenantD_id, $schema);
my $tenants_data_of_d = $tenant_utils_of_d->create_tenants_data_from_db();
#anchestor - no access
ok $tenant_utils_of_d->is_tenant_resource_accessible($tenants_data_of_d, $root_tenant_id) == 0; 
#undef - all have access 
ok $tenant_utils_of_d->is_tenant_resource_accessible($tenants_data_of_d, undef) == 1; 
# parent - no access
ok $tenant_utils_of_d->is_tenant_resource_accessible($tenants_data_of_d, $tenantA_id) == 0; 
# itself - full access
ok $tenant_utils_of_d->is_tenant_resource_accessible($tenants_data_of_d, $tenantD_id) == 1; 
# uncle - no access
ok $tenant_utils_of_d->is_tenant_resource_accessible($tenants_data_of_d, $tenantB_id) == 0; 

#inactive - nothing can do
my $tenant_utils_of_e = Utils::Tenant->new(undef, $tenantE_id, $schema);
my $tenants_data_of_e = $tenant_utils_of_e->create_tenants_data_from_db();
#anchestor - no access
ok $tenant_utils_of_e->is_tenant_resource_accessible($tenants_data_of_e, $root_tenant_id) == 0; 
#undef - all have access 
ok $tenant_utils_of_e->is_tenant_resource_accessible($tenants_data_of_e, undef) == 0; 
# parent - no access
ok $tenant_utils_of_e->is_tenant_resource_accessible($tenants_data_of_e, $tenantA_id) == 0; 
# itself - full access
ok $tenant_utils_of_e->is_tenant_resource_accessible($tenants_data_of_e, $tenantE_id) == 0; 
# uncle - no access
ok $tenant_utils_of_e->is_tenant_resource_accessible($tenants_data_of_e, $tenantB_id) == 0;


#Test disable capabilities
my $useTenancyParamId = &get_param_id('use_tenancy');
ok $t->put_ok('/api/1.2/parameters/' . $useTenancyParamId => {Accept => 'application/json'} => json => {
			'value'      => '0',
		})->status_is(200)
		->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/name" => "use_tenancy" )
		->json_is( "/response/configFile" => "global" )
		->json_is( "/response/value" => "0" )
    , 'Was the disabling paramter set?';

my $tenant_utils_of_d_disabled = Utils::Tenant->new(undef, $tenantD_id, $schema);
my $tenants_data_of_d_disabled = $tenant_utils_of_d_disabled->create_tenants_data_from_db();
#anchestor - now can access
ok $tenant_utils_of_d_disabled->is_tenant_resource_accessible($tenants_data_of_d_disabled, $root_tenant_id) == 1;
#undef - all have access
ok $tenant_utils_of_d_disabled->is_tenant_resource_accessible($tenants_data_of_d_disabled, undef) == 1;
# parent - now can access
ok $tenant_utils_of_d_disabled->is_tenant_resource_accessible($tenants_data_of_d_disabled, $tenantA_id) == 1;
# itself - full access
ok $tenant_utils_of_d_disabled->is_tenant_resource_accessible($tenants_data_of_d_disabled, $tenantD_id) == 1;
# uncle - now can access
ok $tenant_utils_of_d_disabled->is_tenant_resource_accessible($tenants_data_of_d_disabled, $tenantB_id) == 1;

ok $t->put_ok('/api/1.2/parameters/' . $useTenancyParamId => {Accept => 'application/json'} => json => {
			'value'      => '1',
		})->status_is(200)
		->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/name" => "use_tenancy" )
		->json_is( "/response/configFile" => "global" )
		->json_is( "/response/value" => "1" )
    , 'Was the disabling paramter unset?';


#################
#moving A to be the child of B
ok $t->put_ok('/api/1.2/tenants/' . $tenantA_id  => {Accept => 'application/json'} => json => {
			"active" => 1, "parentId" => $tenantB_id, name => "tenantA2"})
			->status_is(200);
			
ok $t->get_ok("/api/1.2/tenants/$tenantA_id")->status_is(200)
	->json_is( "/response/0/parentId", $tenantB_id)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );;
	
ok $t->get_ok("/api/1.2/tenants/$tenantD_id")->status_is(200)
	->json_is( "/response/0/parentId", $tenantA_id)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );;


#Cannot move B to be the child of itself
ok $t->put_ok('/api/1.2/tenants/' . $tenantB_id  => {Accept => 'application/json'} => json => {
			"active" => 1, "parentId" => $tenantB_id, name => "tenant"})
			->json_is( "/alerts/0/text" => "Parent tenant is invalid: same as updated tenant.")
			->status_is(400);
	
#Cannot move B to be the child of A (a descendant)
ok $t->put_ok('/api/1.2/tenants/' . $tenantB_id  => {Accept => 'application/json'} => json => {
			"active" => 1, "parentId" => $tenantA_id, name => "tenant"})
			->json_is( "/alerts/0/text" => "Parent tenant is invalid: a child of the updated tenant.")
			->status_is(400);

#move A back
ok $t->put_ok('/api/1.2/tenants/' . $tenantA_id  => {Accept => 'application/json'} => json => {
			"active" => 1, "parentId" => $root_tenant_id, name => "tenantA2"})
			->status_is(200);

ok $t->get_ok("/api/1.2/tenants/$tenantA_id")->status_is(200)
	->json_is( "/response/0/parentId", $root_tenant_id)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );;
	
#cannot delete a tenant that have children
ok $t->delete_ok('/api/1.2/tenants/' . $tenantA_id)->status_is(400)
	->json_is( "/alerts/0/text" => "Tenant 'tenantA2' has children tenant(s): e.g 'tenantD'. Please update these tenants and retry." )
	->or( sub { diag $t->tx->res->content->asset->{content}; } );


ok $t->delete_ok('/api/1.2/tenants/' . $tenantE_id)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->delete_ok('/api/1.2/tenants/' . $tenantD_id)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->delete_ok('/api/1.2/tenants/' . $tenantA_id)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

#testing the tenant cascade delete
ok $t->post_ok('/api/1.2/tenants' => {Accept => 'application/json'} => json => {
			"name" => "tenant-b1", "active" => 1, "parentId" => $tenantB_id })->status_is(200);
my $num_of_tenants_before_cascade_delete = $t->get_ok('/api/1.2/tenants')->status_is(200)->$responses_counter();
my $tenant_utils_of_root_for_cascade_delete = Utils::Tenant->new(undef, $root_tenant_id, $schema);
my $tenants_data_for_cascade_delete = $tenant_utils_of_root_for_cascade_delete->create_tenants_data_from_db();
$tenant_utils_of_root_for_cascade_delete->cascade_delete_tenants_tree($tenants_data_for_cascade_delete, $tenantB_id);
ok $t->get_ok('/api/1.2/tenants')->status_is(200)->$count_response_test($num_of_tenants_before_cascade_delete-2);

#cannot delete a tenant that have a delivery-service
ok $t->delete_ok('/api/1.2/tenants/' . 10**9)->status_is(400)
	->json_is( "/alerts/0/text" => "Tenant 'root' is assign with delivery-services(s): e.g. 'test-ds1-root'. Please update/delete these delivery-services and retry." )
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->delete_ok('/api/1.2/deliveryservices/' . 2100)->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

#cannot delete a tenant that have a user
ok $t->delete_ok('/api/1.2/tenants/' . 10**9)->status_is(400)
	->json_is( "/alerts/0/text" => "Tenant 'root' is assign with user(s): e.g. 'admin-root'. Please update these users and retry." )
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# verify a null tenant user cannot access tenats resources he shouldn't
ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';

ok $t->post_ok('/api/1.2/tenants' => {Accept => 'application/json'} => json => {
        "name" => "tenantA", "parentId" => $root_tenant_id })->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } );

#no tenants in the list
ok $t->get_ok("/api/1.2/tenants")->status_is(200)
	->json_is( "/response/0/id", undef)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );;

ok $t->get_ok("/api/1.2/tenants/$root_tenant_id")->status_is(403)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );;

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


