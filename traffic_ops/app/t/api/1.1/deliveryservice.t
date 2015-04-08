package main;
#
# Copyright 2015 Comcast Cable Communications Management, LLC
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
use strict;
use warnings;

BEGIN { $ENV{MOJO_MODE} = "test" }

my $dbh    = Schema->database_handle;
my $schema = Schema->connect_to_database;
my $t      = Test::Mojo->new('TrafficOps');

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);

#load core test data
Test::TestHelper->load_core_data($schema);

#login
ok $t->post_ok( '/login', => form => { u => Test::TestHelper::PORTAL_USER, p => Test::TestHelper::PORTAL_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

#get all ds's
ok $t->get_ok('/api/1.1/deliveryservices.json')->status_is(200)->json_has( '/response', 'has a response' ), "Get delivery services";

#validate a ds
ok $t->get_ok('/api/1.1/deliveryservices.json')->status_is(200)->json_is( '/response/0/id' => 1 ), "Has a delivery service";

#get a specific ds
ok $t->get_ok('/api/1.1/deliveryservices/1.json')->status_is(200)->json_has( '/response', 'has a response' ), "Get delivery services";
ok $t->get_ok('/api/1.1/deliveryservices.json')->status_is(200)->json_is( '/response/0/id' => 1 )
	->json_is( '/response/0/longDesc2' => "test-ds1 long_desc_2" )->json_is( '/response/0/globalMaxTps' => "0" )
	->json_is( '/response/0/maxDnsAnswers' => "0" )->json_is( '/response/0/missLat' => "41.881944" )
	->json_is( '/response/0/orgServerFqdn' => "http://test-ds1.edge" )->json_is( '/response/0/checkPath' => "/crossdomain.xml" )
	->json_is( '/response/0/signed' => "0" )->json_is( '/response/0/longDesc1' => "test-ds1 long_desc_1" )->json_is( '/response/0/dscp' => "40" )
	->json_is( '/response/0/longDesc' => "test-ds1 long_desc" )->json_is( '/response/0/geoLimit' => "0" )->json_is( '/response/0/xmlId' => "test-ds1" )
	->json_is( '/response/0/globalMaxMbps' => "0" )->json_is( '/response/0/ccrDnsTtl' => "3600" )->json_is( '/response/0/active' => "1" )
	->json_is( '/response/0/profileDescription' => "ccr description" )->json_is( '/response/0/matchList/0/setNumber' => "0" )
	->json_is( '/response/0/matchList/0/type' => "HOST_REGEXP" )->json_is( '/response/0/matchList/0/pattern' => ".*\\.foo\\..*" )
	->json_is( '/response/0/type' => "EDGE" )->json_is( '/response/0/dnsBypassIp' => "" )->json_is( '/response/0/qstringIgnore' => "0" )
	->json_is( '/response/0/dnsBypassTtl' => undef )->json_is( '/response/0/profileName' => "CCR1" )->json_is( '/response/0/protocol' => "1" )
	->json_is( '/response/0/missLong' => "-87.627778" )->json_is( '/response/0/httpBypassFqdn' => "" )
	->json_is( '/response/0/infoUrl' => "http://test-ds1.edge/info_url" )->json_is( '/response/0/ipv6RoutingEnabled' => '1' )
	->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	"Validate test-ds1";

ok $t->get_ok('/api/1.1/deliveryservices/1/health.json')->status_is(200)->or( sub   { diag $t->tx->res->content->asset->{content}; } );
ok $t->get_ok('/api/1.1/deliveryservices/1/capacity.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->get_ok('/api/1.1/deliveryservices/1/routing.json')->status_is(200)->or( sub  { diag $t->tx->res->content->asset->{content}; } );
ok $t->get_ok('/api/1.1/deliveryservices/1/state.json')->status_is(200)->or( sub    { diag $t->tx->res->content->asset->{content}; } );

#ok $t->get_ok( '/api/1.1/deliveryservices/73/summary/stat/1423343701/1423602901/60/1423343701/1423602901.json' )->status_is(200)
#	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok(
	'/api/1.1/deliveryservices/73/edge/metric_types/kbps/start_date/1423343701/end_date/1423602901/interval/60/window_start/1423343701/window_end/1423602901.json'
)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.1/deliveryservices/73/server_types/mid/metric_types/kbps/start_date/1423343701/end_date/1423602901.json')->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_has("/response/0/stats")->json_has("/response/0/data");

ok $t->get_ok('/api/1.1/deliveryservices/73/server_types/edge/metric_types/tts/start_date/1423844560/end_date/1424103760.json')->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.1/deliveryservices/73/server_types/edge/metric_types/status_codes/start_date/1423844560/end_date/1424103760.json')->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.1/deliveryservices/73/server_types/edge/metric_types/origin_offload/start_date/1423844560/end_date/1424103760.json')->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.1/deliveryservices/73/server_types/edge/metric_types/origin_offload/start_date/1423844560/end_date/1424103760.json')->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.1/deliveryservices/73/server_types/XXX/metric_types/origin_offload/start_date/1423844560/end_date/1424103760.json')->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.1/deliveryservices/73/server_types/mid/metric_types/kbps/start_date/1423343701/end_date/1423602901.json?stats=true')->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_hasnt("/response/0/data");

ok $t->get_ok('/api/1.1/deliveryservices/73/server_types/mid/metric_types/kbps/start_date/1423343701/end_date/1423602901.json?data=true')->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_hasnt("/response/0/stats");

ok $t->post_ok('/api/1.1/user/logout')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# # logout
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

$dbh->disconnect();
done_testing();
