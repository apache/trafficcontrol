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
use Test::TestHelper;

#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

BEGIN { $ENV{MOJO_MODE} = "test" }

my $t      = Test::Mojo->new('TrafficOps');
my $dbh    = Schema->database_handle;
my $schema = Schema->connect_to_database;

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);
Test::TestHelper->teardown( $schema, 'Asn' );

#load core test data
Test::TestHelper->load_core_data($schema);

ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
  ->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok("/api/1.2/asns")->status_is(200)->json_is( "/response/0/id", 100 )->json_is( "/response/0/cachegroup", "mid-northeast-group" )
  ->json_is( "/response/0/asn", 9939 )->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok("/api/1.2/asns?orderby=id")->status_is(200)->json_is( "/response/0/id", 100 )
  ->json_is( "/response/0/cachegroup", "mid-northeast-group" )->json_is( "/response/0/asn", 9939 )
  ->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok("/api/1.2/asns?orderby=cachegroup")->status_is(200)->json_is( "/response/0/id", 100 )
  ->json_is( "/response/0/cachegroup", "mid-northeast-group" )->json_is( "/response/0/asn", 9939 )
  ->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok("/api/1.2/asns/200")->status_is(200)->json_is( "/response/0/id", 200 )
  ->json_is( "/response/0/cachegroup", "mid-northwest-group" )->json_is( "/response/0/asn", 9940 )
  ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok('/api/1.2/asns' => {Accept => 'application/json'} => json => {
            "cachegroupId" => 100
        })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/level" => "error" )
        ->json_is( "/alerts/0/text" => "asn is required" )
    , 'Does ASN create fail because asn is required?';

ok $t->post_ok('/api/1.2/asns' => {Accept => 'application/json'} => json => {
            "asn" => "eightfivetwo",
            "cachegroupId" => 100
        })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/level" => "error" )
        ->json_is( "/alerts/0/text" => "asn must be a positive integer" )
    , 'Does ASN create fail because asn is a string instead of an integer?';

ok $t->post_ok('/api/1.2/asns' => {Accept => 'application/json'} => json => {
            "asn" => 852
        })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/level" => "error" )
        ->json_is( "/alerts/0/text" => "cachegroupId is required" )
    , 'Does ASN create fail because cache group ID is required?';

ok $t->post_ok('/api/1.2/asns' => {Accept => 'application/json'} => json => {
            "asn" => 852,
            "cachegroupId" => "hundred"
        })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/level" => "error" )
        ->json_is( "/alerts/0/text" => "cachegroupId must be a positive integer" )
    , 'Does ASN create fail because cachegroupId is a string instead of an integer?';

ok $t->post_ok('/api/1.2/asns' => {Accept => 'application/json'} => json => {
            "asn" => 852,
            "cachegroupId" => 100
        })
        ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/response/asn" => 852 )
        ->json_is( "/response/cachegroupId" => 100 )
        ->json_is( "/response/cachegroup" => "mid-northeast-group" )
        ->json_is( "/alerts/0/level" => "success" )
        ->json_is( "/alerts/0/text" => "ASN create was successful." )
    , 'Is the ASN successfully created?';

my $asn_id = &get_asn_id(852);

ok $t->put_ok('/api/1.2/asns/' . $asn_id  => {Accept => 'application/json'} => json => {
            "asn" => "eightfivethree",
            "cachegroupId" => 100
        })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/level" => "error" )
        ->json_is( "/alerts/0/text" => "asn must be a positive integer" )
    , 'Does the asn update fail due to bad asn?';

ok $t->put_ok('/api/1.2/asns/' . $asn_id  => {Accept => 'application/json'} => json => {
            "asn" => 853,
            "cachegroupId" => 100
        })
        ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/response/asn" => 853 )
        ->json_is( "/alerts/0/level" => "success" )
    , 'Does the asn details return?';

ok $t->delete_ok('/api/1.2/asns/' . $asn_id)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

$dbh->disconnect();
done_testing();

sub get_asn_id {
    my $asn = shift;
    my $q    = "select id from asn where asn = \'$asn\'";
    my $get_svr = $dbh->prepare($q);
    $get_svr->execute();
    my $p = $get_svr->fetchall_arrayref( {} );
    $get_svr->finish();
    my $id = $p->[0]->{id};
    return $id;
}

