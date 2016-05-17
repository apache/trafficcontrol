package main;

#
# Copyright 2016 Comcast Cable Communications Management, LLC
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
use strict;
use warnings;
use Test::TestHelper;
use Fixtures::TmUser;
use Fixtures::SteeringType;
use Fixtures::SteeringDeliveryservice;
use Fixtures::SteeringDeliveryserviceRegex;
use Fixtures::SteeringTarget;

BEGIN { $ENV{MOJO_MODE} = "test" }

my $t      = Test::Mojo->new("TrafficOps");
my $dbh    = Schema->database_handle;
my $schema = Schema->connect_to_database;

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);
Test::TestHelper->teardown( $schema, "SteeringTarget" );

#load core test data
Test::TestHelper->load_core_data($schema);

my $schema_values = { schema => $schema, no_transactions => 1 };

my $steering_type = Fixtures::SteeringType->new($schema_values);
Test::TestHelper->load_all_fixtures($steering_type);

my $steering_deliveryservice = Fixtures::SteeringDeliveryservice->new($schema_values);
Test::TestHelper->load_all_fixtures($steering_deliveryservice);

my $steering_deliveryservice_regex = Fixtures::SteeringDeliveryserviceRegex->new($schema_values);
Test::TestHelper->load_all_fixtures($steering_deliveryservice_regex);

my $steering_target = Fixtures::SteeringTarget->new($schema_values);
Test::TestHelper->load_all_fixtures($steering_target);

####### Unauthorized User ################################################################################
ok $t->post_ok( "/api/1.2/user/login", => json => { u => Test::TestHelper::CODEBIG_USER, p => Test::TestHelper::CODEBIG_PASSWORD } )
    ->status_is(200)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok("/internal/api/1.2/steering.json")
    ->status_is(401)
    ->or( sub { diag $t->tx->res->headers->to_string(); } );

$t->post_ok("/api/1.2/user/logout")->status_is(200);

####### Administrator ##################################################################################
ok $t->post_ok( "/api/1.2/user/login", => json => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )
        ->status_is(200)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok("/internal/api/1.2/steering.json")->status_is(200)
        ->or( sub { diag $t->tx->res->headers->to_string(); } )
    ->json_is("/response/0/deliveryService", "steering-ds1")
    ->json_is("/response/0/targets/0/deliveryService", "target-ds1")
    ->json_is("/response/0/targets/0/weight", 1000)
    ->json_is("/response/0/targets/0/filters/0", ".*/force-to-one/.*")
    ->json_is("/response/0/targets/0/filters/1", ".*/force-to-one-also/.*")
    ->json_is("/response/0/targets/1/deliveryService", "target-ds2")
    ->json_is("/response/0/targets/1/weight", 7654)
    ->json_is("/response/1/deliveryService", "steering-ds2")
    ->json_is("/response/1/targets/0/deliveryService", "target-ds3")
    ->json_is("/response/1/targets/0/weight", 123)
    ->json_is("/response/1/targets/0/filters/0", ".*/use-three/.*")
    ->json_is("/response/1/targets/1/deliveryService", "target-ds4")
    ->json_is("/response/1/targets/1/weight", 999)
    ->json_is("/response/1/targets/1/filters/0", ".*/go-to-four/.*");

ok $t->post_ok("/internal/api/1.2/steering",
        json => {
            "stuff" => "junk",
        }
    )->status_is(400)
        ->or(sub {diag $t->tx->res->headers->to_string();});

ok $t->post_ok("/internal/api/1.2/steering",
        json => {
            "deliveryService" => "steering-ds1"
        }
    )->status_is(400)
        ->or(sub {diag $t->tx->res->headers->to_string();});

ok $t->post_ok("/internal/api/1.2/steering",
        json => {
            "deliveryService" => "steering-ds1",
            "targets" => "stuff"
        }
    )->status_is(400)
        ->or(sub {diag $t->tx->res->headers->to_string();});

ok $t->post_ok("/internal/api/1.2/steering",
        json => {
            "deliveryService" => "steering-ds1",
            "targets" => [
                {"deliveryService" => "example"},
                {"woops" => "example"},
            ]
        }
    )->status_is(400)
        ->or(sub {diag $t->tx->res->headers->to_string();});

ok $t->post_ok("/internal/api/1.2/steering",
        json => {
            "deliveryService" => "nonexistent-ds",
            "targets" => [
                {"deliveryService" => "target-ds1"},
                {"deliveryService" => "target-ds3"}
            ]
        }
    )->status_is(409)
        ->or(sub {diag $t->tx->res->headers->to_string();});

ok $t->post_ok("/internal/api/1.2/steering",
        json => {
            "deliveryService" => "steering-ds1",
            "targets" => [
                {"deliveryService" => "nonexistent-ds1"},
                {"deliveryService" => "target-ds3"}
            ]
        }
    )->status_is(409)
        ->or(sub {diag $t->tx->res->headers->to_string();});

ok $t->post_ok("/internal/api/1.2/steering",
        json => {
            "deliveryService" => "new-steering-ds",
            "targets" => [
                {"deliveryService" => "target-ds1"},
                {"deliveryService" => "target-ds3"}
            ]
        }
    )->status_is(201)
        ->header_is('Location', "/internal/api/1.2/steering/new-steering-ds.json")
        ->or(sub {diag $t->tx->res->headers->to_string();});

$t->post_ok("/api/1.2/user/logout")->status_is(200);


####### Steering User 1 ################################################################################
ok $t->post_ok( "/api/1.2/user/login", => json => { u => Test::TestHelper::STEERING_USER_1, p => Test::TestHelper::STEERING_PASSWORD_1 } )
        ->status_is(200)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok("/internal/api/1.2/steering.json")->status_is(200)
    ->or( sub { diag $t->tx->res->headers->to_string(); } );

ok $t->post_ok("/internal/api/1.2/steering", json => { "something" => "value" } )->status_is(401)
        ->or(sub {diag $t->tx->res->headers->to_string();});

$t->post_ok("/api/1.2/user/logout")->status_is(200);


######## Steering User 2 ################################################################################
#ok $t->post_ok( "/login", => form => { u => Test::TestHelper::STEERING_USER_2, p => Test::TestHelper::STEERING_PASSWORD_2 } )->status_is(302)
#        ->or( sub { diag $t->tx->res->content->asset->{content}; } );
#
#$t->get_ok("/internal/api/1.2/steering.json")->status_is(200)
#    ->or( sub { diag $t->tx->res->content->asset->{content}; } )
#    ->json_is( "/response/0/deliveryService", "test-steering-ds-2" )
#    ->json_is( "/response/0/bypasses/0/filter", ".*/force-to-three/.*" )
#    ->json_is( "/response/0/bypasses/0/destination", "test-ds3" )
#    ->json_is( "/response/0/bypasses/1/filter", ".*/force-to-four/.*" )
#    ->json_is( "/response/0/bypasses/1/destination", "test-ds4" )
#    ->json_is( "/response/0/steeredDeliveryServices/0/id", "test-ds3" )
#    ->json_is( "/response/0/steeredDeliveryServices/0/weight", "555" )
#    ->json_is( "/response/0/steeredDeliveryServices/1/id", "test-ds4" )
#    ->json_is( "/response/0/steeredDeliveryServices/1/weight", "1234" );

####### Admin User ################################################################################
#ok $t->post_ok( "/login", => form => { u => "admin", p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
#        ->or( sub { diag $t->tx->res->content->asset->{content}; } );
#
#$t->get_ok("/internal/api/1.2/steering.json")->status_is(200)
#    ->or( sub { diag $t->tx->res->content->asset->{content}; } )
#    ->json_is( "/response/0/deliveryService", "test-steering-ds-1" )
#    ->json_is( "/response/0/bypasses/0/filter", ".*/force-to-one/.*" )
#    ->json_is( "/response/0/bypasses/0/destination", "test-ds1" )
#    ->json_is( "/response/0/bypasses/1/filter", ".*/force-to-two/.*" )
#    ->json_is( "/response/0/bypasses/1/destination", "test-ds2" )
#    ->json_is( "/response/0/steeredDeliveryServices/0/id", "test-ds1" )
#    ->json_is( "/response/0/steeredDeliveryServices/0/weight", "9000" )
#    ->json_is( "/response/0/steeredDeliveryServices/1/id", "test-ds2" )
#    ->json_is( "/response/0/steeredDeliveryServices/1/weight", "1000" )
#
#    ->json_is( "/response/1/deliveryService", "test-steering-ds-2" )
#    ->json_is( "/response/1/bypasses/0/filter", ".*/force-to-three/.*" )
#    ->json_is( "/response/1/bypasses/0/destination", "test-ds3" )
#    ->json_is( "/response/1/bypasses/1/filter", ".*/force-to-four/.*" )
#    ->json_is( "/response/1/bypasses/1/destination", "test-ds4" )
#    ->json_is( "/response/1/steeredDeliveryServices/0/id", "test-ds3" )
#    ->json_is( "/response/1/steeredDeliveryServices/0/weight", "555" )
#    ->json_is( "/response/1/steeredDeliveryServices/1/id", "test-ds4" )
#    ->json_is( "/response/1/steeredDeliveryServices/1/weight", "1234" );
#
#ok $t->post_ok("/internal/api/1.2/steering",
#        json => {
#            "id" => "steering-ds-1",
#            "steeredDeliveryServices" => [
#                {"id" => "steering-target-1"},
#                {"id" => "steering-target-2"}
#            ]
#        }
#    )->status_is(201)
#        ->or(sub {diag $t->tx->res->headers->to_string();})
#        ->header_is('Location', '/internal/api/1.2/steering/steering-ds-1');
#
#$t->get_ok("/internal/api/1.2/steering/steering-ds-1")->status_is(200)
#    ->or(sub {diag $t->tx->res->headers->to_string();})
#    ->json_is("deliveryService", "steering-ds-1")
#    ->json_is("steeredDeliveryServices/0/id", "steering-target-1")
#    ->json_is("steeredDeliveryServices/0/weight", "0")
#    ->json_is("steeredDeliveryServices/1/id", "steering-target-2")
#    ->json_is("steeredDeliveryServices/1/weight", "0");

#Test::TestHelper->unload_core_data($schema);
#Test::TestHelper->teardown( $schema, "Steering" );

$dbh->disconnect();
done_testing();

