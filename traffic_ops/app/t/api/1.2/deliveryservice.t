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
ok $t->get_ok("/api/1.2/deliveryservices/list")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } )
		->json_is( "/response/0/xmlId", "steering-ds1" )
		->json_is( "/response/0/logsEnabled", 0 )
		->json_is( "/response/0/ipv6RoutingEnabled", 1 )
		->json_is( "/response/1/xmlId", "steering-ds2" );

ok $t->get_ok("/api/1.2/deliveryservices/list?logsEnabled=true")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } )
		->json_is( "/response/0/xmlId", "test-ds1" )
		->json_is( "/response/0/logsEnabled", 1 )
		->json_is( "/response/0/ipv6RoutingEnabled", 1 )
		->json_is( "/response/1/xmlId", "test-ds1-root" )
		->json_is( "/response/1/tenantId", 10**9 )
		->json_is( "/response/2/xmlId", "test-ds4" );

# It creates new delivery services
ok $t->post_ok('/api/1.2/deliveryservices/create' => {Accept => 'application/json'} => json => {
        "xmlId" => "ds_1",
        "displayName" => "ds_displayname_1",
        "protocol" => "1",
        "orgServerFqdn" => "http://10.75.168.91",
        "cdnName" => "cdn1",
        "tenantId" => $tenant_id,
        "profileName" => "CCR1",
        "type" => "HTTP",
        "multiSiteOrigin" => "0",
        "regionalGeoBlocking" => "1",
        "active" => "false",
        "matchList" => [
            {
                "type" =>  "HOST_REGEXP",
                "setNumber" =>  "0",
                "pattern" => ".*\\.ds_1\\..*"
            },
            {
                "type" =>  "HOST_REGEXP",
                "setNumber" =>  "1",
                "pattern" => ".*\\.my_vod1\\..*"
            }
        ]})
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/xmlId" => "ds_1")
    ->json_is( "/response/displayName" => "ds_displayname_1")
    ->json_is( "/response/orgServerFqdn" => "http://10.75.168.91")
    ->json_is( "/response/cdnName" => "cdn1")
    ->json_is( "/response/tenantId" => $tenant_id)
    ->json_is( "/response/profileName" => "CCR1")
    ->json_is( "/response/protocol" => "1")
    ->json_is( "/response/type" => "HTTP")
    ->json_is( "/response/multiSiteOrigin" => "0")
    ->json_is( "/response/regionalGeoBlocking" => "1")
    ->json_is( "/response/active" => "false")
    ->json_is( "/response/matchList/0/type" => "HOST_REGEXP")
    ->json_is( "/response/matchList/0/setNumber" => "0")
    ->json_is( "/response/matchList/0/pattern" => ".*\\.ds_1\\..*")
    ->json_is( "/response/matchList/1/type" => "HOST_REGEXP")
    ->json_is( "/response/matchList/1/setNumber" => "1")
    ->json_is( "/response/matchList/1/pattern" => ".*\\.my_vod1\\..*")
            , 'Does the deliveryservice details return?';

ok $t->post_ok('/api/1.2/deliveryservices/create' => {Accept => 'application/json'} => json => {
        "xmlId" => "ds_2",
        "displayName" => "ds_displayname_2",
        "protocol" => "1",
        "orgServerFqdn" => "http://10.75.168.91",
        "cdnName" => "cdn1",
        "profileName" => "CCR1",
        "type" => "HTTP",
        "multiSiteOrigin" => "0",
        "regionalGeoBlocking" => "1",
        "active" => "false",
        "matchList" => [
            {
                "type" =>  "HOST_REGEXP",
                "setNumber" =>  "0",
                "pattern" => ".*\\.ds_2\\..*"
            },
            {
                "type" =>  "HOST_REGEXP",
                "setNumber" =>  "1",
                "pattern" => ".*\\.my_vod2\\..*"
            }
        ]})
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/xmlId" => "ds_2")
    ->json_is( "/response/displayName" => "ds_displayname_2")
    ->json_is( "/response/orgServerFqdn" => "http://10.75.168.91")
    ->json_is( "/response/cdnName" => "cdn1")
    ->json_is( "/response/tenantId" => $tenant_id)#tenant id is derived from the current tenant
    ->json_is( "/response/profileName" => "CCR1")
    ->json_is( "/response/protocol" => "1")
    ->json_is( "/response/type" => "HTTP")
    ->json_is( "/response/multiSiteOrigin" => "0")
    ->json_is( "/response/regionalGeoBlocking" => "1")
    ->json_is( "/response/active" => "false")
    ->json_is( "/response/matchList/0/type" => "HOST_REGEXP")
    ->json_is( "/response/matchList/0/setNumber" => "0")
    ->json_is( "/response/matchList/0/pattern" => ".*\\.ds_2\\..*")
    ->json_is( "/response/matchList/1/type" => "HOST_REGEXP")
    ->json_is( "/response/matchList/1/setNumber" => "1")
    ->json_is( "/response/matchList/1/pattern" => ".*\\.my_vod2\\..*")
            , 'Does the deliveryservice details return?';

# It allows adding a delivery service to a different CDN with a matching regex in another CDN's delivery service
ok $t->post_ok('/api/1.2/deliveryservices/create' => {Accept => 'application/json'} => json => {
        "xmlId" => "ds_3",
        "displayName" => "ds_displayname_1",
        "protocol" => "1",
        "orgServerFqdn" => "http://10.75.168.91",
        "cdnName" => "cdn2",
        "tenantId" => $tenant_id,
        "profileName" => "CCR1",
        "type" => "HTTP",
        "multiSiteOrigin" => "0",
        "regionalGeoBlocking" => "1",
        "active" => "false",
        "matchList" => [
            {
                "type" =>  "HOST_REGEXP",
                "setNumber" =>  "0",
                "pattern" => ".*\\.ds_1\\..*"
            },
            {
                "type" =>  "HOST_REGEXP",
                "setNumber" =>  "1",
                "pattern" => ".*\\.my_vod2\\..*"
            }
        ]})
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/xmlId" => "ds_3")
    ->json_is( "/response/displayName" => "ds_displayname_1")
    ->json_is( "/response/orgServerFqdn" => "http://10.75.168.91")
    ->json_is( "/response/cdnName" => "cdn2")
    ->json_is( "/response/tenantId" => $tenant_id)
    ->json_is( "/response/profileName" => "CCR1")
    ->json_is( "/response/protocol" => "1")
    ->json_is( "/response/type" => "HTTP")
    ->json_is( "/response/multiSiteOrigin" => "0")
    ->json_is( "/response/regionalGeoBlocking" => "1")
    ->json_is( "/response/active" => "false")
    ->json_is( "/response/matchList/0/type" => "HOST_REGEXP")
    ->json_is( "/response/matchList/0/setNumber" => "0")
    ->json_is( "/response/matchList/0/pattern" => ".*\\.ds_1\\..*")
    ->json_is( "/response/matchList/1/type" => "HOST_REGEXP")
    ->json_is( "/response/matchList/1/setNumber" => "1")
    ->json_is( "/response/matchList/1/pattern" => ".*\\.my_vod2\\..*")
            , 'Does the deliveryservice details return?';

my $ds_id = &get_ds_id('ds_1');

ok $t->put_ok('/api/1.2/deliveryservices/' . $ds_id . '/update'  => {Accept => 'application/json'} => json => {
        "xmlId" => "ds_1",
        "displayName" => "ds_displayname_11",
        "protocol" => "2",
        "orgServerFqdn" => "http://10.75.168.91",
        "cdnName" => "cdn1",
        "tenantId" => $tenant_id,
        "profileName" => "CCR1",
        "type" => "HTTP",
        "multiSiteOrigin" => "0",
        "regionalGeoBlocking" => "0",
        "active" => "true",
        "matchList" => [
            {
                "type" =>  "HOST_REGEXP",
                "setNumber" =>  "0",
                "pattern" => ".*\\.my_vod1\\..*"
            }
        ]})
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/xmlId" => "ds_1")
    ->json_is( "/response/displayName" => "ds_displayname_11")
    ->json_is( "/response/orgServerFqdn" => "http://10.75.168.91")
    ->json_is( "/response/cdnName" => "cdn1")
    ->json_is( "/response/tenantId" => $tenant_id)
    ->json_is( "/response/profileName" => "CCR1")
    ->json_is( "/response/protocol" => "2")
    ->json_is( "/response/type" => "HTTP")
    ->json_is( "/response/multiSiteOrigin" => "0")
    ->json_is( "/response/regionalGeoBlocking" => "0")
    ->json_is( "/response/active" => "true")
    ->json_is( "/response/matchList/0/type" => "HOST_REGEXP")
    ->json_is( "/response/matchList/0/setNumber" => "0")
    ->json_is( "/response/matchList/0/pattern" => ".*\\.my_vod1\\..*")
            , 'Does the deliveryservice details return?';



#removing tenancy when no set
ok $t->put_ok('/api/1.2/deliveryservices/' . $ds_id . '/update'  => {Accept => 'application/json'} => json => {
        "xmlId" => "ds_1",
        "displayName" => "ds_displayname_11",
        "protocol" => "2",
        "orgServerFqdn" => "http://10.75.168.91",
        "cdnName" => "cdn1",
        "profileName" => "CCR1",
        "type" => "HTTP",
        "multiSiteOrigin" => "0",
        "regionalGeoBlocking" => "0",
        "active" => "true",
        "matchList" => [
            {
                "type" =>  "HOST_REGEXP",
                "setNumber" =>  "0",
                "pattern" => ".*\\.my_vod1\\..*"
            }
        ]})
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/xmlId" => "ds_1")
    ->json_is( "/response/displayName" => "ds_displayname_11")
    ->json_is( "/response/orgServerFqdn" => "http://10.75.168.91")
    ->json_is( "/response/cdnName" => "cdn1")
    ->json_is( "/response/tenantId" => undef)
    ->json_is( "/response/profileName" => "CCR1")
    ->json_is( "/response/protocol" => "2")
    ->json_is( "/response/type" => "HTTP")
    ->json_is( "/response/multiSiteOrigin" => "0")
    ->json_is( "/response/regionalGeoBlocking" => "0")
    ->json_is( "/response/active" => "true")
    ->json_is( "/response/matchList/0/type" => "HOST_REGEXP")
    ->json_is( "/response/matchList/0/setNumber" => "0")
    ->json_is( "/response/matchList/0/pattern" => ".*\\.my_vod1\\..*")
            , 'Does the deliveryservice details return?';


#putting tenancy back
ok $t->put_ok('/api/1.2/deliveryservices/' . $ds_id . '/update'  => {Accept => 'application/json'} => json => {
        "xmlId" => "ds_1",
        "displayName" => "ds_displayname_11",
        "protocol" => "2",
        "orgServerFqdn" => "http://10.75.168.91",
        "cdnName" => "cdn1",
        "tenantId" => $tenant_id,
        "profileName" => "CCR1",
        "type" => "HTTP",
        "multiSiteOrigin" => "0",
        "regionalGeoBlocking" => "0",
        "active" => "true",
        "matchList" => [
            {
                "type" =>  "HOST_REGEXP",
                "setNumber" =>  "0",
                "pattern" => ".*\\.my_vod1\\..*"
            }
        ]})
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/xmlId" => "ds_1")
    ->json_is( "/response/displayName" => "ds_displayname_11")
    ->json_is( "/response/orgServerFqdn" => "http://10.75.168.91")
    ->json_is( "/response/cdnName" => "cdn1")
    ->json_is( "/response/tenantId" => $tenant_id)
    ->json_is( "/response/profileName" => "CCR1")
    ->json_is( "/response/protocol" => "2")
    ->json_is( "/response/type" => "HTTP")
    ->json_is( "/response/multiSiteOrigin" => "0")
    ->json_is( "/response/regionalGeoBlocking" => "0")
    ->json_is( "/response/active" => "true")
    ->json_is( "/response/matchList/0/type" => "HOST_REGEXP")
    ->json_is( "/response/matchList/0/setNumber" => "0")
    ->json_is( "/response/matchList/0/pattern" => ".*\\.my_vod1\\..*")
            , 'Does the deliveryservice details return?';


# It allows updating a delivery service in a CDN with a matching regex in another CDN's delivery service

$ds_id = &get_ds_id('ds_3');

ok $t->put_ok('/api/1.2/deliveryservices/' . $ds_id . '/update' => {Accept => 'application/json'} => json => {
		"xmlId" => "ds_3",
		"displayName" => "ds_displayname_1",
		"protocol" => "1",
		"orgServerFqdn" => "http://10.75.168.91",
		"cdnName" => "cdn2",
	        "tenantId" => $tenant_id,
		"profileName" => "CCR1",
		"type" => "HTTP",
		"multiSiteOrigin" => "0",
		"regionalGeoBlocking" => "1",
		"active" => "false",
		"matchList" => [
			{
				"type" =>  "HOST_REGEXP",
				"setNumber" =>  "0",
				"pattern" => ".*\\.ds_1\\..*"
			},
			{
				"type" =>  "HOST_REGEXP",
				"setNumber" =>  "1",
				"pattern" => ".*\\.my_vod1\\..*"
			}
		]})
	->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/xmlId" => "ds_3")
		->json_is( "/response/displayName" => "ds_displayname_1")
		->json_is( "/response/orgServerFqdn" => "http://10.75.168.91")
		->json_is( "/response/cdnName" => "cdn2")
		->json_is( "/response/tenantId" => $tenant_id)
		->json_is( "/response/profileName" => "CCR1")
		->json_is( "/response/protocol" => "1")
		->json_is( "/response/type" => "HTTP")
		->json_is( "/response/multiSiteOrigin" => "0")
		->json_is( "/response/regionalGeoBlocking" => "1")
		->json_is( "/response/active" => "false")
		->json_is( "/response/matchList/0/type" => "HOST_REGEXP")
		->json_is( "/response/matchList/0/setNumber" => "0")
		->json_is( "/response/matchList/0/pattern" => ".*\\.ds_1\\..*")
		->json_is( "/response/matchList/1/type" => "HOST_REGEXP")
		->json_is( "/response/matchList/1/setNumber" => "1")
		->json_is( "/response/matchList/1/pattern" => ".*\\.my_vod1\\..*")
	, 'Does the deliveryservice details return?';

# It does not allow changing a regex to be the same as an existing one in another DS in the same CDN
$ds_id = &get_ds_id('ds_1');
ok $t->put_ok('/api/1.2/deliveryservices/' . $ds_id . '/update'  => {Accept => 'application/json'} => json => {
			"xmlId" => "ds_1",
			"displayName" => "ds_displayname_2",
			"protocol" => "1",
			"orgServerFqdn" => "http://10.75.168.91",
			"cdnName" => "cdn1",
		        "tenantId" => $tenant_id,
			"profileName" => "CCR1",
			"type" => "HTTP",
			"multiSiteOrigin" => "0",
			"regionalGeoBlocking" => "0",
			"active" => "true",
			"matchList" => [
				{
					"type" =>  "HOST_REGEXP",
					"setNumber" =>  "0",
					"pattern" => ".*\\.ds_2\\..*"
				},
				{
					"type" =>  "HOST_REGEXP",
					"setNumber" =>  "1",
					"pattern" => ".*\\.my_vod2\\..*"
				}
			]})
		->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	, 'Does the deliveryservice details return?';

ok $t->delete_ok('/api/1.2/deliveryservices/' . $ds_id)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->put_ok('/api/1.2/deliveryservices/' . $ds_id . '/update'  => {Accept => 'application/json'} => json => {
        "xmlId" => "ds_1234"
})
    ->status_is(404)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok(
	'/api/1.2/deliveryservices/test-ds1/servers' => { Accept => 'application/json' } => json => {
		"serverNames" => [ "atlanta-edge-01", "atlanta-edge-02" ]
	}
	)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/response/xmlId" => "test-ds1" )
	->json_is( "/response/serverNames/0" => "atlanta-edge-01" )->json_is( "/response/serverNames/1" => "atlanta-edge-02" ),
	'Does the assigned servers return?';

ok $t->get_ok("/api/1.2/deliveryservices.json")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } )
	->json_is( "/response/0/xmlId", "ds_2" )->json_is( "/response/0/logsEnabled", 0 )->json_is( "/response/0/ipv6RoutingEnabled", 0 )
	->json_is( "/response/1/xmlId", "ds_3" );

# Count the 'response number'
my $count_response = sub {
	my ( $t, $count ) = @_;
	my $json = decode_json( $t->tx->res->content->asset->slurp );
	my $r    = $json->{response};
	return $t->success( is( scalar(@$r), $count ) );
};

$t->get_ok('/api/1.2/deliveryservices/list?logsEnabled=true')->status_is(200)->$count_response(3)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->put_ok('/api/1.2/snapshot/cdn1')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

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
