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

ok $t->post_ok('/api/1.2/servers/create' => {Accept => 'application/json'} => json => {
        "hostName" => "tc1_ats2",
        "domainName" => "northbound.com",
        "cachegroup" => "cg1-mid-northeast",
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
    ->json_is( "/response/hostName" => "tc1_ats2")
    ->json_is( "/response/domainName" => "northbound.com")
    ->json_is( "/response/cachegroup" => "cg1-mid-northeast")
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
        "cachegroup" => "cg5-edge_atl_group",
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
    ->json_is( "/response/cachegroup" => "cg5-edge_atl_group")
    ->json_is( "/response/ipNetmask" => "255.255.255.0")
    ->json_is( "/response/interfaceName" => "eth0")
    ->json_is( "/response/ipAddress" => "10.74.27.185")
    ->json_is( "/response/ipGateway" => "10.74.27.1")
    ->json_is( "/response/interfaceMtu" => "1500")
    ->json_is( "/response/physLocation" => "HotAtlanta")
    ->json_is( "/response/type" => "EDGE")
    ->json_is( "/response/profile" => "EDGE1")
            , 'Does the server details return?';

my $svr_id = &get_svr_id('tc1_ats1');

ok $t->put_ok('/api/1.2/servers/' . $svr_id . '/update'  => {Accept => 'application/json'} => json => {
        "hostName" => "tc1_ats3",
        "domainName" => "northbound.com",
        "cachegroup" => "cg5-edge_atl_group",
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
    ->json_is( "/response/cachegroup" => "cg5-edge_atl_group")
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

ok $t->get_ok('/api/1.2/servers?profileId=5' => {Accept => 'application/json'})->status_is(200)
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/id", 5 )
    ->json_is( "/response/1/id", 9 )
    ->json_is( "/response/2/id", 10 )
    ->json_is( "/response/3/id", 11 )
            , "Does the server ids return?";

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
