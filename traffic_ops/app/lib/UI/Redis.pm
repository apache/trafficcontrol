package UI::Redis;
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
use TrafficOps;
use Mojo::Base 'Mojolicious::Controller';
use JSON;
use Redis;
use Data::Dumper;
use Time::HiRes qw(gettimeofday tv_interval);
use Math::Round qw(nearest);
use Utils::Helper;
use Carp qw(cluck confess);

sub stats {
	my $self     = shift;                       # /redis/#match/#start/#end/#interval
	my $match    = $self->param('match');
	my $start    = $self->param('start');       # start time in secs since 1970, or "now" to get latest sample
	my $end      = $self->param('end');         # end time in secs since 1970, or "now" to get latest sample
	my $interval = $self->param('interval');    # the interval between the samples. 10 is minimum, has to be a multiple of 10

	my $j = $self->v11_get_stats( $match, $start, $end, $interval );
	$self->render( json => $j );
}

sub info {
	my $self      = shift;
	my $shortname = $self->param('shortname');

	my $server = $self->db->resultset('Server')->search( { host_name => $shortname } )->single();
	my $ip     = $server->ip_address;
	my $port   = $server->tcp_port;

	my $redis = redis_connect();

	my $data = undef;
	foreach my $sub (qw/Server Clients Memory Persistence Stats Replication CPU Keyspace/) {
		my $subdata = $redis->info($sub);
		$data->{$sub} = $subdata;
	}

	my $i = 0;
	my @slowlist = $redis->slowlog( "get", 1024 );
	foreach my $slowlog (@slowlist) {
		push( @{ $data->{slowlog} }, $slowlog );
		$i++;
	}
	$data->{slowlen} = $i;

	$redis->quit();
	$self->render( json => $data );
}

1;
