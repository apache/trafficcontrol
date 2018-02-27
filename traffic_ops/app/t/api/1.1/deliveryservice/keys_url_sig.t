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
use POSIX ();
use Mojo::Base -strict;
use Test::More;
use Test::Mojo;
use strict;
use warnings;
use Data::Dumper;
use warnings 'all';
use Test::TestHelper;
use Test::MockModule;
use Test::MockObject;
use JSON;
no warnings 'once';
use Schema;

#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

BEGIN { $ENV{MOJO_MODE} = "test" }
my $schema = Schema->connect_to_database;
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');

#load data so we can login
#unload data for a clean test
Test::TestHelper->unload_core_data($schema);

#load core test data
Test::TestHelper->load_core_data($schema);

my $url_sig_keys;
foreach my $i ( 0 .. 15 ) {
	my $v = "value$i";
	my $k = "key$i";
	$url_sig_keys->{$k} = $v;
	$i++;
}

my $fake_lwp = new Test::MockModule( 'LWP::UserAgent', no_auto => 1 );
my $fake_get_200 = HTTP::Response->new( 200, undef, HTTP::Headers->new, encode_json($url_sig_keys) );
$fake_lwp->mock( 'get', sub { return $fake_get_200 } );

my $fake_put_204 = HTTP::Response->new( 204, undef, HTTP::Headers->new, undef );
$fake_lwp->mock( 'put', sub { return $fake_put_204 } );

my $version = "1.1";

# Portal User checks
ok $t->post_ok( '/api/1.1/user/login', json => { u => Test::TestHelper::PORTAL_USER, p => Test::TestHelper::PORTAL_USER_PASSWORD } )->status_is(200),
	'Log into the portal user?';

ok $t->get_ok('/api/1.1/deliveryservices/xmlId/test-ds1/urlkeys.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Can assigned DeliveryService url keys can be viewed?';

ok $t->get_ok('/api/1.1/deliveryservices/xmlId/test-ds2/urlkeys.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Can unassigned DeliveryService url keys can be viewed?';

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# Admin User checks
ok $t->post_ok( '/api/1.1/user/login', json => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(200),
	'Log into the admin user?';

ok $t->post_ok('/api/1.1/deliveryservices/xmlId/XXX/urlkeys/generate')->status_is(400)
		->json_is( "/alerts/0/text/", "Delivery Service 'XXX' does not exist." )->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Can a non existent DeliveryService url keys for the portal user be regenerated?';

ok $t->post_ok('/api/1.1/deliveryservices/xmlId/test-ds1/urlkeys/generate')->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Can an assigned DeliveryService url keys for the portal user be regenerated?';

ok $t->post_ok('/api/1.1/deliveryservices/xmlId/test-ds2/urlkeys/generate')->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Can an unassigned DeliveryService url keys for the portal user be regenerated?';

ok $t->get_ok('/api/1.1/deliveryservices/xmlId/test-ds1/urlkeys.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'DeliveryService Url Keys can be viewed?';

ok $t->post_ok('/api/1.1/deliveryservices/xmlId/XXX/urlkeys/generate')->status_is(400)
	->json_is( "/alerts/0/text/", "Delivery Service 'XXX' does not exist." )->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Can a non-existent DeliveryService url keys for the portal user be regenerated?';

ok $t->get_ok('/api/1.1/deliveryservices/xmlId/test-ds1/urlkeys.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Can assigned DeliveryService url keys can be viewed?';

ok $t->get_ok('/api/1.1/deliveryservices/xmlId/test-ds2/urlkeys.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Can unassigned DeliveryService url keys can be viewed?';
#extract content of previous transaction's response for use verifying copy:
my $tx = $t->tx;
my $jsonKeys = $tx->res->json;

# Test copying of url_sig_keys
# api/$version/deliveryservices/xmlId/:xmlId/fromXmlId/:copyFromXmlId/urlkeys/copy
ok $t->post_ok('/api/1.1/deliveryservices/xmlId/test-ds1/urlkeys/copyFromXmlId/test-ds2')->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Can an unassigned DeliveryService url keys be copied to an assigned DeliveryService url keys?';

#compare contents of call below to stored response body from other ds.
ok $t->get_ok('/api/1.1/deliveryservices/xmlId/test-ds1/urlkeys.json')->status_is(200)->json_is($jsonKeys)->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Are the url sig keys equal after the copy?';


# Out of tenant tests
ok $t->post_ok('/api/1.1/deliveryservices/xmlId/test-ds1-root/urlkeys/generate')->status_is(403)
		->json_is( "/alerts/0/text" => "Forbidden. Delivery-service tenant is not available to the user.")
		->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Cannot generate delivery-service url keys when tenancy not allow?';

ok $t->get_ok('/api/1.1/deliveryservices/xmlId/test-ds1-root/urlkeys.json')->status_is(403)
		->json_is( "/alerts/0/text" => "Forbidden. Delivery-service tenant is not available to the user.")
		->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'DeliveryService Url Keys cannot be viewed out of tenancy?';

ok $t->get_ok('/api/1.1/deliveryservices/xmlId/test-ds1-not-there/urlkeys.json')->status_is(404)
		->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'DeliveryService Url Keys cannot be viewed out of tenancy?';

ok $t->post_ok('/api/1.1/deliveryservices/xmlId/test-ds1-root/urlkeys/copyFromXmlId/test-ds1')->status_is(403)
		->json_is( "/alerts/0/text" => "Forbidden. Delivery-service tenant is not available to the user.")
		->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Can an unassigned DeliveryService url keys be copied to an assigned DeliveryService url keys?';

ok $t->post_ok('/api/1.1/deliveryservices/xmlId/test-ds1/urlkeys/copyFromXmlId/test-ds1-root')->status_is(403)
		->json_is( "/alerts/0/text" => "Forbidden. Source delivery-service tenant is not available to the user.")
		->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Can an unassigned DeliveryService url keys be copied to an assigned DeliveryService url keys?';


# Negative Testing
# With error content
my $fake_put_300 = HTTP::Response->new( 300, undef, HTTP::Headers->new, "You messed it up!" );
$fake_lwp->mock( 'put', sub { return $fake_put_300 } );

ok $t->post_ok('/api/1.1/deliveryservices/xmlId/test-ds1/urlkeys/generate')->status_is(400)
	->json_is( "/alerts/0/text/", "You messed it up!" )->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Can a non-existent DeliveryService url keys for the portal user be regenerated?';

# OFFLINE all riak servers
my $rs = $schema->resultset('Server')->search( { type => 31 } );
$rs->update_all( { status => 1 } );

ok $t->post_ok('/api/1.1/deliveryservices/xmlId/test-ds1/urlkeys/generate')->status_is(400)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_like( "/alerts/0/text" => qr/^No RIAK servers/ ),
	'Can a non-existent DeliveryService url keys for the portal user be regenerated?';

=cut
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
	}
	)->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/alerts/0/level" => "error" )
	->json_like( "/alerts/0/text" => qr/^No Riak servers/ ),

	'Creating ssl after riak servers are all offline should fail?';
=cut

# logout
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
done_testing();
