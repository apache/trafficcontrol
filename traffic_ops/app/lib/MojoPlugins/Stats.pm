package MojoPlugins::Stats;
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
						elsif ( $rascal_data->{$cdn_name}->{config}->{contentServers}->{$cache}->{type} !~ m/^EDGE/ ) {
							next;
						}

						my $key;
						my $c         = $rascal_data->{$cdn_name}->{config}->{contentServers}->{$cache};
						my $r         = $rascal_data->{$cdn_name}->{state}->{$cache};
						my $h         = $health_config->{profiles}->{ $c->{type} }->{ $c->{profile} };
						my $min_avail = $h->{"health.threshold.availableBandwidthInKbps"};
						if ($min_avail) {
							$min_avail =~ s/\D//g;
						}

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
