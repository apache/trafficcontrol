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
use Test::MockModule;
use Test::TestHelper;
use JSON;
use Mojo::UserAgent;
use strict;
use warnings;
use Extensions::Helper::Datasource;
use Common::ReturnCodes qw(SUCCESS ERROR);

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

my $fake_answer = [
	{
		"label" => "vod_cdsis#SERVICES#cim-jitp#EDGE#cdnTTS00 300 raw",
		"data"  => [
			[ 1434918300, 12.0 ],
			[ 1434918600, 14.0 ],
			[ 1434918900, 21.0 ],
			[ 1434919200, 8.0 ],
			[ 1434919500, 8.0 ],
			[ 1434919800, 17.0 ],
			[ 1434920100, 24.0 ],
			[ 1434920400, 2.0 ],
			[ 1434920700, 15.0 ],
			[ 1434921000, 11.0 ],
			[ 1435001100, 1.0 ],
			[ 1435001700, 7.0 ],
			[ 1435002000, 9.0 ],
			[ 1435002300, 2.0 ],
			[ 1435002600, 6.0 ]
		],
		"stats" => {
			"count"          => 233,
			"sum"            => 2699.0,
			"min"            => 1.0,
			"max"            => 61.0,
			"mean"           => 11.583690987124463,
			"stddev"         => 9.667894349106154,
			"median"         => 8.0,
			"98thPercentile" => 36.0,
			"95thPercentile" => 33.0,
			"5thPercentile"  => 2.0
		}
	}
];

#my $ua      = Mojo::UserAgent->new;
#my $ds = new Test::MockModule('Extensions::Helper::Datasource');

my $ds = Test::MockModule->new('Extensions::Helper::Datasource');

my $start_date           = time();
my $end_date             = $start_date + 200;
my $etl_metrics_response = [
	{
		"label" => "EDGE TPS",
		"data"  => [
			[ 1434918300, 12.0 ],
			[ 1434918600, 14.0 ],
			[ 1434918900, 21.0 ],
			[ 1434919200, 8.0 ],
			[ 1434919500, 8.0 ],
			[ 1434919800, 17.0 ],
			[ 1434920100, 24.0 ],
			[ 1434920400, 2.0 ],
			[ 1434920700, 15.0 ],
			[ 1434921000, 11.0 ],
			[ 1435001100, 1.0 ],
			[ 1435001700, 7.0 ],
			[ 1435002000, 9.0 ],
			[ 1435002300, 2.0 ],
			[ 1435002600, 6.0 ]
		],
		"stats" => {
			"count"          => 233,
			"sum"            => 2699.0,
			"min"            => 1.0,
			"max"            => 61.0,
			"mean"           => 11.583690987124463,
			"stddev"         => 9.667894349106154,
			"median"         => 8.0,
			"98thPercentile" => 36.0,
			"95thPercentile" => 33.0,
			"5thPercentile"  => 2.0
		}
	}
];

my $delegate_metrics = Test::MockModule->new('Extensions::Delegate::Metrics');

$delegate_metrics->mock( 'build_etl_metrics_response', sub { SUCCESS, $etl_metrics_response } );

ok $t->get_ok("/api/1.1/deliveryservices/1/server_types/edge/metric_types/tts/start_date/$start_date/end_date/$end_date.json")->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.1/deliveryservices/1/server_types/edge/metric_types/status_codes/start_date/1423844560/end_date/1424103760.json')->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.1/deliveryservices/1/server_types/mid/metric_types/kbps/start_date/1423343701/end_date/1423602901.json')->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_has("/response/0/stats")->json_has("/response/0/data");

ok $t->get_ok(
	"/api/1.1/deliveryservices/1/edge/metric_types/kbps/start_date/$start_date/end_date/$end_date/interval/60/window_start/$start_date/window_end/$end_date.json"
)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.1/deliveryservices/1/health.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->get_ok('/api/1.1/deliveryservices.json')->status_is(200)->json_has( '/response', 'has a response' ), "Get delivery services";
ok $t->get_ok('/api/1.1/deliveryservices/1/capacity.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->get_ok('/api/1.1/deliveryservices/1/routing.json')->status_is(200)->or( sub  { diag $t->tx->res->content->asset->{content}; } );
ok $t->get_ok('/api/1.1/deliveryservices/1/state.json')->status_is(200)->or( sub    { diag $t->tx->res->content->asset->{content}; } );
ok $t->get_ok('/api/1.1/deliveryservices.json')->status_is(200)->json_is( '/response/0/id' => 1 ), "Has a delivery service";

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

ok $t->get_ok('/api/1.1/deliveryservices/1/server_types/edge/metric_types/origin_offload/start_date/1423844560/end_date/1424103760.json')->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.1/deliveryservices/1/server_types/edge/metric_types/origin_offload/start_date/1423844560/end_date/1424103760.json')->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.1/deliveryservices/1/server_types/XXX/metric_types/origin_offload/start_date/1423844560/end_date/1424103760.json')->status_is(400)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

#$delegate_metrics->mock(
#	'get_config',
#	sub {
#		my $config = {
#			url           => undef,
#			interval      => 300,
#			timeout       => 60,
#			convert_to_ms => 0,
#			get_kvp       => sub { [ { key => 'metric', value => 'cdnTTS0p8,300,raw' } ], [ { key => 'start_time', value => $start_date } ] },
#			fixup => sub { label => 'ORIGIN_TPS' },
#		};
#		return $config;
#	}
#);
#
# # logout
ok $t->post_ok('/api/1.1/user/logout')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

$dbh->disconnect();
done_testing();
