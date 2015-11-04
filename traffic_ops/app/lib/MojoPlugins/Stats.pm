package MojoPlugins::Stats;
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
use JSON;
use Utils::CCR;
use Time::HiRes qw(gettimeofday tv_interval);
use Math::Round qw(nearest);

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(

		# start remove when redis is depracated
		get_stats => sub {
			my $self     = shift;
			my $match    = shift;
			my $start    = shift;
			my $end      = shift;
			my $interval = shift;

			# these arguments allow us to grab small windows for summary data while retaining the larger window and short/long term logic below
			my $window_start = shift || $start;
			my $window_end   = shift || $end;

			# remove any trailing .XXX from the times sent to us from angular
			for my $var ( \$start, \$end, \$window_start, \$window_end ) {
				${$var} =~ s/\.\d+$//g;
			}

			my $data;
			my $default_retention_period = 86400;    # one day

			my $retention_period =
				   $self->db->resultset('Parameter')->search( { name => "RetentionPeriod", config_file => "redis.config" } )->get_column('value')->single()
				|| $default_retention_period;

			# numeric start/end only which should be done upstream but let's be extra cautious
			if ( $start =~ /^\d+$/ && $end =~ /^\d+$/ && $window_start < ( time() - $retention_period - 60 ) ) {  # -60 for diff between client and our time
				$data = $self->get_stats_long_term( $match, $start, $end, $interval );
			}
			else {
				# get_usage uses now/now as start/end, so it will pass through to short_term
				$data = $self->get_stats_short_term( $match, $start, $end, $interval );
			}

			if ( defined($data) ) {
				if ( ref($data) eq "HASH" && exists( $data->{series} ) ) {
					$self->normalize_intervals( $data, $interval );
					$self->calc_summary($data);
				}

				return ($data);
			}
			else {
				return (undef);
			}

		}
	);
	$app->renderer->add_helper(
		calc_summary => sub {
			my $self = shift;
			my $data = shift;

			my $interval = $data->{interval} || return (undef);
			my $stat     = $data->{statName} || return (undef);

			my $convert = {
				kbps => sub {
					my $t = shift;
					my $i = shift;
					return ( ( $t / 8 ) * $i );
				},
				tps => sub {
					my $t = shift;
					my $i = shift;
					return ( $t * $i );
				},
				tps_2xx => sub {
					my $t = shift;
					my $i = shift;
					return ( $t * $i );
				},
				tps_3xx => sub {
					my $t = shift;
					my $i = shift;
					return ( $t * $i );
				},
				tps_4xx => sub {
					my $t = shift;
					my $i = shift;
					return ( $t * $i );
				},
				tps_5xx => sub {
					my $t = shift;
					my $i = shift;
					return ( $t * $i );
				},
				tps_total => sub {
					my $t = shift;
					my $i = shift;
					return ( $t * $i );
				},
			};

			my $summary = {
				min         => undef,
				max         => 0,
				average     => 0,
				ninetyFifth => 0,
				total       => 0,
				samples     => []
			};

			for my $series ( @{ $data->{series} } ) {
				for my $sample ( @{ $series->{samples} } ) {
					if ( !defined($sample) ) {
						next;
					}

					if ( !defined( $summary->{min} ) || $sample < $summary->{min} ) {
						$summary->{min} = $sample;
					}

					if ( $sample > $summary->{max} ) {
						$summary->{max} = $sample;
					}

					$summary->{total} += $sample;
					push( @{ $summary->{samples} }, $sample );
				}
			}

			my @sorted = sort { $a <=> $b } @{ $summary->{samples} };
			my $index = ( scalar(@sorted) * .95 ) - 1;    # calc the index of the 95th percentile, subtract one for real index
			$summary->{ninetyFifth} = $sorted[$index];

			if ( $summary->{total} ) {
				if ( scalar( @{ $summary->{samples} } ) > 1 ) {
					$summary->{average} = int( $summary->{total} / scalar( @{ $summary->{samples} } ) );
				}
				else {
					$summary->{average} = $summary->{total};
				}

				if ( exists( $convert->{$stat} ) && $convert->{$stat} ) {
					$summary->{total} = $convert->{$stat}->( $summary->{total}, $interval );
				}
			}

			delete( $summary->{samples} );

			$data->{summary} = $summary;

		}
	);
	$app->renderer->add_helper(
		normalize_intervals => sub {
			my $self     = shift;
			my $data     = shift;
			my $interval = shift;

			# add keys that are "per second" metrics which require special handling for normalization
			my $ps_metrics = {
				kbps      => 1,
				tps       => 1,
				tps_2xx   => 1,
				tps_3xx   => 1,
				tps_4xx   => 1,
				tps_5xx   => 1,
				tps_total => 1,
			};

			if ( $data->{interval} > $interval && $data->{interval} % $interval == 0 ) {
				for my $series ( @{ $data->{series} } ) {
					for my $sample ( @{ $series->{samples} } ) {
						my $slice = $data->{interval} / $interval;

						if ( defined($sample) && !exists( $ps_metrics->{ $data->{statName} } ) ) {
							$sample = $sample / $slice;
						}

						for ( my $i = 0; $i < $slice; $i++ ) {
							push( @{ $series->{new_samples} }, $sample );
						}

					}

					$series->{samples} = delete( $series->{new_samples} );
				}

				$data->{interval} = $interval;
			}
			elsif ( $data->{interval} < $interval && $interval % $data->{interval} == 0 ) {
				for my $series ( @{ $data->{series} } ) {
					my $span    = $interval / $data->{interval};
					my $sum     = 0;
					my $counter = 0;

					for my $sample ( @{ $series->{samples} } ) {
						$counter++;

						if ( defined($sample) ) {
							$sum += $sample;
						}

						if ( $counter == $span ) {
							if ( exists( $ps_metrics->{ $data->{statName} } ) ) {
								$sum = $sum / $counter;
							}

							push( @{ $series->{new_samples} }, $sum );
							$sum     = 0;
							$counter = 0;
						}
					}

					$series->{samples} = delete( $series->{new_samples} );
				}

				$data->{interval} = $interval;
			}

			return ($data);

		}
	);
	$app->renderer->add_helper(
		get_stats_short_term => sub {
			my $self     = shift;
			my $match    = shift;
			my $start    = shift;
			my $end      = shift;
			my $interval = shift;

			my $REDIS_METRIC_INTERVAL = 10;

			my $redis = $self->redis_connect();

			# my $ret = $redis->ping();
			# if ( !defined($ret) || $ret ne "PONG" ) {
			# 	$self->app->log->info("No response from ping: " . Dumper($ret));
			# }
			my $last_one = 0;
			if ( $end eq "now" && $start eq "now" ) {
				$end      = time();
				$start    = time();
				$last_one = 1;
			}

			if ( $interval < 10 ) {    # note 10 is the minimum interval, and the hardcoded interval of getting samples into redis.
				$interval = 10;        # minimum
			}
			my $startts = [gettimeofday];

			my $max_kbps_match = $match;
			my @parts          = split( /:/, $match );
			my $cdn_name       = $parts[0];
			my $ds_name        = $parts[1];
			my $cg_name        = $parts[2];
			my $host_name      = $parts[3];
			my $stat_name      = $parts[4];
			my $capacity;

			if ( $stat_name =~ /daily_/ ) {
				$capacity = 0;
			}
			elsif ( $ds_name ne "all" ) {
				$capacity = $self->db->resultset('Deliveryservice')->search( { xml_id => $ds_name } )->get_column('global_max_mbps')->single() * 1000;
			}
			else {
				my $mkbps_match = $cdn_name . ":" . $ds_name . ":" . $cg_name . ":" . $host_name . ":maxKbps";
				$capacity = ( $redis->zrange( $mkbps_match, -1, -1 ) )[0];
			}

			$capacity = $capacity * 0.85;    # TODO JvD var is not always defined, hardcode
			my $e1 = tv_interval($startts);
			my @vals;
			my @times;
			my $jdata;

			$jdata->{cdnName}             = $cdn_name;
			$jdata->{deliveryServiceName} = $ds_name;
			$jdata->{cacheGroupName}      = $cg_name;
			$jdata->{hostName}            = $host_name;
			$jdata->{statName}            = $stat_name;

			my $first = int( ( $start - time ) / $REDIS_METRIC_INTERVAL );
			my $last  = int( ( $end - time ) / $REDIS_METRIC_INTERVAL );
			if ( $last >= 0 ) {
				$last = -1;
			}
			if ( $stat_name !~ /^daily_/ ) {
				if ( $last_one == 1 ) {
					$jdata->{redisCommand} = "   lrange " . $cdn_name . ":tstamp -1 -1; lrange " . $match . " -1 -1";
					@vals = $redis->lrange( $match, -1, -1 );
					@times = $redis->lrange( $cdn_name . ":tstamp", -1, -1 );
				}
				else {

					$jdata->{redisCommand} =
						"   lrange " . $cdn_name . ":tstamp " . $first . " " . $last . "; lrange " . $match . " " . $first . " " . $last;
					@times = $redis->lrange( $cdn_name . ":tstamp", $first, $last );
					@vals = $redis->lrange( $match, $first, $last );

					# print $#times . " -- " . $#vals . "\n";
				}
				my $prev_tstamp = 0;
				my $j           = 0;

				my $cc    = 1;
				my $total = 0;
				my $i     = 0;
				while ( $i <= $#vals ) {
					my $tstamp = $times[$i];
					my $val    = $vals[$i];

					$i++;

					# occasionally the delta will be +/- 1 from $REDIS_METRIC_INTERVAL
					my $delta = nearest( $REDIS_METRIC_INTERVAL, $tstamp - $prev_tstamp );

					if ( $delta != $REDIS_METRIC_INTERVAL ) {
						if (   exists( $jdata->{series} )
							&& scalar( @{ $jdata->{series} } ) > 0
							&& exists( $jdata->{series}->[$j]->{samples} )
							&& scalar( @{ $jdata->{series}->[$j]->{samples} } ) > 0 )
						{

							# move the index forward one to account for the gap
							$j++;
						}
						$jdata->{series}->[$j]->{timeBase} = int($tstamp);
						$cc                                = 1;
						$total                             = 0;
					}
					$total += $val;
					if ( $cc % ( $interval / $REDIS_METRIC_INTERVAL ) == 0 ) {
						push( @{ $jdata->{series}->[$j]->{samples} }, int( $total / ( $interval / $REDIS_METRIC_INTERVAL ) ) );
						$total = 0;
					}
					$cc++;
					$prev_tstamp = $tstamp;
				}

				# it's possible that the timeBase was inserted but no samples exist; delete the series if so
				if ( exists( $jdata->{series} ) && ref( $jdata->{series} ) eq "ARRAY" ) {
					my @deletes;

					for ( my $x = 0; $x < scalar( @{ $jdata->{series} } ); $x++ ) {
						if (  !exists( $jdata->{series}->[$x]->{samples} )
							|| ref( $jdata->{series}->[$x]->{samples} ) ne "ARRAY"
							|| scalar( @{ $jdata->{series}->[$x]->{samples} } ) == 0 )
						{
							push( @deletes, $x );
						}
					}

					for my $r (@deletes) {
						delete( $jdata->{series}->[$r] );
					}

					# if for some reason we now have no series, delete the key from the structure
					if ( scalar( @{ $jdata->{series} } ) == 0 ) {
						delete( $jdata->{series} );
					}
				}
			}
			else {
				my @daily_vals = $redis->lrange( $match, $first, $last );
				$jdata->{redisCommand} = "   lrange " . $match . " " . $first . " " . $last;
				my $j           = -1;
				my $prev_tstamp = 0;
				foreach my $line (@daily_vals) {
					my ( $tstamp, $val ) = split( /:/, $line );
					if ( $tstamp - $prev_tstamp != $interval ) {
						$j++;
						$jdata->{series}->[$j]->{timeBase} = int($tstamp);
					}
					push( @{ $jdata->{series}->[$j]->{samples} }, int($val) );
					$prev_tstamp = $tstamp;
					push( @times, $tstamp );    # TODO JvD
					push( @vals,  $val );
				}
			}

			my $e2 = tv_interval($startts);
			$jdata->{interval} = int($interval);
			$jdata->{start}    = gmtime( $times[0] );
			$jdata->{end}      = gmtime( $times[ $#vals - 1 ] );
			$jdata->{number}   = $#vals;
			$jdata->{capacity} = $capacity;

			$redis->quit();
			my $e3 = tv_interval($startts);
			$jdata->{elapsed} = $e3 . ' (' . $e1 . ',' . $e2 . ')';
			return $jdata;
		}
	);
	$app->renderer->add_helper(
		get_stats_long_term => sub {
			my $self     = shift;
			my $match    = shift;
			my $start    = shift;
			my $end      = shift;
			my $interval = shift;

			return undef;
		}
	);
	# end rm when redis is depracated
	
	$app->renderer->add_helper(
		get_cache_capacity => sub {
			my $self = shift;
			my $args = shift || {};
			$args->{type}   = "RASCAL";
			$args->{status} = "ONLINE";
			my $rascal_map  = $self->get_host_map($args);
			my $rascal_data = $self->get_rascal_state_data($args);

			my $raw_data = {
				capacity    => 0,
				count       => 0,
				available   => 0,
				unavailable => 0,
				maintenance => 0,
			};

			my $seen = {};

			for my $cdn_name ( keys( %{$rascal_map} ) ) {
				for my $rascal ( keys( %{ $rascal_map->{$cdn_name} } ) ) {
					if ( exists( $seen->{$cdn_name} ) ) {
						next;
					}
					else {
						$seen->{$cdn_name} = 1;
					}

					my $r = $self->get_traffic_monitor_connection( { cdn => $cdn_name } );
					my $stats = $r->get_cache_stats( { stats => "maxKbps,kbps" } );
					my $health_config = $self->get_health_config($cdn_name);

					for my $cache ( keys( %{ $stats->{caches} } ) ) {
						if (   !exists( $rascal_data->{$cdn_name}->{config}->{contentServers}->{$cache} )
							|| !exists( $rascal_data->{$cdn_name}->{state}->{$cache} ) )
						{
							next;
						}
						elsif ( $rascal_data->{$cdn_name}->{config}->{contentServers}->{$cache}->{type} ne "EDGE" ) {
							next;
						}

						my $key;
						my $c         = $rascal_data->{$cdn_name}->{config}->{contentServers}->{$cache};
						my $r         = $rascal_data->{$cdn_name}->{state}->{$cache};
						my $h         = $health_config->{profiles}->{ $c->{type} }->{ $c->{profile} };
						my $min_avail = $h->{"health.threshold.availableBandwidthInKbps"};
						$min_avail =~ s/\D//g;

						if (   ref($args) eq "HASH"
							&& exists( $args->{delivery_service} )
							&& !exists( $c->{deliveryServices}->{ $args->{delivery_service} } ) )
						{
							next;
						}

						if ( $c->{status} eq "REPORTED" || $c->{status} eq "ONLINE" ) {
							if ( $r->{isAvailable} ) {
								$key = "available";
							}
							else {
								$key = "unavailable";
							}
						}
						elsif ( $c->{status} eq "ADMIN_DOWN" ) {
							$key = "maintenance";
						}
						else {
							# skip OFFLINE or any other state
							next;
						}

						$raw_data->{count}++;
						$raw_data->{capacity} += ( $stats->{caches}->{$cache}->{maxKbps}->[0]->{value} - $min_avail );
						$raw_data->{$key} += $stats->{caches}->{$cache}->{kbps}->[0]->{value};
					}
				}
			}

			my $data = {
				utilizedPercent    => 0,
				unavailablePercent => 0,
				maintenancePercent => 0,
				availablePercent   => 0
			};

			if ( $raw_data->{capacity} > 0 ) {
				$data->{utilizedPercent}        = ( $raw_data->{available} / $raw_data->{capacity} ) * 100,
					$data->{unavailablePercent} = ( $raw_data->{unavailable} / $raw_data->{capacity} ) * 100,
					$data->{maintenancePercent} = ( $raw_data->{maintenance} / $raw_data->{capacity} ) * 100,
					$data->{availablePercent} =
					( ( $raw_data->{capacity} - $raw_data->{unavailable} - $raw_data->{maintenance} - $raw_data->{available} ) / $raw_data->{capacity} )
					* 100;
			}

			$self->success($data);
		}
	);
	$app->renderer->add_helper(
		get_routing_stats => sub {

			my $self = shift;
			my $args = shift;

			if ( !exists( $args->{status} ) ) {
				$args->{status} = "ONLINE";
			}

			$args->{type} = "CCR";

			my $ccr_map = $self->get_host_map($args);
			my $data    = {};
			my $stats   = {
				totalCount => 0,
				raw        => {},
			};

			for my $cdn_name ( keys( %{$ccr_map} ) ) {
				for my $ccr ( keys( %{ $ccr_map->{$cdn_name} } ) ) {
					my $ccr_host = $ccr_map->{$cdn_name}->{$ccr}->{host_name} . "." . $ccr_map->{$cdn_name}->{$ccr}->{domain_name};

					# TODO: what happens when the request to CCR times out? -jse
					my $c = $self->get_traffic_router_connection( { hostname => $ccr_host } );
					if ( !defined($c) ) {
						return "Cannot connect to Traffic Router";
					}

					my $s = $c->get_crs_stats();
					if ( !defined($s) ) {
						return ( "No CRS Stats found" );
					}

					if ( exists( $s->{stats} ) ) {
						for my $type ( "httpMap", "dnsMap" ) {
							next if ( exists( $args->{stat_key} ) && $args->{stat_key} ne $type );

							if ( exists( $s->{stats}->{$type} ) ) {
								for my $fqdn ( keys( %{ $s->{stats}->{$type} } ) ) {
									my $count = 1;
									if ( exists( $args->{patterns} ) && ref( $args->{patterns} ) eq "ARRAY" ) {
										$count = 0;

										for my $pattern ( @{ $args->{patterns} } ) {
											if ( $fqdn =~ /$pattern/ ) {
												$count = 1;
												last;
											}
										}
									}
									if ($count) {
										for my $counter ( keys( %{ $s->{stats}->{$type}->{$fqdn} } ) ) {
											if ( !exists( $stats->{raw}->{$counter} ) ) {
												$stats->{raw}->{$counter} = 0;
											}

											$stats->{raw}->{$counter} += $s->{stats}->{$type}->{$fqdn}->{$counter};
											$stats->{totalCount} += $s->{stats}->{$type}->{$fqdn}->{$counter};
										}
									}
								}
							}
						}
					}
				}
			}

			for my $counter ( keys( %{ $stats->{raw} } ) ) {
				my $p = $counter;
				$p =~ s/Count//gi;

				if ( $stats->{totalCount} > 0 ) {
					$data->{$p} = ( $stats->{raw}->{$counter} / $stats->{totalCount} ) * 100;
				}
				else {
					$data->{$p} = 0;
				}
			}

			$self->success($data);
			return (undef);
		}
	);
}

1;
