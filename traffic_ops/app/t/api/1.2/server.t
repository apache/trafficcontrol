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

ok $t->get_ok('/api/1.2/servers/details?hostName=atlanta-edge-01')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/0/ipGateway", "127.0.0.1" )->json_is( "/response/0/deliveryservices/0", "1" ), 'Does the hostname details return?';

ok $t->get_ok('/api/1.2/servers/details?physLocationID=1')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/0/ipGateway", "127.0.0.1" )->json_is( "/response/0/deliveryservices/0", "1" ), 'Does the physLocationID details return?';

ok $t->get_ok('/api/1.2/servers/details')->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Does the validation error occur?';
ok $t->get_ok('/api/1.2/servers/details.json?orderby=hostName')->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Does the orderby work?';

ok $t->get_ok('/api/1.2/servers?type=MID')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
  ->json_is( "/response/0/hostName", "atlanta-mid-01" )
  ->json_is( "/response/0/domainName", "ga.atlanta.kabletown.net" )
  ->json_is( "/response/0/type", "MID" )
  ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.2/servers?type=MID&status=ONLINE')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
  ->json_is( "/response/0/hostName", "atlanta-mid-01" )
  ->json_is( "/response/0/domainName", "ga.atlanta.kabletown.net" )
  ->json_is( "/response/0/type", "MID" )
  ->json_is( "/response/0/status", "ONLINE" )
  ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok('/api/1.2/servers/create' => {Accept => 'application/json'} => json => {
			"hostName" => "server1",
			"domainName" => "example-domain.com",
			"cachegroup" => "mid-northeast-group",
			"cdnName" => "cdn1",
			"ipAddress" => "10.74.27.194",
			"interfaceName" => "bond0",
			"ipNetmask" => "255.255.255.252",
			"ipGateway" => "10.74.27.194",
			"interfaceMtu" => "1500",
			"physLocation" => "Denver",
			"type" => "EDGE",
			"profile" => "EDGE1" })
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	, 'Is a server created when all required fields are provided?';

ok $t->post_ok('/api/1.2/servers/create' => {Accept => 'application/json'} => json => {
			"hostName" => "server2",
			"domainName" => "example-domain.com",
			"cachegroup" => "mid-northeast-group",
			"cdnName" => "cdn1",
			"ipAddress" => "10.74.27.194",
			"interfaceName" => "bond0",
			"ipNetmask" => "255.255.255.252",
			"ipGateway" => "10.74.27.194",
			"interfaceMtu" => "1500",
			"physLocation" => "Denver",
			"type" => "EDGE",
			"profile" => "EDGE1" })
		->status_is(400)
	, 'Does the server creation fail because ip address is already used for the profile?';

ok $t->post_ok('/api/1.2/servers/create' => {Accept => 'application/json'} => json => {
			"hostName" => "server3",
			"domainName" => "example-domain.com",
			"cachegroup" => "mid-northeast-group",
			"cdnName" => "cdn1",
			"ipAddress" => "10.74.27.85",
			"interfaceName" => "bond0",
			"ipNetmask" => "255.255.255.252",
			"ipGateway" => "10.74.27.85",
			"ip6Address" => "2001:852:fe0f:27::2/64",
			"ip6Gateway" => "2001:852:fe0f:27::1",
			"interfaceMtu" => "1500",
			"physLocation" => "Denver",
			"type" => "EDGE",
			"profile" => "EDGE1" })
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	, 'Is a server created when all required fields are provided plus an ip6 address?';

ok $t->post_ok('/api/1.2/servers/create' => {Accept => 'application/json'} => json => {
			"hostName" => "server3",
			"domainName" => "example-domain.com",
			"cachegroup" => "mid-northeast-group",
			"cdnName" => "cdn1",
			"ipAddress" => "10.74.27.77",
			"interfaceName" => "bond0",
			"ipNetmask" => "255.255.255.252",
			"ipGateway" => "10.74.27.77",
			"ip6Address" => "2001:852:fe0f:27::2/64",
			"ip6Gateway" => "2001:852:fe0f:27::1",
			"interfaceMtu" => "1500",
			"physLocation" => "Denver",
			"type" => "EDGE",
			"profile" => "EDGE1" })
		->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	, 'Does the server creation fail because ip6 address is already used for the profile?';

# Count the 'response number'
my $count_response = sub {
	my ( $t, $count ) = @_;
	my $json = decode_json( $t->tx->res->content->asset->slurp );
	my $r    = $json->{response};
	return $t->success( is( scalar(@$r), $count ) );
};

# this is a dns delivery service with 2 edges and 1 mid and since dns ds's DO employ mids, 3 servers return
$t->get_ok('/api/1.2/servers?dsId=5')->status_is(200)->$count_response(3)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

# this is a http_no_cache delivery service with 2 edges and 1 mid and since http_no_cache ds's DON'T employ mids, 2 servers return
$t->get_ok('/api/1.2/servers?dsId=6')->status_is(200)->$count_response(2)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();
