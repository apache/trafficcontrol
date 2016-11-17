package MojoPlugins::Health;
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
use JSON;
use Utils::Rascal;
use Data::Dumper;
use Data::Dump qw(dump);
use File::Slurp;
use Cwd;

#my $text = read_file($file);

sub register {
	my ( $self, $app, $conf ) = @_;
	$app->renderer->add_helper(
		get_health_config => sub {
			my $self = shift;
			my $cdn_name = shift || confess("Please supply a CDN name");
			my $rascal_profile;
			my @cache_profiles;
			my @ccr_profiles;
			my $ccr_profile_id;
			my $data_obj;
			my $profile_to_type;

			my $rs_pp = $self->db->resultset('Server')->search(
				{ 'cdn.name' => $cdn_name },
				{
					prefetch => [ 'cdn', 'profile', 'type' ],
					select   => 'me.profile',
					group_by => [qw/cdn.id profile.id me.profile type.id/],
				}
			);

			while ( my $row = $rs_pp->next ) {
			    my $profile_name = $row->profile->name;
				$self->app->log->debug("profile_name #-> " . Dumper($profile_name));
				if ( $row->profile->name =~ m/^RASCAL/ ) {
					$rascal_profile = $row->profile->name;
				}
				elsif ( $row->profile->name =~ m/^CCR/ ) {
					push( @ccr_profiles, $row->profile->name );

					# TODO MAT: support multiple CCR profiles
					$ccr_profile_id = $row->profile->id;
				}
				elsif ( $row->type->name =~ m/^EDGE/ || $row->type->name =~ m/^MID/ ) {
					push( @cache_profiles, $row->profile->name );
					$profile_to_type->{ $row->profile->name }->{ $row->type->name } = $row->type->name;
				}
			}
			my %condition = ( 'parameter.config_file' => 'rascal-config.txt', 'profile.name' => $rascal_profile );
			$rs_pp = $self->db->resultset('ProfileParameter')->search( \%condition, { prefetch => [ { 'parameter' => undef }, { 'profile' => undef } ] } );
			while ( my $row = $rs_pp->next ) {
				if ( $row->parameter->name =~ m/location/ ) { next; }
				$data_obj->{'rascal-config'}->{ $row->parameter->name } = $row->parameter->value;
			}

			%condition = ( 'parameter.config_file' => 'rascal.properties', 'profile.name' => { -in => \@cache_profiles } );
			$rs_pp = $self->db->resultset('ProfileParameter')->search( \%condition, { prefetch => [ { 'parameter' => undef }, { 'profile' => undef } ] } );
			while ( my $row = $rs_pp->next ) {
				if ( exists( $profile_to_type->{ $row->profile->name } ) ) {
					for my $profile_type ( keys( %{ $profile_to_type->{ $row->profile->name } } ) ) {
						$data_obj->{'profiles'}->{$profile_type}->{ $row->profile->name }->{ $row->parameter->name } = $row->parameter->value;
					}
				}
			}
			foreach my $ccr_profile (@ccr_profiles) {
				$data_obj->{'profiles'}->{'CCR'}->{$ccr_profile} = undef;
			}
			my $rs_ds = $self->db->resultset('Deliveryservice')->search( { 'me.profile' => $ccr_profile_id, 'active' => 1 }, {} );
			while ( my $row = $rs_ds->next ) {

				# MAT: Do we move this to the DB? Rascal needs to know if it should monitor a DS or not, and the status=REPORTED is what we do for caches.
				$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'status'} = "REPORTED";

				# MAT: The > 0 is a hack because MySQL creates phantom 0s on insert
				if ( defined( $row->global_max_mbps ) && $row->global_max_mbps > 0 ) {
					$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'health.threshold.total.kbps'} = $row->global_max_mbps * 1000;
				}
				if ( defined( $row->global_max_tps ) && $row->global_max_tps > 0 ) {
					$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'health.threshold.total.tps_total'} = int( $row->global_max_tps );
				}
			}

			return ($data_obj);

		}
	);

	$app->renderer->add_helper(
		get_host_map => sub {
			my $self = shift;
			my $args = shift;

			if ( !defined($args) || ref($args) ne "HASH" ) {
				confess("Supply a hashref of arguments");
			}
			elsif ( !exists( $args->{type} ) ) {
				confess("Supply a type in the argument hashref");
			}

			my $host_map = {};

			my $rs_type = $self->db->resultset('Type')->search( { -or => [ name => $args->{type} ] } );
			my $rs_data =
				$self->db->resultset('Server')
				->search( { type => { -in => $rs_type->get_column('id')->as_query } }, { prefetch => [ 'cdn', { 'status' => undef } ] } );

			while ( my $row = $rs_data->next ) {
				my $this_cdn_name = $row->cdn->name;

				if ( !defined($this_cdn_name) ) {
					print "cdn name not defined\n";
				}

				next if ( exists( $args->{cdn_name} ) && $args->{cdn_name} ne $this_cdn_name );
				next if ( exists( $args->{status} )   && $args->{status} ne $row->status->name );

				$host_map->{$this_cdn_name}->{ $row->host_name }->{host_name}      = $row->host_name;
				$host_map->{$this_cdn_name}->{ $row->host_name }->{domain_name}    = $row->domain_name;
				$host_map->{$this_cdn_name}->{ $row->host_name }->{ip_address}     = $row->ip_address;
				$host_map->{$this_cdn_name}->{ $row->host_name }->{interface_name} = $row->interface_name;
				$host_map->{$this_cdn_name}->{ $row->host_name }->{tcp_port}       = $row->tcp_port;
				$host_map->{$this_cdn_name}->{ $row->host_name }->{cachegroup}     = $row->cachegroup->name;
				$host_map->{$this_cdn_name}->{ $row->host_name }->{status}         = $row->status->name;
				$host_map->{$this_cdn_name}->{ $row->host_name }->{profile}        = $row->profile->name;
			}

			return ($host_map);
		}
	);
	$app->renderer->add_helper(
		get_cache_health => sub {
			my $self = shift;
			my $args = shift;

			my $rascal_data = $self->get_rascal_state_data($args);
			my $data        = {
				totalOnline  => 0,
				totalOffline => 0,
				cachegroups  => [],
			};

			my $cachegroup_data = {};
			my $capacity_data = { total => 0 };

			for my $cdn_name ( keys( %{$rascal_data} ) ) {
				for my $edge ( keys( %{ $rascal_data->{$cdn_name}->{state} } ) ) {

					if ( exists( $rascal_data->{$cdn_name}->{config}->{contentServers}->{$edge} ) ) {
						my $cache_config = $rascal_data->{$cdn_name}->{config}->{contentServers}->{$edge};
						my $cache_state  = $rascal_data->{$cdn_name}->{state}->{$edge};

						#NOTE: The 'locationId' will need to change when Rascal gets updated to use 'cachegroups' properly.
						my $locationId = $cache_config->{locationId};

						my $status = $cache_config->{status};

						if ( $cache_config->{type} !~ m/^EDGE/ ) {
							next;
						}

						if ( !exists( $capacity_data->{$status} ) ) {
							$capacity_data->{$status} = 0;
						}

						$capacity_data->{total}++;
						$capacity_data->{status}++;

						if ( $status ne "REPORTED" ) {
							next;
						}

						my $count = 0;

						if ( exists( $args->{delivery_service} ) ) {
							if ( exists( $cache_config->{deliveryServices} ) && exists( $cache_config->{deliveryServices}->{ $args->{delivery_service} } ) )
							{
								$count = 1;
							}
						}
						else {
							$count = 1;
						}

						if ($count) {
							if ( !exists( $cachegroup_data->{$locationId} ) ) {
								$cachegroup_data->{$locationId}->{online}  = 0;
								$cachegroup_data->{$locationId}->{offline} = 0;
							}

							if ( $cache_state->{isAvailable} ) {
								$data->{totalOnline}++;
								$cachegroup_data->{$locationId}->{online}++;
							}
							else {
								$data->{totalOffline}++;
								$cachegroup_data->{$locationId}->{offline}++;
							}
						}
					}
				}
			}

			for my $cachegroup ( sort { $cachegroup_data->{$b}->{offline} <=> $cachegroup_data->{$a}->{offline} } keys( %{$cachegroup_data} ) ) {

				# only add cachegroups that have at least one online/offline server
				if ( $cachegroup_data->{$cachegroup}->{online} || $cachegroup_data->{$cachegroup}->{offline} ) {
					push(
						@{ $data->{cachegroups} }, {
							name    => $cachegroup,
							online  => $cachegroup_data->{$cachegroup}->{online},
							offline => $cachegroup_data->{$cachegroup}->{offline},
						}
					);
				}
			}

			$self->success($data);
		}
	);

	$app->renderer->add_helper(
		get_rascal_state_data => sub {
			my $self = shift;

			my $args = shift;
			$args->{status} = "ONLINE";
			$args->{type}   = "RASCAL";
			my $what        = ( defined($args) && ref($args) eq "HASH" && exists( $args->{state_type} ) ) ? $args->{state_type} : "caches";
			my $rascal_map  = $self->get_host_map($args);
			my $rascal_data = {};

			for my $cdn_name ( keys( %{$rascal_map} ) ) {
				for my $rascal ( keys( %{ $rascal_map->{$cdn_name} } ) ) {
					if ( exists( $rascal_data->{$cdn_name} ) ) {
						next;
					}
					my $r      = $self->get_traffic_monitor_connection( { cdn => $cdn_name } );
					my $state  = $r->get_states($what);
					my $config = $r->get_cr_config();
					if ( defined($state) ) {
						$rascal_data->{$cdn_name}->{state}  = $state;
						$rascal_data->{$cdn_name}->{config} = $config;
					}
				}
			}
			return ($rascal_data);
		}
	);

}

1;
