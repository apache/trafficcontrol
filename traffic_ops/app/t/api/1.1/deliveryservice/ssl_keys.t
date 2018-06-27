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
use Schema;
use Test::TestHelper;
use Test::MockModule;
use Test::MockObject;
use strict;
use warnings;
use JSON;

BEGIN { $ENV{MOJO_MODE} = "test" }

my $schema = Schema->connect_to_database;
my $t      = Test::Mojo->new('TrafficOps');

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);

#load core test data
Test::TestHelper->load_core_data($schema);

# create ssl key
my $key      = "test-ds1";
my $version  = 1;
my $country  = "US";
my $state    = "Colorado";
my $city     = "Denver";
my $org      = "KableTown";
my $unit     = "CDN_Eng";
my $hostname = "foober.com";
my $cdn = "cdn1";
my $deliveryservice = "test-ds1";

# PORTAL
#NEGATIVE TESTING -- No Privs
ok $t->post_ok( '/api/1.1/user/login', json => { u => Test::TestHelper::READ_ONLY_ROOT_USER, p => Test::TestHelper::READ_ONLY_ROOT_USER_PASSWORD } )->status_is(200),
	'Log into the portal user?';

#create
ok $t->post_ok(
	'/api/1.1/deliveryservices/sslkeys/generate',
	json => {
		key     => $key,
		version => $version,
		deliveryservice => $deliveryservice,
	}
	)->status_is(403)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

#get_object
ok $t->get_ok("/api/1.1/deliveryservices/xmlId/$key/sslkeys.json")->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

# #delete
ok $t->get_ok("/api/1.1/deliveryservices/xmlId/$key/sslkeys/delete.json")->status_is(403)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

# # logout
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

#login as admin
ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

my $ssl_keys = {
	businessUnit => $unit,
	version      => $version,
	hostname     => $hostname,
	certificate  => {
		key => "some_key",
		csr => "some_cst",
		crt => "some_crt"
	},
	country      => $country,
	organization => $org,
	city         => $city,
	state        => $state,
	cdn 		 => "cdn1",
	deliveryservice => $key
};

my $fake_lwp = new Test::MockModule( 'LWP::UserAgent', no_auto => 1 );

my $fake_get = HTTP::Response->new( 200, undef, HTTP::Headers->new, encode_json($ssl_keys) );
$fake_lwp->mock( 'get', sub { return $fake_get } );

my $fake_put = HTTP::Response->new( 204, undef, HTTP::Headers->new, undef );
$fake_lwp->mock( 'put', sub { return $fake_put } );

my $fake_delete = HTTP::Response->new( 204, undef, HTTP::Headers->new, undef );
$fake_lwp->mock( 'delete', sub { return $fake_delete } );

ok $t->post_ok(
	'/api/1.1/deliveryservices/sslkeys/generate',
	json => {
		key          => $key,
		version      => $version,
		hostname     => $hostname,
		deliveryservice     => $deliveryservice, 
		country      => $country,
		state        => $state,
		city         => $city,
		organization => $org,
		businessUnit => $unit,
	}
	)->status_is(200)

	# ->json_is( "/response" => "Successfully created $key_type keys for $key" )
	->json_has( "/response" => "Successfully created " )
	->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Create ssl key?';

# validate ssl key exists

ok $t->get_ok("/api/1.1/deliveryservices/xmlId/$key/sslkeys.json?version=$version")->json_has("/response")->json_has("/response/certificate/csr")
	->json_has("/response/certificate/key")->json_has("/response/certificate/crt")->json_is( "/response/organization" => $org )
	->json_is( "/response/state" => $state )->json_is( "/response/city" => $city )->json_is( "/response/businessUnit" => $unit )
	->json_is( "/response/version" => $version )->json_is( "/response/country" => $country )->json_is( "/response/hostname" => $hostname )->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

# #get latest key
ok $t->get_ok("/api/1.1/deliveryservices/xmlId/$key/sslkeys.json")->json_has("/response")->json_has("/response/certificate/csr")
	->json_has("/response/certificate/key")->json_has("/response/certificate/crt")->json_is( "/response/organization" => $org )
	->json_is( "/response/state" => $state )->json_is( "/response/city" => $city )->json_is( "/response/businessUnit" => $unit )
	->json_is( "/response/version" => $version )->json_is( "/response/country" => $country )->json_is( "/response/hostname" => $hostname )->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

# #get key with period
ok $t->get_ok("/api/1.1/deliveryservices/xmlId/xxfoo.bar/sslkeys.json")->json_has("/response")->json_has("/response/certificate/csr")
	->json_has("/response/certificate/key")->json_has("/response/certificate/crt")->json_is( "/response/organization" => $org )
	->json_is( "/response/state" => $state )->json_is( "/response/city" => $city )->json_is( "/response/businessUnit" => $unit )
	->json_is( "/response/version" => $version )->json_is( "/response/country" => $country )->json_is( "/response/hostname" => $hostname )->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

#get key by hostname
my $gen_hostname = "edge.foo.cdn1.kabletown.net";
ok $t->get_ok("/api/1.1/deliveryservices/hostname/$gen_hostname/sslkeys.json")->json_has("/response")->json_has("/response/certificate/csr")
	->json_has("/response/certificate/key")->json_has("/response/certificate/crt")->json_is( "/response/organization" => $org )
	->json_is( "/response/state" => $state )->json_is( "/response/city" => $city )->json_is( "/response/businessUnit" => $unit )
	->json_is( "/response/version" => $version )->json_is( "/response/country" => $country )->json_is( "/response/hostname" => $hostname )->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );


#tenancy checks
#get_object
ok $t->get_ok("/api/1.1/deliveryservices/xmlId/test-ds1-root/sslkeys.json")->status_is(403)
		->json_has("Forbidden. Delivery-service tenant is not available to the user.!")->or( sub { diag $t->tx->res->content->asset->{content}; } );

#delete
ok $t->get_ok("/api/1.1/deliveryservices/xmlId/test-ds1-root/sslkeys/delete.json")->status_is(403)
		->json_has("Forbidden. Delivery-service tenant is not available to the user.!")->or( sub { diag $t->tx->res->content->asset->{content}; } );

# #delete ssl key
# #delete version
ok $t->get_ok("/api/1.1/deliveryservices/xmlId/$key/sslkeys/delete.json?version=$version")
	->json_is( "/response" => "Successfully deleted ssl keys for $key-$version" )

	# ->json_has( "Successfully deleted" )
	->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

#delete latest
ok $t->get_ok("/api/1.1/deliveryservices/xmlId/$key/sslkeys/delete.json")->json_is( "/response" => "Successfully deleted ssl keys for $key-latest" )

	# ->json_has( "Successfully deleted" )
	->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

#add ssl keys
ok $t->post_ok(
	'/api/1.1/deliveryservices/sslkeys/add',
	json => {
		key         => $key,
		version     => $version,
		certificate => {
			csr => "csr",
			crt => "crt",
			key => "private key"
		},
		deliveryservice => $key,
		cdn => "foo",
		hostname => "foober.com"
	}
	)->status_is(200)->json_is( "/response" => "Successfully added ssl keys for $key" )
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

#validate keys were added
ok $t->get_ok("/api/1.1/deliveryservices/xmlId/$key/sslkeys.json")
	->json_has("/response")
	->json_has("/response/certificate/csr")
	->json_has("/response/certificate/key")
	->json_has("/response/certificate/crt")
	->json_is( "/response/hostname" => "foober.com" )
	->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

#NEGATIVE TESTING -- Alert handling
#key not found
#get_object by xmlId
my $fake_get_404 = HTTP::Response->new( 404, undef, HTTP::Headers->new, "Not found" );
$fake_lwp->mock( 'get', sub { return $fake_get_404 } );

ok $t->get_ok("/api/1.1/deliveryservices/xmlId/foo/sslkeys.json")->status_is(404)->json_has("A record for ssl key foo could not be found")
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

# TODO: Implement functionality to satisfy this test?
# #get_object by hostname, not a hostname
# ok $t->get_ok("/api/1.1/deliveryservices/hostname/foo/sslkeys.json")->status_is(400)->json_has("foo is not a valid hostname.")
# 	->or( sub { diag $t->tx->res->content->asset->{content}; } );

#get_object by hostname, ds not found
ok $t->get_ok("/api/1.1/deliveryservices/hostname/foo.fake-ds.cdn1.kabletown.net/sslkeys.json")->status_is(400)
	->json_has("A delivery service does not exist for a host with hostanme of foo.fake-ds.cdn1.kabletown.net")->or( sub { diag $t->tx->res->content->asset->{content}; } );

# OFFLINE all riak servers
my $rs = $schema->resultset('Server')->search( { type => 31 } );
$rs->update_all( { status => 3 } );

ok $t->post_ok(
	'/api/1.1/deliveryservices/sslkeys/generate',
	json => {
		key          => $key,
		version      => $version,
		hostname     => $hostname,
		country      => $country,
		state        => $state,
		city         => $city,
		organization => $org,
		businessUnit => $unit,
		cdn             => $cdn,
		deliveryservice => $deliveryservice,
	}
	)->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/alerts/0/level" => "error" )
	->json_like( "/alerts/0/text" => qr/^No RIAK servers/ ),

	'Creating ssl after riak servers are all offline should fail?';

# logout
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

done_testing();
