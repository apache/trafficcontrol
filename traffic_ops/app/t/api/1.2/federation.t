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
use JSON;
use DBI;
use strict;
use warnings;
no warnings 'once';
use warnings 'all';
use Test::TestHelper;
use Fixtures::Federation;
use Fixtures::FederationDeliveryservice;
use Fixtures::FederationResolver;
use Fixtures::FederationFederationResolver;
use Fixtures::FederationTmuser;

#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

BEGIN { $ENV{MOJO_MODE} = "test" }

my $schema = Schema->connect_to_database;
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);
Test::TestHelper->teardown( $schema, "Federation" );
Test::TestHelper->teardown( $schema, "FederationTmuser" );
Test::TestHelper->teardown( $schema, "FederationDeliveryservice" );
Test::TestHelper->teardown( $schema, "FederationFederationResolver" );
Test::TestHelper->teardown( $schema, "FederationResolver" );

#load core test data
Test::TestHelper->load_core_data($schema);

# Count the 'response number'
my $count_response = sub {
    my ( $t, $count ) = @_;
    my $json = decode_json( $t->tx->res->content->asset->slurp );
    my $r    = $json->{response};
    return $t->success( is( scalar(@$r), $count ) );
};

ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';

# create federation: fail (no cname)
ok $t->post_ok('/api/1.2/cdns/cdn1/federations' => {Accept => 'application/json'} => json => {
            "ttl" => 256 })->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/level", "error" )
        ->json_is( "/alerts/0/text", "cname is required" )
    , 'Did federation create fail due to no cname?';

# create federation: fail (no ttl)
ok $t->post_ok('/api/1.2/cdns/cdn1/federations' => {Accept => 'application/json'} => json => {
            "cname" => "my.cname.com." })->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/level", "error" )
        ->json_is( "/alerts/0/text", "ttl is required" )
    , 'Did federation create fail due to no ttl?';

# create federation: fail (bad cname)
ok $t->post_ok('/api/1.2/cdns/cdn1/federations' => {Accept => 'application/json'} => json => {
            "cname" => "my.cname",
            "ttl" => 256 })->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/level", "error" )
        ->json_is( "/alerts/0/text", "cname must contain no spaces and end with a dot" )
    , 'Did federation create fail due to bad cname?';

# create federation: fail (bad ttl)
my $cname = 'my.cname.';
ok $t->post_ok('/api/1.2/cdns/cdn1/federations' => {Accept => 'application/json'} => json => {
            "cname" => $cname,
            "ttl" => "hello" })->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/level", "error" )
        ->json_is( "/alerts/0/text", "ttl must be a number" )
    , 'Did federation create fail due to bad ttl?';

# create federation: success
ok $t->post_ok('/api/1.2/cdns/cdn1/federations' => {Accept => 'application/json'} => json => {
            "cname" => $cname,
            "ttl" => 256 })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/response/cname" => $cname)
        ->json_is( "/response/ttl" => 256)
    , 'Was a federation created?';

# fetch created federation
my $fed_id = &get_fed_id($cname);
ok $t->get_ok('/api/1.2/cdns/cdn1/federations/' . $fed_id)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/response/0/id" => $fed_id)
        ->json_is( "/response/0/cname" => $cname)
        ->json_is( "/response/0/ttl" => 256);

# update federation
my $updated_cname = 'updated.cname.';
ok $t->put_ok('/api/1.2/cdns/cdn1/federations/' . $fed_id  => {Accept => 'application/json'} => json => {
            "cname" => $updated_cname,
            "ttl" => 100,
        })
        ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/response/cname" => $updated_cname)
        ->json_is( "/response/ttl" => 100)
    , 'Was the federation updated?';

# add a user to the federation
ok $t->get_ok("/api/1.2/federations/$fed_id/users")->status_is(200)->$count_response(0)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok("/api/1.2/federations/$fed_id/users" => {Accept => 'application/json'} => json => {
            "userIds" => [ 500 ]})->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/level", "success" )
        ->json_is( "/alerts/0/text", "1 user(s) were assigned to the " . $updated_cname . " federation" )
    , "Was a user added to the federation?";

ok $t->get_ok("/api/1.2/federations/$fed_id/users")->status_is(200)->$count_response(1)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } );

# add a ds to the federation
ok $t->get_ok("/api/1.2/federations/$fed_id/deliveryservices")->status_is(200)->$count_response(0)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok("/api/1.2/federations/$fed_id/deliveryservices" => {Accept => 'application/json'} => json => {
            "dsIds" => [ 100 ]})->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/level", "success" )
        ->json_is( "/alerts/0/text", "1 delivery service(s) were assigned to the " . $updated_cname . " federation" )
    , "Was a delivery service added to the federation?";

ok $t->get_ok("/api/1.2/federations/$fed_id/deliveryservices")->status_is(200)->$count_response(1)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } );

# create a federation resolver
my $ip = '2.2.2.2';
ok $t->post_ok('/api/1.2/federation_resolvers' => {Accept => 'application/json'} => json => {
            "ipAddress" => $ip,
            "typeId" => 33 })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/response/ipAddress" => $ip)
        ->json_is( "/response/typeId" => 33)
    , 'Was a federation resolver created?';

# add a federation resolver to the federation
ok $t->get_ok("/api/1.2/federations/$fed_id/federation_resolvers")->status_is(200)->$count_response(0)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } );

my $fed_res_id = &get_fed_res_id($ip);
ok $t->post_ok("/api/1.2/federations/$fed_id/federation_resolvers" => {Accept => 'application/json'} => json => {
            "fedResolverIds" => [ $fed_res_id ]})->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/level", "success" )
        ->json_is( "/alerts/0/text", "1 resolver(s) were assigned to the " . $updated_cname . " federation" )
    , "Was a federation resolver added to the federation?";

ok $t->get_ok("/api/1.2/federations/$fed_id/federation_resolvers")->status_is(200)->$count_response(1)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } );

# with the ds assigned to the federation, does the cdn now have 1 federation?
ok $t->get_ok('/api/1.2/cdns/cdn1/federations')->status_is(200)->$count_response(1)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } );

# now let's delete the federation resolver
ok $t->delete_ok('/api/1.2/federation_resolvers/' . $fed_res_id)
        ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/level", "success" )
        ->json_is( "/alerts/0/text", "Federation resolver deleted [ IP = " . $ip . " ] with id: " . $fed_res_id )
    , "Was the federation resolver deleted?";

ok $t->get_ok("/api/1.2/federations/$fed_id/federation_resolvers")->status_is(200)->$count_response(0)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } );

# and finally, let's delete the federation
ok $t->delete_ok('/api/1.2/cdns/cdn1/federations/' . $fed_id)
        ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/level", "success" )
        ->json_is( "/alerts/0/text", "Federation deleted [ cname = " . $updated_cname . " ] with id: " . $fed_id )
    , "Was the federation deleted?";

ok $t->get_ok('/api/1.2/cdns/cdn1/federations/' . $fed_id)->status_is(404)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.2/cdns/cdn1/federations')->status_is(200)->$count_response(0)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok("/api/1.2/federations/$fed_id/users")->status_is(200)->$count_response(0)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok("/api/1.2/federations/$fed_id/deliveryservices")->status_is(200)->$count_response(0)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();

sub get_fed_id {
    my $fed_cname = shift;
    my $q      = "select id from federation where cname = \'$fed_cname\'";
    my $get_fed = $dbh->prepare($q);
    $get_fed->execute();
    my $p = $get_fed->fetchall_arrayref( {} );
    $get_fed->finish();
    my $id = $p->[0]->{id};
    return $id;
}

sub get_fed_res_id {
    my $ip = shift;
    my $q      = "select id from federation_resolver where ip_address = \'$ip\'";
    my $get_fed_res = $dbh->prepare($q);
    $get_fed_res->execute();
    my $p = $get_fed_res->fetchall_arrayref( {} );
    $get_fed_res->finish();
    my $id = $p->[0]->{id};
    return $id;
}


