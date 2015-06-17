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
Utils::Helper::Extensions->use;
use Common::ReturnCodes qw(SUCCESS ERROR);

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(
		etl_metrics => sub {
			my $self       = shift;
			my $metric     = $self->param("metric");
			my $start      = $self->param("start");         # start time in secs since 1970
			my $end        = $self->param("end");           # end time in secs since 1970
			my $stats_only = $self->param("stats") || 0;    # stats only
			my $data_only  = $self->param("data") || 0;     # data only
			my $type       = $self->param("type");

			my $m = new Extensions::Delegate::Metrics(
				{
					metricType => $metric,
					startDate  => $start,
					endDate    => $end,
					statsOnly  => $stats_only,
					dataOnly   => $data_only,
					type       => $type
				}
			);
			my ( $rc, $result ) = $m->get_etl_metrics($self);

			if ( $rc == SUCCESS ) {
				return $self->success($result);
			}
			else {
				return $self->alert($result);
			}
		}
	);
}
1;
