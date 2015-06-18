package API::Redis;
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
use UI::Server;
use JSON;
use LWP;
use Redis;
use Data::Dumper;
use Time::HiRes qw(gettimeofday tv_interval);
use Common::ReturnCodes qw(SUCCESS ERROR);

#TODO: drichardson - remove after 1.2 cleaned up
sub stats {
	my $self = shift;

	my $st = new Extensions::Delegate::Statistics($self);
	my ( $rc, $result ) = $st->get_usage_overview();
	if ( $rc == SUCCESS ) {
		$self->success($result);
	}
	else {
		$self->alert($result);
	}
}

sub info {
	my $self      = shift;
	my $host_name = $self->param('host_name');

	my $server = $self->db->resultset('Server')->search( { host_name => $host_name } )->single();
	if ( defined($server) ) {
		my $ip   = $server->ip_address;
		my $port = $server->tcp_port;

		my $redis = $self->redis_connect();

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
	else {
		$self->alert( { error => "Could not find host_name " . $host_name } );
	}
}

sub get_redis_stats {
	my $self  = shift;
	my $what  = $self->param('what');     # what statistic
	my $which = $self->param('which');    # which server, or cachegroup

	my $start = [gettimeofday];
	my $first = -8640;
	my $last  = -1;

	my $jdata = {};
	my $redis = $self->redis_connect();
	if ( defined($redis) ) {

		my $match = $which . " : " . $what;
		my @vals  = $redis->zrange( $match, $first, $last );
		my $e2    = tv_interval($start);

		#my $i = 0;
		$jdata->{which}    = $which;
		$jdata->{what}     = $what;
		$jdata->{interval} = " 10 seconds ";
		$jdata->{start}    = gmtime( ( split( /:/, $vals[0] ) )[0] );
		$jdata->{end}      = gmtime( ( split( /:/, $vals[$#vals] ) )[0] );
		$jdata->{number}   = $#vals;

		my $prev_tstamp = 0;
		my $j           = -1;

		#my $k = 0;
		foreach my $strval (@vals) {
			my ( $tstamp, $val ) = split( /:/, $strval );
			if ( $tstamp - $prev_tstamp != 10 ) {
				$j++;
				$jdata->{series}->[$j]->{time_base} = int($tstamp);

				#$i=0;
			}
			push( @{ $jdata->{series}->[$j]->{samples} }, int($val) );

			#$k++;
			#$i++;
			$prev_tstamp = $tstamp;
		}
		my $e3 = tv_interval($start);
		$jdata->{elapsed} = $e3 . ' (' . $e2 . ') ';
	}
	$self->render( json => $jdata );
}
1;
