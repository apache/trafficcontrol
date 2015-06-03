package MojoPlugins::Metrics;
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

use Mojo::Base 'Mojolicious::Plugin';
use Carp qw(cluck confess);
use Data::Dumper;
use Math::Round qw(nearest);
use POSIX qw(strftime);
use Common::RedisFactory;
use Redis;
use Env;
use Extensions::DatasourceList;

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(
		get_config => sub {
			my $self   = shift;
			my $metric = shift;

			# get the subroutine name for the stats_config from the Extensions::DatasourceList
			my $ext          = new Extensions::DatasourceList();
			my $ext_hash_ref = $ext->hash_ref();
			my $subroutine   = $ext_hash_ref->{get_config};

			# and run it
			my $config = &{ \&{$subroutine} }( $self, $metric );
			return $config;
		}

	);

	$app->renderer->add_helper(
		redis_connection => sub {
			my $self = shift;

			my $redis_connection_string = $self->redis_connection_string();

			my $rm = Common::RedisFactory->new( $self, $redis_connection_string );
			return $rm->connection();
		}
	);

	$app->renderer->add_helper(
		etl_metrics => sub {
			my $self       = shift;
			my $metric     = $self->param("metric");
			my $start      = $self->param("start");         # start time in secs since 1970
			my $end        = $self->param("end");           # end time in secs since 1970
			my $stats_only = $self->param("stats") || 0;    # stats only
			my $data_only  = $self->param("data") || 0;     # data only
			my $type       = $self->param("type");

			my $config = $self->get_config($metric);
			if ( defined($config) ) {
				$start =~ s/\.\d+$//g;
				$end =~ s/\.\d+$//g;

				my $helper = new Utils::Helper::Datasource( { mojo => $self } );
				for my $kvp ( @{ $config->{get_kvp}->( $type, $start, $end ) } ) {
					$helper->kv( $kvp->{key}, $kvp->{value} );
				}

				$self->build_etl_metrics_response( $helper, $config, $start, $end, $stats_only, $data_only );
			}
			else {
				$self->internal_server_error( { " No configuration found for metric : " => $metric } );
			}
		}
	);
	$app->renderer->add_helper(
		build_etl_metrics_response => sub {
			my $self       = shift;
			my $helper     = shift;
			my $config     = shift;
			my $start      = shift;
			my $end        = shift;
			my $stats_only = shift;
			my $data_only  = shift;

			my $data = $helper->get_data( $config->{url}, $config->{convert_to_ms}, $config->{timeout} );

			if ( defined($data) ) {
				if ( exists( $config->{fixup} ) && ref( $config->{fixup} ) eq " CODE " ) {
					$config->{fixup}->($data);
				}

				$helper->pad_and_fill_holes( $data, $start, $end, $config->{interval} );
				$helper->calculate_stats( $data, $stats_only, $data_only );
				if (@$data) {
					$self->success($data);
				}
				else {
					$self->success( get_zero_values( $stats_only, $data_only ) );
				}

			}
			else {
				$self->internal_server_error( { " Internal Server " => " Error 1 " } );
			}
		}
	);
}

sub get_zero_values {
	my $stats_only = shift;
	my $data_only  = shift;

	my $response = ();
	$response->{"stats"}{"95thPercentile"} = 0;
	$response->{"stats"}{"98thPercentile"} = 0;
	$response->{"stats"}{"5thPercentile"}  = 0;
	$response->{"stats"}{"mean"}           = 0;
	$response->{"stats"}{"count"}          = 0;
	$response->{"stats"}{"min"}            = 0;
	$response->{"stats"}{"max"}            = 0;
	$response->{"stats"}{"sum"}            = 0;
	$response->{"data"}                    = [];
	if ($stats_only) {
		delete( $response->{"data"} );
	}
	elsif ($data_only) {
		delete( $response->{"stats"} );
	}
	return [$response];
}

1;
