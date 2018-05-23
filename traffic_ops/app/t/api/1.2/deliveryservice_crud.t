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
	my $tenant_name = defined ($tenant_id) ? $schema->resultset('Tenant')->find( { id => $tenant_id } )->get_column('name') : undef;

	ok $t->post_ok( '/login', => form => { u => $login_user, p => $login_password } )->status_is(302)
		->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';

	my $useTenancyParamId = &get_param_id('use_tenancy');
	ok $t->put_ok('/api/1.2/parameters/' . $useTenancyParamId => {Accept => 'application/json'} => json => {
				'value'      => $use_tenancy,
			})->status_is(200)
			->or( sub { diag $t->tx->res->content->asset->{content}; } )
			->json_is( "/response/name" => "use_tenancy" )
			->json_is( "/response/configFile" => "global" )
			->json_is( "/response/value" => $use_tenancy )
		, 'Was the disabling paramter set?';

	# It gets existing delivery services
	ok $t->get_ok("/api/1.2/deliveryservices")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } )
			->json_is( "/response/0/xmlId", "steering-ds1" )
			->json_is( "/response/0/routingName", "foo" )
			->json_is( "/response/0/deepCachingType", "NEVER")
			->json_is( "/response/0/logsEnabled", 0 )
			->json_is( "/response/0/ipv6RoutingEnabled", 1 )
			->json_is( "/response/1/xmlId", "steering-ds2" );

	# Delivery service create failure due to bad protocol
	ok $t->post_ok('/api/1.2/deliveryservices' => {Accept => 'application/json'} => json => {
				"xmlId" => "ds_1",
				"displayName" => "ds_displayname_1",
				"protocol" => "one",
				"orgServerFqdn" => "http://10.75.168.91",
				"tenantId" => $tenant_id,
				"profileId" => 300,
				"typeId" => 36,
				"multiSiteOrigin" => 0,
				"missLat" => 45,
				"missLong" => 45,
				"regionalGeoBlocking" => 1,
				"anonymousBlockingEnabled" => "0",
				"active" => 0,
				"dscp" => 0,
				"routingName" => "foo",
				"ipv6RoutingEnabled" => 1,
				"logsEnabled" => 1,
				"initialDispersion" => 1,
				"cdnId" => 100,
				"signed" => 0,
				"rangeRequestHandling" => 0,
				"geoLimit" => 0,
				"geoProvider" => 0,
				"qstringIgnore" => 0,
			})
			->status_is(400)
			->json_is( "/alerts/0/text/",
			"protocol invalid. Must be a whole number or null." )->or( sub { diag $t->tx->res->content->asset->{content}; } );

	# Delivery service create failure due to bad missLat
	ok $t->post_ok('/api/1.2/deliveryservices' => {Accept => 'application/json'} => json => {
				"xmlId" => "ds_1",
				"displayName" => "ds_displayname_1",
				"protocol" => 1,
				"orgServerFqdn" => "http://10.75.168.91",
				"tenantId" => $tenant_id,
				"profileId" => 300,
				"typeId" => 36,
				"multiSiteOrigin" => 0,
				"missLat" => "string",
				"missLong" => 45,
				"regionalGeoBlocking" => 1,
				"anonymousBlockingEnabled" => "0",
				"active" => 0,
				"dscp" => 0,
				"routingName" => "foo",
				"ipv6RoutingEnabled" => 1,
				"logsEnabled" => 1,
				"initialDispersion" => 1,
				"cdnId" => 100,
				"signed" => 0,
				"rangeRequestHandling" => 0,
				"geoLimit" => 0,
				"geoProvider" => 0,
				"qstringIgnore" => 0,
			})
			->status_is(400)
			->json_is( "/alerts/0/text/",
			"missLat invalid. Must be a number or null." )->or( sub { diag $t->tx->res->content->asset->{content}; } );

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
        	"missLat" => 45,
        	"missLong" => 45,
        	"regionalGeoBlocking" => "1",
			"anonymousBlockingEnabled" => "0",
        	"active" => "false",
        	"dscp" => 0,
        	"routingName" => "foo",
        	"deepCachingType" => "NEVER",
        	"ipv6RoutingEnabled" => "true",
        	"logsEnabled" => "true",
        	"initialDispersion" => 1,
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
	    ->json_is( "/response/0/routingName" => "foo")
	    ->json_is( "/response/0/deepCachingType" => "NEVER")
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
	    ->json_is( "/response/0/routingName" => "foo")
	    ->json_is( "/response/0/deepCachingType" => "NEVER")
	    ->json_is( "/response/0/tenantId" => $tenant_id)
	    ->json_is( "/response/0/tenant" => $tenant_name)
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
			"missLat" => 45,
			"missLong" => 45,
			"regionalGeoBlocking" => "1",
			"anonymousBlockingEnabled" => "0",
	        "active" => "false",
	        "dscp" => 0,
	        "routingName" => "foo",
	        "deepCachingType" => "NEVER",
	        "ipv6RoutingEnabled" => "true",
	        "logsEnabled" => "true",
	        "initialDispersion" => 1,
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


	if ($use_tenancy) {
		# Removing tenancy by not putting it in the put - should fail
		ok $t->put_ok('/api/1.2/deliveryservices/'.$ds_id => { Accept => 'application/json' } => json => {
					"xmlId"                => "ds_1",
					"displayName"          => "ds_displayname_1",
					"protocol"             => "1",
					"orgServerFqdn"        => "http://10.75.168.92",
					"cdnName"              => "cdn1",
					"profileId"            => 300,
					"typeId"               => "36",
					"multiSiteOrigin"      => "0",
					"missLat" => 45,
					"missLong" => 45,
					"regionalGeoBlocking"  => "1",
					"anonymousBlockingEnabled" => "0",
					"active"               => "false",
					"dscp"                 => 0,
					"routingName"          => "foo",
					"deepCachingType"      => "NEVER",
					"ipv6RoutingEnabled"   => "true",
					"logsEnabled"          => "true",
					"initialDispersion"    => 1,
					"cdnId"                => 100,
					"signed"               => "false",
					"rangeRequestHandling" => 0,
					"geoLimit"             => 0,
					"geoProvider"          => 0,
					"qstringIgnore"        => 0,
				})
				->status_is(400)
				->json_is( "/alerts/0/text/",
				"Invalid tenant. Cannot clear the delivery-service tenancy." )->or( sub { diag $t->tx->res->content->asset->{content}; } )
			,
			, 'Cannot remove tenant by forgetting it?';

		# removing tenant id
		ok $t->put_ok('/api/1.2/deliveryservices/'.$ds_id => { Accept => 'application/json' } => json => {
					"xmlId"                => "ds_1",
					"displayName"          => "ds_displayname_1",
					"protocol"             => "1",
					"orgServerFqdn"        => "http://10.75.168.92",
					"cdnName"              => "cdn1",
					"tenantId"             => undef,
					"profileId"            => 300,
					"typeId"               => "36",
					"multiSiteOrigin"      => "0",
					"missLat" => 45,
					"missLong" => 45,
					"regionalGeoBlocking"  => "1",
					"anonymousBlockingEnabled" => "0",
					"active"               => "false",
					"dscp"                 => 0,
					"routingName"          => "foo",
					"deepCachingType"      => "NEVER",
					"ipv6RoutingEnabled"   => "true",
					"logsEnabled"          => "true",
					"initialDispersion"    => 1,
					"cdnId"                => 100,
					"signed"               => "false",
					"rangeRequestHandling" => 0,
					"geoLimit"             => 0,
					"geoProvider"          => 0,
					"qstringIgnore"        => 0,
				})
				->status_is(400)
				->json_is( "/alerts/0/text/",
				"Invalid tenant. Cannot clear the delivery-service tenancy." )->or( sub { diag $t->tx->res->content->asset->{content}; } )
			,
			, 'Cannot remove tenant?';
	}
	else{
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
				"missLat" => 45,
				"missLong" => 45,
				"regionalGeoBlocking" => "1",
				"anonymousBlockingEnabled" => "0",
       	        "active" => "false",
       	        "dscp" => 0,
       	        "routingName" => "foo",
       	        "deepCachingType" => "NEVER",
       	        "ipv6RoutingEnabled" => "true",
       	        "logsEnabled" => "true",
       	        "initialDispersion" => 1,
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
				"missLat" => 45,
				"missLong" => 45,
				"regionalGeoBlocking" => "1",
				"anonymousBlockingEnabled" => "0",
       	        "active" => "false",
       	        "dscp" => 0,
       	        "routingName" => "foo",
       	        "deepCachingType" => "NEVER",
       	        "ipv6RoutingEnabled" => "true",
       	        "logsEnabled" => "true",
       	        "initialDispersion" => 1,
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
				"missLat" => 45,
				"missLong" => 45,
				"regionalGeoBlocking" => "1",
				"anonymousBlockingEnabled" => "0",
       	        "active" => "false",
       	        "dscp" => 0,
       	        "routingName" => "foo",
       	        "deepCachingType" => "NEVER",
       	        "ipv6RoutingEnabled" => "true",
       	        "logsEnabled" => "true",
       	        "initialDispersion" => 1,
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
	}
	ok $t->delete_ok('/api/1.2/deliveryservices/' . $ds_id)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

	# It creates new delivery services, tenant id is derived from user if not use_tenancy
	ok $t->post_ok('/api/1.2/deliveryservices' => {Accept => 'application/json'} => json => {
	        "xmlId" => "ds_1",
	        "displayName" => "ds_displayname_1",
	        "protocol" => "1",
	        "orgServerFqdn" => "http://10.75.168.91",
	        "cdnName" => "cdn1",
	        "profileId" => 300,
	        "typeId" => "36",
	        "multiSiteOrigin" => "0",
			"missLat" => 45,
			"missLong" => 45,
			"regionalGeoBlocking" => "1",
			"anonymousBlockingEnabled" => "0",
	        "active" => "false",
	        "dscp" => 0,
	        "routingName" => "foo",
	        "deepCachingType" => "NEVER",
	        "ipv6RoutingEnabled" => "true",
	        "logsEnabled" => "true",
	        "initialDispersion" => 1,
	        "cdnId" => 100,
	        "signed" => "false",
	        "rangeRequestHandling" => 0,
	        "geoLimit" => 0,
	        "geoProvider" => 0,
	        "qstringIgnore" => 0,
			"tenantId" => $use_tenancy ? $tenant_id : undef,
	        })
	    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	    ->json_is( "/response/0/xmlId" => "ds_1")->or( sub { diag $t->tx->res->content->asset->{content}; } )
	    ->json_is( "/response/0/tenantId" => $use_tenancy ? $tenant_id : undef)
	            , 'Was the tenant id dervied from the creating user?';

	ok $t->delete_ok('/api/1.2/deliveryservices/' . &get_ds_id('ds_1'))->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

	#change the use_tenancy parameter to 0 (id from Parameters fixture) to test assigned dses table
	ok $t->put_ok('/api/1.2/parameters/67' => {Accept => 'application/json'} => json => 
        {
			'value'       => '0',
        }
	)->status_is(200);

	ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );


	# test safe update route with a portal user
	ok $t->post_ok( '/api/1.2/user/login', json => { u => Test::TestHelper::PORTAL_USER, p => Test::TestHelper::PORTAL_USER_PASSWORD } )->status_is(200),
	'Log into the portal user?';
	
	my $ds_id_portal = &get_ds_id('test-ds1');

	#attempt to change many fields, including the 4 allowed and verify only the 4 actually change
	ok $t->put_ok('/api/1.2/deliveryservices/'.$ds_id_portal.'/safe' => {Accept => 'application/json'} => json => {
        "xmlId" => "test-ds1",
        "displayName" => "ds_displayname_1_new",
        "orgServerFqdn" => "http://10.75.168.91",
        "cdnName" => "cdn1_bad",
        "tenantId" => $tenant_id,
        "profileId" => 300,
        "typeId" => "36",
        "multiSiteOrigin" => "0",
		"missLat" => 45,
		"missLong" => 45,
		"regionalGeoBlocking" => "1",
		"anonymousBlockingEnabled" => "0",
        "active" => "false",
        "dscp" => 0,
        "routingName" => "foo",
        "deepCachingType" => "NEVER",
        "ipv6RoutingEnabled" => "true",
        "logsEnabled" => "true",
        "initialDispersion" => 1,
        "cdnId" => 100,
        "signed" => "false",
        "rangeRequestHandling" => 0,
        "geoLimit" => 0,
        "geoProvider" => 0,
        "qstringIgnore" => 0,
        "infoUrl"    => "http://knutsel-update-new.com",
		"longDesc"   => "long_update_new",
		"longDesc1" => "cust_update_new",
        })
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/xmlId" => "test-ds1")->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/displayName" => "ds_displayname_1_new")
    ->json_is( "/response/0/orgServerFqdn" => "http://test-ds1.edge")
    ->json_is( "/response/0/cdnId" => 100)
    ->json_is( "/response/0/profileId" => 100)
    ->json_is( "/response/0/protocol" => "1")
    ->json_is( "/response/0/typeId" => 21)
    ->json_is( "/response/0/multiSiteOrigin" => "0")
    ->json_is( "/response/0/regionalGeoBlocking" => "1")
    ->json_is( "/response/0/active" => "1")
    ->json_is("/response/0/infoUrl" => "http://knutsel-update-new.com")
    ->json_is("/response/0/longDesc" => "long_update_new")
    ->json_is("/response/0/longDesc1" => "cust_update_new")
            , 'A safe update only changes safe fields';

	my $ds_id_portal_unassigned = &get_ds_id('test-ds2');

	ok $t->put_ok('/api/1.2/deliveryservices/'.$ds_id_portal_unassigned.'/safe' => {Accept => 'application/json'} => json => {
        "xmlId" => "test-ds1",
        "displayName" => "ds_displayname_1_new",
        "orgServerFqdn" => "http://10.75.168.91",
        "cdnName" => "cdn1_bad",
        "tenantId" => $tenant_id,
        "profileId" => 300,
        "typeId" => "36",
        "multiSiteOrigin" => "0",
		"missLat" => 45,
		"missLong" => 45,
		"regionalGeoBlocking" => "1",
		"anonymousBlockingEnabled" => "0",
        "active" => "false",
        "dscp" => 0,
        "routingName" => "foo",
        "deepCachingType" => "NEVER",
        "ipv6RoutingEnabled" => "true",
        "logsEnabled" => "true",
        "initialDispersion" => 1,
        "cdnId" => 100,
        "signed" => "false",
        "rangeRequestHandling" => 0,
        "geoLimit" => 0,
        "geoProvider" => 0,
        "qstringIgnore" => 0,
        "infoUrl"    => "http://knutsel-update-new.com",
		"longDesc"   => "long_update_new",
		"longDesc1" => "cust_update_new",
        })
	->status_is(403)
	->json_is( "/alerts/0/text/", "Forbidden. Delivery service not assigned to user." )->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Can a portal user update an unassigned delivery service?';

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

