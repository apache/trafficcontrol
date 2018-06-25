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
my $schema_values = { schema => $schema, no_transactions => 1 };
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');

Test::TestHelper->unload_core_data($schema);

# Load the test data up until 'cachegroup', because this test case creates
# them.
Test::TestHelper->load_all_fixtures( Fixtures::Tenant->new($schema_values) );
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

ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';

ok $t->post_ok('/api/1.2/cachegroups' => {Accept => 'application/json'} => json => {
            "name" => "cache_group_mid",
            "shortName" => "cg_mid",
            "latitude" => 12,
            "longitude" => 56,
            "typeId" => 6,
            "parentCachegroupId" => "",
        })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/level", "error" )
        ->json_is( "/alerts/0/text", "parentCachegroupId invalid. Must be a positive integer or null.")
    , 'Does the cache group create fail?';

ok $t->post_ok('/api/1.2/cachegroups' => {Accept => 'application/json'} => json => {
            "name" => "cache_group_mid",
            "shortName" => "cg_mid",
            "latitude" => 12,
            "longitude" => 56,
            "typeId" => 6,
            "secondaryParentCachegroupId" => "",
        })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/level", "error" )
        ->json_is( "/alerts/0/text", "secondaryParentCachegroupId invalid. Must be a positive integer or null.")
    , 'Does the cache group create fail?';

ok $t->post_ok('/api/1.2/cachegroups' => {Accept => 'application/json'} => json => {
        "name" => "cache_group_mid",
        "shortName" => "cg_mid",
        "latitude" => 12,
        "longitude" => 56,
        "typeId" => 6 })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "cache_group_mid" )
    ->json_is( "/response/shortName" => "cg_mid")
    ->json_is( "/response/latitude" => "12")
    ->json_is( "/response/longitude" => "56")
    ->json_is( "/response/parentCachegroupId" => undef)
    ->json_is( "/response/secondaryParentCachegroupId" => undef)
            , 'Is the cache group successfully created?';

ok $t->post_ok('/api/1.2/cachegroups' => {Accept => 'application/json'} => json => {
        "name" => "mid-northeast-group",
        "shortName" => "mid-ne-group",
        "latitude" => 44,
        "longitude" => 66,
        "typeId" => 6 })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "mid-northeast-group" )
    ->json_is( "/response/shortName" => "mid-ne-group")
    ->json_is( "/response/latitude" => 44)
    ->json_is( "/response/longitude" => 66)
    ->json_is( "/response/parentCachegroupId" => undef)
    ->json_is( "/response/secondaryParentCachegroupId" => undef)
            , 'Is the cache group successfully created?';

my $cache_group_mid_id = &get_cg_id('cache_group_mid');
my $mid_northeast_group_id = &get_cg_id('mid-northeast-group');
ok $t->post_ok('/api/1.2/cachegroups' => {Accept => 'application/json'} => json => {
        "name" => "cache_group_edge",
        "shortName" => "cg_edge",
        "latitude" => 12,
        "longitude" => 56,
        "parentCachegroupId" => $cache_group_mid_id,
        "secondaryParentCachegroupId" => $mid_northeast_group_id,
        "typeId" => 5 })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "cache_group_edge" )
    ->json_is( "/response/shortName" => "cg_edge")
    ->json_is( "/response/latitude" => 12)
    ->json_is( "/response/longitude" => 56)
    ->json_is( "/response/parentCachegroupId" => $cache_group_mid_id)
    ->json_is( "/response/parentCachegroupName" => "cache_group_mid")
    ->json_is( "/response/secondaryParentCachegroupId" => $mid_northeast_group_id)
    ->json_is( "/response/secondaryParentCachegroupName" => "mid-northeast-group")
            , 'Does the cache group details return?';

ok $t->post_ok('/api/1.2/cachegroups' => {Accept => 'application/json'} => json => {
        "name" => "cache_group_edge1",
        "shortName" => "cg_edge1",
        "latitude" => 23,
        "longitude" => 45,
        "typeId" => 5 })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "cache_group_edge1" )
    ->json_is( "/response/shortName" => "cg_edge1")
    ->json_is( "/response/latitude" => 23)
    ->json_is( "/response/longitude" => 45)
    ->json_is( "/response/parentCachegroupId" => undef)
    ->json_is( "/response/secondaryParentCachegroupId" => undef)
            , 'Does the cache group details return?';

my $mid_northeast_group = &get_cg_id('mid-northeast-group');
ok $t->post_ok('/api/1.2/servers' => {Accept => 'application/json'} => json => {
            "cachegroupId" => $mid_northeast_group,
            "cdnId" => 100,
            "domainName" => "northbound.com",
            "hostName" => "tc1_ats2",
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
    ->json_is( "/response/0/hostName" => "tc1_ats2")
            , 'Does the server details return?';

my $necg_id = &get_cg_id('mid-northeast-group');
ok $t->post_ok('/api/1.2/cachegroups/'. $necg_id .'/queue_update' =>  {Accept => 'application/json'} =>json => {
        'action' => 'queue',
        'cdn' => 'cdn1'})
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/action" => "queue")
    ->json_is( "/response/cdn" => "cdn1")
    ->json_is( "/response/cachegroupName" => "mid-northeast-group")
            , 'Does the queue_update api return?';

ok $t->post_ok('/api/1.2/cachegroups/'. $necg_id .'/queue_update' =>  {Accept => 'application/json'} =>json => {
        'action' => 'dequeue',
        'cdn' => 'cdn1'})
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/action" => "dequeue")
    ->json_is( "/response/cachegroupName" => "mid-northeast-group")
            , 'Does the queue_update api return?';

ok $t->post_ok('/api/1.2/cachegroups/'. $necg_id .'/queue_update' =>  {Accept => 'application/json'} =>json => {
        'action' => 'queue',
        'cdn' => 'cdn'})
    ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
            , 'Does the queueupdate api return?';
ok $t->post_ok('/api/1.2/cachegroups/9999/queue_update' =>  {Accept => 'application/json'} =>json => {
        'action' => 'queue',
        'cdn' => 'cdn1'})
    ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
            , 'Does the queue_update api return?';

my $cg_id = &get_cg_id('cache_group_edge');
ok $t->put_ok('/api/1.2/cachegroups/' . $cg_id => {Accept => 'application/json'} => json => {
            "name" => "cache_group_edge_1",
            "shortName" => "cg_edge_1",
            "latitude" => 23,
            "longitude" => 56,
            "typeId" => 5 })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/response/name" => "cache_group_edge_1" )
        ->json_is( "/response/shortName" => "cg_edge_1")
        ->json_is( "/response/latitude" => 23)
        ->json_is( "/response/longitude" => 56)
        ->json_is( "/response/parentCachegroupId" => undef)
        ->json_is( "/response/parentCachegroupName" => undef)
        ->json_is( "/response/secondaryParentCachegroupName" => undef)
        ->json_is( "/response/secondaryParentCachegroupName" => undef)
    , 'Does the cache group details return?';

ok $t->put_ok('/api/1.2/cachegroups/' . $cg_id => {Accept => 'application/json'} => json => {
        "name" => "cache_group_edge_2",
        "shortName" => "cg_edge_2",
        "latitude" => 0,
        "longitude" => 0,
        "parentCachegroupId" => undef,
        "secondaryParentCachegroupId" => undef,
        "typeId" => 5 })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "cache_group_edge_2" )
    ->json_is( "/response/shortName" => "cg_edge_2")
    ->json_is( "/response/latitude" => 0)
    ->json_is( "/response/longitude" => 0)
    ->json_is( "/response/parentCachegroupId" => undef)
    ->json_is( "/response/parentCachegroupName" => undef)
    ->json_is( "/response/secondaryParentCachegroupId" => undef)
    ->json_is( "/response/secondaryParentCachegroupName" => undef)
            , 'Does the cache group details return?';

ok $t->put_ok('/api/1.2/cachegroups/' . $cg_id => {Accept => 'application/json'} => json => {
        "name" => "cache_group_edge_2",
        "shortName" => "cg_edge_2",
        "latitude" => 0,
        "longitude" => 0,
        "parentCachegroupId" => $cache_group_mid_id,
        "typeId" => 5})->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/parentCachegroupId" => $cache_group_mid_id)
    ->json_is( "/response/parentCachegroupName" => "cache_group_mid")
            , 'Does the cache group details return?';

ok $t->put_ok('/api/1.2/cachegroups/' . $cg_id => {Accept => 'application/json'} => json => {
        "name" => "cache_group_edge_1",
        "typeId" => 5})->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->put_ok('/api/1.2/cachegroups/' . $cg_id => {Accept => 'application/json'} => json => {
        "name" => "cache_group_edge_1",
        "parentCachegroup" => "cache_group_edge_2",
        "typeId" => 5})->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } );


my $cache_group_edge_2_id = &get_cg_id('cache_group_edge_2');
ok $t->post_ok('/api/1.2/servers' => {Accept => 'application/json'} => json => {
            "cachegroupId" => $cache_group_edge_2_id,
            "cdnId" => 100,
            "domainName" => "test.example.com",
            "hostName" => "edge_streamer_1",
            "interfaceMtu" => 1500,
            "interfaceName" => "eth0",
            "ipAddress" => "192.168.100.2",
            "ipNetmask" => "255.255.255.0",
            "ipGateway" => "192.168.100.1",
            "physLocationId" => 300,
            "profileId" => 100,
            "statusId" => 1,
            "typeId" => 1,
            "updPending" => \0,
        })
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/hostName" => "edge_streamer_1")
            , 'Does the server details return?';

ok $t->delete_ok('/api/1.2/cachegroups/' . $cg_id)->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/alerts/0/level", "error" )
    ->json_is( "/alerts/0/text", "This cachegroup is currently used by servers.")
            , "Is the Cachegroup id valid?";

my $midcg_id = &get_cg_id('cache_group_mid');
ok $t->delete_ok('/api/1.2/cachegroups/' . $midcg_id)->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/alerts/0/level", "error" )
    ->json_is( "/alerts/0/text", "This cachegroup is currently used as a parent cachegroup." )
            , "Is the Cachegroup id valid?";

my $svr_id =&get_svr_id('edge_streamer_1');
ok $t->delete_ok('/api/1.2/servers/' . $svr_id)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->delete_ok('/api/1.2/cachegroups/' . $cg_id)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/alerts/0/level", "success" )
    ->json_is( "/alerts/0/text", "Cachegroup deleted." )
            , "Is the Cachegroup id valid?";
ok $t->delete_ok('/api/1.2/cachegroups/' . $cg_id)->status_is(404)->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->put_ok('/api/1.2/cachegroups/' . $cg_id => {Accept => 'application/json'} => json => {
        "name" => "cache_group_edge_1",
        "shortName" => "cg_edge_1",
        "typeId" => 5})->status_is(404)->or( sub { diag $t->tx->res->content->asset->{content}; } );

Test::TestHelper->unload_core_data($schema);
Test::TestHelper->load_core_data($schema);

# Count the 'response number'
my $count_response = sub {
    my ( $t, $count ) = @_;
    my $json = decode_json( $t->tx->res->content->asset->slurp );
    my $r    = $json->{response};
    return $t->success( is( scalar(@$r), $count ) );
};

# there are currently 61 parameters not assigned to cachegroup 100
$t->get_ok('/api/1.2/cachegroups/100/unassigned_parameters')->status_is(200)->$count_response(63)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );

# there are currently 61 parameters not assigned to cachegroup 200
$t->get_ok('/api/1.2/cachegroups/200/unassigned_parameters')->status_is(200)->$count_response(63)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );


ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();

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
