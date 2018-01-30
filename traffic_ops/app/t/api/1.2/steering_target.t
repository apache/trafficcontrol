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

ok $t->post_ok( '/api/1.1/user/login', json => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(200),
    'Log into the admin user?';

ok $t->post_ok('/api/1.2/steering/900/targets' => {Accept => 'application/json'} => json => {
            "targetId" => 1000,
            "value" => 852,
            "typeId" => 41
        })->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    , 'Failed to add on steering target value?';

ok $t->post_ok('/api/1.2/steering/900/targets' => {Accept => 'application/json'} => json => {
            "targetId" => 1000,
            "value" => 852,
            "typeId" => 40
        })->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/response/0/deliveryServiceId" => 900 )
        ->json_is( "/response/0/deliveryService" => "steering-ds3" )
        ->json_is( "/response/0/targetId" => 1000 )
        ->json_is( "/response/0/target" => "steering-target-ds1" )
        ->json_is( "/response/0/value" => 852 )
        ->json_is( "/response/0/typeId" => 40 )
        ->json_is( "/response/0/type" => "STEERING_ORDER" )
    , 'Is the steering target created?';

ok $t->post_ok('/api/1.2/steering/900/targets' => {Accept => 'application/json'} => json => {
            "targetId" => 1000,
            "value" => 6,
            "typeId" => 40
        })->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    , 'Failed to readd steering target?';

ok $t->get_ok("/api/1.2/steering/900/targets")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/response/0/deliveryServiceId" => 900 )
        ->json_is( "/response/0/targetId" => 1000 )
        ->json_is( "/response/0/value" => 852 )
        ->json_is( "/response/0/typeId" => 40 )
    , 'Are steering targets returned?';

$t->get_ok("/api/1.2/steering/900/targets/1000")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/value" => 852 )
    ->json_is( "/response/0/typeId" => 40 )
    , 'Is the steering target returned?';

ok $t->put_ok('/api/1.2/steering/900/targets/1000' => {Accept => 'application/json'} => json => {
            "value" => 999,
            "typeId" => 41
        })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    , 'Failed to change a steering target type to invalid?';

ok $t->put_ok('/api/1.2/steering/900/targets/1000' => {Accept => 'application/json'} => json => {
            "value" => 999,
            "typeId" => 40
        })
        ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/response/deliveryServiceId" => 900 )
        ->json_is( "/response/deliveryService" => "steering-ds3" )
        ->json_is( "/response/targetId" => 1000 )
        ->json_is( "/response/target" => "steering-target-ds1" )
        ->json_is( "/response/value" => 999 )
        ->json_is( "/response/typeId" => 40 )
        ->json_is( "/response/type" => "STEERING_ORDER" )
        ->json_is( "/alerts/0/level" => "success" )
    , 'Did the steering target update?';

ok $t->delete_ok('/api/1.2/steering/900/targets/1000')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    , 'Did the steering target get deleted?';

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

done_testing();
