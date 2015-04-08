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
use POSIX ();
use Mojo::Base -strict;
use Test::More;
use Test::Mojo;
use DBI;
use strict;
use warnings;
use Data::Dumper;
use warnings 'all';
use Schema;
use Test::TestHelper;
use Test::Mock::Redis;
use Common::RedisFactory;
use JSON;
use Test::MockModule;
no warnings 'once';

#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

BEGIN { $ENV{MOJO_MODE} = "test" }
my $schema = Schema->connect_to_database;
my $t      = Test::Mojo->new('TrafficOps');

my $rdf = Common::RedisFactory->new( $t, "redis01.kabletown.net:6379" );
my $redis = $rdf->connection();
$redis->sadd( "cdn1:all:all:all:maxKbps" => 26300000 );
$redis->rpush( "cdn2:tstamp"                => 1422376832 );
$redis->rpush( "cdn1:tstamp"                => 1422376832 );
$redis->rpush( "cdn2:all:all:all:kbps"      => 9188304.37 );
$redis->rpush( "cdn1:all:all:all:tps_total" => 18671.62 );

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);

#load core test data
Test::TestHelper->load_core_data($schema);

ok $t->post_ok( '/login', => form => { u => Test::TestHelper::PORTAL_USER, p => Test::TestHelper::PORTAL_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

my $cr_config = {
	'contentServers' => {
		'atlanta-edge-01' => {
			'fqdn'             => 'atlanta-edge-01.ga.atlanta.kabletown.net',
			'ip'               => '127.0.0.1',
			'ip6'              => '2001:558:FEEC::12/126',
			'status'           => 'REPORTED',
			'interfaceName'    => 'bond0',
			'hashId'           => 'atlanta-edge-01',
			'profile'          => 'EDGE1',
			'cacheGroup'       => 'mid-northeast-group',
			'type'             => 'EDGE',
			'locationId'       => 'mid-northeast',
			'deliveryServices' => {
				'test-ds1' => ['atlanta-edge01.ga.atlanta.kabletown.net'],
				'test-ds2' => [ 'edge.test-ds1.ga.atlanta.kabletown.net', 'edge.test-ds2.ga.atlanta.kabletown.net' ]
			},
			'port' => '80'
		},
		'atlanta-edge-02' => {
			'fqdn'             => 'atlanta-edge-02.ga.atlanta.kabletown.net',
			'ip'               => '127.0.0.1',
			'ip6'              => '2001:558:FEEC::12/126',
			'status'           => 'OFFLINE',
			'interfaceName'    => 'bond0',
			'hashId'           => 'atlanta-edge-02',
			'profile'          => 'EDGE1',
			'cacheGroup'       => 'mid-northeast-group',
			'type'             => 'EDGE',
			'locationId'       => 'mid-northeast',
			'deliveryServices' => {
				'test-ds1' => ['atlanta-edge02.ga.atlanta.kabletown.net'],
				'test-ds2' => [ 'edge.test-ds1.ga.atlanta.kabletown.net', 'edge.test-ds2.ga.atlanta.kabletown.net' ]
			},
			'port' => '80'
		},
		'atlanta-mid-01' => {
			'fqdn'             => 'atlanta-mid-02.ga.atlanta.kabletown.net',
			'ip'               => '127.0.0.2',
			'ip6'              => '2001:558:FEEC::12/126',
			'status'           => 'ONLINE',
			'interfaceName'    => 'bond0',
			'hashId'           => 'atlanta-mid-01',
			'profile'          => 'MID1',
			'cacheGroup'       => 'mid-northwest-group',
			'type'             => 'MID',
			'locationId'       => 'mid-northwest',
			'deliveryServices' => {
				'test-ds1' => ['atlanta-mid-01.ga.atlanta.kabletown.net'],
				'test-ds2' => [ 'edge.test-ds1.ga.atlanta.kabletown.net', 'edge.test-ds2.ga.atlanta.kabletown.net' ]
			},
			'port' => '80'
		},
		'atlanta-mid-02' => {
			'fqdn'             => 'atlanta-mid-02.ga.atlanta.kabletown.net',
			'ip'               => '127.0.0.2',
			'ip6'              => '2001:558:FEEC::12/126',
			'status'           => 'ONLINE',
			'interfaceName'    => 'bond0',
			'hashId'           => 'atlanta-mid-02',
			'profile'          => 'MID1',
			'cacheGroup'       => 'mid-northwest-group',
			'type'             => 'EDGE',
			'locationId'       => 'mid-northwest',
			'deliveryServices' => {
				'test-ds1' => ['atlanta-mid-02.ga.atlanta.kabletown.net'],
				'test-ds2' => [ 'edge.test-ds1.ga.atlanta.kabletown.net', 'edge.test-ds2.ga.atlanta.kabletown.net' ]
			},
			'port' => '80'
		},
	},
	"config" => {
		"geolocation.polling.url"       => "https://to.server.net/MaxMind/GeoLiteCity.dat.gz",
		"geolocation6.polling.url"      => "https://to.server.net/MaxMind/GeoLiteCityv6.dat.gz",
		"geolocation.polling.interval"  => "86400000",
		"geolocation6.polling.interval" => "86400000",
		"domain_name"                   => "edge.server.net",
		"ttls"                          => {
			"AAAA" => "3600",
			"SOA"  => "86400",
			"A"    => "3600",
			"NS"   => "3600"
		},
		"soa" => {
			"expire"  => "604800",
			"minimum" => "86400",
			"admin"   => "admin",
			"retry"   => "7200",
			"refresh" => "28800"
		},
		"coveragezone.polling.url"      => "http://to.server.net/ipcdn/CZF/current/czf.json",
		"coveragezone.polling.interval" => "86400000"
	},
	"edgeLocations" => {
		"atlanta-edge-01" => { "longitude" => 1, "latitude" => 2 },
		"atlanta-mid-01"  => { "longitude" => 1, "latitude" => 2 },
		"atlanta-edge-02" => { "longitude" => 1, "latitude" => 2 },
		"atlanta-mid-02"  => { "longitude" => 1, "latitude" => 2 },
	},
};

my $cr_states = {
	"caches" => {
		"atlanta-edge-01" => { "isAvailable" => "true" },
		"atlanta-mid-01"  => { "isAvailable" => "false" },
		"atlanta-edge-02" => { "isAvailable" => "false" },
		"atlanta-mid-02"  => { "isAvailable" => "false" },
	},
	"deliveryServices" => {
		"test-ds1" => {
			"disabledLocations" => [],
			"isAvailable"       => "true"
		},
		"test-ds2" => {
			"disabledLocations" => [],
			"isAvailable"       => "false"
		},
	},
};

# Why not more cachegroups?
# {"version":"1.1","response":{"cachegroups":[{"name":"mid-northeast","offline":0,"online":4}],"totalOnline":4,"totalOffline":0}}
my $rascal_util = new Test::MockModule( 'Utils::Rascal', no_auto => 1 );
$rascal_util->mock( 'get_cr_config' => sub { return $cr_config } );
$rascal_util->mock( 'get_cr_states' => sub { return $cr_states } );

ok $t->get_ok('/api/1.1/cdns/usage/overview.json')->json_is( "/response/maxGbps", 22.355 )->json_is( "/response/currentGbps", 9.188304 )
	->json_is( "/response/tps", 18671 )->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

use Utils::Helper::Datasource;
my $spdb_util = new Test::MockModule( 'Utils::Helper::Datasource', no_auto => 1 );
my @spdb_response = {
	"stats" => {
		'sum'            => '1744886.99',
		'max'            => '3324.63666666667',
		'min'            => '754.77',
		'5thPercentile'  => '2102.03666666667',
		'count'          => 900,
		'mean'           => '2017.21039306358',
		'95thPercentile' => '3186.65666666667',
		'98thPercentile' => '3242.84333333333'
	},
	"data" => [ [ 1423576800000, 285791 ], [ 1423577100000, 284489 ], [ 1423577400000, 285495 ], [ 1423577700000, 294863 ] ],
	'label'  => '/VOD/ODOL/Delivery_services/mid_all > httpCacheMisses 300 raw > vsum',
	'period' => 300
};

$spdb_util->mock( 'get_data' => sub { return \@spdb_response } );

#api/1.1/metrics/g/kbps/1423343701/1423602901.json
# Mocked out by Utils::Helper::Datasource
ok $t->get_ok('/api/1.1/metrics/server_types/mid/metric_types/origin_tps/start_date/1423343701/end_date/1423602901.json')->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

# Used on the /dailysummary page viewed by EVERYONE if this route changes then /templates/visual_status/daily_summary.html.ep has to change accordingly.
# Mocked out by Utils::Helper::Datasource
ok $t->get_ok('/api/1.1/cdns/peakusage/daily/deliveryservice/all/cachegroup/all/start_date/0/end_date/now/interval/86400.json')->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.1/cdns/usage/overview.json')->json_is( "/response/maxGbps", 22.355 )->json_is( "/response/currentGbps", 9.188304 )
	->json_is( "/response/tps", 18671 )->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->get_ok('/api/1.1/cdns/health.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# Example for mocking out a controller
#my $cdn_controller = new Test::MockModule( 'API::Cdn', no_auto => 1 );
#$cdn_controller->mock( 'capacity' => sub { return 'X' } );
ok $t->get_ok('/api/1.1/cdns/capacity.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/api/1.1/cdns/routing.json')->status_is(200)->or( sub                 { diag $t->tx->res->content->asset->{content}; } );
ok $t->get_ok('/api/1.1/cdns/configs.json')->status_is(200)->or( sub                 { diag $t->tx->res->content->asset->{content}; } );
ok $t->get_ok('/api/1.1/cdns/cdn2/configs/monitoring.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->get_ok('/api/1.1/cdns/cdn1/configs/routing.json')->status_is(200)->or( sub    { diag $t->tx->res->content->asset->{content}; } );

@spdb_response = {
	"data" => [ [ '1424179800000', 287592 ], [ '1424180100000', 293751 ], [ '1424180400000', 303143 ], [ '1424180700000', 308717 ] ],
	'label'  => '/VOD/ODOL/Delivery_services/mid_all > httpCacheMisses 300 raw > vsum',
	'period' => 300
};

$spdb_util->mock( 'get_data' => sub { return \@spdb_response } );
ok $t->get_ok('/api/1.1/cdns/metric_types/origin_tps/start_date/1424205299/end_date/1424206199.json')->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

##DNSSEC KEYS TESTS
ok $t->post_ok( '/api/1.1/user/login', json => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(200),
	'Log into the admin user?';

# create dnssec key
my $key             = "cdn1";
my $name            = "foo.com";
my $ttl             = "60";
my $expiration_days = "365";

my $dnssec_keys = {
	$key => {
		ksk => {
			ttl            => $ttl,
			public         => "public",
			name           => $name . ".",
			inceptionDate  => 1426173464,
			private        => "private key",
			expirationDate => 1457709464
		},
		zsk => {
			private        => "private",
			expirationDate => 1428765464,
			inceptionDate  => 1426173464,
			name           => $name . ".",
			ttl            => $ttl,
			public         => "public"
		},
	},
	"test-ds2" => {
		ksk => {
			ttl            => $ttl,
			public         => "public",
			name           => $name . ".",
			inceptionDate  => 1426173464,
			private        => "private key",
			expirationDate => 1457709464
		},
		zsk => {
			private        => "private",
			expirationDate => 1428765464,
			inceptionDate  => 1426173464,
			name           => $name . ".",
			ttl            => $ttl,
			public         => "public"
		}
	},
	"test-ds1" => {
		ksk => {
			ttl            => $ttl,
			public         => "public",
			name           => $name . ".",
			inceptionDate  => 1426173464,
			private        => "private key",
			expirationDate => 1457709464
		},
		zsk => {
			private        => "private",
			expirationDate => 1428765464,
			inceptionDate  => 1426173464,
			name           => $name . ".",
			ttl            => $ttl,
			public         => "public"
		}
	}
};

my $fake_lwp = new Test::MockModule( 'LWP::UserAgent', no_auto => 1 );

my $fake_get = HTTP::Response->new( 200, undef, HTTP::Headers->new, encode_json($dnssec_keys) );
$fake_lwp->mock( 'get', sub { return $fake_get } );

my $fake_put = HTTP::Response->new( 204, undef, HTTP::Headers->new, undef );
$fake_lwp->mock( 'put', sub { return $fake_put } );

my $fake_delete = HTTP::Response->new( 204, undef, HTTP::Headers->new, undef );
$fake_lwp->mock( 'delete', sub { return $fake_delete } );

ok $t->post_ok(
	'/api/1.1/cdns/dnsseckeys/generate',
	json => {
		key               => $key,
		name              => $name,
		ttl               => $ttl,
		kskExpirationDays => $expiration_days,
		zskExpirationDays => "90",
	}
	)->status_is(200)->json_is( "/response" => "Successfully created dnssec keys for $key" )->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Create dnssec key?';

#validate dnssec zone signing key exist
ok $t->get_ok("/api/1.1/cdns/name/$key/dnsseckeys.json")->json_has("/response")->json_has("/response/$key/ksk/private")
	->json_has("/response/$key/zsk/public")->json_has("/response/$key/ksk/expirationDate")->json_is( "/response/$key/zsk/name" => $name . "." )
	->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

#delete dnssec key
ok $t->get_ok("/api/1.1/cdns/name/$key/dnsseckeys/delete.json")->json_is( "/response" => "Successfully deleted dnssec keys for $key" )->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

#NEGATIVE TESTING -- record doesnt exist.
#get_object
my $fake_get_404 = HTTP::Response->new( 404, undef, HTTP::Headers->new, "Not found" );
$fake_lwp->mock( 'get', sub { return $fake_get_404 } );

ok $t->get_ok("/api/1.1/cdns/name/foo/dnsseckeys.json")->status_is(400)
	->json_is( "/alerts/0/text" => "Error  - Dnssec keys for foo do not exist!  Response was: Not found" )
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

#NEGATIVE TESTING -- no user privs
ok $t->post_ok( '/api/1.1/user/login', json => { u => Test::TestHelper::PORTAL_USER, p => Test::TestHelper::PORTAL_USER_PASSWORD } )->status_is(200),
	'Log into the portal user?';

#create
ok $t->post_ok(
	'/api/1.1/cdns/dnsseckeys/generate',
	json => {
		keyType => "dnssec",
		key     => $key,
	}
	)->status_is(400)->json_has("Error - You do not have permissions to perform this operation!")
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

#get_object
ok $t->get_ok("/api/1.1/cdns/name/$key/dnsseckeys.json")->status_is(400)->json_has("Error - You do not have permissions to perform this operation!")
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

#delete
ok $t->get_ok("/api/1.1/cdns/name/$key/dnsseckeys/delete.json")->status_is(400)->json_has("Error - You do not have permissions to perform this operation!")
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

# logout
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

done_testing();
