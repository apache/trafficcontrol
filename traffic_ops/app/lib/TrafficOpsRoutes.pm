package TrafficOpsRoutes;
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

#
#

sub new {
	my $self  = {};
	my $class = shift;
	return ( bless( $self, $class ) );
}

sub define {
	my $self = shift;
	my $r    = shift;

	my $api_version = "1.2";
	$self->define_influx_routes( $r, $api_version );
}

sub define_influx_routes {
	my $self        = shift;
	my $r           = shift;
	my $api_version = shift;

	$r->get( "/api/$api_version/cdns/usage/overview" => [ format => [qw(json)] ] )
		->to( 'CdnStats#get_usage_overview', namespace => "Extensions::InfluxDB::API" );
	$r->get( "/api/$api_version/deliveryservice_stats" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'DeliveryServiceStats#index', namespace => "Extensions::InfluxDB::API" );
	$r->get( "/api/$api_version/cache_stats" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'CacheStats#index', namespace => "Extensions::InfluxDB::API" );

}

1;
