package Utils::Helper::Datasource;
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

use Carp qw(cluck confess);
use JSON;
use Data::Dumper;
use Utils::Helper;
use Mojo::UserAgent;

our @ISA = ("Utils::Helper");    # inherit our constructor and mojo methods

use constant { TIMEOUT => 30, };

sub kv {
	my $self  = shift || confess("Call on an instance of Utils::Helper::Datasource");
	my $key   = shift;
	my $value = shift;                                                                  # note: can be undef to push a query parameter on with no value

	if ( defined($key) ) {
		if ( !exists( $self->{KV}->{$key} ) ) {
			$self->{KV}->{$key} = [];
		}

		if ( ref($value) eq "ARRAY" ) {
			push( @{ $self->{KV}->{$key} }, @{$value} );
		}
		else {
			push( @{ $self->{KV}->{$key} }, $value );
		}
	}

	return ( $self->{KV} );
}

sub kp {
	my $self = shift || confess("Call on an instance of Utils::Helper::Datasource");

	my $kv = $self->kv();

	if ( !defined($kv) || scalar( keys( %{$kv} ) ) == 0 ) {
		confess("You must supply at least one kv pair");
	}

	my @kp;

	for my $key ( %{$kv} ) {
		for my $value ( @{ $kv->{$key} } ) {
			if ( defined($value) ) {
				push( @kp, sprintf( "%s=%s", $key, $value ) );
			}
			else {
				push( @kp, $key );
			}
		}
	}

	return (@kp);
}

sub get_data {
	my $self          = shift || confess("Call on an instance of Utils::Helper::Datasource");
	my $base_url      = shift || confess("Supply a base URL");
	my $convert_to_ms = shift || 0;
	my $timeout       = shift;

	if ( !defined($timeout) ) {    # could be 0 or > 0
		$timeout = TIMEOUT;
	}

	my $url = sprintf( "%s?%s", $base_url, join( "&", $self->kp() ), );

	my $ua = Mojo::UserAgent->new;
	$ua->request_timeout($timeout);
	$ua->inactivity_timeout($timeout);

	my $result = $ua->get($url)->res->content->asset->slurp;

	if ( defined($result) && $result ne "" ) {
		my $j = undef;

		# decode_json will die() when it's unable to parse $result, so wrap in an eval
		eval { $j = decode_json($result); };

		if ( !$@ && defined($j) && ref($j) eq "ARRAY" ) {
			for my $chunk ( @{$j} ) {
				if ( exists( $chunk->{data} ) && $convert_to_ms ) {
					for ( my $i = 0; $i < scalar( @{ $chunk->{data} } ); $i++ ) {
						if ($convert_to_ms) {

							# multiply time by 1000 to convert to MS
							$chunk->{data}->[$i]->[0] *= 1000;
						}
					}
				}

				if ( exists( $chunk->{series} ) ) {
					delete( $chunk->{series} );
				}
			}
		}

		return ($j);
	}
	else {
		print "Utils::Helper::Datasource::get_data Cannot reach url: " . $url . "\n";
		return (undef);
	}
}

sub pad_and_fill_holes {
	my $self     = shift || confess("Call on an instance of Utils::Helper::Datasource");
	my $data     = shift || confess("Supply an Datasource REST API data structure");
	my $start    = shift || confess("Supply a window start");
	my $end      = shift || confess("Supply a window end");
	my $interval = shift || confess("Supply an interval");

	$start    *= 1000;
	$end      *= 1000;
	$interval *= 1000;

	for my $series ( @{$data} ) {
		if ( exists( $series->{data} ) ) {

			# this all operates under the assumption that our deltas and periods are in MS
			my $ft = undef;
			my $lt = undef;
			my @splice_gaps;

			for ( my $i = 0; $i < scalar( @{ $series->{data} } ); $i++ ) {
				my $ct = $series->{data}->[$i]->[0];

				if ( !defined($ft) ) {
					$ft = $ct;
				}

				if ( !defined($lt) ) {
					$lt = $ct;
				}

				my $delta = $ct - $lt;

				# only fill gaps if we have a delta that divides cleanly
				if ( $delta && $delta != $interval && $delta % $interval == 0 ) {
					my @tdp;
					my $gap_size = $delta / $interval;

					for ( my $j = 1; $j < $gap_size; $j++ ) {
						my $gap_time = $lt + ( $interval * $j );
						push( @tdp, [ $gap_time, undef ] );
					}

					push( @splice_gaps, { index => $i, data => \@tdp } );
				}

				$lt = $ct;
			}

			# fill the gaps in the middle
			my $splice_offset = 0;

			for my $splice (@splice_gaps) {
				splice( @{ $series->{data} }, $splice->{index} + $splice_offset, 0, @{ $splice->{data} } );
				$splice_offset += scalar( @{ $splice->{data} } );
			}

			# null pad the front
			while ( $ft > $start ) {
				$ft -= $interval;
				unshift( @{ $series->{data} }, [ $ft, undef ] );
			}

			# null pad the end
			while ( $lt < $end ) {
				$lt += $interval;
				push( @{ $series->{data} }, [ $lt, undef ] );
			}
		}
	}
}

sub calculate_stats {
	my $self = shift || confess("Call on an instance of Utils::Helper::Datasource");
	my $data = shift || confess("Supply an Datasource REST API data structure");
	my $stats_only = shift;        # this can/will be zero
	my $data_only  = shift;        # this can/will be zero
	my $override   = shift || 0;

	for my $series ( @{$data} ) {
		if ( !exists( $series->{stats} ) || $override ) {
			my $stats = {
				count            => 0,
				min              => undef,
				max              => 0,
				mean             => 0,
				sum              => 0,
				"98thPercentile" => 0,
				"95thPercentile" => 0,
				"5thPercentile"  => 0,
				samples          => []
			};

			for my $item ( @{ $series->{data} } ) {
				my $sample = $item->[1];

				if ( !defined($sample) ) {
					next;
				}

				if ( !defined( $stats->{min} ) || $sample < $stats->{min} ) {
					$stats->{min} = $sample;
				}

				if ( $sample > $stats->{max} ) {
					$stats->{max} = $sample;
				}

				$stats->{sum} += $sample;
				$stats->{count}++;
				push( @{ $stats->{samples} }, $sample );
			}

			my @sorted  = sort { $a <=> $b } @{ $stats->{samples} };
			my $index98 = ( scalar(@sorted) * .98 ) - 1;
			my $index95 = ( scalar(@sorted) * .95 ) - 1;
			my $index5  = ( scalar(@sorted) * .5 ) - 1;

			$stats->{"98thPercentile"} = $sorted[$index98];
			$stats->{"95thPercentile"} = $sorted[$index95];
			$stats->{"5thPercentile"}  = $sorted[$index5];

			if ( $stats->{sum} > 0 ) {
				if ( $stats->{count} > 1 ) {
					$stats->{mean} = $stats->{sum} / $stats->{count};
				}
				else {
					$stats->{mean} = $stats->{sum};
				}
			}

			delete( $stats->{samples} );

			$series->{stats} = $stats;
		}

		if ($stats_only) {
			delete( $series->{data} );
		}
		elsif ($data_only) {
			delete( $series->{stats} );
		}
	}

	return ($data);
}

1;
