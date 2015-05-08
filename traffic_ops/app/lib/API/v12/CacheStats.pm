package API::v12::CacheStats;
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

# JvD Note: you always want to put Utils as the first use. Sh*t don't work if it's after the Mojo lines.
use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use JSON;
use Extensions::Delegate::CacheStatistics;
use Utils::Helper::Extensions;
Utils::Helper::Extensions->use;
use Common::ReturnCodes qw(SUCCESS ERROR);

sub index {
	my $self = shift;

	my $cstats = new Extensions::Delegate::CacheStatistics( $self, $self->get_db_name() );

	my ( $rc, $result ) = $cstats->get_stats();
	if ( $rc == SUCCESS ) {
		return $self->success($result);
	}
	else {
		return $self->alert($result);
	}
}

sub get_db_name {
	my $self      = shift;
	my $mode      = $self->app->mode;
	my $conf_file = MojoPlugins::InfluxDB->INFLUXDB_CONF_FILE_NAME;
	my $conf      = Utils::JsonConfig->load_conf( $mode, $conf_file );
	return $conf->{cache_stats_db_name};
}

1;
