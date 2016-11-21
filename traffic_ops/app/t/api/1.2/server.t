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
my $schema_values = { schema => $schema, no_transactions => 1 };
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');

Test::TestHelper->unload_core_data($schema);
Test::TestHelper->load_all_fixtures( Fixtures::Cdn->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::Role->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::TmUser->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::Status->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::Parameter->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::Profile->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::ProfileParameter->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::Division->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::Region->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::PhysLocation->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::Type->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::Deliveryservice->new($schema_values) );

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


ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';

ok $t->post_ok('/api/1.2/cachegroups/create' => {Accept => 'application/json'} => json => {
        "name" => "cg2-mid-northwest",
        "shortName" => "cg2_mid",
        "latitude" => "12",
        "longitude" => "56",
        "parentCachegroup" => "",
        "secondaryParentCachegroup" => "",
        "typeName" => "MID_LOC" })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "cg2-mid-northwest" )
    ->json_is( "/response/shortName" => "cg2_mid")
    ->json_is( "/response/latitude" => "12")
    ->json_is( "/response/longitude" => "56")
    ->json_is( "/response/parentCachegroup" => "")
    ->json_is( "/response/secondaryParentCachegroup" => "")
            , 'Does the cache group details return?';

ok $t->post_ok('/api/1.2/cachegroups/create' => {Accept => 'application/json'} => json => {
        "name" => "mid-northeast-group",
        "shortName" => "mneg",
        "latitude" => "10",
        "longitude" => "40",
        "parentCachegroup" => "",
        "secondaryParentCachegroup" => "",
        "typeName" => "MID_LOC" })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "mid-northeast-group" )
    ->json_is( "/response/shortName" => "mneg")
    ->json_is( "/response/latitude" => "10")
    ->json_is( "/response/longitude" => "40")
    ->json_is( "/response/parentCachegroup" => "")
    ->json_is( "/response/secondaryParentCachegroup" => "")
            , 'Does the cache group details return?';

ok $t->get_ok('/api/1.2/servers?type=MID')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
  ->json_is( "/response/0/hostName", "atlanta-mid-01" )
  ->json_is( "/response/0/domainName", "ga.atlanta.kabletown.net" )
  ->json_is( "/response/0/type", "MID" )
  ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.2/servers?cdn=1')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
  ->json_is( "/response/0/hostName", "atlanta-edge-01" )
  ->json_is( "/response/0/domainName", "ga.atlanta.kabletown.net" )
  ->json_is( "/response/0/cdnId", 1 )
  ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.2/servers?cachegroup=2')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
  ->json_is( "/response/0/hostName", "atlanta-mid-02" )
  ->json_is( "/response/0/domainName", "ga.atlanta.kabletown.net" )
  ->json_is( "/response/0/cachegroupId", 2 )
  ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.2/servers?type=MID&status=ONLINE')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
  ->json_is( "/response/0/hostName", "atlanta-mid-01" )
  ->json_is( "/response/0/domainName", "ga.atlanta.kabletown.net" )
  ->json_is( "/response/0/type", "MID" )
  ->json_is( "/response/0/status", "ONLINE" )
  ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok('/api/1.2/cachegroups/create' => {Accept => 'application/json'} => json => {
        "name" => "edge_atl_group",
        "shortName" => "eag",
        "latitude" => "22",
        "longitude" => "55",
        "parentCachegroup" => "",
        "secondaryParentCachegroup" => "",
        "typeName" => "MID_LOC" })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "edge_atl_group" )
    ->json_is( "/response/shortName" => "eag")
    ->json_is( "/response/latitude" => "22")
    ->json_is( "/response/longitude" => "55")
    ->json_is( "/response/parentCachegroup" => "")
    ->json_is( "/response/secondaryParentCachegroup" => "")
            , 'Does the cache group details return?';

ok $t->post_ok('/api/1.2/servers/create' => {Accept => 'application/json'} => json => {
			"hostName" => "server1",
			"domainName" => "example-domain.com",
			"cachegroup" => "cg2-mid-northwest",
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
			"cachegroup" => "cg2-mid-northwest",
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
			"cachegroup" => "cg2-mid-northwest",
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
			"cachegroup" => "cg2-mid-northwest",
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

ok $t->post_ok('/api/1.2/servers/create' => {Accept => 'application/json'} => json => {
        "hostName" => "tc1_ats1",
        "domainName" => "northbound.com",
        "cachegroup" => "mid-northeast-group",
        "cdnName" => "cdn1",
        "interfaceName" => "eth0",
        "ipAddress" => "10.74.27.184",
        "ipNetmask" => "255.255.255.0",
        "ipGateway" => "10.74.27.1",
        "interfaceMtu" => "1500",
        "physLocation" => "HotAtlanta",
        "type" => "MID",
        "profile" => "MID1" })
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/hostName" => "tc1_ats1")
    ->json_is( "/response/domainName" => "northbound.com")
    ->json_is( "/response/cachegroup" => "mid-northeast-group")
    ->json_is( "/response/ipNetmask" => "255.255.255.0")
    ->json_is( "/response/interfaceName" => "eth0")
    ->json_is( "/response/ipAddress" => "10.74.27.184")
    ->json_is( "/response/ipGateway" => "10.74.27.1")
    ->json_is( "/response/interfaceMtu" => "1500")
    ->json_is( "/response/physLocation" => "HotAtlanta")
    ->json_is( "/response/type" => "MID")
    ->json_is( "/response/profile" => "MID1")
            , 'Does the server details return?';

ok $t->post_ok('/api/1.2/servers/create' => {Accept => 'application/json'} => json => {
        "hostName" => "tc1_ats1",
        "domainName" => "northbound.com",
        "cachegroup" => "edge_atl_group",
        "cdnName" => "cdn1",
        "interfaceName" => "eth0",
        "ipAddress" => "10.74.27.185",
        "ipNetmask" => "255.255.255.0",
        "ipGateway" => "10.74.27.1",
        "interfaceMtu" => "1500",
        "physLocation" => "HotAtlanta",
        "type" => "EDGE",
        "profile" => "EDGE1" })
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/hostName" => "tc1_ats1")
    ->json_is( "/response/domainName" => "northbound.com")
    ->json_is( "/response/cachegroup" => "edge_atl_group")
    ->json_is( "/response/ipNetmask" => "255.255.255.0")
    ->json_is( "/response/interfaceName" => "eth0")
    ->json_is( "/response/ipAddress" => "10.74.27.185")
    ->json_is( "/response/ipGateway" => "10.74.27.1")
    ->json_is( "/response/interfaceMtu" => "1500")
    ->json_is( "/response/physLocation" => "HotAtlanta")
    ->json_is( "/response/type" => "EDGE")
    ->json_is( "/response/profile" => "EDGE1")
            , 'Does the server details return?';

ok $t->post_ok('/api/1.2/servers/create' => {Accept => 'application/json'} => json => {
        "hostName" => "tc1_ats2",
        "domainName" => "northbound.com",
        "cachegroup" => "edge_atl_group",
        "cdnName" => "cdn1",
        "interfaceName" => "eth0",
        "ipAddress" => "10.74.27.187",
        "ipNetmask" => "255.255.255.0",
        "ipGateway" => "10.74.27.1",
        "interfaceMtu" => "1500",
        "physLocation" => "HotAtlanta",
        "type" => "EDGE",
        "profile" => "EDGE1" })
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/hostName" => "tc1_ats2")
    ->json_is( "/response/domainName" => "northbound.com")
    ->json_is( "/response/cachegroup" => "edge_atl_group")
    ->json_is( "/response/ipNetmask" => "255.255.255.0")
    ->json_is( "/response/interfaceName" => "eth0")
    ->json_is( "/response/ipAddress" => "10.74.27.187")
    ->json_is( "/response/ipGateway" => "10.74.27.1")
    ->json_is( "/response/interfaceMtu" => "1500")
    ->json_is( "/response/physLocation" => "HotAtlanta")
    ->json_is( "/response/type" => "EDGE")
    ->json_is( "/response/profile" => "EDGE1")
            , 'Does the server details return?';

ok $t->post_ok('/api/1.2/servers/create' => {Accept => 'application/json'} => json => {
        "hostName" => "tc2_ats2",
        "domainName" => "northbound.com",
        "cachegroup" => "edge_atl_group",
        "cdnName" => "cdn1",
        "interfaceName" => "eth0",
        "ipAddress" => "10.73.27.187",
        "ipNetmask" => "255.255.255.0",
        "ipGateway" => "10.73.27.1",
        "interfaceMtu" => "1500",
        "physLocation" => "HotAtlanta",
        "type" => "MID",
        "profile" => "MID1" })
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/hostName" => "tc2_ats2")
    ->json_is( "/response/domainName" => "northbound.com")
    ->json_is( "/response/cachegroup" => "edge_atl_group")
    ->json_is( "/response/ipNetmask" => "255.255.255.0")
    ->json_is( "/response/interfaceName" => "eth0")
    ->json_is( "/response/ipAddress" => "10.73.27.187")
    ->json_is( "/response/ipGateway" => "10.73.27.1")
    ->json_is( "/response/interfaceMtu" => "1500")
    ->json_is( "/response/physLocation" => "HotAtlanta")
    ->json_is( "/response/type" => "MID")
    ->json_is( "/response/profile" => "MID1")
            , 'Does the server details return?';



ok $t->post_ok('/api/1.2/deliveryservices/test-ds1/servers' => {Accept => 'application/json'} => json => { "serverNames" => [ 'server1', 'server3' ]})
     ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
     , 'Assign the server to the delivery service?';


ok $t->get_ok('/api/1.2/servers/details.json?hostName=server1')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/0/ipGateway", "10.74.27.194" )->json_is( "/response/0/deliveryservices/0", "100" ), 'Does the hostname details return?';

ok $t->get_ok('/api/1.2/servers/details.json?physLocationID=100')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/0/ipGateway", "10.74.27.194" )->json_is( "/response/0/deliveryservices/0", "100" ), 'Does the physLocationID details return?';

ok $t->get_ok('/api/1.2/servers/details')->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Does the validation error occur?';
ok $t->get_ok('/api/1.2/servers/details.json?orderby=hostName')->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Does the orderby work?';

ok $t->get_ok('/api/1.2/servers?type=MID')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
  ->json_is( "/response/0/hostName", "tc1_ats1" )
  ->json_is( "/response/0/domainName", "northbound.com" )
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
  ->json_is( "/response/0/hostName", "tc1_ats1" )
  ->json_is( "/response/0/domainName", "northbound.com" )
  ->json_is( "/response/0/type", "MID" )
  ->json_is( "/response/0/status", "ONLINE" )
  ->or( sub { diag $t->tx->res->content->asset->{content}; } );

# Count the 'response number'
my $count_response = sub {
	my ( $t, $count ) = @_;
	my $json = decode_json( $t->tx->res->content->asset->slurp );
	my $r    = $json->{response};
	return $t->success( is( scalar(@$r), $count ) );
};

# this is a dns delivery service with 2 edges and 1 mid and since dns ds's DO employ mids, 3 servers return
$t->get_ok('/api/1.2/servers?dsId=100')->status_is(200)->$count_response(2)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

# this is a http_no_cache delivery service with 2 edges and 1 mid and since http_no_cache ds's DON'T employ mids, 2 servers return
$t->get_ok('/api/1.2/servers?dsId=400')->status_is(200)->$count_response(2)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

my $svr_id = &get_svr_id('tc1_ats1');

ok $t->put_ok('/api/1.2/servers/' . $svr_id . '/update'  => {Accept => 'application/json'} => json => {
        "hostName" => "tc1_ats3",
        "domainName" => "northbound.com",
        "cachegroup" => "edge_atl_group",
        "cdnName" => "cdn1",
        "interfaceName" => "eth0",
        "ipAddress" => "10.74.27.186",
        "ipNetmask" => "255.255.255.0",
        "ipGateway" => "10.74.27.1",
        "interfaceMtu" => "1500",
        "physLocation" => "Denver",
        "type" => "EDGE",
        "profile" => "EDGE1" })
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/hostName" => "tc1_ats3")
    ->json_is( "/response/domainName" => "northbound.com")
    ->json_is( "/response/cachegroup" => "edge_atl_group")
    ->json_is( "/response/ipNetmask" => "255.255.255.0")
    ->json_is( "/response/interfaceName" => "eth0")
    ->json_is( "/response/ipAddress" => "10.74.27.186")
    ->json_is( "/response/ipGateway" => "10.74.27.1")
    ->json_is( "/response/interfaceMtu" => "1500")
    ->json_is( "/response/physLocation" => "Denver")
    ->json_is( "/response/type" => "EDGE")
    ->json_is( "/response/profile" => "EDGE1")
            , 'Does the server details return?';

ok $t->put_ok('/api/1.2/servers/' . $svr_id . '/update'  => {Accept => 'application/json'} => json => {
        "ipAddress" => "10.10.10.220",
        "ipGateway" => "111.222.111.1",
        "ipNetmask" => "255.255.255.0" })
    ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
            , 'Does the server details return?';

ok $t->put_ok('/api/1.2/servers/' . $svr_id . '/update'  => {Accept => 'application/json'} => json => {
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

ok $t->put_ok('/api/1.2/servers/' . $svr_id . '/update'  => {Accept => 'application/json'} => json => {
        "hostName" => "tc1_ats1",
        "domainName" => "northbound.com",
        "ipAddress" => "10.74.27.185",
        "physLocation" => "HotAtlanta" })
    ->status_is(404)->or( sub { diag $t->tx->res->content->asset->{content}; } );

$svr_id1 = &get_svr_id('server1');
my $svr_id2 = &get_svr_id('server3');
my $svr_id3 = &get_svr_id('tc1_ats1');
my $svr_id4 = &get_svr_id('tc1_ats2');
ok $t->get_ok('/api/1.2/servers?profileId=100&orderby=id' => {Accept => 'application/json'})->status_is(200)
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/id", $svr_id1 )
    ->json_is( "/response/1/id", $svr_id2 )
    ->json_is( "/response/2/id", $svr_id3 )
    ->json_is( "/response/3/id", $svr_id4 )
            , "Does the server ids return?";

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();
