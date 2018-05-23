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
my $use_tenancy = shift;

Test::TestHelper->unload_core_data($schema);
Test::TestHelper->load_core_data($schema);

my $tenant_id = $schema->resultset('TmUser')->find( { username => $login_user } )->get_column('tenant_id');
my $tenant_name = defined ($tenant_id) ? $schema->resultset('Tenant')->find( { id => $tenant_id } )->get_column('name') : "null";

ok $t->post_ok( '/login', => form => { u => $login_user, p => $login_password } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';
# Count the 'response number'
my $count_response = sub {
	my ( $t, $count ) = @_;
	my $json = decode_json( $t->tx->res->content->asset->slurp );
	my $r    = $json->{response};
	return $t->success( is( scalar(@$r), $count ) );
};

my $useTenancyParamId = &get_param_id('use_tenancy');
ok $t->put_ok('/api/1.2/parameters/' . $useTenancyParamId => {Accept => 'application/json'} => json => {
            'value'      => $use_tenancy,
        })->status_is(200)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/response/name" => "use_tenancy" )
        ->json_is( "/response/configFile" => "global" )
        ->json_is( "/response/value" => $use_tenancy )
    , 'Was the disabling paramter set?';

    # we will assign 2 more servers to ds 100
ok $t->post_ok('/api/1.2/deliveryserviceserver' => {Accept => 'application/json'} => json => {
			"dsId" => 100,
			"servers" => [ 1400, 1600 ]
		})
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/alerts/0/level" => "success" )
		->json_is( "/alerts/0/text" => "Server assignments complete." )
	, 'Are the servers assigned to the delivery service?';

# there are currently 6 servers of type EDGE or ORG that can be assigned to ds 100
$t->get_ok('/api/1.2/deliveryservices/100/servers/eligible')->status_is(200)->$count_response(6)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

# It gets existing delivery services
ok $t->get_ok("/api/1.2/deliveryservices")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } )
		->json_is( "/response/0/xmlId", "steering-ds1" )
		->json_is( "/response/0/routingName", "foo" )
		->json_is( "/response/0/deepCachingType", "NEVER" )
		->json_is( "/response/0/logsEnabled", 0 )
		->json_is( "/response/0/ipv6RoutingEnabled", 1 )
		->json_is( "/response/1/xmlId", "steering-ds2" );

ok $t->get_ok("/api/1.2/deliveryservices?logsEnabled=true")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } )
		->json_is( "/response/0/xmlId", "test-ds1" )
		->json_is( "/response/0/logsEnabled", 1 )
		->json_is( "/response/0/ipv6RoutingEnabled", 1 )
        ->json_is( "/response/1/xmlId", defined($tenant_id) ? "test-ds1-root" : ($use_tenancy ? "test-ds4" : "test-ds1-root"));

	ok $t->post_ok('/api/1.2/deliveryservices' => {Accept => 'application/json'} => json => {
				"active" => \0,
				"anonymousBlockingEnabled" => 0,
				"cdnId" => 100,
				"dscp" => 0,
				"displayName" => "ds_any_map_display_name",
				"geoLimit" => 0,
				"geoProvider" => 0,
				"logsEnabled" => 0,
				"regionalGeoBlocking" => 0,
				"routingName" => "foo",
				"typeId" => 35,
				"tenantId" => $tenant_id,
				"xmlId" => "ds_any_map",
			})->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
			->json_is( "/response/0/active" => 0)
			->json_is( "/response/0/tenantId" => $tenant_id)
			->json_is( "/response/0/cdnName" => "cdn1")
			->json_is( "/response/0/routingName" => "foo")
			->json_is( "/response/0/displayName" => "ds_any_map_display_name")
			->json_is( "/response/0/xmlId" => "ds_any_map")
			->json_is( "/response/0/type" => "ANY_MAP")
		, 'Is the ANY_MAP delivery service created?';

	ok $t->post_ok('/api/1.2/deliveryservices' => {Accept => 'application/json'} => json => {
				"active" => \0,
				"anonymousBlockingEnabled" => 0,
				"cdnId" => 100,
				"dscp" => 0,
				"displayName" => "ds_steering_display_name",
				"geoLimit" => 0,
				"geoProvider" => 0,
				"ipv6RoutingEnabled" => 0,
				"logsEnabled" => 0,
				"protocol" => 1,
				"regionalGeoBlocking" => 0,
				"routingName" => "foo",
				"typeId" => 37,
				"tenantId" => $tenant_id,
				"xmlId" => "ds_steering",
			})->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
			->json_is( "/response/0/active" => 0)
			->json_is( "/response/0/tenantId" => $tenant_id)
			->json_is( "/response/0/cdnName" => "cdn1")
			->json_is( "/response/0/routingName" => "foo")
			->json_is( "/response/0/displayName" => "ds_steering_display_name")
			->json_is( "/response/0/xmlId" => "ds_steering")
			->json_is( "/response/0/type" => "STEERING")
		, 'Is the STEERING delivery service created?';

	ok $t->post_ok('/api/1.2/deliveryservices' => {Accept => 'application/json'} => json => {
				"active" => \0,
				"anonymousBlockingEnabled" => 1,
				"cdnId" => 100,
				"displayName" => "ds_http_display_name",
				"dscp" => 0,
				"geoLimit" => 0,
				"geoProvider" => 0,
				"initialDispersion" => 1,
				"ipv6RoutingEnabled" => 0,
				"logsEnabled" => 0,
				"missLat" => 45,
				"missLong" => 45,
				"multiSiteOrigin" => 0,
				"orgServerFqdn" => "http://10.75.168.91",
				"protocol" => 1,
				"qstringIgnore" => 0,
				"rangeRequestHandling" => 0,
				"regionalGeoBlocking" => 0,
				"routingName" => "foo",
				"typeId" => 36,
				"tenantId" => $tenant_id,
				"xmlId" => "ds_http",
			})->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
			->json_is( "/response/0/active" => 0)
			->json_is( "/response/0/tenantId" => $tenant_id)
			->json_is( "/response/0/cdnName" => "cdn1")
			->json_is( "/response/0/routingName" => "foo")
			->json_is( "/response/0/displayName" => "ds_http_display_name")
			->json_is( "/response/0/xmlId" => "ds_http")
			->json_is( "/response/0/multiSiteOrigin" => 0)
			->json_is( "/response/0/orgServerFqdn" => "http://10.75.168.91")
			->json_is( "/response/0/protocol" => 1)
			->json_is( "/response/0/type" => "HTTP")
			->json_is( "/response/0/anonymousBlockingEnabled" => 1)
		, 'Is the HTTP delivery service created?';

	ok $t->post_ok('/api/1.2/deliveryservices' => {Accept => 'application/json'} => json => {
				"active" => \0,
				"anonymousBlockingEnabled" => 0,
				"cdnId" => 100,
				"displayName" => "ds_http_display_name",
				"dscp" => 0,
				"geoLimit" => 0,
				"geoProvider" => 0,
				"initialDispersion" => 0,
				"ipv6RoutingEnabled" => 1,
				"logsEnabled" => 0,
				"missLat" => 45,
				"missLong" => 45,
				"multiSiteOrigin" => 0,
				"orgServerFqdn" => "http://10.75.168.91",
				"protocol" => 1,
				"qstringIgnore" => 0,
				"rangeRequestHandling" => 0,
				"regionalGeoBlocking" => 0,
				"routingName" => "foo",
				"typeId" => 36,
				"tenantId" => $tenant_id,
				"xmlId" => "ds_http_2",
			})
			->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
			->json_is( "/alerts/0/text" => "initialDispersion invalid. Must be 1 or greater.")
		, 'Does the HTTP deliveryservice create fail when initial dispersion is set to 0?';

	ok $t->post_ok('/api/1.2/deliveryservices' => {Accept => 'application/json'} => json => {
			"active" => \0,
			"anonymousBlockingEnabled" => 0,
			"cdnId" => 100,
			"deepCachingType" => "NEVER",
			"displayName" => "ds_displayname_1",
			"dscp" => 0,
			"geoLimit" => 0,
			"geoProvider" => 0,
			"ipv6RoutingEnabled" => 0,
			"logsEnabled" => 0,
			"missLat" => 45,
			"missLong" => 45,
			"multiSiteOrigin" => 0,
			"orgServerFqdn" => "http://10.75.168.91",
			"protocol" => 1,
			"qstringIgnore" => 0,
			"rangeRequestHandling" => 0,
			"regionalGeoBlocking" => 0,
			"routingName" => "foo",
			"typeId" => 7,
            "tenantId" => $tenant_id,
            "xmlId" => "ds_1",
		})->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/0/active" => 0)
        ->json_is( "/response/0/tenantId" => $tenant_id)
        ->json_is( "/response/0/cdnName" => "cdn1")
		->json_is( "/response/0/routingName" => "foo")
		->json_is( "/response/0/deepCachingType" => "NEVER")
		->json_is( "/response/0/displayName" => "ds_displayname_1")
		->json_is( "/response/0/xmlId" => "ds_1")
		->json_is( "/response/0/multiSiteOrigin" => 0)
		->json_is( "/response/0/orgServerFqdn" => "http://10.75.168.91")
		->json_is( "/response/0/protocol" => 1)
		->json_is( "/response/0/type" => "DNS")
            , 'Is the DNS delivery service created?';

	my $ds_id = &get_ds_id('ds_1');
	ok $t->put_ok('/api/1.2/deliveryservices/' . $ds_id => {Accept => 'application/json'} => json => {
			"active" => \1,
			"anonymousBlockingEnabled" => 0,
			"cdnId" => 100,
			"deepCachingType" => "NEVER",
            "displayName" => "ds_displayname_11",
			"dscp" => 1,
			"geoLimit" => 1,
			"geoProvider" => 1,
			"ipv6RoutingEnabled" => 1,
			"logsEnabled" => 1,
			"missLat" => 45,
			"missLong" => 45,
			"multiSiteOrigin" => 0,
			"orgServerFqdn" => "http://10.75.168.91",
			"protocol" => 2,
			"qstringIgnore" => 1,
			"rangeRequestHandling" => 1,
			"regionalGeoBlocking" => 0,
			"routingName" => "bar",
			"tenantId" => $tenant_id,
            "typeId" => 7,
			"xmlId" => "ds_1",
        })
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/0/active" => 1)
		->json_is( "/response/0/cdnName" => "cdn1")
		->json_is( "/response/0/displayName" => "ds_displayname_11")
		->json_is( "/response/0/xmlId" => "ds_1")
		->json_is( "/response/0/routingName" => "bar")
		->json_is( "/response/0/deepCachingType" => "NEVER")
		->json_is( "/response/0/multiSiteOrigin" => 0)
		->json_is( "/response/0/orgServerFqdn" => "http://10.75.168.91")
		->json_is( "/response/0/protocol" => 2)
		->json_is( "/response/0/type" => "DNS")
            , 'Is the DNS delivery service updated?';

	ok $t->put_ok('/api/1.2/deliveryservices/' . $ds_id => {Accept => 'application/json'} => json => {
				"active" => \1,
				"anonymousBlockingEnabled" => 0,
				"cdnId" => 100,
				"displayName" => "ds_displayname_11",
				"dscp" => 1,
				"routingName" => "baz",
				"deepCachingType" => "NEVER",
				"geoLimit" => 1,
				"geoProvider" => 1,
				"ipv6RoutingEnabled" => 1,
				"initialDispersion" => 0,
				"logsEnabled" => 1,
				"missLat" => 45,
				"missLong" => 45,
				"multiSiteOrigin" => 0,
				"orgServerFqdn" => "http://10.75.168.91",
				"protocol" => 2,
				"qstringIgnore" => 1,
				"rangeRequestHandling" => 1,
				"regionalGeoBlocking" => 0,
				"tenantId" => $tenant_id,
                "typeId" => 7,
				"xmlId" => "ds_2",
			})
			->json_is( "/alerts/0/text" => "A deliveryservice xmlId is immutable.")
			->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		, 'Does the deliveryservice update fail when you try to change the xmlId?';

	ok $t->delete_ok('/api/1.2/deliveryservices/' . $ds_id)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

	ok $t->put_ok('/api/1.2/deliveryservices/' . $ds_id => {Accept => 'application/json'} => json => {
				"active" => \1,
				"anonymousBlockingEnabled" => 0,
				"cdnId" => 100,
				"displayName" => "ds_displayname_11",
				"dscp" => 1,
				"routingName" => "foo",
				"deepCachingType" => "NEVER",
				"geoLimit" => 1,
				"geoProvider" => 1,
				"initialDispersion" => 2,
				"ipv6RoutingEnabled" => 1,
				"logsEnabled" => 1,
				"missLat" => 45,
				"missLong" => 45,
				"multiSiteOrigin" => 0,
				"orgServerFqdn" => "http://10.75.168.91",
				"protocol" => 2,
				"qstringIgnore" => 1,
				"rangeRequestHandling" => 1,
				"regionalGeoBlocking" => 1,
				"signed" => 1,
                "tenantId" => $tenant_id,
                "typeId" => 7,
				"xmlId" => "ds_1",
	})->status_is(404)->or( sub { diag $t->tx->res->content->asset->{content}; } );

	ok $t->post_ok(
		'/api/1.2/deliveryservices/test-ds1/servers' => { Accept => 'application/json' } => json => {
			"serverNames" => [ "atlanta-edge-01", "atlanta-edge-02" ]
		}
		)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/response/xmlId" => "test-ds1" )
		->json_is( "/response/serverNames/0" => "atlanta-edge-01" )->json_is( "/response/serverNames/1" => "atlanta-edge-02" ),
		'Does the assigned servers return?';

	ok $t->get_ok('/api/1.2/deliveryservices?logsEnabled=true')->status_is(200)->$count_response(defined($tenant_id) ? 4 : ($use_tenancy ? 3 : 4));

	ok $t->put_ok('/api/1.2/snapshot/cdn1')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

	# try deleting a deliveryservice that has a job and a staticdnsentry
	ok $t->delete_ok('/api/1.2/deliveryservices/100')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

	ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

}

my $schema = Schema->connect_to_database;
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');
run_ut($t, $schema, Test::TestHelper::ADMIN_USER,  Test::TestHelper::ADMIN_USER_PASSWORD, 0);
run_ut($t, $schema, Test::TestHelper::ADMIN_ROOT_USER,  Test::TestHelper::ADMIN_ROOT_USER_PASSWORD, 0);
run_ut($t, $schema, Test::TestHelper::ADMIN_ROOT_USER,  Test::TestHelper::ADMIN_ROOT_USER_PASSWORD, 1);

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

