package API::Topology;
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
use UI::Render;
use List::Compare;
use JSON;
use Data::Dumper;
use Mojo::Base 'Mojolicious::Controller';
use Time::HiRes qw(gettimeofday);

sub gen_crconfig_json {
	my $self     = shift;
	my $cdn_name = shift;
	my $data_obj;
	my %type_to_name;
	my $ccr_profile_id;
	my $rascal_profile_id;
	my $ccr_domain_name = "";

	#my @cache_rascal_profiles;
	my @cdn_profiles;
	$SIG{__WARN__} = sub { warn $_[0] unless $_[0] =~ m/Prefetching multiple has_many rels deliveryservice_servers/ };

	$data_obj->{'stats'}->{'CDN_name'}   = $cdn_name;
	$data_obj->{'stats'}->{'date'}       = time();
	$data_obj->{'stats'}->{'tm_version'} = &tm_version();
	$data_obj->{'stats'}->{'tm_path'}    = $self->req->url->path->{'path'};
	$data_obj->{'stats'}->{'tm_host'}    = $self->req->headers->host;
	$data_obj->{'stats'}->{'tm_user'}    = $self->current_user()->{username};

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

#@cache_rascal_profiles = $self->db->resultset('Profile')->search( { id => { -in => \@cdn_profiles }, name => [{ like => 'EDGE%'}, {like => 'MID%'}, {like => 'RASCAL%'}, {like => 'CDSIS%'} ] } )->get_column('id')->all();
	}
	else {
		my $e = Mojo::Exception->throw( "Parameter ID not found for CDN_name: " . $cdn_name );
	}

	my %condition = ( 'profile_parameters.profile' => $ccr_profile_id, 'config_file' => 'CRConfig.json' );
	my $rs_config = $self->db->resultset('Parameter')->search( \%condition, { join => 'profile_parameters' } );
	while ( my $row = $rs_config->next ) {
		if ( $row->name eq 'domain_name' ) {
			$ccr_domain_name = $row->value;
			$data_obj->{'config'}->{ $row->name } = $row->value;
		}
		elsif ( $row->name =~ m/^tld/ ) {
			my $param = $row->name;
			$param =~ s/tld\.//;
			( my $top_key, my $second_key ) = split( /\./, $param );
			$data_obj->{'config'}->{$top_key}->{$second_key} = $row->value;
		}
		else {
			$data_obj->{'config'}->{ $row->name } = $row->value;
		}
	}
	my $rs_loc = $self->db->resultset('CachegroupParameter')->search( { 'parameter' => $cdnname_param_id }, { prefetch => 'cachegroup' } );

	while ( my $row = $rs_loc->next ) {
		$data_obj->{'edgeLocations'}->{ $row->cachegroup->name }->{'latitude'}  = $row->cachegroup->latitude + 0;
		$data_obj->{'edgeLocations'}->{ $row->cachegroup->name }->{'longitude'} = $row->cachegroup->longitude + 0;
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
			$data_obj->{'monitors'}->{ $row->host_name }->{'fqdn'}       = $row->host_name . "." . $row->domain_name;
			$data_obj->{'monitors'}->{ $row->host_name }->{'status'}     = $row->status->name;
			$data_obj->{'monitors'}->{ $row->host_name }->{'cachegroup'} = $row->cachegroup->name;
			$data_obj->{'monitors'}->{ $row->host_name }->{'port'}       = $row->tcp_port;
			$data_obj->{'monitors'}->{ $row->host_name }->{'ip'}         = $row->ip_address;
			$data_obj->{'monitors'}->{ $row->host_name }->{'ip6'}        = $row->ip6_address;
			$data_obj->{'monitors'}->{ $row->host_name }->{'profile'}    = $row->profile->name;

		}
		elsif ( $row->type->name eq "CCR" ) {
			my $rs_param = $self->db->resultset('Parameter')
				->search( { 'profile_parameters.profile' => $row->profile->id, 'name' => 'api.port' }, { join => 'profile_parameters' } );
			my $r = $rs_param->single;
			my $port = ( defined($r) && defined( $r->value ) ) ? $r->value : 80;

			$data_obj->{'contentRouters'}->{ $row->host_name }->{'fqdn'}       = $row->host_name . "." . $row->domain_name;
			$data_obj->{'contentRouters'}->{ $row->host_name }->{'status'}     = $row->status->name;
			$data_obj->{'contentRouters'}->{ $row->host_name }->{'cachegroup'} = $row->cachegroup->name;
			$data_obj->{'contentRouters'}->{ $row->host_name }->{'port'}       = $row->tcp_port;
			$data_obj->{'contentRouters'}->{ $row->host_name }->{'api.port'}   = $port;
			$data_obj->{'contentRouters'}->{ $row->host_name }->{'ip'}         = $row->ip_address;
			$data_obj->{'contentRouters'}->{ $row->host_name }->{'ip6'}        = $row->ip6_address;
			$data_obj->{'contentRouters'}->{ $row->host_name }->{'profile'}    = $row->profile->name;
		}
		else {
			if ( !exists $cache_tracker{ $row->id } ) {
				$cache_tracker{ $row->id } = $row->host_name;
			}
			$data_obj->{'contentServers'}->{ $row->host_name }->{'cachegroupId'}  = $row->cachegroup->name;
			$data_obj->{'contentServers'}->{ $row->host_name }->{'fqdn'}          = $row->host_name . "." . $row->domain_name;
			$data_obj->{'contentServers'}->{ $row->host_name }->{'port'}          = $row->tcp_port;
			$data_obj->{'contentServers'}->{ $row->host_name }->{'interfaceName'} = $row->interface_name;
			$data_obj->{'contentServers'}->{ $row->host_name }->{'status'}        = $row->status->name;
			$data_obj->{'contentServers'}->{ $row->host_name }->{'ip'}            = $row->ip_address;
			$data_obj->{'contentServers'}->{ $row->host_name }->{'ip6'}           = ( $row->ip6_address || "" );
			$data_obj->{'contentServers'}->{ $row->host_name }->{'profile'}       = $row->profile->name;
			$data_obj->{'contentServers'}->{ $row->host_name }->{'type'}          = $row->type->name;
			$data_obj->{'contentServers'}->{ $row->host_name }->{'hashId'}        = $row->xmpp_id;
		}
	}
	my $regexps;
	my $rs_ds = $self->db->resultset('Deliveryservice')
		->search( { 'me.profile' => $ccr_profile_id, 'active' => 1 }, { prefetch => [ 'deliveryservice_servers', 'deliveryservice_regexes', 'type' ] } );
	while ( my $row = $rs_ds->next ) {
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
				$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'matchsets'}->[ $subrow->set_number ]->{'protocol'} = $protocol;
				$regex_to_props->{ $subrow->{'_column_data'}->{'regex'} }->{'pattern'} =
					$regex_tracker->{ $subrow->{'_column_data'}->{'regex'} }->{'pattern'};
				$regex_to_props->{ $subrow->{'_column_data'}->{'regex'} }->{'set_number'} = $subrow->set_number;
				$regex_to_props->{ $subrow->{'_column_data'}->{'regex'} }->{'type'} = $regex_tracker->{ $subrow->{'_column_data'}->{'regex'} }->{'type'};
				if ( $regex_to_props->{ $subrow->{'_column_data'}->{'regex'} }->{'type'} eq 'HOST_REGEXP' ) {
					$ds_to_remap{ $row->xml_id }->[ $subrow->set_number ] = $regex_to_props->{ $subrow->{'_column_data'}->{'regex'} }->{'pattern'};
				}
			}
		}
		foreach my $regex ( sort keys %{$regex_to_props} ) {
			my $set_number = $regex_to_props->{$regex}->{'set_number'};
			my $pattern    = $regex_to_props->{$regex}->{'pattern'};
			my $type       = $regex_to_props->{$regex}->{'type'};
			if ( $type eq 'HOST_REGEXP' ) {
				push(
					@{ $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'matchsets'}->[$set_number]->{'matchlist'} },
					{ 'match-type' => 'HOST', 'regex' => $pattern }
				);
			}
			elsif ( $type eq 'PATH_REGEXP' ) {
				push(
					@{ $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'matchsets'}->[$set_number]->{'matchlist'} },
					{ 'match-type' => 'PATH', 'regex' => $pattern }
				);
			}
			elsif ( $type eq 'HEADER_REGEXP' ) {
				push(
					@{ $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'matchsets'}->[$set_number]->{'matchlist'} },
					{ 'match-type' => 'HEADER', 'regex' => $pattern }
				);
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
			foreach my $server ( keys %server_subrow_dedup ) {
				foreach my $host ( @{ $ds_to_remap{ $row->xml_id } } ) {
					my $remap;
					if ( $host =~ m/\.\*$/ ) {
						my $host_copy = $host;
						$host_copy =~ s/$host_regex1//g;
						if ( $protocol eq 'DNS' ) {
							$remap = 'edge' . $host_copy . $ccr_domain_name;
						}
						else {
							$remap = $cache_tracker{$server} . $host_copy . $ccr_domain_name;
						}
					}
					else {
						$remap = $host;
					}
					push( @{ $data_obj->{'contentServers'}->{ $cache_tracker{$server} }->{'deliveryServices'}->{ $row->xml_id } }, $remap );
				}
			}
		}

		$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'ttl'} = $row->ccr_dns_ttl;
		my $geo_limit = $row->geo_limit;
		if ( $geo_limit == 1 ) {
			$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'coverageZoneOnly'} = 'true';
		}
		elsif ( $geo_limit == 2 ) {
			$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'coverageZoneOnly'} = 'false';
			$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'geoEnabled'} = [ { 'countryCode' => 'US' } ];
		}
		elsif ( $geo_limit == 3 ) {
			$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'coverageZoneOnly'} = 'false';
			$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'geoEnabled'} = [ { 'countryCode' => 'CA' } ];
		}
		else {
			$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'coverageZoneOnly'} = 'false';
		}
		if ( $protocol =~ m/DNS/ ) {

			#$data_obj->{'deliveryServices'}->{$row->xml_id}->{'matchsets'}->[0]->{'protocol'} = 'DNS';
			if ( defined( $row->dns_bypass_ip ) && $row->dns_bypass_ip ne "" ) {
				$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'bypassDestination'}->{'DNS'}->{'ip'} = $row->dns_bypass_ip;
			}
			if ( defined( $row->dns_bypass_ip6 ) && $row->dns_bypass_ip6 ne "" ) {
				$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'bypassDestination'}->{'DNS'}->{'ip6'} = $row->dns_bypass_ip6;
			}
			if ( defined( $row->dns_bypass_ttl ) && $row->dns_bypass_ttl ne "" ) {
				$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'bypassDestination'}->{'DNS'}->{'ttl'} = $row->dns_bypass_ttl;
			}
			if ( defined( $row->max_dns_answers ) && $row->max_dns_answers ne "" ) {
				$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'maxDnsIpsForLocation'} = $row->max_dns_answers;
			}
		}
		elsif ( $protocol =~ m/HTTP/ ) {

			#$data_obj->{'deliveryServices'}->{$row->xml_id}->{'matchsets'}->[0]->{'protocol'} = 'HTTP';
			if ( defined( $row->http_bypass_fqdn ) && $row->http_bypass_fqdn ne "" ) {
				my $full = $row->http_bypass_fqdn;
				my $port;
				my $fqdn;
				if ( $full =~ m/\:/ ) {
					( $fqdn, $port ) = split( /\:/, $full );
				}
				else {
					$fqdn = $full;
					$port = '80';
				}
				$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'bypassDestination'}->{'HTTP'}->{'fqdn'} = $fqdn;
				$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'bypassDestination'}->{'HTTP'}->{'port'} = $port;
			}
		}
		if ( defined( $row->miss_lat ) && $row->miss_lat ne "" ) {
			$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'missLocation'}->{'lat'} = $row->miss_lat;
		}
		if ( defined( $row->miss_long ) && $row->miss_long ne "" ) {
			$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'missLocation'}->{'long'} = $row->miss_long;
		}
		$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'ttls'} =
			{ 'A' => $row->ccr_dns_ttl, 'AAAA' => $row->ccr_dns_ttl, 'NS' => "3600", 'SOA' => "86400" };
		$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'soa'}->{'minimum'} = "30";
		$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'soa'}->{'expire'}  = "604800";
		$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'soa'}->{'retry'}   = "7200";
		$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'soa'}->{'refresh'} = "28800";
		$data_obj->{'deliveryServices'}->{ $row->xml_id }->{'soa'}->{'admin'}   = "twelve_monkeys";

		my $rs_dns = $self->db->resultset('Staticdnsentry')->search(
			{ 'deliveryservice.active' => 1, 'deliveryservice.profile' => $ccr_profile_id },
			{ prefetch => [ 'deliveryservice', 'type' ], columns => [ 'host', 'type', 'ttl', 'address' ] }
		);
		while ( my $dns_row = $rs_dns->next ) {
			my $dns_obj;
			$dns_obj->{'name'}  = $dns_row->host;
			$dns_obj->{'ttl'}   = $dns_row->ttl;
			$dns_obj->{'value'} = $dns_row->address;

			my $type = $dns_row->type->name;
			$type =~ s/\_RECORD//g;
			$dns_obj->{'type'} = $type;
			push( @{ $data_obj->{'deliveryServices'}->{ $dns_row->deliveryservice->xml_id }->{'staticDnsEntries'} }, $dns_obj );
		}

	}
	return ($data_obj);
}

1;
