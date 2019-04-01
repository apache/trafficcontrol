package UI::ConfigFiles;

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
use UI::Utils;

use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use Date::Manip;
use NetAddr::IP;
use UI::DeliveryService;
use JSON;
use API::DeliveryService::KeysUrlSig qw(URL_SIG_KEYS_BUCKET);
use URI;

my $dispatch_table ||= {
	"logs_xml.config"         => sub { logs_xml_dot_config(@_) },
	"logging.config"		  => sub { logging_dot_config(@_) },
	"logging.yaml"		  	  => sub { logging_dot_yaml(@_) },
	"cacheurl.config"         => sub { cacheurl_dot_config(@_) },
	"records.config"          => sub { generic_config(@_) },
	"plugin.config"           => sub { generic_config(@_) },
	"astats.config"           => sub { generic_config(@_) },
	"volume.config"           => sub { volume_dot_config(@_) },
	"hosting.config"          => sub { hosting_dot_config(@_) },
	"storage.config"          => sub { storage_dot_config(@_) },
	"50-ats.rules"            => sub { ats_dot_rules(@_) },
	"cache.config"            => sub { cache_dot_config(@_) },
	"remap.config"            => sub { remap_dot_config(@_) },
	"parent.config"           => sub { parent_dot_config(@_) },
	"sysctl.conf"             => sub { generic_config(@_) },
	"ip_allow.config"         => sub { ip_allow_dot_config(@_) },
	"12M_facts"               => sub { facts(@_) },
	"regex_revalidate.config" => sub { regex_revalidate_dot_config(@_) },
	"drop_qstring.config"     => sub { drop_qstring_dot_config(@_) },
	"bg_fetch.config"         => sub { bg_fetch_dot_config(@_) },

	"url_sig_.config"      => sub { url_sig_config(@_) },
	"hdr_rw_.config"       => sub { header_rewrite_dot_config(@_) },
	"set_dscp_.config"     => sub { header_rewrite_dscp_dot_config(@_) },
	"to_ext_.config"       => sub { to_ext_dot_config(@_) },
	"regex_remap_.config"  => sub { regex_remap_dot_config(@_) },
	"cacheurl_.config"     => sub { cacheurl_dot_config(@_) },
	"all"                  => sub { gen_fancybox_data(@_) },
	"ssl_multicert.config" => sub { ssl_multicert_dot_config(@_) },

};

my $separator ||= {
	"records.config"  => " ",
	"plugin.config"   => " ",
	"sysctl.conf"     => " = ",
	"url_sig_.config" => " = ",
	"astats.config"   => "=",
};

sub genfiles {
	my $self = shift;
	my $mode = $self->param('mode');
	my $id   = $self->param('id');
	my $file = $self->param('filename');

	my $org_name = $file;

	$file =~ s/^url_sig_.*\.config$/url_sig_\.config/;
	$file =~ s/^hdr_rw_.*\.config$/hdr_rw_\.config/;
	$file =~ s/^set_dscp_.*\.config$/set_dscp_\.config/;
	$file =~ s/^regex_remap_.*\.config$/regex_remap_\.config/;
	$file =~ s/^cacheurl_.*\.config$/cacheurl_\.config/;
	$file =~ s/^to_ext_.*\.config$/to_ext_\.config/;

	my $text = undef;
	if ( $mode eq 'view' ) {

		if ( defined( $dispatch_table->{$file} ) ) {
			$text = $dispatch_table->{$file}->( $self, $id, $org_name );
		}
		else {
			$text = &take_and_bake( $self, $id, $org_name );
		}
	}
	if ( $text =~ /^Error/ ) {
		$self->internal_server_error( { Error => $text } );
	}

	if ( $file ne "all" ) {
		$self->res->headers->content_type("application/download");
		$self->res->headers->content_disposition("attachment; filename=\"$org_name\"");
		$self->render( text => $text, format => 'txt' );
	}
	else {

		# ignore $text, the good stuff is in the stash
		$self->stash( fbox_layout => 1 );
	}
}

sub gen_fancybox_data {
	my $self     = shift;
	my $id       = shift;
	my $filename = shift;

	my $file_text;
	my $server  = $self->server_data($id);
	my $ds_data = $self->ds_data($server);
	my $rs      = $self->db->resultset('ProfileParameter')->search(
		{ -and => [ profile => $server->profile->id, 'parameter.name' => 'location' ] },
		{ prefetch => [ { parameter => undef }, { profile => undef } ] }
	);
	while ( my $row = $rs->next ) {
		my $file = $row->parameter->config_file;

		# print "Genning $file\n";
		my $org_name = $file;
		$file =~ s/^url_sig_.*\.config$/url_sig_\.config/;
		$file =~ s/^hdr_rw_.*\.config$/hdr_rw_\.config/;
		$file =~ s/^set_dscp_.*\.config$/set_dscp_\.config/;
		$file =~ s/^regex_remap_.*\.config$/regex_remap_\.config/;
		$file =~ s/^cacheurl_.*\.config$/cacheurl_\.config/;
		$file =~ s/^to_ext_.*\.config$/to_ext_\.config/;

		my $text = "boo";
		if ( defined( $dispatch_table->{$file} ) ) {
			$text = $dispatch_table->{$file}->( $self, $id, $org_name, $ds_data );
		}
		else {
			$text = &take_and_bake( $self, $id, $org_name, $ds_data );
		}
		$file_text->{$org_name} = $text;
	}
	$self->stash( file_text => $file_text );
	$self->stash( host_name => $server->host_name );
}

sub server_data {
	my $self = shift;
	my $id   = shift;

	my $server;

	if ( $id =~ /^\d+$/ ) {
		$server = $self->db->resultset('Server')->search( { 'me.id' => $id }, { prefetch => [ 'profile', 'type', 'cachegroup', 'cdn' ] } )->single;
	}
	else {
		$server = $self->db->resultset('Server')->search( { host_name => $id }, { prefetch => [ 'profile', 'type', 'cachegroup', 'cdn' ] } )->single;
	}

	return $server;
}

sub header_comment {
	my $self      = shift;
	my $host_name = shift;

	my $text = "# DO NOT EDIT - Generated for " . $host_name . " by " . &name_version_string($self) . " on " . `date`;
	return $text;
}

sub ds_data {
	my $self   = shift;
	my $server = shift;

	my $dsinfo;

	#	if ( defined( $self->app->session->{dsinfo} ) ) {
	#		$dsinfo = $self->app->session->{dsinfo};
	#		return $dsinfo;
	#	}
	$dsinfo->{host_name}   = $server->host_name;
	$dsinfo->{domain_name} = $server->domain_name;

	my @server_ids = ();
	my $rs;
	if ( $server->type->name =~ m/^MID/ ) {
		# the mids will do all deliveryservices in this CDN
		$rs = $self->db->resultset('DeliveryServiceInfoForDomainList')->search( {}, { bind => [ $server->cdn->name ] } );
	}
	else {
		$rs = $self->db->resultset('DeliveryServiceInfoForServerList')->search( {}, { bind => [ $server->id ] } );
	}

	my $j = 0;
	while ( my $row = $rs->next ) {
		my $org_server                  = $row->org_server_fqdn;
		my $dscp                        = $row->dscp;
		my $re_type                     = $row->re_type;
		my $ds_type                     = $row->ds_type;
		my $signed                      = defined( $row->signing_algorithm ) ? ( $row->signing_algorithm eq "url_sig" ? \1 : \0 ) : \0;
		my $signing_algorithm           = $row->signing_algorithm;
		my $qstring_ignore              = $row->qstring_ignore;
		my $ds_xml_id                   = $row->xml_id;
		my $ds_domain                   = $row->domain_name;
		my $edge_header_rewrite         = $row->edge_header_rewrite;
		my $mid_header_rewrite          = $row->mid_header_rewrite;
		my $regex_remap                 = $row->regex_remap;
		my $protocol                    = $row->protocol;
		my $range_request_handling      = $row->range_request_handling;
		my $origin_shield               = $row->origin_shield;
		my $cacheurl                    = $row->cacheurl;
		my $remap_text                  = $row->remap_text;
		my $multi_site_origin           = $row->multi_site_origin;
		my $multi_site_origin_algorithm = 0;

		if ( $re_type eq 'HOST_REGEXP' ) {
			my $host_re = $row->pattern;
			my $map_to  = $org_server . "/";
			if ( $host_re =~ /\.\*$/ ) {
				my $re = $host_re;
				$re =~ s/\\//g;
				$re =~ s/\.\*//g;
				my $hname = $ds_type =~ /^DNS/ ? $row->routing_name : "__http__";
				my $portstr = "";
				if ( $hname eq "__http__" && $server->tcp_port > 0 && $server->tcp_port != 80 ) {
					$portstr = ":" . $server->tcp_port;
				}
				my $map_from = "http://" . $hname . $re . $ds_domain . $portstr . "/";
				if ( $protocol == 0 ) {
					$dsinfo->{dslist}->[$j]->{"remap_line"}->{$map_from} = $map_to;
				}
				elsif ( $protocol == 1 || $protocol == 3 ) {
					$map_from = "https://" . $hname . $re . $ds_domain . "/";
					$dsinfo->{dslist}->[$j]->{"remap_line"}->{$map_from} = $map_to;
				}
				elsif ( $protocol == 2 ) {

					#add the first one with http
					$dsinfo->{dslist}->[$j]->{"remap_line"}->{$map_from} = $map_to;

					#add the second one for https
					my $map_from2 = "https://" . $hname . $re . $ds_domain . "/";
					$dsinfo->{dslist}->[$j]->{"remap_line2"}->{$map_from2} = $map_to;
				}
			}
			else {
				my $map_from = "http://" . $host_re . "/";
				if ( $protocol == 0 ) {
					$dsinfo->{dslist}->[$j]->{"remap_line"}->{$map_from} = $map_to;
				}
				elsif ( $protocol == 1 || $protocol == 3 ) {
					$map_from = "https://" . $host_re . "/";
					$dsinfo->{dslist}->[$j]->{"remap_line"}->{$map_from} = $map_to;
				}
				elsif ( $protocol == 2 ) {

					#add the first with http
					$dsinfo->{dslist}->[$j]->{"remap_line"}->{$map_from} = $map_to;

					#add the second with https
					my $map_from2 = "https://" . $host_re . "/";
					$dsinfo->{dslist}->[$j]->{"remap_line2"}->{$map_from2} = $map_to;
				}
			}
		}
		$dsinfo->{dslist}->[$j]->{"dscp"}                        = $dscp;
		$dsinfo->{dslist}->[$j]->{"org"}                         = $org_server;
		$dsinfo->{dslist}->[$j]->{"type"}                        = $ds_type;
		$dsinfo->{dslist}->[$j]->{"domain"}                      = $ds_domain;
		$dsinfo->{dslist}->[$j]->{"signed"}                      = $signed;
		$dsinfo->{dslist}->[$j]->{"signing_algorithm"}           = $signing_algorithm;
		$dsinfo->{dslist}->[$j]->{"qstring_ignore"}              = $qstring_ignore;
		$dsinfo->{dslist}->[$j]->{"ds_xml_id"}                   = $ds_xml_id;
		$dsinfo->{dslist}->[$j]->{"edge_header_rewrite"}         = $edge_header_rewrite;
		$dsinfo->{dslist}->[$j]->{"mid_header_rewrite"}          = $mid_header_rewrite;
		$dsinfo->{dslist}->[$j]->{"regex_remap"}                 = $regex_remap;
		$dsinfo->{dslist}->[$j]->{"range_request_handling"}      = $range_request_handling;
		$dsinfo->{dslist}->[$j]->{"origin_shield"}               = $origin_shield;
		$dsinfo->{dslist}->[$j]->{"cacheurl"}                    = $cacheurl;
		$dsinfo->{dslist}->[$j]->{"remap_text"}                  = $remap_text;
		$dsinfo->{dslist}->[$j]->{"multi_site_origin"}           = $multi_site_origin;
		$dsinfo->{dslist}->[$j]->{"multi_site_origin_algorithm"} = $multi_site_origin_algorithm;

		if ( defined($edge_header_rewrite) ) {
			my $fname = "hdr_rw_" . $ds_xml_id . ".config";
			$dsinfo->{dslist}->[$j]->{"hdr_rw_file"} = $fname;
		}
		if ( defined($mid_header_rewrite) ) {
			my $fname = "hdr_rw_mid_" . $ds_xml_id . ".config";
			$dsinfo->{dslist}->[$j]->{"mid_hdr_rw_file"} = $fname;
		}
		if ( defined($cacheurl) ) {
			my $fname = "cacheurl_" . $ds_xml_id . ".config";
			$dsinfo->{dslist}->[$j]->{"cacheurl_file"} = $fname;
		}

		if ( defined( $row->profile ) ) {
			my $dsparamrs = $self->db->resultset('ProfileParameter')->search( { profile => $row->profile }, { prefetch => [ 'profile', 'parameter' ] } );
			while ( my $prow = $dsparamrs->next ) {
				$dsinfo->{dslist}->[$j]->{'param'}->{ $prow->parameter->config_file }->{ $prow->parameter->name } = $prow->parameter->value;
			}
		}

		$j++;
	}

	#	$self->app->session->{dsinfo} = $dsinfo;
	return $dsinfo;
}

sub param_data {
	my $self     = shift;
	my $server   = shift;
	my $filename = shift;
	my $data;

	my $rs = $self->db->resultset('ProfileParameter')->search( { -and => [ profile => $server->profile->id, 'parameter.config_file' => $filename ] },
		{ prefetch => [ { parameter => undef }, { profile => undef } ] } );
	while ( my $row = $rs->next ) {
		if ( $row->parameter->name eq "location" ) {
			next;
		}
		my $value = $row->parameter->value;

		# some files have multiple lines with the same key... handle that with param id.
		my $key = $row->parameter->name;
		if ( defined( $data->{$key} ) ) {
			$key .= "__" . $row->parameter->id;
		}
		if ( $value =~ /^STRING __HOSTNAME__$/ ) {
			$value = "STRING " . $server->host_name . "." . $server->domain_name;
		}
		$data->{$key} = $value;
	}
	return $data;
}

sub profile_param_value {
	my $self       = shift;
	my $pid        = shift;
	my $file       = shift;
	my $param_name = shift;
	my $default    = shift;

	# assign $ds_domain, $weight and $port, and cache the results %profile_cache
	my $param =
		$self->db->resultset('ProfileParameter')
		->search( { -and => [ profile => $pid, 'parameter.config_file' => $file, 'parameter.name' => $param_name ] },
		{ prefetch => [ 'parameter', 'profile' ] } )->first();

	return ( defined $param ? $param->parameter->value : $default );
}

sub by_parent_rank {
	my ($arank) = $a->{"rank"};
	my ($brank) = $b->{"rank"};
	( $arank || 1 ) <=> ( $brank || 1 );
}

sub parent_data {
	my $self   = shift;
	my $server = shift;

	my @parent_cachegroup_ids;
	my @secondary_parent_cachegroup_ids;
	my $org_loc_type_id = &type_id( $self, "ORG_LOC" );
	if ( $server->type->name =~ m/^MID/ ) {

		# multisite origins take all the org groups in to account
		@parent_cachegroup_ids = $self->db->resultset('Cachegroup')->search( { type => $org_loc_type_id } )->get_column('id')->all();
	}
	else {
		@parent_cachegroup_ids =
			grep {defined} $self->db->resultset('Cachegroup')->search( { id => $server->cachegroup->id } )->get_column('parent_cachegroup_id')->all();
		@secondary_parent_cachegroup_ids =
			grep {defined}
			$self->db->resultset('Cachegroup')->search( { id => $server->cachegroup->id } )->get_column('secondary_parent_cachegroup_id')->all();
	}

	# get the server's cdn domain
	my $server_domain = $self->get_cdn_domain_by_profile_id( $server->profile->id );

	my %profile_cache;
	my %deliveryservices;
	my %parent_info;

	$self->cachegroup_profiles( \@parent_cachegroup_ids,           \%profile_cache, \%deliveryservices );
	$self->cachegroup_profiles( \@secondary_parent_cachegroup_ids, \%profile_cache, \%deliveryservices );
	foreach my $prefix ( keys %deliveryservices ) {
		foreach my $row ( @{ $deliveryservices{$prefix} } ) {
			my $pid              = $row->profile->id;
                        if ( $profile_cache{$pid}->{not_a_parent} ne 'false' ) {
                            next;
                        }
			my $ds_domain        = $profile_cache{$pid}->{domain_name};
			my $weight           = $profile_cache{$pid}->{weight};
			my $port             = $profile_cache{$pid}->{port};
			my $use_ip_address   = $profile_cache{$pid}->{use_ip_address};
			my $rank             = $profile_cache{$pid}->{rank};
			my $primary_parent   = $server->cachegroup->parent_cachegroup_id // -1;
			my $secondary_parent = $server->cachegroup->secondary_parent_cachegroup_id // -1;

			if ( defined($ds_domain) && defined($server_domain) && $ds_domain eq $server_domain ) {
				my %p = (
					host_name        => $row->host_name,
					port             => defined($port) ? $port : $row->tcp_port,
					domain_name      => $row->domain_name,
					weight           => $weight,
					use_ip_address   => $use_ip_address,
					rank             => $rank,
					ip_address       => $row->ip_address,
					primary_parent   => ( $primary_parent == $row->cachegroup->id ) ? 1 : 0,
					secondary_parent => ( $secondary_parent == $row->cachegroup->id ) ? 1 : 0,
				);
				push @{ $parent_info{$prefix} }, \%p;
			}
		}
	}
	return \%parent_info;
}

sub cachegroup_profiles {
	my $self             = shift;
	my $ids              = shift;
	my $profile_cache    = shift;
	my $deliveryservices = shift;

	if ( !@$ids ) {
		return;    # nothing to see here..
	}
	my $online   = &admin_status_id( $self, "ONLINE" );
	my $reported = &admin_status_id( $self, "REPORTED" );

	my %condition = (
		status     => { -in => [ $online, $reported ] },
		cachegroup => { -in => $ids }
	);

	my $rs_parent = $self->db->resultset('Server')->search( \%condition, { prefetch => [ 'cachegroup', 'status', 'type', 'profile', 'cdn' ] } );

	while ( my $row = $rs_parent->next ) {

		next unless ( $row->type->name eq 'ORG' || $row->type->name =~ m/^EDGE/ || $row->type->name =~ m/^MID/ );
		if ( $row->type->name eq 'ORG' ) {
			my $rs_ds = $self->db->resultset('DeliveryserviceServer')->search( { server => $row->id }, { prefetch => ['deliveryservice'] } );
			while ( my $ds_row = $rs_ds->next ) {
				my $ds_domain = UI::DeliveryService::compute_org_server_fqdn($self, $ds_row->deliveryservice->id);
				$ds_domain =~ s/https?:\/\/(.*)/$1/;
				push( @{ $deliveryservices->{$ds_domain} }, $row );
			}
		}
		else {
			push( @{ $deliveryservices->{all_parents} }, $row );
		}

		# get the profile info, and cache it in %profile_cache
		my $pid = $row->profile->id;
		if ( !defined( $profile_cache->{$pid} ) ) {

			# assign $ds_domain, $weight and $port, and cache the results %profile_cache
			$profile_cache->{$pid} = {
				domain_name    => $row->cdn->domain_name,
				weight         => $self->profile_param_value( $pid, 'parent.config', 'weight', '0.999' ),
				port           => $self->profile_param_value( $pid, 'parent.config', 'port', undef ),
				use_ip_address => $self->profile_param_value( $pid, 'parent.config', 'use_ip_address', 0 ),
				rank           => $self->profile_param_value( $pid, 'parent.config', 'rank', 1 ),
				not_a_parent   => $self->profile_param_value( $pid, 'parent.config', 'not_a_parent', 'false' ),
			};
		}
	}
}

sub ip_allow_data {
	my $self   = shift;
	my $server = shift;

	my $ipallow;
	$ipallow = ();

	my $i = 0;

	# localhost is trusted.
	$ipallow->[$i]->{src_ip} = '127.0.0.1';
	$ipallow->[$i]->{action} = 'ip_allow';
	$ipallow->[$i]->{method} = "ALL";
	$i++;
	$ipallow->[$i]->{src_ip} = '::1';
	$ipallow->[$i]->{action} = 'ip_allow';
	$ipallow->[$i]->{method} = "ALL";
	$i++;

	# default for coalesce_ipv4 = 24, 5 and for ipv6 48, 5; override with the parameters in the server profile.
	my $coalesce_masklen_v4 = 24;
	my $coalesce_number_v4  = 5;
	my $coalesce_masklen_v6 = 48;
	my $coalesce_number_v6  = 5;
	my $rs_parameter =
		$self->db->resultset('ProfileParameter')->search( { profile => $server->profile->id }, { prefetch => [ "parameter", "profile" ] } );

	while ( my $row = $rs_parameter->next ) {
		if ( $row->parameter->name eq 'purge_allow_ip' && $row->parameter->config_file eq 'ip_allow.config' ) {
			$ipallow->[$i]->{src_ip} = $row->parameter->value;
			$ipallow->[$i]->{action} = "ip_allow";
			$ipallow->[$i]->{method} = "ALL";
			$i++;
		}
		elsif ( $row->parameter->name eq 'coalesce_masklen_v4' && $row->parameter->config_file eq 'ip_allow.config' ) {
			$coalesce_masklen_v4 = $row->parameter->value;
		}
		elsif ( $row->parameter->name eq 'coalesce_number_v4' && $row->parameter->config_file eq 'ip_allow.config' ) {
			$coalesce_number_v4 = $row->parameter->value;
		}
		elsif ( $row->parameter->name eq 'coalesce_masklen_v6' && $row->parameter->config_file eq 'ip_allow.config' ) {
			$coalesce_masklen_v6 = $row->parameter->value;
		}
		elsif ( $row->parameter->name eq 'coalesce_number_v6' && $row->parameter->config_file eq 'ip_allow.config' ) {
			$coalesce_number_v6 = $row->parameter->value;
		}
	}

	if ( $server->type->name =~ m/^MID/ ) {
		my @edge_locs = $self->db->resultset('Cachegroup')->search( { parent_cachegroup_id => $server->cachegroup->id } )->get_column('id')->all();
		my %allow_locs;
		foreach my $loc (@edge_locs) {
			$allow_locs{$loc} = 1;
		}

		# get all the EDGE and RASCAL nets
		my @allowed_netaddrips;
		my @allowed_ipv6_netaddrips;
		my @types;
		push( @types, &type_ids( $self, 'EDGE%', 'server' ) );
		my $rtype = &type_id( $self, 'RASCAL' );
		push( @types, $rtype );
		my $rs_allowed = $self->db->resultset('Server')->search( { 'me.type' => { -in => \@types } }, { prefetch => [ 'type', 'cachegroup' ] } );

		while ( my $allow_row = $rs_allowed->next ) {
			if ( $allow_row->type->id == $rtype
				|| ( defined( $allow_locs{ $allow_row->cachegroup->id } ) && $allow_locs{ $allow_row->cachegroup->id } == 1 ) )
			{
				my $ipv4 = NetAddr::IP->new( $allow_row->ip_address, $allow_row->ip_netmask );

				if ( defined($ipv4) ) {
					push( @allowed_netaddrips, $ipv4 );
				}
				else {
					$self->app->log->error(
						$allow_row->host_name . " has an invalid IPv4 address; excluding from ip_allow data for " . $server->host_name );
				}

				if ( defined $allow_row->ip6_address ) {
					my $ipv6 = NetAddr::IP->new( $allow_row->ip6_address );

					if ( defined($ipv6) ) {
						push( @allowed_ipv6_netaddrips, NetAddr::IP->new( $allow_row->ip6_address ) );
					}
					else {
						$self->app->log->error(
							$allow_row->host_name . " has an invalid IPv6 address; excluding from ip_allow data for " . $server->host_name );
					}
				}
			}
		}

		# compact, coalesce and compact combined list again
		my @compacted_list = NetAddr::IP::Compact(@allowed_netaddrips);
		my $coalesced_list = NetAddr::IP::Coalesce( $coalesce_masklen_v4, $coalesce_number_v4, @allowed_netaddrips );
		my @combined_list  = NetAddr::IP::Compact( @allowed_netaddrips, @{$coalesced_list} );
		foreach my $net (@combined_list) {
			my $range = $net->range();
			$range =~ s/\s+//g;
			$ipallow->[$i]->{src_ip} = $range;
			$ipallow->[$i]->{action} = "ip_allow";
			$ipallow->[$i]->{method} = "ALL";
			$i++;
		}

		# now add IPv6. TODO JvD: paremeterize support enabled on/ofd and /48 and number 5
		my @compacted__ipv6_list = NetAddr::IP::Compact(@allowed_ipv6_netaddrips);
		my $coalesced_ipv6_list  = NetAddr::IP::Coalesce( $coalesce_masklen_v6, $coalesce_number_v6, @allowed_ipv6_netaddrips );
		my @combined_ipv6_list   = NetAddr::IP::Compact( @allowed_ipv6_netaddrips, @{$coalesced_ipv6_list} );
		foreach my $net (@combined_ipv6_list) {
			my $range = $net->range();
			$range =~ s/\s+//g;
			$ipallow->[$i]->{src_ip} = $range;
			$ipallow->[$i]->{action} = "ip_allow";
			$ipallow->[$i]->{method} = "ALL";
			$i++;
		}

		# allow RFC 1918 server space - TODO JvD: parameterize
		$ipallow->[$i]->{src_ip} = '10.0.0.0-10.255.255.255';
		$ipallow->[$i]->{action} = 'ip_allow';
		$ipallow->[$i]->{method} = "ALL";
		$i++;

		$ipallow->[$i]->{src_ip} = '172.16.0.0-172.31.255.255';
		$ipallow->[$i]->{action} = 'ip_allow';
		$ipallow->[$i]->{method} = "ALL";
		$i++;

		$ipallow->[$i]->{src_ip} = '192.168.0.0-192.168.255.255';
		$ipallow->[$i]->{action} = 'ip_allow';
		$ipallow->[$i]->{method} = "ALL";
		$i++;

		# end with a deny
		$ipallow->[$i]->{src_ip} = '0.0.0.0-255.255.255.255';
		$ipallow->[$i]->{action} = 'ip_deny';
		$ipallow->[$i]->{method} = "ALL";
		$i++;
		$ipallow->[$i]->{src_ip} = '::-ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff';
		$ipallow->[$i]->{action} = 'ip_deny';
		$ipallow->[$i]->{method} = "ALL";
		$i++;
	}
	else {

		# for edges deny "PUSH|PURGE|DELETE", allow everything else to everyone.
		$ipallow->[$i]->{src_ip} = '0.0.0.0-255.255.255.255';
		$ipallow->[$i]->{action} = 'ip_deny';
		$ipallow->[$i]->{method} = "PUSH|PURGE|DELETE";
		$i++;
		$ipallow->[$i]->{src_ip} = '::-ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff';
		$ipallow->[$i]->{action} = 'ip_deny';
		$ipallow->[$i]->{method} = "PUSH|PURGE|DELETE";
		$i++;
	}

	return $ipallow;
}

sub facts {
	my $self     = shift;
	my $id       = shift;
	my $filename = shift;

	my $server = $self->server_data($id);
	my $text   = $self->header_comment( $server->host_name );
	$text .= "profile:" . $server->profile->name . "\n";

	return $text;
}

sub logging_dot_yaml {
	my $self     = shift;
	my $id       = shift;
	my $filename = shift;

	my $server = $self->server_data($id);
	my $data   = $self->param_data( $server, $filename );

	my $text   = "# Generated for " . $server->host_name . " by " . &name_version_string($self) . " - Do not edit!! \n";

	my $max_log_objects = 10;
	for ( my $i = 0; $i < $max_log_objects; $i = $i + 1 ) {
		my $log_format_field = "LogFormat";
		my $log_object_field = "LogObject";
		if ( $i > 0 ) {
			$log_format_field = $log_format_field . "$i";
			$log_object_field = $log_object_field . "$i";
		}

		my $log_format_name = $data->{$log_format_field . ".Name"} || "";
		if ( length($log_format_name) > 0 ) {
			my $format = $data->{$log_format_field . ".Format"};
			$text .= "formats: \n";
			$text .= " - name: " . $log_format_name ." \n";
			$text .= "   format: '" . $format . "'\n";
		}

		my $log_object_filename = $data->{$log_object_field . ".Filename"} || "";
		if ( length($log_object_filename) > 0 ) {
			my $log_object_format               = $data->{$log_object_field . ".Format"}             || "";
			my $log_object_rolling_enabled      = $data->{$log_object_field . ".RollingEnabled"}     || "";
			my $log_object_rolling_interval_sec = $data->{$log_object_field . ".RollingIntervalSec"} || "";
			my $log_object_rolling_offset_hr    = $data->{$log_object_field . ".RollingOffsetHr"}    || "";
			my $log_object_rolling_size_mb      = $data->{$log_object_field . ".RollingSizeMb"}      || "";
			my $log_object_header               = $data->{$log_object_field . ".Header"}             || "";

			$text .= "\nlogs:\n";
			$text .= "- mode: ascii\n";
			$text .= "  filename: " . $log_object_filename ."\n";
			$text .= "  format: ". $log_format_name . "\n";
			if ($log_object_rolling_enabled ne ""){
				$text .= "  rolling_enabled: ". $log_object_rolling_enabled . "\n";
			}
			if ($log_object_rolling_interval_sec ne "") {
				$text .= "  rolling_interval_sec: ". $log_object_rolling_interval_sec . "\n";
			}
			if ($log_object_rolling_offset_hr ne "") {
				$text .= "  rolling_offset_hr: ". $log_object_rolling_offset_hr . "\n";
			}
			if ($log_object_rolling_size_mb ne "") {
				$text .= "  rolling_size_mb: ". $log_object_rolling_size_mb . "\n";
			}
		}
	}

	return $text;
}

sub logging_dot_config {
	my $self     = shift;
	my $id       = shift;
	my $filename = shift;

	my $server = $self->server_data($id);
	my $data   = $self->param_data( $server, $filename );

	my $text   = "-- Generated for " . $server->host_name . " by " . &name_version_string($self) . " - Do not edit!! --\n";

	my $max_log_objects = 10;
	for ( my $i = 0; $i < $max_log_objects; $i = $i + 1 ) {
		my $log_format_field = "LogFormat";
		my $log_object_field = "LogObject";
		if ( $i > 0 ) {
			$log_format_field = $log_format_field . "$i";
			$log_object_field = $log_object_field . "$i";
		}

		my $log_format_name = $data->{$log_format_field . ".Name"} || "";
		if ( length($log_format_name) > 0 ) {
			my $format = $data->{$log_format_field . ".Format"};
			$format =~ s/"/\\\"/g;
			$text .= $log_format_name . " = format {\n";
			$text .= "	Format = '" . $format . " '\n";
			$text .= "}\n";
		}

		my $log_object_filename = $data->{$log_object_field . ".Filename"} || "";
		if ( length($log_object_filename) > 0 ) {
			my $log_object_format               = $data->{$log_object_field . ".Format"}             || "";
			my $log_object_rolling_enabled      = $data->{$log_object_field . ".RollingEnabled"}     || "";
			my $log_object_rolling_interval_sec = $data->{$log_object_field . ".RollingIntervalSec"} || "";
			my $log_object_rolling_offset_hr    = $data->{$log_object_field . ".RollingOffsetHr"}    || "";
			my $log_object_rolling_size_mb      = $data->{$log_object_field . ".RollingSizeMb"}      || "";
			my $log_object_header               = $data->{$log_object_field . ".Header"}             || "";

			$text .= "\nlog.ascii {\n";
			$text .= "  Format = " . $log_format_name . ",\n";
			$text .= "  Filename = '" . $log_object_filename . "',\n";
			$text .= "  RollingEnabled = " . $log_object_rolling_enabled . ",\n" unless defined();
			$text .= "  RollingIntervalSec = " . $log_object_rolling_interval_sec . ",\n";
			$text .= "  RollingOffsetHr = " . $log_object_rolling_offset_hr . ",\n";
			$text .= "  RollingSizeMb = " . $log_object_rolling_size_mb . "\n";
			$text .= "}\n";
		}
	}

	return $text;
}

sub logs_xml_dot_config {
	my $self     = shift;
	my $id       = shift;
	my $filename = shift;

	my $server = $self->server_data($id);
	my $data   = $self->param_data( $server, $filename );
	my $text   = "<!-- Generated for " . $server->host_name . " by " . &name_version_string($self) . " - Do not edit!! -->\n";

	my $max_log_objects = 10;
	for ( my $i = 0; $i < $max_log_objects; $i = $i + 1 ) {
		my $log_format_field = "LogFormat";
		my $log_object_field = "LogObject";
		if ( $i > 0 ) {
			$log_format_field = $log_format_field . "$i";
			$log_object_field = $log_object_field . "$i";
		}

		my $log_format_name = $data->{$log_format_field . ".Name"} || "";
		if ( length($log_format_name) > 0 ) {
			my $format = $data->{$log_format_field . ".Format"};
			$format =~ s/"/\\\"/g;
			$text .= "<LogFormat>\n";
			$text .= "  <Name = \"" . $log_format_name . "\"/>\n";
			$text .= "  <Format = \"" . $format . "\"/>\n";
			$text .= "</LogFormat>\n";
		}

		my $log_object_filename = $data->{$log_object_field . ".Filename"} || "";
		if ( length($log_object_filename) > 0 ) {
			my $log_object_format               = $data->{$log_object_field . ".Format"}             || "";
			my $log_object_rolling_enabled      = $data->{$log_object_field . ".RollingEnabled"}     || "";
			my $log_object_rolling_interval_sec = $data->{$log_object_field . ".RollingIntervalSec"} || "";
			my $log_object_rolling_offset_hr    = $data->{$log_object_field . ".RollingOffsetHr"}    || "";
			my $log_object_rolling_size_mb      = $data->{$log_object_field . ".RollingSizeMb"}      || "";
			my $log_object_header               = $data->{$log_object_field . ".Header"}             || "";

			$text .= "<LogObject>\n";
			$text .= "  <Format = \"" . $log_object_format . "\"/>\n";
			$text .= "  <Filename = \"" . $log_object_filename . "\"/>\n";
			$text .= "  <RollingEnabled = " . $log_object_rolling_enabled . "/>\n" unless defined();
			$text .= "  <RollingIntervalSec = " . $log_object_rolling_interval_sec . "/>\n";
			$text .= "  <RollingOffsetHr = " . $log_object_rolling_offset_hr . "/>\n";
			$text .= "  <RollingSizeMb = " . $log_object_rolling_size_mb . "/>\n";
			$text .= "  <Header = \"" . $log_object_header . "\"/>\n" if ( length($log_object_header) > 0 );
			$text .= "</LogObject>\n";
		}
	}

	return $text;
}

sub cacheurl_dot_config {
	my $self     = shift;
	my $id       = shift;
	my $filename = shift;
	my $data     = shift;

	my $server = $self->server_data($id);
	my $text   = $self->header_comment( $server->host_name );
	if ( !defined($data) ) {
		$data = $self->ds_data($server);
	}

	if ( $filename eq "cacheurl_qstring.config" ) {    # This is the per remap drop qstring w cacheurl use case, the file is the same for all remaps
		$text .= "http://([^?]+)(?:\\?|\$)  http://\$1\n";
		$text .= "https://([^?]+)(?:\\?|\$)  https://\$1\n";
	}
	elsif ( $filename =~ /cacheurl_(.*).config/ )
	{    # Yes, it's possibe to have the same plugin invoked multiple times on the same remap line, this is from the remap entry
		my $ds_xml_id = $1;
		my $ds = $self->db->resultset('Deliveryservice')->search( { xml_id => $ds_xml_id }, { prefetch => [ 'type', 'profile' ] } )->single();
		if ($ds) {
			$text .= $ds->cacheurl . "\n";
		}
	}
	elsif ( $filename eq "cacheurl.config" ) {    # this is the global drop qstring w cacheurl use case
		foreach my $remap ( @{ $data->{dslist} } ) {
			if ( $remap->{qstring_ignore} == 1 ) {
				my $org = $remap->{org};
				$org =~ /(https?:\/\/)(.*)/;
				$text .= "$1(" . $2 . "/[^?]+)(?:\\?|\$)  $1\$1\n";
			}
		}

	}

	$text =~ s/\s*__RETURN__\s*/\n/g;
	return $text;
}

# generic key $separator value pairs from the data hash
sub url_sig_config {
	my $self = shift;
	my $id   = shift;
	my $file = shift;

	my $sep    = defined( $separator->{$file} ) ? $separator->{$file} : " = ";
	my $server = $self->server_data($id);
	my $data   = $self->param_data( $server, $file );
	my $text   = $self->header_comment( $server->host_name );

	my $response_container = $self->riak_get( URL_SIG_KEYS_BUCKET, $file );
	my $response = $response_container->{response};
	if ( $response->is_success() ) {
		my $response_json = decode_json( $response->content );
		my $keys          = $response_json;
		foreach my $parameter ( sort keys %{$data} ) {
			if ( !defined($keys) || $parameter !~ /^key\d+/ ) {    # only use key parameters as a fallback (temp, remove me later)
				$text .= $parameter . $sep . $data->{$parameter} . "\n";
			}
		}

		# $self->app->log->debug( "keys #-> " . Dumper($keys) );
		foreach my $parameter ( sort keys %{$keys} ) {
			$text .= $parameter . $sep . $keys->{$parameter} . "\n";
		}
		return $text;
	}
	else {
		my $error = $response->content;
		return "Error: " . $error;
	}
}

# generic key $separator value pairs from the data hash
sub generic_config {
	my $self = shift;
	my $id   = shift;
	my $file = shift;

	my $sep = defined( $separator->{$file} ) ? $separator->{$file} : " = ";

	my $server = $self->server_data($id);
	my $data   = $self->param_data( $server, $file );
	my $text   = $self->header_comment( $server->host_name );
	foreach my $parameter ( sort keys %{$data} ) {
		my $p_name = $parameter;
		$p_name =~ s/__\d+$//;
		$text .= $p_name . $sep . $data->{$parameter} . "\n";
	}
	return $text;
}

sub get_num_volumes {
	my $data = shift;

	my $num            = 0;
	my @drive_prefixes = qw( Drive_Prefix SSD_Drive_Prefix RAM_Drive_Prefix);
	foreach my $pre (@drive_prefixes) {
		if ( exists $data->{$pre} ) {
			$num++;
		}
	}
	return $num;
}

sub volume_dot_config_volume_text {
	my $volume      = shift;
	my $num_volumes = shift;
	my $size        = int( 100 / $num_volumes );
	return "volume=$volume scheme=http size=$size%\n";
}

sub volume_dot_config {
	my $self = shift;
	my $id   = shift;
	my $file = shift;

	my $server = $self->server_data($id);
	my $data   = $self->param_data( $server, "storage.config" );
	my $text   = $self->header_comment( $server->host_name );

	my $num_volumes = get_num_volumes($data);

	my $next_volume = 1;
	$text .= "# 12M NOTE: This is running with forced volumes - the size is irrelevant\n";
	if ( defined( $data->{Drive_Prefix} ) ) {
		$text .= volume_dot_config_volume_text( $next_volume, $num_volumes );
		$next_volume++;
	}
	if ( defined( $data->{RAM_Drive_Prefix} ) ) {
		$text .= volume_dot_config_volume_text( $next_volume, $num_volumes );
		$next_volume++;
	}
	if ( defined( $data->{SSD_Drive_Prefix} ) ) {
		$text .= volume_dot_config_volume_text( $next_volume, $num_volumes );
		$next_volume++;
	}
	return $text;
}

sub hosting_dot_config {
	my $self = shift;
	my $id   = shift;
	my $file = shift;
	my $data = shift;

	my $server       = $self->server_data($id);
	my $storage_data = $self->param_data( $server, "storage.config" );
	my $text         = $self->header_comment( $server->host_name );
	if ( !defined($data) ) {
		$data = $self->ds_data($server);
	}

	if ( defined( $storage_data->{RAM_Drive_Prefix} ) ) {
		my $next_volume = 1;
		if ( defined( $storage_data->{Drive_Prefix} ) ) {
			my $disk_volume = $next_volume;
			$text .= "# 12M NOTE: volume " . $disk_volume . " is the Disk volume\n";
			$next_volume++;
		}
		my $ram_volume = $next_volume;
		$text .= "# 12M NOTE: volume " . $ram_volume . " is the RAM volume\n";

		my %listed = ();
		foreach my $remap ( @{ $data->{dslist} } ) {
			if (   ( ( $remap->{type} =~ /_LIVE$/ || $remap->{type} =~ /_LIVE_NATNL$/ ) && $server->type->name =~ m/^EDGE/ )
				|| ( $remap->{type} =~ /_LIVE_NATNL$/ && $server->type->name =~ m/^MID/ ) )
			{
				if ( defined( $listed{ $remap->{org} } ) ) { next; }
				my $org_fqdn = $remap->{org};
				$org_fqdn =~ s/https?:\/\///;
				$text .= "hostname=" . $org_fqdn . " volume=" . $ram_volume . "\n";
				$listed{ $remap->{org} } = 1;
			}
		}
	}
	my $disk_volume = 1;    # note this will actually be the RAM (RAM_Drive_Prefix) volume if there is no Drive_Prefix parameter.
	$text .= "hostname=*   volume=" . $disk_volume . "\n";

	return $text;
}

sub storage_dot_config_volume_text {
	my $prefix  = shift;
	my $letters = shift;
	my $volume  = shift;

	my $text = "";
	my @postfix = split( /,/, $letters );
	foreach my $l ( sort @postfix ) {
		$text .= $prefix . $l;
		$text .= " volume=" . $volume;
		$text .= "\n";
	}
	return $text;
}

sub storage_dot_config {
	my $self = shift;
	my $id   = shift;
	my $file = shift;

	my $server = $self->server_data($id);
	my $text   = $self->header_comment( $server->host_name );
	my $data   = $self->param_data( $server, $file );

	my $next_volume = 1;
	if ( defined( $data->{Drive_Prefix} ) ) {
		$text .= storage_dot_config_volume_text( $data->{Drive_Prefix}, $data->{Drive_Letters}, $next_volume );
		$next_volume++;
	}

	if ( defined( $data->{RAM_Drive_Prefix} ) ) {
		$text .= storage_dot_config_volume_text( $data->{RAM_Drive_Prefix}, $data->{RAM_Drive_Letters}, $next_volume );
		$next_volume++;
	}

	if ( defined( $data->{SSD_Drive_Prefix} ) ) {
		$text .= storage_dot_config_volume_text( $data->{SSD_Drive_Prefix}, $data->{SSD_Drive_Letters}, $next_volume );
		$next_volume++;
	}
	return $text;
}

sub ats_dot_rules {
	my $self = shift;
	my $id   = shift;
	my $file = shift;

	my $server = $self->server_data($id);
	my $text   = $self->header_comment( $server->host_name );
	my $data   = $self->param_data( $server, "storage.config" );    # ats.rules is based on the storage.config params

	my $drive_prefix = $data->{Drive_Prefix};
	my @drive_postfix = split( /,/, $data->{Drive_Letters} );
	foreach my $l ( sort @drive_postfix ) {
		$drive_prefix =~ s/\/dev\///;
		$text .= "KERNEL==\"" . $drive_prefix . $l . "\", OWNER=\"ats\"\n";
	}
	if ( defined( $data->{RAM_Drive_Prefix} ) ) {
		$drive_prefix = $data->{RAM_Drive_Prefix};
		@drive_postfix = split( /,/, $data->{RAM_Drive_Letters} );
		foreach my $l ( sort @drive_postfix ) {
			$drive_prefix =~ s/\/dev\///;
			$text .= "KERNEL==\"" . $drive_prefix . $l . "\", OWNER=\"ats\"\n";
		}
	}

	return $text;
}

sub cache_dot_config {
	my $self = shift;
	my $id   = shift;
	my $file = shift;
	my $data = shift;

	my $server = $self->server_data($id);
	my $text   = $self->header_comment( $server->host_name );
	if ( !defined($data) ) {
		$data = $self->ds_data($server);
	}

	foreach my $remap ( @{ $data->{dslist} } ) {
		if ( $remap->{type} eq "HTTP_NO_CACHE" ) {
			my $org_fqdn = $remap->{org};
			$org_fqdn =~ s/https?:\/\///;
			$org_fqdn =~ m/(.*?):(\d+).*/;
			my $org_port = $2;

			if (defined($org_port)) {
				$org_fqdn = $1;
				$text .= "dest_domain=" . $org_fqdn . " port=" . $org_port . " scheme=http action=never-cache\n";
			} else { 
				$text .= "dest_domain=" . $org_fqdn . " scheme=http action=never-cache\n";
			}
		}
	}
	return $text;
}

sub remap_dot_config {
	my $self = shift;
	my $id   = shift;
	my $file = shift;
	my $data = shift;

	my $server = $self->server_data($id);
	my $pdata  = $self->param_data( $server, 'package' );
	my $text   = $self->header_comment( $server->host_name );
	if ( !defined($data) ) {
		$data = $self->ds_data($server);
	}

	if ( $server->type->name =~ m/^MID/ ) {
		my %mid_remap;
		foreach my $remap ( @{ $data->{dslist} } ) {
			if ( $remap->{type} =~ /LIVE/ && $remap->{type} !~ /NATNL/ ) {
				next;    # Live local delivery services skip mids
			}
			if ( defined( $remap->{org} ) && defined( $mid_remap{ $remap->{org} } ) ) {
				next;    # skip remap rules from extra HOST_REGEXP entries
			}

			if ( defined( $remap->{mid_header_rewrite} ) && $remap->{mid_header_rewrite} ne "" ) {
				$mid_remap{ $remap->{org} } .= " \@plugin=header_rewrite.so \@pparam=" . $remap->{mid_hdr_rw_file};
			}
			if ( $remap->{qstring_ignore} == 1 ) {
				$mid_remap{ $remap->{org} } .= UI::DeliveryService::get_qstring_ignore_remap(UI::DeliveryService::get_ats_major_version($self, $server), $remap->{range_request_handling});
			}
			if ( defined( $remap->{cacheurl} ) && $remap->{cacheurl} ne "" ) {
				$mid_remap{ $remap->{org} } .= " \@plugin=cacheurl.so \@pparam=" . $remap->{cacheurl_file};
			}
			if ( $remap->{range_request_handling} == 2 ) {
				$mid_remap{ $remap->{org} } .= " \@plugin=cache_range_requests.so";
			}
		}
		foreach my $key ( keys %mid_remap ) {
			$text .= "map " . $key . " " . $key . $mid_remap{$key} . "\n";
		}
		return $text;
	}

	# mids don't get here.
	foreach my $remap ( @{ $data->{dslist} } ) {
		foreach my $map_from ( keys %{ $remap->{remap_line} } ) {
			my $map_to = $remap->{remap_line}->{$map_from};
			$text = $self->build_remap_line( $server, $pdata, $text, $data, $remap, $map_from, $map_to );
		}
		foreach my $map_from ( keys %{ $remap->{remap_line2} } ) {
			my $map_to = $remap->{remap_line2}->{$map_from};
			$text = $self->build_remap_line( $server, $pdata, $text, $data, $remap, $map_from, $map_to );
		}
	}
	return $text;
}

sub build_remap_line {
	my $self     = shift;
	my $server   = shift;
	my $pdata    = shift;
	my $text     = shift;
	my $data     = shift;
	my $remap    = shift;
	my $map_from = shift;
	my $map_to   = shift;

	if ( $remap->{type} eq 'ANY_MAP' ) {
		$text .= $remap->{remap_text} . "\n";
		return $text;
	}

	my $host_name = $data->{host_name};
	my $dscp      = $remap->{dscp};

	$map_from =~ s/__http__/$host_name/;

	if ( defined( $pdata->{'dscp_remap'} ) ) {
		$text .= "map	" . $map_from . "     " . $map_to . " \@plugin=dscp_remap.so \@pparam=" . $dscp;
	}
	else {
		$text .= "map	" . $map_from . "     " . $map_to . " \@plugin=header_rewrite.so \@pparam=dscp/set_dscp_" . $dscp . ".config";
	}
	if ( defined( $remap->{edge_header_rewrite} ) ) {
		$text .= " \@plugin=header_rewrite.so \@pparam=" . $remap->{hdr_rw_file};
	}
	if ( defined($remap->{signing_algorithm})) {
		if ( $remap->{signing_algorithm} eq "url_sig" ) {
			$text .= " \@plugin=url_sig.so \@pparam=url_sig_" . $remap->{ds_xml_id} . ".config";
		}
		elsif ( $remap->{signing_algorithm} eq "uri_signing" ) {
			$text .= " \@plugin=uri_signing.so \@pparam=uri_signing_" . $remap->{ds_xml_id} . ".config";
		}
	}
	if ( $remap->{qstring_ignore} == 2 ) {
		my $dqs_file = "drop_qstring.config";
		$text .= " \@plugin=regex_remap.so \@pparam=" . $dqs_file;
	}
	elsif ( $remap->{qstring_ignore} == 1 ) {
		my $global_exists = $self->profile_param_value( $server->profile->id, 'cacheurl.config', 'location', undef );
		if ($global_exists) {
			$self->app->log->debug(
				"qstring_ignore == 1, but global cacheurl.config param exists, so skipping remap rename config_file=cacheurl.config parameter if you want to change"
			);
		}
		else {
			#If we are on ats 6 and later we want to use the cachekey plugin, otherwise we have to use cacheurl
			$text .= UI::DeliveryService::get_qstring_ignore_remap(UI::DeliveryService::get_ats_major_version($self, $server), $remap->{range_request_handling});
		}
	}
	if ( defined( $remap->{cacheurl} ) && $remap->{cacheurl} ne "" ) {
		$text .= " \@plugin=cacheurl.so \@pparam=" . $remap->{cacheurl_file};
	}

	# Note: should use full path here?
	if ( defined( $remap->{regex_remap} ) && $remap->{regex_remap} ne "" ) {
		$text .= " \@plugin=regex_remap.so \@pparam=regex_remap_" . $remap->{ds_xml_id} . ".config";
	}
	if ( $remap->{range_request_handling} == 1 ) {
		$text .= " \@plugin=background_fetch.so \@pparam=bg_fetch.config";
	}
	elsif ( $remap->{range_request_handling} == 2 ) {
		$text .= " \@plugin=cache_range_requests.so ";
	}
	if ( defined( $remap->{remap_text} ) ) {
		$text .= " " . $remap->{remap_text};
	}
	$text .= "\n";
	return $text;
}

sub format_parent_info {
	my $parent = shift;
	if ( !defined $parent ) {
		return "";    # should never happen..
	}
	my $host =
		( $parent->{use_ip_address} == 1 )
		? $parent->{ip_address}
		: $parent->{host_name} . '.' . $parent->{domain_name};

	my $port   = $parent->{port};
	my $weight = $parent->{weight};
	my $text   = "$host:$port|$weight;";
	return $text;
}

sub parent_dot_config {
	my $self = shift;
	my $id   = shift;
	my $file = shift;
	my $data = shift;

	my $server      = $self->server_data($id);
	my $server_type = $server->type->name;

	my $ats_ver =
		$self->db->resultset('ProfileParameter')
		->search( { 'parameter.name' => 'trafficserver', 'parameter.config_file' => 'package', 'profile.id' => $server->profile->id },
		{ prefetch => [ 'profile', 'parameter' ] } )->get_column('parameter.value')->single();
	if (!defined $ats_ver) {
	        $ats_ver = "5";
                $self->app->log->error("Parameter package.trafficserver missing for profile ".$server->profile->name . ". Assuming version $ats_ver");
        }
	my $ats_major_version = substr( $ats_ver, 0, 1 );

	my $pinfo;
	my $text = $self->header_comment( $server->host_name );
	if ( !defined($data) ) {
		$data = $self->ds_data($server);
	}

	# Origin Shield or Multi Site Origin
	#$self->app->log->debug( "id = $id and server_type = $server_type,  hostname = " . $server->{host_name} );
	if ( $server_type =~ m/^MID/ ) {
		my @unique_origin;
		foreach my $ds ( @{ $data->{dslist} } ) {
			my $xml_id                             = $ds->{ds_xml_id};
			my $os                                 = $ds->{origin_shield};
			my $multi_site_origin                  = $ds->{multi_site_origin} || 0;
			my $mso_algorithm                      = $ds->{'param'}->{'parent.config'}->{'mso.algorithm'} || 0;
			my $parent_retry                       = $ds->{'param'}->{'parent.config'}->{'mso.parent_retry'};
			my $unavailable_server_retry_responses = $ds->{'param'}->{'parent.config'}->{'mso.unavailable_server_retry_responses'};
			my $max_simple_retries                 = $ds->{'param'}->{'parent.config'}->{'mso.max_simple_retries'} || 1;
			my $max_unavailable_server_retries     = $ds->{'param'}->{'parent.config'}->{'mso.max_unavailable_server_retries'} || 1;

			my $qsh            = $ds->{'param'}->{'parent.config'}->{'psel.qstring_handling'};
			my $parent_qstring = "ignore";                                                      # default is ignore, unless for alg consistent_hash
			if ( !defined($qsh) && $mso_algorithm eq 'consistent_hash' && $ds->{qstring_ignore} == 0 ) {
				$parent_qstring = 'consider';
			}
			my $org_uri = URI->new( $ds->{org} );

			# Don't duplicate origin line if multiple seen
			next if ( grep( /^$org_uri$/, @unique_origin ) );
			push @unique_origin, $org_uri;

			if ( defined($os) ) {
				my $pselect_alg = $self->profile_param_value( $server->profile->id, 'parent.config', 'algorithm', undef );
				my $algorithm = "";
				if ( defined($pselect_alg) ) {
					$algorithm = "round_robin=$pselect_alg";
				}
				$text .= "dest_domain=" . $org_uri->host . " port=" . $org_uri->port . " parent=$os $algorithm go_direct=true\n";
			}
			elsif ($multi_site_origin) {
				$text .= "dest_domain=" . $org_uri->host . " port=" . $org_uri->port . " ";

				# If we have multi-site origin, get parent_data once
				if ( not defined($pinfo) ) {
					$pinfo = $self->parent_data($server);
				}

				my @ranked_parents = ();
				if ( exists( $pinfo->{ $org_uri->host } ) ) {
					@ranked_parents = sort by_parent_rank @{ $pinfo->{ $org_uri->host } };
				}
				else {
					$self->app->log->debug( "BUG: Did not match an origin: " . $org_uri );
				}

				my @parent_info;
				my @secondary_parent_info;
				my @null_parent_info;
				foreach my $parent (@ranked_parents) {
					if ( $parent->{primary_parent} ) {
						push @parent_info, format_parent_info($parent);
					}
					elsif ( $parent->{secondary_parent} ) {
						push @secondary_parent_info, format_parent_info($parent);
					}
					else {
						push @null_parent_info, format_parent_info($parent);
					}
				}
				my %seen;
				@parent_info = grep { !$seen{$_}++ } @parent_info;

				if ( scalar @secondary_parent_info > 0 ) {
					my %seen;
					@secondary_parent_info = grep { !$seen{$_}++ } @secondary_parent_info;
				}
				if ( scalar @null_parent_info > 0 ) {
					my %seen;
					@null_parent_info = grep { !$seen{$_}++ } @null_parent_info;
				}
				my $parents = 'parent="' . join( '', @parent_info ) . '' . join( '', @secondary_parent_info ) . '' . join( '', @null_parent_info ) . '"';

				$text .= "$parents round_robin=$mso_algorithm qstring=$parent_qstring go_direct=false parent_is_proxy=false";

				if ( $ats_major_version >= 6 && $parent_retry ne "" ) {
					if ( $unavailable_server_retry_responses ne "") {
						$text .= " parent_retry=$parent_retry unavailable_server_retry_responses=$unavailable_server_retry_responses";
					} else {
						$text .= " parent_retry=$parent_retry";
					}
					$text .= " max_simple_retries=$max_simple_retries max_unavailable_server_retries=$max_unavailable_server_retries";
				}
				$text .= "\n";
			}
		}

		#$text .= "dest_domain=. go_direct=true\n"; # this is implicit.
		#$self->app->log->debug( "MID PARENT.CONFIG:\n" . $text . "\n" );
		return $text;
	}
	else {    #"True" Parent - we are genning a EDGE config that points to a parent proxy.
		$pinfo = $self->parent_data($server);
		my %done = ();

		foreach my $remap ( @{ $data->{dslist} } ) {
			my $org = $remap->{org};
			next if !defined $org || $org eq "";
			next if $done{$org};
			my $org_uri = URI->new($org);
			if ( $remap->{type} eq "HTTP_NO_CACHE" || $remap->{type} eq "HTTP_LIVE" || $remap->{type} eq "DNS_LIVE" ) {
				$text .= "dest_domain=" . $org_uri->host . " port=" . $org_uri->port . " go_direct=true\n";
			}
			else {
				# check for profile psel.qstring_handling.  If this parameter is assigned to the server profile,
				# then edges will use the qstring handling value specified in the parameter for all profiles.
				my $qsh = $self->profile_param_value( $server->profile->id, 'parent.config', 'psel.qstring_handling');
				# If there is no defined parameter in the profile, then check the delivery service profile.
				# If psel.qstring_handling exists in the DS profile, then we use that value for the specified DS only.
				# This is used only if not overridden by a server profile qstring handling parameter.
				if (!defined($qsh)) {
					$qsh = $remap->{'param'}->{'parent.config'}->{'psel.qstring_handling'};
				}
				my $parent_qstring = defined($qsh) ? $qsh : "ignore";
				if ( $remap->{qstring_ignore} == 0 && !defined($qsh) ) {
					$parent_qstring = "consider";
				}

				my @parent_info;
				my @secondary_parent_info;
				foreach my $parent ( @{ $pinfo->{all_parents} } ) {
					my $ptxt = format_parent_info($parent);
					if ( $parent->{primary_parent} ) {
						push @parent_info, $ptxt;
					}
					elsif ( $parent->{secondary_parent} ) {
						push @secondary_parent_info, $ptxt;
					}
				}
				if ( scalar @parent_info == 0  ) {
					@parent_info = @secondary_parent_info;
					@secondary_parent_info = ();
				}
				my %seen;
				@parent_info = grep { !$seen{$_}++ } @parent_info;
				my $parents = 'parent="' . join( '', @parent_info ) . '"';
				my $secparents = '';
				if ( scalar @secondary_parent_info > 0 ) {
					my %seen;
					@secondary_parent_info = grep { !$seen{$_}++ } @secondary_parent_info;
					$secparents = 'secondary_parent="' . join( '', @secondary_parent_info ) . '"';
				}
				my $round_robin = 'round_robin=consistent_hash';
				my $go_direct   = 'go_direct=false';
				$text
					.= "dest_domain="
					. $org_uri->host
					. " port="
					. $org_uri->port
					. " $parents $secparents $round_robin $go_direct qstring=$parent_qstring\n";
			}
			$done{$org} = 1;
		}

		my $pselect_alg = $self->profile_param_value( $server->profile->id, 'parent.config', 'algorithm', undef );
		if ( defined($pselect_alg) && $pselect_alg eq 'consistent_hash' ) {
			my @parent_info;
			foreach my $parent ( @{ $pinfo->{"all_parents"} } ) {
				push @parent_info, $parent->{"host_name"} . "." . $parent->{"domain_name"} . ":" . $parent->{"port"} . "|" . $parent->{"weight"} . ";";
			}
			my %seen;
			@parent_info = grep { !$seen{$_}++ } @parent_info;
			$text .= "dest_domain=.";
			$text .= " parent=\"" . join( '', @parent_info ) . "\"";
			$text .= " round_robin=consistent_hash go_direct=false";
		}
		else {    # default to old situation.
			$text .= "dest_domain=.";
			my @parent_info;
			foreach my $parent ( @{ $pinfo->{"all_parents"} } ) {
				push @parent_info, $parent->{"host_name"} . "." . $parent->{"domain_name"} . ":" . $parent->{"port"} . ";";
			}
			my %seen;
			@parent_info = grep { !$seen{$_}++ } @parent_info;
			$text .= " parent=\"" . join( '', @parent_info ) . "\"";
			$text .= " round_robin=urlhash go_direct=false";
		}

		my $qstring = $self->profile_param_value( $server->profile->id, 'parent.config', 'qstring', undef );
		if ( defined($qstring) ) {
			$text .= " qstring=" . $qstring;
		}

		$text .= "\n";

		# $self->app->log->debug($text);
		return $text;
	}
}

sub ip_allow_dot_config {
	my $self = shift;
	my $id   = shift;
	my $file = shift;

	my $server = $self->server_data($id);
	my $text   = $self->header_comment( $server->host_name );
	my $data   = $self->ip_allow_data( $server, $file );

	foreach my $access ( @{$data} ) {
		$text .= sprintf( "src_ip=%-70s action=%-10s method=%-20s\n", $access->{src_ip}, $access->{action}, $access->{method} );
	}

	return $text;
}

sub regex_revalidate_dot_config {
	my $self = shift;
	my $id   = shift;
	my $file = shift;

	# note: Calling this from outside Configfiles, so $self->method doesn't work. TODO: Be smarter
	# my $server = $self->server_data($id);
	# my $text   = $self->header_comment( $server->host_name );
	my $server = &server_data( $self, $id );

	my $text = "# DO NOT EDIT - Generated for CDN " . $server->cdn->name . " by " . &name_version_string($self) . " on " . `date`;

	my $max_days =
		$self->db->resultset('Parameter')->search( { name => "maxRevalDurationDays" }, { config_file => "regex_revalidate.config" } )->get_column('value')
		->first;
	my $interval = "> now() - interval '$max_days day'";    # postgres

	my %regex_time;
	$max_days =
		$self->db->resultset('Parameter')->search( { name => "maxRevalDurationDays" }, { config_file => "regex_revalidate.config" } )->get_column('value')
		->first;
	my $max_hours = $max_days * 24;
	my $min_hours = 1;

	my $rs = $self->db->resultset('Job')->search( { start_time => \$interval }, { prefetch => 'job_deliveryservice' } );
	while ( my $row = $rs->next ) {
		next unless defined( $row->job_deliveryservice );

		# Purges are CDN - wide, and the job entry has the ds id in it.
		my $parameters = $row->parameters;
		my $ttl;
		if ( $row->keyword eq "PURGE" && ( defined($parameters) && $parameters =~ /TTL:(\d+)h/ ) ) {
			$ttl = $1;
			if ( $ttl < $min_hours ) {
				$ttl = $min_hours;
			}
			elsif ( $ttl > $max_hours ) {
				$ttl = $max_hours;
			}
		}
		else {
			next;
		}

		my $date       = new Date::Manip::Date();
		my $start_time = $row->start_time;
		my $start_date = ParseDate($start_time);
		my $end_date   = DateCalc( $start_date, ParseDateDelta( $ttl . ':00:00' ) );
		my $err        = $date->parse($end_date);
		if ($err) {
			print "ERROR ON DATE CONVERSION:" . $err;
			next;
		}
		my $purge_end = $date->printf("%s");    # this is in secs since the unix epoch

		if ( $purge_end < time() ) {            # skip purges that have an end_time in the past
			next;
		}
		my $asset_url = $row->asset_url;

		my $job_cdn_id = $row->job_deliveryservice->cdn_id;
		if ( $server->cdn_id == $job_cdn_id ) {

			# if there are multipe with same re, pick the longes lasting.
			if ( !defined( $regex_time{ $row->asset_url } )
				|| ( defined( $regex_time{ $row->asset_url } ) && $purge_end > $regex_time{ $row->asset_url } ) )
			{
				$regex_time{ $row->asset_url } = $purge_end;
			}
		}
	}

	foreach my $re ( sort keys %regex_time ) {
		$text .= $re . " " . $regex_time{$re} . "\n";
	}

	return $text;
}

sub take_and_bake {
	my $self = shift;
	my $id   = shift;
	my $file = shift;

	my $server = $self->server_data($id);
	my $data   = $self->param_data( $server, $file );
	my $text   = $self->header_comment( $server->host_name );
	foreach my $parameter ( sort keys %{$data} ) {
		$text .= $data->{$parameter} . "\n";
	}
	return $text;
}

sub drop_qstring_dot_config {
	my $self = shift;
	my $id   = shift;
	my $file = shift;

	my $server = $self->server_data($id);
	my $text   = $self->header_comment( $server->host_name );

	$server = &server_data( $self, $id );
	my $drop_qstring = $self->profile_param_value( $server->profile->id, 'drop_qstring.config', 'content', undef );
	if ($drop_qstring) {
		$text .= $drop_qstring . "\n";
	}
	else {
		$text .= "/([^?]+) \$s://\$t/\$1\n";
	}
	return $text;
}

sub header_rewrite_dot_config {
	my $self = shift;
	my $id   = shift;
	my $file = shift;

	my $server    = $self->server_data($id);
	my $text      = $self->header_comment( $server->host_name );
	my $ds_xml_id = undef;
	if ( $file =~ /^hdr_rw_mid_(.*)\.config$/ ) {
		$ds_xml_id = $1;
		my $ds = $self->db->resultset('Deliveryservice')->search( { xml_id => $ds_xml_id } )->single();
		my $actions = $ds->mid_header_rewrite;
		$text .= $actions . "\n";
	}
	elsif ( $file =~ /^hdr_rw_(.*)\.config$/ ) {
		$ds_xml_id = $1;
		my $ds = $self->db->resultset('Deliveryservice')->search( { xml_id => $ds_xml_id } )->single();
		my $actions = $ds->edge_header_rewrite;
		$text .= $actions . "\n";
	}

	$text =~ s/\s*__RETURN__\s*/\n/g;
	my $ipv4 = $server->ip_address;
	$text =~ s/__CACHE_IPV4__/$ipv4/g;
	return $text;
}

sub regex_remap_dot_config {
	my $self = shift;
	my $id   = shift;
	my $file = shift;

	my $server = $self->server_data($id);
	my $text   = $self->header_comment( $server->host_name );

	if ( $file =~ /^regex_remap_(.*)\.config$/ ) {
		my $ds_xml_id = $1;
		my $ds = $self->db->resultset('Deliveryservice')->search( { xml_id => $ds_xml_id } )->single();
		$text .= $ds->regex_remap . "\n";
	}

	$text =~ s/\s*__RETURN__\s*/\n/g;
	return $text;
}

sub header_rewrite_dscp_dot_config {
	my $self = shift;
	my $id   = shift;
	my $file = shift;

	my $server = $self->server_data($id);
	my $text   = $self->header_comment( $server->host_name );
	my $dscp_decimal;
	if ( $file =~ /^set_dscp_(\d+)\.config$/ ) {
		$dscp_decimal = $1;
	}
	else {
		$text = "An error occured generating the DSCP header rewrite file.";
	}
	$text .= "cond %{REMAP_PSEUDO_HOOK}\n" . "set-conn-dscp " . $dscp_decimal . " [L]\n";

	return $text;
}

sub to_ext_dot_config {
	my $self = shift;
	my $id   = shift;
	my $file = shift;

	my $server = $self->server_data($id);
	my $text   = $self->header_comment( $server->host_name );

	# get the subroutine name for this file from the parameter
	my $subroutine = $self->profile_param_value( $server->profile->id, $file, 'SubRoutine', undef );
	$self->app->log->error( "ToExtDotConfigFile == " . $subroutine );

	# TODO: previous code didn't check for undef -- what to do here?
	if ( defined $subroutine ) {
		my $package;
		( $package = $subroutine ) =~ s/(.*)(::)(.*)/$1/;
		eval "use $package;";

		# And call it - the below calls the subroutine in the var $subroutine.
		no strict 'refs';
		$text .= $subroutine->( $self, $id, $file );

		# $text .= &{ \&{$subroutine} }( $self, $id, $file );
	}

	return $text;
}

sub ssl_multicert_dot_config {
	my $self = shift;
	my $id   = shift;
	my $file = shift;

	#id == hostname
	my $server = $self->server_data($id);
	my $text   = $self->header_comment( $server->host_name );

	# get a list of delivery services for the server
	my $protocol_search = '> 0';
	my @ds_list         = $self->db->resultset('Deliveryservice')->search(
		{ -and => [ 'server.id' => $server->id, 'me.protocol' => \$protocol_search ] },
		{ prefetch => ['cdn'], join => { deliveryservice_servers => { server => undef } }, }
	);
	foreach my $ds (@ds_list) {
		my $ds_id        = $ds->id;
		my $xml_id       = $ds->xml_id;
		my $domain_name  = $ds->cdn->domain_name;
		my $ds_regexes   = UI::DeliveryService::get_regexp_set( $self, $ds_id );
		my @example_urls = UI::DeliveryService::get_example_urls( $self, $ds_id, $ds_regexes, $ds, $domain_name, $ds->protocol );

		#first one is the one we want
		my $hostname = $example_urls[0];
		$hostname =~ /(https?:\/\/)(.*)/;
		my $new_host = $2;
		my $key_name = "$new_host.key";
		$new_host =~ tr/./_/;
		my $cer_name = $new_host . "_cert.cer";

		$text .= "ssl_cert_name=$cer_name\t ssl_key_name=$key_name\n";
	}
	return $text;
}

# This is a temporary workaround until we have real partial object caching support in ATS, so hardcoding for now
sub bg_fetch_dot_config {
	my $self = shift;
	my $id   = shift;

	my $server = $self->server_data($id);
	my $text   = $self->header_comment( $server->host_name );
	$text .= "include User-Agent *\n";

	return $text;
}

1;
