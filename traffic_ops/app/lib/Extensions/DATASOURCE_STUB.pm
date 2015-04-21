package Extensions::DATASOURCE_STUB;
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
use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use Time::HiRes qw(gettimeofday tv_interval);
use Math::Round qw(nearest);
use JSON;
use POSIX qw(strftime);

# NOTE!!! Please don't open source this file!!!

sub info {
	return {
		name        => "DATASOURCE_STUB",
		version     => "0.39",
		info_url    => "",
		description => "Stub for datasource",
		isactive    => 1,
		script_file => "Extensions::DATASOURCE_STUB.pm"
	};
}

sub stats_long_term {
	my $self     = shift;
	my $match    = shift;
	my $start    = shift;
	my $end      = shift;
	my $interval = shift;

	my $start_ts     = [gettimeofday];
	my ( $cdn, $ds, $cache_group_name, $host_name, $metric_name ) = split( /:/, $match );
	my $time_elapsed = tv_interval($start_ts);
	my $data         = {

		statName            => $metric_name,
		cdnName             => $cdn,
		deliveryServiceName => $ds,
		cacheGroupName      => $cache_group_name,
		hostName            => $host_name,
		elapsed             => $time_elapsed,
		start               => 1428599797,
		end                 => 1428599837,
		interval            => 10,
		series              => [
			{
				samples  => [0, 0, 0, 0],
				timeBase => 1428599797
			}
		]
	};

	return ($data);

	# return {};
}

sub get_config {
	my $self   = shift;
	my $metric = shift;

	my $metrics = {
		origin_tps => {
			interval      => 300,
			timeout       => 60,
			url           => "http://your-datasource.stub.kabletown.nearestt/chart",
			convert_to_ms => 0,
		}
	};

	return ($metrics);
}

1;
