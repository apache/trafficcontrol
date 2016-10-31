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
use strict;
use warnings;
use Test::TestHelper;
use Fixtures::TmUser;
use Fixtures::SteeringUsers;
use Fixtures::SteeringType;
use Fixtures::SteeringDeliveryservice;
use Fixtures::SteeringDeliveryserviceRegex;
use Fixtures::SteeringTarget;
use Fixtures::SteeringDeliveryServiceUsers;

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

my $steering_users = Fixtures::SteeringUsers->new($schema_values);
Test::TestHelper->load_all_fixtures($steering_users);

my $steering_type = Fixtures::SteeringType->new($schema_values);
Test::TestHelper->load_all_fixtures($steering_type);

my $steering_deliveryservice = Fixtures::SteeringDeliveryservice->new($schema_values);
Test::TestHelper->load_all_fixtures($steering_deliveryservice);

my $steering_deliveryservice_regex = Fixtures::SteeringDeliveryserviceRegex->new($schema_values);
Test::TestHelper->load_all_fixtures($steering_deliveryservice_regex);

my $steering_target = Fixtures::SteeringTarget->new($schema_values);
Test::TestHelper->load_all_fixtures($steering_target);

my $steering_deliveryservice_users = Fixtures::SteeringDeliveryServiceUsers->new($schema_values);
Test::TestHelper->load_all_fixtures($steering_deliveryservice_users);

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
    ->json_is("/response/0/filters/0/deliveryService", "target-ds1")
    ->json_is("/response/0/filters/0/pattern", ".*/force-to-one-also/.*")
    ->json_is("/response/0/filters/1/deliveryService", "target-ds1")
    ->json_is("/response/0/filters/1/pattern", ".*/force-to-one/.*")
    ->json_is("/response/0/targets/1/deliveryService", "target-ds2")
    ->json_is("/response/0/targets/1/weight", 7654)
    ->json_is("/response/1/deliveryService", "steering-ds2")
    ->json_is("/response/1/targets/0/deliveryService", "target-ds3")
    ->json_is("/response/1/targets/0/weight", 123)
    ->json_is("/response/1/filters/0/pattern", ".*/use-three/.*")
    ->json_is("/response/1/filters/0/deliveryService", "target-ds3")
    ->json_is("/response/1/targets/1/deliveryService", "target-ds4")
    ->json_is("/response/1/targets/1/weight", 999)
    ->json_is("/response/1/filters/1/pattern", ".*/go-to-four/.*")
    ->json_is("/response/1/filters/1/deliveryService", "target-ds4");

ok $t->get_ok("/internal/api/1.2/steering/steering-ds1.json")->status_is(200)
    ->or(sub {diag $t->tx->res->headers->to_string();})
        ->json_is("/response/deliveryService", "steering-ds1")
        ->json_is("/response/targets/0/deliveryService", "target-ds1")
        ->json_is("/response/targets/0/weight", 1000)
        ->json_is("/response/filters/0/pattern", ".*/force-to-one-also/.*")
        ->json_is("/response/filters/1/pattern", ".*/force-to-one/.*");

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
    ->or( sub { diag $t->tx->res->headers->to_string(); } )
    ->json_is("/response/0/deliveryService", "steering-ds1")
    ->json_is("/response/0/targets/0/deliveryService", "target-ds1")
    ->json_is("/response/0/targets/0/weight", 1000)
    ->json_is("/response/0/filters/0/pattern", ".*/force-to-one-also/.*")
    ->json_is("/response/0/filters/1/pattern", ".*/force-to-one/.*")
    ->json_hasnt("/response/0/filters/2/pattern")
    ->json_is("/response/0/targets/1/deliveryService", "target-ds2")
    ->json_is("/response/0/targets/1/weight", 7654)
    ->json_hasnt("/response/0/filters/1/filter/0")
    ->json_hasnt("/response/1");

ok $t->post_ok("/internal/api/1.2/steering", json => { "something" => "value" } )->status_is(401)
    ->or(sub {diag $t->tx->res->headers->to_string();});

ok $t->get_ok("/internal/api/1.2/steering/steering-ds2.json")->status_is(404)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->put_ok("/internal/api/1.2/steering/steering-ds2", json => {"any" => "thing"})->status_is(401)
    ->or( sub { diag $t->tx->res->headers->to_string(); } );

ok $t->put_ok("/internal/api/1.2/steering/steering-ds1",
    json =>  {
            "targets" => [
                { "deliveryService" => "target-ds1", "weight" => 5555 },
                { "deliveryService" => "target-ds2", "weight" => 4444 }
            ],
            "filters" => [
                {
                    "deliveryService" => "target-ds3",
                    "pattern" => ".*/force-to-one/.*"
                },
            ]
        })
    ->status_is(409);

ok $t->put_ok("/internal/api/1.2/steering/steering-ds1",
        json => {
            "targets" => [
                {
                    "deliveryService" => "target-ds1",
                    "weight" => 5555
                },
                {
                    "deliveryService" => "target-ds2",
                    "weight" => 4444
                }
            ],
            "filters" => [
                {
                    "deliveryService" => "target-ds1",
                    "pattern" => ".*/force-to-one/.*"
                },
                {
                    "deliveryService" => "target-ds1",
                    "pattern" => ".*/andnowforsomethingcompletelydifferent/.*"
                },
                {
                    "deliveryService" => "target-ds2",
                    "pattern" => ".*/always-two/.*"
                },
            ]
        })
    ->status_is(200)->or(sub { diag $t->tx->res->headers->to_string(); })
    ->json_is("/response/deliveryService", "steering-ds1")
    ->json_is("/response/targets/0/deliveryService", "target-ds1")
    ->json_is("/response/targets/0/weight", 5555)
    ->json_is("/response/targets/1/deliveryService", "target-ds2")
    ->json_is("/response/targets/1/weight", 4444)
    ->json_is("/response/filters/0/pattern", ".*/andnowforsomethingcompletelydifferent/.*")
    ->json_is("/response/filters/1/pattern", ".*/force-to-one/.*")
    ->json_is("/response/filters/2/pattern", ".*/always-two/.*" );

ok $t->get_ok("/internal/api/1.2/steering/steering-ds1.json")
    ->status_is(200)->or(sub { diag $t->tx->res->headers->to_string(); })
    ->json_is("/response/deliveryService", "steering-ds1")
    ->json_is("/response/targets/0/deliveryService", "target-ds1")
    ->json_is("/response/targets/0/weight", 5555)
    ->json_is("/response/targets/1/deliveryService", "target-ds2")
    ->json_is("/response/targets/1/weight", 4444)
    ->json_is("/response/filters/0/pattern", ".*/andnowforsomethingcompletelydifferent/.*")
    ->json_is("/response/filters/1/pattern", ".*/force-to-one/.*")
    ->json_is("/response/filters/2/pattern", ".*/always-two/.*" );

ok $t->put_ok("/internal/api/1.2/steering/steering-ds1",
        json => {
            "targets" => [
                {
                    "deliveryService" => "target-ds1",
                    "weight" => 1111
                },
                {
                    "deliveryService" => "target-ds2",
                    "weight" => 8888
                }
            ]
        })
        ->status_is(200)->or(sub { diag $t->tx->res->headers->to_string(); })
        ->json_is("/response/deliveryService", "steering-ds1")
        ->json_is("/response/targets/0/deliveryService", "target-ds1")
        ->json_is("/response/targets/0/weight", 1111)
        ->json_hasnt("/response/filter/0/pattern")
        ->json_is("/response/targets/1/deliveryService", "target-ds2")
        ->json_is("/response/targets/1/weight", 8888)
        ->json_is("/response/filters/2/pattern", ".*/always-two/.*" );

ok $t->get_ok("/internal/api/1.2/steering/steering-ds1.json")
        ->status_is(200)->or(sub { diag $t->tx->res->headers->to_string(); })
        ->json_is("/response/deliveryService", "steering-ds1")
        ->json_is("/response/targets/0/deliveryService", "target-ds1")
        ->json_is("/response/targets/0/weight", 1111)
        ->json_hasnt("/response/filter/0/pattern")
        ->json_is("/response/targets/1/deliveryService", "target-ds2")
        ->json_is("/response/targets/1/weight", 8888)
        ->json_is("/response/filters/2/pattern", ".*/always-two/.*" );

#bad json
ok $t->put_ok("/internal/api/1.2/steering/steering-ds1",
    json => {"foo" => "bar"})
    ->status_is(400)
    ->json_is("/message", "please provide a valid json including targets");

#remove filters for single DS
ok $t->put_ok("/internal/api/1.2/steering/steering-ds1",
        json => {
            "targets" => [
                {
                    "deliveryService" => "target-ds1",
                    "weight" => 5555
                },
                {
                    "deliveryService" => "target-ds2",
                    "weight" => 4444
                }
            ],
            "filters" => [
                {
                    "deliveryService" => "target-ds1",
                    "pattern" => ".*/force-to-one/.*"
                }
            ]
        })
    ->status_is(200)->or(sub { diag $t->tx->res->headers->to_string(); })
    ->json_is("/response/deliveryService", "steering-ds1")
    ->json_is("/response/targets/0/deliveryService", "target-ds1")
    ->json_is("/response/targets/0/weight", 5555)
    ->json_is("/response/targets/1/deliveryService", "target-ds2")
    ->json_is("/response/targets/1/weight", 4444)
    ->json_hasnt("/response/filters/1/pattern");

    #remove all filters
ok $t->put_ok("/internal/api/1.2/steering/steering-ds1",
        json => {
            "targets" => [
                {
                    "deliveryService" => "target-ds1",
                    "weight" => 5555
                },
                {
                    "deliveryService" => "target-ds2",
                    "weight" => 4444
                }
            ],
            "filters" => []
        })
    ->status_is(200)->or(sub { diag $t->tx->res->headers->to_string(); })
    ->json_is("/response/deliveryService", "steering-ds1")
    ->json_is("/response/targets/0/deliveryService", "target-ds1")
    ->json_is("/response/targets/0/weight", 5555)
    ->json_is("/response/targets/1/deliveryService", "target-ds2")
    ->json_is("/response/targets/1/weight", 4444)
    ->json_hasnt("/response/filters/0/pattern");

#invalid json
ok $t->put_ok("/internal/api/1.2/steering/steering-ds1",
        json => {
            "targets" => [
                {
                    "deliveryService" => "target-ds1",
                    "weight" => 5555
                },
                {
                    "deliveryService" => "target-ds2",
                    "weight" => 4444
                }
            ],
            "filters" => [
            {
                    "pattern" => ".*/force-to-one/.*"
                }
            ]
        })
    ->status_is(400)
    ->json_is("/message", "please provide a valid json for filters");

ok $t->put_ok("/internal/api/1.2/steering/steering-ds1",
        json => {
            "targets" => [
                {
                    "deliveryService" => "target-ds1",
                    "weight" => 5555
                },
                {
                    "deliveryService" => "target-ds2",
                    "weight" => 4444
                }
            ],
            "filters" => [
            {
                    "deliveryService" => "target-ds1"
                }
            ]
        })
    ->status_is(400)
    ->json_is("/message", "please provide a valid json for filters");

ok $t->put_ok("/internal/api/1.2/steering/steering-ds1",
        json => {
            "filters" => [
            {
                    "deliveryService" => "target-ds1",
                    "pattern" => ".*/force-to-one/.*"
                }
            ]
        })
    ->status_is(400)->json_is("/message", "please provide a valid json including targets");

ok $t->put_ok("/internal/api/1.2/steering/steering-ds1",
        json => {
            "targets" => [
                {
                    "deliveryService" => "target-ds1",
                },
                {
                    "deliveryService" => "target-ds2",
                    "weight" => 4444
                }
            ]
        })
    ->status_is(400)->json_is("/message", "please provide a valid json for targets");


$t->post_ok("/api/1.2/user/logout")->status_is(200);

$dbh->disconnect();
done_testing();

