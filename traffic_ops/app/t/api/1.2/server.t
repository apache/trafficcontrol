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
use Data::Dumper;

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

ok $t->get_ok('/api/1.2/servers/status')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/ONLINE", 16 )
		->json_is( "/response/REPORTED", 1 )
		->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok('/api/1.2/cachegroups' => {Accept => 'application/json'} => json => {
        "name" => "cg2-mid-northwest",
        "shortName" => "cg2_mid",
        "latitude" => 12,
        "longitude" => 56,
        "typeId" => 6 })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "cg2-mid-northwest" )
    ->json_is( "/response/shortName" => "cg2_mid")
    ->json_is( "/response/latitude" => 12)
    ->json_is( "/response/longitude" => 56)
    ->json_is( "/response/parentCachegroupId" => undef)
    ->json_is( "/response/parentCachegroupName" => undef)
    ->json_is( "/response/secondaryParentCachegroupId" => undef)
    ->json_is( "/response/secondaryParentCachegroupName" => undef)
            , 'Does the cache group details return?';

ok $t->post_ok('/api/1.2/cachegroups' => {Accept => 'application/json'} => json => {
        "name" => "cg-mid-northeast",
        "shortName" => "mneg",
        "latitude" => 10,
        "longitude" => 40,
        "typeId" => 6 })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "cg-mid-northeast" )
    ->json_is( "/response/shortName" => "mneg")
    ->json_is( "/response/latitude" => 10)
    ->json_is( "/response/longitude" => 40)
    ->json_is( "/response/parentCachegroupId" => undef)
    ->json_is( "/response/parentCachegroupName" => undef)
    ->json_is( "/response/secondaryParentCachegroupId" => undef)
    ->json_is( "/response/secondaryParentCachegroupName" => undef)
            , 'Does the cache group details return?';

ok $t->get_ok('/api/1.2/servers?type=MID')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
  ->content_like( "/atlanta\-mid\-01/" )
  ->content_like("/ga\.atlanta\.kabletown\.net/" )
  ->content_like("/MID/" )
  ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.2/servers?cdn=100')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
  ->content_like( "/atlanta\-edge\-01/" )
  ->content_like( "/ga\.atlanta\.kabletown\.net/" )
  ->content_like( "/100/" )
  ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.2/servers?cachegroup=200')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
  ->content_like( "/atlanta\-mid\-02/" )
  ->content_like( "/ga\.atlanta\.kabletown\.net/" )
  ->content_like( "/200/" )
  ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.2/servers?type=MID&status=ONLINE')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
  ->content_like(  "/atlanta\-mid\-01/" )
  ->content_like(  "/ga\.atlanta\.kabletown\.net/" )
  ->content_like(  "/MID/" )
  ->content_like(  "/ONLINE/" )
  ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok('/api/1.2/cachegroups' => {Accept => 'application/json'} => json => {
        "name" => "edge_atl_group1",
        "shortName" => "eag1",
        "latitude" => 22,
        "longitude" => 55,
        "typeId" => 6 })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "edge_atl_group1" )
    ->json_is( "/response/shortName" => "eag1")
    ->json_is( "/response/latitude" => 22)
    ->json_is( "/response/longitude" => 55)
    ->json_is( "/response/parentCachegroupId" => undef)
    ->json_is( "/response/parentCachegroupName" => undef)
    ->json_is( "/response/secondaryParentCachegroupId" => undef)
    ->json_is( "/response/secondaryParentCachegroupName" => undef)
            , 'Does the cache group details return?';

my $cg2_mid_northwest = &get_cg_id('cg2-mid-northwest');
ok $t->post_ok('/api/1.2/servers' => {Accept => 'application/json'} => json => {
			"cachegroupId" => $cg2_mid_northwest,
			"cdnId" => 100,
			"domainName" => "example-domain.com",
			"hostName" => "server1",
			"interfaceMtu" => 1500,
			"interfaceName" => "bond0",
			"ipAddress" => "10.74.27.194",
			"ipNetmask" => "255.255.255.252",
			"ipGateway" => "10.74.27.194",
			"physLocationId" => 100,
			"profileId" => 100,
			"statusId" => 1,
			"typeId" => 1,
			"updPending" => \0,
		})
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	, 'Is a server created when all required fields are provided?';

ok $t->post_ok('/api/1.2/servers' => {Accept => 'application/json'} => json => {
			"cachegroupId" => $cg2_mid_northwest,
			"cdnId" => 100,
			"domainName" => "example-domain.com",
			"hostName" => "server2",
			"httpsPort" => "string",
			"interfaceMtu" => 1500,
			"interfaceName" => "bond0",
			"ipAddress" => "10.74.27.255",
			"ipNetmask" => "255.255.255.252",
			"ipGateway" => "10.74.27.194",
			"physLocationId" => 100,
			"profileId" => 100,
			"statusId" => 1,
			"typeId" => 1,
			"updPending" => \0,
		})
		->status_is(400)
		->json_is( "/alerts/0/level", "error" )
		->json_is( "/alerts/0/text", "httpsPort must be an integer." )
	, "Does the server creation fail because httpsPort is a string?";

ok $t->post_ok('/api/1.2/servers' => {Accept => 'application/json'} => json => {
			"cachegroupId" => $cg2_mid_northwest,
			"cdnId" => 100,
			"domainName" => "example-domain.com",
			"hostName" => "server2",
			"interfaceMtu" => 1500,
			"interfaceName" => "bond0",
			"ipAddress" => "10.74.27.255",
			"ipNetmask" => "255.255.255.252",
			"ipGateway" => "10.74.27.194",
			"physLocationId" => 100,
			"profileId" => 100,
			"statusId" => 1,
			"tcpPort" => "string",
			"typeId" => 1,
			"updPending" => \0,
		})
		->status_is(400)
		->json_is( "/alerts/0/level", "error" )
		->json_is( "/alerts/0/text", "tcpPort must be an integer." )
	, "Does the server creation fail because tcpPort is a string?";

ok $t->post_ok('/api/1.2/servers' => {Accept => 'application/json'} => json => {
			"cachegroupId" => $cg2_mid_northwest,
			"cdnId" => 100,
			"domainName" => "example-domain.com",
			"hostName" => "server2",
			"interfaceMtu" => 1500,
			"interfaceName" => "bond0",
			"ipAddress" => "10.74.27.194",
			"ipNetmask" => "255.255.255.252",
			"ipGateway" => "10.74.27.194",
			"physLocationId" => 100,
			"profileId" => 100,
			"statusId" => 1,
			"typeId" => 1,
			"updPending" => \0,
		})
		->status_is(400)
	, 'Does the server creation fail because ip address is already used for the profile?';

ok $t->post_ok('/api/1.2/servers' => {Accept => 'application/json'} => json => {
			"cachegroupId" => $cg2_mid_northwest,
			"cdnId" => 100,
			"domainName" => "example-domain.com",
			"hostName" => "server2",
			"interfaceMtu" => 1500,
			"interfaceName" => "bond0",
			"ipAddress" => "10.74.27.85",
			"ipNetmask" => "255.255.255.252",
			"ipGateway" => "10.74.27.85",
			"ip6Address" => "2001:852:fe0f:27::2/64",
			"ip6Gateway" => "2001:852:fe0f:27::1",
			"physLocationId" => 100,
			"profileId" => 100,
			"statusId" => 1,
			"typeId" => 1,
			"updPending" => \0,
		})
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	, 'Is a server created when all required fields are provided plus an ip6 address?';

ok $t->post_ok('/api/1.2/servers' => {Accept => 'application/json'} => json => {
			"cachegroupId" => $cg2_mid_northwest,
			"cdnId" => 100,
			"domainName" => "example-domain.com",
			"hostName" => "server3",
			"interfaceMtu" => 1500,
			"interfaceName" => "bond0",
			"ipAddress" => "10.74.27.77",
			"ipNetmask" => "255.255.255.252",
			"ipGateway" => "10.74.27.77",
			"ip6Address" => "2001:852:fe0f:27::2/64",
			"ip6Gateway" => "2001:852:fe0f:27::1",
			"physLocationId" => 100,
			"profileId" => 100,
			"statusId" => 1,
			"typeId" => 1,
			"updPending" => \0,
		})
		->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	, 'Does the server creation fail because ip6 address is already used for the profile?';

my $mid_northeast_group = &get_cg_id('mid-northeast-group');
ok $t->post_ok('/api/1.2/servers' => {Accept => 'application/json'} => json => {
			"cachegroupId" => $mid_northeast_group,
			"cdnId" => 100,
			"domainName" => "northbound.com",
			"hostName" => "tc1_ats1",
			"interfaceMtu" => 1500,
			"interfaceName" => "eth0",
			"ipAddress" => "10.74.27.184",
			"ipNetmask" => "255.255.255.0",
			"ipGateway" => "10.74.27.1",
			"ip6Address" => "2001:852:fe0f:27::2/64",
			"ip6Gateway" => "2001:852:fe0f:27::1",
			"physLocationId" => 300,
			"profileId" => 200,
			"statusId" => 1,
			"typeId" => 2,
			"updPending" => \0,
		})
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/0/hostName" => "tc1_ats1")
		->json_is( "/response/0/domainName" => "northbound.com")
		->json_is( "/response/0/cachegroup" => "mid-northeast-group")
		->json_is( "/response/0/ipNetmask" => "255.255.255.0")
		->json_is( "/response/0/interfaceName" => "eth0")
		->json_is( "/response/0/ipAddress" => "10.74.27.184")
		->json_is( "/response/0/ipGateway" => "10.74.27.1")
		->json_is( "/response/0/interfaceMtu" => "1500")
		->json_is( "/response/0/physLocation" => "HotAtlanta")
		->json_is( "/response/0/type" => "MID")
		->json_is( "/response/0/profile" => "MID1")
	, 'Does the server details return?';

my $edge_atl_group = &get_cg_id('edge_atl_group');
ok $t->post_ok('/api/1.2/servers' => {Accept => 'application/json'} => json => {
			"cachegroupId" => $edge_atl_group,
			"cdnId" => 100,
			"domainName" => "northbound.com",
			"hostName" => "tc1_ats1",
			"interfaceMtu" => 1500,
			"interfaceName" => "eth0",
			"ipAddress" => "10.74.27.184",
			"ipNetmask" => "255.255.255.0",
			"ipGateway" => "10.74.27.1",
			"physLocationId" => 300,
			"profileId" => 100,
			"statusId" => 1,
			"typeId" => 1,
			"updPending" => \0,
		})
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/hostName" => "tc1_ats1")
    ->json_is( "/response/0/domainName" => "northbound.com")
    ->json_is( "/response/0/cachegroup" => "edge_atl_group")
    ->json_is( "/response/0/ipNetmask" => "255.255.255.0")
    ->json_is( "/response/0/interfaceName" => "eth0")
    ->json_is( "/response/0/ipAddress" => "10.74.27.184")
    ->json_is( "/response/0/ipGateway" => "10.74.27.1")
    ->json_is( "/response/0/interfaceMtu" => "1500")
    ->json_is( "/response/0/physLocation" => "HotAtlanta")
    ->json_is( "/response/0/type" => "EDGE")
    ->json_is( "/response/0/profile" => "EDGE1")
            , 'Does the server details return?';

ok $t->post_ok('/api/1.2/servers' => {Accept => 'application/json'} => json => {
			"cachegroupId" => $edge_atl_group,
			"cdnId" => 100,
			"domainName" => "northbound.com",
			"hostName" => "tc1_ats2",
			"interfaceMtu" => 1500,
			"interfaceName" => "eth0",
			"ipAddress" => "10.74.27.187",
			"ipNetmask" => "255.255.255.0",
			"ipGateway" => "10.74.27.1",
			"physLocationId" => 300,
			"profileId" => 100,
			"statusId" => 1,
			"typeId" => 1,
			"updPending" => \0,
		})
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/hostName" => "tc1_ats2")
    ->json_is( "/response/0/domainName" => "northbound.com")
    ->json_is( "/response/0/cachegroup" => "edge_atl_group")
    ->json_is( "/response/0/ipNetmask" => "255.255.255.0")
    ->json_is( "/response/0/interfaceName" => "eth0")
    ->json_is( "/response/0/ipAddress" => "10.74.27.187")
    ->json_is( "/response/0/ipGateway" => "10.74.27.1")
    ->json_is( "/response/0/interfaceMtu" => "1500")
    ->json_is( "/response/0/physLocation" => "HotAtlanta")
    ->json_is( "/response/0/type" => "EDGE")
    ->json_is( "/response/0/profile" => "EDGE1")
            , 'Does the server details return?';

ok $t->post_ok('/api/1.2/servers' => {Accept => 'application/json'} => json => {
			"cachegroupId" => $edge_atl_group,
			"cdnId" => 100,
			"domainName" => "northbound.com",
			"hostName" => "tc2_ats2",
			"httpsPort" => 443,
			"iloIpAddress" => "",
			"iloIpNetmask" => "",
			"iloIpGateway" => "",
			"iloUsername" => "",
			"iloPassword" => "",
			"interfaceMtu" => 1500,
			"interfaceName" => "eth0",
			"ip6Address" => "",
			"ip6Gateway" => "",
			"ipAddress" => "10.73.27.187",
			"ipNetmask" => "255.255.255.0",
			"ipGateway" => "10.73.27.1",
			"mgmtIpAddress" => "",
			"mgmtIpNetmask" => "",
			"mgmtIpGateway" => "",
			"offlineReason" => "",
			"physLocationId" => 300,
			"profileId" => 200,
			"rack" => "",
			"routerHostName" => "",
			"routerPortName" => "",
			"tcpPort" => 80,
			"statusId" => 1,
			"typeId" => 2,
			"updPending" => \0,
		})
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/0/cachegroup" => "edge_atl_group")
		->json_is( "/response/0/domainName" => "northbound.com")
		->json_is( "/response/0/hostName" => "tc2_ats2")
		->json_is( "/response/0/httpsPort" =>443)
		->json_is( "/response/0/iloIpAddress" => "")
		->json_is( "/response/0/iloIpNetmask" => "")
		->json_is( "/response/0/iloIpGateway" => "")
		->json_is( "/response/0/iloUsername" => "")
		->json_is( "/response/0/iloPassword" => "")
		->json_is( "/response/0/interfaceMtu" => "1500")
		->json_is( "/response/0/interfaceName" => "eth0")
		->json_is( "/response/0/ip6Address" => undef)
		->json_is( "/response/0/ip6Gateway" => "")
		->json_is( "/response/0/ipNetmask" => "255.255.255.0")
		->json_is( "/response/0/ipAddress" => "10.73.27.187")
		->json_is( "/response/0/ipGateway" => "10.73.27.1")
		->json_is( "/response/0/mgmtIpAddress" => "")
		->json_is( "/response/0/mgmtIpNetmask" => "")
		->json_is( "/response/0/mgmtIpGateway" => "")
		->json_is( "/response/0/offlineReason" => "")
		->json_is( "/response/0/physLocation" => "HotAtlanta")
		->json_is( "/response/0/type" => "MID")
		->json_is( "/response/0/profile" => "MID1")
		->json_is( "/response/0/rack" => "")
		->json_is( "/response/0/routerHostName" => "")
		->json_is( "/response/0/routerPortName" => "")
	, 'Is the server created?';

ok $t->get_ok('/api/1.2/servers/details')->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Does the validation error occur?';
ok $t->get_ok('/api/1.2/servers/details.json?orderby=hostName')->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Does the orderby work?';

ok $t->get_ok('/api/1.2/servers?type=MID')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
  ->json_is( "/response/0/type", "MID" )
  ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok('/api/1.2/deliveryservices/test-ds5/servers' => {Accept => 'application/json'} => json => { "serverNames" => [ 'tc1_ats1' ]})
     ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
     , 'Assign the server to the delivery service?';

ok $t->post_ok('/api/1.2/deliveryservices/test-ds4/servers' => {Accept => 'application/json'} => json => { "serverNames" => [ 'tc1_ats2', 'tc2_ats2' ]})
     ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
     , 'Assign the server to the delivery service?';

# BUG: last one in wins
ok $t->post_ok('/api/1.2/deliveryservices/test-ds4/servers' => {Accept => 'application/json'} => json => { "serverNames" => [ 'tc1_ats2', 'tc2_ats2' ]})
     ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
     , 'Assign the server to the delivery service?';

ok $t->get_ok('/api/1.2/servers?type=MID&status=ONLINE')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
  ->json_is( "/response/0/type", "MID" )
  ->json_is( "/response/0/status", "ONLINE" )
  ->or( sub { diag $t->tx->res->content->asset->{content}; } );

my $svr_id = &get_svr_id('tc1_ats1');
ok $t->put_ok('/api/1.2/servers/' . $svr_id  => {Accept => 'application/json'} => json => {
			"cachegroupId" => $edge_atl_group,
			"cdnId" => 100,
			"domainName" => "northbound.com",
			"hostName" => "tc1_ats3",
			"interfaceMtu" => 1500,
			"interfaceName" => "eth0",
			"ipAddress" => "10.74.27.186",
			"ipNetmask" => "255.255.255.0",
			"ipGateway" => "10.74.27.1",
			"physLocationId" => 100,
			"profileId" => 100,
			"statusId" => 1,
			"typeId" => 1,
			"updPending" => \0,
		})
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/hostName" => "tc1_ats3")
    ->json_is( "/response/0/domainName" => "northbound.com")
    ->json_is( "/response/0/cachegroup" => "edge_atl_group")
    ->json_is( "/response/0/ipNetmask" => "255.255.255.0")
    ->json_is( "/response/0/interfaceName" => "eth0")
    ->json_is( "/response/0/ipAddress" => "10.74.27.186")
    ->json_is( "/response/0/ipGateway" => "10.74.27.1")
    ->json_is( "/response/0/interfaceMtu" => "1500")
    ->json_is( "/response/0/physLocation" => "Denver")
    ->json_is( "/response/0/type" => "EDGE")
    ->json_is( "/response/0/profile" => "EDGE1")
            , 'Does the server details return?';

ok $t->put_ok('/api/1.2/servers/' . $svr_id  => {Accept => 'application/json'} => json => {
			"cachegroupId" => $edge_atl_group,
			"cdnId" => 100,
			"domainName" => "northbound.com",
			"hostName" => "tc1_ats3",
			"httpsPort" => 443,
			"iloIpAddress" => "",
			"iloIpNetmask" => "",
			"iloIpGateway" => "",
			"iloUsername" => "",
			"iloPassword" => "",
			"interfaceMtu" => 1500,
			"interfaceName" => "eth0",
			"ip6Address" => "",
			"ip6Gateway" => "",
			"ipAddress" => "10.74.27.186",
			"ipNetmask" => "255.255.255.0",
			"ipGateway" => "10.74.27.1",
			"mgmtIpAddress" => "",
			"mgmtIpNetmask" => "",
			"mgmtIpGateway" => "",
			"offlineReason" => "",
			"physLocationId" => 100,
			"profileId" => 100,
			"rack" => "",
			"routerHostName" => "",
			"routerPortName" => "",
			"tcpPort" => 80,
			"statusId" => 1,
			"typeId" => 1,
			"updPending" => \0,
		})
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/iloIpAddress" => "")
    ->json_is( "/response/0/iloIpNetmask" => "")
    ->json_is( "/response/0/iloIpGateway" => "")
    ->json_is( "/response/0/iloUsername" => "")
    ->json_is( "/response/0/iloPassword" => "")
    ->json_is( "/response/0/ip6Address" => undef)
    ->json_is( "/response/0/ip6Gateway" => "")
    ->json_is( "/response/0/mgmtIpAddress" => "")
    ->json_is( "/response/0/mgmtIpNetmask" => "")
    ->json_is( "/response/0/mgmtIpGateway" => "")
    ->json_is( "/response/0/offlineReason" => "")
    ->json_is( "/response/0/rack" => "")
    ->json_is( "/response/0/routerHostName" => "")
    ->json_is( "/response/0/routerPortName" => "")
            , 'Are empty strings allowed on a handful of fields?';

ok $t->put_ok('/api/1.2/servers/' . $svr_id => {Accept => 'application/json'} => json => {
        "ipAddress" => "10.10.10.220",
        "ipGateway" => "111.222.111.1",
        "ipNetmask" => "255.255.255.0" })
    ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
            , 'Does the server details return?';

ok $t->put_ok('/api/1.2/servers/' . $svr_id => {Accept => 'application/json'} => json => {
        "ip6Address" => "ee80::1",
        "ip6Gateway" => "fe80::1" })
    ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
            , 'Does the server details return?';

my $svr_id1 = &get_svr_id('tc1_ats3');
ok $t->post_ok('/api/1.2/servers/'. $svr_id1 . '/queue_update' =>  {Accept => 'application/json'} =>json => {
        'action' => 'queue' })
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/action" => "queue")
    ->json_is( "/response/serverId" => "".$svr_id1)
            , 'Does the queue_update api return?';

ok $t->post_ok('/api/1.2/servers/9999/queue_update' =>  {Accept => 'application/json'} =>json => {
        'action' => 'queue' })
    ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
            , 'Does the queue_update api return?';

ok $t->delete_ok('/api/1.2/servers/' . $svr_id)
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/alerts/0/level", "success" )
    ->json_is( "/alerts/0/text", "Server was deleted: tc1_ats3" )
            , "Is the server id valid?";

ok $t->delete_ok('/api/1.2/servers/' . $svr_id)
    ->status_is(404)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->put_ok('/api/1.2/servers/' . $svr_id => {Accept => 'application/json'} => json => {
			"cachegroupId" => $edge_atl_group,
			"cdnId" => 100,
			"domainName" => "northbound.com",
			"hostName" => "tc1_ats3",
			"interfaceMtu" => 1500,
			"interfaceName" => "eth0",
			"ipAddress" => "10.74.27.186",
			"ipNetmask" => "255.255.255.0",
			"ipGateway" => "10.74.27.1",
			"physLocationId" => 100,
			"profileId" => 100,
			"statusId" => 1,
			"typeId" => 1,
			"updPending" => \0,
		})
    ->status_is(404)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok('/api/1.2/servers' => {Accept => 'application/json'} => json => {
			"hostName" => "my-server-host",
			"domainName" => "example-domain.com",
			"cachegroupId" => 100,
			"cdnId" => 100,
			"ipAddress" => "10.74.27.78",
			"interfaceName" => "bond0",
			"ipNetmask" => "255.255.255.252",
			"ipGateway" => "10.74.27.78",
			"interfaceMtu" => 1500,
			"physLocationId" => 100,
			"typeId" => 1,
			"statusId" => 1,
			"updPending" => \0,
			"profileId" => 100 })
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	, 'Does the server creation succeed?';

ok $t->post_ok('/api/1.2/servers' => {Accept => 'application/json'} => json => {
			"hostName" => "my-server-host",
			"domainName" => "example-domain.com",
			"cachegroupId" => 100,
			"cdnId" => 100,
			"ipAddress" => "10.74.27.78",
			"interfaceName" => "bond0",
			"ipNetmask" => "255.255.255.252",
			"ipGateway" => "10.74.27.78",
			"interfaceMtu" => 1500,
			"physLocationId" => 100,
			"typeId" => 1,
			"statusId" => 1,
			"updPending" => \0,
			"profileId" => 100 })
		->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	, 'Does the server creation fail because ipAddress is already used by the profile?';

my $server_id =&get_svr_id('my-server-host');
ok $t->put_ok('/api/1.2/servers/' . $server_id => {Accept => 'application/json'} => json => {
			"hostName" => "my-server-host",
			"domainName" => "example-domain.com",
			"cachegroupId" => 200,
			"cdnId" => 100,
			"ipAddress" => "10.74.27.78",
			"interfaceName" => "bond0",
			"ipNetmask" => "255.255.255.252",
			"ipGateway" => "10.74.27.78",
			"interfaceMtu" => 1500,
			"physLocationId" => 100,
			"typeId" => 1,
			"statusId" => 1,
			"updPending" => \0,
			"profileId" => 100 })
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	, 'Does the server update succeed because ipAddress is already used by the profile but...by this server?';

ok $t->post_ok('/api/1.2/servers' => {Accept => 'application/json'} => json => {
			"hostName" => "my-server-host-ip6",
			"domainName" => "example-domain.com",
			"cachegroupId" => 100,
			"cdnId" => 100,
			"ipAddress" => "10.74.27.79",
			"interfaceName" => "bond0",
			"ipNetmask" => "255.255.255.252",
			"ipGateway" => "10.74.27.79",
			"ip6Address" => "2001:853:fe0f:27::2/64",
			"ip6Gateway" => "2001:853:fe0f:27::1",
			"interfaceMtu" => 1500,
			"physLocationId" => 100,
			"typeId" => 1,
			"statusId" => 1,
			"updPending" => \0,
			"profileId" => 100 })
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	, 'Does the server creation succeed?';

ok $t->post_ok('/api/1.2/servers' => {Accept => 'application/json'} => json => {
			"hostName" => "my-server-host-ip6",
			"domainName" => "example-domain.com",
			"cachegroupId" => 100,
			"cdnId" => 100,
			"ipAddress" => "10.74.27.80",
			"interfaceName" => "bond0",
			"ipNetmask" => "255.255.255.252",
			"ipGateway" => "10.74.27.80",
			"ip6Address" => "2001:853:fe0f:27::2/64",
			"ip6Gateway" => "2001:853:fe0f:27::1",
			"interfaceMtu" => 1500,
			"physLocationId" => 100,
			"typeId" => 1,
			"statusId" => 1,
			"updPending" => \0,
			"profileId" => 100 })
		->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	, 'Does the server creation fail because ip6Address is already used by the profile?';

my $server_id =&get_svr_id('my-server-host-ip6');
ok $t->put_ok('/api/1.2/servers/' . $server_id => {Accept => 'application/json'} => json => {
			"hostName" => "my-server-host-ip6",
			"domainName" => "example-domain.com",
			"cachegroupId" => 200,
			"cdnId" => 100,
			"ipAddress" => "10.74.27.80",
			"interfaceName" => "bond0",
			"ipNetmask" => "255.255.255.252",
			"ipGateway" => "10.74.27.80",
			"ip6Address" => "2001:853:fe0f:27::2/64",
			"ip6Gateway" => "2001:853:fe0f:27::1",
			"interfaceMtu" => 1500,
			"physLocationId" => 100,
			"typeId" => 1,
			"statusId" => 1,
			"updPending" => \0,
			"profileId" => 100 })
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	, 'Does the server update succeed because ip6Address is already used by the profile but...by this server?';

ok $t->put_ok('/api/1.2/servers/' . $server_id . '/status' => {Accept => 'application/json'} => json => {
			"status" => 'CARROT' })
		->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/alerts/0/level", "error" )
		->json_is( "/alerts/0/text", "Invalid status." )
	, 'Does the server status update fail because the status is invalid?';

ok $t->put_ok('/api/1.2/servers/' . $server_id . '/status' => {Accept => 'application/json'} => json => {
			"status" => 'OFFLINE' })
		->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/alerts/0/level", "error" )
		->json_is( "/alerts/0/text", "Offline reason is required for ADMIN_DOWN or OFFLINE status." )
	, 'Does the server status update fail because offline reason was not provided?';

ok $t->put_ok('/api/1.2/servers/' . $server_id . '/status' => {Accept => 'application/json'} => json => {
			"status" => 1,
			"offlineReason" => "taco tuesday" })
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/alerts/0/level", "success" )
	, 'Does the server status update succeed with status ID?';

ok $t->put_ok('/api/1.2/servers/' . $server_id . '/status' => {Accept => 'application/json'} => json => {
			"status" => "OFFLINE",
			"offlineReason" => "wacky wednesday" })
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/alerts/0/level", "success" )
	, 'Does the server status update succeed with status name?';

ok $t->put_ok('/api/1.2/servers/' . $server_id . '/status' => {Accept => 'application/json'} => json => {
			"status" => "ONLINE" })
		->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/alerts/0/level", "success" )
		->json_is( "/alerts/0/text", "Updated status [ ONLINE ] for my-server-host-ip6.example-domain.com [  ] and queued updates on all child caches" )
	, 'Does the server status update succeed and updates are queued when the status is changed on an Edge server?';

ok $t->get_ok('/api/1.2/servers/status')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/ONLINE", 17 )
		->json_is( "/response/REPORTED", 1 )
		->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();

sub get_svr_id {
	my $host_name = shift;
	my $q      = "select id from server where host_name = \'$host_name\'";
	my $get_svr = $dbh->prepare($q);
	$get_svr->execute();
	my $p = $get_svr->fetchall_arrayref( {} );
	$get_svr->finish();
	my $id = $p->[0]->{id};
	return $id;
}

sub get_cg_id {
	my $cg_name = shift;
	my $q      = "select id from cachegroup where name = \'$cg_name\'";
	my $get_cg = $dbh->prepare($q);
	$get_cg->execute();
	my $p = $get_cg->fetchall_arrayref( {} );
	$get_cg->finish();
	my $id = $p->[0]->{id};
	return $id;
}

