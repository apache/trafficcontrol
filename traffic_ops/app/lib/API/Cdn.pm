
package API::Cdn;
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
use Carp qw(cluck confess);
use JSON;
use MIME::Base64;
use UI::DeliveryService;
use MojoPlugins::Response;
use Extensions::v11::Delegate::Metrics;
use Extensions::v11::Delegate::Statistics;
use Common::ReturnCodes qw(SUCCESS ERROR);

################################################################################
# WARNING: This route is unauthenticated!
# Note: we only have a summary route thus far.
################################################################################
sub metrics {
	my $self = shift;
	my $m    = new Extensions::v11::Delegate::Metrics($self);
	my ( $rc, $result ) = $m->get_etl_metrics();
	if ( $rc == SUCCESS ) {
		return ( $self->success($result) );
	}
	else {
		return ( $self->alert($result) );
	}
}

sub usage_overview {
	my $self = shift;

	my $st = new Extensions::v11::Delegate::Statistics($self);
	my ( $rc, $result ) = $st->get_usage_overview();
	if ( $rc == SUCCESS ) {
		$self->success($result);
	}
	else {
		$self->alert($result);
	}
}

# Used by the 'Daily Summary' menu item
sub peakusage {
	my $self = shift;

	my $stats = new Extensions::v11::Delegate::Statistics($self);
	my ( $rc, $result ) = $stats->get_daily_summary();

	if ( $rc == SUCCESS ) {
		return $self->success($result);
	}
	else {
		return $self->alert($result);
	}
}

sub configs_monitoring {
	my $self      = shift;
	my $cdn_name  = $self->param('name');
	my $extension = $self->param('extension');

	my $data_obj = $self->get_traffic_monitor_config($cdn_name);
	$self->success($data_obj);
}

sub get_traffic_monitor_config {
	my $self = shift;
	my $cdn_name = shift || confess("Please supply a CDN name");
	my $rascal_profile;
	my @cache_profiles;
	my @ccr_profiles;
	my $ccr_profile_id;
	my $data_obj;

	my %condition = ( 'parameter.name' => 'CDN_name', 'parameter.value' => $cdn_name );
	my $rs_pp = $self->db->resultset('ProfileParameter')->search( \%condition, { prefetch => [ { 'parameter' => undef }, { 'profile' => undef } ] } );
	while ( my $row = $rs_pp->next ) {
		if ( $row->profile->name =~ m/^RASCAL/ ) {
			$rascal_profile = $row->profile->name;
		}
		elsif ( $row->profile->name =~ m/^CCR/ ) {
			push( @ccr_profiles, $row->profile->name );

			# TODO MAT: support multiple CCR profiles
			$ccr_profile_id = $row->profile->id;
		}
		elsif ( $row->profile->name =~ m/^EDGE/ || $row->profile->name =~ m/^MID/ ) {
			push( @cache_profiles, $row->profile->name );
		}
	}
	%condition = ( 'parameter.config_file' => 'rascal-config.txt', 'profile.name' => $rascal_profile );
	$rs_pp = $self->db->resultset('ProfileParameter')->search( \%condition, { prefetch => [ { 'parameter' => undef }, { 'profile' => undef } ] } );

	while ( my $row = $rs_pp->next ) {
		my $parameter;
		if ( $row->parameter->name =~ m/location/ ) { next; }
		if ( $row->parameter->value =~ m/^\d+$/ ) {
			$data_obj->{'config'}->{ $row->parameter->name } = int( $row->parameter->value );
		}
		else {
			$data_obj->{'config'}->{ $row->parameter->name } = $row->parameter->value;
		}
	}

	%condition = ( 'parameter.config_file' => 'rascal.properties', 'profile.name' => { -in => \@cache_profiles } );
	$rs_pp = $self->db->resultset('ProfileParameter')->search( \%condition, { prefetch => [ { 'parameter' => undef }, { 'profile' => undef } ] } );

	if ( !exists( $data_obj->{'profiles'} ) ) {
		$data_obj->{'profiles'} = undef;
	}
	my $profile_tracker;

	while ( my $row = $rs_pp->next ) {

		my $type;
		if ( $row->profile->name =~ m/^EDGE/ ) {
			$type = "EDGE";
		}
		elsif ( $row->profile->name =~ m/MID/ ) {
			$type = "MID";
		}
		$profile_tracker->{ $row->profile->name }->{'type'} = $type;
		$profile_tracker->{ $row->profile->name }->{'name'} = $row->profile->name;

		if ( $row->parameter->value =~ m/^\d+$/ ) {
			$profile_tracker->{ $row->profile->name }->{'parameters'}->{ $row->parameter->name } = int( $row->parameter->value );
		}
		else {
			$profile_tracker->{ $row->profile->name }->{'parameters'}->{ $row->parameter->name } = $row->parameter->value;
		}
	}

	foreach my $profile ( keys %{$profile_tracker} ) {
		push( @{ $data_obj->{'profiles'} }, $profile_tracker->{$profile} );
	}

	foreach my $ccr_profile (@ccr_profiles) {
		my $profile;
		$profile->{'name'}       = $ccr_profile;
		$profile->{'type'}       = "CCR";
		$profile->{'parameters'} = undef;
		push( @{ $data_obj->{'profiles'} }, $profile );
	}

	my $rs_ds = $self->db->resultset('Deliveryservice')->search( { 'me.profile' => $ccr_profile_id, 'active' => 1 }, {} );
	while ( my $row = $rs_ds->next ) {
		my $delivery_service;

		# MAT: Do we move this to the DB? Rascal needs to know if it should monitor a DS or not, and the status=REPORTED is what we do for caches.
		$delivery_service->{'xmlId'}              = $row->xml_id;
		$delivery_service->{'status'}             = "REPORTED";
		$delivery_service->{'totalKbpsThreshold'} = $row->global_max_mbps * 1000;
		$delivery_service->{'totalTpsThreshold'}  = int( $row->global_max_tps || 0 );
		push( @{ $data_obj->{'deliveryServices'} }, $delivery_service );
	}

	my @cdn_profiles;
	my $cdnname_param_id = $self->db->resultset('Parameter')->search( { name => 'CDN_name', value => $cdn_name } )->get_column('id')->single();
	if ( defined($cdnname_param_id) ) {
		@cdn_profiles = $self->db->resultset('ProfileParameter')->search( { parameter => $cdnname_param_id } )->get_column('profile')->all();
		if ( !scalar(@cdn_profiles) ) {
			my $e = Mojo::Exception->throw( "No profiles found for CDN_name: " . $cdn_name );
		}
	}
	else {
		my $e = Mojo::Exception->throw( "Parameter ID not found for CDN_name: " . $cdn_name );
	}

	my $rs_caches = $self->db->resultset('Server')->search(
		{ 'profile' => { -in => \@cdn_profiles } },
		{
			prefetch => [ 'type',      'status',      'cachegroup', 'profile' ],
			columns  => [ 'host_name', 'domain_name', 'tcp_port',   'interface_name', 'ip_address', 'ip6_address', 'id', 'xmpp_id' ]
		}
	);
	while ( my $row = $rs_caches->next ) {
		if ( $row->type->name eq "RASCAL" ) {
			my $traffic_monitor;
			$traffic_monitor->{'hostName'}   = $row->host_name;
			$traffic_monitor->{'fqdn'}       = $row->host_name . "." . $row->domain_name;
			$traffic_monitor->{'status'}     = $row->status->name;
			$traffic_monitor->{'cachegroup'} = $row->cachegroup->name;
			$traffic_monitor->{'port'}       = int( $row->tcp_port );
			$traffic_monitor->{'ip'}         = $row->ip_address;
			$traffic_monitor->{'ip6'}        = $row->ip6_address;
			$traffic_monitor->{'profile'}    = $row->profile->name;
			push( @{ $data_obj->{'trafficMonitors'} }, $traffic_monitor );

		}
		elsif ( $row->type->name eq "EDGE" || $row->type->name eq "MID" ) {
			my $traffic_server;
			$traffic_server->{'cachegroup'}    = $row->cachegroup->name;
			$traffic_server->{'hostName'}      = $row->host_name;
			$traffic_server->{'fqdn'}          = $row->host_name . "." . $row->domain_name;
			$traffic_server->{'port'}          = int( $row->tcp_port );
			$traffic_server->{'interfaceName'} = $row->interface_name;
			$traffic_server->{'status'}        = $row->status->name;
			$traffic_server->{'ip'}            = $row->ip_address;
			$traffic_server->{'ip6'}           = ( $row->ip6_address || "" );
			$traffic_server->{'profile'}       = $row->profile->name;
			$traffic_server->{'type'}          = $row->type->name;
			$traffic_server->{'hashId'}        = $row->xmpp_id;
			push( @{ $data_obj->{'trafficServers'} }, $traffic_server );
		}

	}

	my $rs_loc = $self->db->resultset('CachegroupParameter')->search( { 'parameter' => $cdnname_param_id }, { prefetch => 'cachegroup' } );
	while ( my $row = $rs_loc->next ) {
		my $cache_group;
		my $latitude  = $row->cachegroup->latitude + 0;
		my $longitude = $row->cachegroup->longitude + 0;
		$cache_group->{'coordinates'}->{'latitude'}  = $latitude;
		$cache_group->{'coordinates'}->{'longitude'} = $longitude;
		$cache_group->{'name'}                       = $row->cachegroup->name;
		push( @{ $data_obj->{'cacheGroups'} }, $cache_group );
	}

	return ($data_obj);
}

sub capacity {
	my $self = shift;

	return $self->get_cache_capacity();
}

sub health {
	my $self = shift;

	return $self->get_cache_health();
}

sub routing {
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
			my $s = $c->get_crs_stats();
			if ( !defined($s) ) {
				return $self->internal_server_error( { "Internal Server" => "Error" } );
			}
			else {

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
}

sub configs_routing {
	my $self     = shift;
	my $cdn_name = $self->param('name');
	my $data_obj;
	my $json = $self->gen_traffic_router_config($cdn_name);
	$self->success($json);
}

sub gen_traffic_router_config {
	my $self     = shift;
	my $cdn_name = shift;
	my $data_obj;
	my @cdn_profiles;
	my $ccr_profile_id;
	my $ccr_domain_name = "";
	$SIG{__WARN__} = sub { warn $_[0] unless $_[0] =~ m/Prefetching multiple has_many rels deliveryservice_servers/ };

	$data_obj->{'stats'}->{'cdnName'}           = $cdn_name;
	$data_obj->{'stats'}->{'date'}              = time();
	$data_obj->{'stats'}->{'trafficOpsVersion'} = &tm_version();
	$data_obj->{'stats'}->{'trafficOpsPath'}    = $self->req->url->path->{'path'};
	$data_obj->{'stats'}->{'trafficOpsHost'}    = $self->req->headers->host;
	$data_obj->{'stats'}->{'trafficOpsUser'}    = $self->current_user()->{username};

	my $cdnname_param_id = $self->db->resultset('Parameter')->search( { name => 'CDN_name', value => $cdn_name } )->get_column('id')->single();
	if ( defined($cdnname_param_id) ) {
		@cdn_profiles = $self->db->resultset('ProfileParameter')->search( { parameter => $cdnname_param_id } )->get_column('profile')->all();
		if ( scalar(@cdn_profiles) ) {
			$ccr_profile_id =
				$self->db->resultset('Profile')->search( { id => { -in => \@cdn_profiles }, name => { -like => 'CCR%' } } )->get_column('id')->single();
			if ( !defined($ccr_profile_id) ) {
				my $e = Mojo::Exception->throw("No CCR profile found in profile IDs: @cdn_profiles ");
			}
		}
		else {
			my $e = Mojo::Exception->throw( "No profiles found for CDN_name: " . $cdn_name );
		}

#@cache_rascal_profiles = $self->db->resultset('Profile')->search( { id => { -in => \@cdn_profiles }, name => [{ like => 'EDGE%'}, {like => 'MID%'}, {like => 'RASCAL%'}, {like => 'CDSIS%'} ] } )->get_column('id')->all();;
	}
	else {
		my $e = Mojo::Exception->throw( "Parameter ID not found for CDN_name: " . $cdn_name );
	}

	my %condition = ( 'profile_parameters.profile' => $ccr_profile_id, 'config_file' => 'CRConfig.json' );
	my $rs_config = $self->db->resultset('Parameter')->search( \%condition, { join => 'profile_parameters' } );
	while ( my $row = $rs_config->next ) {
		my $parameter;
		if ( $row->name eq 'domain_name' ) {
			$ccr_domain_name = $row->value;
		}

		$parameter->{'type'} = "parameter";
		if ( $row->value =~ m/^\d+$/ ) {
			$data_obj->{'config'}->{ $row->name } = int( $row->value );
		}
		else {
			$data_obj->{'config'}->{ $row->name } = $row->value;
		}

		#push (@{$data_obj->{'config'}}, $parameter);

	}
	my $rs_loc = $self->db->resultset('CachegroupParameter')->search( { 'parameter' => $cdnname_param_id }, { prefetch => 'cachegroup' } );

#my $rs_loc = $self->db->resultset('Location')->search( {'servers.profile' => { -in => \@cdn_profiles }, 'type.name' => { -like => 'EDGE%'} }, { join => ['servers', 'type'], group_by => 'short_name' } );
	while ( my $row = $rs_loc->next ) {
		my $cache_group;
		my $latitude  = $row->cachegroup->latitude + 0;
		my $longitude = $row->cachegroup->longitude + 0;
		$cache_group->{'coordinates'}->{'latitude'}  = $latitude;
		$cache_group->{'coordinates'}->{'longitude'} = $longitude;
		$cache_group->{'name'}                       = $row->cachegroup->name;
		push( @{ $data_obj->{'cacheGroups'} }, $cache_group );
	}
	my $regex_tracker;
	my $rs_regexes = $self->db->resultset('Regex')->search( {}, { 'prefetch' => 'type' } );
	while ( my $row = $rs_regexes->next ) {
		$regex_tracker->{ $row->id }->{'type'}    = $row->type->name;
		$regex_tracker->{ $row->id }->{'pattern'} = $row->pattern;
	}
	my %cache_tracker;
	my $rs_caches = $self->db->resultset('Server')->search(
		{ 'profile' => { -in => \@cdn_profiles } },
		{
			prefetch => [ 'type',      'status',      'cachegroup', 'profile' ],
			columns  => [ 'host_name', 'domain_name', 'tcp_port',   'interface_name', 'ip_address', 'ip6_address', 'id', 'xmpp_id' ]
		}
	);
	while ( my $row = $rs_caches->next ) {
		if ( $row->type->name eq "RASCAL" ) {
			my $traffic_monitor;
			$traffic_monitor->{'hostName'} = $row->host_name;
			$traffic_monitor->{'fqdn'}     = $row->host_name . "." . $row->domain_name;
			$traffic_monitor->{'status'}   = $row->status->name;
			$traffic_monitor->{'location'} = $row->cachegroup->name;
			$traffic_monitor->{'port'}     = int( $row->tcp_port );
			$traffic_monitor->{'ip'}       = $row->ip_address;
			$traffic_monitor->{'ip6'}      = $row->ip6_address;
			$traffic_monitor->{'profile'}  = $row->profile->name;
			push( @{ $data_obj->{'trafficMonitors'} }, $traffic_monitor );

		}
		elsif ( $row->type->name eq "CCR" ) {
			my $rs_param = $self->db->resultset('Parameter')
				->search( { 'profile_parameters.profile' => $row->profile->id, 'name' => 'api.port' }, { join => 'profile_parameters' } );
			my $r = $rs_param->single;
			my $api_port = ( defined($r) && defined( $r->value ) ) ? $r->value : 3333;

			my $traffic_router;

			$traffic_router->{'hostName'} = $row->host_name;
			$traffic_router->{'fqdn'}     = $row->host_name . "." . $row->domain_name;
			$traffic_router->{'status'}   = $row->status->name;
			$traffic_router->{'location'} = $row->cachegroup->name;
			$traffic_router->{'port'}     = int( $row->tcp_port );
			$traffic_router->{'apiPort'}  = int($api_port);
			$traffic_router->{'ip'}       = $row->ip_address;
			$traffic_router->{'ip6'}      = $row->ip6_address;
			$traffic_router->{'profile'}  = $row->profile->name;
			push( @{ $data_obj->{'trafficRouters'} }, $traffic_router );
		}
		elsif ( $row->type->name eq "EDGE" || $row->type->name eq "MID" ) {
			if ( !exists $cache_tracker{ $row->id } ) {
				$cache_tracker{ $row->id } = $row->host_name;
			}

			my $traffic_server;
			$traffic_server->{'cacheGroup'}    = $row->cachegroup->name;
			$traffic_server->{'hostName'}      = $row->host_name;
			$traffic_server->{'fqdn'}          = $row->host_name . "." . $row->domain_name;
			$traffic_server->{'port'}          = int( $row->tcp_port );
			$traffic_server->{'interfaceName'} = $row->interface_name;
			$traffic_server->{'status'}        = $row->status->name;
			$traffic_server->{'ip'}            = $row->ip_address;
			$traffic_server->{'ip6'}           = ( $row->ip6_address || "" );
			$traffic_server->{'profile'}       = $row->profile->name;
			$traffic_server->{'type'}          = $row->type->name;
			$traffic_server->{'hashId'}        = $row->xmpp_id;
			push( @{ $data_obj->{'trafficServers'} }, $traffic_server );
		}

	}

	my $ds_regex_tracker;
	my $regexps;
	my $rs_ds =
		$self->db->resultset('Deliveryservice')
		->search( { 'me.profile' => $ccr_profile_id, 'active' => 1 }, { prefetch => [ 'deliveryservice_servers', 'deliveryservice_regexes', 'type' ] } );
	while ( my $row = $rs_ds->next ) {
		my $delivery_service;
		$delivery_service->{'xmlId'} = $row->xml_id;
		my $protocol;
		if ( $row->type->name =~ m/DNS/ ) {
			$protocol = 'DNS';
		}
		else {
			$protocol = 'HTTP';
		}
		my @server_subrows = $row->deliveryservice_servers->all;
		my @regex_subrows  = $row->deliveryservice_regexes->all;
		my $regex_to_props;
		my %ds_to_remap;
		if ( scalar(@regex_subrows) ) {
			foreach my $subrow (@regex_subrows) {
				$delivery_service->{'matchSets'}->[ $subrow->set_number ]->{'protocol'} = $protocol;
				$regex_to_props->{ $subrow->{'_column_data'}->{'regex'} }->{'pattern'} =
					$regex_tracker->{ $subrow->{'_column_data'}->{'regex'} }->{'pattern'};
				$regex_to_props->{ $subrow->{'_column_data'}->{'regex'} }->{'setNumber'} = $subrow->set_number;
				$regex_to_props->{ $subrow->{'_column_data'}->{'regex'} }->{'type'} =
					$regex_tracker->{ $subrow->{'_column_data'}->{'regex'} }->{'type'};
				if ( $regex_to_props->{ $subrow->{'_column_data'}->{'regex'} }->{'type'} eq 'HOST_REGEXP' ) {
					$ds_to_remap{ $row->xml_id }->[ $subrow->set_number ] = $regex_to_props->{ $subrow->{'_column_data'}->{'regex'} }->{'pattern'};
				}
			}
		}
		foreach my $regex ( sort keys %{$regex_to_props} ) {
			my $set_number = $regex_to_props->{$regex}->{'setNumber'};
			my $pattern    = $regex_to_props->{$regex}->{'pattern'};
			my $type       = $regex_to_props->{$regex}->{'type'};
			if ( $type eq 'HOST_REGEXP' ) {
				push( @{ $delivery_service->{'matchSets'}->[$set_number]->{'matchList'} }, { 'matchType' => 'HOST', 'regex' => $pattern } );
			}
			elsif ( $type eq 'PATH_REGEXP' ) {
				push( @{ $delivery_service->{'matchSets'}->[$set_number]->{'matchList'} }, { 'matchType' => 'PATH', 'regex' => $pattern } );
			}
			elsif ( $type eq 'HEADER_REGEXP' ) {
				push( @{ $delivery_service->{'matchSets'}->[$set_number]->{'matchList'} }, { 'matchType' => 'HEADER', 'regex' => $pattern } );
			}
		}
		if ( scalar(@server_subrows) ) {

			#my $host_regex = qr/(^(\.)+\*\\\.)(.*)(\\\.(\.)+\*$)/;
			my $host_regex1 = qr/\\|\.\*/;

			#MAT: Have to do this dedup because @server_subrows contains duplicates (* the # of host regexes)
			my %server_subrow_dedup;
			foreach my $subrow (@server_subrows) {
				$server_subrow_dedup{ $subrow->{'_column_data'}->{'server'} } = $subrow->{'_column_data'}->{'deliveryservice'};
			}
			my $ds_regex->{'xmlId'} = $row->xml_id;
			foreach my $server ( keys %server_subrow_dedup ) {
				my @remaps;
				foreach my $host ( @{ $ds_to_remap{ $row->xml_id } } ) {
					my $remap;
					if ( $host =~ m/\.\*$/ ) {
						my $host_copy = $host;
						$host_copy =~ s/$host_regex1//g;
						if ( $protocol eq 'DNS' ) {
							$remap = 'edge' . $host_copy . $ccr_domain_name;
						}
						else {
							my $cache_tracker_server = $cache_tracker{$server} || "";
							my $host_copy            = $host_copy              || "";
							my $ccr_domain_name      = $ccr_domain_name        || "";
							$remap = $cache_tracker_server . $host_copy . $ccr_domain_name;
						}
					}
					else {
						$remap = $host;
					}
					push( @remaps, $remap );
				}
				my $cache_tracker_server = $cache_tracker{$server} || "";
				push( @{ $ds_regex_tracker->{$cache_tracker_server}->{ $row->xml_id }->{'remaps'} }, @remaps );
			}
		}

		$delivery_service->{'ttl'} = int( $row->ccr_dns_ttl );
		my $geo_limit = $row->geo_limit;
		if ( $geo_limit == 1 ) {

			# Ref to 0 or 1 makes JSON bool value
			$delivery_service->{'coverageZoneOnly'} = \1;
			$delivery_service->{'geoEnabled'}       = [];
		}
		elsif ( $geo_limit == 2 ) {

			# Ref to 0 or 1 makes JSON bool value
			$delivery_service->{'coverageZoneOnly'} = \0;
			$delivery_service->{'geoEnabled'} = [ { 'countryCode' => 'US' } ];
		}
		else {
			# Ref to 0 or 1 makes JSON bool value
			$delivery_service->{'coverageZoneOnly'} = \0;
			$delivery_service->{'geoEnabled'}       = [];
		}
		my $bypass_destination;
		if ( $protocol =~ m/DNS/ ) {
			$bypass_destination->{'type'} = 'DNS';
			if ( defined( $row->dns_bypass_ip ) && $row->dns_bypass_ip ne "" ) {
				$bypass_destination->{'ip'} = $row->dns_bypass_ip;
			}
			if ( defined( $row->dns_bypass_ip6 ) && $row->dns_bypass_ip6 ne "" ) {
				$bypass_destination->{'ip6'} = $row->dns_bypass_ip6;
			}
			if ( defined( $row->dns_bypass_ttl ) && $row->dns_bypass_ttl ne "" ) {
				$bypass_destination->{'ttl'} = int( $row->dns_bypass_ttl );
			}
			if ( defined( $row->max_dns_answers ) && $row->max_dns_answers ne "" ) {
				$bypass_destination->{'maxDnsIpsForLocation'} = int( $row->max_dns_answers );
			}
		}
		elsif ( $protocol =~ m/HTTP/ ) {
			$bypass_destination->{'type'} = 'HTTP';
			if ( defined( $row->http_bypass_fqdn ) && $row->http_bypass_fqdn ne "" ) {
				my $full = $row->http_bypass_fqdn;
				my $port;
				my $fqdn;
				if ( $full =~ m/\:/ ) {
					( $fqdn, $port ) = split( /\:/, $full );
				}
				else {
					$fqdn = $full;
					$port = 80;
				}
				$bypass_destination->{'fqdn'} = $fqdn;
				$bypass_destination->{'port'} = int($port);
			}
		}
		$delivery_service->{'bypassDestination'} = $bypass_destination;

		if ( defined( $row->miss_lat ) && $row->miss_lat ne "" ) {
			$delivery_service->{'missCoordinates'}->{'latitude'} = $row->miss_lat + 0;
		}
		if ( defined( $row->miss_long ) && $row->miss_long ne "" ) {
			$delivery_service->{'missCoordinates'}->{'longitude'} = $row->miss_long + 0;
		}
		$delivery_service->{'ttls'} = { 'A' => int( $row->ccr_dns_ttl ), 'AAAA' => int( $row->ccr_dns_ttl ), 'NS' => 3600, 'SOA' => 86400 };
		$delivery_service->{'soa'}->{'minimum'} = 30;
		$delivery_service->{'soa'}->{'expire'}  = 604800;
		$delivery_service->{'soa'}->{'retry'}   = 7200;
		$delivery_service->{'soa'}->{'refresh'} = 28800;
		$delivery_service->{'soa'}->{'admin'}   = "twelve_monkeys";

		my $rs_dns = $self->db->resultset('Staticdnsentry')->search(
			{ 'deliveryservice.active' => 1, 'deliveryservice.profile' => $ccr_profile_id },
			{ prefetch => [ 'deliveryservice', 'type' ], columns => [ 'host', 'type', 'ttl', 'address' ] }
		);
		while ( my $dns_row = $rs_dns->next ) {
			my $dns_obj;
			$dns_obj->{'name'}  = $dns_row->host;
			$dns_obj->{'ttl'}   = int( $dns_row->ttl );
			$dns_obj->{'value'} = $dns_row->address;

			my $type = $dns_row->type->name;
			$type =~ s/\_RECORD//g;
			$dns_obj->{'type'} = $type;
			push( @{ $delivery_service->{'staticDnsEntries'} }, $dns_obj );
		}

		push( @{ $data_obj->{'deliveryServices'} }, $delivery_service );
	}

	foreach my $cache_hostname ( sort keys %{$ds_regex_tracker} ) {
		my $i = 0;
		my $server_ref;
		foreach my $traffic_server ( @{ $data_obj->{'trafficServers'} } ) {
			$i++;
			my $traffic_server_hostname = $traffic_server->{'hostName'} || "";
			next if ( $traffic_server_hostname ne $cache_hostname );
			$server_ref = $data_obj->{'trafficServers'}->[ $i - 1 ];
		}

		foreach my $xml_id ( sort keys %{ $ds_regex_tracker->{$cache_hostname} } ) {
			my $ds;
			$ds->{'xmlId'}  = $xml_id;
			$ds->{'remaps'} = $ds_regex_tracker->{$cache_hostname}->{$xml_id}->{'remaps'};
			push( @{ $server_ref->{'deliveryServices'} }, $ds );
			$data_obj->{'trafficServers'}->[$i] = $server_ref;
		}
	}

	my @empty_array;
	foreach my $traffic_server ( @{ $data_obj->{'trafficServers'} } ) {
		if ( !defined( $traffic_server->{'deliveryServices'} ) ) {
			push( @{ $traffic_server->{'deliveryServices'} }, @empty_array );
		}
	}

	return ($data_obj);
}

sub get_cdn_name {
	my $self  = shift;
	my $which = shift;

	my $cdn_name = "all";
	my $server =
		$self->db->resultset('Server')->search( { host_name => $which }, { prefetch => [ 'cachegroup', 'type', 'profile', 'status', 'phys_location' ] } )
		->single();
	if ( defined($server) ) {
		my $param =
			$self->db->resultset('ProfileParameter')
			->search( { -and => [ profile => $server->profile->id, 'parameter.config_file' => 'rascal-config.txt', 'parameter.name' => 'CDN_name' ] },
			{ prefetch => [ { parameter => undef }, { profile => undef } ] } )->single();
		$cdn_name = $param->parameter->value;
	}
	return $cdn_name;
}

# Produces a list of Cdns for traversing child links
sub get_cdns {
	my $self = shift;

	my $rs_data = $self->db->resultset("Parameter")->search( { name => 'CDN_name' }, { order_by => "name" } );
	my $json_response = $self->build_cdns_json( $rs_data, "id,name,config_file,value" );

	#push( @{$json_response}, { "links" => [ { "rel" => "configs", "href" => "child" } ] } );
	$self->success($json_response);
}

sub build_cdns_json {
	my $self            = shift;
	my $rs_data         = shift;
	my $default_columns = shift;
	my $columns;

	if ( defined $self->param('columns') ) {
		$columns = $self->param('columns');
	}
	else {
		$columns = $default_columns;
	}

	my (@columns) = split( /,/, $columns );
	my %columns;
	foreach my $col (@columns) {
		$columns{$col} = defined;
	}

	my @data;
	my @cols = grep { exists $columns{$_} } $rs_data->result_source->columns;

	while ( my $row = $rs_data->next ) {
		my %parameter;
		foreach my $col (@cols) {
			$parameter{$col} = $row->$col;
		}
		push( @data, \%parameter );
	}
	return \@data;
}

sub domains {
	my $self = shift;
	my @data;

	my @ccrprofs = $self->db->resultset('Profile')->search( { name => { -like => 'CCR%' } } )->get_column('id')->all();
	my $rs_pp =
		$self->db->resultset('ProfileParameter')
		->search( { profile => { -in => \@ccrprofs }, 'parameter.name' => 'domain_name', 'parameter.config_file' => 'CRConfig.json' },
		{ prefetch => [ 'parameter', 'profile' ] } );
	while ( my $row = $rs_pp->next ) {
		push(
			@data, {
				"domainName"         => $row->parameter->value,
				"parameterId"        => $row->parameter->id,
				"profileId"          => $row->profile->id,
				"profileName"        => $row->profile->name,
				"profileDescription" => $row->profile->description,
			}
		);

	}
	$self->success( \@data );
}

sub dnssec_keys {
	my $self = shift;
	if ( !&is_admin($self) ) {
		$self->alert( { Error => " - You must be an ADMIN to perform this operation!" } );
	}
	else {
		my $cdn_name = $self->param('name');
		my $keys;
		my $response_container = $self->riak_get( "dnssec", $cdn_name );
		my $get_keys = $response_container->{'response'};
		if ( $get_keys->is_success() ) {
			$keys = decode_json( $get_keys->content );
		}
		else {
			return $self->alert( { Error => " - Dnssec keys for $cdn_name do not exist!  Response was: " . $get_keys->content } );
		}
		my %new_keys = ();

		# add TLD keys to new_keys hash.  Remove this if we are checking TLD expiration
		$new_keys{$cdn_name} = $keys->{$cdn_name};
		my $z_update = 0;
		my $k_update = 0;

		#get default expiration days and ttl for DSs from CDN record
		my $default_k_exp_days = "365";
		my $default_z_exp_days = "30";

		my $cdn_ksk = $keys->{$cdn_name}->{ksk};
		foreach my $cdn_krecord (@$cdn_ksk) {
			my $cdn_kstatus = $cdn_krecord->{status};
			if ( $cdn_kstatus eq 'new' ) {    #ignore anything other than the 'new' record
				my $cdn_k_exp   = $cdn_krecord->{expirationDate};
				my $cdn_k_incep = $cdn_krecord->{inceptionDate};
				$default_k_exp_days = ( $cdn_k_exp - $cdn_k_incep ) / 86400;
			}
		}
		my $cdn_zsk = $keys->{$cdn_name}->{zsk};
		foreach my $cdn_zrecord (@$cdn_zsk) {
			my $cdn_zstatus = $cdn_zrecord->{status};
			if ( $cdn_zstatus eq 'new' ) {    #ignore anything other than the 'new' record
				my $cdn_z_exp   = $cdn_zrecord->{expirationDate};
				my $cdn_z_incep = $cdn_zrecord->{inceptionDate};
				$default_z_exp_days = ( $cdn_z_exp - $cdn_z_incep ) / 86400;
			}
		}

		#get DeliveryServices for CDN
		my $profile_id = $self->get_profile_id_by_cdn($cdn_name);
		my %search     = ( profile => $profile_id );
		my @ds_rs      = $self->db->resultset('Deliveryservice')->search( \%search );
		foreach my $ds (@ds_rs) {

			#get DNSKEY ttl for TLD
			my $dnskey_gen_multiplier;
			my $dnskey_ttl;
			my $dnskey_effective_multiplier;
			my $ds_profile = $ds->profile->name;
			my %condition = ( 'parameter.name' => 'tld.ttls.DNSKEY', 'profile.name' => $ds_profile );
			my $rs_pp =
				$self->db->resultset('ProfileParameter')->search( \%condition, { prefetch => [ { 'parameter' => undef }, { 'profile' => undef } ] } )
				->single;
			if ($rs_pp) {
				$dnskey_ttl = $rs_pp->parameter->value;
			}
			else {
				$dnskey_ttl = '60';
			}
			%condition = ( 'parameter.name' => 'DNSKEY.generation.multiplier', 'profile.name' => $ds_profile );
			$rs_pp =
				$self->db->resultset('ProfileParameter')->search( \%condition, { prefetch => [ { 'parameter' => undef }, { 'profile' => undef } ] } )
				->single;
			if ($rs_pp) {
				$dnskey_gen_multiplier = $rs_pp->parameter->value;
			}
			else {
				$dnskey_gen_multiplier = '10';
			}
			%condition = ( 'parameter.name' => 'DNSKEY.effective.multiplier', 'profile.name' => $ds_profile );
			$rs_pp =
				$self->db->resultset('ProfileParameter')->search( \%condition, { prefetch => [ { 'parameter' => undef }, { 'profile' => undef } ] } )
				->single;
			if ($rs_pp) {
				$dnskey_effective_multiplier = $rs_pp->parameter->value;
			}
			else {
				$dnskey_effective_multiplier = '2';
			}
			my $key_expiration = time() - ( $dnskey_ttl * $dnskey_gen_multiplier );

			#check if keys exist for ds
			my $xml_id  = $ds->xml_id;
			my $ds_keys = $keys->{$xml_id};
			if ( !$ds_keys ) {

				#create keys
				$self->app->log->info("Keys do not exist for ds $xml_id");
				my $ds_id = $ds->id;

				#create the ds domain name for dnssec keys
				my $domain_name = UI::DeliveryService::get_cdn_domain( $self, $ds_id );
				my $ds_regexes = UI::DeliveryService::get_regexp_set( $self, $ds_id );
				my $rs_ds = $self->db->resultset('Deliveryservice')
					->search( { 'me.xml_id' => $xml_id }, { prefetch => [ { 'type' => undef }, { 'profile' => undef } ] } );
				my $data = $rs_ds->single;
				my @example_urls = UI::DeliveryService::get_example_urls( $self, $ds_id, $ds_regexes, $data, $domain_name, $data->protocol );

				#first one is the one we want.  period at end for dnssec, substring off stuff we dont want
				my $ds_name = $example_urls[0] . ".";
				my $length = length($ds_name) - index( $ds_name, "." );
				$ds_name = substr( $ds_name, index( $ds_name, "." ) + 1, $length );

				my $inception    = time();
				my $z_expiration = $inception + ( 86400 * $default_z_exp_days );
				my $k_expiration = $inception + ( 86400 * $default_k_exp_days );

				my $zsk = $self->get_dnssec_keys( "zsk", $ds_name, $dnskey_ttl, $inception, $z_expiration, "new", $inception );
				my $ksk = $self->get_dnssec_keys( "ksk", $ds_name, $dnskey_ttl, $inception, $k_expiration, "new", $inception );

				#add to keys hash
				$new_keys{$xml_id} = { zsk => [$zsk], ksk => [$ksk] };

				#update param with current time stamp
				my $param_update = $self->db->resultset('Parameter')->find( { name => $cdn_name . ".dnssec.inception" } );
				$param_update->value($inception);
				$param_update->update();
			}

			#if keys do exist, check expiration
			else {
				my $ksk = $ds_keys->{ksk};
				foreach my $krecord (@$ksk) {
					my $kstatus = $krecord->{status};
					if ( $kstatus eq 'new' ) {    #ignore anything other than the 'new' record
						                          #check if expired
						if ( $krecord->{expirationDate} < $key_expiration ) {

							#if expired create new keys
							$self->app->log->info("The KSK keys for $xml_id are expired!");
							my $new_dnssec_keys = $self->regen_expired_keys( "ksk", $xml_id, $keys, $dnskey_effective_multiplier );
							$new_keys{$xml_id} = $new_dnssec_keys;
							$k_update = 1;
						}
					}
				}
				my $zsk = $ds_keys->{zsk};
				foreach my $zrecord (@$zsk) {
					my $zstatus = $zrecord->{status};
					if ( $zstatus eq 'new' ) {
						if ( $zrecord->{expirationDate} < $key_expiration ) {

							#if expired create new keys
							$self->app->log->info("The ZSK keys for $xml_id are expired!");
							my $new_dnssec_keys = $self->regen_expired_keys( "zsk", $xml_id, $keys, $dnskey_effective_multiplier );
							$new_keys{$xml_id} = $new_dnssec_keys;
							$z_update = 1;
						}
					}
				}

				#if not expired write existing key
				if ( $z_update == 0 && $k_update == 0 ) {
					$new_keys{$xml_id} = $keys->{$xml_id};
				}
			}
		}

		# #convert hash to json and store in Riak
		my $json_data = encode_json( \%new_keys );
		$self->riak_put( "dnssec", $cdn_name, $json_data );

		#get updated record
		$response_container = $self->riak_get( "dnssec", "$cdn_name" );
		my $response = $response_container->{"response"};
		$response->is_success()
			? $self->success( decode_json( $response->content ) )
			: $self->alert( { Error => " - A record for dnssec key $cdn_name could not be found.  Response was: " . $response->content } );
	}
}

sub regen_expired_keys {
	my $self                = shift;
	my $type                = shift;
	my $key                 = shift;
	my $existing_keys       = shift;
	my $effective_multipler = shift;
	my $regen_keys          = {};
	my $old_key;

	my $existing = $existing_keys->{$key}->{$type};
	foreach my $record (@$existing) {
		if ( $record->{status} eq 'new' ) {
			$old_key = $record;
		}
	}
	my $name            = $old_key->{name};
	my $ttl             = $old_key->{ttl};
	my $expiration      = $old_key->{expirationDate};
	my $inception       = $old_key->{inceptionDate};
	my $expiration_days = ( $expiration - $inception ) / 86400;
	my $effective_date  = $expiration - ( $ttl * $effective_multipler );

	#create new expiration and inception time
	my $new_inception = time();
	my $new_expiration = $new_inception + ( 86400 * $expiration_days );

	#generate new keys
	my $new_key = $self->get_dnssec_keys( $type, $name, $ttl, $new_inception, $new_expiration, "new", $effective_date );

	if ( $type eq "ksk" ) {

		#get existing zsk
		my @zsk = $existing_keys->{$key}->{zsk};

		#set existing ksk status to "expired"
		$old_key->{status} = "expired";
		$regen_keys = { zsk => @zsk, ksk => [ $new_key, $old_key ] };
	}
	elsif ( $type eq "zsk" ) {

		#get existing ksk
		my @ksk = $existing_keys->{$key}->{ksk};

		#set existing ksk status to "expired"
		$old_key->{status} = "expired";
		$regen_keys = { zsk => [ $new_key, $old_key ], ksk => @ksk };
	}
	return $regen_keys;
}

sub dnssec_keys_generate {
	my $self = shift;

	if ( !&is_admin($self) ) {
		$self->alert( { Error => " - You must be an ADMIN to perform this operation!" } );
	}
	else {
		my $key_type      = "dnssec";
		my $key           = $self->req->json->{key};
		my $name          = $self->req->json->{name};
		my $ttl           = $self->req->json->{ttl};
		my $k_exp_days    = $self->req->json->{kskExpirationDays};
		my $z_exp_days    = $self->req->json->{zskExpirationDays};
		my $effectiveDate = $self->req->json->{effectiveDate};
		if ( !defined($effectiveDate) ) {
			$effectiveDate = time();
		}
		my $res      = $self->generate_store_dnssec_keys( $key, $name, $ttl, $k_exp_days, $z_exp_days, $effectiveDate );
		my $response = $res->{response};
		my $rc       = $response->{_rc};
		if ( $rc eq "204" ) {
			&log( $self, "Generated dnssec keys for CDN $key", "APICHANGE" );
			$self->success("Successfully created $key_type keys for $key");
		}
		else {
			$self->alert( { Error => " - DNSSEC keys for $key could not be created.  Response was" . $response->content } );
		}
	}
}

sub delete_dnssec_keys {
	my $self     = shift;
	my $key      = $self->param('name');
	my $key_type = "dnssec";
	my $response;
	if ( !&is_admin($self) ) {
		$self->alert( { Error => " - You must be an ADMIN to perform this operation!" } );
	}
	else {
		$self->app->log->info("deleting key_type = $key_type, key = $key");
		my $response_container = $self->riak_delete( $key_type, $key );
		$response = $response_container->{"response"};
		if ( $response->is_success() ) {
			&log( $self, "Deleted dnssec keys for CDN $key", "UICHANGE" );
			$self->success("Successfully deleted $key_type keys for $key");
		}
		else {
			$self->alert( { Error => " - SSL keys for key type $key_type and key $key could not be deleted.  Response was" . $response->content } );
		}
	}
}

sub catch_all {
	my $self     = shift;
	my $mimetype = $self->req->headers->content_type;

	if ( defined( $self->current_user() ) ) {
		return $self->not_found();
	}
	else {
		return $self->unauthorized();
	}
}

1;
