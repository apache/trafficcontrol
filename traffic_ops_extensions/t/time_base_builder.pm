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
use Data::Dumper;
use strict;
use warnings;
use JSON;
use Extensions::InfluxDB::Builder::TimeBaseBuilder;

BEGIN { $ENV{MOJO_MODE} = "test" }

my $spdb_response = [
	{
		'label' => 'ipcdn_rascal#over-the-top#ipvod-cim#all#all#OutBps 300 raw',
		'data'  => [
			[ 1436455200, '86570057139.67' ],
			[ 1436455500, '88491753390' ],
			[ 1436455800, '89594576028' ],
			[ 1436456100, '89497839288.67' ],
			[ 1436456400, '91127984816.67' ],
			[ 1436456700, '92046368254.33' ],
			[ 1436457000, '92400794866.33' ],
			[ 1436457300, '93663994680.33' ],
			[ 1436457600, '95044216970.33' ],
			[ 1436457900, '94923578109.67' ],
			[ 1436458200, '96610808661.33' ],
			[ 1436458500, '97800858183.67' ],
		]
	}
];

my $builder = new Extensions::InfluxDB::Builder::TimeBaseBuilder(
	{
		cdn                      => 'cdn1',
		ds_name                  => 'ds1',
		cache_group_name         => 'all',
		host_name                => 'all',
		interval_for_metric_type => 300,
		interval                 => 60,
		time_elapsed             => 1.146632
	}
);
my ( $rc, $response ) = $builder->to_influx_series_format($spdb_response);
isnt( $response, undef );
#
#print "response #-> (" . Dumper($response) . ")\n";
my $series  = $response->{series};
my $samples = $response->{series}[0]->{samples};
is( @$samples, 55 );

done_testing();
