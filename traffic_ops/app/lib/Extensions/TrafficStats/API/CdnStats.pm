package Extensions::TrafficStats::API::CdnStats;
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
#
#

# JvD Note: you always want to put Utils as the first use. Sh*t don't work if it's after the Mojo lines.
use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use JSON;
my $builder;
use Extensions::TrafficStats::Delegate::CdnStatistics;
use Utils::Helper::Extensions;
use Common::ReturnCodes qw(SUCCESS ERROR);

#TODO: drichardson
#      - Add required fields validation see lib/API/User.pm based on Validate::Tiny
sub index {
	my $self = shift;

	my $cstats = new Extensions::TrafficStats::Delegate::CacheStatistics( $self, $self->get_db_name() );
	my ( $rc, $result ) = $cstats->get_stats();
	if ( $rc == SUCCESS ) {
		return $self->success($result);
	}
	else {
		return $self->alert($result);
	}
}

sub get_usage_overview {
	my $self = shift;

	my $cstats = new Extensions::TrafficStats::Delegate::CdnStatistics(
		$self,
		$self->get_db_name("cache_stats_db_name"),
		$self->get_db_name("deliveryservice_stats_db_name")
	);

	my ( $rc, $result ) = $cstats->get_usage_overview();
	if ( $rc == SUCCESS ) {
		return $self->deprecation_with_no_alternative(200, $result);
	}
	else {
		return $self->deprecation_with_no_alternative(400, $result);
	}
}

sub get_db_name {
	my $self      = shift;
	my $key_name  = shift;
	my $mode      = $self->app->mode;
	my $conf_file = MojoPlugins::InfluxDB->INFLUXDB_CONF_FILE_NAME;
	my $conf      = Utils::JsonConfig->load_conf( $mode, $conf_file );
	return $conf->{$key_name};
}

1;
