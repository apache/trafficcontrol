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
use Data::Dumper;
use DBI;
use JSON;
use strict;
use warnings;
no warnings 'once';
use warnings 'all';
use Test::TestHelper;

#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

BEGIN { $ENV{MOJO_MODE} = "test" }


sub run_ut {

my $t = shift;
my $schema = shift;
my $login_user = shift;
my $login_password = shift;

Test::TestHelper->unload_core_data($schema);
Test::TestHelper->load_core_data($schema);

my $tenant_id = $schema->resultset('TmUser')->find( { username => $login_user } )->get_column('tenant_id');
my $tenant_name = defined ($tenant_id) ? $schema->resultset('Tenant')->find( { id => $tenant_id } )->get_column('name') : "null";

ok $t->post_ok( '/login', => form => { u => $login_user, p => $login_password } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';

# It gets existing delivery services
ok $t->get_ok("/api/1.2/deliveryservices")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } )
		->json_is( "/response/0/xmlId", "steering-ds1" )
		->json_is( "/response/0/logsEnabled", 0 )
		->json_is( "/response/0/ipv6RoutingEnabled", 1 )
		->json_is( "/response/1/xmlId", "steering-ds2" );

# It creates new delivery services
ok $t->post_ok('/api/1.2/deliveryservices' => {Accept => 'application/json'} => json => {
        "xmlId" => "ds_1",
        "displayName" => "ds_displayname_1",
        "protocol" => "1",
        "orgServerFqdn" => "http://10.75.168.91",
        "cdnName" => "cdn1",
        "tenantId" => $tenant_id,
        "profileId" => 300,
        "typeId" => "36",
        "multiSiteOrigin" => "0",
        "regionalGeoBlocking" => "1",
        "active" => "false",
        "dscp" => 0,
        "ipv6RoutingEnabled" => "true",
        "logsEnabled" => "true",
        "initialDispersion" => 0,
        "cdnId" => 100,
        "signed" => "false",
        "rangeRequestHandling" => 0,
        "geoLimit" => 0,
        "geoProvider" => 0,
        "qstringIgnore" => 0,
        })
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/xmlId" => "ds_1")->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/displayName" => "ds_displayname_1")
    ->json_is( "/response/0/orgServerFqdn" => "http://10.75.168.91")
    ->json_is( "/response/0/cdnId" => 100)
    ->json_is( "/response/0/tenantId" => $tenant_id)
    ->json_is( "/response/0/profileId" => 300)
    ->json_is( "/response/0/protocol" => "1")
    ->json_is( "/response/0/typeId" => 36)
    ->json_is( "/response/0/multiSiteOrigin" => "0")
    ->json_is( "/response/0/regionalGeoBlocking" => "1")
    ->json_is( "/response/0/active" => "false")
            , 'Was the DS properly added and reported?';


my $ds_id = &get_ds_id('ds_1');

ok $t->get_ok("/api/1.2/deliveryservices/$ds_id")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } )
    ->json_is( "/response/0/xmlId" => "ds_1")->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/displayName" => "ds_displayname_1")
    ->json_is( "/response/0/orgServerFqdn" => "http://10.75.168.91")
    ->json_is( "/response/0/cdnId" => 100)
    ->json_is( "/response/0/tenantId" => $tenant_id)
    ->json_is( "/response/0/profileId" => 300)
    ->json_is( "/response/0/protocol" => "1")
    ->json_is( "/response/0/typeId" => 36)
    ->json_is( "/response/0/multiSiteOrigin" => "0")
    ->json_is( "/response/0/regionalGeoBlocking" => "1")
    ->json_is( "/response/0/active" => "0")
            , 'Does the deliveryservice details return?';

# A minor change
ok $t->put_ok('/api/1.2/deliveryservices/'.$ds_id => {Accept => 'application/json'} => json => {
        "xmlId" => "ds_1",
        "displayName" => "ds_displayname_1",
        "protocol" => "1",
        "orgServerFqdn" => "http://10.75.168.92",
        "cdnName" => "cdn1",
        "tenantId" => $tenant_id,
        "profileId" => 300,
        "typeId" => "36",
        "multiSiteOrigin" => "0",
        "regionalGeoBlocking" => "1",
        "active" => "false",
        "dscp" => 0,
        "ipv6RoutingEnabled" => "true",
        "logsEnabled" => "true",
        "initialDispersion" => 0,
        "cdnId" => 100,
        "signed" => "false",
        "rangeRequestHandling" => 0,
        "geoLimit" => 0,
        "geoProvider" => 0,
        "qstringIgnore" => 0,
        })
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/xmlId" => "ds_1")->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/displayName" => "ds_displayname_1")
    ->json_is( "/response/0/orgServerFqdn" => "http://10.75.168.92")
    ->json_is( "/response/0/cdnId" => 100)
    ->json_is( "/response/0/tenantId" => $tenant_id)
    ->json_is( "/response/0/profileId" => 300)
    ->json_is( "/response/0/protocol" => "1")
    ->json_is( "/response/0/typeId" => 36)
    ->json_is( "/response/0/multiSiteOrigin" => "0")
    ->json_is( "/response/0/regionalGeoBlocking" => "1")
    ->json_is( "/response/0/active" => "false")
            , 'A minor change';


# Removing tenancy by not putting it in the put
ok $t->put_ok('/api/1.2/deliveryservices/'.$ds_id => {Accept => 'application/json'} => json => {
        "xmlId" => "ds_1",
        "displayName" => "ds_displayname_1",
        "protocol" => "1",
        "orgServerFqdn" => "http://10.75.168.92",
        "cdnName" => "cdn1",
        "profileId" => 300,
        "typeId" => "36",
        "multiSiteOrigin" => "0",
        "regionalGeoBlocking" => "1",
        "active" => "false",
        "dscp" => 0,
        "ipv6RoutingEnabled" => "true",
        "logsEnabled" => "true",
        "initialDispersion" => 0,
        "cdnId" => 100,
        "signed" => "false",
        "rangeRequestHandling" => 0,
        "geoLimit" => 0,
        "geoProvider" => 0,
        "qstringIgnore" => 0,
        })
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/xmlId" => "ds_1")->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/displayName" => "ds_displayname_1")
    ->json_is( "/response/0/tenantId" => undef)
            , 'Was tenant id removed?';

# Putting tenant id back
ok $t->put_ok('/api/1.2/deliveryservices/'.$ds_id => {Accept => 'application/json'} => json => {
        "xmlId" => "ds_1",
        "displayName" => "ds_displayname_1",
        "protocol" => "1",
        "orgServerFqdn" => "http://10.75.168.92",
        "cdnName" => "cdn1",
        "tenantId" => $tenant_id,
        "profileId" => 300,
        "typeId" => "36",
        "multiSiteOrigin" => "0",
        "regionalGeoBlocking" => "1",
        "active" => "false",
        "dscp" => 0,
        "ipv6RoutingEnabled" => "true",
        "logsEnabled" => "true",
        "initialDispersion" => 0,
        "cdnId" => 100,
        "signed" => "false",
        "rangeRequestHandling" => 0,
        "geoLimit" => 0,
        "geoProvider" => 0,
        "qstringIgnore" => 0,
        })
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/xmlId" => "ds_1")->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/displayName" => "ds_displayname_1")
    ->json_is( "/response/0/tenantId" => $tenant_id)
            , 'Was the tenant ID set again?';


# removing tenant id
ok $t->put_ok('/api/1.2/deliveryservices/'.$ds_id => {Accept => 'application/json'} => json => {
        "xmlId" => "ds_1",
        "displayName" => "ds_displayname_1",
        "protocol" => "1",
        "orgServerFqdn" => "http://10.75.168.92",
        "cdnName" => "cdn1",
        "tenantId" => undef,
        "profileId" => 300,
        "typeId" => "36",
        "multiSiteOrigin" => "0",
        "regionalGeoBlocking" => "1",
        "active" => "false",
        "dscp" => 0,
        "ipv6RoutingEnabled" => "true",
        "logsEnabled" => "true",
        "initialDispersion" => 0,
        "cdnId" => 100,
        "signed" => "false",
        "rangeRequestHandling" => 0,
        "geoLimit" => 0,
        "geoProvider" => 0,
        "qstringIgnore" => 0,
        })
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/xmlId" => "ds_1")->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/displayName" => "ds_displayname_1")
    ->json_is( "/response/0/tenantId" => undef)
            , 'Was the tenant ID set again?';

ok $t->delete_ok('/api/1.2/deliveryservices/' . $ds_id)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# It creates new delivery services, tenant id derived from user
ok $t->post_ok('/api/1.2/deliveryservices' => {Accept => 'application/json'} => json => {
        "xmlId" => "ds_1",
        "displayName" => "ds_displayname_1",
        "protocol" => "1",
        "orgServerFqdn" => "http://10.75.168.91",
        "cdnName" => "cdn1",
        "profileId" => 300,
        "typeId" => "36",
        "multiSiteOrigin" => "0",
        "regionalGeoBlocking" => "1",
        "active" => "false",
        "dscp" => 0,
        "ipv6RoutingEnabled" => "true",
        "logsEnabled" => "true",
        "initialDispersion" => 0,
        "cdnId" => 100,
        "signed" => "false",
        "rangeRequestHandling" => 0,
        "geoLimit" => 0,
        "geoProvider" => 0,
        "qstringIgnore" => 0,
        })
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/xmlId" => "ds_1")->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/tenantId" => $tenant_id)
            , 'Was the tenant id dervied from the creating user?';

my $ds_id = &get_ds_id('ds_1');
#ok $t->delete_ok('/api/1.2/deliveryservices/' . $ds_id)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

}

my $schema = Schema->connect_to_database;
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');

run_ut($t, $schema, Test::TestHelper::ADMIN_USER,  Test::TestHelper::ADMIN_USER_PASSWORD);
run_ut($t, $schema, Test::TestHelper::ADMIN_ROOT_USER,  Test::TestHelper::ADMIN_ROOT_USER_PASSWORD);

$dbh->disconnect();
done_testing();

sub get_ds_id {
    my $xml_id = shift;
    my $q      = "select id from deliveryservice where xml_id = \'$xml_id\'";
    my $get_svr = $dbh->prepare($q);
    $get_svr->execute();
    my $p = $get_svr->fetchall_arrayref( {} );
    $get_svr->finish();
    my $id = $p->[0]->{id};
    return $id;
}
